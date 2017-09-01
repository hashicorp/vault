package vault_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/vault/builtin/plugin"
	"github.com/hashicorp/vault/helper/pluginutil"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/logical"
	lplugin "github.com/hashicorp/vault/logical/plugin"
	"github.com/hashicorp/vault/logical/plugin/mock"
	"github.com/hashicorp/vault/vault"
)

func TestSystemBackend_Plugin_secret(t *testing.T) {
	cluster := testSystemBackendMock(t, 1, 1, logical.TypeLogical)
	defer cluster.Cleanup()

	core := cluster.Cores[0]

	// Make a request to lazy load the plugin
	req := logical.TestRequest(t, logical.ReadOperation, "mock-0/internal")
	req.ClientToken = core.Client.Token()
	resp, err := core.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil {
		t.Fatalf("bad: response should not be nil")
	}

	// Seal the cluster
	cluster.EnsureCoresSealed(t)

	// Unseal the cluster
	barrierKeys := cluster.BarrierKeys
	for _, core := range cluster.Cores {
		for _, key := range barrierKeys {
			_, err := core.Unseal(vault.TestKeyCopy(key))
			if err != nil {
				t.Fatal(err)
			}
		}
		sealed, err := core.Sealed()
		if err != nil {
			t.Fatalf("err checking seal status: %s", err)
		}
		if sealed {
			t.Fatal("should not be sealed")
		}
		// Wait for active so post-unseal takes place
		// If it fails, it means unseal process failed
		vault.TestWaitActive(t, core.Core)
	}
}

func TestSystemBackend_Plugin_auth(t *testing.T) {
	cluster := testSystemBackendMock(t, 1, 1, logical.TypeCredential)
	defer cluster.Cleanup()

	core := cluster.Cores[0]

	// Make a request to lazy load the plugin
	req := logical.TestRequest(t, logical.ReadOperation, "auth/mock-0/internal")
	req.ClientToken = core.Client.Token()
	resp, err := core.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil {
		t.Fatalf("bad: response should not be nil")
	}

	// Seal the cluster
	cluster.EnsureCoresSealed(t)

	// Unseal the cluster
	barrierKeys := cluster.BarrierKeys
	for _, core := range cluster.Cores {
		for _, key := range barrierKeys {
			_, err := core.Unseal(vault.TestKeyCopy(key))
			if err != nil {
				t.Fatal(err)
			}
		}
		sealed, err := core.Sealed()
		if err != nil {
			t.Fatalf("err checking seal status: %s", err)
		}
		if sealed {
			t.Fatal("should not be sealed")
		}
		// Wait for active so post-unseal takes place
		// If it fails, it means unseal process failed
		vault.TestWaitActive(t, core.Core)
	}
}

func TestSystemBackend_Plugin_MismatchType(t *testing.T) {
	cluster := testSystemBackendMock(t, 1, 1, logical.TypeLogical)
	defer cluster.Cleanup()

	core := cluster.Cores[0]

	// Replace the plugin with a credential backend
	vault.TestAddTestPlugin(t, core.Core, "mock-plugin", "TestBackend_PluginMainCredentials")

	// Make a request to lazy load the now-credential plugin
	// and expect an error
	req := logical.TestRequest(t, logical.ReadOperation, "mock-0/internal")
	req.ClientToken = core.Client.Token()
	_, err := core.HandleRequest(req)
	if err == nil {
		t.Fatalf("expected error due to mismatch on error type: %s", err)
	}

	// Sleep a bit before cleanup is called
	time.Sleep(1 * time.Second)
}

func TestSystemBackend_Plugin_CatalogRemoved(t *testing.T) {
	t.Run("secret", func(t *testing.T) {
		testPlugin_CatalogRemoved(t, logical.TypeLogical, false)
	})

	t.Run("auth", func(t *testing.T) {
		testPlugin_CatalogRemoved(t, logical.TypeCredential, false)
	})

	t.Run("secret-mount-existing", func(t *testing.T) {
		testPlugin_CatalogRemoved(t, logical.TypeLogical, true)
	})

	t.Run("auth-mount-existing", func(t *testing.T) {
		testPlugin_CatalogRemoved(t, logical.TypeCredential, true)
	})
}

