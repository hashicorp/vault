// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package github_test

import (
	"errors"
	"testing"

	gh "github.com/google/go-github/v83/github"
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/github"
	"github.com/stretchr/testify/require"
)

// TestMarkdownLiteralCells checks that cell text isn't treated as Markdown,
// while our PR links still work and plain-text output stays the same.
func TestMarkdownLiteralCells(t *testing.T) {
	t.Parallel()
	const raw = "literal_*`[x]<tag>\\|next\nline"
	const escaped = `literal&#95;&#42;&#96;&#91;x&#93;&lt;tag&gt;&#92;\|next<br/>line`
	errText := errors.New(raw)
	for name, render := range map[string]func(*testing.T) (string, string){
		"workflow": func(t *testing.T) (string, string) {
			r := &github.CheckWorkflowRunRes{Results: []*github.CheckWorkflowRunBranchResult{{
				Branch: raw, Error: raw,
				MatchedJobs: []*github.CheckWorkflowRunJobMatch{{Name: raw, Pattern: raw}},
				FailureDetails: &github.CheckWorkflowRunFailureDetails{
					FailedJobs:      []*github.CheckWorkflowRunJobMatch{{Name: raw, Pattern: raw}},
					TaintedJobs:     []*github.CheckWorkflowRunJobMatch{{Name: raw}},
					MissingPatterns: []string{raw},
				},
			}}}
			md := r.ToMarkdown()
			require.Contains(t, md, "## Missing Patterns")
			return md, r.ToTable()
		},
		"open PRs": func(t *testing.T) (string, string) {
			r := &github.ListOpenPRsRes{BranchPRs: map[string][]*github.OpenPullRequest{
				raw: {{Milestone: raw, Backports: []string{raw}}},
			}}
			return r.ToMarkdown(), r.ToTable()
		},
		"sync go mod": func(t *testing.T) (string, string) {
			r := &github.SyncGoModRes{ToBranch: raw}
			return r.ToMarkdown(errText), r.ToString(errText)
		},
		"go mod diff": func(t *testing.T) (string, string) {
			r := &github.CheckGoModDiffRes{Diffs: []*github.CheckGoModDiff{{Path: raw, Error: raw}}}
			md, err := r.ToMarkdownTable(nil)
			require.NoError(t, err)
			plain, err := r.ToTable(nil)
			require.NoError(t, err)
			return md.RenderMarkdown(), plain.Render()
		},
		"backport": func(t *testing.T) (string, string) {
			r := &github.CreateBackportRes{Attempts: map[string]*github.CreateBackportAttempt{
				"1": {BaseRef: raw, TargetRef: raw, SkippedReason: raw, Error: errText},
			}}
			return r.CommentBody(), r.ToTable().Render()
		},
		"copy PR": func(t *testing.T) (string, string) {
			r := &github.CopyPullRequestRes{
				Request:           &github.CopyPullRequestReq{FromOwner: "o", FromRepo: "r", ToOwner: "o", ToRepo: "r"},
				OriginPullRequest: &gh.PullRequest{Number: gh.Ptr(1), HTMLURL: gh.Ptr("https://github.com/o/r/pull/1")},
			}
			md := r.ToMarkdownTable(errText).RenderMarkdown()
			require.Contains(t, md, "[o/r#1](https://github.com/o/r/pull/1)")
			require.Contains(t, r.CommentBody(errText), escaped)
			return md, r.ToTable(errText).Render()
		},
		"close PR": func(t *testing.T) (string, string) {
			r := &github.CloseCopiedOriginPullRequestRes{}
			return r.ToMarkdownTable(errText).RenderMarkdown(), r.ToTable(errText).Render()
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			md, plain := render(t)
			require.Contains(t, md, escaped)
			require.NotContains(t, md, "<tag>")
			require.Contains(t, plain, "literal_*`[x]<tag>\\|next")
			require.NotContains(t, plain, "&#95;")
		})
	}
}
