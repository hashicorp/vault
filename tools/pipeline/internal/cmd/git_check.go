// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"github.com/spf13/cobra"
)

func newGitCheckCmd() *cobra.Command {
	checkCmd := &cobra.Command{
		Use:   "check",
		Short: "Git check commands",
		Long:  "Git check commands",
	}

	checkCmd.AddCommand(newGitCheckChangedFilesCmd())

	return checkCmd
}
