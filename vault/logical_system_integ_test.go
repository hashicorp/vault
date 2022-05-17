package vault_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/plugin"
	"github.com/hashicorp/vault/helper/namespace"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/hashicorp/vault/sdk/physical/inmem"
	lplugin "github.com/hashicorp/vault/sdk/plugin"
	"github.com/hashicorp/vault/sdk/plugin/mock"
	"github.com/hashicorp/vault/sdk/version"
	"github.com/hashicorp/vault/vault"
)

const (
	expectedEnvKey   = "FOO"
	expectedEnvValue = "BAR"
)

func TestSystemBackend_Plugin_secret(t *testing.T) {
	cluster := testSystemBackendMock(t, 1, 1, logical.TypeLogical)
	defer cluster.Cleanup()

	core := cluster.Cores[0]

	// Make a request to lazy load the plugin
	req := logical.TestRequest(t, logical.ReadOperation, "mock-0/internal")
	req.ClientToken = core.Client.Token()
	resp, err := core.HandleRequest(namespace.RootContext(nil), req)
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
		if core.Sealed() {
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
	resp, err := core.HandleRequest(namespace.RootContext(nil), req)
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
		if core.Sealed() {
			t.Fatal("should not be sealed")
		}
		// Wait for active so post-unseal takes place
		// If it fails, it means unseal process failed
		vault.TestWaitActive(t, core.Core)
	}
}

func TestSystemBackend_Plugin_MissingBinary(t *testing.T) {
	cluster := testSystemBackendMock(t, 1, 1, logical.TypeLogical)
	defer cluster.Cleanup()

	core := cluster.Cores[0]

	// Make a request to lazy load the plugin
	req := logical.TestRequest(t, logical.ReadOperation, "mock-0/internal")
	req.ClientToken = core.Client.Token()
	resp, err := core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil {
		t.Fatalf("bad: response should not be nil")
	}

	// Seal the cluster
	cluster.EnsureCoresSealed(t)

	// Simulate removal of the plugin binary. Use os.Args to determine file name
	// since that's how we create the file for catalog registration in the test
	// helper.
	pluginFileName := filepath.Base(os.Args[0])
	err = os.Remove(filepath.Join(cluster.TempDir, pluginFileName))
	if err != nil {
		t.Fatal(err)
	}

	// Unseal the cluster
	cluster.UnsealCores(t)

	// Make a request against on tune after it is removed
	req = logical.TestRequest(t, logical.ReadOperation, "sys/mounts/mock-0/tune")
	req.ClientToken = core.Client.Token()
	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestSystemBackend_Plugin_MismatchType(t *testing.T) {
	cluster := testSystemBackendMock(t, 1, 1, logical.TypeLogical)
	defer cluster.Cleanup()

	core := cluster.Cores[0]

	// Add a credential backend with the same name
	vault.TestAddTestPlugin(t, core.Core, "mock-plugin", consts.PluginTypeCredential, "TestBackend_PluginMainCredentials", []string{}, "")

	// Make a request to lazy load the now-credential plugin
	// and expect an error
	req := logical.TestRequest(t, logical.ReadOperation, "mock-0/internal")
	req.ClientToken = core.Client.Token()
	_, err := core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("adding a same-named plugin of a different type should be no problem: %s", err)
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
	req := logical.TestRequest(t, logical.DeleteOperation, "sys/plugins/catalog/database/mock-plugin")
	req.ClientToken = core.Client.Token()
	resp, err := core.HandleRequest(namespace.RootContext(nil), req)
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
		if core.Sealed() {
			t.Fatal("should not be sealed")
		}
	}

	// Wait for active so post-unseal takes place
	// If it fails, it means unseal process failed
	vault.TestWaitActive(t, core.Core)

	if testMount {
		// Mount the plugin at the same path after plugin is re-added to the catalog
		// and expect an error due to existing path.
		var err error
		switch btype {
		case logical.TypeLogical:
			// Add plugin back to the catalog
			vault.TestAddTestPlugin(t, core.Core, "mock-plugin", consts.PluginTypeSecrets, "TestBackend_PluginMainLogical", []string{}, "")
			_, err = core.Client.Logical().Write("sys/mounts/mock-0", map[string]interface{}{
				"type": "test",
			})
		case logical.TypeCredential:
			// Add plugin back to the catalog
			vault.TestAddTestPlugin(t, core.Core, "mock-plugin", consts.PluginTypeCredential, "TestBackend_PluginMainCredentials", []string{}, "")
			_, err = core.Client.Logical().Write("sys/auth/mock-0", map[string]interface{}{
				"type": "test",
			})
		}
		if err == nil {
			t.Fatal("expected error when mounting on existing path")
		}
	}
}

