package mlock

// This should be set by the OS-specific packages to tell whether LockMemory
// is supported or not.
var supported bool

// Supported returns true if LockMemory is functional on this system.
func Supported() bool {
	return supported
}

// LockMemory prevents any memory from being swapped to disk.
func LockMemory() error {
	return lockMemory()
}
