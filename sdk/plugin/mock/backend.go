// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mock

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	MockPluginVersionEnv           = "TESTING_MOCK_VAULT_PLUGIN_VERSION"
	MockPluginDefaultInternalValue = "bar"
)

// New returns a new backend as an interface. This func
// is only necessary for builtin backend plugins.
func New() (interface{}, error) {
	return Backend(), nil
}

// Factory returns a new backend as logical.Backend.
func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := Backend()
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	return b, nil
}

// FactoryType is a wrapper func that allows the Factory func to specify
// the backend type for the mock backend plugin instance.
func FactoryType(backendType logical.BackendType) logical.Factory {
	return func(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
		b := Backend()
		b.BackendType = backendType
		if err := b.Setup(ctx, conf); err != nil {
			return nil, err
		}
		return b, nil
	}
}

// Backend returns a private embedded struct of framework.Backend.
func Backend() *backend {
	var b backend
	b.Backend = &framework.Backend{
		Help: "",
		Paths: framework.PathAppend(
			errorPaths(&b),
			kvPaths(&b),
			[]*framework.Path{
				pathInternal(&b),
				pathSpecial(&b),
				pathRaw(&b),
				pathEnv(&b),
			},
		),
		PathsSpecial: &logical.Paths{
			Unauthenticated: []string{
				"special",
			},
		},
		Secrets:     []*framework.Secret{},
		Invalidate:  b.invalidate,
		BackendType: logical.TypeLogical,
	}
	b.internal = MockPluginDefaultInternalValue
	b.RunningVersion = "v0.0.0+mock"
	if version := os.Getenv(MockPluginVersionEnv); version != "" {
		b.RunningVersion = version
	}
	return &b
}

type backend struct {
	*framework.Backend

	// internal is used to test invalidate and reloads.
	internal string
}

func (b *backend) invalidate(ctx context.Context, key string) {
	switch key {
	case "internal":
		b.internal = ""
	}
}

// WriteInternalValue is a helper to set an in-memory value in the plugin,
// allowing tests to later assert that the plugin either has or hasn't been
// restarted.
func WriteInternalValue(t *testing.T, client *api.Client, mountPath, value string) {
	t.Helper()
	resp, err := client.Logical().Write(fmt.Sprintf("%s/internal", mountPath), map[string]interface{}{
		"value": value,
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %v", resp)
	}
}

// ExpectInternalValue checks the internal in-memory value.
func ExpectInternalValue(t *testing.T, client *api.Client, mountPath, expected string) {
	t.Helper()
	expectInternalValue(t, client, mountPath, expected)
}

func expectInternalValue(t *testing.T, client *api.Client, mountPath, expected string) {
	t.Helper()
	resp, err := client.Logical().Read(fmt.Sprintf("%s/internal", mountPath))
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil {
		t.Fatalf("bad: response should not be nil")
	}
	if resp.Data["value"].(string) != expected {
		t.Fatalf("expected %q but got %q", expected, resp.Data["value"].(string))
	}
}
