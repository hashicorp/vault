// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package billing

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/pluginconsts"
	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	"github.com/hashicorp/vault/helper/testhelpers/minimal"
	"github.com/hashicorp/vault/helper/testhelpers/pluginhelpers"
	"github.com/hashicorp/vault/helper/timeutil"
	sdkconsts "github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/billing"
	"github.com/stretchr/testify/require"
)

// requireAttrCountEqualsMountSum asserts that MetricTypeAttribution.Count equals the sum
// of all per-mount MountAttribution.Count values in the stored attribution blob.  It is
// called at the end of every test that verifies a final stored state so that every metric
// path is covered by the invariant without a separate standalone test.
func requireAttrCountEqualsMountSum(t *testing.T, attr *logical.MetricTypeAttribution) {
	t.Helper()
	attrCount := vault.ToFloat64(attr.Count)
	var mountSum float64
	for _, m := range attr.Mounts {
		mountSum += vault.ToFloat64(m.Count)
	}
	require.InDelta(t, attrCount, mountSum, 1e-9,
		"MetricTypeAttribution.Count (%v) must equal sum of per-mount counts (%v)", attrCount, mountSum)
}

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
// independently. MetricTypeAttribution.Count holds the running total, and is
// always equal to the separately-stored scalar (GetStoredPkiDurationAdjustedCount),
// because production code calls UpdatePkiDurationAdjustedCount and StoreCertAttribution
// with the same countDelta value in every flush.
func TestStoreCertAttribution_PKI(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	core := cluster.Cores[0].Core
	vault.TestWaitActive(t, core)

	ctx := context.Background()
	month := timeutil.StartOfMonth(time.Now().UTC())

	mount1 := logical.MountAttribution{
		MountAccessor:       "pki_aaa",
		MountPath:           "pki/",
		MountType:           "pki",
		NamespaceID:         "root",
		NamespacePath:       "",
		Count:               1.0,
		MountRunningVersion: "version1",
	}
	mount2 := logical.MountAttribution{
		MountAccessor:       "pki_bbb",
		MountPath:           "pki2/",
		MountType:           "pki",
		NamespaceID:         "ns1",
		NamespacePath:       "ns1/",
		Count:               2.0,
		MountRunningVersion: "version2",
	}

	// First flush: mount1 only, delta = 1.0 (== mount1.Count).
	// Mirror production: update scalar and attribution with the same delta.
	err := core.UpdatePkiDurationAdjustedCount(ctx, 1.0, month)
	require.NoError(t, err)
	err = core.StoreCertAttribution(ctx, billing.PkiDurationAdjustedCountPrefix, 1.0,
		map[string]logical.MountAttribution{"pki_aaa": mount1}, month)
	require.NoError(t, err)

	got, err := core.GetStoredAttributionData(ctx, billing.LocalPrefix, month, billing.PkiDurationAdjustedCountPrefix)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, "1", fmt.Sprintf("%v", got.Count))
	require.Len(t, got.Mounts, 1)
	require.Equal(t, "pki_aaa", got.Mounts["pki_aaa"].MountAccessor)
	require.Equal(t, "1", fmt.Sprintf("%v", got.Mounts["pki_aaa"].Count))
	requireAttrCountEqualsMountSum(t, got)

	// Second flush: mount1 (1.5 additional units) + mount2 (2.0 units).
	// delta = 1.5 + 2.0 = 3.5 — identical to sum(incomingMounts.Count).
	mount1v2 := mount1
	mount1v2.Count = 1.5
	mount1v2.MountRunningVersion = "version1b"
	err = core.UpdatePkiDurationAdjustedCount(ctx, 3.5, month)
	require.NoError(t, err)
	err = core.StoreCertAttribution(ctx, billing.PkiDurationAdjustedCountPrefix, 3.5,
		map[string]logical.MountAttribution{"pki_aaa": mount1v2, "pki_bbb": mount2}, month)
	require.NoError(t, err)

	got, err = core.GetStoredAttributionData(ctx, billing.LocalPrefix, month, billing.PkiDurationAdjustedCountPrefix)
	require.NoError(t, err)
	require.NotNil(t, got)
	// Running total: 1.0 + 3.5 = 4.5
	require.Equal(t, "4.5", fmt.Sprintf("%v", got.Count))
	require.Len(t, got.Mounts, 2)
	// mount1 per-mount total: 1.0 + 1.5 = 2.5; version reflects latest flush
	require.Equal(t, "2.5", fmt.Sprintf("%v", got.Mounts["pki_aaa"].Count))
	require.Equal(t, "version1b", got.Mounts["pki_aaa"].MountRunningVersion, "MountRunningVersion should reflect latest flush")
	// mount2 is new: 2.0
	require.Equal(t, "2", fmt.Sprintf("%v", got.Mounts["pki_bbb"].Count))
	require.Equal(t, "ns1", got.Mounts["pki_bbb"].NamespaceID)
	require.Equal(t, "ns1/", got.Mounts["pki_bbb"].NamespacePath)
	require.Equal(t, "version2", got.Mounts["pki_bbb"].MountRunningVersion, "MountRunningVersion should be stored for new mount")
	requireAttrCountEqualsMountSum(t, got)

	// Billing scalar must equal attribution Count (both are running totals of the same deltas).
	pkiScalar, err := core.GetStoredPkiDurationAdjustedCount(ctx, month)
	require.NoError(t, err)
	require.InDelta(t, pkiScalar, vault.ToFloat64(got.Count), 1e-9,
		"PKI scalar (%v) must equal attribution Count (%v)", pkiScalar, got.Count)
}

