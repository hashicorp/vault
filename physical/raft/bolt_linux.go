// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package raft

import (
	"context"
	"os"

	"github.com/shirou/gopsutil/v3/mem"
	"golang.org/x/sys/unix"
)

func init() {
	getMmapFlags = getMmapFlagsLinux
	usingMapPopulate = usingMapPopulateLinux
}

func getMmapFlagsLinux(dbPath string) int {
	if setMapPopulateFlag(dbPath) {
		return unix.MAP_POPULATE
	}

	return 0
}

// setMapPopulateFlag determines whether we should set the MAP_POPULATE flag, which
// prepopulates page tables to be mapped in the virtual memory space,
// helping reduce slowness at runtime caused by page faults.
// We only want to set this flag if we've determined there's enough memory on the system available to do so.
func setMapPopulateFlag(dbPath string) bool {
	if os.Getenv("VAULT_RAFT_DISABLE_MAP_POPULATE") != "" {
		return false
	}
	stat, err := os.Stat(dbPath)
	if err != nil {
		return false
	}
	size := stat.Size()

	v, err := mem.VirtualMemoryWithContext(context.Background())
	if err != nil {
		return false
	}

	// We won't worry about swap, since we already tell people not to use it.
	if v.Total > uint64(size) {
		return true
	}

	return false
}

// the unix.MAP_POPULATE constant only exists on Linux,
// so reference to this constant can only live in a *_linux.go file
func usingMapPopulateLinux(mmapFlag int) bool {
	if mmapFlag == unix.MAP_POPULATE {
		return true
	}
	return false
}
