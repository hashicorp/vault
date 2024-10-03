// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/go-jose/go-jose/v3"
	"github.com/go-jose/go-jose/v3/jwt"
	"github.com/go-test/deep"
	"github.com/hashicorp/go-hclog"
	credUserpass "github.com/hashicorp/vault/builtin/credential/userpass"
	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	gocache "github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestOIDC_Path_OIDC_RoleNoKeyParameter tests that a role cannot be created
// without a key parameter
func TestOIDC_Path_OIDC_RoleNoKeyParameter(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	// Create a test role "test-role1" without a key param -- should fail
	resp, err := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/role/test-role1",
		Operation: logical.CreateOperation,
		Storage:   storage,
	})
	expectError(t, resp, err)
	// validate error message
	expectedStrings := map[string]interface{}{
		"the key parameter is required": true,
	}
	expectStrings(t, []string{resp.Data["error"].(string)}, expectedStrings)
}

// TestOIDC_Path_OIDC_RoleNilKeyEntry tests that a role cannot be created when
// a key parameter is provided but the key does not exist
func TestOIDC_Path_OIDC_RoleNilKeyEntry(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	// Create a test role "test-role1" with a non-existent key -- should fail
	resp, err := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/role/test-role1",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"key": "test-key",
		},
		Storage: storage,
	})
	expectError(t, resp, err)
	// validate error message
	expectedStrings := map[string]interface{}{
		"cannot find key \"test-key\"": true,
	}
	expectStrings(t, []string{resp.Data["error"].(string)}, expectedStrings)
}

// TestOIDC_Path_OIDCRole_UpdateNoKey test that we cannot update a role without
// prividing a key param
func TestOIDC_Path_OIDCRole_UpdateNoKey(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	// Create a test key "test-key"
	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/key/test-key",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"verification_ttl": "2m",
			"rotation_period":  "2m",
		},
		Storage: storage,
	})

	// Create a test role "test-role1" with a valid key -- should succeed
	resp, err := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/role/test-role1",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"key": "test-key",
			"ttl": "1m",
		},
		Storage: storage,
	})
	expectSuccess(t, resp, err)

	// Update "test-role1" without prividing a key param -- should succeed
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/role/test-role1",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"ttl": "2m",
		},
		Storage: storage,
	})
	expectSuccess(t, resp, err)

	// Read "test-role1" again and validate
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/role/test-role1",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)
	expected := map[string]interface{}{
		"key":       "test-key",
		"ttl":       int64(120),
		"template":  "",
		"client_id": resp.Data["client_id"],
	}
	if diff := deep.Equal(expected, resp.Data); diff != nil {
		t.Fatal(diff)
	}
}

// TestOIDC_Path_OIDCRole_UpdateEmptyKey test that we cannot update a role with an
// empty key
func TestOIDC_Path_OIDCRole_UpdateEmptyKey(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	// Create a test key "test-key"
	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/key/test-key",
		Operation: logical.CreateOperation,
		Storage:   storage,
	})

	// Create a test role "test-role1" with a valid key -- should succeed
	resp, err := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/role/test-role1",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"key": "test-key",
		},
		Storage: storage,
	})
	expectSuccess(t, resp, err)

	// Update "test-role1" with valid parameters -- should fail
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/role/test-role1",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"key": "",
		},
		Storage: storage,
	})
	expectError(t, resp, err)

	// Read "test-role1" again and validate
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/role/test-role1",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)
	expected := map[string]interface{}{
		"key":       "test-key",
		"ttl":       int64(86400),
		"template":  "",
		"client_id": resp.Data["client_id"],
	}
	if diff := deep.Equal(expected, resp.Data); diff != nil {
		t.Fatal(diff)
	}
}

// TestOIDC_Path_OIDCRoleRole tests CRUD operations for roles
func TestOIDC_Path_OIDCRoleRole(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	// Create a test key "test-key"
	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/key/test-key",
		Operation: logical.CreateOperation,
		Storage:   storage,
	})

	// Create a test role "test-role1" with a valid key -- should succeed with warning
	resp, err := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/role/test-role1",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"key": "test-key",
		},
		Storage: storage,
	})
	expectSuccess(t, resp, err)

	// Read "test-role1" and validate
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/role/test-role1",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)
	expected := map[string]interface{}{
		"key":       "test-key",
		"ttl":       int64(86400),
		"template":  "",
		"client_id": resp.Data["client_id"],
	}
	if diff := deep.Equal(expected, resp.Data); diff != nil {
		t.Fatal(diff)
	}

	// Update "test-role1" with valid parameters -- should succeed
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/role/test-role1",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"template":  "{\"some-key\":\"some-value\"}",
			"ttl":       "2h",
			"client_id": "my_custom_id",
		},
		Storage: storage,
	})
	expectSuccess(t, resp, err)

	// Read "test-role1" again and validate
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/role/test-role1",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)
	expected = map[string]interface{}{
		"key":       "test-key",
		"ttl":       int64(7200),
		"template":  "{\"some-key\":\"some-value\"}",
		"client_id": "my_custom_id",
	}
	if diff := deep.Equal(expected, resp.Data); diff != nil {
		t.Fatal(diff)
	}

	// Delete "test-role1"
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/role/test-role1",
		Operation: logical.DeleteOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)

	// Read "test-role1"
	respReadTestRole1AfterDelete, err3 := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/role/test-role1",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	// Ensure that "test-role1" has been deleted
	expectSuccess(t, respReadTestRole1AfterDelete, err3)
	if respReadTestRole1AfterDelete != nil {
		t.Fatalf("Expected a nil response but instead got:\n%#v", respReadTestRole1AfterDelete)
	}
	if respReadTestRole1AfterDelete != nil {
		t.Fatalf("Expected role to have been deleted but read response was:\n%#v", respReadTestRole1AfterDelete)
	}
}

// TestOIDC_Path_OIDCRole_InvalidTokenTTL tests the TokenTTL validation
func TestOIDC_Path_OIDCRole_InvalidTokenTTL(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	// Create a test key "test-key"
	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/key/test-key",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"verification_ttl": int64(60),
		},
		Storage: storage,
	})

	// Create a test role "test-role1" with a ttl longer than the
	// verification_ttl -- should fail with error
	resp, err := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/role/test-role1",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"key": "test-key",
			"ttl": int64(3600),
		},
		Storage: storage,
	})
	expectError(t, resp, err)

	// Read "test-role1"
	respReadTestRole1, err3 := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/role/test-role1",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	// Ensure that "test-role1" was not created
	expectSuccess(t, respReadTestRole1, err3)
	if respReadTestRole1 != nil {
		t.Fatalf("Expected a nil response but instead got:\n%#v", respReadTestRole1)
	}
}

