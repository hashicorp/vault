// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"regexp"
	"strconv"
	"strings"
)

// isRetryableError determines if an error is retryable (5xx server errors and rate-limit 429).
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	errStr := err.Error()

	// The shurcooL/graphql client formats non-200 responses as:
	//   "non-200 OK status code: <Status> body: <body>"
	// where <Status> is Go's http.Response.Status, e.g. "503 Service Unavailable".
	// Match the three-digit codes we want to retry on.
	re := regexp.MustCompile(`\b(5\d{2}|429)\b`)
	if matches := re.FindStringSubmatch(errStr); len(matches) > 0 {
		if code, parseErr := strconv.Atoi(matches[1]); parseErr == nil {
			switch code {
			case 429, 500, 502, 503, 504:
				return true
			}
			// if we don't have a match on the retryable statuses, then just continue on to the message text
		}
	}

	// Also match common textual forms that appear in GitHub GraphQL error payloads
	// (HTTP 200 response body with an errors array) or network-level messages.
	errStrLower := strings.ToLower(errStr)
	if strings.Contains(errStrLower, "internal server error") ||
		strings.Contains(errStrLower, "bad gateway") ||
		strings.Contains(errStrLower, "service unavailable") ||
		strings.Contains(errStrLower, "gateway timeout") ||
		strings.Contains(errStrLower, "temporarily unavailable") ||
		strings.Contains(errStrLower, "rate limit") ||
		strings.Contains(errStrLower, "too many requests") {
		return true
	}

	return false
}

// isClientError determines if an error is a non-retryable client error (4xx, excluding 429).
// 429 Too Many Requests is treated as retryable by isRetryableError.
func isClientError(err error) bool {
	if err == nil {
		return false
	}

	errStr := err.Error()

	// Match explicit 4xx codes, but skip 429 (handled as retryable).
	re := regexp.MustCompile(`\b(4\d{2})\b`)
	if matches := re.FindStringSubmatch(errStr); len(matches) > 0 {
		if code, parseErr := strconv.Atoi(matches[1]); parseErr == nil && code >= 400 && code < 500 && code != 429 {
			return true
		}
	}

	// Common textual client error phrases from GitHub GraphQL error payloads.
	// Note: "not found" / "could not resolve" indicate a bad request, not a
	// missing branch — branch absence is detected via empty Oid, not an error.
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
