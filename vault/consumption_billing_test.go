// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/pluginconsts"
	"github.com/hashicorp/vault/internalshared/namespace"
	"github.com/hashicorp/vault/internalshared/timeutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault/billing"
	"github.com/stretchr/testify/require"
)

var secretEngineBackends = map[string]struct {
	mount string
	key   string
}{
	"AWS Dynamic Roles": {
		mount: pluginconsts.SecretEngineAWS,
		key:   "role/",
	},
	"AWS Static Roles": {
		mount: pluginconsts.SecretEngineAWS,
		key:   "static-roles/",
	},
	"Azure Dynamic Roles": {
		mount: pluginconsts.SecretEngineAzure,
		key:   "roles/",
	},
	"Azure Static Roles": {
		mount: pluginconsts.SecretEngineAzure,
		key:   "static-roles/",
	},
	"Database Dynamic Roles": {
		mount: pluginconsts.SecretEngineDatabase,
		key:   "role/",
	},
	"Database Static Roles": {
		mount: pluginconsts.SecretEngineDatabase,
		key:   "static-role/",
	},
	"GCP Impersonated Accounts": {
		mount: pluginconsts.SecretEngineGCP,
		key:   "impersonated-account/",
	},
	"GCP Rolesets": {
		mount: pluginconsts.SecretEngineGCP,
		key:   "roleset/",
	},
	"GCP Static Accounts": {
		mount: pluginconsts.SecretEngineGCP,
		key:   "static-account/",
	},
	"LDAP Dynamic Roles": {
		mount: pluginconsts.SecretEngineLDAP,
		key:   "role/",
	},
	"LDAP Static Roles": {
		mount: pluginconsts.SecretEngineLDAP,
		key:   "static-role/",
	},
	"LDAP Library Sets": {
		mount: pluginconsts.SecretEngineLDAP,
		key:   "library/",
	},
	"OpenLDAP Dynamic Roles": {
		mount: pluginconsts.SecretEngineOpenLDAP,
		key:   "role/",
	},
	"OpenLDAP Static Roles": {
		mount: pluginconsts.SecretEngineOpenLDAP,
		key:   "static-role/",
	},
	"OpenLDAP Library Sets": {
		mount: pluginconsts.SecretEngineOpenLDAP,
		key:   "library/",
	},
	"Alicloud Dynamic Roles": {
		mount: pluginconsts.SecretEngineAlicloud,
		key:   "role/",
	},
	"RabbitMQ Dynamic Roles": {
		mount: pluginconsts.SecretEngineRabbitMQ,
		key:   "role/",
	},
	"Consul Dynamic Roles": {
		mount: pluginconsts.SecretEngineConsul,
		key:   "policy/",
	},
	"Nomad Dynamic Roles": {
		mount: pluginconsts.SecretEngineNomad,
		key:   "role/",
	},
	"Kubernetes Dynamic Roles": {
		mount: pluginconsts.SecretEngineKubernetes,
		key:   "roles/",
	},
	// MongoDB roles, unlike MongoDB Atlas roles, are
	// counted as part of the Database secret engine
	"MongoDB Atlas Dynamic Roles": {
		mount: pluginconsts.SecretEngineMongoDBAtlas,
		key:   "roles/",
	},
	"Terraform Cloud Dynamic Roles": {
		mount: pluginconsts.SecretEngineTerraform,
		key:   "role/",
	},
}

// TestConsumptionBillingMetricsWorker tests that we correctly update the consumption metrics at
// regular intervals
func TestConsumptionBillingMetricsWorker(t *testing.T) {
	coreConfig := &CoreConfig{
		LogicalBackends: roleLogicalBackends,
		BillingConfig: billing.BillingConfig{
			MetricsUpdateCadence: 5 * time.Second,
		},
	}
	core, _, root := TestCoreUnsealedWithConfig(t, coreConfig)
	for mount := range roleLogicalBackends {
		req := logical.TestRequest(t, logical.CreateOperation, fmt.Sprintf("sys/mounts/%v", mount))
		req.Data["type"] = mount
		req.ClientToken = root
		ctx := namespace.RootContext(context.Background())

		resp, err := core.HandleRequest(ctx, req)
		require.NoError(t, err)
		require.Nil(t, resp.Error())
	}

	for _, tc := range secretEngineBackends {
		addRoleToStorage(t, core, tc.mount, tc.key, 5)
	}
	timer := time.NewTimer(5 * time.Second)

	_ = <-timer.C
	// Check that the billing metrics have been updated
	counts, err := core.GetStoredHWMRoleCounts(context.Background(), billing.ReplicatedPrefix, time.Now())
	require.NoError(t, err)

	verifyExpectedRoleCounts(t, counts, 5)

	for _, tc := range secretEngineBackends {
		deleteAllRolesFromStorage(t, core, tc.mount, tc.key)
		addRoleToStorage(t, core, tc.mount, tc.key, 3)
	}

	timer.Reset(5 * time.Second)

	_ = <-timer.C
	// Check that the billing metrics have been updated
	counts, err = core.GetStoredHWMRoleCounts(context.Background(), billing.ReplicatedPrefix, time.Now())
	require.NoError(t, err)

	verifyExpectedRoleCounts(t, counts, 5)
}

