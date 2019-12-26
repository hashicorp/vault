package dockertest

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/cenkalti/backoff"
	dc "github.com/ory/dockertest/docker"
	"github.com/pkg/errors"
)

// Pool represents a connection to the docker API and is used to create and remove docker images.
type Pool struct {
	Client  *dc.Client
	MaxWait time.Duration
}

// Resource represents a docker container.
type Resource struct {
	pool      *Pool
	Container *dc.Container
}

// GetPort returns a resource's published port. You can use it to connect to the service via localhost, e.g. tcp://localhost:1231/
func (r *Resource) GetPort(id string) string {
	if r.Container == nil {
		return ""
	} else if r.Container.NetworkSettings == nil {
		return ""
	}

	m, ok := r.Container.NetworkSettings.Ports[dc.Port(id)]
	if !ok {
		return ""
	} else if len(m) == 0 {
		return ""
	}

	return m[0].HostPort
}

func (r *Resource) GetBoundIP(id string) string {
	if r.Container == nil {
		return ""
	} else if r.Container.NetworkSettings == nil {
		return ""
	}

	m, ok := r.Container.NetworkSettings.Ports[dc.Port(id)]
	if !ok {
		return ""
	} else if len(m) == 0 {
		return ""
	}

	return m[0].HostIP
}

// GetHostPort returns a resource's published port with an address.
func (r *Resource) GetHostPort(portID string) string {
	if r.Container == nil {
		return ""
	} else if r.Container.NetworkSettings == nil {
		return ""
	}

	m, ok := r.Container.NetworkSettings.Ports[dc.Port(portID)]
	if !ok {
		return ""
	} else if len(m) == 0 {
		return ""
	}
	ip := m[0].HostIP
	if ip == "0.0.0.0" {
		ip = "localhost"
	}
	return net.JoinHostPort(ip, m[0].HostPort)
}

// Close removes a container and linked volumes from docker by calling pool.Purge.
func (r *Resource) Close() error {
	return r.pool.Purge(r)
}

// Expire sets a resource's associated container to terminate after a period has passed
func (r *Resource) Expire(seconds uint) error {
	go func() {
		if err := r.pool.Client.StopContainer(r.Container.ID, seconds); err != nil {
			// Error handling?
		}
	}()
	return nil
}

// NewTLSPool creates a new pool given an endpoint and the certificate path. This is required for endpoints that
// require TLS communication.
func NewTLSPool(endpoint, certpath string) (*Pool, error) {
	ca := fmt.Sprintf("%s/ca.pem", certpath)
	cert := fmt.Sprintf("%s/cert.pem", certpath)
	key := fmt.Sprintf("%s/key.pem", certpath)

	client, err := dc.NewTLSClient(endpoint, cert, key, ca)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	return &Pool{
		Client: client,
	}, nil
}

// NewPool creates a new pool. You can pass an empty string to use the default, which is taken from the environment
// variable DOCKER_HOST and DOCKER_URL, or from docker-machine if the environment variable DOCKER_MACHINE_NAME is set,
// or if neither is defined a sensible default for the operating system you are on.
// TLS pools are automatically configured if the DOCKER_CERT_PATH environment variable exists.
func NewPool(endpoint string) (*Pool, error) {
	if endpoint == "" {
		if os.Getenv("DOCKER_MACHINE_NAME") != "" {
			client, err := dc.NewClientFromEnv()
			if err != nil {
				return nil, errors.Wrap(err, "failed to create client from environment")
			}

			return &Pool{Client: client}, nil
		} else if os.Getenv("DOCKER_HOST") != "" {
			endpoint = os.Getenv("DOCKER_HOST")
		} else if os.Getenv("DOCKER_URL") != "" {
			endpoint = os.Getenv("DOCKER_URL")
		} else if runtime.GOOS == "windows" {
			endpoint = "http://localhost:2375"
		} else {
			endpoint = "unix:///var/run/docker.sock"
		}
	}

	if os.Getenv("DOCKER_CERT_PATH") != "" && shouldPreferTls(endpoint) {
		return NewTLSPool(endpoint, os.Getenv("DOCKER_CERT_PATH"))
	}

	client, err := dc.NewClient(endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	return &Pool{
		Client: client,
	}, nil
}

func shouldPreferTls(endpoint string) bool {
	return !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "unix://")
}

