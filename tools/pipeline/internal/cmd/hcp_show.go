// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cmd

import "github.com/spf13/cobra"

func newHCPShowCmd() *cobra.Command {
	showCmd := &cobra.Command{
		Use:   "show",
		Short: "HCP show commands",
		Long:  "HCP show commands",
	}
	showCmd.AddCommand(newHCPShowImageCmd())

	return showCmd
}
