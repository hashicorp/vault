// +build windows

package password

import (
	"os"
	"syscall"
)

var (
	kernel32           = syscall.MustLoadDLL("kernel32.dll")
	setConsoleModeProc = kernel32.MustFindProc("SetConsoleMode")
)

// Magic constant from MSDN to control whether characters read are
// repeated back on the console.
//
// http://msdn.microsoft.com/en-us/library/windows/desktop/ms686033(v=vs.85).aspx
const ENABLE_ECHO_INPUT = 0x0004

func read(f *os.File) (string, error) {
	handle := syscall.Handle(f.Fd())

	// Grab the old console mode so we can reset it. We defer the reset
	// right away because it doesn't matter (it is idempotent).
	var oldMode uint32
	if err := syscall.GetConsoleMode(handle, &oldMode); err != nil {
		return "", err
	}
	defer setConsoleMode(handle, oldMode)

	// The new mode is the old mode WITHOUT the echo input flag set.
	var newMode uint32 = uint32(int(oldMode) & ^ENABLE_ECHO_INPUT)
	if err := setConsoleMode(handle, newMode); err != nil {
		return "", err
	}

	return readline(f)
}

func setConsoleMode(console syscall.Handle, mode uint32) error {
	r, _, err := setConsoleModeProc.Call(uintptr(console), uintptr(mode))
	if r == 0 {
		return err
	}

	return nil
}