// TestStoreCertAttribution_SSHCert verifies the SSH certificate attribution
// round-trip using the SSHCertificateMetric storage key.
// The scalar (GetStoredSSHDurationAdjustedCertCount) is updated with the same
// delta as the attribution blob, mirroring what ConsumeCertCounts does in production.
func TestStoreCertAttribution_SSHCert(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	core := cluster.Cores[0].Core
	vault.TestWaitActive(t, core)

	ctx := context.Background()
	month := timeutil.StartOfMonth(time.Now().UTC())

	mount := logical.MountAttribution{
		MountAccessor:       "ssh_cert_001",
		MountPath:           "ssh/",
		MountType:           "ssh",
		NamespaceID:         "root",
		Count:               0.5,
		MountRunningVersion: "version1",
	}

	// Mirror production: update scalar and attribution with the same delta.
	_, err := core.UpdateStoredSSHDurationAdjustedCertCount(ctx, month, 0.5)
	require.NoError(t, err)
	err = core.StoreCertAttribution(ctx, billing.SSHCertificateMetric, 0.5,
		map[string]logical.MountAttribution{"ssh_cert_001": mount}, month)
	require.NoError(t, err)

	got, err := core.GetStoredAttributionData(ctx, billing.LocalPrefix, month, billing.SSHCertificateMetric)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, "0.5", fmt.Sprintf("%v", got.Count))
	require.Len(t, got.Mounts, 1)
	require.Equal(t, "0.5", fmt.Sprintf("%v", got.Mounts["ssh_cert_001"].Count))
	require.Equal(t, "ssh_cert_001", got.Mounts["ssh_cert_001"].MountAccessor)
	require.Equal(t, "version1", got.Mounts["ssh_cert_001"].MountRunningVersion)
	requireAttrCountEqualsMountSum(t, got)

	// Billing scalar must equal attribution Count.
	sshScalar, err := core.GetStoredSSHDurationAdjustedCertCount(ctx, month)
	require.NoError(t, err)
	require.InDelta(t, sshScalar, vault.ToFloat64(got.Count), 1e-9,
		"SSH cert scalar (%v) must equal attribution Count (%v)", sshScalar, got.Count)
}

