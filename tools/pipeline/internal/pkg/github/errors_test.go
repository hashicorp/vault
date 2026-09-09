// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"errors"
	"testing"
)

// TestIsRetryableError verifies that isRetryableError correctly identifies
// retryable errors (5xx server errors and 429 rate limits).
func TestIsRetryableError(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		err      error
		expected bool
	}{
		"nil error": {
			err:      nil,
			expected: false,
		},
		"500 internal server error": {
			err:      errors.New("non-200 OK status code: 500 Internal Server Error body: something"),
			expected: true,
		},
		"502 bad gateway": {
			err:      errors.New("non-200 OK status code: 502 Bad Gateway body: something"),
			expected: true,
		},
		"503 service unavailable": {
			err:      errors.New("non-200 OK status code: 503 Service Unavailable body: something"),
			expected: true,
		},
		"504 gateway timeout": {
			err:      errors.New("non-200 OK status code: 504 Gateway Timeout body: something"),
			expected: true,
		},
		"429 too many requests (numeric)": {
			err:      errors.New("non-200 OK status code: 429 Too Many Requests body: something"),
			expected: true,
		},
		"internal server error text": {
			err:      errors.New("internal server error occurred"),
			expected: true,
		},
		"temporarily unavailable": {
			err:      errors.New("service temporarily unavailable"),
			expected: true,
		},
		"rate limit text": {
			err:      errors.New("rate limit exceeded"),
			expected: true,
		},
		"too many requests text": {
			err:      errors.New("too many requests"),
			expected: true,
		},
		"404 not found": {
			err:      errors.New("non-200 OK status code: 404 Not Found body: something"),
			expected: false,
		},
		"400 bad request": {
			err:      errors.New("non-200 OK status code: 400 Bad Request body: something"),
			expected: false,
		},
		"generic error": {
			err:      errors.New("something went wrong"),
			expected: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result := isRetryableError(tc.err)
			if result != tc.expected {
				t.Errorf("isRetryableError(%v) = %v, expected %v", tc.err, result, tc.expected)
			}
		})
	}
}

// TestIsClientError verifies that isClientError correctly identifies
// non-retryable client errors (4xx, excluding 429).
func TestIsClientError(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		err      error
		expected bool
	}{
		"nil error": {
			err:      nil,
			expected: false,
		},
		"400 bad request": {
			err:      errors.New("non-200 OK status code: 400 Bad Request body: something"),
			expected: true,
		},
		"401 unauthorized": {
			err:      errors.New("non-200 OK status code: 401 Unauthorized body: something"),
			expected: true,
		},
		"403 forbidden": {
			err:      errors.New("non-200 OK status code: 403 Forbidden body: something"),
			expected: true,
		},
		"404 not found": {
			err:      errors.New("non-200 OK status code: 404 Not Found body: something"),
			expected: true,
		},
		"429 not a client error (retryable)": {
			err:      errors.New("non-200 OK status code: 429 Too Many Requests body: something"),
			expected: false,
		},
		"could not resolve": {
			err:      errors.New("could not resolve to a Repository"),
			expected: true,
		},
		"not found text": {
			err:      errors.New("not found"),
			expected: true,
		},
		"500 server error": {
			err:      errors.New("non-200 OK status code: 500 Internal Server Error body: something"),
			expected: false,
		},
		"generic error": {
			err:      errors.New("something went wrong"),
			expected: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result := isClientError(tc.err)
			if result != tc.expected {
				t.Errorf("isClientError(%v) = %v, expected %v", tc.err, result, tc.expected)
			}
		})
	}
}
