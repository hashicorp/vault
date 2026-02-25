// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package vault

import (
	"context"
	"testing"
	"time"

	logicalKv "github.com/hashicorp/vault-plugin-secrets-kv"
	logicalAws "github.com/hashicorp/vault/builtin/logical/aws"
	logicalDatabase "github.com/hashicorp/vault/builtin/logical/database"
	logicalTransit "github.com/hashicorp/vault/builtin/logical/transit"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/pluginconsts"
	"github.com/hashicorp/vault/helper/timeutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault/billing"
	"github.com/stretchr/testify/require"
)

// TestSystemBackend_BillingOverview verifies that the billing overview endpoint
// returns the correct response structure with current and previous month data.
// It validates the response format, month strings (YYYY-MM), RFC3339 timestamps,
// and ensures both months are present in the response.
func TestSystemBackend_BillingOverview(t *testing.T) {
	_, b, _ := testCoreSystemBackend(t)
	ctx := namespace.RootContext(nil)

	// Make a request to the billing overview endpoint
	req := logical.TestRequest(t, logical.ReadOperation, "billing/overview")
	resp, err := b.HandleRequest(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Data)

	// Verify the response structure
	months, ok := resp.Data["months"].([]interface{})
	require.True(t, ok, "months should be a slice")
	require.Len(t, months, 2, "should have current and previous month")

	// Verify current month structure
	currentMonth, ok := months[0].(map[string]interface{})
	require.True(t, ok, "current month should be a map")
	require.Contains(t, currentMonth, "month")
	require.Contains(t, currentMonth, "updated_at")
	require.Contains(t, currentMonth, "usage_metrics")

	// Verify month format (YYYY-MM)
	monthStr, ok := currentMonth["month"].(string)
	require.True(t, ok)
	require.Regexp(t, `^\d{4}-\d{2}$`, monthStr)

	// Verify updated_at format (RFC3339)
	updatedAt, ok := currentMonth["updated_at"].(string)
	require.True(t, ok)
	_, err = time.Parse(time.RFC3339, updatedAt)
	require.NoError(t, err, "updated_at should be valid RFC3339 timestamp")

	// Verify usage_metrics is a slice
	_, ok = currentMonth["usage_metrics"].([]map[string]interface{})
	require.True(t, ok, "usage_metrics should be a slice of maps")

	// Verify previous month structure
	previousMonth, ok := months[1].(map[string]interface{})
	require.True(t, ok, "previous month should be a map")
	require.Contains(t, previousMonth, "month")
	require.Contains(t, previousMonth, "updated_at")
	require.Contains(t, previousMonth, "usage_metrics")

	// Verify that current month is actually current
	now := time.Now()
	expectedCurrentMonth := now.Format("2006-01")
	require.Equal(t, expectedCurrentMonth, monthStr)

	// Verify that previous month is actually previous
	prevMonthStr, ok := previousMonth["month"].(string)
	require.True(t, ok)
	expectedPreviousMonth := timeutil.StartOfPreviousMonth(now).Format("2006-01")
	require.Equal(t, expectedPreviousMonth, prevMonthStr)
}

