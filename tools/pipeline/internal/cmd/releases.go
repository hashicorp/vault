// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cmd

import "github.com/spf13/cobra"

func newReleasesCmd() *cobra.Command {
	releases := &cobra.Command{
		Use:   "releases",
		Short: "Releases API related tasks",
		Long:  "Releases API related tasks",
	}

	releases.AddCommand(newReleasesVersionsBetweenCmd())

	return releases
}
