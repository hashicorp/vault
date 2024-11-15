//go:build aix && !cgo
// +build aix,!cgo

package disk

import (
	"context"
	"regexp"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/v3/internal/common"
	"golang.org/x/sys/unix"
)

var startBlank = regexp.MustCompile(`^\s+`)

var ignoreFSType = map[string]bool{"procfs": true}
var FSType = map[int]string{
	0: "jfs2", 1: "namefs", 2: "nfs", 3: "jfs", 5: "cdrom", 6: "proc",
	16: "special-fs", 17: "cache-fs", 18: "nfs3", 19: "automount-fs", 20: "pool-fs", 32: "vxfs",
	33: "veritas-fs", 34: "udfs", 35: "nfs4", 36: "nfs4-pseudo", 37: "smbfs", 38: "mcr-pseudofs",
	39: "ahafs", 40: "sterm-nfs", 41: "asmfs",
}

func PartitionsWithContext(ctx context.Context, all bool) ([]PartitionStat, error) {
	var ret []PartitionStat

	out, err := invoke.CommandWithContext(ctx, "mount")
	if err != nil {
		return nil, err
	}

	// parse head lines for column names
	colidx := make(map[string]int)
	lines := strings.Split(string(out), "\n")
	if len(lines) < 3 {
		return nil, common.ErrNotImplementedError
	}

	idx := 0
	start := 0
	finished := false
	for pos, ch := range lines[1] {
		if ch == ' ' && !finished {
			name := strings.TrimSpace(lines[0][start:pos])
			colidx[name] = idx
			finished = true
		} else if ch == '-' && finished {
			idx++
			start = pos
			finished = false
		}
	}
	name := strings.TrimSpace(lines[0][start:len(lines[1])])
	colidx[name] = idx

	for idx := 2; idx < len(lines); idx++ {
		line := lines[idx]
		if startBlank.MatchString(line) {
			line = "localhost" + line
		}
		p := strings.Fields(lines[idx])
		if len(p) < 5 || ignoreFSType[p[colidx["vfs"]]] {
			continue
		}
		d := PartitionStat{
			Device:     p[colidx["mounted"]],
			Mountpoint: p[colidx["mounted over"]],
			Fstype:     p[colidx["vfs"]],
			Opts:       strings.Split(p[colidx["options"]], ","),
		}

		ret = append(ret, d)
	}

	return ret, nil
}

func getFsType(stat unix.Statfs_t) string {
	return FSType[int(stat.Vfstype)]
}

func UsageWithContext(ctx context.Context, path string) (*UsageStat, error) {
	out, err := invoke.CommandWithContext(ctx, "df", "-v")
	if err != nil {
		return nil, err
	}

	ret := &UsageStat{}

	blocksize := uint64(512)
	lines := strings.Split(string(out), "\n")
	if len(lines) < 2 {
		return &UsageStat{}, common.ErrNotImplementedError
	}

	hf := strings.Fields(strings.Replace(lines[0], "Mounted on", "Path", -1)) // headers
	for line := 1; line < len(lines); line++ {
		fs := strings.Fields(lines[line]) // values
		for i, header := range hf {
			// We're done in any of these use cases
			if i >= len(fs) {
				break
			}

			switch header {
			case `Filesystem`:
				// This is not a valid fs for us to parse
				if fs[i] == "/proc" || fs[i] == "/ahafs" || fs[i] != path {
					break
				}

				ret.Fstype, err = GetMountFSTypeWithContext(ctx, fs[i])
				if err != nil {
					return nil, err
				}
			case `Path`:
				ret.Path = fs[i]
			case `512-blocks`:
				total, err := strconv.ParseUint(fs[i], 10, 64)
				ret.Total = total * blocksize
				if err != nil {
					return nil, err
				}
			case `Used`:
				ret.Used, err = strconv.ParseUint(fs[i], 10, 64)
				if err != nil {
					return nil, err
				}
			case `Free`:
				ret.Free, err = strconv.ParseUint(fs[i], 10, 64)
				if err != nil {
					return nil, err
				}
			case `%Used`:
				val, err := strconv.Atoi(strings.Replace(fs[i], "%", "", -1))
				if err != nil {
					return nil, err
				}
				ret.UsedPercent = float64(val) / float64(100)
			case `Ifree`:
				ret.InodesFree, err = strconv.ParseUint(fs[i], 10, 64)
				if err != nil {
					return nil, err
				}
			case `Iused`:
				ret.InodesUsed, err = strconv.ParseUint(fs[i], 10, 64)
				if err != nil {
					return nil, err
				}
			case `%Iused`:
				val, err := strconv.Atoi(strings.Replace(fs[i], "%", "", -1))
				if err != nil {
					return nil, err
				}
				ret.InodesUsedPercent = float64(val) / float64(100)
			}
		}

		// Calculated value, since it isn't returned by the command
		ret.InodesTotal = ret.InodesUsed + ret.InodesFree

		// Valid Usage data, so append it
		return ret, nil
	}

	return ret, nil
}

func GetMountFSTypeWithContext(ctx context.Context, mp string) (string, error) {
	out, err := invoke.CommandWithContext(ctx, "mount")
	if err != nil {
		return "", err
	}

	// Kind of inefficient, but it works
	lines := strings.Split(string(out[:]), "\n")
	for line := 1; line < len(lines); line++ {
		fields := strings.Fields(lines[line])
		if strings.TrimSpace(fields[0]) == mp {
			return fields[2], nil
		}
	}

	return "", nil
}
