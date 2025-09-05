// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"github.com/spf13/cobra"
)

func newGithubFindCmd() *cobra.Command {
	findCmd := &cobra.Command{
		Use:   "find",
		Short: "Github find commands",
		Long:  "Github find commands",
	}
	findCmd.AddCommand(newGithubFindWorkflowArtifactCmd())

	return findCmd
}
