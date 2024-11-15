// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

//go:build windows
// +build windows

package system // import "github.com/ory/dockertest/v3/docker/pkg/system"

// Mknod is not implemented on Windows.
func Mknod(path string, mode uint32, dev int) error {
	return ErrNotSupportedPlatform
}

// Mkdev is not implemented on Windows.
func Mkdev(major int64, minor int64) uint32 {
	panic("Mkdev not implemented on Windows.")
}
