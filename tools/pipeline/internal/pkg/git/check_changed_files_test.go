// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package git

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/changed"
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/config"
	libgit "github.com/hashicorp/vault/tools/pipeline/internal/pkg/git/client"
	"github.com/stretchr/testify/require"
)

// TestCheckChangedFilesReq_validate verifies that every disallowed group has to
// be a changed file group in pipeline.hcl. A misspelled group would never match
// any files, which would let a check that ought to fail quietly pass.
func TestCheckChangedFilesReq_validate(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		disallowedGroups changed.FileGroups
		errHas           string
	}{
		"known groups": {
			disallowedGroups: changed.FileGroups{"changelog", "enterprise"},
		},
		"unknown group": {
			disallowedGroups: changed.FileGroups{"enterprise", "enterprse"},
			errHas:           "disallowed group enterprse is not a changed file group in pipeline.hcl",
		},
		"no groups": {
			errHas: "no disallowed groups or sync policy have been configured",
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := &CheckChangedFilesReq{
				DecodeRes: &config.DecodeRes{
					Config: &config.Config{
						ChangedFiles: &changed.Config{
							Groups: []*changed.GroupConfig{
								{Name: "changelog"},
								{Name: "enterprise"},
							},
						},
					},
				},
				Branch:           "main",
				DisallowedGroups: test.disallowedGroups,
			}

			err := req.validate(t.Context())
			if test.errHas != "" {
				require.ErrorContains(t, err, test.errHas)

				return
			}
			require.NoError(t, err)
		})
	}
}

