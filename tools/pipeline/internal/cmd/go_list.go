// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"github.com/spf13/cobra"
)

func newGoListCmd() *cobra.Command {
	goListCmd := &cobra.Command{
		Use:   "list",
		Short: "Go list commands",
		Long:  "Go list commands",
	}

	goListCmd.AddCommand(newGoListModulesCmd())
	goListCmd.AddCommand(newGoListPackagesCmd())

	return goListCmd
}
