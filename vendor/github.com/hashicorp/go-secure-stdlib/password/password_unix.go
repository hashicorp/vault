// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build linux || darwin || freebsd || netbsd || openbsd || dragonfly
// +build linux darwin freebsd netbsd openbsd dragonfly

package password

import (
	"fmt"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

func read(f *os.File) (string, error) {
	fd := int(f.Fd())
	if !terminal.IsTerminal(fd) {
		return "", fmt.Errorf("file descriptor %d is not a terminal", fd)
	}

	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		return "", err
	}
	defer terminal.Restore(fd, oldState)

	return readline(f)
}
