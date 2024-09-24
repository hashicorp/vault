// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tokenhelper

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// ExternalTokenHelperPath should only be used in dev mode.
// ExternalTokenHelperPath takes the configured path to a helper and expands it to
// a full absolute path that can be executed. As of 0.5, the default token
// helper is internal, to avoid problems running in dev mode (see GH-850 and
// GH-783), so special assumptions of prepending "vault token-" no longer
// apply.
//
// As an additional result, only absolute paths are now allowed. Looking in the
// path or a current directory for an arbitrary executable could allow someone
// to switch the expected binary for one further up the path (or in the current
// directory), potentially opening up execution of an arbitrary binary.
func ExternalTokenHelperPath(path string) (string, error) {
	if !filepath.IsAbs(path) {
		var err error
		path, err = filepath.Abs(path)
		if err != nil {
			return "", err
		}
	}

	if _, err := os.Stat(path); err != nil {
		return "", fmt.Errorf("unknown error getting the external helper path")
	}

	return path, nil
}

var _ TokenHelper = (*ExternalTokenHelper)(nil)

// ExternalTokenHelper should only be used in a dev mode. For all other cases,
// InternalTokenHelper should be used.
// ExternalTokenHelper is the struct that has all the logic for storing and retrieving
// tokens from the token helper. The API for the helpers is simple: the
// BinaryPath is executed within a shell with environment Env. The last argument
// appended will be the operation, which is:
//
//   - "get" - Read the value of the token and write it to stdout.
//   - "store" - Store the value of the token which is on stdin. Output
//     nothing.
//   - "erase" - Erase the contents stored. Output nothing.
//
// Any errors can be written on stdout. If the helper exits with a non-zero
// exit code then the stderr will be made part of the error value.
type ExternalTokenHelper struct {
	BinaryPath string
	Env        []string
}

// Erase deletes the contents from the helper.
func (h *ExternalTokenHelper) Erase() error {
	cmd, err := h.cmd("erase")
	if err != nil {
		return err
	}
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%q: %w", string(output), err)
	}
	return nil
}

// Get gets the token value from the helper.
func (h *ExternalTokenHelper) Get() (string, error) {
	var buf, stderr bytes.Buffer
	cmd, err := h.cmd("get")
	if err != nil {
		return "", err
	}
	cmd.Stdout = &buf
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%q: %w", stderr.String(), err)
	}

	return buf.String(), nil
}

// Store stores the token value into the helper.
func (h *ExternalTokenHelper) Store(v string) error {
	buf := bytes.NewBufferString(v)
	cmd, err := h.cmd("store")
	if err != nil {
		return err
	}
	cmd.Stdin = buf
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%q: %w", string(output), err)
	}

	return nil
}

func (h *ExternalTokenHelper) Path() string {
	return h.BinaryPath
}

func (h *ExternalTokenHelper) cmd(op string) (*exec.Cmd, error) {
	script := strings.ReplaceAll(h.BinaryPath, "\\", "\\\\") + " " + op
	cmd, err := execScript(script)
	if err != nil {
		return nil, err
	}
	cmd.Env = h.Env
	return cmd, nil
}

// execScript returns a command to execute a script
func execScript(script string) (*exec.Cmd, error) {
	var shell, flag string
	if runtime.GOOS == "windows" {
		shell = "cmd"
		flag = "/C"
	} else {
		shell = "/bin/sh"
		flag = "-c"
	}
	if other := os.Getenv("SHELL"); other != "" {
		shell = other
	}
	cmd := exec.Command(shell, flag, script)
	return cmd, nil
}
