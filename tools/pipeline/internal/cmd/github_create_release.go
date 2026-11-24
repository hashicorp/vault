// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"fmt"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/github"
	"github.com/spf13/cobra"
)

var createGithubReleaseReq = &github.CreateReleaseReq{}

func newCreateGithubReleaseCmd() *cobra.Command {
	releaseCmd := &cobra.Command{
		Use:   "release --repo <repo> --tag <tag> [--prerelease]",
		Short: "Create a GitHub release for an existing tag",
		Long: `Create a GitHub release for an existing tag with GitHub generated release
notes. If the tag already has a release we leave it alone.`,
		RunE: runCreateGithubReleaseCmd,
		Args: cobra.NoArgs,
	}

	releaseCmd.PersistentFlags().StringVarP(
		&createGithubReleaseReq.Owner, "owner", "o", "hashicorp",
		"The Github organization",
	)
	releaseCmd.PersistentFlags().StringVarP(
		&createGithubReleaseReq.Repo, "repo", "r", "",
		"The Github repository. Private repositories require auth via a GITHUB_TOKEN env var",
	)
	releaseCmd.PersistentFlags().StringVar(
		&createGithubReleaseReq.Tag, "tag", "",
		"The existing tag to release",
	)
	releaseCmd.PersistentFlags().BoolVar(
		&createGithubReleaseReq.Prerelease, "prerelease", false,
		"Mark the release as a prerelease and don't make it the latest release",
	)
	releaseCmd.PersistentFlags().BoolVar(
		&createGithubReleaseReq.WriteToGithubOutput, "github-output", false,
		"Whether or not to write 'release' to $GITHUB_OUTPUT",
	)

	for _, flag := range []string{"repo", "tag"} {
		if err := releaseCmd.MarkPersistentFlagRequired(flag); err != nil {
			panic(err)
		}
	}

	return releaseCmd
}

func runCreateGithubReleaseCmd(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true // Don't spam the usage on failure

	res, err := createGithubReleaseReq.Run(cmd.Context(), githubCmdState.GithubV3)
	if err != nil {
		return fmt.Errorf("creating release: %w", err)
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

	if createGithubReleaseReq.WriteToGithubOutput {
		return writeToGithubOutput("release", b)
	}

	return nil
}
