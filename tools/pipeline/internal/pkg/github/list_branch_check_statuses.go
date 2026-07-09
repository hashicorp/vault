// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/shurcooL/githubv4"
)

// GraphQLClient is an interface for GitHub GraphQL client operations
type GraphQLClient interface {
	Query(ctx context.Context, q interface{}, variables map[string]interface{}) error
}

// isRetryableError determines if an error is retryable (5xx server errors)
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	errStr := err.Error()

	// Check for explicit 5xx status codes in error message
	// GitHub GraphQL API errors often include status codes in the error string
	if strings.Contains(errStr, "500") ||
		strings.Contains(errStr, "502") ||
		strings.Contains(errStr, "503") ||
		strings.Contains(errStr, "504") {
		return true
	}

	// Check for common server error messages
	if strings.Contains(strings.ToLower(errStr), "internal server error") ||
		strings.Contains(strings.ToLower(errStr), "bad gateway") ||
		strings.Contains(strings.ToLower(errStr), "service unavailable") ||
		strings.Contains(strings.ToLower(errStr), "gateway timeout") ||
		strings.Contains(strings.ToLower(errStr), "temporarily unavailable") {
		return true
	}

	return false
}

// isClientError determines if an error is a client error (4xx)
func isClientError(err error) bool {
	if err == nil {
		return false
	}

	errStr := err.Error()

	// Check for explicit 4xx status codes using regex
	// Match 3-digit numbers starting with 4 (400-499) with word boundaries
	re := regexp.MustCompile(`\b(4\d{2})\b`)
	if matches := re.FindStringSubmatch(errStr); len(matches) > 0 {
		if code, err := strconv.Atoi(matches[1]); err == nil && code >= 400 && code < 500 {
			return true
		}
	}

	// Check for common client error messages (case-insensitive)
	errStrLower := strings.ToLower(errStr)
	if strings.Contains(errStrLower, "could not resolve") ||
		strings.Contains(errStrLower, "not found") ||
		strings.Contains(errStrLower, "bad request") ||
		strings.Contains(errStrLower, "unauthorized") ||
		strings.Contains(errStrLower, "forbidden") {
		return true
	}

	return false
}

type ListBranchCheckStatusesReq struct {
	Owner         string
	Repo          string
	Branches      []string
	RetryWait     int // in minutes
	MaxRetry      int // in minutes
	MaxConcurrent int
}

type BranchCheckResult struct {
	Branch        string `json:"branch"`
	Status        string `json:"status"`
	FailureReason string `json:"failure_reason"`
	CommitHash    string `json:"commit_hash"`
	CommitUrl     string `json:"commit_url"`
}

type ListBranchCheckStatusRes struct {
	Results []BranchCheckResult `json:"results"`
}

func (r *ListBranchCheckStatusesReq) Run(ctx context.Context, client GraphQLClient) (*ListBranchCheckStatusRes, error) {
	if err := r.validate(); err != nil {
		return nil, fmt.Errorf("validating request: %w", err)
	}

	slog.Default().DebugContext(
		ctx, "checking branch statuses",
		"owner", r.Owner,
		"repo", r.Repo,
		"branches", r.Branches,
		"retry-wait", r.RetryWait,
		"max-retry", r.MaxRetry,
		"max-concurrent", r.MaxConcurrent,
	)

	// We add a semaphore control the # of concurrent requests to GH API
	semaphore := make(chan struct{}, r.MaxConcurrent)
	// We create a wait group to wait for all the requests to finish
	var wg sync.WaitGroup
	var mu sync.Mutex
	// Let's allocate a slice to store the results for the full set of branches
	results := make([]BranchCheckResult, len(r.Branches))

	maxRetryAttempts := 1
	if r.MaxRetry > 0 {
		maxRetryAttempts = r.MaxRetry
		if r.RetryWait >= 1 {
			maxRetryAttempts = r.MaxRetry / r.RetryWait
			if maxRetryAttempts == 0 {
				maxRetryAttempts = 1
			}
		}
	}

	for i, branch := range r.Branches {
		// Use the wait group to track our progress through the requests adding 1 for each
		wg.Add(1)
		// run our requests asynchronously ("hence the keyword go")
		go func(i int, branch string) {
			// Wait group to wait for the request to finish before decrementing
			defer wg.Done()
			semaphore <- struct{}{}        // acquire a semaphore token
			defer func() { <-semaphore }() // release the semaphore token

			slog.Default().DebugContext(ctx, "starting branch check", "branch", branch)

			// Send our branch check request using graphql client
			result := r.getBranchCheckResult(ctx, client, branch, maxRetryAttempts)

			// Log result
			if result.Status == "success" {
				shortHash := result.CommitHash
				if len(shortHash) > 8 {
					shortHash = shortHash[:8]
				}
				slog.Default().DebugContext(ctx, "branch check passed", "branch", branch, "commit", shortHash)
			} else {
				slog.Default().DebugContext(ctx, "branch check failed", "branch", branch, "reason", result.FailureReason)
			}

			// Adding our result to the results slice in a thread-safe manner
			mu.Lock()
			results[i] = result
			mu.Unlock()
		}(i, branch)
	}
	// Wait for all the requests to finish
	wg.Wait()

	slog.Default().DebugContext(ctx, "completed branch status checks", "count", len(results))

	return &ListBranchCheckStatusRes{Results: results}, nil
}

