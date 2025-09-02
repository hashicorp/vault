// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plugin

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	log "github.com/hashicorp/go-hclog"
	gplugin "github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/plugin/mock"
	"github.com/stretchr/testify/require"
)

func TestGRPCBackendPlugin_impl(t *testing.T) {
	var _ gplugin.Plugin = new(GRPCBackendPlugin)
	var _ logical.Backend = new(backendGRPCPluginClient)
}

func TestGRPCBackendPlugin_HandleRequest(t *testing.T) {
	b, cleanup := testGRPCBackend(t)
	defer cleanup()

	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "kv/foo",
		Data: map[string]interface{}{
			"value": "bar",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Data["value"] != "bar" {
		t.Fatalf("bad: %#v", resp)
	}
}

// TestGRPCBackendPlugin_SpecialPaths verifies that the special paths returned
// by the plugin are not nil and that the AllowSnapshotRead path is not empty.
func TestGRPCBackendPlugin_SpecialPaths(t *testing.T) {
	b, cleanup := testGRPCBackend(t)
	defer cleanup()

	paths := b.SpecialPaths()
	if paths == nil {
		t.Fatal("SpecialPaths() returned nil")
	}
	if len(paths.AllowSnapshotRead) == 0 {
		t.Fatalf("SpecialPaths() returned empty AllowSnapshotRead")
	}
}

type requireSnapshotCtxStorage struct {
	logical.Storage
	snapshotID string
}

func (r *requireSnapshotCtxStorage) errorIfNoSnapshotID(ctx context.Context) error {
	if snapshotID, _ := logical.ContextSnapshotIDValue(ctx); snapshotID != r.snapshotID {
		return fmt.Errorf("missing snapshot ID")
	}
	return nil
}

func (r *requireSnapshotCtxStorage) Get(ctx context.Context, key string) (*logical.StorageEntry, error) {
	if err := r.errorIfNoSnapshotID(ctx); err != nil {
		return nil, err
	}
	return r.Storage.Get(ctx, key)
}

func (r *requireSnapshotCtxStorage) List(ctx context.Context, prefix string) ([]string, error) {
	if err := r.errorIfNoSnapshotID(ctx); err != nil {
		return nil, err
	}
	return r.Storage.List(ctx, prefix)
}

func (r *requireSnapshotCtxStorage) Delete(ctx context.Context, key string) error {
	if err := r.errorIfNoSnapshotID(ctx); err != nil {
		return err
	}
	return r.Storage.Delete(ctx, key)
}

func (r *requireSnapshotCtxStorage) Put(ctx context.Context, entry *logical.StorageEntry) error {
	if err := r.errorIfNoSnapshotID(ctx); err != nil {
		return err
	}
	return r.Storage.Put(ctx, entry)
}

// TestGRPCBackendPlugin_SnapshotCtx verifies that the backend plugin correctly
// parses the snapshot ID context key and passes it to the storage methods.
func TestGRPCBackendPlugin_SnapshotCtx(t *testing.T) {
	storage := &requireSnapshotCtxStorage{Storage: &logical.InmemStorage{}, snapshotID: "abcd"}
	b, cleanup := testGRPCBackendWithStorage(t, storage)
	defer cleanup()
	ctx := logical.CreateContextWithSnapshotID(context.Background(), "abcd")
	// check put storage method
	resp, err := b.HandleRequest(ctx, &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "kv/foo",
		Data: map[string]interface{}{
			"value": "bar",
		},
	})
	require.NoError(t, err)
	require.NoError(t, resp.Error())

	// check get storage method
	resp, err = b.HandleRequest(ctx, &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "kv/foo",
	})
	require.NoError(t, err)
	require.NoError(t, resp.Error())

	// check list storage method
	resp, err = b.HandleRequest(ctx, &logical.Request{
		Operation: logical.ListOperation,
		Path:      "kv",
	})
	require.NoError(t, err)
	require.NoError(t, resp.Error())

	// check delete storage method
	resp, err = b.HandleRequest(ctx, &logical.Request{
		Operation: logical.DeleteOperation,
		Path:      "kv/foo",
	})
	require.NoError(t, err)
	require.NoError(t, resp.Error())
}

