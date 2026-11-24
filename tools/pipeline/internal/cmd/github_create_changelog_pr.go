// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/github"
	"github.com/spf13/cobra"
)

var createGithubChangelogPRReq = &github.CreateChangelogPRReq{}

func newCreateGithubChangelogPRCmd() *cobra.Command {
	changelogPRCmd := &cobra.Command{
		Use:   "changelog-pr --repo <repo> --tag <tag> --branch <branch> --strategy <strategy>",
		Short: "Open a pull request that adds a tag's entry to a changelog",
		Long: `Open a pull request that adds an entry for an existing tag to a changelog file.
The strategy decides how we make the entry. We create the pull request branch
from the base branch if it does not exist. If the base branch already has an
entry for the tag, or the pull request is already open, we leave it alone.`,
		RunE: runCreateGithubChangelogPRCmd,
		Args: cobra.NoArgs,
	}

	changelogPRCmd.PersistentFlags().StringVarP(
		&createGithubChangelogPRReq.Owner, "owner", "o", "hashicorp",
		"The Github organization",
	)
	changelogPRCmd.PersistentFlags().StringVarP(
		&createGithubChangelogPRReq.Repo, "repo", "r", "",
		"The Github repository. Private repositories require auth via a GITHUB_TOKEN env var",
	)
	changelogPRCmd.PersistentFlags().StringVar(
		&createGithubChangelogPRReq.Tag, "tag", "",
		"The existing tag whose changelog entry we add",
	)
	changelogPRCmd.PersistentFlags().StringVar(
		&createGithubChangelogPRReq.Branch, "branch", "",
		"The pull request branch",
	)
	changelogPRCmd.PersistentFlags().StringVar(
		&createGithubChangelogPRReq.BaseBranch, "base-branch", "main",
		"The pull request base branch",
	)
	changelogPRCmd.PersistentFlags().StringVar(
		&createGithubChangelogPRReq.Path, "path", "CHANGELOG.md",
		"The path to the changelog file",
	)
	changelogPRCmd.PersistentFlags().StringVar(
		(*string)(&createGithubChangelogPRReq.Strategy), "strategy", "",
		"How we make the changelog entry. One of: "+strings.Join(github.ChangelogStrategies(), ", "),
	)
	changelogPRCmd.PersistentFlags().StringArrayVar(
		&createGithubChangelogPRReq.Reviewers, "reviewer", nil,
		"GitHub login to request a review from. Repeatable",
	)
	changelogPRCmd.PersistentFlags().BoolVar(
		&createGithubChangelogPRReq.WriteToGithubOutput, "github-output", false,
		"Whether or not to write 'changelog-pr' to $GITHUB_OUTPUT",
	)

	for _, flag := range []string{"repo", "tag", "branch", "strategy"} {
		if err := changelogPRCmd.MarkPersistentFlagRequired(flag); err != nil {
			panic(err)
		}
	}

	return changelogPRCmd
}

func runCreateGithubChangelogPRCmd(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true // Don't spam the usage on failure

	res, err := createGithubChangelogPRReq.Run(cmd.Context(), githubCmdState.GithubV3)
	if err != nil {
		return fmt.Errorf("creating changelog pull request: %w", err)
	}

	b, err := res.ToJSON()
	if err != nil {
		return err
	}

	switch rootCfg.format {
	case "json":
		fmt.Println(string(b))
	default:
		fmt.Println(res.ToTable())
	}

	if createGithubChangelogPRReq.WriteToGithubOutput {
		return writeToGithubOutput("changelog-pr", b)
	}

	return nil
}
