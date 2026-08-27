// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import "github.com/spf13/cobra"

func newEbomCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ebom",
		Short: "IBM EBOM commands",
		Long:  "Commands for working with IBM Electronic Bill of Materials",
	}
	cmd.AddCommand(newEBOMInspectCmd())
	cmd.AddCommand(newEBOMConvertCmd())
	return cmd
}
