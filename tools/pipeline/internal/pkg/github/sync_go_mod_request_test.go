// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"encoding/json"
	"errors"
	"testing"

	gitpkg "github.com/hashicorp/vault/tools/pipeline/internal/pkg/git"
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/golang"
	"github.com/stretchr/testify/require"
)

// TestGithubSyncGoModReq_validate verifies that validate returns an error for
// each missing required field, accepts valid cross-repo and same-branch-different-repo
// configurations, and rejects same-branch/same-repo combinations.
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
		"empty paths": {
			mutate:     func(r *SyncGoModReq) { r.Paths = []string{} },
			shouldFail: true,
			errHas:     "path",
		},
		"nil paths": {
			mutate:     func(r *SyncGoModReq) { r.Paths = nil },
			shouldFail: true,
			errHas:     "path",
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
		Results: []*gitpkg.GoModPathSyncResult{
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
		Results: []*gitpkg.GoModPathSyncResult{
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
		Results: []*gitpkg.GoModPathSyncResult{
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
		Results: []*gitpkg.GoModPathSyncResult{
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
		Results: []*gitpkg.GoModPathSyncResult{
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
		Results: []*gitpkg.GoModPathSyncResult{
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