// TestHandleEndOfMonthMetrics tests that HandleEndOfMonth cleans up
// billing metrics from billing.DefaultBillingRetentionMonths ago (keeping billing.DefaultBillingRetentionMonths of data) and resets the in memory billing metrics
func TestHandleEndOfMonthMetrics(t *testing.T) {
	coreConfig := &CoreConfig{
		LogicalBackends: roleLogicalBackends,
		BillingConfig: billing.BillingConfig{
			MetricsUpdateCadence: 3 * time.Second,
		},
	}
	core, _, _ := TestCoreUnsealedWithConfig(t, coreConfig)
	// Add some billing metrics to storage for (billing.DefaultBillingRetentionMonths - 1) and billing.DefaultBillingRetentionMonths months ago
	// Use the util functions directly to avoid the need to mount the logical backends
	now := time.Now().UTC()
	oldestRetainedMonth := timeutil.StartOfMonth(now).AddDate(0, -(billing.DefaultBillingRetentionMonths - 1), 0)
	monthToDelete := timeutil.StartOfMonth(now).AddDate(0, -billing.DefaultBillingRetentionMonths, 0)

	for _, month := range []time.Time{monthToDelete, oldestRetainedMonth} {
		for _, localPathPrefix := range []string{billing.ReplicatedPrefix, billing.LocalPrefix} {
			core.storeMaxRoleCountsLocked(context.Background(), &RoleCounts{
				AWSDynamicRoles:      10,
				AWSStaticRoles:       15,
				LDAPDynamicRoles:     8,
				GCPRolesets:          3,
				DatabaseDynamicRoles: 5,
				DatabaseStaticRoles:  7,
				OSLocalAccountRoles:  9,
			}, localPathPrefix, month)
			core.storeMaxKvCountsLocked(context.Background(), 10, localPathPrefix, month)

			// Add KV HWM attribution data
			err := core.StoreAttributionData(context.Background(), localPathPrefix, month, billing.KvHWMCountsHWM, &logical.MetricTypeAttribution{
				Count:       10,
				LastUpdated: month,
				Mounts: map[string]logical.MountAttribution{
					"kv_accessor_1": {
						MountAccessor: "kv_accessor_1",
						MountPath:     "secret/",
						MountType:     "kv",
						NamespaceID:   "root",
						NamespacePath: "",
						Count:         10,
					},
				},
			})
			require.NoError(t, err)

			// Transit, third-party plugins, ssh credential count and OIDC are local aggregated metrics
			// and should only be stored under LocalPrefix
			if localPathPrefix == billing.LocalPrefix {
				core.storeTransitCallCountsLocked(context.Background(), 10, localPathPrefix, month)
				core.storeGcpKmsCallCountsLocked(context.Background(), 10, localPathPrefix, month)
				core.storeThirdPartyPluginCountsLocked(context.Background(), localPathPrefix, month, 10)
				core.storeOidcDurationAdjustedCountLocked(context.Background(), month, 10)
				core.storeSSHOTPCountLocked(context.Background(), localPathPrefix, month, 10)

				// Add transit attribution data (local-only metric)
				err = core.StoreAttributionData(context.Background(), localPathPrefix, month, billing.TransitDataProtectionCallCountsPrefix, &logical.MetricTypeAttribution{
					Count:       10,
					LastUpdated: month,
					Mounts: map[string]logical.MountAttribution{
						"transit_accessor_1": {
							MountAccessor: "transit_accessor_1",
							MountPath:     "transit/",
							MountType:     "transit",
							NamespaceID:   "root",
							NamespacePath: "",
							Count:         10,
						},
					},
				})
				require.NoError(t, err)
			}

			// List the data paths to verify that the billing metrics have been stored
			view, ok := core.GetBillingSubView()
			require.True(t, ok)
			paths, err := view.List(context.Background(), billing.GetMonthlyBillingPath(localPathPrefix, month))
			require.NoError(t, err)
			expectedPaths := 3 // ReplicatedPrefix has roles, kv, and kv attribution (attribution/ is one directory entry)
			if localPathPrefix == billing.LocalPrefix {
				expectedPaths = 8 // LocalPrefix has roles, kv, transit, gcp kms, third-party plugins, ssh, OIDC, and attribution/ (one directory entry for both kv and transit attribution)
			}
			require.Equal(t, expectedPaths, len(paths))
		}
	}

	// Handle the end of the month
	core.HandleStartOfMonth(context.Background(), now)

	for _, localPathPrefix := range []string{billing.ReplicatedPrefix, billing.LocalPrefix} {
		// billing.DefaultBillingRetentionMonths ago should have no billing metrics (deleted)
		view, ok := core.GetBillingSubView()
		require.True(t, ok)
		paths, err := view.List(context.Background(), billing.GetMonthlyBillingPath(localPathPrefix, monthToDelete))
		require.NoError(t, err)
		require.Equal(t, 0, len(paths), "data from billing.DefaultBillingRetentionMonths ago should be deleted")

		// Attribution for monthToDelete must also be gone — HandleStartOfMonth calls
		// DeleteExpiredAttributionData which uses the DefaultAttributionRetentionMonths, which has the same duration window as DefaultBillingRetentionMonths.
		deletedKvAttr, err := core.GetStoredAttributionData(context.Background(), localPathPrefix, monthToDelete, billing.KvHWMCountsHWM)
		require.NoError(t, err)
		require.Empty(t, deletedKvAttr.Mounts, "KV HWM attribution for monthToDelete should be deleted after HandleStartOfMonth")

		// (billing.DefaultBillingRetentionMonths - 1) months ago should still have the billing metrics (kept)
		view, ok = core.GetBillingSubView()
		require.True(t, ok)
		paths, err = view.List(context.Background(), billing.GetMonthlyBillingPath(localPathPrefix, oldestRetainedMonth))
		require.NoError(t, err)
		expectedPaths := 3 // ReplicatedPrefix has roles, kv, and attribution/ (one directory entry)
		if localPathPrefix == billing.LocalPrefix {
			expectedPaths = 8 // LocalPrefix has roles, kv, transit, gcp kms, third-party plugins, ssh, OIDC, and attribution/ (one directory entry)
		}
		require.Equal(t, expectedPaths, len(paths))

		// Verify retained KV HWM attribution has the expected timestamp, namespace, and mount fields.
		kvAttr, err := core.GetStoredAttributionData(context.Background(), localPathPrefix, oldestRetainedMonth, billing.KvHWMCountsHWM)
		require.NoError(t, err)
		require.NotNil(t, kvAttr)
		require.False(t, kvAttr.LastUpdated.IsZero(), "KV HWM attribution LastUpdated must not be zero for retained month")
		require.Equal(t, oldestRetainedMonth.UTC().Truncate(time.Second), kvAttr.LastUpdated.UTC().Truncate(time.Second),
			"KV HWM attribution LastUpdated should match the month it was stored for")
		require.Len(t, kvAttr.Mounts, 1)
		kvMount, ok := kvAttr.Mounts["kv_accessor_1"]
		require.True(t, ok)
		require.Equal(t, "kv_accessor_1", kvMount.MountAccessor)
		require.Equal(t, "secret/", kvMount.MountPath)
		require.Equal(t, "kv", kvMount.MountType)
		require.Equal(t, "root", kvMount.NamespaceID)
		require.Equal(t, "", kvMount.NamespacePath)
		require.Equal(t, "10", fmt.Sprintf("%v", kvMount.Count))

		// Transit attribution is local-only; verify timestamp, namespace, and mount fields for LocalPrefix.
		if localPathPrefix == billing.LocalPrefix {
			// Verify transit attribution for monthToDelete is also gone.
			deletedTransitAttr, err := core.GetStoredAttributionData(context.Background(), localPathPrefix, monthToDelete, billing.TransitDataProtectionCallCountsPrefix)
			require.NoError(t, err)
			require.Empty(t, deletedTransitAttr.Mounts, "transit attribution for monthToDelete should be deleted after HandleStartOfMonth")

			transitAttr, err := core.GetStoredAttributionData(context.Background(), localPathPrefix, oldestRetainedMonth, billing.TransitDataProtectionCallCountsPrefix)
			require.NoError(t, err)
			require.NotNil(t, transitAttr)
			require.False(t, transitAttr.LastUpdated.IsZero(), "transit attribution LastUpdated must not be zero for retained month")
			require.Equal(t, oldestRetainedMonth.UTC().Truncate(time.Second), transitAttr.LastUpdated.UTC().Truncate(time.Second),
				"transit attribution LastUpdated should match the month it was stored for")
			require.Len(t, transitAttr.Mounts, 1)
			transitMount, ok := transitAttr.Mounts["transit_accessor_1"]
			require.True(t, ok)
			require.Equal(t, "transit_accessor_1", transitMount.MountAccessor)
			require.Equal(t, "transit/", transitMount.MountPath)
			require.Equal(t, "transit", transitMount.MountType)
			require.Equal(t, "root", transitMount.NamespaceID)
			require.Equal(t, "", transitMount.NamespacePath)
			require.Equal(t, "10", fmt.Sprintf("%v", transitMount.Count))
		}
	}

	require.Equal(t, uint64(0), core.GetInMemoryTransitDataProtectionCallCounts())
	require.Equal(t, uint64(0), core.GetInMemoryTransformDataProtectionCallCounts())
	require.Equal(t, uint64(0), core.GetInMemoryGcpKmsDataProtectionCallCounts())
	require.Equal(t, float64(0), core.GetInMemoryOidcCounts())
	require.False(t, core.consumptionBilling.KmipSeenEnabledThisMonth.Load())
	require.Equal(t, 0, len(core.GetInMemoryTransitAttribution()))
	require.Equal(t, 0, len(core.GetInMemoryOidcAttribution()))
	require.Equal(t, 0, len(core.GetInMemoryExternalCaAttribution()))
}

