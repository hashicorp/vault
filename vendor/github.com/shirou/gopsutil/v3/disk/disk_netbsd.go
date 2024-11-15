//go:build netbsd
// +build netbsd

package disk

import (
	"context"
	"unsafe"

	"github.com/shirou/gopsutil/v3/internal/common"
	"golang.org/x/sys/unix"
)

const (
	// see sys/fstypes.h and `man 5 statvfs`
	MNT_RDONLY      = 0x00000001 /* read only filesystem */
	MNT_SYNCHRONOUS = 0x00000002 /* file system written synchronously */
	MNT_NOEXEC      = 0x00000004 /* can't exec from filesystem */
	MNT_NOSUID      = 0x00000008 /* don't honor setuid bits on fs */
	MNT_NODEV       = 0x00000010 /* don't interpret special files */
	MNT_ASYNC       = 0x00000040 /* file system written asynchronously */
	MNT_NOATIME     = 0x04000000 /* Never update access times in fs */
	MNT_SOFTDEP     = 0x80000000 /* Use soft dependencies */
)

func PartitionsWithContext(ctx context.Context, all bool) ([]PartitionStat, error) {
	var ret []PartitionStat

	flag := uint64(1) // ST_WAIT/MNT_WAIT, see sys/fstypes.h

	// get required buffer size
	emptyBufSize := 0
	r, _, err := unix.Syscall(
		483, // SYS___getvfsstat90 syscall
		uintptr(unsafe.Pointer(nil)),
		uintptr(unsafe.Pointer(&emptyBufSize)),
		uintptr(unsafe.Pointer(&flag)),
	)
	if err != 0 {
		return ret, err
	}
	mountedFsCount := uint64(r)

	// calculate the buffer size
	bufSize := sizeOfStatvfs * mountedFsCount
	buf := make([]Statvfs, mountedFsCount)

	// request again to get desired mount data
	_, _, err = unix.Syscall(
		483, // SYS___getvfsstat90 syscall
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(&bufSize)),
		uintptr(unsafe.Pointer(&flag)),
	)
	if err != 0 {
		return ret, err
	}

	for _, stat := range buf {
		opts := []string{"rw"}
		if stat.Flag&MNT_RDONLY != 0 {
			opts = []string{"rw"}
		}
		if stat.Flag&MNT_SYNCHRONOUS != 0 {
			opts = append(opts, "sync")
		}
		if stat.Flag&MNT_NOEXEC != 0 {
			opts = append(opts, "noexec")
		}
		if stat.Flag&MNT_NOSUID != 0 {
			opts = append(opts, "nosuid")
		}
		if stat.Flag&MNT_NODEV != 0 {
			opts = append(opts, "nodev")
		}
		if stat.Flag&MNT_ASYNC != 0 {
			opts = append(opts, "async")
		}
		if stat.Flag&MNT_SOFTDEP != 0 {
			opts = append(opts, "softdep")
		}
		if stat.Flag&MNT_NOATIME != 0 {
			opts = append(opts, "noatime")
		}

		d := PartitionStat{
			Device:     common.ByteToString([]byte(stat.Mntfromname[:])),
			Mountpoint: common.ByteToString([]byte(stat.Mntonname[:])),
			Fstype:     common.ByteToString([]byte(stat.Fstypename[:])),
			Opts:       opts,
		}

		ret = append(ret, d)
	}

	return ret, nil
}

func IOCountersWithContext(ctx context.Context, names ...string) (map[string]IOCountersStat, error) {
	ret := make(map[string]IOCountersStat)
	return ret, common.ErrNotImplementedError
}

func UsageWithContext(ctx context.Context, path string) (*UsageStat, error) {
	stat := Statvfs{}
	flag := uint64(1) // ST_WAIT/MNT_WAIT, see sys/fstypes.h

	_path, e := unix.BytePtrFromString(path)
	if e != nil {
		return nil, e
	}

	_, _, err := unix.Syscall(
		484, // SYS___statvfs190, see sys/syscall.h
		uintptr(unsafe.Pointer(_path)),
		uintptr(unsafe.Pointer(&stat)),
		uintptr(unsafe.Pointer(&flag)),
	)
	if err != 0 {
		return nil, err
	}

	// frsize is the real block size on NetBSD. See discuss here: https://bugzilla.samba.org/show_bug.cgi?id=11810
	bsize := stat.Frsize
	ret := &UsageStat{
		Path:        path,
		Fstype:      getFsType(stat),
		Total:       (uint64(stat.Blocks) * uint64(bsize)),
		Free:        (uint64(stat.Bavail) * uint64(bsize)),
		InodesTotal: (uint64(stat.Files)),
		InodesFree:  (uint64(stat.Ffree)),
	}

	ret.InodesUsed = (ret.InodesTotal - ret.InodesFree)
	ret.InodesUsedPercent = (float64(ret.InodesUsed) / float64(ret.InodesTotal)) * 100.0
	ret.Used = (uint64(stat.Blocks) - uint64(stat.Bfree)) * uint64(bsize)
	ret.UsedPercent = (float64(ret.Used) / float64(ret.Total)) * 100.0

	return ret, nil
}

func getFsType(stat Statvfs) string {
	return common.ByteToString(stat.Fstypename[:])
}

func SerialNumberWithContext(ctx context.Context, name string) (string, error) {
	return "", common.ErrNotImplementedError
}

func LabelWithContext(ctx context.Context, name string) (string, error) {
	return "", common.ErrNotImplementedError
}
