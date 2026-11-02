// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	libgithub "github.com/google/go-github/v83/github"
	libgit "github.com/hashicorp/vault/tools/pipeline/internal/pkg/git/client"
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/golang"
	"github.com/stretchr/testify/require"
)

// testGitConfig is the git config our test repositories run with. It gives
// commits an identity, and stops the global and system config of whoever runs
// the tests, such as commit signing or a hooks directory, from changing what
// the tests do. It's set with the git client's WithConfig, which has higher
// precedence than any config file.
var testGitConfig = map[string]string{
	"user.name":      "test",
	"user.email":     "test@example.com",
	"commit.gpgsign": "false",
	"core.hooksPath": "/dev/null",
}

// TestGithubSyncGoModReq_validate verifies that validate returns an error for
// each missing required field, accepts valid cross-repo and same-branch-different-repo
// configurations, rejects same-branch/same-repo combinations, and rejects
// explicit go.mod paths together with workspace modules, and rejects an
// explicit to-branch that's longer than GitHub allows.
func TestGithubSyncGoModReq_validate(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		mutate     func(*SyncGoModReq)
		shouldFail bool
		errHas     string
	}{
		"empty a-owner": {
			mutate:     func(r *SyncGoModReq) { r.AOwner = "" },
			shouldFail: true,
			errHas:     "a-owner",
		},
		"empty a-repo": {
			mutate:     func(r *SyncGoModReq) { r.ARepo = "" },
			shouldFail: true,
			errHas:     "a-repo",
		},
		"empty a-branch": {
			mutate:     func(r *SyncGoModReq) { r.ABranch = "" },
			shouldFail: true,
			errHas:     "a-branch",
		},
		"empty b-owner": {
			mutate:     func(r *SyncGoModReq) { r.BOwner = "" },
			shouldFail: true,
			errHas:     "b-owner",
		},
		"empty b-repo": {
			mutate:     func(r *SyncGoModReq) { r.BRepo = "" },
			shouldFail: true,
			errHas:     "b-repo",
		},
		"empty b-branch": {
			mutate:     func(r *SyncGoModReq) { r.BBranch = "" },
			shouldFail: true,
			errHas:     "b-branch",
		},
		"empty paths default to the root go.mod": {
			mutate: func(r *SyncGoModReq) { r.Paths = []string{} },
		},
		"workspace modules": {
			mutate: func(r *SyncGoModReq) {
				r.Paths = nil
				r.WorkspaceModules = true
			},
		},
		"workspace modules with paths": {
			mutate:     func(r *SyncGoModReq) { r.WorkspaceModules = true },
			shouldFail: true,
			errHas:     "cannot be combined with workspace-modules",
		},
		"to-branch at the length limit": {
			mutate: func(r *SyncGoModReq) { r.ToBranch = strings.Repeat("a", maxBranchNameBytes) },
		},
		"to-branch over the length limit": {
			mutate:     func(r *SyncGoModReq) { r.ToBranch = strings.Repeat("a", maxBranchNameBytes+1) },
			shouldFail: true,
			errHas:     "to-branch is 245 bytes long",
		},
		"same branch same repo": {
			mutate: func(r *SyncGoModReq) {
				r.AOwner = "hashicorp"
				r.ARepo = "vault-enterprise"
				r.ABranch = "main"
				r.BOwner = "hashicorp"
				r.BRepo = "vault-enterprise"
				r.BBranch = "main"
			},
			shouldFail: true,
			errHas:     "must differ",
		},
		"same branch different repo": {
			mutate: func(r *SyncGoModReq) {
				r.ABranch = "main"
				r.BBranch = "main"
			},
			shouldFail: false,
		},
		"valid cross-repo": {
			mutate:     func(r *SyncGoModReq) {},
			shouldFail: false,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := &SyncGoModReq{
				AOwner:  "hashicorp",
				ARepo:   "vault",
				ABranch: "main",
				BOwner:  "hashicorp",
				BRepo:   "vault-enterprise",
				BBranch: "ce/main",
				Paths:   []string{"go.mod"},
			}
			test.mutate(req)
			err := req.validate()
			if test.shouldFail {
				require.Error(t, err)
				require.Contains(t, err.Error(), test.errHas)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// TestGithubSyncGoModReq_Validate_Defaults verifies that after a successful
// validate call, TO* fields default to the B-side values, origin names are
// resolved based on whether A and B share the same repo, and that explicitly
// pre-set fields are not overwritten by the defaults.
func TestGithubSyncGoModReq_Validate_Defaults(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		req           *SyncGoModReq
		wantToOwner   string
		wantToRepo    string
		wantPRBase    string
		wantAOrigin   string
		wantBOrigin   string
		wantToOrigin  string
		checkToOwner  bool
		checkToRepo   bool
		checkPRBase   bool
		checkAOrigin  bool
		checkBOrigin  bool
		checkToOrigin bool
	}{
		"defaults applied — different repo": {
			req: &SyncGoModReq{
				AOwner: "hashicorp", ARepo: "vault", ABranch: "main",
				BOwner: "hashicorp", BRepo: "vault-enterprise", BBranch: "ce/main",
				Paths: []string{"go.mod"},
			},
			wantToOwner:  "hashicorp",
			wantToRepo:   "vault-enterprise",
			wantPRBase:   "ce/main",
			checkToOwner: true,
			checkToRepo:  true,
			checkPRBase:  true,
		},
		"origin defaults — same repo": {
			req: &SyncGoModReq{
				AOwner: "hashicorp", ARepo: "vault-enterprise", ABranch: "main",
				BOwner: "hashicorp", BRepo: "vault-enterprise", BBranch: "ce/main",
				Paths: []string{"go.mod"},
			},
			wantAOrigin:   "origin",
			wantBOrigin:   "origin",
			wantToOrigin:  "origin",
			checkAOrigin:  true,
			checkBOrigin:  true,
			checkToOrigin: true,
		},
		"origin defaults — different repo": {
			req: &SyncGoModReq{
				AOwner: "hashicorp", ARepo: "vault", ABranch: "main",
				BOwner: "hashicorp", BRepo: "vault-enterprise", BBranch: "ce/main",
				Paths: []string{"go.mod"},
			},
			wantAOrigin:   "aorigin",
			wantBOrigin:   "borigin",
			wantToOrigin:  "borigin",
			checkAOrigin:  true,
			checkBOrigin:  true,
			checkToOrigin: true,
		},
		"to-owner overridden": {
			req: &SyncGoModReq{
				AOwner: "hashicorp", ARepo: "vault", ABranch: "main",
				BOwner: "hashicorp", BRepo: "vault-enterprise", BBranch: "ce/main",
				ToOwner: "myorg",
				Paths:   []string{"go.mod"},
			},
			wantToOwner:  "myorg",
			checkToOwner: true,
		},
		"to-repo overridden": {
			req: &SyncGoModReq{
				AOwner: "hashicorp", ARepo: "vault", ABranch: "main",
				BOwner: "hashicorp", BRepo: "vault-enterprise", BBranch: "ce/main",
				ToRepo: "my-fork",
				Paths:  []string{"go.mod"},
			},
			wantToRepo:  "my-fork",
			checkToRepo: true,
		},
		"pr-base overridden": {
			req: &SyncGoModReq{
				AOwner: "hashicorp", ARepo: "vault", ABranch: "main",
				BOwner: "hashicorp", BRepo: "vault-enterprise", BBranch: "ce/main",
				PRBase: "custom-base",
				Paths:  []string{"go.mod"},
			},
			wantPRBase:  "custom-base",
			checkPRBase: true,
		},
		"explicit a-origin preserved — same repo": {
			req: &SyncGoModReq{
				AOwner: "hashicorp", ARepo: "vault-enterprise", ABranch: "main",
				BOwner: "hashicorp", BRepo: "vault-enterprise", BBranch: "ce/main",
				AOrigin: "upstream",
				Paths:   []string{"go.mod"},
			},
			wantAOrigin:  "upstream",
			checkAOrigin: true,
		},
		"explicit b-origin preserved — same repo": {
			req: &SyncGoModReq{
				AOwner: "hashicorp", ARepo: "vault-enterprise", ABranch: "main",
				BOwner: "hashicorp", BRepo: "vault-enterprise", BBranch: "ce/main",
				BOrigin: "fork",
				Paths:   []string{"go.mod"},
			},
			wantBOrigin:  "fork",
			checkBOrigin: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			require.NoError(t, test.req.validate())
			if test.checkToOwner {
				require.Equal(t, test.wantToOwner, test.req.ToOwner)
			}
			if test.checkToRepo {
				require.Equal(t, test.wantToRepo, test.req.ToRepo)
			}
			if test.checkPRBase {
				require.Equal(t, test.wantPRBase, test.req.PRBase)
			}
			if test.checkAOrigin {
				require.Equal(t, test.wantAOrigin, test.req.AOrigin)
			}
			if test.checkBOrigin {
				require.Equal(t, test.wantBOrigin, test.req.BOrigin)
			}
			if test.checkToOrigin {
				require.Equal(t, test.wantToOrigin, test.req.ToOrigin)
			}
		})
	}
}

// TestGithubSyncGoModReq_validate_paths verifies that validate syncs the root
// go.mod when no paths are given, keeps the paths we're given, and leaves the
// paths empty in workspace mode, where Run derives them from go.work.
func TestGithubSyncGoModReq_validate_paths(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		paths            []string
		workspaceModules bool
		want             []string
	}{
		"nil paths":            {want: []string{"go.mod"}},
		"empty paths":          {paths: []string{}, want: []string{"go.mod"}},
		"given paths":          {paths: []string{"sdk/go.mod", "api/go.mod"}, want: []string{"sdk/go.mod", "api/go.mod"}},
		"workspace modules":    {workspaceModules: true},
		"workspace empty list": {paths: []string{}, workspaceModules: true, want: []string{}},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := &SyncGoModReq{
				AOwner:           "hashicorp",
				ARepo:            "vault-enterprise",
				ABranch:          "main",
				BOwner:           "hashicorp",
				BRepo:            "vault-enterprise",
				BBranch:          "ce/main",
				Paths:            test.paths,
				WorkspaceModules: test.workspaceModules,
			}
			require.NoError(t, req.validate())
			require.ElementsMatch(t, test.want, req.Paths)
		})
	}
}

