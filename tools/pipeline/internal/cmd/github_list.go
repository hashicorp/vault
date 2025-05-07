// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"github.com/spf13/cobra"
)

func newGithubListCmd() *cobra.Command {
	github := &cobra.Command{
		Use:   "list",
		Short: "Github list commands",
		Long:  "Github list commands",
	}

	github.AddCommand(newGithubListRunCmd())
	github.AddCommand(newGithubListChangedFilesCmd())

	return github
}