// TestOIDC_Path_OIDCRole tests the List operation for roles
func TestOIDC_Path_OIDCRole(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	// Prepare two roles, test-role1 and test-role2
	// Create a test key "test-key"
	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/key/test-key",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"allowed_client_ids": "test-role1,test-role2",
		},
		Storage: storage,
	})

	// Create "test-role1"
	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/role/test-role1",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"key": "test-key",
		},
		Storage: storage,
	})

	// Create "test-role2"
	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/role/test-role2",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"key": "test-key",
		},
		Storage: storage,
	})

	// list roles
	respListRole, listErr := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/role",
		Operation: logical.ListOperation,
		Storage:   storage,
	})
	expectSuccess(t, respListRole, listErr)

	// validate list response
	expectedStrings := map[string]interface{}{"test-role1": true, "test-role2": true}
	expectStrings(t, respListRole.Data["keys"].([]string), expectedStrings)

	// delete test-role2
	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/role/test-role2",
		Operation: logical.DeleteOperation,
		Storage:   storage,
	})

	// list roles again and validate response
	respListRoleAfterDelete, listErrAfterDelete := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/role",
		Operation: logical.ListOperation,
		Storage:   storage,
	})
	expectSuccess(t, respListRoleAfterDelete, listErrAfterDelete)

	// validate list response
	delete(expectedStrings, "test-role2")
	expectStrings(t, respListRoleAfterDelete.Data["keys"].([]string), expectedStrings)
}

// TestOIDC_DeleteKeyWithMountReference ensures that keys cannot be deleted
// if they're referenced by mounts for plugin identity tokens.
func TestOIDC_DeleteKeyWithMountReference(t *testing.T) {
	ctx := namespace.RootContext(nil)
	core, _, _ := TestCoreUnsealed(t)
	core.credentialBackends["userpass"] = credUserpass.Factory
	idStorage := core.router.MatchingStorageByAPIPath(ctx, mountPathIdentity)
	require.NotNil(t, idStorage)

	tests := []struct {
		name        string
		mountPrefix string
		mountType   string
		keyName     string
	}{
		{
			name:        "delete key referenced by auth mount does not succeed",
			mountPrefix: "auth/",
			mountType:   "userpass/",
			keyName:     "test-key-1",
		},
		{
			name:        "delete key referenced by secret mount does not succeed",
			mountPrefix: "mounts/",
			mountType:   "kv/",
			keyName:     "test-key-2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := core.identityStore.HandleRequest(ctx, testKeyReq(idStorage, tt.keyName,
				[]string{"*"}, "RS256"))
			expectSuccess(t, resp, err)

			createMountEntryWithKey(t, ctx, core.systemBackend, tt.mountPrefix, tt.mountType, tt.keyName)
			require.NoError(t, err)
			require.Nil(t, resp)

			// Deleting the key must not succeed
			resp, err = core.identityStore.HandleRequest(ctx, &logical.Request{
				Path:      fmt.Sprintf("oidc/key/%s", tt.keyName),
				Operation: logical.DeleteOperation,
				Storage:   idStorage,
			})
			expectError(t, resp, err)
			require.Equal(t, fmt.Sprintf(deleteKeyErrorFmt, tt.keyName, "mounts", tt.mountType),
				resp.Error().Error())
		})
	}
}

// TestOIDC_Path_CRUDKey tests CRUD operations for keys
func TestOIDC_Path_CRUDKey(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	// Create a test key "test-key" -- should succeed
	resp, err := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/key/test-key",
		Operation: logical.CreateOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)

	// Read "test-key" and validate
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/key/test-key",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)
	expected := map[string]interface{}{
		"rotation_period":    int64(86400),
		"verification_ttl":   int64(86400),
		"algorithm":          "RS256",
		"allowed_client_ids": []string{},
	}
	if diff := deep.Equal(expected, resp.Data); diff != nil {
		t.Fatal(diff)
	}

	// Update "test-key" -- should succeed
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/key/test-key",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"rotation_period":    "10m",
			"verification_ttl":   "1h",
			"allowed_client_ids": "allowed-test-role",
		},
		Storage: storage,
	})
	expectSuccess(t, resp, err)

	// Read "test-key" again and validate
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/key/test-key",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)
	expected = map[string]interface{}{
		"rotation_period":    int64(600),
		"verification_ttl":   int64(3600),
		"algorithm":          "RS256",
		"allowed_client_ids": []string{"allowed-test-role"},
	}
	if diff := deep.Equal(expected, resp.Data); diff != nil {
		t.Fatal(diff)
	}

	// Create a role that depends on test key
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/role/allowed-test-role",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"key": "test-key",
		},
		Storage: storage,
	})
	expectSuccess(t, resp, err)

	// Delete test-key -- should fail because test-role depends on test-key
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/key/test-key",
		Operation: logical.DeleteOperation,
		Storage:   storage,
	})
	expectError(t, resp, err)
	// validate error message
	expectedStrings := map[string]interface{}{
		"unable to delete key \"test-key\" because it is currently referenced by these roles: allowed-test-role": true,
	}
	expectStrings(t, []string{resp.Data["error"].(string)}, expectedStrings)

	// Delete allowed-test-role
	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/role/allowed-test-role",
		Operation: logical.DeleteOperation,
		Storage:   storage,
	})

	// Delete test-key -- should succeed this time because no roles depend on test-key
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/key/test-key",
		Operation: logical.DeleteOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)
}

// TestOIDC_Path_OIDCKey_InvalidTokenTTL tests the TokenTTL validation
func TestOIDC_Path_OIDCKey_InvalidTokenTTL(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	// Create a test key "test-key" -- should succeed
	resp, err := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/key/test-key",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"verification_ttl": "4m",
		},
		Storage: storage,
	})
	expectSuccess(t, resp, err)

	// Create a role that depends on test key
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/role/allowed-test-role",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"key": "test-key",
			"ttl": "4m",
		},
		Storage: storage,
	})
	expectSuccess(t, resp, err)

	// Update "test-key" -- should fail since allowed-test-role ttl is less than 2m
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/key/test-key",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"rotation_period":    "10m",
			"verification_ttl":   "2m",
			"allowed_client_ids": "allowed-test-role",
		},
		Storage: storage,
	})
	expectError(t, resp, err)

	// Create a client that depends on test key
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/client/test-client",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"key":          "test-key",
			"id_token_ttl": "4m",
		},
		Storage: storage,
	})
	expectSuccess(t, resp, err)

	// Update test key "test-key" -- should fail since id_token_ttl is greater than 2m
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/key/test-key",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"verification_ttl": "2m",
		},
		Storage: storage,
	})
	expectError(t, resp, err)
}

