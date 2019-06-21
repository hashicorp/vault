package vault

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"testing"
	"time"

	"github.com/go-test/deep"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

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

	// Read "test-role1"
	respReadTestRole1, err1 := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/role/test-role1",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	expectSuccess(t, respReadTestRole1, err1)

	// Create a test role "test-role2" witn an invalid key -- should fail
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/role/test-role2",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"key": "a-key-that-does-not-exist",
		},
		Storage: storage,
	})
	expectError(t, resp, err)

	// Update "test-role1" with valid parameters -- should succeed
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/role/test-role1",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"template": "{\"some-key\":\"some-value\"}",
			"ttl":      "2h",
		},
		Storage: storage,
	})
	expectSuccess(t, resp, err)

	// Read "test-role1" again
	respReadTestRole1AfterUpdate, err2 := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/role/test-role1",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	expectSuccess(t, respReadTestRole1AfterUpdate, err2)

	// Compare response for "test-role1" before and after it was updated
	expectedDiff := map[string]interface{}{
		"Data.map[template]:  != {\"some-key\":\"some-value\"}": true,
		"Data.map[ttl]: 86400 != 7200":                          true, // 24h to 2h
	}
	diff := deep.Equal(respReadTestRole1, respReadTestRole1AfterUpdate)
	expectStrings(t, diff, expectedDiff)

	// Delete "test-role1"
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/role/test-role1",
		Operation: logical.DeleteOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)

	// Read "test-role1" again
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
		Storage:   storage,
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

// TestOIDC_Path_OIDCKeyKey tests CRUD operations for keys
func TestOIDC_Path_OIDCKeyKey(t *testing.T) {
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

	// Read "test-key"
	respReadTestKey, err1 := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/key/test-key",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	expectSuccess(t, respReadTestKey, err1)

	// Update "test-key" -- should succeed
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/key/test-key",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"rotation_period":  "10m",
			"verification_ttl": "1h",
		},
		Storage: storage,
	})
	expectSuccess(t, resp, err)

	// Read "test-key" again
	respReadTestKeyAfterUpdate, err2 := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/key/test-key",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	expectSuccess(t, respReadTestKeyAfterUpdate, err2)

	// Compare response for "test-key" before and after it was updated
	expectedDiff := map[string]interface{}{
		"Data.map[rotation_period]: 86400 != 600":   true, // from 24h to 10m
		"Data.map[verification_ttl]: 86400 != 3600": true, // from 24h to 1h
	}
	diff := deep.Equal(respReadTestKey, respReadTestKeyAfterUpdate)
	expectStrings(t, diff, expectedDiff)

	// Create a role that depends on test key
	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/role/test-role",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"key": "test-key",
		},
		Storage: storage,
	})

	// Delete test-key -- should fail because test-role depends on test-key
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/key/test-key",
		Operation: logical.DeleteOperation,
		Storage:   storage,
	})
	expectError(t, resp, err)
	// validate error message
	expectedStrings := map[string]interface{}{
		"unable to delete key \"test-key\" because it is currently referenced by these roles: test-role": true,
	}
	expectStrings(t, []string{resp.Data["error"].(string)}, expectedStrings)

	// Delete test-role
	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/role/test-role",
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

