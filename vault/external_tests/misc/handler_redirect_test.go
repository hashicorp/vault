// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package misc

import (
	"net/http"
	"testing"

	"github.com/hashicorp/vault/helper/testhelpers/minimal"
)

// TestHandler_RedirectPreservesQueryParams verifies that double-slash paths in
// redirect responses preserve query parameters in the Location header.
func TestHandler_RedirectPreservesQueryParams(t *testing.T) {
	cases := []struct {
		name         string
		inputPath    string
		wantLocation string
	}{
		{
			name:         "kv double slash preserves multiple query params",
			inputPath:    "/v1/secret//foo?version=1&list=true",
			wantLocation: "/v1/secret/foo?version=1&list=true",
		},
		{
			name:         "sys policies double slash preserves list param",
			inputPath:    "/v1/sys//policies/acl?list=true",
			wantLocation: "/v1/sys/policies/acl?list=true",
		},
		{
			name:         "auth token accessors double slash preserves list param",
			inputPath:    "/v1/auth//token/accessors?list=true",
			wantLocation: "/v1/auth/token/accessors?list=true",
		},
		{
			name:         "deep path with double slash preserves query params",
			inputPath:    "/v1/secret/data//nested/key?version=2",
			wantLocation: "/v1/secret/data/nested/key?version=2",
		},
		{
			name:         "double slash with no query params",
			inputPath:    "/v1/secret//foo",
			wantLocation: "/v1/secret/foo",
		},
	}

	cluster := minimal.NewTestSoloCluster(t, nil)
	apiClient := cluster.Cores[0].Client

	httpClient := apiClient.CloneConfig().HttpClient
	httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	addr := apiClient.Address()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, addr+tc.inputPath, nil)
			if err != nil {
				t.Fatalf("failed to build request: %v", err)
			}

			resp, err := httpClient.Do(req)
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}
			if resp.StatusCode != http.StatusTemporaryRedirect {
				t.Fatalf("expected %d, got %d", http.StatusTemporaryRedirect, resp.StatusCode)
			}
			if got := resp.Header.Get("Location"); got != tc.wantLocation {
				t.Fatalf("expected Location %q, got %q", tc.wantLocation, got)
			}
		})
	}
}
