// +build !windows

package diagnose

import (
	"io/fs"
	"syscall"
)

// IsOwnedByRoot checks if a file is owned by root
func IsOwnedByRoot(info fs.FileInfo) bool {
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		uid := int(stat.Uid)
		gid := int(stat.Gid)
		if uid == 0 && gid == 0 {
			return true
		}
	}
	return false
}
