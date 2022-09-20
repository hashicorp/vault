package docker

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
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
	"github.com/docker/docker/api/types/strslice"
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
	Entrypoint      []string
	Env             []string
	NetworkID       string
	CopyFromTo      map[string]string
	Ports           []string
	DoNotAutoRemove bool
	AuthUsername    string
	AuthPassword    string
	LogConsumer     func(string)
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

// StartService will start the runner's configured docker container with a
// random UUID suffix appended to the name to make it unique and will return
// either a hostname or local address depending on if a Docker network was given.
//
// Most tests can default to using this.
func (d *Runner) StartService(ctx context.Context, connect ServiceAdapter) (*Service, error) {
	serv, _, err := d.StartNewService(ctx, true, false, connect)

	return serv, err
}

// StartNewService will start the runner's configured docker container but with the
// ability to control adding a name suffix or forcing a local address to be returned.
// 'addSuffix' will add a random UUID to the end of the container name.
// 'forceLocalAddr' will force the container address returned to be in the
// form of '127.0.0.1:1234' where 1234 is the mapped container port.
func (d *Runner) StartNewService(ctx context.Context, addSuffix, forceLocalAddr bool, connect ServiceAdapter) (*Service, string, error) {
	container, hostIPs, containerID, err := d.Start(context.Background(), addSuffix, forceLocalAddr)
	if err != nil {
		return nil, "", err
	}

	cleanup := func() {
		if d.RunOptions.LogConsumer != nil {
			rc, err := d.DockerAPI.ContainerLogs(ctx, container.ID, types.ContainerLogsOptions{
				ShowStdout: true,
				ShowStderr: true,
				Timestamps: true,
				Details:    true,
			})
			if err == nil {
				b, err := ioutil.ReadAll(rc)
				if err != nil {
					d.RunOptions.LogConsumer(fmt.Sprintf("error reading container logs, err=%v, read: %s", err, string(b)))
				} else {
					d.RunOptions.LogConsumer(string(b))
				}
			}
		}

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
	bo.MaxElapsedTime = 2 * time.Minute

	pieces := strings.Split(hostIPs[0], ":")
	portInt, err := strconv.Atoi(pieces[1])
	if err != nil {
		return nil, "", err
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
		return nil, "", err
	}

	return &Service{
		Config:    config,
		Cleanup:   cleanup,
		Container: container,
	}, containerID, nil
}

type Service struct {
	Config    ServiceConfig
	Cleanup   func()
	Container *types.ContainerJSON
}

func (d *Runner) Start(ctx context.Context, addSuffix, forceLocalAddr bool) (*types.ContainerJSON, []string, string, error) {
	name := d.RunOptions.ContainerName
	if addSuffix {
		suffix, err := uuid.GenerateUUID()
		if err != nil {
			return nil, nil, "", err
		}
		name += "-" + suffix
	}

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
	if len(d.RunOptions.Entrypoint) > 0 {
		cfg.Entrypoint = strslice.StrSlice(d.RunOptions.Entrypoint)
	}

	hostConfig := &container.HostConfig{
		AutoRemove:      !d.RunOptions.DoNotAutoRemove,
		PublishAllPorts: true,
	}

	netConfig := &network.NetworkingConfig{}
	if d.RunOptions.NetworkID != "" {
		netConfig.EndpointsConfig = map[string]*network.EndpointSettings{
			d.RunOptions.NetworkID: {},
		}
	}

	// best-effort pull
	var opts types.ImageCreateOptions
	if d.RunOptions.AuthUsername != "" && d.RunOptions.AuthPassword != "" {
		var buf bytes.Buffer
		auth := map[string]string{
			"username": d.RunOptions.AuthUsername,
			"password": d.RunOptions.AuthPassword,
		}
		if err := json.NewEncoder(&buf).Encode(auth); err != nil {
			return nil, nil, "", err
		}
		opts.RegistryAuth = base64.URLEncoding.EncodeToString(buf.Bytes())
	}
	resp, _ := d.DockerAPI.ImageCreate(ctx, cfg.Image, opts)
	if resp != nil {
		_, _ = ioutil.ReadAll(resp)
	}

	c, err := d.DockerAPI.ContainerCreate(ctx, cfg, hostConfig, netConfig, nil, cfg.Hostname)
	if err != nil {
		return nil, nil, "", fmt.Errorf("container create failed: %v", err)
	}

	for from, to := range d.RunOptions.CopyFromTo {
		if err := copyToContainer(ctx, d.DockerAPI, c.ID, from, to); err != nil {
			_ = d.DockerAPI.ContainerRemove(ctx, c.ID, types.ContainerRemoveOptions{})
			return nil, nil, "", err
		}
	}

	err = d.DockerAPI.ContainerStart(ctx, c.ID, types.ContainerStartOptions{})
	if err != nil {
		_ = d.DockerAPI.ContainerRemove(ctx, c.ID, types.ContainerRemoveOptions{})
		return nil, nil, "", fmt.Errorf("container start failed: %v", err)
	}

	inspect, err := d.DockerAPI.ContainerInspect(ctx, c.ID)
	if err != nil {
		_ = d.DockerAPI.ContainerRemove(ctx, c.ID, types.ContainerRemoveOptions{})
		return nil, nil, "", err
	}

	var addrs []string
	for _, port := range d.RunOptions.Ports {
		pieces := strings.Split(port, "/")
		if len(pieces) < 2 {
			return nil, nil, "", fmt.Errorf("expected port of the form 1234/tcp, got: %s", port)
		}
		if d.RunOptions.NetworkID != "" && !forceLocalAddr {
			addrs = append(addrs, fmt.Sprintf("%s:%s", cfg.Hostname, pieces[0]))
		} else {
			mapped, ok := inspect.NetworkSettings.Ports[nat.Port(port)]
			if !ok || len(mapped) == 0 {
				return nil, nil, "", fmt.Errorf("no port mapping found for %s", port)
			}
			addrs = append(addrs, fmt.Sprintf("127.0.0.1:%s", mapped[0].HostPort))
		}
	}

	return &inspect, addrs, c.ID, nil
}

func (d *Runner) Stop(ctx context.Context, containerID string) error {
	timeout := 5 * time.Second
	if err := d.DockerAPI.ContainerStop(ctx, containerID, &timeout); err != nil {
		return err
	}

	return d.DockerAPI.NetworkDisconnect(ctx, d.RunOptions.NetworkID, containerID, true)
}

func (d *Runner) Restart(ctx context.Context, containerID string) error {
	if err := d.DockerAPI.ContainerStart(ctx, containerID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	ends := &network.EndpointSettings{
		NetworkID: d.RunOptions.NetworkID,
	}

	return d.DockerAPI.NetworkConnect(ctx, d.RunOptions.NetworkID, containerID, ends)
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

type RunCmdOpt interface {
	Apply(cfg *types.ExecConfig) error
}

type RunCmdUser string

var _ RunCmdOpt = (*RunCmdUser)(nil)

func (u RunCmdUser) Apply(cfg *types.ExecConfig) error {
	cfg.User = string(u)
	return nil
}

func (d *Runner) RunCmdWithOutput(ctx context.Context, container string, cmd []string, opts ...RunCmdOpt) ([]byte, int, error) {
	runCfg := types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          cmd,
	}

	for index, opt := range opts {
		if err := opt.Apply(&runCfg); err != nil {
			return nil, -1, fmt.Errorf("error applying option (%d / %v): %v", index, opt, err)
		}
	}

	ret, err := d.DockerAPI.ContainerExecCreate(ctx, container, runCfg)
	if err != nil {
		return nil, -1, fmt.Errorf("error creating execution environment: %v\ncfg: %v\n", err, runCfg)
	}

	resp, err := d.DockerAPI.ContainerExecAttach(ctx, ret.ID, types.ExecStartCheck{})
	if err != nil {
		return nil, -1, fmt.Errorf("error attaching to command execution: %v\ncfg: %v\nret: %v\n", err, runCfg, ret)
	}

	result, err := ioutil.ReadAll(resp.Reader)
	resp.Close()
	if err != nil {
		return nil, -1, fmt.Errorf("error reading command response: %v", err)
	}

	// Fetch return code.
	info, err := d.DockerAPI.ContainerExecInspect(ctx, ret.ID)
	if err != nil {
		return result, -1, fmt.Errorf("error reading command exit code: %v", err)
	}

	return result, info.ExitCode, nil
}

func (d *Runner) RunCmdInBackground(ctx context.Context, container string, cmd []string, opts ...RunCmdOpt) (string, error) {
	runCfg := types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          cmd,
	}

	for index, opt := range opts {
		if err := opt.Apply(&runCfg); err != nil {
			return "", fmt.Errorf("error applying option (%d / %v): %v", index, opt, err)
		}
	}

	ret, err := d.DockerAPI.ContainerExecCreate(ctx, container, runCfg)
	if err != nil {
		return "", fmt.Errorf("error creating execution environment: %v\ncfg: %v\n", err, runCfg)
	}

	err = d.DockerAPI.ContainerExecStart(ctx, ret.ID, types.ExecStartCheck{})
	if err != nil {
		return "", fmt.Errorf("error starting command execution: %v\ncfg: %v\nret: %v\n", err, runCfg, ret)
	}

	return ret.ID, nil
}
