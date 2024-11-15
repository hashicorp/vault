// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package dockertest

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	dc "github.com/ory/dockertest/v3/docker"
	options "github.com/ory/dockertest/v3/docker/opts"
)

var (
	ErrNotInContainer = errors.New("not running in container")
)

// Pool represents a connection to the docker API and is used to create and remove docker images.
type Pool struct {
	Client  *dc.Client
	MaxWait time.Duration
}

// Network represents a docker network.
type Network struct {
	pool    *Pool
	Network *dc.Network
}

// Close removes network by calling pool.RemoveNetwork.
func (n *Network) Close() error {
	return n.pool.RemoveNetwork(n)
}

// Resource represents a docker container.
type Resource struct {
	pool      *Pool
	Container *dc.Container
}

// GetPort returns a resource's published port. You can use it to connect to the service via localhost, e.g. tcp://localhost:1231/
func (r *Resource) GetPort(id string) string {
	if r.Container == nil || r.Container.NetworkSettings == nil {
		return ""
	}

	m, ok := r.Container.NetworkSettings.Ports[dc.Port(id)]
	if !ok || len(m) == 0 {
		return ""
	}

	return m[0].HostPort
}

// GetBoundIP returns a resource's published IP address.
func (r *Resource) GetBoundIP(id string) string {
	if r.Container == nil || r.Container.NetworkSettings == nil {
		return ""
	}

	m, ok := r.Container.NetworkSettings.Ports[dc.Port(id)]
	if !ok || len(m) == 0 {
		return ""
	}

	ip := m[0].HostIP
	if ip == "0.0.0.0" || ip == "" {
		return "localhost"
	}
	return ip
}

// GetHostPort returns a resource's published port with an address.
func (r *Resource) GetHostPort(portID string) string {
	if r.Container == nil || r.Container.NetworkSettings == nil {
		return ""
	}

	m, ok := r.Container.NetworkSettings.Ports[dc.Port(portID)]
	if !ok || len(m) == 0 {
		return ""
	}

	ip := m[0].HostIP
	if ip == "0.0.0.0" || ip == "" {
		ip = "localhost"
	}
	return net.JoinHostPort(ip, m[0].HostPort)
}

type ExecOptions struct {
	// Command environment, optional.
	Env []string

	// StdIn will be attached as command stdin if provided.
	StdIn io.Reader

	// StdOut will be attached as command stdout if provided.
	StdOut io.Writer

	// StdErr will be attached as command stdout if provided.
	StdErr io.Writer

	// Allocate TTY for command or not.
	TTY bool
}

// Exec executes command within container.
func (r *Resource) Exec(cmd []string, opts ExecOptions) (exitCode int, err error) {
	exec, err := r.pool.Client.CreateExec(dc.CreateExecOptions{
		Container:    r.Container.ID,
		Cmd:          cmd,
		Env:          opts.Env,
		AttachStderr: true,
		AttachStdout: true,
		AttachStdin:  opts.StdIn != nil,
		Tty:          opts.TTY,
	})
	if err != nil {
		return -1, fmt.Errorf("Create exec failed: %w", err)
	}

	// Always attach stderr/stdout, even if not specified, to ensure that exec
	// waits with opts.Detach as false (default)
	// ref: https://github.com/fsouza/go-dockerclient/issues/838
	if opts.StdErr == nil {
		opts.StdErr = io.Discard
	}
	if opts.StdOut == nil {
		opts.StdOut = io.Discard
	}

	err = r.pool.Client.StartExec(exec.ID, dc.StartExecOptions{
		InputStream:  opts.StdIn,
		OutputStream: opts.StdOut,
		ErrorStream:  opts.StdErr,
		Tty:          opts.TTY,
	})
	if err != nil {
		return -1, fmt.Errorf("Start exec failed: %w", err)
	}

	inspectExec, err := r.pool.Client.InspectExec(exec.ID)
	if err != nil {
		return -1, fmt.Errorf("Inspect exec failed: %w", err)
	}

	return inspectExec.ExitCode, nil
}

// GetIPInNetwork returns container IP address in network.
func (r *Resource) GetIPInNetwork(network *Network) string {
	if r.Container == nil || r.Container.NetworkSettings == nil {
		return ""
	}

	netCfg, ok := r.Container.NetworkSettings.Networks[network.Network.Name]
	if !ok {
		return ""
	}

	return netCfg.IPAddress
}

// ConnectToNetwork connects container to network.
func (r *Resource) ConnectToNetwork(network *Network) error {
	err := r.pool.Client.ConnectNetwork(
		network.Network.ID,
		dc.NetworkConnectionOptions{Container: r.Container.ID},
	)
	if err != nil {
		return fmt.Errorf("Failed to connect container to network: %w", err)
	}

	// refresh internal representation
	r.Container, err = r.pool.Client.InspectContainer(r.Container.ID)
	if err != nil {
		return fmt.Errorf("Failed to refresh container information: %w", err)
	}

	network.Network, err = r.pool.Client.NetworkInfo(network.Network.ID)
	if err != nil {
		return fmt.Errorf("Failed to refresh network information: %w", err)
	}

	return nil
}

