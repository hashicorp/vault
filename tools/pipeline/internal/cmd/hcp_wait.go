// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

func newHCPWaitCmd() *cobra.Command {
	waitCmd := &cobra.Command{
		Use:   "wait",
		Short: "HCP wait commands",
		Long:  "HCP wait commands",
		RunE:  func(*cobra.Command, []string) error { return errors.New("unimplemented") },
	}
	waitCmd.AddCommand(newHCPWaitForImageCmd())

	return waitCmd
}