// TestDeleteExpiredBillingMetrics specifically tests the DeleteExpiredBillingMetrics method
// to ensure it correctly deletes data from billing.DefaultBillingRetentionMonths ago while keeping
// data from (billing.DefaultBillingRetentionMonths - 1) months ago.
func TestDeleteExpiredBillingMetrics(t *testing.T) {
	coreConfig := &CoreConfig{
		LogicalBackends: roleLogicalBackends,
	}
	core, _, _ := TestCoreUnsealedWithConfig(t, coreConfig)

	now := time.Now().UTC()
	currentMonth := timeutil.StartOfMonth(now)
	oldestRetainedMonth := currentMonth.AddDate(0, -(billing.DefaultBillingRetentionMonths - 1), 0)
	monthToDelete := currentMonth.AddDate(0, -billing.DefaultBillingRetentionMonths, 0)

	// Write billing data for multiple months including the month to be deleted and the oldest retained month
	for _, month := range []time.Time{monthToDelete, oldestRetainedMonth, currentMonth} {
		for _, pathPrefix := range []string{billing.ReplicatedPrefix, billing.LocalPrefix} {
			core.storeMaxRoleCountsLocked(context.Background(), &RoleCounts{
				AWSDynamicRoles:     5,
				AWSStaticRoles:      10,
				LDAPDynamicRoles:    3,
				OSLocalAccountRoles: 7,
			}, pathPrefix, month)
			core.storeMaxKvCountsLocked(context.Background(), 20, pathPrefix, month)
			core.storeTransitCallCountsLocked(context.Background(), 15, pathPrefix, month)
			// Add SSH metrics which use subdirectory paths (ssh/normalized-certs-issued, ssh/credential-count)
			core.storeSSHDurationAdjustedCertCountLocked(context.Background(), pathPrefix, month, 10.5)
			core.storeSSHOTPCountLocked(context.Background(), pathPrefix, month, 25.0)

			// Add transit mount/namespace attribution data
			err := core.StoreAttributionData(context.Background(), pathPrefix, month, billing.TransitDataProtectionCallCountsPrefix, &logical.MetricTypeAttribution{
				Count: 10.0,
				Mounts: map[string]logical.MountAttribution{
					"transit_accessor_1": {
						MountAccessor: "transit_accessor_1",
						MountPath:     "transit1/",
						NamespaceID:   "root",
						NamespacePath: "",
						Count:         10.0,
					},
				},
			})
			require.NoError(t, err)

			// Add external CA mount/namespace attribution data
			err = core.StoreAttributionData(context.Background(), pathPrefix, month, billing.ExternalCaDurationAdjustedCountPrefix, &logical.MetricTypeAttribution{
				Count: 9.5,
				Mounts: map[string]logical.MountAttribution{
					"external_ca_accessor_1": {
						MountAccessor: "external_ca_accessor_1",
						MountPath:     "pki-ext1/",
						NamespaceID:   "root",
						NamespacePath: "",
						Count:         9.5,
					},
				},
			})
			require.NoError(t, err)

			// Add OIDC mount/namespace attribution data
			err = core.StoreAttributionData(context.Background(), pathPrefix, month, billing.OidcDurationAdjustedCountPrefix, &logical.MetricTypeAttribution{
				Count: 5.0,
				Mounts: map[string]logical.MountAttribution{
					"identity_accessor_1": {
						MountAccessor: "identity_accessor_1",
						MountPath:     "identity/",
						MountType:     "identity",
						NamespaceID:   "root",
						NamespacePath: "",
						Count:         5.0,
					},
				},
			})
			require.NoError(t, err)

			// Add KV HWM mount/namespace attribution data
			err = core.StoreAttributionData(context.Background(), pathPrefix, month, billing.KvHWMCountsHWM, &logical.MetricTypeAttribution{
				Count:       20,
				LastUpdated: month,
				Mounts: map[string]logical.MountAttribution{
					"kv_accessor_1": {
						MountAccessor: "kv_accessor_1",
						MountPath:     "secret/",
						MountType:     "kv",
						NamespaceID:   "root",
						NamespacePath: "",
						Count:         20,
					},
				},
			})
			require.NoError(t, err)
		}
		// Store updatedAtTimestamp for each month
		testUpdateTime := time.Date(month.Year(), month.Month(), 15, 12, 0, 0, 0, time.UTC)
		err := core.UpdateMetricsLastUpdateTime(context.Background(), month, testUpdateTime)
		require.NoError(t, err)
	}

	// Verify data exists before deletion
	for _, pathPrefix := range []string{billing.ReplicatedPrefix, billing.LocalPrefix} {
		view, ok := core.GetBillingSubView()
		require.True(t, ok)

		// Check month to be deleted has data
		paths, err := view.List(context.Background(), billing.GetMonthlyBillingPath(pathPrefix, monthToDelete))
		require.NoError(t, err)
		require.Greater(t, len(paths), 0, "month to delete should have data before deletion")

		// Check oldest retained month has data
		paths, err = view.List(context.Background(), billing.GetMonthlyBillingPath(pathPrefix, oldestRetainedMonth))
		require.NoError(t, err)
		require.Greater(t, len(paths), 0, "oldest retained month should have data")

		// Verify SSH metrics exist (they use subdirectory paths)
		sshCertPath := billing.GetMonthlyBillingMetricPath(pathPrefix, monthToDelete, billing.SSHCertificateMetric)
		entry, err := view.Get(context.Background(), sshCertPath)
		require.NoError(t, err)
		require.NotNil(t, entry, "SSH cert metric should exist before deletion")

		sshOTPPath := billing.GetMonthlyBillingMetricPath(pathPrefix, monthToDelete, billing.SSHOTPMetric)
		entry, err = view.Get(context.Background(), sshOTPPath)
		require.NoError(t, err)
		require.NotNil(t, entry, "SSH OTP metric should exist before deletion")

		// Verify transit mount breakdown exists
		transitBreakdownPath := billing.GetAttributionMaxPath(pathPrefix, monthToDelete, billing.TransitDataProtectionCallCountsPrefix)
		entry, err = view.Get(context.Background(), transitBreakdownPath)
		require.NoError(t, err)
		require.NotNil(t, entry, "Transit mount breakdown should exist before deletion")

		// Verify OIDC mount breakdown exists
		oidcBreakdownPath := billing.GetAttributionMaxPath(pathPrefix, monthToDelete, billing.OidcDurationAdjustedCountPrefix)
		entry, err = view.Get(context.Background(), oidcBreakdownPath)
		require.NoError(t, err)
		require.NotNil(t, entry, "OIDC mount breakdown should exist before deletion")

		// Verify external CA mount breakdown exists
		externalCaBreakdownPath := billing.GetAttributionMaxPath(pathPrefix, monthToDelete, billing.ExternalCaDurationAdjustedCountPrefix)
		entry, err = view.Get(context.Background(), externalCaBreakdownPath)
		require.NoError(t, err)
		require.NotNil(t, entry, "External CA mount breakdown should exist before deletion")

		// Verify KV HWM mount breakdown exists
		kvBreakdownPath := billing.GetAttributionMaxPath(pathPrefix, monthToDelete, billing.KvHWMCountsHWM)
		entry, err = view.Get(context.Background(), kvBreakdownPath)
		require.NoError(t, err)
		require.NotNil(t, entry, "KV HWM mount breakdown should exist before deletion")
	}

	// Verify updatedAtTimestamp exists for all months before deletion
	for _, month := range []time.Time{monthToDelete, oldestRetainedMonth, currentMonth} {
		timestamp, err := core.GetMetricsLastUpdateTime(context.Background(), month)
		require.NoError(t, err)
		require.False(t, timestamp.IsZero(), "timestamp for month %s should exist before deletion", month.Format("2006-01"))
	}

	// Call DeleteExpiredBillingMetrics directly
	err := core.DeleteExpiredBillingMetrics(context.Background(), currentMonth)
	require.NoError(t, err)

	// Verify deletion results
	for _, pathPrefix := range []string{billing.ReplicatedPrefix, billing.LocalPrefix} {
		view, ok := core.GetBillingSubView()
		require.True(t, ok)

		// Month to delete should have no regular billing data left.
		// Attribution data has its own independent retention and is intentionally skipped
		// by deleteExpiredBillingMetrics, so the attribution/ directory may still be present.
		paths, err := view.List(context.Background(), billing.GetMonthlyBillingPath(pathPrefix, monthToDelete))
		require.NoError(t, err)
		for _, p := range paths {
			require.True(t, strings.HasPrefix(p, "attribution/"), "only attribution/ entries should survive billing deletion for month %s, got unexpected path: %s", monthToDelete.Format("2006-01"), p)
		}

		// Verify SSH metrics are deleted (they use subdirectory paths)
		sshCertPath := billing.GetMonthlyBillingMetricPath(pathPrefix, monthToDelete, billing.SSHCertificateMetric)
		entry, err := view.Get(context.Background(), sshCertPath)
		require.NoError(t, err)
		require.Nil(t, entry, "SSH cert metric should be deleted")

		sshOTPPath := billing.GetMonthlyBillingMetricPath(pathPrefix, monthToDelete, billing.SSHOTPMetric)
		entry, err = view.Get(context.Background(), sshOTPPath)
		require.NoError(t, err)
		require.Nil(t, entry, "SSH OTP metric should be deleted")

		// Verify transit attribution data is NOT deleted by deleteExpiredBillingMetrics —
		// attribution has its own independent retention policy.
		transitBreakdownPath := billing.GetAttributionMaxPath(pathPrefix, monthToDelete, billing.TransitDataProtectionCallCountsPrefix)
		entry, err = view.Get(context.Background(), transitBreakdownPath)
		require.NoError(t, err)
		require.NotNil(t, entry, "Transit attribution should NOT be deleted by deleteExpiredBillingMetrics (independent retention)")

		// Verify OIDC attribution data is also NOT deleted by deleteExpiredBillingMetrics
		oidcBreakdownPath := billing.GetAttributionMaxPath(pathPrefix, monthToDelete, billing.OidcDurationAdjustedCountPrefix)
		entry, err = view.Get(context.Background(), oidcBreakdownPath)
		require.NoError(t, err)
		require.NotNil(t, entry, "OIDC attribution should NOT be deleted by deleteExpiredBillingMetrics (independent retention)")

		// Verify external CA attribution data is NOT deleted by deleteExpiredBillingMetrics
		externalCaBreakdownPath := billing.GetAttributionMaxPath(pathPrefix, monthToDelete, billing.ExternalCaDurationAdjustedCountPrefix)
		entry, err = view.Get(context.Background(), externalCaBreakdownPath)
		require.NoError(t, err)
		require.NotNil(t, entry, "ExternalCA attribution should NOT be deleted by deleteExpiredBillingMetrics (independent retention)")

		// Verify KV HWM attribution data is NOT deleted by deleteExpiredBillingMetrics
		kvBreakdownPath := billing.GetAttributionMaxPath(pathPrefix, monthToDelete, billing.KvHWMCountsHWM)
		entry, err = view.Get(context.Background(), kvBreakdownPath)
		require.NoError(t, err)
		require.NotNil(t, entry, "KV HWM attribution should NOT be deleted by deleteExpiredBillingMetrics (independent retention)")

		// Oldest retained month should still have data
		paths, err = view.List(context.Background(), billing.GetMonthlyBillingPath(pathPrefix, oldestRetainedMonth))
		require.NoError(t, err)
		require.Greater(t, len(paths), 0, "data from (billing.DefaultBillingRetentionMonths - 1) months ago should be kept")

		// Verify SSH metrics are kept for oldest retained month
		sshCertPath = billing.GetMonthlyBillingMetricPath(pathPrefix, oldestRetainedMonth, billing.SSHCertificateMetric)
		entry, err = view.Get(context.Background(), sshCertPath)
		require.NoError(t, err)
		require.NotNil(t, entry, "SSH cert metric should be kept for oldest retained month")

		sshOTPPath = billing.GetMonthlyBillingMetricPath(pathPrefix, oldestRetainedMonth, billing.SSHOTPMetric)
		entry, err = view.Get(context.Background(), sshOTPPath)
		require.NoError(t, err)
		require.NotNil(t, entry, "SSH OTP metric should be kept for oldest retained month")

		// Verify transit mount breakdown is kept for oldest retained month
		transitBreakdownPath = billing.GetAttributionMaxPath(pathPrefix, oldestRetainedMonth, billing.TransitDataProtectionCallCountsPrefix)
		entry, err = view.Get(context.Background(), transitBreakdownPath)
		require.NoError(t, err)
		require.NotNil(t, entry, "Transit mount breakdown should be kept for oldest retained month")

		// Verify OIDC mount breakdown is kept for oldest retained month
		oidcBreakdownPath = billing.GetAttributionMaxPath(pathPrefix, oldestRetainedMonth, billing.OidcDurationAdjustedCountPrefix)
		entry, err = view.Get(context.Background(), oidcBreakdownPath)
		require.NoError(t, err)
		require.NotNil(t, entry, "OIDC mount breakdown should be kept for oldest retained month")

		// Verify external CA mount breakdown is kept for oldest retained month
		externalCaBreakdownPath = billing.GetAttributionMaxPath(pathPrefix, oldestRetainedMonth, billing.ExternalCaDurationAdjustedCountPrefix)
		entry, err = view.Get(context.Background(), externalCaBreakdownPath)
		require.NoError(t, err)
		require.NotNil(t, entry, "External CA mount breakdown should be kept for oldest retained month")

		// Verify KV HWM mount breakdown is kept for oldest retained month
		kvBreakdownPath = billing.GetAttributionMaxPath(pathPrefix, oldestRetainedMonth, billing.KvHWMCountsHWM)
		entry, err = view.Get(context.Background(), kvBreakdownPath)
		require.NoError(t, err)
		require.NotNil(t, entry, "KV HWM mount breakdown should be kept for oldest retained month")

		// Current month should still have data
		paths, err = view.List(context.Background(), billing.GetMonthlyBillingPath(pathPrefix, currentMonth))
		require.NoError(t, err)
		require.Greater(t, len(paths), 0, "current month data should be kept")
	}

	// Verify updatedAtTimestamp deletion
	// Month to delete should have zero timestamp
	deletedTimestamp, err := core.GetMetricsLastUpdateTime(context.Background(), monthToDelete)
	require.NoError(t, err)
	require.True(t, deletedTimestamp.IsZero(), "timestamp for deleted month should be zero")

	// Oldest retained month should still have timestamp
	oldestTimestamp, err := core.GetMetricsLastUpdateTime(context.Background(), oldestRetainedMonth)
	require.NoError(t, err)
	require.False(t, oldestTimestamp.IsZero(), "timestamp for oldest retained month should exist")

	// Current month should still have timestamp
	currentTimestamp, err := core.GetMetricsLastUpdateTime(context.Background(), currentMonth)
	require.NoError(t, err)
	require.False(t, currentTimestamp.IsZero(), "timestamp for current month should exist")
}

