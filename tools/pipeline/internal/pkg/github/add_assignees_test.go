// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	libgithub "github.com/google/go-github/v83/github"
	"github.com/stretchr/testify/require"
)

// Test_addAssignees tests the addAssignees helper function with various input scenarios
// including single/multiple assignees, empty lists, filtering of empty strings, and
// deduplication of logins.
func Test_addAssignees(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		logins        []string
		shouldCall    bool
		expectedError bool
	}{
		"single assignee": {
			logins:     []string{"user1"},
			shouldCall: true,
		},
		"multiple assignees": {
			logins:     []string{"user1", "user2", "user3"},
			shouldCall: true,
		},
		"empty login list": {
			logins:     []string{},
			shouldCall: false,
		},
		"nil login list": {
			logins:     nil,
			shouldCall: false,
		},
		"logins with empty strings": {
			logins:     []string{"user1", "", "user2", ""},
			shouldCall: true, // Empty strings should be filtered out
		},
		"only empty strings": {
			logins:     []string{"", "", ""},
			shouldCall: false, // All empty, should skip
		},
		"duplicate logins": {
			logins:     []string{"user1", "user1", "user2"},
			shouldCall: true, // Duplicates should be compacted
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			called := false
			client, mux, teardown := setupTestClient(t)
			defer teardown()

			if test.shouldCall {
				mux.HandleFunc("/repos/test-owner/test-repo/issues/123/assignees", func(w http.ResponseWriter, r *http.Request) {
					require.Equal(t, http.MethodPost, r.Method)
					called = true
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte(`{"assignees": []}`))
				})
			}

			err := addAssignees(
				context.Background(),
				client,
				"test-owner",
				"test-repo",
				123,
				test.logins,
			)

			if test.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, test.shouldCall, called, "API call expectation mismatch")
		})
	}
}

// Test_addReviewers tests the addReviewers helper function with various input scenarios
// including single/multiple reviewers, empty lists, filtering of empty strings, and
// deduplication of logins.
func Test_addReviewers(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		logins        []string
		shouldCall    bool
		expectedError bool
	}{
		"single reviewer": {
			logins:     []string{"user1"},
			shouldCall: true,
		},
		"multiple reviewers": {
			logins:     []string{"user1", "user2", "user3"},
			shouldCall: true,
		},
		"empty login list": {
			logins:     []string{},
			shouldCall: false,
		},
		"nil login list": {
			logins:     nil,
			shouldCall: false,
		},
		"logins with empty strings": {
			logins:     []string{"user1", "", "user2", ""},
			shouldCall: true, // Empty strings should be filtered out
		},
		"only empty strings": {
			logins:     []string{"", "", ""},
			shouldCall: false, // All empty, should skip
		},
		"duplicate logins": {
			logins:     []string{"user1", "user1", "user2"},
			shouldCall: true, // Duplicates should be compacted
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			called := false
			client, mux, teardown := setupTestClient(t)
			defer teardown()

			if test.shouldCall {
				mux.HandleFunc("/repos/test-owner/test-repo/pulls/123/requested_reviewers", func(w http.ResponseWriter, r *http.Request) {
					require.Equal(t, http.MethodPost, r.Method)
					called = true
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte(`{"users": [], "teams": []}`))
				})
			}

			err := addReviewers(
				context.Background(),
				client,
				"test-owner",
				"test-repo",
				123,
				test.logins,
			)

			if test.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, test.shouldCall, called, "API call expectation mismatch")
		})
	}
}

// setupTestClient creates a test GitHub client with a mock HTTP server for testing.
// It returns the client, the HTTP mux for registering handlers, and a teardown function
// that should be called to clean up the server when the test completes.
func setupTestClient(t *testing.T) (*libgithub.Client, *http.ServeMux, func()) {
	t.Helper()

	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	client := libgithub.NewClient(nil)
	serverURL, err := url.Parse(server.URL + "/")
	require.NoError(t, err)
	client.BaseURL = serverURL

	teardown := func() {
		server.Close()
	}

	return client, mux, teardown
}
