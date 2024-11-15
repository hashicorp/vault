// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build !windows
// +build !windows

package child

import (
	"os/exec"
	"syscall"
)

func setSysProcAttr(cmd *exec.Cmd, setpgid, setsid bool) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: setpgid,
		Setsid:  setsid,
	}
}

func processNotFoundErr(err error) bool {
	// ESRCH == no such process, ie. already exited
	return err == syscall.ESRCH
}
