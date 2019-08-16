package vault

import (
	"runtime"
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

type HostInfoError struct {
	Err error
}

func (e *HostInfoError) WrappedErrors() []error {
	return []error{e.Err}
}

func (e *HostInfoError) Error() string {
	return e.Err.Error()
}

func (c *Core) CollectHostInfo() (*HostInfo, error) {
	var retErr *multierror.Error
	info := &HostInfo{CollectionTime: time.Now()}

	if h, err := host.Info(); err != nil {
		retErr = multierror.Append(retErr, &HostInfoError{err})
	} else {
		info.Host = h
	}

	if v, err := mem.VirtualMemory(); err != nil {
		retErr = multierror.Append(retErr, &HostInfoError{err})
	} else {
		info.Memory = v
	}

	diskPath := "/"
	if runtime.GOOS == "windows" {
		diskPath = "C:"
	}

	if d, err := disk.Usage(diskPath); err != nil {
		retErr = multierror.Append(retErr, &HostInfoError{err})
	} else {
		info.Disk = d
	}

	if c, err := cpu.Info(); err != nil {
		retErr = multierror.Append(retErr, &HostInfoError{err})
	} else {
		info.CPU = c
	}

	return info, retErr.ErrorOrNil()
}
