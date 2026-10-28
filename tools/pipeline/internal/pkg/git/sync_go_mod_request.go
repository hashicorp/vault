// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package git

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	libgit "github.com/hashicorp/vault/tools/pipeline/internal/pkg/git/client"
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/golang"
	"github.com/jedib0t/go-pretty/v6/table"
)

// SyncGoModReq is a request to synchronize go.mod files between two local branches.
type SyncGoModReq struct {
	SourceBranch string
	DestBranch   string

	// Paths are the go.mod paths to sync. Paths are relative to the working
	// directory. When empty and WorkspaceModules isn't set, the go.mod in the
	// working directory is synced.
	Paths []string

	// WorkspaceModules syncs the go.mod of every module in the go.work at the
	// root of SourceBranch instead of Paths, which are then relative to the root
	// of the repository whichever directory we run in. Once they're synced it
	// tidies every module in the go.work of DestBranch with GOWORK=off.
	WorkspaceModules bool

	DiffOpts *golang.DiffOpts

	// NewBranch is the name for the intermediate branch. If empty we auto-generate
	// one in the form <source>-into-<dest> (slashes replaced with dashes).
	NewBranch string

	// Commit controls whether staged changes are committed. Defaults to true in
	// the command layer; pass --commit=false to skip.
	Commit        bool
	CommitMessage string // overrides the auto-generated subject when set
	Ticket        string // bare ticket ID, e.g. "VAULT-1234"
}

// SyncGoModRes is the result of a local-branch go.mod sync operation.
type SyncGoModRes struct {
	SourceBranch string                        `json:"source_branch"`
	DestBranch   string                        `json:"dest_branch"`
	NewBranch    string                        `json:"new_branch"`
	CommitSHA    string                        `json:"commit_sha,omitempty"`
	Results      []*golang.GoModPathSyncResult `json:"results,omitempty"`
	// GoWorkSum is set when we tried to sync the workspace.
	GoWorkSum *golang.GoWorkSumSyncResult `json:"go_work_sum,omitempty"`
}