func TestSystemBackend_Plugin_continueOnError(t *testing.T) {
	t.Run("secret", func(t *testing.T) {
		t.Run("sha256_mismatch", func(t *testing.T) {
			testPlugin_continueOnError(t, logical.TypeLogical, true, "mock-plugin", consts.PluginTypeSecrets)
		})

		t.Run("missing_plugin", func(t *testing.T) {
			testPlugin_continueOnError(t, logical.TypeLogical, false, "mock-plugin", consts.PluginTypeSecrets)
		})
	})

	t.Run("auth", func(t *testing.T) {
		t.Run("sha256_mismatch", func(t *testing.T) {
			testPlugin_continueOnError(t, logical.TypeCredential, true, "mock-plugin", consts.PluginTypeCredential)
		})

		t.Run("missing_plugin", func(t *testing.T) {
			testPlugin_continueOnError(t, logical.TypeCredential, false, "mock-plugin", consts.PluginTypeCredential)
		})

		t.Run("sha256_mismatch", func(t *testing.T) {
			testPlugin_continueOnError(t, logical.TypeCredential, true, "oidc", consts.PluginTypeCredential)
		})

		t.Run("missing_plugin", func(t *testing.T) {
			testPlugin_continueOnError(t, logical.TypeCredential, false, "oidc", consts.PluginTypeCredential)
		})
	})
}

