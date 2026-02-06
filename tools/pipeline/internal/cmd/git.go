// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package cmd

import "github.com/spf13/cobra"

func newGitCmd() *cobra.Command {
	gitCmd := &cobra.Command{
		Use:   "git",
		Short: "Git commands",
		Long:  "Git commands",
	}

	gitCmd.AddCommand(newGitCheckCmd())
	gitCmd.AddCommand(newGitListCmd())

	return gitCmd
}
