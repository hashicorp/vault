// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package system // import "github.com/ory/dockertest/v3/docker/pkg/system"

import "os"

// IsProcessAlive returns true if process with a given pid is running.
func IsProcessAlive(pid int) bool {
	_, err := os.FindProcess(pid)

	return err == nil
}

// KillProcess force-stops a process.
func KillProcess(pid int) {
	p, err := os.FindProcess(pid)
	if err == nil {
		p.Kill()
	}
}
