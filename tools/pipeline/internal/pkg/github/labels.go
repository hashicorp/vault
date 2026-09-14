// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	libgithub "github.com/google/go-github/v83/github"
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

// removeLabelFromIssue removes a given label from the issue or PR.
// If the label is not present, the error is ignored.
func removeLabelFromIssue(
	ctx context.Context,
	github *libgithub.Client,
	owner string,
	repo string,
	number int,
	label string,
) error {
	if label == "" {
		return errors.New("no label was provided for removal")
	}

	ctx = slogctx.Append(
		ctx,
		slog.String("label", label),
		slog.Int("issue-number", number),
	)

	slog.Default().DebugContext(ctx, "removing label from issue or pull request")
	resp, err := github.Issues.RemoveLabelForIssue(ctx, owner, repo, number, label)
	if err != nil {
		// 404 means the label is not present — that's fine.
		if resp != nil && resp.StatusCode == 404 {
			slog.Default().DebugContext(ctx, "label not present on issue or pull request, skipping removal")
			return nil
		}
		return err
	}

	slog.Default().DebugContext(ctx, "successfully removed label from issue or pull request")
	return nil
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

	ctx = slogctx.Append(
		ctx,
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
