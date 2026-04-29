// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package verify

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/vault/sdk/helper/testcluster/blackbox"
)

// TestVaultUndoLogsMetric verifies the vault.core.replication.write_undo_logs gauge metric
// This test runs from CI/GitHub runners and connects to the Vault cluster via API
func TestVaultUndoLogsMetric(t *testing.T) {
	t.Parallel()

	// Read required environment variables
	expectedStateStr := os.Getenv("EXPECTED_STATE")
	if expectedStateStr == "" {
		t.Fatal("EXPECTED_STATE environment variable is required")
	}

	expectedState, err := strconv.ParseFloat(expectedStateStr, 64)
	if err != nil {
		t.Fatalf("Failed to parse EXPECTED_STATE: %v", err)
	}

	// Validate expected state is 0 or 1
	if expectedState != 0 && expectedState != 1 {
		t.Fatalf("EXPECTED_STATE must be 0 or 1, got: %.0f", expectedState)
	}

	timeoutStr := os.Getenv("TIMEOUT_SECONDS")
	if timeoutStr == "" {
		t.Fatal("TIMEOUT_SECONDS environment variable is required")
	}

	timeoutSeconds, err := strconv.Atoi(timeoutStr)
	if err != nil {
		t.Fatalf("Failed to parse TIMEOUT_SECONDS: %v", err)
	}

	retryIntervalStr := os.Getenv("RETRY_INTERVAL")
	if retryIntervalStr == "" {
		t.Fatal("RETRY_INTERVAL environment variable is required")
	}

	retryIntervalSeconds, err := strconv.Atoi(retryIntervalStr)
	if err != nil {
		t.Fatalf("Failed to parse RETRY_INTERVAL: %v", err)
	}

	timeout := time.Duration(timeoutSeconds) * time.Second
	retryInterval := time.Duration(retryIntervalSeconds) * time.Second

	v := blackbox.New(t)

	// Verify the undo logs metric has the expected value
	v.AssertMetricGaugeValue("vault.core.replication.write_undo_logs", expectedState, timeout, retryInterval)
}
