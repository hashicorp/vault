// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	libgithub "github.com/google/go-github/v83/github"
	"github.com/stretchr/testify/require"
)

// newTestGithubClientWithTags builds an httptest server that returns the given
// tag names from the GitHub list-tags endpoint, then returns a *libgithub.Client
// pointed at that server.
func newTestGithubClientWithTags(t *testing.T, tags []string) *libgithub.Client {
	t.Helper()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var items []map[string]any
		for _, name := range tags {
			items = append(items, map[string]any{
				"name":   name,
				"commit": map[string]any{"sha": "abc123", "url": ""},
			})
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(items)
	}))
	t.Cleanup(srv.Close)

	client, err := libgithub.NewClient(nil).WithAuthToken("test").
		WithEnterpriseURLs(srv.URL+"/", srv.URL+"/")
	require.NoError(t, err)

	return client
}

// Test_CheckVersionTagsReq_Validate verifies that the request rejects missing
// required fields before making any API calls.
func Test_CheckVersionTagsReq_Validate(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		req     *CheckVersionTagsReq
		wantErr string
	}{
		"nil request": {
			req:     nil,
			wantErr: "uninitialized",
		},
		"missing owner": {
			req:     &CheckVersionTagsReq{Repo: "vault-enterprise", Versions: []string{"1.21.0"}},
			wantErr: "owner is required",
		},
		"missing repo": {
			req:     &CheckVersionTagsReq{Owner: "hashicorp", Versions: []string{"1.21.0"}},
			wantErr: "repo is required",
		},
		"missing versions": {
			req:     &CheckVersionTagsReq{Owner: "hashicorp", Repo: "vault-enterprise"},
			wantErr: "at least one version is required",
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			_, err := tc.req.Run(context.Background(), nil)
			require.Error(t, err)
			require.Contains(t, err.Error(), tc.wantErr)
		})
	}
}

// Test_CheckVersionTagsReq_Run_Rules verifies each of the four tag sequencing
// rules using a controlled set of existing tags served by a local test server.
func Test_CheckVersionTagsReq_Run_Rules(t *testing.T) {
	t.Parallel()

	// Baseline tags present in the repo for all rule tests.
	existingTags := []string{
		"v1.19.12+ent",    // prior patch for rule 4
		"v1.20.0-rc1+ent", // rc1 for rule 1
		"v1.20.0+ent",     // GA for rule 4 baseline
		"v1.20.1+ent",     // prior patch for rule 4
		"v1.21.0-rc1+ent", // rc1 exists → rc2 is valid (rule 2)
	}

	for name, tc := range map[string]struct {
		versions   []string
		tagSuffix  string
		wantPassed bool
		wantRules  map[string]string // version → expected rule on failure
	}{
		"rule 0: tag already exists fails": {
			versions:   []string{"1.20.0"},
			tagSuffix:  "+ent",
			wantPassed: false,
			wantRules:  map[string]string{"1.20.0": "rule-0"},
		},
		"rule 1: GA patch=0 with prior RC passes": {
			versions:   []string{"1.21.0"},
			tagSuffix:  "+ent",
			wantPassed: true,
		},
		"rule 1: GA patch=0 without prior RC fails": {
			versions:   []string{"1.22.0"},
			tagSuffix:  "+ent",
			wantPassed: false,
			wantRules:  map[string]string{"1.22.0": "rule-1"},
		},
		"rule 2: rcN where N>1 with prior rc exists passes": {
			versions:   []string{"1.21.0-rc2"},
			tagSuffix:  "+ent",
			wantPassed: true,
		},
		"rule 2: rcN where N>1 without prior rc fails": {
			versions:   []string{"1.22.0-rc2"},
			tagSuffix:  "+ent",
			wantPassed: false,
			wantRules:  map[string]string{"1.22.0-rc2": "rule-2"},
		},
		"rule 3: rc1 with no prior tags for that major.minor passes": {
			versions:   []string{"1.22.0-rc1"},
			tagSuffix:  "+ent",
			wantPassed: true,
		},
		"rule 3: rc1 where exact tag already exists — fails rule 0": {
			versions:   []string{"1.20.0-rc1"},
			tagSuffix:  "+ent",
			wantPassed: false,
			wantRules:  map[string]string{"1.20.0-rc1": "rule-0"}, // exact tag exists → rule 0 fires first
		},
		"rule 4: minor release with prior patch exists passes": {
			versions:   []string{"1.20.2"},
			tagSuffix:  "+ent",
			wantPassed: true,
		},
		"rule 4: minor release without prior patch fails": {
			versions:   []string{"1.20.5"},
			tagSuffix:  "+ent",
			wantPassed: false,
			wantRules:  map[string]string{"1.20.5": "rule-4"},
		},
		"no suffix: plain version tags pass rule 4": {
			versions:   []string{"1.19.13"},
			tagSuffix:  "",
			wantPassed: false, // v1.19.12 tag has +ent so plain lookup won't find it
			wantRules:  map[string]string{"1.19.13": "rule-4"},
		},
		"multiple versions mixed pass and fail": {
			versions:   []string{"1.20.2", "1.22.0"},
			tagSuffix:  "+ent",
			wantPassed: false,
			wantRules:  map[string]string{"1.22.0": "rule-1"},
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			client := newTestGithubClientWithTags(t, existingTags)
			req := &CheckVersionTagsReq{
				Owner:     "hashicorp",
				Repo:      "vault-enterprise",
				Versions:  tc.versions,
				TagPrefix: "v",
				TagSuffix: tc.tagSuffix,
			}

			res, err := req.Run(context.Background(), client)
			require.NoError(t, err)
			require.Equal(t, tc.wantPassed, res.Passed)

			for _, result := range res.Results {
				if expectedRule, shouldFail := tc.wantRules[result.Version]; shouldFail {
					require.False(t, result.Passed, "expected version %s to fail", result.Version)
					require.Equal(t, expectedRule, result.Rule,
						"version %s: expected rule %s, got %s (%s)", result.Version, expectedRule, result.Rule, result.Detail)
				}
			}
		})
	}
}

