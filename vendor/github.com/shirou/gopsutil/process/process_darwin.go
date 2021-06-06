// +build darwin

package process

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/internal/common"
	"github.com/shirou/gopsutil/net"
	"golang.org/x/sys/unix"
)

// copied from sys/sysctl.h
const (
	CTLKern          = 1  // "high kernel": proc, limits
	KernProc         = 14 // struct: process entries
	KernProcPID      = 1  // by process id
	KernProcProc     = 8  // only return procs
	KernProcAll      = 0  // everything
	KernProcPathname = 12 // path to executable
)

const (
	ClockTicks = 100 // C.sysconf(C._SC_CLK_TCK)
)

type _Ctype_struct___0 struct {
	Pad uint64
}

func pidsWithContext(ctx context.Context) ([]int32, error) {
	var ret []int32

	pids, err := callPsWithContext(ctx, "pid", 0, false, false)
	if err != nil {
		return ret, err
	}

	for _, pid := range pids {
		v, err := strconv.Atoi(pid[0])
		if err != nil {
			return ret, err
		}
		ret = append(ret, int32(v))
	}

	return ret, nil
}

func (p *Process) PpidWithContext(ctx context.Context) (int32, error) {
	r, err := callPsWithContext(ctx, "ppid", p.Pid, false, false)
	if err != nil {
		return 0, err
	}

	v, err := strconv.Atoi(r[0][0])
	if err != nil {
		return 0, err
	}

	return int32(v), err
}

func (p *Process) NameWithContext(ctx context.Context) (string, error) {
	k, err := p.getKProc()
	if err != nil {
		return "", err
	}
	name := common.IntToString(k.Proc.P_comm[:])

	if len(name) >= 15 {
		cmdName, err := p.cmdNameWithContext(ctx)
		if err != nil {
			return "", err
		}
		if len(cmdName) > 0 {
			extendedName := filepath.Base(cmdName[0])
			if strings.HasPrefix(extendedName, p.name) {
				name = extendedName
			} else {
				name = cmdName[0]
			}
		}
	}

	return name, nil
}

func (p *Process) CmdlineWithContext(ctx context.Context) (string, error) {
	r, err := callPsWithContext(ctx, "command", p.Pid, false, false)
	if err != nil {
		return "", err
	}
	return strings.Join(r[0], " "), err
}

// cmdNameWithContext returns the command name (including spaces) without any arguments
func (p *Process) cmdNameWithContext(ctx context.Context) ([]string, error) {
	r, err := callPsWithContext(ctx, "command", p.Pid, false, true)
	if err != nil {
		return nil, err
	}
	return r[0], err
}

// CmdlineSliceWithContext returns the command line arguments of the process as a slice with each
// element being an argument. Because of current deficiencies in the way that the command
// line arguments are found, single arguments that have spaces in the will actually be
// reported as two separate items. In order to do something better CGO would be needed
// to use the native darwin functions.
func (p *Process) CmdlineSliceWithContext(ctx context.Context) ([]string, error) {
	r, err := callPsWithContext(ctx, "command", p.Pid, false, false)
	if err != nil {
		return nil, err
	}
	return r[0], err
}

func (p *Process) createTimeWithContext(ctx context.Context) (int64, error) {
	r, err := callPsWithContext(ctx, "etime", p.Pid, false, false)
	if err != nil {
		return 0, err
	}

	elapsedSegments := strings.Split(strings.Replace(r[0][0], "-", ":", 1), ":")
	var elapsedDurations []time.Duration
	for i := len(elapsedSegments) - 1; i >= 0; i-- {
		p, err := strconv.ParseInt(elapsedSegments[i], 10, 0)
		if err != nil {
			return 0, err
		}
		elapsedDurations = append(elapsedDurations, time.Duration(p))
	}

	var elapsed = time.Duration(elapsedDurations[0]) * time.Second
	if len(elapsedDurations) > 1 {
		elapsed += time.Duration(elapsedDurations[1]) * time.Minute
	}
	if len(elapsedDurations) > 2 {
		elapsed += time.Duration(elapsedDurations[2]) * time.Hour
	}
	if len(elapsedDurations) > 3 {
		elapsed += time.Duration(elapsedDurations[3]) * time.Hour * 24
	}

	start := time.Now().Add(-elapsed)
	return start.Unix() * 1000, nil
}

