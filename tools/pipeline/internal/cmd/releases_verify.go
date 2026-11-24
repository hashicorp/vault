// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"github.com/spf13/cobra"
)

func newReleasesVerifyCmd() *cobra.Command {
	verifyCmd := &cobra.Command{
		Use:   "verify",
		Short: "Releases verify commands",
		Long:  "Releases verify commands",
	}

	verifyCmd.AddCommand(newReleasesVerifyBinariesCmd())

	return verifyCmd
}