// TestGithubSyncGoModReq_branchName verifies that the auto-generated
// intermediate branch name has the go-mod-sync- prefix, puts the ticket after
// it, and replaces the slashes in the A and B branch names. It also verifies
// that the name never goes over GitHub's limit on ref names, including the
// refs/heads/ prefix, that a cut never splits a multi-byte character, and that
// a cut never leaves a separator at the end of the name.
func TestGithubSyncGoModReq_branchName(t *testing.T) {
	t.Parallel()

	const prefix = "go-mod-sync-"

	for name, test := range map[string]struct {
		req  *SyncGoModReq
		want string
	}{
		"no ticket": {
			req:  &SyncGoModReq{ABranch: "main", BBranch: "ce/main"},
			want: "go-mod-sync-main-into-ce-main",
		},
		"with ticket": {
			req:  &SyncGoModReq{ABranch: "main", BBranch: "ce/main", Ticket: "VAULT-1234"},
			want: "go-mod-sync-VAULT-1234-main-into-ce-main",
		},
		"slashes in both branches": {
			req:  &SyncGoModReq{ABranch: "release/1.18.x", BBranch: "ce/release/1.18.x", Ticket: "VAULT-5678"},
			want: "go-mod-sync-VAULT-5678-release-1.18.x-into-ce-release-1.18.x",
		},
		"exactly at the limit": {
			// 12 for the prefix, 6 for -into-, 1 for the B branch.
			req:  &SyncGoModReq{ABranch: strings.Repeat("a", maxBranchNameBytes-len(prefix)-len("-into-")-1), BBranch: "b"},
			want: prefix + strings.Repeat("a", maxBranchNameBytes-len(prefix)-len("-into-")-1) + "-into-b",
		},
		"long branches are capped": {
			req: &SyncGoModReq{
				ABranch: strings.Repeat("a", 150),
				BBranch: strings.Repeat("b", 150),
				Ticket:  "VAULT-1234",
			},
			want: (prefix + "VAULT-1234-" + strings.Repeat("a", 150) + "-into-" + strings.Repeat("b", 150))[:maxBranchNameBytes],
		},
		"a cut never splits a multi-byte character": {
			// Each € is 3 bytes, and the room after the prefix isn't a multiple of 3.
			req:  &SyncGoModReq{ABranch: strings.Repeat("€", 200), BBranch: "main"},
			want: prefix + strings.Repeat("€", (maxBranchNameBytes-len(prefix))/3),
		},
		"a cut drops the separator it leaves at the end": {
			// The slash becomes a dash that lands on the last byte that fits.
			req:  &SyncGoModReq{ABranch: strings.Repeat("a", maxBranchNameBytes-len(prefix)-1) + "/z", BBranch: "main"},
			want: prefix + strings.Repeat("a", maxBranchNameBytes-len(prefix)-1),
		},
		"a cut never leaves a lock suffix": {
			req:  &SyncGoModReq{ABranch: strings.Repeat("a", 227) + ".lock-more", BBranch: "main"},
			want: prefix + strings.Repeat("a", 227),
		},
		"a cut removes repeated lock suffixes": {
			req:  &SyncGoModReq{ABranch: strings.Repeat("a", 222) + ".lock.lock-more", BBranch: "main"},
			want: prefix + strings.Repeat("a", 222),
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := test.req.branchName()
			require.Equal(t, test.want, got)
			require.LessOrEqual(t, len(got), maxBranchNameBytes)
			require.LessOrEqual(t, len("refs/heads/"+got), 255)
			require.True(t, utf8.ValidString(got), "%q must be valid UTF-8", got)
			out, refErr := exec.CommandContext(t.Context(), "git", "check-ref-format", "--branch", got).CombinedOutput()
			require.NoError(t, refErr, "%s", out)
		})
	}
}

