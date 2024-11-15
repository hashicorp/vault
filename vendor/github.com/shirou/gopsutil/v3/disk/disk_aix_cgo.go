//go:build aix && cgo
// +build aix,cgo

package disk

import (
	"context"
	"fmt"

	"github.com/power-devops/perfstat"
)

var FSType map[int]string

func init() {
	FSType = map[int]string{
		0: "jfs2", 1: "namefs", 2: "nfs", 3: "jfs", 5: "cdrom", 6: "proc",
		16: "special-fs", 17: "cache-fs", 18: "nfs3", 19: "automount-fs", 20: "pool-fs", 32: "vxfs",
		33: "veritas-fs", 34: "udfs", 35: "nfs4", 36: "nfs4-pseudo", 37: "smbfs", 38: "mcr-pseudofs",
		39: "ahafs", 40: "sterm-nfs", 41: "asmfs",
	}
}

func PartitionsWithContext(ctx context.Context, all bool) ([]PartitionStat, error) {
	f, err := perfstat.FileSystemStat()
	if err != nil {
		return nil, err
	}
	ret := make([]PartitionStat, len(f))

	for _, fs := range f {
		fstyp, exists := FSType[fs.FSType]
		if !exists {
			fstyp = "unknown"
		}
		info := PartitionStat{
			Device:     fs.Device,
			Mountpoint: fs.MountPoint,
			Fstype:     fstyp,
		}
		ret = append(ret, info)
	}

	return ret, err
}

func UsageWithContext(ctx context.Context, path string) (*UsageStat, error) {
	f, err := perfstat.FileSystemStat()
	if err != nil {
		return nil, err
	}

	blocksize := uint64(512)
	for _, fs := range f {
		if path == fs.MountPoint {
			fstyp, exists := FSType[fs.FSType]
			if !exists {
				fstyp = "unknown"
			}
			info := UsageStat{
				Path:        path,
				Fstype:      fstyp,
				Total:       uint64(fs.TotalBlocks) * blocksize,
				Free:        uint64(fs.FreeBlocks) * blocksize,
				Used:        uint64(fs.TotalBlocks-fs.FreeBlocks) * blocksize,
				InodesTotal: uint64(fs.TotalInodes),
				InodesFree:  uint64(fs.FreeInodes),
				InodesUsed:  uint64(fs.TotalInodes - fs.FreeInodes),
			}
			info.UsedPercent = (float64(info.Used) / float64(info.Total)) * 100.0
			info.InodesUsedPercent = (float64(info.InodesUsed) / float64(info.InodesTotal)) * 100.0
			return &info, nil
		}
	}
	return nil, fmt.Errorf("mountpoint %s not found", path)
}
