// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package plugin_test

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/plugin"
	"github.com/hashicorp/vault/helper/constants"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	"github.com/hashicorp/vault/helper/testhelpers/teststorage"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/logical"
	lplugin "github.com/hashicorp/vault/sdk/plugin"
	"github.com/hashicorp/vault/sdk/plugin/mock"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/require"
)

// logicalVersionMap is a map of version to test plugin
var logicalVersionMap = map[string]string{
	"v4":             "TestBackend_PluginMain_V4_Logical",
	"v5":             "TestBackend_PluginMainLogical",
	"v5_multiplexed": "TestBackend_PluginMain_Multiplexed_Logical",
}

// credentialVersionMap is a map of version to test plugin
var credentialVersionMap = map[string]string{
	"v4":             "TestBackend_PluginMain_V4_Credentials",
	"v5":             "TestBackend_PluginMainCredentials",
	"v5_multiplexed": "TestBackend_PluginMain_Multiplexed_Credentials",
}

var testCtx = context.TODO()

func TestSystemBackend_Plugin_secret(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		pluginVersion string
	}{
		{
			pluginVersion: "v5_multiplexed",
		},
		{
			pluginVersion: "v5",
		},
		{
			pluginVersion: "v4",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.pluginVersion, func(t *testing.T) {
			t.Parallel()
			cluster := testSystemBackendMock(t, 1, 1, logical.TypeLogical, tc.pluginVersion)
			defer cluster.Cleanup()

			core := cluster.Cores[0]

			// Make a request to lazy load the plugin
			req := logical.TestRequest(t, logical.ReadOperation, "mock-0/internal")
			req.ClientToken = core.Client.Token()
			resp, err := core.HandleRequest(namespace.RootContext(testCtx), req)
			if err != nil {
				t.Fatalf("err: %v", err)
			}
			if resp == nil {
				t.Fatalf("bad: response should not be nil")
			}

			req = logical.TestRequest(t, logical.CreateOperation, "mock-0/config")
			req.ClientToken = core.Client.Token()
			resp, err = core.HandleRequest(namespace.RootContext(testCtx), req)
			if err != nil {
				t.Fatalf("err: %v", err)
			}
			wspost := req.ResponseState()
			if constants.IsEnterprise {
				require.NotZero(t, wspost.ReplicatedIndex)
				require.NotZero(t, wspost.LocalIndex)
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
		})
	}
}

func TestSystemBackend_Plugin_auth(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		pluginVersion string
	}{
		{
			pluginVersion: "v5_multiplexed",
		},
		{
			pluginVersion: "v5",
		},
		{
			pluginVersion: "v4",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.pluginVersion, func(t *testing.T) {
			t.Parallel()
			cluster := testSystemBackendMock(t, 1, 1, logical.TypeCredential, tc.pluginVersion)
			defer cluster.Cleanup()

			core := cluster.Cores[0]

			// Make a request to lazy load the plugin
			req := logical.TestRequest(t, logical.ReadOperation, "auth/mock-0/internal")
			req.ClientToken = core.Client.Token()
			resp, err := core.HandleRequest(namespace.RootContext(testCtx), req)
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
		})
	}
}

func TestSystemBackend_Plugin_MissingBinary(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		pluginVersion string
	}{
		{
			pluginVersion: "v5_multiplexed",
		},
		{
			pluginVersion: "v5",
		},
		{
			pluginVersion: "v4",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.pluginVersion, func(t *testing.T) {
			t.Parallel()
			cluster := testSystemBackendMock(t, 1, 1, logical.TypeLogical, tc.pluginVersion)
			defer cluster.Cleanup()

			core := cluster.Cores[0]

			// Make a request to lazy load the plugin
			req := logical.TestRequest(t, logical.ReadOperation, "mock-0/internal")
			req.ClientToken = core.Client.Token()
			resp, err := core.HandleRequest(namespace.RootContext(testCtx), req)
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
			err = os.Remove(filepath.Join(cluster.Cores[0].CoreConfig.PluginDirectory, pluginFileName))
			if err != nil {
				t.Fatal(err)
			}

			// Unseal the cluster
			cluster.UnsealCores(t)

			// Make a request against on tune after it is removed
			req = logical.TestRequest(t, logical.ReadOperation, "sys/mounts/mock-0/tune")
			req.ClientToken = core.Client.Token()
			_, err = core.HandleRequest(namespace.RootContext(testCtx), req)
			if err == nil {
				t.Fatalf("expected error")
			}
		})
	}
}

