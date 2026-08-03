// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	gh "github.com/google/go-github/v83/github"
	slogctx "github.com/veqryn/slog-context"
)

// getWorkflow attempts to locate the workflow associated with our workflow name.
func getWorkflow(
	ctx context.Context,
	client *gh.Client,
	owner string,
	repo string,
	name string,
) (*gh.Workflow, error) {
	slog.Default().DebugContext(slogctx.Append(
		ctx,
		slog.String("owner", owner),
		slog.String("repo", repo),
		slog.String("name", name),
	), "getting github actions workflow")

	opts := &gh.ListOptions{PerPage: PerPageMax}
	for {
		wfs, res, err := client.Actions.ListWorkflows(ctx, owner, repo, opts)
		if err != nil {
			return nil, err
		}

		for _, wf := range wfs.Workflows {
			if wf.GetName() == name {
				return wf, nil
			}
		}

		if res.NextPage == 0 {
			return nil, fmt.Errorf("no workflow matching %s could be found", name)
		}

		opts.Page = res.NextPage
	}
}

// getWorkflowRuns gets the workflow runs associated with a workflow ID.
// If opts.Status is set, only that status is queried. Otherwise, queries
// multiple statuses.
// maxRuns limits the total number of runs returned (0 = unlimited).
func getWorkflowRuns(
	ctx context.Context,
	client *gh.Client,
	owner string,
	repo string,
	id int64,
	maxRuns int,
	opts *gh.ListWorkflowRunsOptions,
) ([]*WorkflowRun, error) {
	const maxRetries = 3
	var runs []*WorkflowRun

	// If PerPage is not set, use max
	if opts.ListOptions.PerPage == 0 {
		opts.ListOptions = gh.ListOptions{PerPage: PerPageMax}
	}

	// Calculate max pages to fetch if maxRuns is set
	maxPages := 0 // 0 means unlimited
	if maxRuns > 0 {
		maxPages = (maxRuns + PerPageMax - 1) / PerPageMax
	}

	// Determine which statuses to query
	statuses := []string{"", "success", "in_progress"}
	if opts.Status != "" {
		// If status is explicitly set, only query that status
		statuses = []string{opts.Status}
	}

	// Query each status
	for _, status := range statuses {
		var runsForStatus []*WorkflowRun
		statusOpts := *opts
		statusOpts.Status = status
		statusOpts.ListOptions.Page = 0
		pagesFetched := 0

		for {
			slog.Default().DebugContext(slogctx.Append(
				ctx,
				slog.String("owner", owner),
				slog.String("repo", repo),
				slog.Int64("workflow-id", id),
				slog.String("query-status", statusOpts.Status),
			), "getting github actions workflow runs")

			var wfrs *gh.WorkflowRuns
			var res *gh.Response

			err := retryWithBackoff(ctx, maxRetries, func() error {
				var err error
				wfrs, res, err = client.Actions.ListWorkflowRunsByID(ctx, owner, repo, id, &statusOpts)
				return err
			})
			if err != nil {
				return nil, err
			}

			for _, r := range wfrs.WorkflowRuns {
				runsForStatus = append(runsForStatus, &WorkflowRun{Run: r})
				// Stop if we've reached maxRuns limit
				if maxRuns > 0 && len(runs)+len(runsForStatus) >= maxRuns {
					break
				}
			}

			pagesFetched++

			// Stop if: no more pages, reached max pages, or reached max runs
			shouldStop := res.NextPage == 0 ||
				(maxPages > 0 && pagesFetched >= maxPages) ||
				(maxRuns > 0 && len(runs)+len(runsForStatus) >= maxRuns)

			if shouldStop {
				if len(runsForStatus) > 0 {
					slog.Default().DebugContext(slogctx.Append(
						ctx,
						slog.String("owner", owner),
						slog.String("repo", repo),
						slog.Int64("workflow-id", id),
						slog.String("query-status", statusOpts.Status),
						slog.Int("count", len(runsForStatus)),
					), "found github actions workflow runs")
				} else {
					slog.Default().DebugContext(slogctx.Append(
						ctx,
						slog.String("owner", owner),
						slog.String("repo", repo),
						slog.Int64("workflow-id", id),
						slog.String("query-status", statusOpts.Status),
					), "no github actions workflow runs found for status")
				}
				runs = append(runs, runsForStatus...)
				break
			}

			statusOpts.ListOptions.Page = res.NextPage
		}

		// Stop querying other statuses if we've reached maxRuns
		if maxRuns > 0 && len(runs) >= maxRuns {
			break
		}
	}

	// Trim to maxRuns if we exceeded it
	if maxRuns > 0 && len(runs) > maxRuns {
		runs = runs[:maxRuns]
	}

	return runs, nil
}

