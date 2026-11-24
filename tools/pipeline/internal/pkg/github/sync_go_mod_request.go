// Copyright IBM Corp. 2016, 2026
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
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	libgithub "github.com/google/go-github/v83/github"
	gitpkg "github.com/hashicorp/vault/tools/pipeline/internal/pkg/git"
	libgit "github.com/hashicorp/vault/tools/pipeline/internal/pkg/git/client"
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/golang"
	"github.com/jedib0t/go-pretty/v6/table"
)

// SyncGoModReq is a request to sync go.mod files from a remote A branch into a
// new intermediate branch forked from remote B, then push and open a PR on the
// TO remote.
type SyncGoModReq struct {
	// A-side (source — read-only).
	AHost   string
	AOwner  string
	ARepo   string
	ABranch string
	AOrigin string // remote name for A in the local clone

	// B-side (dest — fork point and PR base).
	BHost   string
	BOwner  string
	BRepo   string
	BBranch string
	BOrigin string // remote name for B in the local clone

	// TO-side (push target and PR target — defaults to B values).
	ToHost   string
	ToOwner  string
	ToRepo   string
	ToBranch string // intermediate branch name; auto-generated if empty
	ToOrigin string // remote name for TO in the local clone

	// Repo isolation — if empty a temp dir is created and cleaned up on exit.
	RepoDir string

	// go.mod paths to sync, relative to the repository root. When empty and
	// WorkspaceModules isn't set, the root go.mod is synced.
	Paths []string

	// WorkspaceModules syncs the go.mod of every module in the go.work at the
	// root of the A branch instead of Paths. Once they're synced it tidies every
	// module in the go.work of the clone with GOWORK=off.
	WorkspaceModules bool

	// Options governing which directives to compare and sync.
	DiffOpts *golang.DiffOpts

	// Commit control.
	Commit        bool // default true
	CommitMessage string
	Ticket        string

	// PR control.
	PullRequest bool   // default true
	PRTitle     string // override auto-generated title
	PRBase      string // override PR base branch (default: BBranch)
	PRReviewers []string
	PRAssignees []string
}

// SyncGoModRes is the result of a remote go.mod sync operation.
type SyncGoModRes struct {
	AOwner      string                        `json:"a_owner"`
	ARepo       string                        `json:"a_repo"`
	ABranch     string                        `json:"a_branch"`
	BOwner      string                        `json:"b_owner"`
	BRepo       string                        `json:"b_repo"`
	BBranch     string                        `json:"b_branch"`
	ToBranch    string                        `json:"to_branch"`
	CommitSHA   string                        `json:"commit_sha,omitempty"`
	PullRequest *libgithub.PullRequest        `json:"pull_request,omitempty"`
	Results     []*golang.GoModPathSyncResult `json:"results,omitempty"`
	// GoWorkSum is set when we tried to sync the workspace.
	GoWorkSum *golang.GoWorkSumSyncResult `json:"go_work_sum,omitempty"`
}

