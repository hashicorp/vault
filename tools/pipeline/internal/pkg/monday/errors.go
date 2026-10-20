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
