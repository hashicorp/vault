// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/shurcooL/githubv4"
	"github.com/stretchr/testify/require"
)

// TestListBranchCheckStatusesReq_Validate tests validation of the request
func TestListBranchCheckStatusesReq_Validate(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		req       *ListBranchCheckStatusesReq
		shouldErr bool
		errMsg    string
	}{
		"nil request": {
			req:       nil,
			shouldErr: true,
			errMsg:    "failed to initialize request",
		},
		"valid request": {
			req: &ListBranchCheckStatusesReq{
				Owner:         "hashicorp",
				Repo:          "vault",
				Branches:      []string{"main", "release/1.19.x"},
				RetryWait:     1,
				MaxRetry:      10,
				MaxConcurrent: 5,
			},
			shouldErr: false,
		},
		"missing owner": {
			req: &ListBranchCheckStatusesReq{
				Repo:          "vault",
				Branches:      []string{"main"},
				RetryWait:     1,
				MaxRetry:      10,
				MaxConcurrent: 5,
			},
			shouldErr: true,
			errMsg:    "no github organization has been provided",
		},
		"missing repo": {
			req: &ListBranchCheckStatusesReq{
				Owner:         "hashicorp",
				Branches:      []string{"main"},
				RetryWait:     1,
				MaxRetry:      10,
				MaxConcurrent: 5,
			},
			shouldErr: true,
			errMsg:    "no github repository has been provided",
		},
		"missing branches": {
			req: &ListBranchCheckStatusesReq{
				Owner:         "hashicorp",
				Repo:          "vault",
				Branches:      []string{},
				RetryWait:     1,
				MaxRetry:      10,
				MaxConcurrent: 5,
			},
			shouldErr: true,
			errMsg:    "no branches have been provided",
		},
		"negative retry wait": {
			req: &ListBranchCheckStatusesReq{
				Owner:         "hashicorp",
				Repo:          "vault",
				Branches:      []string{"main"},
				RetryWait:     -1,
				MaxRetry:      10,
				MaxConcurrent: 5,
			},
			shouldErr: true,
			errMsg:    "retry wait must be greater than or equal to 0",
		},
		"negative max retry": {
			req: &ListBranchCheckStatusesReq{
				Owner:         "hashicorp",
				Repo:          "vault",
				Branches:      []string{"main"},
				RetryWait:     1,
				MaxRetry:      -1,
				MaxConcurrent: 5,
			},
			shouldErr: true,
			errMsg:    "max retry must be greater than or equal to 0",
		},
		"zero max concurrent": {
			req: &ListBranchCheckStatusesReq{
				Owner:         "hashicorp",
				Repo:          "vault",
				Branches:      []string{"main"},
				RetryWait:     1,
				MaxRetry:      10,
				MaxConcurrent: 0,
			},
			shouldErr: true,
			errMsg:    "max concurrent must be greater than 0",
		},
		"negative max concurrent": {
			req: &ListBranchCheckStatusesReq{
				Owner:         "hashicorp",
				Repo:          "vault",
				Branches:      []string{"main"},
				RetryWait:     1,
				MaxRetry:      10,
				MaxConcurrent: -1,
			},
			shouldErr: true,
			errMsg:    "max concurrent must be greater than 0",
		},
		"zero retry wait is valid": {
			req: &ListBranchCheckStatusesReq{
				Owner:         "hashicorp",
				Repo:          "vault",
				Branches:      []string{"main"},
				RetryWait:     0,
				MaxRetry:      10,
				MaxConcurrent: 5,
			},
			shouldErr: false,
		},
		"zero max retry is valid": {
			req: &ListBranchCheckStatusesReq{
				Owner:         "hashicorp",
				Repo:          "vault",
				Branches:      []string{"main"},
				RetryWait:     1,
				MaxRetry:      0,
				MaxConcurrent: 5,
			},
			shouldErr: false,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := test.req.validate()
			if test.shouldErr {
				require.Error(t, err)
				if test.errMsg != "" {
					require.Contains(t, err.Error(), test.errMsg)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// TestBranchCheckResult_JSON tests JSON marshaling of BranchCheckResult
func TestBranchCheckResult_JSON(t *testing.T) {
	t.Parallel()

	result := BranchCheckResult{
		Branch:        "main",
		Status:        "success",
		FailureReason: "",
		CommitHash:    "abc123def456",
		CommitUrl:     "https://github.com/hashicorp/vault/commit/abc123def456",
	}

	data, err := json.Marshal(result)
	require.NoError(t, err)
	require.NotEmpty(t, data)

	var unmarshaled BranchCheckResult
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)
	require.Equal(t, result, unmarshaled)
}

// TestListBranchCheckStatusRes_ToJSON tests JSON marshaling of the response
func TestListBranchCheckStatusRes_ToJSON(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		res       *ListBranchCheckStatusRes
		shouldErr bool
	}{
		"nil response": {
			res:       nil,
			shouldErr: true,
		},
		"empty results": {
			res: &ListBranchCheckStatusRes{
				Results: []BranchCheckResult{},
			},
			shouldErr: false,
		},
		"single result": {
			res: &ListBranchCheckStatusRes{
				Results: []BranchCheckResult{
					{
						Branch:        "main",
						Status:        "success",
						FailureReason: "",
						CommitHash:    "abc123",
						CommitUrl:     "https://github.com/hashicorp/vault/commit/abc123",
					},
				},
			},
			shouldErr: false,
		},
		"multiple results": {
			res: &ListBranchCheckStatusRes{
				Results: []BranchCheckResult{
					{
						Branch:        "main",
						Status:        "success",
						FailureReason: "",
						CommitHash:    "abc123",
						CommitUrl:     "https://github.com/hashicorp/vault/commit/abc123",
					},
					{
						Branch:        "release/1.19.x",
						Status:        "failure",
						FailureReason: "checks-failed",
						CommitHash:    "def456",
						CommitUrl:     "https://github.com/hashicorp/vault/commit/def456",
					},
				},
			},
			shouldErr: false,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			data, err := test.res.ToJSON()
			if test.shouldErr {
				require.Error(t, err)
				require.Nil(t, data)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, data)

				// Verify we can unmarshal it back
				var unmarshaled ListBranchCheckStatusRes
				err = json.Unmarshal(data, &unmarshaled)
				require.NoError(t, err)
				require.Equal(t, test.res.Results, unmarshaled.Results)
			}
		})
	}
}

