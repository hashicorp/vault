// Copyright IBM Corp. 2026
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
	"path/filepath"
)

// GoWorkSumSyncResult is the result of bringing the workspace in line with the
// synced modules with `go work sync` and `go mod download`, which record the
// checksums the workspace needs in go.work.sum.
type GoWorkSumSyncResult struct {
	// Updated is true when go.work.sum changed.
	Updated bool `json:"updated"`
	// Error is why the workspace couldn't be synced, including which go command
	// failed.
	Error string `json:"error,omitempty"`
}

// Status describes the result of syncing the workspace in one word.
func (r *GoWorkSumSyncResult) Status() string {
	switch {
	case r.Error != "":
		return "failed"
	case r.Updated:
		return "updated"
	default:
		return "unchanged"
	}
}

// WorkspaceGoModPaths returns the go.mod path of every module in the go.work in
// dir, in go.work order, relative to dir. A dir without a go.work has no
// workspace modules, which returns nil.
func WorkspaceGoModPaths(dir string) ([]string, error) {
	goWork := filepath.Join(dir, "go.work")
	data, err := os.ReadFile(goWork)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", goWork, err)
	}

	paths, err := GoWorkGoModPaths(goWork, data)
	if err != nil {
		return nil, fmt.Errorf("listing workspace modules: %w", err)
	}

	return paths, nil
}

// SyncWorkspace brings the workspace in dir, the root of the repository, in
// line with the modules we synced and tidied, in the same way the user would
// after changing requirements by hand:
//
//  1. go work sync raises each module's requirements to the versions the
//     workspace selects, so that building a module on its own with GOWORK=off
//     doesn't select lower versions than the workspace does.
//  2. go mod download all records the checksums that no module's go.sum holds in
//     go.work.sum.
//  3. go mod download in each of modules, with GOWORK=off, records the checksums
//     of the versions step 1 raised the requirements to in that module's go.sum.
//
// modules are go.mod paths relative to dir. It returns the files to stage,
// as paths that start with dir, that aren't a module's go.mod or go.sum:
// go.work.sum when it changed. A directory without a go.work has nothing to
// sync, which returns a nil result.
func SyncWorkspace(ctx context.Context, dir string, modules []string) (*GoWorkSumSyncResult, []string, error) {
	goWork := filepath.Join(dir, "go.work")

	_, err := os.Stat(goWork)
	if errors.Is(err, fs.ErrNotExist) {
		slog.Default().DebugContext(ctx, "skipping workspace sync, no go.work found", slog.String("dir", dir))
		return nil, nil, nil
	}
	if err != nil {
		return nil, nil, fmt.Errorf("checking for %s: %w", goWork, err)
	}

	// A missing go.work.sum is fine; go creates it when it has something to
	// record.
	goWorkSum := goWork + ".sum"
	before, err := os.ReadFile(goWorkSum)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, nil, fmt.Errorf("reading go.work.sum: %w", err)
	}

	res := &GoWorkSumSyncResult{}
	fail := func(err error) (*GoWorkSumSyncResult, []string, error) {
		res.Error = err.Error()
		return res, nil, fmt.Errorf("syncing the workspace: %w", err)
	}

	err = RunGoWorkSync(ctx, goWork)
	if err != nil {
		return fail(err)
	}

	err = RunGoModDownloadAll(ctx, goWork)
	if err != nil {
		return fail(err)
	}

	var downloadErrs []error
	for _, m := range modules {
		err = RunGoModDownload(ctx, filepath.Join(dir, m))
		if err != nil {
			downloadErrs = append(downloadErrs, err)
		}
	}
	if len(downloadErrs) > 0 {
		return fail(errors.Join(downloadErrs...))
	}

	after, err := os.ReadFile(goWorkSum)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return res, nil, fmt.Errorf("reading go.work.sum: %w", err)
	}

	if bytes.Equal(before, after) {
		slog.Default().DebugContext(ctx, "go.work.sum is up to date")
		return res, nil, nil
	}

	res.Updated = true

	return res, []string{goWorkSum}, nil
}
