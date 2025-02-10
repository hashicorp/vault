// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cmd

import "github.com/spf13/cobra"

type githubCommandFlags struct {
	Format string `json:"format,omitempty"`
}

var githubCmdFlags = &githubCommandFlags{}

func newGithubCmd() *cobra.Command {
	github := &cobra.Command{
		Use:   "github",
		Short: "Github commands",
		Long:  "Github commands",
	}

	github.PersistentFlags().StringVarP(&githubCmdFlags.Format, "format", "f", "table", "The output format. Can be 'json' or 'table'")

	github.AddCommand(newGithubListRunCmd())

	return github
}
