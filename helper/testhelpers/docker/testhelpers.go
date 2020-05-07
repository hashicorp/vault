package docker

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v3"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/go-connections/nat"
	"github.com/hashicorp/go-uuid"
)

type Runner struct {
	DockerAPI  *client.Client
	RunOptions RunOptions
}

type RunOptions struct {
	ImageRepo       string
	ImageTag        string
	ContainerName   string
	Cmd             []string
	Env             []string
	NetworkID       string
	CopyFromTo      map[string]string
	Ports           []string
	DoNotAutoRemove bool
}

func NewServiceRunner(opts RunOptions) (*Runner, error) {
	dapi, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.39"))
	if err != nil {
		return nil, err
	}

	if opts.NetworkID == "" {
		opts.NetworkID = os.Getenv("TEST_DOCKER_NETWORK_ID")
	}
	if opts.ContainerName == "" {
		if strings.Contains(opts.ImageRepo, "/") {
			return nil, fmt.Errorf("ContainerName is required for non-library images")
		}
		// If there's no slash in the repo it's almost certainly going to be
		// a good container name.
		opts.ContainerName = opts.ImageRepo
	}
	return &Runner{
		DockerAPI:  dapi,
		RunOptions: opts,
	}, nil
}

type ServiceConfig interface {
	Address() string
	URL() *url.URL
}

func NewServiceHostPort(host string, port int) *ServiceHostPort {
	return &ServiceHostPort{address: fmt.Sprintf("%s:%d", host, port)}
}

func NewServiceHostPortParse(s string) (*ServiceHostPort, error) {
	pieces := strings.Split(s, ":")
	if len(pieces) != 2 {
		return nil, fmt.Errorf("address must be of the form host:port, got: %v", s)
	}

	port, err := strconv.Atoi(pieces[1])
	if err != nil || port < 1 {
		return nil, fmt.Errorf("address must be of the form host:port, got: %v", s)
	}

	return &ServiceHostPort{s}, nil
}

type ServiceHostPort struct {
	address string
}

func (s ServiceHostPort) Address() string {
	return s.address
}

func (s ServiceHostPort) URL() *url.URL {
	return &url.URL{Host: s.address}
}

func NewServiceURLParse(s string) (*ServiceURL, error) {
	u, err := url.Parse(s)
	if err != nil {
		return nil, err
	}
	return &ServiceURL{u: *u}, nil
}

func NewServiceURL(u url.URL) *ServiceURL {
	return &ServiceURL{u: u}
}

type ServiceURL struct {
	u url.URL
}

func (s ServiceURL) Address() string {
	return s.u.Host
}

func (s ServiceURL) URL() *url.URL {
	return &s.u
}

// ServiceAdapter verifies connectivity to the service, then returns either the
// connection string (typically a URL) and nil, or empty string and an error.
type ServiceAdapter func(ctx context.Context, host string, port int) (ServiceConfig, error)

func (d *Runner) StartService(ctx context.Context, connect ServiceAdapter) (*Service, error) {
	container, hostIPs, err := d.Start(context.Background())
	if err != nil {
		return nil, err
	}

	cleanup := func() {
		for i := 0; i < 10; i++ {
			err := d.DockerAPI.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{Force: true})
			if err == nil {
				return
			}
			time.Sleep(1 * time.Second)
		}
	}

	bo := backoff.NewExponentialBackOff()
	bo.MaxInterval = time.Second * 5
	bo.MaxElapsedTime = time.Minute

	pieces := strings.Split(hostIPs[0], ":")
	portInt, err := strconv.Atoi(pieces[1])
	if err != nil {
		return nil, err
	}

	var config ServiceConfig
	err = backoff.Retry(func() error {
		c, err := connect(ctx, pieces[0], portInt)
		if err != nil {
			return err
		}
		if c == nil {
			return fmt.Errorf("service adapter returned nil error and config")
		}
		config = c
		return nil
	}, bo)

	if err != nil {
		if !d.RunOptions.DoNotAutoRemove {
			cleanup()
		}
		return nil, err
	}

	return &Service{
		Config:  config,
		Cleanup: cleanup,
	}, nil
}

type Service struct {
	Config  ServiceConfig
	Cleanup func()
}

