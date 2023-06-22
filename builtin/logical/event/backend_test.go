// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package event

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/vault/helper/builtinplugins"
	"github.com/hashicorp/vault/helper/namespace"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

func getCluster(t *testing.T) (*vault.TestCluster, logical.SystemView) {
	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"database": Factory,
		},
		BuiltinRegistry: builtinplugins.Registry,
	}

	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	cores := cluster.Cores
	vault.TestWaitActive(t, cores[0].Core)

	os.Setenv(pluginutil.PluginCACertPEMEnv, cluster.CACertPEMFile)

	sys := vault.TestDynamicSystemView(cores[0].Core, nil)
	return cluster, sys
}

// TestBackend_basic tests basic creating a plugin, then setting up, resetting, and deleting a subscription.
func TestBackend_basic(t *testing.T) {
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

	// configure a subscription
	data := map[string]interface{}{
		"plugin_name": "noop-event-plugin",
		"settings": map[string]interface{}{
			"nothing": "here",
		},
	}
	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "subscription/plugin-test",
		Storage:   config.StorageView,
		Data:      data,
	}
	_, err = b.HandleRequest(namespace.RootContext(nil), req)
	if err == nil {
		t.Fatalf("Expected this to fail for now")
	}

	// reset
	req = &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "reset/plugin-test",
		Storage:   config.StorageView,
		Data:      data,
	}
	_, err = b.HandleRequest(namespace.RootContext(nil), req)
	if err == nil {
		t.Fatalf("Expected this to fail for now")
	}

	req = &logical.Request{
		Operation: logical.DeleteOperation,
		Path:      "subscription/plugin-test",
		Storage:   config.StorageView,
		Data:      data,
	}
	resp, err := b.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}
}
