// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"github.com/spf13/cobra"
)

func newGithubCloseCmd() *cobra.Command {
	closeCmd := &cobra.Command{
		Use:   "close",
		Short: "Github close commands",
		Long:  "Github close commands",
	}
	closeCmd.AddCommand(newCloseGithubCopiedPullRequestCmd())

	return closeCmd
}
