// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package osutil

import (
	"fmt"
	"io/fs"
	"os"
)

func IsWriteGroup(mode os.FileMode) bool {
	return mode&0o20 != 0
}

func IsWriteOther(mode os.FileMode) bool {
	return mode&0o02 != 0
}

func checkPathInfo(info fs.FileInfo, path string, uid int, permissions int) error {
	err := FileUidMatch(info, path, uid)
	if err != nil {
		return err
	}
	err = FilePermissionsMatch(info, path, permissions)
	if err != nil {
		return err
	}
	return nil
}

func FilePermissionsMatch(info fs.FileInfo, path string, permissions int) error {
	if permissions != 0 && int(info.Mode().Perm()) != permissions {
		return fmt.Errorf("path %q does not have permissions %o", path, permissions)
	}
	if permissions == 0 && (IsWriteOther(info.Mode()) || IsWriteGroup(info.Mode())) {
		return fmt.Errorf("path %q has insecure permissions %o. Vault expects no write permissions for group or others", path, info.Mode().Perm())
	}

	return nil
}

// OwnerPermissionsMatch checks if vault user is the owner and permissions are secure for input path
func OwnerPermissionsMatch(path string, uid int, permissions int) error {
	if path == "" {
		return fmt.Errorf("could not verify permissions for path. No path provided ")
	}

	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("error stating %q: %w", path, err)
	}
	if info.Mode()&os.ModeSymlink != 0 {
		symLinkInfo, err := os.Lstat(path)
		if err != nil {
			return fmt.Errorf("error stating %q: %w", path, err)
		}
		err = checkPathInfo(symLinkInfo, path, uid, permissions)
		if err != nil {
			return err
		}
	}
	err = checkPathInfo(info, path, uid, permissions)
	if err != nil {
		return err
	}

	return nil
}

// OwnerPermissionsMatchFile checks if vault user is the owner and permissions are secure for the input file
func OwnerPermissionsMatchFile(file *os.File, uid int, permissions int) error {
	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("file stat error on path %q: %w", file.Name(), err)
	}
	err = checkPathInfo(info, file.Name(), uid, permissions)
	if err != nil {
		return err
	}

	return nil
}
