// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"context"
	"fmt"
	"log/slog"

	gh "github.com/google/go-github/v81/github"
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
	slog.Default().DebugContext(slogctx.Append(ctx,
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
func getWorkflowRuns(
	ctx context.Context,
	client *gh.Client,
	owner string,
	repo string,
	id int64,
	opts *gh.ListWorkflowRunsOptions,
) ([]*WorkflowRun, error) {
	var runs []*WorkflowRun
	opts.ListOptions = gh.ListOptions{PerPage: PerPageMax}

	// By default our status will be "success" which elimates in_progress runs.
	// Instead, we'll try both so that we're sure to include what's actually
	// running along with historical runs.
	for _, status := range []string{"", "success", "in_progress"} {
		var runsForStatus []*WorkflowRun
		for {
			opts.Status = status
			slog.Default().DebugContext(slogctx.Append(ctx,
				slog.String("owner", owner),
				slog.String("repo", repo),
				slog.Int64("workflow-id", id),
				slog.String("query-status", opts.Status),
			), "getting github actions workflow runs")

			wfrs, res, err := client.Actions.ListWorkflowRunsByID(ctx, owner, repo, id, opts)
			if err != nil {
				return nil, err
			}

			for _, r := range wfrs.WorkflowRuns {
				runsForStatus = append(runsForStatus, &WorkflowRun{Run: r})
			}

			if res.NextPage == 0 {
				if len(runsForStatus) > 0 {
					slog.Default().DebugContext(slogctx.Append(ctx,
						slog.String("owner", owner),
						slog.String("repo", repo),
						slog.Int64("workflow-id", id),
						slog.String("query-status", opts.Status),
						slog.Int("count", len(runsForStatus)),
					), "found github actions workflow runs")
				} else {
					slog.Default().DebugContext(slogctx.Append(ctx,
						slog.String("owner", owner),
						slog.String("repo", repo),
						slog.Int64("workflow-id", id),
						slog.String("query-status", opts.Status),
					), "no github actions workflow runs found for status")
				}
				runs = append(runs, runsForStatus...)
				break
			}

			opts.ListOptions.Page = res.NextPage
		}
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
	slog.Default().DebugContext(slogctx.Append(ctx,
		slog.String("owner", owner),
		slog.String("repo", repo),
		slog.Int64("run-id", id),
	), "getting github actions workflow run artifacts")

	opts := &gh.ListOptions{PerPage: PerPageMax}
	artifacts := gh.ArtifactList{}

	defer func() {
		if count := artifacts.GetTotalCount(); count > 0 {
			slog.Default().DebugContext(slogctx.Append(ctx,
				slog.String("owner", owner),
				slog.String("repo", repo),
				slog.Int64("run-id", id),
				slog.Int64("count", count),
			), "found workflow run artifacts")
		} else {
			slog.Default().DebugContext(slogctx.Append(ctx,
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
