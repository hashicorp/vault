// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

//go:build !windows
// +build !windows

package system // import "github.com/ory/dockertest/v3/docker/pkg/system"

import (
	"golang.org/x/sys/unix"
)

// Umask sets current process's file mode creation mask to newmask
// and returns oldmask.
func Umask(newmask int) (oldmask int, err error) {
	return unix.Umask(newmask), nil
}
