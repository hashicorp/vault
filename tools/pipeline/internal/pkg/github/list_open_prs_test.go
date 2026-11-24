// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	libgithub "github.com/google/go-github/v83/github"
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/releases"
	"github.com/stretchr/testify/require"
)

// Test_ListOpenPRsReq_resolveBranches verifies that branch resolution uses
// explicit branches when provided, and falls back to deriving active release
// branches from versions.hcl when none are given.
func Test_ListOpenPRsReq_resolveBranches(t *testing.T) {
	t.Parallel()

	activeReleases := &releases.DecodeRes{
		Config: &releases.VersionsConfig{
			ActiveVersion: &releases.ActiveVersion{
				Versions: map[string]*releases.Version{
					"1.19.x": {CEActive: true},
					"1.20.x": {CEActive: true, LTS: true},
				},
			},
		},
	}

	for name, tc := range map[string]struct {
		req            *ListOpenPRsReq
		wantBranches   []string
		wantErrContain string
	}{
		"explicit branches returned as-is": {
			req: &ListOpenPRsReq{
				Branches: []string{"main", "release/1.19.x+ent"},
				Releases: activeReleases,
			},
			wantBranches: []string{"main", "release/1.19.x+ent"},
		},
		"ent repo derives release branches from versions.hcl": {
			req: &ListOpenPRsReq{
				Repo:     "vault-enterprise",
				Releases: activeReleases,
			},
			wantBranches: []string{"release/1.19.x+ent", "release/1.20.x+ent"},
		},
		"ce repo derives release branches from versions.hcl": {
			req: &ListOpenPRsReq{
				Repo:     "vault",
				Releases: activeReleases,
			},
			wantBranches: []string{"release/1.19.x", "release/1.20.x"},
		},
		"nil releases with no branches returns error": {
			req: &ListOpenPRsReq{
				Releases: nil,
			},
			wantErrContain: "no releases configuration available",
		},
		"nil config with no branches returns error": {
			req: &ListOpenPRsReq{
				Releases: &releases.DecodeRes{Config: nil},
			},
			wantErrContain: "no releases configuration available",
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := tc.req.resolveBranches()

			if tc.wantErrContain != "" {
				require.ErrorContains(t, err, tc.wantErrContain,
					"expected resolveBranches to return an error containing %q", tc.wantErrContain)
				return
			}

			require.NoError(t, err)
			// Sort both slices before comparing — map iteration order is non-deterministic.
			require.ElementsMatch(t, tc.wantBranches, got,
				"resolved branches did not match expected set")
		})
	}
}

// Test_ListOpenPRsReq_validate verifies that validate catches missing required
// fields and accepts valid combinations of branches and releases config.
func Test_ListOpenPRsReq_validate(t *testing.T) {
	t.Parallel()

	validReleases := &releases.DecodeRes{
		Config: &releases.VersionsConfig{
			ActiveVersion: &releases.ActiveVersion{
				Versions: map[string]*releases.Version{"1.19.x": {}},
			},
		},
	}

	for name, tc := range map[string]struct {
		req            *ListOpenPRsReq
		wantErrContain string
	}{
		"valid with explicit branches": {
			req: &ListOpenPRsReq{Owner: "hashicorp", Repo: "vault-enterprise", Branches: []string{"main"}},
		},
		"valid with releases and no branches": {
			req: &ListOpenPRsReq{Owner: "hashicorp", Repo: "vault-enterprise", Releases: validReleases},
		},
		"missing owner": {
			req:            &ListOpenPRsReq{Repo: "vault-enterprise", Branches: []string{"main"}},
			wantErrContain: "no github organization",
		},
		"missing repo": {
			req:            &ListOpenPRsReq{Owner: "hashicorp", Branches: []string{"main"}},
			wantErrContain: "no github repository",
		},
		"no branches and no releases": {
			req:            &ListOpenPRsReq{Owner: "hashicorp", Repo: "vault-enterprise"},
			wantErrContain: "either branches or a versions.hcl releases configuration is required",
		},
		"releases with embedded error": {
			req: &ListOpenPRsReq{
				Owner: "hashicorp", Repo: "vault-enterprise",
				Releases: &releases.DecodeRes{Err: fmt.Errorf("parse failure")},
			},
			wantErrContain: "parse failure",
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := tc.req.validate(context.Background())

			if tc.wantErrContain != "" {
				require.ErrorContains(t, err, tc.wantErrContain,
					"expected validate to return an error containing %q", tc.wantErrContain)
				return
			}

			require.NoError(t, err)
		})
	}
}

