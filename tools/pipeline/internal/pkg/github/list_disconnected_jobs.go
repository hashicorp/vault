// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	gh "github.com/google/go-github/v83/github"
	"github.com/jedib0t/go-pretty/v6/table"
)

// ListDisconnectedJobsReq is a request to list the failed jobs of a workflow run
// and whether or not we think their runner disconnected.
type ListDisconnectedJobsReq struct {
	Owner string
	Repo  string
	RunID int64
	// Attempt is the attempt of the run that we'll inspect. If it's not set we'll
	// inspect the latest attempt.
	Attempt             int
	WriteToGithubOutput bool
}

// ListDisconnectedJobsRes is a response to a ListDisconnectedJobsReq.
type ListDisconnectedJobsRes struct {
	Owner               string                   `json:"owner"`
	Repo                string                   `json:"repo"`
	RunID               int64                    `json:"run_id"`
	RunAttempt          int                      `json:"run_attempt"`
	JobsInspected       int                      `json:"jobs_inspected"`
	SuspectedDisconnect bool                     `json:"suspected_disconnect"`
	FailedJobs          []*DisconnectedJobResult `json:"failed_jobs"`
}

// DisconnectedJobResult includes each failed job, even if we don't suspect a disconnect.
type DisconnectedJobResult struct {
	ID                  int64                     `json:"id"`
	Name                string                    `json:"name"`
	URL                 string                    `json:"url"`
	SuspectedDisconnect bool                      `json:"suspected_disconnect"`
	UnfinishedSteps     []*DisconnectedStepResult `json:"unfinished_steps"`
}

