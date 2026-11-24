// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package git

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	libgit "github.com/hashicorp/vault/tools/pipeline/internal/pkg/git/client"
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/golang"
)

// SyncGoModFromRef syncs the directives of the go.mod at p in ref, a branch or
// other revision that git can resolve, into the go.mod at p in dir in the
// working tree, and writes the result there. It doesn't tidy the module.
//
// p is the path of the go.mod, which git resolves from the root of the
// repository, so we only use it to read from ref and from dir if it's
// relative to the root of the repository. dir is the root of the repository for
// a path that comes from a go.work, or "." for a path from the user, which is
// relative to the working directory like the files we read and write.
//
// It's best effort in that it returns a result when the sync fails. The result
// has the error as its Error.
func SyncGoModFromRef(
	ctx context.Context,
	git *libgit.Client,
	ref, dir, p string,
	diffOpts *golang.DiffOpts,
) (*golang.GoModPathSyncResult, error) {
	object := ref + ":" + p
	slog.Default().DebugContext(ctx, "reading source go.mod from ref", slog.String("object", object))
	showRes, err := git.Show(ctx, &libgit.ShowOpts{Object: object})
	if err != nil {
		syncRes := &golang.SyncModRes{
			Path:  p,
			Error: fmt.Sprintf("git show %s: %v", object, err),
		}
		syncRes.Err = errors.New(syncRes.Error)

		return &golang.GoModPathSyncResult{SyncModRes: syncRes}, syncRes.Err
	}
	sourceBytes := showRes.Stdout

	destPath := filepath.Join(dir, p)
	if err := golang.ValidateGoModPath(ctx, dir, p); err != nil {
		syncRes := &golang.SyncModRes{Path: p, Error: err.Error(), Err: err}
		return &golang.GoModPathSyncResult{SyncModRes: syncRes}, err
	}
	slog.Default().DebugContext(ctx, "reading dest go.mod from working tree", slog.String("path", destPath))
	destBytes, err := os.ReadFile(destPath)
	if err != nil {
		syncRes := &golang.SyncModRes{
			Path:  p,
			Error: fmt.Sprintf("reading %s: %v", destPath, err),
		}
		syncRes.Err = errors.New(syncRes.Error)

		return &golang.GoModPathSyncResult{SyncModRes: syncRes}, syncRes.Err
	}

	// B's name is the file the sync writes to, so it has to be the path in dir.
	syncReq := &golang.SyncModReq{
		A:        &golang.ModSource{Name: p, Data: sourceBytes},
		B:        &golang.ModSource{Name: destPath, Data: destBytes},
		DiffOpts: diffOpts,
	}
	syncRes, err := syncReq.Run(ctx)

	// The result takes its path from B's name. We want the path we were given
	// instead, which is how the results are matched to the modules of a go.work
	// and shown, and doesn't depend on where dir is.
	syncRes.Path = p

	return &golang.GoModPathSyncResult{SyncModRes: syncRes}, err
}

// WorkspaceGoModPathsFromRef returns the go.mod path of every module in the
// go.work at the root of ref, a branch or other revision that git can resolve,
// in go.work order. The paths are relative to the root of the repository. It's
// an error for ref to have no go.work, or a go.work without modules, so that we
// never fall back to syncing only the root go.mod while looking like a success.
func WorkspaceGoModPathsFromRef(
	ctx context.Context,
	git *libgit.Client,
	ref string,
) ([]string, error) {
	object := ref + ":go.work"
	slog.Default().DebugContext(ctx, "reading go.work from ref", slog.String("object", object))
	showRes, err := git.Show(ctx, &libgit.ShowOpts{Object: object})
	if err != nil {
		return nil, fmt.Errorf(
			"listing workspace modules: git show %s: %w: %s",
			object, err, strings.TrimSpace(string(showRes.Stderr)),
		)
	}

	paths, err := golang.GoWorkGoModPaths(object, showRes.Stdout)
	if err != nil {
		return nil, fmt.Errorf("listing workspace modules: %w", err)
	}
	if len(paths) == 0 {
		return nil, fmt.Errorf("listing workspace modules: %s has no modules", object)
	}

	return paths, nil
}
