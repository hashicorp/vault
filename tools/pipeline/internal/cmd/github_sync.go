// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"github.com/spf13/cobra"
)

func newGithubSyncCmd() *cobra.Command {
	syncCmd := &cobra.Command{
		Use:   "sync",
		Short: "Github sync commands",
		Long:  "Github sync commands",
	}
	syncCmd.AddCommand(newSyncGithubBranchCmd())

	return syncCmd
}
