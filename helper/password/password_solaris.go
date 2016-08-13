// +build solaris

package password

import (
	"fmt"
	"os"
	"syscall"

	"golang.org/x/sys/unix"
)

func read(f *os.File) (string, error) {
	fd := int(f.Fd())
	if !isTerminal(fd) {
		return "", fmt.Errorf("File descriptor %d is not a terminal", fd)
	}

	oldState, err := makeRaw(fd)
	if err != nil {
		return "", err
	}
	defer unix.IoctlSetTermios(fd, unix.TCSETS, oldState)

	return readline(f)
}

// isTerminal returns true if there is a terminal attached to the given
// file descriptor.
// Source: http://src.illumos.org/source/xref/illumos-gate/usr/src/lib/libbc/libc/gen/common/isatty.c
func isTerminal(fd int) bool {
	var termio unix.Termio
	err := unix.IoctlSetTermio(fd, unix.TCGETA, &termio)
	return err == nil
}

// makeRaw puts the terminal connected to the given file descriptor into raw
// mode and returns the previous state of the terminal so that it can be
// restored.
// Source: http://src.illumos.org/source/xref/illumos-gate/usr/src/lib/libast/common/uwin/getpass.c
func makeRaw(fd int) (*unix.Termios, error) {
	oldTermiosPtr, err := unix.IoctlGetTermios(int(fd), unix.TCGETS)
	if err != nil {
		return nil, err
	}
	oldTermios := *oldTermiosPtr

	newTermios := oldTermios
	newTermios.Lflag &^= syscall.ECHO | syscall.ECHOE | syscall.ECHOK | syscall.ECHONL
	if err := unix.IoctlSetTermios(fd, unix.TCSETS, &newTermios); err != nil {
		return nil, err
	}

	return oldTermiosPtr, nil
}
