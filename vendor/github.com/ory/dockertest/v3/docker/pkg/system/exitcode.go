// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package system // import "github.com/ory/dockertest/v3/docker/pkg/system"

import (
	"fmt"
	"os/exec"
	"syscall"
)

// GetExitCode returns the ExitStatus of the specified error if its type is
// exec.ExitError, returns 0 and an error otherwise.
func GetExitCode(err error) (int, error) {
	exitCode := 0
	if exiterr, ok := err.(*exec.ExitError); ok {
		if procExit, ok := exiterr.Sys().(syscall.WaitStatus); ok {
			return procExit.ExitStatus(), nil
		}
	}
	return exitCode, fmt.Errorf("failed to get exit code")
}
