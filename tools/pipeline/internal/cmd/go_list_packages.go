// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"fmt"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/golang"
	"github.com/spf13/cobra"
)

var goListPackagesReq = &golang.ListPackagesReq{}

func newGoListPackagesCmd() *cobra.Command {
	goListPackagesCmd := &cobra.Command{
		Use:   "packages [--go-work PATH] [--tags TAGS] [--module DIR]...",
		Short: "List the packages of each module in a Go workspace",
		Long: `List the packages of each module in a Go workspace.

Each module's packages are listed from that module's directory with GOWORK=off,
so they match what the module's own go.mod builds. A module that fails to load
fails the command.

Examples:
  # List every package in every workspace module.
  pipeline go list packages

  # List the packages of two modules with build tags.
  pipeline go list packages --tags ent,enterprise --module . --module sdk

  # Print one import path per line.
  pipeline go list packages --format json | jq -r '.modules[].packages[]'`,
		RunE: runGoListPackagesCmd,
		Args: cobra.NoArgs,
	}

	goListPackagesCmd.PersistentFlags().StringVar(&goListPackagesReq.GoWork, "go-work", "", "Path to go.work. Defaults to the first go.work in the current directory or its parents")
	goListPackagesCmd.PersistentFlags().StringVar(&goListPackagesReq.Tags, "tags", "", "Comma- or space-separated build tags")
	goListPackagesCmd.PersistentFlags().StringArrayVar(&goListPackagesReq.Modules, "module", nil, "Repeatable module directory relative to the workspace root. Defaults to every module")

	return goListPackagesCmd
}

func runGoListPackagesCmd(cmd *cobra.Command, _ []string) error {
	cmd.SilenceUsage = true

	res, err := goListPackagesReq.Run(cmd.Context())
	if err != nil {
		return fmt.Errorf("listing go packages: %w", err)
	}

	return printGoCmdResult(cmd, res)
}
