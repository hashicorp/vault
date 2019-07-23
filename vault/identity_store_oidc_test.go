package vault

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/go-test/deep"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/framework"
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
			"template": "{\"some-key\":\"some-value\"}",
			"ttl":      "2h",
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
		"client_id": resp.Data["client_id"],
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
	fmt.Printf("resp is:\n%#v", resp)

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
	expectSuccess(t, resp, err)

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
		t.Fatalf("expected 3 public keys but instead got %d", len(responseJWKS.Keys))
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
		t.Fatalf("expected 1 public keys but instead got %d", len(responseJWKS.Keys))
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

	// Prepare a dummy signing key
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	id, _ := uuid.GenerateUUID()
	jwk := &jose.JSONWebKey{
		Key:       key,
		KeyID:     id,
		Algorithm: "RS256",
		Use:       "sig",
	}

	cyclePeriod := 2 * time.Second

	var testSets = []struct {
		namedKey  *namedKey
		testCases []struct {
			cycle         int
			numKeys       int
			numPublicKeys int
		}
	}{
		{
			&namedKey{
				name:            "test-key",
				Algorithm:       "RS256",
				VerificationTTL: 1 * cyclePeriod,
				RotationPeriod:  1 * cyclePeriod,
				KeyRing:         nil,
				SigningKey:      jwk,
				NextRotation:    time.Now(),
			},
			[]struct {
				cycle         int
				numKeys       int
				numPublicKeys int
			}{
				{1, 1, 1},
				{2, 2, 2},
				{3, 2, 2},
				{4, 2, 2},
			},
		},
	}

	for _, testSet := range testSets {
		// Store namedKey
		storage := c.router.MatchingStorageByAPIPath(ctx, "identity/oidc")
		entry, _ := logical.StorageEntryJSON(namedKeyConfigPath+testSet.namedKey.name, testSet.namedKey)
		if err := storage.Put(ctx, entry); err != nil {
			t.Fatalf("writing to in mem storage failed")
		}

		currentCycle := 1
		numCases := len(testSet.testCases)
		lastCycle := testSet.testCases[numCases-1].cycle
		namedKeySamples := make([]*logical.StorageEntry, numCases)
		publicKeysSamples := make([][]string, numCases)

		i := 0
		// var start time.Time
		for currentCycle <= lastCycle {
			c.identityStore.oidcPeriodicFunc(ctx)
			if currentCycle == testSet.testCases[i].cycle {
				namedKeyEntry, _ := storage.Get(ctx, namedKeyConfigPath+testSet.namedKey.name)
				publicKeysEntry, _ := storage.List(ctx, publicKeysConfigPath)
				namedKeySamples[i] = namedKeyEntry
				publicKeysSamples[i] = publicKeysEntry
				i = i + 1
			}
			currentCycle = currentCycle + 1

			// sleep until we are in the next cycle - where a next run will happen
			v, _ := c.identityStore.oidcCache.Get(nilNamespace, "nextRun")
			nextRun := v.(time.Time)
			now := time.Now()
			diff := nextRun.Sub(now)
			if now.Before(nextRun) {
				time.Sleep(diff)
			}
		}

		// measure collected samples
		for i := range testSet.testCases {
			namedKeySamples[i].DecodeJSON(&testSet.namedKey)
			if len(testSet.namedKey.KeyRing) != testSet.testCases[i].numKeys {
				t.Fatalf("At cycle: %d expected namedKey's KeyRing to be of length %d but was: %d", testSet.testCases[i].cycle, testSet.testCases[i].numKeys, len(testSet.namedKey.KeyRing))
			}
			if len(publicKeysSamples[i]) != testSet.testCases[i].numPublicKeys {
				t.Fatalf("At cycle: %d expected public keys to be of length %d but was: %d", testSet.testCases[i].cycle, testSet.testCases[i].numPublicKeys, len(publicKeysSamples[i]))
			}
		}
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
				"name": &framework.FieldSchema{
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

	// Populte storage with a namedKey
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
				"name": &framework.FieldSchema{
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
				"name": &framework.FieldSchema{
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
				"name": &framework.FieldSchema{
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

func TestOIDC_isNamespacedKey(t *testing.T) {
	tests := []struct {
		nsid     string
		nskey    string
		expected bool
	}{
		{"nsid", "v0:nsid:key", true},
		{"nsid", "v0:nsid:", true},
		{"nsid", "v0:nsid", false},
		{"nsid", "v0:", false},
		{"nsid", "v0", false},
		{"nsid", "", false},
		{"nsid1", "v0:nsid2:key", false},
		{"nsid1", "nsid1:nsid2:nsid1", false},
		{"nsid1", "nsid1:nsid1:nsid1", true},
		{"nsid", "nsid:nsid:nsid:nsid:nsid:nsid", true},
		{"nsid", ":::", false},
		{"", ":::", true}, // "" is a valid key for cache.Set/Get
		{"nsid1", "nsid0:nsid1:nsid0:nsid1:nsid0:nsid1", true},
		{"nsid0", "nsid0:nsid1:nsid0:nsid1:nsid0:nsid1", false},
	}

	for _, test := range tests {
		actual := isNamespacedKey(test.nskey, test.nsid)
		if test.expected != actual {
			t.Fatalf("expected %t but got %t for nsid: %q and nskey: %q", test.expected, actual, test.nsid, test.nskey)
		}
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
