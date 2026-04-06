// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package kv

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/testhelpers/minimal"
	"github.com/stretchr/testify/require"
)

// rawVaultRequest sends an HTTP request to Vault using the exact rawPath
// provided, bypassing the Go API client's path normalization. This is necessary
// for testing URL-encoded path traversal sequences like %2e%2e that would
// otherwise be resolved by path.Join in the standard client.
func rawVaultRequest(t *testing.T, c *api.Client, token, method, rawPath, body string) *http.Response {
	t.Helper()

	baseAddr := c.Address()
	fullURL := baseAddr + rawPath

	// Parse the URL and set Opaque to prevent Go from decoding %2e%2e
	parsed, err := url.Parse(fullURL)
	require.NoError(t, err, "failed to parse URL %q", fullURL)
	parsed.Opaque = fmt.Sprintf("//%s%s", parsed.Host, rawPath)

	var reqBody *strings.Reader
	if body != "" {
		reqBody = strings.NewReader(body)
	}

	var req *http.Request
	if reqBody != nil {
		req, err = http.NewRequest(method, parsed.String(), reqBody)
	} else {
		req, err = http.NewRequest(method, parsed.String(), nil)
	}
	require.NoError(t, err, "failed to create request")

	req.Header.Set("X-Vault-Token", token)
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}

	// Use the same HTTP client as the Vault API client to inherit TLS config
	httpClient := c.CloneConfig().HttpClient
	httpClient.CheckRedirect = nil
	resp, err := httpClient.Do(req)
	require.NoError(t, err, "failed to execute request %s %s", method, rawPath)

	return resp
}

// TestKV_PathTraversal tests a variety of paths to ensure that path traversals
// are not allowed and return the proper error code.
func TestKV_PathTraversal(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	c := cluster.Cores[0].Client

	// Mount a KVv2 backend
	err := c.Sys().Mount("kv-v2", &api.MountInput{
		Type: "kv-v2",
	})
	require.NoError(t, err)

	// Write a protected secret
	kvData := map[string]interface{}{
		"data": map[string]interface{}{
			"password": "THIS IS A SECRET",
		},
	}

	_, err = kvRequestWithRetry(t, func() (interface{}, error) {
		return c.Logical().Write("kv-v2/data/team/protected/dbcreds", kvData)
	})
	require.NoError(t, err)
	// Create a restrictive policy that only allows access to team/public/* paths
	err = c.Sys().PutPolicy("public-policy", `
		path "kv-v2/metadata/team/public/*" { capabilities = ["read", "list"] }
		path "kv-v2/data/team/public/*" { capabilities = ["read", "list"] }
		path "kv-v2/destroy/team/public/*" { capabilities = ["update"] }
		path "kv-v2/team/public/*" { capabilities = ["read"] }
	`)
	require.NoError(t, err)

	// Create an attacker token with only the public-policy and no default policy
	attackerSecret, err := c.Auth().Token().Create(&api.TokenCreateRequest{
		Policies:        []string{"public-policy"},
		NoDefaultPolicy: true,
	})
	require.NoError(t, err)
	attackerToken := attackerSecret.Auth.ClientToken

	// Verify the attacker cannot read protected metadata via canonical path
	attackerClient, err := api.NewClient(c.CloneConfig())
	require.NoError(t, err)
	attackerClient.SetToken(attackerToken)

	testCases := []struct {
		name       string
		path       string
		operation  string
		body       []byte
		expectCode int
	}{
		{
			name:       "read secret",
			path:       "kv-v2/data/team/public/../protected/dbcreds",
			operation:  "GET",
			expectCode: 400,
		},
		{
			name:       "read encoded",
			path:       "kv-v2/data/team/public/%2e%2e/protected/dbcreds",
			operation:  "GET",
			expectCode: 400,
		},
		{
			name:       "destroy",
			path:       "kv-v2/destroy/team/public/../protected/dbcreds",
			body:       []byte(`{"versions":[1]}`),
			operation:  "PUT",
			expectCode: 400,
		},
		{
			name:       "destroy encoded",
			path:       "kv-v2/destroy/team/public/%2e%2e/protected/dbcreds",
			body:       []byte(`{"versions":[1]}`),
			operation:  "PUT",
			expectCode: 400,
		},
		{
			name:       "metadata read",
			path:       "kv-v2/metadata/team/public/../protected/dbcreds",
			operation:  "GET",
			expectCode: 400,
		},
		{
			name:       "metadata read encoded",
			path:       "kv-v2/metadata/team/public/%2e%2e/protected/dbcreds",
			operation:  "GET",
			expectCode: 400,
		},
		{
			name:       "metadata read double encoded",
			path:       "kv-v2/metadata/team/public/%252e%252e/protected/dbcreds",
			operation:  "GET",
			expectCode: 404,
		},
		{
			name:       "metadata read double slash",
			path:       "kv-v2/metadata/team/public////protected/dbcreds",
			operation:  "GET",
			expectCode: 400,
		},
		{
			name:       "metadata read double slash encoded",
			path:       "kv-v2/metadata/team/public/%2F%2F/protected/dbcreds",
			operation:  "GET",
			expectCode: 400,
		},
		{
			name:       "metadata read empty path piece",
			path:       "kv-v2/metadata/team/public//protected/dbcreds",
			operation:  "GET",
			expectCode: 400,
		},
		{
			name:       "ending slash",
			path:       "kv-v2/metadata/team/public/dbcreds/",
			operation:  "GET",
			expectCode: 404,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rawResp := rawVaultRequest(t, c, attackerToken, tc.operation,
				"/v1/"+tc.path, string(tc.body))
			defer rawResp.Body.Close()
			require.Equal(t, tc.expectCode, rawResp.StatusCode)
		})
	}
}