// TestCheckChangedFilesReq_Run verifies how we check changed files for the
// disallowed groups from --disallowed-groups and a sync policy. Each case gets
// a new repository with three commits on main. The first adds base.txt, the
// second adds the files in the case, and the third removes any files the case
// wants removed. We check the whole history of main unless the case asks for
// the range of the last two commits.
//
// Files in a disallowed group fail the check even if they've since been
// removed, because they're still in the history of the branch. The only way
// past the check is to be listed in the sync policy as exempt, and that only
// works for the exact path.
func TestCheckChangedFilesReq_Run(t *testing.T) {
	t.Parallel()

	const pipelineHCL = `
changed_files {
  group "changelog" {
    match {
      base_dir = ["changelog"]
    }
  }

  group "docs" {
    match {
      extension = [".md"]
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

  policy "exempt-changelog" {
    exempt_from_disallowed_groups = [joinpath("changelog", "1.txt")]
  }

  policy "typo" {
    disallowed_groups = ["enterprse"]
  }
}
`

	for name, test := range map[string]struct {
		add              []string
		remove           []string
		useRange         bool
		disallowedGroups changed.FileGroups
		syncPolicy       string
		wantMatched      []string
		wantMatchedGrps  changed.FileGroups
		wantExempt       []string
		wantChecked      changed.FileGroups
		errHas           string
	}{
		"exempt file in history": {
			add:         []string{".github/instructions/index.md"},
			syncPolicy:  "ce",
			wantExempt:  []string{".github/instructions/index.md"},
			wantChecked: changed.FileGroups{"enterprise"},
		},
		"file that isn't exempt": {
			add:             []string{".github/instructions/index.md", "vault_ent/secret.go"},
			syncPolicy:      "ce",
			wantMatched:     []string{"vault_ent/secret.go"},
			wantMatchedGrps: changed.FileGroups{"enterprise"},
			wantExempt:      []string{".github/instructions/index.md"},
			wantChecked:     changed.FileGroups{"enterprise"},
		},
		"new file next to an exempt file": {
			add:             []string{".github/instructions/index.md", ".github/instructions/new.md"},
			syncPolicy:      "ce",
			wantMatched:     []string{".github/instructions/new.md"},
			wantMatchedGrps: changed.FileGroups{"enterprise"},
			wantExempt:      []string{".github/instructions/index.md"},
			wantChecked:     changed.FileGroups{"enterprise"},
		},
		"removed file is still in history": {
			add:             []string{"vault_ent/secret.go"},
			remove:          []string{"vault_ent/secret.go"},
			syncPolicy:      "ce",
			wantMatched:     []string{"vault_ent/secret.go"},
			wantMatchedGrps: changed.FileGroups{"enterprise"},
			wantChecked:     changed.FileGroups{"enterprise"},
		},
		"disallowed groups are added to the policy groups": {
			add:              []string{".github/instructions/index.md", "changelog/1.txt"},
			disallowedGroups: changed.FileGroups{"changelog"},
			syncPolicy:       "ce",
			wantMatched:      []string{"changelog/1.txt"},
			wantMatchedGrps:  changed.FileGroups{"changelog"},
			wantExempt:       []string{".github/instructions/index.md"},
			wantChecked:      changed.FileGroups{"changelog", "enterprise"},
		},
		"exempt from every disallowed group": {
			add:              []string{".github/instructions/index.md"},
			disallowedGroups: changed.FileGroups{"docs"},
			syncPolicy:       "ce",
			wantExempt:       []string{".github/instructions/index.md"},
			wantChecked:      changed.FileGroups{"docs", "enterprise"},
		},
		"policy with only exempt files": {
			add:              []string{"changelog/1.txt", "changelog/2.txt"},
			disallowedGroups: changed.FileGroups{"changelog"},
			syncPolicy:       "exempt-changelog",
			wantMatched:      []string{"changelog/2.txt"},
			wantMatchedGrps:  changed.FileGroups{"changelog"},
			wantExempt:       []string{"changelog/1.txt"},
			wantChecked:      changed.FileGroups{"changelog"},
		},
		"no sync policy means no exempt files": {
			add:              []string{".github/instructions/index.md"},
			disallowedGroups: changed.FileGroups{"enterprise"},
			wantMatched:      []string{".github/instructions/index.md"},
			wantMatchedGrps:  changed.FileGroups{"enterprise"},
			wantChecked:      changed.FileGroups{"enterprise"},
		},
		"range": {
			add:             []string{".github/instructions/index.md", "vault_ent/secret.go"},
			useRange:        true,
			syncPolicy:      "ce",
			wantMatched:     []string{"vault_ent/secret.go"},
			wantMatchedGrps: changed.FileGroups{"enterprise"},
			wantExempt:      []string{".github/instructions/index.md"},
			wantChecked:     changed.FileGroups{"enterprise"},
		},
		"no disallowed groups or sync policy": {
			add:    []string{"vault_ent/secret.go"},
			errHas: "no disallowed groups or sync policy have been configured",
		},
		"sync policy without disallowed groups": {
			add:        []string{"changelog/1.txt"},
			syncPolicy: "exempt-changelog",
			errHas:     "no disallowed groups have been configured",
		},
		"unknown sync policy": {
			add:        []string{"vault_ent/secret.go"},
			syncPolicy: "community",
			errHas:     "no sync policy named community",
		},
		"unknown disallowed group in the sync policy": {
			add:        []string{"vault_ent/secret.go"},
			syncPolicy: "typo",
			errHas:     "disallowed group enterprse is not a changed file group",
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			dir := t.TempDir()
			repo := libgit.NewClient(libgit.WithDir(dir), libgit.WithConfig(testGitConfig))
			write := func(name string) {
				t.Helper()

				require.NoError(t, os.MkdirAll(filepath.Dir(filepath.Join(dir, name)), 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(dir, name), []byte(name), 0o644))
			}
			commit := func(message string) {
				t.Helper()

				gitRes, err := repo.Add(t.Context(), &libgit.AddOpts{All: true})
				require.NoError(t, err, gitRes.String())
				gitRes, err = repo.Commit(t.Context(), &libgit.CommitOpts{Quiet: true, AllowEmpty: true, Message: message})
				require.NoError(t, err, gitRes.String())
			}

			gitRes, err := repo.Init(t.Context(), &libgit.InitOpts{Quiet: true, InitialBranch: "main"})
			require.NoError(t, err, gitRes.String())
			write("base.txt")
			commit("base")
			for _, name := range test.add {
				write(name)
			}
			commit("add")
			for _, name := range test.remove {
				require.NoError(t, os.Remove(filepath.Join(dir, name)))
			}
			commit("remove")

			cfg, err := config.DecodeBytes([]byte(pipelineHCL))
			require.NoError(t, err)
			req := &CheckChangedFilesReq{
				DecodeRes:        &config.DecodeRes{Config: cfg},
				DisallowedGroups: test.disallowedGroups,
				SyncPolicy:       test.syncPolicy,
			}
			if test.useRange {
				req.Range = "HEAD~2..HEAD"
			} else {
				req.Branch = "main"
			}

			res, err := req.Run(t.Context(), repo)
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
		})
	}
}

// TestCheckChangedFilesRes_ToTable verifies that exempt files show up in the
// table so that we can always see what the sync policy let through. When
// nothing was exempt the cells have to be empty rather than <nil>, otherwise
// the exempt column would never be suppressed.
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
		"exempt only": {
			res: &CheckChangedFilesRes{
				CheckedGroups: changed.FileGroups{"enterprise"},
				ExemptFiles: changed.Files{
					{Filename: ".github/instructions/index.md", Groups: changed.FileGroups{"enterprise"}},
				},
			},
			wantContains: []string{".github/instructions/index.md", "EXEMPT", "yes"},
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
