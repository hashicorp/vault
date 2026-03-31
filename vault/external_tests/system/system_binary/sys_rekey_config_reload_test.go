// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package system_binary

import (
	"os"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/testcluster"
	"github.com/hashicorp/vault/sdk/helper/testcluster/docker"
	"github.com/stretchr/testify/require"
)

// waitForRekeyInConfig polls sys/config/state/sanitized until the rekey endpoint
// appears or disappears from enable_unauthenticated_access based on shouldBePresent.
func waitForRekeyInConfig(t *testing.T, client *api.Client, rootToken string, shouldBePresent bool) {
	clientWithAuth, err := client.Clone()
	require.NoError(t, err, "failed to clone client")
	clientWithAuth.SetToken(rootToken)

	require.Eventually(t, func() bool {
		resp, err := clientWithAuth.Logical().Read("sys/config/state/sanitized")
		if err != nil {
			t.Logf("error reading config state: %v", err)
			return false
		}
		if resp == nil || resp.Data == nil {
			t.Logf("nil response or data from config state")
			return false
		}

		override, ok := resp.Data["enable_unauthenticated_access"]
		if !ok {
			// If the field is not present, rekey is not in the override list
			return !shouldBePresent
		}

		// Check if override contains "rekey"
		rekeyFound := false
		if overrideSlice, ok := override.([]interface{}); ok {
			for _, v := range overrideSlice {
				if str, ok := v.(string); ok && str == "rekey" {
					rekeyFound = true
					break
				}
			}
		}

		if shouldBePresent {
			return rekeyFound
		}
		return !rekeyFound
	}, 10*time.Second, 100*time.Millisecond, "rekey presence in enable_unauthenticated_access did not match expected state")
}

// TestSysRekey_ConfigReload tests that the rekey status endpoint can be toggled
// between requiring authentication and not requiring authentication by using
// the enable_unauthenticated_access config option and reloading the config.
func TestSysRekey_ConfigReload(t *testing.T) {
	binary := os.Getenv("VAULT_BINARY")
	if binary == "" {
		t.Skip("only running docker test when $VAULT_BINARY present")
	}

	nodeConfig := &testcluster.VaultNodeConfig{
		LogLevel: "TRACE",
	}
	opts := &docker.DockerClusterOptions{
		ImageRepo:    "hashicorp/vault",
		ImageTag:     "latest",
		VaultBinary:  binary,
		DisableMlock: true,
		ClusterOptions: testcluster.ClusterOptions{
			NumCores:        1,
			VaultNodeConfig: nodeConfig,
		},
	}

	cluster := docker.NewTestDockerCluster(t, opts)
	defer cluster.Cleanup()

	node := cluster.Nodes()[0].(*docker.DockerClusterNode)
	client := node.APIClient()
	rootToken := cluster.GetRootToken()

	// Test 1: Without enable_unauthenticated_access, rekey status should require auth
	t.Run("requires-auth-by-default", func(t *testing.T) {
		// Try without token - should fail
		clientNoAuth, err := client.Clone()
		require.NoError(t, err, "failed to clone client")
		clientNoAuth.SetToken("")

		_, err = clientNoAuth.Logical().Read("sys/rekey/init")
		require.Error(t, err, "expected error when accessing rekey status without token")
		require.Contains(t, err.Error(), "permission denied", "error should indicate permission denied")

		// Try with token - should succeed
		clientWithAuth, err := client.Clone()
		require.NoError(t, err, "failed to clone client")
		clientWithAuth.SetToken(rootToken)

		resp, err := clientWithAuth.Logical().Read("sys/rekey/init")
		require.NoError(t, err, "should succeed with valid token")
		require.NotNil(t, resp, "response should not be nil")
		require.NotNil(t, resp.Data, "response data should not be nil")
		require.False(t, resp.Data["started"].(bool), "rekey should not be started")
	})

	// Test 2: Update config to enable unauthenticated rekey and reload
	t.Run("enable-unauthenticated-via-config-reload", func(t *testing.T) {
		// Create updated config with enable_unauthenticated_access
		nodeConfig.EnableUnauthenticatedAccess = []string{"rekey"}

		// Update the config and copy it to the container
		err := node.UpdateConfig(t.Context(), nodeConfig)
		require.NoError(t, err, "failed to update config")

		// Send SIGHUP to reload the configuration
		err = node.Signal(t.Context(), "SIGHUP")
		require.NoError(t, err, "failed to send SIGHUP")

		// Wait for rekey to appear in enable_unauthenticated_access
		waitForRekeyInConfig(t, client, rootToken, true)

		// Now test that rekey status works without auth
		clientNoAuth, err := client.Clone()
		require.NoError(t, err, "failed to clone client")
		clientNoAuth.SetToken("")

		resp, err := clientNoAuth.Logical().Read("sys/rekey/init")
		require.NoError(t, err, "should succeed without token after config reload")
		require.NotNil(t, resp, "response should not be nil")
		require.NotNil(t, resp.Data, "response data should not be nil")
		require.False(t, resp.Data["started"].(bool), "rekey should not be started")

		// Verify it still works with auth
		clientWithAuth2, err := client.Clone()
		require.NoError(t, err, "failed to clone client")
		clientWithAuth2.SetToken(rootToken)

		resp2, err := clientWithAuth2.Logical().Read("sys/rekey/init")
		require.NoError(t, err, "should still succeed with valid token")
		require.NotNil(t, resp2, "response should not be nil")
		require.NotNil(t, resp2.Data, "response data should not be nil")
		require.False(t, resp2.Data["started"].(bool), "rekey should not be started")
	})

	// Test 3: Remove the override and reload to restore auth requirement
	t.Run("restore-auth-requirement-via-config-reload", func(t *testing.T) {
		// Create config without enable_unauthenticated_access
		nodeConfig.EnableUnauthenticatedAccess = nil

		// Update the config and copy it to the container
		err := node.UpdateConfig(t.Context(), nodeConfig)
		require.NoError(t, err, "failed to update config")

		// Send SIGHUP to reload the configuration
		err = node.Signal(t.Context(), "SIGHUP")
		require.NoError(t, err, "failed to send SIGHUP")

		// Wait for rekey to be removed from enable_unauthenticated_access
		waitForRekeyInConfig(t, client, rootToken, false)

		// Now test that rekey status requires auth again
		clientNoAuth, err := client.Clone()
		require.NoError(t, err, "failed to clone client")
		clientNoAuth.SetToken("")

		_, err = clientNoAuth.Logical().Read("sys/rekey/init")
		require.Error(t, err, "should fail without token after restoring auth requirement")
		require.Contains(t, err.Error(), "permission denied", "error should indicate permission denied")

		// Verify it still works with auth
		clientWithAuth2, err := client.Clone()
		require.NoError(t, err, "failed to clone client")
		clientWithAuth2.SetToken(rootToken)

		resp, err := clientWithAuth2.Logical().Read("sys/rekey/init")
		require.NoError(t, err, "should succeed with valid token")
		require.NotNil(t, resp, "response should not be nil")
		require.NotNil(t, resp.Data, "response data should not be nil")
		require.False(t, resp.Data["started"].(bool), "rekey should not be started")
	})
}
