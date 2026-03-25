// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"github.com/spf13/cobra"
)

func newGenerateCmd() *cobra.Command {
	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate code and dynamic configuration in the context of the pipeline",
		Long:  "Generate code and dynamic configuration in the context of the pipeline",
	}

	// Add subcommands
	generateCmd.AddCommand(newGenerateTemplateCmd())

	return generateCmd
}
