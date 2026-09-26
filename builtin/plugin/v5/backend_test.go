// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package plugin

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/logical"
	sdkplugin "github.com/hashicorp/vault/sdk/plugin"
)

func TestBackend_HandleRequest_BlocksReactiveReloadDuringReloadWindow(t *testing.T) {
	b := &backend{
		Backend: &testReloadingBackend{
			reloading: true,
			reqErr:    sdkplugin.ErrPluginShutdown,
		},
		config: &logical.BackendConfig{
			Config: map[string]string{
				"plugin_name":    "test-plugin",
				"plugin_version": "v5.1.2",
			},
		},
	}

	resp, err := b.HandleRequest(context.Background(), &logical.Request{Storage: &logical.InmemStorage{}})
	if resp != nil {
		t.Fatal("expected nil response")
	}
	if err == nil {
		t.Fatal("expected error")
	}

	codedErr, ok := err.(logical.HTTPCodedError)
	if !ok {
		t.Fatalf("expected logical.HTTPCodedError, got %T", err)
	}
	if codedErr.Code() != 503 {
		t.Fatalf("expected status code 503, got %d", codedErr.Code())
	}
	if !strings.Contains(err.Error(), pluginReloadInProgressErr) {
		t.Fatalf("expected error to mention %q, got %q", pluginReloadInProgressErr, err.Error())
	}
	if !strings.Contains(err.Error(), "test-plugin") {
		t.Fatalf("expected error to include plugin name, got %q", err.Error())
	}
	if !strings.Contains(err.Error(), "version=v5.1.2") {
		t.Fatalf("expected error to include plugin version, got %q", err.Error())
	}
}

func TestBackend_HandleExistenceCheck_BlocksReactiveReloadDuringReloadWindow(t *testing.T) {
	b := &backend{
		Backend: &testReloadingBackend{
			reloading: true,
			existsErr: sdkplugin.ErrPluginShutdown,
		},
		config: &logical.BackendConfig{
			Config: map[string]string{
				"plugin_name":    "test-plugin",
				"plugin_version": "v5.1.2",
			},
		},
	}

	checkFound, exists, err := b.HandleExistenceCheck(context.Background(), &logical.Request{Storage: &logical.InmemStorage{}})
	if checkFound || exists {
		t.Fatal("expected false checkFound/exists")
	}
	if err == nil {
		t.Fatal("expected error")
	}

	codedErr, ok := err.(logical.HTTPCodedError)
	if !ok {
		t.Fatalf("expected logical.HTTPCodedError, got %T", err)
	}
	if codedErr.Code() != 503 {
		t.Fatalf("expected status code 503, got %d", codedErr.Code())
	}
	if !strings.Contains(err.Error(), pluginReloadInProgressErr) {
		t.Fatalf("expected error to mention %q, got %q", pluginReloadInProgressErr, err.Error())
	}
	if !strings.Contains(err.Error(), "test-plugin") {
		t.Fatalf("expected error to include plugin name, got %q", err.Error())
	}
	if !strings.Contains(err.Error(), "version=v5.1.2") {
		t.Fatalf("expected error to include plugin version, got %q", err.Error())
	}
}

type testReloadingBackend struct {
	reloading bool
	reqErr    error
	existsErr error
}

func (b *testReloadingBackend) IsReloading() bool {
	return b.reloading
}

func (b *testReloadingBackend) Cleanup(context.Context) {}

func (b *testReloadingBackend) Setup(context.Context, *logical.BackendConfig) error {
	return nil
}

func (b *testReloadingBackend) Initialize(context.Context, *logical.InitializationRequest) error {
	return nil
}

func (b *testReloadingBackend) Type() logical.BackendType {
	return logical.TypeLogical
}

func (b *testReloadingBackend) SpecialPaths() *logical.Paths {
	return &logical.Paths{}
}

func (b *testReloadingBackend) HandleRequest(context.Context, *logical.Request) (*logical.Response, error) {
	return nil, b.reqErr
}

func (b *testReloadingBackend) HandleExistenceCheck(context.Context, *logical.Request) (bool, bool, error) {
	return false, false, b.existsErr
}

func (b *testReloadingBackend) InvalidateKey(context.Context, string) {}

func (b *testReloadingBackend) System() logical.SystemView {
	return nil
}

func (b *testReloadingBackend) Logger() hclog.Logger {
	return hclog.NewNullLogger()
}
