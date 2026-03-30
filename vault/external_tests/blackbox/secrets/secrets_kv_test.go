// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package secrets

import (
	"testing"

	"github.com/hashicorp/vault/sdk/helper/testcluster/blackbox"
	helpers "github.com/hashicorp/vault/vault/external_tests/blackbox"
)

// testKVSecretsCreate tests KV secrets engine creation
func testKVSecretsCreate(t *testing.T, v *blackbox.Session) {
	// KV secrets engine tests are now in kvv2_test.go - just test basic enablement here
	helpers.SetupKVEngine(v, "kv-create")

	// Write and read test data to verify engine works
	v.MustWriteKV2("kv-create", "test/path", helpers.StandardKVData)
	secret := v.MustReadKV2("kv-create", "test/path")
	helpers.AssertKVData(t, v, secret, helpers.StandardKVData)

	t.Log("Successfully created and tested KV secrets engine")
}

// testKVSecretsRead tests KV secrets engine read operations
func testKVSecretsRead(t *testing.T, v *blackbox.Session) {
	// KV read tests are in kvv2_test.go - test basic read functionality here
	helpers.SetupKVEngine(v, "kv-read")
	v.MustWriteKV2("kv-read", "read/test", helpers.AltKVData)
	secret := v.MustReadKV2("kv-read", "read/test")
	helpers.AssertKVData(t, v, secret, helpers.AltKVData)

	t.Log("Successfully read KV secrets engine data")
}

// testKVSecretsDelete tests KV secrets engine delete operations
func testKVSecretsDelete(t *testing.T, v *blackbox.Session) {
	helpers.SetupKVEngine(v, "kv-delete")
	v.MustWriteKV2("kv-delete", "delete/test", helpers.StandardKVData)
	secret := v.MustReadKV2("kv-delete", "delete/test")
	helpers.AssertKVData(t, v, secret, helpers.StandardKVData)
	v.MustWrite("kv-delete/delete/delete/test", map[string]any{
		"versions": []int{1},
	})
	t.Log("Successfully deleted KV secrets engine data")
}
