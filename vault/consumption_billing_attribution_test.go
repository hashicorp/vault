// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

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
				MountPath:           "secret/",
				MountType:           "kv",
				MountAccessor:       "kv_5d4f8f1c",
				MountRunningVersion: "wre_43",
				NamespaceID:         "root",
				NamespacePath:       "",
				ParentNamespaceID:   "",
				Count:               5,
				BackendAwareUUID:    "wdasd23",
			},
			"kv_be9766a3": {
				MountPath:           "kv/",
				MountType:           "kv",
				MountAccessor:       "kv_be9766a3",
				MountRunningVersion: "wtm_21",
				NamespaceID:         "3bFWF",
				NamespacePath:       "ns1/",
				ParentNamespaceID:   "root",
				Count:               5,
				BackendAwareUUID:    "adwdsd35",
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
	require.Len(t, got.Mounts, 2)

	m1 := got.Mounts["kv_5d4f8f1c"]
	require.Equal(t, "secret/", m1.MountPath)
	require.Equal(t, "kv", m1.MountType)
	require.Equal(t, "wre_43", m1.MountRunningVersion)
	require.Equal(t, "root", m1.NamespaceID)
	require.Equal(t, "", m1.NamespacePath)
	require.Equal(t, "", m1.ParentNamespaceID)
	require.Equal(t, "kv_5d4f8f1c", m1.MountAccessor)
	require.Equal(t, "5", fmt.Sprintf("%v", m1.Count))

	m2 := got.Mounts["kv_be9766a3"]
	require.Equal(t, "kv/", m2.MountPath)
	require.Equal(t, "wtm_21", m2.MountRunningVersion)
	require.Equal(t, "3bFWF", m2.NamespaceID)
	require.Equal(t, "ns1/", m2.NamespacePath)
	require.Equal(t, "", m1.ParentNamespaceID)
	require.Equal(t, "5", fmt.Sprintf("%v", m2.Count))

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
