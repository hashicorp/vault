// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package billing

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/billing"
	"github.com/stretchr/testify/require"
)

// TestGcpKmsDataProtectionCallCounts tests that we correctly store and track
// the GCP KMS data protection call counts by simulating billing operations.
func TestGcpKmsDataProtectionCallCounts(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		BillingConfig: billing.BillingConfig{
			MetricsUpdateCadence: 3 * time.Second,
		},
	}
	core, _, _, _ := vault.TestCoreUnsealedWithMetricsAndConfig(t, coreConfig)

	currentMonth := time.Now()
	ctx := namespace.RootContext(context.Background())

	// Get the consumption billing manager
	cbm := core.GetConsumptionBillingManager()
	require.NotNil(t, cbm)

	// Simulate GCP KMS operations by directly calling the billing manager
	// This tests the Vault-side tracking without needing the actual plugin

	// Simulate encrypt operation
	err := cbm.WriteBillingData(ctx, "gcpkms", map[string]interface{}{"count": uint64(1)})
	require.NoError(t, err)
	require.Equal(t, uint64(1), core.GetInMemoryGcpKmsDataProtectionCallCounts())

	// Wait for storage update
	require.Eventually(t, func() bool {
		counts, err := core.GetStoredGcpKmsCallCounts(context.Background(), currentMonth)
		return err == nil && counts == 1
	}, 5*time.Second, 100*time.Millisecond)
	require.Equal(t, uint64(0), core.GetInMemoryGcpKmsDataProtectionCallCounts())

	// Simulate decrypt operation
	err = cbm.WriteBillingData(ctx, "gcpkms", map[string]interface{}{"count": uint64(1)})
	require.NoError(t, err)
	require.Equal(t, uint64(1), core.GetInMemoryGcpKmsDataProtectionCallCounts())

	// Wait for storage update
	require.Eventually(t, func() bool {
		counts, err := core.GetStoredGcpKmsCallCounts(context.Background(), currentMonth)
		return err == nil && counts == 2
	}, 5*time.Second, 100*time.Millisecond)
	require.Equal(t, uint64(0), core.GetInMemoryGcpKmsDataProtectionCallCounts())

	// Simulate reencrypt operation
	err = cbm.WriteBillingData(ctx, "gcpkms", map[string]interface{}{"count": uint64(1)})
	require.NoError(t, err)
	require.Equal(t, uint64(1), core.GetInMemoryGcpKmsDataProtectionCallCounts())

	// Wait for storage update
	require.Eventually(t, func() bool {
		counts, err := core.GetStoredGcpKmsCallCounts(context.Background(), currentMonth)
		return err == nil && counts == 3
	}, 5*time.Second, 100*time.Millisecond)
	require.Equal(t, uint64(0), core.GetInMemoryGcpKmsDataProtectionCallCounts())

	// Simulate sign operation
	err = cbm.WriteBillingData(ctx, "gcpkms", map[string]interface{}{"count": uint64(1)})
	require.NoError(t, err)
	require.Equal(t, uint64(1), core.GetInMemoryGcpKmsDataProtectionCallCounts())

	// Wait for storage update
	require.Eventually(t, func() bool {
		counts, err := core.GetStoredGcpKmsCallCounts(context.Background(), currentMonth)
		return err == nil && counts == 4
	}, 5*time.Second, 100*time.Millisecond)
	require.Equal(t, uint64(0), core.GetInMemoryGcpKmsDataProtectionCallCounts())

	// Simulate verify operation
	err = cbm.WriteBillingData(ctx, "gcpkms", map[string]interface{}{"count": uint64(1)})
	require.NoError(t, err)
	require.Equal(t, uint64(1), core.GetInMemoryGcpKmsDataProtectionCallCounts())

	// Wait for storage update
	require.Eventually(t, func() bool {
		counts, err := core.GetStoredGcpKmsCallCounts(context.Background(), currentMonth)
		return err == nil && counts == 5
	}, 5*time.Second, 100*time.Millisecond)
	require.Equal(t, uint64(0), core.GetInMemoryGcpKmsDataProtectionCallCounts())

	// Run update again and make sure the value in storage is still 5
	counts, err := core.UpdateGcpKmsCallCounts(context.Background(), currentMonth)
	require.NoError(t, err)
	require.Equal(t, uint64(5), counts)

	// Verify the value in storage is still 5
	counts, err = core.GetStoredGcpKmsCallCounts(context.Background(), currentMonth)
	require.NoError(t, err)
	require.Equal(t, uint64(5), counts)
}
