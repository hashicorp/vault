// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !windows

package osutil

import (
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"strconv"
	"syscall"
)

func FileUIDEqual(info fs.FileInfo, uid int) bool {
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		path_uid := int(stat.Uid)
		if path_uid == uid {
			return true
		}
	}
	return false
}

func FileGIDEqual(info fs.FileInfo, gid int) bool {
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		path_gid := int(stat.Gid)
		if path_gid == gid {
			return true
		}
	}
	return false
}

func FileUidMatch(info fs.FileInfo, path string, uid int) (err error) {
	currentUser, err := user.Current()
	if err != nil {
		return fmt.Errorf("failed to get details of current process owner. The error is: %w", err)
	}
	switch uid {
	case 0:
		currentUserUid, err := strconv.Atoi(currentUser.Uid)
		if err != nil {
			return fmt.Errorf("failed to convert uid %q to int. The error is: %w", currentUser.Uid, err)
		}
		if !FileUIDEqual(info, currentUserUid) {
			return fmt.Errorf("path %q is not owned by my uid %s", path, currentUser.Uid)
		}
	default:
		if !FileUIDEqual(info, uid) {
			return fmt.Errorf("path %q is not owned by uid %d", path, uid)
		}
	}
	return err
}

// Sets new umask and returns old umask
func Umask(newmask int) int {
	return syscall.Umask(newmask)
}

func Chown(f *os.File, owner, group int) error {
	return f.Chown(owner, group)
}