func TestSystemBackend_Plugin_MismatchType(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		pluginVersion string
	}{
		{
			pluginVersion: "v5_multiplexed",
		},
		{
			pluginVersion: "v5",
		},
		{
			pluginVersion: "v4",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.pluginVersion, func(t *testing.T) {
			cluster := testSystemBackendMock(t, 1, 1, logical.TypeLogical, tc.pluginVersion)
			defer cluster.Cleanup()

			core := cluster.Cores[0]

			// Add a credential backend with the same name
			vault.TestAddTestPlugin(t, core.Core, "mock-plugin", consts.PluginTypeCredential, "", "TestBackend_PluginMainCredentials", []string{})

			// Make a request to lazy load the now-credential plugin
			// and expect an error
			req := logical.TestRequest(t, logical.ReadOperation, "mock-0/internal")
			req.ClientToken = core.Client.Token()
			_, err := core.HandleRequest(namespace.RootContext(testCtx), req)
			if err != nil {
				t.Fatalf("adding a same-named plugin of a different type should be no problem: %s", err)
			}
		})
	}
}

func TestSystemBackend_Plugin_CatalogRemoved(t *testing.T) {
	t.Run("secret", func(t *testing.T) {
		t.Parallel()
		testPlugin_CatalogRemoved(t, logical.TypeLogical, false, logicalVersionMap)
	})

	t.Run("auth", func(t *testing.T) {
		t.Parallel()
		testPlugin_CatalogRemoved(t, logical.TypeCredential, false, credentialVersionMap)
	})

	t.Run("secret-mount-existing", func(t *testing.T) {
		t.Parallel()
		testPlugin_CatalogRemoved(t, logical.TypeLogical, true, logicalVersionMap)
	})

	t.Run("auth-mount-existing", func(t *testing.T) {
		t.Parallel()
		testPlugin_CatalogRemoved(t, logical.TypeCredential, true, credentialVersionMap)
	})
}

func testPlugin_CatalogRemoved(t *testing.T, btype logical.BackendType, testMount bool, versionMap map[string]string) {
	testCases := []struct {
		pluginVersion string
	}{
		{
			pluginVersion: "v5_multiplexed",
		},
		{
			pluginVersion: "v5",
		},
		{
			pluginVersion: "v4",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.pluginVersion, func(t *testing.T) {
			t.Parallel()
			cluster := testSystemBackendMock(t, 1, 1, logical.TypeLogical, tc.pluginVersion)
			defer cluster.Cleanup()

			core := cluster.Cores[0]

			// Remove the plugin from the catalog
			req := logical.TestRequest(t, logical.DeleteOperation, "sys/plugins/catalog/database/mock-plugin")
			req.ClientToken = core.Client.Token()
			resp, err := core.HandleRequest(namespace.RootContext(testCtx), req)
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
					vault.TestAddTestPlugin(t, core.Core, "mock-plugin", consts.PluginTypeSecrets, "", logicalVersionMap[tc.pluginVersion], []string{})
					_, err = core.Client.Logical().Write("sys/mounts/mock-0", map[string]interface{}{
						"type": "test",
					})
				case logical.TypeCredential:
					// Add plugin back to the catalog
					vault.TestAddTestPlugin(t, core.Core, "mock-plugin", consts.PluginTypeCredential, "", credentialVersionMap[tc.pluginVersion], []string{})
					_, err = core.Client.Logical().Write("sys/auth/mock-0", map[string]interface{}{
						"type": "test",
					})
				}
				if err == nil {
					t.Fatal("expected error when mounting on existing path")
				}
			}
		})
	}
}

