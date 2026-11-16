// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	gh "github.com/google/go-github/v83/github"
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/changed"
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/config"
	"github.com/stretchr/testify/require"
)

// TestListChangedFilesReq_Run_Disallowed tests that we report disallowed and
// exempt files when we're given a sync policy and/or disallowed groups, and
// that we don't report anything without them. Our build workflow fails pull
// requests against ce/* branches if there are any disallowed files so exempt
// files must never show up as disallowed. Each case gets a fake Github API that
// returns the case's files.
func TestListChangedFilesReq_Run_Disallowed(t *testing.T) {
	t.Parallel()

	const pipelineHCL = `
changed_files {
  group "changelog" {
    match {
      base_dir = ["changelog"]
    }
  }

  group "enterprise" {
    match {
      base_dir = ["vault_ent", joinpath(".github", "instructions")]
    }
  }
}

sync {
  policy "ce" {
    disallowed_groups             = ["enterprise"]
    exempt_from_disallowed_groups = [joinpath(".github", "instructions", "index.md")]
  }
}
`

	for name, test := range map[string]struct {
		files                []*gh.CommitFile
		noGroupFiles         bool
		disallowedGroups     changed.FileGroups
		syncPolicy           string
		wantDisallowed       []string
		wantDisallowedGroups changed.FileGroups
		wantExempt           []string
		errHas               string
	}{
		"no sync policy or disallowed groups": {
			files: []*gh.CommitFile{
				{Filename: gh.Ptr(".github/instructions/index.md"), Status: gh.Ptr("removed")},
				{Filename: gh.Ptr("vault_ent/secret.go"), Status: gh.Ptr("added")},
			},
		},
		"sync policy": {
			files: []*gh.CommitFile{
				{Filename: gh.Ptr(".github/instructions/index.md"), Status: gh.Ptr("removed")},
				{Filename: gh.Ptr("vault_ent/secret.go"), Status: gh.Ptr("added")},
				{Filename: gh.Ptr("README.md"), Status: gh.Ptr("modified")},
			},
			syncPolicy:           "ce",
			wantDisallowed:       []string{"vault_ent/secret.go"},
			wantDisallowedGroups: changed.FileGroups{"enterprise"},
			wantExempt:           []string{".github/instructions/index.md"},
		},
		"sync policy with only exempt files": {
			files: []*gh.CommitFile{
				{Filename: gh.Ptr(".github/instructions/index.md"), Status: gh.Ptr("removed")},
			},
			syncPolicy: "ce",
			wantExempt: []string{".github/instructions/index.md"},
		},
		"disallowed groups are added to the policy groups": {
			files: []*gh.CommitFile{
				{Filename: gh.Ptr(".github/instructions/index.md"), Status: gh.Ptr("removed")},
				{Filename: gh.Ptr("changelog/1.txt"), Status: gh.Ptr("added")},
			},
			disallowedGroups:     changed.FileGroups{"changelog"},
			syncPolicy:           "ce",
			wantDisallowed:       []string{"changelog/1.txt"},
			wantDisallowedGroups: changed.FileGroups{"changelog"},
			wantExempt:           []string{".github/instructions/index.md"},
		},
		"disallowed groups without a sync policy": {
			files: []*gh.CommitFile{
				{Filename: gh.Ptr(".github/instructions/index.md"), Status: gh.Ptr("removed")},
			},
			disallowedGroups:     changed.FileGroups{"enterprise"},
			wantDisallowed:       []string{".github/instructions/index.md"},
			wantDisallowedGroups: changed.FileGroups{"enterprise"},
		},
		"sync policy without grouping files": {
			noGroupFiles: true,
			syncPolicy:   "ce",
			errHas:       "changed file grouping must be enabled to check for disallowed groups",
		},
		"unknown sync policy": {
			syncPolicy: "community",
			errHas:     "no sync policy named community",
		},
		"unknown disallowed group": {
			disallowedGroups: changed.FileGroups{"enterprse"},
			errHas:           "disallowed group enterprse is not a changed file group",
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			mux := http.NewServeMux()
			mux.HandleFunc("GET /repos/hashicorp/vault-enterprise/pulls/1234/files", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(test.files)
			})
			srv := httptest.NewServer(mux)
			t.Cleanup(srv.Close)

			client := gh.NewClient(nil)
			baseURL, err := url.Parse(srv.URL + "/")
			require.NoError(t, err)
			client.BaseURL = baseURL

			cfg, err := config.DecodeBytes([]byte(pipelineHCL))
			require.NoError(t, err)
			req := &ListChangedFilesReq{
				DecodeRes:        &config.DecodeRes{Config: cfg},
				Owner:            "hashicorp",
				Repo:             "vault-enterprise",
				PullNumber:       1234,
				GroupFiles:       !test.noGroupFiles,
				DisallowedGroups: test.disallowedGroups,
				SyncPolicy:       test.syncPolicy,
			}

			res, err := req.Run(t.Context(), client)
			if test.errHas != "" {
				require.ErrorContains(t, err, test.errHas)
				require.Nil(t, res)

				return
			}
			require.NoError(t, err)
			require.Len(t, res.Files, len(test.files))
			require.Equal(t, test.wantDisallowed, res.DisallowedFiles.Names())
			require.Equal(t, test.wantDisallowedGroups, res.DisallowedGroups)
			require.Equal(t, test.wantExempt, res.ExemptFiles.Names())
		})
	}
}

