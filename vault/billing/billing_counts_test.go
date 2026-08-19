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

	data := map[string]interface{}{
		"mountPath":        "transit/",
		"mountAccessor":    accessor,
		"mountType":        "transit",
		"backendAwareUUID": "uuid-xyz",
	}

	tracker := &AttributionTracker{
		MountAttribution: make(map[string]logical.MountAttribution),
	}

	// First call: mount is in ns1.
	ns1 := &namespace.Namespace{ID: "ns1-id", Path: "ns1/"}
	ctx1 := namespace.ContextWithNamespace(context.Background(), ns1)

	err := tracker.AccumulateMountAttributions(ctx1, data, 10.0, nil)
	require.NoError(t, err)

	tracker.MountAttributionLock.RLock()
	entry := tracker.MountAttribution[accessor]
	tracker.MountAttributionLock.RUnlock()

	require.Equal(t, "ns1-id", entry.NamespaceID, "first call: NamespaceID should be ns1")
	require.Equal(t, "ns1/", entry.NamespacePath)
	require.Equal(t, fmt.Sprintf("%v", 10.0), fmt.Sprintf("%v", entry.Count))

	// Second call: same accessor, mount has moved to ns2.
	ns2 := &namespace.Namespace{ID: "ns2-id", Path: "ns2/"}
	ctx2 := namespace.ContextWithNamespace(context.Background(), ns2)

	err = tracker.AccumulateMountAttributions(ctx2, data, 5.0, nil)
	require.NoError(t, err)

	tracker.MountAttributionLock.RLock()
	entry = tracker.MountAttribution[accessor]
	tracker.MountAttributionLock.RUnlock()

	// Metadata must reflect ns2 — the most recent call wins.
	require.Equal(t, "ns2-id", entry.NamespaceID, "second call: NamespaceID should be updated to ns2")
	require.Equal(t, "ns2/", entry.NamespacePath, "second call: NamespacePath should be updated to ns2")
	// Count must be the sum of both calls.
	require.Equal(t, fmt.Sprintf("%v", 15.0), fmt.Sprintf("%v", entry.Count), "count should accumulate: 10 + 5 = 15")
	require.Empty(t, entry.ParentNamespaceID)
}
