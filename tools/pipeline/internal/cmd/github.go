// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/go-github/v68/github"
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/git"
	"github.com/spf13/cobra"
)

type githubCommandState struct {
	Github *github.Client
	Git    *git.Client
}

var githubCmdState = &githubCommandState{
	Github: github.NewClient(nil),
	Git:    git.NewClient(git.WithLoadTokenFromEnv()),
}

func newGithubCmd() *cobra.Command {
	githubCmd := &cobra.Command{
		Use:   "github",
		Short: "Github commands",
		Long:  "Github commands",
	}
	githubCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if token, set := os.LookupEnv("GITHUB_TOKEN"); set {
			githubCmdState.Github = githubCmdState.Github.WithAuthToken(token)
		} else {
			fmt.Println("\x1b[1;33;49mWARNING\x1b[0m: GITHUB_TOKEN has not been set. While not always required for read actions on public repositories you're likely to get throttled without it")
		}
		return nil
	}
	githubCmd.AddCommand(newGithubCopyCmd())
	githubCmd.AddCommand(newGithubCreateCmd())
	githubCmd.AddCommand(newGithubListCmd())

	return githubCmd
}

func writeToGithubOutput(key string, bytes []byte) error {
	devPath, ok := os.LookupEnv("GITHUB_OUTPUT")
	if !ok {
		return errors.New("$GITHUB_OUTPUT has not been set. Cannot write changed files to it")
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

	_, err = dev.Write(append([]byte(key+"="), bytes...))
	if err != nil {
		return fmt.Errorf("failed to write key %s to $GITHUB_OUTPUT: %w", key, err)
	}

	return nil
}