// TestListChangedFilesRes_ToGithubOutput verifies the exact JSON that we write
// to $GITHUB_OUTPUT. Our workflows rely on these keys in conditionals so if you
// need to change this test make sure you update them too. We leave out exempt
// files to keep the output small, and the disallowed keys if nothing was
// disallowed.
func TestListChangedFilesRes_ToGithubOutput(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		res  ListChangedFilesRes
		want string
	}{
		"nothing disallowed": {
			res: ListChangedFilesRes{
				Files: changed.Files{
					{Filename: "README.md", Groups: changed.FileGroups{"docs"}},
					{Filename: ".github/instructions/index.md", Groups: changed.FileGroups{"enterprise"}},
				},
				Groups: changed.FileGroups{"docs", "enterprise"},
				ExemptFiles: changed.Files{
					{Filename: ".github/instructions/index.md", Groups: changed.FileGroups{"enterprise"}},
				},
			},
			want: `{
				"files": ["README.md", ".github/instructions/index.md"],
				"groups": ["docs", "enterprise"]
			}`,
		},
		"disallowed": {
			res: ListChangedFilesRes{
				Files: changed.Files{
					{Filename: "vault_ent/secret.go", Groups: changed.FileGroups{"app", "enterprise"}},
					{Filename: ".github/instructions/index.md", Groups: changed.FileGroups{"enterprise"}},
				},
				Groups: changed.FileGroups{"app", "enterprise"},
				DisallowedFiles: changed.Files{
					{Filename: "vault_ent/secret.go", Groups: changed.FileGroups{"app", "enterprise"}},
				},
				DisallowedGroups: changed.FileGroups{"enterprise"},
				ExemptFiles: changed.Files{
					{Filename: ".github/instructions/index.md", Groups: changed.FileGroups{"enterprise"}},
				},
			},
			want: `{
				"files": ["vault_ent/secret.go", ".github/instructions/index.md"],
				"groups": ["app", "enterprise"],
				"disallowed_files": ["vault_ent/secret.go"],
				"disallowed_groups": ["enterprise"]
			}`,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := test.res.ToGithubOutput()
			require.NoError(t, err)
			require.JSONEq(t, test.want, string(got))
		})
	}
}

// TestListChangedFilesRes_ToTable tests that we only add the disallowed column
// if there are disallowed or exempt files.
func TestListChangedFilesRes_ToTable(t *testing.T) {
	t.Parallel()

	files := changed.Files{
		{Filename: "README.md", Groups: changed.FileGroups{"docs"}},
		{Filename: "vault_ent/secret.go", Groups: changed.FileGroups{"enterprise"}},
		{Filename: ".github/instructions/index.md", Groups: changed.FileGroups{"enterprise"}},
	}

	for name, test := range map[string]struct {
		res      *ListChangedFilesRes
		wantRows [][]string
	}{
		"nothing checked": {
			res: &ListChangedFilesRes{Files: files},
			wantRows: [][]string{
				{"PATH", "GROUPS"},
				{"README.md", "docs"},
				{"vault_ent/secret.go", "enterprise"},
				{".github/instructions/index.md", "enterprise"},
			},
		},
		"disallowed and exempt": {
			res: &ListChangedFilesRes{
				Files:           files,
				DisallowedFiles: changed.Files{files[1]},
				ExemptFiles:     changed.Files{files[2]},
			},
			wantRows: [][]string{
				{"PATH", "GROUPS", "DISALLOWED"},
				{"README.md", "docs"},
				{"vault_ent/secret.go", "enterprise", "yes"},
				{".github/instructions/index.md", "enterprise", "exempt"},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			gotRows := [][]string{}
			for _, line := range strings.Split(test.res.ToTable(true), "\n") {
				gotRows = append(gotRows, strings.Fields(line))
			}
			require.Equal(t, test.wantRows, gotRows)
		})
	}
}
