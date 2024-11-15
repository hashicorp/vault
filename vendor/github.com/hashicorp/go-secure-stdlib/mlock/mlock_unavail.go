// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// +build darwin nacl netbsd plan9 windows

package mlock

func init() {
	supported = false
}

func lockMemory() error {
	// XXX: No good way to do this on Windows. There is the VirtualLock
	// method, but it requires a specific address and offset.
	return nil
}
