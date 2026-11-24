// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	libgithub "github.com/google/go-github/v83/github"
	"github.com/shurcooL/githubv4"
	slogctx "github.com/veqryn/slog-context"
)

// createPullRequestComment creates a comment on a pull request.
func createPullRequestComment(
	ctx context.Context,
	github *libgithub.Client,
	owner string,
	repo string,
	pullNumber int,
	body string,
) (*libgithub.IssueComment, error) {
	ctx = slogctx.Append(
		ctx,
		slog.String("owner", owner),
		slog.String("repo", repo),
		slog.Int("pull-number", pullNumber),
	)
	slog.Default().DebugContext(ctx, "creating pull request comment")

	// Always try and write a comment on the pull request
	limitedBody := limitCharacters(body)
	comment, _, err := github.Issues.CreateComment(
		ctx, owner, repo, pullNumber, &libgithub.IssueComment{
			Body: &limitedBody,
		},
	)
	if err != nil {
		err = fmt.Errorf("creating pull request comment: %w", err)
	}

	return comment, err
}

// findPullRequestCommentByMarker returns the first issue comment on the pull
// request whose body contains marker, or nil if there isn't one.
func findPullRequestCommentByMarker(
	ctx context.Context,
	github *libgithub.Client,
	owner string,
	repo string,
	pullNumber int,
	marker string,
) (*libgithub.IssueComment, error) {
	slog.Default().DebugContext(slogctx.Append(
		ctx,
		slog.String("owner", owner),
		slog.String("repo", repo),
		slog.Int("pull-number", pullNumber),
	), "listing pull request comments")

	opts := &libgithub.IssueListCommentsOptions{
		ListOptions: libgithub.ListOptions{PerPage: PerPageMax},
	}
	for {
		comments, res, err := github.Issues.ListComments(ctx, owner, repo, pullNumber, opts)
		if err != nil {
			return nil, fmt.Errorf("listing pull request comments: %w", err)
		}
		for _, c := range comments {
			if strings.Contains(c.GetBody(), marker) {
				return c, nil
			}
		}

		if res.NextPage == 0 {
			return nil, nil
		}

		opts.Page = res.NextPage
	}
}

// updatePullRequestComment replaces the body of an existing comment.
func updatePullRequestComment(
	ctx context.Context,
	github *libgithub.Client,
	owner string,
	repo string,
	commentID int64,
	body string,
) error {
	slog.Default().DebugContext(slogctx.Append(
		ctx,
		slog.String("owner", owner),
		slog.String("repo", repo),
		slog.Int64("comment-id", commentID),
	), "updating pull request comment")

	limitedBody := limitCharacters(body)
	_, _, err := github.Issues.EditComment(ctx, owner, repo, commentID, &libgithub.IssueComment{
		Body: &limitedBody,
	})
	if err != nil {
		return fmt.Errorf("updating pull request comment: %w", err)
	}

	return nil
}

// getPullRequest gets for the pull request details.
func getPullRequest(
	ctx context.Context,
	github *libgithub.Client,
	owner string,
	repo string,
	pullNumber int,
) (*libgithub.PullRequest, error) {
	ctx = slogctx.Append(
		ctx,
		slog.String("owner", owner),
		slog.String("repo", repo),
		slog.Int("pull-number", pullNumber),
	)
	slog.Default().DebugContext(ctx, "getting pull request details")

	pr, _, err := github.PullRequests.Get(ctx, owner, repo, pullNumber)
	if err != nil {
		return nil, err
	}

	return pr, nil
}

// closePullRequest closes a pull request.
func closePullRequest(
	ctx context.Context,
	github *libgithub.Client,
	owner string,
	repo string,
	pullNumber int,
) error {
	ctx = slogctx.Append(
		ctx,
		slog.String("owner", owner),
		slog.String("repo", repo),
		slog.Int("pull-number", pullNumber),
	)
	slog.Default().DebugContext(ctx, "closing pull request")

	_, _, err := github.PullRequests.Edit(ctx, owner, repo, pullNumber, &libgithub.PullRequest{
		State: libgithub.Ptr("closed"),
	})

	return err
}

// listPullRequestCommits lists all of the commits associated with a pull request.
func listPullRequestCommits(
	ctx context.Context,
	github *libgithub.Client,
	owner string,
	repo string,
	pullNumber int,
) ([]*libgithub.RepositoryCommit, error) {
	ctx = slogctx.Append(
		ctx,
		slog.String("owner", owner),
		slog.String("repo", repo),
		slog.Int("pull-number", pullNumber),
	)
	slog.Default().DebugContext(ctx, "listing pull request commits")

	opts := &libgithub.ListOptions{PerPage: PerPageMax}
	commits := []*libgithub.RepositoryCommit{}
	for {
		cs, res, err := github.PullRequests.ListCommits(ctx, owner, repo, pullNumber, opts)
		if err != nil {
			return nil, err
		}
		commits = append(commits, cs...)

		if res.NextPage == 0 {
			return commits, nil
		}

		opts.Page = res.NextPage
	}
}

// listPullRequestReviews lists all of the reviews associated with a pull request.
func listPullRequestReviews(
	ctx context.Context,
	github *libgithub.Client,
	owner string,
	repo string,
	pullNumber int,
) ([]*libgithub.PullRequestReview, error) {
	ctx = slogctx.Append(
		ctx,
		slog.String("owner", owner),
		slog.String("repo", repo),
		slog.Int("pull-number", pullNumber),
	)
	slog.Default().DebugContext(ctx, "listing pull request reviews")

	opts := &libgithub.ListOptions{PerPage: PerPageMax}
	reviews := []*libgithub.PullRequestReview{}
	for {
		rvs, res, err := github.PullRequests.ListReviews(ctx, owner, repo, pullNumber, opts)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, rvs...)

		if res.NextPage == 0 {
			return reviews, nil
		}

		opts.Page = res.NextPage
	}
}

// listPullRequestClosingIssues lists all "closing issues" associated with
// a pull request. These can be set using keywords associations or manually
// using the "development" sidebar on either an issue or pull request.
func listPullRequestClosingIssues(
	ctx context.Context,
	github *githubv4.Client,
	owner string,
	repo string,
	pullNumber int,
) ([]*ClosingIssueRef, error) {
	slog.Default().DebugContext(slogctx.Append(
		ctx,
		slog.String("owner", owner),
		slog.String("repo", repo),
		slog.Int("pull-number", pullNumber),
	), "getting pull request associated issues")

	oai := ClosingIssueRefs{}
	err := github.Query(ctx, &oai, map[string]any{
		"owner":  githubv4.String(owner),
		"repo":   githubv4.String(repo),
		"number": githubv4.Int(pullNumber),
	})
	if err != nil {
		return nil, err
	}

	ais := []*ClosingIssueRef{}
	for _, ai := range oai.Repository.PullRequest.ClosingIssuesReferences.Edges {
		ais = append(ais, ai.Node)
	}

	return ais, nil
}
