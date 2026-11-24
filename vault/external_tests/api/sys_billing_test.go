// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package api

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
	logicalAws "github.com/hashicorp/vault/builtin/logical/aws"
	logicalDatabase "github.com/hashicorp/vault/builtin/logical/database"
	logicalTransit "github.com/hashicorp/vault/builtin/logical/transit"
	"github.com/hashicorp/vault/helper/pluginconsts"
	"github.com/hashicorp/vault/helper/testhelpers/minimal"
	"github.com/hashicorp/vault/internalshared/timeutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/billing"
	"github.com/stretchr/testify/require"
)

// Test_BillingOverview tests that the BillingOverview API method works correctly
func Test_BillingOverview(t *testing.T) {
	t.Parallel()

	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			pluginconsts.SecretEngineAWS:      logicalAws.Factory,
			pluginconsts.SecretEngineDatabase: logicalDatabase.Factory,
			pluginconsts.SecretEngineTransit:  logicalTransit.Factory,
		},
	}

	cluster := minimal.NewTestSoloCluster(t, coreConfig)
	client := cluster.Cores[0].Client

	// Mount AWS for dynamic roles
	err := client.Sys().Mount("aws", &api.MountInput{
		Type: "aws",
	})
	require.NoError(t, err)

	// Create an AWS role
	_, err = client.Logical().Write("aws/roles/test-role", map[string]interface{}{
		"credential_type": "iam_user",
		"policy_document": `{"Version": "2012-10-17","Statement": [{"Effect": "Allow","Action": "ec2:*","Resource": "*"}]}`,
	})
	require.NoError(t, err)

	// Mount Database for dynamic roles
	err = client.Sys().Mount("database", &api.MountInput{
		Type: "database",
	})
	require.NoError(t, err)

	// Create a database role
	_, err = client.Logical().Write("database/roles/test-db-role", map[string]interface{}{
		"db_name":             "test-db",
		"creation_statements": []string{"CREATE USER '{{name}}'@'%' IDENTIFIED BY '{{password}}';"},
		"default_ttl":         "1h",
		"max_ttl":             "24h",
	})
	require.NoError(t, err)

	// Mount KV for static secrets
	err = client.Sys().Mount("kv-v2", &api.MountInput{
		Type: "kv-v2",
	})
	require.NoError(t, err)

	// Create KV secrets
	secretData := map[string]interface{}{
		"foo": "bar",
	}
	_, err = client.KVv2("kv-v2").Put(context.Background(), "secret1", secretData)
	require.NoError(t, err)

	resp, err := client.Sys().BillingOverview(true)
	require.NoError(t, err)
	require.NotNil(t, resp)

	// Validate response structure
	require.NotNil(t, resp.Months)
	require.Len(t, resp.Months, billing.DefaultBillingRetentionMonths, "should have billing.DefaultBillingRetentionMonths months")

	// Check current month data
	currentMonth := resp.Months[0]
	require.NotEmpty(t, currentMonth.Month)
	require.NotEmpty(t, currentMonth.UpdatedAt)
	require.NotNil(t, currentMonth.UsageMetrics)

	// Verify we have some metrics
	require.NotEmpty(t, currentMonth.UsageMetrics, "should have usage metrics after creating test data")

	// Validate that metrics have the expected structure
	for _, metric := range currentMonth.UsageMetrics {
		require.NotEmpty(t, metric.MetricName)
		require.NotNil(t, metric.MetricData)
		require.NotEmpty(t, metric.MetricData)

		if total, ok := metric.MetricData["total"]; ok {
			_, ok := total.(json.Number)
			require.True(t, ok, "total should be json.Number for metric %s", metric.MetricName)
		}

		if details, ok := metric.MetricData["metric_details"]; ok {
			_, ok := details.([]interface{})
			require.True(t, ok, "metric_details should be []interface{} for metric %s", metric.MetricName)
		}
	}
}

