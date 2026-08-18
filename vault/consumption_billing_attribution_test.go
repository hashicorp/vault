// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/timeutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault/billing"
	"github.com/stretchr/testify/require"
)

// TestToFloat64 verifies that toFloat64 correctly handles all value types that
// can appear in a round-tripped MountAttribution.Count / MetricTypeAttribution.Count.
func TestToFloat64(t *testing.T) {
	// Native float64 — set by in-memory code paths.
	require.Equal(t, 3.14, toFloat64(float64(3.14)))
	require.Equal(t, 0.0, toFloat64(float64(0)))

	// nil — should return 0 safely.
	require.Equal(t, 0.0, toFloat64(nil))

	// json.Number — returned by jsonutil.DecodeJSON for stored numeric values.
	// Verify it is correctly unwrapped via the Float64() interface.
	require.InDelta(t, 2.5, toFloat64(json.Number("2.5")), 0.0001, "Float64()-capable type should be unwrapped")

	// Integer types are coerced to float64.
	require.Equal(t, 5.0, toFloat64(int(5)), "int should coerce to float64")

	// Unsupported types return 0.
	require.Equal(t, 0.0, toFloat64("3.14"), "string should return 0")
}

// TestStoreAndGetAttributionData verifies the round-trip of storeAttributionDataLocked
// and getStoredAttributionDataLocked, and that a second store overwrites the previous
// entry (no implicit merge — callers are responsible for merging before storing).
func TestStoreAndGetAttributionData(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	view := &logical.InmemStorage{}
	now := time.Now().UTC()
	month := timeutil.StartOfMonth(now)
	lastUpdated := time.Date(2026, 5, 14, 18, 7, 23, 0, time.UTC)

	data := &logical.MetricTypeAttribution{
		Count:       10,
		LastUpdated: lastUpdated,
		Mounts: map[string]logical.MountAttribution{
			"kv_5d4f8f1c": {
				MountPath:         "secret/",
				MountType:         "kv",
				MountAccessor:     "kv_5d4f8f1c",
				NamespaceID:       "root",
				NamespacePath:     "",
				ParentNamespaceID: "",
				Count:             5,
				BackendAwareUUID:  "wdasd23",
			},
		},
	}

	err := storeAttributionDataLocked(ctx, view, billing.LocalPrefix, month, billing.KvHWMCountsHWM, data)
	require.NoError(t, err)

	got, err := getStoredAttributionDataLocked(ctx, view, billing.LocalPrefix, month, billing.KvHWMCountsHWM)
	require.NoError(t, err)
	require.NotNil(t, got)

	// Count is interface{} — jsonutil.DecodeJSON deserialises numbers as json.Number.
	// Compare via fmt.Sprintf to avoid type mismatch between int and json.Number.
	require.Equal(t, "10", fmt.Sprintf("%v", got.Count))
	require.Equal(t, data.LastUpdated.UTC(), got.LastUpdated.UTC())
	require.Len(t, got.Mounts, 1)

	m1 := got.Mounts["kv_5d4f8f1c"]
	require.Equal(t, "secret/", m1.MountPath)
	require.Equal(t, "kv", m1.MountType)
	require.Equal(t, "root", m1.NamespaceID)
	require.Equal(t, "", m1.NamespacePath)
	require.Equal(t, "", m1.ParentNamespaceID)
	require.Equal(t, "kv_5d4f8f1c", m1.MountAccessor)
	require.Equal(t, "5", fmt.Sprintf("%v", m1.Count))

	// Overwrite with new data — second store must replace, not merge.
	overwrite := &logical.MetricTypeAttribution{
		Count:       12,
		LastUpdated: time.Now().UTC(),
		Mounts: map[string]logical.MountAttribution{
			"kv_bbb": {Count: 12, MountAccessor: "kv_bbb", MountPath: "new/", MountType: "kv"},
		},
	}
	err = storeAttributionDataLocked(ctx, view, billing.LocalPrefix, month, billing.KvHWMCountsHWM, overwrite)
	require.NoError(t, err)

	got, err = getStoredAttributionDataLocked(ctx, view, billing.LocalPrefix, month, billing.KvHWMCountsHWM)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, "12", fmt.Sprintf("%v", got.Count))
	require.Len(t, got.Mounts, 1, "overwrite should replace all previous mounts")
	_, hasOld := got.Mounts["kv_5d4f8f1c"]
	require.False(t, hasOld, "old mounts should be gone after overwrite")
	_, hasNew := got.Mounts["kv_bbb"]
	require.True(t, hasNew, "new mount should be present after overwrite")
}

