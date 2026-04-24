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
	logicalSsh "github.com/hashicorp/vault/builtin/logical/ssh"
	logicalTransit "github.com/hashicorp/vault/builtin/logical/transit"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/pluginconsts"
	"github.com/hashicorp/vault/helper/timeutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault/billing"
	"github.com/stretchr/testify/require"
)

// TestSystemBackend_BillingOverviewMonthFormat verifies that the billing overview endpoint
// returns the correct response structure with billing.BillingRetentionMonths of data.
// It validates the response format, month strings (YYYY-MM), RFC3339 timestamps,
// and ensures all months are present in the response with correct formatting.
func TestSystemBackend_BillingOverviewMonthFormat(t *testing.T) {
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
	require.Len(t, months, billing.BillingRetentionMonths, "should have billing.BillingRetentionMonths months")

	now := time.Now()
	currentMonthStart := timeutil.StartOfMonth(now)

	// Loop through all months and verify format
	for i := 0; i < billing.BillingRetentionMonths; i++ {
		monthData, ok := months[i].(map[string]interface{})
		require.True(t, ok, "month %d should be a map", i)

		// Verify all required fields are present
		require.Contains(t, monthData, "month", "month %d should have 'month' field", i)
		require.Contains(t, monthData, "updated_at", "month %d should have 'updated_at' field", i)
		require.Contains(t, monthData, "usage_metrics", "month %d should have 'usage_metrics' field", i)

		// Verify month format (YYYY-MM)
		monthStr, ok := monthData["month"].(string)
		require.True(t, ok, "month %d 'month' should be a string", i)
		require.Regexp(t, `^\d{4}-\d{2}$`, monthStr, "month %d should match YYYY-MM format", i)

		// Verify the month string matches expected value
		expectedMonthTime := currentMonthStart.AddDate(0, -i, 0)
		expectedMonth := expectedMonthTime.Format("2006-01")
		require.Equal(t, expectedMonth, monthStr, "month %d should match expected format", i)

		// Verify updated_at format
		updatedAt, ok := monthData["updated_at"].(string)
		require.True(t, ok, "month %d 'updated_at' should be a string", i)
		_, err = time.Parse(time.RFC3339, updatedAt)
		require.NoError(t, err, "month %d updated_at should be valid RFC3339 timestamp", i)

		// Verify usage_metrics is a slice
		_, ok = monthData["usage_metrics"].([]map[string]interface{})
		require.True(t, ok, "month %d usage_metrics should be a slice of maps", i)
	}
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
	require.Len(t, months, billing.BillingRetentionMonths)

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

	// Verify that all previous months (without data) have empty usage_metrics
	currentMonthStart := timeutil.StartOfMonth(currentMonth)
	for i := 1; i < billing.BillingRetentionMonths; i++ {
		monthData, ok := months[i].(map[string]interface{})
		require.True(t, ok, "month %d should be a map", i)

		// Verify month string format
		monthStr, ok := monthData["month"].(string)
		require.True(t, ok, "month %d should have month string", i)
		expectedMonthTime := currentMonthStart.AddDate(0, -i, 0)
		expectedMonth := expectedMonthTime.Format("2006-01")
		require.Equal(t, expectedMonth, monthStr, "month %d should match expected format", i)

		usageMetrics, ok := monthData["usage_metrics"].([]map[string]interface{})
		require.True(t, ok, "month %d should have usage_metrics", i)

		// Previous months without data should have empty metrics or all zeros
		if len(usageMetrics) > 0 {
			// If metrics exist, verify they are all zero
			for _, metric := range usageMetrics {
				metricData, ok := metric["metric_data"].(map[string]interface{})
				if ok {
					total, ok := metricData["total"].(int)
					if ok {
						require.Equal(t, 0, total, "month %d metric total should be 0", i)
					}
				}
			}
		}
	}
}