func testPlugin_continueOnError(t *testing.T, btype logical.BackendType, mismatch bool, mountPoint string, pluginType consts.PluginType) {
	cluster := testSystemBackendMock(t, 1, 1, btype)
	defer cluster.Cleanup()

	core := cluster.Cores[0]

	// Get the registered plugin
	req := logical.TestRequest(t, logical.ReadOperation, fmt.Sprintf("sys/plugins/catalog/%s/mock-plugin", pluginType))
	req.ClientToken = core.Client.Token()
	resp, err := core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || resp == nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	command, ok := resp.Data["command"].(string)
	if !ok || command == "" {
		t.Fatal("invalid command")
	}

	// Mount credential type plugins
	switch btype {
	case logical.TypeCredential:
		vault.TestAddTestPlugin(t, core.Core, mountPoint, consts.PluginTypeCredential, "TestBackend_PluginMainCredentials", []string{}, cluster.TempDir)
		_, err = core.Client.Logical().Write(fmt.Sprintf("sys/auth/%s", mountPoint), map[string]interface{}{
			"type": "mock-plugin",
		})
		if err != nil {
			t.Fatalf("err:%v", err)
		}
	}

	// Trigger a sha256 mismatch or missing plugin error
	if mismatch {
		req = logical.TestRequest(t, logical.UpdateOperation, fmt.Sprintf("sys/plugins/catalog/%s/mock-plugin", pluginType))
		req.Data = map[string]interface{}{
			"sha256":  "d17bd7334758e53e6fbab15745d2520765c06e296f2ce8e25b7919effa0ac216",
			"command": filepath.Base(command),
		}
		req.ClientToken = core.Client.Token()
		resp, err = core.HandleRequest(namespace.RootContext(nil), req)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("err:%v resp:%#v", err, resp)
		}
	} else {
		err := os.Remove(filepath.Join(cluster.TempDir, filepath.Base(command)))
		if err != nil {
			t.Fatal(err)
		}
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
		if core.Sealed() {
			t.Fatal("should not be sealed")
		}
	}

	// Wait for active so post-unseal takes place
	// If it fails, it means unseal process failed
	vault.TestWaitActive(t, core.Core)

	// Re-add the plugin to the catalog
	switch btype {
	case logical.TypeLogical:
		vault.TestAddTestPlugin(t, core.Core, "mock-plugin", consts.PluginTypeSecrets, "TestBackend_PluginMainLogical", []string{}, cluster.TempDir)
	case logical.TypeCredential:
		vault.TestAddTestPlugin(t, core.Core, "mock-plugin", consts.PluginTypeCredential, "TestBackend_PluginMainCredentials", []string{}, cluster.TempDir)
	}

	// Reload the plugin
	req = logical.TestRequest(t, logical.UpdateOperation, "sys/plugins/reload/backend")
	req.Data = map[string]interface{}{
		"plugin": "mock-plugin",
	}
	req.ClientToken = core.Client.Token()
	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	// Make a request to lazy load the plugin
	var reqPath string
	switch btype {
	case logical.TypeLogical:
		reqPath = "mock-0/internal"
	case logical.TypeCredential:
		reqPath = "auth/mock-0/internal"
	}

	req = logical.TestRequest(t, logical.ReadOperation, reqPath)
	req.ClientToken = core.Client.Token()
	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil {
		t.Fatalf("bad: response should not be nil")
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
	resp, err := core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %v", resp)
	}

	// Call errors/rpc endpoint to trigger reload
	req = logical.TestRequest(t, logical.ReadOperation, "mock-0/errors/rpc")
	req.ClientToken = core.Client.Token()
	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err == nil {
		t.Fatalf("expected error from error/rpc request")
	}

	// Check internal value to make sure it's reset
	req = logical.TestRequest(t, logical.ReadOperation, "mock-0/internal")
	req.ClientToken = core.Client.Token()
	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
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
		if core.Sealed() {
			t.Fatal("should not be sealed")
		}
	}

	// Wait for active so post-unseal takes place
	// If it fails, it means unseal process failed
	vault.TestWaitActive(t, cluster.Cores[0].Core)
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
	if resp == nil {
		t.Fatalf("bad: %v", resp)
	}
	if resp.Data["reload_id"] == nil {
		t.Fatal("no reload_id in response")
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
// of mounted mock plugin backends. numMounts alternates between different
// ways of providing the plugin_name.
//
// The mounts are mounted at sys/mounts/mock-[numMounts] or sys/auth/mock-[numMounts]
func testSystemBackendMock(t *testing.T, numCores, numMounts int, backendType logical.BackendType) *vault.TestCluster {
	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"plugin": plugin.Factory,
		},
		CredentialBackends: map[string]logical.Factory{
			"plugin": plugin.Factory,
		},
	}

	// Create a tempdir, cluster.Cleanup will clean up this directory
	tempDir, err := ioutil.TempDir("", "vault-test-cluster")
	if err != nil {
		t.Fatal(err)
	}

	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc:        vaulthttp.Handler,
		KeepStandbysSealed: true,
		NumCores:           numCores,
		TempDir:            tempDir,
	})
	cluster.Start()

	core := cluster.Cores[0]
	vault.TestWaitActive(t, core.Core)
	client := core.Client

	os.Setenv(pluginutil.PluginCACertPEMEnv, cluster.CACertPEMFile)

	switch backendType {
	case logical.TypeLogical:
		vault.TestAddTestPlugin(t, core.Core, "mock-plugin", consts.PluginTypeSecrets, "TestBackend_PluginMainLogical", []string{}, tempDir)
		for i := 0; i < numMounts; i++ {
			// Alternate input styles for plugin_name on every other mount
			options := map[string]interface{}{
				"type": "mock-plugin",
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
		vault.TestAddTestPlugin(t, core.Core, "mock-plugin", consts.PluginTypeCredential, "TestBackend_PluginMainCredentials", []string{}, tempDir)
		for i := 0; i < numMounts; i++ {
			// Alternate input styles for plugin_name on every other mount
			options := map[string]interface{}{
				"type": "mock-plugin",
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

func TestSystemBackend_Plugin_Env(t *testing.T) {
	kvPair := fmt.Sprintf("%s=%s", expectedEnvKey, expectedEnvValue)
	cluster := testSystemBackend_SingleCluster_Env(t, []string{kvPair})
	defer cluster.Cleanup()
}

// testSystemBackend_SingleCluster_Env is a helper func that returns a single
// cluster and a single mounted plugin logical backend.
func testSystemBackend_SingleCluster_Env(t *testing.T, env []string) *vault.TestCluster {
	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"test": plugin.Factory,
		},
	}
	// Create a tempdir, cluster.Cleanup will clean up this directory
	tempDir, err := ioutil.TempDir("", "vault-test-cluster")
	if err != nil {
		t.Fatal(err)
	}

	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc:        vaulthttp.Handler,
		KeepStandbysSealed: true,
		NumCores:           1,
		TempDir:            tempDir,
	})
	cluster.Start()

	core := cluster.Cores[0]
	vault.TestWaitActive(t, core.Core)
	client := core.Client

	os.Setenv(pluginutil.PluginCACertPEMEnv, cluster.CACertPEMFile)

	vault.TestAddTestPlugin(t, core.Core, "mock-plugin", consts.PluginTypeSecrets, "TestBackend_PluginMainEnv", env, tempDir)
	options := map[string]interface{}{
		"type": "mock-plugin",
	}

	resp, err := client.Logical().Write("sys/mounts/mock", options)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %v", resp)
	}

	return cluster
}

