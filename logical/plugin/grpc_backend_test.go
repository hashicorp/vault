package plugin

import (
	"context"
	"os"
	"testing"
	"time"

	log "github.com/hashicorp/go-hclog"
	gplugin "github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/helper/logging"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/plugin/mock"
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

func TestGRPCBackendPlugin_SpecialPaths(t *testing.T) {
	b, cleanup := testGRPCBackend(t)
	defer cleanup()

	paths := b.SpecialPaths()
	if paths == nil {
		t.Fatal("SpecialPaths() returned nil")
	}
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

func testGRPCBackend(t *testing.T) (logical.Backend, func()) {
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
	client, _ := gplugin.TestPluginGRPCConn(t, pluginMap)
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
		StorageView: &logical.InmemStorage{},
	})
	if err != nil {
		t.Fatal(err)
	}

	return b, cleanup
}
