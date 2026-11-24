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

// CreateTagReq is a request to create a lightweight tag at the head of a branch.
type CreateTagReq struct {
	Owner               string
	Repo                string
	Tag                 string
	Branch              string
	WriteToGithubOutput bool
}

// CreateTagRes is a CreateTagReq response.
type CreateTagRes struct {
	Tag     string `json:"tag"`
	SHA     string `json:"sha"`
	Created bool   `json:"created"`
}

// Run creates the tag at the head of the branch. If the tag already exists we
// leave it alone as long as the branch contains its commit, since the branch
// may have moved on since we tagged it. If it doesn't we fail.
func (r *CreateTagReq) Run(ctx context.Context, client *gh.Client) (*CreateTagRes, error) {
	if err := r.validate(client); err != nil {
		return nil, err
	}

	tagSHA, err := getTagCommitSHA(ctx, client, r.Owner, r.Repo, r.Tag)
	if err != nil {
		return nil, err
	}
	if tagSHA != "" {
		if err := checkBranchContainsCommit(ctx, client, r.Owner, r.Repo, r.Branch, tagSHA); err != nil {
			return nil, fmt.Errorf("tag %s already exists: %w", r.Tag, err)
		}
		return &CreateTagRes{Tag: r.Tag, SHA: tagSHA}, nil
	}

	sha, err := getBranchHeadSHA(ctx, client, r.Owner, r.Repo, r.Branch)
	if err != nil {
		return nil, err
	}
	res := &CreateTagRes{Tag: r.Tag, SHA: sha}

	slog.Default().DebugContext(ctx, "creating tag",
		slog.String("owner", r.Owner),
		slog.String("repo", r.Repo),
		slog.String("tag", r.Tag),
		slog.String("sha", sha),
	)
	if _, _, err := client.Git.CreateRef(ctx, r.Owner, r.Repo, gh.CreateRef{
		Ref: "refs/tags/" + r.Tag,
		SHA: sha,
	}); err != nil {
		return nil, fmt.Errorf("create tag %s: %w", r.Tag, err)
	}
	res.Created = true

	return res, nil
}

func (r *CreateTagReq) validate(client *gh.Client) error {
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
	case strings.TrimSpace(r.Branch) == "":
		return errors.New("branch is required")
	}

	return nil
}

// ToJSON marshals the response to JSON.
func (r *CreateTagRes) ToJSON() ([]byte, error) {
	if r == nil {
		return nil, errors.New("uninitialized response")
	}
	b, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("marshaling create tag response to JSON: %w", err)
	}

	return b, nil
}

// ToTable marshals the response to a text table.
func (r *CreateTagRes) ToTable() string {
	if r == nil {
		return ""
	}
	t := table.NewWriter()
	t.Style().Options = table.Options{}
	t.AppendHeader(table.Row{"Tag", "SHA", "Status"})
	t.AppendRow(table.Row{r.Tag, r.SHA, createStatus(r.Created)})
	t.SuppressTrailingSpaces()

	return t.Render()
}

// createStatus describes whether a create request made something or found it
// already in place.
func createStatus(created bool) string {
	if created {
		return "created"
	}

	return "skipped"
}

// getBranchHeadSHA returns the commit SHA at the head of a branch.
func getBranchHeadSHA(ctx context.Context, client *gh.Client, owner, repo, branch string) (string, error) {
	slog.Default().DebugContext(ctx, "resolving branch head",
		slog.String("owner", owner),
		slog.String("repo", repo),
		slog.String("branch", branch),
	)
	res, _, err := client.Repositories.GetBranch(ctx, owner, repo, branch, 0)
	if err != nil {
		return "", fmt.Errorf("get branch %s: %w", branch, err)
	}
	sha := res.GetCommit().GetSHA()
	if sha == "" {
		return "", fmt.Errorf("branch %s did not resolve to a commit", branch)
	}

	return sha, nil
}

// checkBranchContainsCommit returns an error unless the commit is the head of
// the branch or one of its ancestors.
func checkBranchContainsCommit(ctx context.Context, client *gh.Client, owner, repo, branch, sha string) error {
	slog.Default().DebugContext(ctx, "comparing commit to branch",
		slog.String("owner", owner),
		slog.String("repo", repo),
		slog.String("branch", branch),
		slog.String("sha", sha),
	)
	comparison, _, err := client.Repositories.CompareCommits(ctx, owner, repo, sha, branch, nil)
	if err != nil {
		return fmt.Errorf("compare %s to branch %s: %w", sha, branch, err)
	}

	switch status := comparison.GetStatus(); status {
	case "ahead", "identical":
		return nil
	default:
		return fmt.Errorf("commit %s is not on branch %s, compare status is %q", sha, branch, status)
	}
}

// getTagCommitSHA returns the commit SHA that a lightweight or annotated tag
// points to, or an empty string if the tag does not exist.
func getTagCommitSHA(ctx context.Context, client *gh.Client, owner, repo, tag string) (string, error) {
	exists, err := tagExists(ctx, client, owner, repo, tag)
	if err != nil || !exists {
		return "", err
	}

	slog.Default().DebugContext(ctx, "resolving tag commit",
		slog.String("owner", owner),
		slog.String("repo", repo),
		slog.String("tag", tag),
	)
	commit, _, err := client.Repositories.GetCommit(ctx, owner, repo, "refs/tags/"+tag, nil)
	if err != nil {
		return "", fmt.Errorf("resolve tag %s commit: %w", tag, err)
	}

	return commit.GetSHA(), nil
}

// tagExists determines if a tag exists.
func tagExists(ctx context.Context, client *gh.Client, owner, repo, tag string) (bool, error) {
	slog.Default().DebugContext(ctx, "checking for tag",
		slog.String("owner", owner),
		slog.String("repo", repo),
		slog.String("tag", tag),
	)
	_, _, err := client.Git.GetRef(ctx, owner, repo, "tags/"+tag)
	if isNotFound(err) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("get tag %s: %w", tag, err)
	}

	return true, nil
}
