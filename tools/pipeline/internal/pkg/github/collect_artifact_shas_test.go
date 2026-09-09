// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestCollectArtifactShasReq_Validate tests validation of the request.
func TestCollectArtifactShasReq_Validate(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		req       *CollectArtifactShasReq
		shouldErr bool
		errMsg    string
	}{
		"valid request": {
			req: &CollectArtifactShasReq{
				PullRequests: []PullRequestEntry{
					{TargetBranch: "release/2.0.x+ent", URL: "https://github.com/hashicorp/vault-enterprise/pull/100"},
				},
				Versions:     []string{"2.0.3+ent"},
				SlackChannel: "C12345",
				DryRunShas:   true,
			},
			shouldErr: false,
		},
		"no pull requests": {
			req: &CollectArtifactShasReq{
				PullRequests: []PullRequestEntry{},
				Versions:     []string{"2.0.3+ent"},
				SlackChannel: "C12345",
			},
			shouldErr: true,
			errMsg:    "no pull requests provided",
		},
		"no versions": {
			req: &CollectArtifactShasReq{
				PullRequests: []PullRequestEntry{
					{TargetBranch: "release/2.0.x+ent", URL: "https://github.com/hashicorp/vault-enterprise/pull/100"},
				},
				Versions:     []string{},
				SlackChannel: "C12345",
			},
			shouldErr: true,
			errMsg:    "no versions provided",
		},
		"no slack channel": {
			req: &CollectArtifactShasReq{
				PullRequests: []PullRequestEntry{
					{TargetBranch: "release/2.0.x+ent", URL: "https://github.com/hashicorp/vault-enterprise/pull/100"},
				},
				Versions:     []string{"2.0.3+ent"},
				SlackChannel: "",
			},
			shouldErr: true,
			errMsg:    "no slack channel provided",
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			err := test.req.validate()
			if test.shouldErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), test.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// TestMapBranchesToVersions tests the branch-to-version mapping logic.
func TestMapBranchesToVersions(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	for name, test := range map[string]struct {
		branches    []string
		versions    []string
		shouldErr   bool
		errMsg      string
		expectedMap map[string]string
	}{
		"plain patch version ent": {
			branches: []string{"release/2.0.x+ent"},
			versions: []string{"2.0.3+ent"},
			expectedMap: map[string]string{
				"release/2.0.x+ent": "2.0.3+ent",
			},
		},
		"rc version maps to same x branch": {
			// 2.1.0-rc1 lives on release/2.1.x — the core invariant
			branches: []string{"release/2.1.x+ent"},
			versions: []string{"2.1.0-rc1+ent"},
			expectedMap: map[string]string{
				"release/2.1.x+ent": "2.1.0-rc1+ent",
			},
		},
		"ce branch matches ce version": {
			branches: []string{"release/2.0.x"},
			versions: []string{"2.0.3"},
			expectedMap: map[string]string{
				"release/2.0.x": "2.0.3",
			},
		},
		"multiple branches map to correct versions from mixed list": {
			branches: []string{"release/2.0.x+ent", "release/2.1.x+ent"},
			versions: []string{"2.0.3+ent", "2.1.0-rc1+ent", "2.0.3"},
			expectedMap: map[string]string{
				"release/2.0.x+ent": "2.0.3+ent",
				"release/2.1.x+ent": "2.1.0-rc1+ent",
			},
		},
		"major minor does not cross match": {
			// release/2.0.x must not match 2.1.3+ent
			branches:  []string{"release/2.0.x+ent"},
			versions:  []string{"2.1.3+ent"},
			shouldErr: true,
			errMsg:    "could not find matching version for branch",
		},
		"no matching version for branch": {
			branches:  []string{"release/2.0.x+ent"},
			versions:  []string{"2.1.3+ent"},
			shouldErr: true,
			errMsg:    `"release/2.0.x+ent"`,
		},
		"invalid branch format": {
			branches:  []string{"main"},
			versions:  []string{"2.0.3+ent"},
			shouldErr: true,
			errMsg:    `could not extract version pattern from branch "main"`,
		},
		"invalid branch format feature branch": {
			branches:  []string{"feature/my-feature"},
			versions:  []string{"2.0.3+ent"},
			shouldErr: true,
			errMsg:    "could not extract version pattern from branch",
		},
		"malformed version with no second dot is skipped and branch goes unmatched": {
			// A version like "2.1" has no second dot — the secondDot guard skips it,
			// so the branch never matches and mapBranchesToVersions returns an error.
			branches:  []string{"release/2.1.x+ent"},
			versions:  []string{"2.1"},
			shouldErr: true,
			errMsg:    "could not find matching version for branch",
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			result, err := mapBranchesToVersions(ctx, test.branches, test.versions)
			if test.shouldErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), test.errMsg)
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.expectedMap, result)
			}
		})
	}
}

