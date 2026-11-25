// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	gh "github.com/google/go-github/v83/github"
)

// RerunIfDisconnectedReq is a request to rerun the failed jobs of a workflow run
// if any of them look like their runner disconnected.
type RerunIfDisconnectedReq struct {
	Owner string
	Repo  string
	RunID int64
	// Wait for the run to complete by polling it every PollInterval. Github will
	// only rerun completed runs. Use the context to limit how long we'll wait.
	Wait         bool
	PollInterval time.Duration
	// Attempt is the attempt of the run that we expect to rerun. If it's set and
	// the run has moved on to another attempt we'll assume someone else has
	// already rerun it and do nothing.
	Attempt int
}

// RerunIfDisconnectedRes is a response to a RerunIfDisconnectedReq.
type RerunIfDisconnectedRes struct {
	*ListDisconnectedJobsRes
	RerunRequested bool `json:"rerun_requested"`
	// SkipReason is why we didn't request a rerun of a run that might have needed one
	SkipReason string `json:"skip_reason,omitempty"`
}

// Run waits for the run if necessary and requests a single rerun of all of its
// failed jobs if we suspect any of them were disconnected. We only rerun failed
// runs that are on the expected attempt and haven't been superseded by a newer
// run.
func (r *RerunIfDisconnectedReq) Run(ctx context.Context, client *gh.Client) (*RerunIfDisconnectedRes, error) {
	if err := r.validate(client); err != nil {
		return nil, err
	}

	logger := slog.Default().With(
		slog.String("owner", r.Owner),
		slog.String("repo", r.Repo),
		slog.Int64("run_id", r.RunID),
	)

	run, err := r.waitForRun(ctx, client, logger)
	if err != nil {
		return nil, err
	}

	res := &RerunIfDisconnectedRes{
		ListDisconnectedJobsRes: &ListDisconnectedJobsRes{
			Owner:      r.Owner,
			Repo:       r.Repo,
			RunID:      r.RunID,
			RunAttempt: run.GetRunAttempt(),
			FailedJobs: []*DisconnectedJobResult{},
		},
	}
	res.SkipReason = r.skipReason(run)
	if res.SkipReason != "" {
		logger.DebugContext(ctx, "not inspecting the run", slog.String("reason", res.SkipReason))
		return res, nil
	}

	// Inspect the attempt that we checked. If someone reruns it in the meantime we
	// don't want to look at the jobs of the new attempt.
	listReq := &ListDisconnectedJobsReq{Owner: r.Owner, Repo: r.Repo, RunID: r.RunID, Attempt: run.GetRunAttempt()}
	res.ListDisconnectedJobsRes, err = listReq.Run(ctx, client)
	if err != nil {
		return nil, err
	}
	if !res.SuspectedDisconnect {
		logger.DebugContext(ctx, "no suspected runner disconnects found")
		return res, nil
	}

	newer, err := r.findNewerWorkflowRun(ctx, client, run)
	if err != nil {
		return nil, fmt.Errorf("checking for runs newer than %s/%s run %d: %w", r.Owner, r.Repo, r.RunID, err)
	}
	if newer != nil {
		res.SkipReason = fmt.Sprintf("newer run %d has superseded it", newer.GetID())
		logger.DebugContext(ctx, "not rerunning a superseded run", slog.Int64("newer_run_id", newer.GetID()))
		return res, nil
	}

	logger.DebugContext(ctx, "requesting rerun of all failed workflow jobs and their dependents")
	if _, err := client.Actions.RerunFailedJobsByID(ctx, r.Owner, r.Repo, r.RunID); err != nil {
		return nil, fmt.Errorf("rerunning failed jobs for %s/%s run %d: %w", r.Owner, r.Repo, r.RunID, err)
	}

	res.RerunRequested = true
	logger.DebugContext(ctx, "Github accepted our request to rerun failed jobs")
	return res, nil
}

// skipReason returns why we won't inspect the run, or nothing if we will.
func (r *RerunIfDisconnectedReq) skipReason(run *gh.WorkflowRun) string {
	switch {
	case r.Attempt > 0 && run.GetRunAttempt() != r.Attempt:
		return fmt.Sprintf("we expected attempt %d but the run is on attempt %d", r.Attempt, run.GetRunAttempt())
	case run.GetConclusion() != "failure":
		return fmt.Sprintf("the run didn't fail, it concluded with %s", run.GetConclusion())
	default:
		return ""
	}
}

