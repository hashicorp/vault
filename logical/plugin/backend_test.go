package plugin

import (
	"testing"
	"time"

	gplugin "github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/helper/logformat"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/plugin/mock"
	log "github.com/mgutz/logxi/v1"
)

func TestBackendPlugin_impl(t *testing.T) {
	var _ gplugin.Plugin = new(BackendPlugin)
	var _ logical.Backend = new(backendPluginClient)
}

func TestBackendPlugin_HandleRequest(t *testing.T) {
	b, cleanup := testBackend(t)
	defer cleanup()

	resp, err := b.HandleRequest(&logical.Request{
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

func TestBackendPlugin_SpecialPaths(t *testing.T) {
	b, cleanup := testBackend(t)
	defer cleanup()

	paths := b.SpecialPaths()
	if paths == nil {
		t.Fatal("SpecialPaths() returned nil")
	}
}

func TestBackendPlugin_System(t *testing.T) {
	b, cleanup := testBackend(t)
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

func TestBackendPlugin_Logger(t *testing.T) {
	b, cleanup := testBackend(t)
	defer cleanup()

	logger := b.Logger()
	if logger == nil {
		t.Fatal("Logger() returned nil")
	}
}

func TestBackendPlugin_HandleExistenceCheck(t *testing.T) {
	b, cleanup := testBackend(t)
	defer cleanup()

	checkFound, exists, err := b.HandleExistenceCheck(&logical.Request{
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

func TestBackendPlugin_Cleanup(t *testing.T) {
	b, cleanup := testBackend(t)
	defer cleanup()

	b.Cleanup()
}

func TestBackendPlugin_Initialize(t *testing.T) {
	b, cleanup := testBackend(t)
	defer cleanup()

	err := b.Initialize()
	if err != nil {
		t.Fatal(err)
	}
}

func TestBackendPlugin_InvalidateKey(t *testing.T) {
	b, cleanup := testBackend(t)
	defer cleanup()

	resp, err := b.HandleRequest(&logical.Request{
		Operation: logical.ReadOperation,
		Path:      "internal",
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Data["value"] == "" {
		t.Fatalf("bad: %#v, expected non-empty value", resp)
	}

	b.InvalidateKey("internal")

	resp, err = b.HandleRequest(&logical.Request{
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

func TestBackendPlugin_Setup(t *testing.T) {
	_, cleanup := testBackend(t)
	defer cleanup()
}

func testBackend(t *testing.T) (logical.Backend, func()) {
	// Create a mock provider
	pluginMap := map[string]gplugin.Plugin{
		"backend": &BackendPlugin{
			Factory: mock.Factory,
		},
	}
	client, _ := gplugin.TestPluginRPCConn(t, pluginMap)
	cleanup := func() {
		client.Close()
	}

	// Request the backend
	raw, err := client.Dispense(BackendPluginName)
	if err != nil {
		t.Fatal(err)
	}
	b := raw.(logical.Backend)

	err = b.Setup(&logical.BackendConfig{
		Logger: logformat.NewVaultLogger(log.LevelTrace),
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
