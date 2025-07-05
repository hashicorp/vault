// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/github"
	"github.com/spf13/cobra"
)

var copyGithubPullRequestReq = github.CopyPullRequestReq{}

func newCopyGithubPullRequestCmd() *cobra.Command {
	copyPRCmd := &cobra.Command{
		Use:   "pr [number]",
		Short: "Copy a pull request",
		Long:  "Copy a pull request from the Community repository to the Enterprise repository",
		RunE:  runCopyGithubPullRequestCmd,
		Args: func(cmd *cobra.Command, args []string) error {
			switch len(args) {
			case 1:
				pr, err := strconv.ParseUint(args[0], 10, 0)
				if err != nil {
					return fmt.Errorf("invalid pull number: %s: %w", args[0], err)
				}
				if pr <= math.MaxUint32 {
					copyGithubPullRequestReq.PullNumber = uint(pr)
				} else {
					return fmt.Errorf("invalid pull number: %s: number is too large", args[0])
				}
				return nil
			case 0:
				return errors.New("no pull request number has been provided")
			default:
				return fmt.Errorf("invalid arguments: only pull request number is expected, received %d arguments: %v", len(args), args)
			}
		},
	}

	copyPRCmd.PersistentFlags().StringVar(&copyGithubPullRequestReq.FromOrigin, "from-origin", "ce", "The name to use for the base remote origin")
	copyPRCmd.PersistentFlags().StringVar(&copyGithubPullRequestReq.FromOwner, "from-owner", "hashicorp", "The Github organization")
	copyPRCmd.PersistentFlags().StringVar(&copyGithubPullRequestReq.FromRepo, "from-repo", "vault", "The CE Github repository to copy the PR from")
	copyPRCmd.PersistentFlags().StringVar(&copyGithubPullRequestReq.ToOrigin, "to-origin", "origin", "The name to use for the base remote origin")
	copyPRCmd.PersistentFlags().StringVar(&copyGithubPullRequestReq.ToOwner, "to-owner", "hashicorp", "The Github organization")
	copyPRCmd.PersistentFlags().StringVar(&copyGithubPullRequestReq.ToRepo, "to-repo", "vault-enterprise", "The Github repository. Private repositories require auth via a GITHUB_TOKEN env var")
	copyPRCmd.PersistentFlags().StringVarP(&copyGithubPullRequestReq.RepoDir, "repo-dir", "d", "", "The path to the vault repository dir. If not set a temporary directory will be used")
	copyPRCmd.PersistentFlags().StringVar(&copyGithubPullRequestReq.EntBranchSuffix, "ent-branch-suffix", "+ent", "The release branch suffix for enterprise branches")

	return copyPRCmd
}

func runCopyGithubPullRequestCmd(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true // Don't spam the usage on failure

	res, err := copyGithubPullRequestReq.Run(context.TODO(), githubCmdState.Github, githubCmdState.Git)

	switch rootCfg.format {
	case "json":
		b, err1 := res.ToJSON()
		if err1 != nil {
			return errors.Join(err, err1)
		}
		fmt.Println(string(b))
	default:
		fmt.Println(res.ToTable(err))
	}

	return err
}
