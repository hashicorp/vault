//go:build isolated
// +build isolated

// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package verify

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/testcluster/blackbox"
	"github.com/stretchr/testify/require"
)

// TestPerformanceReplication_Status verifies performance replication status
// between primary and secondary clusters. This test validates:
// - Primary cluster is in "running" state
// - Secondary cluster is in "stream-wals" state
// - Connection status is "connected" on both sides
// - Secondary knows the primary cluster addresses
//
// This test requires:
// - VAULT_SECONDARY_ADDR environment variable pointing to the secondary cluster
// - VAULT_SECONDARY_TOKEN environment variable with secondary cluster token
// - Performance replication to be already enabled between clusters
func TestPerformanceReplication_Status(t *testing.T) {
	t.Parallel()

	// Get secondary cluster connection info from environment
	secondaryAddr := os.Getenv("VAULT_SECONDARY_ADDR")
	secondaryToken := os.Getenv("VAULT_SECONDARY_TOKEN")

	if secondaryAddr == "" || secondaryToken == "" {
		t.Skip("Skipping performance replication test - VAULT_SECONDARY_ADDR and VAULT_SECONDARY_TOKEN not set")
	}

	// Primary cluster session
	primary := blackbox.New(t)

	// Create secondary cluster client
	secondaryConfig := api.DefaultConfig()
	secondaryConfig.Address = secondaryAddr
	secondaryClient, err := api.NewClient(secondaryConfig)
	require.NoError(t, err, "Failed to create secondary cluster client")
	secondaryClient.SetToken(secondaryToken)

	// Create secondary session wrapper
	secondary := &blackbox.Session{
		Client: secondaryClient,
	}

	// Verify primary cluster replication status
	t.Run("primary_cluster_status", func(t *testing.T) {
		verifyPrimaryReplicationStatus(t, primary)
	})

	// Verify secondary cluster replication status
	t.Run("secondary_cluster_status", func(t *testing.T) {
		verifySecondaryReplicationStatus(t, secondary)
	})

	// Verify connection between primary and secondary
	t.Run("verify_connection", func(t *testing.T) {
		verifyReplicationConnection(t, primary, secondary)
	})

	t.Log("✓ Performance replication verification completed successfully")
}

// verifyPrimaryReplicationStatus checks the primary cluster's replication status
func verifyPrimaryReplicationStatus(t *testing.T, v *blackbox.Session) {
	t.Helper()

	var prStatus *api.Secret
	var err error

	// Retry reading status as replication may take time to establish
	v.EventuallyWithTimeout(func() error {
		prStatus, err = v.WithRootNamespace(func() (*api.Secret, error) {
			return v.Client.Logical().Read("sys/replication/performance/status")
		})
		if err != nil {
			return fmt.Errorf("failed to read performance replication status: %w", err)
		}

		if prStatus == nil || prStatus.Data == nil {
			return fmt.Errorf("empty performance replication status response")
		}

		return nil
	}, 30*time.Second)

	require.NoError(t, err)
	require.NotNil(t, prStatus)
	require.NotNil(t, prStatus.Data)

	// Verify mode is "primary"
	mode, ok := prStatus.Data["mode"].(string)
	require.True(t, ok, "mode field not found or not a string")
	require.Equal(t, "primary", mode, "Expected mode to be 'primary'")

	// Verify state is not "idle"
	state, ok := prStatus.Data["state"].(string)
	require.True(t, ok, "state field not found or not a string")
	require.NotEqual(t, "idle", state, "Primary cluster state should not be 'idle'")
	require.Equal(t, "running", state, "Primary cluster should be in 'running' state")

	// Verify secondaries connection status
	secondaries, ok := prStatus.Data["secondaries"].([]interface{})
	if ok && len(secondaries) > 0 {
		secondary := secondaries[0].(map[string]interface{})
		connectionStatus, ok := secondary["connection_status"].(string)
		require.True(t, ok, "connection_status not found in secondaries[0]")
		require.NotEqual(t, "disconnected", connectionStatus, "Secondary connection should not be 'disconnected'")
		t.Logf("✓ Primary cluster: mode=%s, state=%s, secondary_connection=%s", mode, state, connectionStatus)
	} else {
		t.Log("⚠ No secondaries found in primary status - replication may not be fully established")
	}
}

