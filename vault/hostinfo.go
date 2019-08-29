package vault

import (
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

// HostInfo holds all the information that gets captured on the host.
type HostInfo struct {
	Timestamp time.Time              `json:"timestamp"`
	CPU       []cpu.InfoStat         `json:"cpu"`
	Disk      []*disk.UsageStat      `json:"disk"`
	Host      *host.InfoStat         `json:"host"`
	Memory    *mem.VirtualMemoryStat `json:"memory"`
}

// HostInfoError is a typed error for more convenient error checking.
type HostInfoError struct {
	Err error
}

func (e *HostInfoError) WrappedErrors() []error {
	return []error{e.Err}
}

func (e *HostInfoError) Error() string {
	return e.Err.Error()
}

// CollectHostInfo returns information on the host, which includes general
// host status, CPU, memory, and disk utilization.
//
// The function does a best-effort capture on the most information possible,
// continuing on capture errors encountered and appending them to a resulting
// multierror.Error that gets returned at the end.
func (c *Core) CollectHostInfo() (*HostInfo, error) {
	var retErr *multierror.Error
	info := &HostInfo{Timestamp: time.Now().UTC()}

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

	parts, err := disk.Partitions(false)
	if err != nil {
		retErr = multierror.Append(retErr, &HostInfoError{err})
	} else {
		var usage []*disk.UsageStat
		for _, part := range parts {
			u, err := disk.Usage(part.Mountpoint)
			if err != nil {
				retErr = multierror.Append(retErr, &HostInfoError{err})
			}
			usage = append(usage, u)

		}
		info.Disk = usage
	}

	if c, err := cpu.Info(); err != nil {
		retErr = multierror.Append(retErr, &HostInfoError{err})
	} else {
		info.CPU = c
	}

	return info, retErr.ErrorOrNil()
}
