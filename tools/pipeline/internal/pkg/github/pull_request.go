// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"context"
	"fmt"
	"log/slog"

	libgithub "github.com/google/go-github/v68/github"
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
	ctx = slogctx.Append(ctx,
		slog.String("owner", owner),
		slog.String("repo", repo),
		slog.Int("pull-number", pullNumber),
	)
	slog.Default().DebugContext(ctx, "creating pull request comment")

	// Always try and write a comment on the pull request
	comment, _, err := github.Issues.CreateComment(
		ctx, owner, repo, pullNumber, &libgithub.IssueComment{
			Body: &body,
		},
	)
	if err != nil {
		err = fmt.Errorf("creating pull request comment: %w", err)
	}

	return comment, err
}

// getPullRequest gets for the pull request details.
func getPullRequest(
	ctx context.Context,
	github *libgithub.Client,
	owner string,
	repo string,
	pullNumber int,
) (*libgithub.PullRequest, error) {
	ctx = slogctx.Append(ctx,
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

// listPullRequestCommits lists all of the commits associated with a pull request.
func listPullRequestCommits(
	ctx context.Context,
	github *libgithub.Client,
	owner string,
	repo string,
	pullNumber int,
) ([]*libgithub.RepositoryCommit, error) {
	ctx = slogctx.Append(ctx,
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
	ctx = slogctx.Append(ctx,
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
