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

	libgithub "github.com/google/go-github/v81/github"
	libgit "github.com/hashicorp/vault/tools/pipeline/internal/pkg/git"
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/golang"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/pmezard/go-difflib/difflib"
	slogctx "github.com/veqryn/slog-context"
)

// CheckGoModDiffReq is a request to create a diff between one-or-more go.mod
// files in two different Github hosted branches.
type CheckGoModDiffReq struct {
	// The A side of the diff
	AOwner  string `json:"a_owner,omitempty"`
	ARepo   string `json:"a_repo,omitempty"`
	ABranch string `json:"a_branch,omitempty"`

	// The B side of the diff
	BOwner  string `json:"b_owner,omitempty"`
	BRepo   string `json:"b_repo,omitempty"`
	BBranch string `json:"b_branch,omitempty"`

	// Paths to go.mod files to compare. The paths must be the same on both
	// sides of the diff.
	Paths []string `json:"paths,omitempty"`

	// DiffOpts the option to pass to the Go module diff.
	DiffOpts *golang.DiffOpts `json:"diff_opts,omitempty"`
}

// CheckGoModDiffRes is a response of checking remote Go modules for diffs.
type CheckGoModDiffRes struct {
	Diffs []*CheckGoModDiff `json:"diffs,omitempty"`
}

// CheckGoModDiff is one instance of a checked Go module diff.
type CheckGoModDiff struct {
	Err     error          `json:"err,omitempty"`
	Error   string         `json:"error,omitempty"`
	Path    string         `json:"path,omitempty"`
	ModDiff golang.ModDiff `json:"diff,omitempty"`
}

// Run runs the request check go.mod's for diffs.
func (r *CheckGoModDiffReq) Run(
	ctx context.Context,
	github *libgithub.Client,
	git *libgit.Client,
) (*CheckGoModDiffRes, error) {
	var err error
	res := &CheckGoModDiffRes{}

	slog.Default().DebugContext(slogctx.Append(ctx,
		slog.String("from-owner", r.AOwner),
		slog.String("from-repo", r.ARepo),
		slog.String("from-branch", r.ABranch),
		slog.String("to-owner", r.BOwner),
		slog.String("to-repo", r.BRepo),
		slog.String("to-branch", r.BBranch),
	), "checking go.mod diffs")

	if err = r.validate(); err != nil {
		return nil, fmt.Errorf("validating request: %w", err)
	}

	// Create a temp repository directory where we can keep both branches we're
	// comparing.
	repoDir, err, tmpDir := ensureGitRepoDir(ctx, "")
	if err != nil {
		return res, err
	}
	if tmpDir {
		defer os.RemoveAll(repoDir)
	}

	// It is entirely possible that we're comparing the same code but different
	// branches, or completely different code bases. It doesn't matter, really.
	// All that needs to be true is that the given paths match on each side of the
	// diff.

	// Determine what we'll consider our origin. If we've been given the same
	// owner and repo on both sides we'll use a single origin.
	aOrigin := "aorigin"
	bOrigin := "borigin"
	if r.AOwner == r.BOwner && r.ARepo == r.BRepo {
		aOrigin = "origin"
		bOrigin = "origin"
	}

	// First, get our A branch. This will also change our working directory into
	// the repoDir.
	err = initializeNewRepo(
		ctx, git, repoDir, r.AOwner, r.ARepo, aOrigin, r.ABranch,
	)
	if err != nil {
		return res, err
	}

	// Get our B branch. Start by adding our second upstream if necessary.
	if bOrigin != "origin" {
		// Our origin differs. Fetch the B origin and track the branch we care about.
		slog.Default().DebugContext(ctx, "adding B upstream remote")
		remoteRes, err := git.Remote(ctx, &libgit.RemoteOpts{
			Command: libgit.RemoteCommandAdd,
			Track:   []string{r.BBranch},
			Fetch:   true,
			Name:    bOrigin,
			URL:     fmt.Sprintf("https://github.com/%s/%s.git", r.BOwner, r.BRepo),
		})
		if err != nil {
			err = fmt.Errorf("adding B upstream remote: %s, %w", remoteRes.String(), err)
			return res, err
		}
	}

	// Determine our local B branch name. It's entirely possible and perhaps likely
	// that the branches will be named the same on both sides of the diff. Handle
	// that by using a unique name on the B side if necessary.
	bBranch := r.BBranch
	if r.ABranch == bBranch {
		// Make sure we don't try and use the same branch name
		bBranch = "b-" + r.BBranch
	}

	// Create a local branch of B
	slog.Default().DebugContext(ctx, "fetching B branch")
	fetchRes, err := git.Fetch(ctx, &libgit.FetchOpts{
		Refspec:   []string{bOrigin, r.BBranch + ":" + bBranch},
		Porcelain: true,
	})
	if err != nil {
		err = fmt.Errorf("fetching B branch: %s, %w", fetchRes.String(), err)
		return res, err
	}

	// Diff each path that we've been configured with.
	if len(r.Paths) < 1 {
		slog.Default().DebugContext(ctx, "No go.mod paths have been given. Assuming go.mod in the root directory of the repositories")
		r.Paths = []string{"go.mod"}
	}

	for _, path := range r.Paths {
		slog.Default().DebugContext(
			slogctx.Append(ctx, slog.String("path", path)),
			"creating module diff",
		)
		diffCheck := &CheckGoModDiff{
			Path: path,
		}

		// Checkout just the go.mod file we care about
		aCheckoutRes, err := git.Checkout(ctx, &libgit.CheckoutOpts{
			Branch:   r.ABranch,
			PathSpec: []string{path},
		})
		if err != nil {
			diffCheck.Err = fmt.Errorf("checking out %s on A branch: %s: %w", path, aCheckoutRes.String(), err)
			diffCheck.Error = diffCheck.Err.Error()
			res.Diffs = append(res.Diffs, diffCheck)

			continue
		}

		// Read in the contents
		aSource := &golang.ModSource{Name: r.AOwner + "/" + r.ARepo + "/" + r.ABranch + ":" + path}
		slog.Default().DebugContext(
			slogctx.Append(ctx, slog.String("path", path)),
			"reading module from A repository",
		)
		aSource.Data, err = os.ReadFile(path)
		if err != nil {
			diffCheck.Err = fmt.Errorf("reading %s on A branch: %w", path, err)
			diffCheck.Error = diffCheck.Err.Error()
			res.Diffs = append(res.Diffs, diffCheck)

			continue
		}

		// Checkout just the go.mod file we care about
		bCheckoutRes, err := git.Checkout(ctx, &libgit.CheckoutOpts{
			Branch:   bBranch,
			PathSpec: []string{path},
		})
		if err != nil {
			diffCheck.Err = fmt.Errorf("checking out %s on B branch: %s: %w", path, bCheckoutRes.String(), err)
			diffCheck.Error = diffCheck.Err.Error()
			res.Diffs = append(res.Diffs, diffCheck)

			continue
		}

		// Read in the contents
		bSource := &golang.ModSource{Name: r.BOwner + "/" + r.BRepo + "/" + r.BBranch + ":" + path}
		slog.Default().DebugContext(
			slogctx.Append(ctx, slog.String("path", path)),
			"reading module from B repository",
		)
		bSource.Data, err = os.ReadFile(path)
		if err != nil {
			diffCheck.Err = fmt.Errorf("reading %s on B branch: %w", path, err)
			diffCheck.Error = diffCheck.Err.Error()
			res.Diffs = append(res.Diffs, diffCheck)

			continue
		}

		// Diff the contents
		diffCheck.ModDiff, err = golang.DiffModFiles(aSource, bSource, r.DiffOpts)
		if err != nil {
			diffCheck.Err = fmt.Errorf("checking diff on %s: %w", path, err)
			diffCheck.Error = diffCheck.Err.Error()
		}
		res.Diffs = append(res.Diffs, diffCheck)
	}

	return res, nil
}

