// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package delegated_auth

import (
	"context"
	"fmt"
	paths "path"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/credential/userpass"
	"github.com/hashicorp/vault/builtin/logical/totp"
	"github.com/hashicorp/vault/helper/testhelpers"
	"github.com/hashicorp/vault/helper/testhelpers/teststorage"
	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/require"
)

func buildDaError(defaults map[string]string, d *framework.FieldData) *logical.RequestDelegatedAuth {
	fieldDataOrDefault := func(field string, d *framework.FieldData) string {
		if val, ok := d.GetOk(field); ok {
			return val.(string)
		}

		return defaults[field]
	}

	accessor := fieldDataOrDefault("accessor", d)
	path := fieldDataOrDefault("path", d)
	username := fieldDataOrDefault("username", d)
	password := fieldDataOrDefault("password", d)
	var errorHandler logical.DelegatedAuthErrorHandler
	if handleErrorRaw, ok := d.GetOk("handle_error"); ok {
		if handleErrorRaw.(bool) {
			errorHandler = func(ctx context.Context, initiatingRequest, authRequest *logical.Request, authResponse *logical.Response, err error) (*logical.Response, error) {
				return logical.ErrorResponse(fmt.Sprintf("my custom handler: %v", err)), nil
			}
		}
	}

	loginPath := paths.Join(path, username)
	data := map[string]interface{}{"password": password}

	return logical.NewDelegatedAuthenticationRequest(accessor, loginPath, data, errorHandler)
}

func buildDelegatedAuthFactory(defaults map[string]string) logical.Factory {
	opHandler := func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		daError := buildDaError(defaults, d)
		if req.ClientToken == "" || req.ClientTokenSource != logical.ClientTokenFromInternalAuth {
			return nil, daError
		}

		if req.Operation == logical.ListOperation {
			return logical.ListResponse([]string{"success", req.ClientToken}), nil
		}

		if d.Get("loop").(bool) {
			return nil, daError
		}

		if d.Get("perform_write").(bool) {
			entry, err := logical.StorageEntryJSON("test", map[string]string{"test": "value"})
			if err != nil {
				return nil, err
			}
			if err = req.Storage.Put(ctx, entry); err != nil {
				return nil, err
			}
		}

		return &logical.Response{
			Data: map[string]interface{}{
				"success": true,
				"token":   req.ClientToken,
			},
		}, nil
	}

	return func(ctx context.Context, config *logical.BackendConfig) (logical.Backend, error) {
		b := new(framework.Backend)
		b.BackendType = logical.TypeLogical
		b.Paths = []*framework.Path{
			{
				Pattern: "preauth-test/list/?",
				Operations: map[logical.Operation]framework.OperationHandler{
					logical.ListOperation: &framework.PathOperation{Callback: opHandler},
				},
			},
			{
				Pattern: "preauth-test",
				Operations: map[logical.Operation]framework.OperationHandler{
					logical.ReadOperation:   &framework.PathOperation{Callback: opHandler},
					logical.PatchOperation:  &framework.PathOperation{Callback: opHandler},
					logical.UpdateOperation: &framework.PathOperation{Callback: opHandler},
					logical.DeleteOperation: &framework.PathOperation{Callback: opHandler},
				},
				Fields: map[string]*framework.FieldSchema{
					"accessor":      {Type: framework.TypeString},
					"path":          {Type: framework.TypeString},
					"username":      {Type: framework.TypeString},
					"password":      {Type: framework.TypeString},
					"loop":          {Type: framework.TypeBool},
					"perform_write": {Type: framework.TypeBool},
					"handle_error":  {Type: framework.TypeBool},
				},
			},
		}
		b.PathsSpecial = &logical.Paths{Unauthenticated: []string{"preauth-test", "preauth-test/*"}}
		err := b.Setup(ctx, config)
		return b, err
	}
}

