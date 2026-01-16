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

	// If we are the primary, then combine the replicated and local max role counts. Else just output the local
	// max role counts. replicatedMaxRoleCounts will be empty if we are not a primary, so this is taken care of for us.
	combinedMaxRoleCounts := combineRoleCounts(ctx, replicatedMaxRoleCounts, localMaxRoleCounts)
	combinedMaxKvCounts := replicatedKvHWMCounts + localKvHWMCounts

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

	combinedPreviousMonthRoleCounts := combineRoleCounts(ctx, replicatedPreviousMonthRoleCounts, localPreviousMonthRoleCounts)
	combinedPreviousMonthKvHWMCounts := replicatedPreviousMonthKvHWMCounts + localPreviousMonthKvHWMCounts

	resp := map[string]interface{}{
		"current_month": map[string]interface{}{
			"timestamp":           timeutil.StartOfMonth(currentMonth),
			"maximum_role_counts": combinedMaxRoleCounts,
			"maximum_kv_counts":   combinedMaxKvCounts,
		},
		"previous_month": map[string]interface{}{
			"timestamp":           previousMonth,
			"maximum_role_counts": combinedPreviousMonthRoleCounts,
			"maximum_kv_counts":   combinedPreviousMonthKvHWMCounts,
		},
	}

	return &logical.Response{
		Data: resp,
	}, nil
}
