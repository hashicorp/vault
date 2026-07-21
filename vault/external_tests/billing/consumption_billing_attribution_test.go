// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package billing

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/testhelpers/minimal"
	"github.com/hashicorp/vault/helper/timeutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/billing"
	"github.com/stretchr/testify/require"
)

// TestDeleteExpiredAttributionData verifies that DeleteExpiredAttributionData removes
// attribution data older than DefaultAttributionRetentionMonths while preserving
// newer data and leaving regular billing metrics untouched.
func TestDeleteExpiredAttributionData(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	core := cluster.Cores[0].Core
	vault.TestWaitActive(t, core)

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

	view, ok := core.GetBillingSubView()
	require.True(t, ok)

	// Store attribution data for all three months under both prefixes
	for _, month := range []time.Time{monthToDelete, oldestRetainedMonth, currentMonth} {
		for _, prefix := range []string{billing.ReplicatedPrefix, billing.LocalPrefix} {
			err := core.StoreAttributionData(ctx, prefix, month, billing.KvHWMCountsHWM, attrData)
			require.NoError(t, err)
		}

		// Also store regular billing metrics alongside to verify they are not deleted
		err := view.Put(ctx, &logical.StorageEntry{
			Key:   billing.GetMonthlyBillingMetricPath(billing.LocalPrefix, month, billing.KvHWMCountsHWM),
			Value: []byte("20"),
		})
		require.NoError(t, err)
	}

	// Verify all attribution data exists before deletion
	for _, month := range []time.Time{monthToDelete, oldestRetainedMonth, currentMonth} {
		attrPath := billing.GetAttributionMaxPath(billing.LocalPrefix, month, billing.KvHWMCountsHWM)
		entry, err := view.Get(ctx, attrPath)
		require.NoError(t, err)
		require.NotNil(t, entry, "attribution should exist for month %s before deletion", month.Format("2006-01"))
	}

	// Call DeleteExpiredAttributionData
	err := core.DeleteExpiredAttributionData(ctx, currentMonth)
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

	// Regular billing metrics for the deleted month should be untouched by DeleteExpiredAttributionData
	kvCounts, err := core.GetStoredHWMKvCounts(ctx, billing.LocalPrefix, monthToDelete)
	require.NoError(t, err)
	require.Equal(t, 20, kvCounts, "regular billing metrics should not be affected by attribution deletion")

	// Now verify the inverse: DeleteExpiredBillingMetrics must not delete attribution data.
	// Store regular billing metrics for the billing-retention boundary month and re-run the
	// billing deletion to confirm attribution survives.
	billingMonthToDelete := currentMonth.AddDate(0, -billing.DefaultBillingRetentionMonths, 0)
	err = view.Put(ctx, &logical.StorageEntry{
		Key:   billing.GetMonthlyBillingMetricPath(billing.LocalPrefix, billingMonthToDelete, billing.KvHWMCountsHWM),
		Value: []byte("99"),
	})
	require.NoError(t, err)
	err = core.StoreAttributionData(ctx, billing.LocalPrefix, billingMonthToDelete, billing.KvHWMCountsHWM, attrData)
	require.NoError(t, err)

	err = core.DeleteExpiredBillingMetrics(ctx, currentMonth)
	require.NoError(t, err)

	// Regular billing metric at the billing boundary should be deleted
	billingKvCounts, err := core.GetStoredHWMKvCounts(ctx, billing.LocalPrefix, billingMonthToDelete)
	require.NoError(t, err)
	require.Equal(t, 0, billingKvCounts, "regular billing metric at billing boundary should be deleted")

	// Attribution at the billing boundary should still be present (independent retention)
	billingAttrPath := billing.GetAttributionMaxPath(billing.LocalPrefix, billingMonthToDelete, billing.KvHWMCountsHWM)
	billingAttrEntry, err := view.Get(ctx, billingAttrPath)
	require.NoError(t, err)
	require.NotNil(t, billingAttrEntry, "attribution data should NOT be deleted by DeleteExpiredBillingMetrics")
}