// TestSystemBackend_BillingOverview_MetricTypeFormat validates that different metric types
// in the billing overview response have the correct data structure.
func TestSystemBackend_BillingOverview_MetricTypeFormat(t *testing.T) {
	c, _, root, _ := TestCoreUnsealedWithMetricsAndConfig(t, &CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			pluginconsts.SecretEngineKV:       logicalKv.Factory,
			pluginconsts.SecretEngineAWS:      logicalAws.Factory,
			pluginconsts.SecretEngineDatabase: logicalDatabase.Factory,
			pluginconsts.SecretEngineTransit:  logicalTransit.Factory,
			pluginconsts.SecretEngineSsh:      logicalSsh.Factory,
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

	// Create SSH certificate and OTP
	req = logical.TestRequest(t, logical.CreateOperation, "sys/mounts/ssh")
	req.Data["type"] = "ssh"
	req.ClientToken = root
	resp, err = c.HandleRequest(ctx, req)
	require.NoError(t, err)

	req = logical.TestRequest(t, logical.CreateOperation, "ssh/config/ca")
	req.ClientToken = root
	resp, err = c.HandleRequest(ctx, req)
	require.NoError(t, err)

	req = logical.TestRequest(t, logical.CreateOperation, "ssh/roles/test-cert")
	req.ClientToken = root
	req.Data["key_type"] = "ca"
	req.Data["allow_user_certificates"] = true
	req.Data["allow_empty_principals"] = true
	req.Data["ttl"] = "1d"
	_, err = c.HandleRequest(ctx, req)
	require.NoError(t, err)

	req = logical.TestRequest(t, logical.UpdateOperation, "ssh/issue/test-cert")
	req.ClientToken = root
	resp, err = c.HandleRequest(ctx, req)
	require.NoError(t, err)
	require.Nil(t, resp.Error())

	req = logical.TestRequest(t, logical.CreateOperation, "ssh/roles/test-otp")
	req.ClientToken = root
	req.Data["key_type"] = "otp"
	req.Data["default_user"] = "user"
	_, err = c.HandleRequest(ctx, req)
	require.NoError(t, err)

	req = logical.TestRequest(t, logical.CreateOperation, "ssh/config/zeroaddress")
	req.ClientToken = root
	req.Data["roles"] = "test-otp"
	_, err = c.HandleRequest(ctx, req)
	require.NoError(t, err)

	req = logical.TestRequest(t, logical.UpdateOperation, "ssh/creds/test-otp")
	req.ClientToken = root
	req.Data["ip"] = "1.2.3.4"
	resp, err = c.HandleRequest(ctx, req)
	require.NoError(t, err)
	require.Nil(t, resp.Error())

	// Update all metrics
	currentMonth := time.Now()
	_, err = c.UpdateMaxKvCounts(ctx, billing.LocalPrefix, currentMonth)
	require.NoError(t, err)

	_, _, err = c.UpdateMaxRoleAndManagedKeyCounts(ctx, billing.LocalPrefix, currentMonth)
	require.NoError(t, err)

	_, err = c.UpdateTransitCallCounts(ctx, currentMonth)
	require.NoError(t, err)

	_, err = c.UpdateStoredSSHDurationAdjustedCertCount(ctx, currentMonth, c.certCountManager.GetCounts().SSHIssuedCerts)
	require.NoError(t, err)

	_, err = c.UpdateStoredSSHOTPCount(ctx, currentMonth, c.certCountManager.GetCounts().SSHIssuedOTPs)
	require.NoError(t, err)

	// Write GCP KMS count directly to storage
	c.consumptionBilling.BillingStorageLock.Lock()
	err = c.storeGcpKmsCallCountsLocked(ctx, uint64(5), billing.LocalPrefix, currentMonth)
	c.consumptionBilling.BillingStorageLock.Unlock()
	require.NoError(t, err)

	// Make a request to the billing overview endpoint
	req = logical.TestRequest(t, logical.ReadOperation, "billing/overview")
	req.Data["refresh_data"] = true
	resp, err = b.HandleRequest(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	months, ok := resp.Data["months"].([]interface{})
	require.True(t, ok)
	require.Len(t, months, billing.BillingRetentionMonths)

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
			require.Equal(t, total, uint64(6)) // 1 transit + 5 gcpkms

			require.Contains(t, metricData, "metric_details")
			metricDetails, ok := metricData["metric_details"].([]map[string]interface{})
			require.True(t, ok, "metric_details should be []map[string]interface{}")
			require.NotEmpty(t, metricDetails)

			foundTransit := false
			foundGcpKms := false
			for _, detail := range metricDetails {
				if detail["type"] == "transit" {
					foundTransit = true
					count, ok := detail["count"].(uint64)
					require.True(t, ok)
					require.Equal(t, count, uint64(1))
				}
				if detail["type"] == "gcpkms" {
					foundGcpKms = true
					count, ok := detail["count"].(uint64)
					require.True(t, ok)
					require.Equal(t, count, uint64(5))
				}
			}
			require.True(t, foundTransit, "should have transit type in metric_details")
			require.True(t, foundGcpKms, "should have gcpkms type in metric_details")

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

		case "ssh_units":
			require.Contains(t, metricData, "total")
			total, ok := metricData["total"].(float64)
			require.True(t, ok, "ssh_units total should be float64")
			require.GreaterOrEqual(t, total, float64(0))

			require.Contains(t, metricData, "metric_details")
			metricDetails, ok := metricData["metric_details"].([]map[string]interface{})
			require.True(t, ok, "metric_details should be []map[string]interface{}")
			require.NotEmpty(t, metricDetails)
			require.Equal(t, len(metricDetails), 2)

			require.Equal(t, metricDetails[0]["type"], "otp_units")
			require.GreaterOrEqual(t, metricDetails[0]["count"], float64(0))

			require.Equal(t, metricDetails[1]["type"], "certificate_units")
			require.GreaterOrEqual(t, metricDetails[1]["count"], float64(0))
		}
	}

	// Verify we found the expected metrics
	require.True(t, metricsFound["static_secrets"], "should have static_secrets metric")
	require.True(t, metricsFound["dynamic_roles"], "should have dynamic_roles metric")
	require.True(t, metricsFound["auto_rotated_roles"], "should have auto_rotated_roles metric")
	require.True(t, metricsFound["data_protection_calls"], "should have data_protection_calls metric")
	require.True(t, metricsFound["pki_units"], "should have pki_units metric")
	require.True(t, metricsFound["managed_keys"], "should have managed_keys metric")
	require.True(t, metricsFound["ssh_units"], "should have ssh_units metric")
}

