// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fileutils // import "github.com/ory/dockertest/v3/docker/pkg/fileutils"

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// GetTotalUsedFds returns the number of used File Descriptors by
// executing `lsof -p PID`
func GetTotalUsedFds() int {
	pid := os.Getpid()

	cmd := exec.Command("lsof", "-p", strconv.Itoa(pid))

	output, err := cmd.CombinedOutput()
	if err != nil {
		return -1
	}

	outputStr := strings.TrimSpace(string(output))

	fds := strings.Split(outputStr, "\n")

	return len(fds) - 1
}
