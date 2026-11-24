// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"fmt"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/github"
	"github.com/spf13/cobra"
)

var listGithubOpenPRs = &github.ListOpenPRsReq{}

func newGithubListOpenPRsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "open-prs [branches...]",
		Short: "List open pull requests across one or more branches",
		Long: `List open (non-draft) pull requests for the given branches.

For each branch, retrieves open PRs including author, PR number, milestone, and
any backport labels. Results can be written to $GITHUB_OUTPUT as a JSON map of
branch to PR list.

If no branches are provided, the active release branches are read from versions.hcl.

Examples:
  # List open PRs on a single branch
  pipeline github list open-prs main \
    --owner hashicorp \
    --repo vault-enterprise

  # List open PRs across multiple branches
  pipeline github list open-prs main release/1.19.x+ent release/1.20.x+ent \
    --owner hashicorp \
    --repo vault-enterprise

  # List open PRs for all active branches from versions.hcl
  pipeline github list open-prs

  # Write results to $GITHUB_OUTPUT
  pipeline github list open-prs main release/1.19.x+ent \
    --github-output

  # Output as JSON
  pipeline github list open-prs main \
    --format json`,
		RunE: runGithubListOpenPRsCmd,
		Args: cobra.ArbitraryArgs,
	}

	cmd.PersistentFlags().StringVarP(&listGithubOpenPRs.Owner, "owner", "o", "hashicorp", "GitHub repository owner")
	cmd.PersistentFlags().StringVarP(&listGithubOpenPRs.Repo, "repo", "r", "vault-enterprise", "GitHub repository name")
	cmd.PersistentFlags().BoolVar(&listGithubOpenPRs.IncludeDrafts, "include-drafts", false, "Include draft pull requests")
	cmd.PersistentFlags().BoolVar(&listGithubOpenPRs.WriteToGithubOutput, "github-output", false, "Write result to $GITHUB_OUTPUT")
	cmd.PersistentFlags().StringVar(&listGithubOpenPRs.GithubOutputKey, "github-output-key", "open-prs-map", "The output key to use when writing to $GITHUB_OUTPUT")

	return cmd
}

func runGithubListOpenPRsCmd(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	listGithubOpenPRs.Branches = args
	listGithubOpenPRs.Releases = rootCfg.versionsDecodeRes

	res, err := listGithubOpenPRs.Run(cmd.Context(), githubCmdState.GithubV3)
	if err != nil {
		return fmt.Errorf("list open pull requests: %w", err)
	}

	switch rootCfg.format {
	case "json":
		b, err := res.ToJSON()
		if err != nil {
			return err
		}
		fmt.Println(string(b))
	case "markdown":
		fmt.Println(res.ToMarkdown())
	default:
		fmt.Println(res.ToTable())
	}

	if listGithubOpenPRs.WriteToGithubOutput {
		output, err := res.ToGithubOutput()
		if err != nil {
			return err
		}
		return writeToGithubOutput(listGithubOpenPRs.GithubOutputKey, output)
	}

	return nil
}
