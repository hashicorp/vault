// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	libgithub "github.com/google/go-github/v81/github"
	"github.com/jedib0t/go-pretty/v6/table"
)

// CheckCommitStatusReq is a request to check commit statuses for an expected
// state, context, and/or creator.
type CheckCommitStatusReq struct {
	Owner   string
	Repo    string
	Commit  string
	PR      int
	Context string
	Creator string
	State   string
}

// CheckCommitStatusRes is a list workflows response.
type CheckCommitStatusRes struct {
	CheckSuccessful bool                 `json:"check_success,omitempty"`
	ExpectedContext string               `json:"expected_context,omitempty"`
	ExpectedState   string               `json:"expected_state,omitempty"`
	ExpectedCreator string               `json:"expected_creator,omitempty"`
	Statuses        []*CheckCommitStatus `json:"statuses,omitempty"`
}

// CheckCommitStatus is an instance of one commit status.
type CheckCommitStatus struct {
	Status       *libgithub.RepoStatus `json:"statuses,omitempty"`
	CheckSuccess bool                  `json:"success,omitempty"`
}

// String returns the response as a string
func (r *CheckCommitStatusRes) String() string {
	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("success:%t", r.CheckSuccessful))
	b.WriteString(" context:" + r.ExpectedContext)
	b.WriteString(" state:" + r.ExpectedState)
	b.WriteString(" creator:" + r.ExpectedCreator)

	return b.String()
}

// Run runs the request to check the commit statuses of a Pull Request.
func (r *CheckCommitStatusReq) Run(ctx context.Context, client *libgithub.Client) (*CheckCommitStatusRes, error) {
	var err error
	res := &CheckCommitStatusRes{
		CheckSuccessful: false,
		ExpectedCreator: r.Creator,
		ExpectedState:   r.State,
		ExpectedContext: r.Context,
		Statuses:        []*CheckCommitStatus{},
	}

	if err = r.validate(); err != nil {
		return nil, fmt.Errorf("validating request: %w", err)
	}

	statusesReq := &ListCommitStatusesReq{
		Owner:  r.Owner,
		Repo:   r.Repo,
		Commit: r.Commit,
		PR:     r.PR,
	}
	statuses, err := statusesReq.Run(ctx, client)
	if err != nil {
		return nil, err
	}

	for _, status := range statuses.Statuses {
		if status.GetContext() != r.Context {
			continue
		}

		if r.Creator != "" && status.GetCreator().GetLogin() != r.Creator {
			continue
		}

		// There can be many statuses with the same context and creator. Keep track
		// of them all but only update our success if we get a match.
		res.Statuses = append(res.Statuses, &CheckCommitStatus{
			Status:       status,
			CheckSuccess: status.GetState() == r.State,
		})

		if status.GetState() == r.State {
			res.CheckSuccessful = true
		}
	}

	return res, nil
}

// validate ensures that we've been given the arguments necessary to complete a
// request.
func (r *CheckCommitStatusReq) validate() error {
	if r == nil {
		return errors.New("failed to initialize request")
	}

	if r.Owner == "" {
		return errors.New("no github organization has been provided")
	}

	if r.Repo == "" {
		return errors.New("no github repository has been provided")
	}

	if r.Context == "" {
		return errors.New("no status context has been provided")
	}

	if r.PR == 0 && r.Commit == "" {
		return errors.New("no commit or Pull Request number has been provided")
	}

	allowedStates := []string{"error", "failure", "pending", "success"}
	if !slices.Contains(allowedStates, r.State) {
		return fmt.Errorf("invalid state, got: %s, expected one of: %v+", r.State, allowedStates)
	}

	return nil
}

// ToTable marshals the response to a text table.
func (r *CheckCommitStatusRes) ToTable() table.Writer {
	t := table.NewWriter()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false
	t.Style().Options.SeparateFooter = false
	t.Style().Options.SeparateHeader = false
	t.Style().Options.SeparateRows = false
	t.AppendHeader(table.Row{"context", "creator", "date", "state", "check success"})
	for _, status := range r.Statuses {
		t.AppendRow(table.Row{
			status.Status.GetContext(),
			status.Status.GetCreator().GetLogin(),
			status.Status.GetUpdatedAt(),
			status.Status.GetState(),
			status.CheckSuccess,
		})
	}
	t.SuppressEmptyColumns()
	t.SuppressTrailingSpaces()

	return t
}