func testPlugin_CatalogRemoved(t *testing.T, btype logical.BackendType, testMount bool) {
	cluster := testSystemBackendMock(t, 1, 1, btype)
	defer cluster.Cleanup()

	core := cluster.Cores[0]

	// Remove the plugin from the catalog
	req := logical.TestRequest(t, logical.DeleteOperation, "sys/plugins/catalog/mock-plugin")
	req.ClientToken = core.Client.Token()
	resp, err := core.HandleRequest(req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	// Seal the cluster
	cluster.EnsureCoresSealed(t)

	// Unseal the cluster
	barrierKeys := cluster.BarrierKeys
	for _, core := range cluster.Cores {
		for _, key := range barrierKeys {
			_, err := core.Unseal(vault.TestKeyCopy(key))
			if err != nil {
				t.Fatal(err)
			}
		}
		sealed, err := core.Sealed()
		if err != nil {
			t.Fatalf("err checking seal status: %s", err)
		}
		if sealed {
			t.Fatal("should not be sealed")
		}
		// Wait for active so post-unseal takes place
		// If it fails, it means unseal process failed
		vault.TestWaitActive(t, core.Core)
	}

	if testMount {
		// Add plugin back to the catalog
		vault.TestAddTestPlugin(t, core.Core, "mock-plugin", "TestBackend_PluginMainLogical")

		// Mount the plugin at the same path after plugin is re-added to the catalog
		// and expect an error due to existing path.
		var err error
		switch btype {
		case logical.TypeLogical:
			_, err = core.Client.Logical().Write("sys/mounts/mock-0", map[string]interface{}{
				"type": "plugin",
				"config": map[string]interface{}{
					"plugin_name": "mock-plugin",
				},
			})
		case logical.TypeCredential:
			_, err = core.Client.Logical().Write("sys/auth/mock-0", map[string]interface{}{
				"type":        "plugin",
				"plugin_name": "mock-plugin",
			})
		}
		if err == nil {
			t.Fatal("expected error when mounting on existing path")
		}
	}
}

func TestSystemBackend_Plugin_autoReload(t *testing.T) {
	cluster := testSystemBackendMock(t, 1, 1, logical.TypeLogical)
	defer cluster.Cleanup()

	core := cluster.Cores[0]

	// Update internal value
	req := logical.TestRequest(t, logical.UpdateOperation, "mock-0/internal")
	req.ClientToken = core.Client.Token()
	req.Data["value"] = "baz"
	resp, err := core.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %v", resp)
	}

	// Call errors/rpc endpoint to trigger reload
	req = logical.TestRequest(t, logical.ReadOperation, "mock-0/errors/rpc")
	req.ClientToken = core.Client.Token()
	resp, err = core.HandleRequest(req)
	if err == nil {
		t.Fatalf("expected error from error/rpc request")
	}

	// Check internal value to make sure it's reset
	req = logical.TestRequest(t, logical.ReadOperation, "mock-0/internal")
	req.ClientToken = core.Client.Token()
	resp, err = core.HandleRequest(req)
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

func TestSystemBackend_Plugin_SealUnseal(t *testing.T) {
	cluster := testSystemBackendMock(t, 1, 1, logical.TypeLogical)
	defer cluster.Cleanup()

	// Seal the cluster
	cluster.EnsureCoresSealed(t)

	// Unseal the cluster
	barrierKeys := cluster.BarrierKeys
	for _, core := range cluster.Cores {
		for _, key := range barrierKeys {
			_, err := core.Unseal(vault.TestKeyCopy(key))
			if err != nil {
				t.Fatal(err)
			}
		}
		sealed, err := core.Sealed()
		if err != nil {
			t.Fatalf("err checking seal status: %s", err)
		}
		if sealed {
			t.Fatal("should not be sealed")
		}
		// Wait for active so post-unseal takes place
		// If it fails, it means unseal process failed
		vault.TestWaitActive(t, core.Core)
	}
}

func TestSystemBackend_Plugin_reload(t *testing.T) {
	data := map[string]interface{}{
		"plugin": "mock-plugin",
	}
	t.Run("plugin", func(t *testing.T) { testSystemBackend_PluginReload(t, data) })

	data = map[string]interface{}{
		"mounts": "mock-0/,mock-1/",
	}
	t.Run("mounts", func(t *testing.T) { testSystemBackend_PluginReload(t, data) })
}

// Helper func to test different reload methods on plugin reload endpoint
func testSystemBackend_PluginReload(t *testing.T, reqData map[string]interface{}) {
	cluster := testSystemBackendMock(t, 1, 2, logical.TypeLogical)
	defer cluster.Cleanup()

	core := cluster.Cores[0]
	client := core.Client

	for i := 0; i < 2; i++ {
		// Update internal value in the backend
		resp, err := client.Logical().Write(fmt.Sprintf("mock-%d/internal", i), map[string]interface{}{
			"value": "baz",
		})
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if resp != nil {
			t.Fatalf("bad: %v", resp)
		}
	}

	// Perform plugin reload
	resp, err := client.Logical().Write("sys/plugins/reload/backend", reqData)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %v", resp)
	}

	for i := 0; i < 2; i++ {
		// Ensure internal backed value is reset
		resp, err := client.Logical().Read(fmt.Sprintf("mock-%d/internal", i))
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
func testSystemBackendMock(t *testing.T, numCores, numMounts int, backendType logical.BackendType) *vault.TestCluster {
	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"plugin": plugin.Factory,
		},
		CredentialBackends: map[string]logical.Factory{
			"plugin": plugin.Factory,
		},
	}

	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc:        vaulthttp.Handler,
		KeepStandbysSealed: true,
		NumCores:           numCores,
	})
	cluster.Start()

	core := cluster.Cores[0]
	vault.TestWaitActive(t, core.Core)
	client := core.Client

	os.Setenv(pluginutil.PluginCACertPEMEnv, cluster.CACertPEMFile)

	switch backendType {
	case logical.TypeLogical:
		vault.TestAddTestPlugin(t, core.Core, "mock-plugin", "TestBackend_PluginMainLogical")
		for i := 0; i < numMounts; i++ {
			// Alternate input styles for plugin_name on every other mount
			options := map[string]interface{}{
				"type": "plugin",
			}
			if (i+1)%2 == 0 {
				options["config"] = map[string]interface{}{
					"plugin_name": "mock-plugin",
				}
			} else {
				options["plugin_name"] = "mock-plugin"
			}
			resp, err := client.Logical().Write(fmt.Sprintf("sys/mounts/mock-%d", i), options)
			if err != nil {
				t.Fatalf("err: %v", err)
			}
			if resp != nil {
				t.Fatalf("bad: %v", resp)
			}
		}
	case logical.TypeCredential:
		vault.TestAddTestPlugin(t, core.Core, "mock-plugin", "TestBackend_PluginMainCredentials")
		for i := 0; i < numMounts; i++ {
			// Alternate input styles for plugin_name on every other mount
			options := map[string]interface{}{
				"type": "plugin",
			}
			if (i+1)%2 == 0 {
				options["config"] = map[string]interface{}{
					"plugin_name": "mock-plugin",
				}
			} else {
				options["plugin_name"] = "mock-plugin"
			}
			resp, err := client.Logical().Write(fmt.Sprintf("sys/auth/mock-%d", i), options)
			if err != nil {
				t.Fatalf("err: %v", err)
			}
			if resp != nil {
				t.Fatalf("bad: %v", resp)
			}
		}
	default:
		t.Fatal("unknown backend type provided")
	}

	return cluster
}

