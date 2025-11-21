// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/github"
	"github.com/spf13/cobra"
)

var githubCheckCommitStatusReq = github.CheckCommitStatusReq{}

func newGithubCheckCommitStatusCmd() *cobra.Command {
	checkCommitStatusCmd := &cobra.Command{
		Use:   "commit-status --context <context> --creator <creator> --state <state> [--pr 1234 | --commit 1a61dbe57a7538f22f7fa3fb65900d9d41c2ba14]",
		Short: "Check a commit status for a PR or commit",
		Long:  "Check a commit status for a specific context and creator matches an expected state. When given provided the --pr flag the HEAD sha will be used as the commit",
		RunE:  runGithubCheckCommitStatusCmd,
	}

	checkCommitStatusCmd.PersistentFlags().StringVarP(&githubCheckCommitStatusReq.Owner, "owner", "o", "hashicorp", "The Github organization")
	checkCommitStatusCmd.PersistentFlags().StringVarP(&githubCheckCommitStatusReq.Repo, "repo", "r", "vault", "The Github repository. Private repositories require auth via a GITHUB_TOKEN env var")
	checkCommitStatusCmd.PersistentFlags().StringVarP(&githubCheckCommitStatusReq.Commit, "commit", "c", "", "The commit you wish to list the statuses of")
	checkCommitStatusCmd.PersistentFlags().IntVarP(&githubCheckCommitStatusReq.PR, "pr", "p", 0, "The Pull Request number you wish to use as the source. The HEAD commit will be used")
	checkCommitStatusCmd.PersistentFlags().StringVar(&githubCheckCommitStatusReq.Context, "context", "", "The context of the status. This usually maps to the name of the check that shows up in the Pull Request status box")
	checkCommitStatusCmd.PersistentFlags().StringVar(&githubCheckCommitStatusReq.Creator, "creator", "", "The github login of the creator of the status")
	checkCommitStatusCmd.PersistentFlags().StringVarP(&githubCheckCommitStatusReq.State, "state", "s", "success", "The expected state of the status. Can be one of 'error', 'failure', 'pending', 'success'")

	err := checkCommitStatusCmd.MarkPersistentFlagRequired("context")
	if err != nil {
		panic(err)
	}

	return checkCommitStatusCmd
}

func runGithubCheckCommitStatusCmd(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true // Don't spam the usage on failure

	res, err := githubCheckCommitStatusReq.Run(context.TODO(), githubCmdState.GithubV3)
	if err != nil {
		return err
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

	if !res.CheckSuccessful {
		return fmt.Errorf("no statuses matched expected criteria: %s", res.String())
	}

	return nil
}
