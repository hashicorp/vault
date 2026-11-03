// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package golang

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
)

// maxTidyPasses is the most passes RunGoModTidyUntilStable makes over the
// modules before giving up on them settling.
const maxTidyPasses = 5

// RunGoModTidy executes `go mod tidy` in the directory containing path.
// The caller is responsible for GOPRIVATE / GONOSUMDB / GOFLAGS; this function
// passes the caller's environment through unchanged.
func RunGoModTidy(ctx context.Context, path string) error {
	slog.Default().DebugContext(ctx, "checking module file", slog.String("path", path))
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("reading %s: %w", path, err)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("module file %s is not a regular file", path)
	}

	cmd := exec.CommandContext(ctx, "go", "mod", "tidy")
	cmd.Dir = filepath.Dir(path)
	cmd.Env = append(os.Environ(), "GOWORK=off")
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

// RunGoModTidyUntilStable runs `go mod tidy` for each of paths, in order, and
// repeats until a pass changes none of the modules' go.mod or go.sum files.
//
// Tidying a module reads the go.mod of the local modules it replaces, so one
// pass in any fixed order isn't enough: a module that's tidied before a module
// it depends on is tidied against that module's untidy go.mod and can need
// another tidy once that module is done. A single module is tidied once, as
// that's all it needs.
//
// A module that fails to tidy is left out of the remaining passes. It returns
// the error of each module that failed, by path, and an error when the modules
// are still changing after the most passes we'll make. The caller's environment
// is passed through as it is with RunGoModTidy.
func RunGoModTidyUntilStable(ctx context.Context, paths []string) (map[string]error, error) {
	return tidyUntilStable(ctx, paths, RunGoModTidy)
}

// tidyUntilStable is RunGoModTidyUntilStable with the tidy of a module
// injected, which lets tests control what each pass changes.
func tidyUntilStable(
	ctx context.Context,
	paths []string,
	tidy func(ctx context.Context, path string) error,
) (map[string]error, error) {
	failed := map[string]error{}
	active := append([]string{}, paths...)

	for pass := 1; pass <= maxTidyPasses; pass++ {
		before, err := readModuleFiles(active)
		if err != nil {
			return failed, err
		}

		var tidied []string
		for _, p := range active {
			err := tidy(ctx, p)
			if err != nil {
				failed[p] = err
				continue
			}
			tidied = append(tidied, p)
		}
		active = tidied

		// One module is stable after one tidy, as no other module can change what
		// it needs, and with none left there's nothing more to tidy.
		if len(active) <= 1 {
			return failed, nil
		}

		after, err := readModuleFiles(active)
		if err != nil {
			return failed, err
		}

		changed := false
		for _, p := range active {
			if !bytes.Equal(before[p], after[p]) {
				changed = true
				slog.Default().DebugContext(ctx, "go mod tidy changed a module",
					slog.Int("pass", pass), slog.String("path", p),
				)
			}
		}
		if !changed {
			return failed, nil
		}
	}

	return failed, fmt.Errorf("go mod tidy was still changing modules after %d passes", maxTidyPasses)
}

// readModuleFiles returns the go.mod and go.sum contents of each module, by
// path, as a single value to compare. A missing go.sum reads as empty, because a
// module with no dependencies has no go.sum; go.mod must exist.
func readModuleFiles(paths []string) (map[string][]byte, error) {
	files := make(map[string][]byte, len(paths))
	for _, p := range paths {
		var contents []byte
		for _, name := range []string{p, filepath.Join(filepath.Dir(p), "go.sum")} {
			data, err := os.ReadFile(name)
			if err != nil && (name == p || !errors.Is(err, fs.ErrNotExist)) {
				return nil, fmt.Errorf("reading %s: %w", name, err)
			}
			// Delimit the files so that a move of text from one to the other
			// is a change.
			contents = append(contents, []byte(name+"\n")...)
			contents = append(contents, data...)
			contents = append(contents, 0)
		}
		files[p] = contents
	}

	return files, nil
}
