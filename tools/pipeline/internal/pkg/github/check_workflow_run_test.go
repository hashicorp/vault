// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"context"
	"regexp"
	"testing"

	gh "github.com/google/go-github/v83/github"
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/releases"
	"github.com/stretchr/testify/require"
)

// TestCheckWorkflowRunReq_validate verifies request validation logic.
func TestCheckWorkflowRunReq_validate(t *testing.T) {
	tests := map[string]struct {
		req     *CheckWorkflowRunReq
		wantErr bool
		errMsg  string
	}{
		"nil request": {
			req:     nil,
			wantErr: true,
			errMsg:  "uninitialized",
		},
		"missing owner": {
			req: &CheckWorkflowRunReq{
				Repo:     "vault",
				Workflow: "build.yml",
				Branches: []string{"main"},
				Patterns: []string{"^build-.*$"},
			},
			wantErr: true,
			errMsg:  "owner is required",
		},
		"missing repo": {
			req: &CheckWorkflowRunReq{
				Owner:    "hashicorp",
				Workflow: "build.yml",
				Branches: []string{"main"},
				Patterns: []string{"^build-.*$"},
			},
			wantErr: true,
			errMsg:  "repo is required",
		},
		"missing workflow": {
			req: &CheckWorkflowRunReq{
				Owner:    "hashicorp",
				Repo:     "vault",
				Branches: []string{"main"},
				Patterns: []string{"^build-.*$"},
			},
			wantErr: true,
			errMsg:  "workflow is required",
		},
		"missing branches and releases": {
			req: &CheckWorkflowRunReq{
				Owner:    "hashicorp",
				Repo:     "vault",
				Workflow: "build.yml",
				Patterns: []string{"^build-.*$"},
			},
			wantErr: true,
			errMsg:  "either branches or releases configuration is required",
		},
		"missing patterns": {
			req: &CheckWorkflowRunReq{
				Owner:    "hashicorp",
				Repo:     "vault",
				Workflow: "build.yml",
				Branches: []string{"main"},
				Patterns: []string{},
			},
			wantErr: true,
			errMsg:  "at least one pattern is required",
		},
		"invalid regex pattern": {
			req: &CheckWorkflowRunReq{
				Owner:    "hashicorp",
				Repo:     "vault",
				Workflow: "build.yml",
				Branches: []string{"main"},
				Patterns: []string{"[invalid"},
			},
			wantErr: true,
			errMsg:  "compiling pattern",
		},
		"valid request": {
			req: &CheckWorkflowRunReq{
				Owner:    "hashicorp",
				Repo:     "vault",
				Workflow: "build.yml",
				Branches: []string{"main"},
				Patterns: []string{"^build-.*$", "^test-.*$"},
			},
			wantErr: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := tt.req.validate(context.Background())
			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// TestCheckWorkflowRunReq_determineBranches verifies branch resolution logic
// for different repository types and CEActive settings.
func TestCheckWorkflowRunReq_determineBranches(t *testing.T) {
	t.Parallel()

	mockReleases := &releases.DecodeRes{
		Config: &releases.VersionsConfig{
			Schema: 1,
			ActiveVersion: &releases.ActiveVersion{
				Versions: map[string]*releases.Version{
					"1.19.x": {CEActive: true, LTS: true},
					"1.18.x": {CEActive: true, LTS: false},
					"1.17.x": {CEActive: false, LTS: false},
					"1.16.x": {CEActive: false, LTS: true},
				},
			},
		},
	}

	tests := map[string]struct {
		repo             string
		onlyCE           bool
		onlyLTS          bool
		includeMain      bool
		ceBranchPrefix   string
		expectedBranches []string
	}{
		"enterprise repo without filters": {
			repo:           "vault-enterprise",
			onlyCE:         false,
			onlyLTS:        false,
			ceBranchPrefix: "",
			expectedBranches: []string{
				"release/1.19.x+ent",
				"release/1.18.x+ent",
				"release/1.17.x+ent",
				"release/1.16.x+ent",
			},
		},
		"enterprise repo with CE prefix": {
			repo:           "vault-enterprise",
			onlyCE:         false,
			onlyLTS:        false,
			ceBranchPrefix: "ce",
			expectedBranches: []string{
				"ce/release/1.19.x",
				"ce/release/1.18.x",
				"release/1.17.x+ent",
				"release/1.16.x+ent",
			},
		},
		"enterprise repo with only-ce flag": {
			repo:           "vault-enterprise",
			onlyCE:         true,
			onlyLTS:        false,
			ceBranchPrefix: "ce",
			expectedBranches: []string{
				"ce/release/1.19.x",
				"ce/release/1.18.x",
			},
		},
		"enterprise repo with only-lts flag": {
			repo:           "vault-enterprise",
			onlyCE:         false,
			onlyLTS:        true,
			ceBranchPrefix: "",
			expectedBranches: []string{
				"release/1.19.x+ent",
				"release/1.16.x+ent",
			},
		},
		"CE repo without filters - only CEActive branches": {
			repo:           "vault",
			onlyCE:         false,
			onlyLTS:        false,
			ceBranchPrefix: "",
			expectedBranches: []string{
				"release/1.19.x",
				"release/1.18.x",
			},
		},
		"CE repo with only-lts flag - only CEActive and LTS": {
			repo:           "vault",
			onlyCE:         false,
			onlyLTS:        true,
			ceBranchPrefix: "",
			expectedBranches: []string{
				"release/1.19.x",
			},
		},
		"CE repo with only-ce flag - only CEActive branches": {
			repo:           "vault",
			onlyCE:         true,
			onlyLTS:        false,
			ceBranchPrefix: "",
			expectedBranches: []string{
				"release/1.19.x",
				"release/1.18.x",
			},
		},
		"enterprise repo with include-main": {
			repo:           "vault-enterprise",
			onlyCE:         false,
			onlyLTS:        false,
			ceBranchPrefix: "",
			includeMain:    true,
			expectedBranches: []string{
				"main",
				"release/1.19.x+ent",
				"release/1.18.x+ent",
				"release/1.17.x+ent",
				"release/1.16.x+ent",
			},
		},
		"enterprise repo with include-main and CE prefix": {
			repo:           "vault-enterprise",
			onlyCE:         false,
			onlyLTS:        false,
			ceBranchPrefix: "ce",
			includeMain:    true,
			expectedBranches: []string{
				"main",
				"ce/main",
				"ce/release/1.19.x",
				"ce/release/1.18.x",
				"release/1.17.x+ent",
				"release/1.16.x+ent",
			},
		},
		"enterprise repo with include-main and only-ce": {
			repo:           "vault-enterprise",
			onlyCE:         true,
			onlyLTS:        false,
			ceBranchPrefix: "ce",
			includeMain:    true,
			expectedBranches: []string{
				"ce/main",
				"ce/release/1.19.x",
				"ce/release/1.18.x",
			},
		},
		"CE repo with include-main": {
			repo:           "vault",
			onlyCE:         false,
			onlyLTS:        false,
			ceBranchPrefix: "",
			includeMain:    true,
			expectedBranches: []string{
				"main",
				"release/1.19.x",
				"release/1.18.x",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			req := &CheckWorkflowRunReq{
				Owner:          "hashicorp",
				Repo:           tt.repo,
				Workflow:       "build.yml",
				Patterns:       []string{"^build-.*$"},
				Releases:       mockReleases,
				OnlyCE:         tt.onlyCE,
				OnlyLTS:        tt.onlyLTS,
				IncludeMain:    tt.includeMain,
				CEBranchPrefix: tt.ceBranchPrefix,
			}

			branches, err := req.determineBranches(context.Background())
			require.NoError(t, err)
			require.ElementsMatch(t, tt.expectedBranches, branches)
		})
	}
}

// TestEvaluateWorkflowRun verifies job evaluation logic.
func TestEvaluateWorkflowRun(t *testing.T) {
	tests := map[string]struct {
		jobs                []*gh.WorkflowJob
		patterns            []string
		expectedComplete    bool
		expectedSuccess     bool
		expectedTainted     bool
		expectedMatched     int
		expectedMissing     int
		expectedFailed      int
		expectedTaintedJobs int
	}{
		"all patterns matched and successful": {
			jobs: []*gh.WorkflowJob{
				{Name: gh.Ptr("build-linux-amd64"), Status: gh.Ptr("completed"), Conclusion: gh.Ptr("success")},
				{Name: gh.Ptr("test-go"), Status: gh.Ptr("completed"), Conclusion: gh.Ptr("success")},
			},
			patterns:            []string{"^build-.*-amd64$", "^test-go$"},
			expectedComplete:    true,
			expectedSuccess:     true,
			expectedTainted:     false,
			expectedMatched:     2,
			expectedMissing:     0,
			expectedFailed:      0,
			expectedTaintedJobs: 0,
		},
		"pattern missing": {
			jobs: []*gh.WorkflowJob{
				{Name: gh.Ptr("build-linux-amd64"), Status: gh.Ptr("completed"), Conclusion: gh.Ptr("success")},
			},
			patterns:            []string{"^build-.*-amd64$", "^test-go$"},
			expectedComplete:    false,
			expectedSuccess:     true,
			expectedTainted:     false,
			expectedMatched:     1,
			expectedMissing:     1,
			expectedFailed:      0,
			expectedTaintedJobs: 0,
		},
		"job failed": {
			jobs: []*gh.WorkflowJob{
				{Name: gh.Ptr("build-linux-amd64"), Status: gh.Ptr("completed"), Conclusion: gh.Ptr("failure")},
				{Name: gh.Ptr("test-go"), Status: gh.Ptr("completed"), Conclusion: gh.Ptr("success")},
			},
			patterns:            []string{"^build-.*-amd64$", "^test-go$"},
			expectedComplete:    true,
			expectedSuccess:     false,
			expectedTainted:     false,
			expectedMatched:     2,
			expectedMissing:     0,
			expectedFailed:      1,
			expectedTaintedJobs: 0,
		},
		"multiple jobs match same pattern": {
			jobs: []*gh.WorkflowJob{
				{Name: gh.Ptr("build-linux-amd64"), Status: gh.Ptr("completed"), Conclusion: gh.Ptr("success")},
				{Name: gh.Ptr("build-darwin-amd64"), Status: gh.Ptr("completed"), Conclusion: gh.Ptr("success")},
			},
			patterns:            []string{"^build-.*-amd64$"},
			expectedComplete:    true,
			expectedSuccess:     true,
			expectedTainted:     false,
			expectedMatched:     2,
			expectedMissing:     0,
			expectedFailed:      0,
			expectedTaintedJobs: 0,
		},
		"skipped jobs are successful": {
			jobs: []*gh.WorkflowJob{
				{Name: gh.Ptr("build-linux-amd64"), Status: gh.Ptr("completed"), Conclusion: gh.Ptr("skipped")},
			},
			patterns:            []string{"^build-.*-amd64$"},
			expectedComplete:    true,
			expectedSuccess:     true,
			expectedTainted:     false,
			expectedMatched:     1,
			expectedMissing:     0,
			expectedFailed:      0,
			expectedTaintedJobs: 0,
		},
		"incomplete with tainted jobs": {
			jobs: []*gh.WorkflowJob{
				{Name: gh.Ptr("build-linux"), Status: gh.Ptr("completed"), Conclusion: gh.Ptr("success")},
				{Name: gh.Ptr("setup"), Status: gh.Ptr("completed"), Conclusion: gh.Ptr("failure")},
			},
			patterns:            []string{"^build-.*$", "^test-.*$"},
			expectedComplete:    false,
			expectedSuccess:     true,
			expectedTainted:     true,
			expectedMatched:     1,
			expectedMissing:     1,
			expectedFailed:      0,
			expectedTaintedJobs: 1,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			req := &CheckWorkflowRunReq{
				Patterns:         tt.patterns,
				compiledPatterns: make([]*regexp.Regexp, len(tt.patterns)),
			}
			for i, pattern := range tt.patterns {
				req.compiledPatterns[i] = regexp.MustCompile(pattern)
			}

			eval := req.evaluateWorkflowRun(tt.jobs)

			require.Equal(t, tt.expectedComplete, eval.IsComplete, "IsComplete mismatch")
			require.Equal(t, tt.expectedSuccess, eval.IsSuccessful, "IsSuccessful mismatch")
			require.Equal(t, tt.expectedTainted, eval.IsTainted, "IsTainted mismatch")
			require.Len(t, eval.MatchedJobs, tt.expectedMatched, "MatchedJobs count mismatch")
			require.Len(t, eval.MissingPatterns, tt.expectedMissing, "MissingPatterns count mismatch")
			require.Len(t, eval.FailedJobs, tt.expectedFailed, "FailedJobs count mismatch")
			require.Len(t, eval.TaintedJobs, tt.expectedTaintedJobs, "TaintedJobs count mismatch")
		})
	}
}

// TestCheckWorkflowRunRes_ToJSON verifies JSON output.
func TestCheckWorkflowRunRes_ToJSON(t *testing.T) {
	tests := map[string]struct {
		res     *CheckWorkflowRunRes
		wantErr bool
	}{
		"nil response": {
			res:     nil,
			wantErr: true,
		},
		"valid response": {
			res: &CheckWorkflowRunRes{
				Success: true,
				Results: []*CheckWorkflowRunBranchResult{
					{
						Branch:        "main",
						Success:       true,
						WorkflowRunID: 12345,
						MatchedJobs: []*CheckWorkflowRunJobMatch{
							{Name: "build-linux", Pattern: "^build-.*$", Status: "completed", Conclusion: "success"},
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			b, err := tt.res.ToJSON()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, b)
			}
		})
	}
}

// TestCheckWorkflowRunRes_ToTable verifies table output.
func TestCheckWorkflowRunRes_ToTable(t *testing.T) {
	tests := map[string]struct {
		res      *CheckWorkflowRunRes
		contains string
	}{
		"nil response": {
			res:      nil,
			contains: "",
		},
		"success result": {
			res: &CheckWorkflowRunRes{
				Success: true,
				Results: []*CheckWorkflowRunBranchResult{
					{
						Branch:        "main",
						Success:       true,
						WorkflowRunID: 12345,
					},
				},
			},
			contains: "main",
		},
		"failure result": {
			res: &CheckWorkflowRunRes{
				Success: false,
				Results: []*CheckWorkflowRunBranchResult{
					{
						Branch:  "main",
						Success: false,
						FailureDetails: &CheckWorkflowRunFailureDetails{
							Reason: "Jobs failed",
						},
					},
				},
			},
			contains: "Branch Summary",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			output := tt.res.ToTable()
			if tt.contains != "" {
				require.Contains(t, output, tt.contains)
			}
		})
	}
}

// TestCheckWorkflowRunRes_ToMarkdown verifies markdown output.
func TestCheckWorkflowRunRes_ToMarkdown(t *testing.T) {
	tests := map[string]struct {
		res      *CheckWorkflowRunRes
		contains string
	}{
		"nil response": {
			res:      nil,
			contains: "",
		},
		"valid response": {
			res: &CheckWorkflowRunRes{
				Success: true,
				Results: []*CheckWorkflowRunBranchResult{
					{
						Branch:        "main",
						Success:       true,
						WorkflowRunID: 12345,
					},
				},
			},
			contains: "| Branch | Success |",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			output := tt.res.ToMarkdown()
			if tt.contains != "" {
				require.Contains(t, output, tt.contains)
			}
		})
	}
}

// TestCheckWorkflowRunRes_String verifies string output.
func TestCheckWorkflowRunRes_String(t *testing.T) {
	tests := map[string]struct {
		res      *CheckWorkflowRunRes
		contains string
	}{
		"nil response": {
			res:      nil,
			contains: "No result",
		},
		"success result": {
			res: &CheckWorkflowRunRes{
				Success: true,
				Results: []*CheckWorkflowRunBranchResult{
					{
						Branch:  "main",
						Success: true,
					},
				},
			},
			contains: "All 1 branches passed",
		},
		"failure result": {
			res: &CheckWorkflowRunRes{
				Success: false,
				Results: []*CheckWorkflowRunBranchResult{
					{
						Branch:  "main",
						Success: false,
					},
				},
			},
			contains: "1 of 1 branches failed",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			output := tt.res.String()
			require.Contains(t, output, tt.contains)
		})
	}
}

