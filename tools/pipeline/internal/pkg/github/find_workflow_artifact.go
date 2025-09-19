// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"slices"

	gh "github.com/google/go-github/v74/github"
	"github.com/jedib0t/go-pretty/v6/table"
)

// FindWorkflowArtifactReq is a request to find an artifact associated with a
// workflow run.
type FindWorkflowArtifactReq struct {
	ArtifactName        string
	ArtifactPattern     string
	Owner               string
	PullNumber          int
	Repo                string
	WorkflowName        string
	WriteToGithubOutput bool
	compiledPattern     *regexp.Regexp
}

// FindWorkflowArtifactRes is a FindWorkflowArtifactReq response.
type FindWorkflowArtifactRes struct {
	PR       *gh.PullRequest `json:"pr,omitempty"`
	Workflow *gh.Workflow    `json:"workflow,omitempty"`
	Run      *WorkflowRun    `json:"runs,omitempty"`
	Artifact *gh.Artifact    `json:"artifact,omitempty"`
}

// Run performs the search to find an artifact associated with a workflow.
func (r *FindWorkflowArtifactReq) Run(ctx context.Context, client *gh.Client) (*FindWorkflowArtifactRes, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	var err error
	res := &FindWorkflowArtifactRes{}

	// Validate our request. This also ensures that any pattern we've been given
	// is a valid regex.
	if err = r.validate(); err != nil {
		return nil, fmt.Errorf("validating request: %w", err)
	}

	// Get the workflow details for the repo
	res.Workflow, err = getWorkflow(ctx, client, r.Owner, r.Repo, r.WorkflowName)
	if err != nil {
		return nil, fmt.Errorf("getting workflow: %w", err)
	}

	// Get the pull request we're searching
	res.PR, err = getPullRequest(ctx, client, r.Owner, r.Repo, r.PullNumber)
	if err != nil {
		return nil, fmt.Errorf("getting pull request: %w", err)
	}

	// Get the workflow runs associated with the workflow and the PR
	opts := &gh.ListWorkflowRunsOptions{
		Branch:              res.PR.GetHead().GetRef(),
		ExcludePullRequests: false,
		HeadSHA:             res.PR.GetHead().GetSHA(),
		ListOptions:         gh.ListOptions{PerPage: PerPageMax},
		Status:              "success",
	}
	runs, err := getWorkflowRuns(ctx, client, r.Owner, r.Repo, res.Workflow.GetID(), opts)
	if err != nil {
		return nil, fmt.Errorf("getting workflow runs: %w", err)
	}

	if len(runs) < 1 {
		return nil, fmt.Errorf("no matching workflow runs are associated with the pull request")
	}

	// In instances where we have more than one run we want to get the artifact
	// from the most recent run if possible. Search our runs in reverse order to
	// find the most recent artifact.
	slices.SortFunc(runs, func(a, b *WorkflowRun) int {
		return cmp.Compare(*b.Run.RunAttempt, *a.Run.RunAttempt)
	})

	var artifacts gh.ArtifactList
	for _, run := range runs {
		artifacts, err = getWorkflowRunArtifacts(ctx, client, r.Owner, r.Repo, *run.Run.ID)
		if err != nil {
			return nil, fmt.Errorf("getting artifacts for workflow run %d: %w", *run.Run.ID, err)
		}

		for _, art := range artifacts.Artifacts {
			// If we've been given a name locate it by that
			if r.ArtifactName != "" {
				if art.GetName() == r.ArtifactName {
					res.Artifact = art

					return res, nil
				}
			} else {
				// Find it by regex
				if r.compiledPattern.MatchString(art.GetName()) {
					res.Artifact = art

					return res, nil
				}
			}
		}
	}

	return nil, errors.New("unable to find artifact matching given criteria")
}

// validate ensures that we've been given the request configuration to perform
// the request.
func (r *FindWorkflowArtifactReq) validate() error {
	if r == nil {
		return errors.New("failed to initialize request")
	}

	if r.Owner == "" {
		return errors.New("no github organization has been provided")
	}

	if r.Repo == "" {
		return errors.New("no github repository has been provided")
	}

	if r.PullNumber == 0 {
		return errors.New("no github pull request number has been provided")
	}

	if r.WorkflowName == "" {
		return errors.New("no workflow name has been provided")
	}

	if r.ArtifactName == "" && r.ArtifactPattern == "" {
		return errors.New("no artifact name or pattern has been provided")
	}

	if r.ArtifactName != "" && r.ArtifactPattern != "" {
		return errors.New("you must provide only an artifact name or pattern")
	}

	if r.ArtifactPattern != "" {
		var err error
		r.compiledPattern, err = regexp.Compile(r.ArtifactPattern)
		if err != nil {
			return fmt.Errorf("invalid artifact pattern: %w", err)
		}
	}

	return nil
}

// ToJSON marshals the response to JSON.
func (r *FindWorkflowArtifactRes) ToJSON() ([]byte, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("marshaling find workflow artifact to JSON: %w", err)
	}

	return b, nil
}

// ToGithubOutput marshals just the artifact response to JSON.
func (r *FindWorkflowArtifactRes) ToGithubOutput() ([]byte, error) {
	b, err := json.Marshal(r.Artifact)
	if err != nil {
		return nil, fmt.Errorf("marshaling find workflow artifact to GITHUB_OUTPUT JSON: %w", err)
	}

	return b, nil
}

// ToTable marshals the response to a text table.
func (r *FindWorkflowArtifactRes) ToTable() string {
	t := table.NewWriter()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false
	t.Style().Options.SeparateFooter = false
	t.Style().Options.SeparateHeader = false
	t.Style().Options.SeparateRows = false
	t.AppendHeader(table.Row{"name", "run id", "artifact id", "url"})
	t.AppendRow(table.Row{
		r.Artifact.GetName(),
		r.Artifact.GetWorkflowRun().GetID(),
		r.Artifact.GetID(),
		r.Artifact.GetArchiveDownloadURL(),
	})
	return t.Render()
}