// This is the meat of this command the actual work of getting the branch check results
func (r *ListBranchCheckStatusesReq) getBranchCheckResult(ctx context.Context, client GraphQLClient, branch string, maxRetries int) BranchCheckResult {
	// Build out our graphql query
	var query struct {
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
	}

	variables := map[string]interface{}{
		"branch": githubv4.String(fmt.Sprintf("refs/heads/%s", branch)),
		"owner":  githubv4.String(r.Owner),
		"repo":   githubv4.String(r.Repo),
	}

	// Execute query for the branch at most MaxRetries times.
	// We track query attempts separately from status check retries
	const maxQueryRetries = 3
	var lastQueryErr error

	for requestAttempts := 0; requestAttempts < maxRetries; requestAttempts++ {
		// Try the query with retries for server errors
		for queryAttempt := 0; queryAttempt < maxQueryRetries; queryAttempt++ {
			err := client.Query(ctx, &query, variables)
			if err == nil {
				// Query succeeded, break out of retry loop
				lastQueryErr = nil
				break
			}

			lastQueryErr = err

			// If it's a client error (4xx), don't retry - fail immediately
			if isClientError(err) {
				slog.Default().DebugContext(ctx, "client error, not retrying", "branch", branch, "error", err)
				return BranchCheckResult{
					Branch:        branch,
					Status:        "failure",
					FailureReason: fmt.Sprintf("query-error (client): %s", err.Error()),
					CommitHash:    "unknown",
					CommitUrl:     "",
				}
			}

			// If it's a retryable server error (5xx), retry with backoff
			if isRetryableError(err) {
				if queryAttempt < maxQueryRetries-1 {
					// retry delay progression: 2s, 4s, 8s
					retryDelay := time.Duration(1<<uint(queryAttempt)) * 2 * time.Second
					slog.Default().DebugContext(
						ctx, "server error, retrying",
						"branch", branch,
						"retry-delay", retryDelay,
						"attempt", queryAttempt+1,
						"max-attempts", maxQueryRetries,
						"error", err,
					)
					time.Sleep(retryDelay)
					continue
				}
				// Exhausted retries for this query
				slog.Default().DebugContext(
					ctx, "server error persists after retries",
					"branch", branch,
					"retries", maxQueryRetries,
					"error", err,
				)
				return BranchCheckResult{
					Branch:        branch,
					Status:        "failure",
					FailureReason: fmt.Sprintf("query-error (server, retried %d times): %s", maxQueryRetries, err.Error()),
					CommitHash:    "unknown",
					CommitUrl:     "",
				}
			}

			// Unknown error type, don't retry
			slog.Default().DebugContext(ctx, "unknown error, not retrying", "branch", branch, "error", err)
			return BranchCheckResult{
				Branch:        branch,
				Status:        "failure",
				FailureReason: fmt.Sprintf("query-error: %s", err.Error()),
				CommitHash:    "unknown",
				CommitUrl:     "",
			}
		}

		// If we still have an error after retries, return failure
		if lastQueryErr != nil {
			return BranchCheckResult{
				Branch:        branch,
				Status:        "failure",
				FailureReason: fmt.Sprintf("query-error: %s", lastQueryErr.Error()),
				CommitHash:    "unknown",
				CommitUrl:     "",
			}
		}

		// We got back a valid response from query so let's process it
		commitHash := string(query.Repository.Ref.Target.Commit.Oid)
		status := string(query.Repository.Ref.Target.Commit.StatusCheckRollup.State)
		commitURL := fmt.Sprintf("https://github.com/%s/%s/commit/%s", r.Owner, r.Repo, commitHash)

		// Helper to safely get short commit hash
		shortHash := commitHash
		if len(commitHash) > 8 {
			shortHash = commitHash[:8]
		}

		switch status {
		case "SUCCESS":
			slog.Default().DebugContext(ctx, "checks passed", "branch", branch, "commit", shortHash)
			return BranchCheckResult{
				Branch:        branch,
				Status:        "success",
				FailureReason: "",
				CommitHash:    commitHash,
				CommitUrl:     commitURL,
			}
		case "FAILURE", "ERROR":
			slog.Default().DebugContext(ctx, "checks failed", "branch", branch, "commit", shortHash)
			return BranchCheckResult{
				Branch:        branch,
				Status:        "failure",
				FailureReason: "checks-failed",
				CommitHash:    commitHash,
				CommitUrl:     commitURL,
			}
		case "PENDING":
			// Not done yet, so wait and retry
			remainingAttempts := maxRetries - requestAttempts - 1
			slog.Default().DebugContext(
				ctx, "checks pending, retrying",
				"branch", branch,
				"commit", shortHash,
				"retry-wait", r.RetryWait,
				"attempt", requestAttempts+1,
				"max-attempts", maxRetries,
				"remaining", remainingAttempts,
			)
			time.Sleep(time.Duration(r.RetryWait) * time.Minute)
			continue
		default:
			// In case of an unexpected status
			slog.Default().DebugContext(ctx, "unexpected check status", "branch", branch, "status", status, "commit", shortHash)
			return BranchCheckResult{
				Branch:        branch,
				Status:        "failure",
				FailureReason: fmt.Sprintf("unexpected-status: %s", status),
				CommitHash:    commitHash,
				CommitUrl:     commitURL,
			}
		}
	}

	// All retries have been exhausted, so we have a time-out failure
	slog.Default().DebugContext(
		ctx, "checks timed out",
		"branch", branch,
		"attempts", maxRetries,
		"total-minutes", maxRetries*r.RetryWait,
	)
	return BranchCheckResult{
		Branch:        branch,
		Status:        "failure",
		FailureReason: "checks-timed-out",
		CommitHash:    "unknown",
		CommitUrl:     "",
	}
}

