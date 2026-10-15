// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package monday

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	graphql "github.com/hasura/go-graphql-client"
	"github.com/stretchr/testify/require"
)

// TestIsRetryableError verifies that isRetryableError correctly classifies
// network-level HTTP status codes and GraphQL application errors as retryable
// or non-retryable. Retryable conditions are rate limits (429), server errors
// (500/502/503/504), and Monday.com complexity-budget exhaustion. Auth errors,
// bad requests, and unknown errors are non-retryable.
func TestIsRetryableError(t *testing.T) {
	t.Parallel()

	// networkErrFor returns an error produced by the hasura client against a
	// minimal httptest server that replies with statusCode, exercising the
	// graphql.NetworkError path without requiring a real Monday.com token.
	networkErrFor := func(t *testing.T, statusCode int) error {
		t.Helper()
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(statusCode)
		}))
		t.Cleanup(srv.Close)
		c := graphql.NewClient(srv.URL, srv.Client())
		var q struct{ Dummy string }
		return c.Query(context.Background(), &q, nil)
	}

	for name, tc := range map[string]struct {
		errFn    func(t *testing.T) error
		wantTrue bool
	}{
		"nil":           {errFn: func(*testing.T) error { return nil }, wantTrue: false},
		"generic error": {errFn: func(*testing.T) error { return errors.New("something") }, wantTrue: false},
		"network 429":   {errFn: func(t *testing.T) error { return networkErrFor(t, 429) }, wantTrue: true},
		"network 500":   {errFn: func(t *testing.T) error { return networkErrFor(t, 500) }, wantTrue: true},
		"network 502":   {errFn: func(t *testing.T) error { return networkErrFor(t, 502) }, wantTrue: true},
		"network 503":   {errFn: func(t *testing.T) error { return networkErrFor(t, 503) }, wantTrue: true},
		"network 504":   {errFn: func(t *testing.T) error { return networkErrFor(t, 504) }, wantTrue: true},
		"network 400":   {errFn: func(t *testing.T) error { return networkErrFor(t, 400) }, wantTrue: false},
		"network 401":   {errFn: func(t *testing.T) error { return networkErrFor(t, 401) }, wantTrue: false},
		"network 403":   {errFn: func(t *testing.T) error { return networkErrFor(t, 403) }, wantTrue: false},
		"gql complexity budget exhausted": {errFn: func(*testing.T) error {
			return graphql.Errors{{Message: "Complexity budget exhausted, retry after 10 seconds"}}
		}, wantTrue: true},
		"gql complexity budget exceeded": {errFn: func(*testing.T) error {
			return graphql.Errors{{Message: "complexity budget exceeded"}}
		}, wantTrue: true},
		"gql not authenticated": {errFn: func(*testing.T) error {
			return graphql.Errors{{Message: "Not authenticated"}}
		}, wantTrue: false},
		"gql item not found": {errFn: func(*testing.T) error {
			return graphql.Errors{{Message: "Item not found"}}
		}, wantTrue: false},
		"gql user unauthorized": {errFn: func(*testing.T) error {
			return graphql.Errors{{Message: "User unauthorized to perform action", Extensions: map[string]any{"code": "UserUnauthorizedException", "status_code": 403}}}
		}, wantTrue: false},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := isRetryableError(tc.errFn(t))
			if tc.wantTrue {
				require.True(t, got, "expected retryable")
			} else {
				require.False(t, got, "expected non-retryable")
			}
		})
	}
}