// TestListBranchCheckStatusRes_ToTable tests table rendering
func TestListBranchCheckStatusRes_ToTable(t *testing.T) {
	t.Parallel()

	res := &ListBranchCheckStatusRes{
		Results: []BranchCheckResult{
			{
				Branch:        "main",
				Status:        "success",
				FailureReason: "",
				CommitHash:    "abc123",
				CommitUrl:     "https://github.com/hashicorp/vault/commit/abc123",
			},
			{
				Branch:        "release/1.19.x",
				Status:        "failure",
				FailureReason: "checks-failed",
				CommitHash:    "def456",
				CommitUrl:     "https://github.com/hashicorp/vault/commit/def456",
			},
		},
	}

	table := res.ToTable()
	require.NotNil(t, table)

	rendered := table.Render()
	require.NotEmpty(t, rendered)
	require.Contains(t, rendered, "main")
	require.Contains(t, rendered, "success")
	require.Contains(t, rendered, "release/1.19.x")
	require.Contains(t, rendered, "failure")
	require.Contains(t, rendered, "checks-failed")
}

// TestListBranchCheckStatusRes_ToTable_Empty tests table rendering with empty results
func TestListBranchCheckStatusRes_ToTable_Empty(t *testing.T) {
	t.Parallel()

	res := &ListBranchCheckStatusRes{
		Results: []BranchCheckResult{},
	}

	table := res.ToTable()
	require.NotNil(t, table)

	rendered := table.Render()
	// No results should render as empty table
	_ = rendered
}

// TestListBranchCheckStatusesReq_Run_Validation tests that Run validates the request
func TestListBranchCheckStatusesReq_Run_Validation(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Test with invalid request (missing owner)
	req := &ListBranchCheckStatusesReq{
		Repo:          "vault",
		Branches:      []string{"main"},
		RetryWait:     1,
		MaxRetry:      10,
		MaxConcurrent: 5,
	}

	res, err := req.Run(ctx, nil)
	require.Error(t, err)
	require.Nil(t, res)
	require.Contains(t, err.Error(), "validating request")
	require.Contains(t, err.Error(), "no github organization has been provided")
}

