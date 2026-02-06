// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/git"
	"github.com/spf13/cobra"
)

var listGitChangedFiles = &git.ListChangedFilesReq{}

func newGitListChangedFilesCmd() *cobra.Command {
	changedFilesCmd := &cobra.Command{
		Use:   "changed-files [--branch <branch> | --range <range> | --commit <sha>]",
		Short: "List changed files using git",
		Long:  "List changed files using git by specifying a branch, range, or commit",
		RunE:  runGitListChangedFilesCmd,
		Args:  cobra.NoArgs,
	}

	changedFilesCmd.PersistentFlags().StringVarP(&listGitChangedFiles.Branch, "branch", "b", "", "The branch to compare against")
	changedFilesCmd.PersistentFlags().StringVarP(&listGitChangedFiles.Range, "range", "r", "", "The commit range to compare (e.g., HEAD~5..HEAD)")
	changedFilesCmd.PersistentFlags().StringVarP(&listGitChangedFiles.Commit, "commit", "c", "", "The specific commit SHA to analyze")
	changedFilesCmd.PersistentFlags().BoolVarP(&listGitChangedFiles.GroupFiles, "group", "g", true, "Whether or not to determine changed file groups")
	changedFilesCmd.PersistentFlags().BoolVar(&listGitChangedFiles.WriteToGithubOutput, "github-output", false, "Whether or not to write 'changed-files' to $GITHUB_OUTPUT")

	return changedFilesCmd
}

func runGitListChangedFilesCmd(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true // Don't spam the usage on failure

	res, err := listGitChangedFiles.Run(context.TODO(), githubCmdState.Git)
	if err != nil {
		return fmt.Errorf("listing changed files: %w", err)
	}

	switch rootCfg.format {
	case "json":
		b, err := res.ToJSON()
		if err != nil {
			return err
		}
		fmt.Println(string(b))
	default:
		fmt.Println(res.ToTable(listGithubChangedFiles.GroupFiles))
	}

	if listGitChangedFiles.WriteToGithubOutput {
		output, err := res.ToGithubOutput()
		if err != nil {
			return err
		}

		return writeToGithubOutput("changed-files", output)
	}

	return nil
}
