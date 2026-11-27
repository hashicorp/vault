// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"testing"

	gh "github.com/google/go-github/v83/github"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

// TestGithubListDisconnectedJobsCmd checks that we reject bad run IDs, attempts,
// and formats before we ever call the Github API.
func TestGithubListDisconnectedJobsCmd(t *testing.T) {
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
		"negative attempt": {args: []string{"1", "--attempt", "-1"}, expectErr: "attempt can't be negative"},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			parent := &cobra.Command{Use: "list"}
			parent.PersistentFlags().StringP("format", "f", "table", "")
			parent.AddCommand(newGithubListDisconnectedJobsCmd(&githubCommandState{Owner: "o", Repo: "r", GithubV3: gh.NewClient(nil)}))
			suppressOutput(t, parent)
			parent.SetArgs(append([]string{"disconnected-jobs"}, test.args...))
			require.ErrorContains(t, parent.Execute(), test.expectErr)
		})
	}
}

// TestGithubListDisconnectedJobsCmd_Registration checks that the command lives
// under 'github list' and inherits the repository flags from 'github'. We run it
// sequentially because the older commands bind flags to globals.
func TestGithubListDisconnectedJobsCmd_Registration(t *testing.T) {
	cmd, _, err := newGithubCmd().Find([]string{"list", "disconnected-jobs"})
	require.NoError(t, err)
	require.Equal(t, "disconnected-jobs", cmd.Name())
	require.NotPanics(t, func() { cmd.InheritedFlags() })
	require.NotNil(t, cmd.InheritedFlags().Lookup("owner"))
	require.NotNil(t, cmd.InheritedFlags().Lookup("repo"))
	require.NotNil(t, cmd.Flags().Lookup("github-output"))
	require.NotNil(t, cmd.Flags().Lookup("attempt"))
}
