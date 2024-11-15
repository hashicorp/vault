package process

import (
	"bytes"
	"context"
	"os"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/internal/common"
	"github.com/shirou/gopsutil/v3/net"
)

type MemoryMapsStat struct {
	Path         string `json:"path"`
	Rss          uint64 `json:"rss"`
	Size         uint64 `json:"size"`
	Pss          uint64 `json:"pss"`
	SharedClean  uint64 `json:"sharedClean"`
	SharedDirty  uint64 `json:"sharedDirty"`
	PrivateClean uint64 `json:"privateClean"`
	PrivateDirty uint64 `json:"privateDirty"`
	Referenced   uint64 `json:"referenced"`
	Anonymous    uint64 `json:"anonymous"`
	Swap         uint64 `json:"swap"`
}

type MemoryInfoExStat struct{}

func pidsWithContext(ctx context.Context) ([]int32, error) {
	return readPidsFromDir(common.HostProcWithContext(ctx))
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

func (p *Process) PpidWithContext(ctx context.Context) (int32, error) {
	return 0, common.ErrNotImplementedError
}

func (p *Process) NameWithContext(ctx context.Context) (string, error) {
	return "", common.ErrNotImplementedError
}

func (p *Process) TgidWithContext(ctx context.Context) (int32, error) {
	return 0, common.ErrNotImplementedError
}

func (p *Process) ExeWithContext(ctx context.Context) (string, error) {
	exe, err := p.fillFromPathAOutWithContext(ctx)
	if os.IsNotExist(err) {
		exe, err = p.fillFromExecnameWithContext(ctx)
	}
	return exe, err
}

func (p *Process) CmdlineWithContext(ctx context.Context) (string, error) {
	return p.fillFromCmdlineWithContext(ctx)
}

func (p *Process) CmdlineSliceWithContext(ctx context.Context) ([]string, error) {
	return p.fillSliceFromCmdlineWithContext(ctx)
}

func (p *Process) createTimeWithContext(ctx context.Context) (int64, error) {
	return 0, common.ErrNotImplementedError
}

func (p *Process) CwdWithContext(ctx context.Context) (string, error) {
	return p.fillFromPathCwdWithContext(ctx)
}

func (p *Process) StatusWithContext(ctx context.Context) ([]string, error) {
	return []string{""}, common.ErrNotImplementedError
}

func (p *Process) ForegroundWithContext(ctx context.Context) (bool, error) {
	return false, common.ErrNotImplementedError
}

func (p *Process) UidsWithContext(ctx context.Context) ([]int32, error) {
	return nil, common.ErrNotImplementedError
}

func (p *Process) GidsWithContext(ctx context.Context) ([]int32, error) {
	return nil, common.ErrNotImplementedError
}

func (p *Process) GroupsWithContext(ctx context.Context) ([]int32, error) {
	return nil, common.ErrNotImplementedError
}

func (p *Process) TerminalWithContext(ctx context.Context) (string, error) {
	return "", common.ErrNotImplementedError
}

func (p *Process) NiceWithContext(ctx context.Context) (int32, error) {
	return 0, common.ErrNotImplementedError
}

func (p *Process) IOniceWithContext(ctx context.Context) (int32, error) {
	return 0, common.ErrNotImplementedError
}

func (p *Process) RlimitWithContext(ctx context.Context) ([]RlimitStat, error) {
	return nil, common.ErrNotImplementedError
}

func (p *Process) RlimitUsageWithContext(ctx context.Context, gatherUsed bool) ([]RlimitStat, error) {
	return nil, common.ErrNotImplementedError
}

func (p *Process) IOCountersWithContext(ctx context.Context) (*IOCountersStat, error) {
	return nil, common.ErrNotImplementedError
}

func (p *Process) NumCtxSwitchesWithContext(ctx context.Context) (*NumCtxSwitchesStat, error) {
	return nil, common.ErrNotImplementedError
}

func (p *Process) NumFDsWithContext(ctx context.Context) (int32, error) {
	_, fnames, err := p.fillFromfdListWithContext(ctx)
	return int32(len(fnames)), err
}

func (p *Process) NumThreadsWithContext(ctx context.Context) (int32, error) {
	return 0, common.ErrNotImplementedError
}

func (p *Process) ThreadsWithContext(ctx context.Context) (map[int32]*cpu.TimesStat, error) {
	return nil, common.ErrNotImplementedError
}

func (p *Process) TimesWithContext(ctx context.Context) (*cpu.TimesStat, error) {
	return nil, common.ErrNotImplementedError
}

func (p *Process) CPUAffinityWithContext(ctx context.Context) ([]int32, error) {
	return nil, common.ErrNotImplementedError
}

func (p *Process) MemoryInfoWithContext(ctx context.Context) (*MemoryInfoStat, error) {
	return nil, common.ErrNotImplementedError
}

func (p *Process) MemoryInfoExWithContext(ctx context.Context) (*MemoryInfoExStat, error) {
	return nil, common.ErrNotImplementedError
}

func (p *Process) PageFaultsWithContext(ctx context.Context) (*PageFaultsStat, error) {
	return nil, common.ErrNotImplementedError
}

func (p *Process) ChildrenWithContext(ctx context.Context) ([]*Process, error) {
	return nil, common.ErrNotImplementedError
}

func (p *Process) OpenFilesWithContext(ctx context.Context) ([]OpenFilesStat, error) {
	return nil, common.ErrNotImplementedError
}

func (p *Process) ConnectionsWithContext(ctx context.Context) ([]net.ConnectionStat, error) {
	return nil, common.ErrNotImplementedError
}

func (p *Process) ConnectionsMaxWithContext(ctx context.Context, max int) ([]net.ConnectionStat, error) {
	return nil, common.ErrNotImplementedError
}

func (p *Process) MemoryMapsWithContext(ctx context.Context, grouped bool) (*[]MemoryMapsStat, error) {
	return nil, common.ErrNotImplementedError
}

func (p *Process) EnvironWithContext(ctx context.Context) ([]string, error) {
	return nil, common.ErrNotImplementedError
}

/**
** Internal functions
**/

func (p *Process) fillFromfdListWithContext(ctx context.Context) (string, []string, error) {
	pid := p.Pid
	statPath := common.HostProcWithContext(ctx, strconv.Itoa(int(pid)), "fd")
	d, err := os.Open(statPath)
	if err != nil {
		return statPath, []string{}, err
	}
	defer d.Close()
	fnames, err := d.Readdirnames(-1)
	return statPath, fnames, err
}

func (p *Process) fillFromPathCwdWithContext(ctx context.Context) (string, error) {
	pid := p.Pid
	cwdPath := common.HostProcWithContext(ctx, strconv.Itoa(int(pid)), "path", "cwd")
	cwd, err := os.Readlink(cwdPath)
	if err != nil {
		return "", err
	}
	return cwd, nil
}

func (p *Process) fillFromPathAOutWithContext(ctx context.Context) (string, error) {
	pid := p.Pid
	cwdPath := common.HostProcWithContext(ctx, strconv.Itoa(int(pid)), "path", "a.out")
	exe, err := os.Readlink(cwdPath)
	if err != nil {
		return "", err
	}
	return exe, nil
}

func (p *Process) fillFromExecnameWithContext(ctx context.Context) (string, error) {
	pid := p.Pid
	execNamePath := common.HostProcWithContext(ctx, strconv.Itoa(int(pid)), "execname")
	exe, err := os.ReadFile(execNamePath)
	if err != nil {
		return "", err
	}
	return string(exe), nil
}

func (p *Process) fillFromCmdlineWithContext(ctx context.Context) (string, error) {
	pid := p.Pid
	cmdPath := common.HostProcWithContext(ctx, strconv.Itoa(int(pid)), "cmdline")
	cmdline, err := os.ReadFile(cmdPath)
	if err != nil {
		return "", err
	}
	ret := strings.FieldsFunc(string(cmdline), func(r rune) bool {
		if r == '\u0000' {
			return true
		}
		return false
	})

	return strings.Join(ret, " "), nil
}

func (p *Process) fillSliceFromCmdlineWithContext(ctx context.Context) ([]string, error) {
	pid := p.Pid
	cmdPath := common.HostProcWithContext(ctx, strconv.Itoa(int(pid)), "cmdline")
	cmdline, err := os.ReadFile(cmdPath)
	if err != nil {
		return nil, err
	}
	if len(cmdline) == 0 {
		return nil, nil
	}
	if cmdline[len(cmdline)-1] == 0 {
		cmdline = cmdline[:len(cmdline)-1]
	}
	parts := bytes.Split(cmdline, []byte{0})
	var strParts []string
	for _, p := range parts {
		strParts = append(strParts, string(p))
	}

	return strParts, nil
}

func readPidsFromDir(path string) ([]int32, error) {
	var ret []int32

	d, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer d.Close()

	fnames, err := d.Readdirnames(-1)
	if err != nil {
		return nil, err
	}
	for _, fname := range fnames {
		pid, err := strconv.ParseInt(fname, 10, 32)
		if err != nil {
			// if not numeric name, just skip
			continue
		}
		ret = append(ret, int32(pid))
	}

	return ret, nil
}
