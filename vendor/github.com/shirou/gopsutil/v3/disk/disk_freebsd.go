//go:build freebsd
// +build freebsd

package disk

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/v3/internal/common"
	"golang.org/x/sys/unix"
)

// PartitionsWithContext returns disk partition.
// 'all' argument is ignored, see: https://github.com/giampaolo/psutil/issues/906
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
		if stat.Flags&unix.MNT_SUIDDIR != 0 {
			opts = append(opts, "suiddir")
		}
		if stat.Flags&unix.MNT_SOFTDEP != 0 {
			opts = append(opts, "softdep")
		}
		if stat.Flags&unix.MNT_NOSYMFOLLOW != 0 {
			opts = append(opts, "nosymfollow")
		}
		if stat.Flags&unix.MNT_GJOURNAL != 0 {
			opts = append(opts, "gjournal")
		}
		if stat.Flags&unix.MNT_MULTILABEL != 0 {
			opts = append(opts, "multilabel")
		}
		if stat.Flags&unix.MNT_ACLS != 0 {
			opts = append(opts, "acls")
		}
		if stat.Flags&unix.MNT_NOATIME != 0 {
			opts = append(opts, "noatime")
		}
		if stat.Flags&unix.MNT_NOCLUSTERR != 0 {
			opts = append(opts, "noclusterr")
		}
		if stat.Flags&unix.MNT_NOCLUSTERW != 0 {
			opts = append(opts, "noclusterw")
		}
		if stat.Flags&unix.MNT_NFS4ACLS != 0 {
			opts = append(opts, "nfsv4acls")
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

func IOCountersWithContext(ctx context.Context, names ...string) (map[string]IOCountersStat, error) {
	// statinfo->devinfo->devstat
	// /usr/include/devinfo.h
	ret := make(map[string]IOCountersStat)

	r, err := unix.Sysctl("kern.devstat.all")
	if err != nil {
		return nil, err
	}
	buf := []byte(r)
	length := len(buf)

	count := int(uint64(length) / uint64(sizeOfdevstat))

	buf = buf[8:] // devstat.all has version in the head.
	// parse buf to devstat
	for i := 0; i < count; i++ {
		b := buf[i*sizeOfdevstat : i*sizeOfdevstat+sizeOfdevstat]
		d, err := parsedevstat(b)
		if err != nil {
			continue
		}
		un := strconv.Itoa(int(d.Unit_number))
		name := common.IntToString(d.Device_name[:]) + un

		if len(names) > 0 && !common.StringsHas(names, name) {
			continue
		}

		ds := IOCountersStat{
			ReadCount:  d.Operations[devstat_READ],
			WriteCount: d.Operations[devstat_WRITE],
			ReadBytes:  d.Bytes[devstat_READ],
			WriteBytes: d.Bytes[devstat_WRITE],
			ReadTime:   uint64(d.Duration[devstat_READ].Compute() * 1000),
			WriteTime:  uint64(d.Duration[devstat_WRITE].Compute() * 1000),
			IoTime:     uint64(d.Busy_time.Compute() * 1000),
			Name:       name,
		}
		ds.SerialNumber, _ = SerialNumberWithContext(ctx, name)
		ret[name] = ds
	}

	return ret, nil
}

func (b bintime) Compute() float64 {
	BINTIME_SCALE := 5.42101086242752217003726400434970855712890625e-20
	return float64(b.Sec) + float64(b.Frac)*BINTIME_SCALE
}

// BT2LD(time)     ((long double)(time).sec + (time).frac * BINTIME_SCALE)

func parsedevstat(buf []byte) (devstat, error) {
	var ds devstat
	br := bytes.NewReader(buf)
	//	err := binary.Read(br, binary.LittleEndian, &ds)
	err := common.Read(br, binary.LittleEndian, &ds)
	if err != nil {
		return ds, err
	}

	return ds, nil
}

func getFsType(stat unix.Statfs_t) string {
	return common.ByteToString(stat.Fstypename[:])
}

func SerialNumberWithContext(ctx context.Context, name string) (string, error) {
	geomOut, err := invoke.CommandWithContext(ctx, "geom", "disk", "list", name)
	if err != nil {
		return "", fmt.Errorf("exec geom: %w", err)
	}
	s := bufio.NewScanner(bytes.NewReader(geomOut))
	serial := ""
	for s.Scan() {
		flds := strings.Fields(s.Text())
		if len(flds) == 2 && flds[0] == "ident:" {
			if flds[1] != "(null)" {
				serial = flds[1]
			}
			break
		}
	}
	if err = s.Err(); err != nil {
		return "", err
	}
	return serial, nil
}

func LabelWithContext(ctx context.Context, name string) (string, error) {
	return "", common.ErrNotImplementedError
}
