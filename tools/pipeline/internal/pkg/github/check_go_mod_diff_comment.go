// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"context"
	"fmt"

	libgithub "github.com/google/go-github/v83/github"
)

// goModDiffCommentMarker is a hidden marker that identifies comments created
// for go.mod diffs so that they can be found later.
const goModDiffCommentMarker = "<!-- pipeline:go-mod-diff -->"

// goModNoDriftComment is the body of a previously posted comment once the
// drift has been resolved.
const goModNoDriftComment = goModDiffCommentMarker + "\n**Go Mod Diff:** No go.mod drift."

// CommentMarkdown renders the response as a pull request comment. It returns
// an empty string when there is no drift to report.
func (r *CheckGoModDiffRes) CommentMarkdown(err error) (string, error) {
	if err == nil && r.Err() == nil {
		return "", nil
	}

	tbl, tblErr := r.ToMarkdownTable(err)
	if tblErr != nil {
		return "", tblErr
	}
	tbl.SetTitle("Go Mod Diff")

	return goModDiffCommentMarker + "\n" + tbl.RenderMarkdown(), nil
}

// CommentOnPullRequest creates or updates the go.mod drift comment on the pull
// request. With no drift an existing comment is updated to say so and nothing
// is posted otherwise.
func (r *CheckGoModDiffRes) CommentOnPullRequest(
	ctx context.Context,
	github *libgithub.Client,
	owner string,
	repo string,
	pullNumber int,
	err error,
) error {
	body, bodyErr := r.CommentMarkdown(err)
	if bodyErr != nil {
		return fmt.Errorf("rendering go.mod diff comment: %w", bodyErr)
	}

	existing, findErr := findPullRequestCommentByMarker(ctx, github, owner, repo, pullNumber, goModDiffCommentMarker)
	if findErr != nil {
		return findErr
	}

	switch {
	case existing != nil && body == "":
		return updatePullRequestComment(ctx, github, owner, repo, existing.GetID(), goModNoDriftComment)
	case existing != nil:
		return updatePullRequestComment(ctx, github, owner, repo, existing.GetID(), body)
	case body == "":
		return nil
	}

	_, createErr := createPullRequestComment(ctx, github, owner, repo, pullNumber, body)

	return createErr
}
