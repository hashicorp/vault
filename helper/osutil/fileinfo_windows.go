// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build windows

package osutil

import (
	"io/fs"
	"os"
)

func FileUidMatch(info fs.FileInfo, path string, uid int) error {
	return nil
}

// Umask does nothing for windows for now
func Umask(newmask int) int {
	return 0
}

func Chown(f *os.File, owner, group int) error {
	return nil
}