// Test_BillingOverview_WithoutUpdateCounts tests that BillingOverview works with updateCounts=false
func Test_BillingOverview_WithoutUpdateCounts(t *testing.T) {
	t.Parallel()

	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			pluginconsts.SecretEngineAWS: logicalAws.Factory,
		},
	}

	cluster := minimal.NewTestSoloCluster(t, coreConfig)
	client := cluster.Cores[0].Client

	// Call BillingOverview without updating counts
	resp, err := client.Sys().BillingOverview(false)
	require.NoError(t, err)
	require.NotNil(t, resp)

	// Validate basic response structure
	require.NotNil(t, resp.Months)
	require.Len(t, resp.Months, billing.DefaultBillingRetentionMonths, "should have billing.DefaultBillingRetentionMonths months")

	// Check that months are properly formatted
	for _, month := range resp.Months {
		require.NotEmpty(t, month.Month)
		require.NotEmpty(t, month.UpdatedAt)
		require.NotNil(t, month.UsageMetrics)
	}
}

// Test_BillingOverview_EmptyCluster tests BillingOverview on a cluster with no mounts.
// Verifies that all metrics are present with zero values
func Test_BillingOverview_EmptyCluster(t *testing.T) {
	t.Parallel()

	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	resp, err := client.Sys().BillingOverview(true)
	require.NoError(t, err)
	require.NotNil(t, resp)

	require.NotNil(t, resp.Months)
	require.Len(t, resp.Months, billing.DefaultBillingRetentionMonths)

	currentMonth := resp.Months[0]
	require.NotEmpty(t, currentMonth.Month)
	require.NotEmpty(t, currentMonth.UpdatedAt)
	require.NotNil(t, currentMonth.UsageMetrics)

	// Verify all expected metrics are present even with no usage
	require.NotEmpty(t, currentMonth.UsageMetrics, "should have all metrics even with zero values")

	expectedMetrics := map[string]bool{
		"static_secrets":        false,
		"dynamic_roles":         false,
		"auto_rotated_roles":    false,
		"kmip":                  false,
		"external_plugins":      false,
		"data_protection_calls": false,
		"pki_units":             false,
		"managed_keys":          false,
		"ssh_units":             false,
		"id_token_units":        false,
		"external_ca_pki_units": false,
	}

	for _, metric := range currentMonth.UsageMetrics {
		require.NotEmpty(t, metric.MetricName)
		require.Contains(t, expectedMetrics, metric.MetricName, "unexpected metric: %s", metric.MetricName)
		expectedMetrics[metric.MetricName] = true
		require.NotNil(t, metric.MetricData)
	}

	// Verify all expected metrics were found
	for metricName, found := range expectedMetrics {
		require.True(t, found, "metric %s should be present", metricName)
	}
}

// Test_BillingOverview_MonthFormat tests that month strings are in correct format
func Test_BillingOverview_MonthFormat(t *testing.T) {
	t.Parallel()

	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	resp, err := client.Sys().BillingOverview(false)
	require.NoError(t, err)
	require.NotNil(t, resp)

	// Verify month format (YYYY-MM)
	for _, month := range resp.Months {
		require.Regexp(t, `^\d{4}-\d{2}$`, month.Month, "month should be in YYYY-MM format")
		require.Regexp(t, `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}`, month.UpdatedAt, "updated_at should be in ISO 8601 format")
	}

	// Verify months are in descending order (current, then previous months)
	require.Greater(t, resp.Months[0].Month, resp.Months[1].Month, "first month should be more recent than second")
	// Verify we have billing.DefaultBillingRetentionMonths months
	require.Len(t, resp.Months, billing.DefaultBillingRetentionMonths)

	// Verify the oldest month is exactly (billing.DefaultBillingRetentionMonths - 1) months before the current month
	currentMonthTime, err := time.Parse("2006-01", resp.Months[0].Month)
	require.NoError(t, err, "should parse current month")
	oldestMonthTime, err := time.Parse("2006-01", resp.Months[billing.DefaultBillingRetentionMonths-1].Month)
	require.NoError(t, err, "should parse oldest month")

	expectedOldestMonth := currentMonthTime.AddDate(0, -(billing.DefaultBillingRetentionMonths - 1), 0)
	require.Equal(t, expectedOldestMonth.Format("2006-01"), oldestMonthTime.Format("2006-01"),
		"oldest month should be exactly %d months before current month", billing.DefaultBillingRetentionMonths-1)
}