// TestOIDC_Path_OIDCKey tests the List operation for keys
func TestOIDC_Path_OIDCKey(t *testing.T) {
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

	// .well-known/keys should contain 1 public key
	resp, err := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/.well-known/keys",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)
	// parse response
	responseJWKS := &jose.JSONWebKeySet{}
	json.Unmarshal(resp.Data["http_raw_body"].([]byte), responseJWKS)
	if len(responseJWKS.Keys) != 1 {
		t.Fatalf("expected 1 public key but instead got %d", len(responseJWKS.Keys))
	}

	// rotate test-key a few times, each rotate should increase the length of public keys returned
	// by the .well-known endpoint
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
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

	// .well-known/keys should contain 3 public keys
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/.well-known/keys",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)
	// parse response
	json.Unmarshal(resp.Data["http_raw_body"].([]byte), responseJWKS)
	if len(responseJWKS.Keys) != 3 {
		t.Fatalf("expected 3 public keya but instead got %d", len(responseJWKS.Keys))
	}

	// create another named key
	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/key/test-key2",
		Operation: logical.CreateOperation,
		Storage:   storage,
	})

	// delete test key
	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/key/test-key",
		Operation: logical.DeleteOperation,
		Storage:   storage,
	})

	// .well-known/keys should contain 1 public key, all of the public keys
	// from named key "test-key" should have been deleted
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/.well-known/keys",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)
	// parse response
	json.Unmarshal(resp.Data["http_raw_body"].([]byte), responseJWKS)
	if len(responseJWKS.Keys) != 1 {
		t.Fatalf("expected 1 public keya but instead got %d", len(responseJWKS.Keys))
	}
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

	// Generate a token against the role "test-role" -- should succeed
	resp, err := c.identityStore.HandleRequest(ctx, &logical.Request{
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
	responseJWKS := &jose.JSONWebKeySet{}
	json.Unmarshal(resp.Data["http_raw_body"].([]byte), responseJWKS)

	// Validate the signature
	claims := &jwt.Claims{}
	if err := parsedToken.Claims(responseJWKS.Keys[0], claims); err != nil {
		t.Fatalf("unable to validate signed token, err:\n%#v", err)
	}
}

// TestOIDC_PeriodicFunc tests timing logic for running key
// rotations and expiration actions.
func TestOIDC_PeriodicFunc(t *testing.T) {
	// Prepare a storage to run through periodicFunc
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	// populate storage with a named key
	period := 2 * time.Second
	keyName := "test-key"
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	id, _ := uuid.GenerateUUID()
	jwk := &jose.JSONWebKey{
		Key:       key,
		KeyID:     id,
		Algorithm: "RS256",
		Use:       "sig",
	}
	namedKey := &namedKey{
		Name:            keyName,
		Algorithm:       "RS256",
		VerificationTTL: 1 * period,
		RotationPeriod:  1 * period,
		KeyRing:         nil,
		SigningKey:      jwk,
		NextRotation:    time.Now().Add(1 * time.Second),
	}

	// Store namedKey
	entry, _ := logical.StorageEntryJSON(namedKeyConfigPath+keyName, namedKey)
	if err := storage.Put(ctx, entry); err != nil {
		t.Fatalf("writing to in mem storage failed")
	}

	// Time 0 - 1 Period
	// PeriodicFunc should set nextRun - nothing else
	c.identityStore.oidcPeriodicFunc(ctx, storage)
	entry, _ = storage.Get(ctx, namedKeyConfigPath+keyName)
	entry.DecodeJSON(&namedKey)
	if len(namedKey.KeyRing) != 0 {
		t.Fatalf("expected namedKey's KeyRing to be of length 0 but was: %#v", len(namedKey.KeyRing))
	}
	// There should be no public keys yet
	publicKeys, _ := storage.List(ctx, publicKeysConfigPath)
	if len(publicKeys) != 0 {
		t.Fatalf("expected publicKeys to be of length 0 but was: %#v", len(publicKeys))
	}
	// Next run should be set
	v, _ := c.identityStore.oidcCache.Get("nextRun")
	if v == nil {
		t.Fatalf("Expected nextRun to be set but it was nil")
	}
	earlierNextRun := v.(time.Time)

	// Time 1 - 2 Period
	// PeriodicFunc should rotate namedKey and update nextRun
	time.Sleep(period)
	c.identityStore.oidcPeriodicFunc(ctx, storage)
	entry, _ = storage.Get(ctx, namedKeyConfigPath+keyName)
	entry.DecodeJSON(&namedKey)
	if len(namedKey.KeyRing) != 1 {
		t.Fatalf("expected namedKey's KeyRing to be of length 1 but was: %#v", len(namedKey.KeyRing))
	}
	// There should be one public key
	publicKeys, _ = storage.List(ctx, publicKeysConfigPath)
	if len(publicKeys) != 1 {
		t.Fatalf("expected publicKeys to be of length 1 but was: %#v", len(publicKeys))
	}
	// nextRun should have been updated
	v, _ = c.identityStore.oidcCache.Get("nextRun")
	laterNextRun := v.(time.Time)
	if !laterNextRun.After(earlierNextRun) {
		t.Fatalf("laterNextRun: %#v is not after earlierNextRun: %#v", laterNextRun.String(), earlierNextRun.String())
	}

	// Time 2-3
	// PeriodicFunc should rotate namedKey and expire 1 public key
	time.Sleep(period)
	c.identityStore.oidcPeriodicFunc(ctx, storage)
	entry, _ = storage.Get(ctx, namedKeyConfigPath+keyName)
	entry.DecodeJSON(&namedKey)
	if len(namedKey.KeyRing) != 2 {
		t.Fatalf("expected namedKey's KeyRing to be of length 2 but was: %#v", len(namedKey.KeyRing))
	}
	// There should be two public keys
	publicKeys, _ = storage.List(ctx, publicKeysConfigPath)
	if len(publicKeys) != 2 {
		t.Fatalf("expected publicKeys to be of length 2 but was: %#v", len(publicKeys))
	}

	// Time 3-4
	// PeriodicFunc should rotate namedKey and expire 1 public key
	time.Sleep(period)
	c.identityStore.oidcPeriodicFunc(ctx, storage)
	entry, _ = storage.Get(ctx, namedKeyConfigPath+keyName)
	entry.DecodeJSON(&namedKey)
	if len(namedKey.KeyRing) != 2 {
		t.Fatalf("expected namedKey's KeyRing to be of length 1 but was: %#v", len(namedKey.KeyRing))
	}
	// There should be two public keys
	publicKeys, _ = storage.List(ctx, publicKeysConfigPath)
	if len(publicKeys) != 2 {
		t.Fatalf("expected publicKeys to be of length 1 but was: %#v", len(publicKeys))
	}
}

// some helpers
func expectSuccess(t *testing.T, resp *logical.Response, err error) {
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("expected success but got error:\n%v\nresp: %#v", err, resp)
	}
}

func expectError(t *testing.T, resp *logical.Response, err error) {
	if err == nil {
		if resp == nil || !resp.IsError() {
			t.Fatalf("expected error but got success; error:\n%v\nresp: %#v", err, resp)
		}
	}
}

// expectString fails unless every string in actualStrings is also included in expectedStrings and
// the length of actualStrings and expectedStrings are the same
func expectStrings(t *testing.T, actualStrings []string, expectedStrings map[string]interface{}) {
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
