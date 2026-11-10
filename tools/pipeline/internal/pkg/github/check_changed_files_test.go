// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	gh "github.com/google/go-github/v83/github"
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/changed"
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/config"
	"github.com/stretchr/testify/require"
)

// TestCheckChangedFilesReq_Run tests checking the files changed in a pull
// request or commit for disallowed groups from --disallowed-groups and a sync
// policy. Files that the sync policy exempts shouldn't fail the check. Each
// case gets a fake Github API that returns the case's files.
func TestCheckChangedFilesReq_Run(t *testing.T) {
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

  policy "typo" {
    disallowed_groups = ["enterprse"]
  }
}
`

	for name, test := range map[string]struct {
		files            []*gh.CommitFile
		commit           bool
		noTarget         bool
		disallowedGroups changed.FileGroups
		syncPolicy       string
		wantMatched      []string
		wantMatchedGrps  changed.FileGroups
		wantExempt       []string
		wantChecked      changed.FileGroups
		errHas           string
	}{
		"removing exempt files": {
			files: []*gh.CommitFile{
				{Filename: gh.Ptr(".github/instructions/index.md"), Status: gh.Ptr("removed")},
				{Filename: gh.Ptr("README.md"), Status: gh.Ptr("modified")},
			},
			syncPolicy:  "ce",
			wantExempt:  []string{".github/instructions/index.md"},
			wantChecked: changed.FileGroups{"enterprise"},
		},
		"file that isn't exempt": {
			files: []*gh.CommitFile{
				{Filename: gh.Ptr(".github/instructions/index.md"), Status: gh.Ptr("modified")},
				{Filename: gh.Ptr("vault_ent/secret.go"), Status: gh.Ptr("added")},
			},
			syncPolicy:      "ce",
			wantMatched:     []string{"vault_ent/secret.go"},
			wantMatchedGrps: changed.FileGroups{"enterprise"},
			wantExempt:      []string{".github/instructions/index.md"},
			wantChecked:     changed.FileGroups{"enterprise"},
		},
		"new file next to an exempt file": {
			files: []*gh.CommitFile{
				{Filename: gh.Ptr(".github/instructions/new.md"), Status: gh.Ptr("added")},
			},
			syncPolicy:      "ce",
			wantMatched:     []string{".github/instructions/new.md"},
			wantMatchedGrps: changed.FileGroups{"enterprise"},
			wantChecked:     changed.FileGroups{"enterprise"},
		},
		"disallowed groups are added to the policy groups": {
			files: []*gh.CommitFile{
				{Filename: gh.Ptr(".github/instructions/index.md"), Status: gh.Ptr("removed")},
				{Filename: gh.Ptr("changelog/1.txt"), Status: gh.Ptr("added")},
			},
			disallowedGroups: changed.FileGroups{"changelog"},
			syncPolicy:       "ce",
			wantMatched:      []string{"changelog/1.txt"},
			wantMatchedGrps:  changed.FileGroups{"changelog"},
			wantExempt:       []string{".github/instructions/index.md"},
			wantChecked:      changed.FileGroups{"changelog", "enterprise"},
		},
		"no sync policy means no exempt files": {
			files: []*gh.CommitFile{
				{Filename: gh.Ptr(".github/instructions/index.md"), Status: gh.Ptr("removed")},
			},
			disallowedGroups: changed.FileGroups{"enterprise"},
			wantMatched:      []string{".github/instructions/index.md"},
			wantMatchedGrps:  changed.FileGroups{"enterprise"},
			wantChecked:      changed.FileGroups{"enterprise"},
		},
		"commit": {
			files: []*gh.CommitFile{
				{Filename: gh.Ptr(".github/instructions/index.md"), Status: gh.Ptr("removed")},
				{Filename: gh.Ptr("vault_ent/secret.go"), Status: gh.Ptr("added")},
			},
			commit:          true,
			syncPolicy:      "ce",
			wantMatched:     []string{"vault_ent/secret.go"},
			wantMatchedGrps: changed.FileGroups{"enterprise"},
			wantExempt:      []string{".github/instructions/index.md"},
			wantChecked:     changed.FileGroups{"enterprise"},
		},
		"no disallowed groups or sync policy": {
			errHas: "no disallowed groups or sync policy have been configured",
		},
		"unknown sync policy": {
			syncPolicy: "community",
			errHas:     "no sync policy named community",
		},
		"unknown disallowed group in the sync policy": {
			syncPolicy: "typo",
			errHas:     "disallowed group enterprse is not a changed file group",
		},
		"no pull request or commit": {
			noTarget:   true,
			syncPolicy: "ce",
			errHas:     "no pull request number or commit SHA has been provided",
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			mux := http.NewServeMux()
			mux.HandleFunc("GET /repos/hashicorp/vault-enterprise/pulls/1234/files", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(test.files)
			})
			mux.HandleFunc("GET /repos/hashicorp/vault-enterprise/commits/abc123", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(&gh.RepositoryCommit{SHA: gh.Ptr("abc123"), Files: test.files})
			})
			srv := httptest.NewServer(mux)
			t.Cleanup(srv.Close)

			client := gh.NewClient(nil)
			baseURL, err := url.Parse(srv.URL + "/")
			require.NoError(t, err)
			client.BaseURL = baseURL

			cfg, err := config.DecodeBytes([]byte(pipelineHCL))
			require.NoError(t, err)
			req := &CheckChangedFilesReq{
				DecodeRes:        &config.DecodeRes{Config: cfg},
				Owner:            "hashicorp",
				Repo:             "vault-enterprise",
				DisallowedGroups: test.disallowedGroups,
				SyncPolicy:       test.syncPolicy,
			}
			switch {
			case test.noTarget:
			case test.commit:
				req.CommitSHA = "abc123"
			default:
				req.PullNumber = 1234
			}

			res, err := req.Run(t.Context(), client)
			if test.errHas != "" {
				require.ErrorContains(t, err, test.errHas)
				require.Nil(t, res)

				return
			}
			require.NoError(t, err)
			require.Equal(t, test.wantMatched, res.MatchedFiles.Names())
			require.Equal(t, test.wantMatchedGrps, res.MatchedGroups)
			require.Equal(t, test.wantExempt, res.ExemptFiles.Names())
			require.Equal(t, test.wantChecked, res.CheckedGroups)
			if test.syncPolicy == "" {
				require.Nil(t, res.SyncPolicy)
			} else {
				require.Equal(t, test.syncPolicy, res.SyncPolicy.Name)
			}
		})
	}
}

// TestCheckChangedFilesRes_ToTable tests that exempt files show up in the table
// and that we don't render anything if no changed files were disallowed or
// exempt.
func TestCheckChangedFilesRes_ToTable(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		res          *CheckChangedFilesRes
		wantContains []string
		wantMissing  []string
		wantEmpty    bool
	}{
		"nothing matched or exempt": {
			res: &CheckChangedFilesRes{
				CheckedGroups: changed.FileGroups{"enterprise"},
			},
			wantEmpty: true,
		},
		"matched only": {
			res: &CheckChangedFilesRes{
				CheckedGroups: changed.FileGroups{"enterprise"},
				MatchedFiles: changed.Files{
					{Filename: "vault_ent/secret.go", Groups: changed.FileGroups{"enterprise"}},
				},
				MatchedGroups: changed.FileGroups{"enterprise"},
			},
			wantContains: []string{"vault_ent/secret.go"},
			wantMissing:  []string{"EXEMPT", "yes", "<nil>"},
		},
		"matched and exempt": {
			res: &CheckChangedFilesRes{
				CheckedGroups: changed.FileGroups{"enterprise"},
				MatchedFiles: changed.Files{
					{Filename: "vault_ent/secret.go", Groups: changed.FileGroups{"enterprise"}},
				},
				MatchedGroups: changed.FileGroups{"enterprise"},
				ExemptFiles: changed.Files{
					{Filename: ".github/instructions/index.md", Groups: changed.FileGroups{"enterprise"}},
				},
			},
			wantContains: []string{"vault_ent/secret.go", ".github/instructions/index.md", "EXEMPT", "yes"},
			wantMissing:  []string{"<nil>"},
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := test.res.ToTable()
			if test.wantEmpty {
				require.Empty(t, got)

				return
			}
			for _, want := range test.wantContains {
				require.Contains(t, got, want)
			}
			for _, missing := range test.wantMissing {
				require.NotContains(t, got, missing)
			}
		})
	}
}
