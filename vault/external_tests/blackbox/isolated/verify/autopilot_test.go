//go:build isolated
// +build isolated

// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package verify

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/vault/sdk/helper/testcluster/blackbox"
	"github.com/stretchr/testify/require"
)

// TestAutopilotUpgradeStatus verifies that the autopilot upgrade status and target version
// match expected values. This test polls the autopilot state endpoint with retries.
func TestAutopilotUpgradeStatus(t *testing.T) {
	expectedStatus := os.Getenv("VAULT_AUTOPILOT_UPGRADE_STATUS")
	require.NotEmpty(t, expectedStatus, "VAULT_AUTOPILOT_UPGRADE_STATUS must be set")

	expectedVersion := os.Getenv("VAULT_AUTOPILOT_UPGRADE_VERSION")
	require.NotEmpty(t, expectedVersion, "VAULT_AUTOPILOT_UPGRADE_VERSION must be set")

	timeoutStr := os.Getenv("TIMEOUT_SECONDS")
	if timeoutStr == "" {
		timeoutStr = "180" // Default timeout
	}
	timeoutSeconds, err := strconv.Atoi(timeoutStr)
	require.NoError(t, err, "failed to parse TIMEOUT_SECONDS")

	retryIntervalStr := os.Getenv("RETRY_INTERVAL")
	if retryIntervalStr == "" {
		retryIntervalStr = "5" // Default retry interval
	}
	retryIntervalSeconds, err := strconv.Atoi(retryIntervalStr)
	require.NoError(t, err, "failed to parse RETRY_INTERVAL")

	timeout := time.Duration(timeoutSeconds) * time.Second
	retryInterval := time.Duration(retryIntervalSeconds) * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	session := blackbox.New(t, blackbox.WithoutNamespace())

	// Debug: Log connection details
	t.Logf("DEBUG: VAULT_ADDR from env: %s", os.Getenv("VAULT_ADDR"))
	t.Logf("DEBUG: Client address: %s", session.Client.Address())
	t.Logf("DEBUG: VAULT_TOKEN set: %v", os.Getenv("VAULT_TOKEN") != "")
	t.Logf("DEBUG: VAULT_NAMESPACE from env: %s", os.Getenv("VAULT_NAMESPACE"))

	// Debug: Verify token is valid before attempting autopilot API call
	t.Logf("DEBUG: Attempting to verify token validity...")
	tokenLookup, err := session.Client.Auth().Token().LookupSelf()
	if err != nil {
		t.Fatalf("FATAL: Token lookup failed - token is invalid or expired: %v", err)
	}
	t.Logf("DEBUG: Token is valid. Policies: %v, TTL: %v, Renewable: %v",
		tokenLookup.Data["policies"],
		tokenLookup.Data["ttl"],
		tokenLookup.Data["renewable"])

	ticker := time.NewTicker(retryInterval)
	defer ticker.Stop()

	var lastStatus, lastVersion string
	for {
		select {
		case <-ctx.Done():
			t.Fatalf("Timeout waiting for autopilot status. Expected status=%s (got %s), expected version=%s (got %s)",
				expectedStatus, lastStatus, expectedVersion, lastVersion)
		case <-ticker.C:
			resp, err := session.Client.Logical().Read("sys/storage/raft/autopilot/state")
			require.NoError(t, err, "failed to read autopilot state")
			require.NotNil(t, resp, "autopilot state response is nil")
			require.NotNil(t, resp.Data, "autopilot state data is nil")

			upgradeInfo, ok := resp.Data["upgrade_info"].(map[string]interface{})
			if !ok {
				t.Logf("upgrade_info not found or invalid type in response")
				continue
			}

			status, _ := upgradeInfo["status"].(string)
			targetVersion, _ := upgradeInfo["target_version"].(string)

			lastStatus = status
			lastVersion = targetVersion

			if status == expectedStatus && targetVersion == expectedVersion {
				t.Logf("Autopilot status verified: status=%s, target_version=%s", status, targetVersion)
				return
			}

			t.Logf("Waiting for autopilot status. Current: status=%s (expected %s), target_version=%s (expected %s)",
				status, expectedStatus, targetVersion, expectedVersion)
		}
	}
}

