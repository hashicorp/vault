// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package cmd

import "github.com/spf13/cobra"

func newGitSyncCmd() *cobra.Command {
	gitSyncCmd := &cobra.Command{
		Use:   "sync",
		Short: "git sync subcommands",
		Long:  "git sync subcommands",
	}

	gitSyncCmd.AddCommand(newGitSyncGoModCmd())

	return gitSyncCmd
}
