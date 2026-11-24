// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package golang

import (
	"context"
	"errors"
	"os"
	"path/filepath"
)

// GoModPathSyncResult wraps SyncModRes with the result of tidying the module.
type GoModPathSyncResult struct {
	*SyncModRes
	TidyRan   bool   `json:"tidy_ran"`
	TidyError string `json:"tidy_error,omitempty"`
	TidyErr   error  `json:"-"`
}

// TidySyncedModulesRes is the result of TidySyncedModules.
type TidySyncedModulesRes struct {
	// Results are the results we were given, plus one for each module we tidied
	// that we hadn't synced, as with a module that's in the workspace of the
	// branch we synced into but not the one we synced from.
	Results []*GoModPathSyncResult
	// GoWorkSum is set when we tried to sync the workspace.
	GoWorkSum *GoWorkSumSyncResult
	// Stage are the files that changed and that exist, as paths that start with
	// the directory we ran in. A module that has no dependencies has no go.sum,
	// and git add fails on a path that doesn't exist.
	Stage []string
}

// TidySyncedModules tidies the modules after their directives have been synced
// and returns the files to stage. It's the part of a go.mod sync that follows
// syncing the directives, whichever branches or repositories we synced between.
//
// dir is the root of the repository, the directory that the paths of results,
// and the go.work when workspace is set, are relative to.
//
// We sync the directives of every module before we tidy any of them. Tidying a
// module reads the go.mod of the local modules it replaces, so a module that we
// tidy before the modules it depends on are synced is tidied against stale
// requirements and needs another pass. We tidy the modules until they settle
// and stage the ones that tidied.
//
// When workspace is set we also tidy every module in the go.work in dir, without
// raising their requirements to workspace-selected versions. A dir without a
// go.work has no workspace modules.
//
// Like the sync, it's best effort. It always returns the result, and an error
// that has every problem we ran into.
func TidySyncedModules(
	ctx context.Context,
	dir string,
	results []*GoModPathSyncResult,
	workspace bool,
) (*TidySyncedModulesRes, error) {
	res := &TidySyncedModulesRes{Results: results}

	var errs []error
	var stage []string

	var workspacePaths []string
	if workspace {
		var err error
		workspacePaths, err = WorkspaceGoModPaths(dir)
		if err != nil {
			errs = append(errs, err)
			res.GoWorkSum = &GoWorkSumSyncResult{Error: err.Error()}
		}
	}

	tidyList := modulesToTidy(results, workspacePaths)
	var validList, tidyPaths []string
	for _, p := range tidyList {
		if err := ValidateGoModPath(ctx, dir, p); err != nil {
			pr, updated := resultFor(res.Results, p)
			res.Results = updated
			pr.TidyErr = err
			pr.TidyError = err.Error()
			errs = append(errs, err)
			continue
		}
		validList = append(validList, p)
		tidyPaths = append(tidyPaths, filepath.Join(dir, p))
	}

	ran := map[string]bool{}
	failedTidy, err := tidyUntilStable(ctx, tidyPaths, func(ctx context.Context, p string) error {
		ran[p] = true
		return RunGoModTidy(ctx, p)
	})
	if err != nil {
		errs = append(errs, err)
	}

	for i, p := range validList {
		var pr *GoModPathSyncResult
		pr, res.Results = resultFor(res.Results, p)
		pr.TidyRan = ran[tidyPaths[i]]

		tidyErr := failedTidy[tidyPaths[i]]
		if tidyErr != nil {
			errs = append(errs, tidyErr)
		} else {
			tidyErr = err
		}
		if tidyErr != nil {
			pr.TidyErr = tidyErr
			pr.TidyError = tidyErr.Error()
			continue
		}

		modDir := filepath.Dir(tidyPaths[i])
		stage = append(stage, filepath.Join(modDir, "go.mod"), filepath.Join(modDir, "go.sum"))
	}

	res.Stage = existingFiles(stage)

	return res, errors.Join(errs...)
}

// modulesToTidy returns the go.mod paths to tidy, in the order to tidy them,
// given the results of syncing the modules. Nothing is tidied when the sync
// changed nothing, so that we leave the working tree alone. Otherwise we tidy
// every module that changed and, when workspacePaths is set, every module in the
// workspace, in go.work order, because tidying a module whose own go.mod didn't
// change can still be needed when a module it replaces did. A module that
// failed to sync isn't tidied.
func modulesToTidy(results []*GoModPathSyncResult, workspacePaths []string) []string {
	changed := map[string]struct{}{}
	failed := map[string]struct{}{}
	var changedOrder []string
	for _, pr := range results {
		if pr == nil || pr.SyncModRes == nil {
			continue
		}
		if pr.Error != "" {
			failed[pr.Path] = struct{}{}
			continue
		}
		if len(pr.Changes) > 0 {
			changed[pr.Path] = struct{}{}
			changedOrder = append(changedOrder, pr.Path)
		}
	}

	if len(changed) == 0 {
		return nil
	}

	var paths []string
	seen := map[string]struct{}{}
	for _, p := range append(append([]string{}, workspacePaths...), changedOrder...) {
		if _, ok := seen[p]; ok {
			continue
		}
		seen[p] = struct{}{}

		if _, ok := failed[p]; ok {
			continue
		}
		paths = append(paths, p)
	}

	return paths
}

// resultFor returns the result for the go.mod path p from results, adding one
// when we haven't synced p. It returns results, which has the new result when we
// added one.
func resultFor(results []*GoModPathSyncResult, p string) (*GoModPathSyncResult, []*GoModPathSyncResult) {
	for _, pr := range results {
		if pr != nil && pr.SyncModRes != nil && pr.Path == p {
			return pr, results
		}
	}

	pr := &GoModPathSyncResult{SyncModRes: &SyncModRes{Path: p}}

	return pr, append(results, pr)
}

// existingFiles returns the files in paths that exist.
func existingFiles(paths []string) []string {
	var existing []string
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			existing = append(existing, p)
		}
	}

	return existing
}