// TestSystemBackend_BillingOverview_WithMetrics tests the billing overview endpoint
// with actual KV secrets created to generate billing metrics. It verifies that KV v2
// secrets are properly counted in billing, the static_secrets metric appears in the
// response, the metric_data structure contains total and metric_details and the
// metric_details include the correct type and count.
func TestSystemBackend_BillingOverview_WithMetrics(t *testing.T) {
	c, b, root := testCoreSystemBackend(t)
	ctx := namespace.RootContext(nil)

	// Create some KV secrets to generate metrics
	req := logical.TestRequest(t, logical.UpdateOperation, "mounts/testkv")
	req.Data = map[string]interface{}{
		"type": "kv-v2",
	}
	resp, err := b.HandleRequest(ctx, req)
	require.NoError(t, err)
	require.Nil(t, resp)

	kvReq := logical.TestRequest(t, logical.CreateOperation, "testkv/data/test")
	kvReq.Data["data"] = map[string]interface{}{
		"foo": "bar",
	}
	kvReq.ClientToken = root
	kvResp, err := c.HandleRequest(ctx, kvReq)
	require.NoError(t, err)
	require.NotNil(t, kvResp)

	currentMonth := time.Now()
	_, err = c.UpdateMaxKvCounts(ctx, billing.LocalPrefix, currentMonth)
	require.NoError(t, err)

	req = logical.TestRequest(t, logical.ReadOperation, "billing/overview")
	req.Data["refresh_data"] = true
	resp, err = b.HandleRequest(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	// Verify the response contains metrics
	months, ok := resp.Data["months"].([]interface{})
	require.True(t, ok)
	require.Len(t, months, 2)

	currentMonthData, ok := months[0].(map[string]interface{})
	require.True(t, ok)

	usageMetrics, ok := currentMonthData["usage_metrics"].([]map[string]interface{})
	require.True(t, ok)

	// Check if static_secrets metric exists
	foundStaticSecrets := false
	for _, metric := range usageMetrics {
		if metricName, ok := metric["metric_name"].(string); ok && metricName == "static_secrets" {
			foundStaticSecrets = true

			// Verify metric_data structure
			metricData, ok := metric["metric_data"].(map[string]interface{})
			require.True(t, ok)
			require.Contains(t, metricData, "total")
			require.Contains(t, metricData, "metric_details")

			// Verify total is greater than 0
			total, ok := metricData["total"].(int)
			require.True(t, ok)
			require.Greater(t, total, 0)

			// Verify metric_details structure
			metricDetails, ok := metricData["metric_details"].([]map[string]interface{})
			require.True(t, ok)
			require.NotEmpty(t, metricDetails)

			// Verify first detail has type and count
			firstDetail := metricDetails[0]
			require.Contains(t, firstDetail, "type")
			require.Contains(t, firstDetail, "count")
			require.Equal(t, "kv", firstDetail["type"])

			break
		}
	}
	require.True(t, foundStaticSecrets, "static_secrets metric should be present")
}

// TestSystemBackend_BillingOverview_MetricFormats validates that different metric types
// in the billing overview response have the correct data structure.
func TestSystemBackend_BillingOverview_MetricFormats(t *testing.T) {
	c, _, root, _ := TestCoreUnsealedWithMetricsAndConfig(t, &CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			pluginconsts.SecretEngineKV:       logicalKv.Factory,
			pluginconsts.SecretEngineAWS:      logicalAws.Factory,
			pluginconsts.SecretEngineDatabase: logicalDatabase.Factory,
			pluginconsts.SecretEngineTransit:  logicalTransit.Factory,
		},
	})
	b := c.systemBackend
	ctx := namespace.RootContext(context.Background())

	// Create KV secrets for static_secrets metric
	req := logical.TestRequest(t, logical.UpdateOperation, "mounts/kv-v1")
	req.Data = map[string]interface{}{"type": "kv-v1"}
	resp, err := b.HandleRequest(ctx, req)
	require.NoError(t, err)
	require.Nil(t, resp)

	req = logical.TestRequest(t, logical.UpdateOperation, "mounts/kv-v2")
	req.Data = map[string]interface{}{"type": "kv-v2"}
	resp, err = b.HandleRequest(ctx, req)
	require.NoError(t, err)
	require.Nil(t, resp)

	addKvSecretToStorage(t, ctx, c, "kv-v1", root, "secret1", "kv-v1")
	addKvSecretToStorage(t, ctx, c, "kv-v2", root, "secret2", "kv-v2")

	// Create roles for dynamic_roles and auto_rotated_roles metrics
	req = logical.TestRequest(t, logical.UpdateOperation, "mounts/aws")
	req.Data = map[string]interface{}{"type": "aws"}
	resp, err = b.HandleRequest(ctx, req)
	require.NoError(t, err)
	require.Nil(t, resp)

	addRoleToStorage(t, c, "aws", "role/", 2)
	addRoleToStorage(t, c, "aws", "static-roles/", 1)

	req = logical.TestRequest(t, logical.UpdateOperation, "mounts/database")
	req.Data = map[string]interface{}{"type": "database"}
	resp, err = b.HandleRequest(ctx, req)
	require.NoError(t, err)
	require.Nil(t, resp)

	addRoleToStorage(t, c, "database", "role/", 1)
	addRoleToStorage(t, c, "database", "static-role/", 1)

	// Mount transit backend
	req = logical.TestRequest(t, logical.CreateOperation, "sys/mounts/transit")
	req.Data["type"] = "transit"
	req.ClientToken = root
	_, err = c.HandleRequest(ctx, req)
	require.NoError(t, err)

	// Create an encryption key
	req = logical.TestRequest(t, logical.CreateOperation, "transit/keys/foo")
	req.Data["type"] = "aes256-gcm96"
	req.ClientToken = root
	_, err = c.HandleRequest(ctx, req)
	require.NoError(t, err)

	// Perform encryption on the key
	req = logical.TestRequest(t, logical.UpdateOperation, "transit/encrypt/foo")
	req.Data["plaintext"] = "dGhlIHF1aWNrIGJyb3duIGZveA=="
	req.ClientToken = root
	_, err = c.HandleRequest(ctx, req)
	require.NoError(t, err)

	// Update all metrics
	currentMonth := time.Now()
	_, err = c.UpdateMaxKvCounts(ctx, billing.LocalPrefix, currentMonth)
	require.NoError(t, err)

	_, _, err = c.UpdateMaxRoleAndManagedKeyCounts(ctx, billing.LocalPrefix, currentMonth)
	require.NoError(t, err)

	_, err = c.UpdateTransitCallCounts(ctx, currentMonth)
	require.NoError(t, err)

	// Make a request to the billing overview endpoint
	req = logical.TestRequest(t, logical.ReadOperation, "billing/overview")
	req.Data["refresh_data"] = true
	resp, err = b.HandleRequest(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	months, ok := resp.Data["months"].([]interface{})
	require.True(t, ok)
	require.Len(t, months, 2)

	currentMonthData, ok := months[0].(map[string]interface{})
	require.True(t, ok)

	usageMetrics, ok := currentMonthData["usage_metrics"].([]map[string]interface{})
	require.True(t, ok)
	require.NotEmpty(t, usageMetrics, "usage_metrics should not be empty after creating test data")

	metricsFound := make(map[string]bool)

	// Verify each metric has the correct structure
	for _, metric := range usageMetrics {
		metricName, ok := metric["metric_name"].(string)
		require.True(t, ok, "metric_name should be a string")
		require.NotEmpty(t, metricName)

		metricsFound[metricName] = true

		metricData, ok := metric["metric_data"].(map[string]interface{})
		require.True(t, ok, "metric_data should be a map")
		require.NotEmpty(t, metricData)

		// Different metrics have different structures
		switch metricName {
		case "static_secrets":
			// Should have total and metric_details
			require.Contains(t, metricData, "total")
			require.Contains(t, metricData, "metric_details")

			total, ok := metricData["total"].(int)
			require.True(t, ok)
			require.Equal(t, total, 2)

			metricDetails, ok := metricData["metric_details"].([]map[string]interface{})
			require.True(t, ok)
			require.NotEmpty(t, metricDetails)
			for _, detail := range metricDetails {
				require.Contains(t, detail, "type")
				require.Contains(t, detail, "count")
				count, ok := detail["count"].(int)
				require.True(t, ok)
				require.Equal(t, count, 2)
			}

		case "dynamic_roles":
			// Should have total and metric_details
			require.Contains(t, metricData, "total")
			require.Contains(t, metricData, "metric_details")

			total, ok := metricData["total"].(int)
			require.True(t, ok)
			require.Equal(t, total, 3)

			metricDetails, ok := metricData["metric_details"].([]map[string]interface{})
			require.True(t, ok)
			require.NotEmpty(t, metricDetails)
			for _, detail := range metricDetails {
				require.Contains(t, detail, "type")
				require.Contains(t, detail, "count")
			}

		case "auto_rotated_roles":
			// Should have total and metric_details
			require.Contains(t, metricData, "total")
			require.Contains(t, metricData, "metric_details")

			total, ok := metricData["total"].(int)
			require.True(t, ok)
			require.Equal(t, total, 2)

			metricDetails, ok := metricData["metric_details"].([]map[string]interface{})
			require.True(t, ok)
			require.NotEmpty(t, metricDetails)
			for _, detail := range metricDetails {
				require.Contains(t, detail, "type")
				require.Contains(t, detail, "count")
			}

		case "data_protection_calls":
			require.Contains(t, metricData, "total")
			total, ok := metricData["total"].(uint64)
			require.True(t, ok)
			require.Equal(t, total, uint64(1))

			require.Contains(t, metricData, "metric_details")
			metricDetails, ok := metricData["metric_details"].([]map[string]interface{})
			require.True(t, ok, "metric_details should be []map[string]interface{}")
			require.NotEmpty(t, metricDetails)

			foundTransit := false
			for _, detail := range metricDetails {
				if detail["type"] == "transit" {
					foundTransit = true
					count, ok := detail["count"].(uint64)
					require.True(t, ok)
					require.Equal(t, count, uint64(1))
				}
			}
			require.True(t, foundTransit, "should have transit type in metric_details")

		case "pki_units":
			require.Contains(t, metricData, "total")
			total, ok := metricData["total"].(float64)
			require.True(t, ok, "pki_units total should be float64")
			require.GreaterOrEqual(t, total, float64(0))

		case "managed_keys":
			require.Contains(t, metricData, "total")
			total, ok := metricData["total"].(int)
			require.True(t, ok, "managed_keys total should be int")
			require.GreaterOrEqual(t, total, 0)
			require.Contains(t, metricData, "metric_details")
		}
	}

	// Verify we found the expected metrics
	require.True(t, metricsFound["static_secrets"], "should have static_secrets metric")
	require.True(t, metricsFound["dynamic_roles"], "should have dynamic_roles metric")
	require.True(t, metricsFound["auto_rotated_roles"], "should have auto_rotated_roles metric")
	require.True(t, metricsFound["data_protection_calls"], "should have data_protection_calls metric")
	require.True(t, metricsFound["pki_units"], "should have pki_units metric")
	require.True(t, metricsFound["managed_keys"], "should have managed_keys metric")
}