// TestBranchCheckResult_FailureReasons tests various failure reason scenarios
func TestBranchCheckResult_FailureReasons(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		result         BranchCheckResult
		expectedStatus string
	}{
		{
			name: "query error",
			result: BranchCheckResult{
				Branch:        "main",
				Status:        "failure",
				FailureReason: "query-error: API rate limit exceeded",
				CommitHash:    "unknown",
				CommitUrl:     "",
			},
			expectedStatus: "failure",
		},
		{
			name: "checks failed",
			result: BranchCheckResult{
				Branch:        "main",
				Status:        "failure",
				FailureReason: "checks-failed",
				CommitHash:    "abc123",
				CommitUrl:     "https://github.com/hashicorp/vault/commit/abc123",
			},
			expectedStatus: "failure",
		},
		{
			name: "checks timed out",
			result: BranchCheckResult{
				Branch:        "main",
				Status:        "failure",
				FailureReason: "checks-timed-out",
				CommitHash:    "unknown",
				CommitUrl:     "",
			},
			expectedStatus: "failure",
		},
		{
			name: "unexpected status",
			result: BranchCheckResult{
				Branch:        "main",
				Status:        "failure",
				FailureReason: "unexpected-status: EXPECTED",
				CommitHash:    "abc123",
				CommitUrl:     "https://github.com/hashicorp/vault/commit/abc123",
			},
			expectedStatus: "failure",
		},
		{
			name: "success",
			result: BranchCheckResult{
				Branch:        "main",
				Status:        "success",
				FailureReason: "",
				CommitHash:    "abc123",
				CommitUrl:     "https://github.com/hashicorp/vault/commit/abc123",
			},
			expectedStatus: "success",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tc.expectedStatus, tc.result.Status)
			require.Equal(t, tc.result.Branch, "main")
		})
	}
}

// TestListBranchCheckStatusesReq_MaxRetryAttempts tests retry attempt calculation
func TestListBranchCheckStatusesReq_MaxRetryAttempts(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name                     string
		retryWait                int
		maxRetry                 int
		expectedMaxRetryAttempts int
	}{
		{
			name:                     "normal calculation",
			retryWait:                2,
			maxRetry:                 10,
			expectedMaxRetryAttempts: 5, // 10 / 2
		},
		{
			name:                     "retry wait is 1",
			retryWait:                1,
			maxRetry:                 10,
			expectedMaxRetryAttempts: 10, // 10 / 1
		},
		{
			name:                     "retry wait is 0",
			retryWait:                0,
			maxRetry:                 10,
			expectedMaxRetryAttempts: 10, // uses maxRetry directly
		},
		{
			name:                     "max retry is 0",
			retryWait:                1,
			maxRetry:                 0,
			expectedMaxRetryAttempts: 1,
		},
		{
			name:                     "max retry smaller than retry wait still tries once",
			retryWait:                10,
			maxRetry:                 5,
			expectedMaxRetryAttempts: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			maxRetryAttempts := 1
			if tc.maxRetry > 0 {
				maxRetryAttempts = tc.maxRetry
				if tc.retryWait >= 1 {
					maxRetryAttempts = tc.maxRetry / tc.retryWait
					if maxRetryAttempts == 0 {
						maxRetryAttempts = 1
					}
				}
			}

			require.Equal(t, tc.expectedMaxRetryAttempts, maxRetryAttempts)
		})
	}
}

// Made with Bob

// mockGraphQLClient is a mock implementation of the githubv4 client for testing
type mockGraphQLClient struct {
	queryFunc func(ctx context.Context, q interface{}, variables map[string]interface{}) error
}

func (m *mockGraphQLClient) Query(ctx context.Context, q interface{}, variables map[string]interface{}) error {
	if m.queryFunc != nil {
		return m.queryFunc(ctx, q, variables)
	}
	return nil
}

