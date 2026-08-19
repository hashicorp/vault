// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package billing

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/pluginconsts"
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
			Count:         5,
			MountAccessor: "kv_abc123",
			MountPath:     "secret/",
			MountType:     "kv",
			NamespaceID:   "root",
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

	// Second call with lower count — HWM must remain at 5 and attribution must not change.
	lowerAttribution := vault.MountAttributionMap{
		"kv_lower": logical.MountAttribution{
			Count:         3,
			MountAccessor: "kv_lower",
			MountPath:     "lower/",
			MountType:     "kv",
			NamespaceID:   "root",
		},
	}
	max, err = core.UpdateMaxKvCounts(ctx, billing.ReplicatedPrefix, month, 3, lowerAttribution)
	require.NoError(t, err)
	require.Equal(t, 5, max, "HWM must not decrease")

	stored, err = core.GetStoredAttributionData(ctx, billing.ReplicatedPrefix, month, billing.KvHWMCountsHWM)
	require.NoError(t, err)
	require.Len(t, stored.Mounts, 1)
	require.Equal(t, "kv_abc123", stored.Mounts["kv_abc123"].MountAccessor)

	// Third call with higher count — should update HWM and replace attribution.
	higherAttribution := vault.MountAttributionMap{
		"kv_higher": logical.MountAttribution{
			Count:         9,
			MountAccessor: "kv_higher",
			MountPath:     "higher/",
			MountType:     "kv",
			NamespaceID:   "root",
		},
	}
	max, err = core.UpdateMaxKvCounts(ctx, billing.ReplicatedPrefix, month, 9, higherAttribution)
	require.NoError(t, err)
	require.Equal(t, 9, max, "HWM must increase to new maximum")

	stored, err = core.GetStoredAttributionData(ctx, billing.ReplicatedPrefix, month, billing.KvHWMCountsHWM)
	require.NoError(t, err)
	require.Len(t, stored.Mounts, 1)
	require.Equal(t, "kv_higher", stored.Mounts["kv_higher"].MountAccessor)
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
		Count:         5,
		MountAccessor: "aws_aaa",
		MountPath:     "aws/",
		MountType:     pluginconsts.SecretEngineAWS,
		NamespaceID:   "root",
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

	// Now pass a lower count with different attribution
	lowerAttribution := map[string]vault.MountAttributionMap{
		billing.AWSDynamicRoles: {
			"aws_bbb": logical.MountAttribution{Count: 3, MountAccessor: "aws_bbb", MountPath: "aws/", MountType: pluginconsts.SecretEngineAWS, NamespaceID: "root"},
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
		Count:         4,
		MountAccessor: "totp_t1",
		MountPath:     "totp/",
		MountType:     pluginconsts.SecretEngineTOTP,
		NamespaceID:   "root",
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

	// Second call with higher count: hwmUpdated is true and attribution must be stored.
	higherAttribution := map[string]vault.MountAttributionMap{
		billing.TotpKeys: {
			"totp_t2": logical.MountAttribution{
				Count:         7,
				MountAccessor: "totp_t2",
				MountPath:     "totp2/",
				MountType:     pluginconsts.SecretEngineTOTP,
				NamespaceID:   "root",
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

	// Third call with a lower count — HWM stays at 7 and attribution must not change.
	lowerAttribution := map[string]vault.MountAttributionMap{
		billing.TotpKeys: {
			"totp_t3": logical.MountAttribution{
				Count:         5,
				MountAccessor: "totp_t3",
				MountPath:     "totp3/",
				MountType:     pluginconsts.SecretEngineTOTP,
				NamespaceID:   "root",
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
}
