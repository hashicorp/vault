//go:build !windows

package osutil

import (
	"syscall"
)

// Sets new umask and returns old umask
func Umask(newmask int) int {
	return syscall.Umask(newmask)
}
