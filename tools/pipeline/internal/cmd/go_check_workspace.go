// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"fmt"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/golang"
	"github.com/spf13/cobra"
)

var goCheckWorkspaceReq = &golang.CheckWorkspaceReq{}

func newGoCheckWorkspaceCmd() *cobra.Command {
	goCheckWorkspaceCmd := &cobra.Command{
		Use:   "workspace [--go-work PATH] [--published PATTERN]...",
		Short: "Check go.work and the replace directives of workspace modules",
		Long: `Check that go.work and the committed go.mod files agree, and that every
workspace module builds against the right version of the workspace modules it
requires when built on its own with GOWORK=off.

Unpublished modules must replace every workspace module they require with that
module's directory, so they build against the code on the branch. Published
modules, selected with --published, are imported by other modules, and Go
ignores replace directives outside the main module. They must not replace
workspace modules, so they build against the released versions their importers
get.

It reports:
  unregistered-module    a go.mod in the git index whose directory isn't used in
                         go.work. Directories the go command ignores (testdata,
                         vendor, and names starting with '.' or '_') are skipped.
  missing-local-replace  an unpublished module that requires another workspace
                         module without a replace to that module's directory.
  invalid-local-replace  a replace of a workspace module, in an unpublished
                         module, that only applies to one version, isn't a
                         directory replace, or points at another directory.
  published-replace      a published module that replaces another workspace
                         module.`,
		RunE: runGoCheckWorkspaceCmd,
		Args: cobra.NoArgs,
	}

	goCheckWorkspaceCmd.PersistentFlags().StringVar(&goCheckWorkspaceReq.GoWork, "go-work", "", "Path to go.work. Defaults to the first go.work in the current directory or its parents")
	goCheckWorkspaceCmd.PersistentFlags().StringArrayVar(&goCheckWorkspaceReq.Published, "published", nil, "Repeatable glob of published module directories relative to the workspace root, e.g. 'api/auth/*'. Each must match a module")

	return goCheckWorkspaceCmd
}

func runGoCheckWorkspaceCmd(cmd *cobra.Command, _ []string) error {
	cmd.SilenceUsage = true

	res, err := goCheckWorkspaceReq.Run(cmd.Context())
	if err != nil {
		return fmt.Errorf("checking go workspace: %w", err)
	}

	if err := printGoCmdResult(cmd, res); err != nil {
		return err
	}

	if !res.Success() {
		return fmt.Errorf("check failed: %s", res.String())
	}

	return nil
}
