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

// TestSystemBackend_BillingOverview_StartEndMonthParams tests the billing overview
// endpoint with different combinations of start_month and end_month parameters. It
// verifies that the correct range of months is returned along with any expected warnings
// or errors.
func TestSystemBackend_BillingOverview_StartEndMonthParams(t *testing.T) {
	now := time.Now().UTC()
	currentMonth := now.Format("2006-01")
	previousMonth := timeutil.StartOfPreviousMonth(now).Format("2006-01")
	nextMonth := timeutil.StartOfNextMonth(now).Format("2006-01")
	twoMonthsAfterCurrent := timeutil.StartOfMonth(now).AddDate(0, 2, 0).Format("2006-01")
	retentionStart := timeutil.StartOfMonth(now).AddDate(0, -billing.BillingRetentionMonths+1, 0).Format("2006-01")
	beforeRetentionStart := timeutil.StartOfMonth(now).AddDate(0, -billing.BillingRetentionMonths, 0).Format("2006-01")
	twoMonthsBeforeRetentionStart := timeutil.StartOfMonth(now).AddDate(0, -billing.BillingRetentionMonths-1, 0).Format("2006-01")

	testCases := []struct {
		name            string
		startMonth      interface{}
		endMonth        interface{}
		expectedMonths  int
		expectedWarning string
		expectedError   string
	}{
		{
			name:           "start and end in retention period",
			startMonth:     previousMonth,
			endMonth:       currentMonth,
			expectedMonths: 2,
		},
		{
			name:            "start before retention period, default end",
			startMonth:      beforeRetentionStart,
			expectedMonths:  billing.BillingRetentionMonths + 1,
			expectedWarning: WarningStartEndMonthOutOfRetentionRange,
		},
		{
			name:            "end after retention period, default start",
			endMonth:        nextMonth,
			expectedMonths:  billing.BillingRetentionMonths + 1,
			expectedWarning: WarningStartEndMonthOutOfRetentionRange,
		},
		{
			name:           "start is exactly start of retention period",
			startMonth:     retentionStart,
			endMonth:       previousMonth,
			expectedMonths: billing.BillingRetentionMonths - 1,
		},
		{
			name:            "start and end after retention period",
			startMonth:      nextMonth,
			endMonth:        twoMonthsAfterCurrent,
			expectedMonths:  2,
			expectedWarning: WarningStartEndMonthOutOfRetentionRange,
		},
		{
			name:            "start and end before retention period",
			startMonth:      twoMonthsBeforeRetentionStart,
			endMonth:        beforeRetentionStart,
			expectedMonths:  2,
			expectedWarning: WarningStartEndMonthOutOfRetentionRange,
		},
		{
			name:          "start after retention period, default end",
			startMonth:    nextMonth,
			expectedError: "start_month is later than end_month",
		},
		{
			name:           "no parameters, default start and end",
			expectedMonths: billing.BillingRetentionMonths,
		},
		{
			name:          "start after end",
			startMonth:    previousMonth,
			endMonth:      retentionStart,
			expectedError: "start_month is later than end_month",
		},
		{
			name:           "same month",
			startMonth:     currentMonth,
			endMonth:       currentMonth,
			expectedMonths: 1,
		},
		{
			name:          "invalid date format",
			startMonth:    "2023/01",
			endMonth:      previousMonth,
			expectedError: "invalid start_month format",
		},
		{
			name:          "invalid month",
			startMonth:    "2023-13",
			endMonth:      previousMonth,
			expectedError: "invalid start_month format",
		},
		{
			name:          "invalid data type",
			startMonth:    previousMonth,
			endMonth:      45,
			expectedError: "invalid end_month format",
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			_, b, _ := testCoreSystemBackend(t)
			ctx := namespace.RootContext(nil)

			req := logical.TestRequest(t, logical.ReadOperation, "billing/overview")
			req.Data["start_month"] = test.startMonth
			req.Data["end_month"] = test.endMonth
			resp, err := b.HandleRequest(ctx, req)

			if test.expectedError != "" {
				require.Nil(t, resp)
				require.Error(t, err)
				require.Contains(t, err.Error(), test.expectedError)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)

			if test.expectedWarning != "" {
				require.NotEmpty(t, resp.Warnings)
				require.Contains(t, resp.Warnings, test.expectedWarning)
			} else {
				require.Empty(t, resp.Warnings)
			}

			// Verify the correct number of months are returned
			months := resp.Data["months"].([]interface{})
			require.Len(t, months, test.expectedMonths)

			// expected start and end months are the test parameters if specified,
			// or default to the retention start and current month
			var expectedStartMonth, expectedEndMonth string
			if test.startMonth != nil {
				expectedStartMonth = test.startMonth.(string)
			} else {
				expectedStartMonth = retentionStart
			}
			if test.endMonth != nil {
				expectedEndMonth = test.endMonth.(string)
			} else {
				expectedEndMonth = currentMonth
			}

			// Months are ordered from most recent to oldest, so the first month returned
			// should be the expected endMonth and the last month the expected startMonth
			firstMonth, ok := months[0].(map[string]interface{})
			require.True(t, ok)
			firstMonthStr, ok := firstMonth["month"].(string)
			require.True(t, ok)
			require.Equal(t, expectedEndMonth, firstMonthStr)

			lastMonth, ok := months[len(months)-1].(map[string]interface{})
			require.True(t, ok)
			lastMonthStr, ok := lastMonth["month"].(string)
			require.True(t, ok)
			require.Equal(t, expectedStartMonth, lastMonthStr)
		})
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
		case "static_secrets":
			total, ok := metricData["total"].(int)
			require.True(t, ok, "%s total should be int", metricName)
			require.Equal(t, 0, total, "%s total should be 0", metricName)

			details, ok := metricData["metric_details"].([]map[string]interface{})
			require.True(t, ok, "%s metric_details should be array", metricName)
			require.NotEmpty(t, details, "%s metric_details should always be present", metricName)
			// Verify kv type is present with zero count
			require.Len(t, details, 1)
			require.Equal(t, "kv", details[0]["type"])
			require.Equal(t, 0, details[0]["count"])

		case "dynamic_roles", "auto_rotated_roles":
			total, ok := metricData["total"].(int)
			require.True(t, ok, "%s total should be int", metricName)
			require.Equal(t, 0, total, "%s total should be 0", metricName)

			details, ok := metricData["metric_details"].([]map[string]interface{})
			require.True(t, ok, "%s metric_details should be array", metricName)
			require.NotEmpty(t, details, "%s metric_details should always be present", metricName)
			// Verify all role types are present with zero counts
			for _, detail := range details {
				require.Contains(t, detail, "type")
				require.Contains(t, detail, "count")
				require.Equal(t, 0, detail["count"])
			}

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
			require.NotEmpty(t, details, "data_protection_calls metric_details should always be present")
			// Verify all data protection types are present with zero counts
			require.Len(t, details, 3)
			expectedTypes := map[string]bool{"transit": false, "transform": false, "gcpkms": false}
			for _, detail := range details {
				detailType := detail["type"].(string)
				expectedTypes[detailType] = true
				require.Equal(t, uint64(0), detail["count"])
			}
			for typeName, found := range expectedTypes {
				require.True(t, found, "type %s should be present", typeName)
			}

		case "pki_units":
			total, ok := metricData["total"].(float64)
			require.True(t, ok, "pki_units total should be float64")
			require.Equal(t, float64(0), total, "pki units total should be 0")

		case "managed_keys":
			total, ok := metricData["total"].(int)
			require.True(t, ok, "managed_keys total should be int")
			require.Equal(t, int(0), total, "managed keys total should be 0")
			details, ok := metricData["metric_details"].([]map[string]interface{})
			require.True(t, ok, "%s metric_details should be array", metricName)
			require.NotEmpty(t, details, "%s metric_details should always be present", metricName)
			// Verify both managed key types are present with zero counts
			require.Len(t, details, 2)
			expectedTypes := map[string]bool{"totp": false, "kmse": false}
			for _, detail := range details {
				detailType := detail["type"].(string)
				expectedTypes[detailType] = true
				require.Equal(t, 0, detail["count"])
			}
			for typeName, found := range expectedTypes {
				require.True(t, found, "type %s should be present", typeName)
			}

		case "ssh_units":
			total, ok := metricData["total"].(float64)
			require.True(t, ok, "ssh_units total should be float64")
			require.Equal(t, float64(0), total, "ssh_units total should be 0")

		case "id_token_units":
			total, ok := metricData["total"].(float64)
			require.True(t, ok, "id_token_units total should be float64")
			require.Equal(t, float64(0), total, "id_token_units total should be 0")

			details, ok := metricData["metric_details"].([]map[string]interface{})
			require.True(t, ok, "id_token_units metric_details should be array")
			require.NotEmpty(t, details, "id_token_units metric_details should always be present")
			// Verify both token types are present with zero counts
			require.Len(t, details, 2)
			expectedTypes := map[string]bool{"oidc": false, "spiffe": false}
			for _, detail := range details {
				detailType := detail["type"].(string)
				expectedTypes[detailType] = true
				require.Equal(t, float64(0), detail["count"])
			}
			for typeName, found := range expectedTypes {
				require.True(t, found, "type %s should be present", typeName)
			}
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

// TestSystemBackend_BillingOverview_AllMetricTypesPresent verifies that all metric types
// are always present in the response, even when their counts are zero. This test specifically
// validates that metric_details arrays contain all expected types for each metric category.
func TestSystemBackend_BillingOverview_AllMetricTypesPresent(t *testing.T) {
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

	// Check current month has all metrics
	currentMonth, ok := months[0].(map[string]interface{})
	require.True(t, ok)
	require.Contains(t, currentMonth, "usage_metrics")

	usageMetrics, ok := currentMonth["usage_metrics"].([]map[string]interface{})
	require.True(t, ok)
	require.NotNil(t, usageMetrics)
	require.NotEmpty(t, usageMetrics, "usage_metrics should contain all metrics even with zero values")

	// Build a map of metrics for easy lookup
	metricsMap := make(map[string]map[string]interface{})
	for _, metric := range usageMetrics {
		metricName, ok := metric["metric_name"].(string)
		require.True(t, ok, "metric_name should be a string")
		metricsMap[metricName] = metric
	}

	// Verify static_secrets has kv type
	staticSecretsMetric, exists := metricsMap["static_secrets"]
	require.True(t, exists, "static_secrets metric should be present")
	staticSecretsData := staticSecretsMetric["metric_data"].(map[string]interface{})
	staticSecretsDetails := staticSecretsData["metric_details"].([]map[string]interface{})
	require.Len(t, staticSecretsDetails, 1, "static_secrets should have 1 type")
	require.Equal(t, "kv", staticSecretsDetails[0]["type"])
	require.Equal(t, 0, staticSecretsDetails[0]["count"])

	// Verify dynamic_roles has all 13 types
	dynamicRolesMetric, exists := metricsMap["dynamic_roles"]
	require.True(t, exists, "dynamic_roles metric should be present")
	dynamicRolesData := dynamicRolesMetric["metric_data"].(map[string]interface{})
	dynamicRolesDetails := dynamicRolesData["metric_details"].([]map[string]interface{})
	require.Len(t, dynamicRolesDetails, 13, "dynamic_roles should have 13 types")

	expectedDynamicTypes := []string{
		"aws_dynamic", "azure_dynamic", "database_dynamic", "gcp_dynamic",
		"ldap_dynamic", "openldap_dynamic", "alicloud_dynamic", "rabbitmq_dynamic",
		"consul_dynamic", "nomad_dynamic", "kubernetes_dynamic", "mongodbatlas_dynamic",
		"terraform_dynamic",
	}
	for i, expectedType := range expectedDynamicTypes {
		require.Equal(t, expectedType, dynamicRolesDetails[i]["type"], "dynamic role type at index %d should be %s", i, expectedType)
		require.Equal(t, 0, dynamicRolesDetails[i]["count"], "dynamic role count at index %d should be 0", i)
	}

	// Verify auto_rotated_roles has all 8 types
	autoRotatedMetric, exists := metricsMap["auto_rotated_roles"]
	require.True(t, exists, "auto_rotated_roles metric should be present")
	autoRotatedData := autoRotatedMetric["metric_data"].(map[string]interface{})
	autoRotatedDetails := autoRotatedData["metric_details"].([]map[string]interface{})
	require.Len(t, autoRotatedDetails, 8, "auto_rotated_roles should have 8 types")

	expectedAutoRotatedTypes := []string{
		"aws_static", "azure_static", "database_static", "gcp_static",
		"gcp_impersonated", "ldap_static", "openldap_static", "os_local_account_static",
	}
	for i, expectedType := range expectedAutoRotatedTypes {
		require.Equal(t, expectedType, autoRotatedDetails[i]["type"], "auto-rotated role type at index %d should be %s", i, expectedType)
		require.Equal(t, 0, autoRotatedDetails[i]["count"], "auto-rotated role count at index %d should be 0", i)
	}

	// Verify data_protection_calls has all 3 types
	dataProtectionMetric, exists := metricsMap["data_protection_calls"]
	require.True(t, exists, "data_protection_calls metric should be present")
	dataProtectionData := dataProtectionMetric["metric_data"].(map[string]interface{})
	dataProtectionDetails := dataProtectionData["metric_details"].([]map[string]interface{})
	require.Len(t, dataProtectionDetails, 3, "data_protection_calls should have 3 types")

	expectedDataProtectionTypes := []string{"transit", "transform", "gcpkms"}
	for i, expectedType := range expectedDataProtectionTypes {
		require.Equal(t, expectedType, dataProtectionDetails[i]["type"], "data protection type at index %d should be %s", i, expectedType)
		require.Equal(t, uint64(0), dataProtectionDetails[i]["count"], "data protection count at index %d should be 0", i)
	}

	// Verify managed_keys has both types
	managedKeysMetric, exists := metricsMap["managed_keys"]
	require.True(t, exists, "managed_keys metric should be present")
	managedKeysData := managedKeysMetric["metric_data"].(map[string]interface{})
	managedKeysDetails := managedKeysData["metric_details"].([]map[string]interface{})
	require.Len(t, managedKeysDetails, 2, "managed_keys should have 2 types")

	expectedManagedKeyTypes := []string{"totp", "kmse"}
	for i, expectedType := range expectedManagedKeyTypes {
		require.Equal(t, expectedType, managedKeysDetails[i]["type"], "managed key type at index %d should be %s", i, expectedType)
		require.Equal(t, 0, managedKeysDetails[i]["count"], "managed key count at index %d should be 0", i)
	}

	// Verify ssh_units has both types
	sshMetric, exists := metricsMap["ssh_units"]
	require.True(t, exists, "ssh_units metric should be present")
	sshData := sshMetric["metric_data"].(map[string]interface{})
	sshDetails := sshData["metric_details"].([]map[string]interface{})
	require.Len(t, sshDetails, 2, "ssh_units should have 2 types")
	require.Equal(t, "otp_units", sshDetails[0]["type"])
	require.Equal(t, "certificate_units", sshDetails[1]["type"])

	// Verify id_token_units has both types
	idTokenMetric, exists := metricsMap["id_token_units"]
	require.True(t, exists, "id_token_units metric should be present")
	idTokenData := idTokenMetric["metric_data"].(map[string]interface{})
	idTokenDetails := idTokenData["metric_details"].([]map[string]interface{})
	require.Len(t, idTokenDetails, 2, "id_token_units should have 2 types")

	expectedIdTokenTypes := []string{"oidc", "spiffe"}
	for i, expectedType := range expectedIdTokenTypes {
		require.Equal(t, expectedType, idTokenDetails[i]["type"], "id token type at index %d should be %s", i, expectedType)
		require.Equal(t, float64(0), idTokenDetails[i]["count"], "id token count at index %d should be 0", i)
	}
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

// TestRoundToFour tests the roundToFour function
func TestRoundToFour(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		input    float64
		expected float64
	}{
		{"Round up", 1.23456, 1.2346},
		{"Round down", 1.23454, 1.2345},
		{"Exactly four decimals", 1.1111, 1.1111},
		{"Fewer than four decimals", 1.2, 1.2000},
		{"Zero value", 0.0, 0.0},
		{"Large values", 0.189900000000, 0.1899},
		{"Large values with round up", 0.189990000000, 0.1900},
		{"Large values with round down", 0.189920000000, 0.1899},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := roundToFour(tt.input)
			require.Equal(t, tt.expected, got)
		})
	}
}

