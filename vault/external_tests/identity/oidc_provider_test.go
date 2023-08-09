// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package identity

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/cap/oidc"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/credential/userpass"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/require"
)

const (
	testPassword           = "testpassword"
	testRedirectURI        = "https://127.0.0.1:8251/callback"
	testGroupScopeTemplate = `
		{
			"groups": {{identity.entity.groups.names}}
		}
	`
	testUserScopeTemplate = `
		{
			"username": {{identity.entity.aliases.%s.name}},
			"contact": {
				"email": {{identity.entity.metadata.email}},
				"phone_number": {{identity.entity.metadata.phone_number}}
			}
		}
	`
)

// TestOIDC_Auth_Code_Flow_Default_Resources tests the authorization
// code flow using the default OIDC provider, default key, and allow_all
// assignment. This ensures that the resources are created and usable with
// an initial setup of Vault.
func TestOIDC_Auth_Code_Flow_Default_Resources(t *testing.T) {
	cluster := setupOIDCTestCluster(t, 2)
	defer cluster.Cleanup()
	active := cluster.Cores[0].Client
	standby := cluster.Cores[1].Client

	// Enable userpass auth and create a user
	err := active.Sys().EnableAuthWithOptions("userpass", &api.EnableAuthOptions{
		Type: "userpass",
	})
	require.NoError(t, err)
	_, err = active.Logical().Write("auth/userpass/users/end-user", map[string]interface{}{
		"password": testPassword,
	})
	require.NoError(t, err)

	// Create a confidential client
	_, err = active.Logical().Write("identity/oidc/client/confidential", map[string]interface{}{
		"redirect_uris":    []string{testRedirectURI},
		"assignments":      []string{"allow_all"},
		"id_token_ttl":     "1h",
		"access_token_ttl": "30m",
	})
	require.NoError(t, err)

	// Read the client ID and secret in order to configure the OIDC client
	resp, err := active.Logical().Read("identity/oidc/client/confidential")
	require.NoError(t, err)
	clientID := resp.Data["client_id"].(string)
	clientSecret := resp.Data["client_secret"].(string)

	// We aren't going to open up a browser to facilitate the login and redirect
	// from this test, so we'll log in via userpass and set the client's token as
	// the token that results from the authentication.
	resp, err = active.Logical().Write("auth/userpass/login/end-user", map[string]interface{}{
		"password": testPassword,
	})
	require.NoError(t, err)
	clientToken := resp.Auth.ClientToken
	entityID := resp.Auth.EntityID

	// Look up the token to get its creation time. This will be used for test
	// cases that make assertions on the max_age parameter and auth_time claim.
	resp, err = active.Logical().Write("auth/token/lookup", map[string]interface{}{
		"token": clientToken,
	})
	require.NoError(t, err)
	expectedAuthTime, err := strconv.Atoi(string(resp.Data["creation_time"].(json.Number)))
	require.NoError(t, err)

	// Read the issuer from the OIDC provider's discovery document
	var discovery struct {
		Issuer string `json:"issuer"`
	}
	decodeRawRequest(t, active, http.MethodGet,
		"/v1/identity/oidc/provider/default/.well-known/openid-configuration",
		nil, &discovery)

	// Create the client-side OIDC provider config
	pc, err := oidc.NewConfig(discovery.Issuer, clientID,
		oidc.ClientSecret(clientSecret), []oidc.Alg{oidc.RS256},
		[]string{testRedirectURI}, oidc.WithProviderCA(string(cluster.CACertPEM)))
	require.NoError(t, err)

	// Create the client-side OIDC provider
	p, err := oidc.NewProvider(pc)
	require.NoError(t, err)
	defer p.Done()

	// Create the client-side PKCE code verifier
	v, err := oidc.NewCodeVerifier()
	require.NoError(t, err)

	type args struct {
		useStandby bool
		options    []oidc.Option
	}
	tests := []struct {
		name     string
		args     args
		expected string
	}{
		{
			name: "active: authorization code flow",
			args: args{
				options: []oidc.Option{
					oidc.WithScopes("openid"),
				},
			},
			expected: fmt.Sprintf(`{
					"iss": "%s",
					"aud": "%s",
					"sub": "%s",
					"namespace": "root"
				}`, discovery.Issuer, clientID, entityID),
		},
		{
			name: "active: authorization code flow with max_age parameter",
			args: args{
				options: []oidc.Option{
					oidc.WithScopes("openid"),
					oidc.WithMaxAge(60),
				},
			},
			expected: fmt.Sprintf(`{
				"iss": "%s",
				"aud": "%s",
				"sub": "%s",
				"namespace": "root",
				"auth_time": %d
			}`, discovery.Issuer, clientID, entityID, expectedAuthTime),
		},
		{
			name: "active: authorization code flow with Proof Key for Code Exchange (PKCE)",
			args: args{
				options: []oidc.Option{
					oidc.WithScopes("openid"),
					oidc.WithPKCE(v),
				},
			},
			expected: fmt.Sprintf(`{
				"iss": "%s",
				"aud": "%s",
				"sub": "%s",
				"namespace": "root"
			}`, discovery.Issuer, clientID, entityID),
		},
		{
			name: "standby: authorization code flow",
			args: args{
				useStandby: true,
				options: []oidc.Option{
					oidc.WithScopes("openid"),
				},
			},
			expected: fmt.Sprintf(`{
				"iss": "%s",
				"aud": "%s",
				"sub": "%s",
				"namespace": "root"
			}`, discovery.Issuer, clientID, entityID),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := active
			if tt.args.useStandby {
				client = standby
			}
			client.SetToken(clientToken)

			// Create the client-side OIDC request state
			oidcRequest, err := oidc.NewRequest(10*time.Minute, testRedirectURI, tt.args.options...)
			require.NoError(t, err)

			// Get the URL for the authorization endpoint from the OIDC client
			authURL, err := p.AuthURL(context.Background(), oidcRequest)
			require.NoError(t, err)
			parsedAuthURL, err := url.Parse(authURL)
			require.NoError(t, err)

			// This replace only occurs because we're not using the browser in this test
			authURLPath := strings.Replace(parsedAuthURL.Path, "/ui/vault/", "/v1/", 1)

			// Kick off the authorization code flow
			var authResp struct {
				Code  string `json:"code"`
				State string `json:"state"`
			}
			decodeRawRequest(t, client, http.MethodGet, authURLPath, parsedAuthURL.Query(), &authResp)

			// The returned state must match the OIDC client state
			require.Equal(t, oidcRequest.State(), authResp.State)

			// Exchange the authorization code for an ID token and access token.
			// The ID token signature is verified using the provider's public keys after
			// the exchange takes place. The ID token is also validated according to the
			// client-side requirements of the OIDC spec. See the validation code at:
			// - https://github.com/hashicorp/cap/blob/main/oidc/provider.go#L240
			// - https://github.com/hashicorp/cap/blob/main/oidc/provider.go#L441
			token, err := p.Exchange(context.Background(), oidcRequest, authResp.State, authResp.Code)
			require.NoError(t, err)
			require.NotNil(t, token)
			idToken := token.IDToken()
			accessToken := token.StaticTokenSource()

			// Get the ID token claims
			allClaims := make(map[string]interface{})
			require.NoError(t, idToken.Claims(&allClaims))

			// Get the sub claim for userinfo validation
			require.NotEmpty(t, allClaims["sub"])
			subject := allClaims["sub"].(string)

			// Request userinfo using the access token
			err = p.UserInfo(context.Background(), accessToken, subject, &allClaims)
			require.NoError(t, err)

			// Assert that claims computed during the flow (i.e., not known
			// ahead of time in this test) are present as top-level keys
			for _, claim := range []string{"iat", "exp", "nonce", "at_hash", "c_hash"} {
				_, ok := allClaims[claim]
				require.True(t, ok)
			}

			// Assert that all other expected claims are populated
			expectedClaims := make(map[string]interface{})
			require.NoError(t, json.Unmarshal([]byte(tt.expected), &expectedClaims))
			for k, expectedVal := range expectedClaims {
				actualVal, ok := allClaims[k]
				require.True(t, ok)
				require.EqualValues(t, expectedVal, actualVal)
			}
		})
	}
}

