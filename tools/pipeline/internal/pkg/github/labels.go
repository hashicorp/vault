// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"context"
	"log/slog"
	"strings"

	libgithub "github.com/google/go-github/v74/github"
	slogctx "github.com/veqryn/slog-context"
)

// filterNonBackportLabels returns a slice of label names that do not have the
// specified backport prefix, filtering out backport labels from the input labels
func filterNonBackportLabels(labels Labels, backportPrefix string) []string {
	var labelsToAdd []string
	for _, label := range labels {
		if label.GetName() != "" && !strings.HasPrefix(label.GetName(), backportPrefix+"/") {
			labelsToAdd = append(labelsToAdd, label.GetName())
		}
	}
	return labelsToAdd
}

// addLabelsToIssue adds the given labels to the issue or pull request
func addLabelsToIssue(
	ctx context.Context,
	github *libgithub.Client,
	owner string,
	repo string,
	number int,
	labels []string,
) error {
	if len(labels) < 1 {
		slog.Default().DebugContext(ctx, "skipping label assignment because no labels were provided")
		return nil
	}

	ctx = slogctx.Append(ctx,
		slog.String("labels", strings.Join(labels, ", ")),
		slog.Int("issue-number", number),
	)

	slog.Default().DebugContext(ctx, "adding labels to issue or pull request")
	_, _, err := github.Issues.AddLabelsToIssue(ctx, owner, repo, number, labels)
	if err != nil {
		return err
	}

	slog.Default().DebugContext(ctx, "successfully added labels to issue or pull request")
	return nil
}