// TestOIDC_Path_ListKey tests the List operation for keys
func TestOIDC_Path_ListKey(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	// Prepare two keys, test-key1 and test-key2
	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/key/test-key1",
		Operation: logical.CreateOperation,
		Storage:   storage,
	})

	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/key/test-key2",
		Operation: logical.CreateOperation,
		Storage:   storage,
	})

	// list keys
	respListKey, listErr := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/key",
		Operation: logical.ListOperation,
		Storage:   storage,
	})
	expectSuccess(t, respListKey, listErr)

	// validate list response
	expectedStrings := map[string]interface{}{"test-key1": true, "test-key2": true}
	expectStrings(t, respListKey.Data["keys"].([]string), expectedStrings)

	// delete test-key2
	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/key/test-key2",
		Operation: logical.DeleteOperation,
		Storage:   storage,
	})

	// list keyes again and validate response
	respListKeyAfterDelete, listErrAfterDelete := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/key",
		Operation: logical.ListOperation,
		Storage:   storage,
	})
	expectSuccess(t, respListKeyAfterDelete, listErrAfterDelete)

	// validate list response
	delete(expectedStrings, "test-key2")
	expectStrings(t, respListKeyAfterDelete.Data["keys"].([]string), expectedStrings)
}

// TestOIDC_Path_OIDCKey_DeleteWithExistingClient tests that a key cannot be
// deleted if it is referenced by an existing client
func TestOIDC_Path_OIDCKey_DeleteWithExistingClient(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	// Prepare test key test-key
	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/key/test-key",
		Operation: logical.CreateOperation,
		Storage:   storage,
	})

	// Create a test client "test-client" -- should succeed
	resp, err := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/client/test-client",
		Operation: logical.CreateOperation,
		Storage:   storage,
		Data: map[string]interface{}{
			"key": "test-key",
		},
	})
	expectSuccess(t, resp, err)

	// Delete test key "test-key" -- should fail
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/key/test-key",
		Operation: logical.DeleteOperation,
		Storage:   storage,
	})
	expectError(t, resp, err)
}

// TestOIDC_PublicKeys_NoRole tests that public keys are not returned by the
// oidc/.well-known/keys endpoint when they are not associated with a role
func TestOIDC_PublicKeys_NoRole(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	s := &logical.InmemStorage{}

	// Create a test key "test-key"
	resp, err := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/key/test-key",
		Operation: logical.CreateOperation,
		Storage:   s,
	})
	expectSuccess(t, resp, err)

	// .well-known/keys should contain 0 public keys
	assertPublicKeyCount(t, ctx, s, c, 0)
}

func assertPublicKeyCount(t *testing.T, ctx context.Context, s logical.Storage, c *Core, keyCount int) {
	t.Helper()

	// .well-known/keys should contain keyCount public keys
	resp, err := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/.well-known/keys",
		Operation: logical.ReadOperation,
		Storage:   s,
	})
	expectSuccess(t, resp, err)

	assertRespPublicKeyCount(t, resp, keyCount)
}

func assertRespPublicKeyCount(t *testing.T, resp *logical.Response, keyCount int) {
	t.Helper()

	// parse response
	responseJWKS := &jose.JSONWebKeySet{}
	json.Unmarshal(resp.Data["http_raw_body"].([]byte), responseJWKS)
	if len(responseJWKS.Keys) != keyCount {
		t.Fatalf("expected %d public keys but instead got %d", keyCount, len(responseJWKS.Keys))
	}
}

// TestOIDC_PublicKeys tests that public keys are updated by
// key creation, rotation, and deletion
func TestOIDC_PublicKeys(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	// Create a test key "test-key"
	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/key/test-key",
		Operation: logical.CreateOperation,
		Storage:   storage,
	})

	// Create a test role "test-role"
	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/role/test-role",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"key": "test-key",
		},
		Storage: storage,
	})

	// .well-known/keys should contain 2 public keys
	assertPublicKeyCount(t, ctx, storage, c, 2)

	// rotate test-key a few times, each rotate should increase the length of public keys returned
	// by the .well-known endpoint
	resp, err := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/key/test-key/rotate",
		Operation: logical.UpdateOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)
	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/key/test-key/rotate",
		Operation: logical.UpdateOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)

	// .well-known/keys should contain 4 public keys
	assertPublicKeyCount(t, ctx, storage, c, 4)

	// create another named key "test-key2"
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/key/test-key2",
		Operation: logical.CreateOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)

	// Create a test role "test-role2"
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/role/test-role2",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"key": "test-key2",
		},
		Storage: storage,
	})
	expectSuccess(t, resp, err)
	// .well-known/keys should contain 6 public keys
	assertPublicKeyCount(t, ctx, storage, c, 6)

	// delete test role that references "test-key"
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/role/test-role",
		Operation: logical.DeleteOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)
	// delete test key
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/key/test-key",
		Operation: logical.DeleteOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)

	// .well-known/keys should contain 2 public keys, all of the public keys
	// from named key "test-key" should have been deleted
	assertPublicKeyCount(t, ctx, storage, c, 2)
}

// TestOIDC_PublicKeys tests that public keys are updated by
// key creation, rotation, and deletion
func TestOIDC_SharedPublicKeysByRoles(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	// Create a test key "test-key"
	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/key/test-key",
		Operation: logical.CreateOperation,
		Storage:   storage,
	})

	// Create a test role "test-role"
	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/role/test-role",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"key": "test-key",
		},
		Storage: storage,
	})

	// Create a test role "test-role2"
	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/role/test-role2",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"key": "test-key",
		},
		Storage: storage,
	})

	// Create a test role "test-role3"
	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/role/test-role3",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"key": "test-key",
		},
		Storage: storage,
	})

	// .well-known/keys should contain 2 public keys
	assertPublicKeyCount(t, ctx, storage, c, 2)
}