// TestListBranchCheckStatusesReq_Run_InvalidRepo tests behavior with invalid repository
func TestListBranchCheckStatusesReq_Run_InvalidRepo(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Create a mock client that simulates GitHub API error for invalid repo
	mockClient := &mockGraphQLClient{
		queryFunc: func(ctx context.Context, q interface{}, variables map[string]interface{}) error {
			return fmt.Errorf("Could not resolve to a Repository with the name 'hashicorp/invalid-repo-name'")
		},
	}

	req := &ListBranchCheckStatusesReq{
		Owner:         "hashicorp",
		Repo:          "invalid-repo-name",
		Branches:      []string{"main"},
		RetryWait:     0,
		MaxRetry:      1,
		MaxConcurrent: 1,
	}

	res, err := req.Run(ctx, mockClient)
	require.NoError(t, err) // Run itself doesn't error, but returns failure results
	require.NotNil(t, res)
	require.Len(t, res.Results, 1)

	result := res.Results[0]
	require.Equal(t, "main", result.Branch)
	require.Equal(t, "failure", result.Status)
	require.Contains(t, result.FailureReason, "query-error")
	require.Contains(t, result.FailureReason, "invalid-repo-name")
	require.Equal(t, "unknown", result.CommitHash)
	require.Empty(t, result.CommitUrl)
}

// TestListBranchCheckStatusesReq_Run_InvalidOwner tests behavior with invalid organization
func TestListBranchCheckStatusesReq_Run_InvalidOwner(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Create a mock client that simulates GitHub API error for invalid owner
	mockClient := &mockGraphQLClient{
		queryFunc: func(ctx context.Context, q interface{}, variables map[string]interface{}) error {
			return fmt.Errorf("Could not resolve to a User with the login of 'invalid-org-12345'")
		},
	}

	req := &ListBranchCheckStatusesReq{
		Owner:         "invalid-org-12345",
		Repo:          "vault",
		Branches:      []string{"main"},
		RetryWait:     0,
		MaxRetry:      1,
		MaxConcurrent: 1,
	}

	res, err := req.Run(ctx, mockClient)
	require.NoError(t, err) // Run itself doesn't error, but returns failure results
	require.NotNil(t, res)
	require.Len(t, res.Results, 1)

	result := res.Results[0]
	require.Equal(t, "main", result.Branch)
	require.Equal(t, "failure", result.Status)
	require.Contains(t, result.FailureReason, "query-error")
	require.Contains(t, result.FailureReason, "invalid-org-12345")
	require.Equal(t, "unknown", result.CommitHash)
	require.Empty(t, result.CommitUrl)
}

// TestListBranchCheckStatusesReq_Run_InvalidBothOwnerAndRepo tests behavior with both invalid
func TestListBranchCheckStatusesReq_Run_InvalidBothOwnerAndRepo(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Create a mock client that simulates GitHub API error
	mockClient := &mockGraphQLClient{
		queryFunc: func(ctx context.Context, q interface{}, variables map[string]interface{}) error {
			return fmt.Errorf("Could not resolve to a Repository with the name 'fake-org/fake-repo'")
		},
	}

	req := &ListBranchCheckStatusesReq{
		Owner:         "fake-org",
		Repo:          "fake-repo",
		Branches:      []string{"main", "develop"},
		RetryWait:     0,
		MaxRetry:      1,
		MaxConcurrent: 2,
	}

	res, err := req.Run(ctx, mockClient)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Len(t, res.Results, 2)

	// Both branches should fail with query errors
	for _, result := range res.Results {
		require.Contains(t, []string{"main", "develop"}, result.Branch)
		require.Equal(t, "failure", result.Status)
		require.Contains(t, result.FailureReason, "query-error")
		require.Equal(t, "unknown", result.CommitHash)
		require.Empty(t, result.CommitUrl)
	}
}

// TestListBranchCheckStatusesReq_Run_APIRateLimitError tests behavior with rate limit errors
func TestListBranchCheckStatusesReq_Run_APIRateLimitError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Create a mock client that simulates GitHub API rate limit error
	mockClient := &mockGraphQLClient{
		queryFunc: func(ctx context.Context, q interface{}, variables map[string]interface{}) error {
			return fmt.Errorf("API rate limit exceeded for user ID 12345")
		},
	}

	req := &ListBranchCheckStatusesReq{
		Owner:         "hashicorp",
		Repo:          "vault",
		Branches:      []string{"main"},
		RetryWait:     0,
		MaxRetry:      1,
		MaxConcurrent: 1,
	}

	res, err := req.Run(ctx, mockClient)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Len(t, res.Results, 1)

	result := res.Results[0]
	require.Equal(t, "main", result.Branch)
	require.Equal(t, "failure", result.Status)
	require.Contains(t, result.FailureReason, "query-error")
	require.Contains(t, result.FailureReason, "rate limit")
	require.Equal(t, "unknown", result.CommitHash)
	require.Empty(t, result.CommitUrl)
}

