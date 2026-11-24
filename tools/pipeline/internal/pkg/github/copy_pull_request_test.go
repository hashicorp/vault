// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"errors"
	"testing"

	libgithub "github.com/google/go-github/v83/github"
	"github.com/stretchr/testify/require"
)

func Test_CopyPullRequest_getCoAuthoredByTrailers(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		commits  []*libgithub.RepositoryCommit
		expected []string
	}{
		"no commits": {
			nil,
			nil,
		},
		"one author one commit": {
			[]*libgithub.RepositoryCommit{
				{Commit: &libgithub.Commit{Author: &libgithub.CommitAuthor{
					Name:  libgithub.Ptr("John Doe"),
					Email: libgithub.Ptr("john@example.com"),
				}}},
			},
			[]string{"Co-Authored-By: John Doe <john@example.com>"},
		},
		"one author multiple commits": {
			[]*libgithub.RepositoryCommit{
				{Commit: &libgithub.Commit{Author: &libgithub.CommitAuthor{
					Name:  libgithub.Ptr("John Doe"),
					Email: libgithub.Ptr("john@example.com"),
				}}},
				{Commit: &libgithub.Commit{Author: &libgithub.CommitAuthor{
					Name:  libgithub.Ptr("John Doe"),
					Email: libgithub.Ptr("john@example.com"),
				}}},
			},
			[]string{"Co-Authored-By: John Doe <john@example.com>"},
		},
		"multiple authors with one commit each": {
			[]*libgithub.RepositoryCommit{
				{Commit: &libgithub.Commit{Author: &libgithub.CommitAuthor{
					Name:  libgithub.Ptr("John Doe"),
					Email: libgithub.Ptr("john@example.com"),
				}}},
				{Commit: &libgithub.Commit{Author: &libgithub.CommitAuthor{
					Name:  libgithub.Ptr("Jane Doe"),
					Email: libgithub.Ptr("jane@example.com"),
				}}},
			},
			[]string{"Co-Authored-By: John Doe <john@example.com>", "Co-Authored-By: Jane Doe <jane@example.com>"},
		},
		"multiple authors with multiple commits": {
			[]*libgithub.RepositoryCommit{
				{Commit: &libgithub.Commit{Author: &libgithub.CommitAuthor{
					Name:  libgithub.Ptr("John Doe"),
					Email: libgithub.Ptr("john@example.com"),
				}}},
				{Commit: &libgithub.Commit{Author: &libgithub.CommitAuthor{
					Name:  libgithub.Ptr("Jane Doe"),
					Email: libgithub.Ptr("jane@example.com"),
				}}},
				{Commit: &libgithub.Commit{Author: &libgithub.CommitAuthor{
					Name:  libgithub.Ptr("Jane Doe"),
					Email: libgithub.Ptr("jane@example.com"),
				}}},
				{Commit: &libgithub.Commit{Author: &libgithub.CommitAuthor{
					Name:  libgithub.Ptr("John Doe"),
					Email: libgithub.Ptr("john@example.com"),
				}}},
			},
			[]string{"Co-Authored-By: John Doe <john@example.com>", "Co-Authored-By: Jane Doe <jane@example.com>"},
		},
	} {
		t.Run(name, func(t *testing.T) {
			req := &CopyPullRequestReq{}
			require.EqualValues(t, test.expected, req.getCoAuthoredByTrailers(test.commits))
		})
	}
}

// TestCopyPullRequestRes_ToTable verifies that empty cells are empty strings
// rather than <nil>, which lets us hide empty columns and keeps <nil> out of
// the comment we leave on the pull request. It also makes sure an error without
// a request lands in the error column rather than in skipped merge commits.
func TestCopyPullRequestRes_ToTable(t *testing.T) {
	t.Parallel()

	req := &CopyPullRequestReq{
		FromOwner: "hashicorp",
		FromRepo:  "vault-enterprise",
		ToOwner:   "hashicorp",
		ToRepo:    "vault",
	}

	for name, test := range map[string]struct {
		res          *CopyPullRequestRes
		err          error
		wantContains []string
		wantMissing  []string
	}{
		"request without an error": {
			res:          &CopyPullRequestRes{Request: req},
			wantContains: []string{"hashicorp/vault-enterprise", "hashicorp/vault"},
			wantMissing:  []string{"<nil>", "ERROR"},
		},
		"request with an error": {
			res:          &CopyPullRequestRes{Request: req},
			err:          errors.New("copying failed"),
			wantContains: []string{"hashicorp/vault-enterprise", "ERROR", "copying failed"},
			wantMissing:  []string{"<nil>"},
		},
		"error without a request": {
			res:          &CopyPullRequestRes{},
			err:          errors.New("copying failed"),
			wantContains: []string{"ERROR", "copying failed"},
			wantMissing:  []string{"<nil>", "FROM", "SKIPPED MERGE COMMITS"},
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := test.res.ToTable(test.err).Render()
			for _, want := range test.wantContains {
				require.Contains(t, got, want)
			}
			for _, missing := range test.wantMissing {
				require.NotContains(t, got, missing)
			}

			md := test.res.ToMarkdownTable(test.err).RenderMarkdown()
			require.NotContains(t, md, "<nil>")
			require.NotContains(t, md, "&lt;nil&gt;")
		})
	}
}
