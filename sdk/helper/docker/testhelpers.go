// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package docker

import (
	"archive/tar"
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v3"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
	"github.com/hashicorp/go-uuid"
)

const DockerAPIVersion = "1.40"

type Runner struct {
	DockerAPI  *client.Client
	RunOptions RunOptions
}

type RunOptions struct {
	ImageRepo              string
	ImageTag               string
	ContainerName          string
	Cmd                    []string
	Entrypoint             []string
	Env                    []string
	NetworkName            string
	NetworkID              string
	CopyFromTo             map[string]string
	Ports                  []string
	DoNotAutoRemove        bool
	AuthUsername           string
	AuthPassword           string
	OmitLogTimestamps      bool
	LogConsumer            func(string)
	Capabilities           []string
	PreDelete              bool
	PostStart              func(string, string) error
	LogStderr              io.Writer
	LogStdout              io.Writer
	VolumeNameToMountPoint map[string]string
}

func NewDockerAPI() (*client.Client, error) {
	return client.NewClientWithOpts(client.FromEnv, client.WithVersion(DockerAPIVersion))
}

func NewServiceRunner(opts RunOptions) (*Runner, error) {
	dapi, err := NewDockerAPI()
	if err != nil {
		return nil, err
	}

	if opts.NetworkName == "" {
		opts.NetworkName = os.Getenv("TEST_DOCKER_NETWORK_NAME")
	}
	if opts.NetworkName != "" {
		nets, err := dapi.NetworkList(context.TODO(), types.NetworkListOptions{
			Filters: filters.NewArgs(filters.Arg("name", opts.NetworkName)),
		})
		if err != nil {
			return nil, err
		}
		if len(nets) != 1 {
			return nil, fmt.Errorf("expected exactly one docker network named %q, got %d", opts.NetworkName, len(nets))
		}
		opts.NetworkID = nets[0].ID
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

type LogConsumerWriter struct {
	consumer func(string)
}

func (l LogConsumerWriter) Write(p []byte) (n int, err error) {
	// TODO this assumes that we're never passed partial log lines, which
	// seems a safe assumption for now based on how docker looks to implement
	// logging, but might change in the future.
	scanner := bufio.NewScanner(bytes.NewReader(p))
	scanner.Buffer(make([]byte, 64*1024), bufio.MaxScanTokenSize)
	for scanner.Scan() {
		l.consumer(scanner.Text())
	}
	return len(p), nil
}

var _ io.Writer = &LogConsumerWriter{}

// StartNewService will start the runner's configured docker container but with the
// ability to control adding a name suffix or forcing a local address to be returned.
// 'addSuffix' will add a random UUID to the end of the container name.
// 'forceLocalAddr' will force the container address returned to be in the
// form of '127.0.0.1:1234' where 1234 is the mapped container port.
func (d *Runner) StartNewService(ctx context.Context, addSuffix, forceLocalAddr bool, connect ServiceAdapter) (*Service, string, error) {
	if d.RunOptions.PreDelete {
		name := d.RunOptions.ContainerName
		matches, err := d.DockerAPI.ContainerList(ctx, types.ContainerListOptions{
			All: true,
			// TODO use labels to ensure we don't delete anything we shouldn't
			Filters: filters.NewArgs(
				filters.Arg("name", name),
			),
		})
		if err != nil {
			return nil, "", fmt.Errorf("failed to list containers named %q", name)
		}
		for _, cont := range matches {
			err = d.DockerAPI.ContainerRemove(ctx, cont.ID, types.ContainerRemoveOptions{Force: true})
			if err != nil {
				return nil, "", fmt.Errorf("failed to pre-delete container named %q", name)
			}
		}
	}
	result, err := d.Start(context.Background(), addSuffix, forceLocalAddr)
	if err != nil {
		return nil, "", err
	}

	// The waitgroup wg is used here to support some stuff in NewDockerCluster.
	// We can't generate the PKI cert for the https listener until we know the
	// container's address, meaning we must first start the container, then
	// generate the cert, then copy it into the container, then signal Vault
	// to reload its config/certs.  However, if we SIGHUP Vault before Vault
	// has installed its signal handler, that will kill Vault, since the default
	// behaviour for HUP is termination.  So the PostStart that NewDockerCluster
	// passes in (which does all that PKI cert stuff) waits to see output from
	// Vault on stdout/stderr before it sends the signal, and we don't want to
	// run the PostStart until we've hooked into the docker logs.
	var wg sync.WaitGroup
	logConsumer := d.createLogConsumer(result.Container.ID, &wg)

	if logConsumer != nil {
		wg.Add(1)
		go logConsumer()
	}
	wg.Wait()

	if d.RunOptions.PostStart != nil {
		if err := d.RunOptions.PostStart(result.Container.ID, result.RealIP); err != nil {
			return nil, "", fmt.Errorf("poststart failed: %w", err)
		}
	}

	cleanup := func() {
		for i := 0; i < 10; i++ {
			err := d.DockerAPI.ContainerRemove(ctx, result.Container.ID, types.ContainerRemoveOptions{Force: true})
			if err == nil || client.IsErrNotFound(err) {
				return
			}
			time.Sleep(1 * time.Second)
		}
	}

	bo := backoff.NewExponentialBackOff()
	bo.MaxInterval = time.Second * 5
	bo.MaxElapsedTime = 2 * time.Minute

	pieces := strings.Split(result.Addrs[0], ":")
	portInt, err := strconv.Atoi(pieces[1])
	if err != nil {
		return nil, "", err
	}

	var config ServiceConfig
	err = backoff.Retry(func() error {
		container, err := d.DockerAPI.ContainerInspect(ctx, result.Container.ID)
		if err != nil || !container.State.Running {
			return backoff.Permanent(fmt.Errorf("failed inspect or container %q not running: %w", result.Container.ID, err))
		}

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
		Config:      config,
		Cleanup:     cleanup,
		Container:   result.Container,
		StartResult: result,
	}, result.Container.ID, nil
}

// createLogConsumer returns a function to consume the logs of the container with the given ID.
// If a wait group is given, `WaitGroup.Done()` will be called as soon as the call to the
// ContainerLogs Docker API call is done.
// The returned function will block, so it should be run on a goroutine.
func (d *Runner) createLogConsumer(containerId string, wg *sync.WaitGroup) func() {
	if d.RunOptions.LogStdout != nil && d.RunOptions.LogStderr != nil {
		return func() {
			d.consumeLogs(containerId, wg, d.RunOptions.LogStdout, d.RunOptions.LogStderr)
		}
	}
	if d.RunOptions.LogConsumer != nil {
		return func() {
			d.consumeLogs(containerId, wg, &LogConsumerWriter{d.RunOptions.LogConsumer}, &LogConsumerWriter{d.RunOptions.LogConsumer})
		}
	}
	return nil
}

// consumeLogs is the function called by the function returned by createLogConsumer.
func (d *Runner) consumeLogs(containerId string, wg *sync.WaitGroup, logStdout, logStderr io.Writer) {
	// We must run inside a goroutine because we're using Follow:true,
	// and StdCopy will block until the log stream is closed.
	stream, err := d.DockerAPI.ContainerLogs(context.Background(), containerId, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Timestamps: !d.RunOptions.OmitLogTimestamps,
		Details:    true,
		Follow:     true,
	})
	wg.Done()
	if err != nil {
		d.RunOptions.LogConsumer(fmt.Sprintf("error reading container logs: %v", err))
	} else {
		_, err := stdcopy.StdCopy(logStdout, logStderr, stream)
		if err != nil {
			d.RunOptions.LogConsumer(fmt.Sprintf("error demultiplexing docker logs: %v", err))
		}
	}
}

type Service struct {
	Config      ServiceConfig
	Cleanup     func()
	Container   *types.ContainerJSON
	StartResult *StartResult
}

type StartResult struct {
	Container *types.ContainerJSON
	Addrs     []string
	RealIP    string
}

func (d *Runner) Start(ctx context.Context, addSuffix, forceLocalAddr bool) (*StartResult, error) {
	name := d.RunOptions.ContainerName
	if addSuffix {
		suffix, err := uuid.GenerateUUID()
		if err != nil {
			return nil, err
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
	if len(d.RunOptions.Capabilities) > 0 {
		hostConfig.CapAdd = d.RunOptions.Capabilities
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
			return nil, err
		}
		opts.RegistryAuth = base64.URLEncoding.EncodeToString(buf.Bytes())
	}
	resp, _ := d.DockerAPI.ImageCreate(ctx, cfg.Image, opts)
	if resp != nil {
		_, _ = ioutil.ReadAll(resp)
	}

	for vol, mtpt := range d.RunOptions.VolumeNameToMountPoint {
		hostConfig.Mounts = append(hostConfig.Mounts, mount.Mount{
			Type:     "volume",
			Source:   vol,
			Target:   mtpt,
			ReadOnly: false,
		})
	}

	c, err := d.DockerAPI.ContainerCreate(ctx, cfg, hostConfig, netConfig, nil, cfg.Hostname)
	if err != nil {
		return nil, fmt.Errorf("container create failed: %v", err)
	}

	for from, to := range d.RunOptions.CopyFromTo {
		if err := copyToContainer(ctx, d.DockerAPI, c.ID, from, to); err != nil {
			_ = d.DockerAPI.ContainerRemove(ctx, c.ID, types.ContainerRemoveOptions{})
			return nil, err
		}
	}

	err = d.DockerAPI.ContainerStart(ctx, c.ID, types.ContainerStartOptions{})
	if err != nil {
		_ = d.DockerAPI.ContainerRemove(ctx, c.ID, types.ContainerRemoveOptions{})
		return nil, fmt.Errorf("container start failed: %v", err)
	}

	inspect, err := d.DockerAPI.ContainerInspect(ctx, c.ID)
	if err != nil {
		_ = d.DockerAPI.ContainerRemove(ctx, c.ID, types.ContainerRemoveOptions{})
		return nil, err
	}

	var addrs []string
	for _, port := range d.RunOptions.Ports {
		pieces := strings.Split(port, "/")
		if len(pieces) < 2 {
			return nil, fmt.Errorf("expected port of the form 1234/tcp, got: %s", port)
		}
		if d.RunOptions.NetworkID != "" && !forceLocalAddr {
			addrs = append(addrs, fmt.Sprintf("%s:%s", cfg.Hostname, pieces[0]))
		} else {
			mapped, ok := inspect.NetworkSettings.Ports[nat.Port(port)]
			if !ok || len(mapped) == 0 {
				return nil, fmt.Errorf("no port mapping found for %s", port)
			}
			addrs = append(addrs, fmt.Sprintf("127.0.0.1:%s", mapped[0].HostPort))
		}
	}

	var realIP string
	if d.RunOptions.NetworkID == "" {
		if len(inspect.NetworkSettings.Networks) > 1 {
			return nil, fmt.Errorf("Set d.RunOptions.NetworkName instead for container with multiple networks: %v", inspect.NetworkSettings.Networks)
		}
		for _, network := range inspect.NetworkSettings.Networks {
			realIP = network.IPAddress
			break
		}
	} else {
		realIP = inspect.NetworkSettings.Networks[d.RunOptions.NetworkName].IPAddress
	}

	return &StartResult{
		Container: &inspect,
		Addrs:     addrs,
		RealIP:    realIP,
	}, nil
}

func (d *Runner) RefreshFiles(ctx context.Context, containerID string) error {
	for from, to := range d.RunOptions.CopyFromTo {
		if err := copyToContainer(ctx, d.DockerAPI, containerID, from, to); err != nil {
			// TODO too drastic?
			_ = d.DockerAPI.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{})
			return err
		}
	}
	return d.DockerAPI.ContainerKill(ctx, containerID, "SIGHUP")
}

func (d *Runner) Stop(ctx context.Context, containerID string) error {
	if d.RunOptions.NetworkID != "" {
		if err := d.DockerAPI.NetworkDisconnect(ctx, d.RunOptions.NetworkID, containerID, true); err != nil {
			return fmt.Errorf("error disconnecting network (%v): %v", d.RunOptions.NetworkID, err)
		}
	}

	// timeout in seconds
	timeout := 5
	options := container.StopOptions{
		Timeout: &timeout,
	}
	if err := d.DockerAPI.ContainerStop(ctx, containerID, options); err != nil {
		return fmt.Errorf("error stopping container: %v", err)
	}

	return nil
}

func (d *Runner) RestartContainerWithTimeout(ctx context.Context, containerID string, timeout int) error {
	err := d.DockerAPI.ContainerRestart(ctx, containerID, container.StopOptions{Timeout: &timeout})
	if err != nil {
		return fmt.Errorf("failed to restart container: %s", err)
	}
	var wg sync.WaitGroup
	logConsumer := d.createLogConsumer(containerID, &wg)
	if logConsumer != nil {
		wg.Add(1)
		go logConsumer()
	}
	// we don't really care about waiting for logs to start showing up, do we?
	return nil
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

func (d *Runner) RunCmdWithOutput(ctx context.Context, container string, cmd []string, opts ...RunCmdOpt) ([]byte, []byte, int, error) {
	return RunCmdWithOutput(d.DockerAPI, ctx, container, cmd, opts...)
}

func RunCmdWithOutput(api *client.Client, ctx context.Context, container string, cmd []string, opts ...RunCmdOpt) ([]byte, []byte, int, error) {
	runCfg := types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          cmd,
	}

	for index, opt := range opts {
		if err := opt.Apply(&runCfg); err != nil {
			return nil, nil, -1, fmt.Errorf("error applying option (%d / %v): %w", index, opt, err)
		}
	}

	ret, err := api.ContainerExecCreate(ctx, container, runCfg)
	if err != nil {
		return nil, nil, -1, fmt.Errorf("error creating execution environment: %v\ncfg: %v\n", err, runCfg)
	}

	resp, err := api.ContainerExecAttach(ctx, ret.ID, types.ExecStartCheck{})
	if err != nil {
		return nil, nil, -1, fmt.Errorf("error attaching to command execution: %v\ncfg: %v\nret: %v\n", err, runCfg, ret)
	}
	defer resp.Close()

	var stdoutB bytes.Buffer
	var stderrB bytes.Buffer
	if _, err := stdcopy.StdCopy(&stdoutB, &stderrB, resp.Reader); err != nil {
		return nil, nil, -1, fmt.Errorf("error reading command output: %v", err)
	}

	stdout := stdoutB.Bytes()
	stderr := stderrB.Bytes()

	// Fetch return code.
	info, err := api.ContainerExecInspect(ctx, ret.ID)
	if err != nil {
		return stdout, stderr, -1, fmt.Errorf("error reading command exit code: %v", err)
	}

	return stdout, stderr, info.ExitCode, nil
}

func (d *Runner) RunCmdInBackground(ctx context.Context, container string, cmd []string, opts ...RunCmdOpt) (string, error) {
	return RunCmdInBackground(d.DockerAPI, ctx, container, cmd, opts...)
}

func RunCmdInBackground(api *client.Client, ctx context.Context, container string, cmd []string, opts ...RunCmdOpt) (string, error) {
	runCfg := types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          cmd,
	}

	for index, opt := range opts {
		if err := opt.Apply(&runCfg); err != nil {
			return "", fmt.Errorf("error applying option (%d / %v): %w", index, opt, err)
		}
	}

	ret, err := api.ContainerExecCreate(ctx, container, runCfg)
	if err != nil {
		return "", fmt.Errorf("error creating execution environment: %w\ncfg: %v\n", err, runCfg)
	}

	err = api.ContainerExecStart(ctx, ret.ID, types.ExecStartCheck{})
	if err != nil {
		return "", fmt.Errorf("error starting command execution: %w\ncfg: %v\nret: %v\n", err, runCfg, ret)
	}

	return ret.ID, nil
}

// Mapping of path->contents
type PathContents interface {
	UpdateHeader(header *tar.Header) error
	Get() ([]byte, error)
	SetMode(mode int64)
	SetOwners(uid int, gid int)
}

type FileContents struct {
	Data []byte
	Mode int64
	UID  int
	GID  int
}

func (b FileContents) UpdateHeader(header *tar.Header) error {
	header.Mode = b.Mode
	header.Uid = b.UID
	header.Gid = b.GID
	return nil
}

func (b FileContents) Get() ([]byte, error) {
	return b.Data, nil
}

func (b *FileContents) SetMode(mode int64) {
	b.Mode = mode
}

func (b *FileContents) SetOwners(uid int, gid int) {
	b.UID = uid
	b.GID = gid
}

func PathContentsFromBytes(data []byte) PathContents {
	return &FileContents{
		Data: data,
		Mode: 0o644,
	}
}

func PathContentsFromString(data string) PathContents {
	return PathContentsFromBytes([]byte(data))
}

type BuildContext map[string]PathContents

func NewBuildContext() BuildContext {
	return BuildContext{}
}

func BuildContextFromTarball(reader io.Reader) (BuildContext, error) {
	archive := tar.NewReader(reader)
	bCtx := NewBuildContext()

	for true {
		header, err := archive.Next()
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, fmt.Errorf("failed to parse provided tarball: %v", err)
		}

		data := make([]byte, int(header.Size))
		read, err := archive.Read(data)
		if err != nil {
			return nil, fmt.Errorf("failed to parse read from provided tarball: %v", err)
		}

		if read != int(header.Size) {
			return nil, fmt.Errorf("unexpectedly short read on tarball: %v of %v", read, header.Size)
		}

		bCtx[header.Name] = &FileContents{
			Data: data,
			Mode: header.Mode,
			UID:  header.Uid,
			GID:  header.Gid,
		}
	}

	return bCtx, nil
}

