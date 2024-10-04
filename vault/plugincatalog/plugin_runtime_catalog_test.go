// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package plugincatalog

import (
	"context"
	"reflect"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/pluginruntimeutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/physical/inmem"
)

func testPluginRuntimeCatalog(t *testing.T) *PluginRuntimeCatalog {
	logger := hclog.Default()
	storage, err := inmem.NewInmem(nil, logger)
	if err != nil {
		t.Fatal(err)
	}
	pluginRuntimeCatalog, err := SetupPluginRuntimeCatalog(context.Background(), logger, logical.NewLogicalStorage(storage))
	if err != nil {
		t.Fatal(err)
	}
	return pluginRuntimeCatalog
}

func TestPluginRuntimeCatalog_CRUD(t *testing.T) {
	pluginRuntimeCatalog := testPluginRuntimeCatalog(t)
	ctx := context.Background()

	expected := &pluginruntimeutil.PluginRuntimeConfig{
		Name:         "gvisor",
		OCIRuntime:   "runsc",
		CgroupParent: "/cpulimit/",
		CPU:          1,
		Memory:       10000,
	}

	// Set new plugin runtime
	err := pluginRuntimeCatalog.Set(ctx, expected)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Get plugin runtime
	runner, err := pluginRuntimeCatalog.Get(ctx, expected.Name, expected.Type)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !reflect.DeepEqual(expected, runner) {
		t.Fatalf("expected did not match actual, got %#v\n expected %#v\n", runner, expected)
	}

	// Set existing plugin runtime
	expected.CgroupParent = "memorylimit-cgroup"
	expected.CPU = 2
	expected.Memory = 5000
	err = pluginRuntimeCatalog.Set(ctx, expected)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Get plugin runtime again
	runner, err = pluginRuntimeCatalog.Get(ctx, expected.Name, expected.Type)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if !reflect.DeepEqual(expected, runner) {
		t.Fatalf("expected did not match actual, got %#v\n expected %#v\n", runner, expected)
	}

	configs, err := pluginRuntimeCatalog.List(ctx, expected.Type)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(configs) != 1 {
		t.Fatalf("expected plugin runtime catalog to have 1 container runtime but got %d", len(configs))
	}

	// Delete plugin runtime
	err = pluginRuntimeCatalog.Delete(ctx, expected.Name, expected.Type)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Assert the plugin runtime catalog is empty
	configs, err = pluginRuntimeCatalog.List(ctx, expected.Type)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(configs) != 0 {
		t.Fatalf("expected plugin runtime catalog to have 0 container runtimes but got %d", len(configs))
	}
}
