// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"github.com/spf13/cobra"
)

func newGitListCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "Git list commands",
		Long:  "Git list commands",
	}

	listCmd.AddCommand(newGitListChangedFilesCmd())

	return listCmd
}
