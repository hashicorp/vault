// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	gh "github.com/google/go-github/v81/github"
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/changed"
	"github.com/jedib0t/go-pretty/v6/table"
)

type (
	// ListChangedFilesReq is a request to list workflows runs. The fields represent
	// various criteria we can use to filter.
	ListChangedFilesReq struct {
		Owner               string
		Repo                string
		PullNumber          int
		CommitSHA           string
		GroupFiles          bool
		WriteToGithubOutput bool
	}

	// ListChangedFilesRes is a list workflows response.
	ListChangedFilesRes struct {
		Files  changed.Files      `json:"files,omitempty"`
		Groups changed.FileGroups `json:"groups,omitempty"`
	}

	// ListChangedFilesGithubOutput is out GITHUB_OUTPUT type. It's a slimmed down
	// type that only include file names and groups.
	ListChangedFilesGithubOutput struct {
		Files  []string           `json:"files,omitempty"`
		Groups changed.FileGroups `json:"groups,omitempty"`
	}
)

// Run runs the request to gather all instances of the workflow that match
// our filter criteria.
func (r *ListChangedFilesReq) Run(ctx context.Context, client *gh.Client) (*ListChangedFilesRes, error) {
	var err error
	res := &ListChangedFilesRes{}

	if err = r.validate(); err != nil {
		return nil, fmt.Errorf("validating request: %w", err)
	}

	if r.CommitSHA != "" {
		files, err := r.getCommitFiles(ctx, client)
		if err != nil {
			return nil, err
		}
		res.Files = append(res.Files, files...)
	}

	if r.PullNumber != 0 {
		files, err := r.getPullFiles(ctx, client)
		if err != nil {
			return nil, err
		}
		res.Files = append(res.Files, files...)
	}

	if r.GroupFiles {
		changed.GroupFiles(ctx, res.Files, changed.DefaultFileGroupCheckers...)
		res.Groups = changed.FileGroups{}
		for _, file := range res.Files {
			for _, group := range file.Groups {
				res.Groups = res.Groups.Add(group)
			}
		}
	}

	return res, nil
}

// validate ensures that we've been given the minimum filter arguments necessary to complete a
// request. It is always recommended that additional fitlers be given to reduce the response size
// and not exhaust API limits.
func (r *ListChangedFilesReq) validate() error {
	if r == nil {
		return errors.New("failed to initialize request")
	}

	if r.Owner == "" {
		return errors.New("no github organization has been provided")
	}

	if r.Repo == "" {
		return errors.New("no github repository has been provided")
	}

	if r.PullNumber == 0 && r.CommitSHA == "" {
		return errors.New("no pull request number or commit SHA has been provided")
	}

	return nil
}

// getCommitFiles attempts to locate the workflow associated with our workflow name.
func (r *ListChangedFilesReq) getCommitFiles(ctx context.Context, client *gh.Client) (changed.Files, error) {
	opts := &gh.ListOptions{PerPage: PerPageMax}
	files := changed.Files{}
	for {
		commit, res, err := client.Repositories.GetCommit(ctx, r.Owner, r.Repo, r.CommitSHA, opts)
		if err != nil {
			return nil, err
		}

		for _, f := range commit.Files {
			files = append(files, &changed.File{File: f})
		}

		if res.NextPage == 0 {
			return files, nil
		}

		opts.Page = res.NextPage
	}
}

// getPullFiles attempts to locate the workflow associated with our workflow name.
func (r *ListChangedFilesReq) getPullFiles(ctx context.Context, client *gh.Client) (changed.Files, error) {
	opts := &gh.ListOptions{PerPage: PerPageMax}
	files := changed.Files{}
	for {
		fl, res, err := client.PullRequests.ListFiles(ctx, r.Owner, r.Repo, r.PullNumber, opts)
		if err != nil {
			return nil, err
		}

		for _, f := range fl {
			files = append(files, &changed.File{File: f})
		}

		if res.NextPage == 0 {
			return files, nil
		}

		opts.Page = res.NextPage
	}
}

// ToJSON marshals the response to JSON.
func (r *ListChangedFilesRes) ToJSON() ([]byte, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("marshaling list changed files to JSON: %w", err)
	}

	return b, nil
}

// ToGithubOutput writes a simplified list of changed files to be used $GITHUB_OUTPUT
func (r ListChangedFilesRes) ToGithubOutput() ([]byte, error) {
	res := &ListChangedFilesGithubOutput{
		Groups: r.Groups,
	}
	if f := r.Files; f != nil {
		res.Files = f.Names()
	}

	b, err := json.Marshal(res)
	if err != nil {
		return nil, fmt.Errorf("marshaling list changed files GITHUB_OUTPUT to JSON: %w", err)
	}

	return b, nil
}

// ToTable marshals the response to a text table.
func (r *ListChangedFilesRes) ToTable(groups bool) string {
	if !groups {
		w := strings.Builder{}
		for _, name := range r.Files.Names() {
			w.WriteString(name + "\n")
		}

		return w.String()
	}

	t := table.NewWriter()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false
	t.Style().Options.SeparateFooter = false
	t.Style().Options.SeparateHeader = false
	t.Style().Options.SeparateRows = false
	t.AppendHeader(table.Row{"path", "groups"})
	for _, file := range r.Files {
		t.AppendRow(table.Row{file.Name(), file.Groups.String()})
	}
	return t.Render()
}