// TestSystemBackend_BillingOverview_PreviousMonth verifies that the billing overview
// endpoint correctly retrieves and formats data for the previous month. It stores
// billing data for the previous month, validates the previous month string format,
// enures the updated_at timestamp is set to the end of the previous month, and confirms
// the previous month data is included in the response.
func TestSystemBackend_BillingOverview_PreviousMonth(t *testing.T) {
	c, b, _ := testCoreSystemBackend(t)
	ctx := namespace.RootContext(nil)

	// Store some data for previous month
	previousMonth := timeutil.StartOfPreviousMonth(time.Now())

	// Manually store some counts for previous month
	c.consumptionBilling.BillingStorageLock.Lock()
	err := c.storeMaxKvCountsLocked(ctx, 5, "local/", previousMonth)
	c.consumptionBilling.BillingStorageLock.Unlock()
	require.NoError(t, err)

	// Make a request to the billing overview endpoint
	req := logical.TestRequest(t, logical.ReadOperation, "billing/overview")
	resp, err := b.HandleRequest(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	months, ok := resp.Data["months"].([]interface{})
	require.True(t, ok)
	require.Len(t, months, 2)

	// Check previous month data
	previousMonthData, ok := months[1].(map[string]interface{})
	require.True(t, ok)

	monthStr, ok := previousMonthData["month"].(string)
	require.True(t, ok)
	expectedMonth := previousMonth.Format("2006-01")
	require.Equal(t, expectedMonth, monthStr)

	// Verify updated_at is end of previous month
	updatedAt, ok := previousMonthData["updated_at"].(string)
	require.True(t, ok)
	parsedTime, err := time.Parse(time.RFC3339, updatedAt)
	require.NoError(t, err)

	// The updated_at for previous month should be at the end of that month
	expectedEndOfMonth := timeutil.StartOfMonth(previousMonth.AddDate(0, 1, 0)).Add(-time.Second)
	require.WithinDuration(t, expectedEndOfMonth, parsedTime, time.Minute)
}

// TestSystemBackend_BillingOverview_EmptyMetrics verifies that the billing overview
// endpoint returns all metrics with zero values when no billing data exists.
func TestSystemBackend_BillingOverview_EmptyMetrics(t *testing.T) {
	_, b, _ := testCoreSystemBackend(t)
	ctx := namespace.RootContext(nil)

	// Make a request without creating any billable resources
	req := logical.TestRequest(t, logical.ReadOperation, "billing/overview")
	resp, err := b.HandleRequest(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Data)

	// Verify the response structure exists
	months, ok := resp.Data["months"].([]interface{})
	require.True(t, ok)
	require.Len(t, months, 2)

	// Check current month has all metrics with zero values
	currentMonth, ok := months[0].(map[string]interface{})
	require.True(t, ok)
	require.Contains(t, currentMonth, "usage_metrics")

	usageMetrics, ok := currentMonth["usage_metrics"].([]map[string]interface{})
	require.True(t, ok)
	require.NotNil(t, usageMetrics)
	require.NotEmpty(t, usageMetrics, "usage_metrics should contain all metrics even with zero values")

	// Verify all expected metrics are present with zero/false values
	expectedMetrics := map[string]bool{
		"static_secrets":        false,
		"dynamic_roles":         false,
		"auto_rotated_roles":    false,
		"kmip":                  false,
		"external_plugins":      false,
		"data_protection_calls": false,
		"pki_units":             false,
		"managed_keys":          false,
	}

	for _, metric := range usageMetrics {
		metricName, ok := metric["metric_name"].(string)
		require.True(t, ok, "metric_name should be a string")
		require.Contains(t, expectedMetrics, metricName, "unexpected metric: %s", metricName)
		expectedMetrics[metricName] = true

		metricData, ok := metric["metric_data"].(map[string]interface{})
		require.True(t, ok, "metric_data should be a map")

		// Verify each metric has appropriate zero value
		switch metricName {
		case "static_secrets", "dynamic_roles", "auto_rotated_roles":
			total, ok := metricData["total"].(int)
			require.True(t, ok, "%s total should be int", metricName)
			require.Equal(t, 0, total, "%s total should be 0", metricName)

			details, ok := metricData["metric_details"].([]map[string]interface{})
			require.True(t, ok, "%s metric_details should be array", metricName)
			require.Empty(t, details, "%s metric_details should be empty when total is 0", metricName)

		case "kmip":
			used, ok := metricData["used_in_month"].(bool)
			require.True(t, ok, "kmip used_in_month should be bool")
			require.False(t, used, "kmip should be false when not used")

		case "external_plugins":
			total, ok := metricData["total"].(int)
			require.True(t, ok, "external_plugins total should be int")
			require.Equal(t, 0, total, "external_plugins total should be 0")

		case "data_protection_calls":
			total, ok := metricData["total"].(uint64)
			require.True(t, ok, "data_protection_calls total should be uint64")
			require.Equal(t, uint64(0), total, "data_protection_calls total should be 0")

			details, ok := metricData["metric_details"].([]map[string]interface{})
			require.True(t, ok, "data_protection_calls metric_details should be array")
			require.Empty(t, details, "data_protection_calls metric_details should be empty when total is 0")

		case "pki_units":
			total, ok := metricData["total"].(float64)
			require.True(t, ok, "pki_units total should be float64")
			require.Equal(t, float64(0), total, "data_protection_calls total should be 0")

		case "managed_keys":
			total, ok := metricData["total"].(int)
			require.True(t, ok, "managed_keys total should be float64")
			require.Equal(t, int(0), total, "data_protection_calls total should be 0")
			details, ok := metricData["metric_details"].([]map[string]interface{})
			require.True(t, ok, "%s metric_details should be array", metricName)
			require.Empty(t, details, "%s metric_details should be empty when total is 0", metricName)
		}
	}

	// Verify all expected metrics were found
	for metricName, found := range expectedMetrics {
		require.True(t, found, "metric %s should be present", metricName)
	}
}

// TestSystemBackend_BillingOverview_MultipleMetricTypes tests the billing overview
// endpoint with multiple different metric types to ensure they all appear correctly
// in the response with their respective data structures.
func TestSystemBackend_BillingOverview_MultipleMetricTypes(t *testing.T) {
	c, b, root := testCoreSystemBackend(t)
	ctx := namespace.RootContext(nil)
	currentMonth := time.Now()

	// Create KV secrets
	req := logical.TestRequest(t, logical.UpdateOperation, "mounts/testkv")
	req.Data = map[string]interface{}{
		"type": "kv-v2",
	}
	resp, err := b.HandleRequest(ctx, req)
	require.NoError(t, err)
	require.Nil(t, resp)

	kvReq := logical.TestRequest(t, logical.CreateOperation, "testkv/data/test")
	kvReq.Data["data"] = map[string]interface{}{"foo": "bar"}
	kvReq.ClientToken = root
	_, err = c.HandleRequest(ctx, kvReq)
	require.NoError(t, err)

	_, err = c.UpdateMaxKvCounts(ctx, billing.LocalPrefix, currentMonth)
	require.NoError(t, err)

	// Make request to billing overview
	req = logical.TestRequest(t, logical.ReadOperation, "billing/overview")
	req.Data["refresh_data"] = true
	resp, err = b.HandleRequest(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	months, ok := resp.Data["months"].([]interface{})
	require.True(t, ok)
	require.Len(t, months, 2)

	currentMonthData, ok := months[0].(map[string]interface{})
	require.True(t, ok)

	usageMetrics, ok := currentMonthData["usage_metrics"].([]map[string]interface{})
	require.True(t, ok)
	require.NotEmpty(t, usageMetrics)

	// Verify each metric has proper structure
	for _, metric := range usageMetrics {
		metricName, ok := metric["metric_name"].(string)
		require.True(t, ok, "metric_name should be a string")
		require.NotEmpty(t, metricName, "metric_name should not be empty")

		metricData, ok := metric["metric_data"].(map[string]interface{})
		require.True(t, ok, "metric_data should be a map")
		require.NotEmpty(t, metricData, "metric_data should not be empty")
	}

	// Verify we have at least the static_secrets metric from our KV secret
	foundStaticSecrets := false
	for _, metric := range usageMetrics {
		if metricName, ok := metric["metric_name"].(string); ok && metricName == "static_secrets" {
			foundStaticSecrets = true
			break
		}
	}
	require.True(t, foundStaticSecrets, "should have static_secrets metric from KV secret")
}

// TestSystemBackend_BillingOverview_UpdatedAtTimestamp verifies that the updated_at
// timestamp behaves correctly based on whether data was refreshed.
func TestSystemBackend_BillingOverview_UpdatedAtTimestamp(t *testing.T) {
	c, b, _ := testCoreSystemBackend(t)
	ctx := namespace.RootContext(nil)

	// First, call with refresh_data set to set the LastMetricsUpdate timestamp
	req := logical.TestRequest(t, logical.ReadOperation, "billing/overview")
	req.Data["refresh_data"] = true
	resp, err := b.HandleRequest(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	months, ok := resp.Data["months"].([]interface{})
	require.True(t, ok)
	require.Len(t, months, 2)

	currentMonth, ok := months[0].(map[string]interface{})
	require.True(t, ok)

	// Get the updated_at timestamp from the first call
	firstUpdatedAt, ok := currentMonth["updated_at"].(string)
	require.True(t, ok)
	firstTime, err := time.Parse(time.RFC3339, firstUpdatedAt)
	require.NoError(t, err)

	// Verify LastMetricsUpdate was set
	lastUpdate := c.consumptionBilling.LastMetricsUpdate.Load()
	require.NotNil(t, lastUpdate, "LastMetricsUpdate should be set after refresh")
	storedTime, ok := lastUpdate.(time.Time)
	require.True(t, ok)
	require.WithinDuration(t, firstTime, storedTime, time.Second, "stored timestamp should match response timestamp")

	// Wait a moment to ensure time difference
	time.Sleep(100 * time.Millisecond)

	// Now call without refresh_data
	req = logical.TestRequest(t, logical.ReadOperation, "billing/overview")
	req.Data["refresh_data"] = false
	resp, err = b.HandleRequest(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	months, ok = resp.Data["months"].([]interface{})
	require.True(t, ok)
	require.Len(t, months, 2)

	currentMonth, ok = months[0].(map[string]interface{})
	require.True(t, ok)

	// Get the updated_at timestamp from the second call
	secondUpdatedAt, ok := currentMonth["updated_at"].(string)
	require.True(t, ok)
	secondTime, err := time.Parse(time.RFC3339, secondUpdatedAt)
	require.NoError(t, err)

	// The timestamp should be the same as the first call because we didn't refresh the data
	require.WithinDuration(t, firstTime, secondTime, time.Second,
		"updated_at without refresh should use stored LastMetricsUpdate timestamp")

	// Verify the timestamps are equal
	require.Equal(t, firstUpdatedAt, secondUpdatedAt,
		"updated_at without refresh should be identical to the stored timestamp")
}
