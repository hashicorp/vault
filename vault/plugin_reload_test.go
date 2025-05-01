// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/logical"
)

const reloadFailurePrefix = "cannot reload"

// A goodReloadPluginTestBackend does nothing but always succeeds at Initialization (how good!)
//
// The name is long to avoid collisions in the vault package.
type goodReloadPluginTestBackend struct{}

func (g goodReloadPluginTestBackend) Initialize(ctx context.Context, request *logical.InitializationRequest) error {
	return nil
}

func (g goodReloadPluginTestBackend) HandleRequest(ctx context.Context, request *logical.Request) (*logical.Response, error) {
	return nil, nil
}
func (g goodReloadPluginTestBackend) SpecialPaths() *logical.Paths { return nil }
func (g goodReloadPluginTestBackend) System() logical.SystemView   { return nil }
func (g goodReloadPluginTestBackend) Logger() hclog.Logger         { return nil }
func (g goodReloadPluginTestBackend) HandleExistenceCheck(ctx context.Context, request *logical.Request) (bool, bool, error) {
	return false, false, nil
}
func (g goodReloadPluginTestBackend) Cleanup(ctx context.Context)                 {}
func (g goodReloadPluginTestBackend) InvalidateKey(ctx context.Context, s string) {}
func (g goodReloadPluginTestBackend) Setup(ctx context.Context, config *logical.BackendConfig) error {
	return nil
}
func (g goodReloadPluginTestBackend) Type() logical.BackendType { return logical.TypeLogical }

// An incrementallyBadReloadPluginTestBackend is a plugin that will fail after a specified number of Initialize calls.
// The first time it is initialized (at a given path), the check value is put in storage and persisted. When the plugin
// has been initialized that many times, the next call will fail. Note that every plugin needs to succeed at initialization
// once, or it will just fail to mount.
type incrementallyBadReloadPluginTestBackend struct {
	check byte
}

func (i incrementallyBadReloadPluginTestBackend) Initialize(ctx context.Context, r *logical.InitializationRequest) error {
	countEntry, err := r.Storage.Get(ctx, "count")
	if err != nil {
		return fmt.Errorf("failed to initialize for the wrong reason: %s", err)
	}
	if countEntry == nil {
		r.Storage.Put(ctx, &logical.StorageEntry{
			Key:   "count",
			Value: []byte{1},
		})
		// also write check to storage, because this is our real value
		r.Storage.Put(ctx, &logical.StorageEntry{
			Key:   "check",
			Value: []byte{i.check},
		})
		return nil
	}

	// read check value from storage
	checkEntry, err := r.Storage.Get(ctx, "check")
	if err != nil {
		return fmt.Errorf("failed to initialize for the wrong reason: %s", err)
	}

	i.check = checkEntry.Value[0]

	if countEntry.Value[0] >= i.check {
		return fmt.Errorf("initialized error requested, %d was greater than or equal to %d", countEntry.Value[0], i.check)
	}
	r.Storage.Put(ctx, &logical.StorageEntry{
		Key:   "count",
		Value: []byte{countEntry.Value[0] + 1},
	})

	return nil
}

func (i incrementallyBadReloadPluginTestBackend) HandleRequest(ctx context.Context, r *logical.Request) (*logical.Response, error) {
	return nil, nil
}

func (i incrementallyBadReloadPluginTestBackend) SpecialPaths() *logical.Paths {
	return nil
}

func (i incrementallyBadReloadPluginTestBackend) System() logical.SystemView {
	return nil
}

func (i incrementallyBadReloadPluginTestBackend) Logger() hclog.Logger {
	return nil
}

func (i incrementallyBadReloadPluginTestBackend) HandleExistenceCheck(ctx context.Context, request *logical.Request) (bool, bool, error) {
	return false, false, nil
}

func (i incrementallyBadReloadPluginTestBackend) Cleanup(ctx context.Context)                 {}
func (i incrementallyBadReloadPluginTestBackend) InvalidateKey(ctx context.Context, s string) {}
func (i incrementallyBadReloadPluginTestBackend) Setup(ctx context.Context, config *logical.BackendConfig) error {
	return nil
}

func (i incrementallyBadReloadPluginTestBackend) Type() logical.BackendType {
	return logical.TypeLogical
}

