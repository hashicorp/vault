// +build linux

package mlock

import "syscall"

func init() {
	supported = true
}

func lockMemory() error {
	// Mlockall prevents all current and future pages from being swapped out.
	return syscall.Mlockall(syscall.MCL_CURRENT | syscall.MCL_FUTURE)
}
