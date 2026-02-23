// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"fmt"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/releases"
	"github.com/spf13/cobra"
)

var listReleaseActiveVersionsReq = &releases.ListActiveVersionsReq{}

func newReleasesListActiveVersionsCmd() *cobra.Command {
	activeVersionsCmd := &cobra.Command{
		Use:   "active-versions [--versions-config .release/versions.hcl]",
		Short: "List the active versions from .release/versions.hcl",
		Long:  "List the active versions from .release/versions.hcl",
		RunE:  runListActiveVersionsReq,
		Args:  cobra.NoArgs,
	}

	activeVersionsCmd.PersistentFlags().BoolVar(&listReleaseActiveVersionsReq.WriteToGithubOutput, "github-output", false, "Whether or not to write 'active-versions' to $GITHUB_OUTPUT")

	return activeVersionsCmd
}

func runListActiveVersionsReq(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true // Don't spam the usage on failure

	listReleaseActiveVersionsReq.VersionsDecodeRes = rootCfg.versionsDecodeRes
	res, err := listReleaseActiveVersionsReq.Run(cmd.Context(), rootCfg.git)
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

	if listReleaseActiveVersionsReq.WriteToGithubOutput {
		output, err := res.ToGithubOutput()
		if err != nil {
			return err
		}

		return writeToGithubOutput("active-versions", output)
	}

	return err
}
