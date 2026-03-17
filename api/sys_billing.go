// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/mitchellh/mapstructure"
)

// BillingOverview returns billing metrics for the current and previous month.
// If updateCounts is true, the current month's counts will be updated before returning.
// This is an expensive operation that holds locks and should be used sparingly.
func (c *Sys) BillingOverview(updateCounts bool) (*BillingOverviewResponse, error) {
	return c.BillingOverviewWithContext(context.Background(), updateCounts)
}

// BillingOverviewWithContext returns billing metrics for the current and previous month.
func (c *Sys) BillingOverviewWithContext(ctx context.Context, updateCounts bool) (*BillingOverviewResponse, error) {
	ctx, cancelFunc := c.c.withConfiguredTimeout(ctx)
	defer cancelFunc()

	r := c.c.NewRequest(http.MethodGet, "/v1/sys/billing/overview")
	if updateCounts {
		r.Params.Set("refresh_data", "true")
	}

	resp, err := c.c.rawRequestWithContext(ctx, r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	secret, err := ParseSecret(resp.Body)
	if err != nil {
		return nil, err
	}
	if secret == nil || secret.Data == nil {
		return nil, errors.New("data from server response is empty")
	}

	var result BillingOverviewResponse
	err = mapstructure.Decode(secret.Data, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// BillingOverviewResponse represents the response from the billing overview endpoint.
type BillingOverviewResponse struct {
	Months []BillingMonth `json:"months" mapstructure:"months"`
}

// BillingMonth represents billing data for a single month.
type BillingMonth struct {
	Month        string        `json:"month" mapstructure:"month"`
	UpdatedAt    string        `json:"updated_at" mapstructure:"updated_at"`
	UsageMetrics []UsageMetric `json:"usage_metrics" mapstructure:"usage_metrics"`
}

// UsageMetric represents a single usage metric with its data.
type UsageMetric struct {
	MetricName string                 `json:"metric_name" mapstructure:"metric_name"`
	MetricData map[string]interface{} `json:"metric_data" mapstructure:"metric_data"`
}