// TestStoreCertAttribution_PKI verifies the PKI attribution merge round-trip:
// two flushes to the same mount accumulate counts, and a second mount is added
// independently. MetricTypeAttribution.Count holds the running total.
func TestStoreCertAttribution_PKI(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	core := cluster.Cores[0].Core
	vault.TestWaitActive(t, core)

	ctx := context.Background()
	month := timeutil.StartOfMonth(time.Now().UTC())

	mount1 := logical.MountAttribution{
		MountAccessor: "pki_aaa",
		MountPath:     "pki/",
		MountType:     "pki",
		NamespaceID:   "root",
		NamespacePath: "",
		Count:         1.0,
	}
	mount2 := logical.MountAttribution{
		MountAccessor: "pki_bbb",
		MountPath:     "pki2/",
		MountType:     "pki",
		NamespaceID:   "ns1",
		NamespacePath: "ns1/",
		Count:         2.0,
	}

	// First flush: mount1 only, delta = 1.0
	err := core.StoreCertAttribution(ctx, billing.PkiDurationAdjustedCountPrefix, 1.0,
		map[string]logical.MountAttribution{"pki_aaa": mount1}, month)
	require.NoError(t, err)

	got, err := core.GetStoredAttributionData(ctx, billing.LocalPrefix, month, billing.PkiDurationAdjustedCountPrefix)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, "1", fmt.Sprintf("%v", got.Count))
	require.Len(t, got.Mounts, 1)
	require.Equal(t, "pki_aaa", got.Mounts["pki_aaa"].MountAccessor)
	require.Equal(t, "1", fmt.Sprintf("%v", got.Mounts["pki_aaa"].Count))

	// Second flush: mount1 again + mount2, delta = 3.0
	mount1v2 := mount1
	mount1v2.Count = 1.5 // additional units for mount1 this flush
	err = core.StoreCertAttribution(ctx, billing.PkiDurationAdjustedCountPrefix, 3.0,
		map[string]logical.MountAttribution{"pki_aaa": mount1v2, "pki_bbb": mount2}, month)
	require.NoError(t, err)

	got, err = core.GetStoredAttributionData(ctx, billing.LocalPrefix, month, billing.PkiDurationAdjustedCountPrefix)
	require.NoError(t, err)
	require.NotNil(t, got)
	// Running total: 1.0 + 3.0 = 4.0
	require.Equal(t, "4", fmt.Sprintf("%v", got.Count))
	require.Len(t, got.Mounts, 2)
	// mount1 per-mount total: 1.0 + 1.5 = 2.5
	require.Equal(t, "2.5", fmt.Sprintf("%v", got.Mounts["pki_aaa"].Count))
	// mount2 is new: 2.0
	require.Equal(t, "2", fmt.Sprintf("%v", got.Mounts["pki_bbb"].Count))
	require.Equal(t, "ns1", got.Mounts["pki_bbb"].NamespaceID)
	require.Equal(t, "ns1/", got.Mounts["pki_bbb"].NamespacePath)
}

// TestStoreCertAttribution_SSHCert verifies the SSH certificate attribution
// round-trip using the SSHCertificateMetric storage key.
func TestStoreCertAttribution_SSHCert(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	core := cluster.Cores[0].Core
	vault.TestWaitActive(t, core)

	ctx := context.Background()
	month := timeutil.StartOfMonth(time.Now().UTC())

	mount := logical.MountAttribution{
		MountAccessor: "ssh_cert_001",
		MountPath:     "ssh/",
		MountType:     "ssh",
		NamespaceID:   "root",
		Count:         0.5,
	}

	err := core.StoreCertAttribution(ctx, billing.SSHCertificateMetric, 0.5,
		map[string]logical.MountAttribution{"ssh_cert_001": mount}, month)
	require.NoError(t, err)

	got, err := core.GetStoredAttributionData(ctx, billing.LocalPrefix, month, billing.SSHCertificateMetric)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, "0.5", fmt.Sprintf("%v", got.Count))
	require.Len(t, got.Mounts, 1)
	require.Equal(t, "0.5", fmt.Sprintf("%v", got.Mounts["ssh_cert_001"].Count))
	require.Equal(t, "ssh_cert_001", got.Mounts["ssh_cert_001"].MountAccessor)
}

