// +build linux darwin freebsd openbsd solaris netbsd

package signals

import (
	"os"
	"syscall"
)

var SignalLookup = map[string]os.Signal{
	"SIGABRT":  syscall.SIGABRT,
	"SIGALRM":  syscall.SIGALRM,
	"SIGBUS":   syscall.SIGBUS,
	"SIGCHLD":  syscall.SIGCHLD,
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
	"SIGURG":   syscall.SIGURG,
	"SIGUSR1":  syscall.SIGUSR1,
	"SIGUSR2":  syscall.SIGUSR2,
	"SIGWINCH": syscall.SIGWINCH,
	"SIGXCPU":  syscall.SIGXCPU,
	"SIGXFSZ":  syscall.SIGXFSZ,
}
