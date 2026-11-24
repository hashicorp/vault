// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package github_test

import (
	"context"
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

// TestCreateReleaseReq_Run verifies that we create a release with generated
// notes when the tag has none, return the existing release so retries succeed,
// and refuse to release a missing tag since GitHub would create it from the
// default branch.
func TestCreateReleaseReq_Run(t *testing.T) {
	t.Parallel()

	const releaseURL = "https://github.com/o/r/releases/tag/v1.2.3"
	for name, test := range map[string]struct {
		releaseExists bool
		tagMissing    bool
		created       bool
		wantErr       string
	}{
		"missing release":  {created: true},
		"existing release": {releaseExists: true},
		"missing tag":      {tagMissing: true, wantErr: "tag v1.2.3 does not exist"},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var creates atomic.Int32
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.Method + " " + r.URL.Path {
				case "GET /repos/o/r/releases/tags/v1.2.3":
					if !test.releaseExists {
						http.NotFound(w, r)
						return
					}
					_, err := fmt.Fprintf(w, `{"html_url":%q}`, releaseURL)
					assert.NoError(t, err)
				case "GET /repos/o/r/git/ref/tags/v1.2.3":
					if test.tagMissing {
						http.NotFound(w, r)
						return
					}
					_, err := fmt.Fprint(w, `{"ref":"refs/tags/v1.2.3"}`)
					assert.NoError(t, err)
				case "POST /repos/o/r/releases/generate-notes":
					_, err := fmt.Fprint(w, `{"name":"v1.2.3","body":"Generated notes"}`)
					assert.NoError(t, err)
				case "POST /repos/o/r/releases":
					creates.Add(1)
					var release gh.RepositoryRelease
					assert.NoError(t, json.NewDecoder(r.Body).Decode(&release))
					assert.Equal(t, "v1.2.3", release.GetTagName())
					assert.Equal(t, "Generated notes", release.GetBody())
					assert.Equal(t, "legacy", release.GetMakeLatest(), "older release lines must not become latest")
					w.WriteHeader(http.StatusCreated)
					_, err := fmt.Fprintf(w, `{"html_url":%q}`, releaseURL)
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

			req := &github.CreateReleaseReq{Owner: "o", Repo: "r", Tag: "v1.2.3"}
			res, err := req.Run(context.Background(), client)
			if test.wantErr != "" {
				require.ErrorContains(t, err, test.wantErr)
				require.Zero(t, creates.Load(), "we must not create a release for a missing tag")
				return
			}
			require.NoError(t, err)
			want := &github.CreateReleaseRes{Tag: "v1.2.3", URL: releaseURL, Created: test.created}
			require.Equal(t, want, res)
			require.Equal(t, test.created, creates.Load() > 0, "we must only create missing releases")
		})
	}
}