// TestConsumptionBillingMetricsWorkerWithCustomClock tests that we correctly delete data older than billing.DefaultBillingRetentionMonths
// and reset the in memory billing metrics when the clock is overridden for testing purposes.
// It also verifies that mount/namespace attribution accumulated during the ending month is flushed
// to storage with the correct timestamps, namespace IDs/paths, and mount fields at the month cutover,
// and that in-memory attribution maps are cleared after the boundary is crossed.
func TestConsumptionBillingMetricsWorkerWithCustomClock(t *testing.T) {
	// 10 seconds until a new month (leave buffer for require.Eventually timeout)
	now := time.Date(2021, 1, 31, 23, 59, 50, 0, time.UTC)
	coreConfig := &CoreConfig{
		LogicalBackends: roleLogicalBackends,
		BillingConfig: billing.BillingConfig{
			TestOverrideClock: newMockTimeNowClock(now),
		},
	}
	core, _, _ := TestCoreUnsealedWithConfig(t, coreConfig)

	// Add some billing metrics to storage for (billing.DefaultBillingRetentionMonths - 1) and billing.DefaultBillingRetentionMonths months ago
	// Use the util functions directly to avoid the need to mount the logical backends
	// The worker's "end of month" path calls HandleEndOfMonth with the *current* month,
	// which will be the next month once we cross the boundary. So the months should be
	// calculated relative to that boundary.
	currentMonthAtBoundary := timeutil.StartOfNextMonth(now)
	endingMonth := timeutil.StartOfMonth(now)
	oldestRetainedMonth := timeutil.StartOfMonth(currentMonthAtBoundary).AddDate(0, -(billing.DefaultBillingRetentionMonths - 1), 0)
	monthToDelete := timeutil.StartOfMonth(currentMonthAtBoundary).AddDate(0, -billing.DefaultBillingRetentionMonths, 0)
	view, ok := core.GetBillingSubView()
	require.True(t, ok)
	roleCounts := &RoleCounts{
		AWSDynamicRoles:            10,
		AWSStaticRoles:             15,
		AzureDynamicRoles:          10,
		AzureStaticRoles:           15,
		DatabaseDynamicRoles:       5,
		DatabaseStaticRoles:        7,
		LDAPDynamicRoles:           8,
		LDAPStaticRoles:            10,
		OpenLDAPDynamicRoles:       5,
		OpenLDAPStaticRoles:        7,
		AlicloudDynamicRoles:       10,
		RabbitMQDynamicRoles:       5,
		ConsulDynamicRoles:         7,
		NomadDynamicRoles:          10,
		KubernetesDynamicRoles:     5,
		MongoDBAtlasDynamicRoles:   7,
		TerraformCloudDynamicRoles: 10,
		OSLocalAccountRoles:        11,
	}

	verifyMonthlyBillingMetrics := func(month time.Time, localPathPrefix string) {
		maxRoleCounts, err := core.GetStoredHWMRoleCounts(context.Background(), localPathPrefix, month)
		require.NoError(t, err)
		require.Equal(t, roleCounts, maxRoleCounts)
		kvCounts, err := core.GetStoredHWMKvCounts(context.Background(), localPathPrefix, month)
		require.NoError(t, err)
		require.Equal(t, 10, kvCounts)
		if localPathPrefix == billing.LocalPrefix {
			transitCounts, err := core.GetStoredTransitCallCounts(context.Background(), month)
			require.NoError(t, err)
			require.Equal(t, uint64(10), transitCounts)
			gcpKmsCounts, err := core.GetStoredGcpKmsCallCounts(context.Background(), month)
			require.NoError(t, err)
			require.Equal(t, uint64(10), gcpKmsCounts)
			thirdPartyPluginCounts, err := core.GetStoredThirdPartyPluginCounts(context.Background(), month)
			require.NoError(t, err)
			require.Equal(t, 10, thirdPartyPluginCounts)
		}
	}

	for _, month := range []time.Time{monthToDelete, oldestRetainedMonth} {

		core.storeTransitCallCountsLocked(context.Background(), uint64(10), billing.LocalPrefix, month)
		core.storeGcpKmsCallCountsLocked(context.Background(), uint64(10), billing.LocalPrefix, month)
		core.storeThirdPartyPluginCountsLocked(context.Background(), billing.LocalPrefix, month, 10)

		for _, localPathPrefix := range []string{billing.ReplicatedPrefix, billing.LocalPrefix} {
			core.storeMaxRoleCountsLocked(context.Background(), roleCounts, localPathPrefix, month)
			core.storeMaxKvCountsLocked(context.Background(), 10, localPathPrefix, month)
		}

		verifyMonthlyBillingMetrics(month, billing.LocalPrefix)
		verifyMonthlyBillingMetrics(month, billing.ReplicatedPrefix)
	}

	// Seed in-memory transit attribution for the ending month (January 2021).
	// The periodic worker will flush these to storage before HandleStartOfMonth clears them.
	cb := core.GetCoreConsumptionBillingManager()
	require.NotNil(t, cb)

	cb.SecretEngineCounts.Transit.MountAttributionLock.Lock()
	cb.SecretEngineCounts.Transit.MountAttribution["transit_accessor_1"] = logical.MountAttribution{
		MountAccessor: "transit_accessor_1",
		MountPath:     "transit/",
		MountType:     "transit",
		NamespaceID:   "root",
		NamespacePath: "",
		Count:         float64(42),
	}
	cb.SecretEngineCounts.Transit.MountAttributionLock.Unlock()

	// Seed KV HWM attribution for the ending month (January 2021) directly, since KV has no in-memory tracker.
	for _, localPathPrefix := range []string{billing.ReplicatedPrefix, billing.LocalPrefix} {
		err := core.StoreAttributionData(context.Background(), localPathPrefix, endingMonth, billing.KvHWMCountsHWM, &logical.MetricTypeAttribution{
			Count:       10,
			LastUpdated: endingMonth,
			Mounts: map[string]logical.MountAttribution{
				"kv_accessor_1": {
					MountAccessor: "kv_accessor_1",
					MountPath:     "secret/",
					MountType:     "kv",
					NamespaceID:   "root",
					NamespacePath: "",
					Count:         10,
				},
			},
		})
		require.NoError(t, err)
	}

	for _, localPathPrefix := range []string{billing.ReplicatedPrefix, billing.LocalPrefix} {
		// billing.DefaultBillingRetentionMonths ago should eventually have no billing metrics (deleted)
		require.Eventually(t, func() bool {
			paths, err := view.List(context.Background(), billing.GetMonthlyBillingPath(localPathPrefix, monthToDelete))
			return err == nil && len(paths) == 0
		}, 20*time.Second, 100*time.Millisecond)

		// All values from billing.DefaultBillingRetentionMonths ago should be 0
		maxRoleCounts, _ := core.GetStoredHWMRoleCounts(context.Background(), localPathPrefix, monthToDelete)
		require.Equal(t, &RoleCounts{}, maxRoleCounts)
		kvCounts, _ := core.GetStoredHWMKvCounts(context.Background(), localPathPrefix, monthToDelete)
		require.Equal(t, 0, kvCounts)
		if localPathPrefix == billing.LocalPrefix {
			transitCounts, _ := core.GetStoredTransitCallCounts(context.Background(), monthToDelete)
			require.Equal(t, uint64(0), transitCounts)
			gcpKmsCounts, _ := core.GetStoredGcpKmsCallCounts(context.Background(), monthToDelete)
			require.Equal(t, uint64(0), gcpKmsCounts)
			thirdPartyPluginCounts, _ := core.GetStoredThirdPartyPluginCounts(context.Background(), monthToDelete)
			require.Equal(t, 0, thirdPartyPluginCounts)
		}

		// (billing.DefaultBillingRetentionMonths - 1) months ago should still have the billing metrics (kept)
		verifyMonthlyBillingMetrics(oldestRetainedMonth, localPathPrefix)
	}

	// After the cutover, verify that transit attribution for the ending month was flushed to
	// storage with the correct timestamp, namespace, and mount fields.
	transitAttr, err := core.GetStoredAttributionData(context.Background(), billing.LocalPrefix, endingMonth, billing.TransitDataProtectionCallCountsPrefix)
	require.NoError(t, err)
	require.NotNil(t, transitAttr, "transit attribution should be stored for the ending month after cutover")
	require.False(t, transitAttr.LastUpdated.IsZero(), "transit attribution LastUpdated must not be zero after cutover")
	require.Equal(t, endingMonth.UTC(), transitAttr.LastUpdated.UTC(),
		"transit attribution LastUpdated should equal the ending month timestamp")
	require.Len(t, transitAttr.Mounts, 1, "transit attribution should have one mount entry after cutover")
	transitMount, ok := transitAttr.Mounts["transit_accessor_1"]
	require.True(t, ok, "transit attribution should be keyed by transit_accessor_1")
	require.Equal(t, "transit_accessor_1", transitMount.MountAccessor)
	require.Equal(t, "transit/", transitMount.MountPath)
	require.Equal(t, "transit", transitMount.MountType)
	require.Equal(t, "root", transitMount.NamespaceID)
	require.Equal(t, "", transitMount.NamespacePath)

	// Verify in-memory transit attribution is cleared after the cutover.
	require.Empty(t, core.GetInMemoryTransitAttribution(), "transit in-memory attribution should be cleared after cutover")

	// Verify KV HWM attribution for the ending month is present under both prefixes.
	// KV attribution is stored directly (no in-memory tracker), so it should survive the cutover untouched.
	for _, localPathPrefix := range []string{billing.ReplicatedPrefix, billing.LocalPrefix} {
		kvAttr, err := core.GetStoredAttributionData(context.Background(), localPathPrefix, endingMonth, billing.KvHWMCountsHWM)
		require.NoError(t, err)
		require.NotNil(t, kvAttr, "KV HWM attribution should be stored for the ending month under %s", localPathPrefix)
		require.Len(t, kvAttr.Mounts, 1, "KV HWM attribution should have one mount entry under %s", localPathPrefix)
		kvMount, ok := kvAttr.Mounts["kv_accessor_1"]
		require.True(t, ok, "KV HWM attribution should be keyed by kv_accessor_1 under %s", localPathPrefix)
		require.Equal(t, "kv_accessor_1", kvMount.MountAccessor)
		require.Equal(t, "secret/", kvMount.MountPath)
		require.Equal(t, "kv", kvMount.MountType)
		require.Equal(t, "root", kvMount.NamespaceID)
		require.Equal(t, "", kvMount.NamespacePath)
	}

	require.Equal(t, uint64(0), core.GetInMemoryTransitDataProtectionCallCounts())
	require.Equal(t, uint64(0), core.GetInMemoryTransformDataProtectionCallCounts())
	require.Equal(t, uint64(0), core.GetInMemoryGcpKmsDataProtectionCallCounts())
	require.False(t, core.consumptionBilling.KmipSeenEnabledThisMonth.Load())
}