func TestBackend_PluginMainLogical(t *testing.T) {
	args := []string{}
	if os.Getenv(pluginutil.PluginUnwrapTokenEnv) == "" && os.Getenv(pluginutil.PluginMetadaModeEnv) != "true" {
		return
	}

	caPEM := os.Getenv(pluginutil.PluginCACertPEMEnv)
	if caPEM == "" {
		t.Fatal("CA cert not passed in")
	}
	args = append(args, fmt.Sprintf("--ca-cert=%s", caPEM))

	apiClientMeta := &pluginutil.APIClientMeta{}
	flags := apiClientMeta.FlagSet()
	flags.Parse(args)
	tlsConfig := apiClientMeta.GetTLSConfig()
	tlsProviderFunc := pluginutil.VaultPluginTLSProvider(tlsConfig)

	factoryFunc := mock.FactoryType(logical.TypeLogical)

	err := lplugin.Serve(&lplugin.ServeOpts{
		BackendFactoryFunc: factoryFunc,
		TLSProviderFunc:    tlsProviderFunc,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestBackend_PluginMainCredentials(t *testing.T) {
	args := []string{}
	if os.Getenv(pluginutil.PluginUnwrapTokenEnv) == "" && os.Getenv(pluginutil.PluginMetadaModeEnv) != "true" {
		return
	}

	caPEM := os.Getenv(pluginutil.PluginCACertPEMEnv)
	if caPEM == "" {
		t.Fatal("CA cert not passed in")
	}
	args = append(args, fmt.Sprintf("--ca-cert=%s", caPEM))

	apiClientMeta := &pluginutil.APIClientMeta{}
	flags := apiClientMeta.FlagSet()
	flags.Parse(args)
	tlsConfig := apiClientMeta.GetTLSConfig()
	tlsProviderFunc := pluginutil.VaultPluginTLSProvider(tlsConfig)

	factoryFunc := mock.FactoryType(logical.TypeCredential)

	err := lplugin.Serve(&lplugin.ServeOpts{
		BackendFactoryFunc: factoryFunc,
		TLSProviderFunc:    tlsProviderFunc,
	})
	if err != nil {
		t.Fatal(err)
	}
}
