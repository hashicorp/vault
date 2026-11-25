// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package billing

import (
	"testing"

	"github.com/hashicorp/vault/helper/testhelpers/minimal"
	"github.com/hashicorp/vault/vault/billing"
	"github.com/stretchr/testify/require"
)

// TestBillingLimits_LimitableMetricsMatchOverview fails when the list of
// limitable metrics and the metrics reported by sys/billing/overview drift
// apart. Every overview metric must be limitable except kmip, which is a boolean
// flag rather than a quantity, so adding a metric to the overview forces a
// decision about whether it can be limited.
func TestBillingLimits_LimitableMetricsMatchOverview(t *testing.T) {
	t.Parallel()

	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	resp, err := client.Sys().BillingOverview(false)
	require.NoError(t, err)
	require.NotEmpty(t, resp.Months)

	var overviewMetrics []string
	for _, metric := range resp.Months[0].UsageMetrics {
		if metric.MetricName == "kmip" {
			continue
		}
		overviewMetrics = append(overviewMetrics, metric.MetricName)
	}

	require.ElementsMatch(t, overviewMetrics, billing.LimitableMetrics)
	require.Len(t, billing.LimitableMetrics, 10)
}