// RunOptions is used to pass in optional parameters when running a container.
type RunOptions struct {
	Hostname     string
	Name         string
	Repository   string
	Tag          string
	Env          []string
	Entrypoint   []string
	Cmd          []string
	Mounts       []string
	Links        []string
	ExposedPorts []string
	ExtraHosts   []string
	CapAdd       []string
	SecurityOpt  []string
	DNS          []string
	WorkingDir   string
	NetworkID    string
	Labels       map[string]string
	Auth         dc.AuthConfiguration
	PortBindings map[dc.Port][]dc.PortBinding
	Privileged   bool
}

// BuildOptions is used to pass in optional parameters when building a container
type BuildOptions struct {
	Dockerfile string
	ContextDir string
}

// BuildAndRunWithBuildOptions builds and starts a docker container.
// Optional modifier functions can be passed in order to change the hostconfig values not covered in RunOptions
func (d *Pool) BuildAndRunWithBuildOptions(buildOpts *BuildOptions, runOpts *RunOptions, hcOpts ...func(*dc.HostConfig)) (*Resource, error) {
	err := d.Client.BuildImage(dc.BuildImageOptions{
		Name:         runOpts.Name,
		Dockerfile:   buildOpts.Dockerfile,
		OutputStream: ioutil.Discard,
		ContextDir:   buildOpts.ContextDir,
	})

	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	runOpts.Repository = runOpts.Name

	return d.RunWithOptions(runOpts, hcOpts...)
}

// BuildAndRunWithOptions builds and starts a docker container.
// Optional modifier functions can be passed in order to change the hostconfig values not covered in RunOptions
func (d *Pool) BuildAndRunWithOptions(dockerfilePath string, opts *RunOptions, hcOpts ...func(*dc.HostConfig)) (*Resource, error) {
	// Set the Dockerfile folder as build context
	dir, file := filepath.Split(dockerfilePath)
	buildOpts := BuildOptions{Dockerfile:file, ContextDir:dir}
	return d.BuildAndRunWithBuildOptions(&buildOpts, opts, hcOpts...)
}


// BuildAndRun builds and starts a docker container
func (d *Pool) BuildAndRun(name, dockerfilePath string, env []string) (*Resource, error) {
	return d.BuildAndRunWithOptions(dockerfilePath, &RunOptions{Name: name, Env: env})
}

