// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"fmt"
	"maps"
	"net/http"
	"slices"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/github"
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/releases"
	"github.com/spf13/cobra"
)

var checkGoVersionReq = &github.CheckGoVersionReq{}

var includeMainGoVersion bool

func newGithubCheckGoVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "go-version [branches...]",
		Short: "Check the .go-version file on each branch against the latest stable Go versions",
		Long: `Fetches the .go-version file from each branch via the GitHub API and
compares it against the two most recent stable Go releases fetched from go.dev.
A branch is "up-to-date" if it uses either the latest (N) or previous (N-1)
stable release. Reports each branch as "up-to-date", "out-of-date", or "file-missing".

Branches can be specified as positional arguments. When no branches are
provided, the command checks all active release branches derived from
.release/versions.hcl. Use --include-main to also check the main branch.

When --github-output is set, writes vault_branch_go_statuses (JSON array),
version_mismatch_count, and file_missing_count to $GITHUB_OUTPUT.`,
		RunE: runGithubCheckGoVersionCmd,
		Args: cobra.ArbitraryArgs,
	}

	cmd.PersistentFlags().StringVarP(&checkGoVersionReq.Owner, "owner", "o", "hashicorp", "GitHub organization")
	cmd.PersistentFlags().StringVarP(&checkGoVersionReq.Repo, "repo", "r", "", "GitHub repository")
	cmd.PersistentFlags().BoolVar(&checkGoVersionReq.WriteToGithubOutput, "github-output", false, "Write outputs to $GITHUB_OUTPUT")
	cmd.PersistentFlags().BoolVar(&includeMainGoVersion, "include-main", false, "Include 'main' branch when auto-resolving branches from versions.hcl")

	if err := cmd.MarkPersistentFlagRequired("repo"); err != nil {
		panic(err)
	}

	return cmd
}

func runGithubCheckGoVersionCmd(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	checkGoVersionReq.HTTPClient = http.DefaultClient

	checkGoVersionReq.Branches = branchesToCheck(args, rootCfg.versionsDecodeRes, includeMainGoVersion)

	res, err := checkGoVersionReq.Run(cmd.Context(), githubCmdState.GithubV3)
	if err != nil {
		return fmt.Errorf("check go-version: %w", err)
	}

	switch rootCfg.format {
	case "json":
		b, err := res.ToJSON()
		if err != nil {
			return err
		}
		fmt.Println(string(b))
	default:
		fmt.Println(res.ToTable())
	}

	if checkGoVersionReq.WriteToGithubOutput {
		output, err := res.ToGithubOutput()
		if err != nil {
			return err
		}
		for _, k := range slices.Sorted(maps.Keys(output)) {
			if err := writeToGithubOutput(k, []byte(output[k])); err != nil {
				return err
			}
		}
	}

	return nil
}

// branchesToCheck returns the positional args as-is when provided; otherwise it
// derives the branches from the decoded versions config.
func branchesToCheck(args []string, dec *releases.DecodeRes, includeMain bool) []string {
	if len(args) > 0 {
		return args
	}

	var cfg *releases.VersionsConfig
	if dec != nil {
		cfg = dec.Config
	}

	return releases.ActiveEnterpriseBranches(cfg, includeMain)
}
