// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"github.com/spf13/cobra"
)

func newGithubCopyCmd() *cobra.Command {
	copyCmd := &cobra.Command{
		Use:   "copy",
		Short: "Github copy commands",
		Long:  "Github copy commands",
	}
	copyCmd.AddCommand(newCopyGithubPullRequestCmd())

	return copyCmd
}