// TestStoreCertAttribution_SSHOTP verifies the SSH OTP attribution round-trip
// using the SSHOTPMetric storage key.
// The scalar (GetStoredSSHOTPCount) is updated with the same delta as the
// attribution blob, mirroring what ConsumeCertCounts does in production.
func TestStoreCertAttribution_SSHOTP(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	core := cluster.Cores[0].Core
	vault.TestWaitActive(t, core)

	ctx := context.Background()
	month := timeutil.StartOfMonth(time.Now().UTC())

	mount := logical.MountAttribution{
		MountAccessor:       "ssh_otp_001",
		MountPath:           "ssh/",
		MountType:           "ssh",
		NamespaceID:         "root",
		Count:               0.0014,
		MountRunningVersion: "version1",
	}

	// Mirror production: update scalar and attribution with the same delta.
	_, err := core.UpdateStoredSSHOTPCount(ctx, month, 0.0014)
	require.NoError(t, err)
	err = core.StoreCertAttribution(ctx, billing.SSHOTPMetric, 0.0014,
		map[string]logical.MountAttribution{"ssh_otp_001": mount}, month)
	require.NoError(t, err)

	got, err := core.GetStoredAttributionData(ctx, billing.LocalPrefix, month, billing.SSHOTPMetric)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, "0.0014", fmt.Sprintf("%v", got.Count))
	require.Len(t, got.Mounts, 1)
	require.Equal(t, "version1", got.Mounts["ssh_otp_001"].MountRunningVersion, "MountRunningVersion should be stored from first flush")

	// A second OTP on the same mount: count accumulates, version upgraded — confirms version is overwritten by latest flush
	mount2 := mount
	mount2.Count = 0.0014
	mount2.MountRunningVersion = "version2"
	_, err = core.UpdateStoredSSHOTPCount(ctx, month, 0.0014)
	require.NoError(t, err)
	err = core.StoreCertAttribution(ctx, billing.SSHOTPMetric, 0.0014,
		map[string]logical.MountAttribution{"ssh_otp_001": mount2}, month)
	require.NoError(t, err)

	got, err = core.GetStoredAttributionData(ctx, billing.LocalPrefix, month, billing.SSHOTPMetric)
	require.NoError(t, err)
	require.Equal(t, "0.0028", fmt.Sprintf("%v", got.Count))
	require.Equal(t, "0.0028", fmt.Sprintf("%v", got.Mounts["ssh_otp_001"].Count))
	require.Equal(t, "version2", got.Mounts["ssh_otp_001"].MountRunningVersion, "MountRunningVersion should reflect latest flush")
	requireAttrCountEqualsMountSum(t, got)

	// Billing scalar must equal attribution Count.
	otpScalar, err := core.GetStoredSSHOTPCount(ctx, month)
	require.NoError(t, err)
	require.InDelta(t, otpScalar, vault.ToFloat64(got.Count), 1e-9,
		"SSH OTP scalar (%v) must equal attribution Count (%v)", otpScalar, got.Count)
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
		MountAccessor:       "pki_consume",
		MountPath:           "pki/",
		MountType:           "pki",
		NamespaceID:         "root",
		Count:               1.0,
		MountRunningVersion: "version1",
	}
	sshMount := logical.MountAttribution{
		MountAccessor:       "ssh_consume",
		MountPath:           "ssh/",
		MountType:           "ssh",
		NamespaceID:         "root",
		Count:               0.5,
		MountRunningVersion: "version1",
	}
	otpMount := logical.MountAttribution{
		MountAccessor:       "otp_consume",
		MountPath:           "ssh/",
		MountType:           "ssh",
		NamespaceID:         "root",
		Count:               0.0014,
		MountRunningVersion: "version1",
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
	core.ConsumeCertCounts(inc, true)

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

// TestUpdateMaxKvCounts_StoresAttributionOnHWMUpdate verifies that UpdateMaxKvCounts stores
// attribution data when a new HWM is reached, and does not store attribution when the
// current count is below the stored maximum.
func TestUpdateMaxKvCounts_StoresAttributionOnHWMUpdate(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	core := cluster.Cores[0].Core
	vault.TestWaitActive(t, core)

	ctx := context.Background()
	month := timeutil.StartOfMonth(time.Now().UTC())

	attribution := vault.MountAttributionMap{
		"kv_abc123": logical.MountAttribution{
			Count:               5,
			MountAccessor:       "kv_abc123",
			MountPath:           "secret/",
			MountType:           "kv",
			NamespaceID:         "root",
			MountRunningVersion: "v1.0.0",
		},
	}

	// First call: no previous HWM — should set HWM and store attribution.
	max, err := core.UpdateMaxKvCounts(ctx, billing.ReplicatedPrefix, month, 5, attribution)
	require.NoError(t, err)
	require.Equal(t, 5, max)

	stored, err := core.GetStoredAttributionData(ctx, billing.ReplicatedPrefix, month, billing.KvHWMCountsHWM)
	require.NoError(t, err)
	require.Len(t, stored.Mounts, 1)
	require.Equal(t, "kv_abc123", stored.Mounts["kv_abc123"].MountAccessor)
	require.Equal(t, "v1.0.0", stored.Mounts["kv_abc123"].MountRunningVersion)

	// Second call with lower count — HWM must remain at 5 and attribution must not change.
	lowerAttribution := vault.MountAttributionMap{
		"kv_lower": logical.MountAttribution{
			Count:               3,
			MountAccessor:       "kv_lower",
			MountPath:           "lower/",
			MountType:           "kv",
			NamespaceID:         "root",
			MountRunningVersion: "v2.0.0",
		},
	}
	max, err = core.UpdateMaxKvCounts(ctx, billing.ReplicatedPrefix, month, 3, lowerAttribution)
	require.NoError(t, err)
	require.Equal(t, 5, max, "HWM must not decrease")

	stored, err = core.GetStoredAttributionData(ctx, billing.ReplicatedPrefix, month, billing.KvHWMCountsHWM)
	require.NoError(t, err)
	require.Len(t, stored.Mounts, 1)
	require.Equal(t, "kv_abc123", stored.Mounts["kv_abc123"].MountAccessor)
	require.Equal(t, "v1.0.0", stored.Mounts["kv_abc123"].MountRunningVersion, "attribution must not change when HWM not exceeded")

	// Third call with higher count — should update HWM and replace attribution.
	higherAttribution := vault.MountAttributionMap{
		"kv_higher": logical.MountAttribution{
			Count:               9,
			MountAccessor:       "kv_higher",
			MountPath:           "higher/",
			MountType:           "kv",
			NamespaceID:         "root",
			MountRunningVersion: "v3.0.0",
		},
	}
	max, err = core.UpdateMaxKvCounts(ctx, billing.ReplicatedPrefix, month, 9, higherAttribution)
	require.NoError(t, err)
	require.Equal(t, 9, max, "HWM must increase to new maximum")

	stored, err = core.GetStoredAttributionData(ctx, billing.ReplicatedPrefix, month, billing.KvHWMCountsHWM)
	require.NoError(t, err)
	require.Len(t, stored.Mounts, 1)
	require.Equal(t, "kv_higher", stored.Mounts["kv_higher"].MountAccessor)
	require.Equal(t, "v3.0.0", stored.Mounts["kv_higher"].MountRunningVersion, "attribution must reflect new HWM mount")
	requireAttrCountEqualsMountSum(t, stored)

	// The billing scalar (HWM) must equal the attribution Count.
	kvScalar, err := core.GetStoredHWMKvCounts(ctx, billing.ReplicatedPrefix, month)
	require.NoError(t, err)
	require.InDelta(t, float64(kvScalar), vault.ToFloat64(stored.Count), 1e-9,
		"HWM scalar (%d) must equal attribution Count (%v)", kvScalar, stored.Count)
}

// TestUpdateMaxKvCounts_NoAttributionWhenEmpty verifies that when attribution is nil or empty,
// UpdateMaxKvCounts sets the HWM correctly but does not store any attribution data.
func TestUpdateMaxKvCounts_NoAttributionWhenEmpty(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	core := cluster.Cores[0].Core
	vault.TestWaitActive(t, core)

	ctx := context.Background()
	month := timeutil.StartOfMonth(time.Now().UTC())

	max, err := core.UpdateMaxKvCounts(ctx, billing.LocalPrefix, month, 7, nil)
	require.NoError(t, err)
	require.Equal(t, 7, max)

	// Attribution entry should be absent (empty attributions → nothing stored)
	stored, err := core.GetStoredAttributionData(ctx, billing.LocalPrefix, month, billing.KvHWMCountsHWM)
	require.NoError(t, err)
	require.Empty(t, stored.Mounts, "no attribution should be stored when input attributions are empty")
}

// TestUpdateMaxRoleAndManagedKeyCounts_StoresRoleAttributionPerType verifies that when a new HWM
// is reached for a role type, the attribution is stored under the correct per-type key, and no
// attribution is stored when the HWM is not reached.
func TestUpdateMaxRoleAndManagedKeyCounts_StoresRoleAttributionPerType(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	core := cluster.Cores[0].Core
	vault.TestWaitActive(t, core)

	ctx := context.Background()
	month := timeutil.StartOfMonth(time.Now().UTC())

	awsEntry := logical.MountAttribution{
		Count:               5,
		MountAccessor:       "aws_aaa",
		MountPath:           "aws/",
		MountType:           pluginconsts.SecretEngineAWS,
		NamespaceID:         "root",
		MountRunningVersion: "v1.0.0",
	}
	roleAttribution := map[string]vault.MountAttributionMap{
		billing.AWSDynamicRoles: {"aws_aaa": awsEntry},
	}

	roleCounts := &vault.RoleCounts{AWSDynamicRoles: 5}
	managedKeyCounts := &vault.ManagedKeyCounts{}
	managedKeyAttribution := map[string]vault.MountAttributionMap{}

	_, _, err := core.UpdateMaxRoleAndManagedKeyCounts(ctx, billing.ReplicatedPrefix, month, roleCounts, managedKeyCounts, roleAttribution, managedKeyAttribution)
	require.NoError(t, err)

	// HWM updated - attribution stored under correct role type
	stored, err := core.GetStoredAttributionData(ctx, billing.ReplicatedPrefix, month, billing.RoleHWMCountsHWM+billing.AWSDynamicRoles)
	require.NoError(t, err)
	require.Len(t, stored.Mounts, 1)
	require.Equal(t, "aws_aaa", stored.Mounts["aws_aaa"].MountAccessor)
	// Count is stored as JSON and deserialised as json.Number; compare via string to avoid type mismatch.
	require.Equal(t, "5", fmt.Sprintf("%v", stored.Mounts["aws_aaa"].Count))
	require.Equal(t, "v1.0.0", stored.Mounts["aws_aaa"].MountRunningVersion)

	// Now pass a lower count with different attribution
	lowerAttribution := map[string]vault.MountAttributionMap{
		billing.AWSDynamicRoles: {
			"aws_bbb": logical.MountAttribution{Count: 3, MountAccessor: "aws_bbb", MountPath: "aws/", MountType: pluginconsts.SecretEngineAWS, NamespaceID: "root", MountRunningVersion: "v2.0.0"},
		},
	}
	lowerCounts := &vault.RoleCounts{AWSDynamicRoles: 3}
	_, _, err = core.UpdateMaxRoleAndManagedKeyCounts(ctx, billing.ReplicatedPrefix, month, lowerCounts, managedKeyCounts, lowerAttribution, managedKeyAttribution)
	require.NoError(t, err)

	// HWM not reached, attribution in storage should remain same as before
	stored, err = core.GetStoredAttributionData(ctx, billing.ReplicatedPrefix, month, billing.RoleHWMCountsHWM+billing.AWSDynamicRoles)
	require.NoError(t, err)
	require.Len(t, stored.Mounts, 1)
	require.Equal(t, "aws_aaa", stored.Mounts["aws_aaa"].MountAccessor)
	// Count is stored as JSON and deserialised as json.Number; compare via string to avoid type mismatch.
	require.Equal(t, "5", fmt.Sprintf("%v", stored.Mounts["aws_aaa"].Count))
	require.Equal(t, "v1.0.0", stored.Mounts["aws_aaa"].MountRunningVersion, "attribution must not change when HWM not exceeded")

	requireAttrCountEqualsMountSum(t, stored)

	// The billing scalar for AWSDynamicRoles must equal the attribution Count.
	roleScalars, err := core.GetStoredHWMRoleCounts(ctx, billing.ReplicatedPrefix, month)
	require.NoError(t, err)
	require.InDelta(t, float64(roleScalars.AWSDynamicRoles), vault.ToFloat64(stored.Count), 1e-9,
		"AWSDynamicRoles scalar (%d) must equal attribution Count (%v)", roleScalars.AWSDynamicRoles, stored.Count)

	// A role type that did not reach a new HWM should have no attribution stored
	storedDB, err := core.GetStoredAttributionData(ctx, billing.ReplicatedPrefix, month, billing.RoleHWMCountsHWM+billing.DatabaseDynamicRoles)
	require.NoError(t, err)
	require.Empty(t, storedDB.Mounts, "no attribution should be stored for role types at zero")
}

// TestUpdateMaxRoleAndManagedKeyCounts_TotpAttributionStoredOnHWM verifies that UpdateMaxRoleAndManagedKeyCounts
// stores attribution data when a new HWM is reached, and does not store attribution when the
// current count is below the stored maximum.
func TestUpdateMaxRoleAndManagedKeyCounts_TotpAttributionStoredOnHWM(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	core := cluster.Cores[0].Core
	vault.TestWaitActive(t, core)

	ctx := context.Background()
	month := timeutil.StartOfMonth(time.Now().UTC())

	totpEntry := logical.MountAttribution{
		Count:               4,
		MountAccessor:       "totp_t1",
		MountPath:           "totp/",
		MountType:           pluginconsts.SecretEngineTOTP,
		NamespaceID:         "root",
		MountRunningVersion: "v1.0.0",
	}
	managedKeyAttribution := map[string]vault.MountAttributionMap{
		billing.TotpKeys: {"totp_t1": totpEntry},
	}
	roleCounts := &vault.RoleCounts{}
	roleAttribution := map[string]vault.MountAttributionMap{}

	managedKeyCounts := &vault.ManagedKeyCounts{TotpKeys: 4}
	_, _, err := core.UpdateMaxRoleAndManagedKeyCounts(ctx, billing.ReplicatedPrefix, month, roleCounts, managedKeyCounts, roleAttribution, managedKeyAttribution)
	require.NoError(t, err)

	// HWM updated - attribution should be stored for TOTP
	stored, err := core.GetStoredAttributionData(ctx, billing.ReplicatedPrefix, month, billing.TotpHWMCountsHWM)
	require.NoError(t, err)
	require.Len(t, stored.Mounts, 1)
	require.Equal(t, "totp_t1", stored.Mounts["totp_t1"].MountAccessor)
	// Count is stored as JSON and deserialised as json.Number; compare via string to avoid type mismatch.
	require.Equal(t, "4", fmt.Sprintf("%v", stored.Mounts["totp_t1"].Count))
	require.Equal(t, "v1.0.0", stored.Mounts["totp_t1"].MountRunningVersion)

	// Second call with higher count: hwmUpdated is true and attribution must be stored.
	higherAttribution := map[string]vault.MountAttributionMap{
		billing.TotpKeys: {
			"totp_t2": logical.MountAttribution{
				Count:               7,
				MountAccessor:       "totp_t2",
				MountPath:           "totp2/",
				MountType:           pluginconsts.SecretEngineTOTP,
				NamespaceID:         "root",
				MountRunningVersion: "v2.0.0",
			},
		},
	}
	higherManagedKeyCounts := &vault.ManagedKeyCounts{TotpKeys: 7}
	_, _, err = core.UpdateMaxRoleAndManagedKeyCounts(ctx, billing.ReplicatedPrefix, month, roleCounts, higherManagedKeyCounts, roleAttribution, higherAttribution)
	require.NoError(t, err)

	stored, err = core.GetStoredAttributionData(ctx, billing.ReplicatedPrefix, month, billing.TotpHWMCountsHWM)
	require.NoError(t, err)
	require.Len(t, stored.Mounts, 1)
	require.Equal(t, "totp_t2", stored.Mounts["totp_t2"].MountAccessor)
	require.Equal(t, "v2.0.0", stored.Mounts["totp_t2"].MountRunningVersion, "attribution must reflect new HWM mount")

	// Third call with a lower count — HWM stays at 7 and attribution must not change.
	lowerAttribution := map[string]vault.MountAttributionMap{
		billing.TotpKeys: {
			"totp_t3": logical.MountAttribution{
				Count:               5,
				MountAccessor:       "totp_t3",
				MountPath:           "totp3/",
				MountType:           pluginconsts.SecretEngineTOTP,
				NamespaceID:         "root",
				MountRunningVersion: "v3.0.0",
			},
		},
	}
	lowerManagedKeyCounts := &vault.ManagedKeyCounts{TotpKeys: 5}
	_, _, err = core.UpdateMaxRoleAndManagedKeyCounts(ctx, billing.ReplicatedPrefix, month, roleCounts, lowerManagedKeyCounts, roleAttribution, lowerAttribution)
	require.NoError(t, err)

	stored, err = core.GetStoredAttributionData(ctx, billing.ReplicatedPrefix, month, billing.TotpHWMCountsHWM)
	require.NoError(t, err)
	require.Len(t, stored.Mounts, 1)
	require.Equal(t, "totp_t2", stored.Mounts["totp_t2"].MountAccessor)
	require.Equal(t, "v2.0.0", stored.Mounts["totp_t2"].MountRunningVersion, "attribution must not change when HWM not exceeded")
	requireAttrCountEqualsMountSum(t, stored)

	// The billing scalar (TOTP HWM) must equal the attribution Count.
	totpScalar, err := core.GetStoredHWMTotpCounts(ctx, billing.ReplicatedPrefix, month)
	require.NoError(t, err)
	require.InDelta(t, float64(totpScalar), vault.ToFloat64(stored.Count), 1e-9,
		"TOTP HWM scalar (%d) must equal attribution Count (%v)", totpScalar, stored.Count)
}

// TestUpdateMaxThirdPartyPluginCounts_StoresAttributionOnHWMUpdate verifies that when a new
// HWM is reached, attribution is stored with all correct mount metadata fields including
// MountRunningVersion, that a second call with the same count does not overwrite attribution,
// and that two mounts of the same plugin+version deduplicate to a single count.
func TestUpdateMaxThirdPartyPluginCounts_StoresAttributionOnHWMUpdate(t *testing.T) {
	t.Parallel()
	pluginDir := corehelpers.MakeTestPluginDir(t)
	cluster := minimal.NewTestSoloCluster(t, &vault.CoreConfig{
		PluginDirectory: pluginDir,
	})
	core := cluster.Cores[0].Core
	vault.TestWaitActive(t, core)
	client := cluster.Cores[0].Client

	ctx := context.Background()
	month := timeutil.StartOfMonth(time.Now().UTC())

	// Register a plugin and mount it at two paths to verify deduplication.
	secretPlugin := pluginhelpers.CompilePlugin(t, sdkconsts.PluginTypeSecrets, "v1.0.0", pluginDir)
	_, err := client.Sys().RegisterPluginDetailed(&api.RegisterPluginInput{
		Name:    secretPlugin.Name,
		Type:    api.PluginType(secretPlugin.Typ),
		Command: secretPlugin.FileName,
		SHA256:  secretPlugin.Sha256,
		Version: secretPlugin.Version,
	})
	require.NoError(t, err)
	require.NoError(t, client.Sys().Mount(secretPlugin.Name, &api.MountInput{
		Type: secretPlugin.Name,
		Config: api.MountConfigInput{
			PluginVersion: secretPlugin.Version,
		},
	}))
	// Mount the same plugin at a second path; both mounts share the same plugin+version and
	// must deduplicate to a single count.
	require.NoError(t, client.Sys().Mount("mount-b", &api.MountInput{
		Type:   secretPlugin.Name,
		Config: api.MountConfigInput{PluginVersion: secretPlugin.Version},
	}))
	// A request is needed to start the plugin process and populate RunningVersion.
	_, _ = client.Logical().Read(secretPlugin.Name + "/random")

	mounts, err := core.ListMounts()
	require.NoError(t, err)
	var mountEntry *vault.MountEntry
	for _, me := range mounts {
		if me.Path == secretPlugin.Name+"/" {
			mountEntry = me
			break
		}
	}
	require.NotNil(t, mountEntry, "mount entry not found for %q", secretPlugin.Name)

	// First call: no previous HWM — two mounts of the same plugin+version deduplicate to HWM=1.
	max, err := core.UpdateMaxThirdPartyPluginCounts(ctx, month)
	require.NoError(t, err)
	require.Equal(t, 1, max, "two mounts of the same plugin+version must deduplicate to 1")

	stored, err := core.GetStoredAttributionData(ctx, billing.LocalPrefix, month, billing.ThirdPartyPluginsPrefix)
	require.NoError(t, err)
	require.Len(t, stored.Mounts, 1)
	attr := stored.Mounts[mountEntry.Accessor]
	require.Equal(t, json.Number("1"), attr.Count)
	require.Equal(t, mountEntry.Accessor, attr.MountAccessor)
	require.Equal(t, mountEntry.Path, attr.MountPath)
	require.Equal(t, secretPlugin.Name, attr.MountType)
	require.Equal(t, "v1.0.0", attr.MountRunningVersion, "MountRunningVersion must be stored in attribution")
	require.Equal(t, "root", attr.NamespaceID)
	require.Empty(t, attr.NamespacePath)
	require.Equal(t, mountEntry.BackendAwareUUID, attr.BackendAwareUUID)

	// Second call with the same count — HWM must not change and attribution must not be overwritten.
	max, err = core.UpdateMaxThirdPartyPluginCounts(ctx, month)
	require.NoError(t, err)
	require.Equal(t, 1, max, "HWM must not change when count is the same")

	storedAgain, err := core.GetStoredAttributionData(ctx, billing.LocalPrefix, month, billing.ThirdPartyPluginsPrefix)
	require.NoError(t, err)
	require.Len(t, storedAgain.Mounts, 1, "attribution must not change when HWM is not updated")
	requireAttrCountEqualsMountSum(t, storedAgain)

	// The billing scalar (third-party plugin HWM) must equal the attribution Count.
	tpScalar, err := core.GetStoredThirdPartyPluginCounts(ctx, month)
	require.NoError(t, err)
	require.InDelta(t, float64(tpScalar), vault.ToFloat64(storedAgain.Count), 1e-9,
		"third-party plugin HWM scalar (%d) must equal attribution Count (%v)", tpScalar, storedAgain.Count)
}

// TestDeleteExpiredAttributionData_CustomRetention verifies that DeleteExpiredAttributionData
// uses the configured attribution retention period instead of the default.
func TestDeleteExpiredAttributionData_CustomRetention(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	core := cluster.Cores[0].Core
	vault.TestWaitActive(t, core)

	ctx := context.Background()
	now := time.Now().UTC()
	currentMonth := timeutil.StartOfMonth(now)

	// Configure a shorter retention: 3 months.
	customRetention := 3
	err := core.UpdateAttributionRetentionMonths(ctx, customRetention)
	require.NoError(t, err)

	// Three months ago should be deleted; two months ago should be kept.
	monthToDelete := currentMonth.AddDate(0, -customRetention, 0)
	oldestRetained := currentMonth.AddDate(0, -(customRetention - 1), 0)

	attrData := &logical.MetricTypeAttribution{
		Count:       1,
		LastUpdated: currentMonth,
		Mounts: map[string]logical.MountAttribution{
			"kv_a": {Count: 1, MountAccessor: "kv_a", MountPath: "secret/", MountType: "kv"},
		},
	}

	view, ok := core.GetBillingSubView()
	require.True(t, ok)

	for _, month := range []time.Time{monthToDelete, oldestRetained, currentMonth} {
		require.NoError(t, core.StoreAttributionData(ctx, billing.LocalPrefix, month, billing.KvHWMCountsHWM, attrData))
	}

	require.NoError(t, core.DeleteExpiredAttributionData(ctx, currentMonth))

	// monthToDelete must be gone.
	entry, err := view.Get(ctx, billing.GetAttributionMaxPath(billing.LocalPrefix, monthToDelete, billing.KvHWMCountsHWM))
	require.NoError(t, err)
	require.Nil(t, entry, "attribution older than custom retention must be deleted")

	// oldestRetained and currentMonth must be present.
	for _, month := range []time.Time{oldestRetained, currentMonth} {
		entry, err := view.Get(ctx, billing.GetAttributionMaxPath(billing.LocalPrefix, month, billing.KvHWMCountsHWM))
		require.NoError(t, err)
		require.NotNil(t, entry, "attribution within custom retention period must be kept: %s", month.Format("2006-01"))
	}
}

// TestDeleteExpiredAttributionData_ZeroRetentionWipes verifies that when attribution
// retention is configured to 0, DeleteExpiredAttributionData wipes all existing attribution
// data across all months and prefixes.
func TestDeleteExpiredAttributionData_ZeroRetentionWipes(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	core := cluster.Cores[0].Core
	vault.TestWaitActive(t, core)

	ctx := context.Background()
	now := time.Now().UTC()
	currentMonth := timeutil.StartOfMonth(now)

	view, ok := core.GetBillingSubView()
	require.True(t, ok)

	attrData := &logical.MetricTypeAttribution{
		Count:       1,
		LastUpdated: currentMonth,
		Mounts: map[string]logical.MountAttribution{
			"kv_b": {Count: 1, MountAccessor: "kv_b", MountPath: "secret/", MountType: "kv"},
		},
	}

	// Store attribution for several months.
	months := []time.Time{
		currentMonth.AddDate(0, -2, 0),
		currentMonth.AddDate(0, -1, 0),
		currentMonth,
	}
	for _, m := range months {
		require.NoError(t, core.StoreAttributionData(ctx, billing.LocalPrefix, m, billing.KvHWMCountsHWM, attrData))
	}

	// Configure retention to 0 (disable).
	err := core.UpdateAttributionRetentionMonths(ctx, billing.MinAttributionRetentionMonths)
	require.NoError(t, err)
	require.True(t, core.IsAttributionDisabled(ctx))

	require.NoError(t, core.DeleteExpiredAttributionData(ctx, currentMonth))

	// All attribution entries must be gone.
	for _, m := range months {
		entry, err := view.Get(ctx, billing.GetAttributionMaxPath(billing.LocalPrefix, m, billing.KvHWMCountsHWM))
		require.NoError(t, err)
		require.Nil(t, entry, "attribution for %s must be wiped when retention=0", m.Format("2006-01"))
	}
}

// TestAttributionDisabled_SkipsAllAttributionStorage verifies that every attribution
// write path is suppressed when attribution storage is disabled (retention = 0).
// The test covers all 8 production call sites:
//   - StoreCertAttribution: PKI, SSH cert, SSH OTP
//   - UpdateMountAttribution (via Update*Attribution): Transit, Transform, GcpKms, Spiffe, OIDC, ExternalCA
//   - UpdateMaxKvCounts
//   - UpdateMaxRoleAndManagedKeyCounts (roles + TOTP managed keys)
//   - UpdateMaxThirdPartyPluginCounts
//   - UpdateKmipEnabled
func TestAttributionDisabled_SkipsAllAttributionStorage(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	core := cluster.Cores[0].Core
	vault.TestWaitActive(t, core)

	ctx := context.Background()
	month := timeutil.StartOfMonth(time.Now().UTC())

	// Disable attribution.
	require.NoError(t, core.UpdateAttributionRetentionMonths(ctx, billing.MinAttributionRetentionMonths))

	// Helper that asserts no attribution was stored for a given metric key.
	assertEmpty := func(prefix, metricKey string) {
		t.Helper()
		got, err := core.GetStoredAttributionData(ctx, prefix, month, metricKey)
		require.NoError(t, err)
		require.Empty(t, got.Mounts, "attribution must not be stored for %s when disabled", metricKey)
	}

	mount := func(accessor, path, mountType string) logical.MountAttribution {
		return logical.MountAttribution{
			MountAccessor: accessor,
			MountPath:     path,
			MountType:     mountType,
			NamespaceID:   "root",
			Count:         1.0,
		}
	}

	// --- StoreCertAttribution (PKI, SSH cert, SSH OTP) ---
	require.NoError(t, core.StoreCertAttribution(ctx, billing.PkiDurationAdjustedCountPrefix, 1.0,
		map[string]logical.MountAttribution{"pki_a": mount("pki_a", "pki/", "pki")}, month))
	assertEmpty(billing.LocalPrefix, billing.PkiDurationAdjustedCountPrefix)

	require.NoError(t, core.StoreCertAttribution(ctx, billing.SSHCertificateMetric, 0.5,
		map[string]logical.MountAttribution{"ssh_a": mount("ssh_a", "ssh/", "ssh")}, month))
	assertEmpty(billing.LocalPrefix, billing.SSHCertificateMetric)

	require.NoError(t, core.StoreCertAttribution(ctx, billing.SSHOTPMetric, 0.0014,
		map[string]logical.MountAttribution{"otp_a": mount("otp_a", "ssh/", "ssh")}, month))
	assertEmpty(billing.LocalPrefix, billing.SSHOTPMetric)

	// --- UpdateMountAttribution (in-memory tracker → storage) ---
	// Seed each in-memory tracker directly, then call the corresponding Update*Attribution.
	// If attribution is disabled the flush must write nothing.
	cbTyped := core.GetCoreConsumptionBillingManager()
	require.NotNil(t, cbTyped)

	seedTracker := func(tracker *billing.AttributionTracker, accessor, path, mountType string) {
		tracker.MountAttributionLock.Lock()
		tracker.MountAttribution[accessor] = logical.MountAttribution{
			MountAccessor: accessor, MountPath: path, MountType: mountType,
			NamespaceID: "root", Count: 1.0,
		}
		tracker.MountAttributionLock.Unlock()
	}

	seedTracker(&cbTyped.SecretEngineCounts.Transit.AttributionTracker, "transit_a", "transit/", "transit")
	require.NoError(t, core.UpdateTransitAttribution(ctx, month))
	assertEmpty(billing.LocalPrefix, billing.TransitDataProtectionCallCountsPrefix)

	seedTracker(&cbTyped.SecretEngineCounts.Transform.AttributionTracker, "transform_a", "transform/", "transform")
	require.NoError(t, core.UpdateTransformAttribution(ctx, month))
	assertEmpty(billing.LocalPrefix, billing.TransformDataProtectionCallCountsPrefix)

	seedTracker(&cbTyped.SecretEngineCounts.GcpKms.AttributionTracker, "gcpkms_a", "gcpkms/", "gcpkms")
	require.NoError(t, core.UpdateGcpKmsAttribution(ctx, month))
	assertEmpty(billing.LocalPrefix, billing.GcpKmsDataProtectionCallCountsPrefix)

	seedTracker(&cbTyped.SecretEngineCounts.Spiffe.AttributionTracker, "spiffe_a", "spiffe/", "spiffe")
	require.NoError(t, core.UpdateSpiffeAttribution(ctx, month))
	assertEmpty(billing.LocalPrefix, billing.SpiffeJwtNormalizedTokenUnits)

	seedTracker(&cbTyped.SecretEngineCounts.Oidc.AttributionTracker, "oidc_a", "oidc/", "oidc")
	require.NoError(t, core.UpdateOidcAttribution(ctx, month))
	assertEmpty(billing.LocalPrefix, billing.OidcDurationAdjustedCountPrefix)

	seedTracker(&cbTyped.SecretEngineCounts.ExternalCa.AttributionTracker, "exca_a", "pki/", "external-ca")
	require.NoError(t, core.UpdateExternalCaAttribution(ctx, month))
	assertEmpty(billing.LocalPrefix, billing.ExternalCaDurationAdjustedCountPrefix)

	// --- UpdateMaxKvCounts ---
	_, err := core.UpdateMaxKvCounts(ctx, billing.LocalPrefix, month, 5, vault.MountAttributionMap{
		"kv_a": mount("kv_a", "secret/", "kv"),
	})
	require.NoError(t, err)
	assertEmpty(billing.LocalPrefix, billing.KvHWMCountsHWM)

	// --- UpdateMaxRoleAndManagedKeyCounts (AWS dynamic role + TOTP managed key) ---
	roleAttr := map[string]vault.MountAttributionMap{
		billing.AWSDynamicRoles: {"aws_a": mount("aws_a", "aws/", "aws")},
	}
	managedKeyAttr := map[string]vault.MountAttributionMap{
		billing.TotpKeys: {"totp_a": mount("totp_a", "totp/", "totp")},
	}
	_, _, err = core.UpdateMaxRoleAndManagedKeyCounts(ctx, billing.LocalPrefix, month,
		&vault.RoleCounts{AWSDynamicRoles: 1},
		&vault.ManagedKeyCounts{TotpKeys: 1},
		roleAttr, managedKeyAttr)
	require.NoError(t, err)
	assertEmpty(billing.LocalPrefix, billing.RoleHWMCountsHWM+billing.AWSDynamicRoles)
	assertEmpty(billing.LocalPrefix, billing.TotpHWMCountsHWM)

	// --- UpdateMaxThirdPartyPluginCounts ---
	// No real plugin mounts exist in this minimal cluster so the HWM will be 0
	// and no attribution block is entered; the call must simply not error.
	_, err = core.UpdateMaxThirdPartyPluginCounts(ctx, month)
	require.NoError(t, err)
	assertEmpty(billing.LocalPrefix, billing.ThirdPartyPluginsPrefix)

	// --- UpdateKmipEnabled ---
	// No KMIP mounts in the minimal cluster; call must not error and nothing stored.
	_, err = core.UpdateKmipEnabled(ctx, month)
	require.NoError(t, err)
	assertEmpty(billing.LocalPrefix, billing.KmipEnabledPrefix)
}
