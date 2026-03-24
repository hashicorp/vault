// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package api

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/testhelpers/minimal"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/require"
)

// TestPassthroughHeaders_Authorization tests that vault tokens are removed from the Authorization
// header before it is passed through to the backend
func TestPassthroughHeaders_Authorization(t *testing.T) {
	secretNoop := &vault.NoopBackend{
		Response: &logical.Response{Data: map[string]interface{}{}},
	}
	authNoop := &vault.NoopBackend{
		Login:       []string{"login"},
		Response:    &logical.Response{},
		BackendType: logical.TypeCredential,
	}
	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"noop": func(ctx context.Context, config *logical.BackendConfig) (logical.Backend, error) {
				return secretNoop, nil
			},
		},
		CredentialBackends: map[string]logical.Factory{
			"noop": func(ctx context.Context, config *logical.BackendConfig) (logical.Backend, error) {
				return authNoop, nil
			},
		},
	}
	cluster := minimal.NewTestSoloCluster(t, coreConfig)
	client := cluster.Cores[0].Client

	// Enable logical backend with Authorization passthrough.
	err := client.Sys().Mount("foo", &api.MountInput{
		Type: "noop",
		Config: api.MountConfigInput{
			PassthroughRequestHeaders: []string{"Authorization"},
		},
	})
	require.NoError(t, err)

	// Enable credential backend with Authorization passthrough.
	err = client.Sys().EnableAuthWithOptions("bar", &api.EnableAuthOptions{
		Type: "noop",
		Config: api.AuthConfigInput{
			PassthroughRequestHeaders: []string{"Authorization"},
		},
	})
	require.NoError(t, err)

	token := cluster.RootToken
	httpClient := client.CloneConfig().HttpClient
	address := client.Address()
	u, err := url.Parse(address)
	require.NoError(t, err)
	doRequest := func(path string, operation logical.Operation, authHeaderValue []string) error {
		httpOp := "PUT"
		if operation == logical.ReadOperation {
			httpOp = "GET"
		}
		u.Path = "/v1/" + path
		req := &http.Request{
			Method: httpOp,
			URL:    u,
			Header: map[string][]string{
				"Authorization": authHeaderValue,
			},
		}

		resp, err := httpClient.Do(req)
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
		return err
	}
	t.Run("authenticated request strips bearer from authz header", func(t *testing.T) {
		err := doRequest("foo/test", logical.ReadOperation, []string{"Bearer " + token, "Basic dXNlcjpwYXNz"})
		require.NoError(t, err)
		require.NotEmpty(t, secretNoop.Requests)
		headers := secretNoop.Requests[len(secretNoop.Requests)-1].Headers
		require.Equal(t, []string{"Basic dXNlcjpwYXNz"}, headers["Authorization"])
	})

	t.Run("login request with valid token strips bearer", func(t *testing.T) {
		err = doRequest("auth/bar/login", logical.UpdateOperation, []string{"Bearer " + token, "Basic dXNlcjpwYXNz"})
		require.NoError(t, err)
		require.NotEmpty(t, authNoop.Requests)

		headers := authNoop.Requests[len(authNoop.Requests)-1].Headers
		require.Equal(t, []string{"Basic dXNlcjpwYXNz"}, headers["Authorization"])
	})

	t.Run("login request with non-vault bearer keeps header", func(t *testing.T) {
		err = doRequest("auth/bar/login", logical.UpdateOperation, []string{"Bearer not-a-valid-vault-token", "Basic dXNlcjpwYXNz"})
		require.NoError(t, err)
		require.NotEmpty(t, authNoop.Requests)

		headers := authNoop.Requests[len(authNoop.Requests)-1].Headers
		require.Equal(t, []string{"Bearer not-a-valid-vault-token", "Basic dXNlcjpwYXNz"}, headers["Authorization"])
	})
}
