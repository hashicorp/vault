// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"testing"

	libgithub "github.com/google/go-github/v81/github"
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