// Run forks a new branch from DestBranch, reads source go.mod files from
// SourceBranch via git show, applies golang.SyncModReq per path, runs go mod
// tidy for the modules that changed, and optionally commits the result. In
// workspace mode it syncs every module of the go.work of SourceBranch, tidies
// every module of the go.work of DestBranch with GOWORK=off.
func (r *SyncGoModReq) Run(ctx context.Context, git *libgit.Client) (*SyncGoModRes, error) {
	slog.Default().DebugContext(ctx, "starting git sync go-mod",
		slog.String("source_branch", r.SourceBranch),
		slog.String("dest_branch", r.DestBranch),
		slog.Int("paths", len(r.Paths)),
		slog.Bool("workspace-modules", r.WorkspaceModules),
	)

	if err := r.validate(); err != nil {
		return nil, err
	}

	slog.Default().DebugContext(ctx, "checking working tree status")
	statusRes, err := git.Status(ctx, &libgit.StatusOpts{Porcelain: true})
	if err != nil {
		return nil, fmt.Errorf("git status: %w", err)
	}
	if len(strings.TrimSpace(string(statusRes.Stdout))) > 0 {
		return nil, errors.New(
			"working tree has uncommitted changes; commit, stash, or discard them first",
		)
	}

	// The go.mod paths to sync, and the directory they're relative to. A path
	// from the user is relative to the working directory. The modules of a go.work
	// are relative to the root of the repository, which isn't always where we
	// run, so in workspace mode we work from the root without changing directory.
	// We list the modules before we create the branch so that a source without a
	// usable go.work leaves nothing behind. A failure is fatal so that we never
	// fall back to syncing only the root go.mod.
	paths := r.Paths
	dir := "."
	if r.WorkspaceModules {
		topRes, err := git.RevParse(ctx, &libgit.RevParseOpts{ShowTopLevel: true})
		if err != nil {
			return nil, fmt.Errorf("finding the root of the repository: %w", err)
		}
		dir = strings.TrimSuffix(string(topRes.Stdout), "\n")

		paths, err = WorkspaceGoModPathsFromRef(ctx, git, r.SourceBranch)
		if err != nil {
			return nil, err
		}
	}

	newBranch, err := r.resolveBranchName(ctx, git)
	if err != nil {
		return nil, fmt.Errorf("resolving branch name: %w", err)
	}
	slog.Default().DebugContext(ctx, "creating intermediate branch",
		slog.String("new_branch", newBranch),
		slog.String("start_point", r.DestBranch),
	)

	// Fork and check out the new intermediate branch.
	coRes, coErr := git.Checkout(ctx, &libgit.CheckoutOpts{
		NewBranch:  newBranch,
		StartPoint: r.DestBranch,
	})
	if coErr != nil {
		return nil, fmt.Errorf("git checkout -b %s %s: %w\n%s", newBranch, r.DestBranch, coErr, coRes.String())
	}

	res := &SyncGoModRes{
		SourceBranch: r.SourceBranch,
		DestBranch:   r.DestBranch,
		NewBranch:    newBranch,
	}

	// A collection of errors. We do this so that our sync is best effort for each
	// go.mod. We'll bubble up any issues we run into along the way.
	var errs []error

	// Sync the directives of every module before we tidy any of them, as tidying
	// a module needs the modules it depends on to be synced.
	for _, p := range paths {
		pathRes, pathErr := SyncGoModFromRef(ctx, git, r.SourceBranch, dir, p, r.DiffOpts)
		res.Results = append(res.Results, pathRes)
		if pathErr != nil {
			errs = append(errs, pathErr)
		}
	}

	// Tidy the modules until they settle, and bring the workspace in line when
	// we're syncing the workspace.
	tidyRes, tidyErr := golang.TidySyncedModules(ctx, dir, res.Results, r.WorkspaceModules)
	res.Results = tidyRes.Results
	res.GoWorkSum = tidyRes.GoWorkSum
	if tidyErr != nil {
		errs = append(errs, tidyErr)
	}
	stageFiles := tidyRes.Stage

	// Stage and commit when requested and there's something to commit.
	if r.Commit && len(stageFiles) > 0 {
		slog.Default().DebugContext(ctx, "staging files", slog.Int("count", len(stageFiles)))
		if _, addErr := git.Add(ctx, &libgit.AddOpts{PathSpec: stageFiles}); addErr != nil {
			errs = append(errs, fmt.Errorf("git add: %w", addErr))
		} else {
			subject := r.commitSubject()
			slog.Default().DebugContext(ctx, "committing changes", slog.String("subject", subject))
			if _, commitErr := git.Commit(ctx, &libgit.CommitOpts{Message: subject}); commitErr != nil {
				errs = append(errs, fmt.Errorf("git commit: %w", commitErr))
			} else {
				slog.Default().DebugContext(ctx, "capturing commit SHA")
				shaRes, shaErr := git.RevParse(ctx, &libgit.RevParseOpts{Args: []string{"HEAD"}})
				if shaErr != nil {
					errs = append(errs, fmt.Errorf("git rev-parse HEAD after commit: %w", shaErr))
				} else {
					res.CommitSHA = strings.TrimSpace(string(shaRes.Stdout))
				}
			}
		}
	}

	slog.Default().DebugContext(ctx, "completed git sync go-mod",
		slog.Int("paths", len(paths)),
		slog.Int("errors", len(errs)),
	)

	return res, errors.Join(errs...)
}

