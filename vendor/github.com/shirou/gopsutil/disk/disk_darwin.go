// +build darwin

package disk

import (
	"context"
	"path"

	"github.com/shirou/gopsutil/internal/common"
	"golang.org/x/sys/unix"
)

func Partitions(all bool) ([]PartitionStat, error) {
	return PartitionsWithContext(context.Background(), all)
}

func PartitionsWithContext(ctx context.Context, all bool) ([]PartitionStat, error) {
	var ret []PartitionStat

	count, err := unix.Getfsstat(nil, unix.MNT_WAIT)
	if err != nil {
		return ret, err
	}
	fs := make([]unix.Statfs_t, count)
	if _, err = unix.Getfsstat(fs, unix.MNT_WAIT); err != nil {
		return ret, err
	}
	for _, stat := range fs {
		opts := "rw"
		if stat.Flags&unix.MNT_RDONLY != 0 {
			opts = "ro"
		}
		if stat.Flags&unix.MNT_SYNCHRONOUS != 0 {
			opts += ",sync"
		}
		if stat.Flags&unix.MNT_NOEXEC != 0 {
			opts += ",noexec"
		}
		if stat.Flags&unix.MNT_NOSUID != 0 {
			opts += ",nosuid"
		}
		if stat.Flags&unix.MNT_UNION != 0 {
			opts += ",union"
		}
		if stat.Flags&unix.MNT_ASYNC != 0 {
			opts += ",async"
		}
		if stat.Flags&unix.MNT_DONTBROWSE != 0 {
			opts += ",nobrowse"
		}
		if stat.Flags&unix.MNT_AUTOMOUNTED != 0 {
			opts += ",automounted"
		}
		if stat.Flags&unix.MNT_JOURNALED != 0 {
			opts += ",journaled"
		}
		if stat.Flags&unix.MNT_MULTILABEL != 0 {
			opts += ",multilabel"
		}
		if stat.Flags&unix.MNT_NOATIME != 0 {
			opts += ",noatime"
		}
		if stat.Flags&unix.MNT_NODEV != 0 {
			opts += ",nodev"
		}
		d := PartitionStat{
			Device:     common.IntToString(stat.Mntfromname[:]),
			Mountpoint: common.IntToString(stat.Mntonname[:]),
			Fstype:     common.IntToString(stat.Fstypename[:]),
			Opts:       opts,
		}
		if all == false {
			if !path.IsAbs(d.Device) || !common.PathExists(d.Device) {
				continue
			}
		}

		ret = append(ret, d)
	}

	return ret, nil
}

func getFsType(stat unix.Statfs_t) string {
	return common.IntToString(stat.Fstypename[:])
}