// waitForRun gets the run. If we're waiting we'll keep polling until the run
// has completed or moved on to an attempt that we don't expect.
func (r *RerunIfDisconnectedReq) waitForRun(
	ctx context.Context,
	client *gh.Client,
	logger *slog.Logger,
) (*gh.WorkflowRun, error) {
	for {
		logger.DebugContext(ctx, "getting workflow run")
		var run *gh.WorkflowRun
		err := retryWithBackoff(ctx, 3, func() error {
			var err error
			run, _, err = client.Actions.GetWorkflowRunByID(ctx, r.Owner, r.Repo, r.RunID)
			return err
		})
		if err != nil {
			if ctx.Err() != nil {
				err = context.Cause(ctx)
			}
			return nil, fmt.Errorf("getting %s/%s run %d: %w", r.Owner, r.Repo, r.RunID, err)
		}
		logger.DebugContext(ctx, "got workflow run",
			slog.String("status", run.GetStatus()),
			slog.String("conclusion", run.GetConclusion()),
			slog.Int("run_attempt", run.GetRunAttempt()),
		)

		if run.GetStatus() == "completed" || (r.Attempt > 0 && run.GetRunAttempt() != r.Attempt) {
			return run, nil
		}
		if !r.Wait {
			return nil, fmt.Errorf(
				"%s/%s run %d is %s and Github will only rerun completed runs",
				r.Owner, r.Repo, r.RunID, run.GetStatus(),
			)
		}

		logger.DebugContext(ctx, "waiting for workflow run to complete", slog.Duration("poll_interval", r.PollInterval))
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf(
				"waiting for %s/%s run %d to complete: %w", r.Owner, r.Repo, r.RunID, context.Cause(ctx),
			)
		case <-time.After(r.PollInterval):
		}
	}
}

// findNewerWorkflowRun returns the latest run of the same workflow for the same
// branch and event if it's newer than our run. There's no point in rerunning a
// superseded run, and if the workflow cancels in-progress runs for a branch it
// might even cancel the newer run.
func (r *RerunIfDisconnectedReq) findNewerWorkflowRun(
	ctx context.Context,
	client *gh.Client,
	run *gh.WorkflowRun,
) (*gh.WorkflowRun, error) {
	opts := &gh.ListWorkflowRunsOptions{
		Branch:      run.GetHeadBranch(),
		Event:       run.GetEvent(),
		ListOptions: gh.ListOptions{PerPage: 1},
	}
	slog.Default().DebugContext(ctx, "getting latest workflow run",
		slog.String("owner", r.Owner),
		slog.String("repo", r.Repo),
		slog.Int64("workflow_id", run.GetWorkflowID()),
		slog.String("branch", opts.Branch),
		slog.String("event", opts.Event),
	)

	var runs *gh.WorkflowRuns
	err := retryWithBackoff(ctx, 3, func() error {
		var err error
		runs, _, err = client.Actions.ListWorkflowRunsByID(ctx, r.Owner, r.Repo, run.GetWorkflowID(), opts)
		return err
	})
	if err != nil {
		return nil, err
	}

	return r.newerWorkflowRun(run, runs.WorkflowRuns), nil
}

// newerWorkflowRun returns the first of the runs that isn't our run and was
// created after it, if there is one.
func (r *RerunIfDisconnectedReq) newerWorkflowRun(run *gh.WorkflowRun, runs []*gh.WorkflowRun) *gh.WorkflowRun {
	for _, other := range runs {
		if other.GetID() != run.GetID() && other.GetCreatedAt().After(run.GetCreatedAt().Time) {
			return other
		}
	}

	return nil
}

// validate ensures that we've been given everything we need to rerun the jobs.
func (r *RerunIfDisconnectedReq) validate(client *gh.Client) error {
	if r == nil {
		return errors.New("uninitialized request")
	}

	listReq := &ListDisconnectedJobsReq{Owner: r.Owner, Repo: r.Repo, RunID: r.RunID}
	if err := listReq.validate(client); err != nil {
		return err
	}

	switch {
	case r.Attempt < 0:
		return errors.New("attempt can't be negative")
	case r.Wait && r.PollInterval <= 0:
		return errors.New("poll interval must be positive when waiting")
	}

	return nil
}

// ToJSON marshals the response to JSON.
func (r *RerunIfDisconnectedRes) ToJSON() ([]byte, error) {
	if r == nil {
		return nil, errors.New("uninitialized response")
	}

	return json.Marshal(r)
}

// ToTable marshals the response to plain text tables.
func (r *RerunIfDisconnectedRes) ToTable() string {
	return r.render(false)
}

// ToMarkdown marshals the response to Markdown tables.
func (r *RerunIfDisconnectedRes) ToMarkdown() string {
	return r.render(true)
}

func (r *RerunIfDisconnectedRes) render(markdown bool) string {
	if r == nil {
		return ""
	}

	output := r.ListDisconnectedJobsRes.render(markdown)
	switch {
	case r.RerunRequested:
		output += "\n\nRerun of failed jobs requested"
	case r.SkipReason != "":
		output += "\n\nNo rerun requested: " + r.SkipReason
	default:
		output += "\n\nNo suspected disconnects found. No rerun requested."
	}

	return output
}
