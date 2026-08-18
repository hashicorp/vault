// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"testing"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/github"
	"github.com/stretchr/testify/require"
)

// TestGithubCheckWorkflowRunCmd_TableOutput verifies that table output
// is generated correctly.
func TestGithubCheckWorkflowRunCmd_TableOutput(t *testing.T) {
	res := &github.CheckWorkflowRunRes{
		Success: true,
		Results: []*github.CheckWorkflowRunBranchResult{
			{
				Branch:        "main",
				Success:       true,
				WorkflowRunID: 12345,
				Status:        "completed",
				Conclusion:    "success",
				WorkflowURL:   "https://github.com/hashicorp/vault/actions/runs/12345",
				MatchedJobs: []*github.CheckWorkflowRunJobMatch{
					{
						Name:       "build-linux-amd64",
						Pattern:    "^build-.*-amd64$",
						Status:     "completed",
						Conclusion: "success",
					},
				},
			},
		},
	}

	output := res.ToTable()
	require.NotEmpty(t, output)
	require.Contains(t, output, "Branch Summary:")
	require.Contains(t, output, "main")
	require.Contains(t, output, "true")
	require.Contains(t, output, "completed")
	require.Contains(t, output, "success")
	require.Contains(t, output, "Matched Jobs:")
	require.Contains(t, output, "build-linux-amd64")
}

// TestGithubCheckWorkflowRunCmd_MarkdownOutput verifies that markdown output
// is generated correctly.
func TestGithubCheckWorkflowRunCmd_MarkdownOutput(t *testing.T) {
	res := &github.CheckWorkflowRunRes{
		Success: true,
		Results: []*github.CheckWorkflowRunBranchResult{
			{
				Branch:        "main",
				Success:       true,
				WorkflowRunID: 12345,
				Status:        "completed",
				Conclusion:    "success",
				MatchedJobs: []*github.CheckWorkflowRunJobMatch{
					{
						Name:       "build-linux-amd64",
						Pattern:    "^build-.*-amd64$",
						Status:     "completed",
						Conclusion: "success",
					},
				},
			},
		},
	}

	output := res.ToMarkdown()
	require.NotEmpty(t, output)
	require.Contains(t, output, "## Branch Summary")
	require.Contains(t, output, "| Branch | Success | Depth |")
	require.Contains(t, output, "| main | true | 0 |")
	require.Contains(t, output, "completed | success")
	require.Contains(t, output, "## Matched Jobs")
	require.Contains(t, output, "build-linux-amd64")
}

// TestGithubCheckWorkflowRunCmd_FailureWithBlockingInfo verifies that failure
// responses include blocking information.
func TestGithubCheckWorkflowRunCmd_FailureWithBlockingInfo(t *testing.T) {
	res := &github.CheckWorkflowRunRes{
		Success: false,
		Results: []*github.CheckWorkflowRunBranchResult{
			{
				Branch:          "main",
				Success:         false,
				WorkflowRunID:   12345,
				Status:          "completed",
				Conclusion:      "failure",
				MissingPatterns: []string{"test-integration"},
				FailedJobs: []*github.CheckWorkflowRunJobMatch{
					{
						Name:       "build-linux-amd64",
						Pattern:    "^build-.*-amd64$",
						Status:     "completed",
						Conclusion: "failure",
					},
				},
				FailureDetails: &github.CheckWorkflowRunFailureDetails{
					Reason:          "Required jobs failed or missing",
					MissingPatterns: []string{"test-integration"},
					FailedJobs: []*github.CheckWorkflowRunJobMatch{
						{
							Name:       "build-linux-amd64",
							Pattern:    "^build-.*-amd64$",
							Status:     "completed",
							Conclusion: "failure",
						},
					},
				},
			},
		},
	}

	output := res.ToTable()
	require.NotEmpty(t, output)
	require.Contains(t, output, "Branch Summary:")
	require.Contains(t, output, "main")
	require.Contains(t, output, "false")
	require.Contains(t, output, "Failed Jobs:")
	require.Contains(t, output, "Missing Patterns:")
	require.Contains(t, output, "test-integration")
	require.Contains(t, output, "build-linux-amd64")
}