// TestCollectArtifactShasReq_Run_DryRun tests the dry-run path of Run, which
// exercises all branching logic without any GitHub API calls.
func TestCollectArtifactShasReq_Run_DryRun(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	for name, test := range map[string]struct {
		req               *CollectArtifactShasReq
		shouldErr         bool
		errMsg            string
		expectedVersions  []string
		expectedSHAPrefix string
		expectTestMode    bool
	}{
		"single ent branch plain patch": {
			req: &CollectArtifactShasReq{
				PullRequests: []PullRequestEntry{
					{TargetBranch: "release/2.0.x+ent", URL: "https://github.com/hashicorp/vault-enterprise/pull/100"},
				},
				Versions:     []string{"2.0.3+ent"},
				SlackChannel: "C12345",
				DryRunShas:   true,
			},
			expectedVersions:  []string{"2.0.3+ent"},
			expectedSHAPrefix: "fakeshaforbranch_",
			expectTestMode:    true,
		},
		"single ent branch rc version": {
			req: &CollectArtifactShasReq{
				PullRequests: []PullRequestEntry{
					{TargetBranch: "release/2.1.x+ent", URL: "https://github.com/hashicorp/vault-enterprise/pull/200"},
				},
				Versions:     []string{"2.1.0-rc1+ent"},
				SlackChannel: "C12345",
				DryRunShas:   true,
			},
			expectedVersions:  []string{"2.1.0-rc1+ent"},
			expectedSHAPrefix: "fakeshaforbranch_",
			expectTestMode:    true,
		},
		"output order follows versions slice not pull requests order": {
			req: &CollectArtifactShasReq{
				PullRequests: []PullRequestEntry{
					{TargetBranch: "release/2.1.x+ent", URL: "https://github.com/hashicorp/vault-enterprise/pull/200"},
					{TargetBranch: "release/2.0.x+ent", URL: "https://github.com/hashicorp/vault-enterprise/pull/100"},
				},
				// Versions deliberately in ascending order, opposite of PR order above
				Versions:     []string{"2.0.3+ent", "2.1.0-rc1+ent"},
				SlackChannel: "C12345",
				DryRunShas:   true,
			},
			expectedVersions: []string{"2.0.3+ent", "2.1.0-rc1+ent"},
			expectTestMode:   true,
		},
		"validation fails missing slack channel": {
			req: &CollectArtifactShasReq{
				PullRequests: []PullRequestEntry{
					{TargetBranch: "release/2.0.x+ent", URL: "https://github.com/hashicorp/vault-enterprise/pull/100"},
				},
				Versions:   []string{"2.0.3+ent"},
				DryRunShas: true,
			},
			shouldErr: true,
			errMsg:    "no slack channel provided",
		},
		"versions and branches are disjoint returns error": {
			req: &CollectArtifactShasReq{
				PullRequests: []PullRequestEntry{
					{TargetBranch: "release/2.0.x+ent", URL: "https://github.com/hashicorp/vault-enterprise/pull/100"},
				},
				// Version belongs to a different minor — no overlap
				Versions:     []string{"2.1.0+ent"},
				SlackChannel: "C12345",
				DryRunShas:   true,
			},
			shouldErr: true,
			errMsg:    "could not find matching version",
		},
		"slack message contains success footer on successful run": {
			req: &CollectArtifactShasReq{
				PullRequests: []PullRequestEntry{
					{TargetBranch: "release/2.0.x+ent", URL: "https://github.com/hashicorp/vault-enterprise/pull/100"},
				},
				Versions:     []string{"2.0.3+ent"},
				SlackChannel: "C12345",
				DryRunShas:   true,
			},
			expectedVersions: []string{"2.0.3+ent"},
			expectTestMode:   true,
		},
		"duplicate target branch last write wins": {
			// Two PRs targeting the same branch — second SHA overwrites first.
			// This documents the accepted behaviour: one SHAMap entry is produced.
			req: &CollectArtifactShasReq{
				PullRequests: []PullRequestEntry{
					{TargetBranch: "release/2.0.x+ent", URL: "https://github.com/hashicorp/vault-enterprise/pull/100"},
					{TargetBranch: "release/2.0.x+ent", URL: "https://github.com/hashicorp/vault-enterprise/pull/101"},
				},
				Versions:     []string{"2.0.3+ent"},
				SlackChannel: "C12345",
				DryRunShas:   true,
			},
			expectedVersions: []string{"2.0.3+ent"},
			expectTestMode:   true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res, err := test.req.Run(ctx, nil) // nil client is safe in dry-run
			if test.shouldErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), test.errMsg)
				require.Nil(t, res)
			} else {
				require.NoError(t, err)
				require.NotNil(t, res)

				// SHAMap order must match Versions order
				require.Len(t, res.SHAMap, len(test.expectedVersions))
				for i, version := range test.expectedVersions {
					require.Equal(t, version, res.SHAMap[i].Version)
					if test.expectedSHAPrefix != "" {
						require.True(t, strings.HasPrefix(res.SHAMap[i].SHA, test.expectedSHAPrefix),
							"expected SHA to start with %q, got %q", test.expectedSHAPrefix, res.SHAMap[i].SHA)
					}
				}

				// Slack message header reflects mode
				if test.expectTestMode {
					require.Contains(t, res.SlackMsg, "(Test Mode)")
				} else {
					require.NotContains(t, res.SlackMsg, "(Test Mode)")
				}

				// Slack message must contain every version
				for _, version := range test.expectedVersions {
					require.Contains(t, res.SlackMsg, version)
				}

				// Success footer must always be present on a successful run
				require.Contains(t, res.SlackMsg, "PRs and SHAs pulled successfully")
			}
		})
	}
}

