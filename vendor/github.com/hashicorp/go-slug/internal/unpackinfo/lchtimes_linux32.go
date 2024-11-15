// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build linux_amd || linux_arm
// +build linux_amd linux_arm

package unpackinfo

import (
	"golang.org/x/sys/unix"
)

// Lchtimes modifies the access and modified timestamps on a target path
// This capability is only available on Linux and Darwin as of now.
func (i UnpackInfo) Lchtimes() error {
	return unix.Lutimes(i.Path, []unix.Timeval{
		{Sec: i.OriginalAccessTime.Unix(), Usec: int32(i.OriginalAccessTime.Nanosecond() / 1e6 % 1e6)},
		{Sec: i.OriginalModTime.Unix(), Usec: int32(i.OriginalModTime.Nanosecond() / 1e6 % 1e6)}},
	)
}

// CanMaintainSymlinkTimestamps determines whether is is possible to change
// timestamps on symlinks for the the current platform. For regular files
// and directories, attempts are made to restore permissions and timestamps
// after extraction. But for symbolic links, go's cross-platform
// packages (Chmod and Chtimes) are not capable of changing symlink info
// because those methods follow the symlinks. However, a platform-dependent option
// is provided for linux and darwin (see Lchtimes)
func CanMaintainSymlinkTimestamps() bool {
	return true
}
