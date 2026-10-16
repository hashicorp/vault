// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"fmt"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/golang"
	"github.com/spf13/cobra"
)

var goCheckTidyReq = &golang.CheckTidyReq{}

func newGoCheckTidyCmd() *cobra.Command {
	goCheckTidyCmd := &cobra.Command{
		Use:   "tidy [--go-work PATH] [--module DIR]...",
		Short: "Check that Go workspace modules are tidy",
		Long: `Check that Go workspace modules are tidy.

Runs 'go mod tidy -diff' with GOWORK=off in each module's directory, one module
at a time, and reports the changes tidy would make. Nothing is modified. Tidy
may need network access to fetch go.mod files that aren't in the module cache.`,
		RunE: runGoCheckTidyCmd,
		Args: cobra.NoArgs,
	}

	goCheckTidyCmd.PersistentFlags().StringVar(&goCheckTidyReq.GoWork, "go-work", "", "Path to go.work. Defaults to the first go.work in the current directory or its parents")
	goCheckTidyCmd.PersistentFlags().StringArrayVar(&goCheckTidyReq.Modules, "module", nil, "Repeatable module directory relative to the workspace root. Defaults to every module")

	return goCheckTidyCmd
}

func runGoCheckTidyCmd(cmd *cobra.Command, _ []string) error {
	cmd.SilenceUsage = true

	res, err := goCheckTidyReq.Run(cmd.Context())
	if err != nil {
		return fmt.Errorf("checking go mod tidy: %w", err)
	}

	if err := printGoCmdResult(cmd, res); err != nil {
		return err
	}

	if !res.Success() {
		return fmt.Errorf("check failed: %s", res.String())
	}

	return nil
}