// TestAutopilotUpgradeStatusOutput verifies autopilot upgrade status and outputs the full state as JSON.
// This test is useful for debugging and provides detailed autopilot state information.
func TestAutopilotUpgradeStatusOutput(t *testing.T) {
	expectedStatus := os.Getenv("VAULT_AUTOPILOT_UPGRADE_STATUS")
	require.NotEmpty(t, expectedStatus, "VAULT_AUTOPILOT_UPGRADE_STATUS must be set")

	expectedVersion := os.Getenv("VAULT_AUTOPILOT_UPGRADE_VERSION")
	require.NotEmpty(t, expectedVersion, "VAULT_AUTOPILOT_UPGRADE_VERSION must be set")

	timeoutStr := os.Getenv("TIMEOUT_SECONDS")
	if timeoutStr == "" {
		timeoutStr = "180" // Default timeout
	}
	timeoutSeconds, err := strconv.Atoi(timeoutStr)
	require.NoError(t, err, "failed to parse TIMEOUT_SECONDS")

	retryIntervalStr := os.Getenv("RETRY_INTERVAL")
	if retryIntervalStr == "" {
		retryIntervalStr = "5" // Default retry interval
	}
	retryIntervalSeconds, err := strconv.Atoi(retryIntervalStr)
	require.NoError(t, err, "failed to parse RETRY_INTERVAL")

	timeout := time.Duration(timeoutSeconds) * time.Second
	retryInterval := time.Duration(retryIntervalSeconds) * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	session := blackbox.New(t, blackbox.WithoutNamespace())

	// Debug: Log connection details
	t.Logf("DEBUG: VAULT_ADDR from env: %s", os.Getenv("VAULT_ADDR"))
	t.Logf("DEBUG: Client address: %s", session.Client.Address())
	t.Logf("DEBUG: VAULT_TOKEN set: %v", os.Getenv("VAULT_TOKEN") != "")
	t.Logf("DEBUG: VAULT_NAMESPACE from env: %s", os.Getenv("VAULT_NAMESPACE"))

	// Debug: Verify token is valid before attempting autopilot API call
	t.Logf("DEBUG: Attempting to verify token validity...")
	tokenLookup, err := session.Client.Auth().Token().LookupSelf()
	if err != nil {
		t.Fatalf("FATAL: Token lookup failed - token is invalid or expired: %v", err)
	}
	t.Logf("DEBUG: Token is valid. Policies: %v, TTL: %v, Renewable: %v",
		tokenLookup.Data["policies"],
		tokenLookup.Data["ttl"],
		tokenLookup.Data["renewable"])

	ticker := time.NewTicker(retryInterval)
	defer ticker.Stop()

	var lastResp *map[string]interface{}
	for {
		select {
		case <-ctx.Done():
			if lastResp != nil {
				stateJSON, _ := json.MarshalIndent(lastResp, "", "  ")
				t.Logf("Final autopilot state:\n%s", string(stateJSON))
			}
			t.Fatalf("Timeout waiting for autopilot status. Expected status=%s, expected version=%s",
				expectedStatus, expectedVersion)
		case <-ticker.C:
			resp, err := session.Client.Logical().Read("sys/storage/raft/autopilot/state")
			require.NoError(t, err, "failed to read autopilot state")
			require.NotNil(t, resp, "autopilot state response is nil")
			require.NotNil(t, resp.Data, "autopilot state data is nil")

			lastResp = &resp.Data

			upgradeInfo, ok := resp.Data["upgrade_info"].(map[string]interface{})
			if !ok {
				t.Logf("upgrade_info not found or invalid type in response")
				continue
			}

			status, _ := upgradeInfo["status"].(string)
			targetVersion, _ := upgradeInfo["target_version"].(string)

			if status == expectedStatus && targetVersion == expectedVersion {
				stateJSON, err := json.MarshalIndent(resp.Data, "", "  ")
				require.NoError(t, err, "failed to marshal autopilot state")
				fmt.Printf("Autopilot state verified:\n%s\n", string(stateJSON))
				return
			}

			t.Logf("Waiting for autopilot status. Current: status=%s (expected %s), target_version=%s (expected %s)",
				status, expectedStatus, targetVersion, expectedVersion)
		}
	}
}
