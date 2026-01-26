// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/go-github/v81/github"
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/git"
	"github.com/shurcooL/githubv4"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

type githubCommandState struct {
	Git      *git.Client
	GithubV3 *github.Client
	GithubV4 *githubv4.Client
}

var githubCmdState = &githubCommandState{
	GithubV3: github.NewClient(nil),
	GithubV4: githubv4.NewClient(nil),
	Git:      git.NewClient(git.WithLoadTokenFromEnv()),
}

func newGithubCmd() *cobra.Command {
	githubCmd := &cobra.Command{
		Use:   "github",
		Short: "Github commands",
		Long:  "Github commands",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if token, set := os.LookupEnv("GITHUB_TOKEN"); set {
				githubCmdState.GithubV3 = githubCmdState.GithubV3.WithAuthToken(token)
				githubCmdState.GithubV4 = githubv4.NewClient(
					oauth2.NewClient(context.Background(),
						oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}),
					),
				)
			} else {
				fmt.Println("\x1b[1;33;49mWARNING\x1b[0m: GITHUB_TOKEN has not been set. While not always required for read actions on public repositories you're likely to get throttled without it")
			}

			return nil
		},
	}

	githubCmd.AddCommand(newGithubCheckCmd())
	githubCmd.AddCommand(newGithubCloseCmd())
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