// TestOIDC_SignIDToken tests acquiring a signed token and verifying the public portion
// of the signing key
func TestOIDC_SignIDToken(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	// Create and load an entity, an entity is required to generate an ID token
	testEntity := &identity.Entity{
		Name:      "test-entity-name",
		ID:        "test-entity-id",
		BucketKey: "test-entity-bucket-key",
	}

	txn := c.identityStore.db.Txn(true)
	defer txn.Abort()
	err := c.identityStore.upsertEntityInTxn(ctx, txn, testEntity, nil, true)
	if err != nil {
		t.Fatal(err)
	}
	txn.Commit()

	// Create a test key "test-key"
	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/key/test-key",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"allowed_client_ids": "*",
		},
		Storage: storage,
	})

	// Create a test role "test-role" -- expect no warning
	resp, err := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/role/test-role",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"key": "test-key",
		},
		Storage: storage,
	})
	expectSuccess(t, resp, err)
	if resp != nil {
		t.Fatalf("was expecting a nil response but instead got: %#v", resp)
	}

	// Determine test-role's client_id
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/role/test-role",
		Operation: logical.ReadOperation,
		Data: map[string]interface{}{
			"allowed_client_ids": "",
		},
		Storage: storage,
	})
	expectSuccess(t, resp, err)
	clientID := resp.Data["client_id"].(string)

	// remove test-role as an allowed role from test-key
	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/key/test-key",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"allowed_client_ids": "",
		},
		Storage: storage,
	})

	// Generate a token against the role "test-role" -- should fail
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/token/test-role",
		Operation: logical.ReadOperation,
		Storage:   storage,
		EntityID:  "test-entity-id",
	})
	expectError(t, resp, err)
	// validate error message
	expectedStrings := map[string]interface{}{
		"the key \"test-key\" does not list the client ID of the role \"test-role\" as an allowed client ID": true,
	}
	expectStrings(t, []string{resp.Data["error"].(string)}, expectedStrings)

	// add test-role as an allowed role from test-key
	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/key/test-key",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"allowed_client_ids": clientID,
		},
		Storage: storage,
	})

	// Generate a token against the role "test-role" -- should succeed
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/token/test-role",
		Operation: logical.ReadOperation,
		Storage:   storage,
		EntityID:  "test-entity-id",
	})
	expectSuccess(t, resp, err)
	parsedToken, err := jwt.ParseSigned(resp.Data["token"].(string))
	if err != nil {
		t.Fatalf("error parsing token: %s", err.Error())
	}

	// Acquire the public parts of the key that signed parsedToken
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/.well-known/keys",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	responseJWKS := &jose.JSONWebKeySet{}
	json.Unmarshal(resp.Data["http_raw_body"].([]byte), responseJWKS)

	keyCount := len(responseJWKS.Keys)
	errorCount := 0
	for _, key := range responseJWKS.Keys {
		// Validate the signature
		claims := &jwt.Claims{}
		if err := parsedToken.Claims(key, claims); err != nil {
			t.Logf("unable to validate signed token, err:\n%#v", err)
			errorCount += 1
		}
	}
	if errorCount == keyCount {
		t.Fatalf("unable to validate signed token with any of the .well-known keys")
	}
}

// TestOIDC_SignIDToken_NilSigningKey tests that an error is returned when
// attempting to sign an ID token with a nil signing key
func TestOIDC_SignIDToken_NilSigningKey(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)

	// Create and load an entity, an entity is required to generate an ID token
	testEntity := &identity.Entity{
		Name:      "test-entity-name",
		ID:        "test-entity-id",
		BucketKey: "test-entity-bucket-key",
	}

	txn := c.identityStore.db.Txn(true)
	defer txn.Abort()
	err := c.identityStore.upsertEntityInTxn(ctx, txn, testEntity, nil, true)
	if err != nil {
		t.Fatal(err)
	}
	txn.Commit()

	// Create a test key "test-key" with a nil SigningKey
	namedKey := &namedKey{
		name:             "test-key",
		AllowedClientIDs: []string{"*"},
		Algorithm:        "RS256",
		VerificationTTL:  60 * time.Second,
		RotationPeriod:   60 * time.Second,
		KeyRing:          nil,
		SigningKey:       nil,
		NextSigningKey:   nil,
		NextRotation:     time.Now(),
	}
	s := c.router.MatchingStorageByAPIPath(ctx, "identity/oidc")
	if err := namedKey.generateAndSetNextKey(ctx, hclog.NewNullLogger(), s); err != nil {
		t.Fatalf("failed to set next signing key")
	}
	// Store namedKey
	entry, _ := logical.StorageEntryJSON(namedKeyConfigPath+namedKey.name, namedKey)
	if err := s.Put(ctx, entry); err != nil {
		t.Fatalf("writing to in mem storage failed")
	}

	// Create a test role "test-role" -- expect no warning
	resp, err := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/role/test-role",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"key": "test-key",
			"ttl": "1m",
		},
		Storage: s,
	})
	expectSuccess(t, resp, err)
	if resp != nil {
		t.Fatalf("was expecting a nil response but instead got: %#v", resp)
	}

	// Generate a token against the role "test-role" -- should fail
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/token/test-role",
		Operation: logical.ReadOperation,
		Storage:   s,
		EntityID:  "test-entity-id",
	})
	expectError(t, resp, err)
	// validate error message
	expectedStrings := map[string]interface{}{
		"error signing OIDC token: signing key is nil; rotate the key and try again": true,
	}
	expectStrings(t, []string{err.Error()}, expectedStrings)
}

func testNamedKey(name string) *namedKey {
	return &namedKey{
		name:            name,
		Algorithm:       "RS256",
		VerificationTTL: 1 * time.Second,
		RotationPeriod:  2 * time.Second,
		KeyRing:         nil,
		SigningKey:      nil,
		NextSigningKey:  nil,
		NextRotation:    time.Now(),
	}
}