// ToTable returns a table writer populated with the per-path summary rows.
func (r *SyncGoModRes) ToTable() table.Writer {
	if r == nil {
		return table.NewWriter()
	}
	t := table.NewWriter()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false
	t.Style().Options.SeparateFooter = false
	t.Style().Options.SeparateHeader = false
	t.Style().Options.SeparateRows = false
	t.SuppressEmptyColumns()
	t.SuppressTrailingSpaces()

	t.AppendHeader(table.Row{"path", "changes", "tidy", "error"})
	for _, pr := range r.Results {
		if pr == nil {
			continue
		}
		changes := 0
		errStr := ""
		if pr.SyncModRes != nil {
			changes = len(pr.SyncModRes.Changes)
			if pr.SyncModRes.Error != "" {
				errStr = pr.SyncModRes.Error
			}
		}
		if pr.TidyError != "" {
			if errStr != "" {
				errStr = errStr + "; " + pr.TidyError
			} else {
				errStr = pr.TidyError
			}
		}
		tidyStr := "no"
		if pr.TidyRan {
			tidyStr = "yes"
		}
		path := ""
		if pr.SyncModRes != nil {
			path = pr.SyncModRes.Path
		}
		t.AppendRow(table.Row{path, changes, tidyStr, errStr})
	}

	// go.work.sum is a file we changed when we synced the workspace. It has one
	// change when it was updated, and the reason when we couldn't sync it.
	if r.GoWorkSum != nil {
		changes := 0
		if r.GoWorkSum.Updated {
			changes = 1
		}
		t.AppendRow(table.Row{"go.work.sum", changes, "no", r.GoWorkSum.Error})
	}

	return t
}

// ToString renders the result as a plain-text summary table.
func (r *SyncGoModRes) ToString() string {
	if r == nil {
		return ""
	}

	return r.ToTable().Render()
}

// ToMarkdown renders the result as Markdown.
func (r *SyncGoModRes) ToMarkdown() string {
	if r == nil {
		return ""
	}

	return r.ToTable().RenderMarkdown()
}

// ToJSON marshals the result to JSON.
func (r *SyncGoModRes) ToJSON() ([]byte, error) {
	if r == nil {
		return nil, errors.New("sync go-mod: nil result")
	}
	b, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("marshaling sync go-mod result to JSON: %w", err)
	}
	return b, nil
}

// validate checks the request fields before any git operations run.
func (r *SyncGoModReq) validate() error {
	if r.SourceBranch == "" {
		return errors.New("source branch is required")
	}
	if r.DestBranch == "" {
		return errors.New("dest branch is required")
	}
	if r.SourceBranch == r.DestBranch {
		return fmt.Errorf("source and dest branches must differ: both are %q", r.SourceBranch)
	}
	if r.WorkspaceModules && len(r.Paths) > 0 {
		return errors.New("paths cannot be combined with workspace-modules")
	}

	// Without workspace modules we sync the go.mod in the working directory unless
	// told otherwise.
	if !r.WorkspaceModules && len(r.Paths) == 0 {
		r.Paths = []string{"go.mod"}
	}

	return nil
}

// resolveBranchName returns the name to use for the new intermediate branch.
// If r.NewBranch is already set we use it as-is. Otherwise we generate a name
// from the ticket (if any), source, and dest branch names and append a unix
// timestamp suffix when that name already exists locally.
func (r *SyncGoModReq) resolveBranchName(ctx context.Context, git *libgit.Client) (string, error) {
	name := r.NewBranch
	if name == "" {
		name = r.branchName()
	}

	// Check whether the branch already exists locally.
	_, err := git.RevParse(ctx, &libgit.RevParseOpts{
		Verify: true,
		Quiet:  true,
		Args:   []string{name},
	})
	if err != nil {
		return name, nil
	}

	// Branch already exists — append a unix timestamp to make the name unique.
	return name[:min(len(name), 230)] + "-" + strconv.FormatInt(time.Now().Unix(), 10), nil
}

func (r *SyncGoModReq) branchName() string {
	base := strings.ReplaceAll(r.SourceBranch, "/", "-") + "-into-" + strings.ReplaceAll(r.DestBranch, "/", "-")
	// GitHub limits refs/heads/<name> to 255 characters, so we cap at 240 to
	// leave headroom for the refs/heads/ prefix.
	if r.Ticket != "" {
		full := r.Ticket + "-" + base
		return full[:min(len(full), 240)]
	}
	return base[:min(len(base), 240)]
}

// commitSubject returns the auto-generated commit subject line.
func (r *SyncGoModReq) commitSubject() string {
	if r.CommitMessage != "" {
		return r.CommitMessage
	}
	if r.Ticket != "" {
		return fmt.Sprintf("[%s] go: sync go.mod from %s to %s", r.Ticket, r.SourceBranch, r.DestBranch)
	}
	return fmt.Sprintf("go: sync go.mod from %s to %s", r.SourceBranch, r.DestBranch)
}