func (bCtx *BuildContext) ToTarball() (io.Reader, error) {
	var err error
	buffer := new(bytes.Buffer)
	tarBuilder := tar.NewWriter(buffer)
	defer tarBuilder.Close()

	now := time.Now()
	for filepath, contents := range *bCtx {
		fileHeader := &tar.Header{
			Name:       filepath,
			ModTime:    now,
			AccessTime: now,
			ChangeTime: now,
		}
		if contents == nil && !strings.HasSuffix(filepath, "/") {
			return nil, fmt.Errorf("expected file path (%v) to have trailing / due to nil contents, indicating directory", filepath)
		}

		if err := contents.UpdateHeader(fileHeader); err != nil {
			return nil, fmt.Errorf("failed to update tar header entry for %v: %w", filepath, err)
		}

		var rawContents []byte
		if contents != nil {
			rawContents, err = contents.Get()
			if err != nil {
				return nil, fmt.Errorf("failed to get file contents for %v: %w", filepath, err)
			}

			fileHeader.Size = int64(len(rawContents))
		}

		if err := tarBuilder.WriteHeader(fileHeader); err != nil {
			return nil, fmt.Errorf("failed to write tar header entry for %v: %w", filepath, err)
		}

		if contents != nil {
			if _, err := tarBuilder.Write(rawContents); err != nil {
				return nil, fmt.Errorf("failed to write tar file entry for %v: %w", filepath, err)
			}
		}
	}

	return bytes.NewReader(buffer.Bytes()), nil
}

