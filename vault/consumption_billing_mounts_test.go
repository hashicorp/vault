// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	logicalKv "github.com/hashicorp/vault-plugin-secrets-kv"
	"github.com/hashicorp/vault/builtin/logical/totp"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/pluginconsts"
	"github.com/hashicorp/vault/helper/testhelpers/pluginhelpers"
	sdkconsts "github.com/hashicorp/vault/sdk/helper/consts"
	sdkpluginutil "github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

// TestCountMetricsFromMounts_LocalReplicatedMounts tests that CountMetricsFromMounts correctly separates
// local and replicated mount metrics and returns the expected counts.
func TestCountMetricsFromMounts_LocalReplicatedMounts(t *testing.T) {
	// Set up core with 2 role backends, 1 managed key backend (TOTP), and KV backend
	coreConfig := &CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			pluginconsts.SecretEngineAWS:      roleLogicalBackends[pluginconsts.SecretEngineAWS],
			pluginconsts.SecretEngineDatabase: roleLogicalBackends[pluginconsts.SecretEngineDatabase],
			pluginconsts.SecretEngineKV:       logicalKv.Factory,
			pluginconsts.SecretEngineTOTP:     totp.Factory,
		},
	}
	core, _, root := TestCoreUnsealedWithConfig(t, coreConfig)
	ctx := namespace.RootContext(context.Background())

	// Create replicated and local mounts
	mounts := []struct {
		mountType string
		path      string
		keyName   string
		local     bool
		numKeys   int
	}{
		{pluginconsts.SecretEngineAWS, "aws", "role/", false, 10},
		{pluginconsts.SecretEngineDatabase, "db", "role/", false, 10},
		{pluginconsts.SecretEngineAWS, "local-aws", "static-roles/", true, 10},
		{pluginconsts.SecretEngineDatabase, "local-db", "static-role/", true, 10},
		{pluginconsts.SecretEngineTOTP, "totp", "my-key-%d", false, 10},
		{pluginconsts.SecretEngineTOTP, "local-totp", "my-key-%d", true, 10},
		{"kv-v1", "kv", "secret-%d", false, 10},
		{"kv-v2", "local-kv", "secret-%d", true, 10},
	}

	for _, mount := range mounts {
		createMount(t, ctx, core, root, mount.mountType, mount.path, mount.local)
	}

	// Sleep to prevent race conditions during the role initialization
	time.Sleep(1 * time.Second)

	core.mountsLock.RLock()
	defer core.mountsLock.RUnlock()

	// Add roles, managed keys, and kv secrets
	for _, mount := range mounts {
		switch mount.mountType {
		case pluginconsts.SecretEngineAWS, pluginconsts.SecretEngineDatabase:
			addRoleToStorage(t, core, mount.path, mount.keyName, mount.numKeys)
		case pluginconsts.SecretEngineTOTP:
			for i := 0; i < mount.numKeys; i++ {
				keyName := fmt.Sprintf(mount.keyName, i)
				addTotpKeyToStorage(t, ctx, core, mount.path, root, keyName)
			}
		case "kv-v1", "kv-v2":
			for i := 0; i < mount.numKeys; i++ {
				secretName := fmt.Sprintf(mount.keyName, i)
				addKvSecretToStorage(t, ctx, core, mount.path, root, secretName, mount.mountType)
			}
		}
	}

	metrics, err := core.CountMetricsFromMounts(true)
	require.NoError(t, err)
	require.NotNil(t, metrics)

	// Only fields for which the test added metrics should have counts, all others should be 0
	expectedReplicatedRoleCounts := &RoleCounts{
		AWSDynamicRoles:      10,
		DatabaseDynamicRoles: 10,
	}
	expectedLocalRoleCounts := &RoleCounts{
		AWSStaticRoles:      10,
		DatabaseStaticRoles: 10,
	}
	expectedReplicatedManagedKeys := &ManagedKeyCounts{
		TotpKeys: 10,
	}
	expectedLocalManagedKeys := &ManagedKeyCounts{
		TotpKeys: 10,
	}

	// Verify counts
	require.Equal(t, expectedReplicatedRoleCounts, metrics.ReplicatedRoleCounts)
	require.Equal(t, expectedLocalRoleCounts, metrics.LocalRoleCounts)
	require.Equal(t, expectedReplicatedManagedKeys, metrics.ReplicatedManagedKeys)
	require.Equal(t, expectedLocalManagedKeys, metrics.LocalManagedKeys)
	require.Equal(t, 10, metrics.ReplicatedKvCounts)
	require.Equal(t, 10, metrics.LocalKvCounts)
}

