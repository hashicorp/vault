// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package plugin_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/api/auth/approle"
	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/builtin/logical/database"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/testhelpers/consul"
	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	"github.com/hashicorp/vault/helper/testhelpers/pluginhelpers"
	postgreshelper "github.com/hashicorp/vault/helper/testhelpers/postgresql"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func getCluster(t *testing.T, numCores int, types ...consts.PluginType) *vault.TestCluster {
	pluginDir := corehelpers.MakeTestPluginDir(t)
	coreConfig := &vault.CoreConfig{
		PluginDirectory: pluginDir,
		LogicalBackends: map[string]logical.Factory{
			"database": database.Factory,
		},
		AuditBackends: map[string]audit.Factory{
			"file": audit.NewFileBackend,
		},
	}

	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		TempDir:  pluginDir,
		NumCores: numCores,
		Plugins: func() []*vault.TestPluginConfig {
			var plugins []*vault.TestPluginConfig
			for _, typ := range types {
				plugins = append(plugins, &vault.TestPluginConfig{
					Typ:      typ,
					Versions: []string{""},
				})
			}
			return plugins
		}(),
		HandlerFunc: vaulthttp.Handler,
	})

	cluster.Start()
	vault.TestWaitActive(t, cluster.Cores[0].Core)

	return cluster
}

// TestExternalPlugin_RollbackAndReload ensures that we can successfully
// rollback and reload a plugin without triggering race conditions by the go
// race detector
func TestExternalPlugin_RollbackAndReload(t *testing.T) {
	pluginDir := corehelpers.MakeTestPluginDir(t)
	coreConfig := &vault.CoreConfig{
		// set rollback period to a short interval to make conditions more "racy"
		RollbackPeriod:  1 * time.Second,
		PluginDirectory: pluginDir,
	}

	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		TempDir:  pluginDir,
		NumCores: 1,
		Plugins: []*vault.TestPluginConfig{
			{
				Typ:      consts.PluginTypeSecrets,
				Versions: []string{""},
			},
		},
		HandlerFunc: vaulthttp.Handler,
	})

	cluster.Start()
	vault.TestWaitActive(t, cluster.Cores[0].Core)

	core := cluster.Cores[0]
	plugin := cluster.Plugins[0]
	client := core.Client
	client.SetToken(cluster.RootToken)

	testRegisterAndEnable(t, client, plugin)
	if _, err := client.Sys().ReloadPlugin(&api.ReloadPluginInput{
		Plugin: plugin.Name,
	}); err != nil {
		t.Fatal(err)
	}
}

func testRegisterVersion(t *testing.T, client *api.Client, plugin pluginhelpers.TestPlugin, version string) {
	t.Helper()
	if err := client.Sys().RegisterPlugin(&api.RegisterPluginInput{
		Name:    plugin.Name,
		Type:    api.PluginType(plugin.Typ),
		Command: plugin.Name,
		SHA256:  plugin.Sha256,
		Version: version,
	}); err != nil {
		t.Fatal(err)
	}
}

func testEnableVersion(t *testing.T, client *api.Client, plugin pluginhelpers.TestPlugin, version string) {
	t.Helper()
	switch plugin.Typ {
	case consts.PluginTypeSecrets:
		if err := client.Sys().Mount(plugin.Name, &api.MountInput{
			Type: plugin.Name,
			Config: api.MountConfigInput{
				PluginVersion: version,
			},
		}); err != nil {
			t.Fatal(err)
		}
	case consts.PluginTypeCredential:
		if err := client.Sys().EnableAuthWithOptions(plugin.Name, &api.EnableAuthOptions{
			Type: plugin.Name,
			Config: api.MountConfigInput{
				PluginVersion: version,
			},
		}); err != nil {
			t.Fatal(err)
		}
	}
}

func testRegisterAndEnable(t *testing.T, client *api.Client, plugin pluginhelpers.TestPlugin) {
	t.Helper()
	testRegisterVersion(t, client, plugin, plugin.Version)
	testEnableVersion(t, client, plugin, plugin.Version)
}

