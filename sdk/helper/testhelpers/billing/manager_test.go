// Copyright IBM Corp. 2026
// SPDX-License-Identifier: MPL-2.0

package billing_test

import (
	"context"
	"sync"
	"testing"

	"github.com/hashicorp/vault/sdk/helper/testhelpers/billing"
	"github.com/stretchr/testify/require"
)

// TestMockConsumptionBillingManager_WriteAndAttribution verifies that calls to WriteBillingData
// aggregate MountAttribution per mountAccessor and track the total count.
func TestMockConsumptionBillingManager_WriteAndAttribution(t *testing.T) {
	t.Parallel()

	m := billing.NewMockConsumptionBillingManager()
	ctx := context.Background()

	require.Equal(t, uint64(0), m.TotalCount())

	data1 := map[string]interface{}{
		"count":         uint64(5),
		"mountAccessor": "transit_abc123",
		"mountPath":     "transit/",
		"mountType":     "transit",
	}
	require.NoError(t, m.WriteBillingData(ctx, "transit", data1))

	data2 := map[string]interface{}{
		"count":         uint64(10),
		"mountAccessor": "gcpkms_def456",
		"mountPath":     "gcpkms/",
		"mountType":     "gcpkms",
	}
	require.NoError(t, m.WriteBillingData(ctx, "gcpkms", data2))

	data3 := map[string]interface{}{
		"count":         uint64(15),
		"mountAccessor": "transit_abc123",
		"mountPath":     "transit/",
		"mountType":     "transit",
	}
	require.NoError(t, m.WriteBillingData(ctx, "transit", data3))

	require.Equal(t, uint64(30), m.TotalCount())

	// Verify aggregated MountAttribution for transit (5 + 15 = 20)
	transitAttr, ok := m.GetAttribution("transit_abc123")
	require.True(t, ok)
	require.Equal(t, uint64(20), transitAttr.Count)
	require.Equal(t, "transit/", transitAttr.MountPath)
	require.Equal(t, "transit_abc123", transitAttr.MountAccessor)
	require.Equal(t, "transit", transitAttr.MountType)

	// Verify aggregated MountAttribution for gcpkms (10)
	gcpkmsAttr, ok := m.GetAttribution("gcpkms_def456")
	require.True(t, ok)
	require.Equal(t, uint64(10), gcpkmsAttr.Count)
	require.Equal(t, "gcpkms/", gcpkmsAttr.MountPath)

	// Verify non-existent mount
	_, ok = m.GetAttribution("unknown")
	require.False(t, ok)
}

// TestMockConsumptionBillingManager_Units verifies that WriteBillingData accumulates
// float64 units correctly for plugins like external-ca.
func TestMockConsumptionBillingManager_Units(t *testing.T) {
	t.Parallel()

	m := billing.NewMockConsumptionBillingManager()
	ctx := context.Background()

	require.Equal(t, float64(0), m.TotalUnits())

	data1 := map[string]interface{}{
		"units":               1.3,
		"mountAccessor":       "external_ca_1",
		"mountPath":           "external-ca/",
		"mountType":           "external-ca",
		"backendAwareUUID":    "uuid-123",
		"mountRunningVersion": "v1.0.0",
	}
	require.NoError(t, m.WriteBillingData(ctx, "external-ca", data1))

	data2 := map[string]interface{}{
		"units":               2.7,
		"mountAccessor":       "external_ca_1",
		"mountPath":           "external-ca/",
		"mountType":           "external-ca",
		"backendAwareUUID":    "uuid-123",
		"mountRunningVersion": "v1.0.0",
	}
	require.NoError(t, m.WriteBillingData(ctx, "external-ca", data2))

	require.InDelta(t, 4.0, m.TotalUnits(), 0.00001)

	attr, ok := m.GetAttribution("external_ca_1")
	require.True(t, ok)
	require.InDelta(t, 4.0, attr.Count.(float64), 0.00001)
	require.Equal(t, "external-ca/", attr.MountPath)
	require.Equal(t, "external-ca", attr.MountType)
	require.Equal(t, "external_ca_1", attr.MountAccessor)
	require.Equal(t, "uuid-123", attr.BackendAwareUUID)
	require.Equal(t, "v1.0.0", attr.MountRunningVersion)
}

// TestMockConsumptionBillingManager_Concurrency verifies that MockConsumptionBillingManager
// safely supports concurrent writes and reads without races.
func TestMockConsumptionBillingManager_Concurrency(t *testing.T) {
	t.Parallel()

	m := billing.NewMockConsumptionBillingManager()
	ctx := context.Background()

	const numGoroutines = 50
	var wg sync.WaitGroup
	wg.Add(numGoroutines * 2)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			_ = m.WriteBillingData(ctx, "transit", map[string]interface{}{
				"count":         uint64(1),
				"mountAccessor": "transit_concurrent",
			})
		}()
		go func() {
			defer wg.Done()
			_ = m.TotalCount()
			_, _ = m.GetAttribution("transit_concurrent")
		}()
	}

	wg.Wait()
	require.Equal(t, uint64(numGoroutines), m.TotalCount())
}
