// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cmd

import "github.com/spf13/cobra"

func newGenerateCmd() *cobra.Command {
	releases := &cobra.Command{
		Use:   "generate",
		Short: "Pipeline configuration generation tasks",
		Long:  "Pipeline configuration generation tasks",
	}

	releases.AddCommand(newGenerateEnosDynamicConfigCmd())

	return releases
}