// Run forks a new branch from B, syncs go.mod files from A into it, pushes to
// the TO remote, and optionally opens a PR.
func (r *SyncGoModReq) Run(
	ctx context.Context,
	github *libgithub.Client,
	git *libgit.Client,
) (*SyncGoModRes, error) {
	slog.Default().DebugContext(ctx, "starting github sync go-mod",
		slog.String("a-owner", r.AOwner),
		slog.String("a-repo", r.ARepo),
		slog.String("a-branch", r.ABranch),
		slog.String("b-owner", r.BOwner),
		slog.String("b-repo", r.BRepo),
		slog.String("b-branch", r.BBranch),
		slog.Int("mod-count", len(r.Paths)),
		slog.Bool("workspace-modules", r.WorkspaceModules),
	)

	if err := r.validate(); err != nil {
		return nil, err
	}

	// The go.mod paths to sync. In workspace mode we derive them from the A branch
	// once we've fetched it.
	paths := r.Paths

	// Save the initial working directory so we can restore it on exit.
	initialDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("getting working directory: %w", err)
	}
	defer func() { _ = os.Chdir(initialDir) }()

	// Paths must be repo-root-relative (e.g. "go.mod", "sdk/go.mod"). After we
	// os.Chdir into the clone the CWD is the repo root, so the same path works
	// for git show <ref>:<path>, os.ReadFile, and go mod tidy. We reject absolute
	// paths here — they can't be used as git object paths.
	for _, p := range paths {
		if filepath.IsAbs(p) {
			return nil, fmt.Errorf("path %q must be repo-root-relative, not absolute", p)
		}
	}

	// Set up the repo dir, creating a temp dir when RepoDir is empty.
	repoDir, err, tmpDir := ensureGitRepoDir(ctx, r.RepoDir)
	if err != nil {
		return nil, err
	}
	if tmpDir {
		defer os.RemoveAll(repoDir)
	}

	// Initialize the repo — reuse an existing clone or create a fresh one from B.
	_, statErr := os.Stat(filepath.Join(repoDir, ".git"))
	if statErr == nil {
		// Check before pulling or checking out, so neither autostash nor staging
		// can carry pre-existing edits into the branch we publish.
		repoGit := *git
		repoGit.Dir = repoDir
		slog.Default().DebugContext(ctx, "checking reused repository working tree")
		statusRes, statusErr := repoGit.Status(ctx, &libgit.StatusOpts{
			Porcelain: true, UntrackedFiles: libgit.UntrackedFilesAll,
		})
		if statusErr != nil {
			return nil, fmt.Errorf("git status: %w\n%s", statusErr, statusRes.String())
		}
		if len(statusRes.Stdout) > 0 {
			return nil, errors.New("working tree has uncommitted changes; commit, stash, or discard them first")
		}
		err = initializeExistingRepo(ctx, git, repoDir, r.BOrigin, r.BBranch)
	} else {
		err = initializeNewRepo(ctx, git, repoDir, r.BOwner, r.BRepo, r.BOrigin, r.BBranch)
	}
	if err != nil {
		return nil, err
	}

	// Check out the B branch for a clean starting point.
	slog.Default().DebugContext(ctx, "checking out b-branch", slog.String("b-branch", r.BBranch))
	coRes, coErr := git.Checkout(ctx, &libgit.CheckoutOpts{Branch: r.BBranch})
	if coErr != nil {
		return nil, fmt.Errorf("git checkout %s: %w\n%s", r.BBranch, coErr, coRes.String())
	}

	// Add the A remote and fetch the A branch. When A and B are the same repo we
	// skip adding a new remote and just fetch through the existing origin.
	if r.AOwner != r.BOwner || r.ARepo != r.BRepo {
		slog.Default().DebugContext(ctx, "adding a-remote and fetching a-branch",
			slog.String("a-origin", r.AOrigin),
			slog.String("a-branch", r.ABranch),
		)
		remoteRes, remoteErr := git.Remote(ctx, &libgit.RemoteOpts{
			Command: libgit.RemoteCommandAdd,
			Track:   []string{r.ABranch},
			Fetch:   true,
			Name:    r.AOrigin,
			URL:     fmt.Sprintf("https://%s/%s/%s.git", r.AHost, r.AOwner, r.ARepo),
		})
		if remoteErr != nil {
			return nil, fmt.Errorf("adding a-remote and fetching: %s: %w", remoteRes.String(), remoteErr)
		}
	} else {
		slog.Default().DebugContext(ctx, "fetching a-branch via existing origin",
			slog.String("a-origin", r.AOrigin),
			slog.String("a-branch", r.ABranch),
		)
		// Use an explicit src:dst refspec so git creates the remote-tracking ref
		// refs/remotes/<origin>/<branch> even on a --single-branch clone, which
		// locks the fetch refspec to only the B branch. Without the explicit dst,
		// git fetch fetches the objects but doesn't update any remote-tracking ref,
		// so subsequent git show <origin>/<branch>:<path> calls fail.
		refspec := fmt.Sprintf("refs/heads/%s:refs/remotes/%s/%s", r.ABranch, r.AOrigin, r.ABranch)
		fetchRes, fetchErr := git.Fetch(ctx, &libgit.FetchOpts{
			Repository: r.AOrigin,
			Refspec:    []string{refspec},
		})
		if fetchErr != nil {
			return nil, fmt.Errorf("git fetch %s %s: %w\n%s", r.AOrigin, r.ABranch, fetchErr, fetchRes.String())
		}
	}

	// Add the TO remote when it's a different repo from B.
	if r.ToOwner != r.BOwner || r.ToRepo != r.BRepo {
		slog.Default().DebugContext(ctx, "adding to-remote",
			slog.String("to-origin", r.ToOrigin),
		)
		remoteRes, remoteErr := git.Remote(ctx, &libgit.RemoteOpts{
			Command: libgit.RemoteCommandAdd,
			Name:    r.ToOrigin,
			URL:     fmt.Sprintf("https://%s/%s/%s.git", r.ToHost, r.ToOwner, r.ToRepo),
		})
		if remoteErr != nil {
			return nil, fmt.Errorf("adding to-remote: %s: %w", remoteRes.String(), remoteErr)
		}
	}

	// Derive the modules to sync from the A branch's go.work, so we sync the
	// modules that A has rather than whichever go.work is in the working tree. A
	// failure is fatal so that we never fall back to syncing only the root go.mod.
	if r.WorkspaceModules {
		paths, err = gitpkg.WorkspaceGoModPathsFromRef(ctx, git, r.AOrigin+"/"+r.ABranch)
		if err != nil {
			return nil, err
		}
	}

	toBranch, err := r.resolveToBranch(ctx, git)
	if err != nil {
		return nil, fmt.Errorf("resolving to branch name: %w", err)
	}

	slog.Default().DebugContext(ctx, "creating to branch",
		slog.String("to-branch", toBranch),
		slog.String("start-point", r.BOrigin+"/"+r.BBranch),
	)
	newCoRes, err := git.Checkout(ctx, &libgit.CheckoutOpts{
		NewBranch:  toBranch,
		StartPoint: r.BOrigin + "/" + r.BBranch,
	})
	if err != nil {
		return nil, fmt.Errorf(
			"git checkout -b %s %s/%s: %w\n%s",
			toBranch, r.BOrigin, r.BBranch, err, newCoRes.String(),
		)
	}

	res := &SyncGoModRes{
		AOwner:   r.AOwner,
		ARepo:    r.ARepo,
		ABranch:  r.ABranch,
		BOwner:   r.BOwner,
		BRepo:    r.BRepo,
		BBranch:  r.BBranch,
		ToBranch: toBranch,
	}

	// A collection of errors. We do this so that our sync is best effort for each
	// go.mod. We'll bubble up any issues we run into along the way.
	var errs []error

	// Sync the directives of every Go module before we tidy any of them. Tidying
	// a module reads the go.mod of the local modules it replaces, so a module that
	// we tidy before the modules it depends on are synced is tidied against stale
	// requirements and needs another pass.
	for _, p := range paths {
		pathRes, err := gitpkg.SyncGoModFromRef(ctx, git, r.AOrigin+"/"+r.ABranch, ".", p, r.DiffOpts)
		res.Results = append(res.Results, pathRes)
		if err != nil {
			errs = append(errs, err)
		}
	}

	// Tidy the modules until they settle, and bring the workspace in line when
	// we're syncing the workspace. In workspace mode the modules to tidy come from
	// the go.work in the clone, which is the B branch's rather than A's. A clone
	// without a go.work has no workspace modules.
	tidyRes, err := golang.TidySyncedModules(ctx, ".", res.Results, r.WorkspaceModules)
	res.Results = tidyRes.Results
	res.GoWorkSum = tidyRes.GoWorkSum
	if err != nil {
		errs = append(errs, err)
	}
	stageFiles := tidyRes.Stage

	// Stage and commit changed files.
	stageAndCommitFiles := func() error {
		if !r.Commit {
			slog.Default().DebugContext(ctx, "skipping file stage, commit is disabled")
			return nil
		}
		if len(stageFiles) < 1 {
			slog.Default().DebugContext(ctx, "skipping file stage, no files modified")
			return nil
		}

		slog.Default().DebugContext(ctx, "staging files", slog.Int("count", len(stageFiles)))

		addRes, err := git.Add(ctx, &libgit.AddOpts{PathSpec: stageFiles})
		if err != nil {
			return fmt.Errorf("git add: %w\n%s", err, addRes.String())
		}

		// Tidy can undo every synced directive, leaving no net changes to commit.
		slog.Default().DebugContext(ctx, "checking staged changes")
		diffRes, err := git.Diff(ctx, &libgit.DiffOpts{Cached: true, NameOnly: true})
		if err != nil {
			return fmt.Errorf("git diff --cached: %w\n%s", err, diffRes.String())
		}
		if len(diffRes.Stdout) == 0 {
			slog.Default().DebugContext(ctx, "skipping commit, no staged changes")
			return nil
		}

		subject := r.commitSubject()
		slog.Default().DebugContext(ctx, "committing changes", slog.String("subject", subject))
		commitRes, err := git.Commit(ctx, &libgit.CommitOpts{Message: subject})
		if err != nil {
			return fmt.Errorf("git commit: %w\n%s", err, commitRes.String())
		}

		slog.Default().DebugContext(ctx, "capturing commit SHA") //nolint:govet
		shaRes, err := git.RevParse(ctx, &libgit.RevParseOpts{Args: []string{"HEAD"}})
		if err != nil {
			return fmt.Errorf("git rev-parse HEAD after commit: %w", err)
		}

		res.CommitSHA = strings.TrimSpace(string(shaRes.Stdout))

		return nil
	}
	err = stageAndCommitFiles()
	if err != nil {
		errs = append(errs, err)
	}

	// Push the new branch to the TO remote once we have a commit.
	if res.CommitSHA != "" {
		slog.Default().DebugContext(ctx, "pushing intermediate branch",
			slog.String("to-origin", r.ToOrigin),
			slog.String("to-branch", toBranch),
		)
		pushRes, err := git.Push(ctx, &libgit.PushOpts{
			Repository: r.ToOrigin,
			Refspec:    []string{toBranch},
		})
		if err != nil {
			return res, fmt.Errorf("pushing intermediate branch: %s: %w", pushRes.String(), err)
		}
	}

	// Open a PR when requested and we have something pushed. We try and handle
	// errors in the various side quests like assigment and requesting review
	// gracefully.
	if r.PullRequest && res.CommitSHA != "" {
		prTitle := r.prTitle()

		hasPathErrors := func(results []*golang.GoModPathSyncResult) bool {
			for _, pr := range results {
				if pr == nil {
					continue
				}
				if pr.SyncModRes != nil && pr.SyncModRes.Error != "" {
					return true
				}
				if pr.TidyError != "" {
					return true
				}
			}
			return false
		}

		prBody, err := renderEmbeddedTemplate("sync-go-mod-pr-message.tmpl", struct {
			Ticket      string
			AOwner      string
			ARepo       string
			ABranch     string
			BBranch     string
			ToBranch    string
			PathResults []*golang.GoModPathSyncResult
			HasErrors   bool
			GoWorkSum   *golang.GoWorkSumSyncResult
		}{
			Ticket:      r.Ticket,
			AOwner:      r.AOwner,
			ARepo:       r.ARepo,
			ABranch:     r.ABranch,
			BBranch:     r.BBranch,
			ToBranch:    toBranch,
			PathResults: res.Results,
			HasErrors: hasPathErrors(res.Results) ||
				(res.GoWorkSum != nil && res.GoWorkSum.Error != ""),
			GoWorkSum: res.GoWorkSum,
		})

		if err != nil {
			errs = append(errs, fmt.Errorf("rendering PR body: %w", err))
		} else {
			limitedBody := limitCharacters(prBody)
			prBase := r.PRBase
			slog.Default().DebugContext(ctx, "creating pull request",
				slog.String("title", prTitle),
				slog.String("head", toBranch),
				slog.String("base", prBase),
			)
			pr, _, err := github.PullRequests.Create(
				ctx, r.ToOwner, r.ToRepo, &libgithub.NewPullRequest{
					Title: &prTitle,
					Head:  &toBranch,
					Base:  &prBase,
					Body:  &limitedBody,
				},
			)
			if err != nil {
				errs = append(errs, fmt.Errorf("creating pull request: %w", err))
			} else {
				res.PullRequest = pr
				prNum := int(pr.GetNumber())
				if len(r.PRAssignees) > 0 {
					if err := addAssignees(ctx, github, r.ToOwner, r.ToRepo, prNum, r.PRAssignees); err != nil {
						errs = append(errs, fmt.Errorf("adding PR assignees: %w", err))
					}
				}
				if len(r.PRReviewers) > 0 {
					if revErr := addReviewers(ctx, github, r.ToOwner, r.ToRepo, prNum, r.PRReviewers); revErr != nil {
						errs = append(errs, fmt.Errorf("adding PR reviewers: %w", revErr))
					}
				}
			}
		}
	}

	slog.Default().DebugContext(ctx, "completed github sync go-mod",
		slog.Int("paths", len(paths)),
		slog.Int("errors", len(errs)),
	)

	return res, errors.Join(errs...)
}