// TestDeleteExpiredBillingMetrics_CustomRetention tests that DeleteExpiredBillingMetrics
// respects custom retention configuration. It verifies that when a custom retention period
// is set (e.g., 13 months), data is deleted according to that configuration rather than
// the default 37 months.
func TestDeleteExpiredBillingMetrics_CustomRetention(t *testing.T) {
	coreConfig := &CoreConfig{
		LogicalBackends: roleLogicalBackends,
	}
	core, _, _ := TestCoreUnsealedWithConfig(t, coreConfig)
	ctx := namespace.RootContext(context.Background())

	// Set custom retention to minimum (13 months)
	customRetention := billing.MinBillingRetentionMonths
	err := core.UpdateBillingRetentionMonths(ctx, customRetention)
	require.NoError(t, err)

	// Verify the custom retention was set
	retentionMonths, err := core.GetBillingRetentionMonths(ctx)
	require.NoError(t, err)
	require.Equal(t, customRetention, retentionMonths)

	now := time.Now().UTC()
	currentMonth := timeutil.StartOfMonth(now)

	// With 13 months retention:
	// - Month 12 months ago (index 12) should be kept (oldest retained)
	// - Month 13 months ago (index 13) should be deleted
	oldestRetainedMonth := currentMonth.AddDate(0, -(customRetention - 1), 0)
	monthToDelete := currentMonth.AddDate(0, -customRetention, 0)

	// Write billing data for multiple months
	for _, month := range []time.Time{monthToDelete, oldestRetainedMonth, currentMonth} {
		for _, pathPrefix := range []string{billing.ReplicatedPrefix, billing.LocalPrefix} {
			core.storeMaxRoleCountsLocked(context.Background(), &RoleCounts{
				AWSDynamicRoles:     5,
				AWSStaticRoles:      10,
				LDAPDynamicRoles:    3,
				OSLocalAccountRoles: 7,
			}, pathPrefix, month)
			core.storeMaxKvCountsLocked(context.Background(), 20, pathPrefix, month)
			core.storeTransitCallCountsLocked(context.Background(), 15, pathPrefix, month)
			core.storeSSHDurationAdjustedCertCountLocked(context.Background(), pathPrefix, month, 10.5)
			core.storeSSHOTPCountLocked(context.Background(), pathPrefix, month, 25.0)
		}
		// Store updatedAtTimestamp for each month
		testUpdateTime := time.Date(month.Year(), month.Month(), 15, 12, 0, 0, 0, time.UTC)
		err := core.UpdateMetricsLastUpdateTime(context.Background(), month, testUpdateTime)
		require.NoError(t, err)
	}

	// Verify data exists before deletion
	for _, pathPrefix := range []string{billing.ReplicatedPrefix, billing.LocalPrefix} {
		view, ok := core.GetBillingSubView()
		require.True(t, ok)

		// Check month to be deleted has data
		paths, err := view.List(context.Background(), billing.GetMonthlyBillingPath(pathPrefix, monthToDelete))
		require.NoError(t, err)
		require.Greater(t, len(paths), 0, "month to delete should have data before deletion")

		// Check oldest retained month has data
		paths, err = view.List(context.Background(), billing.GetMonthlyBillingPath(pathPrefix, oldestRetainedMonth))
		require.NoError(t, err)
		require.Greater(t, len(paths), 0, "oldest retained month should have data")
	}

	// Call DeleteExpiredBillingMetrics - it should use the custom retention
	err = core.DeleteExpiredBillingMetrics(context.Background(), currentMonth)
	require.NoError(t, err)

	// Verify deletion results with custom retention
	for _, pathPrefix := range []string{billing.ReplicatedPrefix, billing.LocalPrefix} {
		view, ok := core.GetBillingSubView()
		require.True(t, ok)

		// Month to delete (13 months ago with custom retention) should have no data
		paths, err := view.List(context.Background(), billing.GetMonthlyBillingPath(pathPrefix, monthToDelete))
		require.NoError(t, err)
		require.Equal(t, 0, len(paths), "data from %d months ago should be deleted with custom retention", customRetention)

		// Oldest retained month (12 months ago) should still have data
		paths, err = view.List(context.Background(), billing.GetMonthlyBillingPath(pathPrefix, oldestRetainedMonth))
		require.NoError(t, err)
		require.Greater(t, len(paths), 0, "data from %d months ago should be kept with custom retention", customRetention-1)

		// Current month should still have data
		paths, err = view.List(context.Background(), billing.GetMonthlyBillingPath(pathPrefix, currentMonth))
		require.NoError(t, err)
		require.Greater(t, len(paths), 0, "current month data should be kept")
	}

	// Verify updatedAtTimestamp deletion with custom retention
	deletedTimestamp, err := core.GetMetricsLastUpdateTime(context.Background(), monthToDelete)
	require.NoError(t, err)
	require.True(t, deletedTimestamp.IsZero(), "timestamp for deleted month should be zero")

	oldestTimestamp, err := core.GetMetricsLastUpdateTime(context.Background(), oldestRetainedMonth)
	require.NoError(t, err)
	require.False(t, oldestTimestamp.IsZero(), "timestamp for oldest retained month should exist")

	currentTimestamp, err := core.GetMetricsLastUpdateTime(context.Background(), currentMonth)
	require.NoError(t, err)
	require.False(t, currentTimestamp.IsZero(), "timestamp for current month should exist")
}

