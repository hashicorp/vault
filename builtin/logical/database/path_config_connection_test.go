// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package database

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/versions"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/logical"
)

// TestDisplayNameTemplateWarning tests whether a warning is
// issued when a custom username_template references .DisplayName without a truncate function.
func TestDisplayNameTemplateWarning(t *testing.T) {
	tests := map[string]struct {
		connectionDetails map[string]interface{}
		wantWarning       bool
	}{
		"no username_template key": {
			connectionDetails: map[string]interface{}{},
			wantWarning:       false,
		},
		"username_template is nil": {
			connectionDetails: map[string]interface{}{"username_template": nil},
			wantWarning:       false,
		},
		"username_template is empty string": {
			connectionDetails: map[string]interface{}{"username_template": ""},
			wantWarning:       false,
		},
		"username_template is non-string": {
			connectionDetails: map[string]interface{}{"username_template": 42},
			wantWarning:       false,
		},
		"template without DisplayName": {
			connectionDetails: map[string]interface{}{"username_template": "{{.RoleName}}-{{random 10}}"},
			wantWarning:       false,
		},
		"DisplayName with truncate": {
			connectionDetails: map[string]interface{}{"username_template": "{{.DisplayName | truncate 8}}-{{random 10}}"},
			wantWarning:       false,
		},
		"DisplayName with wrapped pipeline truncate": {
			connectionDetails: map[string]interface{}{"username_template": "{{printf \"%s\" .DisplayName | truncate 8}}-{{random 10}}"},
			wantWarning:       false,
		},
		"truncate only applies to RoleName": {
			connectionDetails: map[string]interface{}{"username_template": "{{.RoleName | truncate 8}}-{{.DisplayName}}-{{random 10}}"},
			wantWarning:       true,
		},
		"DisplayName without truncate in action": {
			connectionDetails: map[string]interface{}{"username_template": "{{.DisplayName | upper}}-{{random 10}}"},
			wantWarning:       true,
		},
		"DisplayName without truncate": {
			connectionDetails: map[string]interface{}{"username_template": "{{.DisplayName}}-{{random 10}}"},
			wantWarning:       true,
		},
		"DisplayName only in comment, no warning": {
			connectionDetails: map[string]interface{}{"username_template": "{{/* .DisplayName is intentionally not used */}}{{.RoleName | truncate 8}}"},
			wantWarning:       false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := displayNameTemplateWarning(tc.connectionDetails)
			if tc.wantWarning && got == "" {
				t.Fatalf("expected a warning, got none")
			}
			if !tc.wantWarning && got != "" {
				t.Fatalf("expected no warning, got %q", got)
			}
		})
	}
}

func TestWriteConfig_PluginVersionInStorage(t *testing.T) {
	_, sys := getCluster(t)

	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	config.System = sys

	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}
	defer b.Cleanup(context.Background())

	const hdb = "hana-database-plugin"
	hdbBuiltin := versions.GetBuiltinVersion(consts.PluginTypeDatabase, hdb)

	// Configure a connection
	writePluginVersion := func() {
		t.Helper()
		req := &logical.Request{
			Operation: logical.UpdateOperation,
			Path:      "config/plugin-test",
			Storage:   config.StorageView,
			Data: map[string]interface{}{
				"connection_url":    "test",
				"plugin_name":       hdb,
				"plugin_version":    hdbBuiltin,
				"verify_connection": false,
			},
		}
		resp, err := b.HandleRequest(namespace.RootContext(nil), req)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("err:%s resp:%#v\n", err, resp)
		}
	}
	writePluginVersion()

	getPluginVersionFromAPI := func() string {
		t.Helper()
		req := &logical.Request{
			Operation: logical.ReadOperation,
			Path:      "config/plugin-test",
			Storage:   config.StorageView,
		}

		resp, err := b.HandleRequest(namespace.RootContext(nil), req)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("err:%s resp:%#v\n", err, resp)
		}

		return resp.Data["plugin_version"].(string)
	}
	pluginVersion := getPluginVersionFromAPI()
	if pluginVersion != "" {
		t.Fatalf("expected plugin_version empty but got %s", pluginVersion)
	}

	// Directly store config to get the builtin plugin version into storage,
	// simulating a write that happened before upgrading to 1.12.2+
	err = storeConfig(context.Background(), config.StorageView, "plugin-test", &DatabaseConfig{
		PluginName:    hdb,
		PluginVersion: hdbBuiltin,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Now replay the read request, and we still shouldn't get the builtin version back.
	pluginVersion = getPluginVersionFromAPI()
	if pluginVersion != "" {
		t.Fatalf("expected plugin_version empty but got %s", pluginVersion)
	}

	// Check the underlying data, which should still have the version in storage.
	getPluginVersionFromStorage := func() string {
		t.Helper()
		entry, err := config.StorageView.Get(context.Background(), "config/plugin-test")
		if err != nil {
			t.Fatal(err)
		}
		if entry == nil {
			t.Fatal()
		}

		var config DatabaseConfig
		if err := entry.DecodeJSON(&config); err != nil {
			t.Fatal(err)
		}
		return config.PluginVersion
	}

	storagePluginVersion := getPluginVersionFromStorage()
	if storagePluginVersion != hdbBuiltin {
		t.Fatalf("Expected %s, got: %s", hdbBuiltin, storagePluginVersion)
	}

	// Trigger a write to storage, which should clean up plugin version in the storage entry.
	writePluginVersion()

	storagePluginVersion = getPluginVersionFromStorage()
	if storagePluginVersion != "" {
		t.Fatalf("Expected empty, got: %s", storagePluginVersion)
	}

	// Finally, confirm API requests still return empty plugin version too
	pluginVersion = getPluginVersionFromAPI()
	if pluginVersion != "" {
		t.Fatalf("expected plugin_version empty but got %s", pluginVersion)
	}
}

func TestWriteConfig_HelpfulErrorMessageWhenBuiltinOverridden(t *testing.T) {
	_, sys := getClusterPostgresDB(t)

	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	config.System = sys

	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}
	defer b.Cleanup(context.Background())

	const pg = "postgresql-database-plugin"
	pgBuiltin := versions.GetBuiltinVersion(consts.PluginTypeDatabase, pg)

	// Configure a connection
	data := map[string]interface{}{
		"connection_url":    "test",
		"plugin_name":       pg,
		"plugin_version":    pgBuiltin,
		"verify_connection": false,
	}
	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/plugin-test",
		Storage:   config.StorageView,
		Data:      data,
	}
	resp, err := b.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || !resp.IsError() {
		t.Fatalf("resp:%#v", resp)
	}
	if !strings.Contains(resp.Error().Error(), "overridden by an unversioned plugin") {
		t.Fatalf("expected overridden error but got: %s", resp.Error())
	}
}
