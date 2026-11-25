// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/github"
	"github.com/spf13/cobra"
)

func newGithubRerunIfDisconnectedCmd(state *githubCommandState) *cobra.Command {
	req := &github.RerunIfDisconnectedReq{}
	var timeout time.Duration

	cmd := &cobra.Command{
		Use:   "rerun-if-disconnected <run-id>",
		Short: "Rerun the failed jobs of a workflow run if their runner disconnected",
		Long: `Rerun the failed jobs of a workflow run if any of them look like their runner
disconnected. We use the same heuristic as 'github list disconnected-jobs'. If we
find any we'll request a single rerun of all failed jobs and their dependents,
which includes any unrelated failures.

Github will only rerun a completed run so you can't rerun a run from inside of it.
Use --wait to poll the run until it has completed. We'll skip the rerun if the run
didn't fail, if it's no longer on the --attempt that you gave us, or if a newer
run of the same workflow exists for the same branch and event. We don't wait for
the rerun to finish.`,
		Args: runIDArg,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true // Don't spam the usage on failure

			format, err := cmd.Flags().GetString("format")
			if err != nil {
				return err
			}
			switch format {
			case "table", "markdown", "json":
			default:
				return fmt.Errorf("unsupported format: %s (expected table, markdown, or json)", format)
			}
			if req.Wait && timeout <= 0 {
				return errors.New("timeout must be positive when waiting")
			}

			req.Owner, req.Repo = state.Owner, state.Repo
			req.RunID, err = parseGithubWorkflowRunIDArg(args[0])
			if err != nil {
				return err
			}

			ctx := cmd.Context()
			if req.Wait {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeoutCause(
					ctx, timeout, fmt.Errorf("timed out after %s", timeout),
				)
				defer cancel()
			}

			res, err := req.Run(ctx, state.GithubV3)
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

	cmd.Flags().BoolVar(&req.Wait, "wait", false, "Wait for the run to complete before we check it")
	cmd.Flags().DurationVar(&req.PollInterval, "poll-interval", 30*time.Second, "How often to check on the run while we wait")
	cmd.Flags().DurationVar(&timeout, "timeout", 30*time.Minute, "How long to wait for the run to complete")
	cmd.Flags().IntVar(&req.Attempt, "attempt", 0, "Only rerun the run if it's still on this attempt. E.g. 1")

	return cmd
}