// TestHandleStartOfMonth_CustomRetention tests that HandleStartOfMonth respects
// custom retention configuration when deleting expired billing metrics.
func TestHandleStartOfMonth_CustomRetention(t *testing.T) {
	coreConfig := &CoreConfig{
		LogicalBackends: roleLogicalBackends,
		BillingConfig: billing.BillingConfig{
			MetricsUpdateCadence: 3 * time.Second,
		},
	}
	core, _, _ := TestCoreUnsealedWithConfig(t, coreConfig)
	ctx := namespace.RootContext(context.Background())

	// Set custom retention to 20 months
	customRetention := 20
	err := core.UpdateBillingRetentionMonths(ctx, customRetention)
	require.NoError(t, err)

	// Add billing metrics for months based on custom retention
	now := time.Now().UTC()
	oldestRetainedMonth := timeutil.StartOfMonth(now).AddDate(0, -(customRetention - 1), 0)
	monthToDelete := timeutil.StartOfMonth(now).AddDate(0, -customRetention, 0)

	for _, month := range []time.Time{monthToDelete, oldestRetainedMonth} {
		for _, pathPrefix := range []string{billing.ReplicatedPrefix, billing.LocalPrefix} {
			core.storeMaxRoleCountsLocked(context.Background(), &RoleCounts{
				AWSDynamicRoles:      10,
				AWSStaticRoles:       15,
				LDAPDynamicRoles:     8,
				GCPRolesets:          3,
				DatabaseDynamicRoles: 5,
				DatabaseStaticRoles:  7,
				OSLocalAccountRoles:  9,
			}, pathPrefix, month)
			core.storeMaxKvCountsLocked(context.Background(), 10, pathPrefix, month)

			if pathPrefix == billing.LocalPrefix {
				core.storeTransitCallCountsLocked(context.Background(), 10, pathPrefix, month)
				core.storeGcpKmsCallCountsLocked(context.Background(), 10, pathPrefix, month)
				core.storeThirdPartyPluginCountsLocked(context.Background(), pathPrefix, month, 10)
				core.storeOidcDurationAdjustedCountLocked(context.Background(), month, 10)
				core.storeSSHOTPCountLocked(context.Background(), pathPrefix, month, 10)
			}

			// Verify data was stored
			view, ok := core.GetBillingSubView()
			require.True(t, ok)
			paths, err := view.List(context.Background(), billing.GetMonthlyBillingPath(pathPrefix, month))
			require.NoError(t, err)
			expectedPaths := 2 // ReplicatedPrefix has roles and kv
			if pathPrefix == billing.LocalPrefix {
				expectedPaths = 7 // LocalPrefix has roles, kv, transit, gcp kms, third-party plugins, ssh and OIDC
			}
			require.Equal(t, expectedPaths, len(paths))
		}
	}

	// Handle the start of the month - should delete based on custom retention
	core.HandleStartOfMonth(context.Background(), now)

	for _, pathPrefix := range []string{billing.ReplicatedPrefix, billing.LocalPrefix} {
		// Month to delete (customRetention months ago) should have no billing metrics
		view, ok := core.GetBillingSubView()
		require.True(t, ok)
		paths, err := view.List(context.Background(), billing.GetMonthlyBillingPath(pathPrefix, monthToDelete))
		require.NoError(t, err)
		require.Equal(t, 0, len(paths), "data from %d months ago should be deleted with custom retention", customRetention)

		// Oldest retained month should still have the billing metrics
		paths, err = view.List(context.Background(), billing.GetMonthlyBillingPath(pathPrefix, oldestRetainedMonth))
		require.NoError(t, err)
		expectedPaths := 2 // ReplicatedPrefix has roles and kv
		if pathPrefix == billing.LocalPrefix {
			expectedPaths = 7 // LocalPrefix has roles, kv, transit, gcp kms, third-party plugins, ssh and OIDC
		}
		require.Equal(t, expectedPaths, len(paths), "data from %d months ago should be kept with custom retention", customRetention-1)
	}

	require.Equal(t, uint64(0), core.GetInMemoryTransitDataProtectionCallCounts())
	require.Equal(t, uint64(0), core.GetInMemoryTransformDataProtectionCallCounts())
	require.Equal(t, uint64(0), core.GetInMemoryGcpKmsDataProtectionCallCounts())
	require.Equal(t, float64(0), core.GetInMemoryOidcCounts())
	require.False(t, core.consumptionBilling.KmipSeenEnabledThisMonth.Load())
}
