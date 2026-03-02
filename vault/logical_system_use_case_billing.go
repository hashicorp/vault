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

const (
	WarningRefreshIgnoredOnStandby = "refresh_data parameter is supported only on the active node. " +
		"Since this parameter was set on a performance standby, the billing data was not refreshed " +
		"and retrieved from storage without update."
)

func (b *SystemBackend) useCaseConsumptionBillingPaths() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "billing/overview$",
			Fields: map[string]*framework.FieldSchema{
				"refresh_data": {
					Type:        framework.TypeBool,
					Description: "If set, updates the billing counts for the current month before returning. This is an expensive operation with potential performance impact and should be used sparingly.",
					Query:       true,
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handleUseCaseConsumption,
					Summary:  "Reports consumption billing metrics for the current and previous months.",
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: http.StatusText(http.StatusOK),
							Fields: map[string]*framework.FieldSchema{
								"months": {
									Type:        framework.TypeSlice,
									Description: "List of monthly billing data, including the current and previous months.",
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
	refreshData := data.Get("refresh_data").(bool)

	currentMonth := time.Now().UTC()
	previousMonth := timeutil.StartOfPreviousMonth(currentMonth)

	warnings := make([]string, 0)

	// Check if this is a performance standby and if refreshData is true,
	// and add a warning that refresh will be ignored in this case.
	// We do not need to hold stateLock here since HandleRequest is already holding this lock.
	if refreshData && b.Core.perfStandby {
		warnings = append(warnings, WarningRefreshIgnoredOnStandby)
		refreshData = false
	}

	// Refresh data only if explicitly requested and for current month
	currentMonthData, err := b.buildMonthBillingData(ctx, currentMonth, refreshData)
	if err != nil {
		return nil, fmt.Errorf("error building current month billing data: %w", err)
	}

	previousMonthData, err := b.buildMonthBillingData(ctx, previousMonth, false)
	if err != nil {
		return nil, fmt.Errorf("error building previous month billing data: %w", err)
	}

	resp := map[string]interface{}{
		"months": []interface{}{
			currentMonthData,
			previousMonthData,
		},
	}

	return &logical.Response{
		Data:     resp,
		Warnings: warnings,
	}, nil
}

// buildMonthBillingData constructs billing data for a specific month
func (b *SystemBackend) buildMonthBillingData(ctx context.Context, month time.Time, refreshData bool) (map[string]interface{}, error) {
	currentMonth := timeutil.StartOfMonth(time.Now().UTC())
	// Check if the billing metrics need to be refreshed. We're running
	// under the core stateLock during request handling,so call the no-lock helper to
	// avoid recursive locking.
	if refreshData {
		if err := b.Core.updateBillingMetricsLocked(ctx, currentMonth); err != nil {
			return nil, fmt.Errorf("error refreshing billing metrics: %w", err)
		}
	}

	// Retrieve all billing metrics
	combinedRoleCounts, combinedManagedKeyCounts, err := b.Core.getRoleAndManagedKeyCounts(ctx, month)
	if err != nil {
		return nil, err
	}

	combinedKvCounts, err := b.Core.getKvCounts(ctx, month)
	if err != nil {
		return nil, err
	}

	transitCounts, transformCounts, err := b.Core.getDataProtectionCounts(ctx, month)
	if err != nil {
		return nil, err
	}

	kmipEnabled, err := b.Core.getKmipStatus(ctx, month)
	if err != nil {
		return nil, err
	}

	thirdPartyPluginCounts, err := b.Core.getThirdPartyPluginCounts(ctx, month)
	if err != nil {
		return nil, err
	}

	// Build the usage metrics
	usageMetrics := []map[string]interface{}{}

	kvDetails := []map[string]interface{}{}
	if combinedKvCounts > 0 {
		kvDetails = append(kvDetails, map[string]interface{}{"type": "kv", "count": combinedKvCounts})
	}
	usageMetrics = append(usageMetrics, map[string]interface{}{
		"metric_name": "static_secrets",
		"metric_data": map[string]interface{}{
			"total":          combinedKvCounts,
			"metric_details": kvDetails,
		},
	})

	usageMetrics = append(usageMetrics, buildDynamicRolesMetric(combinedRoleCounts))

	usageMetrics = append(usageMetrics, buildAutoRotatedRolesMetric(combinedRoleCounts))

	usageMetrics = append(usageMetrics, map[string]interface{}{
		"metric_name": "kmip",
		"metric_data": map[string]interface{}{
			"used_in_month": kmipEnabled,
		},
	})

	usageMetrics = append(usageMetrics, map[string]interface{}{
		"metric_name": "external_plugins",
		"metric_data": map[string]interface{}{
			"total": thirdPartyPluginCounts,
		},
	})

	dataProtectionDetails := []map[string]interface{}{}
	if transitCounts > 0 {
		dataProtectionDetails = append(dataProtectionDetails, map[string]interface{}{"type": "transit", "count": transitCounts})
	}
	if transformCounts > 0 {
		dataProtectionDetails = append(dataProtectionDetails, map[string]interface{}{"type": "transform", "count": transformCounts})
	}

	usageMetrics = append(usageMetrics, map[string]interface{}{
		"metric_name": "data_protection_calls",
		"metric_data": map[string]interface{}{
			"total":          transitCounts + transformCounts,
			"metric_details": dataProtectionDetails,
		},
	})

	pkiMetric, err := b.buildPkiBillingMetric(ctx, month)
	if err != nil {
		return nil, err
	}
	usageMetrics = append(usageMetrics, pkiMetric)

	managedKeysDetails := []map[string]interface{}{}
	if combinedManagedKeyCounts.TotpKeys > 0 {
		managedKeysDetails = append(managedKeysDetails, map[string]interface{}{"type": "totp", "count": combinedManagedKeyCounts.TotpKeys})
	}
	if combinedManagedKeyCounts.KmseKeys > 0 {
		managedKeysDetails = append(managedKeysDetails, map[string]interface{}{"type": "kmse", "count": combinedManagedKeyCounts.KmseKeys})
	}
	usageMetrics = append(usageMetrics, map[string]interface{}{
		"metric_name": "managed_keys",
		"metric_data": map[string]interface{}{
			"total":          combinedManagedKeyCounts.TotpKeys + combinedManagedKeyCounts.KmseKeys,
			"metric_details": managedKeysDetails,
		},
	})

	dataUpdatedAt := b.Core.computeUpdatedAt(ctx, month, currentMonth)

	monthStr := month.Format("2006-01")

	return map[string]interface{}{
		"month":         monthStr,
		"updated_at":    dataUpdatedAt.Format(time.RFC3339),
		"usage_metrics": usageMetrics,
	}, nil
}

// computeUpdatedAt determines the appropriate updated_at timestamp for billing data
func (c *Core) computeUpdatedAt(ctx context.Context, month, currentMonth time.Time) time.Time {
	var dataUpdatedAt time.Time
	isCurrentMonth := timeutil.StartOfMonth(month).Equal(currentMonth)
	if isCurrentMonth {
		// Use the last time metrics were updated. If it is zero, it means the data has not
		// been updated yet for the current month.
		lastUpdate, err := c.GetMetricsLastUpdateTime(ctx, currentMonth)
		if err != nil {
			// Avoid logging raw error contents which may include sensitive information.
			c.logger.Error("error retrieving last metrics update time")
			return time.Time{}
		}
		dataUpdatedAt = lastUpdate
	} else {
		// Check presence of a stored metrics timestamp for the previous month.
		// If present, return the canonical end-of-month for the requested
		// `month`. The stored timestamp acts strictly as a
		// presence indicator.
		previousMonthStart := timeutil.StartOfPreviousMonth(currentMonth)
		previousMonthTimestamp, err := c.GetMetricsLastUpdateTime(ctx, previousMonthStart)

		// The previous month has not been updated yet.
		if err != nil || previousMonthTimestamp.IsZero() {
			return time.Time{}
		}

		// Use requested month's canonical end-of-month.
		dataUpdatedAt = timeutil.EndOfMonth(month.UTC())
	}

	return dataUpdatedAt
}

// buildDynamicRolesMetric creates the dynamic_roles metric from role counts.
func buildDynamicRolesMetric(counts *RoleCounts) map[string]interface{} {
	total := 0
	if counts != nil {
		total = counts.AWSDynamicRoles +
			counts.AzureDynamicRoles +
			counts.DatabaseDynamicRoles +
			counts.GCPRolesets +
			counts.LDAPDynamicRoles +
			counts.OpenLDAPDynamicRoles +
			counts.AlicloudDynamicRoles +
			counts.RabbitMQDynamicRoles +
			counts.ConsulDynamicRoles +
			counts.NomadDynamicRoles +
			counts.KubernetesDynamicRoles +
			counts.MongoDBAtlasDynamicRoles +
			counts.TerraformCloudDynamicRoles
	}

	details := []map[string]interface{}{}
	if counts != nil {
		if counts.AWSDynamicRoles > 0 {
			details = append(details, map[string]interface{}{"type": "aws_dynamic", "count": counts.AWSDynamicRoles})
		}
		if counts.AzureDynamicRoles > 0 {
			details = append(details, map[string]interface{}{"type": "azure_dynamic", "count": counts.AzureDynamicRoles})
		}
		if counts.DatabaseDynamicRoles > 0 {
			details = append(details, map[string]interface{}{"type": "database_dynamic", "count": counts.DatabaseDynamicRoles})
		}
		if counts.GCPRolesets > 0 {
			details = append(details, map[string]interface{}{"type": "gcp_dynamic", "count": counts.GCPRolesets})
		}
		if counts.LDAPDynamicRoles > 0 {
			details = append(details, map[string]interface{}{"type": "ldap_dynamic", "count": counts.LDAPDynamicRoles})
		}
		if counts.OpenLDAPDynamicRoles > 0 {
			details = append(details, map[string]interface{}{"type": "openldap_dynamic", "count": counts.OpenLDAPDynamicRoles})
		}
		if counts.AlicloudDynamicRoles > 0 {
			details = append(details, map[string]interface{}{"type": "alicloud_dynamic", "count": counts.AlicloudDynamicRoles})
		}
		if counts.RabbitMQDynamicRoles > 0 {
			details = append(details, map[string]interface{}{"type": "rabbitmq_dynamic", "count": counts.RabbitMQDynamicRoles})
		}
		if counts.ConsulDynamicRoles > 0 {
			details = append(details, map[string]interface{}{"type": "consul_dynamic", "count": counts.ConsulDynamicRoles})
		}
		if counts.NomadDynamicRoles > 0 {
			details = append(details, map[string]interface{}{"type": "nomad_dynamic", "count": counts.NomadDynamicRoles})
		}
		if counts.KubernetesDynamicRoles > 0 {
			details = append(details, map[string]interface{}{"type": "kubernetes_dynamic", "count": counts.KubernetesDynamicRoles})
		}
		if counts.MongoDBAtlasDynamicRoles > 0 {
			details = append(details, map[string]interface{}{"type": "mongodbatlas_dynamic", "count": counts.MongoDBAtlasDynamicRoles})
		}
		if counts.TerraformCloudDynamicRoles > 0 {
			details = append(details, map[string]interface{}{"type": "terraform_dynamic", "count": counts.TerraformCloudDynamicRoles})
		}
	}

	return map[string]interface{}{
		"metric_name": "dynamic_roles",
		"metric_data": map[string]interface{}{
			"total":          total,
			"metric_details": details,
		},
	}
}

// buildAutoRotatedRolesMetric creates the auto_rotated_roles metric from role counts.
func buildAutoRotatedRolesMetric(counts *RoleCounts) map[string]interface{} {
	total := 0
	if counts != nil {
		total = counts.AWSStaticRoles +
			counts.AzureStaticRoles +
			counts.DatabaseStaticRoles +
			counts.GCPStaticAccounts +
			counts.GCPImpersonatedAccounts +
			counts.LDAPStaticRoles +
			counts.OpenLDAPStaticRoles
	}

	details := []map[string]interface{}{}
	if counts != nil {
		if counts.AWSStaticRoles > 0 {
			details = append(details, map[string]interface{}{"type": "aws_static", "count": counts.AWSStaticRoles})
		}
		if counts.AzureStaticRoles > 0 {
			details = append(details, map[string]interface{}{"type": "azure_static", "count": counts.AzureStaticRoles})
		}
		if counts.DatabaseStaticRoles > 0 {
			details = append(details, map[string]interface{}{"type": "database_static", "count": counts.DatabaseStaticRoles})
		}
		if counts.GCPStaticAccounts > 0 {
			details = append(details, map[string]interface{}{"type": "gcp_static", "count": counts.GCPStaticAccounts})
		}
		if counts.GCPImpersonatedAccounts > 0 {
			details = append(details, map[string]interface{}{"type": "gcp_impersonated", "count": counts.GCPImpersonatedAccounts})
		}
		if counts.LDAPStaticRoles > 0 {
			details = append(details, map[string]interface{}{"type": "ldap_static", "count": counts.LDAPStaticRoles})
		}
		if counts.OpenLDAPStaticRoles > 0 {
			details = append(details, map[string]interface{}{"type": "openldap_static", "count": counts.OpenLDAPStaticRoles})
		}
	}

	return map[string]interface{}{
		"metric_name": "auto_rotated_roles",
		"metric_data": map[string]interface{}{
			"total":          total,
			"metric_details": details,
		},
	}
}

// buildPkiBillingMetric creates the billing metric for PKI duration-adjusted certificate counts.
func (b *SystemBackend) buildPkiBillingMetric(ctx context.Context, month time.Time) (map[string]interface{}, error) {
	count, err := b.Core.GetStoredPkiDurationAdjustedCount(ctx, month)
	if err != nil {
		return nil, fmt.Errorf("error retrieving PKI duration-adjusted count for month: %w", err)
	}

	return map[string]interface{}{
		"metric_name": "pki_units",
		"metric_data": map[string]interface{}{
			"total": count,
		},
	}, nil
}

// getRoleCounts retrieves and combines role and managed key counts from replicated and local storage
func (c *Core) getRoleAndManagedKeyCounts(ctx context.Context, month time.Time) (*RoleCounts, *ManagedKeyCounts, error) {
	var replicatedRoleCounts *RoleCounts
	replicatedTotpHWMValue := 0
	replicatedKmseHWMValue := 0
	var err error

	if c.isPrimary() {
		replicatedRoleCounts, err = c.GetStoredHWMRoleCounts(ctx, billing.ReplicatedPrefix, month)
		if err != nil {
			return nil, nil, fmt.Errorf("error retrieving replicated max role counts: %w", err)
		}
		replicatedTotpHWMValue, err = c.GetStoredHWMTotpCounts(ctx, billing.ReplicatedPrefix, month)
		if err != nil {
			return nil, nil, fmt.Errorf("error retrieving replicated max managed key count: %w", err)
		}
		replicatedKmseHWMValue, err = c.GetStoredHWMKmseCounts(ctx, billing.ReplicatedPrefix, month)
		if err != nil {
			return nil, nil, fmt.Errorf("error retrieving replicated max kmse key count: %w", err)
		}
	}

	localRoleCounts, err := c.GetStoredHWMRoleCounts(ctx, billing.LocalPrefix, month)
	if err != nil {
		return nil, nil, fmt.Errorf("error retrieving local max role counts: %w", err)
	}
	localTotpHWMValue, err := c.GetStoredHWMTotpCounts(ctx, billing.LocalPrefix, month)
	if err != nil {
		return nil, nil, fmt.Errorf("error retrieving local max totp key count: %w", err)
	}
	localKmseHWMValue, err := c.GetStoredHWMKmseCounts(ctx, billing.LocalPrefix, month)
	if err != nil {
		return nil, nil, fmt.Errorf("error retrieving local max kmse key count: %w", err)
	}

	combinedManagedKeyCounts := &ManagedKeyCounts{
		TotpKeys: localTotpHWMValue + replicatedTotpHWMValue,
		KmseKeys: localKmseHWMValue + replicatedKmseHWMValue,
	}

	return combineRoleCounts(replicatedRoleCounts, localRoleCounts), combinedManagedKeyCounts, nil
}

// getKvCounts retrieves and combines KV secret counts from replicated and local storage
func (c *Core) getKvCounts(ctx context.Context, month time.Time) (int, error) {
	var replicatedKvCounts int
	var err error

	if c.isPrimary() {
		replicatedKvCounts, err = c.GetStoredHWMKvCounts(ctx, billing.ReplicatedPrefix, month)
		if err != nil {
			return 0, fmt.Errorf("error retrieving replicated max kv counts: %w", err)
		}
	}

	localKvCounts, err := c.GetStoredHWMKvCounts(ctx, billing.LocalPrefix, month)
	if err != nil {
		return 0, fmt.Errorf("error retrieving local max kv counts: %w", err)
	}

	return replicatedKvCounts + localKvCounts, nil
}

// getDataProtectionCounts retrieves Transit and Transform call counts
// Data protection call counts are stored at local path only
// Each cluster tracks its own total requests to avoid double counting
func (c *Core) getDataProtectionCounts(ctx context.Context, month time.Time) (uint64, uint64, error) {
	transitCounts, err := c.GetStoredTransitCallCounts(ctx, month)
	if err != nil {
		return 0, 0, fmt.Errorf("error retrieving local transit call counts: %w", err)
	}
	transformCounts, err := c.GetStoredTransformCallCounts(ctx, month)
	if err != nil {
		return 0, 0, fmt.Errorf("error retrieving local transform call counts: %w", err)
	}

	return transitCounts, transformCounts, nil
}

// getKmipStatus retrieves KMIP enabled status (always stored at local path)
func (c *Core) getKmipStatus(ctx context.Context, month time.Time) (bool, error) {
	kmipEnabled, err := c.GetStoredKmipEnabled(ctx, month)
	if err != nil {
		return false, fmt.Errorf("error retrieving KMIP enabled status: %w", err)
	}

	return kmipEnabled, nil
}

// getThirdPartyPluginCounts retrieves third-party plugin counts (always stored at local path)
func (c *Core) getThirdPartyPluginCounts(ctx context.Context, month time.Time) (int, error) {
	thirdPartyPluginCounts, err := c.GetStoredThirdPartyPluginCounts(ctx, month)
	if err != nil {
		return 0, fmt.Errorf("error retrieving third-party plugin counts: %w", err)
	}

	return thirdPartyPluginCounts, nil
}
