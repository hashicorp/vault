// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build !windows
// +build !windows

package child

import (
	"os/exec"
	"strings"
)

// Evaluates the command slice for the different possible formats.
// Returns the command slice ready to pass to exec.Command.
// Returns a boolean 'true' if it wrapped the call in 'sh -c' so the caller
// knows it needs to setpgid to get signal propagation.
func CommandPrep(command []string) ([]string, bool, error) {
	switch {
	case len(command) == 1 && len(strings.Fields(command[0])) > 1:
		// command is []string{"command using arguments or shell features"}
		shell := "sh"
		// default to 'sh' on path, else try a couple common absolute paths
		if _, err := exec.LookPath(shell); err != nil {
			shell = ""
			for _, sh := range []string{"/bin/sh", "/usr/bin/sh"} {
				if absPath, err := exec.LookPath(sh); err == nil {
					shell = absPath
					break
				}
			}
		}
		if shell == "" {
			return []string{}, false, exec.ErrNotFound
		}
		cmd := []string{shell, "-c", command[0]}
		return cmd, true, nil
	case len(command) >= 1 && len(strings.TrimSpace(command[0])) > 0:
		// command is already good ([]string{"foo"}, []string{"foo", "bar"}, ..)
		return command, false, nil
	default:
		// command is []string{} or []string{""}
		return []string{}, false, exec.ErrNotFound
	}
}
