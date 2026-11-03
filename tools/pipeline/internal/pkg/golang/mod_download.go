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

// RunGoModDownloadAll executes `go mod download all` in the workspace described
// by goWork. In workspace mode the go command records the checksums that no
// module's go.sum holds in go.work.sum, so this brings go.work.sum up to date
// after the workspace modules' requirements change.
//
// We set GOWORK to goWork explicitly, so the go command uses that workspace
// even when the caller's environment sets GOWORK=off. As with RunGoModTidy, the
// caller is responsible for GOPRIVATE / GONOSUMDB / GOFLAGS; they're passed
// through unchanged.
func RunGoModDownloadAll(ctx context.Context, goWork string) error {
	if goWork == "" {
		return errors.New("go.work path is required")
	}

	goWork, err := filepath.Abs(goWork)
	if err != nil {
		return fmt.Errorf("resolving go.work path: %w", err)
	}

	slog.Default().DebugContext(ctx, "running go mod download all",
		slog.String("go_work", goWork),
	)
	res, err := runGo(ctx, filepath.Dir(goWork), []string{"GOWORK=" + goWork}, "mod", "download", "all")
	if err != nil {
		slog.Default().ErrorContext(ctx, "go mod download all failed",
			slog.String("go_work", goWork),
			slog.String("output", string(res.Stderr)),
		)

		// runGo's error names the command and directory but not why it failed.
		return fmt.Errorf("%w: %s", err, strings.TrimSpace(string(res.Stderr)))
	}

	return nil
}

// RunGoModDownload executes `go mod download` for the module at path, with
// GOWORK=off so the module is resolved from its own go.mod and go.sum rather
// than the workspace's. The go command records the checksums of every module in
// the module's build list in its go.sum, including the go.mod checksums of
// versions a workspace sync raised the requirements to. As with RunGoModTidy,
// the caller is responsible for GOPRIVATE / GONOSUMDB / GOFLAGS; they're passed
// through unchanged.
func RunGoModDownload(ctx context.Context, path string) error {
	if path == "" {
		return errors.New("go.mod path is required")
	}

	dir := filepath.Dir(path)
	slog.Default().DebugContext(ctx, "running go mod download", slog.String("dir", dir))
	res, err := runGo(ctx, dir, nil, "mod", "download")
	if err != nil {
		slog.Default().ErrorContext(ctx, "go mod download failed",
			slog.String("dir", dir),
			slog.String("output", string(res.Stderr)),
		)

		// runGo's error names the command and directory but not why it failed.
		return fmt.Errorf("%w: %s", err, strings.TrimSpace(string(res.Stderr)))
	}

	return nil
}
