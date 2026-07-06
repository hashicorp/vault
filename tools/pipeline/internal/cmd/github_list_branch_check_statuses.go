// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"fmt"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/github"
	"github.com/spf13/cobra"
)

var listBranchCheckStatusesReq = &github.ListBranchCheckStatusesReq{}

func newGithubListBranchCheckStatusesCmd() *cobra.Command {
	listBranchCheckStatusesCmd := &cobra.Command{
		Use:   "branch-check-statuses [--owner hashicorp --repo vault-enterprise --branches release/2.x.x release/1.22.x --retry-wait 5 --max-retry-minutes 75 --max-concurrent 5]",
		Short: "List branch check statuses",
		Long:  "List branch check status once completed for all branches. Use --log=info or --log=debug for detailed progress logging.",
		RunE:  runListBranchCheckStatusesCmd,
	}

	// I've seen this pattern used in other cobra commands, but I'm not sure if it's the best practice. Is it? No plans for subcommands at the moment
	listBranchCheckStatusesCmd.PersistentFlags().StringVar(&listBranchCheckStatusesReq.Owner, "owner", "hashicorp", "owner of the repository")
	listBranchCheckStatusesCmd.PersistentFlags().StringVar(&listBranchCheckStatusesReq.Repo, "repo", "vault-enterprise", "repository name")
	listBranchCheckStatusesCmd.PersistentFlags().StringSliceVar(&listBranchCheckStatusesReq.Branches, "branches", nil, "branches to check")
	listBranchCheckStatusesCmd.PersistentFlags().IntVar(&listBranchCheckStatusesReq.RetryWait, "retry-wait", 5, "retry wait time in minutes")
	listBranchCheckStatusesCmd.PersistentFlags().IntVar(&listBranchCheckStatusesReq.MaxRetry, "max-retry-minutes", 75, "max total retry time in minutes")
	listBranchCheckStatusesCmd.PersistentFlags().IntVar(&listBranchCheckStatusesReq.MaxConcurrent, "max-concurrent", 5, "max concurrent requests")

	return listBranchCheckStatusesCmd
}

func runListBranchCheckStatusesCmd(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true // we don't want the usage continually printed out

	res, err := listBranchCheckStatusesReq.Run(cmd.Context(), githubCmdState.GithubV4)
	if err != nil {
		return fmt.Errorf("error getting branch check statuses: %w", err)
	}

	switch rootCfg.format {
	case "json":
		json, err := res.ToJSON()
		if err != nil {
			return fmt.Errorf("error converting response to JSON: %w", err)
		}
		fmt.Println(string(json))
	default:
		fmt.Println(res.ToTable().Render())
	}

	return nil
}