// TestOIDC_PeriodicFunc tests timing logic for running key
// rotations and expiration actions.
func TestOIDC_PeriodicFunc(t *testing.T) {
	type testCase struct {
		minKeyRingLen int
		maxKeyRingLen int
	}
	testSets := []struct {
		namedKey          *namedKey
		setSigningKey     bool
		setNextSigningKey bool
		testCases         []testCase
	}{
		{
			namedKey:          testNamedKey("test-key"),
			setSigningKey:     true,
			setNextSigningKey: true,
			testCases: []testCase{
				// we must always have at least 2 keys and at most 3 through rotation cycles, since
				// each rotation cycle results in a key going in/out of its verification_ttl period.
				{2, 3},
				{2, 3},
				{2, 3},
				{2, 3},
			},
		},
		{
			// don't set SigningKey to ensure its non-existence can be handled
			namedKey:          testNamedKey("test-key-nil-signing-key"),
			setSigningKey:     false,
			setNextSigningKey: true,
			testCases: []testCase{
				{1, 1},

				// key counts jump from 1 to 2 because the next signing key becomes
				// the signing key, and no key is in its verification_ttl period
				{2, 2},
			},
		},
		{
			// don't set NextSigningKey to ensure its non-existence can be handled
			namedKey:          testNamedKey("test-key-nil-next-signing-key"),
			setSigningKey:     true,
			setNextSigningKey: false,
			testCases: []testCase{
				{1, 1},

				// max key counts jump from 1 to 3 because the original signing
				// key could still be within its verification_ttl period
				{2, 3},
			},
		},
		{
			// don't set keys to ensure non-existence can be handled
			namedKey:          testNamedKey("test-key-nil-signing-and-next-signing-key"),
			setSigningKey:     false,
			setNextSigningKey: false,
			testCases: []testCase{
				{0, 0},

				// First rotation populates both current/next signing keys
				{2, 2},
			},
		},
	}

	for _, testSet := range testSets {
		testSet := testSet
		t.Run(testSet.namedKey.name, func(t *testing.T) {
			t.Parallel()

			c, _, _ := TestCoreUnsealed(t)
			ctx := namespace.RootContext(nil)
			storage := c.router.MatchingStorageByAPIPath(ctx, "identity/oidc")

			// Stop the core's rollback manager so that periodic function testing
			// doesn't race with the rollback manager.
			c.rollback.StopTicker()

			// Generate current and next keys as needed by the test. This ensures
			// we can rotate a key when either are unset.
			if testSet.setSigningKey {
				require.NoError(t, testSet.namedKey.generateAndSetKey(ctx, hclog.NewNullLogger(), storage))
			}
			if testSet.setNextSigningKey {
				require.NoError(t, testSet.namedKey.generateAndSetNextKey(ctx, hclog.NewNullLogger(), storage))
			}
			testSet.namedKey.NextRotation = time.Now().Add(testSet.namedKey.RotationPeriod)

			// Store the named key so it can be rotated by the periodic func
			entry, _ := logical.StorageEntryJSON(namedKeyConfigPath+testSet.namedKey.name, testSet.namedKey)
			require.NoError(t, storage.Put(ctx, entry))
			t.Cleanup(func() {
				require.NoError(t, storage.Delete(ctx, namedKeyConfigPath+testSet.namedKey.name))
			})

			// Manually execute the periodic func to rotate keys and collect
			// both the key ring and public keys.
			namedKeySamples := make([]*logical.StorageEntry, len(testSet.testCases))
			publicKeysSamples := make([][]string, len(testSet.testCases))
			for i := range testSet.testCases {
				c.identityStore.oidcPeriodicFunc(ctx)
				namedKeyEntry, _ := storage.Get(ctx, namedKeyConfigPath+testSet.namedKey.name)
				publicKeysEntry, _ := storage.List(ctx, publicKeysConfigPath)
				namedKeySamples[i] = namedKeyEntry
				publicKeysSamples[i] = publicKeysEntry

				// sleep until we are in the next cycle - where a next run will happen
				v, _, _ := c.identityStore.oidcCache.Get(noNamespace, "nextRun")
				nextRun := v.(time.Time)
				now := time.Now()
				diff := nextRun.Sub(now)
				if now.Before(nextRun) {
					time.Sleep(diff + 100*time.Millisecond)
				}
			}

			// Assert that the key lengths through each rotation match expectations
			for i, tc := range testSet.testCases {
				require.NoError(t, namedKeySamples[i].DecodeJSON(&testSet.namedKey))

				actualKeyRingLen := len(testSet.namedKey.KeyRing)
				assert.LessOrEqual(t, actualKeyRingLen, tc.maxKeyRingLen,
					"test case index %d: key ring length must be at most %d", i, tc.maxKeyRingLen)
				assert.GreaterOrEqual(t, actualKeyRingLen, tc.minKeyRingLen,
					"test case index %d: key ring length must be at least %d", i, tc.minKeyRingLen)

				actualPubKeysLen := len(publicKeysSamples[i])
				assert.LessOrEqual(t, actualPubKeysLen, tc.maxKeyRingLen,
					"test case index %d: public key ring length must be at most %d", i, tc.maxKeyRingLen)
				assert.GreaterOrEqual(t, actualPubKeysLen, tc.minKeyRingLen,
					"test case index %d: public key ring length must be at least %d", i, tc.minKeyRingLen)
			}
		})
	}
}

// TestOIDC_Config tests CRUD operations for configuring the OIDC backend
func TestOIDC_Config(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	testIssuer := "https://example.com:1234"

	// Read Config - expect defaults
	resp, err := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/config",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)
	// issuer should not be set
	if resp.Data["issuer"].(string) != "" {
		t.Fatalf("Expected issuer to not be set but found %q instead", resp.Data["issuer"].(string))
	}

	// Update Config
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/config",
		Operation: logical.UpdateOperation,
		Storage:   storage,
		Data: map[string]interface{}{
			"issuer": testIssuer,
		},
	})
	expectSuccess(t, resp, err)

	// Read Config - expect updated issuer value
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/config",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)
	// issuer should be set
	if resp.Data["issuer"].(string) != testIssuer {
		t.Fatalf("Expected issuer to be %q but found %q instead", testIssuer, resp.Data["issuer"].(string))
	}

	// Test bad issuers
	for _, iss := range []string{"asldfk", "ftp://a.com", "a.com", "http://a.com/", "https://a.com/foo", "http:://a.com"} {
		resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
			Path:      "oidc/config",
			Operation: logical.UpdateOperation,
			Storage:   storage,
			Data: map[string]interface{}{
				"issuer": iss,
			},
		})
		if err != nil {
			t.Fatal(err)
		}
		if resp == nil || !resp.IsError() {
			t.Fatalf("Expected issuer %q to fail but it succeeded.", iss)
		}

	}
}