// DisconnectedStepResult contains an unfinished step with its original number and status.
type DisconnectedStepResult struct {
	Number int64  `json:"number"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

// Run lists the failed jobs of an attempt of the run. Github doesn't tell us when
// a runner disconnects but we can usually spot it. If a job has completed and
// failed but some of its steps are still queued, pending, or in progress, the
// runner most likely went away before it could finish them.
//
// The run doesn't have to be completed so it's safe to list the jobs of a run
// that is still in progress.
func (r *ListDisconnectedJobsReq) Run(ctx context.Context, client *gh.Client) (*ListDisconnectedJobsRes, error) {
	if err := r.validate(client); err != nil {
		return nil, err
	}

	var jobs []*gh.WorkflowJob
	var err error
	if r.Attempt > 0 {
		jobs, err = getWorkflowJobsForRunAttempt(ctx, client, r.Owner, r.Repo, r.RunID, int64(r.Attempt))
	} else {
		jobs, err = getWorkflowJobsForRun(ctx, client, r.Owner, r.Repo, r.RunID)
	}
	if err != nil {
		return nil, fmt.Errorf("inspecting %s/%s run %d: %w", r.Owner, r.Repo, r.RunID, err)
	}

	return r.inspectJobs(ctx, jobs), nil
}

// inspectJobs looks for failed jobs that still have unfinished steps.
func (r *ListDisconnectedJobsReq) inspectJobs(ctx context.Context, jobs []*gh.WorkflowJob) *ListDisconnectedJobsRes {
	logger := slog.Default().With(
		slog.String("owner", r.Owner),
		slog.String("repo", r.Repo),
		slog.Int64("run_id", r.RunID),
		slog.Int("run_attempt", r.Attempt),
	)
	logger.DebugContext(ctx, "checking workflow jobs for runner disconnects")

	res := &ListDisconnectedJobsRes{
		Owner:         r.Owner,
		Repo:          r.Repo,
		RunID:         r.RunID,
		RunAttempt:    r.Attempt,
		JobsInspected: len(jobs),
		FailedJobs:    []*DisconnectedJobResult{},
	}
	suspectedCount := 0

	for _, job := range jobs {
		if res.RunAttempt == 0 {
			// Every job of the latest attempt has the same attempt so we'll take the first
			res.RunAttempt = int(job.GetRunAttempt())
		}
		if job.GetStatus() != "completed" || job.GetConclusion() != "failure" {
			continue
		}

		result := &DisconnectedJobResult{
			ID:              job.GetID(),
			Name:            job.GetName(),
			URL:             job.GetHTMLURL(),
			UnfinishedSteps: []*DisconnectedStepResult{},
		}
		for _, step := range job.Steps {
			switch step.GetStatus() {
			case "queued", "pending", "in_progress":
				logger.DebugContext(ctx, "found unfinished step in a completed, failed job",
					slog.Int64("job_id", job.GetID()),
					slog.String("job_name", job.GetName()),
					slog.Int64("step_number", step.GetNumber()),
					slog.String("step_name", step.GetName()),
					slog.String("step_status", step.GetStatus()),
				)
				result.UnfinishedSteps = append(result.UnfinishedSteps, &DisconnectedStepResult{
					Number: step.GetNumber(), Name: step.GetName(), Status: step.GetStatus(),
				})
			}
		}
		result.SuspectedDisconnect = len(result.UnfinishedSteps) > 0
		if result.SuspectedDisconnect {
			suspectedCount++
		}
		logger.DebugContext(ctx, "checked failed job for a runner disconnect",
			slog.Int64("job_id", result.ID),
			slog.String("job_name", result.Name),
			slog.Bool("suspected_disconnect", result.SuspectedDisconnect),
			slog.Int("unfinished_steps", len(result.UnfinishedSteps)),
		)
		res.SuspectedDisconnect = res.SuspectedDisconnect || result.SuspectedDisconnect
		res.FailedJobs = append(res.FailedJobs, result)
	}

	logger.DebugContext(ctx, "finished checking workflow jobs",
		slog.Int("jobs_inspected", res.JobsInspected),
		slog.Int("failed_jobs", len(res.FailedJobs)),
		slog.Int("suspected_disconnects", suspectedCount),
	)

	return res
}

// validate ensures that we've been given everything we need to list the jobs.
func (r *ListDisconnectedJobsReq) validate(client *gh.Client) error {
	switch {
	case r == nil:
		return errors.New("uninitialized request")
	case client == nil:
		return errors.New("GitHub client is required")
	case strings.TrimSpace(r.Owner) == "":
		return errors.New("owner is required")
	case strings.TrimSpace(r.Repo) == "":
		return errors.New("repo is required")
	case r.RunID <= 0:
		return errors.New("run ID must be positive")
	case r.Attempt < 0:
		return errors.New("attempt can't be negative")
	}

	return nil
}

// ToGithubOutput marshals the response to JSON for $GITHUB_OUTPUT.
func (r *ListDisconnectedJobsRes) ToGithubOutput() ([]byte, error) {
	return r.ToJSON()
}

// ToJSON marshals the response to JSON.
func (r *ListDisconnectedJobsRes) ToJSON() ([]byte, error) {
	if r == nil {
		return nil, errors.New("uninitialized response")
	}

	return json.Marshal(r)
}

// ToTable marshals the response to plain text tables.
func (r *ListDisconnectedJobsRes) ToTable() string {
	return r.render(false)
}

// ToMarkdown marshals the response to Markdown tables.
func (r *ListDisconnectedJobsRes) ToMarkdown() string {
	return r.render(true)
}

func (r *ListDisconnectedJobsRes) render(markdown bool) string {
	if r == nil {
		return ""
	}
	cell := func(s string) string {
		if !markdown {
			return s
		}
		return escapeMarkdownTableCell(s)
	}
	newTable := func(header table.Row) table.Writer {
		t := table.NewWriter()
		t.Style().Options = table.Options{}
		t.AppendHeader(header)
		t.SuppressTrailingSpaces()
		return t
	}
	summary := newTable(table.Row{"Repository", "Run ID", "Attempt", "Jobs inspected", "Failed jobs", "Suspected disconnects"})
	failed := newTable(table.Row{"Job ID", "Name", "Suspected disconnect", "URL"})
	steps := newTable(table.Row{"Job ID", "Step", "Name", "Status"})
	suspected, unfinished := 0, 0
	for _, job := range r.FailedJobs {
		if job.SuspectedDisconnect {
			suspected++
		}
		failed.AppendRow(table.Row{job.ID, cell(job.Name), job.SuspectedDisconnect, cell(job.URL)})
		for _, step := range job.UnfinishedSteps {
			steps.AppendRow(table.Row{job.ID, step.Number, cell(step.Name), cell(step.Status)})
			unfinished++
		}
	}
	summary.AppendRow(table.Row{cell(r.Owner + "/" + r.Repo), r.RunID, r.RunAttempt, r.JobsInspected, len(r.FailedJobs), suspected})
	var output strings.Builder
	section := func(name string, t table.Writer) {
		if output.Len() > 0 {
			output.WriteString("\n\n")
		}
		if markdown {
			output.WriteString("## " + name + "\n\n" + t.RenderMarkdown())
		} else {
			output.WriteString(name + ":\n" + t.Render())
		}
	}
	section("Summary", summary)
	if len(r.FailedJobs) > 0 {
		section("Failed Jobs", failed)
	}
	if unfinished > 0 {
		section("Unfinished Steps", steps)
	}

	return output.String()
}