// DisconnectFromNetwork disconnects container from network.
func (r *Resource) DisconnectFromNetwork(network *Network) error {
	err := r.pool.Client.DisconnectNetwork(
		network.Network.ID,
		dc.NetworkConnectionOptions{Container: r.Container.ID},
	)
	if err != nil {
		return fmt.Errorf("Failed to connect container to network: %w", err)
	}

	// refresh internal representation
	r.Container, err = r.pool.Client.InspectContainer(r.Container.ID)
	if err != nil {
		return fmt.Errorf("Failed to refresh container information: %w", err)
	}

	network.Network, err = r.pool.Client.NetworkInfo(network.Network.ID)
	if err != nil {
		return fmt.Errorf("Failed to refresh network information: %w", err)
	}

	return nil
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
		return nil, err
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
				return nil, fmt.Errorf("failed to create client from environment: %w", err)
			}

			return &Pool{Client: client}, nil
		}
		if os.Getenv("DOCKER_HOST") != "" {
			endpoint = os.Getenv("DOCKER_HOST")
		} else if os.Getenv("DOCKER_URL") != "" {
			endpoint = os.Getenv("DOCKER_URL")
		} else if runtime.GOOS == "windows" {
			if _, err := os.Stat(`\\.\pipe\docker_engine`); err == nil {
				endpoint = "npipe:////./pipe/docker_engine"
			} else {
				endpoint = "http://localhost:2375"
			}
		} else {
			endpoint = options.DefaultHost
		}
	}

	if os.Getenv("DOCKER_CERT_PATH") != "" && shouldPreferTLS(endpoint) {
		return NewTLSPool(endpoint, os.Getenv("DOCKER_CERT_PATH"))
	}

	client, err := dc.NewClient(endpoint)
	if err != nil {
		return nil, err
	}

	return &Pool{
		Client: client,
	}, nil
}

func shouldPreferTLS(endpoint string) bool {
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
	Networks     []*Network // optional networks to join
	Labels       map[string]string
	Auth         dc.AuthConfiguration
	PortBindings map[dc.Port][]dc.PortBinding
	Privileged   bool
	User         string
	Tty          bool
	Platform     string
}

// BuildOptions is used to pass in optional parameters when building a container
type BuildOptions struct {
	Dockerfile string
	ContextDir string
	BuildArgs  []dc.BuildArg
	Platform   string
}

// BuildAndRunWithBuildOptions builds and starts a docker container.
// Optional modifier functions can be passed in order to change the hostconfig values not covered in RunOptions
func (d *Pool) BuildAndRunWithBuildOptions(buildOpts *BuildOptions, runOpts *RunOptions, hcOpts ...func(*dc.HostConfig)) (*Resource, error) {
	err := d.Client.BuildImage(dc.BuildImageOptions{
		Name:         runOpts.Name,
		Dockerfile:   buildOpts.Dockerfile,
		OutputStream: io.Discard,
		ContextDir:   buildOpts.ContextDir,
		BuildArgs:    buildOpts.BuildArgs,
		Platform:     buildOpts.Platform,
	})

	if err != nil {
		return nil, err
	}

	runOpts.Repository = runOpts.Name

	return d.RunWithOptions(runOpts, hcOpts...)
}

// BuildAndRunWithOptions builds and starts a docker container.
// Optional modifier functions can be passed in order to change the hostconfig values not covered in RunOptions
func (d *Pool) BuildAndRunWithOptions(dockerfilePath string, opts *RunOptions, hcOpts ...func(*dc.HostConfig)) (*Resource, error) {
	// Set the Dockerfile folder as build context
	dir, file := filepath.Split(dockerfilePath)
	buildOpts := BuildOptions{Dockerfile: file, ContextDir: dir}
	return d.BuildAndRunWithBuildOptions(&buildOpts, opts, hcOpts...)
}

// BuildAndRun builds and starts a docker container
func (d *Pool) BuildAndRun(name, dockerfilePath string, env []string) (*Resource, error) {
	return d.BuildAndRunWithOptions(dockerfilePath, &RunOptions{Name: name, Env: env})
}