// TestOIDC_Auth_Code_Flow_Confidential_CAP_Client tests the authorization code
// flow using a Vault OIDC provider. The test uses the CAP OIDC client to verify
// that the Vault OIDC provider's responses pass the various client-side validation
// requirements of the OIDC spec. This test uses a confidential client which has
// a client secret and authenticates to the token endpoint.
func TestOIDC_Auth_Code_Flow_Confidential_CAP_Client(t *testing.T) {
	cluster := setupOIDCTestCluster(t, 2)
	defer cluster.Cleanup()
	active := cluster.Cores[0].Client
	standby := cluster.Cores[1].Client

	// Create an entity with some metadata
	resp, err := active.Logical().Write("identity/entity", map[string]interface{}{
		"name": "test-entity",
		"metadata": map[string]string{
			"email":        "test@hashicorp.com",
			"phone_number": "123-456-7890",
		},
	})
	require.NoError(t, err)
	entityID := resp.Data["id"].(string)

	// Create a group
	resp, err = active.Logical().Write("identity/group", map[string]interface{}{
		"name":              "engineering",
		"member_entity_ids": []string{entityID},
	})
	require.NoError(t, err)
	groupID := resp.Data["id"].(string)

	// Create a policy that allows updating the provider
	err = active.Sys().PutPolicy("test-policy", `
		path "identity/oidc/provider/test-provider" {
			capabilities = ["update"]
		}
	`)
	require.NoError(t, err)

	// Enable userpass auth and create a user
	err = active.Sys().EnableAuthWithOptions("userpass", &api.EnableAuthOptions{
		Type: "userpass",
	})
	require.NoError(t, err)
	_, err = active.Logical().Write("auth/userpass/users/end-user", map[string]interface{}{
		"password":       testPassword,
		"token_policies": "test-policy",
	})
	require.NoError(t, err)

	// Get the userpass mount accessor
	mounts, err := active.Sys().ListAuth()
	require.NoError(t, err)
	var mountAccessor string
	for k, v := range mounts {
		if k == "userpass/" {
			mountAccessor = v.Accessor
			break
		}
	}
	require.NotEmpty(t, mountAccessor)

	// Create an entity alias
	_, err = active.Logical().Write("identity/entity-alias", map[string]interface{}{
		"name":           "end-user",
		"canonical_id":   entityID,
		"mount_accessor": mountAccessor,
	})
	require.NoError(t, err)

	// Create some custom scopes
	_, err = active.Logical().Write("identity/oidc/scope/groups", map[string]interface{}{
		"template": testGroupScopeTemplate,
	})
	require.NoError(t, err)
	_, err = active.Logical().Write("identity/oidc/scope/user", map[string]interface{}{
		"template": fmt.Sprintf(testUserScopeTemplate, mountAccessor),
	})
	require.NoError(t, err)

	// Create a key
	_, err = active.Logical().Write("identity/oidc/key/test-key", map[string]interface{}{
		"allowed_client_ids": []string{"*"},
		"algorithm":          "RS256",
	})
	require.NoError(t, err)

	// Create an assignment
	_, err = active.Logical().Write("identity/oidc/assignment/test-assignment", map[string]interface{}{
		"entity_ids": []string{entityID},
		"group_ids":  []string{groupID},
	})
	require.NoError(t, err)

	// Create a confidential client
	_, err = active.Logical().Write("identity/oidc/client/confidential", map[string]interface{}{
		"key":              "test-key",
		"redirect_uris":    []string{testRedirectURI},
		"assignments":      []string{"test-assignment"},
		"id_token_ttl":     "1h",
		"access_token_ttl": "30m",
	})
	require.NoError(t, err)

	// Read the client ID and secret in order to configure the OIDC client
	resp, err = active.Logical().Read("identity/oidc/client/confidential")
	require.NoError(t, err)
	clientID := resp.Data["client_id"].(string)
	clientSecret := resp.Data["client_secret"].(string)

	// Create the OIDC provider
	_, err = active.Logical().Write("identity/oidc/provider/test-provider", map[string]interface{}{
		"allowed_client_ids": []string{clientID},
		"scopes_supported":   []string{"user", "groups"},
	})
	require.NoError(t, err)

	// We aren't going to open up a browser to facilitate the login and redirect
	// from this test, so we'll log in via userpass and set the client's token as
	// the token that results from the authentication.
	resp, err = active.Logical().Write("auth/userpass/login/end-user", map[string]interface{}{
		"password": testPassword,
	})
	require.NoError(t, err)
	clientToken := resp.Auth.ClientToken

	// Look up the token to get its creation time. This will be used for test
	// cases that make assertions on the max_age parameter and auth_time claim.
	resp, err = active.Logical().Write("auth/token/lookup", map[string]interface{}{
		"token": clientToken,
	})
	require.NoError(t, err)
	expectedAuthTime, err := strconv.Atoi(string(resp.Data["creation_time"].(json.Number)))
	require.NoError(t, err)

	// Read the issuer from the OIDC provider's discovery document
	var discovery struct {
		Issuer string `json:"issuer"`
	}
	decodeRawRequest(t, active, http.MethodGet,
		"/v1/identity/oidc/provider/test-provider/.well-known/openid-configuration",
		nil, &discovery)

	// Create the client-side OIDC provider config
	pc, err := oidc.NewConfig(discovery.Issuer, clientID,
		oidc.ClientSecret(clientSecret), []oidc.Alg{oidc.RS256},
		[]string{testRedirectURI}, oidc.WithProviderCA(string(cluster.CACertPEM)))
	require.NoError(t, err)

	// Create the client-side OIDC provider
	p, err := oidc.NewProvider(pc)
	require.NoError(t, err)
	defer p.Done()

	// Create the client-side PKCE code verifier
	v, err := oidc.NewCodeVerifier()
	require.NoError(t, err)

	type args struct {
		useStandby bool
		options    []oidc.Option
	}
	tests := []struct {
		name     string
		args     args
		expected string
	}{
		{
			name: "active: authorization code flow",
			args: args{
				options: []oidc.Option{
					oidc.WithScopes("openid user"),
				},
			},
			expected: fmt.Sprintf(`{
					"iss": "%s",
					"aud": "%s",
					"sub": "%s",
					"namespace": "root",
					"username": "end-user",
					"contact": {
						"email": "test@hashicorp.com",
						"phone_number": "123-456-7890"
					}
				}`, discovery.Issuer, clientID, entityID),
		},
		{
			name: "active: authorization code flow with additional scopes",
			args: args{
				options: []oidc.Option{
					oidc.WithScopes("openid user groups"),
				},
			},
			expected: fmt.Sprintf(`{
				"iss": "%s",
				"aud": "%s",
				"sub": "%s",
				"namespace": "root",
				"username": "end-user",
				"contact": {
					"email": "test@hashicorp.com",
					"phone_number": "123-456-7890"
				},
				"groups": ["engineering"]
			}`, discovery.Issuer, clientID, entityID),
		},
		{
			name: "active: authorization code flow with max_age parameter",
			args: args{
				options: []oidc.Option{
					oidc.WithScopes("openid"),
					oidc.WithMaxAge(60),
				},
			},
			expected: fmt.Sprintf(`{
				"iss": "%s",
				"aud": "%s",
				"sub": "%s",
				"namespace": "root",
				"auth_time": %d
			}`, discovery.Issuer, clientID, entityID, expectedAuthTime),
		},
		{
			name: "active: authorization code flow with Proof Key for Code Exchange (PKCE)",
			args: args{
				options: []oidc.Option{
					oidc.WithScopes("openid"),
					oidc.WithPKCE(v),
				},
			},
			expected: fmt.Sprintf(`{
				"iss": "%s",
				"aud": "%s",
				"sub": "%s",
				"namespace": "root"
			}`, discovery.Issuer, clientID, entityID),
		},
		{
			name: "standby: authorization code flow with additional scopes",
			args: args{
				useStandby: true,
				options: []oidc.Option{
					oidc.WithScopes("openid user groups"),
				},
			},
			expected: fmt.Sprintf(`{
				"iss": "%s",
				"aud": "%s",
				"sub": "%s",
				"namespace": "root",
				"username": "end-user",
				"contact": {
					"email": "test@hashicorp.com",
					"phone_number": "123-456-7890"
				},
				"groups": ["engineering"]
			}`, discovery.Issuer, clientID, entityID),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := active
			if tt.args.useStandby {
				client = standby
			}
			client.SetToken(clientToken)

			// Update allowed client IDs before the authentication flow
			_, err = client.Logical().Write("identity/oidc/provider/test-provider", map[string]interface{}{
				"allowed_client_ids": []string{clientID},
			})
			require.NoError(t, err)

			// Create the client-side OIDC request state
			oidcRequest, err := oidc.NewRequest(10*time.Minute, testRedirectURI, tt.args.options...)
			require.NoError(t, err)

			// Get the URL for the authorization endpoint from the OIDC client
			authURL, err := p.AuthURL(context.Background(), oidcRequest)
			require.NoError(t, err)
			parsedAuthURL, err := url.Parse(authURL)
			require.NoError(t, err)

			// This replace only occurs because we're not using the browser in this test
			authURLPath := strings.Replace(parsedAuthURL.Path, "/ui/vault/", "/v1/", 1)

			// Kick off the authorization code flow
			var authResp struct {
				Code  string `json:"code"`
				State string `json:"state"`
			}
			decodeRawRequest(t, client, http.MethodGet, authURLPath, parsedAuthURL.Query(), &authResp)

			// The returned state must match the OIDC client state
			require.Equal(t, oidcRequest.State(), authResp.State)

			// Exchange the authorization code for an ID token and access token.
			// The ID token signature is verified using the provider's public keys after
			// the exchange takes place. The ID token is also validated according to the
			// client-side requirements of the OIDC spec. See the validation code at:
			// - https://github.com/hashicorp/cap/blob/main/oidc/provider.go#L240
			// - https://github.com/hashicorp/cap/blob/main/oidc/provider.go#L441
			token, err := p.Exchange(context.Background(), oidcRequest, authResp.State, authResp.Code)
			require.NoError(t, err)
			require.NotNil(t, token)
			idToken := token.IDToken()
			accessToken := token.StaticTokenSource()

			// Get the ID token claims
			allClaims := make(map[string]interface{})
			require.NoError(t, idToken.Claims(&allClaims))

			// Get the sub claim for userinfo validation
			require.NotEmpty(t, allClaims["sub"])
			subject := allClaims["sub"].(string)

			// Request userinfo using the access token
			err = p.UserInfo(context.Background(), accessToken, subject, &allClaims)
			require.NoError(t, err)

			// Assert that claims computed during the flow (i.e., not known
			// ahead of time in this test) are present as top-level keys
			for _, claim := range []string{"iat", "exp", "nonce", "at_hash", "c_hash"} {
				_, ok := allClaims[claim]
				require.True(t, ok)
			}

			// Assert that all other expected claims are populated
			expectedClaims := make(map[string]interface{})
			require.NoError(t, json.Unmarshal([]byte(tt.expected), &expectedClaims))
			for k, expectedVal := range expectedClaims {
				actualVal, ok := allClaims[k]
				require.True(t, ok)
				require.EqualValues(t, expectedVal, actualVal)
			}

			// Assert that the access token is no longer able to obtain user info
			// after removing the client from the provider's allowed client ids
			_, err = client.Logical().Write("identity/oidc/provider/test-provider", map[string]interface{}{
				"allowed_client_ids": []string{},
			})
			require.NoError(t, err)
			err = p.UserInfo(context.Background(), accessToken, subject, &allClaims)
			require.Error(t, err)
			require.Equal(t, `Provider.UserInfo: provider UserInfo request failed: 403 Forbidden: {"error":"access_denied","error_description":"client is not authorized to use the provider"}`,
				err.Error())
		})
	}
}