// TestGithubSyncGoModReq_resolveToBranch verifies that an explicit --to-branch
// is used as it is, that otherwise the auto-generated name is used, and that a
// generated name that already exists in the clone gets a timestamp suffix
// rather than colliding with it, and that the suffixed name still fits within
// GitHub's limit on branch names even when the original is at the limit.
func TestGithubSyncGoModReq_resolveToBranch(t *testing.T) {
	t.Parallel()

	// A unix timestamp has 10 digits until the year 2286, plus the dash.
	suffixLen := len("-" + strconv.FormatInt(time.Now().Unix(), 10))
	longReq := &SyncGoModReq{ABranch: strings.Repeat("a", 150), BBranch: strings.Repeat("b", 150)}

	for name, test := range map[string]struct {
		req      *SyncGoModReq
		existing string // a branch that already exists in the clone
		want     string // the expected name, or its prefix when suffixed
		suffixed bool
	}{
		"explicit name": {
			req:  &SyncGoModReq{ToBranch: "my-sync", ABranch: "main", BBranch: "ce/main"},
			want: "my-sync",
		},
		"generated name": {
			req:  &SyncGoModReq{ABranch: "main", BBranch: "ce/main"},
			want: "go-mod-sync-main-into-ce-main",
		},
		"generated name with ticket": {
			req:  &SyncGoModReq{ABranch: "main", BBranch: "ce/main", Ticket: "VAULT-1234"},
			want: "go-mod-sync-VAULT-1234-main-into-ce-main",
		},
		"generated name that exists": {
			req:      &SyncGoModReq{ABranch: "main", BBranch: "ce/main"},
			existing: "go-mod-sync-main-into-ce-main",
			want:     "go-mod-sync-main-into-ce-main-",
			suffixed: true,
		},
		"generated name at the limit that exists": {
			// The suffix doesn't fit after a name of the maximum length, so the name
			// is cut to make room and the result still fits.
			req:      longReq,
			existing: longReq.branchName(),
			want:     longReq.branchName()[:maxBranchNameBytes-suffixLen] + "-",
			suffixed: true,
		},
		"explicit name at the limit that exists": {
			req:      &SyncGoModReq{ToBranch: strings.Repeat("x", maxBranchNameBytes), ABranch: "main", BBranch: "ce/main"},
			existing: strings.Repeat("x", maxBranchNameBytes),
			want:     strings.Repeat("x", maxBranchNameBytes-suffixLen) + "-",
			suffixed: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			dir := t.TempDir()
			repo := libgit.NewClient(libgit.WithDir(dir), libgit.WithConfig(testGitConfig))

			gitRes, err := repo.Init(t.Context(), &libgit.InitOpts{Quiet: true, InitialBranch: "main"})
			require.NoError(t, err, gitRes.String())
			gitRes, err = repo.Commit(t.Context(), &libgit.CommitOpts{Quiet: true, AllowEmpty: true, Message: "main"})
			require.NoError(t, err, gitRes.String())
			if test.existing != "" {
				gitRes, err = repo.Branch(t.Context(), &libgit.BranchOpts{BranchName: test.existing})
				require.NoError(t, err, gitRes.String())
			}

			got, err := test.req.resolveToBranch(t.Context(), repo)
			require.NoError(t, err)
			require.LessOrEqual(t, len(got), maxBranchNameBytes)
			if !test.suffixed {
				require.Equal(t, test.want, got)

				return
			}

			require.True(t, strings.HasPrefix(got, test.want), got)
			_, err = strconv.ParseInt(strings.TrimPrefix(got, test.want), 10, 64)
			require.NoError(t, err, "the suffix of %q must be a unix timestamp", got)
		})
	}
}

// TestGithubSyncGoModReq_PRTitle verifies that prTitle returns the correct
// auto-generated pull request title with and without a ticket, and that a
// custom PRTitle override takes precedence over the auto-generated title.
func TestGithubSyncGoModReq_PRTitle(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		req  *SyncGoModReq
		want string
	}{
		"no ticket": {
			req: &SyncGoModReq{
				AOwner: "hashicorp", ARepo: "vault",
				ABranch: "main", BBranch: "ce/main",
			},
			want: "go: sync go.mod from hashicorp/vault/main to ce/main",
		},
		"with ticket": {
			req: &SyncGoModReq{
				AOwner: "hashicorp", ARepo: "vault",
				ABranch: "main", BBranch: "ce/main",
				Ticket: "VAULT-1234",
			},
			want: "[VAULT-1234] go: sync go.mod from hashicorp/vault/main to ce/main",
		},
		"custom title overrides ticket": {
			req: &SyncGoModReq{
				AOwner: "hashicorp", ARepo: "vault",
				ABranch: "main", BBranch: "ce/main",
				Ticket:  "VAULT-1234",
				PRTitle: "custom PR title",
			},
			want: "custom PR title",
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, test.want, test.req.prTitle())
		})
	}
}