// TestRoundUsageMetrics tests the roundUsageMetrics function.
func TestRoundUsageMetrics(t *testing.T) {
	tests := []struct {
		name     string
		input    []map[string]interface{}
		expected []map[string]interface{}
	}{
		{
			name: "Round float64 totals and counts in metric_details",
			input: []map[string]interface{}{
				{
					"metric_name": "pki_units",
					"metric_data": map[string]interface{}{
						"total": 123.456789,
					},
				},
				{
					"metric_name": "ssh_units",
					"metric_data": map[string]interface{}{
						"total": 98.765432,
						"metric_details": []map[string]interface{}{
							{"type": "otp_units", "count": 45.678901},
							{"type": "certificate_units", "count": 53.086531},
						},
					},
				},
			},
			expected: []map[string]interface{}{
				{
					"metric_name": "pki_units",
					"metric_data": map[string]interface{}{
						"total": 123.4568,
					},
				},
				{
					"metric_name": "ssh_units",
					"metric_data": map[string]interface{}{
						"total": 98.7654,
						"metric_details": []map[string]interface{}{
							{"type": "otp_units", "count": 45.6789},
							{"type": "certificate_units", "count": 53.0865},
						},
					},
				},
			},
		},
		{
			name: "Handle integer counts (should not be modified)",
			input: []map[string]interface{}{
				{
					"metric_name": "static_secrets",
					"metric_data": map[string]interface{}{
						"total": 100,
						"metric_details": []map[string]interface{}{
							{"type": "kv", "count": 100},
						},
					},
				},
			},
			expected: []map[string]interface{}{
				{
					"metric_name": "static_secrets",
					"metric_data": map[string]interface{}{
						"total": 100,
						"metric_details": []map[string]interface{}{
							{"type": "kv", "count": 100},
						},
					},
				},
			},
		},
		{
			name: "Handle mixed float64 and integer values",
			input: []map[string]interface{}{
				{
					"metric_name": "id_token_units",
					"metric_data": map[string]interface{}{
						"total": 150.123456,
						"metric_details": []map[string]interface{}{
							{"type": "oidc", "count": 100.987654},
							{"type": "spiffe", "count": 49.135802},
						},
					},
				},
				{
					"metric_name": "dynamic_roles",
					"metric_data": map[string]interface{}{
						"total": 50,
						"metric_details": []map[string]interface{}{
							{"type": "aws_dynamic", "count": 25},
							{"type": "azure_dynamic", "count": 25},
						},
					},
				},
			},
			expected: []map[string]interface{}{
				{
					"metric_name": "id_token_units",
					"metric_data": map[string]interface{}{
						"total": 150.1235,
						"metric_details": []map[string]interface{}{
							{"type": "oidc", "count": 100.9877},
							{"type": "spiffe", "count": 49.1358},
						},
					},
				},
				{
					"metric_name": "dynamic_roles",
					"metric_data": map[string]interface{}{
						"total": 50,
						"metric_details": []map[string]interface{}{
							{"type": "aws_dynamic", "count": 25},
							{"type": "azure_dynamic", "count": 25},
						},
					},
				},
			},
		},
		{
			name: "Handle metrics without metric_details",
			input: []map[string]interface{}{
				{
					"metric_name": "kmip",
					"metric_data": map[string]interface{}{
						"used_in_month": true,
					},
				},
				{
					"metric_name": "external_plugins",
					"metric_data": map[string]interface{}{
						"total": 5,
					},
				},
			},
			expected: []map[string]interface{}{
				{
					"metric_name": "kmip",
					"metric_data": map[string]interface{}{
						"used_in_month": true,
					},
				},
				{
					"metric_name": "external_plugins",
					"metric_data": map[string]interface{}{
						"total": 5,
					},
				},
			},
		},
		{
			name:     "Handle empty metrics slice",
			input:    []map[string]interface{}{},
			expected: []map[string]interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make a deep copy of input to avoid modifying the test case
			inputCopy := make([]map[string]interface{}, len(tt.input))
			for i, metric := range tt.input {
				inputCopy[i] = make(map[string]interface{})
				for k, v := range metric {
					if k == "metric_data" {
						metricData := v.(map[string]interface{})
						metricDataCopy := make(map[string]interface{})
						for mk, mv := range metricData {
							if mk == "metric_details" {
								details := mv.([]map[string]interface{})
								detailsCopy := make([]map[string]interface{}, len(details))
								for di, detail := range details {
									detailsCopy[di] = make(map[string]interface{})
									for dk, dv := range detail {
										detailsCopy[di][dk] = dv
									}
								}
								metricDataCopy[mk] = detailsCopy
							} else {
								metricDataCopy[mk] = mv
							}
						}
						inputCopy[i][k] = metricDataCopy
					} else {
						inputCopy[i][k] = v
					}
				}
			}

			// Apply rounding
			roundUsageMetrics(inputCopy)

			// Verify the results
			require.Equal(t, len(tt.expected), len(inputCopy))
			for i, expectedMetric := range tt.expected {
				actualMetric := inputCopy[i]
				require.Equal(t, expectedMetric["metric_name"], actualMetric["metric_name"])

				expectedData := expectedMetric["metric_data"].(map[string]interface{})
				actualData := actualMetric["metric_data"].(map[string]interface{})

				// Check total
				if expectedTotal, ok := expectedData["total"]; ok {
					require.Equal(t, expectedTotal, actualData["total"])
				}

				// Check metric_details
				if expectedDetails, ok := expectedData["metric_details"].([]map[string]interface{}); ok {
					actualDetails := actualData["metric_details"].([]map[string]interface{})
					require.Equal(t, len(expectedDetails), len(actualDetails))
					for j, expectedDetail := range expectedDetails {
						actualDetail := actualDetails[j]
						require.Equal(t, expectedDetail["type"], actualDetail["type"])
						require.Equal(t, expectedDetail["count"], actualDetail["count"])
					}
				}

				// Check other fields (like used_in_month)
				for key, expectedValue := range expectedData {
					if key != "total" && key != "metric_details" {
						require.Equal(t, expectedValue, actualData[key])
					}
				}
			}
		})
	}
}
