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

// TestHandleEndOfMonthMetrics tests that that HandleEndOfMonth cleans up
// previous month billing metrics and resets the in memory billing metrics
func TestHandleEndOfMonthMetrics(t *testing.T) {
	coreConfig := &CoreConfig{
		LogicalBackends: roleLogicalBackends,
		BillingConfig: billing.BillingConfig{
			MetricsUpdateCadence: 3 * time.Second,
		},
	}
	core, _, _ := TestCoreUnsealedWithConfig(t, coreConfig)
	// Add some billing metrics to storage for the previous two months
	// Use the util functions directly to avoid the need to mount the logical backends
	previousMonth := timeutil.StartOfPreviousMonth(time.Now().UTC())
	twoMonthsAgo := previousMonth.AddDate(0, -1, 0)
	for _, month := range []time.Time{timeutil.StartOfPreviousMonth(time.Now().UTC()), timeutil.StartOfPreviousMonth(time.Now().UTC()).AddDate(0, -1, 0)} {
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
			core.storeTransitCallCountsLocked(context.Background(), 10, localPathPrefix, month)
			core.storeThirdPartyPluginCountsLocked(context.Background(), localPathPrefix, month, 10)

			// List the data paths to verify that the billing metrics have been stored
			view, ok := core.GetBillingSubView()
			require.True(t, ok)
			paths, err := view.List(context.Background(), billing.GetMonthlyBillingPath(localPathPrefix, month))
			require.NoError(t, err)
			require.Equal(t, 4, len(paths))
		}
	}

	// Handle the end of the month
	core.HandleEndOfMonth(context.Background(), time.Now().UTC())

	for _, localPathPrefix := range []string{billing.ReplicatedPrefix, billing.LocalPrefix} {
		// Two months ago should have no billing metrics
		view, ok := core.GetBillingSubView()
		require.True(t, ok)
		paths, err := view.List(context.Background(), billing.GetMonthlyBillingPath(localPathPrefix, twoMonthsAgo))
		require.NoError(t, err)
		require.Equal(t, 0, len(paths))

		// Previous month should have the billing metrics
		view, ok = core.GetBillingSubView()
		require.True(t, ok)
		paths, err = view.List(context.Background(), billing.GetMonthlyBillingPath(localPathPrefix, previousMonth))
		require.NoError(t, err)
		require.Equal(t, 4, len(paths))
	}

	require.Equal(t, uint64(0), core.GetInMemoryTransitDataProtectionCallCounts())
	require.Equal(t, uint64(0), core.GetInMemoryTransformDataProtectionCallCounts())
	require.False(t, core.consumptionBilling.KmipSeenEnabledThisMonth.Load())
}

// TestConsumptionBillingMetricsWorkerWithCustomClock tests that we correctly delete n-2 month billing metrics
// and reset the in memory billing metrics when the clock is overridden for testing purposes
func TestConsumptionBillingMetricsWorkerWithCustomClock(t *testing.T) {
	// 10 seconds until a new month (leave buffer for require.Eventually timeout)
	now := time.Date(2021, 1, 31, 23, 59, 50, 0, time.UTC)
	coreConfig := &CoreConfig{
		LogicalBackends: roleLogicalBackends,
		BillingConfig: billing.BillingConfig{
			Clock: newMockTimeNowClock(now),
		},
	}
	core, _, _ := TestCoreUnsealedWithConfig(t, coreConfig)

	// Add some billing metrics to storage for the previous two months
	// Use the util functions directly to avoid the need to mount the logical backends
	// The worker's "end of month" path calls HandleEndOfMonth with the *current* month,
	// which will be the next month once we cross the boundary. So "previousMonth" and
	// "twoMonthsAgo" should be calculated relative to that boundary.
	currentMonthAtBoundary := timeutil.StartOfNextMonth(now)
	previousMonth := timeutil.StartOfPreviousMonth(currentMonthAtBoundary)
	twoMonthsAgo := previousMonth.AddDate(0, -1, 0)
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
			thirdPartyPluginCounts, err := core.GetStoredThirdPartyPluginCounts(context.Background(), month)
			require.NoError(t, err)
			require.Equal(t, 10, thirdPartyPluginCounts)
		}
	}

	for _, month := range []time.Time{twoMonthsAgo, previousMonth} {

		core.storeTransitCallCountsLocked(context.Background(), uint64(10), billing.LocalPrefix, month)
		core.storeThirdPartyPluginCountsLocked(context.Background(), billing.LocalPrefix, month, 10)

		for _, localPathPrefix := range []string{billing.ReplicatedPrefix, billing.LocalPrefix} {
			core.storeMaxRoleCountsLocked(context.Background(), roleCounts, localPathPrefix, month)
			core.storeMaxKvCountsLocked(context.Background(), 10, localPathPrefix, month)
		}

		verifyMonthlyBillingMetrics(month, billing.LocalPrefix)
		verifyMonthlyBillingMetrics(month, billing.ReplicatedPrefix)
	}

	for _, localPathPrefix := range []string{billing.ReplicatedPrefix, billing.LocalPrefix} {
		// Two months ago should have no billing metricss should eventually should have no billing metrics
		require.Eventually(t, func() bool {
			paths, err := view.List(context.Background(), billing.GetMonthlyBillingPath(localPathPrefix, twoMonthsAgo))
			return err == nil && len(paths) == 0
		}, 20*time.Second, 100*time.Millisecond)

		// All values n-2 months ago should be 0
		maxRoleCounts, _ := core.GetStoredHWMRoleCounts(context.Background(), localPathPrefix, twoMonthsAgo)
		require.Equal(t, &RoleCounts{}, maxRoleCounts)
		kvCounts, _ := core.GetStoredHWMKvCounts(context.Background(), localPathPrefix, twoMonthsAgo)
		require.Equal(t, 0, kvCounts)
		if localPathPrefix == billing.LocalPrefix {
			transitCounts, _ := core.GetStoredTransitCallCounts(context.Background(), twoMonthsAgo)
			require.Equal(t, uint64(0), transitCounts)
			thirdPartyPluginCounts, _ := core.GetStoredThirdPartyPluginCounts(context.Background(), twoMonthsAgo)
			require.Equal(t, 0, thirdPartyPluginCounts)
		}

		// Previous month should have the billing metrics
		verifyMonthlyBillingMetrics(previousMonth, localPathPrefix)
	}

	require.Equal(t, uint64(0), core.GetInMemoryTransitDataProtectionCallCounts())
	require.Equal(t, uint64(0), core.GetInMemoryTransformDataProtectionCallCounts())
	require.False(t, core.consumptionBilling.KmipSeenEnabledThisMonth.Load())
}