// Test_ListOpenPRsReq_Run_happyPath verifies that Run correctly maps open PRs
// to their target branches, populating author, URL, milestone, and backport labels.
func Test_ListOpenPRsReq_Run_happyPath(t *testing.T) {
	t.Parallel()

	client, mux, teardown := setupTestClient(t)
	defer teardown()

	mux.HandleFunc("/repos/hashicorp/vault-enterprise/pulls", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "open", r.URL.Query().Get("state"))
		require.Equal(t, "main", r.URL.Query().Get("base"))

		prs := []*libgithub.PullRequest{
			{
				Number:    libgithub.Ptr(101),
				User:      &libgithub.User{Login: libgithub.Ptr("alice")},
				Milestone: &libgithub.Milestone{Title: libgithub.Ptr("v1.20.0")},
				Labels: []*libgithub.Label{
					{Name: libgithub.Ptr("backport/release/1.19.x+ent")},
				},
				Draft: libgithub.Ptr(false),
			},
		}
		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(prs))
	})

	req := &ListOpenPRsReq{
		Owner:    "hashicorp",
		Repo:     "vault-enterprise",
		Branches: []string{"main"},
	}

	res, err := req.Run(context.Background(), client)
	require.NoError(t, err)
	require.Len(t, res.BranchPRs["main"], 1, "expected one PR for branch main")

	pr := res.BranchPRs["main"][0]
	require.Equal(t, "https://github.com/hashicorp/vault-enterprise/pull/101", pr.URL)
	require.Equal(t, "alice", pr.Author)
	require.Equal(t, "v1.20.0", pr.Milestone)
	require.Equal(t, []string{"release/1.19.x+ent"}, pr.Backports)
}

// Test_ListOpenPRsReq_Run_draftFiltering verifies that draft PRs are excluded
// by default and included only when IncludeDrafts is true.
func Test_ListOpenPRsReq_Run_draftFiltering(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		includeDrafts bool
		wantCount     int
	}{
		"drafts excluded by default":    {includeDrafts: false, wantCount: 1},
		"drafts included when flag set": {includeDrafts: true, wantCount: 2},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			client, mux, teardown := setupTestClient(t)
			defer teardown()

			mux.HandleFunc("/repos/hashicorp/vault-enterprise/pulls", func(w http.ResponseWriter, r *http.Request) {
				prs := []*libgithub.PullRequest{
					{Number: libgithub.Ptr(1), User: &libgithub.User{Login: libgithub.Ptr("alice")}, Draft: libgithub.Ptr(false)},
					{Number: libgithub.Ptr(2), User: &libgithub.User{Login: libgithub.Ptr("bob")}, Draft: libgithub.Ptr(true)},
				}
				w.Header().Set("Content-Type", "application/json")
				require.NoError(t, json.NewEncoder(w).Encode(prs))
			})

			req := &ListOpenPRsReq{
				Owner:         "hashicorp",
				Repo:          "vault-enterprise",
				Branches:      []string{"main"},
				IncludeDrafts: tc.includeDrafts,
			}

			res, err := req.Run(context.Background(), client)
			require.NoError(t, err)
			require.Len(t, res.BranchPRs["main"], tc.wantCount,
				"unexpected PR count with includeDrafts=%v", tc.includeDrafts)
		})
	}
}

