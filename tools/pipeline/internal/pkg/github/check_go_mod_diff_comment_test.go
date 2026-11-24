// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// Test_CheckGoModDiffRes_CommentOnPullRequest verifies that a single marked
// comment is created, updated in place, or left alone depending on whether
// there is drift and whether a marked comment already exists. Existing
// comments are spread across pages to exercise pagination.
func Test_CheckGoModDiffRes_CommentOnPullRequest(t *testing.T) {
	t.Parallel()

	drift := &CheckGoModDiffRes{Diffs: []*CheckGoModDiff{{
		Path:  "go.mod",
		Err:   errors.New("go.mod is missing"),
		Error: "go.mod is missing",
	}}}
	clean := &CheckGoModDiffRes{}

	for name, test := range map[string]struct {
		res            *CheckGoModDiffRes
		existingOnPage int // 0 for none, otherwise the page holding the marked comment
		expectCreate   bool
		expectUpdate   bool
		expectContains string
	}{
		"drift creates":                    {res: drift, expectCreate: true, expectContains: "go.mod is missing"},
		"drift updates existing":           {res: drift, existingOnPage: 1, expectUpdate: true, expectContains: "go.mod is missing"},
		"drift updates existing on page 2": {res: drift, existingOnPage: 2, expectUpdate: true, expectContains: "go.mod is missing"},
		"no drift without comment":         {res: clean},
		"no drift updates existing":        {res: clean, existingOnPage: 2, expectUpdate: true, expectContains: "No go.mod drift"},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			client, mux, teardown := setupTestClient(t)
			defer teardown()

			var created, updated string
			var updatedID int64
			mux.HandleFunc("/repos/hashicorp/vault-enterprise/issues/42/comments", func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost {
					created = decodeCommentBody(t, r)
					w.WriteHeader(http.StatusCreated)
					_, _ = w.Write([]byte(`{}`))
					return
				}

				page := r.URL.Query().Get("page")
				if page == "" || page == "1" {
					w.Header().Set("Link", `<http://`+r.Host+`/repos/hashicorp/vault-enterprise/issues/42/comments?page=2>; rel="next"`)
					body := "unrelated"
					if test.existingOnPage == 1 {
						body = goModDiffCommentMarker + " old"
					}
					_, _ = w.Write([]byte(`[{"id":1,"body":"` + body + `"}]`))
					return
				}
				body := "also unrelated"
				if test.existingOnPage == 2 {
					body = goModDiffCommentMarker + " old"
				}
				_, _ = w.Write([]byte(`[{"id":2,"body":"` + body + `"}]`))
			})
			mux.HandleFunc("/repos/hashicorp/vault-enterprise/issues/comments/", func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, http.MethodPatch, r.Method)
				updated = decodeCommentBody(t, r)
				updatedID = int64(test.existingOnPage)
				require.True(t, strings.HasSuffix(r.URL.Path, "/"+strconv.Itoa(test.existingOnPage)))
				_, _ = w.Write([]byte(`{}`))
			})

			err := test.res.CommentOnPullRequest(context.Background(), client, "hashicorp", "vault-enterprise", 42, nil)
			require.NoError(t, err)

			require.Equal(t, test.expectCreate, created != "")
			require.Equal(t, test.expectUpdate, updated != "")
			if test.expectUpdate {
				require.Equal(t, int64(test.existingOnPage), updatedID)
			}
			posted := created + updated
			if test.expectContains != "" {
				require.Contains(t, posted, goModDiffCommentMarker)
				require.Contains(t, posted, test.expectContains)
			}
		})
	}
}

// Test_CheckGoModDiffRes_CommentOnPullRequest_APIError verifies API failures
// are returned to the caller, which decides whether they are fatal.
func Test_CheckGoModDiffRes_CommentOnPullRequest_APIError(t *testing.T) {
	t.Parallel()

	client, mux, teardown := setupTestClient(t)
	defer teardown()
	mux.HandleFunc("/repos/hashicorp/vault-enterprise/issues/42/comments", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	})

	res := &CheckGoModDiffRes{Diffs: []*CheckGoModDiff{{Path: "go.mod", Err: errors.New("boom"), Error: "boom"}}}
	require.Error(t, res.CommentOnPullRequest(context.Background(), client, "hashicorp", "vault-enterprise", 42, nil))
}

func decodeCommentBody(t *testing.T, r *http.Request) string {
	t.Helper()

	var body struct {
		Body string `json:"body"`
	}
	require.NoError(t, json.NewDecoder(r.Body).Decode(&body))

	return body.Body
}