// validate ensures that we've been given the minimum arguments necessary to complete a
// request.
func (r *CheckGoModDiffReq) validate() error {
	if r == nil {
		return errors.New("failed to initialize request")
	}

	if r.AOwner == "" {
		return errors.New("no A Github Organization has been provided")
	}

	if r.ARepo == "" {
		return errors.New("no A Github Repository has been provided")
	}

	if r.ABranch == "" {
		return errors.New("no A Branch has been provided")
	}

	if r.BOwner == "" {
		return errors.New("no B Github Organization has been provided")
	}

	if r.BRepo == "" {
		return errors.New("no B Github Repository has been provided")
	}

	if r.BBranch == "" {
		return errors.New("no B Branch has been provided")
	}

	if r.AOwner == r.BOwner && r.ARepo == r.BRepo && r.ABranch == r.BBranch {
		return errors.New("cannot use same repository and branch on both sides of the diff")
	}

	return nil
}

// ToTable marshals the response to a text table.
func (r *CheckGoModDiffRes) ToTable(err error) (table.Writer, error) {
	t := table.NewWriter()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false
	t.Style().Options.SeparateFooter = false
	t.Style().Options.SeparateHeader = false
	t.Style().Options.SeparateRows = false

	if r == nil || len(r.Diffs) == 0 || err != nil {
		if err != nil {
			t.AppendHeader(table.Row{"error"})
			t.AppendRow(table.Row{err.Error()})
		}
		t.SuppressEmptyColumns()
		t.SuppressTrailingSpaces()

		return t, err
	}

	t.AppendHeader(table.Row{"path", "explanation", "diff"})
	for _, check := range r.Diffs {
		for _, diff := range check.ModDiff {
			if diff == nil {
				continue
			}
			if diff.Diff == nil {
				return nil, fmt.Errorf("missing unified diff: %v", diff)
			}

			diffText, err := difflib.GetUnifiedDiffString(*diff.Diff)
			if err != nil {
				return nil, err
			}
			t.AppendRow(table.Row{check.Path, diff.Directive.Explanation(), diffText})
		}
	}
	t.SuppressEmptyColumns()
	t.SuppressTrailingSpaces()

	return t, nil
}

// ToJSON marshals the response to JSON.
func (r *CheckGoModDiffRes) ToJSON() ([]byte, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("marshaling latest HCP image response to JSON: %w", err)
	}

	return b, nil
}