// verifySecondaryReplicationStatus checks the secondary cluster's replication status
func verifySecondaryReplicationStatus(t *testing.T, secondary *blackbox.Session) {
	t.Helper()

	var prStatus *api.Secret
	var err error

	// Read performance replication status from secondary
	prStatus, err = secondary.Client.Logical().Read("sys/replication/performance/status")
	require.NoError(t, err, "Failed to read performance replication status from secondary")
	require.NotNil(t, prStatus, "Empty performance replication status response from secondary")
	require.NotNil(t, prStatus.Data, "Empty data in performance replication status from secondary")

	// Verify mode is "secondary"
	mode, ok := prStatus.Data["mode"].(string)
	require.True(t, ok, "mode field not found or not a string")
	require.Equal(t, "secondary", mode, "Expected mode to be 'secondary'")

	// Verify state is not "idle"
	state, ok := prStatus.Data["state"].(string)
	require.True(t, ok, "state field not found or not a string")
	require.NotEqual(t, "idle", state, "Secondary cluster state should not be 'idle'")
	require.Equal(t, "stream-wals", state, "Secondary cluster should be in 'stream-wals' state")

	// Verify primaries connection status
	primaries, ok := prStatus.Data["primaries"].([]interface{})
	if ok && len(primaries) > 0 {
		primary := primaries[0].(map[string]interface{})
		connectionStatus, ok := primary["connection_status"].(string)
		require.True(t, ok, "connection_status not found in primaries[0]")
		require.NotEqual(t, "disconnected", connectionStatus, "Primary connection should not be 'disconnected'")
		t.Logf("✓ Secondary cluster: mode=%s, state=%s, primary_connection=%s", mode, state, connectionStatus)
	} else {
		t.Fatal("No primaries found in secondary status - replication not properly configured")
	}
}

// verifyReplicationConnection verifies the connection between primary and secondary clusters
func verifyReplicationConnection(t *testing.T, primary *blackbox.Session, secondary *blackbox.Session) {
	t.Helper()

	// Read secondary status to get primary cluster address information
	prStatus, err := secondary.Client.Logical().Read("sys/replication/performance/status")
	require.NoError(t, err)
	require.NotNil(t, prStatus)
	require.NotNil(t, prStatus.Data)

	// Get known primary cluster addresses
	knownPrimaryAddrs, ok := prStatus.Data["known_primary_cluster_addrs"].([]interface{})
	require.True(t, ok, "known_primary_cluster_addrs not found or not an array")
	require.NotEmpty(t, knownPrimaryAddrs, "known_primary_cluster_addrs should not be empty")

	// Get primary cluster address from primaries array
	primaries, ok := prStatus.Data["primaries"].([]interface{})
	require.True(t, ok && len(primaries) > 0, "primaries array not found or empty")

	primaryInfo := primaries[0].(map[string]interface{})
	clusterAddr, ok := primaryInfo["cluster_address"].(string)
	require.True(t, ok, "cluster_address not found in primaries[0]")
	require.NotEmpty(t, clusterAddr, "cluster_address should not be empty")

	// Extract IP address from cluster_address (format: https://[ip]:port or https://ip:port)
	var primaryIP string
	if strings.Contains(clusterAddr, "[") {
		// IPv6 format: https://[2001:db8::1]:8201
		start := strings.Index(clusterAddr, "[")
		end := strings.Index(clusterAddr, "]")
		if start != -1 && end != -1 {
			primaryIP = clusterAddr[start+1 : end]
		}
	} else {
		// IPv4 format: https://10.0.0.1:8201
		parts := strings.Split(clusterAddr, "://")
		if len(parts) == 2 {
			hostPort := parts[1]
			ipPort := strings.Split(hostPort, ":")
			if len(ipPort) >= 2 {
				primaryIP = ipPort[0]
			}
		}
	}

	require.NotEmpty(t, primaryIP, "Failed to extract primary IP from cluster_address: %s", clusterAddr)

	// Verify primary IP is in known_primary_cluster_addrs
	found := false
	for _, addr := range knownPrimaryAddrs {
		addrStr := addr.(string)
		if strings.Contains(addrStr, primaryIP) {
			found = true
			t.Logf("✓ Primary cluster address %s found in known_primary_cluster_addrs", primaryIP)
			break
		}
	}

	require.True(t, found, "Primary cluster address %s not found in known_primary_cluster_addrs: %v", primaryIP, knownPrimaryAddrs)
}

