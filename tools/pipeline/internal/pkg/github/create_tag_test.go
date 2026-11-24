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

// TestCreateTagReq_Run verifies that we tag the head of the branch when the tag
// is missing, and leave an existing tag alone so retries succeed, even after
// the branch has moved past it. We refuse to use a tag that isn't on the
// branch, and never move an existing tag.
func TestCreateTagReq_Run(t *testing.T) {
	t.Parallel()

	const (
		head   = "1111111111111111111111111111111111111111"
		tagged = "2222222222222222222222222222222222222222"
	)
	for name, test := range map[string]struct {
		tagSHA  string
		status  string
		wantSHA string
		created bool
		wantErr string
	}{
		"missing tag":               {wantSHA: head, created: true},
		"tag at head":               {tagSHA: head, status: "identical", wantSHA: head},
		"branch moved past the tag": {tagSHA: tagged, status: "ahead", wantSHA: tagged},
		"tag not on branch":         {tagSHA: tagged, status: "diverged", wantErr: "not on branch main"},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var creates atomic.Int32
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.Method + " " + r.URL.Path {
				case "GET /repos/o/r/branches/main":
					_, err := fmt.Fprintf(w, `{"commit":{"sha":%q}}`, head)
					assert.NoError(t, err)
				case "GET /repos/o/r/git/ref/tags/v1.2.3":
					if test.tagSHA == "" {
						http.NotFound(w, r)
						return
					}
					_, err := fmt.Fprint(w, `{"ref":"refs/tags/v1.2.3"}`)
					assert.NoError(t, err)
				case "GET /repos/o/r/commits/refs/tags/v1.2.3":
					_, err := fmt.Fprintf(w, `{"sha":%q}`, test.tagSHA)
					assert.NoError(t, err)
				case "GET /repos/o/r/compare/" + test.tagSHA + "...main":
					_, err := fmt.Fprintf(w, `{"status":%q}`, test.status)
					assert.NoError(t, err)
				case "POST /repos/o/r/git/refs":
					creates.Add(1)
					var ref gh.CreateRef
					assert.NoError(t, json.NewDecoder(r.Body).Decode(&ref))
					assert.Equal(t, gh.CreateRef{Ref: "refs/tags/v1.2.3", SHA: head}, ref)
					w.WriteHeader(http.StatusCreated)
					_, err := fmt.Fprint(w, `{}`)
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

			req := &github.CreateTagReq{Owner: "o", Repo: "r", Tag: "v1.2.3", Branch: "main"}
			res, err := req.Run(context.Background(), client)
			if test.wantErr != "" {
				require.ErrorContains(t, err, test.wantErr)
				require.Zero(t, creates.Load(), "we must never move an existing tag")
				return
			}
			require.NoError(t, err)
			require.Equal(t, &github.CreateTagRes{Tag: "v1.2.3", SHA: test.wantSHA, Created: test.created}, res)
			require.Equal(t, test.created, creates.Load() > 0, "we must only create missing tags")
		})
	}
}
