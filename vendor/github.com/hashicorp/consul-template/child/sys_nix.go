// +build !windows

package child

import (
	"os/exec"
	"syscall"
)

func setSetpgid(cmd *exec.Cmd, value bool) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: value}
}

func processNotFoundErr(err error) bool {
	// ESRCH == no such process, ie. already exited
	return err == syscall.ESRCH
}