// ToString renders the result as plain text: metadata block followed by the
// per-path results table.
func (r *SyncGoModRes) ToString(err error) string {
	if r == nil {
		return r.metaTable(err).Render()
	}
	meta := r.metaTable(err)
	paths := r.pathTable()
	return meta.Render() + "\n" + paths.Render()
}

// ToMarkdown renders the result as Markdown: metadata block followed by the
// per-path results table. Same layout as ToString.
func (r *SyncGoModRes) ToMarkdown(err error) string {
	if r == nil {
		return escapeMarkdownTableColumns(r.metaTable(err), 2).RenderMarkdown()
	}
	meta := r.metaTable(err)
	paths := r.pathTable()
	return escapeMarkdownTableColumns(meta, 2).RenderMarkdown() + "\n" +
		escapeMarkdownTableColumns(paths, 1, 4).RenderMarkdown()
}

// ToJSON marshals the result to JSON.
func (r *SyncGoModRes) ToJSON() ([]byte, error) {
	if r == nil {
		return nil, errors.New("github sync go-mod result is nil")
	}
	b, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("marshaling github sync go-mod result to JSON: %w", err)
	}
	return b, nil
}

// validate checks the request fields and resolves defaults before any operations.
func (r *SyncGoModReq) validate() error {
	if r.AOwner == "" {
		return errors.New("a-owner is required")
	}
	if r.ARepo == "" {
		return errors.New("a-repo is required")
	}
	if r.ABranch == "" {
		return errors.New("a-branch is required")
	}
	if r.BOwner == "" {
		return errors.New("b-owner is required")
	}
	if r.BRepo == "" {
		return errors.New("b-repo is required")
	}
	if r.BBranch == "" {
		return errors.New("b-branch is required")
	}
	if r.WorkspaceModules && len(r.Paths) > 0 {
		return errors.New("paths cannot be combined with workspace-modules")
	}

	// Without workspace modules we sync the root go.mod unless told otherwise.
	if !r.WorkspaceModules && len(r.Paths) == 0 {
		r.Paths = []string{"go.mod"}
	}

	// TO-side defaults from B.
	if r.ToOwner == "" {
		r.ToOwner = r.BOwner
	}
	if r.ToRepo == "" {
		r.ToRepo = r.BRepo
	}
	if r.ToHost == "" {
		r.ToHost = r.BHost
	}

	// Host defaults.
	if r.AHost == "" {
		r.AHost = "github.com"
	}
	if r.BHost == "" {
		r.BHost = "github.com"
	}
	if r.ToHost == "" {
		r.ToHost = "github.com"
	}

	// Resolve origin names. When A and B are the same repo we share "origin" for
	// both — consistent with how check go-mod-diff behaves.
	if r.AOwner == r.BOwner && r.ARepo == r.BRepo {
		if r.AOrigin == "" {
			r.AOrigin = "origin"
		}
		if r.BOrigin == "" {
			r.BOrigin = "origin"
		}
	} else {
		if r.AOrigin == "" {
			r.AOrigin = "aorigin"
		}
		if r.BOrigin == "" {
			r.BOrigin = "borigin"
		}
	}
	if r.ToOrigin == "" {
		r.ToOrigin = r.BOrigin
	}

	// We use an explicit to-branch as it is, so one that GitHub won't accept
	// has to be an error rather than something we quietly shorten.
	if len(r.ToBranch) > maxBranchNameBytes {
		return fmt.Errorf(
			"to-branch is %d bytes long, but GitHub limits a branch name to %d bytes",
			len(r.ToBranch), maxBranchNameBytes,
		)
	}

	// PRBase defaults to the B branch.
	if r.PRBase == "" {
		r.PRBase = r.BBranch
	}

	// Syncing a branch into itself in the same repo makes no sense.
	if r.ABranch == r.BBranch && r.AOwner == r.BOwner && r.ARepo == r.BRepo {
		return fmt.Errorf(
			"source and dest branches must differ when using the same repo: both are %q",
			r.ABranch,
		)
	}

	return nil
}

