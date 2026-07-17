// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"github.com/spf13/cobra"
)

func newSarifCmd() *cobra.Command {
	sarifCmd := &cobra.Command{
		Use:   "sarif",
		Short: "Sarif commands",
		Long:  "Commands for working with Static Analysis Results Interchange Format files",
	}
	sarifCmd.AddCommand(newSarifConvertZapCmd())
	sarifCmd.AddCommand(newSarifCVEReportCmd())

	return sarifCmd
}
