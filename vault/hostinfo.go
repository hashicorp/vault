package vault

import (
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

type HostInfo struct {
	CollectionTime time.Time              `json:"collection_time"`
	CPU            []cpu.InfoStat         `json:"cpu"`
	Disk           *disk.UsageStat        `json:"disk"`
	Host           *host.InfoStat         `json:"host"`
	Memory         *mem.VirtualMemoryStat `json:"memory"`
}

func (c *Core) CollectHostInfo() (*HostInfo, error) {
	var retErr error
	info := &HostInfo{CollectionTime: time.Now()}

	if h, err := host.Info(); err != nil {
		retErr = multierror.Append(retErr, err)
	} else {
		info.Host = h
	}

	if v, err := mem.VirtualMemory(); err != nil {
		retErr = multierror.Append(retErr, err)
	} else {
		info.Memory = v
	}

	if d, err := disk.Usage("/"); err != nil {
		retErr = multierror.Append(retErr, err)
	} else {
		info.Disk = d
	}

	if c, err := cpu.Info(); err != nil {
		retErr = multierror.Append(retErr, err)
	} else {
		info.CPU = c
	}

	return info, retErr
}
