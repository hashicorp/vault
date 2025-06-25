// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"

	libgit "github.com/hashicorp/vault/tools/pipeline/internal/pkg/git"
	slogctx "github.com/veqryn/slog-context"
)

// ensureGitRepoDir repoDir verifies that the `dir` exists and is a directory.
// If the `dir` is unset, a temporary directory will be created. A boolean
// is returned which can be used to determine whether or not the path returned
// is a temporary directory.
func ensureGitRepoDir(ctx context.Context, dir string) (string, error, bool) {
	if dir == "" {
		slog.Default().DebugContext(ctx, "creating repository directory")
		dir, err := os.MkdirTemp("", "pipeline-cmd-git-rep-dir")
		return dir, err, true
	}

	ctx = slogctx.Append(ctx, slog.String("repo-dir", dir))
	slog.Default().DebugContext(ctx, "verifying repository directory")
	info, err := os.Stat(dir)
	if err != nil {
		return dir, fmt.Errorf("stating repository directory: %w", err), false
	}

	if !info.IsDir() {
		return dir, errors.New("repo dir must be a directory"), false
	}

	return dir, nil, false
}

// initializeExistingRepo initializes an existing repository. It assumes that
// at least one remote origin exists and that some branch is checked out. If
// the current branch is our baseRef we'll pull in the latest changes, otherwise
// we'll fetch baseRef.
func initializeExistingRepo(
	ctx context.Context,
	git *libgit.Client,
	repoDir string,
	baseOrigin string,
	baseRef string,
) error {
	ctx = slogctx.Append(ctx,
		slog.String("repo-dir", repoDir),
		slog.String("base-origin", baseOrigin),
		slog.String("base-ref", baseRef),
	)
	// We've been given an already initialized git directory. We'll have to
	// assume it's the correct repo that has been cloned.
	slog.Default().WarnContext(ctx, "using an already initialized git repository")

	slog.Default().DebugContext(ctx, "changing working directory to repository dir")
	err := os.Chdir(repoDir)
	if err != nil {
		return fmt.Errorf("changing directory to the repository dir: %w", err)
	}

	// Determine if we're on the correct branch. If we are, pull it, otherwise
	// fetch it.
	slog.Default().DebugContext(ctx, "getting existing repo current branch")
	res, err := git.Branch(ctx, &libgit.BranchOpts{
		NoColor:     true,
		ShowCurrent: true,
	})
	if err != nil {
		return fmt.Errorf("getting existing repo current branch: %w", err)
	}

	if strings.TrimSpace(string(res.Stdout)) == baseRef {
		// Our existing repo is already checked out to correct branch.
		// Fetch the base ref to make sure our existing repository has the necessary
		// objects and references we'll need to cherry-pick the commits to our
		// branch.
		slog.Default().DebugContext(ctx, "pulling in latest changes")
		res, err := git.Pull(ctx, &libgit.PullOpts{
			Refspec:     []string{baseOrigin, baseRef},
			Autostash:   true,
			SetUpstream: true,
			Rebase:      libgit.RebaseStrategyTrue,
		})
		if err != nil {
			return fmt.Errorf("pulling repo current branch: %s: %w", res.String(), err)
		}

		return nil
	}

	// Fetch the base ref to make sure our existing repository has the necessary
	// objects and references.
	slog.Default().DebugContext(ctx, "fetching repository base ref")
	res, err = git.Fetch(ctx, &libgit.FetchOpts{
		// Fetch the ref but also provide a local tracking branch of the same name
		// e.g. "git fetch origin main:main"
		Refspec:     []string{baseOrigin, fmt.Sprintf("%s:%s", baseRef, baseRef)},
		SetUpstream: true,
		Porcelain:   true,
	})
	if err != nil {
		return fmt.Errorf("fetching base ref: %s: %w", res.String(), err)
	}

	return nil
}

// initializeNewRepo initializes a new repository by cloning the repo fetching
// the `baseRef`.
func initializeNewRepo(
	ctx context.Context,
	git *libgit.Client,
	repoDir string,
	owner string,
	repo string,
	baseOrigin string,
	baseRef string,
) error {
	cloneURL := fmt.Sprintf("https://github.com/%s/%s.git", owner, repo)
	slog.Default().DebugContext(slogctx.Append(ctx,
		slog.String("owner", owner),
		slog.String("repo", repo),
		slog.String("base-origin", baseOrigin),
		slog.String("base-ref", baseRef),
		slog.String("repo-dir", repoDir),
		slog.String("repo-url", cloneURL),
	), "initializing new clone of repository")

	res, err := git.Clone(ctx, &libgit.CloneOpts{
		Repository:   cloneURL,
		Directory:    repoDir,
		Origin:       baseOrigin,
		Branch:       baseRef,
		SingleBranch: true,
		NoCheckout:   true,
	})
	if err != nil {
		return fmt.Errorf("cloning repository: %s: %w", res.String(), err)
	}

	slog.Default().DebugContext(ctx, "changing working directory to repo-dir")
	err = os.Chdir(repoDir)
	if err != nil {
		return fmt.Errorf("changing directory to the repository dir: %w", err)
	}

	return nil
}
