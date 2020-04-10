package acctest

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

func (d *Runner) Start(ctx context.Context) (*types.ContainerJSON, error) {
	hostConfig := &container.HostConfig{
		PublishAllPorts: true,
		// AutoRemove:      false,
		// TODO: configure auto remove
		AutoRemove: true,
	}

	networkingConfig := &network.NetworkingConfig{}
	switch d.NetName {
	case "":
	case "host":
		hostConfig.NetworkMode = "host"
	default:
		es := &network.EndpointSettings{
			//Links:               nil,
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

	// best-effort pull
	resp, _ := d.dockerAPI.ImageCreate(ctx, d.ContainerConfig.Image, types.ImageCreateOptions{})
	if resp != nil {
		_, _ = ioutil.ReadAll(resp)
	}

	cfg := *d.ContainerConfig
	hostConfig.CapAdd = strslice.StrSlice{"IPC_LOCK"}
	cfg.Hostname = d.ContainerName
	//fullName := d.NetName + "." + d.ContainerName
	fullName := d.ContainerName
	container, err := d.dockerAPI.ContainerCreate(ctx, &cfg, hostConfig, networkingConfig, fullName)
	if err != nil {
		return nil, fmt.Errorf("container create failed: %v", err)
	}

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
