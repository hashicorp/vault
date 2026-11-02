// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package golang

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"
)

// RunGoWorkSync executes `go work sync` in the workspace described by goWork.
// It raises the requirements in each workspace module's go.mod to the versions
// the workspace build list selects, so that a module that's built on its own
// with GOWORK=off doesn't select lower versions than the workspace does.
//
// We set GOWORK to goWork explicitly, so the go command uses that workspace
// even when the caller's environment sets GOWORK=off. As with RunGoModTidy, the
// caller is responsible for GOPRIVATE / GONOSUMDB / GOFLAGS; they're passed
// through unchanged.
func RunGoWorkSync(ctx context.Context, goWork string) error {
	if goWork == "" {
		return errors.New("go.work path is required")
	}

	goWork, err := filepath.Abs(goWork)
	if err != nil {
		return fmt.Errorf("resolving go.work path: %w", err)
	}

	slog.Default().DebugContext(ctx, "running go work sync", slog.String("go_work", goWork))
	res, err := runGo(ctx, filepath.Dir(goWork), []string{"GOWORK=" + goWork}, "work", "sync")
	if err != nil {
		slog.Default().ErrorContext(ctx, "go work sync failed",
			slog.String("go_work", goWork),
			slog.String("output", string(res.Stderr)),
		)

		// runGo's error names the command and directory but not why it failed.
		return fmt.Errorf("%w: %s", err, strings.TrimSpace(string(res.Stderr)))
	}

	return nil
}
