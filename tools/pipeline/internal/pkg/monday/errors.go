// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package monday

import (
	"errors"
	"net/http"
	"strings"

	graphql "github.com/hasura/go-graphql-client"
)

// isRetryableError reports whether err is a transient error that warrants a
// retry. Retryable conditions:
//   - HTTP 429 (rate-limited), 500, 502, 503, 504
//   - Monday.com GraphQL complexity-budget exhaustion ("complexity" in the message)
//
// Use this policy only for idempotent read queries. For write mutations that
// create data, use isCreateRetryableError instead to avoid double-creation.
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	// Network-level HTTP error (4xx / 5xx from the transport).
	var netErr graphql.NetworkError
	if errors.As(err, &netErr) {
		switch netErr.StatusCode() {
		case http.StatusTooManyRequests,
			http.StatusInternalServerError,
			http.StatusBadGateway,
			http.StatusServiceUnavailable,
			http.StatusGatewayTimeout:
			return true
		}
		return false
	}

	// GraphQL application error (HTTP 200 + errors array).
	var gqlErrs graphql.Errors
	if errors.As(err, &gqlErrs) {
		for _, e := range gqlErrs {
			if strings.Contains(strings.ToLower(e.Message), "complexity") {
				return true
			}
		}
		return false
	}

	return false
}

// isCreateRetryableError reports whether err warrants retrying a non-idempotent
// create mutation. The policy is intentionally narrower than isRetryableError:
// 5xx gateway codes (500, 502, 503, 504) are excluded because the server may
// have committed the item before the gateway returned the error, so retrying
// would create a duplicate. Only pre-execution failures that guarantee no data
// was written are safe to retry:
//   - HTTP 429 (rate-limited): Monday.com rejects the request before execution.
//   - GraphQL complexity budget exhausted: the API returns an errors array before
//     running the mutation, so no item is created.
func isCreateRetryableError(err error) bool {
	if err == nil {
		return false
	}

	// HTTP 429: Monday.com rate-limits the request before executing it.
	var netErr graphql.NetworkError
	if errors.As(err, &netErr) {
		return netErr.StatusCode() == http.StatusTooManyRequests
	}

	// GraphQL complexity error: returned before the mutation runs.
	var gqlErrs graphql.Errors
	if errors.As(err, &gqlErrs) {
		for _, e := range gqlErrs {
			if strings.Contains(strings.ToLower(e.Message), "complexity") {
				return true
			}
		}
		return false
	}

	return false
}
