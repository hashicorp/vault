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

// RerunIfDisconnectedReq is a request to check a workflow run for runner disconnects.
type RerunIfDisconnectedReq struct {
	Owner string
	Repo  string
	RunID int64
}

// RerunIfDisconnectedRes contains the failed jobs we checked and whether Github
// accepted our rerun request.
type RerunIfDisconnectedRes struct {
	Owner          string                   `json:"owner"`
	Repo           string                   `json:"repo"`
	RunID          int64                    `json:"run_id"`
	JobsInspected  int                      `json:"jobs_inspected"`
	RerunRequested bool                     `json:"rerun_requested"`
	FailedJobs     []*DisconnectedJobResult `json:"failed_jobs"`
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

// Run checks every job before requesting one rerun if we suspect a disconnect.
func (r *RerunIfDisconnectedReq) Run(ctx context.Context, client *gh.Client) (*RerunIfDisconnectedRes, error) {
	if err := r.validate(client); err != nil {
		return nil, err
	}

	logger := slog.Default().With(
		slog.String("owner", r.Owner),
		slog.String("repo", r.Repo),
		slog.Int64("run_id", r.RunID),
	)
	logger.DebugContext(ctx, "checking workflow jobs for runner disconnects")

	jobs, err := getWorkflowJobsForRun(ctx, client, r.Owner, r.Repo, r.RunID)
	if err != nil {
		return nil, fmt.Errorf("inspecting %s/%s run %d: %w", r.Owner, r.Repo, r.RunID, err)
	}

	res := &RerunIfDisconnectedRes{
		Owner:         r.Owner,
		Repo:          r.Repo,
		RunID:         r.RunID,
		JobsInspected: len(jobs),
		FailedJobs:    []*DisconnectedJobResult{},
	}
	suspected := false
	suspectedCount := 0

	for _, job := range jobs {
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
		suspected = suspected || result.SuspectedDisconnect
		res.FailedJobs = append(res.FailedJobs, result)
	}
	logger.DebugContext(ctx, "finished checking workflow jobs",
		slog.Int("jobs_inspected", res.JobsInspected),
		slog.Int("failed_jobs", len(res.FailedJobs)),
		slog.Int("suspected_disconnects", suspectedCount),
	)
	if !suspected {
		logger.DebugContext(ctx, "no suspected runner disconnects found")
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

// ToTable formats the results as plain-text tables.
func (r *RerunIfDisconnectedRes) ToTable() string {
	return r.render(false)
}

// ToMarkdown formats the results as Markdown tables.
func (r *RerunIfDisconnectedRes) ToMarkdown() string {
	return r.render(true)
}

// ToJSON formats the results as JSON.
func (r *RerunIfDisconnectedRes) ToJSON() ([]byte, error) {
	if r == nil {
		return nil, errors.New("uninitialized response")
	}
	return json.Marshal(r)
}

func (r *RerunIfDisconnectedReq) validate(client *gh.Client) error {
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
	}

	return nil
}

func (r *RerunIfDisconnectedRes) render(markdown bool) string {
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
	summary := newTable(table.Row{"Repository", "Run ID", "Jobs inspected", "Failed jobs", "Suspected disconnects", "Rerun requested"})
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
	summary.AppendRow(table.Row{cell(r.Owner + "/" + r.Repo), r.RunID, r.JobsInspected, len(r.FailedJobs), suspected, r.RerunRequested})
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
	if r.RerunRequested {
		output.WriteString("\n\nRerun of failed jobs requested")
	} else {
		output.WriteString("\n\nNo suspected disconnects found. No rerun requested.")
	}
	return output.String()
}