// TestGithubSyncGoModRes_ToTable_WithResults verifies that metaTable and
// pathTable produce tables that include A/B metadata and path result rows
// with changes and tidy information.
func TestGithubSyncGoModRes_ToTable_WithResults(t *testing.T) {
	t.Parallel()

	res := &SyncGoModRes{
		AOwner:   "hashicorp",
		ARepo:    "vault",
		ABranch:  "main",
		BOwner:   "hashicorp",
		BRepo:    "vault-enterprise",
		BBranch:  "ce/main",
		ToBranch: "main-into-ce-main",
		Results: []*golang.GoModPathSyncResult{
			{
				SyncModRes: &golang.SyncModRes{
					Path: "go.mod",
					Changes: []*golang.SyncChange{
						{
							Directive:  golang.DirectiveRequire,
							Module:     "github.com/foo/bar",
							OldVersion: "v1.0.0",
							NewVersion: "v1.2.0",
							Action:     golang.SyncActionUpdated,
						},
					},
				},
				TidyRan: true,
			},
		},
	}

	meta := res.metaTable(nil)
	require.NotNil(t, meta)
	metaOut := meta.Render()
	require.Contains(t, metaOut, "hashicorp/vault/main")       // A row
	require.Contains(t, metaOut, "hashicorp/vault-enterprise") // B row
	require.Contains(t, metaOut, "main-into-ce-main")          // branch row

	paths := res.pathTable()
	require.NotNil(t, paths)
	pathsOut := paths.Render()
	require.Contains(t, pathsOut, "go.mod") // path column
	require.Contains(t, pathsOut, "1")      // changes count
	require.Contains(t, pathsOut, "yes")    // tidy ran
}

// TestGithubSyncGoModRes_ToTable_WithErrors verifies that the path table shows
// both sync and tidy errors, and that a non-nil err appears in the meta table.
func TestGithubSyncGoModRes_ToTable_WithErrors(t *testing.T) {
	t.Parallel()

	syncRes := &golang.SyncModRes{
		Path:  "go.mod",
		Error: "sync failed",
	}
	syncRes.Err = errors.New(syncRes.Error)

	res := &SyncGoModRes{
		AOwner:  "hashicorp",
		ARepo:   "vault",
		ABranch: "main",
		BOwner:  "hashicorp",
		BRepo:   "vault-enterprise",
		BBranch: "ce/main",
		Results: []*golang.GoModPathSyncResult{
			{
				SyncModRes: syncRes,
				TidyRan:    true,
				TidyError:  "tidy failed",
			},
		},
	}

	topLevelErr := errors.New("push failed")
	metaOut := res.metaTable(topLevelErr).Render()
	require.Contains(t, metaOut, "push failed") // top-level error row

	pathsOut := res.pathTable().Render()
	require.Contains(t, pathsOut, "go.mod")
	require.Contains(t, pathsOut, "sync failed")
	require.Contains(t, pathsOut, "tidy failed")
}

// TestGithubSyncGoModRes_ToTable_NilPathResult verifies that a nil entry in the
// Results slice is silently skipped and does not panic.
func TestGithubSyncGoModRes_ToTable_NilPathResult(t *testing.T) {
	t.Parallel()

	res := &SyncGoModRes{
		AOwner:  "hashicorp",
		ARepo:   "vault",
		ABranch: "main",
		BOwner:  "hashicorp",
		BRepo:   "vault-enterprise",
		BBranch: "ce/main",
		Results: []*golang.GoModPathSyncResult{
			nil,
			{
				SyncModRes: &golang.SyncModRes{Path: "go.mod"},
			},
		},
	}

	out := res.pathTable().Render()
	require.Contains(t, out, "go.mod")
}

// TestGithubSyncGoModRes_ToTable_NilSyncModRes verifies that a PathResult with
// a nil SyncModRes still renders a path row with zero changes and tidy=no.
func TestGithubSyncGoModRes_ToTable_NilSyncModRes(t *testing.T) {
	t.Parallel()

	res := &SyncGoModRes{
		AOwner:  "hashicorp",
		ARepo:   "vault",
		ABranch: "main",
		BOwner:  "hashicorp",
		BRepo:   "vault-enterprise",
		BBranch: "ce/main",
		Results: []*golang.GoModPathSyncResult{
			{SyncModRes: nil, TidyRan: false},
		},
	}

	out := res.pathTable().Render()
	require.NotEmpty(t, out)
	require.Contains(t, out, "0")  // changes=0
	require.Contains(t, out, "no") // tidy=no
}

// TestGithubSyncGoModRes_ToMarkdown_WithResults verifies that ToMarkdown
// produces pipe-delimited markdown tables for metadata and path results.
func TestGithubSyncGoModRes_ToMarkdown_WithResults(t *testing.T) {
	t.Parallel()

	res := &SyncGoModRes{
		AOwner:  "hashicorp",
		ARepo:   "vault",
		ABranch: "main",
		BOwner:  "hashicorp",
		BRepo:   "vault-enterprise",
		BBranch: "ce/main",
		Results: []*golang.GoModPathSyncResult{
			{
				SyncModRes: &golang.SyncModRes{
					Path: "go.mod",
				},
				TidyRan: false,
			},
		},
	}

	out := res.ToMarkdown(nil)
	require.NotEmpty(t, out)
	require.Contains(t, out, "|")
	require.Contains(t, out, "hashicorp/vault/main")
	require.Contains(t, out, "go.mod")
	require.Contains(t, out, "no") // tidy=no
}

// TestGithubSyncGoModRes_ToJSON_Roundtrip verifies that ToJSON marshals a
// result and the output can be unmarshaled back into the expected fields using
// the new a_/b_ field naming convention.
func TestGithubSyncGoModRes_ToJSON_Roundtrip(t *testing.T) {
	t.Parallel()

	res := &SyncGoModRes{
		AOwner:    "hashicorp",
		ARepo:     "vault",
		ABranch:   "main",
		BOwner:    "hashicorp",
		BRepo:     "vault-enterprise",
		BBranch:   "ce/main",
		ToBranch:  "main-into-ce-main",
		CommitSHA: "abc123def456",
		Results: []*golang.GoModPathSyncResult{
			{
				SyncModRes: &golang.SyncModRes{
					Path: "go.mod",
					Changes: []*golang.SyncChange{
						{
							Directive:  golang.DirectiveRequire,
							Module:     "github.com/foo/bar",
							OldVersion: "v1.0.0",
							NewVersion: "v1.2.0",
							Action:     golang.SyncActionUpdated,
						},
					},
				},
				TidyRan: true,
			},
		},
	}

	b, err := res.ToJSON()
	require.NoError(t, err)
	require.NotEmpty(t, b)

	var parsed map[string]any
	require.NoError(t, json.Unmarshal(b, &parsed))
	require.Equal(t, "hashicorp", parsed["a_owner"])
	require.Equal(t, "vault", parsed["a_repo"])
	require.Equal(t, "main", parsed["a_branch"])
	require.Equal(t, "hashicorp", parsed["b_owner"])
	require.Equal(t, "vault-enterprise", parsed["b_repo"])
	require.Equal(t, "ce/main", parsed["b_branch"])
	require.Equal(t, "main-into-ce-main", parsed["to_branch"])
	require.Equal(t, "abc123def456", parsed["commit_sha"])

	pathResults, ok := parsed["results"].([]any)
	require.True(t, ok, "results must be a JSON array")
	require.Len(t, pathResults, 1)

	pr := pathResults[0].(map[string]any)
	require.Equal(t, "go.mod", pr["path"])
	require.Equal(t, true, pr["tidy_ran"])
}

