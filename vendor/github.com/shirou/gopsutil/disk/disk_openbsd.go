// +build openbsd

package disk

import (
	"bytes"
	"context"
	"encoding/binary"

	"github.com/shirou/gopsutil/internal/common"
	"golang.org/x/sys/unix"
)

func PartitionsWithContext(ctx context.Context, all bool) ([]PartitionStat, error) {
	var ret []PartitionStat

	// get length
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
		if stat.F_flags&unix.MNT_RDONLY != 0 {
			opts = "ro"
		}
		if stat.F_flags&unix.MNT_SYNCHRONOUS != 0 {
			opts += ",sync"
		}
		if stat.F_flags&unix.MNT_NOEXEC != 0 {
			opts += ",noexec"
		}
		if stat.F_flags&unix.MNT_NOSUID != 0 {
			opts += ",nosuid"
		}
		if stat.F_flags&unix.MNT_NODEV != 0 {
			opts += ",nodev"
		}
		if stat.F_flags&unix.MNT_ASYNC != 0 {
			opts += ",async"
		}
		if stat.F_flags&unix.MNT_SOFTDEP != 0 {
			opts += ",softdep"
		}
		if stat.F_flags&unix.MNT_NOATIME != 0 {
			opts += ",noatime"
		}
		if stat.F_flags&unix.MNT_WXALLOWED != 0 {
			opts += ",wxallowed"
		}

		d := PartitionStat{
			Device:     common.IntToString(stat.F_mntfromname[:]),
			Mountpoint: common.IntToString(stat.F_mntonname[:]),
			Fstype:     common.IntToString(stat.F_fstypename[:]),
			Opts:       opts,
		}

		ret = append(ret, d)
	}

	return ret, nil
}

func IOCountersWithContext(ctx context.Context, names ...string) (map[string]IOCountersStat, error) {
	ret := make(map[string]IOCountersStat)

	r, err := unix.SysctlRaw("hw.diskstats")
	if err != nil {
		return nil, err
	}
	buf := []byte(r)
	length := len(buf)

	count := int(uint64(length) / uint64(sizeOfDiskstats))

	// parse buf to Diskstats
	for i := 0; i < count; i++ {
		b := buf[i*sizeOfDiskstats : i*sizeOfDiskstats+sizeOfDiskstats]
		d, err := parseDiskstats(b)
		if err != nil {
			continue
		}
		name := common.IntToString(d.Name[:])

		if len(names) > 0 && !common.StringsHas(names, name) {
			continue
		}

		ds := IOCountersStat{
			ReadCount:  d.Rxfer,
			WriteCount: d.Wxfer,
			ReadBytes:  d.Rbytes,
			WriteBytes: d.Wbytes,
			Name:       name,
		}
		ret[name] = ds
	}

	return ret, nil
}

// BT2LD(time)     ((long double)(time).sec + (time).frac * BINTIME_SCALE)

func parseDiskstats(buf []byte) (Diskstats, error) {
	var ds Diskstats
	br := bytes.NewReader(buf)
	//	err := binary.Read(br, binary.LittleEndian, &ds)
	err := common.Read(br, binary.LittleEndian, &ds)
	if err != nil {
		return ds, err
	}

	return ds, nil
}

func UsageWithContext(ctx context.Context, path string) (*UsageStat, error) {
	stat := unix.Statfs_t{}
	err := unix.Statfs(path, &stat)
	if err != nil {
		return nil, err
	}
	bsize := stat.F_bsize

	ret := &UsageStat{
		Path:        path,
		Fstype:      getFsType(stat),
		Total:       (uint64(stat.F_blocks) * uint64(bsize)),
		Free:        (uint64(stat.F_bavail) * uint64(bsize)),
		InodesTotal: (uint64(stat.F_files)),
		InodesFree:  (uint64(stat.F_ffree)),
	}

	ret.InodesUsed = (ret.InodesTotal - ret.InodesFree)
	ret.InodesUsedPercent = (float64(ret.InodesUsed) / float64(ret.InodesTotal)) * 100.0
	ret.Used = (uint64(stat.F_blocks) - uint64(stat.F_bfree)) * uint64(bsize)
	ret.UsedPercent = (float64(ret.Used) / float64(ret.Total)) * 100.0

	return ret, nil
}

func getFsType(stat unix.Statfs_t) string {
	return common.IntToString(stat.F_fstypename[:])
}
