// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"github.com/spf13/cobra"
)

func newGoDiffCmd() *cobra.Command {
	goModCmd := &cobra.Command{
		Use:   "diff",
		Short: "Go diff commands",
		Long:  "Go diff commands",
	}

	goModCmd.AddCommand(newGoDiffModCmd())

	return goModCmd
}
