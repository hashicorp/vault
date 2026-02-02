// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"context"
	"errors"
	"fmt"

	libgithub "github.com/google/go-github/v81/github"
	"github.com/jedib0t/go-pretty/v6/table"
)

// ListCommitStatusesReq is a request to list workflows runs. The fields represent
// various criteria we can use to filter.
type ListCommitStatusesReq struct {
	Owner  string
	Repo   string
	Commit string
	PR     int
}

// ListCommitStatusesRes is a list workflows response.
type ListCommitStatusesRes struct {
	Statuses []*libgithub.RepoStatus `json:"statuses,omitempty"`
}

// Run runs the request to gather all instances of the workflow that match
// our filter criteria.
func (r *ListCommitStatusesReq) Run(ctx context.Context, client *libgithub.Client) (*ListCommitStatusesRes, error) {
	var err error
	res := &ListCommitStatusesRes{}

	if err = r.validate(); err != nil {
		return nil, fmt.Errorf("validating request: %w", err)
	}

	if r.Commit == "" {
		var pr *libgithub.PullRequest
		pr, err = getPullRequest(ctx, client, r.Owner, r.Repo, r.PR)
		if err != nil {
			return nil, fmt.Errorf("getting pull request: %w", err)
		}
		r.Commit = pr.GetHead().GetSHA()
	}

	res.Statuses, err = listCommitStatuses(ctx, client, r.Owner, r.Repo, r.Commit)
	if err != nil {
		return nil, fmt.Errorf("getting commit statuses: %w", err)
	}

	return res, nil
}

// validate ensures that we've been given the minimum filter arguments necessary to complete a
// request. It is always recommended that additional fitlers be given to reduce the response size
// and not exhaust API limits.
func (r *ListCommitStatusesReq) validate() error {
	if r == nil {
		return errors.New("failed to initialize request")
	}

	if r.Owner == "" {
		return errors.New("no github organization has been provided")
	}

	if r.Repo == "" {
		return errors.New("no github repository has been provided")
	}

	if r.PR == 0 && r.Commit == "" {
		return errors.New("no commit or Pull Request number has been provided")
	}

	return nil
}

// ToTable marshals the response to a text table.
func (r *ListCommitStatusesRes) ToTable() table.Writer {
	t := table.NewWriter()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false
	t.Style().Options.SeparateFooter = false
	t.Style().Options.SeparateHeader = false
	t.Style().Options.SeparateRows = false
	t.AppendHeader(table.Row{"context", "creator", "date", "type", "state"})
	for _, status := range r.Statuses {
		t.AppendRow(table.Row{
			status.GetContext(),
			status.GetCreator().GetLogin(),
			status.GetUpdatedAt(),
			status.GetCreator().GetType(),
			status.GetState(),
		})
	}

	return t
}