func TestSystemBackend_Plugin_autoReload(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		pluginVersion string
	}{
		{
			pluginVersion: "v5_multiplexed",
		},
		{
			pluginVersion: "v5",
		},
		{
			pluginVersion: "v4",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.pluginVersion, func(t *testing.T) {
			t.Parallel()
			cluster := testSystemBackendMock(t, 1, 1, logical.TypeLogical, tc.pluginVersion)
			defer cluster.Cleanup()

			core := cluster.Cores[0]

			// Update internal value
			req := logical.TestRequest(t, logical.UpdateOperation, "mock-0/internal")
			req.ClientToken = core.Client.Token()
			req.Data["value"] = "baz"
			resp, err := core.HandleRequest(namespace.RootContext(testCtx), req)
			if err != nil {
				t.Fatalf("err: %v", err)
			}
			if resp != nil {
				t.Fatalf("bad: %v", resp)
			}

			// Call errors/rpc endpoint to trigger reload
			req = logical.TestRequest(t, logical.ReadOperation, "mock-0/errors/rpc")
			req.ClientToken = core.Client.Token()
			_, err = core.HandleRequest(namespace.RootContext(testCtx), req)
			if err == nil {
				t.Fatalf("expected error from error/rpc request")
			}

			// Check internal value to make sure it's reset
			req = logical.TestRequest(t, logical.ReadOperation, "mock-0/internal")
			req.ClientToken = core.Client.Token()
			resp, err = core.HandleRequest(namespace.RootContext(testCtx), req)
			if err != nil {
				t.Fatalf("err: %v", err)
			}
			if resp == nil {
				t.Fatalf("bad: response should not be nil")
			}
			if resp.Data["value"].(string) == "baz" {
				t.Fatal("did not expect backend internal value to be 'baz'")
			}
		})
	}
}

func TestSystemBackend_Plugin_SealUnseal(t *testing.T) {
	testCases := []struct {
		pluginVersion string
	}{
		{
			pluginVersion: "v5_multiplexed",
		},
		{
			pluginVersion: "v5",
		},
		{
			pluginVersion: "v4",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.pluginVersion, func(t *testing.T) {
			t.Parallel()
			cluster := testSystemBackendMock(t, 1, 1, logical.TypeLogical, tc.pluginVersion)
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
		})
	}
}

func TestSystemBackend_Plugin_reload(t *testing.T) {
	// Paths being tested.
	const (
		reloadBackendPath = "sys/plugins/reload/backend"
		rootReloadPath    = "sys/plugins/reload/%s/%s"
	)
	testCases := []struct {
		name        string
		backendType logical.BackendType
		path        string
		data        map[string]interface{}
	}{
		{
			name:        "test plugin reload for type credential",
			backendType: logical.TypeCredential,
			path:        reloadBackendPath,
			data: map[string]interface{}{
				"plugin": "mock-plugin",
			},
		},
		{
			name:        "test mount reload for type credential",
			backendType: logical.TypeCredential,
			path:        reloadBackendPath,
			data: map[string]interface{}{
				"mounts": "sys/auth/mock-0/,auth/mock-1/",
			},
		},
		{
			name:        "test plugin reload for type secret",
			backendType: logical.TypeLogical,
			path:        reloadBackendPath,
			data: map[string]interface{}{
				"plugin": "mock-plugin",
			},
		},
		{
			name:        "test mount reload for type secret",
			backendType: logical.TypeLogical,
			path:        reloadBackendPath,
			data: map[string]interface{}{
				"mounts": "mock-0/,mock-1",
			},
		},
		{
			name:        "root plugin reload for type auth",
			backendType: logical.TypeCredential,
			path:        fmt.Sprintf(rootReloadPath, "auth", "mock-plugin"),
		},
		{
			name:        "root plugin reload for type secret",
			backendType: logical.TypeLogical,
			path:        fmt.Sprintf(rootReloadPath, "secret", "mock-plugin"),
		},
		{
			name:        "root plugin reload for unknown type",
			backendType: logical.TypeUnknown,
			path:        fmt.Sprintf(rootReloadPath, "unknown", "mock-plugin"),
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			testSystemBackend_PluginReload(t, tc.path, tc.data, tc.backendType)
		})
	}
}

// Helper func to test different reload methods on plugin reload endpoint
func testSystemBackend_PluginReload(t *testing.T, path string, reqData map[string]interface{}, backendType logical.BackendType) {
	testCases := []struct {
		pluginVersion string
	}{
		{
			pluginVersion: "v5_multiplexed",
		},
		{
			pluginVersion: "v5",
		},
		{
			pluginVersion: "v4",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.pluginVersion, func(t *testing.T) {
			cluster := testSystemBackendMock(t, 1, 2, backendType, tc.pluginVersion)
			defer cluster.Cleanup()

			core := cluster.Cores[0]
			client := core.Client

			var mountPaths []string
			if backendType == logical.TypeCredential {
				mountPaths = []string{"auth/mock-0", "auth/mock-1"}
			} else {
				mountPaths = []string{"mock-0", "mock-1"}
			}

			for _, mountPath := range mountPaths {
				// Update internal value in the backend
				mock.WriteInternalValue(t, client, mountPath, "baz")
			}

			// Verify our precondition that the write succeeded.
			for _, mountPath := range mountPaths {
				mock.ExpectInternalValue(t, client, mountPath, "baz")
			}

			// Perform plugin reload which should reset the value.
			resp, err := client.Logical().Write(path, reqData)
			if err != nil {
				t.Fatalf("err: %v", err)
			}
			if resp == nil {
				t.Fatalf("bad: %v", resp)
			}
			if resp.Data["reload_id"] == nil {
				t.Fatal("no reload_id in response")
			}
			if len(resp.Warnings) != 0 {
				t.Fatal(resp.Warnings)
			}

			// Ensure internal backed value is reset
			for _, mountPath := range mountPaths {
				mock.ExpectInternalValue(t, client, mountPath, mock.MockPluginDefaultInternalValue)
			}
		})
	}
}

