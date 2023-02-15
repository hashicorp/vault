package api

import (
	"net/http"
	"net/url"
	"testing"
)

func TestBuildSamplePolicy(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		req      *OutputPolicyError
		expected string
		err      error
	}{
		{
			&OutputPolicyError{
				method: http.MethodGet,
				path:   "/something",
			},
			fmtOutput("/something", []string{"read"}),
			nil,
		},
		{ // test included to clear up some confusion around the sanitize comment
			&OutputPolicyError{
				method: http.MethodGet,
				path:   "http://vault.test/v1/something",
			},
			fmtOutput("http://vault.test/v1/something", []string{"read"}),
			nil,
		},
		{ // test that list is properly returned
			&OutputPolicyError{
				method: http.MethodGet,
				path:   "/something",
				params: url.Values{
					listKey: []string{"true"},
				},
			},
			fmtOutput("/something", []string{"list"}),
			nil,
		},
		{
			&OutputPolicyError{
				method: http.MethodGet,
				path:   "/sys/config/ui/headers/",
			},
			fmtOutput("/sys/config/ui/headers/", []string{"read", "sudo"}),
			nil,
		},
		{ // ensure that a formatted path that trims the trailing slash as the code does still works for recognizing a sudo path
			&OutputPolicyError{
				method: http.MethodGet,
				path:   "/sys/config/ui/headers",
			},
			fmtOutput("/sys/config/ui/headers", []string{"read", "sudo"}),
			nil,
		},
	}

	for _, tc := range testCases {
		result, err := tc.req.buildSamplePolicy()
		if tc.err != err {
			t.Fatalf("expected for the error to be %v instead got %v\n", tc.err, err)
		}

		if tc.expected != result {
			t.Fatalf("expected for the policy string to be %v instead got %v\n", tc.expected, result)
		}
	}
}
