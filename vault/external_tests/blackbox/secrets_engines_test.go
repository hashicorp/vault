// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package blackbox

import (
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/testcluster/blackbox"
)

// TestSecretsEngineCreate tests creation/setup of various secrets engines
func TestSecretsEngineCreate(t *testing.T) {
	v := blackbox.New(t)

	// Verify we have a healthy cluster first
	v.AssertClusterHealthy()

	// KV secrets engine tests are now in kvv2_test.go - just test basic enablement here
	t.Run("KVSecrets", func(t *testing.T) {
		SetupKVEngine(v, "kv-create")
		// Write and read test data to verify engine works
		v.MustWriteKV2("kv-create", "test/path", StandardKVData)
		secret := v.MustReadKV2("kv-create", "test/path")
		AssertKVData(t, v, secret, StandardKVData)
		t.Log("Successfully created and tested KV secrets engine")
	})

	// Stub out remaining secret engine creation tests
	t.Run("PKISecrets", func(t *testing.T) {
		t.Skip("PKI secrets engine create test - implementation pending")
	})

	t.Run("SSHSecrets", func(t *testing.T) {
		t.Skip("SSH secrets engine create test - implementation pending")
	})

	t.Run("IdentitySecrets", func(t *testing.T) {
		t.Skip("Identity secrets engine create test - implementation pending")
	})

	t.Run("AWSSecrets", func(t *testing.T) {
		t.Skip("AWS secrets engine create test - implementation pending")
	})

	t.Run("LDAPSecrets", func(t *testing.T) {
		t.Skip("LDAP secrets engine create test - implementation pending")
	})

	t.Run("KMIPSecrets", func(t *testing.T) {
		t.Skip("KMIP secrets engine create test - implementation pending")
	})

	t.Run("DatabaseSecrets", func(t *testing.T) {
		t.Skip("Database secrets engine create test - implementation pending")
	})

	t.Run("TransitSecrets", func(t *testing.T) {
		testTransitSecretsCreate(t, v)
	})
}

// TestSecretsEngineRead tests read operations for various secrets engines
func TestSecretsEngineRead(t *testing.T) {
	v := blackbox.New(t)

	// Verify we have a healthy cluster first
	v.AssertClusterHealthy()

	// KV read tests are in kvv2_test.go - test basic read functionality here
	t.Run("KVSecrets", func(t *testing.T) {
		SetupKVEngine(v, "kv-read")
		v.MustWriteKV2("kv-read", "read/test", AltKVData)
		secret := v.MustReadKV2("kv-read", "read/test")
		AssertKVData(t, v, secret, AltKVData)
		t.Log("Successfully read KV secrets engine data")
	})

	// Stub out remaining secret engine read tests
	t.Run("PKISecrets", func(t *testing.T) {
		t.Skip("PKI secrets engine read test - implementation pending")
	})

	t.Run("SSHSecrets", func(t *testing.T) {
		t.Skip("SSH secrets engine read test - implementation pending")
	})

	t.Run("IdentitySecrets", func(t *testing.T) {
		t.Skip("Identity secrets engine read test - implementation pending")
	})

	t.Run("AWSSecrets", func(t *testing.T) {
		t.Skip("AWS secrets engine read test - implementation pending")
	})

	t.Run("LDAPSecrets", func(t *testing.T) {
		t.Skip("LDAP secrets engine read test - implementation pending")
	})

	t.Run("KMIPSecrets", func(t *testing.T) {
		t.Skip("KMIP secrets engine read test - implementation pending")
	})

	t.Run("DatabaseSecrets", func(t *testing.T) {
		t.Skip("Database secrets engine read test - implementation pending")
	})

	t.Run("TransitSecrets", func(t *testing.T) {
		t.Skip("Transit secrets engine read test - implementation pending")
	})
}

// TestSecretsEngineDelete tests delete operations for various secrets engines
func TestSecretsEngineDelete(t *testing.T) {
	v := blackbox.New(t)

	// Verify we have a healthy cluster first
	v.AssertClusterHealthy()

	// KV delete tests are in kvv2_test.go - test basic delete functionality here
	t.Run("KVSecrets", func(t *testing.T) {
		SetupKVEngine(v, "kv-delete")

		// Write test data
		v.MustWriteKV2("kv-delete", "delete/test", StandardKVData)

		// Verify it exists
		secret := v.MustReadKV2("kv-delete", "delete/test")
		AssertKVData(t, v, secret, StandardKVData)

		// Delete using KV v2 delete endpoint
		v.MustWrite("kv-delete/delete/delete/test", map[string]any{
			"versions": []int{1},
		})

		t.Log("Successfully deleted KV secrets engine data")
	})

	// Stub out remaining secret engine delete tests
	t.Run("PKISecrets", func(t *testing.T) {
		t.Skip("PKI secrets engine delete test - implementation pending")
	})

	t.Run("SSHSecrets", func(t *testing.T) {
		t.Skip("SSH secrets engine delete test - implementation pending")
	})

	t.Run("IdentitySecrets", func(t *testing.T) {
		t.Skip("Identity secrets engine delete test - implementation pending")
	})

	t.Run("AWSSecrets", func(t *testing.T) {
		t.Skip("AWS secrets engine delete test - implementation pending")
	})

	t.Run("LDAPSecrets", func(t *testing.T) {
		t.Skip("LDAP secrets engine delete test - implementation pending")
	})

	t.Run("KMIPSecrets", func(t *testing.T) {
		t.Skip("KMIP secrets engine delete test - implementation pending")
	})

	t.Run("DatabaseSecrets", func(t *testing.T) {
		t.Skip("Database secrets engine delete test - implementation pending")
	})

	t.Run("TransitSecrets", func(t *testing.T) {
		testTransitSecretsDelete(t, v)
	})
}

