// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package blackbox

import (
	"testing"

	"github.com/hashicorp/vault/sdk/helper/testcluster/blackbox"
)

// testKVSecretsCreate tests KV secrets engine creation
func testKVSecretsCreate(t *testing.T, v *blackbox.Session) {
	// KV secrets engine tests are now in kvv2_test.go - just test basic enablement here
	SetupKVEngine(v, "kv-create")

	// Write and read test data to verify engine works
	v.MustWriteKV2("kv-create", "test/path", StandardKVData)
	secret := v.MustReadKV2("kv-create", "test/path")
	AssertKVData(t, v, secret, StandardKVData)

	t.Log("Successfully created and tested KV secrets engine")
}

// testKVSecretsRead tests KV secrets engine read operations
func testKVSecretsRead(t *testing.T, v *blackbox.Session) {
	// KV read tests are in kvv2_test.go - test basic read functionality here
	SetupKVEngine(v, "kv-read")
	v.MustWriteKV2("kv-read", "read/test", AltKVData)
	secret := v.MustReadKV2("kv-read", "read/test")
	AssertKVData(t, v, secret, AltKVData)

	t.Log("Successfully read KV secrets engine data")
}

// testKVSecretsDelete tests KV secrets engine delete operations
func testKVSecretsDelete(t *testing.T, v *blackbox.Session) {
	SetupKVEngine(v, "kv-delete")
	v.MustWriteKV2("kv-delete", "delete/test", StandardKVData)
	secret := v.MustReadKV2("kv-delete", "delete/test")
	AssertKVData(t, v, secret, StandardKVData)
	v.MustWrite("kv-delete/delete/delete/test", map[string]any{
		"versions": []int{1},
	})
	t.Log("Successfully deleted KV secrets engine data")
}
