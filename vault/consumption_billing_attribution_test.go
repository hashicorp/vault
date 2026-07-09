// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/timeutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault/billing"
	"github.com/stretchr/testify/require"
)

// TestStoreAndGetAttributionData verifies the round-trip of storeAttributionDataLocked
// and getStoredAttributionDataLocked, and the public GetStoredAttributionData wrapper.
// It also verifies that a second store overwrites the previous entry (no implicit merge —
// callers are responsible for merging before storing).
func TestStoreAndGetAttributionData(t *testing.T) {
	core, _, _ := TestCoreUnsealed(t)
	ctx := context.Background()

	now := time.Now().UTC()
	month := timeutil.StartOfMonth(now)
	lastUpdated := time.Date(2026, 5, 14, 18, 7, 23, 0, time.UTC)

	data := &logical.MetricTypeAttribution{
		Count:       10,
		LastUpdated: lastUpdated,
		Mounts: map[string]logical.MountAttribution{
			"kv_5d4f8f1c": {
				MountPath:        "secret/",
				MountType:        "kv",
				MountAccessor:    "kv_5d4f8f1c",
				NamespaceID:      "root",
				NamespacePath:    "",
				Count:            5,
				BackendAwareUUID: "wdasd23",
			},
			"kv_be9766a3": {
				MountPath:        "kv/",
				MountType:        "kv",
				MountAccessor:    "kv_be9766a3",
				NamespaceID:      "3bFWF",
				NamespacePath:    "ns1/",
				Count:            5,
				BackendAwareUUID: "adwdsd35",
			},
		},
	}

	// Store via the locked helper (lock is not held in tests since there's no contention)
	err := core.storeAttributionDataLocked(ctx, billing.LocalPrefix, month, billing.KvHWMCountsHWM, data)
	require.NoError(t, err)

	// Retrieve via the locked helper
	got, err := core.getStoredAttributionDataLocked(ctx, billing.LocalPrefix, month, billing.KvHWMCountsHWM)
	require.NoError(t, err)
	require.NotNil(t, got)

	// Count is interface{} — jsonutil.DecodeJSON deserialises numbers as json.Number.
	// Compare via fmt.Sprintf to avoid type mismatch between int and json.Number.
	require.Equal(t, "10", fmt.Sprintf("%v", got.Count))
	require.Equal(t, data.LastUpdated.UTC(), got.LastUpdated.UTC())
	require.Len(t, got.Mounts, 2)

	m1 := got.Mounts["kv_5d4f8f1c"]
	require.Equal(t, "secret/", m1.MountPath)
	require.Equal(t, "kv", m1.MountType)
	require.Equal(t, "root", m1.NamespaceID)
	require.Equal(t, "", m1.NamespacePath)
	require.Equal(t, "kv_5d4f8f1c", m1.MountAccessor)
	require.Equal(t, "5", fmt.Sprintf("%v", m1.Count))

	m2 := got.Mounts["kv_be9766a3"]
	require.Equal(t, "kv/", m2.MountPath)
	require.Equal(t, "3bFWF", m2.NamespaceID)
	require.Equal(t, "ns1/", m2.NamespacePath)
	require.Equal(t, "5", fmt.Sprintf("%v", m2.Count))

	// Overwrite with new data — second store must replace, not merge.
	overwrite := &logical.MetricTypeAttribution{
		Count:       12,
		LastUpdated: time.Now().UTC(),
		Mounts: map[string]logical.MountAttribution{
			"kv_bbb": {Count: 12, MountAccessor: "kv_bbb", MountPath: "new/", MountType: "kv"},
		},
	}
	err = core.storeAttributionDataLocked(ctx, billing.LocalPrefix, month, billing.KvHWMCountsHWM, overwrite)
	require.NoError(t, err)

	got, err = core.getStoredAttributionDataLocked(ctx, billing.LocalPrefix, month, billing.KvHWMCountsHWM)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, "12", fmt.Sprintf("%v", got.Count))
	require.Len(t, got.Mounts, 1, "overwrite should replace all previous mounts")
	_, hasOld := got.Mounts["kv_5d4f8f1c"]
	require.False(t, hasOld, "old mounts should be gone after overwrite")
	_, hasNew := got.Mounts["kv_bbb"]
	require.True(t, hasNew, "new mount should be present after overwrite")
}

