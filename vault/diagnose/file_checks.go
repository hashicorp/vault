package diagnose

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
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

// HasCorrectFilePerms checks if the specified file has owner rw perms
// and no other permissions
func HasCorrectFilePerms(info fs.FileInfo) (bool, []string) {
	var errors []string
	mode := info.Mode()

	//check that owners have read and write permissions
	if mode&(1<<7) == 0 || mode&(1<<8) == 0 {
		errors = append(errors, fmt.Sprintf(FilePermissionsMissingWarning+": perms are %s", mode.String()))
	}

	// Check user rw and group rw for overpermissions
	if mode&(1<<1) != 0 || mode&(1<<2) != 0 || mode&(1<<4) != 0 || mode&(1<<5) != 0 {
		errors = append(errors, fmt.Sprintf(FileTooPermissiveWarning+": perms are %s", mode.String()))
	}

	if mode&os.ModeSymlink != 0 {
		errors = append(errors, FileIsSymlinkWarning)
	}

	if len(errors) > 0 {
		return false, errors
	}
	return true, nil
}