// TestExternalPlugin_ContinueOnError tests that vault can recover from a
// sha256 mismatch or missing plugin binary scenario
func TestExternalPlugin_ContinueOnError(t *testing.T) {
	t.Run("secret", func(t *testing.T) {
		t.Parallel()
		t.Run("sha256_mismatch", func(t *testing.T) {
			t.Parallel()
			testExternalPlugin_ContinueOnError(t, true, consts.PluginTypeSecrets)
		})

		t.Run("missing_plugin", func(t *testing.T) {
			t.Parallel()
			testExternalPlugin_ContinueOnError(t, false, consts.PluginTypeSecrets)
		})
	})

	t.Run("auth", func(t *testing.T) {
		t.Parallel()
		t.Run("sha256_mismatch", func(t *testing.T) {
			t.Parallel()
			testExternalPlugin_ContinueOnError(t, true, consts.PluginTypeCredential)
		})

		t.Run("missing_plugin", func(t *testing.T) {
			t.Parallel()
			testExternalPlugin_ContinueOnError(t, false, consts.PluginTypeCredential)
		})
	})
}

func testExternalPlugin_ContinueOnError(t *testing.T, mismatch bool, pluginType consts.PluginType) {
	cluster := getCluster(t, 1, pluginType)
	defer cluster.Cleanup()

	core := cluster.Cores[0]
	plugin := cluster.Plugins[0]
	client := core.Client
	client.SetToken(cluster.RootToken)

	testRegisterAndEnable(t, client, plugin)
	pluginPath := fmt.Sprintf("sys/plugins/catalog/%s/%s", pluginType, plugin.Name)
	// Get the registered plugin
	req := logical.TestRequest(t, logical.ReadOperation, pluginPath)
	req.ClientToken = core.Client.Token()
	resp, err := core.HandleRequest(namespace.RootContext(testCtx), req)
	if err != nil || resp == nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	command, ok := resp.Data["command"].(string)
	if !ok || command == "" {
		t.Fatal("invalid command")
	}

	// Trigger a sha256 mismatch or missing plugin error
	if mismatch {
		req = logical.TestRequest(t, logical.UpdateOperation, pluginPath)
		req.Data = map[string]interface{}{
			"sha256":  "d17bd7334758e53e6fbab15745d2520765c06e296f2ce8e25b7919effa0ac216",
			"command": filepath.Base(command),
		}
		req.ClientToken = core.Client.Token()
		resp, err = core.HandleRequest(namespace.RootContext(testCtx), req)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("err:%v resp:%#v", err, resp)
		}
	} else {
		err := os.Remove(filepath.Join(cluster.Cores[0].CoreConfig.PluginDirectory, filepath.Base(command)))
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

	// unmount
	switch pluginType {
	case consts.PluginTypeSecrets:
		if err := client.Sys().Unmount(plugin.Name); err != nil {
			t.Fatal(err)
		}
	case consts.PluginTypeCredential:
		if err := client.Sys().DisableAuth(plugin.Name); err != nil {
			t.Fatal(err)
		}
	}

	// Re-compile plugin
	var plugins []pluginhelpers.TestPlugin
	plugins = append(plugins, pluginhelpers.CompilePlugin(t, pluginType, "", core.CoreConfig.PluginDirectory))
	cluster.Plugins = plugins

	// Re-add the plugin to the catalog
	testRegisterAndEnable(t, client, plugin)

	// Reload the plugin
	req = logical.TestRequest(t, logical.UpdateOperation, "sys/plugins/reload/backend")
	req.Data = map[string]interface{}{
		"plugin": plugin.Name,
	}
	req.ClientToken = core.Client.Token()
	resp, err = core.HandleRequest(namespace.RootContext(testCtx), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	req = logical.TestRequest(t, logical.ReadOperation, pluginPath)
	req.ClientToken = core.Client.Token()
	resp, err = core.HandleRequest(namespace.RootContext(testCtx), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil {
		t.Fatalf("bad: response should not be nil")
	}
}

// TestExternalPlugin_AuthMethod tests that we can build, register and use an
// external auth method
func TestExternalPlugin_AuthMethod(t *testing.T) {
	cluster := getCluster(t, 5, consts.PluginTypeCredential)
	defer cluster.Cleanup()

	plugin := cluster.Plugins[0]
	client := cluster.Cores[0].Client
	client.SetToken(cluster.RootToken)

	// Register
	testRegisterVersion(t, client, plugin, plugin.Version)

	// define a group of parallel tests so we wait for their execution before
	// continuing on to cleanup
	// see: https://go.dev/blog/subtests
	t.Run("parallel execution group", func(t *testing.T) {
		// loop to mount 5 auth methods that will each share a single
		// plugin process
		for i := 0; i < 5; i++ {
			i := i
			pluginPath := fmt.Sprintf("%s-%d", plugin.Name, i)
			client := cluster.Cores[i].Client
			t.Run(pluginPath, func(t *testing.T) {
				t.Parallel()
				client.SetToken(cluster.RootToken)
				// Enable
				if err := client.Sys().EnableAuthWithOptions(pluginPath, &api.EnableAuthOptions{
					Type: plugin.Name,
				}); err != nil {
					t.Fatal(err)
				}

				// Configure
				_, err := client.Logical().Write("auth/"+pluginPath+"/role/role1", map[string]interface{}{
					"bind_secret_id": "true",
					"period":         "300",
				})
				if err != nil {
					t.Fatal(err)
				}

				secret, err := client.Logical().Write("auth/"+pluginPath+"/role/role1/secret-id", nil)
				if err != nil {
					t.Fatal(err)
				}
				secretID := secret.Data["secret_id"].(string)

				secret, err = client.Logical().Read("auth/" + pluginPath + "/role/role1/role-id")
				if err != nil {
					t.Fatal(err)
				}
				roleID := secret.Data["role_id"].(string)

				// Login - expect SUCCESS
				authMethod, err := approle.NewAppRoleAuth(
					roleID,
					&approle.SecretID{FromString: secretID},
					approle.WithMountPath(pluginPath),
				)
				if err != nil {
					t.Fatal(err)
				}
				_, err = client.Auth().Login(context.Background(), authMethod)
				if err != nil {
					t.Fatal(err)
				}

				// Renew
				resp, err := client.Auth().Token().RenewSelf(30)
				if err != nil {
					t.Fatal(err)
				}

				// Login - expect SUCCESS
				resp, err = client.Auth().Login(context.Background(), authMethod)
				if err != nil {
					t.Fatal(err)
				}

				revokeToken := resp.Auth.ClientToken
				// Revoke
				if err = client.Auth().Token().RevokeSelf(revokeToken); err != nil {
					t.Fatal(err)
				}

				// Reset root token
				client.SetToken(cluster.RootToken)

				// Lookup - expect FAILURE
				resp, err = client.Auth().Token().Lookup(revokeToken)
				if err == nil {
					t.Fatalf("expected error, got nil")
				}

				// Reset root token
				client.SetToken(cluster.RootToken)
			})
		}
	})

	// Deregister
	if err := client.Sys().DeregisterPlugin(&api.DeregisterPluginInput{
		Name:    plugin.Name,
		Type:    api.PluginType(plugin.Typ),
		Version: plugin.Version,
	}); err != nil {
		t.Fatal(err)
	}
}

// TestExternalPlugin_AuthMethodReload tests that we can use an external auth
// method after reload
func TestExternalPlugin_AuthMethodReload(t *testing.T) {
	cluster := getCluster(t, 1, consts.PluginTypeCredential)
	defer cluster.Cleanup()

	plugin := cluster.Plugins[0]
	client := cluster.Cores[0].Client
	client.SetToken(cluster.RootToken)

	testRegisterAndEnable(t, client, plugin)

	// Configure
	_, err := client.Logical().Write("auth/"+plugin.Name+"/role/role1", map[string]interface{}{
		"bind_secret_id": "true",
		"period":         "300",
	})
	if err != nil {
		t.Fatal(err)
	}

	secret, err := client.Logical().Write("auth/"+plugin.Name+"/role/role1/secret-id", nil)
	if err != nil {
		t.Fatal(err)
	}
	secretID := secret.Data["secret_id"].(string)

	secret, err = client.Logical().Read("auth/" + plugin.Name + "/role/role1/role-id")
	if err != nil {
		t.Fatal(err)
	}
	roleID := secret.Data["role_id"].(string)

	// Login - expect SUCCESS
	authMethod, err := approle.NewAppRoleAuth(
		roleID,
		&approle.SecretID{FromString: secretID},
		approle.WithMountPath(plugin.Name),
	)
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Auth().Login(context.Background(), authMethod)
	if err != nil {
		t.Fatal(err)
	}

	// Reset root token
	client.SetToken(cluster.RootToken)

	// Reload plugin
	if _, err := client.Sys().ReloadPlugin(&api.ReloadPluginInput{
		Plugin: plugin.Name,
	}); err != nil {
		t.Fatal(err)
	}

	_, err = client.Auth().Login(context.Background(), authMethod)
	if err != nil {
		t.Fatal(err)
	}

	// Reset root token
	client.SetToken(cluster.RootToken)

	// Deregister
	if err := client.Sys().DeregisterPlugin(&api.DeregisterPluginInput{
		Name:    plugin.Name,
		Type:    api.PluginType(plugin.Typ),
		Version: plugin.Version,
	}); err != nil {
		t.Fatal(err)
	}
}

// TestExternalPlugin_SecretsEngine tests that we can build, register and use an
// external secrets engine
func TestExternalPlugin_SecretsEngine(t *testing.T) {
	cluster := getCluster(t, 1, consts.PluginTypeSecrets)
	defer cluster.Cleanup()

	plugin := cluster.Plugins[0]
	client := cluster.Cores[0].Client
	client.SetToken(cluster.RootToken)

	// Register
	testRegisterVersion(t, client, plugin, plugin.Version)

	// define a group of parallel tests so we wait for their execution before
	// continuing on to cleanup
	// see: https://go.dev/blog/subtests
	t.Run("parallel execution group", func(t *testing.T) {
		// loop to mount 5 secrets engines that will each share a single
		// plugin process
		for i := 0; i < 5; i++ {
			pluginPath := fmt.Sprintf("%s-%d", plugin.Name, i)
			t.Run(pluginPath, func(t *testing.T) {
				t.Parallel()
				// Enable
				if err := client.Sys().Mount(pluginPath, &api.MountInput{
					Type: plugin.Name,
				}); err != nil {
					t.Fatal(err)
				}

				// Configure
				cleanupConsul, consulConfig := consul.PrepareTestContainer(t, "", false, true)
				defer cleanupConsul()

				_, err := client.Logical().Write(pluginPath+"/config/access", map[string]interface{}{
					"address": consulConfig.Address(),
					"token":   consulConfig.Token,
				})
				if err != nil {
					t.Fatal(err)
				}

				_, err = client.Logical().Write(pluginPath+"/roles/test", map[string]interface{}{
					"consul_policies": []string{"test"},
					"ttl":             "6h",
					"local":           false,
				})
				if err != nil {
					t.Fatal(err)
				}

				resp, err := client.Logical().Read(pluginPath + "/creds/test")
				if err != nil {
					t.Fatal(err)
				}
				if resp == nil {
					t.Fatal("read creds response is nil")
				}
			})
		}
	})

	// Deregister
	if err := client.Sys().DeregisterPlugin(&api.DeregisterPluginInput{
		Name:    plugin.Name,
		Type:    api.PluginType(plugin.Typ),
		Version: plugin.Version,
	}); err != nil {
		t.Fatal(err)
	}
}

// TestExternalPlugin_SecretsEngineReload tests that we can use an external
// secrets engine after reload
func TestExternalPlugin_SecretsEngineReload(t *testing.T) {
	cluster := getCluster(t, 1, consts.PluginTypeSecrets)
	defer cluster.Cleanup()

	plugin := cluster.Plugins[0]
	client := cluster.Cores[0].Client
	client.SetToken(cluster.RootToken)

	testRegisterAndEnable(t, client, plugin)

	// Configure
	cleanupConsul, consulConfig := consul.PrepareTestContainer(t, "", false, true)
	defer cleanupConsul()

	_, err := client.Logical().Write(plugin.Name+"/config/access", map[string]interface{}{
		"address": consulConfig.Address(),
		"token":   consulConfig.Token,
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Write(plugin.Name+"/roles/test", map[string]interface{}{
		"consul_policies": []string{"test"},
		"ttl":             "6h",
		"local":           false,
	})
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.Logical().Read(plugin.Name + "/creds/test")
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("read creds response is nil")
	}

	// Reload plugin
	if _, err := client.Sys().ReloadPlugin(&api.ReloadPluginInput{
		Plugin: plugin.Name,
	}); err != nil {
		t.Fatal(err)
	}

	resp, err = client.Logical().Read(plugin.Name + "/creds/test")
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("read creds response is nil")
	}

	// Deregister
	if err := client.Sys().DeregisterPlugin(&api.DeregisterPluginInput{
		Name:    plugin.Name,
		Type:    api.PluginType(plugin.Typ),
		Version: plugin.Version,
	}); err != nil {
		t.Fatal(err)
	}
}

// TestExternalPlugin_Database tests that we can build, register and use an
// external database secrets engine
func TestExternalPlugin_Database(t *testing.T) {
	cluster := getCluster(t, 1, consts.PluginTypeDatabase)
	defer cluster.Cleanup()

	plugin := cluster.Plugins[0]
	client := cluster.Cores[0].Client
	client.SetToken(cluster.RootToken)

	// Register
	testRegisterVersion(t, client, plugin, plugin.Version)

	// Enable
	if err := client.Sys().Mount(consts.PluginTypeDatabase.String(), &api.MountInput{
		Type: consts.PluginTypeDatabase.String(),
	}); err != nil {
		t.Fatal(err)
	}

	// define a group of parallel tests so we wait for their execution before
	// continuing on to cleanup
	// see: https://go.dev/blog/subtests
	t.Run("parallel execution group", func(t *testing.T) {
		// loop to mount 5 database connections that will each share a single
		// plugin process
		for i := 0; i < 5; i++ {
			dbName := fmt.Sprintf("%s-%d", plugin.Name, i)
			t.Run(dbName, func(t *testing.T) {
				t.Parallel()
				roleName := "test-role-" + dbName

				cleanupContainer, connURL := postgreshelper.PrepareTestContainerWithVaultUser(t, context.Background(), "13.4-buster")
				defer cleanupContainer()

				_, err := client.Logical().Write("database/config/"+dbName, map[string]interface{}{
					"connection_url": connURL,
					"plugin_name":    plugin.Name,
					"allowed_roles":  []string{roleName},
					"username":       "vaultadmin",
					"password":       "vaultpass",
				})
				if err != nil {
					t.Fatal(err)
				}

				_, err = client.Logical().Write("database/rotate-root/"+dbName, map[string]interface{}{})
				if err != nil {
					t.Fatal(err)
				}

				_, err = client.Logical().Write("database/roles/"+roleName, map[string]interface{}{
					"db_name":             dbName,
					"creation_statements": testRole,
					"max_ttl":             "10m",
				})
				if err != nil {
					t.Fatal(err)
				}

				// Generate credentials
				resp, err := client.Logical().Read("database/creds/" + roleName)
				if err != nil {
					t.Fatal(err)
				}
				if resp == nil {
					t.Fatal("read creds response is nil")
				}

				_, err = client.Logical().Write("database/reset/"+dbName, map[string]interface{}{})
				if err != nil {
					t.Fatal(err)
				}

				// Generate credentials
				resp, err = client.Logical().Read("database/creds/" + roleName)
				if err != nil {
					t.Fatal(err)
				}
				if resp == nil {
					t.Fatal("read creds response is nil")
				}

				resp, err = client.Logical().Read("database/creds/" + roleName)
				if err != nil {
					t.Fatal(err)
				}
				if resp == nil {
					t.Fatal("read creds response is nil")
				}

				revokeLease := resp.LeaseID
				// Lookup - expect SUCCESS
				resp, err = client.Sys().Lookup(revokeLease)
				if err != nil {
					t.Fatal(err)
				}
				if resp == nil {
					t.Fatalf("lease lookup response is nil")
				}

				// Revoke
				if err = client.Sys().Revoke(revokeLease); err != nil {
					t.Fatal(err)
				}

				// Reset root token
				client.SetToken(cluster.RootToken)

				// Lookup - expect FAILURE
				_, err = client.Sys().Lookup(revokeLease)
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
			})
		}
	})

	// Deregister
	if err := client.Sys().DeregisterPlugin(&api.DeregisterPluginInput{
		Name:    plugin.Name,
		Type:    api.PluginType(plugin.Typ),
		Version: plugin.Version,
	}); err != nil {
		t.Fatal(err)
	}
}

// TestExternalPlugin_DatabaseReload tests that we can use an external database
// secrets engine after reload
func TestExternalPlugin_DatabaseReload(t *testing.T) {
	cluster := getCluster(t, 1, consts.PluginTypeDatabase)
	defer cluster.Cleanup()

	plugin := cluster.Plugins[0]
	client := cluster.Cores[0].Client
	client.SetToken(cluster.RootToken)

	// Register
	testRegisterVersion(t, client, plugin, plugin.Version)

	// Enable
	if err := client.Sys().Mount(consts.PluginTypeDatabase.String(), &api.MountInput{
		Type: consts.PluginTypeDatabase.String(),
	}); err != nil {
		t.Fatal(err)
	}

	dbName := fmt.Sprintf("%s-%d", plugin.Name, 0)
	roleName := "test-role-" + dbName

	cleanupContainer, connURL := postgreshelper.PrepareTestContainerWithVaultUser(t, context.Background(), "13.4-buster")
	defer cleanupContainer()

	_, err := client.Logical().Write("database/config/"+dbName, map[string]interface{}{
		"connection_url": connURL,
		"plugin_name":    plugin.Name,
		"allowed_roles":  []string{roleName},
		"username":       "vaultadmin",
		"password":       "vaultpass",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Write("database/roles/"+roleName, map[string]interface{}{
		"db_name":             dbName,
		"creation_statements": testRole,
		"max_ttl":             "10m",
	})
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.Logical().Read("database/creds/" + roleName)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("read creds response is nil")
	}

	// Reload plugin
	if _, err := client.Sys().ReloadPlugin(&api.ReloadPluginInput{
		Plugin: plugin.Name,
	}); err != nil {
		t.Fatal(err)
	}

	// Generate credentials after reload
	resp, err = client.Logical().Read("database/creds/" + roleName)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("read creds response is nil")
	}

	// Deregister
	if err := client.Sys().DeregisterPlugin(&api.DeregisterPluginInput{
		Name:    plugin.Name,
		Type:    api.PluginType(plugin.Typ),
		Version: plugin.Version,
	}); err != nil {
		t.Fatal(err)
	}
}

const testRole = `
CREATE ROLE "{{name}}" WITH
  LOGIN
  PASSWORD '{{password}}'
  VALID UNTIL '{{expiration}}';
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO "{{name}}";
`

func testExternalPluginMetadataAuditLog(t *testing.T, log map[string]interface{}, expectedMountClass string) {
	if mountClass, ok := log["mount_class"].(string); !ok {
		t.Fatalf("mount_class should be a string, not %T", log["mount_class"])
	} else if mountClass != expectedMountClass {
		t.Fatalf("bad: mount_class should be %s, not %s", expectedMountClass, mountClass)
	}

	if mountIsExternalPlugin, ok := log["mount_is_external_plugin"].(bool); !ok {
		t.Fatalf("mount_is_external_plugin should be a bool, not %T", log["mount_is_external_plugin"])
	} else if !mountIsExternalPlugin {
		t.Fatalf("bad: mount_is_external_plugin should be true, not %t", mountIsExternalPlugin)
	}

	if _, ok := log["mount_running_sha256"].(string); !ok {
		t.Fatalf("mount_running_sha256 should be a string, not %T", log["mount_running_sha256"])
	}
}

// TestExternalPlugin_AuditEnabled_ShouldLogPluginMetadata_Auth tests that we have plugin metadata of an auth plugin
// in audit log when it is enabled
func TestExternalPlugin_AuditEnabled_ShouldLogPluginMetadata_Auth(t *testing.T) {
	cluster := getCluster(t, 1, consts.PluginTypeCredential)
	defer cluster.Cleanup()

	plugin := cluster.Plugins[0]
	client := cluster.Cores[0].Client
	client.SetToken(cluster.RootToken)

	testRegisterAndEnable(t, client, plugin)

	// Enable the audit backend
	tempDir := t.TempDir()
	auditLogFile, err := os.CreateTemp(tempDir, "")
	if err != nil {
		t.Fatal(err)
	}
	defer auditLogFile.Close()

	err = client.Sys().EnableAuditWithOptions("file", &api.EnableAuditOptions{
		Type: "file",
		Options: map[string]string{
			"file_path": auditLogFile.Name(),
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Write("auth/"+plugin.Name+"/role/role1", map[string]interface{}{
		"bind_secret_id": "true",
		"period":         "300",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check the audit trail on request and response
	decoder := json.NewDecoder(auditLogFile)
	var auditRecord map[string]interface{}
	for decoder.Decode(&auditRecord) == nil {
		auditRequest := map[string]interface{}{}
		if req, ok := auditRecord["request"]; ok {
			auditRequest = req.(map[string]interface{})
			if auditRequest["path"] != "auth/"+plugin.Name+"/role/role1" {
				continue
			}
		}
		testExternalPluginMetadataAuditLog(t, auditRequest, consts.PluginTypeCredential.String())

		auditResponse := map[string]interface{}{}
		if req, ok := auditRecord["response"]; ok {
			auditRequest = req.(map[string]interface{})
			if auditResponse["path"] != "auth/"+plugin.Name+"/role/role1" {
				continue
			}
		}
		testExternalPluginMetadataAuditLog(t, auditResponse, consts.PluginTypeCredential.String())
	}

	// Deregister
	if err := client.Sys().DeregisterPlugin(&api.DeregisterPluginInput{
		Name:    plugin.Name,
		Type:    api.PluginType(plugin.Typ),
		Version: plugin.Version,
	}); err != nil {
		t.Fatal(err)
	}
}

// TestExternalPlugin_AuditEnabled_ShouldLogPluginMetadata_Secret tests that we have plugin metadata of a secret plugin
// in audit log when it is enabled
func TestExternalPlugin_AuditEnabled_ShouldLogPluginMetadata_Secret(t *testing.T) {
	cluster := getCluster(t, 1, consts.PluginTypeSecrets)
	defer cluster.Cleanup()

	plugin := cluster.Plugins[0]
	client := cluster.Cores[0].Client
	client.SetToken(cluster.RootToken)

	testRegisterAndEnable(t, client, plugin)

	// Enable the audit backend
	tempDir := t.TempDir()
	auditLogFile, err := os.CreateTemp(tempDir, "")
	if err != nil {
		t.Fatal(err)
	}
	defer auditLogFile.Close()

	err = client.Sys().EnableAuditWithOptions("file", &api.EnableAuditOptions{
		Type: "file",
		Options: map[string]string{
			"file_path": auditLogFile.Name(),
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	// Configure
	cleanupConsul, consulConfig := consul.PrepareTestContainer(t, "", false, true)
	defer cleanupConsul()
	_, err = client.Logical().Write(plugin.Name+"/config/access", map[string]interface{}{
		"address": consulConfig.Address(),
		"token":   consulConfig.Token,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check the audit trail on request and response
	decoder := json.NewDecoder(auditLogFile)
	var auditRecord map[string]interface{}
	for decoder.Decode(&auditRecord) == nil {
		auditRequest := map[string]interface{}{}
		if req, ok := auditRecord["request"]; ok {
			auditRequest = req.(map[string]interface{})
			if auditRequest["path"] != plugin.Name+"/config/access" {
				continue
			}
		}
		testExternalPluginMetadataAuditLog(t, auditRequest, consts.PluginTypeSecrets.String())

		auditResponse := map[string]interface{}{}
		if req, ok := auditRecord["response"]; ok {
			auditRequest = req.(map[string]interface{})
			if auditResponse["path"] != plugin.Name+"/config/access" {
				continue
			}
		}
		testExternalPluginMetadataAuditLog(t, auditResponse, consts.PluginTypeSecrets.String())
	}

	// Deregister
	if err := client.Sys().DeregisterPlugin(&api.DeregisterPluginInput{
		Name:    plugin.Name,
		Type:    api.PluginType(plugin.Typ),
		Version: plugin.Version,
	}); err != nil {
		t.Fatal(err)
	}
}

func testPin(t *testing.T, client *api.Client, op logical.Operation, pin *pluginutil.PinnedVersion) *api.Secret {
	t.Helper()
	switch op {
	case logical.CreateOperation, logical.UpdateOperation:
		resp, err := client.Logical().Write(fmt.Sprintf("sys/plugins/pins/%s/%s", pin.Type.String(), pin.Name), map[string]any{
			"version": pin.Version,
		})
		if err != nil {
			t.Fatal(err)
		}
		return resp
	case logical.DeleteOperation:
		resp, err := client.Logical().Delete(fmt.Sprintf("sys/plugins/pins/%s/%s", pin.Type.String(), pin.Name))
		if err != nil {
			t.Fatal(err)
		}
		return resp
	default:
		t.Fatal("unsupported operation")
		// Satisfy the compiler that there's no escape from the switch statement.
		return nil
	}
}

func testReload(t *testing.T, client *api.Client, plugin pluginhelpers.TestPlugin) {
	_, err := client.Sys().RootReloadPlugin(context.Background(), &api.RootReloadPluginInput{
		Plugin: plugin.Name,
		Type:   api.PluginType(plugin.Typ),
	})
	if err != nil {
		t.Fatal(err)
	}
}

func expectRunningVersion(t *testing.T, client *api.Client, plugin pluginhelpers.TestPlugin, expectedVersion string) {
	t.Helper()
	switch plugin.Typ {
	case consts.PluginTypeCredential:
		auth, err := client.Logical().Read("sys/auth/" + plugin.Name)
		if err != nil {
			t.Fatal(err)
		}
		if auth.Data["running_plugin_version"] != expectedVersion {
			t.Fatalf("expected running_plugin_version to be %s but got %s", expectedVersion, auth.Data["running_plugin_version"])
		}
	case consts.PluginTypeSecrets:
		mount, err := client.Logical().Read("sys/mounts/" + plugin.Name)
		if err != nil {
			t.Fatal(err)
		}
		if mount.Data["running_plugin_version"] != expectedVersion {
			t.Fatalf("expected running_plugin_version to be %s but got %s", expectedVersion, mount.Data["running_plugin_version"])
		}
	case consts.PluginTypeDatabase:
		resp, err := client.Logical().Read("database/config/" + plugin.Name)
		if err != nil {
			t.Fatal(err)
		}
		if resp.Data["running_plugin_version"] != expectedVersion {
			t.Fatalf("expected running_plugin_version to be %s but got %s", expectedVersion, resp.Data["running_plugin_version"])
		}
		expectedWarnings := 0
		if resp.Data["plugin_version"] != resp.Data["running_plugin_version"] {
			expectedWarnings = 1
		}

		if expectedWarnings != len(resp.Warnings) {
			t.Fatalf("expected %d warning(s) but got %v", expectedWarnings, resp.Warnings)
		}
	default:
		t.Fatal("unsupported plugin type")
	}
}

// TestCore_UpgradePluginUsingPinnedVersion_AuthAndSecret tests a full workflow
// of upgrading an external plugin gated by pinned versions.
func TestCore_UpgradePluginUsingPinnedVersion_AuthAndSecret(t *testing.T) {
	cluster := getCluster(t, 1, consts.PluginTypeCredential, consts.PluginTypeSecrets)
	t.Cleanup(cluster.Cleanup)

	client := cluster.Cores[0].Client

	for name, idx := range map[string]int{
		"credential plugin": 0,
		"secrets plugin":    1,
	} {
		t.Run(name, func(t *testing.T) {
			plugin := cluster.Plugins[idx]

			// Register the same plugin with two versions.
			for _, version := range []string{"v1.0.0", "v1.0.1"} {
				testRegisterVersion(t, client, plugin, version)
			}

			pin101 := &pluginutil.PinnedVersion{
				Name:    plugin.Name,
				Type:    plugin.Typ,
				Version: "v1.0.1",
			}

			// Mount 1.0.0 then pin to 1.0.1
			testEnableVersion(t, client, plugin, "v1.0.0")
			testPin(t, client, logical.CreateOperation, pin101)
			expectRunningVersion(t, client, plugin, "v1.0.0")

			// Pinned version should be in effect after reloading.
			testReload(t, client, plugin)
			expectRunningVersion(t, client, plugin, "v1.0.1")

			// Deregistering a pinned plugin should fail.
			if err := client.Sys().DeregisterPlugin(&api.DeregisterPluginInput{
				Name:    plugin.Name,
				Type:    api.PluginType(plugin.Typ),
				Version: "v1.0.1",
			}); err == nil {
				t.Fatal("expected error, got nil")
			}

			// Now delete, reload, and we should be back to 1.0.0
			testPin(t, client, logical.DeleteOperation, pin101)
			testReload(t, client, plugin)
			expectRunningVersion(t, client, plugin, "v1.0.0")
		})
	}
}

// TestCore_UpgradePluginUsingPinnedVersion_Database tests a full workflow
// of upgrading an external database plugin gated by pinned versions.
func TestCore_UpgradePluginUsingPinnedVersion_Database(t *testing.T) {
	cluster := getCluster(t, 3, consts.PluginTypeDatabase)
	t.Cleanup(cluster.Cleanup)

	client := cluster.Cores[0].Client
	plugin := cluster.Plugins[0]

	// Register the same plugin with two versions.
	for _, version := range []string{"v1.0.0", "v1.0.1"} {
		testRegisterVersion(t, client, plugin, version)
	}

	pin101 := &pluginutil.PinnedVersion{
		Name:    plugin.Name,
		Type:    plugin.Typ,
		Version: "v1.0.1",
	}

	// Enable the combined db engine first.
	if err := client.Sys().Mount(consts.PluginTypeDatabase.String(), &api.MountInput{
		Type: consts.PluginTypeDatabase.String(),
	}); err != nil {
		t.Fatal(err)
	}

	cleanupPG, connURL := postgreshelper.PrepareTestContainerWithVaultUser(t, context.Background(), "13.4-buster")
	t.Cleanup(cleanupPG)

	// Mount 1.0.0 then pin to 1.0.1
	_, err := client.Logical().Write("database/config/"+plugin.Name, map[string]interface{}{
		"plugin_name":    plugin.Name,
		"plugin_version": "v1.0.0",
		"connection_url": connURL,
		"username":       "vaultadmin",
		"password":       "vaultpass",
	})
	if err != nil {
		t.Fatal(err)
	}
	testPin(t, client, logical.CreateOperation, pin101)
	expectRunningVersion(t, client, plugin, "v1.0.0")

	// Pinned version should be in effect after reloading.
	testReload(t, client, plugin)
	// All nodes in the cluster should report the same info, because although
	// the running_plugin_version info is local to the leader, the standbys
	// should forward the request to the leader.
	for i := 0; i < 3; i++ {
		expectRunningVersion(t, cluster.Cores[i].Client, plugin, "v1.0.1")
	}

	// Now delete, reload, and we should be back to 1.0.0
	testPin(t, client, logical.DeleteOperation, pin101)
	testReload(t, client, plugin)
	for i := 0; i < 3; i++ {
		expectRunningVersion(t, client, plugin, "v1.0.0")
	}
}
