//go:build !openbsd

package hostutil

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

// HostInfo holds all the information that gets captured on the host. The
// set of information captured depends on the host operating system. For more
// information, refer to: https://github.com/shirou/gopsutil#current-status
type HostInfo struct {
	// Timestamp returns the timestamp in UTC on the collection time.
	Timestamp time.Time `json:"timestamp"`
	// CPU returns information about the CPU such as family, model, cores, etc.
	CPU []cpu.InfoStat `json:"cpu"`
	// CPUTimes returns statistics on CPU usage represented in Jiffies.
	CPUTimes []cpu.TimesStat `json:"cpu_times"`
	// Disk returns statitics on disk usage for all accessible partitions.
	Disk []*disk.UsageStat `json:"disk"`
	// Host returns general host information such as hostname, platform, uptime,
	// kernel version, etc.
	Host *HostInfoStat `json:"host"`
	// Memory contains statistics about the memory such as total, available, and
	// used memory in number of bytes.
	Memory *VirtualMemoryStat `json:"memory"`
}

// CollectHostInfo returns information on the host, which includes general
// host status, CPU, memory, and disk utilization.
//
// The function does a best-effort capture on the most information possible,
// continuing on capture errors encountered and appending them to a resulting
// multierror.Error that gets returned at the end.
func CollectHostInfo(ctx context.Context) (*HostInfo, error) {
	var retErr *multierror.Error
	info := &HostInfo{Timestamp: time.Now().UTC()}

	if h, err := CollectHostInfoStat(ctx); err != nil {
		retErr = multierror.Append(retErr, &HostInfoError{"host", err})
	} else {
		info.Host = h
	}

	if v, err := CollectHostMemory(ctx); err != nil {
		retErr = multierror.Append(retErr, &HostInfoError{"memory", err})
	} else {
		info.Memory = v
	}

	parts, err := disk.PartitionsWithContext(ctx, false)
	if err != nil {
		retErr = multierror.Append(retErr, &HostInfoError{"disk", err})
	} else {
		var usage []*disk.UsageStat
		for i, part := range parts {
			u, err := disk.UsageWithContext(ctx, part.Mountpoint)
			if err != nil {
				retErr = multierror.Append(retErr, &HostInfoError{fmt.Sprintf("disk.%d", i), err})
				continue
			}
			usage = append(usage, u)

		}
		info.Disk = usage
	}

	if c, err := cpu.InfoWithContext(ctx); err != nil {
		retErr = multierror.Append(retErr, &HostInfoError{"cpu", err})
	} else {
		info.CPU = c
	}

	t, err := cpu.TimesWithContext(ctx, true)
	if err != nil {
		retErr = multierror.Append(retErr, &HostInfoError{"cpu_times", err})
	} else {
		info.CPUTimes = t
	}

	return info, retErr.ErrorOrNil()
}

func CollectHostMemory(ctx context.Context) (*VirtualMemoryStat, error) {
	m, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		return nil, err
	}

	return &VirtualMemoryStat{
		VirtualMemoryStat: m,
		// Fields below are added to maintain backwards compatibility with gopsutil v2
		CommitLimit:    m.CommitLimit,
		CommittedAS:    m.CommittedAS,
		HighFree:       m.HighFree,
		HighTotal:      m.HighTotal,
		HugePagesFree:  m.HugePagesFree,
		HugePageSize:   m.HugePageSize,
		HugePagesTotal: m.HugePagesTotal,
		LowFree:        m.LowFree,
		LowTotal:       m.LowTotal,
		PageTables:     m.PageTables,
		SwapCached:     m.SwapCached,
		SwapFree:       m.SwapFree,
		SwapTotal:      m.SwapTotal,
		VMallocChunk:   m.VmallocChunk,
		VMallocTotal:   m.VmallocTotal,
		VMallocUsed:    m.VmallocUsed,
		Writeback:      m.WriteBack,
		WritebackTmp:   m.WriteBackTmp,
	}, nil
}

func CollectHostInfoStat(ctx context.Context) (*HostInfoStat, error) {
	h, err := host.InfoWithContext(ctx)
	if err != nil {
		return nil, err
	}

	return &HostInfoStat{
		InfoStat: h,
		// Fields below are added to maintain backwards compatibility with gopsutil v2
		HostID: h.HostID,
	}, nil
}
