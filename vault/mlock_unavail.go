// +build windows plan9 darwin freebsd openbsd solaris

package vault

// LockMemory is used to prevent any memory from being swapped to disk
func LockMemory() error {
	// XXX: No good way to do this on Windows. There is the VirtualLock
	// method, but it requires a specific address and offset.
	return nil
}
