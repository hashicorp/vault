// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package plugin

import (
	"context"
	"errors"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/plugin"
)

func TestBackend_lazyLoad(t *testing.T) {
	// normal load
	var invocations int
	b := testLazyLoad(t, func() error {
		invocations++
		return nil
	})
	if invocations != 1 {
		t.Fatalf("expected 1 invocation")
	}
	if b.canary != "" {
		t.Fatalf("expected empty canary")
	}

	// load with plugin shutdown
	invocations = 0
	b = testLazyLoad(t, func() error {
		invocations++
		if invocations == 1 {
			return plugin.ErrPluginShutdown
		}
		return nil
	})
	if invocations != 2 {
		t.Fatalf("expected 2 invocations")
	}
	if b.canary == "" {
		t.Fatalf("expected canary")
	}
}

func testLazyLoad(t *testing.T, methodWrapper func() error) *PluginBackend {
	sysView := newTestSystemView()

	ctx := context.Background()
	config := &logical.BackendConfig{
		Logger: logging.NewVaultLogger(hclog.Trace),
		System: sysView,
		Config: map[string]string{
			"plugin_name": "test-plugin",
			"plugin_type": "secret",
		},
	}

	// this is a dummy plugin that hasn't really been loaded yet
	orig, err := plugin.NewBackend(ctx, "test-plugin", consts.PluginTypeSecrets, sysView, config, true)
	if err != nil {
		t.Fatal(err)
	}

	b := &PluginBackend{
		Backend: orig,
		config:  config,
	}

	// lazy load
	err = b.lazyLoadBackend(ctx, &logical.InmemStorage{}, methodWrapper)
	if err != nil {
		t.Fatal(err)
	}
	if !b.loaded {
		t.Fatalf("not loaded")
	}

	// make sure dummy plugin was handled properly
	ob := orig.(*testBackend)
	if !ob.cleaned {
		t.Fatalf("not cleaned")
	}
	if ob.setup {
		t.Fatalf("setup")
	}
	if ob.initialized {
		t.Fatalf("initialized")
	}

	// make sure our newly initialized plugin was handled properly
	nb := b.Backend.(*testBackend)
	if nb.cleaned {
		t.Fatalf("cleaned")
	}
	if !nb.setup {
		t.Fatalf("not setup")
	}
	if !nb.initialized {
		t.Fatalf("not initialized")
	}

	return b
}

//------------------------------------------------------------------

type testBackend struct {
	cleaned     bool
	setup       bool
	initialized bool
}

var _ logical.Backend = (*testBackend)(nil)

func (b *testBackend) Cleanup(context.Context) {
	b.cleaned = true
}

func (b *testBackend) Setup(context.Context, *logical.BackendConfig) error {
	b.setup = true
	return nil
}

func (b *testBackend) Initialize(context.Context, *logical.InitializationRequest) error {
	b.initialized = true
	return nil
}

func (b *testBackend) Type() logical.BackendType {
	return logical.TypeLogical
}

func (b *testBackend) SpecialPaths() *logical.Paths {
	return &logical.Paths{
		Root: []string{"test-root"},
	}
}

func (b *testBackend) Logger() hclog.Logger {
	return logging.NewVaultLogger(hclog.Trace)
}

func (b *testBackend) HandleRequest(context.Context, *logical.Request) (*logical.Response, error) {
	panic("not needed")
}

func (b *testBackend) System() logical.SystemView {
	panic("not needed")
}

func (b *testBackend) HandleExistenceCheck(context.Context, *logical.Request) (bool, bool, error) {
	panic("not needed")
}

func (b *testBackend) InvalidateKey(context.Context, string) {
	panic("not needed")
}

//------------------------------------------------------------------

type testSystemView struct {
	logical.StaticSystemView
	factory logical.Factory
}

func newTestSystemView() testSystemView {
	return testSystemView{
		factory: func(_ context.Context, _ *logical.BackendConfig) (logical.Backend, error) {
			return &testBackend{}, nil
		},
	}
}

func (v testSystemView) LookupPlugin(context.Context, string, consts.PluginType) (*pluginutil.PluginRunner, error) {
	return &pluginutil.PluginRunner{
		Name:    "test-plugin-runner",
		Builtin: true,
		BuiltinFactory: func() (interface{}, error) {
			return v.factory, nil
		},
	}, nil
}

func (v testSystemView) LookupPluginVersion(context.Context, string, consts.PluginType, string) (*pluginutil.PluginRunner, error) {
	return &pluginutil.PluginRunner{
		Name:    "test-plugin-runner",
		Builtin: true,
		BuiltinFactory: func() (interface{}, error) {
			return v.factory, nil
		},
	}, nil
}

func (v testSystemView) ListVersionedPlugins(_ context.Context, _ consts.PluginType) ([]pluginutil.VersionedPlugin, error) {
	return nil, errors.New("ListVersionedPlugins not implemented for testSystemView")
}