func TestSystemBackend_PluginReload_WarningIfNoneReloaded(t *testing.T) {
	cluster := testSystemBackendMock(t, 1, 2, logical.TypeLogical, "v5")
	defer cluster.Cleanup()

	core := cluster.Cores[0]
	client := core.Client

	for _, backendType := range []logical.BackendType{logical.TypeLogical, logical.TypeCredential} {
		t.Run(backendType.String(), func(t *testing.T) {
			// Perform plugin reload
			resp, err := client.Logical().Write("sys/plugins/reload/backend", map[string]any{
				"plugin": "does-not-exist",
			})
			if err != nil {
				t.Fatalf("err: %v", err)
			}
			if resp == nil {
				t.Fatalf("bad: %v", resp)
			}
			if resp.Data["reload_id"] == nil {
				t.Fatal("no reload_id in response")
			}
			if len(resp.Warnings) == 0 {
				t.Fatal("expected warning")
			}
		})
	}
}

// testSystemBackendMock returns a systemBackend with the desired number
// of mounted mock plugin backends. numMounts alternates between different
// ways of providing the plugin_name.
//
// The mounts are mounted at sys/mounts/mock-[numMounts] or sys/auth/mock-[numMounts]
func testSystemBackendMock(t *testing.T, numCores, numMounts int, backendType logical.BackendType, pluginVersion string) *vault.TestCluster {
	t.Helper()
	pluginDir := corehelpers.MakeTestPluginDir(t)
	conf, opts := teststorage.ClusterSetup(&vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"plugin": plugin.Factory,
		},
		CredentialBackends: map[string]logical.Factory{
			"plugin": plugin.Factory,
		},
		PluginDirectory: pluginDir,
	}, &vault.TestClusterOptions{
		KeepStandbysSealed: true,
		NumCores:           numCores,
		TempDir:            pluginDir,
	}, teststorage.InmemBackendSetup)

	cluster := vault.NewTestCluster(t, conf, opts)

	core := cluster.Cores[0]
	vault.TestWaitActive(t, core.Core)
	client := core.Client

	env := []string{pluginutil.PluginCACertPEMEnv + "=" + cluster.CACertPEMFile}

	switch backendType {
	case logical.TypeLogical, logical.TypeUnknown:
		plugin := logicalVersionMap[pluginVersion]
		vault.TestAddTestPlugin(t, core.Core, "mock-plugin", consts.PluginTypeSecrets, "", plugin, env)
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
		plugin := credentialVersionMap[pluginVersion]
		vault.TestAddTestPlugin(t, core.Core, "mock-plugin", consts.PluginTypeCredential, "", plugin, env)
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

// TestSystemBackend_Plugin_Env ensures we use env vars specified during plugin
// registration, and get the priority between OS and plugin env vars correct.
func TestSystemBackend_Plugin_Env(t *testing.T) {
	pluginDir := corehelpers.MakeTestPluginDir(t)
	coreConfig := &vault.CoreConfig{
		PluginDirectory: pluginDir,
	}

	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
		NumCores:    1,
		TempDir:     pluginDir,
	})
	cluster.Start()
	t.Cleanup(cluster.Cleanup)

	core := cluster.Cores[0]
	vault.TestWaitActive(t, core.Core)
	client := core.Client

	key := t.Name() + "_FOO"
	osValue := "bar"
	pluginValue := "baz"
	t.Setenv(key, osValue)
	env := []string{
		fmt.Sprintf("%s=%s", key, pluginValue),
		pluginutil.PluginCACertPEMEnv + "=" + cluster.CACertPEMFile,
	}
	vault.TestAddTestPlugin(t, core.Core, "mock-plugin", consts.PluginTypeSecrets, "", "TestBackend_PluginMainLogical", env)

	err := client.Sys().Mount("mock", &api.MountInput{
		Type: "mock-plugin",
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Plugin env should take precedence by default.
	resp, err := client.Logical().Read("mock/env/" + key)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.Data["key"] != pluginValue {
		t.Fatal(resp)
	}

	// Now set the flag that reverts to legacy behavior and reload the plugin.
	t.Setenv(pluginutil.PluginUseLegacyEnvLayering, "true")
	_, err = client.Sys().RootReloadPlugin(context.Background(), &api.RootReloadPluginInput{
		Plugin: "mock-plugin",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Now the OS value should take precedence.
	resp, err = client.Logical().Read("mock/env/" + key)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.Data["key"] != osValue {
		t.Fatal(resp)
	}
}

// TestSystemBackend_PluginCatalog_ReadVersionedOnly verifies the behavior when
// a user runs 'vault plugin info <type> <name>' without specifying -version:
//   - 1 version registered  → auto-select it, return plugin info with a warning
//   - 2+ versions registered → error listing all available versions
//   - 0 versions registered  → silent 404 (nil, nil), same as any unregistered plugin
func TestSystemBackend_PluginCatalog_ReadVersionedOnly(t *testing.T) {
	t.Parallel()

	// MakeTestPluginDir creates a real temp directory on disk.
	// Vault requires plugin binaries to live inside PluginDirectory.
	pluginDir := corehelpers.MakeTestPluginDir(t)

	// NewTestCluster starts a full in-memory Vault with a real HTTP listener —
	// same as "vault server -dev" but automated. NumCores:1 = single node.
	cluster := vault.NewTestCluster(t, &vault.CoreConfig{
		PluginDirectory: pluginDir,
	}, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
		NumCores:    1,
		TempDir:     pluginDir,
	})
	defer cluster.Cleanup()

	// core.Client is a pre-configured api.Client pointed at the cluster —
	// the same type the CLI uses when you run vault commands.
	client := cluster.Cores[0].Client

	// Create a real empty file in pluginDir to stand in as the plugin binary.
	// Vault checks that the command path exists within PluginDirectory.
	pluginFile, err := os.CreateTemp(pluginDir, "test-versioned-plugin")
	if err != nil {
		t.Fatalf("failed to create test plugin file: %v", err)
	}
	pluginFile.Close()
	pluginName := "test-versioned-plugin"
	pluginCommand := filepath.Base(pluginFile.Name())
	// hex.EncodeToString just needs to produce a valid hex string for registration.
	fakeSHA := hex.EncodeToString([]byte("fake-sha"))

	// --- Case: 1 version registered — should auto-select and warn ---

	// Register one versioned entry (no unversioned entry).
	// Equivalent to: vault plugin register -version=v1.2.3 secret <name>
	err = client.Sys().RegisterPlugin(&api.RegisterPluginInput{
		Name:    pluginName,
		Type:    api.PluginTypeSecrets,
		Version: "v1.2.3",
		Command: pluginCommand,
		SHA256:  fakeSHA,
	})
	if err != nil {
		t.Fatalf("unexpected error registering versioned plugin: %v", err)
	}

	// Read without -version — should succeed and include a warning.
	// We use client.Logical().Read() instead of client.Sys().GetPlugin() here
	// because the typed GetPluginResponse does not expose Warnings; the raw
	// Secret response does.
	// Equivalent to: vault plugin info secret <name>
	secret, err := client.Logical().Read("sys/plugins/catalog/secret/" + pluginName)
	if err != nil {
		t.Fatalf("expected auto-select to succeed, got error: %v", err)
	}
	if secret == nil {
		t.Fatal("expected non-nil response for auto-selected plugin")
	}
	if secret.Data["version"] != "v1.2.3" {
		t.Errorf("expected auto-selected version v1.2.3, got %q", secret.Data["version"])
	}
	// Assert the warning is present so the user knows a version was auto-selected.
	if len(secret.Warnings) == 0 {
		t.Error("expected a warning about auto-selected version, got none")
	} else if !strings.Contains(secret.Warnings[0], "v1.2.3") {
		t.Errorf("expected warning to mention v1.2.3, got: %q", secret.Warnings[0])
	}

	// --- Case: 2+ versions registered — should error listing both ---

	// Register a second version.
	err = client.Sys().RegisterPlugin(&api.RegisterPluginInput{
		Name:    pluginName,
		Type:    api.PluginTypeSecrets,
		Version: "v1.3.0",
		Command: pluginCommand,
		SHA256:  fakeSHA,
	})
	if err != nil {
		t.Fatalf("unexpected error registering second version: %v", err)
	}

	// Read without -version — should error listing both versions.
	_, err = client.Sys().GetPlugin(&api.GetPluginInput{
		Name: pluginName,
		Type: api.PluginTypeSecrets,
	})
	if err == nil {
		t.Fatal("expected error with multiple versions registered, got nil")
	}
	errMsg := err.Error()
	if !strings.Contains(errMsg, "v1.2.3") || !strings.Contains(errMsg, "v1.3.0") {
		t.Errorf("expected error to list both versions, got: %q", errMsg)
	}
	if !strings.Contains(errMsg, "-version") {
		t.Errorf("expected error to mention -version flag, got: %q", errMsg)
	}

	// --- Case: 0 versions registered — plugin does not exist at all ---

	// Deregister both versioned entries so nothing is left.
	err = client.Sys().DeregisterPlugin(&api.DeregisterPluginInput{
		Name:    pluginName,
		Type:    api.PluginTypeSecrets,
		Version: "v1.2.3",
	})
	if err != nil {
		t.Fatalf("unexpected error deregistering v1.2.3: %v", err)
	}
	err = client.Sys().DeregisterPlugin(&api.DeregisterPluginInput{
		Name:    pluginName,
		Type:    api.PluginTypeSecrets,
		Version: "v1.3.0",
	})
	if err != nil {
		t.Fatalf("unexpected error deregistering v1.3.0: %v", err)
	}

	// Read without -version — should return a 404 (nil, nil from the server),
	// same as any other unregistered plugin. Must not contain "Available versions".
	_, err = client.Sys().GetPlugin(&api.GetPluginInput{
		Name: pluginName,
		Type: api.PluginTypeSecrets,
	})
	if err == nil {
		t.Fatal("expected 404 error for non-existent plugin, got nil")
	}
	if strings.Contains(err.Error(), "Available versions") {
		t.Errorf("expected no 'Available versions' in not-found error, got: %q", err.Error())
	}
}

func TestBackend_PluginMain_V4_Logical(t *testing.T) {
	args := []string{}
	// don't run as a standalone unit test
	if os.Getenv(pluginutil.PluginVaultVersionEnv) == "" {
		return
	}

	// don't run as a V5 plugin
	if os.Getenv(pluginutil.PluginAutoMTLSEnv) == "true" {
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

	// V4 does not support AutoMTLS so we set a TLSConfig via TLSProviderFunc
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

func TestBackend_PluginMain_Multiplexed_Logical(t *testing.T) {
	args := []string{}
	if os.Getenv(pluginutil.PluginVaultVersionEnv) == "" {
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

	factoryFunc := mock.FactoryType(logical.TypeLogical)

	err := lplugin.ServeMultiplex(&lplugin.ServeOpts{
		BackendFactoryFunc: factoryFunc,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestBackend_PluginMainLogical(t *testing.T) {
	args := []string{}
	if os.Getenv(pluginutil.PluginVaultVersionEnv) == "" {
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

	factoryFunc := mock.FactoryType(logical.TypeLogical)

	err := lplugin.Serve(&lplugin.ServeOpts{
		BackendFactoryFunc: factoryFunc,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestBackend_PluginMain_V4_Credentials(t *testing.T) {
	args := []string{}
	// don't run as a standalone unit test
	if os.Getenv(pluginutil.PluginVaultVersionEnv) == "" {
		return
	}

	// don't run as a V5 plugin
	if os.Getenv(pluginutil.PluginAutoMTLSEnv) == "true" {
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

	// V4 does not support AutoMTLS so we set a TLSConfig via TLSProviderFunc
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

func TestBackend_PluginMain_Multiplexed_Credentials(t *testing.T) {
	args := []string{}
	if os.Getenv(pluginutil.PluginVaultVersionEnv) == "" {
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

	factoryFunc := mock.FactoryType(logical.TypeCredential)

	err := lplugin.ServeMultiplex(&lplugin.ServeOpts{
		BackendFactoryFunc: factoryFunc,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestBackend_PluginMainCredentials(t *testing.T) {
	args := []string{}
	if os.Getenv(pluginutil.PluginVaultVersionEnv) == "" {
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

	factoryFunc := mock.FactoryType(logical.TypeCredential)

	err := lplugin.Serve(&lplugin.ServeOpts{
		BackendFactoryFunc: factoryFunc,
	})
	if err != nil {
		t.Fatal(err)
	}
}