// TestDeleteExpiredAttributionData verifies that deleteExpiredAttributionData removes
// attribution data older than DefaultAttributionRetentionMonths while preserving
// newer data and leaving regular billing metrics untouched.
func TestDeleteExpiredAttributionData(t *testing.T) {
	coreConfig := &CoreConfig{
		LogicalBackends: roleLogicalBackends,
	}
	core, _, _ := TestCoreUnsealedWithConfig(t, coreConfig)
	ctx := context.Background()

	now := time.Now().UTC()
	currentMonth := timeutil.StartOfMonth(now)
	oldestRetainedMonth := currentMonth.AddDate(0, -(billing.DefaultAttributionRetentionMonths - 1), 0)
	monthToDelete := currentMonth.AddDate(0, -billing.DefaultAttributionRetentionMonths, 0)

	attrData := &logical.MetricTypeAttribution{
		Count:       7,
		LastUpdated: time.Now().UTC(),
		Mounts: map[string]logical.MountAttribution{
			"kv_test": {Count: 7, MountAccessor: "kv_test", MountPath: "secret/", MountType: "kv"},
		},
	}

	// Store attribution data for all three months under both prefixes
	for _, month := range []time.Time{monthToDelete, oldestRetainedMonth, currentMonth} {
		for _, prefix := range []string{billing.ReplicatedPrefix, billing.LocalPrefix} {
			err := core.storeAttributionDataLocked(ctx, prefix, month, billing.KvHWMCountsHWM, attrData)
			require.NoError(t, err)
		}

		// Also store regular billing metrics alongside to verify they are not deleted
		core.storeMaxKvCountsLocked(ctx, 20, billing.LocalPrefix, month)
	}

	// Verify all attribution data exists before deletion
	view, ok := core.GetBillingSubView()
	require.True(t, ok)
	for _, month := range []time.Time{monthToDelete, oldestRetainedMonth, currentMonth} {
		attrPath := billing.GetAttributionMaxPath(billing.LocalPrefix, month, billing.KvHWMCountsHWM)
		entry, err := view.Get(ctx, attrPath)
		require.NoError(t, err)
		require.NotNil(t, entry, "attribution should exist for month %s before deletion", month.Format("2006-01"))
	}

	// Call deleteExpiredAttributionData
	err := core.deleteExpiredAttributionData(ctx, currentMonth)
	require.NoError(t, err)

	// Month to delete: attribution should be gone
	for _, prefix := range []string{billing.ReplicatedPrefix, billing.LocalPrefix} {
		attrPath := billing.GetAttributionMaxPath(prefix, monthToDelete, billing.KvHWMCountsHWM)
		entry, err := view.Get(ctx, attrPath)
		require.NoError(t, err)
		require.Nil(t, entry, "attribution for %s should be deleted", monthToDelete.Format("2006-01"))
	}

	// Oldest retained month: attribution should still exist
	for _, prefix := range []string{billing.ReplicatedPrefix, billing.LocalPrefix} {
		attrPath := billing.GetAttributionMaxPath(prefix, oldestRetainedMonth, billing.KvHWMCountsHWM)
		entry, err := view.Get(ctx, attrPath)
		require.NoError(t, err)
		require.NotNil(t, entry, "attribution for %s should be kept", oldestRetainedMonth.Format("2006-01"))
	}

	// Current month: attribution should still exist
	for _, prefix := range []string{billing.ReplicatedPrefix, billing.LocalPrefix} {
		attrPath := billing.GetAttributionMaxPath(prefix, currentMonth, billing.KvHWMCountsHWM)
		entry, err := view.Get(ctx, attrPath)
		require.NoError(t, err)
		require.NotNil(t, entry, "attribution for current month should be kept")
	}

	// Regular billing metrics for the deleted month should be untouched by deleteExpiredAttributionData
	kvCounts, err := core.GetStoredHWMKvCounts(ctx, billing.LocalPrefix, monthToDelete)
	require.NoError(t, err)
	require.Equal(t, 20, kvCounts, "regular billing metrics should not be affected by attribution deletion")

	// Now verify the inverse: deleteExpiredBillingMetrics must not delete attribution data.
	// The attribution data for monthToDelete is still present (deleteExpiredAttributionData only
	// deletes at DefaultAttributionRetentionMonths boundary, not DefaultBillingRetentionMonths).
	// Store regular billing metrics for the billing-retention boundary month and re-run the
	// billing deletion to confirm attribution survives.
	billingMonthToDelete := currentMonth.AddDate(0, -billing.DefaultBillingRetentionMonths, 0)
	core.storeMaxKvCountsLocked(ctx, 99, billing.LocalPrefix, billingMonthToDelete)
	err = core.storeAttributionDataLocked(ctx, billing.LocalPrefix, billingMonthToDelete, billing.KvHWMCountsHWM, attrData)
	require.NoError(t, err)

	err = core.deleteExpiredBillingMetrics(ctx, currentMonth)
	require.NoError(t, err)

	// Regular billing metric at the billing boundary should be deleted
	billingKvCounts, err := core.GetStoredHWMKvCounts(ctx, billing.LocalPrefix, billingMonthToDelete)
	require.NoError(t, err)
	require.Equal(t, 0, billingKvCounts, "regular billing metric at billing boundary should be deleted")

	// Attribution at the billing boundary should still be present (independent retention)
	billingAttrPath := billing.GetAttributionMaxPath(billing.LocalPrefix, billingMonthToDelete, billing.KvHWMCountsHWM)
	billingAttrEntry, err := view.Get(ctx, billingAttrPath)
	require.NoError(t, err)
	require.NotNil(t, billingAttrEntry, "attribution data should NOT be deleted by deleteExpiredBillingMetrics")
}
