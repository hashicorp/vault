// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package golang

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
)

// RunGoModTidy executes `go mod tidy` in the directory containing path.
// The caller is responsible for GOPRIVATE / GONOSUMDB / GOFLAGS; this function
// passes the caller's environment through unchanged.
func RunGoModTidy(ctx context.Context, path string) error {
	cmd := exec.CommandContext(ctx, "go", "mod", "tidy")
	cmd.Dir = filepath.Dir(path)
	cmd.Env = os.Environ()
	slog.Default().DebugContext(ctx, "running go mod tidy",
		slog.String("dir", cmd.Dir),
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		slog.Default().ErrorContext(ctx, "go mod tidy failed",
			slog.String("dir", cmd.Dir),
			slog.String("output", string(out)),
		)
		return fmt.Errorf("go mod tidy in %s: %w", cmd.Dir, err)
	}
	return nil
}
