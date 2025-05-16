// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cmd

import "github.com/spf13/cobra"

func newReleasesListCmd() *cobra.Command {
	releases := &cobra.Command{
		Use:   "list",
		Short: "Releases list commands",
		Long:  "Releases list commands",
	}

	releases.AddCommand(newReleasesVersionsBetweenCmd())
	releases.AddCommand(newReleasesListActiveVersionsCmd())

	return releases
}