// Test_BillingConfig tests the GetBillingConfig and SetBillingConfig API methods
func Test_BillingConfig(t *testing.T) {
	t.Parallel()

	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	// Test GetBillingConfig - should return default value
	config, err := client.Sys().GetBillingConfig()
	require.NoError(t, err)
	require.NotNil(t, config)
	require.Equal(t, billing.DefaultBillingRetentionMonths, config.RetentionMonths)

	// Test SetBillingConfig with valid value
	err = client.Sys().SetBillingConfig(24)
	require.NoError(t, err)

	// Verify the value was updated
	config, err = client.Sys().GetBillingConfig()
	require.NoError(t, err)
	require.Equal(t, 24, config.RetentionMonths)

	// Test SetBillingConfig with minimum value
	err = client.Sys().SetBillingConfig(billing.MinBillingRetentionMonths)
	require.NoError(t, err)

	config, err = client.Sys().GetBillingConfig()
	require.NoError(t, err)
	require.Equal(t, billing.MinBillingRetentionMonths, config.RetentionMonths)

	// Test SetBillingConfig with maximum value
	err = client.Sys().SetBillingConfig(billing.MaxBillingRetentionMonths)
	require.NoError(t, err)

	config, err = client.Sys().GetBillingConfig()
	require.NoError(t, err)
	require.Equal(t, billing.MaxBillingRetentionMonths, config.RetentionMonths)
}

// Test_BillingConfig_InvalidValues tests that invalid retention values are rejected
func Test_BillingConfig_InvalidValues(t *testing.T) {
	t.Parallel()

	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	// Test below minimum
	err := client.Sys().SetBillingConfig(billing.MinBillingRetentionMonths - 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "must be between")

	// Test above maximum
	err = client.Sys().SetBillingConfig(billing.MaxBillingRetentionMonths + 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "must be between")

	// Test zero
	err = client.Sys().SetBillingConfig(0)
	require.Error(t, err)
	require.Contains(t, err.Error(), "must be between")

	// Test negative
	err = client.Sys().SetBillingConfig(-1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "must be between")
}

// Test_BillingConfig_AffectsOverview tests that config changes affect billing overview
func Test_BillingConfig_AffectsOverview(t *testing.T) {
	t.Parallel()

	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	// Set retention to minimum
	err := client.Sys().SetBillingConfig(billing.MinBillingRetentionMonths)
	require.NoError(t, err)

	// Get billing overview
	resp, err := client.Sys().BillingOverview(false)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.Months, billing.MinBillingRetentionMonths)

	// Set retention to maximum
	err = client.Sys().SetBillingConfig(billing.MaxBillingRetentionMonths)
	require.NoError(t, err)

	// Get billing overview again
	resp, err = client.Sys().BillingOverview(false)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.Months, billing.MaxBillingRetentionMonths)

	// Set back to default
	err = client.Sys().SetBillingConfig(billing.DefaultBillingRetentionMonths)
	require.NoError(t, err)

	// Verify default is restored
	resp, err = client.Sys().BillingOverview(false)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.Months, billing.DefaultBillingRetentionMonths)
}

