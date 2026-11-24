// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"github.com/spf13/cobra"
)

func newGoGroupCmd() *cobra.Command {
	goGroupCmd := &cobra.Command{
		Use:   "group",
		Short: "Go group commands",
		Long:  "Go group commands",
	}

	goGroupCmd.AddCommand(newGoGroupPackagesCmd())

	return goGroupCmd
}