// Test_ListOpenPRsReq_Run_pagination verifies that Run fetches all pages when
// the GitHub API returns a Link header indicating additional pages exist.
func Test_ListOpenPRsReq_Run_pagination(t *testing.T) {
	t.Parallel()

	client, mux, teardown := setupTestClient(t)
	defer teardown()

	mux.HandleFunc("/repos/hashicorp/vault-enterprise/pulls", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.URL.Query().Get("page") == "2" {
			prs := []*libgithub.PullRequest{
				{Number: libgithub.Ptr(2), User: &libgithub.User{Login: libgithub.Ptr("bob")}, Draft: libgithub.Ptr(false)},
			}
			require.NoError(t, json.NewEncoder(w).Encode(prs))
			return
		}

		// Page 1 — set Link header pointing to page 2.
		serverURL := "http://" + r.Host
		w.Header().Set("Link", fmt.Sprintf(`<%s/repos/hashicorp/vault-enterprise/pulls?page=2>; rel="next"`, serverURL))
		prs := []*libgithub.PullRequest{
			{Number: libgithub.Ptr(1), User: &libgithub.User{Login: libgithub.Ptr("alice")}, Draft: libgithub.Ptr(false)},
		}
		require.NoError(t, json.NewEncoder(w).Encode(prs))
	})

	req := &ListOpenPRsReq{
		Owner:    "hashicorp",
		Repo:     "vault-enterprise",
		Branches: []string{"main"},
	}

	res, err := req.Run(context.Background(), client)
	require.NoError(t, err)
	require.Len(t, res.BranchPRs["main"], 2,
		"expected PRs from both pages to be collected")
}

// Test_ListOpenPRsRes_ToTable verifies that ToTable renders branches in
// alphabetical order and shows "no open PRs" for branches with empty lists.
func Test_ListOpenPRsRes_ToTable(t *testing.T) {
	t.Parallel()

	res := &ListOpenPRsRes{
		Owner: "hashicorp",
		Repo:  "vault-enterprise",
		BranchPRs: map[string][]*OpenPullRequest{
			"release/1.20.x+ent": {
				{URL: "https://github.com/hashicorp/vault-enterprise/pull/200", Author: "carol"},
			},
			"main": {},
			"release/1.19.x+ent": {
				{URL: "https://github.com/hashicorp/vault-enterprise/pull/100", Author: "alice"},
			},
		},
	}

	out := res.ToTable()

	// Verify alphabetical branch ordering: main < release/1.19.x+ent < release/1.20.x+ent
	mainIdx := strings.Index(out, "main")
	r119Idx := strings.Index(out, "release/1.19.x+ent")
	r120Idx := strings.Index(out, "release/1.20.x+ent")
	require.True(t, mainIdx < r119Idx && r119Idx < r120Idx,
		"branches should appear in alphabetical order in table output")

	require.Contains(t, out, "no open PRs",
		"branch with no PRs should show 'no open PRs'")
	require.Contains(t, out, "alice")
	require.Contains(t, out, "carol")
}

// Test_ListOpenPRsRes_ToMarkdown verifies that ToMarkdown renders branches in
// alphabetical order and includes the expected markdown table markers.
func Test_ListOpenPRsRes_ToMarkdown(t *testing.T) {
	t.Parallel()

	res := &ListOpenPRsRes{
		Owner: "hashicorp",
		Repo:  "vault-enterprise",
		BranchPRs: map[string][]*OpenPullRequest{
			"release/1.20.x+ent": {
				{
					URL:       "https://github.com/hashicorp/vault-enterprise/pull/200",
					Author:    "carol",
					Backports: []string{"release/1.19.x+ent"},
				},
			},
			"release/1.19.x+ent": {},
		},
	}

	out := res.ToMarkdown()

	r119Idx := strings.Index(out, "release/1.19.x+ent")
	r120Idx := strings.Index(out, "release/1.20.x+ent")
	require.True(t, r119Idx < r120Idx,
		"branches should appear in alphabetical order in markdown output")

	require.Contains(t, out, "|", "markdown output should contain table separators")
	require.Contains(t, out, "carol")
	require.Contains(t, out, "no open PRs")
}
