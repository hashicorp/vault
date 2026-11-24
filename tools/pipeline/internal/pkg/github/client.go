// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"fmt"
	"net/http"
	"net/url"

	gh "github.com/google/go-github/v83/github"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

// NewClient constructs an authenticated *gh.Client for the given host and
// static bearer token.
//
// When host is "github.com" the standard github.com endpoints are used.
// Any other value is treated as a GitHub Enterprise Server host and
// WithEnterpriseURLs is applied automatically.
//
// base is the underlying HTTP transport. Pass nil to use
// http.DefaultTransport. Pass a custom transport in tests to avoid
// touching the network (e.g. srv.Client().Transport from an httptest.Server).
func NewClient(host, token string, base http.RoundTripper) (*gh.Client, error) {
	httpClient := oauth2HTTPClient(token, base)

	if host == "github.com" || host == "" {
		return gh.NewClient(httpClient), nil
	}

	apiURL, err := url.JoinPath("https://", host, "/api/v3/")
	if err != nil {
		return nil, fmt.Errorf("determining GitHub Enterprise API URL for %s: %w", host, err)
	}
	uploadURL, err := url.JoinPath("https://", host, "/api/uploads/")
	if err != nil {
		return nil, fmt.Errorf("determining GitHub Enterprise upload URL for %s: %w", host, err)
	}

	client, err := gh.NewClient(httpClient).WithEnterpriseURLs(apiURL, uploadURL)
	if err != nil {
		return nil, fmt.Errorf("configuring GitHub Enterprise URLs for %s: %w", host, err)
	}
	return client, nil
}

// NewV4Client constructs an authenticated *githubv4.Client for the given host
// and static bearer token using the same oauth2.Transport as NewClient.
//
// base follows the same nil-means-DefaultTransport contract as NewClient.
func NewV4Client(host, token string, base http.RoundTripper) (*githubv4.Client, error) {
	httpClient := oauth2HTTPClient(token, base)

	if host == "github.com" || host == "" {
		return githubv4.NewClient(httpClient), nil
	}

	graphqlURL, err := url.JoinPath("https://", host, "/api/graphql")
	if err != nil {
		return nil, fmt.Errorf("determining GitHub Enterprise GraphQL URL for %s: %w", host, err)
	}
	return githubv4.NewEnterpriseClient(graphqlURL, httpClient), nil
}

// oauth2HTTPClient returns an *http.Client whose transport injects
// "Authorization: Bearer <token>" on every request via oauth2.Transport.
// When token is empty the transport is still constructed but will send an
// empty Bearer header; callers that need unauthenticated access should
// pass a nil *http.Client to the upstream constructors directly.
func oauth2HTTPClient(token string, base http.RoundTripper) *http.Client {
	return &http.Client{
		Transport: &oauth2.Transport{
			Source: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}),
			Base:   base, // nil is safe: oauth2.Transport falls back to http.DefaultTransport
		},
	}
}
