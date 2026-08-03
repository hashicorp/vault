// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"context"
	"log/slog"
	"slices"

	libgithub "github.com/google/go-github/v83/github"
	slogctx "github.com/veqryn/slog-context"
)

// addReviewers requests reviews from the given logins on the pull request
func addReviewers(
	ctx context.Context,
	github *libgithub.Client,
	owner string,
	repo string,
	number int,
	logins []string,
) error {
	logins = slices.Compact(slices.DeleteFunc(logins, func(a string) bool {
		return a == ""
	}))
	ctx = slogctx.Append(ctx, slog.Any("reviewer-logins", logins))

	if len(logins) < 1 {
		slog.Default().InfoContext(ctx, "skipping pull request review requests because no logins were provided")
		return nil
	}

	slog.Default().DebugContext(ctx, "requesting reviews on pull request")
	_, _, err := github.PullRequests.RequestReviewers(ctx, owner, repo, number, libgithub.ReviewersRequest{
		Reviewers: logins,
	})
	return err
}
