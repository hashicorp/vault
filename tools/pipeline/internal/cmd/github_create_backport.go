// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
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
		Use:   "backport 1234 [--release.release/pipeline.hcl]",
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

	backportCmd.PersistentFlags().StringSliceVarP(&createGithubBackportState.ceAllowInactive, "ce-allow-inactive-groups", "a", []string{}, "Change file groups that should be allowed to backport to inactive CE branches")
	backportCmd.PersistentFlags().StringVar(&createGithubBackportState.req.CEBranchPrefix, "ce-branch-prefix", "ce", "The branch name prefix")
	backportCmd.PersistentFlags().StringSliceVarP(&createGithubBackportState.ceExclude, "ce-exclude-groups", "e", []string{"enterprise"}, "Change file groups that should be excluded from the backporting to CE branches")
	backportCmd.PersistentFlags().StringVar(&createGithubBackportState.req.BaseOrigin, "base-origin", "origin", "The name to use for the base remote origin")
	backportCmd.PersistentFlags().StringVarP(&createGithubBackportState.req.Owner, "owner", "o", "hashicorp", "The Github organization")
	backportCmd.PersistentFlags().StringVarP(&createGithubBackportState.req.Repo, "repo", "r", "vault-enterprise", "The Github repository. Private repositories require auth via a GITHUB_TOKEN env var")
	backportCmd.PersistentFlags().StringVarP(&createGithubBackportState.req.RepoDir, "repo-dir", "d", "", "The path to the vault repository dir. If not set a temporary directory will be used")

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

	// Pass along our configuration decode responses. The request will handle
	// scenarios as necessary.
	createGithubBackportState.req.ConfigDecodeRes = rootCfg.configDecodeRes
	createGithubBackportState.req.VersionsDecodeRes = rootCfg.versionsDecodeRes

	if createGithubBackportState.req.CEAllowInactiveGroups == nil {
		createGithubBackportState.req.CEAllowInactiveGroups = changed.FileGroups{}
	}
	for _, ig := range createGithubBackportState.ceAllowInactive {
		createGithubBackportState.req.CEAllowInactiveGroups = createGithubBackportState.req.CEAllowInactiveGroups.Add(changed.FileGroup(ig))
	}

	if createGithubBackportState.req.CEExclude == nil {
		createGithubBackportState.req.CEExclude = changed.FileGroups{}
	}
	for _, eg := range createGithubBackportState.ceExclude {
		createGithubBackportState.req.CEExclude = createGithubBackportState.req.CEExclude.Add(changed.FileGroup(eg))
	}

	res := createGithubBackportState.req.Run(cmd.Context(), githubCmdState.GithubV3, rootCfg.git)
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