// getWorkflowRunArtifacts gets the artifacts associated with a workflow run
func getWorkflowRunArtifacts(
	ctx context.Context,
	client *gh.Client,
	owner string,
	repo string,
	id int64,
) (gh.ArtifactList, error) {
	slog.Default().DebugContext(slogctx.Append(
		ctx,
		slog.String("owner", owner),
		slog.String("repo", repo),
		slog.Int64("run-id", id),
	), "getting github actions workflow run artifacts")

	opts := &gh.ListOptions{PerPage: PerPageMax}
	artifacts := gh.ArtifactList{}

	defer func() {
		if count := artifacts.GetTotalCount(); count > 0 {
			slog.Default().DebugContext(slogctx.Append(
				ctx,
				slog.String("owner", owner),
				slog.String("repo", repo),
				slog.Int64("run-id", id),
				slog.Int64("count", count),
			), "found workflow run artifacts")
		} else {
			slog.Default().DebugContext(slogctx.Append(
				ctx,
				slog.String("owner", owner),
				slog.String("repo", repo),
				slog.Int64("run-id", id),
			), "no workflow run artifacts found")
		}
	}()

	for {
		arts, res, err := client.Actions.ListWorkflowRunArtifacts(ctx, owner, repo, id, opts)
		if err != nil {
			return artifacts, err
		}

		newTotal := artifacts.GetTotalCount() + arts.GetTotalCount()
		artifacts.TotalCount = &newTotal
		artifacts.Artifacts = append(artifacts.Artifacts, arts.Artifacts...)

		if res.NextPage == 0 {
			return artifacts, nil
		}

		opts.Page = res.NextPage
	}
}

// getWorkflowJobsForRun gets all jobs for a workflow run with pagination.
func getWorkflowJobsForRun(
	ctx context.Context,
	client *gh.Client,
	owner string,
	repo string,
	runID int64,
) ([]*gh.WorkflowJob, error) {
	slog.Default().DebugContext(slogctx.Append(
		ctx,
		slog.String("owner", owner),
		slog.String("repo", repo),
		slog.Int64("run_id", runID),
	), "fetching workflow run jobs")

	opts := &gh.ListWorkflowJobsOptions{
		Filter:      "latest",
		ListOptions: gh.ListOptions{PerPage: PerPageMax},
	}

	var allJobs []*gh.WorkflowJob

	for {
		jobs, res, err := client.Actions.ListWorkflowJobs(ctx, owner, repo, runID, opts)
		if err != nil {
			return nil, fmt.Errorf("listing workflow jobs: %w", err)
		}

		allJobs = append(allJobs, jobs.Jobs...)

		if res.NextPage == 0 {
			break
		}

		opts.ListOptions.Page = res.NextPage
	}

	slog.Default().DebugContext(slogctx.Append(
		ctx,
		slog.String("owner", owner),
		slog.String("repo", repo),
		slog.Int64("run_id", runID),
		slog.Int("count", len(allJobs)),
	), "fetched workflow run jobs")

	return allJobs, nil
}

// retryWithBackoff executes a function with exponential backoff retry logic.
// It retries on transient errors (5xx, 429) but fails fast on client errors (4xx).
func retryWithBackoff(
	ctx context.Context,
	maxRetries int,
	operation func() error,
) error {
	for attempt := range maxRetries {
		err := operation()
		if err == nil {
			return nil
		}

		if isClientError(err) {
			return fmt.Errorf("client error (not retrying): %w", err)
		}

		if !isRetryableError(err) {
			return fmt.Errorf("non-retryable error: %w", err)
		}

		if attempt < maxRetries-1 {
			retryDelay := time.Duration(1<<uint(min(attempt, 10))) * 2 * time.Second
			slog.Default().DebugContext(
				ctx, "retrying after error",
				slog.String("error", err.Error()),
				slog.Duration("retry_delay", retryDelay),
				slog.Int("attempt", attempt+1),
			)
			time.Sleep(retryDelay)
			continue
		}

		return fmt.Errorf("max retries exceeded: %w", err)
	}

	return nil
}
