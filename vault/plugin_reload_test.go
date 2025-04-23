// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
)

type goodPluginTestBackend struct{}

func (g goodPluginTestBackend) Initialize(ctx context.Context, request *logical.InitializationRequest) error {
	return nil
}

func (g goodPluginTestBackend) HandleRequest(ctx context.Context, request *logical.Request) (*logical.Response, error) {
	return nil, nil
}
func (g goodPluginTestBackend) SpecialPaths() *logical.Paths { return nil }
func (g goodPluginTestBackend) System() logical.SystemView   { return nil }
func (g goodPluginTestBackend) Logger() hclog.Logger         { return nil }
func (g goodPluginTestBackend) HandleExistenceCheck(ctx context.Context, request *logical.Request) (bool, bool, error) {
	return false, false, nil
}
func (g goodPluginTestBackend) Cleanup(ctx context.Context)                 {}
func (g goodPluginTestBackend) InvalidateKey(ctx context.Context, s string) {}
func (g goodPluginTestBackend) Setup(ctx context.Context, config *logical.BackendConfig) error {
	return nil
}
func (g goodPluginTestBackend) Type() logical.BackendType { return logical.TypeLogical }

type badPluginTestBackend struct{}

func (b badPluginTestBackend) Initialize(ctx context.Context, r *logical.InitializationRequest) error {
	se, err := r.Storage.Get(ctx, "tag")
	if err != nil {
		return fmt.Errorf("couldn't read storage: %s", err)
	}
	if se == nil {
		r.Storage.Put(ctx, &logical.StorageEntry{
			Key:   "tag",
			Value: []byte("boo"),
		})
		return nil
	}

	return errors.New("already initialized")
}

func (b badPluginTestBackend) HandleRequest(_ context.Context, _ *logical.Request) (*logical.Response, error) {
	return nil, nil
}
func (b badPluginTestBackend) SpecialPaths() *logical.Paths { return nil }
func (b badPluginTestBackend) System() logical.SystemView   { return nil }
func (b badPluginTestBackend) Logger() hclog.Logger         { return nil }
func (b badPluginTestBackend) HandleExistenceCheck(ctx context.Context, request *logical.Request) (bool, bool, error) {
	return false, false, nil
}
func (b badPluginTestBackend) Cleanup(ctx context.Context)                 {}
func (b badPluginTestBackend) InvalidateKey(ctx context.Context, s string) {}
func (b badPluginTestBackend) Setup(ctx context.Context, config *logical.BackendConfig) error {
	return nil
}
func (b badPluginTestBackend) Type() logical.BackendType { return logical.TypeLogical }

type BadBackend struct{}

// TestReload vets the reload functionality of Core by adding plugins of varying kinds
// and finds out if they reload as expected.
func TestReload(t *testing.T) {
	ctx := context.Background()
	nsCTX := namespace.RootContext(ctx)

	core, _, _ := TestCoreUnsealed(t)

	cases := []struct {
		name        string
		backends    map[string]logical.Factory
		errCount    int
		errPrefixes []string
	}{
		{
			name: "reload one",
			backends: map[string]logical.Factory{
				"work": func(ctx context.Context, config *logical.BackendConfig) (logical.Backend, error) {
					return goodPluginTestBackend{}, nil
				},
			},
		},
		{
			name: "reload two",
			backends: map[string]logical.Factory{
				"work": func(ctx context.Context, config *logical.BackendConfig) (logical.Backend, error) {
					return goodPluginTestBackend{}, nil
				},
				"work2": func(_ context.Context, _ *logical.BackendConfig) (logical.Backend, error) {
					return goodPluginTestBackend{}, nil
				},
			},
		},
		{
			name:     "two broken",
			errCount: 2,
			backends: map[string]logical.Factory{
				"bad1": func(_ context.Context, _ *logical.BackendConfig) (logical.Backend, error) {
					return badPluginTestBackend{}, nil
				},
				"bad2": func(_ context.Context, _ *logical.BackendConfig) (logical.Backend, error) {
					return badPluginTestBackend{}, nil
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var mes []*MountEntry
			var names []string

			for k, v := range c.backends {
				core.logicalBackends[k] = v
				// meUUID, _ := uuid.GenerateUUID()
				mes = append(mes, &MountEntry{
					Path:  k,
					Type:  k,
					Table: mountTableType,
				})
				names = append(names, k)
			}

			for _, m := range mes {
				err := core.mount(nsCTX, m)
				if err != nil {
					t.Fatalf("%s", err)
				}
			}

			err := core.reloadMatchingPluginMounts(nsCTX, namespace.RootNamespace, names)
			if c.errCount == 0 {
				if err != nil {
					t.Fatalf("expected no errors but got: %s", err)
				}
				// otherwise good
			} else {
				// expected error
				if err == nil {
					t.Fatal("expected an error but got none")
				} else {
					var merr *multierror.Error
					errors.As(err, &merr)
					if merr.Len() != c.errCount {
						t.Fatalf("didn't get the right number of errors, expected %d but got %d", c.errCount, merr.Len())
					}
				}
			}
		})
	}
}
