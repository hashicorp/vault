// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestSys_BillingOverview tests the BillingOverview API client method and structure of the response
func TestSys_BillingOverview(t *testing.T) {
	mockVaultServer := httptest.NewServer(http.HandlerFunc(mockVaultBillingHandler))
	defer mockVaultServer.Close()

	// Create API client pointing to mock server
	cfg := DefaultConfig()
	cfg.Address = mockVaultServer.URL
	client, err := NewClient(cfg)
	require.NoError(t, err)

	resp, err := client.Sys().BillingOverview(false)
	require.NoError(t, err)
	require.NotNil(t, resp)

	// Verify we have 2 months (current and previous)
	require.Len(t, resp.Months, 2)

	// Verify current month structure
	currentMonth := resp.Months[0]
	require.Equal(t, "2026-01", currentMonth.Month)
	require.Equal(t, "2026-01-14T10:49:00Z", currentMonth.UpdatedAt)
	require.Len(t, currentMonth.UsageMetrics, 4)

	// Verify static_secrets metric
	staticSecretsMetric := currentMonth.UsageMetrics[0]
	require.Equal(t, "static_secrets", staticSecretsMetric.MetricName)
	require.NotNil(t, staticSecretsMetric.MetricData)
	require.Contains(t, staticSecretsMetric.MetricData, "total")
	require.Contains(t, staticSecretsMetric.MetricData, "metric_details")

	// Verify kmip metric
	kmipMetric := currentMonth.UsageMetrics[1]
	require.Equal(t, "kmip", kmipMetric.MetricName)
	require.NotNil(t, kmipMetric.MetricData)
	require.Contains(t, kmipMetric.MetricData, "used_in_month")
	require.Equal(t, true, kmipMetric.MetricData["used_in_month"])

	// Verify pki_units metric
	pkiMetric := currentMonth.UsageMetrics[2]
	require.Equal(t, "pki_units", pkiMetric.MetricName)
	require.NotNil(t, pkiMetric.MetricData)
	require.Contains(t, pkiMetric.MetricData, "total")

	// Verify managed_keys metric
	managedKeysMetric := currentMonth.UsageMetrics[3]
	require.Equal(t, "managed_keys", managedKeysMetric.MetricName)
	require.NotNil(t, managedKeysMetric.MetricData)
	require.Contains(t, managedKeysMetric.MetricData, "total")
	require.Contains(t, managedKeysMetric.MetricData, "metric_details")

	// Verify previous month structure
	previousMonth := resp.Months[1]
	require.Equal(t, "2025-12", previousMonth.Month)
	require.Equal(t, "2025-12-31T23:59:59Z", previousMonth.UpdatedAt)
	require.Len(t, previousMonth.UsageMetrics, 1)

	// Verify external_plugins metric in previous month
	externalPluginsMetric := previousMonth.UsageMetrics[0]
	require.Equal(t, "external_plugins", externalPluginsMetric.MetricName)
	require.NotNil(t, externalPluginsMetric.MetricData)
	require.Contains(t, externalPluginsMetric.MetricData, "total")
}

func mockVaultBillingHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte(billingOverviewResponse))
}

const billingOverviewResponse = `{
  "request_id": "d8d3e6e1-4e5f-6a7b-8c9d-0e1f2a3b4c5d",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "months": [
      {
        "month": "2026-01",
        "updated_at": "2026-01-14T10:49:00Z",
        "usage_metrics": [
          {
            "metric_name": "static_secrets",
            "metric_data": {
              "total": 10,
              "metric_details": [
                {
                  "type": "kv",
                  "count": 10
                }
              ]
            }
          },
          {
            "metric_name": "kmip",
            "metric_data": {
              "used_in_month": true
            }
          },
          {
            "metric_name": "pki_units",
            "metric_data": {
              "total": 100.5
            }
          },
          {
            "metric_name": "managed_keys",
            "metric_data": {
              "total": 10,
              "metric_details": [
                {
                  "type": "totp",
                  "count": 5
                },
				{
				  "type": "kmse",
				  "count": 5
				}
              ]
            }
          }
        ]
      },
      {
        "month": "2025-12",
        "updated_at": "2025-12-31T23:59:59Z",
        "usage_metrics": [
          {
            "metric_name": "external_plugins",
            "metric_data": {
              "total": 5
            }
          }
        ]
      }
    ]
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}`
