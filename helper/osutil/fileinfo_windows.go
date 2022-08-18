//go:build windows

package osutil

// Umask does nothing for windows for now
func Umask(newmask int) int {
	return 0
}