// TestReloadMatchingMounts vets the reload functionality of Core by adding plugins of varying kinds
// and finds out if they reload as expected.
func TestReloadMatchingMounts(t *testing.T) {
	ctx := context.Background()
	nsCTX := namespace.RootContext(ctx)

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
					return goodReloadPluginTestBackend{}, nil
				},
			},
		},
		{
			name: "reload two",
			backends: map[string]logical.Factory{
				"work": func(ctx context.Context, config *logical.BackendConfig) (logical.Backend, error) {
					return goodReloadPluginTestBackend{}, nil
				},
				"work2": func(_ context.Context, _ *logical.BackendConfig) (logical.Backend, error) {
					return goodReloadPluginTestBackend{}, nil
				},
			},
		},
		{
			name:     "two broken",
			errCount: 2,
			backends: map[string]logical.Factory{
				"bad1": func(_ context.Context, _ *logical.BackendConfig) (logical.Backend, error) {
					return incrementallyBadReloadPluginTestBackend{1}, nil
				},
				"bad2": func(_ context.Context, _ *logical.BackendConfig) (logical.Backend, error) {
					return incrementallyBadReloadPluginTestBackend{1}, nil
				},
			},
		},
		{
			name:     "one broken, one good",
			errCount: 1,
			backends: map[string]logical.Factory{
				"bad": func(_ context.Context, _ *logical.BackendConfig) (logical.Backend, error) {
					return incrementallyBadReloadPluginTestBackend{1}, nil
				},
				"good": func(_ context.Context, _ *logical.BackendConfig) (logical.Backend, error) {
					return goodReloadPluginTestBackend{}, nil
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			core, _, _ := TestCoreUnsealed(t)
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
					for _, e := range merr.Errors {
						if !strings.HasPrefix(e.Error(), reloadFailurePrefix) {
							t.Fatalf("got a different error than expected: %s", e)
						}
					}
				}
			}
		})
	}
}

// TestProgressiveReloadErrorsByPluginType tests core.reloadMatchingPlugin by creating a plugin of a single type
// that fails after 1, 2, and 3 initializations (mounting is the first).
func TestProgressiveReloadErrorsByPluginType(t *testing.T) {
	nsCTX := namespace.RootContext(nil)

	core, _, _ := TestCoreUnsealed(t)

	failBar := byte(1)

	core.logicalBackends["incr"] = func(_ context.Context, _ *logical.BackendConfig) (logical.Backend, error) {
		b := incrementallyBadReloadPluginTestBackend{
			check: failBar,
		}
		failBar++
		return b, nil
	}

	core.mount(nsCTX, &MountEntry{
		Path:  "a1",
		Type:  "incr",
		Table: mountTableType,
	})
	core.mount(nsCTX, &MountEntry{
		Path:  "a2",
		Type:  "incr",
		Table: mountTableType,
	})
	core.mount(nsCTX, &MountEntry{
		Path:  "a3",
		Type:  "incr",
		Table: mountTableType,
	})

	num, err := core.reloadMatchingPlugin(nsCTX, namespace.RootNamespace, consts.PluginTypeSecrets, "incr")
	var merr *multierror.Error
	errors.As(err, &merr)
	if num != 2 {
		t.Fatalf("expected 2 successes but got %d", num)
	}
	if merr.Len() != 1 {
		t.Fatalf("expected 1 reload error but got %d", merr.Len())
	}

	num, err = core.reloadMatchingPlugin(nsCTX, namespace.RootNamespace, consts.PluginTypeSecrets, "incr")
	errors.As(err, &merr)
	if num != 1 {
		t.Fatalf("expected 2 successes but got %d", num)
	}
	if merr.Len() != 2 {
		t.Fatalf("expected 1 reload error but got %d", merr.Len())
	}

	num, err = core.reloadMatchingPlugin(nsCTX, namespace.RootNamespace, consts.PluginTypeSecrets, "incr")
	errors.As(err, &merr)
	if num != 0 {
		t.Fatalf("expected 2 successes but got %d", num)
	}
	if merr.Len() != 3 {
		t.Fatalf("expected 1 reload error but got %d", merr.Len())
	}
	for _, e := range merr.Errors {
		if !strings.HasPrefix(e.Error(), reloadFailurePrefix) {
			t.Fatalf("got a different error than expected: %s", e)
		}
	}
}
