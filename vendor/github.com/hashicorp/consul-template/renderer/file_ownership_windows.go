// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build windows
// +build windows

package renderer

import (
	"fmt"
)

var notSupportedError error = fmt.Errorf("managing file ownership is not supported on Windows")

func setFileOwnership(path string, uid, gid int) error {
	if uid == -1 && gid == -1 {
		return nil
	}
	return notSupportedError
}

func isChownNeeded(path string, wantedUid, wantedGid int) (bool, error) {
	if wantedUid == -1 && wantedGid == -1 {
		return false, nil
	}
	return false, notSupportedError
}

func lookupUser(user string) (int, error) {
	if user == "" {
		return -1, nil
	}
	return 0, notSupportedError
}

func lookupGroup(group string) (int, error) {
	if group == "" {
		return -1, nil
	}
	return 0, notSupportedError
}