// TestCollectArtifactShasReq_Run_MergeCommitRetry tests the merge commit SHA
// retry logic using a mock HTTP server, with the retry wait zeroed out so the
// suite stays fast.
func TestCollectArtifactShasReq_Run_MergeCommitRetry(t *testing.T) {
	t.Parallel()

	// Zero out the wait so tests don't actually sleep.
	original := mergeCommitRetryWait
	mergeCommitRetryWait = 0
	t.Cleanup(func() { mergeCommitRetryWait = original })

	ctx := context.Background()

	for name, test := range map[string]struct {
		// responses is the ordered list of merge_commit_sha values the mock
		// server will return on successive GET /repos/.../pulls/100 calls.
		responses []string
		shouldErr bool
		errMsg    string
		expectSHA string
	}{
		"sha available on first attempt": {
			responses: []string{"abc123merged"},
			expectSHA: "abc123merged",
		},
		"sha empty on first attempt then available on retry": {
			responses: []string{"", "abc123merged"},
			expectSHA: "abc123merged",
		},
		"sha empty on all attempts exhausts retries": {
			responses: []string{"", "", ""},
			shouldErr: true,
			errMsg:    "merge commit SHA not available for PR #100",
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			callCount := 0
			client, mux, teardown := setupTestClient(t)
			defer teardown()

			mux.HandleFunc("/repos/hashicorp/vault-enterprise/pulls/100", func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, http.MethodGet, r.Method)
				sha := ""
				if callCount < len(test.responses) {
					sha = test.responses[callCount]
				}
				callCount++
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				fmt.Fprintf(w, `{"number":100,"merge_commit_sha":%q,"head":{"sha":"headsha"}}`, sha)
			})

			req := &CollectArtifactShasReq{
				PullRequests: []PullRequestEntry{
					{TargetBranch: "release/2.0.x+ent", URL: "https://github.com/hashicorp/vault-enterprise/pull/100"},
				},
				Versions:     []string{"2.0.3+ent"},
				SlackChannel: "C12345",
				DryRunShas:   false,
			}

			res, err := req.Run(ctx, client)
			if test.shouldErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), test.errMsg)
				require.Nil(t, res)
			} else {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.Len(t, res.SHAMap, 1)
				require.Equal(t, test.expectSHA, res.SHAMap[0].SHA)
			}
		})
	}
}

// TestRepoForBranch tests that the correct repo constant is selected based on branch suffix.
func TestRepoForBranch(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		branch   string
		expected string
	}{
		"ent branch returns ent repo": {
			branch:   "release/2.0.x+ent",
			expected: ENTRepo,
		},
		"ce branch returns ce repo": {
			branch:   "release/2.0.x",
			expected: CERepo,
		},
		"main branch returns ce repo": {
			branch:   "main",
			expected: CERepo,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, test.expected, repoForBranch(test.branch))
		})
	}
}

// TestParsePRNumber tests extraction of the PR number from a GitHub URL.
func TestParsePRNumber(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		url       string
		expected  int
		shouldErr bool
	}{
		"valid url": {
			url:      "https://github.com/hashicorp/vault-enterprise/pull/13558",
			expected: 13558,
		},
		"non-numeric suffix": {
			url:       "https://github.com/hashicorp/vault-enterprise/pull/abc",
			shouldErr: true,
		},
		"empty string": {
			url:       "",
			shouldErr: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			n, err := parsePRNumber(test.url)
			if test.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.expected, n)
			}
		})
	}
}

// TestSplitOwnerRepo tests splitting an "owner/repo" string.
func TestSplitOwnerRepo(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		input         string
		expectedOwner string
		expectedRepo  string
		shouldErr     bool
	}{
		"valid": {
			input:         "hashicorp/vault-enterprise",
			expectedOwner: "hashicorp",
			expectedRepo:  "vault-enterprise",
		},
		"no slash": {
			input:     "hashicorp",
			shouldErr: true,
		},
		"empty owner": {
			input:     "/vault",
			shouldErr: true,
		},
		"empty repo": {
			input:     "hashicorp/",
			shouldErr: true,
		},
		"empty string": {
			input:     "",
			shouldErr: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			owner, repo, err := splitOwnerRepo(test.input)
			if test.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.expectedOwner, owner)
				require.Equal(t, test.expectedRepo, repo)
			}
		})
	}
}
