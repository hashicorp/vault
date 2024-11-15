// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

//go:build !linux
// +build !linux

package system // import "github.com/ory/dockertest/v3/docker/pkg/system"

// Lgetxattr is not supported on platforms other than linux.
func Lgetxattr(path string, attr string) ([]byte, error) {
	return nil, ErrNotSupportedPlatform
}

// Lsetxattr is not supported on platforms other than linux.
func Lsetxattr(path string, attr string, data []byte, flags int) error {
	return ErrNotSupportedPlatform
}
