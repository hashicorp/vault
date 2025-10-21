// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/github"
	"github.com/spf13/cobra"
)

var closeCopiedOriginPullRequestReq = github.CloseCopiedOriginPullRequestReq{}

func newCloseGithubCopiedPullRequestCmd() *cobra.Command {
	closeCopiedOriginPRCmd := &cobra.Command{
		Use:   "origin-pull-request [number]",
		Short: "Close the origin pull request of a copied pull request",
		RunE:  runCloseGithubCopiedPullRequestCmd,
		Args:  argsOnlyPRNumber(&closeCopiedOriginPullRequestReq.PullNumber),
	}

	closeCopiedOriginPRCmd.PersistentFlags().StringVarP(&closeCopiedOriginPullRequestReq.Owner, "owner", "o", "hashicorp", "The Github organization")
	closeCopiedOriginPRCmd.PersistentFlags().StringVarP(&closeCopiedOriginPullRequestReq.Repo, "repo", "r", "vault-enterprise", "The Github repository. Private repositories require auth via a GITHUB_TOKEN env var")

	return closeCopiedOriginPRCmd
}

func runCloseGithubCopiedPullRequestCmd(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true // Don't spam the usage on failure

	res, err := closeCopiedOriginPullRequestReq.Run(context.TODO(), githubCmdState.GithubV3, githubCmdState.GithubV4)
	switch rootCfg.format {
	case "json":
		b, err1 := res.ToJSON()
		if err != nil {
			return errors.Join(err, err1)
		}
		fmt.Println(string(b))
	case "markdown":
		tbl := res.ToTable(err)
		tbl.SetTitle("Close Origin Pull Request")
		fmt.Println(tbl.RenderMarkdown())
	default:
		fmt.Println(res.ToTable(err).Render())
	}

	return err
}