// TestTransitUpdateAndGetAttribution verifies that the store operation correctly writes the cumulative counts
// of transit operations for the current month by. This also verifies that the retrieve operations correctly
// returns a map by mount that contains the correct cumulative transit operation counts for each month.
func TestTransitUpdateAndGetAttribution(t *testing.T) {
	core, _, _ := TestCoreUnsealedWithConfig(t, &CoreConfig{})
	ctx := namespace.RootContext(context.Background())
	currentMonth := time.Now()

	mountAccessor := "transit-accessor"
	mountAccessor2 := "transit-accessor-2"

	// Get consumption billing reference
	core.consumptionBillingLock.RLock()
	cb := core.consumptionBilling
	core.consumptionBillingLock.RUnlock()
	require.NotNil(t, cb)

	// Test Case 1: Simple update with a couple of mounts
	t.Log("Test Case 1: Initial update with two mounts")
	testBreakdown1 := logical.MountAttribution{
		MountAccessor:     mountAccessor,
		MountPath:         "transit/",
		NamespaceID:       "root",
		NamespacePath:     "",
		ParentNamespaceID: "",
		Count:             10.0,
	}

	testBreakdown2 := logical.MountAttribution{
		MountAccessor:     mountAccessor2,
		MountPath:         "transit2/",
		NamespaceID:       "ns1-id",
		NamespacePath:     "ns1/",
		ParentNamespaceID: "root",
		Count:             25.5,
	}

	// Store the breakdowns in the map
	cb.SecretEngineCounts.Transit.MountAttributionLock.Lock()
	cb.SecretEngineCounts.Transit.MountAttribution[mountAccessor] = testBreakdown1
	cb.SecretEngineCounts.Transit.MountAttribution[mountAccessor2] = testBreakdown2
	cb.SecretEngineCounts.Transit.MountAttributionLock.Unlock()

	err := core.UpdateTransitAttribution(ctx, currentMonth)
	require.NoError(t, err, "First UpdateTransitAttribution should succeed")

	retrievedAttribution, err := core.GetStoredAttributionData(ctx, billing.LocalPrefix, currentMonth, billing.TransitDataProtectionCallCountsPrefix)
	require.NoError(t, err)
	require.NotNil(t, retrievedAttribution)
	require.Len(t, retrievedAttribution.Mounts, 2, "Should have 2 mounts after first update")
	// Top-level count must equal the sum of all per-mount counts (10 + 25.5 = 35.5)
	require.Equal(t, "35.5", fmt.Sprintf("%v", retrievedAttribution.Count), "Top-level count should be sum of all mount counts")

	retrieved1, ok := retrievedAttribution.Mounts[mountAccessor]
	require.True(t, ok, "Should find breakdown for first mount accessor")
	verifyMountAttributionBreakdowns(t, testBreakdown1, retrieved1)

	retrieved2, ok := retrievedAttribution.Mounts[mountAccessor2]
	require.True(t, ok, "Should find breakdown for transit-accessor-2")
	verifyMountAttributionBreakdowns(t, testBreakdown2, retrieved2)

	// Test Case 2: Update with no mounts (empty map) - should keep existing counts
	t.Log("Test Case 2: Update with no mounts (empty map)")
	// TransitAttributions was already cleared by UpdateTransitAttribution above; nothing to add.

	err = core.UpdateTransitAttribution(ctx, currentMonth)
	require.NoError(t, err, "UpdateTransitAttribution with empty map should succeed")

	// Retrieve and verify - stored counts must be unchanged since nothing was flushed
	retrievedAttribution, err = core.GetStoredAttributionData(ctx, billing.LocalPrefix, currentMonth, billing.TransitDataProtectionCallCountsPrefix)
	require.NoError(t, err, "GetTransitAttribution after empty update should succeed")
	require.NotNil(t, retrievedAttribution, "Retrieved attribution should not be nil")
	require.Len(t, retrievedAttribution.Mounts, 2, "Should still have 2 mounts (counts preserved)")

	retrieved1, ok = retrievedAttribution.Mounts[mountAccessor]
	require.True(t, ok, "Should still find breakdown for first mount accessor")
	require.Equal(t, "10", fmt.Sprintf("%v", retrieved1.Count), "Count should remain unchanged")

	retrieved2, ok = retrievedAttribution.Mounts[mountAccessor2]
	require.True(t, ok, "Should still find breakdown for transit-accessor-2")
	require.Equal(t, "25.5", fmt.Sprintf("%v", retrieved2.Count), "Count should remain unchanged")

	// Test Case 3: Update with only one of the mounts - should accumulate counts correctly
	t.Log("Test Case 3: Update with only one mount (cumulative count)")
	testBreakdown1Updated := logical.MountAttribution{
		MountAccessor:     mountAccessor,
		MountPath:         "transit/",
		NamespaceID:       "root",
		NamespacePath:     "",
		ParentNamespaceID: "",
		Count:             15.0, // Adding 15 more to the existing 10
	}

	cb.SecretEngineCounts.Transit.MountAttributionLock.Lock()
	cb.SecretEngineCounts.Transit.MountAttribution[mountAccessor] = testBreakdown1Updated
	cb.SecretEngineCounts.Transit.MountAttributionLock.Unlock()

	err = core.UpdateTransitAttribution(ctx, currentMonth)
	require.NoError(t, err, "UpdateTransitAttribution with one mount should succeed")

	retrievedAttribution, err = core.GetStoredAttributionData(ctx, billing.LocalPrefix, currentMonth, billing.TransitDataProtectionCallCountsPrefix)
	require.NoError(t, err)
	require.Len(t, retrievedAttribution.Mounts, 2, "Should still have 2 mounts")
	// Top-level count: 25 (mount1) + 25.5 (mount2) = 50.5
	require.Equal(t, "50.5", fmt.Sprintf("%v", retrievedAttribution.Count), "Top-level count should be sum of all mount counts")

	// Verify the first mount has cumulative count (10 + 15 = 25)
	retrieved1, ok = retrievedAttribution.Mounts[mountAccessor]
	require.True(t, ok, "Should find breakdown for first mount accessor")
	require.Equal(t, "25", fmt.Sprintf("%v", retrieved1.Count), "Count should be cumulative: 10 + 15 = 25")
	require.Equal(t, mountAccessor, retrieved1.MountAccessor)
	require.Equal(t, "transit/", retrieved1.MountPath)

	// Second mount is unchanged
	retrieved2, ok = retrievedAttribution.Mounts[mountAccessor2]
	require.True(t, ok, "Should still find breakdown for transit-accessor-2")
	require.Equal(t, "25.5", fmt.Sprintf("%v", retrieved2.Count), "Count should remain unchanged for mount not in update")

	// Additional Test: Update both mounts with new counts to verify cumulative behavior
	t.Log("Additional Test: Update both mounts with new counts")
	testBreakdown1NewCount := logical.MountAttribution{
		MountAccessor:     mountAccessor,
		MountPath:         "transit/",
		NamespaceID:       "root",
		NamespacePath:     "",
		ParentNamespaceID: "",
		Count:             5.0, // Adding 5 more to the existing 25
	}
	testBreakdown2NewCount := logical.MountAttribution{
		MountAccessor:     mountAccessor2,
		MountPath:         "transit2/",
		NamespaceID:       "ns1-id",
		NamespacePath:     "ns1/",
		ParentNamespaceID: "root",
		Count:             10.5, // Adding 10.5 more to the existing 25.5
	}

	cb.SecretEngineCounts.Transit.MountAttributionLock.Lock()
	cb.SecretEngineCounts.Transit.MountAttribution[mountAccessor] = testBreakdown1NewCount
	cb.SecretEngineCounts.Transit.MountAttribution[mountAccessor2] = testBreakdown2NewCount
	cb.SecretEngineCounts.Transit.MountAttributionLock.Unlock()

	err = core.UpdateTransitAttribution(ctx, currentMonth)
	require.NoError(t, err, "UpdateTransitAttribution with both mounts should succeed")

	retrievedAttribution, err = core.GetStoredAttributionData(ctx, billing.LocalPrefix, currentMonth, billing.TransitDataProtectionCallCountsPrefix)
	require.NoError(t, err)
	require.Len(t, retrievedAttribution.Mounts, 2, "Should have 2 mounts")
	// Top-level count: 30 (mount1) + 36 (mount2) = 66
	require.Equal(t, "66", fmt.Sprintf("%v", retrievedAttribution.Count), "Top-level count should be sum of all mount counts")

	// Verify cumulative counts: mount1 = 25 + 5 = 30, mount2 = 25.5 + 10.5 = 36
	retrieved1, ok = retrievedAttribution.Mounts[mountAccessor]
	require.True(t, ok, "Should find breakdown for first mount accessor")
	require.Equal(t, "30", fmt.Sprintf("%v", retrieved1.Count), "Count should be cumulative: 25 + 5 = 30")

	retrieved2, ok = retrievedAttribution.Mounts[mountAccessor2]
	require.True(t, ok, "Should find breakdown for transit-accessor-2")
	require.Equal(t, "36", fmt.Sprintf("%v", retrieved2.Count), "Count should be cumulative: 25.5 + 10.5 = 36")

	// Verify in-memory map is empty after update (atomic swap behaviour)
	cb.SecretEngineCounts.Transit.MountAttributionLock.RLock()
	inMemoryCount := len(cb.SecretEngineCounts.Transit.MountAttribution)
	cb.SecretEngineCounts.Transit.MountAttributionLock.RUnlock()
	require.Equal(t, 0, inMemoryCount, "In-memory map should be empty after update (atomic swap)")

	// Test retrieval for a different month (should return empty mounts map)
	t.Log("Test: Retrieval for different month")
	differentMonth := currentMonth.AddDate(0, 1, 0)
	retrievedAttribution, err = core.GetStoredAttributionData(ctx, billing.LocalPrefix, differentMonth, billing.TransitDataProtectionCallCountsPrefix)
	require.NoError(t, err, "GetTransitAttribution for different month should succeed")
	require.Empty(t, retrievedAttribution.Mounts, "Should have no mounts for a month with no data")
}

