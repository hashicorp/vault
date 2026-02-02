// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"context"
	"log/slog"
	"slices"

	libgithub "github.com/google/go-github/v81/github"
	slogctx "github.com/veqryn/slog-context"
)

// addAssignees assigns the given logins to the issue or pull request
func addAssignees(
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
	ctx = slogctx.Append(ctx, slog.Any("assignee-logins", logins))

	if len(logins) < 1 {
		slog.Default().InfoContext(ctx, "skipping pull request actor assignments because no logins were provided")
		return nil
	}

	slog.Default().DebugContext(ctx, "adding assignees to pull request")
	_, _, err := github.Issues.AddAssignees(ctx, owner, repo, number, logins)
	return err
}