// TestGithubSyncGoModRes_ToJSON_CommitSHAOmittedWhenEmpty verifies that
// commit_sha is absent from JSON output when the field is empty.
func TestGithubSyncGoModRes_ToJSON_CommitSHAOmittedWhenEmpty(t *testing.T) {
	t.Parallel()

	res := &SyncGoModRes{
		AOwner:   "hashicorp",
		ARepo:    "vault",
		ABranch:  "main",
		BOwner:   "hashicorp",
		BRepo:    "vault-enterprise",
		BBranch:  "ce/main",
		ToBranch: "main-into-ce-main",
	}

	b, err := res.ToJSON()
	require.NoError(t, err)

	var parsed map[string]any
	require.NoError(t, json.Unmarshal(b, &parsed))
	require.Equal(t, "main-into-ce-main", parsed["to_branch"])
	_, hasCommitSHA := parsed["commit_sha"]
	require.False(t, hasCommitSHA, "commit_sha must be absent when empty (omitempty)")
}

// TestGithubSyncGoModRes_ToJSON_WithPullRequestOmitted verifies that
// pull_request is absent from JSON output when PullRequest is nil (omitempty).
func TestGithubSyncGoModRes_ToJSON_WithPullRequestOmitted(t *testing.T) {
	t.Parallel()

	res := &SyncGoModRes{
		AOwner:      "hashicorp",
		ARepo:       "vault",
		ABranch:     "main",
		BOwner:      "hashicorp",
		BRepo:       "vault-enterprise",
		BBranch:     "ce/main",
		ToBranch:    "main-into-ce-main",
		PullRequest: nil,
	}

	b, err := res.ToJSON()
	require.NoError(t, err)

	var parsed map[string]any
	require.NoError(t, json.Unmarshal(b, &parsed))
	_, hasPR := parsed["pull_request"]
	require.False(t, hasPR, "pull_request must be absent when nil (omitempty)")
}