// TestCertAttributionNamespaceMove verifies that when a cert mount's namespace
// metadata changes between two storeCertAttributionLocked calls (i.e. across
// flush boundaries), the stored entry reflects the new namespace while the
// count remains cumulative. It exercises storeCertAttributionLocked directly,
// covering PKI, SSH cert, and SSH OTP metric names in a single parameterised pass.
func TestCertAttributionNamespaceMove(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	view := &logical.InmemStorage{}
	month := time.Now().UTC()

	const accessor = "pki_ns_move"

	for _, metricName := range []string{
		billing.PkiDurationAdjustedCountPrefix,
		billing.SSHCertificateMetric,
		billing.SSHOTPMetric,
	} {
		metricName := metricName // capture
		t.Run(metricName, func(t *testing.T) {
			// Flush 1: mount is in ns1 with count 10.
			inNS1 := map[string]logical.MountAttribution{
				accessor: {
					MountAccessor: accessor,
					MountPath:     "cert/",
					MountType:     "pki",
					NamespaceID:   "ns1-id",
					NamespacePath: "ns1/",
					Count:         10.0,
				},
			}
			err := storeCertAttributionLocked(ctx, view, billing.LocalPrefix, metricName, 10.0, inNS1, month)
			require.NoError(t, err)

			stored, err := getStoredAttributionDataLocked(ctx, view, billing.LocalPrefix, month, metricName)
			require.NoError(t, err)
			require.Len(t, stored.Mounts, 1)
			entry := stored.Mounts[accessor]
			require.Equal(t, "ns1-id", entry.NamespaceID, "flush 1: NamespaceID should be ns1")
			require.Equal(t, "ns1/", entry.NamespacePath)
			require.Equal(t, "10", fmt.Sprintf("%v", entry.Count))

			// Flush 2: same accessor, mount has moved to ns2, count delta 5.
			inNS2 := map[string]logical.MountAttribution{
				accessor: {
					MountAccessor: accessor,
					MountPath:     "cert/",
					MountType:     "pki",
					NamespaceID:   "ns2-id",
					NamespacePath: "ns2/",
					Count:         5.0,
				},
			}
			err = storeCertAttributionLocked(ctx, view, billing.LocalPrefix, metricName, 5.0, inNS2, month)
			require.NoError(t, err)

			stored, err = getStoredAttributionDataLocked(ctx, view, billing.LocalPrefix, month, metricName)
			require.NoError(t, err)
			require.Len(t, stored.Mounts, 1, "still one entry — same accessor")

			entry = stored.Mounts[accessor]
			// Metadata must reflect ns2.
			require.Equal(t, "ns2-id", entry.NamespaceID, "flush 2: NamespaceID should be updated to ns2")
			require.Equal(t, "ns2/", entry.NamespacePath, "flush 2: NamespacePath should be updated to ns2")
			// Count must be cumulative: 10 + 5 = 15.
			require.Equal(t, "15", fmt.Sprintf("%v", entry.Count), "count should accumulate: 10 + 5 = 15")
			// Top-level total must also accumulate.
			require.Equal(t, "15", fmt.Sprintf("%v", stored.Count))
		})
	}
}
