// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package cmd

import "github.com/spf13/cobra"

func newConfigCmd() *cobra.Command {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Config commands",
		Long:  "Pipeline configuration commands",
	}

	configCmd.AddCommand(newConfigValidateCmd())

	return configCmd
}
