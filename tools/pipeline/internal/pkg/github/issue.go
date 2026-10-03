// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"context"
	"log/slog"

	libgithub "github.com/google/go-github/v81/github"
	slogctx "github.com/veqryn/slog-context"
)

// closeIssue closes an issue.
func closeIssue(
	ctx context.Context,
	github *libgithub.Client,
	owner string,
	repo string,
	issueNumber int,
) error {
	ctx = slogctx.Append(ctx,
		slog.String("owner", owner),
		slog.String("repo", repo),
		slog.Int("issue-number", issueNumber),
	)
	slog.Default().DebugContext(ctx, "closing issue")

	_, _, err := github.Issues.Edit(ctx, owner, repo, issueNumber, &libgithub.IssueRequest{
		State: libgithub.Ptr("closed"),
	})

	return err
}
