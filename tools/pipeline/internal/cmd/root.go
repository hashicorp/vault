// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pipeline",
	Short: "Execute pipeline tasks",
	Long:  "Pipeline automation tasks",
}

// Execute executes the root pipeline command.
func Execute() {
	rootCmd.SilenceErrors = true // We handle this below

	rootCmd.AddCommand(newGenerateCmd())
	rootCmd.AddCommand(newReleasesCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
