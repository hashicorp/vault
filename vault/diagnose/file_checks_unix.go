// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build !windows

package diagnose

import (
	"io/fs"

	"github.com/hashicorp/vault/helper/osutil"
)

// IsOwnedByRoot checks if a file is owned by root
func IsOwnedByRoot(info fs.FileInfo) bool {
	if !osutil.FileUIDEqual(info, 0) {
		return false
	}
	if !osutil.FileGIDEqual(info, 0) {
		return false
	}
	return true
}
