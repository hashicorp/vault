// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cmd

import "github.com/spf13/cobra"

func newGoCmd() *cobra.Command {
	goCmd := &cobra.Command{
		Use:   "go",
		Short: "Go commands",
		Long:  "Go commands",
	}

	goCmd.AddCommand(newGoDiffCmd())

	return goCmd
}
