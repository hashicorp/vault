// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"github.com/spf13/cobra"
)

func newGithubCheckCmd() *cobra.Command {
	checkCmd := &cobra.Command{
		Use:   "check",
		Short: "Github check commands",
		Long:  "Github check commands",
	}
	checkCmd.AddCommand(newCheckGithubCommitStatusCmd())

	return checkCmd
}
