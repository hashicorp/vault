// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"fmt"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/golang"
	"github.com/spf13/cobra"
)

var goListModulesReq = &golang.ListModulesReq{}

func newGoListModulesCmd() *cobra.Command {
	goListModulesCmd := &cobra.Command{
		Use:   "modules [--go-work PATH]",
		Short: "List the modules in a Go workspace",
		Long: `List the modules in the use directives of a go.work file, in go.work order.

The JSON output includes each module's path, its directory relative to the
workspace root, and its go.mod path relative to the workspace root.

Examples:
  # List the modules in the go.work of the current directory or a parent.
  pipeline go list modules

  # Print every go.mod path.
  pipeline go list modules --format json | jq -r '.modules[].go_mod'`,
		RunE: runGoListModulesCmd,
		Args: cobra.NoArgs,
	}

	goListModulesCmd.PersistentFlags().StringVar(&goListModulesReq.GoWork, "go-work", "", "Path to go.work. Defaults to the first go.work in the current directory or its parents")

	return goListModulesCmd
}

func runGoListModulesCmd(cmd *cobra.Command, _ []string) error {
	cmd.SilenceUsage = true

	res, err := goListModulesReq.Run(cmd.Context())
	if err != nil {
		return fmt.Errorf("listing go modules: %w", err)
	}

	return printGoCmdResult(cmd, res)
}
