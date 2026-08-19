// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	gh "github.com/google/go-github/v83/github"
	ghpkg "github.com/hashicorp/vault/tools/pipeline/internal/pkg/github"
	"github.com/shurcooL/githubv4"
	"github.com/spf13/cobra"
)

type githubCommandState struct {
	GithubV3   *gh.Client
	GithubV4   *githubv4.Client
	GithubHost string
}

var githubCmdState = &githubCommandState{
	GithubV3: gh.NewClient(nil),
	GithubV4: githubv4.NewClient(nil),
}

func newGithubCmd() *cobra.Command {
	githubCmd := &cobra.Command{
		Use:   "github",
		Short: "Github commands",
		Long:  "Github commands",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			token, tokenSet := os.LookupEnv("GITHUB_TOKEN")
			if !tokenSet {
				slog.Default().WarnContext(cmd.Context(), "GITHUB_TOKEN has not been set. While not always required for read actions on public repositories you're likely to get throttled without it")
			}

			// Both clients are constructed via the shared ghpkg helpers so
			// auth transport setup (oauth2.Transport) is never duplicated here.
			var err error
			githubCmdState.GithubV3, err = ghpkg.NewClient(githubCmdState.GithubHost, token, nil)
			if err != nil {
				return fmt.Errorf("creating GitHub v3 REST client: %w", err)
			}

			githubCmdState.GithubV4, err = ghpkg.NewV4Client(githubCmdState.GithubHost, token, nil)
			if err != nil {
				return fmt.Errorf("creating GitHub v4 GraphQL client: %w", err)
			}

			return nil
		},
	}

	githubCmd.PersistentFlags().StringVar(&githubCmdState.GithubHost, "github-host", "github.com", "GitHub API host (use for GitHub Enterprise Server)")

	githubCmd.AddCommand(newGithubCheckCmd())
	githubCmd.AddCommand(newGithubCloseCmd())
	githubCmd.AddCommand(newCollectArtifactShasCmd())
	githubCmd.AddCommand(newGithubCopyCmd())
	githubCmd.AddCommand(newGithubCreateCmd())
	githubCmd.AddCommand(newGithubFindCmd())
	githubCmd.AddCommand(newGithubListCmd())
	githubCmd.AddCommand(newGithubSyncCmd())

	return githubCmd
}

func writeToGithubOutput(key string, bytes []byte) error {
	devPath, ok := os.LookupEnv("GITHUB_OUTPUT")
	if !ok {
		return fmt.Errorf("$GITHUB_OUTPUT has not been set. Cannot write %s to it", key)
	}

	expanded, err := filepath.Abs(devPath)
	if err != nil {
		return fmt.Errorf("failed to expand $GITHUB_OUTPUT path: %w", err)
	}

	dev, err := os.OpenFile(expanded, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("failed to open $GITHUB_OUTPUT for writing: %w", err)
	}
	defer func() { _ = dev.Close() }()

	_, err = dev.Write(append(append([]byte(key+"="), bytes...), []byte("\n")...))
	if err != nil {
		return fmt.Errorf("failed to write key %s to $GITHUB_OUTPUT: %w", key, err)
	}

	return nil
}
