package docker

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/strslice"
	docker "github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
)

// Runner manages the lifecycle of the Docker container
type Runner struct {
	dockerAPI       *docker.Client
	ContainerConfig *container.Config
	ContainerName   string
	NetName         string
	IP              string
	CopyFromTo      map[string]string
}

// Start is responsible for executing the Vault container. It consists of
// pulling the specified Vault image, creating the container, and copies the
// plugin binary into the container file system before starting the container
// itself.
func (d *Runner) Start(ctx context.Context) (*types.ContainerJSON, error) {
	hostConfig := &container.HostConfig{
		PublishAllPorts: true,
		AutoRemove:      true,
	}

	networkingConfig := &network.NetworkingConfig{}
	switch d.NetName {
	case "":
	case "host":
		hostConfig.NetworkMode = "host"
	default:
		es := &network.EndpointSettings{
			Aliases: []string{d.ContainerName},
		}
		if len(d.IP) != 0 {
			es.IPAMConfig = &network.EndpointIPAMConfig{
				IPv4Address: d.IP,
			}
		}
		networkingConfig.EndpointsConfig = map[string]*network.EndpointSettings{
			d.NetName: es,
		}
	}

	// Best-effort pull. ImageCreate here will use a matching image from the local
	// Docker library, or if not found pull the matching image from docker hub. If
	// not found on docker hub, returns an error. The response must be read in
	// order for the local image.
	resp, err := d.dockerAPI.ImageCreate(ctx, d.ContainerConfig.Image, types.ImageCreateOptions{})
	if err != nil {
		return nil, err
	}
	if resp != nil {
		_, _ = ioutil.ReadAll(resp)
	}

	cfg := *d.ContainerConfig
	hostConfig.CapAdd = strslice.StrSlice{"IPC_LOCK"}
	cfg.Hostname = d.ContainerName
	fullName := d.ContainerName
	container, err := d.dockerAPI.ContainerCreate(ctx, &cfg, hostConfig, networkingConfig, fullName)
	if err != nil {
		return nil, fmt.Errorf("container create failed: %v", err)
	}

	// copies the plugin binary into the Docker file system. This copy is only
	// allowed before the container is started
	for from, to := range d.CopyFromTo {
		srcInfo, err := archive.CopyInfoSourcePath(from, false)
		if err != nil {
			return nil, fmt.Errorf("error copying from source %q: %v", from, err)
		}

		srcArchive, err := archive.TarResource(srcInfo)
		if err != nil {
			return nil, fmt.Errorf("error creating tar from source %q: %v", from, err)
		}
		defer srcArchive.Close()

		dstInfo := archive.CopyInfo{Path: to}

		dstDir, content, err := archive.PrepareArchiveCopy(srcArchive, srcInfo, dstInfo)
		if err != nil {
			return nil, fmt.Errorf("error preparing copy from %q -> %q: %v", from, to, err)
		}
		defer content.Close()
		err = d.dockerAPI.CopyToContainer(ctx, container.ID, dstDir, content, types.CopyToContainerOptions{})
		if err != nil {
			return nil, fmt.Errorf("error copying from %q -> %q: %v", from, to, err)
		}
	}

	err = d.dockerAPI.ContainerStart(ctx, container.ID, types.ContainerStartOptions{})
	if err != nil {
		return nil, fmt.Errorf("container start failed: %v", err)
	}

	inspect, err := d.dockerAPI.ContainerInspect(ctx, container.ID)
	if err != nil {
		return nil, err
	}
	return &inspect, nil
}