// TestGithubSyncGoModReq_Run_workspaceModules verifies the workspace mode end
// to end against a local origin. Every module in the A branch's go.work is
// synced without a --path for each, every module in the clone's go.work is
// tidied with GOWORK=off, and the changed go.mod and go.sum files are committed
// without raising requirements to workspace-selected versions.
//
// The fixture has these modules:
//   - root requires dep and replaces sdk with the local sdk module. Its dep
//     requirement differs between the branches.
//   - sdk requires dep v1.1.0 on the A branch but a dep version that doesn't
//     exist on the B branch. The root module is synced before sdk and can't be
//     tidied until sdk is synced, which shows that we sync every module before
//     tidying any of them.
//   - tool requires dep v1.0.0 on A and v1.1.0 on B. Its downgrade must survive
//     even though the workspace selects dep v1.1.0.
//   - leaf has no dependencies, so tidying it leaves it without a go.sum, which
//     we can't stage.
//   - extra is in the A branch's go.work but not the B branch's, which shows
//     that the modules to sync come from A rather than from the branch we fork
//     and sync into. It isn't tidied because it isn't in the workspace.
//
// The test serves dep from a file:// module proxy so it needs no network.
//
// Run changes the working directory and the test sets environment variables,
// so it can't run in parallel.
func TestGithubSyncGoModReq_Run_workspaceModules(t *testing.T) {
	root := t.TempDir()

	writeFile := func(name, content string) {
		t.Helper()

		require.NoError(t, os.MkdirAll(filepath.Dir(name), 0o755))
		require.NoError(t, os.WriteFile(name, []byte(content), 0o644))
	}

	// A module proxy with dep v1.0.0 and v1.1.0.
	proxyDir := filepath.Join(root, "proxy")
	for _, version := range []string{"v1.0.0", "v1.1.0"} {
		var zipped bytes.Buffer
		zw := zip.NewWriter(&zipped)
		for name, content := range map[string]string{
			"example.com/dep@" + version + "/go.mod": "module example.com/dep\n\ngo 1.21\n",
			"example.com/dep@" + version + "/dep.go": "package dep\n",
		} {
			w, err := zw.Create(name)
			require.NoError(t, err)
			_, err = w.Write([]byte(content))
			require.NoError(t, err)
		}
		require.NoError(t, zw.Close())

		versionDir := filepath.Join(proxyDir, "example.com/dep/@v")
		writeFile(filepath.Join(versionDir, version+".info"), `{"Version":"`+version+`","Time":"2024-01-01T00:00:00Z"}`)
		writeFile(filepath.Join(versionDir, version+".mod"), "module example.com/dep\n\ngo 1.21\n")
		writeFile(filepath.Join(versionDir, version+".zip"), zipped.String())
	}
	writeFile(filepath.Join(proxyDir, "example.com/dep/@v/list"), "v1.0.0\nv1.1.0\n")

	t.Setenv("GOPROXY", "file://"+filepath.ToSlash(proxyDir))
	t.Setenv("GOMODCACHE", filepath.Join(root, "modcache"))
	t.Setenv("GOFLAGS", "-modcacherw")
	t.Setenv("GOSUMDB", "off")
	t.Setenv("GOWORK", "off")

	// Seed the origin with ce/main and then main on top of it. Main differs in
	// the root, sdk and tool modules' requirements on dep, and in listing extra
	// in its go.work.
	originDir := filepath.Join(root, "origin.git")
	seedDir := filepath.Join(root, "seed")
	require.NoError(t, os.MkdirAll(seedDir, 0o755))
	origin := libgit.NewClient(libgit.WithDir(originDir), libgit.WithConfig(testGitConfig))
	seed := libgit.NewClient(libgit.WithDir(seedDir), libgit.WithConfig(testGitConfig))
	gitRes, err := libgit.NewClient(libgit.WithConfig(testGitConfig)).Init(t.Context(), &libgit.InitOpts{
		Quiet:         true,
		Bare:          true,
		InitialBranch: "main",
		Directory:     originDir,
	})
	require.NoError(t, err, gitRes.String())
	gitRes, err = seed.Init(t.Context(), &libgit.InitOpts{Quiet: true, InitialBranch: "ce/main"})
	require.NoError(t, err, gitRes.String())
	gitRes, err = seed.Remote(t.Context(), &libgit.RemoteOpts{
		Command: libgit.RemoteCommandAdd,
		Name:    "origin",
		URL:     originDir,
	})
	require.NoError(t, err, gitRes.String())

	// The B branch comes first, and then main, the A branch, is made from it.
	writeFile(filepath.Join(seedDir, "go.work"), "go 1.21\n\nuse (\n\t.\n\t./leaf\n\t./sdk\n\t./tool\n)\n")
	writeFile(filepath.Join(seedDir, "go.mod"), "module example.com/root\n\ngo 1.21\n\n"+
		"require (\n\texample.com/dep v1.0.0\n\texample.com/root/sdk v0.0.0\n)\n\n"+
		"replace example.com/root/sdk => ./sdk\n")
	writeFile(filepath.Join(seedDir, "root.go"), "package root\n\nimport (\n\t_ \"example.com/dep\"\n\t_ \"example.com/root/sdk\"\n)\n")
	writeFile(filepath.Join(seedDir, "sdk/go.mod"), "module example.com/root/sdk\n\ngo 1.21\n\nrequire example.com/dep v1.9.9\n")
	writeFile(filepath.Join(seedDir, "sdk/sdk.go"), "package sdk\n\nimport _ \"example.com/dep\"\n")
	writeFile(filepath.Join(seedDir, "tool/go.mod"), "module example.com/root/tool\n\ngo 1.21\n\nrequire example.com/dep v1.1.0\n")
	writeFile(filepath.Join(seedDir, "tool/tool.go"), "package tool\n\nimport _ \"example.com/dep\"\n")
	writeFile(filepath.Join(seedDir, "leaf/go.mod"), "module example.com/root/leaf\n\ngo 1.21\n")
	writeFile(filepath.Join(seedDir, "leaf/leaf.go"), "package leaf\n")
	writeFile(filepath.Join(seedDir, "extra/go.mod"), "module example.com/root/extra\n\ngo 1.21\n")
	writeFile(filepath.Join(seedDir, "extra/extra.go"), "package extra\n")
	gitRes, err = seed.Add(t.Context(), &libgit.AddOpts{All: true})
	require.NoError(t, err, gitRes.String())
	gitRes, err = seed.Commit(t.Context(), &libgit.CommitOpts{Quiet: true, Message: "ce/main"})
	require.NoError(t, err, gitRes.String())
	gitRes, err = seed.Checkout(t.Context(), &libgit.CheckoutOpts{Quiet: true, NewBranch: "main"})
	require.NoError(t, err, gitRes.String())
	writeFile(filepath.Join(seedDir, "go.work"), "go 1.21\n\nuse (\n\t.\n\t./extra\n\t./leaf\n\t./sdk\n\t./tool\n)\n")
	writeFile(filepath.Join(seedDir, "go.mod"), "module example.com/root\n\ngo 1.21\n\n"+
		"require (\n\texample.com/dep v1.1.0\n\texample.com/root/sdk v0.0.0\n)\n\n"+
		"replace example.com/root/sdk => ./sdk\n")
	writeFile(filepath.Join(seedDir, "sdk/go.mod"), "module example.com/root/sdk\n\ngo 1.21\n\nrequire example.com/dep v1.1.0\n")
	writeFile(filepath.Join(seedDir, "tool/go.mod"), "module example.com/root/tool\n\ngo 1.21\n\nrequire example.com/dep v1.0.0\n")
	gitRes, err = seed.Commit(t.Context(), &libgit.CommitOpts{Quiet: true, All: true, Message: "main"})
	require.NoError(t, err, gitRes.String())
	gitRes, err = seed.Push(t.Context(), &libgit.PushOpts{
		Quiet:      true,
		Repository: "origin",
		Refspec:    []string{"main", "ce/main"},
	})
	require.NoError(t, err, gitRes.String())

	// The repo dir is a clone of the origin that we sync into.
	repoDir := filepath.Join(root, "repo")
	gitRes, err = libgit.NewClient(libgit.WithConfig(testGitConfig)).Clone(t.Context(), &libgit.CloneOpts{
		Quiet:      true,
		Repository: originDir,
		Directory:  repoDir,
	})
	require.NoError(t, err, gitRes.String())
	repo := libgit.NewClient(libgit.WithDir(repoDir), libgit.WithConfig(testGitConfig))

	req := &SyncGoModReq{
		AOwner:           "hashicorp",
		ARepo:            "vault-enterprise",
		ABranch:          "main",
		BOwner:           "hashicorp",
		BRepo:            "vault-enterprise",
		BBranch:          "ce/main",
		ToBranch:         "sync-test",
		RepoDir:          repoDir,
		WorkspaceModules: true,
		DiffOpts:         golang.DefaultSyncDiffOpts(),
		Commit:           true,
	}

	res, err := req.Run(t.Context(), nil, libgit.NewClient(libgit.WithConfig(testGitConfig)))
	require.NoError(t, err)

	// Every module in the go.work of main was synced. The root and sdk modules
	// changed, as did tool. We tidied every module in the clone's go.work,
	// including leaf though it didn't change, but not extra.
	require.Len(t, res.Results, 5)
	for i, want := range []struct {
		path    string
		changes int
		tidy    bool
	}{
		{"go.mod", 1, true},
		{"extra/go.mod", 0, false},
		{"leaf/go.mod", 0, true},
		{"sdk/go.mod", 1, true},
		{"tool/go.mod", 1, true},
	} {
		require.Equal(t, want.path, res.Results[i].Path)
		require.Len(t, res.Results[i].Changes, want.changes, want.path)
		require.Equal(t, want.tidy, res.Results[i].TidyRan, want.path)
		require.Empty(t, res.Results[i].TidyError, want.path)
	}

	require.Nil(t, res.GoWorkSum)
	require.NotEmpty(t, res.CommitSHA)

	// Tool keeps its own requirement even though the workspace selects v1.1.0.
	diffRes, err := origin.Diff(t.Context(), &libgit.DiffOpts{
		NameOnly: true,
		Commits:  []string{"ce/main", "sync-test"},
	})
	require.NoError(t, err, diffRes.String())
	require.ElementsMatch(t, []string{
		"go.mod", "go.sum", "sdk/go.mod", "sdk/go.sum", "tool/go.mod", "tool/go.sum",
	}, strings.Fields(string(diffRes.Stdout)))
	for _, mod := range []string{"go.mod", "sdk/go.mod"} {
		showRes, err := origin.Show(t.Context(), &libgit.ShowOpts{Object: "sync-test:" + mod})
		require.NoError(t, err, showRes.String())
		require.Contains(t, string(showRes.Stdout), "example.com/dep v1.1.0", mod)
		require.NotContains(t, string(showRes.Stdout), "example.com/dep v1.0.0", mod)
	}
	showRes, err := origin.Show(t.Context(), &libgit.ShowOpts{Object: "sync-test:tool/go.mod"})
	require.NoError(t, err, showRes.String())
	require.Contains(t, string(showRes.Stdout), "example.com/dep v1.0.0")

	// Everything we changed was committed.
	statusRes, err := repo.Status(t.Context(), &libgit.StatusOpts{Porcelain: true})
	require.NoError(t, err, statusRes.String())
	require.Empty(t, string(statusRes.Stdout))

	// Tidying each module again leaves its requirements unchanged.
	goCmd := func(dir string, env []string, args ...string) {
		t.Helper()

		cmd := exec.Command("go", args...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(), env...)
		out, err := cmd.CombinedOutput()
		require.NoError(t, err, "go %v in %s: %s", args, dir, out)
	}
	for _, mod := range []string{".", "sdk", "tool", "leaf"} {
		goCmd(filepath.Join(repoDir, mod), []string{"GOWORK=off"}, "mod", "tidy")
	}
	statusRes, err = repo.Status(t.Context(), &libgit.StatusOpts{Porcelain: true})
	require.NoError(t, err, statusRes.String())
	require.Empty(t, string(statusRes.Stdout))
}