// RunWithOptions starts a docker container.
// Optional modifier functions can be passed in order to change the hostconfig values not covered in RunOptions
//
//	 pool.RunWithOptions(&RunOptions{Repository: "mongo", Cmd: []string{"mongod", "--smallfiles"}})
//	 pool.RunWithOptions(&RunOptions{Repository: "mongo", Cmd: []string{"mongod", "--smallfiles"}}, func(hostConfig *dc.HostConfig) {
//				hostConfig.ShmSize = shmemsize
//			})
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
		s, d, err := options.MountParser(m)
		if err != nil {
			return nil, err
		}
		mounts = append(mounts, dc.Mount{
			Source:      s,
			Destination: d,
			RW:          true,
		})
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
	for _, network := range opts.Networks {
		networkingConfig.EndpointsConfig[network.Network.ID] = &dc.EndpointConfig{}
	}

	_, err := d.Client.InspectImage(fmt.Sprintf("%s:%s", repository, tag))
	if err != nil {
		if err := d.Client.PullImage(dc.PullImageOptions{
			Repository: repository,
			Tag:        tag,
			Platform:   opts.Platform,
		}, opts.Auth); err != nil {
			return nil, err
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
			User:         opts.User,
			Tty:          opts.Tty,
		},
		HostConfig:       &hostConfig,
		NetworkingConfig: &networkingConfig,
	})
	if err != nil {
		return nil, err
	}

	if err := d.Client.StartContainer(c.ID, nil); err != nil {
		return nil, err
	}

	c, err = d.Client.InspectContainer(c.ID)
	if err != nil {
		return nil, err
	}

	for _, network := range opts.Networks {
		network.Network, err = d.Client.NetworkInfo(network.Network.ID)
		if err != nil {
			return nil, err
		}
	}

	return &Resource{
		pool:      d,
		Container: c,
	}, nil
}

// Run starts a docker container.
//
//	pool.Run("mysql", "5.3", []string{"FOO=BAR", "BAR=BAZ"})
func (d *Pool) Run(repository, tag string, env []string) (*Resource, error) {
	return d.RunWithOptions(&RunOptions{Repository: repository, Tag: tag, Env: env})
}

// ContainerByName finds a container with the given name and returns it if present
func (d *Pool) ContainerByName(containerName string) (*Resource, bool) {
	containers, err := d.Client.ListContainers(dc.ListContainersOptions{
		All: true,
		Filters: map[string][]string{
			"name": {containerName},
		},
	})

	if err != nil {
		return nil, false
	}

	if len(containers) == 0 {
		return nil, false
	}

	c, err := d.Client.InspectContainer(containers[0].ID)
	if err != nil {
		return nil, false
	}

	return &Resource{
		pool:      d,
		Container: c,
	}, true
}

// RemoveContainerByName find a container with the given name and removes it if present
func (d *Pool) RemoveContainerByName(containerName string) error {
	containers, err := d.Client.ListContainers(dc.ListContainersOptions{
		All: true,
		Filters: map[string][]string{
			"name": {containerName},
		},
	})
	if err != nil {
		return fmt.Errorf("Error while listing containers with name %s: %w", containerName, err)
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
		return fmt.Errorf("Error while removing container with name %s: %w", containerName, err)
	}

	return nil
}

// Purge removes a container and linked volumes from docker.
func (d *Pool) Purge(r *Resource) error {
	if err := d.Client.RemoveContainer(dc.RemoveContainerOptions{ID: r.Container.ID, Force: true, RemoveVolumes: true}); err != nil {
		return err
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
	if err := backoff.Retry(op, bo); err != nil {
		if bo.NextBackOff() == backoff.Stop {
			return fmt.Errorf("reached retry deadline: %w", err)
		}

		return err
	}

	return nil
}

// CurrentContainer returns current container descriptor if this function called within running container.
// It returns ErrNotInContainer as error if this function running not in container.
func (d *Pool) CurrentContainer() (*Resource, error) {
	// docker daemon puts short container id into hostname
	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("Get hostname failed: %w", err)
	}

	container, err := d.Client.InspectContainer(hostname)
	switch err.(type) {
	case nil:
		return &Resource{
			pool:      d,
			Container: container,
		}, nil
	case *dc.NoSuchContainer:
		return nil, ErrNotInContainer
	default:
		return nil, err
	}
}

// CreateNetwork creates docker network. It's useful for linking multiple containers.
func (d *Pool) CreateNetwork(name string, opts ...func(config *dc.CreateNetworkOptions)) (*Network, error) {
	var cfg dc.CreateNetworkOptions
	cfg.Name = name
	for _, opt := range opts {
		opt(&cfg)
	}

	network, err := d.Client.CreateNetwork(cfg)
	if err != nil {
		return nil, err
	}

	return &Network{
		pool:    d,
		Network: network,
	}, nil
}

// NetworksByName returns a list of docker networks filtered by name
func (d *Pool) NetworksByName(name string) ([]Network, error) {
	networks, err := d.Client.ListNetworks()
	if err != nil {
		return nil, err
	}

	var foundNetworks []Network
	for idx := range networks {
		if networks[idx].Name == name {
			foundNetworks = append(foundNetworks,
				Network{
					pool:    d,
					Network: &networks[idx],
				},
			)
		}
	}

	return foundNetworks, nil
}

// RemoveNetwork disconnects containers and removes provided network.
func (d *Pool) RemoveNetwork(network *Network) error {
	for container := range network.Network.Containers {
		_ = d.Client.DisconnectNetwork(
			network.Network.ID,
			dc.NetworkConnectionOptions{Container: container, Force: true},
		)
	}

	return d.Client.RemoveNetwork(network.Network.ID)
}
