// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tokenhelper

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

var _ TokenHelper = new(ExternalTokenHelper)

// ExternalTokenHelper should only be used in a dev mode. For all other cases,
// InternalTokenHelper should be used.
// ExternalTokenHelper is the struct that has all the logic for storing and retrieving
// tokens from the token helper. The API for the helpers is simple: the
// BinaryPath is executed directly with arguments Args and environment Env.
// The last argument appended to Args will be the operation, which is:
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
	Args       []string
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
	binPath := strings.ReplaceAll(h.BinaryPath, "\\", "\\\\")

	args := make([]string, len(h.Args))
	copy(args, h.Args)
	args = append(args, op)

	cmd := exec.Command(binPath, args...)
	cmd.Env = h.Env
	return cmd, nil
}
