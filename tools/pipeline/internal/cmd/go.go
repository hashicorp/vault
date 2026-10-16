// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newGoCmd() *cobra.Command {
	goCmd := &cobra.Command{
		Use:   "go",
		Short: "Go commands",
		Long:  "Go commands",
	}

	goCmd.AddCommand(newGoCheckCmd())
	goCmd.AddCommand(newGoDiffCmd())
	goCmd.AddCommand(newGoGroupCmd())
	goCmd.AddCommand(newGoListCmd())
	goCmd.AddCommand(newGoSyncCmd())

	return goCmd
}

// goCmdResult is a go sub-command result that can be printed in every output
// format.
type goCmdResult interface {
	ToJSON() ([]byte, error)
	ToTable() string
	ToMarkdown() string
}

// printGoCmdResult writes res to cmd.OutOrStdout() in the format requested via
// rootCfg.format.
func printGoCmdResult(cmd *cobra.Command, res goCmdResult) error {
	out := cmd.OutOrStdout()

	switch rootCfg.format {
	case "json":
		b, err := res.ToJSON()
		if err != nil {
			return err
		}
		fmt.Fprintln(out, string(b))
	case "markdown":
		fmt.Fprintln(out, res.ToMarkdown())
	default:
		fmt.Fprintln(out, res.ToTable())
	}

	return nil
}