// maxBranchNameBytes is the longest branch name we'll create. GitHub limits a
// ref name to 255 bytes, and that includes the refs/heads/ prefix.
const maxBranchNameBytes = 255 - len("refs/heads/")

// branchName returns the auto-generated intermediate branch name, in the form
// go-mod-sync-[<ticket>-]<a-branch>-into-<b-branch> with slashes replaced by
// dashes. It's never longer than maxBranchNameBytes.
func (r *SyncGoModReq) branchName() string {
	name := "go-mod-sync-"
	if r.Ticket != "" {
		name += r.Ticket + "-"
	}
	name += strings.ReplaceAll(r.ABranch, "/", "-") + "-into-" + strings.ReplaceAll(r.BBranch, "/", "-")

	return truncateBranchName(name, maxBranchNameBytes)
}

// truncateBranchName cuts name on a UTF-8 boundary, dropping trailing separators
// and .lock suffixes that Git rejects. A name that already fits is unchanged.
func truncateBranchName(name string, limit int) string {
	if len(name) <= limit {
		return name
	}

	// name[limit] is the first byte we drop. If it's a continuation byte then
	// the cut would split a character, so move it back to where it starts.
	for limit > 0 && !utf8.RuneStart(name[limit]) {
		limit--
	}

	name = strings.TrimRight(name[:limit], "-./")
	for strings.HasSuffix(name, ".lock") {
		name = strings.TrimRight(strings.TrimSuffix(name, ".lock"), "-./")
	}
	return name
}

