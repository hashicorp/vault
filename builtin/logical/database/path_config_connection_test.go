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

func TestWriteConfig_PluginVersionInStorage(t *testing.T) {
	cluster, sys := getCluster(t)
	defer cluster.Cleanup()

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
	data := map[string]interface{}{
		"connection_url":    "test",
		"plugin_name":       hdb,
		"plugin_version":    hdbBuiltin,
		"verify_connection": false,
	}
	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/plugin-test",
		Storage:   config.StorageView,
		Data:      data,
	}
	resp, err := b.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	req = &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "config/plugin-test",
		Storage:   config.StorageView,
		Data:      data,
	}

	resp, err = b.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	if resp.Data["plugin_version"] != "" {
		t.Fatalf("expected plugin_version empty but got %s", resp.Data["plugin_version"])
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
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	if resp.Data["plugin_version"] != "" {
		t.Fatalf("expected plugin_version empty but got %s", resp.Data["plugin_version"])
	}
}

func TestWriteConfig_HelpfulErrorMessageWhenBuiltinOverridden(t *testing.T) {
	cluster, sys := getCluster(t)
	defer cluster.Cleanup()

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
