// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build !darwin && !linux_amd64 && !linux_arm64 && !linux_amd && !linux_arm
// +build !darwin,!linux_amd64,!linux_arm64,!linux_amd,!linux_arm

package unpackinfo

import (
	"errors"
)

// Lchtimes modifies the access and modified timestamps on a target path
// This capability is only available on Linux and Darwin as of now.
func (i UnpackInfo) Lchtimes() error {
	return errors.New("Lchtimes is not supported on this platform")
}

// CanMaintainSymlinkTimestamps determines whether is is possible to change
// timestamps on symlinks for the the current platform. For regular files
// and directories, attempts are made to restore permissions and timestamps
// after extraction. But for symbolic links, go's cross-platform
// packages (Chmod and Chtimes) are not capable of changing symlink info
// because those methods follow the symlinks. However, a platform-dependent option
// is provided for linux and darwin (see Lchtimes)
func CanMaintainSymlinkTimestamps() bool {
	return false
}
