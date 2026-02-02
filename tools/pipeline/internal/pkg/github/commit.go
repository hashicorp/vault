// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"context"
	"log/slog"

	libgithub "github.com/google/go-github/v81/github"
	slogctx "github.com/veqryn/slog-context"
)

// listCommitStatuses lists all of the statuses associated with a commit.
func listCommitStatuses(
	ctx context.Context,
	github *libgithub.Client,
	owner string,
	repo string,
	ref string,
) ([]*libgithub.RepoStatus, error) {
	ctx = slogctx.Append(ctx,
		slog.String("owner", owner),
		slog.String("repo", repo),
		slog.String("commit-ref", ref),
	)
	slog.Default().DebugContext(ctx, "listing commit statuses")

	opts := &libgithub.ListOptions{PerPage: PerPageMax}
	statuses := []*libgithub.RepoStatus{}
	for {
		ss, res, err := github.Repositories.ListStatuses(ctx, owner, repo, ref, opts)
		if err != nil {
			return nil, err
		}
		statuses = append(statuses, ss...)

		if res.NextPage == 0 {
			return statuses, nil
		}

		opts.Page = res.NextPage
	}
}
