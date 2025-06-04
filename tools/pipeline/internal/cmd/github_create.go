// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"github.com/spf13/cobra"
)

func newGithubCreateCmd() *cobra.Command {
	create := &cobra.Command{
		Use:   "create",
		Short: "Github create commands",
		Long:  "Github create commands",
	}
	create.AddCommand(newGithubCreateBackportCmd())

	return create
}
