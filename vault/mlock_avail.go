// +build linux

package vault

import "syscall"

// LockMemory is used to prevent any memory from being swapped to disk
func LockMemory() error {
	// Mlockall prevents all current and future pages from being
	// swapped out.
	err := syscall.Mlockall(syscall.MCL_CURRENT | syscall.MCL_FUTURE)

	// Catch if this is not implemented (darwin for example)
	if err == syscall.ENOSYS {
		err = nil
	}
	return err
}
