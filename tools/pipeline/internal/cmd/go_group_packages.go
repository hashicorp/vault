// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"fmt"
	"io"
	"strings"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/golang"
	"github.com/spf13/cobra"
)

var goGroupPackagesReq = &golang.GroupPackagesReq{}

func newGoGroupPackagesCmd() *cobra.Command {
	goGroupPackagesCmd := &cobra.Command{
		Use:   "packages {PACKAGE|PATH|-}... [--go-work PATH]",
		Short: "Group Go packages by Go workspace module",
		Long: `Group Go packages by the Go workspace module that contains them.

Each argument is a package import path, a package directory or a .go file,
which stands for the package in its directory. An argument that exists on disk
is a path, relative to the current directory; anything else is an import path.
'-' reads whitespace-separated arguments from stdin.

A package belongs to the module whose directory contains it, so nested modules
like api/auth/approle are grouped on their own. Packages are printed as import
paths. The command fails if any argument isn't a package directory in a module
in go.work. The go command isn't run.

Examples:
  # Group packages by file, directory or import path.
  pipeline go group packages sdk/logical/token.go api/auth/approle \
    github.com/hashicorp/vault/vault

  # Group a whitespace-separated list from stdin.
  echo "github.com/hashicorp/vault/sdk/helper/tpm vault" |
    pipeline go group packages --format json -`,
		RunE: runGoGroupPackagesCmd,
		Args: cobra.MinimumNArgs(1),
	}

	goGroupPackagesCmd.PersistentFlags().StringVar(&goGroupPackagesReq.GoWork, "go-work", "", "Path to go.work. Defaults to the first go.work in the current directory or its parents")

	return goGroupPackagesCmd
}

func runGoGroupPackagesCmd(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	goGroupPackagesReq.Packages = []string{}
	for _, arg := range args {
		if arg != "-" {
			goGroupPackagesReq.Packages = append(goGroupPackagesReq.Packages, arg)
			continue
		}

		b, err := io.ReadAll(cmd.InOrStdin())
		if err != nil {
			return fmt.Errorf("reading packages from stdin: %w", err)
		}
		goGroupPackagesReq.Packages = append(goGroupPackagesReq.Packages, strings.Fields(string(b))...)
	}

	res, err := goGroupPackagesReq.Run(cmd.Context())
	if err != nil {
		return fmt.Errorf("grouping go packages: %w", err)
	}

	return printGoCmdResult(cmd, res)
}
