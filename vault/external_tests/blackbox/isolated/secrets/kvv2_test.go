//go:build isolated
// +build isolated

// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package secrets

import (
	"testing"

	"github.com/hashicorp/vault/sdk/helper/testcluster/blackbox"
	helpers "github.com/hashicorp/vault/vault/external_tests/blackbox"
)

// TestKVv2_SoftDeleteAndRestore_Workflow tests the complete workflow of KV v2 soft delete and restore
func TestKVv2_SoftDeleteAndRestore_Workflow(t *testing.T) {
	v := blackbox.New(t)

	// Setup KV engine
	helpers.SetupKVEngine(v, "secret")

	// Write initial data using standard test data
	testData := map[string]any{
		"api_key":     "A1B2-C3D4",
		"is_active":   true,
		"retry_count": 3,
	}
	v.MustWriteKV2("secret", "app-config", testData)

	// Verify data can be read
	secret := v.MustReadKV2("secret", "app-config")
	helpers.AssertKVData(t, v, secret, testData)

	// Perform soft delete
	v.MustWrite("secret/delete/app-config", map[string]any{
		"versions": []int{1},
	})

	// Verify data is deleted
	deletedSecret := v.MustReadRequired("secret/data/app-config")
	if deletedSecret.Data["data"] != nil {
		t.Fatal("Expected secret data to be nil after soft delete, but got data")
	}

	// Restore the data
	v.MustWrite("secret/undelete/app-config", map[string]any{
		"versions": []int{1},
	})

	// Verify data is restored
	restoredSecret := v.MustReadRequired("secret/data/app-config")
	v.AssertSecret(restoredSecret).
		KV2().
		HasKey("api_key", "A1B2-C3D4")
}

// TestKVv2_BasicOperations tests basic KV v2 create, read, update operations
func TestKVv2_BasicOperations(t *testing.T) {
	v := blackbox.New(t)

	// Setup KV engine
	helpers.SetupKVEngine(v, "kv-basic")

	// Test create
	testData := map[string]any{
		"api_key":     "test-key-123",
		"is_active":   true,
		"retry_count": 3,
	}
	v.MustWriteKV2("kv-basic", "test/data", testData)

	// Test read
	secret := v.MustReadKV2("kv-basic", "test/data")
	helpers.AssertKVData(t, v, secret, testData)

	// Test update
	updatedData := map[string]any{
		"api_key":     "updated-key-456",
		"is_active":   false,
		"retry_count": 5,
	}
	v.MustWriteKV2("kv-basic", "test/data", updatedData)

	// Verify update
	updatedSecret := v.MustReadKV2("kv-basic", "test/data")
	helpers.AssertKVData(t, v, updatedSecret, updatedData)
}