func (p *Process) ParentWithContext(ctx context.Context) (*Process, error) {
	out, err := common.CallLsofWithContext(ctx, invoke, p.Pid, "-FR")
	if err != nil {
		return nil, err
	}
	for _, line := range out {
		if len(line) >= 1 && line[0] == 'R' {
			v, err := strconv.Atoi(line[1:])
			if err != nil {
				return nil, err
			}
			return NewProcessWithContext(ctx, int32(v))
		}
	}
	return nil, fmt.Errorf("could not find parent line")
}

func (p *Process) StatusWithContext(ctx context.Context) (string, error) {
	r, err := callPsWithContext(ctx, "state", p.Pid, false, false)
	if err != nil {
		return "", err
	}

	return r[0][0][0:1], err
}

func (p *Process) ForegroundWithContext(ctx context.Context) (bool, error) {
	// see https://github.com/shirou/gopsutil/issues/596#issuecomment-432707831 for implementation details
	pid := p.Pid
	ps, err := exec.LookPath("ps")
	if err != nil {
		return false, err
	}
	out, err := invoke.CommandWithContext(ctx, ps, "-o", "stat=", "-p", strconv.Itoa(int(pid)))
	if err != nil {
		return false, err
	}
	return strings.IndexByte(string(out), '+') != -1, nil
}

func (p *Process) UidsWithContext(ctx context.Context) ([]int32, error) {
	k, err := p.getKProc()
	if err != nil {
		return nil, err
	}

	// See: http://unix.superglobalmegacorp.com/Net2/newsrc/sys/ucred.h.html
	userEffectiveUID := int32(k.Eproc.Ucred.UID)

	return []int32{userEffectiveUID}, nil
}

func (p *Process) GidsWithContext(ctx context.Context) ([]int32, error) {
	k, err := p.getKProc()
	if err != nil {
		return nil, err
	}

	gids := make([]int32, 0, 3)
	gids = append(gids, int32(k.Eproc.Pcred.P_rgid), int32(k.Eproc.Ucred.Ngroups), int32(k.Eproc.Pcred.P_svgid))

	return gids, nil
}

func (p *Process) GroupsWithContext(ctx context.Context) ([]int32, error) {
	return nil, common.ErrNotImplementedError
	// k, err := p.getKProc()
	// if err != nil {
	// 	return nil, err
	// }

	// groups := make([]int32, k.Eproc.Ucred.Ngroups)
	// for i := int16(0); i < k.Eproc.Ucred.Ngroups; i++ {
	// 	groups[i] = int32(k.Eproc.Ucred.Groups[i])
	// }

	// return groups, nil
}

func (p *Process) TerminalWithContext(ctx context.Context) (string, error) {
	return "", common.ErrNotImplementedError
	/*
		k, err := p.getKProc()
		if err != nil {
			return "", err
		}

		ttyNr := uint64(k.Eproc.Tdev)
		termmap, err := getTerminalMap()
		if err != nil {
			return "", err
		}

		return termmap[ttyNr], nil
	*/
}

func (p *Process) NiceWithContext(ctx context.Context) (int32, error) {
	k, err := p.getKProc()
	if err != nil {
		return 0, err
	}
	return int32(k.Proc.P_nice), nil
}

func (p *Process) IOCountersWithContext(ctx context.Context) (*IOCountersStat, error) {
	return nil, common.ErrNotImplementedError
}

func (p *Process) NumThreadsWithContext(ctx context.Context) (int32, error) {
	r, err := callPsWithContext(ctx, "utime,stime", p.Pid, true, false)
	if err != nil {
		return 0, err
	}
	return int32(len(r)), nil
}

func convertCPUTimes(s string) (ret float64, err error) {
	var t int
	var _tmp string
	if strings.Contains(s, ":") {
		_t := strings.Split(s, ":")
		switch len(_t) {
		case 3:
			hour, err := strconv.Atoi(_t[0])
			if err != nil {
				return ret, err
			}
			t += hour * 60 * 60 * ClockTicks

			mins, err := strconv.Atoi(_t[1])
			if err != nil {
				return ret, err
			}
			t += mins * 60 * ClockTicks
			_tmp = _t[2]
		case 2:
			mins, err := strconv.Atoi(_t[0])
			if err != nil {
				return ret, err
			}
			t += mins * 60 * ClockTicks
			_tmp = _t[1]
		case 1, 0:
			_tmp = s
		default:
			return ret, fmt.Errorf("wrong cpu time string")
		}
	} else {
		_tmp = s
	}

	_t := strings.Split(_tmp, ".")
	if err != nil {
		return ret, err
	}
	h, err := strconv.Atoi(_t[0])
	t += h * ClockTicks
	h, err = strconv.Atoi(_t[1])
	t += h
	return float64(t) / ClockTicks, nil
}