// Test_CheckVersionTagsRes_ToJSON verifies the JSON output contains the
// expected top-level keys and per-result fields.
func Test_CheckVersionTagsRes_ToJSON(t *testing.T) {
	t.Parallel()

	res := &CheckVersionTagsRes{
		Passed: false,
		Results: []*VersionTagCheckResult{
			{Version: "1.22.0", Passed: false, Rule: "rule-1", Detail: "no RC tag found for 1.22.0"},
		},
	}

	b, err := res.ToJSON()
	require.NoError(t, err)

	var parsed map[string]any
	require.NoError(t, json.Unmarshal(b, &parsed))
	require.Contains(t, parsed, "passed")
	require.Contains(t, parsed, "results")
}

// Test_CheckVersionTagsRes_ToGithubOutput verifies that failures are written
// in the version_failure_map multiline format expected by workflow steps, and
// that a fully passing result omits the failure map block.
func Test_CheckVersionTagsRes_ToGithubOutput(t *testing.T) {
	t.Parallel()

	t.Run("failures produce version_failure_map block", func(t *testing.T) {
		t.Parallel()

		res := &CheckVersionTagsRes{
			Passed: false,
			Results: []*VersionTagCheckResult{
				{Version: "1.22.0", Passed: false, Rule: "rule-1", Detail: "no RC tag found for 1.22.0"},
				{Version: "1.20.2", Passed: true},
			},
		}

		b, err := res.ToGithubOutput()
		require.NoError(t, err)
		output := string(b)

		require.Contains(t, output, "passed=false\n")
		require.Contains(t, output, "version_failure_map<<EOF\n")
		require.Contains(t, output, "1.22.0:no RC tag found for 1.22.0\n")
		require.NotContains(t, output, "1.20.2")
	})

	t.Run("all passing omits failure map", func(t *testing.T) {
		t.Parallel()

		res := &CheckVersionTagsRes{
			Passed:  true,
			Results: []*VersionTagCheckResult{{Version: "1.20.2", Passed: true}},
		}

		b, err := res.ToGithubOutput()
		require.NoError(t, err)
		output := string(b)

		require.Contains(t, output, "passed=true\n")
		require.NotContains(t, output, "version_failure_map")
	})
}

// Test_parseVersion verifies the internal version parser handles all expected
// formats including edge cases like rc10+.
func Test_parseVersion(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		input     string
		wantMajor int
		wantMinor int
		wantPatch int
		wantRC    string
		wantOK    bool
	}{
		"plain":          {"1.20.4", 1, 20, 4, "", true},
		"rc1":            {"1.21.0-rc1", 1, 21, 0, "-rc1", true},
		"rc10":           {"1.21.0-rc10", 1, 21, 0, "-rc10", true},
		"invalid format": {"not-a-version", 0, 0, 0, "", false},
		"with ent":       {"1.20.4+ent", 0, 0, 0, "", false}, // +ent not handled by parseVersion
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			major, minor, patch, rc, ok := parseVersion(tc.input)
			require.Equal(t, tc.wantOK, ok)
			if tc.wantOK {
				require.Equal(t, tc.wantMajor, major)
				require.Equal(t, tc.wantMinor, minor)
				require.Equal(t, tc.wantPatch, patch)
				require.Equal(t, tc.wantRC, rc)
			}
		})
	}
}

// Test_CheckVersionTagsReq_Run_EntSuffixDerivedFromRepo documents that the
// +ent suffix must be provided on the request when using vault-enterprise
// tags — the pkg layer is repo-agnostic and does not derive it automatically.
// Without the suffix the plain tag v1.19.12 is found, so rule 4 passes.
// With the suffix required (e.g. +ent) the plain tag is ignored and the check fails.
func Test_CheckVersionTagsReq_Run_EntSuffixDerivedFromRepo(t *testing.T) {
	t.Parallel()

	// Only the plain tag exists (no +ent).
	client := newTestGithubClientWithTags(t, []string{"v1.19.12"})

	t.Run("no suffix finds plain tag and passes", func(t *testing.T) {
		t.Parallel()
		req := &CheckVersionTagsReq{
			Owner:     "hashicorp",
			Repo:      "some-oss-project",
			Versions:  []string{"1.19.13"},
			TagPrefix: "v",
			TagSuffix: "",
		}
		res, err := req.Run(context.Background(), client)
		require.NoError(t, err)
		require.True(t, res.Passed)
	})

	t.Run("with +ent suffix plain tag is ignored and rule 4 fails", func(t *testing.T) {
		t.Parallel()
		req := &CheckVersionTagsReq{
			Owner:     "hashicorp",
			Repo:      "vault-enterprise",
			Versions:  []string{"1.19.13"},
			TagPrefix: "v",
			TagSuffix: "+ent",
		}
		res, err := req.Run(context.Background(), client)
		require.NoError(t, err)
		require.False(t, res.Passed)
		require.Equal(t, "rule-4", res.Results[0].Rule)
		require.Contains(t, res.Results[0].Detail, fmt.Sprintf("v1.19.12+ent"))
	})
}
