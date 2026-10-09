// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package golang

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
)

// GoExecResponse is the result of running a go command.
type GoExecResponse struct {
	Cmd    string
	Dir    string
	Stdout []byte
	Stderr []byte
}

// runGo runs the go command in dir with GOWORK=off, so the module in dir is
// built with its own go.mod and go.sum rather than the workspace's. env is
// added after the process environment and GOWORK; later entries win.
func runGo(ctx context.Context, dir string, env []string, args ...string) (*GoExecResponse, error) {
	cmd := exec.CommandContext(ctx, "go", args...)
	cmd.Dir = dir
	cmd.Env = append(append(os.Environ(), "GOWORK=off"), env...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	res := &GoExecResponse{Cmd: cmd.String(), Dir: dir}
	// Only log the extra environment. The process environment can hold secrets.
	slog.Default().DebugContext(ctx, "executing go command",
		slog.String("cmd", res.Cmd),
		slog.String("dir", dir),
		slog.Any("env", env),
	)

	err := cmd.Run()
	res.Stdout = stdout.Bytes()
	res.Stderr = stderr.Bytes()
	if err != nil {
		return res, fmt.Errorf("running %s in %s: %w", res.Cmd, dir, err)
	}

	return res, nil
}