// TestGithubSyncGoModReq_Run_commit verifies that only successfully tidied changes
// are published and that workspace failures are included in pull request warnings.
// Run changes the working directory, so this test must run sequentially.
func TestGithubSyncGoModReq_Run_commit(t *testing.T) {
	for name, test := range map[string]struct {
		rejectCommit       bool
		malformedWorkspace bool
		snapshotFailure    bool
	}{
		"tidy removes synced changes": {},
		"commit hook rejects changes": {rejectCommit: true},
		"workspace error warns in PR": {malformedWorkspace: true},
		"snapshot error blocks push":  {snapshotFailure: true},
	} {
		t.Run(name, func(t *testing.T) {
			root := t.TempDir()
			seedDir := filepath.Join(root, "seed")
			originDir := filepath.Join(root, "origin.git")
			repoDir := filepath.Join(root, "repo")
			require.NoError(t, os.MkdirAll(seedDir, 0o755))
			seed := libgit.NewClient(libgit.WithDir(seedDir), libgit.WithConfig(testGitConfig))
			gitRes, err := seed.Init(t.Context(), &libgit.InitOpts{InitialBranch: "ce/main"})
			require.NoError(t, err, gitRes.String())
			modPath := filepath.Join(seedDir, "go.mod")
			mod := "module example.com/root\n\ngo 1.21\n\nreplace example.com/unused => ./unused\n"
			require.NoError(t, os.WriteFile(modPath, []byte(mod), 0o644))
			require.NoError(t, os.WriteFile(filepath.Join(seedDir, "go.work"), []byte("go 1.22\n\nuse .\n"), 0o644))
			if test.malformedWorkspace {
				require.NoError(t, os.WriteFile(filepath.Join(seedDir, "go.work"), []byte("use (\n"), 0o644))
			}
			if test.snapshotFailure {
				require.NoError(t, os.MkdirAll(filepath.Join(seedDir, "bad/go.sum"), 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(seedDir, "bad/go.sum/.keep"), nil, 0o644))
				require.NoError(t, os.WriteFile(filepath.Join(seedDir, "bad/go.mod"), []byte("module example.com/bad\n\ngo 1.21\n"), 0o644))
				require.NoError(t, os.WriteFile(filepath.Join(seedDir, "go.work"), []byte("go 1.22\n\nuse (\n .\n ./bad\n)\n"), 0o644))
			}
			require.NoError(t, os.WriteFile(filepath.Join(seedDir, "root.go"), []byte("package root\n"), 0o644))
			gitRes, err = seed.Add(t.Context(), &libgit.AddOpts{All: true})
			require.NoError(t, err, gitRes.String())
			gitRes, err = seed.Commit(t.Context(), &libgit.CommitOpts{Message: "ce/main"})
			require.NoError(t, err, gitRes.String())
			gitRes, err = seed.Checkout(t.Context(), &libgit.CheckoutOpts{NewBranch: "main"})
			require.NoError(t, err, gitRes.String())
			if test.rejectCommit || test.malformedWorkspace {
				mod = strings.Replace(mod, "go 1.21", "go 1.22", 1)
			} else {
				mod += "\nrequire example.com/unused v1.0.0\n"
			}
			require.NoError(t, os.WriteFile(modPath, []byte(mod), 0o644))
			require.NoError(t, os.WriteFile(filepath.Join(seedDir, "go.work"), []byte("go 1.22\n\nuse .\n"), 0o644))
			gitRes, err = seed.Commit(t.Context(), &libgit.CommitOpts{All: true, Message: "main"})
			require.NoError(t, err, gitRes.String())
			gitRes, err = seed.Init(t.Context(), &libgit.InitOpts{Bare: true, Directory: originDir})
			require.NoError(t, err, gitRes.String())
			gitRes, err = seed.Push(t.Context(), &libgit.PushOpts{Repository: originDir, Refspec: []string{"main", "ce/main"}})
			require.NoError(t, err, gitRes.String())
			gitRes, err = seed.Clone(t.Context(), &libgit.CloneOpts{Branch: "main", Repository: originDir, Directory: repoDir})
			require.NoError(t, err, gitRes.String())

			config := map[string]string{}
			for k, v := range testGitConfig {
				config[k] = v
			}
			if test.rejectCommit {
				hooksDir := filepath.Join(root, "hooks")
				require.NoError(t, os.MkdirAll(hooksDir, 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(hooksDir, "pre-commit"),
					[]byte("#!/bin/sh\nprintf 'commit rejected by test hook\\n' >&2\nexit 1\n"), 0o755))
				config["core.hooksPath"] = hooksDir
			}
			req := &SyncGoModReq{
				AOwner: "hashicorp", ARepo: "vault-enterprise", ABranch: "main",
				BOwner: "hashicorp", BRepo: "vault-enterprise", BBranch: "ce/main",
				ToBranch: "sync-test", RepoDir: repoDir, WorkspaceModules: true,
				DiffOpts: golang.DefaultSyncDiffOpts(), Commit: true, PullRequest: true,
			}
			req.DiffOpts.Go = true
			var github *libgithub.Client
			if test.malformedWorkspace {
				client, mux, teardown := setupTestClient(t)
				t.Cleanup(teardown)
				github = client
				mux.HandleFunc("POST /repos/hashicorp/vault-enterprise/pulls", func(w http.ResponseWriter, r *http.Request) {
					var pr libgithub.NewPullRequest
					if decodeErr := json.NewDecoder(r.Body).Decode(&pr); decodeErr != nil {
						http.Error(w, decodeErr.Error(), http.StatusBadRequest)
						return
					}
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusCreated)
					if encodeErr := json.NewEncoder(w).Encode(&libgithub.PullRequest{Number: libgithub.Ptr(1), Body: pr.Body}); encodeErr != nil {
						t.Error(encodeErr)
					}
				})
			}
			res, err := req.Run(t.Context(), github, libgit.NewClient(libgit.WithConfig(config)))
			if test.malformedWorkspace {
				require.ErrorContains(t, err, "listing workspace modules")
				require.NotEmpty(t, res.CommitSHA)
				require.NotNil(t, res.PullRequest)
				require.Contains(t, res.PullRequest.GetBody(), "Warning")
				require.Contains(t, res.PullRequest.GetBody(), "Workspace sync error: listing workspace modules")
				return
			}
			if test.rejectCommit {
				require.ErrorContains(t, err, "commit rejected by test hook")
			} else if test.snapshotFailure {
				require.ErrorContains(t, err, "reading")
				require.Len(t, res.Results, 2)
				for _, result := range res.Results {
					require.False(t, result.TidyRan)
					require.NotEmpty(t, result.TidyError)
				}
			} else {
				require.NoError(t, err)
				statusRes, statusErr := libgit.NewClient(libgit.WithDir(repoDir)).Status(t.Context(), &libgit.StatusOpts{Porcelain: true})
				require.NoError(t, statusErr, statusRes.String())
				require.Empty(t, string(statusRes.Stdout))
			}
			require.NotNil(t, res)
			require.Empty(t, res.CommitSHA)
			require.Nil(t, res.PullRequest)
			origin := libgit.NewClient(libgit.WithDir(originDir))
			_, err = origin.RevParse(t.Context(), &libgit.RevParseOpts{Verify: true, Args: []string{"refs/heads/sync-test"}})
			require.Error(t, err, "a sync without a commit must not push a branch")
		})
	}
}