func (d *Runner) Start(ctx context.Context) (*types.ContainerJSON, []string, error) {
	suffix, err := uuid.GenerateUUID()
	if err != nil {
		return nil, nil, err
	}
	name := d.RunOptions.ContainerName + "-" + suffix

	cfg := &container.Config{
		Hostname: name,
		Image:    fmt.Sprintf("%s:%s", d.RunOptions.ImageRepo, d.RunOptions.ImageTag),
		Env:      d.RunOptions.Env,
		Cmd:      d.RunOptions.Cmd,
	}
	if len(d.RunOptions.Ports) > 0 {
		cfg.ExposedPorts = make(map[nat.Port]struct{})
		for _, p := range d.RunOptions.Ports {
			cfg.ExposedPorts[nat.Port(p)] = struct{}{}
		}
	}

	hostConfig := &container.HostConfig{
		AutoRemove:      !d.RunOptions.DoNotAutoRemove,
		PublishAllPorts: true,
	}

	netConfig := &network.NetworkingConfig{}
	if d.RunOptions.NetworkID != "" {
		netConfig.EndpointsConfig = map[string]*network.EndpointSettings{
			d.RunOptions.NetworkID: &network.EndpointSettings{},
		}
	}

	// best-effort pull
	resp, _ := d.DockerAPI.ImageCreate(ctx, cfg.Image, types.ImageCreateOptions{})
	if resp != nil {
		_, _ = ioutil.ReadAll(resp)
	}

	container, err := d.DockerAPI.ContainerCreate(ctx, cfg, hostConfig, netConfig, cfg.Hostname)
	if err != nil {
		return nil, nil, fmt.Errorf("container create failed: %v", err)
	}

	for from, to := range d.RunOptions.CopyFromTo {
		if err := copyToContainer(ctx, d.DockerAPI, container.ID, from, to); err != nil {
			_ = d.DockerAPI.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{})
			return nil, nil, err
		}
	}

	err = d.DockerAPI.ContainerStart(ctx, container.ID, types.ContainerStartOptions{})
	if err != nil {
		_ = d.DockerAPI.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{})
		return nil, nil, fmt.Errorf("container start failed: %v", err)
	}

	inspect, err := d.DockerAPI.ContainerInspect(ctx, container.ID)
	if err != nil {
		_ = d.DockerAPI.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{})
		return nil, nil, err
	}

	var addrs []string
	for _, port := range d.RunOptions.Ports {
		pieces := strings.Split(port, "/")
		if len(pieces) < 2 {
			return nil, nil, fmt.Errorf("expected port of the form 1234/tcp, got: %s", port)
		}
		if d.RunOptions.NetworkID != "" {
			addrs = append(addrs, fmt.Sprintf("%s:%s", cfg.Hostname, pieces[0]))
		} else {
			mapped, ok := inspect.NetworkSettings.Ports[nat.Port(port)]
			if !ok || len(mapped) == 0 {
				return nil, nil, fmt.Errorf("no port mapping found for %s", port)
			}

			addrs = append(addrs, fmt.Sprintf("127.0.0.1:%s", mapped[0].HostPort))
		}
	}

	return &inspect, addrs, nil
}

func copyToContainer(ctx context.Context, dapi *client.Client, containerID, from, to string) error {
	srcInfo, err := archive.CopyInfoSourcePath(from, false)
	if err != nil {
		return fmt.Errorf("error copying from source %q: %v", from, err)
	}

	srcArchive, err := archive.TarResource(srcInfo)
	if err != nil {
		return fmt.Errorf("error creating tar from source %q: %v", from, err)
	}
	defer srcArchive.Close()

	dstInfo := archive.CopyInfo{Path: to}

	dstDir, content, err := archive.PrepareArchiveCopy(srcArchive, srcInfo, dstInfo)
	if err != nil {
		return fmt.Errorf("error preparing copy from %q -> %q: %v", from, to, err)
	}
	defer content.Close()
	err = dapi.CopyToContainer(ctx, containerID, dstDir, content, types.CopyToContainerOptions{})
	if err != nil {
		return fmt.Errorf("error copying from %q -> %q: %v", from, to, err)
	}

	return nil
}
