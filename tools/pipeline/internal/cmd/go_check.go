// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"github.com/spf13/cobra"
)

func newGoCheckCmd() *cobra.Command {
	goCheckCmd := &cobra.Command{
		Use:   "check",
		Short: "Go check commands",
		Long:  "Go check commands",
	}

	goCheckCmd.AddCommand(newGoCheckGenerateCmd())
	goCheckCmd.AddCommand(newGoCheckTidyCmd())
	goCheckCmd.AddCommand(newGoCheckWorkspaceCmd())

	return goCheckCmd
}
