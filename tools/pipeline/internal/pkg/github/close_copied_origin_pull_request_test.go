// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"context"
	"os"
	"testing"

	"github.com/shurcooL/githubv4"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

// TestAcc_AssociatedIssues does an actual request against a known PR to check
// whether or not or query works as expected.
func TestAcc_AssociatedIssues(t *testing.T) {
	t.Parallel()

	token, setToken := os.LookupEnv("GITHUB_TOKEN")
	_, setACC := os.LookupEnv("PIPELINE_ACC")
	if !setACC && setToken {
		t.Skip("GITHUB_TOKEN and PIPELINE_ACC are not set")
	}

	github := githubv4.NewClient(
		oauth2.NewClient(context.Background(),
			oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}),
		),
	)

	ai := &ClosingIssueRefs{}
	require.NoError(t, github.Query(t.Context(), ai, map[string]any{
		"owner":  githubv4.String("hashicorp"),
		"repo":   githubv4.String("vault-enterprise"),
		"number": githubv4.Int(9484),
	}))

	require.Equal(t, "hashicorp/vault-enterprise", ai.Repository.PullRequest.Repository.NameWithOwner)
	require.Equal(t, 9484, ai.Repository.PullRequest.Number)
	require.Len(t, ai.Repository.PullRequest.ClosingIssuesReferences.Edges, 1)
	require.True(t, ai.Repository.PullRequest.ClosingIssuesReferences.Edges[0].Node.Closed)
	require.Equal(t, 31545, ai.Repository.PullRequest.ClosingIssuesReferences.Edges[0].Node.Number)
	require.Equal(t,
		"https://github.com/hashicorp/vault/issues/31545",
		ai.Repository.PullRequest.ClosingIssuesReferences.Edges[0].Node.URL,
	)
}
