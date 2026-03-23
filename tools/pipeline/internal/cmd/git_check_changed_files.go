// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"fmt"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/git"
	"github.com/spf13/cobra"
)

var checkGitChangedFiles = &git.CheckChangedFilesReq{}

func newGitCheckChangedFilesCmd() *cobra.Command {
	changedFilesCmd := &cobra.Command{
		Use:   "changed-files (--branch <branch> | --range <range> | --commit <sha>) --disallowed-group <group>... [--pipeline-config .release/pipeline.hcl]",
		Short: "Check if any changed files are matching disallowed groups",
		Long:  "Check if any changed files are matching disallowed groups. You must specify a path to a pipeline.hcl of enable auto recusive search for one.",
		RunE:  runGitCheckChangedFilesCmd,
		Args:  cobra.NoArgs,
	}

	changedFilesCmd.PersistentFlags().StringVarP(&checkGitChangedFiles.Branch, "branch", "b", "", "The branch to compare against")
	changedFilesCmd.PersistentFlags().StringVarP(&checkGitChangedFiles.Range, "range", "r", "", "The commit range to compare (e.g., HEAD~5..HEAD)")
	changedFilesCmd.PersistentFlags().StringVarP(&checkGitChangedFiles.Commit, "commit", "c", "", "The specific commit SHA to analyze")
	changedFilesCmd.PersistentFlags().StringSliceVarP(&checkGitChangedFiles.DisallowedGroups, "disallowed-groups", "g", nil, "File group(s) to check changed files for")
	changedFilesCmd.PersistentFlags().BoolVar(&checkGitChangedFiles.WriteToGithubOutput, "github-output", false, "Whether or not to write 'changed-files' to $GITHUB_OUTPUT")

	return changedFilesCmd
}

func runGitCheckChangedFilesCmd(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true // Don't spam the usage on failure

	checkGitChangedFiles.DecodeRes = rootCfg.configDecodeRes
	res, err := checkGitChangedFiles.Run(cmd.Context(), rootCfg.git)
	if err != nil {
		return fmt.Errorf("checking changed files: %w", err)
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

	if checkGitChangedFiles.WriteToGithubOutput {
		output, err := res.ToGithubOutput()
		if err != nil {
			return err
		}

		return writeToGithubOutput("changed-files", output)
	}

	if len(res.MatchedGroups) > 0 {
		return fmt.Errorf("one-or-more changed files matched disallowed groups: %s", res.MatchedGroups.String())
	}

	return nil
}
