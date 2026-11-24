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

// CreateReleaseReq is a request to create a GitHub release for an existing tag.
type CreateReleaseReq struct {
	Owner               string
	Repo                string
	Tag                 string
	Prerelease          bool
	WriteToGithubOutput bool
}

// CreateReleaseRes is a CreateReleaseReq response.
type CreateReleaseRes struct {
	Tag     string `json:"tag"`
	URL     string `json:"url"`
	Created bool   `json:"created"`
}

// Run creates a release with GitHub generated notes for the tag, unless the tag
// already has a release.
func (r *CreateReleaseReq) Run(ctx context.Context, client *gh.Client) (*CreateReleaseRes, error) {
	if err := r.validate(client); err != nil {
		return nil, err
	}
	res := &CreateReleaseRes{Tag: r.Tag}

	slog.Default().DebugContext(ctx, "checking for release",
		slog.String("owner", r.Owner),
		slog.String("repo", r.Repo),
		slog.String("tag", r.Tag),
	)
	release, _, err := client.Repositories.GetReleaseByTag(ctx, r.Owner, r.Repo, r.Tag)
	switch {
	case err == nil:
		res.URL = release.GetHTMLURL()
		return res, nil
	case !isNotFound(err):
		return nil, fmt.Errorf("get release for %s: %w", r.Tag, err)
	}

	notes, err := generateReleaseNotes(ctx, client, r.Owner, r.Repo, r.Tag)
	if err != nil {
		return nil, err
	}

	// "legacy" lets GitHub keep the highest semver as latest when we release an
	// older release line.
	makeLatest := "legacy"
	if r.Prerelease {
		makeLatest = "false"
	}
	slog.Default().DebugContext(ctx, "creating release",
		slog.String("owner", r.Owner),
		slog.String("repo", r.Repo),
		slog.String("tag", r.Tag),
		slog.Bool("prerelease", r.Prerelease),
	)
	release, _, err = client.Repositories.CreateRelease(ctx, r.Owner, r.Repo, &gh.RepositoryRelease{
		TagName:    gh.Ptr(r.Tag),
		Name:       gh.Ptr(r.Tag),
		Body:       gh.Ptr(notes),
		Prerelease: gh.Ptr(r.Prerelease),
		MakeLatest: gh.Ptr(makeLatest),
	})
	if err != nil {
		return nil, fmt.Errorf("create release for %s: %w", r.Tag, err)
	}
	res.URL = release.GetHTMLURL()
	res.Created = true

	return res, nil
}

func (r *CreateReleaseReq) validate(client *gh.Client) error {
	switch {
	case r == nil:
		return errors.New("uninitialized request")
	case client == nil:
		return errors.New("GitHub client is required")
	case strings.TrimSpace(r.Owner) == "":
		return errors.New("owner is required")
	case strings.TrimSpace(r.Repo) == "":
		return errors.New("repo is required")
	case strings.TrimSpace(r.Tag) == "":
		return errors.New("tag is required")
	}

	return nil
}

// ToJSON marshals the response to JSON.
func (r *CreateReleaseRes) ToJSON() ([]byte, error) {
	if r == nil {
		return nil, errors.New("uninitialized response")
	}
	b, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("marshaling create release response to JSON: %w", err)
	}

	return b, nil
}

// ToTable marshals the response to a text table.
func (r *CreateReleaseRes) ToTable() string {
	if r == nil {
		return ""
	}
	t := table.NewWriter()
	t.Style().Options = table.Options{}
	t.AppendHeader(table.Row{"Tag", "URL", "Status"})
	t.AppendRow(table.Row{r.Tag, r.URL, createStatus(r.Created)})
	t.SuppressTrailingSpaces()

	return t.Render()
}

// generateReleaseNotes returns GitHub generated release notes for an existing tag.
func generateReleaseNotes(ctx context.Context, client *gh.Client, owner, repo, tag string) (string, error) {
	// GitHub would create a missing tag from the default branch when we create
	// the release, so we only generate notes for tags that already exist.
	exists, err := tagExists(ctx, client, owner, repo, tag)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", fmt.Errorf("tag %s does not exist", tag)
	}

	slog.Default().DebugContext(ctx, "generating release notes",
		slog.String("owner", owner),
		slog.String("repo", repo),
		slog.String("tag", tag),
	)
	notes, _, err := client.Repositories.GenerateReleaseNotes(ctx, owner, repo, &gh.GenerateNotesOptions{
		TagName: tag,
	})
	if err != nil {
		return "", fmt.Errorf("generate release notes for %s: %w", tag, err)
	}

	return notes.Body, nil
}
