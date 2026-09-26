// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"testing"

	githubpkg "github.com/hashicorp/vault/tools/pipeline/internal/pkg/github"
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/golang"
	"github.com/stretchr/testify/require"
)

// TestGithubSyncGoModCmd_MissingRequiredFlags verifies that omitting each
// individually-required flag causes Cobra to return a required-flag error
// before RunE is called, and that the error message names the missing flag.
func TestGithubSyncGoModCmd_MissingRequiredFlags(t *testing.T) {
	for name, tc := range map[string]struct {
		args        []string
		errContains string
	}{
		"missing a-branch": {
			args:        []string{"--b-branch", "ce/main"},
			errContains: "a-branch",
		},
		"missing b-branch": {
			args:        []string{"--a-branch", "main"},
			errContains: "b-branch",
		},
	} {
		t.Run(name, func(t *testing.T) {
			githubSyncGoModReq = &githubpkg.SyncGoModReq{
				DiffOpts:    golang.DefaultSyncDiffOpts(),
				Commit:      true,
				PullRequest: true,
			}

			cmd := newGithubSyncGoModCmd()
			suppressOutput(t, cmd)
			cmd.SetArgs(tc.args)
			err := cmd.Execute()
			require.Error(t, err)
			require.Contains(t, err.Error(), tc.errContains)
		})
	}
}

// TestGithubSyncGoModCmd_FlagDefaults verifies that all documented flag defaults
// are set correctly when the command is created.
func TestGithubSyncGoModCmd_FlagDefaults(t *testing.T) {
	githubSyncGoModReq = &githubpkg.SyncGoModReq{
		DiffOpts:    golang.DefaultSyncDiffOpts(),
		Commit:      true,
		PullRequest: true,
	}

	cmd := newGithubSyncGoModCmd()
	require.NotNil(t, cmd)

	aOwner, err := cmd.PersistentFlags().GetString("a-owner")
	require.NoError(t, err)
	require.Equal(t, "hashicorp", aOwner)

	aRepo, err := cmd.PersistentFlags().GetString("a-repo")
	require.NoError(t, err)
	require.Equal(t, "vault-enterprise", aRepo)

	bOwner, err := cmd.PersistentFlags().GetString("b-owner")
	require.NoError(t, err)
	require.Equal(t, "hashicorp", bOwner)

	bRepo, err := cmd.PersistentFlags().GetString("b-repo")
	require.NoError(t, err)
	require.Equal(t, "vault-enterprise", bRepo)

	commit, err := cmd.PersistentFlags().GetBool("commit")
	require.NoError(t, err)
	require.True(t, commit)

	pr, err := cmd.PersistentFlags().GetBool("pull-request")
	require.NoError(t, err)
	require.True(t, pr)

	paths, err := cmd.PersistentFlags().GetStringArray("path")
	require.NoError(t, err)
	require.Equal(t, []string{"go.mod"}, paths)

	requireFlag, err := cmd.PersistentFlags().GetBool("require")
	require.NoError(t, err)
	require.True(t, requireFlag)

	replace, err := cmd.PersistentFlags().GetBool("replace")
	require.NoError(t, err)
	require.True(t, replace)

	goFlag, err := cmd.PersistentFlags().GetBool("go")
	require.NoError(t, err)
	require.False(t, goFlag)

	toolchain, err := cmd.PersistentFlags().GetBool("toolchain")
	require.NoError(t, err)
	require.False(t, toolchain)

	strictRequire, err := cmd.PersistentFlags().GetBool("strict-require")
	require.NoError(t, err)
	require.True(t, strictRequire)

	strictReplace, err := cmd.PersistentFlags().GetBool("strict-replace")
	require.NoError(t, err)
	require.True(t, strictReplace)
}

// TestGithubSyncGoModCmd_SameBranchSameRepoError verifies that providing
// --a-branch=main --b-branch=main with the same default owner/repo returns a
// validation error from the implementation layer.
func TestGithubSyncGoModCmd_SameBranchSameRepoError(t *testing.T) {
	githubSyncGoModReq = &githubpkg.SyncGoModReq{
		DiffOpts:    golang.DefaultSyncDiffOpts(),
		Commit:      true,
		PullRequest: true,
	}

	cmd := newGithubSyncGoModCmd()
	suppressOutput(t, cmd)
	cmd.SetArgs([]string{"--a-branch", "main", "--b-branch", "main"})
	err := cmd.Execute()
	require.Error(t, err)
	require.Contains(t, err.Error(), "must differ")
}

// TestGithubSyncGoModCmd_PositionalArgsRejected verifies that providing
// positional arguments is rejected because the command takes no args.
func TestGithubSyncGoModCmd_PositionalArgsRejected(t *testing.T) {
	githubSyncGoModReq = &githubpkg.SyncGoModReq{
		DiffOpts:    golang.DefaultSyncDiffOpts(),
		Commit:      true,
		PullRequest: true,
	}

	cmd := newGithubSyncGoModCmd()
	suppressOutput(t, cmd)
	cmd.SetArgs([]string{"main", "ce/main"})
	err := cmd.Execute()
	require.Error(t, err)
	require.Contains(t, err.Error(), "unknown command")
}
