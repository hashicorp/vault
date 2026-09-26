// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"io"
	"testing"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/git"
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/golang"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

// suppressOutput directs a cobra command's stdout and stderr to io.Discard so
// that Cobra's usage and error messages do not appear in test output. It must
// be called after the command is constructed and before Execute is called.
// Using cmd.SetOut/SetErr avoids any writes to the os.Stdout/os.Stderr globals,
// making the helper safe to use in parallel tests.
func suppressOutput(_ *testing.T, cmd *cobra.Command) {
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
}

// TestGitSyncGoModCmd_WrongArgCount verifies that Cobra rejects any argument
// count other than exactly two, returning an error that contains "accepts 2 arg"
// before RunE is called.
func TestGitSyncGoModCmd_WrongArgCount(t *testing.T) {
	for name, tc := range map[string]struct {
		args []string
	}{
		"no args":    {args: []string{}},
		"one arg":    {args: []string{"main"}},
		"three args": {args: []string{"main", "ce/main", "extra"}},
	} {
		t.Run(name, func(t *testing.T) {
			gitSyncGoModReq = &git.SyncGoModReq{
				DiffOpts: golang.DefaultSyncDiffOpts(),
				Commit:   true,
			}

			cmd := newGitSyncGoModCmd()
			suppressOutput(t, cmd)
			cmd.SetArgs(tc.args)
			err := cmd.Execute()
			require.Error(t, err)
			require.Contains(t, err.Error(), "accepts 2 arg")
		})
	}
}

// TestGitSyncGoModCmd_SameBranchError verifies that providing identical source
// and dest branches returns a validation error from the implementation layer.
func TestGitSyncGoModCmd_SameBranchError(t *testing.T) {
	gitSyncGoModReq = &git.SyncGoModReq{
		DiffOpts: golang.DefaultSyncDiffOpts(),
		Commit:   true,
	}

	cmd := newGitSyncGoModCmd()
	suppressOutput(t, cmd)
	cmd.SetArgs([]string{"main", "main"})
	err := cmd.Execute()
	require.Error(t, err)
	require.Contains(t, err.Error(), "must differ")
}
