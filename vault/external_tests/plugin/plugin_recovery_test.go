// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package plugin_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/plugin/mock"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/require"
)

// TestPlugin_ReloadAfterFailedSetup verifies that a mount skipped during unseal
// retains its data and becomes writable after its binary or registration is restored.
func TestPlugin_ReloadAfterFailedSetup(t *testing.T) {
	t.Parallel()

	for _, backend := range []struct {
		name       string
		pluginType consts.PluginType
		main       string
		mountPath  string
	}{
		{"secrets", consts.PluginTypeSecrets, "TestBackend_PluginMain_Multiplexed_Logical", "recoverable/"},
		{"auth", consts.PluginTypeCredential, "TestBackend_PluginMain_Multiplexed_Credentials", "auth/recoverable/"},
	} {
		for _, failure := range []string{"missing_binary", "missing_registration"} {
			t.Run(backend.name+"/"+failure, func(t *testing.T) {
				t.Parallel()
				ctx := t.Context()
				pluginDir := corehelpers.MakeTestPluginDir(t)
				cluster := vault.NewTestCluster(t, &vault.CoreConfig{
					PluginDirectory: pluginDir,
				}, &vault.TestClusterOptions{
					NumCores:    1,
					TempDir:     pluginDir,
					HandlerFunc: vaulthttp.Handler,
				})
				client := cluster.Cores[0].Client
				const pluginName = "recoverable-plugin"
				const pluginVersion = "v1.0.0"
				env := []string{
					pluginutil.PluginCACertPEMEnv + "=" + cluster.CACertPEMFile,
					mock.MockPluginVersionEnv + "=" + pluginVersion,
				}
				vault.TestAddTestPlugin(t, cluster.Cores[0].Core, pluginName, backend.pluginType, pluginVersion, backend.main, env)
				pluginType := api.PluginType(backend.pluginType)
				registered, catalogErr := client.Sys().GetPluginWithContext(ctx, &api.GetPluginInput{
					Name: pluginName, Type: pluginType, Version: pluginVersion,
				})
				require.NoError(t, catalogErr, "read the plugin registration before disrupting setup")
				require.NotNil(t, registered, "the test plugin must be registered")

				switch backend.pluginType {
				case consts.PluginTypeSecrets:
					require.NoError(t, client.Sys().MountWithContext(ctx, "recoverable", &api.MountInput{
						Type: pluginName, Config: api.MountConfigInput{PluginVersion: pluginVersion},
					}), "enable the secrets plugin before disrupting setup")
				case consts.PluginTypeCredential:
					require.NoError(t, client.Sys().EnableAuthWithOptionsWithContext(ctx, "recoverable", &api.EnableAuthOptions{
						Type: pluginName, Config: api.MountConfigInput{PluginVersion: pluginVersion},
					}), "enable the auth plugin before disrupting setup")
				}

				dataPath := backend.mountPath + "kv/preserved"
				_, seedErr := client.Logical().WriteWithContext(ctx, dataPath, map[string]interface{}{"value": "before"})
				require.NoError(t, seedErr, "the plugin must accept storage writes before the setup failure")

				binaryPath := filepath.Join(pluginDir, filepath.Base(registered.Command))
				movedPath := binaryPath + ".unavailable"
				switch failure {
				case "missing_binary":
					require.NoError(t, os.Rename(binaryPath, movedPath), "make the registered plugin binary unavailable")
				case "missing_registration":
					require.NoError(t, client.Sys().DeregisterPluginWithContext(ctx, &api.DeregisterPluginInput{
						Name: pluginName, Type: pluginType, Version: pluginVersion,
					}), "remove the registration while preserving the mount")
				}

				cluster.EnsureCoresSealed(t)
				cluster.UnsealCores(t)

				_, unavailableErr := client.Logical().WriteWithContext(ctx, dataPath, map[string]interface{}{"value": "unavailable"})
				require.ErrorContains(t, unavailableErr, "backend is nil", "unseal must retain the route without instantiating the unavailable plugin")
				_, failedReloadErr := client.Sys().ReloadPluginWithContext(ctx, &api.ReloadPluginInput{Mounts: []string{backend.mountPath}})
				require.Error(t, failedReloadErr, "reload must fail while the binary or registration is still unavailable")

				switch failure {
				case "missing_binary":
					require.NoError(t, os.Rename(movedPath, binaryPath), "restore the original plugin binary")
				case "missing_registration":
					_, registerErr := client.Sys().RegisterPluginWithContextDetailed(ctx, &api.RegisterPluginInput{
						Name: pluginName, Type: pluginType, Version: pluginVersion,
						Command: filepath.Base(registered.Command), SHA256: registered.SHA256,
						Args: registered.Args, Env: env,
					})
					require.NoError(t, registerErr, "restore the plugin registration")
				}

				_, reloadErr := client.Sys().ReloadPluginWithContext(ctx, &api.ReloadPluginInput{Mounts: []string{backend.mountPath}})
				require.NoError(t, reloadErr, "recover the existing mount without another seal or restart")
				preserved, readErr := client.Logical().ReadWithContext(ctx, dataPath)
				require.NoError(t, readErr, "read data persisted before the setup failure")
				require.NotNil(t, preserved, "reloading must preserve the mount's existing data")
				require.Equal(t, "before", preserved.Data["value"], "failed writes must not change the preserved data")

				_, writeErr := client.Logical().WriteWithContext(ctx, dataPath, map[string]interface{}{"value": "after"})
				require.NoError(t, writeErr, "a recovered mount must not retain the setup read-only guard")
				updated, verifyErr := client.Logical().ReadWithContext(ctx, dataPath)
				require.NoError(t, verifyErr, "read the recovered mount's updated data")
				require.NotNil(t, updated, "the recovered write must persist data")
				require.Equal(t, "after", updated.Data["value"], "the recovered write must update storage")
				_, deleteErr := client.Logical().DeleteWithContext(ctx, dataPath)
				require.NoError(t, deleteErr, "a recovered mount must also allow storage deletes")
				deleted, finalReadErr := client.Logical().ReadWithContext(ctx, dataPath)
				require.NoError(t, finalReadErr, "read the deleted key")
				require.Nil(t, deleted, "the recovered delete must remove the stored key")
			})
		}
	}
}
