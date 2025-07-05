// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"github.com/spf13/cobra"
)

func newGithubCreateCmd() *cobra.Command {
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Github create commands",
		Long:  "Github create commands",
	}
	createCmd.AddCommand(newCreateGithubBackportCmd())

	return createCmd
}
