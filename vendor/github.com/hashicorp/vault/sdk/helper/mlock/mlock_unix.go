// +build dragonfly freebsd linux openbsd solaris

package mlock

import (
	"syscall"

	"golang.org/x/sys/unix"
)

func init() {
	supported = true
}

func lockMemory() error {
	// Mlockall prevents all current and future pages from being swapped out.
	return unix.Mlockall(syscall.MCL_CURRENT | syscall.MCL_FUTURE)
}