// resolveToBranch returns the intermediate branch name. If ToBranch is already
// set it's used as-is. Otherwise we auto-generate one and append a unix timestamp
// suffix when that name already exists locally.
func (r *SyncGoModReq) resolveToBranch(ctx context.Context, git *libgit.Client) (string, error) {
	name := r.ToBranch
	if name == "" {
		name = r.branchName()
	}

	_, err := git.RevParse(ctx, &libgit.RevParseOpts{
		Verify: true,
		Quiet:  true,
		Args:   []string{name},
	})
	if err != nil {
		return name, nil
	}

	// Branch already exists, append a unix timestamp to make the name unique.
	// Make room for the suffix so we stay within GitHub's limit on branch names.
	suffix := "-" + strconv.FormatInt(time.Now().Unix(), 10)

	return truncateBranchName(name, maxBranchNameBytes-len(suffix)) + suffix, nil
}

// commitSubject returns the auto-generated commit subject line.
func (r *SyncGoModReq) commitSubject() string {
	if r.CommitMessage != "" {
		return r.CommitMessage
	}
	if r.Ticket != "" {
		return fmt.Sprintf("[%s] go: sync go.mod from %s to %s", r.Ticket, r.ABranch, r.BBranch)
	}
	return fmt.Sprintf("go: sync go.mod from %s to %s", r.ABranch, r.BBranch)
}