// validate ensures the request has the minimum required inputs and safe values.
func (r *ListBranchCheckStatusesReq) validate() error {
	if r == nil {
		return errors.New("failed to initialize request")
	}

	if r.Owner == "" {
		return errors.New("no github organization has been provided")
	}

	if r.Repo == "" {
		return errors.New("no github repository has been provided")
	}

	if len(r.Branches) == 0 {
		return errors.New("no branches have been provided")
	}

	if r.RetryWait < 0 {
		return errors.New("retry wait must be greater than or equal to 0")
	}

	if r.MaxRetry < 0 {
		return errors.New("max retry must be greater than or equal to 0")
	}

	if r.MaxConcurrent <= 0 {
		return errors.New("max concurrent must be greater than 0")
	}

	return nil
}

// ToJSON marshals the response to JSON.
func (r *ListBranchCheckStatusRes) ToJSON() ([]byte, error) {
	if r == nil {
		return nil, errors.New("uninitialized")
	}

	b, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("marshaling branch check statuses to JSON: %w", err)
	}

	return b, nil
}

// ToTable marshals the response to a text table.
func (r *ListBranchCheckStatusRes) ToTable() table.Writer {
	t := table.NewWriter()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false
	t.Style().Options.SeparateFooter = false
	t.Style().Options.SeparateHeader = false
	t.Style().Options.SeparateRows = false
	t.AppendHeader(table.Row{"branch", "status", "failure reason", "commit hash", "commit url"})

	for _, result := range r.Results {
		t.AppendRow(table.Row{
			result.Branch,
			result.Status,
			result.FailureReason,
			result.CommitHash,
			result.CommitUrl,
		})
	}

	t.SuppressEmptyColumns()
	t.SuppressTrailingSpaces()

	return t
}
