// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/pluginconsts"
	"github.com/hashicorp/vault/helper/timeutil"
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
	"OpenLDAP Dynamic Roles": {
		mount: pluginconsts.SecretEngineOpenLDAP,
		key:   "role/",
	},
	"OpenLDAP Static Roles": {
		mount: pluginconsts.SecretEngineOpenLDAP,
		key:   "static-role/",
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
// billing metrics from billing.BillingRetentionMonths ago (keeping billing.BillingRetentionMonths of data) and resets the in memory billing metrics
func TestHandleEndOfMonthMetrics(t *testing.T) {
	coreConfig := &CoreConfig{
		LogicalBackends: roleLogicalBackends,
		BillingConfig: billing.BillingConfig{
			MetricsUpdateCadence: 3 * time.Second,
		},
	}
	core, _, _ := TestCoreUnsealedWithConfig(t, coreConfig)
	// Add some billing metrics to storage for (billing.BillingRetentionMonths - 1) and billing.BillingRetentionMonths months ago
	// Use the util functions directly to avoid the need to mount the logical backends
	now := time.Now().UTC()
	oldestRetainedMonth := timeutil.StartOfMonth(now).AddDate(0, -(billing.BillingRetentionMonths - 1), 0)
	monthToDelete := timeutil.StartOfMonth(now).AddDate(0, -billing.BillingRetentionMonths, 0)

	for _, month := range []time.Time{monthToDelete, oldestRetainedMonth} {
		for _, localPathPrefix := range []string{billing.ReplicatedPrefix, billing.LocalPrefix} {
			core.storeMaxRoleCountsLocked(context.Background(), &RoleCounts{
				AWSDynamicRoles:      10,
				AWSStaticRoles:       15,
				LDAPDynamicRoles:     8,
				GCPRolesets:          3,
				DatabaseDynamicRoles: 5,
				DatabaseStaticRoles:  7,
			}, localPathPrefix, month)
			core.storeMaxKvCountsLocked(context.Background(), 10, localPathPrefix, month)

			// Transit, third-party plugins, ssh credential count and OIDC are local aggregated metrics
			// and should only be stored under LocalPrefix
			if localPathPrefix == billing.LocalPrefix {
				core.storeTransitCallCountsLocked(context.Background(), 10, localPathPrefix, month)
				core.storeGcpKmsCallCountsLocked(context.Background(), 10, localPathPrefix, month)
				core.storeThirdPartyPluginCountsLocked(context.Background(), localPathPrefix, month, 10)
				core.storeOidcDurationAdjustedCountLocked(context.Background(), month, 10)
				core.storeSSHOTPCountLocked(context.Background(), localPathPrefix, month, 10)
			}

			// List the data paths to verify that the billing metrics have been stored
			view, ok := core.GetBillingSubView()
			require.True(t, ok)
			paths, err := view.List(context.Background(), billing.GetMonthlyBillingPath(localPathPrefix, month))
			require.NoError(t, err)
			expectedPaths := 2 // ReplicatedPrefix has roles and kv
			if localPathPrefix == billing.LocalPrefix {
				expectedPaths = 7 // LocalPrefix has roles, kv, transit, gcp kms, third-party plugins, ssh and OIDC
			}
			require.Equal(t, expectedPaths, len(paths))
		}
	}

	// Handle the end of the month
	core.HandleStartOfMonth(context.Background(), now)

	for _, localPathPrefix := range []string{billing.ReplicatedPrefix, billing.LocalPrefix} {
		// billing.BillingRetentionMonths ago should have no billing metrics (deleted)
		view, ok := core.GetBillingSubView()
		require.True(t, ok)
		paths, err := view.List(context.Background(), billing.GetMonthlyBillingPath(localPathPrefix, monthToDelete))
		require.NoError(t, err)
		require.Equal(t, 0, len(paths), "data from billing.BillingRetentionMonths ago should be deleted")

		// (billing.BillingRetentionMonths - 1) months ago should still have the billing metrics (kept)
		view, ok = core.GetBillingSubView()
		require.True(t, ok)
		paths, err = view.List(context.Background(), billing.GetMonthlyBillingPath(localPathPrefix, oldestRetainedMonth))
		require.NoError(t, err)
		expectedPaths := 2 // ReplicatedPrefix has roles and kv
		if localPathPrefix == billing.LocalPrefix {
			expectedPaths = 7 // LocalPrefix has roles, kv, transit, gcp kms, third-party plugins, ssh and OIDC
		}
		require.Equal(t, expectedPaths, len(paths))
	}

	require.Equal(t, uint64(0), core.GetInMemoryTransitDataProtectionCallCounts())
	require.Equal(t, uint64(0), core.GetInMemoryTransformDataProtectionCallCounts())
	require.Equal(t, uint64(0), core.GetInMemoryGcpKmsDataProtectionCallCounts())
	require.Equal(t, float64(0), core.GetInMemoryOidcCounts())
	require.False(t, core.consumptionBilling.KmipSeenEnabledThisMonth.Load())
}

// TestDeleteExpiredBillingMetrics specifically tests the deleteExpiredBillingMetrics method
// to ensure it correctly deletes data from billing.BillingRetentionMonths ago while keeping
// data from (billing.BillingRetentionMonths - 1) months ago.
func TestDeleteExpiredBillingMetrics(t *testing.T) {
	coreConfig := &CoreConfig{
		LogicalBackends: roleLogicalBackends,
	}
	core, _, _ := TestCoreUnsealedWithConfig(t, coreConfig)

	now := time.Now().UTC()
	currentMonth := timeutil.StartOfMonth(now)
	oldestRetainedMonth := currentMonth.AddDate(0, -(billing.BillingRetentionMonths - 1), 0)
	monthToDelete := currentMonth.AddDate(0, -billing.BillingRetentionMonths, 0)

	// Write billing data for multiple months including the month to be deleted and the oldest retained month
	for _, month := range []time.Time{monthToDelete, oldestRetainedMonth, currentMonth} {
		for _, pathPrefix := range []string{billing.ReplicatedPrefix, billing.LocalPrefix} {
			core.storeMaxRoleCountsLocked(context.Background(), &RoleCounts{
				AWSDynamicRoles:  5,
				AWSStaticRoles:   10,
				LDAPDynamicRoles: 3,
			}, pathPrefix, month)
			core.storeMaxKvCountsLocked(context.Background(), 20, pathPrefix, month)
			core.storeTransitCallCountsLocked(context.Background(), 15, pathPrefix, month)
			// Add SSH metrics which use subdirectory paths (ssh/normalized-certs-issued, ssh/credential-count)
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

		// Verify SSH metrics exist (they use subdirectory paths)
		sshCertPath := billing.GetMonthlyBillingMetricPath(pathPrefix, monthToDelete, billing.SSHCertificateMetric)
		entry, err := view.Get(context.Background(), sshCertPath)
		require.NoError(t, err)
		require.NotNil(t, entry, "SSH cert metric should exist before deletion")

		sshOTPPath := billing.GetMonthlyBillingMetricPath(pathPrefix, monthToDelete, billing.SSHOTPMetric)
		entry, err = view.Get(context.Background(), sshOTPPath)
		require.NoError(t, err)
		require.NotNil(t, entry, "SSH OTP metric should exist before deletion")
	}

	// Verify updatedAtTimestamp exists for all months before deletion
	for _, month := range []time.Time{monthToDelete, oldestRetainedMonth, currentMonth} {
		timestamp, err := core.GetMetricsLastUpdateTime(context.Background(), month)
		require.NoError(t, err)
		require.False(t, timestamp.IsZero(), "timestamp for month %s should exist before deletion", month.Format("2006-01"))
	}

	// Call deleteExpiredBillingMetrics directly
	err := core.deleteExpiredBillingMetrics(context.Background(), currentMonth)
	require.NoError(t, err)

	// Verify deletion results
	for _, pathPrefix := range []string{billing.ReplicatedPrefix, billing.LocalPrefix} {
		view, ok := core.GetBillingSubView()
		require.True(t, ok)

		// Month to delete should have no data
		paths, err := view.List(context.Background(), billing.GetMonthlyBillingPath(pathPrefix, monthToDelete))
		require.NoError(t, err)
		require.Equal(t, 0, len(paths), "data from billing.BillingRetentionMonths ago should be deleted")

		// Verify SSH metrics are deleted (they use subdirectory paths)
		sshCertPath := billing.GetMonthlyBillingMetricPath(pathPrefix, monthToDelete, billing.SSHCertificateMetric)
		entry, err := view.Get(context.Background(), sshCertPath)
		require.NoError(t, err)
		require.Nil(t, entry, "SSH cert metric should be deleted")

		sshOTPPath := billing.GetMonthlyBillingMetricPath(pathPrefix, monthToDelete, billing.SSHOTPMetric)
		entry, err = view.Get(context.Background(), sshOTPPath)
		require.NoError(t, err)
		require.Nil(t, entry, "SSH OTP metric should be deleted")

		// Oldest retained month should still have data
		paths, err = view.List(context.Background(), billing.GetMonthlyBillingPath(pathPrefix, oldestRetainedMonth))
		require.NoError(t, err)
		require.Greater(t, len(paths), 0, "data from (billing.BillingRetentionMonths - 1) months ago should be kept")

		// Verify SSH metrics are kept for oldest retained month
		sshCertPath = billing.GetMonthlyBillingMetricPath(pathPrefix, oldestRetainedMonth, billing.SSHCertificateMetric)
		entry, err = view.Get(context.Background(), sshCertPath)
		require.NoError(t, err)
		require.NotNil(t, entry, "SSH cert metric should be kept for oldest retained month")

		sshOTPPath = billing.GetMonthlyBillingMetricPath(pathPrefix, oldestRetainedMonth, billing.SSHOTPMetric)
		entry, err = view.Get(context.Background(), sshOTPPath)
		require.NoError(t, err)
		require.NotNil(t, entry, "SSH OTP metric should be kept for oldest retained month")

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

// TestConsumptionBillingMetricsWorkerWithCustomClock tests that we correctly delete data older than billing.BillingRetentionMonths
// and reset the in memory billing metrics when the clock is overridden for testing purposes
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

	// Add some billing metrics to storage for (billing.BillingRetentionMonths - 1) and billing.BillingRetentionMonths months ago
	// Use the util functions directly to avoid the need to mount the logical backends
	// The worker's "end of month" path calls HandleEndOfMonth with the *current* month,
	// which will be the next month once we cross the boundary. So the months should be
	// calculated relative to that boundary.
	currentMonthAtBoundary := timeutil.StartOfNextMonth(now)
	oldestRetainedMonth := timeutil.StartOfMonth(currentMonthAtBoundary).AddDate(0, -(billing.BillingRetentionMonths - 1), 0)
	monthToDelete := timeutil.StartOfMonth(currentMonthAtBoundary).AddDate(0, -billing.BillingRetentionMonths, 0)
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

	for _, localPathPrefix := range []string{billing.ReplicatedPrefix, billing.LocalPrefix} {
		// billing.BillingRetentionMonths ago should eventually have no billing metrics (deleted)
		require.Eventually(t, func() bool {
			paths, err := view.List(context.Background(), billing.GetMonthlyBillingPath(localPathPrefix, monthToDelete))
			return err == nil && len(paths) == 0
		}, 20*time.Second, 100*time.Millisecond)

		// All values from billing.BillingRetentionMonths ago should be 0
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

		// (billing.BillingRetentionMonths - 1) months ago should still have the billing metrics (kept)
		verifyMonthlyBillingMetrics(oldestRetainedMonth, localPathPrefix)
	}

	require.Equal(t, uint64(0), core.GetInMemoryTransitDataProtectionCallCounts())
	require.Equal(t, uint64(0), core.GetInMemoryTransformDataProtectionCallCounts())
	require.Equal(t, uint64(0), core.GetInMemoryGcpKmsDataProtectionCallCounts())
	require.False(t, core.consumptionBilling.KmipSeenEnabledThisMonth.Load())
}