// TestOIDC_Auth_Code_Flow_Public_CAP_Client tests the authorization code flow using
// a Vault OIDC provider. The test uses the CAP OIDC client to verify that the Vault
// OIDC provider's responses pass the various client-side validation requirements of
// the OIDC spec. This test uses a public client which does not have a client secret
// and always uses proof key for code exchange (PKCE).
func TestOIDC_Auth_Code_Flow_Public_CAP_Client(t *testing.T) {
	cluster := setupOIDCTestCluster(t, 2)
	defer cluster.Cleanup()
	active := cluster.Cores[0].Client
	standby := cluster.Cores[1].Client

	// Create an entity with some metadata
	resp, err := active.Logical().Write("identity/entity", map[string]interface{}{
		"name": "test-entity",
		"metadata": map[string]string{
			"email":        "test@hashicorp.com",
			"phone_number": "123-456-7890",
		},
	})
	require.NoError(t, err)
	entityID := resp.Data["id"].(string)

	// Create a group
	resp, err = active.Logical().Write("identity/group", map[string]interface{}{
		"name":              "engineering",
		"member_entity_ids": []string{entityID},
	})
	require.NoError(t, err)
	groupID := resp.Data["id"].(string)

	// Create a policy that allows updating the provider
	err = active.Sys().PutPolicy("test-policy", `
		path "identity/oidc/provider/test-provider" {
			capabilities = ["update"]
		}
	`)
	require.NoError(t, err)

	// Enable userpass auth and create a user
	err = active.Sys().EnableAuthWithOptions("userpass", &api.EnableAuthOptions{
		Type: "userpass",
	})
	require.NoError(t, err)
	_, err = active.Logical().Write("auth/userpass/users/end-user", map[string]interface{}{
		"password":       testPassword,
		"token_policies": "test-policy",
	})
	require.NoError(t, err)

	// Get the userpass mount accessor
	mounts, err := active.Sys().ListAuth()
	require.NoError(t, err)
	var mountAccessor string
	for k, v := range mounts {
		if k == "userpass/" {
			mountAccessor = v.Accessor
			break
		}
	}
	require.NotEmpty(t, mountAccessor)

	// Create an entity alias
	_, err = active.Logical().Write("identity/entity-alias", map[string]interface{}{
		"name":           "end-user",
		"canonical_id":   entityID,
		"mount_accessor": mountAccessor,
	})
	require.NoError(t, err)

	// Create some custom scopes
	_, err = active.Logical().Write("identity/oidc/scope/groups", map[string]interface{}{
		"template": testGroupScopeTemplate,
	})
	require.NoError(t, err)
	_, err = active.Logical().Write("identity/oidc/scope/user", map[string]interface{}{
		"template": fmt.Sprintf(testUserScopeTemplate, mountAccessor),
	})
	require.NoError(t, err)

	// Create a key
	_, err = active.Logical().Write("identity/oidc/key/test-key", map[string]interface{}{
		"allowed_client_ids": []string{"*"},
		"algorithm":          "RS256",
	})
	require.NoError(t, err)

	// Create an assignment
	_, err = active.Logical().Write("identity/oidc/assignment/test-assignment", map[string]interface{}{
		"entity_ids": []string{entityID},
		"group_ids":  []string{groupID},
	})
	require.NoError(t, err)

	// Create a public client
	_, err = active.Logical().Write("identity/oidc/client/public", map[string]interface{}{
		"key":              "test-key",
		"redirect_uris":    []string{testRedirectURI},
		"assignments":      []string{"test-assignment"},
		"id_token_ttl":     "1h",
		"access_token_ttl": "30m",
		"client_type":      "public",
	})
	require.NoError(t, err)

	// Read the client ID in order to configure the OIDC client
	resp, err = active.Logical().Read("identity/oidc/client/public")
	require.NoError(t, err)
	clientID := resp.Data["client_id"].(string)

	// Create the OIDC provider
	_, err = active.Logical().Write("identity/oidc/provider/test-provider", map[string]interface{}{
		"allowed_client_ids": []string{clientID},
		"scopes_supported":   []string{"user", "groups"},
	})
	require.NoError(t, err)

	// We aren't going to open up a browser to facilitate the login and redirect
	// from this test, so we'll log in via userpass and set the client's token as
	// the token that results from the authentication.
	resp, err = active.Logical().Write("auth/userpass/login/end-user", map[string]interface{}{
		"password": testPassword,
	})
	require.NoError(t, err)
	clientToken := resp.Auth.ClientToken

	// Look up the token to get its creation time. This will be used for test
	// cases that make assertions on the max_age parameter and auth_time claim.
	resp, err = active.Logical().Write("auth/token/lookup", map[string]interface{}{
		"token": clientToken,
	})
	require.NoError(t, err)
	expectedAuthTime, err := strconv.Atoi(string(resp.Data["creation_time"].(json.Number)))
	require.NoError(t, err)

	// Read the issuer from the OIDC provider's discovery document
	var discovery struct {
		Issuer string `json:"issuer"`
	}
	decodeRawRequest(t, active, http.MethodGet,
		"/v1/identity/oidc/provider/test-provider/.well-known/openid-configuration",
		nil, &discovery)

	// Create the client-side OIDC provider config with client secret intentionally empty
	clientSecret := oidc.ClientSecret("")
	pc, err := oidc.NewConfig(discovery.Issuer, clientID, clientSecret, []oidc.Alg{oidc.RS256},
		[]string{testRedirectURI}, oidc.WithProviderCA(string(cluster.CACertPEM)))
	require.NoError(t, err)

	// Create the client-side OIDC provider
	p, err := oidc.NewProvider(pc)
	require.NoError(t, err)
	defer p.Done()

	type args struct {
		useStandby bool
		options    []oidc.Option
	}
	tests := []struct {
		name     string
		args     args
		expected string
	}{
		{
			name: "active: authorization code flow",
			args: args{
				options: []oidc.Option{
					oidc.WithScopes("openid user"),
				},
			},
			expected: fmt.Sprintf(`{
					"iss": "%s",
					"aud": "%s",
					"sub": "%s",
					"namespace": "root",
					"username": "end-user",
					"contact": {
						"email": "test@hashicorp.com",
						"phone_number": "123-456-7890"
					}
				}`, discovery.Issuer, clientID, entityID),
		},
		{
			name: "active: authorization code flow with additional scopes",
			args: args{
				options: []oidc.Option{
					oidc.WithScopes("openid user groups"),
				},
			},
			expected: fmt.Sprintf(`{
				"iss": "%s",
				"aud": "%s",
				"sub": "%s",
				"namespace": "root",
				"username": "end-user",
				"contact": {
					"email": "test@hashicorp.com",
					"phone_number": "123-456-7890"
				},
				"groups": ["engineering"]
			}`, discovery.Issuer, clientID, entityID),
		},
		{
			name: "active: authorization code flow with max_age parameter",
			args: args{
				options: []oidc.Option{
					oidc.WithScopes("openid"),
					oidc.WithMaxAge(60),
				},
			},
			expected: fmt.Sprintf(`{
				"iss": "%s",
				"aud": "%s",
				"sub": "%s",
				"namespace": "root",
				"auth_time": %d
			}`, discovery.Issuer, clientID, entityID, expectedAuthTime),
		},
		{
			name: "standby: authorization code flow with additional scopes",
			args: args{
				useStandby: true,
				options: []oidc.Option{
					oidc.WithScopes("openid user groups"),
				},
			},
			expected: fmt.Sprintf(`{
				"iss": "%s",
				"aud": "%s",
				"sub": "%s",
				"namespace": "root",
				"username": "end-user",
				"contact": {
					"email": "test@hashicorp.com",
					"phone_number": "123-456-7890"
				},
				"groups": ["engineering"]
			}`, discovery.Issuer, clientID, entityID),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := active
			if tt.args.useStandby {
				client = standby
			}
			client.SetToken(clientToken)

			// Update allowed client IDs before the authentication flow
			_, err = client.Logical().Write("identity/oidc/provider/test-provider", map[string]interface{}{
				"allowed_client_ids": []string{clientID},
			})
			require.NoError(t, err)

			// Create the required client-side PKCE code verifier.
			v, err := oidc.NewCodeVerifier()
			require.NoError(t, err)
			options := append([]oidc.Option{oidc.WithPKCE(v)}, tt.args.options...)

			// Create the client-side OIDC request state
			oidcRequest, err := oidc.NewRequest(10*time.Minute, testRedirectURI, options...)
			require.NoError(t, err)

			// Get the URL for the authorization endpoint from the OIDC client
			authURL, err := p.AuthURL(context.Background(), oidcRequest)
			require.NoError(t, err)
			parsedAuthURL, err := url.Parse(authURL)
			require.NoError(t, err)

			// This replace only occurs because we're not using the browser in this test
			authURLPath := strings.Replace(parsedAuthURL.Path, "/ui/vault/", "/v1/", 1)

			// Kick off the authorization code flow
			var authResp struct {
				Code  string `json:"code"`
				State string `json:"state"`
			}
			decodeRawRequest(t, client, http.MethodGet, authURLPath, parsedAuthURL.Query(), &authResp)

			// The returned state must match the OIDC client state
			require.Equal(t, oidcRequest.State(), authResp.State)

			// Exchange the authorization code for an ID token and access token.
			// The ID token signature is verified using the provider's public keys after
			// the exchange takes place. The ID token is also validated according to the
			// client-side requirements of the OIDC spec. See the validation code at:
			// - https://github.com/hashicorp/cap/blob/main/oidc/provider.go#L240
			// - https://github.com/hashicorp/cap/blob/main/oidc/provider.go#L441
			token, err := p.Exchange(context.Background(), oidcRequest, authResp.State, authResp.Code)
			require.NoError(t, err)
			require.NotNil(t, token)
			idToken := token.IDToken()
			accessToken := token.StaticTokenSource()

			// Get the ID token claims
			allClaims := make(map[string]interface{})
			require.NoError(t, idToken.Claims(&allClaims))

			// Get the sub claim for userinfo validation
			require.NotEmpty(t, allClaims["sub"])
			subject := allClaims["sub"].(string)

			// Request userinfo using the access token
			err = p.UserInfo(context.Background(), accessToken, subject, &allClaims)
			require.NoError(t, err)

			// Assert that claims computed during the flow (i.e., not known
			// ahead of time in this test) are present as top-level keys
			for _, claim := range []string{"iat", "exp", "nonce", "at_hash", "c_hash"} {
				_, ok := allClaims[claim]
				require.True(t, ok)
			}

			// Assert that all other expected claims are populated
			expectedClaims := make(map[string]interface{})
			require.NoError(t, json.Unmarshal([]byte(tt.expected), &expectedClaims))
			for k, expectedVal := range expectedClaims {
				actualVal, ok := allClaims[k]
				require.True(t, ok)
				require.EqualValues(t, expectedVal, actualVal)
			}

			// Assert that the access token is no longer able to obtain user info
			// after removing the client from the provider's allowed client ids
			_, err = client.Logical().Write("identity/oidc/provider/test-provider", map[string]interface{}{
				"allowed_client_ids": []string{},
			})
			require.NoError(t, err)
			err = p.UserInfo(context.Background(), accessToken, subject, &allClaims)
			require.Error(t, err)
			require.Equal(t, `Provider.UserInfo: provider UserInfo request failed: 403 Forbidden: {"error":"access_denied","error_description":"client is not authorized to use the provider"}`,
				err.Error())
		})
	}
}

func setupOIDCTestCluster(t *testing.T, numCores int) *vault.TestCluster {
	t.Helper()

	coreConfig := &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"userpass": userpass.Factory,
		},
	}
	clusterOptions := &vault.TestClusterOptions{
		NumCores:    numCores,
		HandlerFunc: vaulthttp.Handler,
	}
	cluster := vault.NewTestCluster(t, coreConfig, clusterOptions)
	cluster.Start()
	vault.TestWaitActive(t, cluster.Cores[0].Core)

	return cluster
}

func decodeRawRequest(t *testing.T, client *api.Client, method, path string, params url.Values, v interface{}) {
	t.Helper()

	// Create the request and add query params if provided
	req := client.NewRequest(method, path)
	req.Params = params

	// Send the raw request
	r, err := client.RawRequest(req)
	require.NoError(t, err)
	require.NotNil(t, r)
	require.Equal(t, http.StatusOK, r.StatusCode)
	defer r.Body.Close()

	// Decode the body into v
	require.NoError(t, json.NewDecoder(r.Body).Decode(v))
}