// TestBuildBranchResult verifies branch result builder.
func TestBuildBranchResult(t *testing.T) {
	tests := map[string]struct {
		branch  string
		run     *gh.WorkflowRun
		eval    *CheckWorkflowRunEvaluation
		success bool
		check   func(*testing.T, *CheckWorkflowRunBranchResult)
	}{
		"success result": {
			branch: "main",
			run: &gh.WorkflowRun{
				ID:         gh.Ptr(int64(12345)),
				Status:     gh.Ptr("completed"),
				Conclusion: gh.Ptr("success"),
			},
			eval: &CheckWorkflowRunEvaluation{
				IsComplete:   true,
				IsSuccessful: true,
				IsTainted:    false,
				MatchedJobs: []*CheckWorkflowRunJobMatch{
					{Name: "build-linux-amd64", Pattern: "^build-.*$", Status: "completed", Conclusion: "success"},
				},
			},
			success: true,
			check: func(t *testing.T, br *CheckWorkflowRunBranchResult) {
				require.True(t, br.Success)
				require.Equal(t, "main", br.Branch)
				require.Equal(t, int64(12345), br.WorkflowRunID)
				require.Len(t, br.MatchedJobs, 1)
				require.Nil(t, br.FailureDetails)
			},
		},
		"failure result": {
			branch: "main",
			run: &gh.WorkflowRun{
				ID:         gh.Ptr(int64(12345)),
				Status:     gh.Ptr("completed"),
				Conclusion: gh.Ptr("failure"),
			},
			eval: &CheckWorkflowRunEvaluation{
				IsComplete:      false,
				IsSuccessful:    false,
				IsTainted:       false,
				MissingPatterns: []string{"test-integration"},
				FailedJobs: []*CheckWorkflowRunJobMatch{
					{Name: "build-linux-amd64", Pattern: "^build-.*$", Status: "completed", Conclusion: "failure"},
				},
			},
			success: false,
			check: func(t *testing.T, br *CheckWorkflowRunBranchResult) {
				require.False(t, br.Success)
				require.Equal(t, "main", br.Branch)
				require.Len(t, br.MissingPatterns, 1)
				require.Len(t, br.FailedJobs, 1)
				require.NotNil(t, br.FailureDetails)
				require.Equal(t, "Required jobs failed or missing", br.FailureDetails.Reason)
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			req := &CheckWorkflowRunReq{}
			result := req.buildBranchResult(tt.branch, tt.run, tt.eval, tt.success)
			tt.check(t, result)
		})
	}
}