// TestOIDC_pathOIDCKeyExistenceCheck tests pathOIDCKeyExistenceCheck
func TestOIDC_pathOIDCKeyExistenceCheck(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	keyName := "test"

	// Expect nil with empty storage
	exists, err := c.identityStore.pathOIDCKeyExistenceCheck(
		ctx,
		&logical.Request{
			Storage: storage,
		},
		&framework.FieldData{
			Raw: map[string]interface{}{"name": keyName},
			Schema: map[string]*framework.FieldSchema{
				"name": {
					Type: framework.TypeString,
				},
			},
		},
	)
	if err != nil {
		t.Fatalf("Error during existence check on an expected nil entry, err:\n%#v", err)
	}
	if exists {
		t.Fatalf("Expected existence check to return false but instead returned: %t", exists)
	}

	// Populate storage with a namedKey
	namedKey := &namedKey{}
	entry, _ := logical.StorageEntryJSON(namedKeyConfigPath+keyName, namedKey)
	if err := storage.Put(ctx, entry); err != nil {
		t.Fatalf("writing to in mem storage failed")
	}

	// Expect true with a populated storage
	exists, err = c.identityStore.pathOIDCKeyExistenceCheck(
		ctx,
		&logical.Request{
			Storage: storage,
		},
		&framework.FieldData{
			Raw: map[string]interface{}{"name": keyName},
			Schema: map[string]*framework.FieldSchema{
				"name": {
					Type: framework.TypeString,
				},
			},
		},
	)
	if err != nil {
		t.Fatalf("Error during existence check on an expected nil entry, err:\n%#v", err)
	}
	if !exists {
		t.Fatalf("Expected existence check to return true but instead returned: %t", exists)
	}
}

// TestOIDC_pathOIDCRoleExistenceCheck tests pathOIDCRoleExistenceCheck
func TestOIDC_pathOIDCRoleExistenceCheck(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	roleName := "test"

	// Expect nil with empty storage
	exists, err := c.identityStore.pathOIDCRoleExistenceCheck(
		ctx,
		&logical.Request{
			Storage: storage,
		},
		&framework.FieldData{
			Raw: map[string]interface{}{"name": roleName},
			Schema: map[string]*framework.FieldSchema{
				"name": {
					Type: framework.TypeString,
				},
			},
		},
	)
	if err != nil {
		t.Fatalf("Error during existence check on an expected nil entry, err:\n%#v", err)
	}
	if exists {
		t.Fatalf("Expected existence check to return false but instead returned: %t", exists)
	}

	// Populate storage with a role
	role := &role{}
	entry, _ := logical.StorageEntryJSON(roleConfigPath+roleName, role)
	if err := storage.Put(ctx, entry); err != nil {
		t.Fatalf("writing to in mem storage failed")
	}

	// Expect true with a populated storage
	exists, err = c.identityStore.pathOIDCRoleExistenceCheck(
		ctx,
		&logical.Request{
			Storage: storage,
		},
		&framework.FieldData{
			Raw: map[string]interface{}{"name": roleName},
			Schema: map[string]*framework.FieldSchema{
				"name": {
					Type: framework.TypeString,
				},
			},
		},
	)
	if err != nil {
		t.Fatalf("Error during existence check on an expected nil entry, err:\n%#v", err)
	}
	if !exists {
		t.Fatalf("Expected existence check to return true but instead returned: %t", exists)
	}
}

// TestOIDC_Path_OpenIDConfig tests read operations for the openid-configuration path
func TestOIDC_Path_OpenIDConfig(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	// Expect defaults from .well-known/openid-configuration
	resp, err := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/.well-known/openid-configuration",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)
	// Validate configurable parts - for now just issuer
	discoveryResp := &discovery{}
	json.Unmarshal(resp.Data["http_raw_body"].([]byte), discoveryResp)
	expected := "/v1/identity/oidc"
	if discoveryResp.Issuer != expected {
		t.Fatalf("Expected Issuer path to be %q but found %q instead", expected, discoveryResp.Issuer)
	}

	// Update issuer config
	testIssuer := "https://example.com:1234"
	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/config",
		Operation: logical.UpdateOperation,
		Storage:   storage,
		Data: map[string]interface{}{
			"issuer": testIssuer,
		},
	})

	// Expect updates from .well-known/openid-configuration
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/.well-known/openid-configuration",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)
	// Validate configurable parts - for now just issuer
	json.Unmarshal(resp.Data["http_raw_body"].([]byte), discoveryResp)
	expected = "https://example.com:1234/v1/identity/oidc"
	if discoveryResp.Issuer != expected {
		t.Fatalf("Expected Issuer path to be %q but found %q instead", expected, discoveryResp.Issuer)
	}
}

// TestOIDC_Path_Introspect tests update operations on the introspect path
func TestOIDC_Path_Introspect(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	// Expect active false and an error from a malformed token
	resp, err := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/introspect/",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"token": "not-a-valid-token",
		},
		Storage: storage,
	})
	expectSuccess(t, resp, err)
	type introspectResponse struct {
		Active bool   `json:"active"`
		Error  string `json:"error"`
	}
	iresp := &introspectResponse{}
	json.Unmarshal(resp.Data["http_raw_body"].([]byte), iresp)
	if iresp.Active {
		t.Fatalf("expected active state of a malformed token to be false but what was found to be: %t", iresp.Active)
	}
	if iresp.Error == "" {
		t.Fatalf("expected a malformed token to return an error message but instead returned %q", iresp.Error)
	}

	// Populate backend with a valid token ---
	// Create and load an entity, an entity is required to generate an ID token
	testEntity := &identity.Entity{
		Name:      "test-entity-name",
		ID:        "test-entity-id",
		BucketKey: "test-entity-bucket-key",
	}

	txn := c.identityStore.db.Txn(true)
	defer txn.Abort()
	err = c.identityStore.upsertEntityInTxn(ctx, txn, testEntity, nil, true)
	if err != nil {
		t.Fatal(err)
	}
	txn.Commit()

	for _, alg := range []string{"RS256", "RS384", "RS512", "ES256", "ES384", "ES512", "EdDSA"} {
		key := "test-key-" + alg
		role := "test-role-" + alg

		// Create a test key "test-key"
		resp, err := c.identityStore.HandleRequest(ctx, &logical.Request{
			Path:      "oidc/key/" + key,
			Operation: logical.CreateOperation,
			Storage:   storage,
			Data: map[string]interface{}{
				"algorithm":          alg,
				"allowed_client_ids": "*",
			},
		})
		expectSuccess(t, resp, err)

		// Create a test role "test-role"
		resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
			Path:      "oidc/role/" + role,
			Operation: logical.CreateOperation,
			Data: map[string]interface{}{
				"key": key,
			},
			Storage: storage,
		})
		expectSuccess(t, resp, err)

		// Generate a token against the role "test-role" -- should succeed
		resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
			Path:      "oidc/token/" + role,
			Operation: logical.ReadOperation,
			Storage:   storage,
			EntityID:  "test-entity-id",
		})
		expectSuccess(t, resp, err)

		validToken := resp.Data["token"].(string)

		//	Expect active true and no error from a valid token
		resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
			Path:      "oidc/introspect/",
			Operation: logical.UpdateOperation,
			Data: map[string]interface{}{
				"token": validToken,
			},
			Storage: storage,
		})
		expectSuccess(t, resp, err)
		iresp2 := &introspectResponse{}
		json.Unmarshal(resp.Data["http_raw_body"].([]byte), iresp2)
		if !iresp2.Active {
			t.Fatalf("expected active state of a valid token to be true but what was found to be: %t", iresp2.Active)
		}
		if iresp2.Error != "" {
			t.Fatalf("expected a valid token to return an empty error message but instead got %q", iresp.Error)
		}
	}
}

