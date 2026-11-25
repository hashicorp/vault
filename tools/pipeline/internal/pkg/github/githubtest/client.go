// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

// Package githubtest has helpers for testing code that talks to the Github API.
//
// NOTE: It's a separate package so that only our tests import it. That way we don't
// build testing or httptest into the pipeline binary.
package githubtest

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	gh "github.com/google/go-github/v83/github"
	"github.com/stretchr/testify/require"
)

// SetupTestClient creates a Github client that talks to a mock HTTP server. It
// returns the client, the mux to register handlers on, and a teardown func that
// ought to be called to shut down the server when the test completes.
func SetupTestClient(t *testing.T) (*gh.Client, *http.ServeMux, func()) {
	t.Helper()

	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	client := gh.NewClient(nil)
	serverURL, err := url.Parse(server.URL + "/")
	require.NoError(t, err)
	client.BaseURL = serverURL

	teardown := func() {
		server.Close()
	}

	return client, mux, teardown
}
