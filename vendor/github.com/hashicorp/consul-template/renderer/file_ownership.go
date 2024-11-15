// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build !windows
// +build !windows

package renderer

import (
	"os"
	"os/user"
	"strconv"
	"syscall"
)

func getFileOwnership(path string) (int, int, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return 0, 0, err
	}

	fileSys := fileInfo.Sys()
	st := fileSys.(*syscall.Stat_t)
	return int(st.Uid), int(st.Gid), nil
}

func setFileOwnership(path string, uid, gid int) error {
	if uid == -1 && gid == -1 {
		return nil // noop
	}
	return os.Chown(path, uid, gid)
}

func isChownNeeded(path string, uid, gid int) (bool, error) {
	if uid == -1 && gid == -1 {
		return false, nil
	}

	currUid, currGid, err := getFileOwnership(path)
	if err != nil {
		return false, err
	}

	switch {
	case uid == -1:
		currUid = -1
	case gid == -1:
		currGid = -1
	}

	return uid != currUid || gid != currGid, nil
}

// parseUidGid parses the uid/gid so that it can be input to os.Chown
func parseUidGid(s string) (int, error) {
	if s == "" {
		return -1, nil
	}
	return strconv.Atoi(s)
}

func lookupUser(s string) (int, error) {
	if id, err := parseUidGid(s); err == nil {
		return id, nil
	}
	u, err := user.Lookup(s)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(u.Uid)
}

func lookupGroup(s string) (int, error) {
	if id, err := parseUidGid(s); err == nil {
		return id, nil
	}
	u, err := user.LookupGroup(s)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(u.Gid)
}
