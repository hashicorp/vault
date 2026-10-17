// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/github"
	"github.com/spf13/cobra"
)

func newGithubRerunIfDisconnectedCmd(state *githubCommandState) *cobra.Command {
	return &cobra.Command{
		Use:   "rerun-if-disconnected <run-id>",
		Short: "Rerun failed jobs when unfinished steps suggest a runner disconnect",
		Long: `Check failed workflow jobs for steps that are still queued, pending,
or in progress even though the job has completed. This usually means the runner
disconnected before it could finish the job.

If we find any matches, request one rerun of all failed jobs and their dependents,
including unrelated failures. Otherwise exit 0 without requesting a rerun.
The workflow controls retry limits. We don't wait for the rerun to finish.`,
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.ExactArgs(1)(cmd, args); err != nil {
				return err
			}
			runID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil || runID <= 0 {
				return fmt.Errorf("run ID must be a positive int64: %q", args[0])
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			format, err := cmd.Flags().GetString("format")
			if err != nil {
				return err
			}
			switch format {
			case "table", "markdown", "json":
			default:
				return fmt.Errorf("unsupported format: %s (expected table, markdown, or json)", format)
			}
			runID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("parsing run ID: %w", err)
			}
			req := &github.RerunIfDisconnectedReq{Owner: state.Owner, Repo: state.Repo, RunID: runID}
			res, err := req.Run(cmd.Context(), state.GithubV3)
			if err != nil {
				return fmt.Errorf("rerun if disconnected: %w", err)
			}
			var output string
			switch format {
			case "json":
				b, err := res.ToJSON()
				if err != nil {
					return err
				}
				output = string(b)
			case "markdown":
				output = res.ToMarkdown()
			default:
				output = res.ToTable()
			}
			_, err = fmt.Fprintln(cmd.OutOrStdout(), output)
			return err
		},
	}
}
