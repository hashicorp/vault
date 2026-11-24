// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/config"
	gitpkg "github.com/hashicorp/vault/tools/pipeline/internal/pkg/git"
	libgit "github.com/hashicorp/vault/tools/pipeline/internal/pkg/git/client"
	"github.com/stretchr/testify/require"
)

// TestGitCheckChangedFilesCmd_GithubOutput verifies that we write to
// $GITHUB_OUTPUT and still fail if any changed files are disallowed. We used to
// return right after writing the output, which meant that --github-output would
// never fail.
func TestGitCheckChangedFilesCmd_GithubOutput(t *testing.T) {
	const pipelineHCL = `
changed_files {
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
		file   string
		errHas string
	}{
		"exempt file": {
			file: ".github/instructions/index.md",
		},
		"disallowed file": {
			file:   "vault_ent/secret.go",
			errHas: "one-or-more changed files matched disallowed groups: enterprise",
		},
	} {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			repo := libgit.NewClient(libgit.WithDir(dir), libgit.WithConfig(map[string]string{
				"user.name":      "test",
				"user.email":     "test@example.com",
				"commit.gpgsign": "false",
				"core.hooksPath": "/dev/null",
			}))
			gitRes, err := repo.Init(t.Context(), &libgit.InitOpts{Quiet: true, InitialBranch: "main"})
			require.NoError(t, err, gitRes.String())
			require.NoError(t, os.MkdirAll(filepath.Dir(filepath.Join(dir, test.file)), 0o755))
			require.NoError(t, os.WriteFile(filepath.Join(dir, test.file), []byte(test.file), 0o644))
			gitRes, err = repo.Add(t.Context(), &libgit.AddOpts{All: true})
			require.NoError(t, err, gitRes.String())
			gitRes, err = repo.Commit(t.Context(), &libgit.CommitOpts{Quiet: true, Message: "add"})
			require.NoError(t, err, gitRes.String())

			cfg, err := config.DecodeBytes([]byte(pipelineHCL))
			require.NoError(t, err)
			oldGit, oldDecodeRes := rootCfg.git, rootCfg.configDecodeRes
			rootCfg.git = repo
			rootCfg.configDecodeRes = &config.DecodeRes{Config: cfg}
			t.Cleanup(func() {
				rootCfg.git, rootCfg.configDecodeRes = oldGit, oldDecodeRes
			})

			githubOutput := filepath.Join(t.TempDir(), "github_output")
			t.Setenv("GITHUB_OUTPUT", githubOutput)

			checkGitChangedFiles = &gitpkg.CheckChangedFilesReq{}
			cmd := newGitCheckChangedFilesCmd()
			suppressOutput(t, cmd)
			cmd.SetArgs([]string{"--branch", "main", "--sync-policy", "ce", "--github-output"})

			err = cmd.Execute()
			if test.errHas != "" {
				require.ErrorContains(t, err, test.errHas)
			} else {
				require.NoError(t, err)
			}

			output, err := os.ReadFile(githubOutput)
			require.NoError(t, err)
			require.Contains(t, string(output), "changed-files=")
			require.Contains(t, string(output), test.file)
		})
	}
}