// TestCountMetricsFromMounts_OfficialUnofficialMounts tests that CountMetricsFromMounts correctly
// distinguishes between official and unofficial mounts and returns the expected counts depending
// on if officialOnly is set.
func TestCountMetricsFromMounts_OfficialUnofficialMounts(t *testing.T) {
	pluginDir, err := filepath.EvalSymlinks(t.TempDir())
	require.NoError(t, err)
	coreConfig := &CoreConfig{
		PluginDirectory: pluginDir,
		LogicalBackends: map[string]logical.Factory{
			pluginconsts.SecretEngineAWS:        roleLogicalBackends[pluginconsts.SecretEngineAWS],
			pluginconsts.SecretEngineAzure:      roleLogicalBackends[pluginconsts.SecretEngineAzure],
			pluginconsts.SecretEngineKubernetes: roleLogicalBackends[pluginconsts.SecretEngineKubernetes],
		},
	}
	core, _, root := TestCoreUnsealedWithConfig(t, coreConfig)
	ctx := namespace.RootContext(context.Background())

	// Create official plugin mounts
	createMount(t, ctx, core, root, pluginconsts.SecretEngineAWS, "aws", false)
	createMount(t, ctx, core, root, pluginconsts.SecretEngineAzure, "azure", false)

	// Register an unofficial kubernetes plugin
	k8sPlugin := pluginhelpers.CompilePlugin(t, sdkconsts.PluginTypeSecrets, "v1.0.0", pluginDir)
	k8sPlugin.Name = "kubernetes"
	shaBytes, err := hex.DecodeString(k8sPlugin.Sha256)
	require.NoError(t, err)

	require.NoError(t, core.pluginCatalog.Set(ctx, sdkpluginutil.SetPluginInput{
		Name:    k8sPlugin.Name,
		Type:    k8sPlugin.Typ,
		Version: k8sPlugin.Version,
		Command: k8sPlugin.FileName,
		Sha256:  shaBytes,
	}))

	// Add unofficial mount
	mountEntry := &MountEntry{
		Table:          mountTableType,
		Path:           "kubernetes/",
		namespace:      &namespace.Namespace{Path: ""},
		Type:           k8sPlugin.Name,
		Version:        k8sPlugin.Version,
		RunningVersion: k8sPlugin.Version,
	}
	err = core.mount(ctx, mountEntry)
	require.NoError(t, err)

	// Sleep to prevent race conditions during the role initialization
	time.Sleep(1 * time.Second)

	core.mountsLock.RLock()
	defer core.mountsLock.RUnlock()

	// Add roles to mounts
	addRoleToStorage(t, core, "aws", "role/", 5)
	addRoleToStorage(t, core, "azure", "roles/", 5)
	addRoleToStorage(t, core, "kubernetes", "roles/", 5)

	// Get official-only counts, should not include kubernetes
	metrics, err := core.CountMetricsFromMounts(true)
	require.NoError(t, err)
	require.NotNil(t, metrics)
	require.Equal(t, 5, metrics.ReplicatedRoleCounts.AWSDynamicRoles)
	require.Equal(t, 5, metrics.ReplicatedRoleCounts.AzureDynamicRoles)
	require.Equal(t, 0, metrics.ReplicatedRoleCounts.KubernetesDynamicRoles)

	// Get counts including unofficial mounts
	metrics, err = core.CountMetricsFromMounts(false)
	require.NoError(t, err)
	require.NotNil(t, metrics)
	require.Equal(t, 5, metrics.ReplicatedRoleCounts.AWSDynamicRoles)
	require.Equal(t, 5, metrics.ReplicatedRoleCounts.AzureDynamicRoles)
	require.Equal(t, 5, metrics.ReplicatedRoleCounts.KubernetesDynamicRoles)
}

func createMount(t *testing.T, ctx context.Context, core *Core, token string, mountType string, mountPath string, local bool) {
	req := logical.TestRequest(t, logical.CreateOperation, fmt.Sprintf("sys/mounts/%s", mountPath))
	req.Data["type"] = mountType
	req.Data["local"] = local
	req.ClientToken = token
	resp, err := core.HandleRequest(ctx, req)
	require.NoError(t, err)
	require.Nil(t, resp.Error())
}