// TestStoreCertAttribution_SSHOTP verifies the SSH OTP attribution round-trip
// using the SSHOTPMetric storage key.
func TestStoreCertAttribution_SSHOTP(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	core := cluster.Cores[0].Core
	vault.TestWaitActive(t, core)

	ctx := context.Background()
	month := timeutil.StartOfMonth(time.Now().UTC())

	mount := logical.MountAttribution{
		MountAccessor: "ssh_otp_001",
		MountPath:     "ssh/",
		MountType:     "ssh",
		NamespaceID:   "root",
		Count:         0.0014,
	}

	err := core.StoreCertAttribution(ctx, billing.SSHOTPMetric, 0.0014,
		map[string]logical.MountAttribution{"ssh_otp_001": mount}, month)
	require.NoError(t, err)

	got, err := core.GetStoredAttributionData(ctx, billing.LocalPrefix, month, billing.SSHOTPMetric)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, "0.0014", fmt.Sprintf("%v", got.Count))
	require.Len(t, got.Mounts, 1)

	// A second OTP on the same mount accumulates
	mount2 := mount
	mount2.Count = 0.0014
	err = core.StoreCertAttribution(ctx, billing.SSHOTPMetric, 0.0014,
		map[string]logical.MountAttribution{"ssh_otp_001": mount2}, month)
	require.NoError(t, err)

	got, err = core.GetStoredAttributionData(ctx, billing.LocalPrefix, month, billing.SSHOTPMetric)
	require.NoError(t, err)
	require.Equal(t, "0.0028", fmt.Sprintf("%v", got.Count))
	require.Equal(t, "0.0028", fmt.Sprintf("%v", got.Mounts["ssh_otp_001"].Count))
}

// TestConsumeCertCounts_StoresAttribution verifies the full Active-node path:
// a CertCount with attribution maps fed into ConsumeCertCounts results in
// those attributions being persisted to the billing store.
func TestConsumeCertCounts_StoresAttribution(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	core := cluster.Cores[0].Core
	vault.TestWaitActive(t, core)

	ctx := context.Background()
	month := timeutil.StartOfMonth(time.Now().UTC())

	pkiMount := logical.MountAttribution{
		MountAccessor: "pki_consume",
		MountPath:     "pki/",
		MountType:     "pki",
		NamespaceID:   "root",
		Count:         1.0,
	}
	sshMount := logical.MountAttribution{
		MountAccessor: "ssh_consume",
		MountPath:     "ssh/",
		MountType:     "ssh",
		NamespaceID:   "root",
		Count:         0.5,
	}
	otpMount := logical.MountAttribution{
		MountAccessor: "otp_consume",
		MountPath:     "ssh/",
		MountType:     "ssh",
		NamespaceID:   "root",
		Count:         0.0014,
	}

	inc := logical.CertCount{
		IssuedCerts:              1,
		StoredCerts:              1,
		PkiDurationAdjustedCerts: 1.0,
		SSHIssuedCerts:           0.5,
		SSHIssuedOTPs:            0.0014,
		PkiMountAttributions:     map[string]logical.MountAttribution{"pki_consume": pkiMount},
		SshCertMountAttributions: map[string]logical.MountAttribution{"ssh_consume": sshMount},
		SshOtpMountAttributions:  map[string]logical.MountAttribution{"otp_consume": otpMount},
	}

	// ConsumeCertCounts checks HAState; the cluster core is Active.
	core.ConsumeCertCounts(inc)

	// PKI attribution
	pkiAttr, err := core.GetStoredAttributionData(ctx, billing.LocalPrefix, month, billing.PkiDurationAdjustedCountPrefix)
	require.NoError(t, err)
	require.NotNil(t, pkiAttr)
	require.Len(t, pkiAttr.Mounts, 1)
	require.Contains(t, pkiAttr.Mounts, "pki_consume")

	// SSH cert attribution
	sshAttr, err := core.GetStoredAttributionData(ctx, billing.LocalPrefix, month, billing.SSHCertificateMetric)
	require.NoError(t, err)
	require.NotNil(t, sshAttr)
	require.Len(t, sshAttr.Mounts, 1)
	require.Contains(t, sshAttr.Mounts, "ssh_consume")

	// SSH OTP attribution
	otpAttr, err := core.GetStoredAttributionData(ctx, billing.LocalPrefix, month, billing.SSHOTPMetric)
	require.NoError(t, err)
	require.NotNil(t, otpAttr)
	require.Len(t, otpAttr.Mounts, 1)
	require.Contains(t, otpAttr.Mounts, "otp_consume")
}
