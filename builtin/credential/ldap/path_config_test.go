// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package ldap

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/vault/sdk/helper/ldaputil"
	"github.com/hashicorp/vault/sdk/logical"
)

// TestConfig_UpgradePath_Schema verifies that a ConfigEntry written to storage
// before the "schema" field existed (i.e. Schema == "") is transparently
// upgraded to the default value (SchemaOpenLDAP) when read back via Config().
// This mirrors the upgrade handling already in place for CaseSensitiveNames
// and UsePre111GroupCNBehavior.
func TestConfig_UpgradePath_Schema(t *testing.T) {
	ctx := context.Background()
	b, store := createBackendWithStorage(t)

	// Simulate a pre-existing stored config that has no "schema" key —
	// i.e. data written before the field was introduced. We do this by
	// marshalling a map that deliberately omits the field.
	oldEntry := map[string]interface{}{
		"url":      "ldap://127.0.0.1",
		"userdn":   "ou=People,dc=example,dc=org",
		"binddn":   "cn=admin,dc=example,dc=org",
		"bindpass": "secret",
		// "schema" intentionally absent
		"CaseSensitiveNames":           false,
		"use_pre111_group_cn_behavior": true,
	}
	raw, err := json.Marshal(oldEntry)
	if err != nil {
		t.Fatalf("failed to marshal old config entry: %s", err)
	}
	entry := &logical.StorageEntry{
		Key:   "config",
		Value: raw,
	}
	if err := store.Put(ctx, entry); err != nil {
		t.Fatalf("failed to write legacy config to storage: %s", err)
	}

	// Read back via Config() — this should trigger the upgrade path.
	req := &logical.Request{Storage: store}
	cfg, err := b.Config(ctx, req)
	if err != nil {
		t.Fatalf("Config() returned unexpected error: %s", err)
	}
	if cfg == nil {
		t.Fatal("Config() returned nil; expected a valid ConfigEntry")
	}

	// Schema must be defaulted to SchemaOpenLDAP, not left as "".
	if cfg.Schema != ldaputil.SchemaOpenLDAP {
		t.Errorf("expected Schema=%q after upgrade, got %q", ldaputil.SchemaOpenLDAP, cfg.Schema)
	}

	// Verify the migrated value was persisted to storage.
	persisted, err := store.Get(ctx, "config")
	if err != nil {
		t.Fatalf("failed to read persisted config: %s", err)
	}
	if persisted == nil {
		t.Fatal("expected config to be re-persisted after upgrade, but nothing found in storage")
	}
	var persistedMap map[string]interface{}
	if err := json.Unmarshal(persisted.Value, &persistedMap); err != nil {
		t.Fatalf("failed to decode persisted config: %s", err)
	}
	if persistedMap["schema"] != ldaputil.SchemaOpenLDAP {
		t.Errorf("persisted schema=%q; expected %q", persistedMap["schema"], ldaputil.SchemaOpenLDAP)
	}
}

// TestConfig_UpgradePath_Schema_NoOverwrite ensures that an already-set schema
// value is not overwritten by the upgrade logic.
func TestConfig_UpgradePath_Schema_NoOverwrite(t *testing.T) {
	ctx := context.Background()
	b, store := createBackendWithStorage(t)

	oldEntry := map[string]interface{}{
		"url":                          "ldap://127.0.0.1",
		"userdn":                       "ou=People,dc=example,dc=org",
		"binddn":                       "cn=admin,dc=example,dc=org",
		"bindpass":                     "secret",
		"schema":                       ldaputil.SchemaAD, // explicitly set to AD
		"CaseSensitiveNames":           false,
		"use_pre111_group_cn_behavior": true,
	}
	raw, err := json.Marshal(oldEntry)
	if err != nil {
		t.Fatalf("failed to marshal config entry: %s", err)
	}
	entry := &logical.StorageEntry{Key: "config", Value: raw}
	if err := store.Put(ctx, entry); err != nil {
		t.Fatalf("failed to write config to storage: %s", err)
	}

	req := &logical.Request{Storage: store}
	cfg, err := b.Config(ctx, req)
	if err != nil {
		t.Fatalf("Config() returned unexpected error: %s", err)
	}

	if cfg.Schema != ldaputil.SchemaAD {
		t.Errorf("expected Schema=%q to be preserved, got %q", ldaputil.SchemaAD, cfg.Schema)
	}
}