// prTitle returns the PR title, using PRTitle as an override when set.
func (r *SyncGoModReq) prTitle() string {
	if r.PRTitle != "" {
		return r.PRTitle
	}
	base := fmt.Sprintf("go: sync go.mod from %s/%s/%s to %s", r.AOwner, r.ARepo, r.ABranch, r.BBranch)
	if r.Ticket != "" {
		return fmt.Sprintf("[%s] %s", r.Ticket, base)
	}
	return base
}

// metaTable returns a table writer with A/B/branch/commit/PR metadata rows.
// Column 0 is the label, column 1 is the value.
func (r *SyncGoModRes) metaTable(err error) table.Writer {
	if r == nil {
		t := r.newResultTable()
		if err != nil {
			t.AppendRow(table.Row{"error", err.Error()})
		}
		return t
	}
	t := r.newResultTable()

	t.AppendRow(table.Row{"a", r.AOwner + "/" + r.ARepo + "/" + r.ABranch})
	t.AppendRow(table.Row{"b", r.BOwner + "/" + r.BRepo + "/" + r.BBranch})
	if r.ToBranch != "" {
		t.AppendRow(table.Row{"branch", r.ToBranch})
	}
	if r.CommitSHA != "" {
		t.AppendRow(table.Row{"commit", r.CommitSHA})
	}
	if r.PullRequest != nil {
		t.AppendRow(table.Row{"pr", r.PullRequest.GetHTMLURL()})
	}
	if r.GoWorkSum != nil {
		t.AppendRow(table.Row{"go.work.sum", r.GoWorkSum.Status()})
	}
	if err != nil {
		t.AppendRow(table.Row{"error", err.Error()})
	}

	return t
}

// pathTable returns a table writer with the per-path sync results.
func (r *SyncGoModRes) pathTable() table.Writer {
	t := r.newResultTable()
	t.AppendHeader(table.Row{"path", "changes", "tidy", "error"})

	if r == nil {
		return t
	}

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

	return t
}

// newResultTable returns a borderless table writer for metadata key/value rows.
func (r *SyncGoModRes) newResultTable() table.Writer {
	t := table.NewWriter()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false
	t.Style().Options.SeparateFooter = false
	t.Style().Options.SeparateHeader = false
	t.Style().Options.SeparateRows = false
	t.SuppressEmptyColumns()
	t.SuppressTrailingSpaces()
	return t
}
