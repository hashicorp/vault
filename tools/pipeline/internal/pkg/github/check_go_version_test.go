// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	gh "github.com/google/go-github/v83/github"
	"github.com/stretchr/testify/require"
)

// newTestGithubClientForGoVersion builds an httptest server that serves a
// .go-version file with the given content for any branch, or a 404 when
// the branch is not in the map. Returns a GitHub client pointed at the test server.
func newTestGithubClientForGoVersion(t *testing.T, branchContents map[string]string) *gh.Client {
	t.Helper()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		branch := r.URL.Query().Get("ref")
		content, ok := branchContents[branch]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			_, _ = fmt.Fprintf(w, `{"message":"Not Found","documentation_url":""}`)
			return
		}
		// go-github's GetContent() only handles base64 encoding.
		encoded := base64.StdEncoding.EncodeToString([]byte(content))
		resp := gh.RepositoryContent{
			Type:     gh.Ptr("file"),
			Name:     gh.Ptr(".go-version"),
			Encoding: gh.Ptr("base64"),
			Content:  gh.Ptr(encoded),
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	t.Cleanup(srv.Close)

	client, err := gh.NewClient(nil).WithAuthToken("test").
		WithEnterpriseURLs(srv.URL+"/", srv.URL+"/")
	require.NoError(t, err)

	return client
}

// newGoDevServerForCheckGoVersion builds an httptest server that mimics the
// go.dev release API and returns an *http.Client pre-wired to hit it.
func newGoDevServerForCheckGoVersion(t *testing.T, latest, previous string) *http.Client {
	t.Helper()

	releases := []map[string]any{
		{"version": "go" + latest, "stable": true},
		{"version": "go" + previous, "stable": true},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(releases)
	}))
	t.Cleanup(srv.Close)

	return &http.Client{
		Transport: &redirectTransport{base: http.DefaultTransport, target: srv.URL},
	}
}

// Test_CheckGoVersionReq_Validate verifies that missing required fields are
// rejected before any API calls are made.
func Test_CheckGoVersionReq_Validate(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		req     *CheckGoVersionReq
		wantErr string
	}{
		"nil request":      {req: nil, wantErr: "uninitialized"},
		"missing owner":    {req: &CheckGoVersionReq{Repo: "vault-enterprise", Branches: []string{"main"}}, wantErr: "owner is required"},
		"missing repo":     {req: &CheckGoVersionReq{Owner: "hashicorp", Branches: []string{"main"}}, wantErr: "repo is required"},
		"missing branches": {req: &CheckGoVersionReq{Owner: "hashicorp", Repo: "vault-enterprise"}, wantErr: "at least one branch"},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			_, err := tc.req.Run(context.Background(), nil)
			require.Error(t, err)
			require.Contains(t, err.Error(), tc.wantErr)
		})
	}
}

// Test_CheckGoVersionReq_Run verifies status classification and count
// aggregation. A branch using either the latest (N) or previous (N-1) release
// is "up-to-date"; a branch on any other version is "out-of-date"; a branch
// whose .go-version is absent (404) or blank is "file-missing". Every row must
// carry latest_go because release-procedure-ent.yml compares it to go_version.
func Test_CheckGoVersionReq_Run(t *testing.T) {
	t.Parallel()

	const latest = "1.23.4"
	const previous = "1.22.9"

	for name, tc := range map[string]struct {
		files        map[string]string // branch -> .go-version content; absent branch -> 404
		branches     []string
		wantStatuses map[string]string // branch -> status
		wantVersions map[string]string // branch -> go_version
		wantMismatch int
		wantMissing  int
	}{
		"mixed statuses": {
			files: map[string]string{
				"main":               latest,
				"release/1.21.x+ent": previous,
				"release/1.20.x+ent": "1.21.7",
			},
			branches: []string{"main", "release/1.21.x+ent", "release/1.20.x+ent", "release/1.19.x+ent"},
			wantStatuses: map[string]string{
				"main":               "up-to-date",
				"release/1.21.x+ent": "up-to-date",
				"release/1.20.x+ent": "out-of-date",
				"release/1.19.x+ent": "file-missing",
			},
			wantVersions: map[string]string{
				"main":               latest,
				"release/1.21.x+ent": previous,
				"release/1.20.x+ent": "1.21.7",
				"release/1.19.x+ent": "N/A",
			},
			wantMismatch: 1,
			wantMissing:  1,
		},
		"all up to date": {
			files:        map[string]string{"main": latest, "release/1.21.x+ent": previous},
			branches:     []string{"main", "release/1.21.x+ent"},
			wantStatuses: map[string]string{"main": "up-to-date", "release/1.21.x+ent": "up-to-date"},
			wantVersions: map[string]string{"main": latest, "release/1.21.x+ent": previous},
		},
		"whitespace-only file is file-missing": {
			files:        map[string]string{"main": "   \n"},
			branches:     []string{"main"},
			wantStatuses: map[string]string{"main": "file-missing"},
			wantVersions: map[string]string{"main": "N/A"},
			wantMissing:  1,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := &CheckGoVersionReq{
				Owner:      "hashicorp",
				Repo:       "vault-enterprise",
				Branches:   tc.branches,
				HTTPClient: newGoDevServerForCheckGoVersion(t, latest, previous),
			}

			res, err := req.Run(context.Background(), newTestGithubClientForGoVersion(t, tc.files))
			require.NoError(t, err)

			require.Len(t, res.Statuses, len(tc.branches))
			require.Equal(t, tc.wantMismatch, res.MismatchCount)
			require.Equal(t, tc.wantMissing, res.FileMissingCount)

			for _, s := range res.Statuses {
				require.Equal(t, tc.wantStatuses[s.Branch], s.Status, s.Branch)
				require.Equal(t, tc.wantVersions[s.Branch], s.GoVersion, s.Branch)
				require.Equal(t, "1.23.4 / 1.22.9", s.AcceptedVersions, s.Branch)
				require.Equal(t, latest, s.LatestGo, s.Branch)
			}
		})
	}
}

// Test_CheckGoVersionRes_ToGithubOutput verifies that each output is returned
// under the exact key the workflow reads via steps.<id>.outputs.<key>, and that
// the statuses value is a JSON array whose rows include latest_go.
func Test_CheckGoVersionRes_ToGithubOutput(t *testing.T) {
	t.Parallel()

	res := &CheckGoVersionRes{
		Statuses: []*BranchGoVersionStatus{
			{Branch: "main", GoVersion: "1.23.4", Status: "up-to-date", AcceptedVersions: "1.23.4 / 1.22.9", LatestGo: "1.23.4"},
			{Branch: "release/1.20.x+ent", GoVersion: "1.21.7", Status: "out-of-date", AcceptedVersions: "1.23.4 / 1.22.9", LatestGo: "1.23.4"},
		},
		MismatchCount:    1,
		FileMissingCount: 0,
	}

	out, err := res.ToGithubOutput()
	require.NoError(t, err)

	require.Len(t, out, 3)
	require.Equal(t, "1", out["version_mismatch_count"])
	require.Equal(t, "0", out["file_missing_count"])

	var parsed []map[string]any
	require.NoError(t, json.Unmarshal([]byte(out["vault_branch_go_statuses"]), &parsed))
	require.Len(t, parsed, 2)
	require.Equal(t, "1.23.4", parsed[1]["latest_go"])
}
