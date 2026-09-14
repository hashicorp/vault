// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"github.com/spf13/cobra"
)

func newGoSyncCmd() *cobra.Command {
	goSyncCmd := &cobra.Command{
		Use:   "sync",
		Short: "go sync subcommands",
		Long:  "go sync subcommands",
	}

	goSyncCmd.AddCommand(newGoSyncModCmd())

	return goSyncCmd
}