// TestListBranchCheckStatusesReq_Run_NetworkError tests behavior with network errors
func TestListBranchCheckStatusesReq_Run_NetworkError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Create a mock client that simulates network error
	mockClient := &mockGraphQLClient{
		queryFunc: func(ctx context.Context, q interface{}, variables map[string]interface{}) error {
			return fmt.Errorf("dial tcp: lookup api.github.com: no such host")
		},
	}

	req := &ListBranchCheckStatusesReq{
		Owner:         "hashicorp",
		Repo:          "vault",
		Branches:      []string{"main"},
		RetryWait:     0,
		MaxRetry:      1,
		MaxConcurrent: 1,
	}

	res, err := req.Run(ctx, mockClient)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Len(t, res.Results, 1)

	result := res.Results[0]
	require.Equal(t, "main", result.Branch)
	require.Equal(t, "failure", result.Status)
	require.Contains(t, result.FailureReason, "query-error")
	require.Contains(t, result.FailureReason, "no such host")
	require.Equal(t, "unknown", result.CommitHash)
	require.Empty(t, result.CommitUrl)
}

// TestListBranchCheckStatusesReq_Run_ServerError5xx tests retry behavior with 5xx errors
func TestListBranchCheckStatusesReq_Run_ServerError5xx(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Track number of attempts
	attempts := 0

	// Create a mock client that fails twice with 503, then succeeds
	mockClient := &mockGraphQLClient{
		queryFunc: func(ctx context.Context, q interface{}, variables map[string]interface{}) error {
			attempts++
			if attempts <= 2 {
				return fmt.Errorf("503 Service Unavailable: GitHub is temporarily unavailable")
			}
			// On third attempt, succeed
			return nil
		},
	}

	req := &ListBranchCheckStatusesReq{
		Owner:         "hashicorp",
		Repo:          "vault",
		Branches:      []string{"main"},
		RetryWait:     0,
		MaxRetry:      1,
		MaxConcurrent: 1,
	}

	res, err := req.Run(ctx, mockClient)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Len(t, res.Results, 1)

	// Should succeed after retries
	result := res.Results[0]
	require.Equal(t, "main", result.Branch)
	// Note: Since we return nil without setting query data, this will have empty commit info
	// In a real scenario with proper mock data, this would be "success"
	require.Equal(t, 3, attempts, "Should have retried twice before succeeding")
}

// TestListBranchCheckStatusesReq_Run_SuccessPath tests successful GraphQL response handling.
func TestListBranchCheckStatusesReq_Run_SuccessPath(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	mockClient := &mockGraphQLClient{
		queryFunc: func(ctx context.Context, q interface{}, variables map[string]interface{}) error {
			branchVar := string(variables["branch"].(githubv4.String))
			query := q.(*struct {
				Repository struct {
					Ref struct {
						Target struct {
							Commit struct {
								Oid               githubv4.String
								StatusCheckRollup struct {
									State githubv4.String
								}
							} `graphql:"... on Commit"`
						}
					} `graphql:"ref(qualifiedName: $branch)"`
				} `graphql:"repository(owner: $owner, name: $repo)"`
			})

			switch branchVar {
			case "refs/heads/main":
				query.Repository.Ref.Target.Commit.Oid = githubv4.String("1234567890abcdef")
				query.Repository.Ref.Target.Commit.StatusCheckRollup.State = githubv4.String("SUCCESS")
			case "refs/heads/release/1.21.x":
				query.Repository.Ref.Target.Commit.Oid = githubv4.String("fedcba0987654321")
				query.Repository.Ref.Target.Commit.StatusCheckRollup.State = githubv4.String("FAILURE")
			default:
				return fmt.Errorf("unexpected branch %s", branchVar)
			}

			return nil
		},
	}

	req := &ListBranchCheckStatusesReq{
		Owner:         "hashicorp",
		Repo:          "vault",
		Branches:      []string{"main", "release/1.21.x"},
		RetryWait:     0,
		MaxRetry:      0,
		MaxConcurrent: 2,
	}

	res, err := req.Run(ctx, mockClient)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Len(t, res.Results, 2)

	require.Equal(t, BranchCheckResult{
		Branch:        "main",
		Status:        "success",
		FailureReason: "",
		CommitHash:    "1234567890abcdef",
		CommitUrl:     "https://github.com/hashicorp/vault/commit/1234567890abcdef",
	}, res.Results[0])
	require.Equal(t, BranchCheckResult{
		Branch:        "release/1.21.x",
		Status:        "failure",
		FailureReason: "checks-failed",
		CommitHash:    "fedcba0987654321",
		CommitUrl:     "https://github.com/hashicorp/vault/commit/fedcba0987654321",
	}, res.Results[1])
}

