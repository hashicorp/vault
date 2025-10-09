// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/github"
	"github.com/spf13/cobra"
)

var listGithubCommitStatuses = &github.ListCommitStatusesReq{}

func newGithubListCommitStatusesCmd() *cobra.Command {
	listCommitStatusesCmd := &cobra.Command{
		Use:   "commit-statuses [--pr 1234 | --commit 1a61dbe57a7538f22f7fa3fb65900d9d41c2ba14]",
		Short: "List the statuses associated with a commit",
		Long:  "List the statuses associated with a commit. When given passed the --pr flag it will use the head ref of the PR as the commit",
		RunE:  runListGithubCommitStatuses,
	}

	listCommitStatusesCmd.PersistentFlags().StringVarP(&listGithubCommitStatuses.Owner, "owner", "o", "hashicorp", "The Github organization")
	listCommitStatusesCmd.PersistentFlags().StringVarP(&listGithubCommitStatuses.Repo, "repo", "r", "vault", "The Github repository. Private repositories require auth via a GITHUB_TOKEN env var")
	listCommitStatusesCmd.PersistentFlags().StringVarP(&listGithubCommitStatuses.Commit, "commit", "c", "", "The commit you wish to list the statuses of")
	listCommitStatusesCmd.PersistentFlags().IntVarP(&listGithubCommitStatuses.PR, "pr", "p", 0, "The Pull Request number you wish to use as the source. The HEAD commit will be used")

	return listCommitStatusesCmd
}

func runListGithubCommitStatuses(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true // Don't spam the usage on failure

	res, err := listGithubCommitStatuses.Run(context.TODO(), githubCmdState.Github)
	if err != nil {
		return fmt.Errorf("listing github commit statuses: %w", err)
	}

	switch rootCfg.format {
	case "json":
		b, err := json.Marshal(res)
		if err != nil {
			return fmt.Errorf("marshaling response to JSON: %w", err)
		}
		fmt.Println(string(b))
	default:
		fmt.Println(res.ToTable().Render())
	}

	return err
}
