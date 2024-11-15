//go:build darwin
// +build darwin

package disk

import (
	"context"

	"golang.org/x/sys/unix"

	"github.com/shirou/gopsutil/v3/internal/common"
)

// PartitionsWithContext returns disk partition.
// 'all' argument is ignored, see: https://github.com/giampaolo/psutil/issues/906
func PartitionsWithContext(ctx context.Context, all bool) ([]PartitionStat, error) {
	var ret []PartitionStat

	count, err := unix.Getfsstat(nil, unix.MNT_WAIT)
	if err != nil {
		return ret, err
	}
	fs := make([]unix.Statfs_t, count)
	count, err = unix.Getfsstat(fs, unix.MNT_WAIT)
	if err != nil {
		return ret, err
	}
	// On 10.14, and possibly other OS versions, the actual count may
	// be less than from the first call. Truncate to the returned count
	// to prevent accessing uninitialized entries.
	// https://github.com/shirou/gopsutil/issues/1390
	fs = fs[:count]
	for _, stat := range fs {
		opts := []string{"rw"}
		if stat.Flags&unix.MNT_RDONLY != 0 {
			opts = []string{"ro"}
		}
		if stat.Flags&unix.MNT_SYNCHRONOUS != 0 {
			opts = append(opts, "sync")
		}
		if stat.Flags&unix.MNT_NOEXEC != 0 {
			opts = append(opts, "noexec")
		}
		if stat.Flags&unix.MNT_NOSUID != 0 {
			opts = append(opts, "nosuid")
		}
		if stat.Flags&unix.MNT_UNION != 0 {
			opts = append(opts, "union")
		}
		if stat.Flags&unix.MNT_ASYNC != 0 {
			opts = append(opts, "async")
		}
		if stat.Flags&unix.MNT_DONTBROWSE != 0 {
			opts = append(opts, "nobrowse")
		}
		if stat.Flags&unix.MNT_AUTOMOUNTED != 0 {
			opts = append(opts, "automounted")
		}
		if stat.Flags&unix.MNT_JOURNALED != 0 {
			opts = append(opts, "journaled")
		}
		if stat.Flags&unix.MNT_MULTILABEL != 0 {
			opts = append(opts, "multilabel")
		}
		if stat.Flags&unix.MNT_NOATIME != 0 {
			opts = append(opts, "noatime")
		}
		if stat.Flags&unix.MNT_NODEV != 0 {
			opts = append(opts, "nodev")
		}
		d := PartitionStat{
			Device:     common.ByteToString(stat.Mntfromname[:]),
			Mountpoint: common.ByteToString(stat.Mntonname[:]),
			Fstype:     common.ByteToString(stat.Fstypename[:]),
			Opts:       opts,
		}

		ret = append(ret, d)
	}

	return ret, nil
}

func getFsType(stat unix.Statfs_t) string {
	return common.ByteToString(stat.Fstypename[:])
}

func SerialNumberWithContext(ctx context.Context, name string) (string, error) {
	return "", common.ErrNotImplementedError
}

func LabelWithContext(ctx context.Context, name string) (string, error) {
	return "", common.ErrNotImplementedError
}