// RunWithOptions starts a docker container.
// Optional modifier functions can be passed in order to change the hostconfig values not covered in RunOptions
//
// pool.Run(&RunOptions{Repository: "mongo", Cmd: []string{"mongod", "--smallfiles"}})
// pool.Run(&RunOptions{Repository: "mongo", Cmd: []string{"mongod", "--smallfiles"}}, func(hostConfig *dc.HostConfig) {
//			hostConfig.ShmSize = shmemsize
//		})
func (d *Pool) RunWithOptions(opts *RunOptions, hcOpts ...func(*dc.HostConfig)) (*Resource, error) {
	repository := opts.Repository
	tag := opts.Tag
	env := opts.Env
	cmd := opts.Cmd
	ep := opts.Entrypoint
	wd := opts.WorkingDir
	var exp map[dc.Port]struct{}

	if len(opts.ExposedPorts) > 0 {
		exp = map[dc.Port]struct{}{}
		for _, p := range opts.ExposedPorts {
			exp[dc.Port(p)] = struct{}{}
		}
	}

	mounts := []dc.Mount{}

	for _, m := range opts.Mounts {
		sd := strings.Split(m, ":")
		if len(sd) == 2 {
			mounts = append(mounts, dc.Mount{
				Source:      sd[0],
				Destination: sd[1],
				RW:          true,
			})
		} else {
			return nil, errors.Wrap(fmt.Errorf("invalid mount format: got %s, expected <src>:<dst>", m), "")
		}
	}

	if tag == "" {
		tag = "latest"
	}

	networkingConfig := dc.NetworkingConfig{
		EndpointsConfig: map[string]*dc.EndpointConfig{},
	}
	if opts.NetworkID != "" {
		networkingConfig.EndpointsConfig[opts.NetworkID] = &dc.EndpointConfig{}
	}

	_, err := d.Client.InspectImage(fmt.Sprintf("%s:%s", repository, tag))
	if err != nil {
		if err := d.Client.PullImage(dc.PullImageOptions{
			Repository: repository,
			Tag:        tag,
		}, opts.Auth); err != nil {
			return nil, errors.Wrap(err, "")
		}
	}

	hostConfig := dc.HostConfig{
		PublishAllPorts: true,
		Binds:           opts.Mounts,
		Links:           opts.Links,
		PortBindings:    opts.PortBindings,
		ExtraHosts:      opts.ExtraHosts,
		CapAdd:          opts.CapAdd,
		SecurityOpt:     opts.SecurityOpt,
		Privileged:      opts.Privileged,
		DNS:             opts.DNS,
	}

	for _, hostConfigOption := range hcOpts {
		hostConfigOption(&hostConfig)
	}

	c, err := d.Client.CreateContainer(dc.CreateContainerOptions{
		Name: opts.Name,
		Config: &dc.Config{
			Hostname:     opts.Hostname,
			Image:        fmt.Sprintf("%s:%s", repository, tag),
			Env:          env,
			Entrypoint:   ep,
			Cmd:          cmd,
			Mounts:       mounts,
			ExposedPorts: exp,
			WorkingDir:   wd,
			Labels:       opts.Labels,
			StopSignal:   "SIGWINCH", // to support timeouts
		},
		HostConfig:       &hostConfig,
		NetworkingConfig: &networkingConfig,
	})
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	if err := d.Client.StartContainer(c.ID, nil); err != nil {
		return nil, errors.Wrap(err, "")
	}

	c, err = d.Client.InspectContainer(c.ID)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	return &Resource{
		pool:      d,
		Container: c,
	}, nil
}

// Run starts a docker container.
//
// pool.Run("mysql", "5.3", []string{"FOO=BAR", "BAR=BAZ"})
func (d *Pool) Run(repository, tag string, env []string) (*Resource, error) {
	return d.RunWithOptions(&RunOptions{Repository: repository, Tag: tag, Env: env})
}

// RemoveContainerByName find a container with the given name and removes it if present
func (d *Pool) RemoveContainerByName(containerName string) error {
	containers, err := d.Client.ListContainers(dc.ListContainersOptions{
		All: true,
		Filters: map[string][]string{
			"name": []string{containerName},
		},
	})
	if err != nil {
		return errors.Wrapf(err, "Error while listing containers with name %s", containerName)
	}

	if len(containers) == 0 {
		return nil
	}

	err = d.Client.RemoveContainer(dc.RemoveContainerOptions{
		ID:            containers[0].ID,
		Force:         true,
		RemoveVolumes: true,
	})
	if err != nil {
		return errors.Wrapf(err, "Error while removing container with name %s", containerName)
	}

	return nil
}

// Purge removes a container and linked volumes from docker.
func (d *Pool) Purge(r *Resource) error {
	if err := d.Client.RemoveContainer(dc.RemoveContainerOptions{ID: r.Container.ID, Force: true, RemoveVolumes: true}); err != nil {
		return errors.Wrap(err, "")
	}

	return nil
}

// Retry is an exponential backoff retry helper. You can use it to wait for e.g. mysql to boot up.
func (d *Pool) Retry(op func() error) error {
	if d.MaxWait == 0 {
		d.MaxWait = time.Minute
	}
	bo := backoff.NewExponentialBackOff()
	bo.MaxInterval = time.Second * 5
	bo.MaxElapsedTime = d.MaxWait
	return backoff.Retry(op, bo)
}