// Test_AttributionRetentionConfig tests GetBillingConfig and SetAttributionRetentionConfig
// for the attribution_retention_months field via the real Vault API.
func Test_AttributionRetentionConfig(t *testing.T) {
	t.Parallel()

	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	// Default: attribution_retention_months should be DefaultAttributionRetentionMonths
	config, err := client.Sys().GetBillingConfig()
	require.NoError(t, err)
	require.NotNil(t, config)
	require.Equal(t, billing.DefaultAttributionRetentionMonths, config.AttributionRetentionMonths)

	// Set to a valid mid-range value
	err = client.Sys().SetAttributionRetentionConfig(24)
	require.NoError(t, err)

	config, err = client.Sys().GetBillingConfig()
	require.NoError(t, err)
	require.Equal(t, 24, config.AttributionRetentionMonths)
	// billing retention_months must be unaffected
	require.Equal(t, billing.DefaultBillingRetentionMonths, config.RetentionMonths)

	// Set to minimum (0) — disables attribution
	err = client.Sys().SetAttributionRetentionConfig(billing.MinAttributionRetentionMonths)
	require.NoError(t, err)

	config, err = client.Sys().GetBillingConfig()
	require.NoError(t, err)
	require.Equal(t, billing.MinAttributionRetentionMonths, config.AttributionRetentionMonths)

	// Set to maximum (72)
	err = client.Sys().SetAttributionRetentionConfig(billing.MaxAttributionRetentionMonths)
	require.NoError(t, err)

	config, err = client.Sys().GetBillingConfig()
	require.NoError(t, err)
	require.Equal(t, billing.MaxAttributionRetentionMonths, config.AttributionRetentionMonths)
}

// Test_AttributionRetentionConfig_InvalidValues tests that invalid
// attribution_retention_months values are rejected.
func Test_AttributionRetentionConfig_InvalidValues(t *testing.T) {
	t.Parallel()

	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	// Below minimum: -1 is rejected
	err := client.Sys().SetAttributionRetentionConfig(-1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "must be between")

	// Above maximum: 73 is rejected
	err = client.Sys().SetAttributionRetentionConfig(billing.MaxAttributionRetentionMonths + 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "must be between")
}

// Test_AttributionRetentionConfig_DisablesCollection verifies the end-to-end behaviour
// when attribution is disabled (retention = 0):
//   - attribution data written before disabling is wiped immediately
//   - after disabling, new UpdateMaxKvCounts calls do not store attribution
//   - after re-enabling, UpdateMaxKvCounts stores attribution again
func Test_AttributionRetentionConfig_DisablesCollection(t *testing.T) {
	t.Parallel()

	cluster := minimal.NewTestSoloCluster(t, nil)
	core := cluster.Cores[0].Core
	vault.TestWaitActive(t, core)
	client := cluster.Cores[0].Client

	ctx := context.Background()
	month := timeutil.StartOfMonth(time.Now().UTC())

	attribution := vault.MountAttributionMap{
		"kv_test": logical.MountAttribution{
			Count:         5,
			MountAccessor: "kv_test",
			MountPath:     "secret/",
			MountType:     "kv",
			NamespaceID:   "root",
		},
	}

	// Seed attribution data for the current month.
	_, err := core.UpdateMaxKvCounts(ctx, billing.LocalPrefix, month, 5, attribution)
	require.NoError(t, err)

	// Verify data was stored.
	stored, err := core.GetStoredAttributionData(ctx, billing.LocalPrefix, month, billing.KvHWMCountsHWM)
	require.NoError(t, err)
	require.Len(t, stored.Mounts, 1, "attribution should be present before disabling")

	// Disable attribution (sets retention to 0, wipes existing data).
	err = client.Sys().SetAttributionRetentionConfig(billing.MinAttributionRetentionMonths)
	require.NoError(t, err)

	// Existing attribution data must have been wiped.
	wipedStored, err := core.GetStoredAttributionData(ctx, billing.LocalPrefix, month, billing.KvHWMCountsHWM)
	require.NoError(t, err)
	require.Empty(t, wipedStored.Mounts, "attribution should be wiped after disabling")

	// New UpdateMaxKvCounts with a higher HWM must NOT store attribution.
	higherAttribution := vault.MountAttributionMap{
		"kv_higher": logical.MountAttribution{
			Count:         10,
			MountAccessor: "kv_higher",
			MountPath:     "higher/",
			MountType:     "kv",
			NamespaceID:   "root",
		},
	}
	_, err = core.UpdateMaxKvCounts(ctx, billing.LocalPrefix, month, 10, higherAttribution)
	require.NoError(t, err)

	afterDisable, err := core.GetStoredAttributionData(ctx, billing.LocalPrefix, month, billing.KvHWMCountsHWM)
	require.NoError(t, err)
	require.Empty(t, afterDisable.Mounts, "attribution must not be stored while disabled")

	// Re-enable attribution and verify the next update stores attribution again.
	err = client.Sys().SetAttributionRetentionConfig(billing.DefaultAttributionRetentionMonths)
	require.NoError(t, err)

	reEnabledAttribution := vault.MountAttributionMap{
		"kv_reenabled": logical.MountAttribution{
			Count:         20,
			MountAccessor: "kv_reenabled",
			MountPath:     "reenabled/",
			MountType:     "kv",
			NamespaceID:   "root",
		},
	}
	_, err = core.UpdateMaxKvCounts(ctx, billing.LocalPrefix, month, 20, reEnabledAttribution)
	require.NoError(t, err)

	afterReEnable, err := core.GetStoredAttributionData(ctx, billing.LocalPrefix, month, billing.KvHWMCountsHWM)
	require.NoError(t, err)
	require.Len(t, afterReEnable.Mounts, 1, "attribution should be stored after re-enabling")
	require.Contains(t, afterReEnable.Mounts, "kv_reenabled")
}

