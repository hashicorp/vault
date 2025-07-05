// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/changed"
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/github"
	"github.com/spf13/cobra"
)

var createGithubBackportState struct {
	req             github.CreateBackportReq
	ceExclude       []string
	ceAllowInactive []string
}

func newCreateGithubBackportCmd() *cobra.Command {
	backportCmd := &cobra.Command{
		Use:   "backport 1234",
		Short: "Create a backport pull request from another pull request",
		Long:  "Create a backport pull request from another pull request",
		RunE:  runCreateGithubBackportCmd,
		Args: func(cmd *cobra.Command, args []string) error {
			switch len(args) {
			case 1:
				pr, err := strconv.ParseUint(args[0], 10, 0)
				if err != nil {
					return fmt.Errorf("invalid pull number: %s: %w", args[0], err)
				}
				if pr <= math.MaxUint32 {
					createGithubBackportState.req.PullNumber = uint(pr)
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

	backportCmd.PersistentFlags().StringSliceVarP(&createGithubBackportState.ceAllowInactive, "ce-allow-inactive-groups", "a", []string{"docs", "changelog", "pipeline"}, "Change file groups that should be allowed to backport to inactive CE branches")
	backportCmd.PersistentFlags().StringVar(&createGithubBackportState.req.CEBranchPrefix, "ce-branch-prefix", "ce", "The branch name prefix")
	backportCmd.PersistentFlags().StringSliceVarP(&createGithubBackportState.ceExclude, "ce-exclude-groups", "e", []string{"enterprise"}, "Change file groups that should be excluded from the backporting to CE branches")
	backportCmd.PersistentFlags().StringVar(&createGithubBackportState.req.BaseOrigin, "base-origin", "origin", "The name to use for the base remote origin")
	backportCmd.PersistentFlags().StringVarP(&createGithubBackportState.req.Owner, "owner", "o", "hashicorp", "The Github organization")
	backportCmd.PersistentFlags().StringVarP(&createGithubBackportState.req.Repo, "repo", "r", "vault-enterprise", "The Github repository. Private repositories require auth via a GITHUB_TOKEN env var")
	backportCmd.PersistentFlags().StringVarP(&createGithubBackportState.req.RepoDir, "repo-dir", "d", "", "The path to the vault repository dir. If not set a temporary directory will be used")
	backportCmd.PersistentFlags().StringVarP(&createGithubBackportState.req.ReleaseVersionConfigPath, "releases-version-path", "m", "", "The path to .release/versions.hcl")
	backportCmd.PersistentFlags().UintVar(&createGithubBackportState.req.ReleaseRecurseDepth, "recurse", 3, "If no path to a config file is given, recursively search backwards for it and stop at root or until we've his the configured depth.")

	// NOTE: The following are technically flags but they only for testing testing
	// the command before we cut over to new utility.
	backportCmd.PersistentFlags().StringVar(&createGithubBackportState.req.EntBranchPrefix, "ent-branch-prefix", "", "The ent branch name prefix. Only used for testing before migration to the new workflow")
	backportCmd.PersistentFlags().StringVar(&createGithubBackportState.req.BackportLabelPrefix, "backport-label-prefix", "backport", "The name to use for the base remote origin")

	err := backportCmd.PersistentFlags().MarkHidden("ent-branch-prefix")
	if err != nil {
		panic(err)
	}

	err = backportCmd.PersistentFlags().MarkHidden("backport-label-prefix")
	if err != nil {
		panic(err)
	}

	return backportCmd
}

func runCreateGithubBackportCmd(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true // Don't spam the usage on failure

	for i, ig := range createGithubBackportState.ceAllowInactive {
		if i == 0 && createGithubBackportState.req.CEAllowInactiveGroups == nil {
			createGithubBackportState.req.CEAllowInactiveGroups = changed.FileGroups{}
		}
		createGithubBackportState.req.CEAllowInactiveGroups = createGithubBackportState.req.CEAllowInactiveGroups.Add(changed.FileGroup(ig))
	}

	for i, eg := range createGithubBackportState.ceExclude {
		if i == 0 && createGithubBackportState.req.CEExclude == nil {
			createGithubBackportState.req.CEExclude = changed.FileGroups{}
		}
		createGithubBackportState.req.CEExclude = createGithubBackportState.req.CEExclude.Add(changed.FileGroup(eg))
	}

	res := createGithubBackportState.req.Run(context.TODO(), githubCmdState.Github, githubCmdState.Git)
	if res == nil {
		res = &github.CreateBackportRes{}
	}
	if err := res.Err(); err != nil {
		res.ErrorMessage = err.Error()
	}

	switch rootCfg.format {
	case "json":
		b, err := res.ToJSON()
		if err != nil {
			return errors.Join(res.Err(), err)
		}
		fmt.Println(string(b))
	default:
		fmt.Println(res.ToTable().Render())
	}

	return res.Err()
}
