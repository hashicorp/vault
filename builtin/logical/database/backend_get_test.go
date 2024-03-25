// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package database

import (
	"context"
	"sync"
	"testing"

	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/queue"
)

func newSystemViewWrapper(view logical.SystemView) logical.SystemView {
	return &systemViewWrapper{
		view,
	}
}

type systemViewWrapper struct {
	logical.SystemView
}

var _ logical.ExtendedSystemView = (*systemViewWrapper)(nil)

func (s *systemViewWrapper) RequestWellKnownRedirect(ctx context.Context, src, dest string) error {
	panic("nope")
}

func (s *systemViewWrapper) DeregisterWellKnownRedirect(ctx context.Context, src string) bool {
	panic("nope")
}

func (s *systemViewWrapper) Auditor() logical.Auditor {
	panic("nope")
}

func (s *systemViewWrapper) ForwardGenericRequest(ctx context.Context, request *logical.Request) (*logical.Response, error) {
	panic("nope")
}

func (s *systemViewWrapper) APILockShouldBlockRequest() (bool, error) {
	panic("nope")
}

func (s *systemViewWrapper) GetPinnedPluginVersion(ctx context.Context, pluginType consts.PluginType, pluginName string) (*pluginutil.PinnedVersion, error) {
	return nil, pluginutil.ErrPinnedVersionNotFound
}

func (s *systemViewWrapper) LookupPluginVersion(ctx context.Context, pluginName string, pluginType consts.PluginType, version string) (*pluginutil.PluginRunner, error) {
	return &pluginutil.PluginRunner{
		Name:           "mockv5",
		Type:           consts.PluginTypeDatabase,
		Builtin:        true,
		BuiltinFactory: New,
	}, nil
}

func getDbBackend(t *testing.T) (*databaseBackend, logical.Storage) {
	t.Helper()
	config := logical.TestBackendConfig()
	config.System = newSystemViewWrapper(config.System)
	config.StorageView = &logical.InmemStorage{}
	// Create and init the backend ourselves instead of using a Factory because
	// the factory function kicks off threads that cause racy tests.
	b := Backend(config)
	if err := b.Setup(context.Background(), config); err != nil {
		t.Fatal(err)
	}
	b.schedule = &TestSchedule{}
	b.credRotationQueue = queue.New()
	b.populateQueue(context.Background(), config.StorageView)

	return b, config.StorageView
}

// TestGetConnectionRaceCondition checks that GetConnection always returns the same instance, even when asked
// by multiple goroutines in parallel.
func TestGetConnectionRaceCondition(t *testing.T) {
	ctx := context.Background()
	b, s := getDbBackend(t)
	defer b.Cleanup(ctx)
	configureDBMount(t, s)

	goroutines := 16

	wg := sync.WaitGroup{}
	wg.Add(goroutines)
	dbis := make([]*dbPluginInstance, goroutines)
	errs := make([]error, goroutines)
	for i := 0; i < goroutines; i++ {
		go func(i int) {
			dbis[i], errs[i] = b.GetConnection(ctx, s, "mockv5")
			wg.Done()
		}(i)
	}
	wg.Wait()
	for i := 0; i < goroutines; i++ {
		if errs[i] != nil {
			t.Fatal(errs[i])
		}
		if dbis[0] != dbis[i] {
			t.Fatal("Error: database instances did not match")
		}
	}
}
