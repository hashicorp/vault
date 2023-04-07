// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"net/http"
	"net/url"
	"testing"
)

func TestBuildSamplePolicy(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		req      *OutputPolicyError
		expected string
		err      error
	}{
		{
			"happy path",
			&OutputPolicyError{
				method: http.MethodGet,
				path:   "/something",
			},
			formatOutputPolicy("/something", []string{"read"}),
			nil,
		},
		{ // test included to clear up some confusion around the sanitize comment
			"demonstrate that this function does not format fully",
			&OutputPolicyError{
				method: http.MethodGet,
				path:   "http://vault.test/v1/something",
			},
			formatOutputPolicy("http://vault.test/v1/something", []string{"read"}),
			nil,
		},
		{ // test that list is properly returned
			"list over read returned",
			&OutputPolicyError{
				method: http.MethodGet,
				path:   "/something",
				params: url.Values{
					"list": []string{"true"},
				},
			},
			formatOutputPolicy("/something", []string{"list"}),
			nil,
		},
		{
			"valid protected path",
			&OutputPolicyError{
				method: http.MethodGet,
				path:   "/sys/config/ui/headers/",
			},
			formatOutputPolicy("/sys/config/ui/headers/", []string{"read", "sudo"}),
			nil,
		},
		{ // ensure that a formatted path that trims the trailing slash as the code does still works for recognizing a sudo path
			"valid protected path no trailing /",
			&OutputPolicyError{
				method: http.MethodGet,
				path:   "/sys/config/ui/headers",
			},
			formatOutputPolicy("/sys/config/ui/headers", []string{"read", "sudo"}),
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := tc.req.buildSamplePolicy()
			if tc.err != err {
				t.Fatalf("expected for the error to be %v instead got %v\n", tc.err, err)
			}

			if tc.expected != result {
				t.Fatalf("expected for the policy string to be %v instead got %v\n", tc.expected, result)
			}
		})
	}
}
