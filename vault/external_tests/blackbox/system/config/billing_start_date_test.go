//go:build system
// +build system

// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package config

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/hashicorp/vault/sdk/helper/testcluster/blackbox"
	"github.com/stretchr/testify/require"
)

// TestBillingStartDate verifies that the billing start date has successfully
// rolled over to the latest billing year if needed.
//
// This test replicates the behavior of the enos module:
// enos/modules/vault_verify_billing_start_date/
//
// The test validates that:
// 1. The billing start timestamp exists in sys/internal/counters/config
// 2. The timestamp is within the last 12 months (current billing year)
// 3. The timestamp is not in the future
//
// This is used in the upgrade scenario to ensure billing dates properly
// roll over after cluster upgrades.
func TestBillingStartDate(t *testing.T) {
	t.Parallel()
	v := blackbox.New(t)

	// Verify cluster is healthy and unsealed
	v.AssertUnsealedAny()

	// Read the billing configuration
	secret, err := v.Client.Logical().Read("sys/internal/counters/config")
	require.NoError(t, err, "failed to read sys/internal/counters/config")
	require.NotNil(t, secret, "expected response from sys/internal/counters/config")
	require.NotNil(t, secret.Data, "expected data in response")

	// Extract billing_start_timestamp
	billingStartRaw, ok := secret.Data["billing_start_timestamp"]
	require.True(t, ok, "billing_start_timestamp not found in response")

	// Parse the timestamp
	var billingStart time.Time
	switch v := billingStartRaw.(type) {
	case string:
		var parseErr error
		billingStart, parseErr = time.Parse(time.RFC3339, v)
		require.NoError(t, parseErr, "failed to parse billing_start_timestamp as RFC3339 string")
	case time.Time:
		billingStart = v
	default:
		// Try unmarshaling through JSON as fallback
		bytes, marshalErr := json.Marshal(secret.Data)
		require.NoError(t, marshalErr, "failed to marshal response data")

		var configResp struct {
			BillingStartTimestamp time.Time `json:"billing_start_timestamp"`
		}
		unmarshalErr := json.Unmarshal(bytes, &configResp)
		require.NoError(t, unmarshalErr, "failed to unmarshal response data")
		billingStart = configResp.BillingStartTimestamp
	}

	// Calculate one year ago from now (this is the cutoff for the current billing year)
	oneYearAgo := time.Now().AddDate(-1, 0, 0)

	// Verify the billing start date is within the current billing year
	// (not more than 1 year old)
	require.False(t, billingStart.Before(oneYearAgo),
		"billing start date %s is not in the current billing year (more than 1 year old, cutoff: %s)",
		billingStart.Format(time.RFC3339),
		oneYearAgo.Format(time.RFC3339))

	// Verify the billing start date is not in the future
	require.False(t, billingStart.After(time.Now()),
		"billing start date %s is in the future",
		billingStart.Format(time.RFC3339))

	t.Logf("✓ Billing start date %s is within the current billing year", billingStart.Format(time.RFC3339))
}

// TestBillingStartDateRollover verifies that the billing start date is within
// the current billing year, confirming that automatic rollover has occurred if needed.
//
// Note: This test verifies the current state of a running cluster. The actual
// rollover behavior happens during cluster startup. In the upgrade scenario,
// this verification is done via the shell script module vault_verify_billing_start_date
// which has version-specific skip conditions (<=1.16.6 or 1.17.0-1.17.2).
//
// This blackbox test runs in the cloud-ent scenario without version restrictions,
// as cloud environments always run supported versions.
func TestBillingStartDateRollover(t *testing.T) {
	t.Parallel()
	v := blackbox.New(t)

	// Verify cluster is healthy and unsealed
	v.AssertUnsealedAny()

	// Read the billing configuration
	secret, err := v.Client.Logical().Read("sys/internal/counters/config")
	require.NoError(t, err, "failed to read sys/internal/counters/config")
	require.NotNil(t, secret, "expected response from sys/internal/counters/config")

	// Extract and validate billing_start_timestamp
	billingStartRaw, ok := secret.Data["billing_start_timestamp"]
	require.True(t, ok, "billing_start_timestamp not found in response")

	// Parse the timestamp
	var billingStart time.Time
	switch v := billingStartRaw.(type) {
	case string:
		billingStart, err = time.Parse(time.RFC3339, v)
		require.NoError(t, err, "failed to parse billing_start_timestamp")
	default:
		bytes, _ := json.Marshal(secret.Data)
		var configResp struct {
			BillingStartTimestamp time.Time `json:"billing_start_timestamp"`
		}
		json.Unmarshal(bytes, &configResp)
		billingStart = configResp.BillingStartTimestamp
	}

	// Verify the billing start date is within the current billing year
	oneYearAgo := time.Now().AddDate(-1, 0, 0)
	require.False(t, billingStart.Before(oneYearAgo),
		"billing start date %s should be within the current billing year (cutoff: %s)",
		billingStart.Format(time.RFC3339),
		oneYearAgo.Format(time.RFC3339))

	t.Logf("✓ Billing start date %s is within the current billing year",
		billingStart.Format(time.RFC3339))
}
