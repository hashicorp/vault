// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"fmt"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/github"
	"github.com/spf13/cobra"
)

var createGithubTagReq = &github.CreateTagReq{}

func newCreateGithubTagCmd() *cobra.Command {
	tagCmd := &cobra.Command{
		Use:   "tag --repo <repo> --tag <tag> --branch <branch>",
		Short: "Create a tag at the head of a branch",
		Long: `Create a lightweight tag at the head of a branch if it does not exist.
If the tag already exists we leave it alone as long as the branch contains its
commit, even if the branch has moved on since. Otherwise we fail.`,
		RunE: runCreateGithubTagCmd,
		Args: cobra.NoArgs,
	}

	tagCmd.PersistentFlags().StringVarP(
		&createGithubTagReq.Owner, "owner", "o", "hashicorp",
		"The Github organization",
	)
	tagCmd.PersistentFlags().StringVarP(
		&createGithubTagReq.Repo, "repo", "r", "",
		"The Github repository. Private repositories require auth via a GITHUB_TOKEN env var",
	)
	tagCmd.PersistentFlags().StringVar(
		&createGithubTagReq.Tag, "tag", "",
		"The tag to create",
	)
	tagCmd.PersistentFlags().StringVar(
		&createGithubTagReq.Branch, "branch", "",
		"The branch whose head we tag",
	)
	tagCmd.PersistentFlags().BoolVar(
		&createGithubTagReq.WriteToGithubOutput, "github-output", false,
		"Whether or not to write 'tag' to $GITHUB_OUTPUT",
	)

	for _, flag := range []string{"repo", "tag", "branch"} {
		if err := tagCmd.MarkPersistentFlagRequired(flag); err != nil {
			panic(err)
		}
	}

	return tagCmd
}

func runCreateGithubTagCmd(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true // Don't spam the usage on failure

	res, err := createGithubTagReq.Run(cmd.Context(), githubCmdState.GithubV3)
	if err != nil {
		return fmt.Errorf("creating tag: %w", err)
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

	if createGithubTagReq.WriteToGithubOutput {
		return writeToGithubOutput("tag", b)
	}

	return nil
}