func TestBackend_PluginMainLogical(t *testing.T) {
	args := []string{}
	if os.Getenv(pluginutil.PluginUnwrapTokenEnv) == "" && os.Getenv(pluginutil.PluginMetadataModeEnv) != "true" {
		return
	}

	caPEM := os.Getenv(pluginutil.PluginCACertPEMEnv)
	if caPEM == "" {
		t.Fatal("CA cert not passed in")
	}
	args = append(args, fmt.Sprintf("--ca-cert=%s", caPEM))

	apiClientMeta := &api.PluginAPIClientMeta{}
	flags := apiClientMeta.FlagSet()
	flags.Parse(args)
	tlsConfig := apiClientMeta.GetTLSConfig()
	tlsProviderFunc := api.VaultPluginTLSProvider(tlsConfig)

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
	if os.Getenv(pluginutil.PluginUnwrapTokenEnv) == "" && os.Getenv(pluginutil.PluginMetadataModeEnv) != "true" {
		return
	}

	caPEM := os.Getenv(pluginutil.PluginCACertPEMEnv)
	if caPEM == "" {
		t.Fatal("CA cert not passed in")
	}
	args = append(args, fmt.Sprintf("--ca-cert=%s", caPEM))

	apiClientMeta := &api.PluginAPIClientMeta{}
	flags := apiClientMeta.FlagSet()
	flags.Parse(args)
	tlsConfig := apiClientMeta.GetTLSConfig()
	tlsProviderFunc := api.VaultPluginTLSProvider(tlsConfig)

	factoryFunc := mock.FactoryType(logical.TypeCredential)

	err := lplugin.Serve(&lplugin.ServeOpts{
		BackendFactoryFunc: factoryFunc,
		TLSProviderFunc:    tlsProviderFunc,
	})
	if err != nil {
		t.Fatal(err)
	}
}

