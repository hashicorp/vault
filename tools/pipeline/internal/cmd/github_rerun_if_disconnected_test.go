// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"testing"

	gh "github.com/google/go-github/v83/github"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

// TestGithubRerunIfDisconnectedCmd checks that we reject bad run IDs, attempts,
// formats, and timeouts before we ever call the Github API.
func TestGithubRerunIfDisconnectedCmd(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		args      []string
		expectErr string
	}{
		"no ID":            {expectErr: "accepts 1 arg"},
		"extra ID":         {args: []string{"1", "2"}, expectErr: "accepts 1 arg"},
		"invalid ID":       {args: []string{"abc"}, expectErr: "positive int64"},
		"zero ID":          {args: []string{"0"}, expectErr: "positive int64"},
		"negative ID":      {args: []string{"--", "-1"}, expectErr: "positive int64"},
		"overflow ID":      {args: []string{"9223372036854775808"}, expectErr: "positive int64"},
		"csv":              {args: []string{"1", "-f", "csv"}, expectErr: "unsupported format"},
		"unknown format":   {args: []string{"1", "-f", "xml"}, expectErr: "unsupported format"},
		"no timeout":       {args: []string{"1", "--wait", "--timeout", "0s"}, expectErr: "timeout must be positive"},
		"negative attempt": {args: []string{"1", "--attempt", "-1"}, expectErr: "attempt can't be negative"},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			parent := &cobra.Command{Use: "github"}
			parent.PersistentFlags().StringP("format", "f", "table", "")
			parent.AddCommand(newGithubRerunIfDisconnectedCmd(&githubCommandState{Owner: "o", Repo: "r", GithubV3: gh.NewClient(nil)}))
			suppressOutput(t, parent)
			parent.SetArgs(append([]string{"rerun-if-disconnected"}, test.args...))
			require.ErrorContains(t, parent.Execute(), test.expectErr)
		})
	}
}

// TestGithubRerunIfDisconnectedCmd_Registration checks that inherited flags don't
// conflict with existing commands or change their defaults. We run it sequentially
// because the older commands bind flags to globals.
func TestGithubRerunIfDisconnectedCmd_Registration(t *testing.T) {
	parent := newGithubCmd()
	var visit func(*cobra.Command)
	visit = func(cmd *cobra.Command) {
		require.NotPanics(t, func() { cmd.InheritedFlags() }, cmd.CommandPath())
		for _, child := range cmd.Commands() {
			visit(child)
		}
	}
	visit(parent)
	cmd, _, err := parent.Find([]string{"rerun-if-disconnected"})
	require.NoError(t, err)
	require.Equal(t, "rerun-if-disconnected", cmd.Name())
	require.Equal(t, "hashicorp", cmd.InheritedFlags().Lookup("owner").DefValue)
	require.Equal(t, "vault-enterprise", cmd.InheritedFlags().Lookup("repo").DefValue)
	for _, tc := range []struct {
		path []string
		repo string
	}{
		{[]string{"check", "workflow-run"}, "vault"},
		{[]string{"create", "backport"}, "vault-enterprise"},
	} {
		child, _, err := parent.Find(tc.path)
		require.NoError(t, err)
		require.Equal(t, tc.repo, child.PersistentFlags().Lookup("repo").DefValue)
	}
}