func TestGRPCBackendPlugin_System(t *testing.T) {
	b, cleanup := testGRPCBackend(t)
	defer cleanup()

	sys := b.System()
	if sys == nil {
		t.Fatal("System() returned nil")
	}

	actual := sys.DefaultLeaseTTL()
	expected := 300 * time.Second

	if actual != expected {
		t.Fatalf("bad: %v, expected %v", actual, expected)
	}
}

func TestGRPCBackendPlugin_Logger(t *testing.T) {
	b, cleanup := testGRPCBackend(t)
	defer cleanup()

	logger := b.Logger()
	if logger == nil {
		t.Fatal("Logger() returned nil")
	}
}

func TestGRPCBackendPlugin_HandleExistenceCheck(t *testing.T) {
	b, cleanup := testGRPCBackend(t)
	defer cleanup()

	checkFound, exists, err := b.HandleExistenceCheck(context.Background(), &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "kv/foo",
		Data:      map[string]interface{}{"value": "bar"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !checkFound {
		t.Fatal("existence check not found for path 'kv/foo")
	}
	if exists {
		t.Fatal("existence check should have returned 'false' for 'kv/foo'")
	}
}

func TestGRPCBackendPlugin_Cleanup(t *testing.T) {
	b, cleanup := testGRPCBackend(t)
	defer cleanup()

	b.Cleanup(context.Background())
}

func TestGRPCBackendPlugin_InvalidateKey(t *testing.T) {
	b, cleanup := testGRPCBackend(t)
	defer cleanup()

	ctx := context.Background()

	resp, err := b.HandleRequest(ctx, &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "internal",
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Data["value"] == "" {
		t.Fatalf("bad: %#v, expected non-empty value", resp)
	}

	b.InvalidateKey(ctx, "internal")

	resp, err = b.HandleRequest(ctx, &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "internal",
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Data["value"] != "" {
		t.Fatalf("bad: expected empty response data, got %#v", resp)
	}
}

func TestGRPCBackendPlugin_Setup(t *testing.T) {
	_, cleanup := testGRPCBackend(t)
	defer cleanup()
}

func TestGRPCBackendPlugin_Initialize(t *testing.T) {
	b, cleanup := testGRPCBackend(t)
	defer cleanup()

	err := b.Initialize(context.Background(), &logical.InitializationRequest{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestGRPCBackendPlugin_Version(t *testing.T) {
	b, cleanup := testGRPCBackend(t)
	defer cleanup()

	versioner, ok := b.(logical.PluginVersioner)
	if !ok {
		t.Fatalf("Expected %T to implement logical.PluginVersioner interface", b)
	}

	version := versioner.PluginVersion().Version
	if version != "v0.0.0+mock" {
		t.Fatalf("Got version %s, expected 'mock'", version)
	}
}

func testGRPCBackend(t *testing.T) (logical.Backend, func()) {
	return testGRPCBackendWithStorage(t, &logical.InmemStorage{})
}

func testGRPCBackendWithStorage(t *testing.T, storage logical.Storage) (logical.Backend, func()) {
	// Create a mock provider
	pluginMap := map[string]gplugin.Plugin{
		"backend": &GRPCBackendPlugin{
			Factory: mock.Factory,
			Logger: log.New(&log.LoggerOptions{
				Level:      log.Debug,
				Output:     os.Stderr,
				JSONFormat: true,
			}),
		},
	}
	client, _ := gplugin.TestPluginGRPCConn(t, false, pluginMap)
	cleanup := func() {
		client.Close()
	}

	// Request the backend
	raw, err := client.Dispense(BackendPluginName)
	if err != nil {
		t.Fatal(err)
	}
	b := raw.(logical.Backend)

	err = b.Setup(context.Background(), &logical.BackendConfig{
		Logger: logging.NewVaultLogger(log.Debug),
		System: &logical.StaticSystemView{
			DefaultLeaseTTLVal: 300 * time.Second,
			MaxLeaseTTLVal:     1800 * time.Second,
		},
		StorageView: storage,
	})
	if err != nil {
		t.Fatal(err)
	}

	return b, cleanup
}
