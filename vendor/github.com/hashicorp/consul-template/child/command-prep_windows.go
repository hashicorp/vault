// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build windows
// +build windows

package child

import (
	"os/exec"
	"strings"
)

// Simplified command prep for windows. Only supports single or slice commands
// and doesn't wrap anything in a shell call (so bool is always false).
func CommandPrep(command []string) ([]string, bool, error) {
	switch {
	case len(command) == 1 && len(strings.Fields(command[0])) == 1:
		// command is []string{"foo"}
		return []string{command[0]}, false, nil
	case len(command) > 1:
		// command is []string{"foo", "bar"}
		return command, false, nil
	default:
		// command is []string{}, []string{""}, []string{"foo bar"}
		return []string{}, false, exec.ErrNotFound
	}
}