type BuildOpt interface {
	Apply(cfg *types.ImageBuildOptions) error
}

type BuildRemove bool

var _ BuildOpt = (*BuildRemove)(nil)

func (u BuildRemove) Apply(cfg *types.ImageBuildOptions) error {
	cfg.Remove = bool(u)
	return nil
}

type BuildForceRemove bool

var _ BuildOpt = (*BuildForceRemove)(nil)

func (u BuildForceRemove) Apply(cfg *types.ImageBuildOptions) error {
	cfg.ForceRemove = bool(u)
	return nil
}

type BuildPullParent bool

var _ BuildOpt = (*BuildPullParent)(nil)

func (u BuildPullParent) Apply(cfg *types.ImageBuildOptions) error {
	cfg.PullParent = bool(u)
	return nil
}

type BuildArgs map[string]*string

var _ BuildOpt = (*BuildArgs)(nil)

func (u BuildArgs) Apply(cfg *types.ImageBuildOptions) error {
	cfg.BuildArgs = u
	return nil
}

type BuildTags []string

var _ BuildOpt = (*BuildTags)(nil)

func (u BuildTags) Apply(cfg *types.ImageBuildOptions) error {
	cfg.Tags = u
	return nil
}

const containerfilePath = "_containerfile"

func (d *Runner) BuildImage(ctx context.Context, containerfile string, containerContext BuildContext, opts ...BuildOpt) ([]byte, error) {
	return BuildImage(ctx, d.DockerAPI, containerfile, containerContext, opts...)
}