func TestOIDC_isTargetNamespacedKey(t *testing.T) {
	tests := []struct {
		nsTargets []string
		nskey     string
		expected  bool
	}{
		{[]string{"nsid"}, "v0:nsid:key", true},
		{[]string{"nsid"}, "v0:nsid:", true},
		{[]string{"nsid"}, "v0:nsid", false},
		{[]string{"nsid"}, "v0:", false},
		{[]string{"nsid"}, "v0", false},
		{[]string{"nsid"}, "", false},
		{[]string{"nsid1"}, "v0:nsid2:key", false},
		{[]string{"nsid1"}, "nsid1:nsid2:nsid1", false},
		{[]string{"nsid1"}, "nsid1:nsid1:nsid1", true},
		{[]string{"nsid"}, "nsid:nsid:nsid:nsid:nsid:nsid", true},
		{[]string{"nsid"}, ":::", false},
		{[]string{""}, ":::", true}, // "" is a valid key for cache.Set/Get
		{[]string{"nsid1"}, "nsid0:nsid1:nsid0:nsid1:nsid0:nsid1", true},
		{[]string{"nsid0"}, "nsid0:nsid1:nsid0:nsid1:nsid0:nsid1", false},
		{[]string{"nsid0", "nsid1"}, "v0:nsid2:key", false},
		{[]string{"nsid0", "nsid1", "nsid2", "nsid3", "nsid4"}, "v0:nsid3:key", true},
		{[]string{"nsid0", "nsid1", "nsid2", "nsid3", "nsid4"}, "nsid0:nsid1:nsid2:nsid3:nsid4:nsid5", true},
		{[]string{"nsid0", "nsid1", "nsid2", "nsid3", "nsid4"}, "nsid4:nsid5:nsid6:nsid7:nsid8:nsid9", false},
		{[]string{"nsid0", "nsid0", "nsid0", "nsid0", "nsid0"}, "nsid0:nsid0:nsid0:nsid0:nsid0:nsid0", true},
		{[]string{"nsid1", "nsid1", "nsid2", "nsid2"}, "nsid0:nsid0:nsid0:nsid0:nsid0:nsid0", false},
		{[]string{"nsid1", "nsid1", "nsid2", "nsid2"}, "nsid0:nsid0:nsid0:nsid0:nsid0:nsid0", false},
	}

	for _, test := range tests {
		actual := isTargetNamespacedKey(test.nskey, test.nsTargets)
		if test.expected != actual {
			t.Fatalf("expected %t but got %t for nstargets: %q and nskey: %q", test.expected, actual, test.nsTargets, test.nskey)
		}
	}
}

func TestOIDC_Flush(t *testing.T) {
	c := newOIDCCache(gocache.NoExpiration, gocache.NoExpiration)
	ns := []*namespace.Namespace{
		noNamespace, // ns[0] is nilNamespace
		{ID: "ns1"},
		{ID: "ns2"},
	}

	// populateNs populates cache by ns with some data
	populateNs := func() {
		for i := range ns {
			for _, val := range []string{"keyA", "keyB", "keyC"} {
				if err := c.SetDefault(ns[i], val, struct{}{}); err != nil {
					t.Fatal(err)
				}
			}
		}
	}

	// validate verifies that cache items exist or do not exist based on their namespaced key
	verify := func(items map[string]gocache.Item, expect, doNotExpect []*namespace.Namespace) {
		for _, expectNs := range expect {
			found := false
			for i := range items {
				if isTargetNamespacedKey(i, []string{expectNs.ID}) {
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("Expected cache to contain an entry with a namespaced key for namespace: %q but did not find one", expectNs.ID)
			}
		}

		for _, doNotExpectNs := range doNotExpect {
			for i := range items {
				if isTargetNamespacedKey(i, []string{doNotExpectNs.ID}) {
					t.Fatalf("Did not expect cache to contain an entry with a namespaced key for namespace: %q but found the key: %q", doNotExpectNs.ID, i)
				}
			}
		}
	}

	// flushing ns1 should flush ns1 and nilNamespace but not ns2
	populateNs()
	if err := c.Flush(ns[1]); err != nil {
		t.Fatal(err)
	}
	items := c.c.Items()
	verify(items, []*namespace.Namespace{ns[2]}, []*namespace.Namespace{ns[0], ns[1]})

	// flushing nilNamespace should flush nilNamespace but not ns1 or ns2
	populateNs()
	if err := c.Flush(ns[0]); err != nil {
		t.Fatal(err)
	}
	items = c.c.Items()
	verify(items, []*namespace.Namespace{ns[1], ns[2]}, []*namespace.Namespace{ns[0]})
}

func TestOIDC_CacheNamespaceNilCheck(t *testing.T) {
	cache := newOIDCCache(gocache.NoExpiration, gocache.NoExpiration)

	if _, _, err := cache.Get(nil, "foo"); err == nil {
		t.Fatal("expected error, got nil")
	}

	if err := cache.SetDefault(nil, "foo", 42); err == nil {
		t.Fatal("expected error, got nil")
	}

	if err := cache.Flush(nil); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestOIDC_GetKeysCacheControlHeader(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)

	// get default value
	header, err := c.identityStore.getKeysCacheControlHeader()
	if err != nil {
		t.Fatalf("expected success, got error:\n%v", err)
	}

	expectedHeader := ""
	if header != expectedHeader {
		t.Fatalf("expected %s, got %s", expectedHeader, header)
	}

	// set nextRun
	nextRun := time.Now().Add(24 * time.Hour)
	if err = c.identityStore.oidcCache.SetDefault(noNamespace, "nextRun", nextRun); err != nil {
		t.Fatal(err)
	}

	header, err = c.identityStore.getKeysCacheControlHeader()
	if err != nil {
		t.Fatalf("expected success, got error:\n%v", err)
	}

	expectedNextRun := "max-age=86400"
	if header != expectedNextRun {
		t.Fatalf("expected %s, got %s", expectedNextRun, header)
	}

	// set jwksCacheControlMaxAge
	durationSeconds := 60
	jwksCacheControlMaxAge := time.Duration(durationSeconds) * time.Second
	if err = c.identityStore.oidcCache.SetDefault(noNamespace, "jwksCacheControlMaxAge", jwksCacheControlMaxAge); err != nil {
		t.Fatal(err)
	}

	header, err = c.identityStore.getKeysCacheControlHeader()
	if err != nil {
		t.Fatalf("expected success, got error:\n%v", err)
	}

	if header == "" {
		t.Fatalf("expected header to be set, got %s", header)
	}

	maxAgeValue := strings.Split(header, "=")[1]
	headerVal, err := strconv.Atoi(maxAgeValue)
	if err != nil {
		t.Fatal(err)
	}
	// headerVal will be a random value between 0 and jwksCacheControlMaxAge (in seconds)
	if headerVal > durationSeconds {
		t.Logf("jwksCacheControlMaxAge: %d", int(jwksCacheControlMaxAge))
		t.Fatalf("unexpected header value, got %d expected less than %d", headerVal, durationSeconds)
	}
}

// some helpers
func expectSuccess(t *testing.T, resp *logical.Response, err error) {
	t.Helper()
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("expected success but got error:\n%v\nresp: %#v", err, resp)
	}
}