// TestGithubSyncGoModReq_Run_dirtyRepo verifies that reused clones with staged,
// unstaged, or untracked edits are rejected before pulling or creating a branch.
// Run restores the process working directory, so this test runs sequentially.
func TestGithubSyncGoModReq_Run_dirtyRepo(t *testing.T) {
	for _, state := range []string{"unstaged", "staged", "untracked"} {
		t.Run(state, func(t *testing.T) {
			dir := t.TempDir()
			repo := libgit.NewClient(libgit.WithDir(dir), libgit.WithConfig(testGitConfig))
			gitRes, initErr := repo.Init(t.Context(), &libgit.InitOpts{Quiet: true, InitialBranch: "dest"})
			require.NoError(t, initErr, gitRes.String())
			modPath := filepath.Join(dir, "go.mod")
			require.NoError(t, os.WriteFile(modPath, []byte("module example.com/root\n\ngo 1.21\n"), 0o644))
			gitRes, addErr := repo.Add(t.Context(), &libgit.AddOpts{All: true})
			require.NoError(t, addErr, gitRes.String())
			gitRes, commitErr := repo.Commit(t.Context(), &libgit.CommitOpts{Quiet: true, Message: "dest"})
			require.NoError(t, commitErr, gitRes.String())

			path := modPath
			if state == "untracked" {
				path = filepath.Join(dir, "private.txt")
			}
			require.NoError(t, os.WriteFile(path, []byte("private local edit\n"), 0o644))
			if state == "staged" {
				gitRes, stageErr := repo.Add(t.Context(), &libgit.AddOpts{All: true})
				require.NoError(t, stageErr, gitRes.String())
			}
			before, beforeErr := repo.Status(t.Context(), &libgit.StatusOpts{Porcelain: true, UntrackedFiles: libgit.UntrackedFilesAll})
			require.NoError(t, beforeErr, before.String())
			req := &SyncGoModReq{
				AOwner: "test", ARepo: "test", ABranch: "source",
				BOwner: "test", BRepo: "test", BBranch: "dest",
				RepoDir: dir, WorkspaceModules: true, Commit: true,
			}
			git := libgit.NewClient(libgit.WithConfig(testGitConfig), libgit.WithConfig(map[string]string{
				"status.showUntrackedFiles": "no",
			}))
			res, runErr := req.Run(t.Context(), nil, git)
			require.ErrorContains(t, runErr, "working tree has uncommitted changes")
			require.Nil(t, res)
			after, afterErr := repo.Status(t.Context(), &libgit.StatusOpts{Porcelain: true, UntrackedFiles: libgit.UntrackedFilesAll})
			require.NoError(t, afterErr, after.String())
			require.Equal(t, before.Stdout, after.Stdout)
			branches, branchErr := repo.Branch(t.Context(), &libgit.BranchOpts{Format: "%(refname:short)"})
			require.NoError(t, branchErr, branches.String())
			require.Equal(t, "dest", strings.TrimSpace(string(branches.Stdout)))
		})
	}
}

// TestGithubSyncGoModRes_GoWorkSum verifies that syncing the workspace is
// reported in the metadata table as the state of go.work.sum, only when we
// tried to sync it, and in the JSON output.
func TestGithubSyncGoModRes_GoWorkSum(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		goWorkSum *golang.GoWorkSumSyncResult
		want      string
	}{
		"not synced": {},
		"unchanged":  {goWorkSum: &golang.GoWorkSumSyncResult{}, want: "unchanged"},
		"updated":    {goWorkSum: &golang.GoWorkSumSyncResult{Updated: true}, want: "updated"},
		"failed":     {goWorkSum: &golang.GoWorkSumSyncResult{Error: "download failed"}, want: "failed"},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			res := &SyncGoModRes{
				AOwner:    "hashicorp",
				ARepo:     "vault-enterprise",
				ABranch:   "main",
				BOwner:    "hashicorp",
				BRepo:     "vault-enterprise",
				BBranch:   "ce/main",
				GoWorkSum: test.goWorkSum,
			}

			metaOut := res.metaTable(nil).Render()
			b, err := res.ToJSON()
			require.NoError(t, err)

			var parsed map[string]any
			require.NoError(t, json.Unmarshal(b, &parsed))

			if test.goWorkSum == nil {
				require.NotContains(t, metaOut, "go.work.sum")
				require.NotContains(t, parsed, "go_work_sum")
				return
			}

			require.Contains(t, metaOut, "go.work.sum")
			require.Contains(t, metaOut, test.want)
			require.Equal(t, test.goWorkSum.Updated, parsed["go_work_sum"].(map[string]any)["updated"])
		})
	}
}
