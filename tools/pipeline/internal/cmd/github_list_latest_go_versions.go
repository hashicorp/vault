// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"fmt"
	"maps"
	"slices"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/github"
	"github.com/spf13/cobra"
)

var listLatestGoVersionsReq = &github.ListLatestGoVersionsReq{}

func newGithubListLatestGoVersionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "latest-go-versions",
		Short: "Fetch the latest stable Go versions from go.dev",
		Long: `Fetches the two most recent stable Go releases from https://go.dev/dl/?mode=json
and outputs the latest and previous stable version strings.

When --github-output is set, writes latest_stable_version and
previous_stable_version to $GITHUB_OUTPUT.`,
		RunE: runGithubListLatestGoVersionsCmd,
		Args: cobra.NoArgs,
	}

	cmd.PersistentFlags().BoolVar(&listLatestGoVersionsReq.WriteToGithubOutput, "github-output", false, "Write outputs to $GITHUB_OUTPUT")

	return cmd
}

func runGithubListLatestGoVersionsCmd(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	res, err := listLatestGoVersionsReq.Run(cmd.Context())
	if err != nil {
		return fmt.Errorf("list latest-go-versions: %w", err)
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

	if listLatestGoVersionsReq.WriteToGithubOutput {
		output, err := res.ToGithubOutput()
		if err != nil {
			return err
		}
		for _, k := range slices.Sorted(maps.Keys(output)) {
			if err := writeToGithubOutput(k, []byte(output[k])); err != nil {
				return err
			}
		}
	}

	return nil
}
