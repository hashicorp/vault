// +build solaris,!appengine

package logrus

import (
	"os"

	"golang.org/x/sys/unix"
)

// IsTerminal returns true if the given file descriptor is a terminal.
func IsTerminal(f io.Writer) bool {
	var termios Termios
	switch v := f.(type) {
	case *os.File:
		_, err := unix.IoctlGetTermios(int(v.Fd()), unix.TCGETA)
		return err == nil
	default:
		return false
	}
}
