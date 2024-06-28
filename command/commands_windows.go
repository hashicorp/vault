// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build windows

package command

// MakeSigUSR2Ch does nothing useful on Windows.
func MakeSigUSR2Ch() chan struct{} {
	return make(chan struct{})
}