// Transit Secrets Engine Test Implementation Functions

func testTransitSecretsCreate(t *testing.T, v *blackbox.Session) {
	// Enable transit secrets engine
	v.MustEnableSecretsEngine("transit", &api.MountInput{Type: "transit"})

	// Create an encryption key
	keyName := "test-key"
	v.MustWrite("transit/keys/"+keyName, map[string]any{
		"type": "aes256-gcm96",
	})

	// Verify the key was created by reading it
	keyInfo := v.MustRead("transit/keys/" + keyName)
	if keyInfo.Data == nil {
		t.Fatal("Expected to read key configuration")
	}

	// Verify key type
	if keyType, ok := keyInfo.Data["type"]; !ok || keyType != "aes256-gcm96" {
		t.Fatalf("Expected key type 'aes256-gcm96', got: %v", keyInfo.Data["type"])
	}

	// Test encryption
	plaintext := "dGhlIHF1aWNrIGJyb3duIGZveA==" // base64 encoded "the quick brown fox"
	encryptResp := v.MustWrite("transit/encrypt/"+keyName, map[string]any{
		"plaintext": plaintext,
	})

	if encryptResp.Data == nil || encryptResp.Data["ciphertext"] == nil {
		t.Fatal("Expected ciphertext in encryption response")
	}

	ciphertext := encryptResp.Data["ciphertext"].(string)
	t.Logf("Encrypted ciphertext: %s", ciphertext[:20]+"...")

	// Test decryption
	decryptResp := v.MustWrite("transit/decrypt/"+keyName, map[string]any{
		"ciphertext": ciphertext,
	})

	if decryptResp.Data == nil || decryptResp.Data["plaintext"] == nil {
		t.Fatal("Expected plaintext in decryption response")
	}

	decryptedText := decryptResp.Data["plaintext"].(string)
	if decryptedText != plaintext {
		t.Fatalf("Decrypted text doesn't match original. Expected: %s, Got: %s", plaintext, decryptedText)
	}

	t.Log("Successfully created transit secrets engine and tested encryption/decryption")
}

func testTransitSecretsRead(t *testing.T, v *blackbox.Session) {
	// Enable transit secrets engine with unique mount
	v.MustEnableSecretsEngine("transit-read", &api.MountInput{Type: "transit"})

	// Create an encryption key
	keyName := "read-test-key"
	v.MustWrite("transit-read/keys/"+keyName, map[string]any{
		"type":       "aes256-gcm96",
		"exportable": false,
	})

	// Read the key configuration
	keyInfo := v.MustRead("transit-read/keys/" + keyName)
	if keyInfo.Data == nil {
		t.Fatal("Expected to read key configuration")
	}

	// Verify key properties
	assertions := v.AssertSecret(keyInfo)
	assertions.Data().
		HasKey("type", "aes256-gcm96").
		HasKey("exportable", false).
		HasKeyExists("keys").
		HasKeyExists("latest_version")

	t.Log("Successfully read transit secrets engine key configuration")
}

func testTransitSecretsDelete(t *testing.T, v *blackbox.Session) {
	// Enable transit secrets engine with unique mount
	v.MustEnableSecretsEngine("transit-delete", &api.MountInput{Type: "transit"})

	// Create an encryption key
	keyName := "delete-test-key"
	v.MustWrite("transit-delete/keys/"+keyName, map[string]any{
		"type": "aes256-gcm96",
	})

	// Verify the key exists
	keyInfo := v.MustRead("transit-delete/keys/" + keyName)
	if keyInfo.Data == nil {
		t.Fatal("Expected key to exist before deletion")
	}

	// Configure the key to allow deletion (transit keys require this)
	v.MustWrite("transit-delete/keys/"+keyName+"/config", map[string]any{
		"deletion_allowed": true,
	})

	// Delete the key
	_, err := v.Client.Logical().Delete("transit-delete/keys/" + keyName)
	if err != nil {
		t.Fatalf("Failed to delete transit key: %v", err)
	}

	// Verify the key is deleted by attempting to read it
	readSecret, err := v.Client.Logical().Read("transit-delete/keys/" + keyName)
	if err == nil && readSecret != nil {
		t.Fatal("Expected key to be deleted, but it still exists")
	}

	t.Logf("Successfully deleted transit key: %s", keyName)
}
