// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/github"
	"github.com/spf13/cobra"
)

var collectArtifactShasReq = github.CollectArtifactShasReq{}

// pullRequestsJSON holds the raw JSON string from --pull-requests before it is
// decoded into collectArtifactShasReq.PullRequests.
var pullRequestsJSON string

// versionsRaw holds the raw space/comma-separated version string from
// --versions before it is split into collectArtifactShasReq.Versions.
var versionsRaw string

func newCollectArtifactShasCmd() *cobra.Command {
	collectArtifactShasCmd := &cobra.Command{
		Use:   "collect-artifact-shas",
		Short: "Collect artifact SHAs for merged release PRs",
		Long:  "Collect the merge-commit SHAs for a set of release pull requests and build the Slack notification message.",
		RunE:  runCollectArtifactShasCmd,
	}

	collectArtifactShasCmd.PersistentFlags().StringVar(
		&pullRequestsJSON, "pull-requests", "",
		`JSON array of pull-request objects, e.g. '[{"target_branch":"origin/release/2.0.x+ent","url":"https://github.com/hashicorp/vault-enterprise/pull/13558"}]'`,
	)
	collectArtifactShasCmd.PersistentFlags().BoolVar(
		&collectArtifactShasReq.DryRunShas, "dry-run", true,
		"Defaults to true (test mode). Set to false to make real GitHub API calls and collect actual merge SHAs.",
	)
	collectArtifactShasCmd.PersistentFlags().StringVar(
		&collectArtifactShasReq.SlackChannel, "slack-channel", "",
		"Slack channel ID for the SHA notification",
	)
	collectArtifactShasCmd.PersistentFlags().StringVar(
		&versionsRaw, "versions", "",
		`Space or comma-separated list of versions corresponding to the pull requests (e.g. "1.22.0-rc1 1.21.3+ent")`,
	)

	return collectArtifactShasCmd
}

func runCollectArtifactShasCmd(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	// Decode --pull-requests JSON into the request struct.
	if pullRequestsJSON != "" {
		if err := json.Unmarshal([]byte(pullRequestsJSON), &collectArtifactShasReq.PullRequests); err != nil {
			return fmt.Errorf("decoding --pull-requests JSON: %w", err)
		}
	}

	// Split --versions into a slice, tolerating spaces and commas.
	if versionsRaw != "" {
		normalized := strings.ReplaceAll(versionsRaw, ",", " ")
		collectArtifactShasReq.Versions = append(collectArtifactShasReq.Versions, strings.Fields(normalized)...)
	}

	res, err := collectArtifactShasReq.Run(context.TODO(), githubCmdState.GithubV3)
	if err != nil {
		return fmt.Errorf("collecting artifact SHAs: %w", err)
	}

	switch rootCfg.format {
	case "json":
		b, err := json.Marshal(res)
		if err != nil {
			return fmt.Errorf("marshaling response to JSON: %w", err)
		}
		fmt.Println(string(b))
	default:
		fmt.Println(res.SlackMsg)
	}

	return nil
}
