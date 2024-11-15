// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

//go:build linux || freebsd || darwin
// +build linux freebsd darwin

package system // import "github.com/ory/dockertest/v3/docker/pkg/system"

import (
	"syscall"

	"golang.org/x/sys/unix"
)

// IsProcessAlive returns true if process with a given pid is running.
func IsProcessAlive(pid int) bool {
	err := unix.Kill(pid, syscall.Signal(0))
	if err == nil || err == unix.EPERM {
		return true
	}

	return false
}

// KillProcess force-stops a process.
func KillProcess(pid int) {
	unix.Kill(pid, unix.SIGKILL)
}
