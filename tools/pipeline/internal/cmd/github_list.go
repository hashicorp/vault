// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"github.com/spf13/cobra"
)

func newGithubListCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "Github list commands",
		Long:  "Github list commands",
	}

	listCmd.AddCommand(newGithubListRunCmd())
	listCmd.AddCommand(newGithubListChangedFilesCmd())

	return listCmd
}
