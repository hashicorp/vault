// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/releases"
	"github.com/spf13/cobra"
)

var listReleaseActiveVersionsReq = &releases.ListActiveVersionsReq{}

func newReleasesListActiveVersionsCmd() *cobra.Command {
	activeVersionsCmd := &cobra.Command{
		Use:   "active-versions [.release/versions.hcl]",
		Short: "List the active versions from .release/versions.hcl",
		Long:  "List the active versions from .release/versions.hcl",
		RunE:  runListActiveVersionsReq,
		Args:  cobra.MaximumNArgs(1), // path to .release/versions.hcl
	}

	activeVersionsCmd.PersistentFlags().UintVarP(&listReleaseActiveVersionsReq.Recurse, "recurse", "r", 0, "If no path to a config file is given, recursively search backwards for it and stop at root or until we've his the configured depth.")

	return activeVersionsCmd
}

func runListActiveVersionsReq(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true // Don't spam the usage on failure

	if len(args) > 0 {
		listReleaseActiveVersionsReq.ReleaseVersionConfigPath = args[0]
	}

	res, err := listReleaseActiveVersionsReq.Run(context.TODO())
	if err != nil {
		return err
	}

	switch rootCfg.format {
	case "json":
		b, err := res.ToJSON()
		if err != nil {
			return err
		}
		fmt.Println(string(b))
	default:
		fmt.Println(res.ToTable())
	}

	return err
}
