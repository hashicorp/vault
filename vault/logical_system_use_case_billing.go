// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package vault

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/vault/helper/timeutil"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault/billing"
)

const pkiDurationAjustedCountMetricName = "pki_units"

func (b *SystemBackend) useCaseConsumptionBillingPaths() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "billing/overview$",
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handleUseCaseConsumption,
					Summary:  "Report the count of secrets and roles for the purposes of use case billing.",
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: http.StatusText(http.StatusOK),
							Fields: map[string]*framework.FieldSchema{
								"high_watermark_role_counts": {
									Type:        framework.TypeMap,
									Description: "High watermark (for this month) role counts for this cluster.",
								},
								"data_protection_call_counts": {
									Type:        framework.TypeMap,
									Description: "Count of data protection calls on this cluster.",
								},
							},
						}},
						http.StatusNoContent: {{
							Description: http.StatusText(http.StatusNoContent),
						}},
						http.StatusBadRequest: {{
							Description: http.StatusText(http.StatusBadRequest),
						}},
						http.StatusInternalServerError: {{
							Description: http.StatusText(http.StatusInternalServerError),
						}},
					},
				},
			},
		},
	}
}

func (b *SystemBackend) handleUseCaseConsumption(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Get HWM role counts
	replicatedMaxRoleCounts := &RoleCounts{}
	replicatedKvHWMCounts := 0
	var err error
	currentMonth := time.Now()
	previousMonth := timeutil.StartOfPreviousMonth(currentMonth)

	// If we are the primary, then we want to get the replicated max role counts. Else we shouldn't retrieve them.
	if b.Core.isPrimary() {
		// We use update instead of Get so that the counts are up to date.
		replicatedMaxRoleCounts, err = b.Core.UpdateMaxRoleCounts(ctx, billing.ReplicatedPrefix, currentMonth)
		if err != nil {
			return nil, fmt.Errorf("error retrieving replicated max role counts: %w", err)
		}
		replicatedKvHWMCounts, err = b.Core.UpdateMaxKvCounts(ctx, billing.ReplicatedPrefix, currentMonth)
		if err != nil {
			return nil, fmt.Errorf("error retrieving replicated max kv counts: %w", err)
		}
	}

	// We always want to get the local max role counts
	// We use update instead of Get so that the counts are up to date.
	localMaxRoleCounts, err := b.Core.UpdateMaxRoleCounts(ctx, billing.LocalPrefix, currentMonth)
	if err != nil {
		return nil, fmt.Errorf("error retrieving local max role counts: %w", err)
	}
	localKvHWMCounts, err := b.Core.UpdateMaxKvCounts(ctx, billing.LocalPrefix, currentMonth)
	if err != nil {
		return nil, fmt.Errorf("error retrieving local max kv counts: %w", err)
	}

	// Data protection call counts are stored to local path only
	// Each cluster tracks its own total requests to avoid double counting
	localTransitCallCounts, err := b.Core.UpdateTransitCallCounts(ctx, currentMonth)
	if err != nil {
		return nil, fmt.Errorf("error retrieving local transit call counts: %w", err)
	}
	localTransformCallCounts, err := b.Core.UpdateTransformCallCounts(ctx, currentMonth)
	if err != nil {
		return nil, fmt.Errorf("error retrieving local transform call counts: %w", err)
	}

	// If we are the primary, then combine the replicated and local max role counts. Else just output the local
	// max role counts. replicatedMaxRoleCounts will be empty if we are not a primary, so this is taken care of for us.
	combinedMaxRoleCounts := combineRoleCounts(replicatedMaxRoleCounts, localMaxRoleCounts)
	combinedMaxKvCounts := replicatedKvHWMCounts + localKvHWMCounts
	// Data protection counts are not combined - each cluster reports its own total
	combinedMaxDataProtectionCallCounts := map[string]interface{}{
		"transit":   localTransitCallCounts,
		"transform": localTransformCallCounts,
	}

	var replicatedPreviousMonthRoleCounts *RoleCounts
	replicatedPreviousMonthKvHWMCounts := 0
	if b.Core.isPrimary() {
		replicatedPreviousMonthRoleCounts, err = b.Core.GetStoredHWMRoleCounts(ctx, billing.ReplicatedPrefix, previousMonth)
		if err != nil {
			return nil, fmt.Errorf("error retrieving replicated max role counts for previous month: %w", err)
		}
		replicatedPreviousMonthKvHWMCounts, err = b.Core.GetStoredHWMKvCounts(ctx, billing.ReplicatedPrefix, previousMonth)
		if err != nil {
			return nil, fmt.Errorf("error retrieving replicated max kv counts for previous month: %w", err)
		}
	}
	localPreviousMonthRoleCounts, err := b.Core.GetStoredHWMRoleCounts(ctx, billing.LocalPrefix, previousMonth)
	if err != nil {
		return nil, fmt.Errorf("error retrieving local max role counts for previous month: %w", err)
	}
	localPreviousMonthKvHWMCounts, err := b.Core.GetStoredHWMKvCounts(ctx, billing.LocalPrefix, previousMonth)
	if err != nil {
		return nil, fmt.Errorf("error retrieving local max kv counts for previous month: %w", err)
	}

	// Data protection counts for previous month
	localPreviousMonthTransitCallCounts, err := b.Core.GetStoredTransitCallCounts(ctx, previousMonth)
	if err != nil {
		return nil, fmt.Errorf("error retrieving local transit call counts for previous month: %w", err)
	}
	localPreviousMonthTransformCallCounts, err := b.Core.GetStoredTransformCallCounts(ctx, previousMonth)
	if err != nil {
		return nil, fmt.Errorf("error retrieving local transform call counts for previous month: %w", err)
	}

	combinedPreviousMonthRoleCounts := combineRoleCounts(replicatedPreviousMonthRoleCounts, localPreviousMonthRoleCounts)
	combinedPreviousMonthKvHWMCounts := replicatedPreviousMonthKvHWMCounts + localPreviousMonthKvHWMCounts
	// Data protection counts are not combined - each cluster reports its own total
	combinedPreviousMonthDataProtectionCallCounts := map[string]interface{}{
		"transit":   localPreviousMonthTransitCallCounts,
		"transform": localPreviousMonthTransformCallCounts,
	}

	resp := map[string]interface{}{
		"current_month": map[string]interface{}{
			"timestamp":                   timeutil.StartOfMonth(currentMonth),
			"maximum_role_counts":         combinedMaxRoleCounts,
			"maximum_kv_counts":           combinedMaxKvCounts,
			"data_protection_call_counts": combinedMaxDataProtectionCallCounts,
		},
		"previous_month": map[string]interface{}{
			"timestamp":                   previousMonth,
			"maximum_role_counts":         combinedPreviousMonthRoleCounts,
			"maximum_kv_counts":           combinedPreviousMonthKvHWMCounts,
			"data_protection_call_counts": combinedPreviousMonthDataProtectionCallCounts,
		},
	}

	return &logical.Response{
		Data: resp,
	}, nil
}

// generatePkiBillingMetric generates the billing metric for PKI duration-adjusted certificate counts.
func (b *SystemBackend) generatePkiBillingMetric(ctx context.Context, month time.Time) (map[string]interface{}, error) {
	count, err := b.Core.GetStoredPkiDurationAdjustedCount(ctx, month)
	if err != nil {
		return nil, fmt.Errorf("error retrieving PKI duration-adjusted count for month: %w", err)
	}

	return map[string]interface{}{
		"metric_name": pkiDurationAjustedCountMetricName,
		"metric_data": map[string]interface{}{
			"total": count,
		},
	}, nil
}
