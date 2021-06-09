// +build !windows

package diagnose

import (
	"fmt"
	"io/fs"
	"os"
	"syscall"
)

// IsOwnedByRoot checks if a particular file is owned by root or not
func IsOwnedByRoot(info fs.FileInfo) (bool, error) {
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		uid := int(stat.Uid)
		gid := int(stat.Gid)
		if uid == 0 && gid == 0 {
			return true, nil
		}
		return false, nil
	}
	return false, fmt.Errorf("permissions could not be determined")
}

// IsProcRoot checks if vault is running as root.
func IsProcRoot() bool {
	vaultUid := os.Getuid()
	vaultGid := os.Getgid()
	if vaultUid == 0 && vaultGid == 0 {
		return true
	}
	return false
}