// Test_BillingOverview_RefreshInProgressFlag verifies that a refresh_data=true request
// covering the current month executes the scan and releases the in-progress lock upon completion,
// and that historical-only date ranges do not acquire or contend for the lock.
func Test_BillingOverview_RefreshInProgressFlag(t *testing.T) {
	t.Parallel()

	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client
	core := cluster.Cores[0].Core
	cb := core.GetCoreConsumptionBillingManager()
	require.NotNil(t, cb)

	// Happy path: the scan runs, records a non-zero timestamp, and releases the flag.
	resp, err := client.Logical().ReadWithData("sys/billing/overview", map[string][]string{
		"refresh_data": {"true"},
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotContains(t, resp.Warnings, vault.WarningRefreshAlreadyInProgress)

	months, ok := resp.Data["months"].([]interface{})
	require.True(t, ok)
	currentMonth, ok := months[0].(map[string]interface{})
	require.True(t, ok)
	updatedAt, err := time.Parse(time.RFC3339, currentMonth["updated_at"].(string))
	require.NoError(t, err)
	require.False(t, updatedAt.IsZero(), "the scan should have recorded a last-update timestamp")
	require.True(t, cb.RefreshInProgress.TryLock(), "lock should be released after the scan")
	cb.RefreshInProgress.Unlock()

	// Historical-only range: requests that do not include the current month should not acquire
	// or contend for the lock, even if refresh_data=true is passed while another scan is running.
	cb.RefreshInProgress.Lock()
	// Release before cluster cleanup so a blocked background worker can't stall shutdown.
	t.Cleanup(cb.RefreshInProgress.Unlock)
	currentMonthStart := timeutil.StartOfMonth(time.Now().UTC())
	lastMonth := currentMonthStart.AddDate(0, -1, 0).Format("2006-01")
	twoMonthsAgo := currentMonthStart.AddDate(0, -2, 0).Format("2006-01")

	resp, err = client.Logical().ReadWithData("sys/billing/overview", map[string][]string{
		"refresh_data": {"true"},
		"start_month":  {twoMonthsAgo},
		"end_month":    {lastMonth},
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotContains(t, resp.Warnings, vault.WarningRefreshAlreadyInProgress,
		"historical-only query should not receive in-progress warning")
}

// Test_BillingOverview_RefreshInProgressReturnsWarningAndData verifies that when a refresh
// is already in progress, a refresh_data=true request to sys/billing/overview does not fail
// or return an empty response: it returns the in-progress warning and still returns the
// stored billing data for the requested month, without scanning, and leaves the lock held
// by the refresh that is running.
func Test_BillingOverview_RefreshInProgressReturnsWarningAndData(t *testing.T) {
	t.Parallel()

	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client
	core := cluster.Cores[0].Core
	cb := core.GetCoreConsumptionBillingManager()
	require.NotNil(t, cb)

	// Seed a distinctive stored timestamp: a scan would overwrite it with the current time,
	// so getting it back unchanged shows the request returned stored data without scanning.
	seededAt := time.Date(2020, time.January, 2, 3, 4, 5, 0, time.UTC)
	require.NoError(t, core.UpdateMetricsLastUpdateTime(t.Context(), time.Now().UTC(), seededAt))

	// Simulate a refresh that is already running.
	cb.RefreshInProgress.Lock()
	// Release before cluster cleanup so a blocked background worker can't stall shutdown.
	t.Cleanup(cb.RefreshInProgress.Unlock)

	resp, err := client.Logical().ReadWithData("sys/billing/overview", map[string][]string{
		"refresh_data": {"true"},
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Contains(t, resp.Warnings, vault.WarningRefreshAlreadyInProgress)

	// The request must still return data for the current month despite the warning.
	require.NotNil(t, resp.Data)
	months, ok := resp.Data["months"].([]interface{})
	require.True(t, ok)
	require.NotEmpty(t, months)

	currentMonth, ok := months[0].(map[string]interface{})
	require.True(t, ok)
	require.Equal(t, timeutil.StartOfMonth(time.Now().UTC()).Format("2006-01"), currentMonth["month"])
	updatedAt, err := time.Parse(time.RFC3339, currentMonth["updated_at"].(string))
	require.NoError(t, err)
	require.True(t, seededAt.Equal(updatedAt), "skipped request must return stored data, not a refreshed timestamp")
	usageMetrics, ok := currentMonth["usage_metrics"].([]interface{})
	require.True(t, ok)
	require.NotEmpty(t, usageMetrics)

	require.False(t, cb.RefreshInProgress.TryLock(), "skipped request must leave the lock held for the running refresh")
}

// Test_BillingMetricsWorker_WaitsForAPIScanThenRuns verifies that the background billing
// worker never skips its scan when the refresh-in-progress lock is held (as by an
// API-triggered scan): it waits for the lock and then completes the contended flush
// without waiting for another tick.
func Test_BillingMetricsWorker_WaitsForAPIScanThenRuns(t *testing.T) {
	t.Parallel()

	const cadence = 2 * time.Second

	cluster := minimal.NewTestSoloCluster(t, &vault.CoreConfig{
		BillingConfig: billing.BillingConfig{MetricsUpdateCadence: cadence},
	})
	core := cluster.Cores[0].Core
	cb := core.GetCoreConsumptionBillingManager()
	require.NotNil(t, cb)

	// Simulate an API scan holding the lock before the first worker tick.
	require.Eventually(t, cb.RefreshInProgress.TryLock,
		15*time.Second, 10*time.Millisecond, "simulated API scan should acquire the refresh lock")
	locked := true
	t.Cleanup(func() {
		if locked {
			cb.RefreshInProgress.Unlock()
		}
	})

	// Seed a distinctive timestamp that only a worker scan would overwrite.
	seededAt := time.Date(2020, time.January, 2, 3, 4, 5, 0, time.UTC)
	require.NoError(t, core.UpdateMetricsLastUpdateTime(t.Context(), time.Now().UTC(), seededAt))

	// Hold the lock for longer than one tick cadence so a tick fires and blocks on the lock.
	require.Never(t, func() bool {
		updatedAt, err := core.GetMetricsLastUpdateTime(t.Context(), time.Now().UTC())
		return err == nil && !updatedAt.Equal(seededAt)
	}, cadence+500*time.Millisecond, 100*time.Millisecond, "worker should wait while the refresh lock is held")

	// Release the lock. The blocked tick must complete immediately.
	// The assertion timeout (cadence/2 = 1s) is strictly shorter than the next tick interval (2s),
	// proving that the contended flush itself unblocked rather than waiting for another tick.
	locked = false
	cb.RefreshInProgress.Unlock()
	require.Eventually(t, func() bool {
		updatedAt, err := core.GetMetricsLastUpdateTime(t.Context(), time.Now().UTC())
		return err == nil && !updatedAt.Equal(seededAt)
	}, cadence/2, 50*time.Millisecond, "blocked flush should complete immediately without waiting for another tick")
}
