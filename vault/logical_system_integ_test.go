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

func TestSystemBackend_PluginReload_plugin(t *testing.T) {
	b, cleanup := testSystemBackendMock(t, 2)
	defer cleanup()

	req := logical.TestRequest(t, logical.UpdateOperation, "plugin/backend/reload")
	req.Data["plugin"] = "mock-plugin"
	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %v", resp)
	}
}
func TestSystemBackend_PluginReload_mounts(t *testing.T) {
	b, cleanup := testSystemBackendMock(t, 2)
	defer cleanup()

	req := logical.TestRequest(t, logical.UpdateOperation, "plugin/backend/reload")
	req.Data["mounts"] = "mock-0/,mock-1/"
	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %v", resp)
	}
}

// testSystemBackendMock returns a systemBackend, with deired number
// of mounted mock plugin backends
func testSystemBackendMock(t *testing.T, numMounts int) (b *vault.SystemBackend, cleanup func()) {
	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"plugin": plugin.Factory,
		},
	}

	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	cleanup = func() {
		cluster.Cleanup()
	}

	core := cluster.Cores[0].Core
	vault.TestWaitActive(t, core)

	b = vault.NewSystemBackend(core)
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

	// Mount plugin in two mount points
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

	return b, cleanup
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