func TestDelegatedAuth(t *testing.T) {
	t.Parallel()

	// A map of success values to be populated once and used in request
	// operations that can't pass in values
	delegatedReqDefaults := map[string]string{
		"username": "allowed-est",
		"password": "test",
		"path":     "login",
	}

	delegatedAuthFactory := buildDelegatedAuthFactory(delegatedReqDefaults)
	coreConfig := &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"userpass":  userpass.Factory,
			"userpass2": userpass.Factory,
		},
		LogicalBackends: map[string]logical.Factory{
			"delegateauthtest": delegatedAuthFactory,
			"totp":             totp.Factory,
		},
	}

	conf, opts := teststorage.ClusterSetup(coreConfig, &vault.TestClusterOptions{
		HandlerFunc: http.Handler,
		NumCores:    1,
	}, teststorage.InmemBackendSetup)

	cluster := vault.NewTestCluster(t, conf, opts)
	cluster.Start()
	defer cluster.Cleanup()

	client := testhelpers.WaitForActiveNode(t, cluster).Client

	// Setup two users, one with an allowed policy, another without a policy within userpass
	err := client.Sys().PutPolicy("allow-est",
		`path "dat/preauth-test" { capabilities = ["read", "create", "update", "patch", "delete"] }
               path "dat/preauth-test/*" { capabilities = ["read","list"] }`)
	require.NoError(t, err, "Failed to write policy allow-est")

	err = client.Sys().EnableAuthWithOptions("userpass", &api.EnableAuthOptions{
		Type: "userpass",
	})
	require.NoError(t, err, "failed mounting userpass endpoint")

	_, err = client.Logical().Write("auth/userpass/users/allowed-est", map[string]interface{}{
		"password":   "test",
		"policies":   "allow-est",
		"token_type": "batch",
	})
	require.NoError(t, err, "failed to create allowed-est user")

	_, err = client.Logical().Write("auth/userpass/users/not-allowed-est", map[string]interface{}{
		"password":   "test",
		"token_type": "batch",
	})
	require.NoError(t, err, "failed to create not-allowed-est user")

	_, err = client.Logical().Write("auth/userpass/users/bad-token-type-est", map[string]interface{}{
		"password":   "test",
		"token_type": "service",
	})
	require.NoError(t, err, "failed to create bad-token-type-est user")

	// Setup another auth mount so we can test multiple accessors in mount tuning works later
	err = client.Sys().EnableAuthWithOptions("userpass2", &api.EnableAuthOptions{
		Type: "userpass",
	})
	require.NoError(t, err, "failed mounting userpass2")

	_, err = client.Logical().Write("auth/userpass2/users/allowed-est-2", map[string]interface{}{
		"password":   "test",
		"policies":   "allow-est",
		"token_type": "batch",
	})
	require.NoError(t, err, "failed to create allowed-est-2 user")

	// Fetch the userpass auth accessors
	resp, err := client.Logical().Read("/sys/mounts/auth/userpass")
	require.NoError(t, err, "failed to query for mount accessor")
	require.NotNil(t, resp, "received nil response from mount accessor query")
	require.NotNil(t, resp.Data, "received response with nil Data")
	require.NotEmpty(t, resp.Data["accessor"], "Accessor field was empty: %v", resp)
	upAccessor := resp.Data["accessor"].(string)

	resp, err = client.Logical().Read("/sys/mounts/auth/userpass2")
	require.NoError(t, err, "failed to query for mount accessor for userpass2")
	require.NotNil(t, resp, "received nil response from mount accessor query for userpass2")
	require.NotNil(t, resp.Data, "received response with nil Data")
	require.NotEmpty(t, resp.Data["accessor"], "Accessor field was empty: %v", resp)
	upAccessor2 := resp.Data["accessor"].(string)

	resp, err = client.Logical().Read("/sys/mounts/cubbyhole")
	require.NoError(t, err, "failed to query for mount accessor for cubbyhole")
	require.NotNil(t, resp, "received nil response from mount accessor query for cubbyhole")
	require.NotNil(t, resp.Data, "received response with nil Data")
	require.NotEmpty(t, resp.Data["accessor"], "Accessor field was empty: %v", resp)
	cubbyAccessor := resp.Data["accessor"].(string)

	// Set up our backend mount that will delegate its auth to the userpass mount
	err = client.Sys().Mount("dat", &api.MountInput{
		Type: "delegateauthtest",
		Config: api.MountConfigInput{
			DelegatedAuthAccessors: []string{upAccessor, "an-accessor-that-does-not-exist", cubbyAccessor},
		},
	})
	require.NoError(t, err, "failed mounting delegated auth endpoint")

	delegatedReqDefaults["accessor"] = upAccessor

	// We want a client without any previous tokens set to make sure we aren't using
	// the other token.
	clientNoToken, err := client.Clone()
	require.NoError(t, err, "failed cloning client")
	clientNoToken.ClearToken()

	// Happy path test for the various operation types we want to support, make sure
	// for each one that we don't error out and we get back a token value from the backend
	// call.
	for _, test := range []string{"delete", "read", "list", "write"} {
		t.Run("op-"+test, func(st *testing.T) {
			switch test {
			case "delete":
				resp, err = clientNoToken.Logical().Delete("dat/preauth-test")
			case "read":
				resp, err = clientNoToken.Logical().Read("dat/preauth-test")
			case "list":
				resp, err = clientNoToken.Logical().List("dat/preauth-test/list/")
			case "write":
				resp, err = clientNoToken.Logical().Write("dat/preauth-test", map[string]interface{}{
					"accessor": delegatedReqDefaults["accessor"],
					"path":     delegatedReqDefaults["path"],
					"username": delegatedReqDefaults["username"],
					"password": delegatedReqDefaults["password"],
				})
			}

			require.NoErrorf(st, err, "failed making %s pre-auth call with allowed-est", test)
			require.NotNilf(st, resp, "pre-auth %s call returned nil", test)
			require.NotNil(t, resp.Data, "received response with nil Data")
			if test != "list" {
				require.Equalf(st, true, resp.Data["success"], "Got an incorrect response from %s call in success field", test)
				require.NotEmptyf(st, resp.Data["token"], "no token returned by %s handler", test)
			} else {
				require.NotEmpty(st, resp.Data["keys"], "list operation did not contain keys in response")
				keys := resp.Data["keys"].([]interface{})
				require.Equal(st, 2, len(keys), "keys field did not contain expected 2 elements")
				require.Equal(st, "success", keys[0], "the first keys field did not contain expected value")
				require.NotEmpty(st, keys[1], "the second keys field did not contain a token")
			}
		})
	}

	// Test various failure scenarios
	failureTests := []struct {
		name          string
		accessor      string
		path          string
		username      string
		password      string
		errorContains string
		forceLoop     bool
	}{
		{
			name:          "policy-denies-user",
			accessor:      upAccessor,
			path:          "login",
			username:      "not-allowed-est",
			password:      "test",
			errorContains: "permission denied",
		},
		{
			name:          "bad-password",
			accessor:      upAccessor,
			path:          "login",
			username:      "allowed-est",
			password:      "bad-password",
			errorContains: "invalid credentials",
		},
		{
			name:          "unknown-user",
			accessor:      upAccessor,
			path:          "login",
			username:      "non-existant-user",
			password:      "test",
			errorContains: "invalid username or password",
		},
		{
			name:          "missing-user",
			accessor:      upAccessor,
			path:          "login",
			username:      "",
			password:      "test",
			errorContains: "was not considered a login request",
		},
		{
			name:          "missing-password",
			accessor:      upAccessor,
			path:          "login",
			username:      "allowed-est",
			password:      "",
			errorContains: "missing password",
		},
		{
			name:          "bad-path-within-delegated-auth-error",
			accessor:      upAccessor,
			path:          "not-the-login-path",
			username:      "allowed-est",
			password:      "test",
			errorContains: "was not considered a login request",
		},
		{
			name:          "empty-path-within-delegated-auth-error",
			accessor:      upAccessor,
			path:          "",
			username:      "allowed-est",
			password:      "test",
			errorContains: "was not considered a login request",
		},
		{
			name:          "empty-accessor-within-delegated-auth-error",
			accessor:      "",
			path:          "login",
			username:      "allowed-est",
			password:      "test",
			errorContains: "backend returned an invalid mount accessor",
		},
		{
			name:          "accessor-does-not-exist-within-delegated-auth-error",
			accessor:      "an-accessor-that-does-not-exist",
			path:          "login",
			username:      "allowed-est",
			password:      "test",
			errorContains: "requested delegate authentication accessor 'an-accessor-that-does-not-exist' was not found",
		},
		{
			name:          "non-allowed-accessor-within-delegated-auth-error",
			accessor:      upAccessor2,
			path:          "login",
			username:      "allowed-est-2",
			password:      "test",
			errorContains: fmt.Sprintf("delegated auth to accessor %s not permitted", upAccessor2),
		},
		{
			name:          "force-constant-login-request-loop",
			accessor:      upAccessor,
			path:          "login",
			username:      "allowed-est",
			password:      "test",
			forceLoop:     true,
			errorContains: "delegated authentication requested but authentication token present",
		},
		{
			name:          "non-auth-mount-accessor",
			accessor:      cubbyAccessor,
			path:          "login",
			username:      "allowed-est",
			password:      "test",
			errorContains: fmt.Sprintf("requested delegate authentication mount '%s' was not an auth mount", cubbyAccessor),
		},
		{
			name:          "fails-on-non-batch-token",
			accessor:      upAccessor,
			path:          "login",
			username:      "bad-token-type-est",
			password:      "test",
			errorContains: "delegated auth requests must be configured to issue batch tokens",
		},
	}
	for _, test := range failureTests {
		t.Run(test.name, func(st *testing.T) {
			resp, err = clientNoToken.Logical().Write("dat/preauth-test", map[string]interface{}{
				"accessor": test.accessor,
				"path":     test.path,
				"username": test.username,
				"password": test.password,
				"loop":     test.forceLoop,
			})
			if test.errorContains != "" {
				require.ErrorContains(st, err, test.errorContains,
					"pre-auth call should have failed due to policy restriction got resp: %v err: %v", resp, err)
			} else {
				require.Error(st, err, "Expected failure got resp: %v err: %v", resp, err)
			}
		})
	}

	// Make sure we can add an accessor to the mount that previously failed above,
	// and the request handling code does use both accessor values.
	t.Run("multiple-accessors", func(st *testing.T) {
		err = client.Sys().TuneMount("dat", api.MountConfigInput{DelegatedAuthAccessors: []string{upAccessor, upAccessor2}})
		require.NoError(t, err, "Failed to tune mount to update delegated auth accessors")

		resp, err = clientNoToken.Logical().Write("dat/preauth-test", map[string]interface{}{
			"accessor": upAccessor2,
			"path":     "login",
			"username": "allowed-est-2",
			"password": "test",
		})

		require.NoError(st, err, "failed making pre-auth call with allowed-est-2")
		require.NotNil(st, resp, "pre-auth %s call returned nil with allowed-est-2")
		require.NotNil(t, resp.Data, "received response with nil Data")
		require.Equal(st, true, resp.Data["success"], "Got an incorrect response from call in success field with allowed-est-2")
		require.NotEmpty(st, resp.Data["token"], "no token returned with allowed-est-2 user")
	})

	// Test we can delegate a permission denied error back to the originating
	// backend for processing/response to the client
	t.Run("backend-handles-permission-denied", func(st *testing.T) {
		resp, err = clientNoToken.Logical().Write("dat/preauth-test", map[string]interface{}{
			"accessor":     upAccessor,
			"path":         "login",
			"username":     "allowed-est",
			"password":     "test2",
			"handle_error": true,
		})

		require.ErrorContains(st, err, "my custom handler: invalid credentials")
	})

	// Test we can delegate a permission denied error back to the originating
	// backend for processing/response to the client
	t.Run("mfa-request-is-denied", func(st *testing.T) {
		// Mount the totp secrets engine
		testhelpers.SetupTOTPMount(st, client)

		// Create a test entity and alias
		totpUser := "test-totp"
		testhelpers.CreateEntityAndAlias(st, client, upAccessor, "entity1", totpUser)

		// Configure a default TOTP method
		methodID := testhelpers.SetupTOTPMethod(st, client, map[string]interface{}{
			"issuer":                  "yCorp",
			"period":                  5,
			"algorithm":               "SHA256",
			"digits":                  6,
			"skew":                    1,
			"key_size":                20,
			"qr_size":                 200,
			"max_validation_attempts": 5,
			"method_name":             "foo",
		})

		// Configure a default login enforcement
		enforcementConfig := map[string]interface{}{
			"auth_method_types": []string{"userpass"},
			"name":              methodID[0:4],
			"mfa_method_ids":    []string{methodID},
		}

		testhelpers.SetupMFALoginEnforcement(st, client, enforcementConfig)

		_, err = clientNoToken.Logical().Write("dat/preauth-test", map[string]interface{}{
			"accessor":     upAccessor,
			"path":         "login",
			"username":     totpUser,
			"password":     "testpassword",
			"handle_error": true,
		})

		require.ErrorContains(st, err, "delegated auth request requiring MFA is not supported")
	})

	// Test the behavior around receiving a request asking for response wrapping and
	// being delegated to the secondary query we do
	t.Run("response-wrapping-test", func(st *testing.T) {
		resWrapClient, err := client.Clone()
		require.NoError(st, err, "failed cloning client for response wrapping")

		resWrapClient.SetWrappingLookupFunc(func(operation, path string) string {
			if path == "dat/preauth-test" {
				return "15s"
			}
			return ""
		})

		resp, err = resWrapClient.Logical().Write("dat/preauth-test", map[string]interface{}{
			"accessor": upAccessor,
			"path":     "login",
			"username": "allowed-est",
			"password": "test",
		})
		require.NoError(st, err, "failed calling preauth-test api with response wrapping")
		require.NotNil(st, resp, "Got nil, nil response from preauth-test api with response wrapping")
		require.NotNil(st, resp.WrapInfo, "response object didn't contain wrapped info")

		unwrapClient, err := client.Clone()
		require.NoError(st, err, "failed cloning client for lookups")
		wrapToken := resp.WrapInfo.Token
		unwrapClient.SetToken(wrapToken)

		unwrapResp, err := unwrapClient.Logical().Write("sys/wrapping/unwrap", map[string]interface{}{})
		require.NoError(st, err, "failed unwrap call")
		require.NotNil(st, unwrapResp, "unwrap response was nil")
		require.NotNil(st, unwrapResp.Data, "unwrap response did not contain Data")
		require.Contains(st, unwrapResp.Data, "success", "unwrap response data did not contain success field")
		require.Contains(st, unwrapResp.Data, "token", "unwrap response data did not contain token field")
		require.Equal(st, true, unwrapResp.Data["success"], "Got an incorrect response in success field within unwrap call")
		require.NotEmptyf(st, unwrapResp.Data["token"], "no token returned by handler within unwrap call")
	})
}