func (p *Process) TimesWithContext(ctx context.Context) (*cpu.TimesStat, error) {
	r, err := callPsWithContext(ctx, "utime,stime", p.Pid, false, false)

	if err != nil {
		return nil, err
	}

	utime, err := convertCPUTimes(r[0][0])
	if err != nil {
		return nil, err
	}
	stime, err := convertCPUTimes(r[0][1])
	if err != nil {
		return nil, err
	}

	ret := &cpu.TimesStat{
		CPU:    "cpu",
		User:   utime,
		System: stime,
	}
	return ret, nil
}

func (p *Process) MemoryInfoWithContext(ctx context.Context) (*MemoryInfoStat, error) {
	r, err := callPsWithContext(ctx, "rss,vsize,pagein", p.Pid, false, false)
	if err != nil {
		return nil, err
	}
	rss, err := strconv.Atoi(r[0][0])
	if err != nil {
		return nil, err
	}
	vms, err := strconv.Atoi(r[0][1])
	if err != nil {
		return nil, err
	}
	pagein, err := strconv.Atoi(r[0][2])
	if err != nil {
		return nil, err
	}

	ret := &MemoryInfoStat{
		RSS:  uint64(rss) * 1024,
		VMS:  uint64(vms) * 1024,
		Swap: uint64(pagein),
	}

	return ret, nil
}

func (p *Process) ChildrenWithContext(ctx context.Context) ([]*Process, error) {
	pids, err := common.CallPgrepWithContext(ctx, invoke, p.Pid)
	if err != nil {
		return nil, err
	}
	ret := make([]*Process, 0, len(pids))
	for _, pid := range pids {
		np, err := NewProcessWithContext(ctx, pid)
		if err != nil {
			return nil, err
		}
		ret = append(ret, np)
	}
	return ret, nil
}

func (p *Process) ConnectionsWithContext(ctx context.Context) ([]net.ConnectionStat, error) {
	return net.ConnectionsPidWithContext(ctx, "all", p.Pid)
}

func (p *Process) ConnectionsMaxWithContext(ctx context.Context, max int) ([]net.ConnectionStat, error) {
	return net.ConnectionsPidMaxWithContext(ctx, "all", p.Pid, max)
}

func ProcessesWithContext(ctx context.Context) ([]*Process, error) {
	out := []*Process{}

	pids, err := PidsWithContext(ctx)
	if err != nil {
		return out, err
	}

	for _, pid := range pids {
		p, err := NewProcessWithContext(ctx, pid)
		if err != nil {
			continue
		}
		out = append(out, p)
	}

	return out, nil
}

// Returns a proc as defined here:
// http://unix.superglobalmegacorp.com/Net2/newsrc/sys/kinfo_proc.h.html
func (p *Process) getKProc() (*KinfoProc, error) {
	buf, err := unix.SysctlRaw("kern.proc.pid", int(p.Pid))
	if err != nil {
		return nil, err
	}
	k, err := parseKinfoProc(buf)
	if err != nil {
		return nil, err
	}

	return &k, nil
}

// call ps command.
// Return value deletes Header line(you must not input wrong arg).
// And split by space. Caller have responsibility to manage.
// If passed arg pid is 0, get information from all process.
func callPsWithContext(ctx context.Context, arg string, pid int32, threadOption bool, nameOption bool) ([][]string, error) {
	bin, err := exec.LookPath("ps")
	if err != nil {
		return [][]string{}, err
	}

	var cmd []string
	if pid == 0 { // will get from all processes.
		cmd = []string{"-ax", "-o", arg}
	} else if threadOption {
		cmd = []string{"-x", "-o", arg, "-M", "-p", strconv.Itoa(int(pid))}
	} else {
		cmd = []string{"-x", "-o", arg, "-p", strconv.Itoa(int(pid))}
	}

	if nameOption {
		cmd = append(cmd, "-c")
	}
	out, err := invoke.CommandWithContext(ctx, bin, cmd...)
	if err != nil {
		return [][]string{}, err
	}
	lines := strings.Split(string(out), "\n")

	var ret [][]string
	for _, l := range lines[1:] {

		var lr []string
		if nameOption {
			lr = append(lr, l)
		} else {
			for _, r := range strings.Split(l, " ") {
				if r == "" {
					continue
				}
				lr = append(lr, strings.TrimSpace(r))
			}
		}

		if len(lr) != 0 {
			ret = append(ret, lr)
		}
	}

	return ret, nil
}
