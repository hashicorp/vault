// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build linux || darwin || freebsd || openbsd || solaris || netbsd
// +build linux darwin freebsd openbsd solaris netbsd

package signals

import (
	"os"
	"syscall"
)

//// Ignored Signals
// SIGCHLD - don't propagate these to child process as we manage it instead
// SIGURG  - used by the golang scheduler for parallel runtime.

var SignalLookup = map[string]os.Signal{
	"SIGNULL":  SIGNULL,
	"SIGABRT":  syscall.SIGABRT,
	"SIGALRM":  syscall.SIGALRM,
	"SIGBUS":   syscall.SIGBUS,
	"SIGCONT":  syscall.SIGCONT,
	"SIGFPE":   syscall.SIGFPE,
	"SIGHUP":   syscall.SIGHUP,
	"SIGILL":   syscall.SIGILL,
	"SIGINT":   syscall.SIGINT,
	"SIGIO":    syscall.SIGIO,
	"SIGIOT":   syscall.SIGIOT,
	"SIGKILL":  syscall.SIGKILL,
	"SIGPIPE":  syscall.SIGPIPE,
	"SIGPROF":  syscall.SIGPROF,
	"SIGQUIT":  syscall.SIGQUIT,
	"SIGSEGV":  syscall.SIGSEGV,
	"SIGSTOP":  syscall.SIGSTOP,
	"SIGSYS":   syscall.SIGSYS,
	"SIGTERM":  syscall.SIGTERM,
	"SIGTRAP":  syscall.SIGTRAP,
	"SIGTSTP":  syscall.SIGTSTP,
	"SIGTTIN":  syscall.SIGTTIN,
	"SIGTTOU":  syscall.SIGTTOU,
	"SIGUSR1":  syscall.SIGUSR1,
	"SIGUSR2":  syscall.SIGUSR2,
	"SIGWINCH": syscall.SIGWINCH,
	"SIGXCPU":  syscall.SIGXCPU,
	"SIGXFSZ":  syscall.SIGXFSZ,
}
