// +build !windows

package manager

import (
	"os/exec"
	"strings"
)

func prepCommand(command string) ([]string, error) {
	switch len(strings.Fields(command)) {
	case 0:
		return []string{}, nil
	case 1:
		return []string{command}, nil
	}

	// default to 'sh' on path, else try a couple common absolute paths
	shell := "sh"
	if _, err := exec.LookPath(shell); err != nil {
		for _, sh := range []string{"/bin/sh", "/usr/bin/sh"} {
			if sh, err := exec.LookPath(sh); err == nil {
				shell = sh
				break
			}
		}
	}
	if shell == "" {
		return []string{}, exec.ErrNotFound
	}

	cmd := []string{shell, "-c", command}
	return cmd, nil
}