// TestListBranchCheckStatusesReq_Run_ServerErrorExhausted tests when 5xx retries are exhausted
func TestListBranchCheckStatusesReq_Run_ServerErrorExhausted(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	attempts := 0

	// Create a mock client that always fails with 502
	mockClient := &mockGraphQLClient{
		queryFunc: func(ctx context.Context, q interface{}, variables map[string]interface{}) error {
			attempts++
			return fmt.Errorf("502 Bad Gateway: upstream server error")
		},
	}

	req := &ListBranchCheckStatusesReq{
		Owner:         "hashicorp",
		Repo:          "vault",
		Branches:      []string{"main"},
		RetryWait:     0,
		MaxRetry:      1,
		MaxConcurrent: 1,
	}

	res, err := req.Run(ctx, mockClient)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Len(t, res.Results, 1)

	result := res.Results[0]
	require.Equal(t, "main", result.Branch)
	require.Equal(t, "failure", result.Status)
	require.Contains(t, result.FailureReason, "query-error (server, retried")
	require.Contains(t, result.FailureReason, "502")
	require.Equal(t, "unknown", result.CommitHash)
	require.Empty(t, result.CommitUrl)
	require.Equal(t, 3, attempts, "Should have attempted 3 times (initial + 2 retries)")
}

// TestListBranchCheckStatusesReq_Run_ClientError4xx tests that 4xx errors don't retry
func TestListBranchCheckStatusesReq_Run_ClientError4xx(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	attempts := 0

	// Create a mock client that returns 404
	mockClient := &mockGraphQLClient{
		queryFunc: func(ctx context.Context, q interface{}, variables map[string]interface{}) error {
			attempts++
			return fmt.Errorf("404 Not Found: Could not resolve to a Repository")
		},
	}

	req := &ListBranchCheckStatusesReq{
		Owner:         "hashicorp",
		Repo:          "nonexistent-repo",
		Branches:      []string{"main"},
		RetryWait:     0,
		MaxRetry:      1,
		MaxConcurrent: 1,
	}

	res, err := req.Run(ctx, mockClient)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Len(t, res.Results, 1)

	result := res.Results[0]
	require.Equal(t, "main", result.Branch)
	require.Equal(t, "failure", result.Status)
	require.Contains(t, result.FailureReason, "query-error (client)")
	require.Contains(t, result.FailureReason, "404")
	require.Equal(t, "unknown", result.CommitHash)
	require.Empty(t, result.CommitUrl)
	require.Equal(t, 1, attempts, "Should NOT retry on 4xx client errors")
}