func expectError(t *testing.T, resp *logical.Response, err error) {
	t.Helper()
	if err == nil {
		if resp == nil || !resp.IsError() {
			t.Fatalf("expected error but got success; error:\n%v\nresp: %#v", err, resp)
		}
	}
}

// expectString fails unless every string in actualStrings is also included in expectedStrings and
// the length of actualStrings and expectedStrings are the same
func expectStrings(t *testing.T, actualStrings []string, expectedStrings map[string]interface{}) {
	t.Helper()
	if len(actualStrings) != len(expectedStrings) {
		t.Fatalf("expectStrings mismatch:\nactual strings:\n%#v\nexpected strings:\n%#v\n", actualStrings, expectedStrings)
	}
	for _, actualString := range actualStrings {
		_, ok := expectedStrings[actualString]
		if !ok {
			t.Fatalf("the string %q was not expected", actualString)
		}
	}
}

// Test_oidcConfig_fullIssuer tests that the full issuer matches expectations
// given different issuer bases and children.
func Test_oidcConfig_fullIssuer(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	tests := []struct {
		name    string
		issuer  string
		child   string
		want    string
		wantErr bool
	}{
		{
			name:   "issuer with valid empty child",
			issuer: "https://vault.dev",
			child:  baseIdentityTokenIssuer,
			want:   fmt.Sprintf("https://vault.dev/v1/%s", issuerPath),
		},
		{
			name:    "issuer with invalid child",
			issuer:  "http://127.0.0.1:8200",
			child:   "invalid",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := c.identityStore.HandleRequest(ctx, &logical.Request{
				Path:      "oidc/config",
				Operation: logical.UpdateOperation,
				Storage:   storage,
				Data: map[string]interface{}{
					"issuer": tt.issuer,
				},
			})
			expectSuccess(t, resp, err)

			config, err := c.identityStore.getOIDCConfig(ctx, storage)
			require.NoError(t, err)

			got, err := config.fullIssuer(tt.child)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equalf(t, tt.want, got, "fullIssuer(%v)", tt.child)
		})
	}
}

// Test_validChildIssuer tests that only valid child issuers are accepted.
func Test_validChildIssuer(t *testing.T) {
	tests := []struct {
		name  string
		child string
		want  bool
	}{
		{
			name:  "valid child issuer",
			child: baseIdentityTokenIssuer,
			want:  true,
		},
		{
			name:  "invalid child issuer",
			child: "test",
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equalf(t, tt.want, validChildIssuer(tt.child), "validChildIssuer(%v)", tt.child)
		})
	}
}

// Test_optionalChildIssuerRegex tests that the regex returned from
// optionalChildIssuerRegex produces the expected captures given different
// input paths.
func Test_optionalChildIssuerRegex(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		path     string
		captures map[string]string
	}{
		{
			name:     "valid match with capture",
			pattern:  "oidc" + optionalChildIssuerRegex("child") + "/.well-known/keys",
			path:     "oidc/test/.well-known/keys",
			captures: map[string]string{"child": "test"},
		},
		{
			name:     "valid match with capture name, segment, and path change",
			pattern:  "oidc" + optionalChildIssuerRegex("name") + "/.well-known/openid-configuration",
			path:     "oidc/test/.well-known/openid-configuration",
			captures: map[string]string{"name": "test"},
		},
		{
			name:     "valid match with empty capture",
			pattern:  "oidc" + optionalChildIssuerRegex("child") + "/.well-known/keys",
			path:     "oidc/.well-known/keys",
			captures: map[string]string{"child": ""},
		},
		{
			name:     "invalid match with multiple path segments",
			pattern:  "oidc" + optionalChildIssuerRegex("child") + "/.well-known/keys",
			path:     "oidc/test/invalid/.well-known/keys",
			captures: map[string]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			re := regexp.MustCompile(tt.pattern)
			matches := re.FindStringSubmatch(tt.path)
			actualCaptures := make(map[string]string)
			for i, name := range re.SubexpNames() {
				if name != "" && i < len(matches) {
					actualCaptures[name] = matches[i]
				}
			}
			require.Equal(t, tt.captures, actualCaptures)
		})
	}
}

func createMountEntryWithKey(t *testing.T, ctx context.Context, sys *SystemBackend, mountPrefix, mountType, key string) {
	t.Helper()

	resp, err := sys.HandleRequest(ctx, &logical.Request{
		Path:      mountPrefix + mountType,
		Operation: logical.UpdateOperation,
		Storage:   new(logical.InmemStorage),
		Data: map[string]interface{}{
			"type": strings.TrimSuffix(mountType, "/"),
			"config": map[string]interface{}{
				"identity_token_key": key,
			},
		},
	})
	expectSuccess(t, resp, err)
}
