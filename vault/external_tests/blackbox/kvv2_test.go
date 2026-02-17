// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package blackbox

import (
	"testing"

	"github.com/hashicorp/vault/sdk/helper/testcluster/blackbox"
)

// TestKVv2_SoftDeleteAndRestore_Workflow tests the complete workflow of KV v2 soft delete and restore
func TestKVv2_SoftDeleteAndRestore_Workflow(t *testing.T) {
	v := blackbox.New(t)

	// Setup KV engine and authenticated user using common utilities
	bob := SetupStandardKVUserpass(v, "secret", "bob", "lol")

	// Write initial data using standard test data
	testData := map[string]any{
		"api_key":     "A1B2-C3D4",
		"is_active":   true,
		"retry_count": 3,
	}
	v.MustWriteKV2("secret", "app-config", testData)

	// Verify data can be read
	secret := bob.MustReadKV2("secret", "app-config")
	AssertKVData(t, bob, secret, testData)

	// Perform soft delete
	bob.MustWrite("secret/delete/app-config", map[string]any{
		"versions": []int{1},
	})

	// Verify data is deleted
	deletedSecret := bob.MustReadRequired("secret/data/app-config")
	if deletedSecret.Data["data"] != nil {
		t.Fatal("Expected secret data to be nil after soft delete, but got data")
	}

	// Restore the data
	bob.MustWrite("secret/undelete/app-config", map[string]any{
		"versions": []int{1},
	})

	// Verify data is restored
	restoredSecret := bob.MustReadRequired("secret/data/app-config")
	bob.AssertSecret(restoredSecret).
		KV2().
		HasKey("api_key", "A1B2-C3D4")
}

// TestKVv2_BasicOperations tests basic KV v2 create, read, update operations
func TestKVv2_BasicOperations(t *testing.T) {
	v := blackbox.New(t)

	// Setup using common utilities
	user := SetupStandardKVUserpass(v, "kv-basic", "testuser", "testpass")

	// Test create
	user.MustWriteKV2("kv-basic", "test/data", StandardKVData)

	// Test read
	secret := user.MustReadKV2("kv-basic", "test/data")
	AssertKVData(t, user, secret, StandardKVData)

	// Test update
	updatedData := map[string]any{
		"api_key":     "updated-key-456",
		"is_active":   false,
		"retry_count": 5,
	}
	user.MustWriteKV2("kv-basic", "test/data", updatedData)

	// Verify update
	updatedSecret := user.MustReadKV2("kv-basic", "test/data")
	AssertKVData(t, user, updatedSecret, updatedData)
}