// TestBackend_PluginMainEnv is a mock plugin that simply checks for the existence of FOO env var.
func TestBackend_PluginMainEnv(t *testing.T) {
	args := []string{}
	if os.Getenv(pluginutil.PluginUnwrapTokenEnv) == "" && os.Getenv(pluginutil.PluginMetadataModeEnv) != "true" {
		return
	}

	// Check on actual vs expected env var
	actual := os.Getenv(expectedEnvKey)
	if actual != expectedEnvValue {
		t.Fatalf("expected: %q, got: %q", expectedEnvValue, actual)
	}

	caPEM := os.Getenv(pluginutil.PluginCACertPEMEnv)
	if caPEM == "" {
		t.Fatal("CA cert not passed in")
	}
	args = append(args, fmt.Sprintf("--ca-cert=%s", caPEM))

	apiClientMeta := &api.PluginAPIClientMeta{}
	flags := apiClientMeta.FlagSet()
	flags.Parse(args)
	tlsConfig := apiClientMeta.GetTLSConfig()
	tlsProviderFunc := api.VaultPluginTLSProvider(tlsConfig)

	factoryFunc := mock.FactoryType(logical.TypeLogical)

	err := lplugin.Serve(&lplugin.ServeOpts{
		BackendFactoryFunc: factoryFunc,
		TLSProviderFunc:    tlsProviderFunc,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestSystemBackend_InternalUIResultantACL(t *testing.T) {
	cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()
	client := cluster.Cores[0].Client

	resp, err := client.Auth().Token().Create(&api.TokenCreateRequest{
		Policies: []string{"default"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("nil response")
	}
	if resp.Auth == nil {
		t.Fatal("nil auth")
	}
	if resp.Auth.ClientToken == "" {
		t.Fatal("empty client token")
	}

	client.SetToken(resp.Auth.ClientToken)

	resp, err = client.Logical().Read("sys/internal/ui/resultant-acl")
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("nil response")
	}
	if resp.Data == nil {
		t.Fatal("nil data")
	}

	exp := map[string]interface{}{
		"exact_paths": map[string]interface{}{
			"auth/token/lookup-self": map[string]interface{}{
				"capabilities": []interface{}{
					"read",
				},
			},
			"auth/token/renew-self": map[string]interface{}{
				"capabilities": []interface{}{
					"update",
				},
			},
			"auth/token/revoke-self": map[string]interface{}{
				"capabilities": []interface{}{
					"update",
				},
			},
			"sys/capabilities-self": map[string]interface{}{
				"capabilities": []interface{}{
					"update",
				},
			},
			"sys/control-group/request": map[string]interface{}{
				"capabilities": []interface{}{
					"update",
				},
			},
			"sys/internal/ui/resultant-acl": map[string]interface{}{
				"capabilities": []interface{}{
					"read",
				},
			},
			"sys/leases/lookup": map[string]interface{}{
				"capabilities": []interface{}{
					"update",
				},
			},
			"sys/leases/renew": map[string]interface{}{
				"capabilities": []interface{}{
					"update",
				},
			},
			"sys/renew": map[string]interface{}{
				"capabilities": []interface{}{
					"update",
				},
			},
			"sys/tools/hash": map[string]interface{}{
				"capabilities": []interface{}{
					"update",
				},
			},
			"sys/wrapping/lookup": map[string]interface{}{
				"capabilities": []interface{}{
					"update",
				},
			},
			"sys/wrapping/unwrap": map[string]interface{}{
				"capabilities": []interface{}{
					"update",
				},
			},
			"sys/wrapping/wrap": map[string]interface{}{
				"capabilities": []interface{}{
					"update",
				},
			},
		},
		"glob_paths": map[string]interface{}{
			"cubbyhole/": map[string]interface{}{
				"capabilities": []interface{}{
					"create",
					"delete",
					"list",
					"read",
					"update",
				},
			},
			"sys/tools/hash/": map[string]interface{}{
				"capabilities": []interface{}{
					"update",
				},
			},
		},
		"root": false,
	}

	if diff := deep.Equal(resp.Data, exp); diff != nil {
		t.Fatal(diff)
	}
}

func TestSystemBackend_HAStatus(t *testing.T) {
	logger := logging.NewVaultLogger(hclog.Trace)
	inm, err := inmem.NewTransactionalInmem(nil, logger)
	if err != nil {
		t.Fatal(err)
	}
	inmha, err := inmem.NewInmemHA(nil, logger)
	if err != nil {
		t.Fatal(err)
	}

	conf := &vault.CoreConfig{
		Physical:   inm,
		HAPhysical: inmha.(physical.HABackend),
	}
	opts := &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	}
	cluster := vault.NewTestCluster(t, conf, opts)
	cluster.Start()
	defer cluster.Cleanup()

	vault.RetryUntil(t, 15*time.Second, func() error {
		// Use standby deliberately to make sure it forwards
		client := cluster.Cores[1].Client
		resp, err := client.Sys().HAStatus()
		if err != nil {
			t.Fatal(err)
		}

		if len(resp.Nodes) != len(cluster.Cores) {
			return fmt.Errorf("expected %d nodes, got %d", len(cluster.Cores), len(resp.Nodes))
		}
		return nil
	})
}

// TestSystemBackend_VersionHistory_unauthenticated tests the sys/version-history
// endpoint without providing a token. Requests to the endpoint must be
// authenticated and thus a 403 response is expected.
func TestSystemBackend_VersionHistory_unauthenticated(t *testing.T) {
	cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()
	client := cluster.Cores[0].Client

	client.SetToken("")
	resp, err := client.Logical().List("sys/version-history")

	if resp != nil {
		t.Fatalf("expected nil response, resp: %#v", resp)
	}

	respErr, ok := err.(*api.ResponseError)
	if !ok {
		t.Fatalf("unexpected error type: err: %#v", err)
	}

	if respErr.StatusCode != 403 {
		t.Fatalf("expected response status to be 403, actual: %d", respErr.StatusCode)
	}
}

// TestSystemBackend_VersionHistory_authenticated tests the sys/version-history
// endpoint with authentication. Without synthetically altering the underlying
// core/versions storage entries, a single version entry should exist.
func TestSystemBackend_VersionHistory_authenticated(t *testing.T) {
	cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()
	client := cluster.Cores[0].Client

	resp, err := client.Logical().List("sys/version-history")
	if err != nil || resp == nil {
		t.Fatalf("request failed, err: %v, resp: %#v", err, resp)
	}

	var ok bool
	var keys []interface{}
	var keyInfo map[string]interface{}

	if keys, ok = resp.Data["keys"].([]interface{}); !ok {
		t.Fatalf("expected keys to be array, actual: %#v", resp.Data["keys"])
	}

	if keyInfo, ok = resp.Data["key_info"].(map[string]interface{}); !ok {
		t.Fatalf("expected key_info to be map, actual: %#v", resp.Data["key_info"])
	}

	if len(keys) != 1 {
		t.Fatalf("expected single version history entry for %q", version.Version)
	}

	if keyInfo[version.Version] == nil {
		t.Fatalf("expected version %s to be present in key_info, actual: %#v", version.Version, keyInfo)
	}
}
