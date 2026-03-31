// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package api

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/vault/api"
	logicalAws "github.com/hashicorp/vault/builtin/logical/aws"
	logicalDatabase "github.com/hashicorp/vault/builtin/logical/database"
	logicalTransit "github.com/hashicorp/vault/builtin/logical/transit"
	"github.com/hashicorp/vault/helper/pluginconsts"
	"github.com/hashicorp/vault/helper/testhelpers/minimal"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
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
	require.Len(t, resp.Months, 2, "should have current and previous month")

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
	require.Len(t, resp.Months, 2, "should have current and previous month")

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
	require.Len(t, resp.Months, 2)

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

	// Verify months are in descending order (current, then previous)
	require.Greater(t, resp.Months[0].Month, resp.Months[1].Month, "first month should be more recent than second")
}
