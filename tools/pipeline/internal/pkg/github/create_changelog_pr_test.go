// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package github_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync/atomic"
	"testing"

	gh "github.com/google/go-github/v83/github"
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCreateChangelogPRReq_Run verifies that we open a pull request that adds
// the tag's generated notes to the changelog, return the open pull request so
// retries succeed, do nothing once the entry has merged into the base branch,
// and fail on an unknown strategy before we change anything.
func TestCreateChangelogPRReq_Run(t *testing.T) {
	t.Parallel()

	const (
		baseSHA = "1111111111111111111111111111111111111111"
		branch  = "changelog/v1.2.3"
		pullURL = "https://github.com/o/r/pull/1"
	)
	for name, test := range map[string]struct {
		strategy github.ChangelogStrategy
		merged   bool
		open     bool
		want     *github.CreateChangelogPRRes
		err      string
		writes   int32
	}{
		"missing pull request": {
			strategy: github.ChangelogStrategyGithubReleaseNotes,
			want:     &github.CreateChangelogPRRes{Tag: "v1.2.3", URL: pullURL, Created: true},
			writes:   3,
		},
		"open pull request": {
			strategy: github.ChangelogStrategyGithubReleaseNotes,
			open:     true,
			want:     &github.CreateChangelogPRRes{Tag: "v1.2.3", URL: pullURL},
		},
		"merged entry": {
			strategy: github.ChangelogStrategyGithubReleaseNotes,
			merged:   true,
			want:     &github.CreateChangelogPRRes{Tag: "v1.2.3"},
		},
		"unknown strategy": {
			strategy: "unknown",
			err:      `unsupported changelog strategy "unknown"`,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// A retry finds the branch with our entry and the open pull request.
			// Once the pull request merges, the base branch has our entry.
			hasEntry := map[string]bool{baseSHA: test.merged, branch: test.open}
			branchStatus, pulls := http.StatusNotFound, `[]`
			if test.open {
				branchStatus, pulls = http.StatusOK, fmt.Sprintf(`[{"html_url":%q}]`, pullURL)
			}
			var writes atomic.Int32
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.Method + " " + r.URL.Path {
				case "GET /repos/o/r/branches/main":
					_, err := fmt.Fprintf(w, `{"commit":{"sha":%q}}`, baseSHA)
					assert.NoError(t, err)
				case "GET /repos/o/r/contents/CHANGELOG.md":
					content := "# Changelog\n"
					if hasEntry[r.URL.Query().Get("ref")] {
						content = "## v1.2.3\n\nGenerated notes\n\n" + content
					}
					assert.NoError(t, json.NewEncoder(w).Encode(map[string]string{
						"type":     "file",
						"sha":      "blob",
						"encoding": "base64",
						"content":  base64.StdEncoding.EncodeToString([]byte(content)),
					}))
				case "GET /repos/o/r/git/ref/tags/v1.2.3":
					_, err := fmt.Fprint(w, `{"ref":"refs/tags/v1.2.3"}`)
					assert.NoError(t, err)
				case "POST /repos/o/r/releases/generate-notes":
					_, err := fmt.Fprint(w, `{"name":"v1.2.3","body":"Generated notes"}`)
					assert.NoError(t, err)
				case "GET /repos/o/r/git/ref/heads/" + branch:
					w.WriteHeader(branchStatus)
					_, err := fmt.Fprint(w, `{}`)
					assert.NoError(t, err)
				case "POST /repos/o/r/git/refs":
					writes.Add(1)
					var ref gh.CreateRef
					assert.NoError(t, json.NewDecoder(r.Body).Decode(&ref))
					assert.Equal(t, gh.CreateRef{Ref: "refs/heads/" + branch, SHA: baseSHA}, ref)
					w.WriteHeader(http.StatusCreated)
					_, err := fmt.Fprint(w, `{}`)
					assert.NoError(t, err)
				case "PUT /repos/o/r/contents/CHANGELOG.md":
					writes.Add(1)
					var file gh.RepositoryContentFileOptions
					assert.NoError(t, json.NewDecoder(r.Body).Decode(&file))
					assert.Equal(t, branch, file.GetBranch())
					assert.Equal(t, "## v1.2.3\n\nGenerated notes\n\n# Changelog\n", string(file.Content))
					_, err := fmt.Fprint(w, `{}`)
					assert.NoError(t, err)
				case "GET /repos/o/r/pulls":
					assert.Equal(t, "o:"+branch, r.URL.Query().Get("head"))
					_, err := fmt.Fprint(w, pulls)
					assert.NoError(t, err)
				case "POST /repos/o/r/pulls":
					writes.Add(1)
					w.WriteHeader(http.StatusCreated)
					_, err := fmt.Fprintf(w, `{"number":1,"html_url":%q}`, pullURL)
					assert.NoError(t, err)
				default:
					t.Errorf("unexpected request: %s %s", r.Method, r.URL)
					http.NotFound(w, r)
				}
			}))
			t.Cleanup(srv.Close)

			client := gh.NewClient(srv.Client())
			var err error
			client.BaseURL, err = url.Parse(srv.URL + "/")
			require.NoError(t, err)

			req := &github.CreateChangelogPRReq{
				Owner:      "o",
				Repo:       "r",
				Tag:        "v1.2.3",
				Branch:     branch,
				BaseBranch: "main",
				Path:       "CHANGELOG.md",
				Strategy:   test.strategy,
			}
			res, err := req.Run(context.Background(), client)
			if test.err != "" {
				require.ErrorContains(t, err, test.err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, test.want, res)
			require.Equal(t, test.writes, writes.Load(), "we must only create what is missing")
		})
	}
}