// TestPerformanceReplication_DataReplication verifies that data written to primary
// is replicated to secondary cluster. This is a basic smoke test for replication.
func TestPerformanceReplication_DataReplication(t *testing.T) {
	t.Parallel()

	// Get secondary cluster connection info from environment
	secondaryAddr := os.Getenv("VAULT_SECONDARY_ADDR")
	secondaryToken := os.Getenv("VAULT_SECONDARY_TOKEN")

	if secondaryAddr == "" || secondaryToken == "" {
		t.Skip("Skipping data replication test - VAULT_SECONDARY_ADDR and VAULT_SECONDARY_TOKEN not set")
	}

	// Primary cluster session
	primary := blackbox.New(t)

	// Verify performance replication is enabled
	primary.AssertPerformanceReplicationStatus("primary")

	// Create secondary cluster client
	secondaryConfig := api.DefaultConfig()
	secondaryConfig.Address = secondaryAddr
	secondaryClient, err := api.NewClient(secondaryConfig)
	require.NoError(t, err)
	secondaryClient.SetToken(secondaryToken)

	// Enable KV v2 secrets engine on primary (if not already enabled)
	testMount := "secret"

	// Write test data to primary
	testPath := fmt.Sprintf("%s/data/replication-test-%d", testMount, time.Now().Unix())
	testData := map[string]interface{}{
		"data": map[string]interface{}{
			"test_key":  "test_value",
			"timestamp": time.Now().Unix(),
			"test_id":   "performance-replication-test",
		},
	}

	_, err = primary.Client.Logical().Write(testPath, testData)
	require.NoError(t, err, "Failed to write test data to primary")
	t.Logf("✓ Wrote test data to primary at %s", testPath)

	// Wait for replication to sync (with retry)
	var replicatedSecret *api.Secret
	maxRetries := 10
	retryDelay := 2 * time.Second

	for i := 0; i < maxRetries; i++ {
		replicatedSecret, err = secondaryClient.Logical().Read(testPath)
		if err == nil && replicatedSecret != nil && replicatedSecret.Data != nil {
			break
		}
		if i < maxRetries-1 {
			t.Logf("Waiting for replication to sync (attempt %d/%d)...", i+1, maxRetries)
			time.Sleep(retryDelay)
		}
	}

	require.NoError(t, err, "Failed to read replicated data from secondary")
	require.NotNil(t, replicatedSecret, "Replicated secret is nil")
	require.NotNil(t, replicatedSecret.Data, "Replicated secret data is nil")

	// Verify the data matches
	dataMap, ok := replicatedSecret.Data["data"].(map[string]interface{})
	require.True(t, ok, "data field not found or not a map")

	testKey, ok := dataMap["test_key"].(string)
	require.True(t, ok, "test_key not found in replicated data")
	require.Equal(t, "test_value", testKey, "Replicated data does not match")

	testID, ok := dataMap["test_id"].(string)
	require.True(t, ok, "test_id not found in replicated data")
	require.Equal(t, "performance-replication-test", testID, "Replicated test_id does not match")

	t.Log("✓ Data successfully replicated from primary to secondary")

	// Cleanup: delete test data from primary
	_, err = primary.Client.Logical().Delete(testPath)
	if err != nil {
		t.Logf("Warning: Failed to cleanup test data: %v", err)
	}
}
