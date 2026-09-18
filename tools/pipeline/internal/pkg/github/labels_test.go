// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// Test_addLabelsToIssue tests that addLabelsToIssue calls the (mocked) GitHub API
// with the expected labels, and skips the call when no labels are provided.
func Test_addLabelsToIssue(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		labels     []string
		shouldCall bool
	}{
		"single label": {
			labels:     []string{"backport-failed"},
			shouldCall: true,
		},
		"multiple labels": {
			labels:     []string{"backport-failed", "bug"},
			shouldCall: true,
		},
		"empty label list": {
			labels:     []string{},
			shouldCall: false,
		},
		"nil label list": {
			labels:     nil,
			shouldCall: false,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			called := false
			var receivedLabels []string
			client, mux, teardown := setupTestClient(t)
			defer teardown()

			if test.shouldCall {
				mux.HandleFunc("/repos/test-owner/test-repo/issues/42/labels", func(w http.ResponseWriter, r *http.Request) {
					require.Equal(t, http.MethodPost, r.Method)
					called = true

					var body []string
					require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
					receivedLabels = body

					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`[]`))
				})
			}

			err := addLabelsToIssue(
				context.Background(),
				client,
				"test-owner",
				"test-repo",
				42,
				test.labels,
			)

			require.NoError(t, err)
			require.Equal(t, test.shouldCall, called, "API call expectation mismatch")
			if test.shouldCall {
				require.Equal(t, test.labels, receivedLabels, "labels sent to API do not match")
			}
		})
	}
}

// Test_removeLabelFromIssue tests that removeLabelFromIssue calls the (mocked)
// GitHub API with the expected label, errors when no label provided, and
// ignores a 404 response (label not present on original PR).
func Test_removeLabelFromIssue(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		label          string
		responseStatus int
		shouldCall     bool
		expectError    bool
	}{
		"label present": {
			label:          "backport-failed",
			responseStatus: http.StatusOK,
			shouldCall:     true,
		},
		"label not present - 404 ignored": {
			label:          "backport-failed",
			responseStatus: http.StatusNotFound,
			shouldCall:     true,
		},
		"empty label - error": {
			label:       "",
			shouldCall:  false,
			expectError: true,
		},
		"server error": {
			label:          "backport-failed",
			responseStatus: http.StatusInternalServerError,
			shouldCall:     true,
			expectError:    true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			called := false
			client, mux, teardown := setupTestClient(t)
			defer teardown()

			if test.shouldCall {
				mux.HandleFunc("/repos/test-owner/test-repo/issues/42/labels/backport-failed", func(w http.ResponseWriter, r *http.Request) {
					require.Equal(t, http.MethodDelete, r.Method)
					called = true
					w.WriteHeader(test.responseStatus)
				})
			}

			err := removeLabelFromIssue(
				context.Background(),
				client,
				"test-owner",
				"test-repo",
				42,
				test.label,
			)

			if test.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, test.shouldCall, called, "API call expectation mismatch")
		})
	}
}
