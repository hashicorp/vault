// +build windows

package signals

import (
	"os"
	"syscall"
)

// RuntimeSig is set to nil on windows as it doesn't support the signal (SIGURG)
var RuntimeSig = os.Signal(nil)

var SignalLookup = map[string]os.Signal{
	"SIGABRT": syscall.SIGABRT,
	"SIGALRM": syscall.SIGALRM,
	"SIGBUS":  syscall.SIGBUS,
	"SIGFPE":  syscall.SIGFPE,
	"SIGHUP":  syscall.SIGHUP,
	"SIGILL":  syscall.SIGILL,
	"SIGINT":  syscall.SIGINT,
	"SIGKILL": syscall.SIGKILL,
	"SIGPIPE": syscall.SIGPIPE,
	"SIGQUIT": syscall.SIGQUIT,
	"SIGSEGV": syscall.SIGSEGV,
	"SIGTERM": syscall.SIGTERM,
	"SIGTRAP": syscall.SIGTRAP,
}
