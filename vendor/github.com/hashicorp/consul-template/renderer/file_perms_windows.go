// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build windows
// +build windows

package renderer

import "os"

// Not done as Windows doedsn't realiably support permissions setting.
// https://github.com/google/renameio/issues/17
func preserveFilePermissions(path string, fileInfo os.FileInfo) error {
	return nil
}
