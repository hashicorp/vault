// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	libgithub "github.com/google/go-github/v81/github"
	libgit "github.com/hashicorp/vault/tools/pipeline/internal/pkg/git"
	"github.com/jedib0t/go-pretty/v6/table"
	slogctx "github.com/veqryn/slog-context"
)

// SyncBranchReq is a request to synchronize two github hosted branches with
// a git merge from one into another.
//
// NOTE: We require that both branches exist for the operation to succeed.
type SyncBranchReq struct {
	FromOwner  string
	FromRepo   string
	FromOrigin string
	FromBranch string
	ToOwner    string
	ToRepo     string
	ToOrigin   string
	ToBranch   string
	RepoDir    string
}

// SyncBranchRes is a copy pull request response.
type SyncBranchRes struct {
	Error   error          `json:"error,omitempty"`
	Request *SyncBranchReq `json:"request,omitempty"`
}

// Run runs the request to synchronize a branches via a merge.
func (r *SyncBranchReq) Run(
	ctx context.Context,
	github *libgithub.Client,
	git *libgit.Client,
) (*SyncBranchRes, error) {
	var err error
	res := &SyncBranchRes{Request: r}

	slog.Default().DebugContext(slogctx.Append(ctx,
		slog.String("from-owner", r.FromOwner),
		slog.String("from-repo", r.FromRepo),
		slog.String("from-origin", r.FromOrigin),
		slog.String("from-branch", r.FromBranch),
		slog.String("to-owner", r.ToOwner),
		slog.String("to-repo", r.ToRepo),
		slog.String("to-origin", r.ToOrigin),
		slog.String("to-branch", r.ToBranch),
		slog.String("repo-dir", r.RepoDir),
	), "synchronizing branches")

	// Make sure we have required and valid fields
	err = r.Validate(ctx)
	if err != nil {
		return res, err
	}

	// Make sure we've been given a valid location for a repo and/or create a
	// temporary one
	var tmpDir bool
	r.RepoDir, err, tmpDir = ensureGitRepoDir(ctx, r.RepoDir)
	if err != nil {
		return res, err
	}
	if tmpDir {
		defer os.RemoveAll(r.RepoDir)
	}

	// Clone the remote repository and fetch the branch we're going to merge into.
	// These will change our working directory into RepoDir.
	_, err = os.Stat(filepath.Join(r.RepoDir, ".git"))
	if err == nil {
		err = initializeExistingRepo(
			ctx, git, r.RepoDir, r.ToOrigin, r.ToBranch,
		)
	} else {
		err = initializeNewRepo(
			ctx, git, r.RepoDir, r.ToOwner, r.ToRepo, r.ToOrigin, r.ToBranch,
		)
	}
	if err != nil {
		return res, err
	}

	// Check out our branch. Our intialization above will ensure we have a local
	// reference.
	slog.Default().DebugContext(ctx, "checking out to-branch")
	checkoutRes, err := git.Checkout(ctx, &libgit.CheckoutOpts{
		Branch: r.ToBranch,
	})
	if err != nil {
		return res, fmt.Errorf("checking out to-branch: %s: %w", checkoutRes.String(), err)
	}

	// Add our from upstream as a remote and fetch our from branch.
	slog.Default().DebugContext(ctx, "adding from upstream and fetching from-branch")
	remoteRes, err := git.Remote(ctx, &libgit.RemoteOpts{
		Command: libgit.RemoteCommandAdd,
		Track:   []string{r.FromBranch},
		Fetch:   true,
		Name:    r.FromOrigin,
		URL:     fmt.Sprintf("https://github.com/%s/%s.git", r.FromOwner, r.FromRepo),
	})
	if err != nil {
		err = fmt.Errorf("fetching from branch: %s, %w", remoteRes.String(), err)
		return res, err
	}

	// Use our remote reference as we haven't created a local reference.
	fromBranch := "remotes/" + r.FromOrigin + "/" + r.FromBranch
	slog.Default().DebugContext(ctx, "merging from-branch into to-branch")
	mergeRes, err := git.Merge(ctx, &libgit.MergeOpts{
		NoVerify: true,
		Strategy: libgit.MergeStrategyORT,
		StrategyOptions: []libgit.MergeStrategyOption{
			libgit.MergeStrategyOptionTheirs,
			libgit.MergeStrategyOptionIgnoreSpaceChange,
		},
		IntoName: r.ToBranch,
		Commit:   fromBranch,
	})
	if err != nil {
		return res, fmt.Errorf("merging from-branch into to-branch: %s: %w", mergeRes.String(), err)
	}

	slog.Default().DebugContext(ctx, "pushing to-branch")
	pushRes, err := git.Push(ctx, &libgit.PushOpts{
		Repository: r.ToOrigin,
		Refspec:    []string{r.ToBranch},
	})
	if err != nil {
		return res, fmt.Errorf("pushing to-branch: %s: %w", pushRes.String(), err)
	}

	return res, nil
}

// Validate ensures that we've been given the minimum filter arguments necessary to complete a
// request. It is always recommended that additional fitlers be given to reduce the response size
// and not exhaust API limits.
func (r *SyncBranchReq) Validate(ctx context.Context) error {
	if r == nil {
		return errors.New("failed to initialize request")
	}

	if r.FromOrigin == "" {
		return errors.New("no github from origin has been provided")
	}

	if r.FromOwner == "" {
		return errors.New("no github from owner has been provided")
	}

	if r.FromRepo == "" {
		return errors.New("no github from repository has been provided")
	}

	if r.FromBranch == "" {
		return errors.New("no github from branch has been provided")
	}

	if r.ToOrigin == "" {
		return errors.New("no github to origin has been provided")
	}

	if r.ToOwner == "" {
		return errors.New("no github to owner has been provided")
	}

	if r.ToRepo == "" {
		return errors.New("no github to repository has been provided")
	}

	if r.ToBranch == "" {
		return errors.New("no github to branch has been provided")
	}

	return nil
}

// ToJSON marshals the response to JSON.
func (r *SyncBranchRes) ToJSON() ([]byte, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("marshaling list changed files to JSON: %w", err)
	}

	return b, nil
}

// ToTable marshals the response to a text table.
func (r *SyncBranchRes) ToTable(err error) table.Writer {
	t := table.NewWriter()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false
	t.Style().Options.SeparateFooter = false
	t.Style().Options.SeparateHeader = false
	t.Style().Options.SeparateRows = false
	t.AppendHeader(table.Row{
		"From", "To", "Error",
	})

	if r.Request != nil {
		from := r.Request.FromOwner + "/" + r.Request.FromRepo + "/" + r.Request.FromBranch
		to := r.Request.ToOwner + "/" + r.Request.ToRepo + "/" + r.Request.ToBranch
		row := table.Row{from, to}
		if err != nil {
			row = append(row, err.Error())
		} else {
			row = append(row, nil)
		}
		t.AppendRow(row)
	}

	t.SuppressEmptyColumns()
	t.SuppressTrailingSpaces()

	return t
}
