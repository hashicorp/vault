package vault_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/vault/builtin/plugin"
	"github.com/hashicorp/vault/helper/logformat"
	"github.com/hashicorp/vault/helper/pluginutil"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/logical"
	lplugin "github.com/hashicorp/vault/logical/plugin"
	"github.com/hashicorp/vault/logical/plugin/mock"
	"github.com/hashicorp/vault/vault"
	log "github.com/mgutz/logxi/v1"
)

func TestSystemBackend_enableAuth_plugin(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"plugin": plugin.Factory,
		},
	}

	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()
	core := cluster.Cores[0].Core
	vault.TestWaitActive(t, core)

	b := vault.NewSystemBackend(core)
	logger := logformat.NewVaultLogger(log.LevelTrace)
	bc := &logical.BackendConfig{
		Logger: logger,
		System: logical.StaticSystemView{
			DefaultLeaseTTLVal: time.Hour * 24,
			MaxLeaseTTLVal:     time.Hour * 24 * 32,
		},
	}

	err := b.Backend.Setup(bc)
	if err != nil {
		t.Fatal(err)
	}

	os.Setenv(pluginutil.PluginCACertPEMEnv, cluster.CACertPEMFile)

	vault.TestAddTestPlugin(t, core, "mock-plugin", "TestBackend_PluginMainCredentials")

	req := logical.TestRequest(t, logical.UpdateOperation, "auth/mock-plugin")
	req.Data["type"] = "plugin"
	req.Data["plugin_name"] = "mock-plugin"

	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %v", resp)
	}
}

func TestSystemBackend_PluginReload(t *testing.T) {
	data := map[string]interface{}{
		"plugin": "mock-plugin",
	}
	t.Run("plugin", func(t *testing.T) { testSystemBackend_PluginReload(t, data) })

	data = map[string]interface{}{
		"mounts": "mock-0/,mock-1/",
	}
	t.Run("mounts", func(t *testing.T) { testSystemBackend_PluginReload(t, data) })
}

func testSystemBackend_PluginReload(t *testing.T, reqData map[string]interface{}) {
	cluster, b := testSystemBackendMock(t, 2)
	defer cluster.Cleanup()

	core := cluster.Cores[0]

	for i := 0; i < 2; i++ {
		// Update internal value in the backend
		req := logical.TestRequest(t, logical.UpdateOperation, fmt.Sprintf("mock-%d/internal", i))
		req.ClientToken = core.Client.Token()
		req.Data["value"] = "baz"
		resp, err := core.HandleRequest(req)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if resp != nil {
			t.Fatalf("bad: %v", resp)
		}
	}

	// Perform plugin reload
	req := logical.TestRequest(t, logical.UpdateOperation, "plugins/backend/reload")
	req.ClientToken = core.Client.Token()
	req.Data = reqData
	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %v", resp)
	}

	for i := 0; i < 2; i++ {
		// Ensure internal backed value is reset
		req := logical.TestRequest(t, logical.ReadOperation, "mock-1/internal")
		req.ClientToken = core.Client.Token()
		resp, err := core.HandleRequest(req)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if resp == nil {
			t.Fatalf("bad: response should not be nil")
		}
		if resp.Data["value"].(string) == "baz" {
			t.Fatal("did not expect backend internal value to be 'baz'")
		}
	}
}

// testSystemBackendMock returns a systemBackend with the desired number
// of mounted mock plugin backends
func testSystemBackendMock(t *testing.T, numMounts int) (*vault.TestCluster, *vault.SystemBackend) {
	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"plugin": plugin.Factory,
		},
	}

	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()

	core := cluster.Cores[0].Core
	vault.TestWaitActive(t, core)

	b := vault.NewSystemBackend(core)
	logger := logformat.NewVaultLogger(log.LevelTrace)
	bc := &logical.BackendConfig{
		Logger: logger,
		System: logical.StaticSystemView{
			DefaultLeaseTTLVal: time.Hour * 24,
			MaxLeaseTTLVal:     time.Hour * 24 * 32,
		},
	}

	err := b.Backend.Setup(bc)
	if err != nil {
		t.Fatal(err)
	}

	os.Setenv(pluginutil.PluginCACertPEMEnv, cluster.CACertPEMFile)

	vault.TestAddTestPlugin(t, core, "mock-plugin", "TestBackend_PluginMainLogical")

	for i := 0; i < numMounts; i++ {
		req := logical.TestRequest(t, logical.UpdateOperation, fmt.Sprintf("mounts/mock-%d/", i))
		req.Data["type"] = "plugin"
		req.Data["config"] = map[string]interface{}{
			"plugin_name": "mock-plugin",
		}

		resp, err := b.HandleRequest(req)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if resp != nil {
			t.Fatalf("bad: %v", resp)
		}
	}

	return cluster, b
}

func TestBackend_PluginMainLogical(t *testing.T) {
	if os.Getenv(pluginutil.PluginUnwrapTokenEnv) == "" {
		return
	}

	caPEM := os.Getenv(pluginutil.PluginCACertPEMEnv)
	if caPEM == "" {
		t.Fatal("CA cert not passed in")
	}

	factoryFunc := mock.FactoryType(logical.TypeLogical)

	args := []string{"--ca-cert=" + caPEM}

	apiClientMeta := &pluginutil.APIClientMeta{}
	flags := apiClientMeta.FlagSet()
	flags.Parse(args)
	tlsConfig := apiClientMeta.GetTLSConfig()
	tlsProviderFunc := pluginutil.VaultPluginTLSProvider(tlsConfig)
	err := lplugin.Serve(&lplugin.ServeOpts{
		BackendFactoryFunc: factoryFunc,
		TLSProviderFunc:    tlsProviderFunc,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestBackend_PluginMainCredentials(t *testing.T) {
	if os.Getenv(pluginutil.PluginUnwrapTokenEnv) == "" {
		return
	}

	caPEM := os.Getenv(pluginutil.PluginCACertPEMEnv)
	if caPEM == "" {
		t.Fatal("CA cert not passed in")
	}

	factoryFunc := mock.FactoryType(logical.TypeCredential)

	args := []string{"--ca-cert=" + caPEM}

	apiClientMeta := &pluginutil.APIClientMeta{}
	flags := apiClientMeta.FlagSet()
	flags.Parse(args)
	tlsConfig := apiClientMeta.GetTLSConfig()
	tlsProviderFunc := pluginutil.VaultPluginTLSProvider(tlsConfig)
	err := lplugin.Serve(&lplugin.ServeOpts{
		BackendFactoryFunc: factoryFunc,
		TLSProviderFunc:    tlsProviderFunc,
	})
	if err != nil {
		t.Fatal(err)
	}
}
