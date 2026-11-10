// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"fmt"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/changed"
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/github"
	"github.com/spf13/cobra"
)

var (
	checkGithubChangedFiles                 = &github.CheckChangedFilesReq{}
	checkGithubChangedFilesDisallowedGroups []string
)

func newGithubCheckChangedFilesCmd() *cobra.Command {
	changedFilesCmd := &cobra.Command{
		Use:   "changed-files (--pr <number> | --commit <sha>) (--sync-policy <policy> | --disallowed-groups <group>...) [--pipeline-config .release/pipeline.hcl]",
		Short: "Check if any files changed in a pull request or commit are matching disallowed groups",
		Long:  "Check if any files changed in a pull request or commit are matching disallowed groups. We check the changed files for disallowed groups from the --sync-policy in pipeline.hcl and any --disallowed-groups. Files that the sync policy exempts from disallowed groups are shown but won't fail the check. You must specify a path to a pipeline.hcl or enable auto recursive search for one.",
		RunE:  runGithubCheckChangedFilesCmd,
		Args:  cobra.NoArgs,
	}

	changedFilesCmd.PersistentFlags().StringVarP(&checkGithubChangedFiles.Owner, "owner", "o", "hashicorp", "The Github organization")
	changedFilesCmd.PersistentFlags().StringVarP(&checkGithubChangedFiles.Repo, "repo", "r", "vault-enterprise", "The Github repository. Private repositories require auth via a GITHUB_TOKEN env var")
	changedFilesCmd.PersistentFlags().StringVarP(&checkGithubChangedFiles.CommitSHA, "commit", "c", "", "The commit SHA to use as a changed file source")
	changedFilesCmd.PersistentFlags().IntVarP(&checkGithubChangedFiles.PullNumber, "pr", "p", 0, "The pull request to use as a changed file source")
	changedFilesCmd.PersistentFlags().StringSliceVarP(&checkGithubChangedFilesDisallowedGroups, "disallowed-groups", "g", nil, "File group(s) to check changed files for. They're added to the disallowed groups of the sync policy")
	changedFilesCmd.PersistentFlags().StringVar(&checkGithubChangedFiles.SyncPolicy, "sync-policy", "", "The sync policy in pipeline.hcl to check changed files with")

	return changedFilesCmd
}

func runGithubCheckChangedFilesCmd(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true // Don't spam the usage on failure

	checkGithubChangedFiles.DecodeRes = rootCfg.configDecodeRes
	for _, group := range checkGithubChangedFilesDisallowedGroups {
		checkGithubChangedFiles.DisallowedGroups = checkGithubChangedFiles.DisallowedGroups.Add(changed.FileGroup(group))
	}

	res, err := checkGithubChangedFiles.Run(cmd.Context(), githubCmdState.GithubV3)
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

	if len(res.MatchedGroups) > 0 {
		return fmt.Errorf("one-or-more changed files matched disallowed groups: %s", res.MatchedGroups.String())
	}

	return nil
}
