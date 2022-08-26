package diagnose

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

const (
	FileIsSymlinkWarning          = "raft storage backend file is a symlink"
	FileTooPermissiveWarning      = "too many permissions"
	FilePermissionsMissingWarning = "owner or group needs read and write permissions"
)

func IsDir(info fs.FileInfo) bool {
	if info.Mode().IsDir() {
		return true
	}
	return false
}

func HasDB(path string) bool {
	dbPath := filepath.Join(path, DatabaseFilename)
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return false
	}
	return true
}

// CheckFilePerms checks if the specified file does not have other permissions, and
// whether the specified file just has owner rw permissions.
func CheckFilePerms(info fs.FileInfo) (bool, []string) {
	var errors []string
	mode := info.Mode()
	hasOnlyOwnerRW := false
	hasOwnerRead := false
	hasOwnerWrite := false
	hasSomeRead := false
	hasSomeWrite := false

	// Check owner perms
	if mode&0o400 != 0 {
		hasSomeRead = true
		hasOwnerRead = true
	}
	if mode&0o200 != 0 {
		hasSomeWrite = true
		hasOwnerWrite = true
	}

	if hasOwnerRead && hasOwnerWrite {
		hasOnlyOwnerRW = true
	}

	// These are "other" perms.
	// These don't count has "some read" or "some write" permissions because there should
	// never be a case when these permissions are set.
	if mode&0o007 != 0 {
		hasOnlyOwnerRW = false
		errors = append(errors, fmt.Sprintf(FileTooPermissiveWarning+": perms are %s", mode.String()))
	}

	// Check group permissions
	if mode&0o040 != 0 {
		hasOnlyOwnerRW = false
		hasSomeRead = true
	}
	if mode&0o020 != 0 {
		hasOnlyOwnerRW = false
		hasSomeWrite = true
	}

	// check that owners have read and write permissions
	if !hasSomeRead || !hasSomeWrite {
		errors = append(errors, fmt.Sprintf(FilePermissionsMissingWarning+": perms are %s", mode.String()))
	}

	if mode&os.ModeSymlink != 0 {
		errors = append(errors, FileIsSymlinkWarning)
	}

	if hasOnlyOwnerRW {
		return true, errors
	}
	return false, errors
}
