// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	ghclient "github.com/google/go-github/v68/github"
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/github"
	"github.com/spf13/cobra"
)

var listGithubWorkflowRuns = &github.ListWorkflowRunsReq{}

func newGithubListRunCmd() *cobra.Command {
	listRuns := &cobra.Command{
		Use:   "list-workflow-runs [WORKFLOW_NAME]",
		Short: "List workflow runs",
		Long:  "List Github Actions workflow runs for a given workflow. Be sure to use filter arguments to reduce the search, otherwise you'll likely hit your API limit.",
		RunE:  runListGithubWorkflowsCmd,
		Args: func(cmd *cobra.Command, args []string) error {
			switch len(args) {
			case 1:
				listGithubWorkflowRuns.WorkflowName = args[0]
				return nil
			case 0:
				return errors.New("no workflow name argument has been provided")
			default:
				return fmt.Errorf("expected a single workflow name as an argument, received (%d): %v", len(args), args)
			}
		},
	}

	listRuns.PersistentFlags().StringVarP(&listGithubWorkflowRuns.Actor, "actor", "a", "", "Filter using a specific Github actor")
	listRuns.PersistentFlags().StringVarP(&listGithubWorkflowRuns.Branch, "branch", "b", "", "Filter using a specific Github branch")
	listRuns.PersistentFlags().Int64VarP(&listGithubWorkflowRuns.CheckSuiteID, "check-suite-id", "c", 0, "Filter using a specific Github check suite")
	listRuns.PersistentFlags().BoolVar(&listGithubWorkflowRuns.Compact, "compact", true, "When given a status filter, only fetch data for workflows, jobs, checks, and annotations that match our status and/or conclusion. Disabling compact mode with a large query range might result in Github throttling the requests.")
	listRuns.PersistentFlags().StringVarP(&listGithubWorkflowRuns.DateQuery, "date-query", "d", fmt.Sprintf("%s..*", time.Now().Add(-168*time.Hour).Format(time.DateOnly)), "Filter using a date range query. It supports the Github ISO8601-ish date range query format. Default is newer than one week ago")
	listRuns.PersistentFlags().StringVarP(&listGithubWorkflowRuns.Event, "event", "e", "", "Filter using a workflow triggered by an event type. E.g. push, pull_request, issue")
	listRuns.PersistentFlags().BoolVarP(&listGithubWorkflowRuns.IncludePRs, "include-prs", "p", false, "Include workflow runs triggered via pull requests")
	listRuns.PersistentFlags().StringVarP(&listGithubWorkflowRuns.Owner, "owner", "o", "hashicorp", "The Github organization")
	listRuns.PersistentFlags().StringVarP(&listGithubWorkflowRuns.Repo, "repo", "r", "vault", "The Github repository. Private repositories require auth via a GITHUB_TOKEN env var")
	listRuns.PersistentFlags().StringVar(&listGithubWorkflowRuns.Sha, "sha", "", "Filter based on the HEAD SHA associated with the workflow run")
	listRuns.PersistentFlags().StringVar(&listGithubWorkflowRuns.Status, "status", "", "Filter by a given run status. For example: completed, cancelled, failure, skipped, success, in_progress")

	return listRuns
}

func runListGithubWorkflowsCmd(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true // Don't spam the usage on failure

	client := ghclient.NewClient(nil)
	if token, set := os.LookupEnv("GITHUB_TOKEN"); set {
		client = client.WithAuthToken(token)
	} else {
		fmt.Println("\x1b[1;33;49mWARNING\x1b[0m: GITHUB_TOKEN has not been set. While not required for public repositories you're likely to get throttled without it")
	}

	res, err := listGithubWorkflowRuns.Run(context.TODO(), client)
	if err != nil {
		return fmt.Errorf("listing github workflow failures: %w", err)
	}

	switch githubCmdFlags.Format {
	case "json":
		b, err := json.Marshal(res)
		if err != nil {
			return fmt.Errorf("marshaling response to JSON: %w", err)
		}
		fmt.Println(string(b))
	default:
		for _, run := range res.Runs {
			summary, err := run.Summary()
			if err != nil {
				return fmt.Errorf("generating workflow run response summary: %w", err)
			}
			fmt.Println(summary)
		}
	}

	return err
}
