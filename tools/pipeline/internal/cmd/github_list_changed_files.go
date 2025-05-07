// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/github"
	"github.com/spf13/cobra"
)

var listGithubChangedFiles = &github.ListChangedFilesReq{}

func newGithubListChangedFilesCmd() *cobra.Command {
	listRuns := &cobra.Command{
		Use:   "changed-files [--pr 1234 | --commit abcd1234 ]",
		Short: "List changed files in a pull request or commit",
		Long:  "List changed files in a pull request or commit",
		RunE:  runListGithubChangedFilesCmd,
	}

	listRuns.PersistentFlags().StringVarP(&listGithubChangedFiles.Owner, "owner", "o", "hashicorp", "The Github organization")
	listRuns.PersistentFlags().StringVarP(&listGithubChangedFiles.Repo, "repo", "r", "vault", "The Github repository. Private repositories require auth via a GITHUB_TOKEN env var")
	listRuns.PersistentFlags().StringVarP(&listGithubChangedFiles.CommitSHA, "commit", "c", "", "The commit SHA to use as a changed file source")
	listRuns.PersistentFlags().IntVarP(&listGithubChangedFiles.PullNumber, "pr", "p", 0, "The pull request to use as a changed file source")
	listRuns.PersistentFlags().BoolVarP(&listGithubChangedFiles.GroupFiles, "group", "g", true, "Whether or not to determine changed file groups")
	listRuns.PersistentFlags().BoolVar(&listGithubChangedFiles.WriteToGithubOutput, "github-output", false, "Whether or not to write 'changed-files' to $GITHUB_OUTPUT")

	return listRuns
}

func runListGithubChangedFilesCmd(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true // Don't spam the usage on failure

	res, err := listGithubChangedFiles.Run(context.TODO(), githubCmdState.Client)
	if err != nil {
		return fmt.Errorf("listing github workflow failures: %w", err)
	}

	switch githubCmdFlags.Format {
	case "json":
		b, err := res.ToJSON()
		if err != nil {
			return err
		}
		fmt.Println(string(b))
	default:
		fmt.Println(res.ToTable(listGithubChangedFiles.GroupFiles))
	}

	if listGithubChangedFiles.WriteToGithubOutput {
		output, err := res.ToGithubOutput()
		if err != nil {
			return err
		}

		return writeToGithubOutput("changed-files", output)
	}

	return nil
}
