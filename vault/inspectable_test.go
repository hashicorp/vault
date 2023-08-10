// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package vault

import (
	"strings"
	"testing"

	"github.com/hashicorp/vault/command/server"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
)

func TestInspectRouter(t *testing.T) {
	// Verify that all the expected tables exist when we inspect the router
	coreConfig := &CoreConfig{
		EnableIntrospection: true,
	}
	c, _, root := TestCoreUnsealedWithConfig(t, coreConfig)

	rootCtx := namespace.RootContext(nil)
	subTrees := map[string][]string{
		"routeEntry": {"root", "storage"},
		"mountEntry": {"uuid", "accessor"},
	}
	subTreeFields := map[string][]string{
		"routeEntry": {"tainted", "storage_prefix", "accessor", "mount_namespace", "mount_path", "mount_type", "uuid"},
		"mountEntry": {"accessor", "mount_namespace", "mount_path", "mount_type", "uuid"},
	}
	for subTreeType, subTreeArray := range subTrees {
		for _, tag := range subTreeArray {
			resp, err := c.HandleRequest(rootCtx, &logical.Request{
				ClientToken: root,
				Operation:   logical.ReadOperation,
				Path:        "sys/internal/inspect/router/" + tag,
			})
			if err != nil || (resp != nil && resp.IsError()) {
				t.Fatalf("bad: resp: %#v\n err: %v", resp, err)
			}
			// Verify that data exists
			data, ok := resp.Data[tag].([]map[string]interface{})
			if !ok {
				t.Fatalf("Router data is malformed")
			}
			for _, entry := range data {
				for _, field := range subTreeFields[subTreeType] {
					if _, ok := entry[field]; !ok {
						t.Fatalf("Field %s not found in %s", field, tag)
					}
				}
			}

		}
	}
}

func TestInvalidInspectRouterPath(t *testing.T) {
	// Verify that attempting to inspect an invalid tree in the router fails
	coreConfig := &CoreConfig{
		EnableIntrospection: true,
	}
	core, _, rootToken := TestCoreUnsealedWithConfig(t, coreConfig)
	rootCtx := namespace.RootContext(nil)
	_, err := core.HandleRequest(rootCtx, &logical.Request{
		ClientToken: rootToken,
		Operation:   logical.ReadOperation,
		Path:        "sys/internal/inspect/router/random",
	})
	if !strings.Contains(err.Error(), logical.ErrUnsupportedPath.Error()) {
		t.Fatal("expected unsupported path error")
	}
}

func TestInspectAPIDisabled(t *testing.T) {
	// Verify that the Inspect API is turned off by default
	core, _, rootToken := testCoreSystemBackend(t)
	rootCtx := namespace.RootContext(nil)
	resp, err := core.HandleRequest(rootCtx, &logical.Request{
		ClientToken: rootToken,
		Operation:   logical.ReadOperation,
		Path:        "sys/internal/inspect/router/root",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !resp.IsError() || !strings.Contains(resp.Error().Error(), ErrIntrospectionNotEnabled.Error()) {
		t.Fatal("expected invalid configuration error")
	}
}

func TestInspectAPISudoProtect(t *testing.T) {
	// Verify that the Inspect API path is sudo protected
	core, _, rootToken := testCoreSystemBackend(t)
	testMakeServiceTokenViaBackend(t, core.tokenStore, rootToken, "tokenid", "", []string{"secret"})
	rootCtx := namespace.RootContext(nil)
	_, err := core.HandleRequest(rootCtx, &logical.Request{
		ClientToken: "tokenid",
		Operation:   logical.ReadOperation,
		Path:        "sys/internal/inspect/router/root",
	})
	if !strings.Contains(err.Error(), logical.ErrPermissionDenied.Error()) {
		t.Fatal("expected permission denied error")
	}
}

func TestInspectAPIReload(t *testing.T) {
	// Verify that the Inspect API is turned off by default
	core, _, rootToken := testCoreSystemBackend(t)
	rootCtx := namespace.RootContext(nil)
	resp, err := core.HandleRequest(rootCtx, &logical.Request{
		ClientToken: rootToken,
		Operation:   logical.ReadOperation,
		Path:        "sys/internal/inspect/router/root",
	})
	if err != nil {
		t.Fatal("Unexpected error")
	}
	if !resp.IsError() {
		t.Fatal("expected invalid configuration error")
	}
	if !strings.Contains(resp.Error().Error(), ErrIntrospectionNotEnabled.Error()) {
		t.Fatalf("expected invalid configuration error but recieved: %s", resp.Error())
	}

	originalConfig := core.rawConfig.Load().(*server.Config)
	newConfig := originalConfig
	newConfig.EnableIntrospectionEndpointRaw = true
	newConfig.EnableIntrospectionEndpoint = true

	// Reload Endpoint
	core.SetConfig(newConfig)
	core.ReloadIntrospectionEndpointEnabled()

	resp, err = core.HandleRequest(rootCtx, &logical.Request{
		ClientToken: rootToken,
		Operation:   logical.ReadOperation,
		Path:        "sys/internal/inspect/router/root",
	})
	if err != nil || resp.IsError() {
		t.Fatal("Unexpected error after reload")
	}
}