func BuildImage(ctx context.Context, api *client.Client, containerfile string, containerContext BuildContext, opts ...BuildOpt) ([]byte, error) {
	var cfg types.ImageBuildOptions

	// Build container context tarball, provisioning containerfile in.
	containerContext[containerfilePath] = PathContentsFromBytes([]byte(containerfile))
	tar, err := containerContext.ToTarball()
	if err != nil {
		return nil, fmt.Errorf("failed to create build image context tarball: %w", err)
	}
	cfg.Dockerfile = "/" + containerfilePath

	// Apply all given options
	for index, opt := range opts {
		if err := opt.Apply(&cfg); err != nil {
			return nil, fmt.Errorf("failed to apply option (%d / %v): %w", index, opt, err)
		}
	}

	resp, err := api.ImageBuild(ctx, tar, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to build image: %v", err)
	}

	output, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read image build output: %w", err)
	}

	return output, nil
}

func (d *Runner) CopyTo(container string, destination string, contents BuildContext) error {
	// XXX: currently we use the default options but we might want to allow
	// modifying cfg.CopyUIDGID in the future.
	var cfg types.CopyToContainerOptions

	// Convert our provided contents to a tarball to ship up.
	tar, err := contents.ToTarball()
	if err != nil {
		return fmt.Errorf("failed to build contents into tarball: %v", err)
	}

	return d.DockerAPI.CopyToContainer(context.Background(), container, destination, tar, cfg)
}

func (d *Runner) CopyFrom(container string, source string) (BuildContext, *types.ContainerPathStat, error) {
	reader, stat, err := d.DockerAPI.CopyFromContainer(context.Background(), container, source)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read %v from container: %v", source, err)
	}

	result, err := BuildContextFromTarball(reader)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to build archive from result: %v", err)
	}

	return result, &stat, nil
}

func (d *Runner) GetNetworkAndAddresses(container string) (map[string]string, error) {
	response, err := d.DockerAPI.ContainerInspect(context.Background(), container)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch container inspection data: %v", err)
	}

	if response.NetworkSettings == nil || len(response.NetworkSettings.Networks) == 0 {
		return nil, fmt.Errorf("container (%v) had no associated network settings: %v", container, response)
	}

	ret := make(map[string]string)
	ns := response.NetworkSettings.Networks
	for network, data := range ns {
		if data == nil {
			continue
		}

		ret[network] = data.IPAddress
	}

	if len(ret) == 0 {
		return nil, fmt.Errorf("no valid network data for container (%v): %v", container, response)
	}

	return ret, nil
}
