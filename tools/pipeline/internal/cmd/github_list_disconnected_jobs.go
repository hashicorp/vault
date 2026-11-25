// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/github"
	"github.com/spf13/cobra"
)

func newGithubListDisconnectedJobsCmd(state *githubCommandState) *cobra.Command {
	req := &github.ListDisconnectedJobsReq{}

	cmd := &cobra.Command{
		Use:   "disconnected-jobs <run-id>",
		Short: "List the failed jobs of a workflow run and whether or not their runner disconnected",
		Long: `List the failed jobs of a workflow run and whether or not their runner disconnected.

Github doesn't tell us when a runner disconnects but we can usually spot it. If a
job has completed and failed but some of its steps are still queued, pending, or
in progress, the runner most likely went away before it could finish them.

The run doesn't have to be completed so it's safe to list the jobs of a run that
is still in progress. Use 'github rerun-if-disconnected' if you want to rerun them.

We'll inspect the latest attempt of the run unless you give us an --attempt. E.g.
if the run has already been rerun use --attempt 1 to inspect the first attempt.`,
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

			req.Owner, req.Repo = state.Owner, state.Repo
			req.RunID, err = parseGithubWorkflowRunIDArg(args[0])
			if err != nil {
				return err
			}

			res, err := req.Run(cmd.Context(), state.GithubV3)
			if err != nil {
				return fmt.Errorf("listing disconnected jobs: %w", err)
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
			if _, err = fmt.Fprintln(cmd.OutOrStdout(), output); err != nil {
				return err
			}

			if req.WriteToGithubOutput {
				b, err := res.ToGithubOutput()
				if err != nil {
					return err
				}

				return writeToGithubOutput("disconnected-jobs", b)
			}

			return nil
		},
	}

	cmd.Flags().IntVar(&req.Attempt, "attempt", 0, "The attempt of the run to inspect. Defaults to the latest attempt")
	cmd.Flags().BoolVar(&req.WriteToGithubOutput, "github-output", false, "Whether or not to write 'disconnected-jobs' to $GITHUB_OUTPUT")

	return cmd
}

// runIDArg ensures that we've been given a single valid workflow run ID.
func runIDArg(cmd *cobra.Command, args []string) error {
	if err := cobra.ExactArgs(1)(cmd, args); err != nil {
		return err
	}
	_, err := parseGithubWorkflowRunIDArg(args[0])

	return err
}

// parseGithubWorkflowRunIDArg parses a Github Actions workflow run ID from a command argument.
func parseGithubWorkflowRunIDArg(arg string) (int64, error) {
	runID, err := strconv.ParseInt(arg, 10, 64)
	if err != nil || runID <= 0 {
		return 0, fmt.Errorf("run ID must be a positive int64: %q", arg)
	}

	return runID, nil
}