// TestSystemBackend_BillingOverview_HistoricalMonths verifies that the billing overview
// endpoint correctly retrieves and formats data for all months. It stores
// billing data for the previous month, validates month string formats,
// ensures the updated_at timestamp is set correctly for each month, and confirms
// all month data is included in the response.
func TestSystemBackend_BillingOverview_HistoricalMonths(t *testing.T) {
	c, b, _ := testCoreSystemBackend(t)
	ctx := namespace.RootContext(nil)
	now := time.Now().UTC()
	currentMonth := timeutil.StartOfMonth(now)

	// Store some data for previous month
	previousMonth := timeutil.StartOfPreviousMonth(now)

	// Manually store some counts for previous month
	c.consumptionBilling.BillingStorageLock.Lock()
	err := c.storeMaxKvCountsLocked(ctx, 5, "local/", previousMonth)
	c.consumptionBilling.BillingStorageLock.Unlock()
	require.NoError(t, err)

	// Store metrics last update timestamp for previous month so it's detected as having data
	testUpdateTime := time.Date(previousMonth.Year(), previousMonth.Month(), 15, 12, 0, 0, 0, time.UTC)
	err = c.UpdateMetricsLastUpdateTime(ctx, previousMonth, testUpdateTime)
	require.NoError(t, err)

	// Store metrics last update timestamp for current month
	currentMonthUpdateTime := time.Date(now.Year(), now.Month(), now.Day(), 10, 30, 0, 0, time.UTC)
	err = c.UpdateMetricsLastUpdateTime(ctx, currentMonth, currentMonthUpdateTime)
	require.NoError(t, err)

	// Make a request to the billing overview endpoint
	req := logical.TestRequest(t, logical.ReadOperation, "billing/overview")
	resp, err := b.HandleRequest(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	months, ok := resp.Data["months"].([]interface{})
	require.True(t, ok)
	require.Len(t, months, billing.BillingRetentionMonths)

	// Loop through all months and verify timestamps
	for i := 0; i < billing.BillingRetentionMonths; i++ {
		monthData, ok := months[i].(map[string]interface{})
		require.True(t, ok, "month %d should be a map", i)

		// Verify month string format
		monthStr, ok := monthData["month"].(string)
		require.True(t, ok, "month %d should have month string", i)
		expectedMonthTime := currentMonth.AddDate(0, -i, 0)
		expectedMonth := expectedMonthTime.Format("2006-01")
		require.Equal(t, expectedMonth, monthStr, "month %d should match expected format", i)

		// Verify updated_at timestamp
		updatedAt, ok := monthData["updated_at"].(string)
		require.True(t, ok, "month %d should have updated_at", i)
		parsedTime, err := time.Parse(time.RFC3339, updatedAt)
		require.NoError(t, err, "month %d updated_at should parse", i)

		// Verify timestamps based on which month we're checking
		if i == 0 {
			// Current month should have the timestamp we set
			require.Equal(t, currentMonthUpdateTime, parsedTime, "current month updated_at should match set time")
		} else if i == 1 {
			// Previous month should have timestamp at end of month
			expectedEndOfMonth := timeutil.EndOfMonth(expectedMonthTime)
			require.Equal(t, expectedEndOfMonth, parsedTime, "previous month updated_at should be end of month")
		} else {
			// Older months without data should have zero time
			require.True(t, parsedTime.IsZero() || parsedTime.Equal(time.Time{}),
				"month %d without data should have zero timestamp", i)
		}
	}
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
	require.Len(t, months, billing.BillingRetentionMonths)

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
		"ssh_units":             false,
		"id_token_units":        false,
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
			require.Equal(t, float64(0), total, "pki units total should be 0")

		case "managed_keys":
			total, ok := metricData["total"].(int)
			require.True(t, ok, "managed_keys total should be float64")
			require.Equal(t, int(0), total, "managed keys total should be 0")
			details, ok := metricData["metric_details"].([]map[string]interface{})
			require.True(t, ok, "%s metric_details should be array", metricName)
			require.Empty(t, details, "%s metric_details should be empty when total is 0", metricName)

		case "ssh_units":
			total, ok := metricData["total"].(float64)
			require.True(t, ok, "ssh_units total should be float64")
			require.Equal(t, float64(0), total, "ssh_units total should be 0")

		case "id_token_units":
			total, ok := metricData["total"].(float64)
			require.True(t, ok, "id_token_units total should be float64")
			require.Equal(t, float64(0), total, "id_token_units total should be 0")
		}
	}

	// Verify all expected metrics were found
	for metricName, found := range expectedMetrics {
		require.True(t, found, "metric %s should be present", metricName)
	}

	// Verify all previous months also have zero values
	for i := 1; i < billing.BillingRetentionMonths; i++ {
		monthData, ok := months[i].(map[string]interface{})
		require.True(t, ok, "month %d should be a map", i)
		require.Contains(t, monthData, "usage_metrics", "month %d should have usage_metrics", i)

		usageMetrics, ok := monthData["usage_metrics"].([]map[string]interface{})
		require.True(t, ok, "month %d usage_metrics should be array", i)
		require.NotNil(t, usageMetrics, "month %d usage_metrics should not be nil", i)

		// Verify all metrics in previous months have zero values
		for _, metric := range usageMetrics {
			metricName, ok := metric["metric_name"].(string)
			require.True(t, ok, "month %d metric_name should be a string", i)

			metricData, ok := metric["metric_data"].(map[string]interface{})
			require.True(t, ok, "month %d metric_data should be a map", i)

			// Verify each metric has appropriate zero value
			switch metricName {
			case "static_secrets", "dynamic_roles", "auto_rotated_roles":
				total, ok := metricData["total"].(int)
				require.True(t, ok)
				require.Equal(t, 0, total)

			case "kmip":
				used, ok := metricData["used_in_month"].(bool)
				require.True(t, ok)
				require.False(t, used)

			case "external_plugins":
				total, ok := metricData["total"].(int)
				require.True(t, ok)
				require.Equal(t, 0, total)

			case "data_protection_calls":
				total, ok := metricData["total"].(uint64)
				require.True(t, ok)
				require.Equal(t, uint64(0), total)

			case "pki_units":
				total, ok := metricData["total"].(float64)
				require.True(t, ok)
				require.Equal(t, float64(0), total)

			case "managed_keys":
				total, ok := metricData["total"].(int)
				require.True(t, ok)
				require.Equal(t, 0, total)

			case "ssh_units":
				total, ok := metricData["total"].(float64)
				require.True(t, ok)
				require.Equal(t, float64(0), total)
			}
		}
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
	require.Len(t, months, billing.BillingRetentionMonths)

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

	// First, call with refresh_data set to set the metrics last update timestamp
	req := logical.TestRequest(t, logical.ReadOperation, "billing/overview")
	req.Data["refresh_data"] = true
	resp, err := b.HandleRequest(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	months, ok := resp.Data["months"].([]interface{})
	require.True(t, ok)
	require.Len(t, months, billing.BillingRetentionMonths)

	currentMonth, ok := months[0].(map[string]interface{})
	require.True(t, ok)

	// Get the updated_at timestamp from the first call (current month)
	firstUpdatedAt, ok := currentMonth["updated_at"].(string)
	require.True(t, ok)
	firstTime, err := time.Parse(time.RFC3339, firstUpdatedAt)
	require.NoError(t, err)

	// Verify the metrics last update time was set
	lastUpdate, err := c.GetMetricsLastUpdateTime(ctx, time.Now().UTC())
	require.NoError(t, err)
	require.Equal(t, firstTime, lastUpdate, "stored timestamp should match response timestamp")

	// Verify all previous months have zero timestamp (no data stored for them)
	for i := 1; i < billing.BillingRetentionMonths; i++ {
		prevMonth, ok := months[i].(map[string]interface{})
		require.True(t, ok, "month %d should be a map", i)

		prevMonthUpdatedAt, ok := prevMonth["updated_at"].(string)
		require.True(t, ok, "month %d should have updated_at", i)
		prevMonthTime, err := time.Parse(time.RFC3339, prevMonthUpdatedAt)
		require.NoError(t, err, "month %d updated_at should parse", i)

		// All previous months should be zero time since we haven't stored any data for them
		require.True(t, prevMonthTime.IsZero(),
			"month %d updated_at should be zero time when no data is stored", i)
	}

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
	require.Len(t, months, billing.BillingRetentionMonths)

	currentMonth, ok = months[0].(map[string]interface{})
	require.True(t, ok)

	// Get the updated_at timestamp from the second call (current month)
	secondUpdatedAt, ok := currentMonth["updated_at"].(string)
	require.True(t, ok)
	secondTime, err := time.Parse(time.RFC3339, secondUpdatedAt)
	require.NoError(t, err)

	// The timestamp should be the same as the first call because we didn't refresh the data
	require.Equal(t, firstTime, secondTime,
		"updated_at without refresh should use stored metrics last update timestamp")

	// Verify the timestamps are equal
	require.Equal(t, firstUpdatedAt, secondUpdatedAt,
		"updated_at without refresh should be identical to the stored timestamp")

	// Verify all previous months' timestamps remain the same (zero time)
	for i := 1; i < billing.BillingRetentionMonths; i++ {
		prevMonth, ok := months[i].(map[string]interface{})
		require.True(t, ok, "month %d should be a map", i)

		secondPrevMonthUpdatedAt, ok := prevMonth["updated_at"].(string)
		require.True(t, ok, "month %d should have updated_at", i)
		secondPrevMonthTime, err := time.Parse(time.RFC3339, secondPrevMonthUpdatedAt)
		require.NoError(t, err, "month %d updated_at should parse", i)

		require.True(t, secondPrevMonthTime.IsZero(),
			"month %d updated_at should remain zero time", i)
	}
}

// TestSystemBackend_BillingOverview_UpdatedAtTimestamp_NoStoredTimestamp tests the behavior
// when the metrics last update time is zero time (background worker hasn't run yet)
func TestSystemBackend_BillingOverview_UpdatedAtTimestamp_NoStoredTimestamp(t *testing.T) {
	c, b, _ := testCoreSystemBackend(t)
	ctx := namespace.RootContext(nil)

	// Verify the metrics last update time is zero time initially
	lastUpdate, err := c.GetMetricsLastUpdateTime(ctx, time.Now().UTC())
	require.NoError(t, err)
	require.True(t, lastUpdate.IsZero(), "metrics last update time should be zero time initially")

	// Call without refresh_data when timestamp is zero
	req := logical.TestRequest(t, logical.ReadOperation, "billing/overview")
	req.Data["refresh_data"] = false
	resp, err := b.HandleRequest(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	months, ok := resp.Data["months"].([]interface{})
	require.True(t, ok)
	require.Len(t, months, billing.BillingRetentionMonths)

	currentMonth, ok := months[0].(map[string]interface{})
	require.True(t, ok)

	// Get the updated_at timestamp
	updatedAt, ok := currentMonth["updated_at"].(string)
	require.True(t, ok)
	updatedTime, err := time.Parse(time.RFC3339, updatedAt)
	require.NoError(t, err)

	// Verify it's zero time to indicate data hasn't been updated yet
	require.True(t, updatedTime.IsZero(),
		"updated_at should be zero time when the metrics last update time is zero")

	// Verify previous month is also zero time (no stored timestamp for previous month)
	previousMonth, ok := months[1].(map[string]interface{})
	require.True(t, ok)
	prevMonthUpdatedAt, ok := previousMonth["updated_at"].(string)
	require.True(t, ok)
	prevMonthTime, err := time.Parse(time.RFC3339, prevMonthUpdatedAt)
	require.NoError(t, err)

	// Previous month should also be zero time since no timestamp is stored
	require.True(t, prevMonthTime.IsZero(),
		"previous month updated_at should be zero time when no stored timestamp exists")
}

// TestSystemBackend_BillingOverview_PreviousMonth_WithError tests the behavior
// when retrieving the previous month's timestamp fails with an error.
// This ensures the endpoint gracefully handles storage errors by returning zero time.
func TestSystemBackend_BillingOverview_PreviousMonth_WithError(t *testing.T) {
	c, b, _ := testCoreSystemBackend(t)
	ctx := namespace.RootContext(nil)

	// Store some data for previous month
	previousMonth := timeutil.StartOfPreviousMonth(time.Now())

	// Store counts but intentionally do NOT store the metrics last update timestamp
	// This simulates a scenario where data exists but timestamp retrieval might fail
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
	require.Len(t, months, billing.BillingRetentionMonths)

	// Check previous month data
	previousMonthData, ok := months[1].(map[string]interface{})
	require.True(t, ok)

	// Verify updated_at is zero time when no timestamp is stored
	updatedAt, ok := previousMonthData["updated_at"].(string)
	require.True(t, ok)
	parsedTime, err := time.Parse(time.RFC3339, updatedAt)
	require.NoError(t, err)

	// Should be zero time since no timestamp was stored for previous month
	require.True(t, parsedTime.IsZero(),
		"previous month updated_at should be zero time when timestamp is not stored")
}
