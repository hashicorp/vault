//go:build windows

package osutil

import (
	"io/fs"
)

// Umask does nothing for windows for now
func Umask(newmask int) int {
	return 0
}
