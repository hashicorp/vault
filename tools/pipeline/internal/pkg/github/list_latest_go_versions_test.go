// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

// newGoDevServer builds an httptest server that mimics the go.dev release API.
func newGoDevServer(t *testing.T, releases []map[string]any) *httptest.Server {
	t.Helper()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(releases)
	}))
	t.Cleanup(srv.Close)

	return srv
}

// redirectTransport redirects all requests to a test server URL, preserving
// the path and query so the handler can inspect them normally.
type redirectTransport struct {
	base   http.RoundTripper
	target string // full URL including scheme, e.g. "http://127.0.0.1:PORT"
}

func (rt *redirectTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(req.Context())
	req.URL.Scheme = "http"
	req.URL.Host = rt.target[len("http://"):]
	return rt.base.RoundTrip(req)
}

// newGoDevClient returns an *http.Client whose transport redirects all requests
// to srv, making tests parallel-safe (no global transport mutation).
func newGoDevClient(srv *httptest.Server) *http.Client {
	return &http.Client{
		Transport: &redirectTransport{base: http.DefaultTransport, target: srv.URL},
	}
}

// Test_ListLatestGoVersionsReq_Run verifies that the latest and previous stable
// versions are correctly extracted from the go.dev response, and that the
// "go" prefix is stripped from version strings.
func Test_ListLatestGoVersionsReq_Run(t *testing.T) {
	t.Parallel()

	releases := []map[string]any{
		{"version": "go1.23.4", "stable": true},
		{"version": "go1.22.9", "stable": true},
		{"version": "go1.23.4rc1", "stable": false}, // non-stable, must be skipped
	}

	srv := newGoDevServer(t, releases)

	req := &ListLatestGoVersionsReq{HTTPClient: newGoDevClient(srv)}
	res, err := req.Run(context.Background())
	require.NoError(t, err)

	require.Equal(t, "1.23.4", res.LatestStableVersion)
	require.Equal(t, "1.22.9", res.PreviousStableVersion)
}

// Test_ListLatestGoVersionsReq_Run_InsufficientReleases verifies an error is
// returned when fewer than two stable releases are available.
func Test_ListLatestGoVersionsReq_Run_InsufficientReleases(t *testing.T) {
	t.Parallel()

	releases := []map[string]any{
		{"version": "go1.23.4", "stable": true},
	}

	srv := newGoDevServer(t, releases)

	req := &ListLatestGoVersionsReq{HTTPClient: newGoDevClient(srv)}
	_, err := req.Run(context.Background())
	require.Error(t, err)
	require.Contains(t, err.Error(), "at least two stable")
}

// Test_ListLatestGoVersionsRes_ToGithubOutput verifies the output key names
// match what the workflow's get-latest-go-version step expects.
func Test_ListLatestGoVersionsRes_ToGithubOutput(t *testing.T) {
	t.Parallel()

	res := &ListLatestGoVersionsRes{
		LatestStableVersion:   "1.23.4",
		PreviousStableVersion: "1.22.9",
	}

	out, err := res.ToGithubOutput()
	require.NoError(t, err)

	require.Equal(t, map[string]string{
		"latest_stable_version":   "1.23.4",
		"previous_stable_version": "1.22.9",
	}, out)
}
