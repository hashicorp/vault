// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package billing

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

// TestGetMonthlyBillingMetricPath verifies the TestGetMonthlyBillingMetricPath function
// returns the correct billing metric path for the given product area and month
func TestGetMonthlyBillingMetricPath(t *testing.T) {
	ts := time.Date(2026, time.January, 15, 12, 0, 0, 0, time.UTC)

	got := GetMonthlyBillingMetricPath(ReplicatedPrefix, ts, KvHWMCountsHWM)
	want := "replicated/2026/01/maxKvCounts/"
	require.Equal(t, got, want)
}

// TestGetMonthlyBillingPath verifies the GetMonthlyBillingPath function
// returns the correct billing path for the given product area and month
func TestGetMonthlyBillingPath(t *testing.T) {
	ts := time.Date(2026, time.January, 15, 12, 0, 0, 0, time.UTC)

	got := GetMonthlyBillingPath(ReplicatedPrefix, ts)
	want := "replicated/2026/01/"
	require.Equal(t, got, want)
}

// TestGetAttributionPath verifies the GetAttributionPath function
// returns the correct attribution path for the given product area and month
func TestGetAttributionPath(t *testing.T) {
	ts := time.Date(2026, time.January, 15, 12, 0, 0, 0, time.UTC)

	got := GetAttributionMaxPath(ReplicatedPrefix, ts, RoleHWMCountsHWM)
	want := "replicated/2026/01/attribution/maximum/maxRoleCounts/"
	require.Equal(t, got, want)
}

// TestAccumulateMountAttributions_MetadataUpdate verifies that when the same mount
// accessor appears twice within the same accumulation window (before a flush),
// the second call's namespace/path metadata overwrites the first, while the count
// is accumulated across both calls.
func TestAccumulateMountAttributions_MetadataUpdate(t *testing.T) {
	const accessor = "transit_abc123"

	data1 := map[string]interface{}{
		"mountPath":           "transit-v1/",
		"mountAccessor":       accessor,
		"mountType":           "transit",
		"backendAwareUUID":    "uuid-aaa",
		"mountRunningVersion": "v1.0.0",
	}

	tracker := &AttributionTracker{
		MountAttribution: make(map[string]logical.MountAttribution),
	}

	// First call: mount is in ns1.
	ns1 := &namespace.Namespace{ID: "ns1-id", Path: "ns1/"}
	ctx1 := namespace.ContextWithNamespace(context.Background(), ns1)

	err := tracker.AccumulateMountAttributions(ctx1, data1, 10.0, nil)
	require.NoError(t, err)

	tracker.MountAttributionLock.RLock()
	entry := tracker.MountAttribution[accessor]
	tracker.MountAttributionLock.RUnlock()

	require.Equal(t, "transit-v1/", entry.MountPath, "first call: MountPath should be transit-v1/")
	require.Equal(t, "transit", entry.MountType, "first call: MountType should be transit")
	require.Equal(t, "uuid-aaa", entry.BackendAwareUUID, "first call: BackendAwareUUID should be uuid-aaa")
	require.Equal(t, "v1.0.0", entry.MountRunningVersion, "first call: MountRunningVersion should be v1.0.0")
	require.Equal(t, "ns1-id", entry.NamespaceID, "first call: NamespaceID should be ns1")
	require.Equal(t, "ns1/", entry.NamespacePath, "first call: NamespacePath should be ns1/")
	require.Equal(t, fmt.Sprintf("%v", 10.0), fmt.Sprintf("%v", entry.Count))

	// Second call: same accessor, all metadata has changed (mount renamed, plugin
	// upgraded, namespace moved). Every field must reflect the newer values.
	data2 := map[string]interface{}{
		"mountPath":           "transit-v2/",
		"mountAccessor":       accessor,
		"mountType":           "transit",
		"backendAwareUUID":    "uuid-bbb",
		"mountRunningVersion": "v2.0.0",
	}
	ns2 := &namespace.Namespace{ID: "ns2-id", Path: "ns2/"}
	ctx2 := namespace.ContextWithNamespace(context.Background(), ns2)

	err = tracker.AccumulateMountAttributions(ctx2, data2, 5.0, nil)
	require.NoError(t, err)

	tracker.MountAttributionLock.RLock()
	entry = tracker.MountAttribution[accessor]
	tracker.MountAttributionLock.RUnlock()

	// All metadata fields must reflect the second call — the most recent call wins.
	require.Equal(t, "transit-v2/", entry.MountPath, "second call: MountPath should be updated to transit-v2/")
	require.Equal(t, "transit", entry.MountType, "second call: MountType should be transit")
	require.Equal(t, "uuid-bbb", entry.BackendAwareUUID, "second call: BackendAwareUUID should be updated to uuid-bbb")
	require.Equal(t, "v2.0.0", entry.MountRunningVersion, "second call: MountRunningVersion should be updated to v2.0.0")
	require.Equal(t, "ns2-id", entry.NamespaceID, "second call: NamespaceID should be updated to ns2")
	require.Equal(t, "ns2/", entry.NamespacePath, "second call: NamespacePath should be updated to ns2/")
	// Count must be accumulated across both calls.
	require.Equal(t, fmt.Sprintf("%v", 15.0), fmt.Sprintf("%v", entry.Count), "count should accumulate: 10 + 5 = 15")
	require.Empty(t, entry.ParentNamespaceID)
}