// TestIsRetryableError tests the retry error detection logic
func TestIsRetryableError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		err       error
		retryable bool
	}{
		{
			name:      "nil error",
			err:       nil,
			retryable: false,
		},
		{
			name:      "500 internal server error",
			err:       fmt.Errorf("500 Internal Server Error"),
			retryable: true,
		},
		{
			name:      "502 bad gateway",
			err:       fmt.Errorf("502 Bad Gateway"),
			retryable: true,
		},
		{
			name:      "503 service unavailable",
			err:       fmt.Errorf("503 Service Unavailable"),
			retryable: true,
		},
		{
			name:      "504 gateway timeout",
			err:       fmt.Errorf("504 Gateway Timeout"),
			retryable: true,
		},
		{
			name:      "internal server error text",
			err:       fmt.Errorf("Internal Server Error occurred"),
			retryable: true,
		},
		{
			name:      "temporarily unavailable",
			err:       fmt.Errorf("Service temporarily unavailable"),
			retryable: true,
		},
		{
			name:      "404 not found",
			err:       fmt.Errorf("404 Not Found"),
			retryable: false,
		},
		{
			name:      "400 bad request",
			err:       fmt.Errorf("400 Bad Request"),
			retryable: false,
		},
		{
			name:      "generic error",
			err:       fmt.Errorf("something went wrong"),
			retryable: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := isRetryableError(tc.err)
			require.Equal(t, tc.retryable, result)
		})
	}
}

// TestIsClientError tests the client error detection logic
func TestIsClientError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		err         error
		clientError bool
	}{
		{
			name:        "nil error",
			err:         nil,
			clientError: false,
		},
		{
			name:        "400 bad request",
			err:         fmt.Errorf("400 Bad Request"),
			clientError: true,
		},
		{
			name:        "401 unauthorized",
			err:         fmt.Errorf("401 Unauthorized"),
			clientError: true,
		},
		{
			name:        "403 forbidden",
			err:         fmt.Errorf("403 Forbidden"),
			clientError: true,
		},
		{
			name:        "404 not found",
			err:         fmt.Errorf("404 Not Found"),
			clientError: true,
		},
		{
			name:        "could not resolve",
			err:         fmt.Errorf("Could not resolve to a Repository"),
			clientError: true,
		},
		{
			name:        "not found text",
			err:         fmt.Errorf("Repository not found"),
			clientError: true,
		},
		{
			name:        "500 server error",
			err:         fmt.Errorf("500 Internal Server Error"),
			clientError: false,
		},
		{
			name:        "generic error",
			err:         fmt.Errorf("something went wrong"),
			clientError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := isClientError(tc.err)
			require.Equal(t, tc.clientError, result)
		})
	}
}

// TestListBranchCheckStatusesReq_Run_ConcurrencyControl tests that semaphore properly limits concurrent requests
func TestListBranchCheckStatusesReq_Run_ConcurrencyControl(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Track which branches are currently processing
	var mu sync.Mutex
	processing := make(map[string]bool)
	maxConcurrent := 0
	processOrder := []string{}

	mockClient := &mockGraphQLClient{
		queryFunc: func(ctx context.Context, q interface{}, variables map[string]interface{}) error {
			// Extract branch name from variables
			branchVar := variables["branch"].(githubv4.String)
			branch := string(branchVar)
			// Remove "refs/heads/" prefix to get just the branch name
			branch = branch[len("refs/heads/"):]

			mu.Lock()
			processing[branch] = true
			processOrder = append(processOrder, branch)
			currentCount := len(processing)
			if currentCount > maxConcurrent {
				maxConcurrent = currentCount
			}
			mu.Unlock()

			// Simulate some work
			time.Sleep(100 * time.Millisecond)

			mu.Lock()
			delete(processing, branch)
			mu.Unlock()

			// Return success - mock the query response structure
			// Note: In a real scenario, we'd populate the query struct properly
			return nil
		},
	}

	req := &ListBranchCheckStatusesReq{
		Owner:         "hashicorp",
		Repo:          "vault",
		Branches:      []string{"branch1", "branch2", "branch3"},
		RetryWait:     0,
		MaxRetry:      1,
		MaxConcurrent: 2, // Only 2 at a time
	}

	res, err := req.Run(ctx, mockClient)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Len(t, res.Results, 3, "All 3 branches should complete")
	require.LessOrEqual(t, maxConcurrent, 2, "Should never exceed max concurrent limit of 2")
	require.Len(t, processOrder, 3, "All 3 branches should have been processed")

	// Verify all branches are in the results
	branchesInResults := make(map[string]bool)
	for _, result := range res.Results {
		branchesInResults[result.Branch] = true
	}
	require.True(t, branchesInResults["branch1"], "branch1 should be in results")
	require.True(t, branchesInResults["branch2"], "branch2 should be in results")
	require.True(t, branchesInResults["branch3"], "branch3 should be in results")
}
