// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build windows

package diagnose

import "io/fs"

// IsOwnedByRoot does nothing on windows for now.
// TODO: find an equivalent check for file ownership in windows.
func IsOwnedByRoot(info fs.FileInfo) bool {
	return false
}
