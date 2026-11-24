// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package vault

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/hashicorp/vault/internalshared/timeutil"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault/billing"
)

const (
	WarningRefreshIgnoredOnStandby = "refresh_data parameter is supported only on the active node. " +
		"Since this parameter was set on a performance standby, the billing data was not refreshed " +
		"and retrieved from storage without update."

	WarningRefreshAlreadyInProgress = "a billing data refresh is already in progress. " +
		"Since refresh_data was set, the billing data was not refreshed by this request " +
		"and retrieved from storage without update."

	WarningStartEndMonthOutOfRetentionRange = "the specified start_month and/or end_month fall outside the range of the current billing data retention period." +
		"Months that are not covered in the retention period will show a zero updated_at timestamp and no metrics."
)

func (b *SystemBackend) useCaseConsumptionBillingPaths() []*framework.Path {
	return []*framework.Path{
		b.billingOverviewPath(),
		b.billingConfigPath(),
	}
}

func (b *SystemBackend) billingOverviewPath() *framework.Path {
	return &framework.Path{
		Pattern: "billing/overview$",
		Fields: map[string]*framework.FieldSchema{
			"refresh_data": {
				Type:        framework.TypeBool,
				Description: "If set, updates the billing counts for the current month before returning. This is an expensive operation with potential performance impact and should be used sparingly.",
				Query:       true,
			},
			"start_month": {
				Type:        framework.TypeString,
				Description: "Start month in YYYY-MM format (inclusive). If not specified, defaults to the oldest available month within BillingRetentionMonths.",
				Query:       true,
			},
			"end_month": {
				Type:        framework.TypeString,
				Description: "End month in YYYY-MM format (inclusive). If not specified, defaults to the current month.",
				Query:       true,
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.handleBillingOverview,
				Summary:  "Reports consumption billing metrics on a monthly granularity.",
				Responses: map[int][]framework.Response{
					http.StatusOK: {{
						Description: http.StatusText(http.StatusOK),
						Fields: map[string]*framework.FieldSchema{
							"months": {
								Type:        framework.TypeSlice,
								Description: "List of monthly billing data.",
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
	}
}

func (b *SystemBackend) billingConfigPath() *framework.Path {
	return &framework.Path{
		Pattern: "billing/config$",
		Fields: map[string]*framework.FieldSchema{
			"retention_months": {
				Type:        framework.TypeInt,
				Description: fmt.Sprintf("Number of months to retain billing data. Must be between %d and %d months. Defaults to %d months.", billing.MinBillingRetentionMonths, billing.MaxBillingRetentionMonths, billing.DefaultBillingRetentionMonths),
			},
			"attribution_retention_months": {
				Type:        framework.TypeInt,
				Description: fmt.Sprintf("Number of months to retain mount attribution data for billing metrics. Must be between %d and %d months. Defaults to %d months. Setting to 0 disables attribution storage and wipes all existing attribution data.", billing.MinAttributionRetentionMonths, billing.MaxAttributionRetentionMonths, billing.DefaultAttributionRetentionMonths),
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.handleBillingConfigRead,
				Summary:  "Read the billing data retention configuration.",
				Responses: map[int][]framework.Response{
					http.StatusOK: {{
						Description: http.StatusText(http.StatusOK),
						Fields: map[string]*framework.FieldSchema{
							"retention_months": {
								Type:        framework.TypeInt,
								Description: "Number of months of billing data to retain.",
							},
							"attribution_retention_months": {
								Type:        framework.TypeInt,
								Description: "Number of months of attribution data to retain. 0 means attribution storage is disabled.",
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
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.handleBillingConfigWrite,
				Summary:  "Configure the billing data retention period.",
				Responses: map[int][]framework.Response{
					http.StatusOK: {{
						Description: http.StatusText(http.StatusOK),
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
	}
}

func (b *SystemBackend) handleBillingConfigRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	retentionMonths, err := b.Core.GetBillingRetentionMonths(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get billing retention configuration: %w", err)
	}

	attributionRetentionMonths, err := b.Core.GetAttributionRetentionMonths(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get attribution retention configuration: %w", err)
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"retention_months":             retentionMonths,
			"attribution_retention_months": attributionRetentionMonths,
		},
	}, nil
}

func (b *SystemBackend) handleBillingConfigWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	resp := &logical.Response{}

	// Handle retention_months if provided
	if rawRetention, ok := data.GetOk("retention_months"); ok {
		retentionMonths := rawRetention.(int)
		if retentionMonths < billing.MinBillingRetentionMonths || retentionMonths > billing.MaxBillingRetentionMonths {
			return logical.ErrorResponse(fmt.Sprintf("retention_months must be between %d and %d months", billing.MinBillingRetentionMonths, billing.MaxBillingRetentionMonths)), logical.ErrInvalidRequest
		}

		// Get current retention to check if it's being increased
		currentRetention, err := b.Core.GetBillingRetentionMonths(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get current billing retention configuration: %w", err)
		}

		// Store the configuration
		if err := b.Core.UpdateBillingRetentionMonths(ctx, retentionMonths); err != nil {
			return nil, fmt.Errorf("failed to set billing retention configuration: %w", err)
		}

		// Add warning if retention period is being increased
		if retentionMonths > currentRetention {
			resp.Warnings = append(resp.Warnings, fmt.Sprintf(
				"Retention period increased from %d to %d months. Historical data will only be available for months within the previous retention period. Older months outside the previous retention range will not have data.",
				currentRetention, retentionMonths))
		}
	}

	// Handle attribution_retention_months if provided
	if rawAttrRetention, ok := data.GetOk("attribution_retention_months"); ok {
		attributionRetentionMonths := rawAttrRetention.(int)
		if attributionRetentionMonths < billing.MinAttributionRetentionMonths || attributionRetentionMonths > billing.MaxAttributionRetentionMonths {
			return logical.ErrorResponse(fmt.Sprintf("attribution_retention_months must be between %d and %d months", billing.MinAttributionRetentionMonths, billing.MaxAttributionRetentionMonths)), logical.ErrInvalidRequest
		}

		currentAttributionRetention, err := b.Core.GetAttributionRetentionMonths(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get current attribution retention configuration: %w", err)
		}

		// Store the configuration
		if err := b.Core.UpdateAttributionRetentionMonths(ctx, attributionRetentionMonths); err != nil {
			return nil, fmt.Errorf("failed to set attribution retention configuration: %w", err)
		}

		// Add warning if retention period is being increased
		if attributionRetentionMonths > currentAttributionRetention {
			resp.Warnings = append(resp.Warnings, fmt.Sprintf(
				"Attribution retention period increased from %d to %d months. Historical data will only be available for months within the previous retention period. Older months outside the previous retention range will not have data.",
				currentAttributionRetention, attributionRetentionMonths))
		}

		// When attribution is disabled (0 months), wipe all existing attribution data
		if attributionRetentionMonths == billing.MinAttributionRetentionMonths {
			if err := b.Core.DeleteExpiredAttributionData(ctx, time.Now().UTC()); err != nil {
				b.Core.logger.Error("failed to wipe attribution data after disabling attribution", "error", err)
				resp.Warnings = append(resp.Warnings, "Attribution storage has been disabled but existing data could not be fully wiped. It will be cleaned up on the next billing cycle.")
			} else {
				resp.Warnings = append(resp.Warnings, "Attribution storage has been disabled and all existing attribution data has been wiped.")
			}
		}
	}

	return resp, nil
}

func (b *SystemBackend) handleBillingOverview(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	refreshData := data.Get("refresh_data").(bool)

	currentMonth := time.Now().UTC()

	warnings := make([]string, 0)

	// Check if this is a performance standby and if refreshData is true,
	// and add a warning that refresh will be ignored in this case.
	// We do not need to hold stateLock here since HandleRequest is already holding this lock.
	if refreshData && b.Core.perfStandby {
		warnings = append(warnings, WarningRefreshIgnoredOnStandby)
		refreshData = false
	}

	// Get the configured retention period
	retentionMonths, err := b.Core.GetBillingRetentionMonths(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get billing retention configuration: %w", err)
	}

	startMonth, endMonth, isOutOfRetention, err := parseStartEndMonths(data, currentMonth, retentionMonths)
	if err != nil {
		return nil, err
	}

	if isOutOfRetention {
		warnings = append(warnings, WarningStartEndMonthOutOfRetentionRange)
	}

	// Build list of months to retrieve (from end to start, newest first)
	monthsToRetrieve := []time.Time{}
	for month := endMonth; !month.Before(startMonth); month = month.AddDate(0, -1, 0) {
		monthsToRetrieve = append(monthsToRetrieve, month)
	}

	// Build billing data for requested months
	months := make([]interface{}, 0, len(monthsToRetrieve))

	for _, month := range monthsToRetrieve {
		// Only refresh current month if refresh_data is true
		shouldRefresh := refreshData && month.Equal(timeutil.StartOfMonth(currentMonth))

		monthData, err := b.buildMonthBillingData(ctx, month, shouldRefresh, &warnings)
		if err != nil {
			return nil, fmt.Errorf("error building billing data for month %s: %w", month.Format("2006-01"), err)
		}

		months = append(months, monthData)
	}

	resp := map[string]interface{}{
		"months": months,
	}

	return &logical.Response{
		Data:     resp,
		Warnings: warnings,
	}, nil
}

// parseStartEndMonths parses the start and end month parameters from the request and validates if they are valid.
// If they are outside of the retention range, it returns a warning. If no parameter is specified,
// the start and end defaults to the start of the retention range and the current month, respectively.
func parseStartEndMonths(data *framework.FieldData, currentMonth time.Time, retentionMonths int) (time.Time, time.Time, bool, error) {
	defaultStartMonth := timeutil.StartOfMonth(currentMonth).AddDate(0, -retentionMonths+1, 0)
	defaultEndMonth := timeutil.StartOfMonth(currentMonth)

	parseMonth := func(key string, defaultMonth time.Time) (time.Time, error) {
		if monthStr := data.Get(key).(string); monthStr != "" {
			return time.Parse("2006-01", monthStr)
		}
		return defaultMonth, nil
	}

	var startMonth, endMonth time.Time
	var isOutOfRetention bool
	var err error

	startMonth, err = parseMonth("start_month", defaultStartMonth)
	if err != nil {
		return time.Time{}, time.Time{}, false, fmt.Errorf("invalid start_month format: %w", err)
	}

	endMonth, err = parseMonth("end_month", defaultEndMonth)
	if err != nil {
		return time.Time{}, time.Time{}, false, fmt.Errorf("invalid end_month format: %w", err)
	}

	if startMonth.After(endMonth) {
		return time.Time{}, time.Time{}, false, fmt.Errorf("start_month is later than end_month")
	}

	// We don't need to check for startMonth after the current month because either an even later endMonth is
	// specified which would be caught by the second condition, or no end was set and it defaulted to the current month,
	// which would have been caught in the check above. Vice versa for endMonth before the default start month.
	if startMonth.Before(defaultStartMonth) || endMonth.After(defaultEndMonth) {
		isOutOfRetention = true
	}

	return startMonth, endMonth, isOutOfRetention, nil
}

// buildMonthBillingData constructs billing data for a specific month
func (b *SystemBackend) buildMonthBillingData(ctx context.Context, month time.Time, refreshData bool, warnings *[]string) (map[string]interface{}, error) {
	currentMonth := timeutil.StartOfMonth(time.Now().UTC())
	// Check if the billing metrics need to be refreshed. We're running
	// under the core stateLock during request handling,so call the no-lock helper to
	// avoid recursive locking.
	if refreshData {
		b.Core.consumptionBillingLock.RLock()
		cb := b.Core.consumptionBilling
		b.Core.consumptionBillingLock.RUnlock()

		// Avoid starting a redundant scan from the API. If one is already in progress, skip the
		// refresh, warn, and return the stored (non-updated) data. If we can grab the lock, hold
		// it for the refresh; it is released when this function returns, on success or error.
		if cb != nil && !cb.RefreshInProgress.TryLock() {
			*warnings = append(*warnings, WarningRefreshAlreadyInProgress)
		} else {
			if cb != nil {
				defer cb.RefreshInProgress.Unlock()
			}
			if err := b.Core.updateBillingMetricsLocked(ctx, currentMonth); err != nil {
				return nil, fmt.Errorf("error refreshing billing metrics: %w", err)
			}
		}
	}

	m, err := b.Core.getMonthlyConsumption(ctx, month)
	if err != nil {
		return nil, err
	}

	usageMetrics := formatUsageMetrics(m)

	// Round all float64 values in usageMetrics to 4 decimal places.
	// Rounding time for usage metrics is insignificant, so we can keep it centralized here.
	// This prevents us from having to do it in each individual metric.
	roundUsageMetrics(usageMetrics)

	dataUpdatedAt := b.Core.computeUpdatedAt(ctx, month, currentMonth)

	monthStr := month.Format("2006-01")

	return map[string]interface{}{
		"month":         monthStr,
		"updated_at":    dataUpdatedAt.Format(time.RFC3339),
		"usage_metrics": usageMetrics,
	}, nil
}

// formatUsageMetrics shapes the stored monthly consumption into the usage_metrics
// list of the overview response. It performs no storage reads.
func formatUsageMetrics(m *monthlyConsumption) []map[string]interface{} {
	usageMetrics := []map[string]interface{}{}

	usageMetrics = append(usageMetrics, map[string]interface{}{
		"metric_name": "static_secrets",
		"metric_data": map[string]interface{}{
			"total": m.kv,
			"metric_details": []map[string]interface{}{
				{"type": "kv", "count": m.kv},
			},
		},
	})

	usageMetrics = append(usageMetrics, buildDynamicRolesMetric(m.roles))

	usageMetrics = append(usageMetrics, buildAutoRotatedRolesMetric(m.roles))

	usageMetrics = append(usageMetrics, map[string]interface{}{
		"metric_name": "kmip",
		"metric_data": map[string]interface{}{
			"used_in_month": m.kmipUsed,
		},
	})

	usageMetrics = append(usageMetrics, map[string]interface{}{
		"metric_name": "external_plugins",
		"metric_data": map[string]interface{}{
			"total": m.externalPlugins,
		},
	})

	usageMetrics = append(usageMetrics, map[string]interface{}{
		"metric_name": "data_protection_calls",
		"metric_data": map[string]interface{}{
			"total": m.transitCalls + m.transformCalls + m.gcpKmsCalls,
			"metric_details": []map[string]interface{}{
				{"type": "transit", "count": m.transitCalls},
				{"type": "transform", "count": m.transformCalls},
				{"type": "gcpkms", "count": m.gcpKmsCalls},
			},
		},
	})

	usageMetrics = append(usageMetrics, formatPkiMetric(m.pkiUnits))

	usageMetrics = append(usageMetrics, map[string]interface{}{
		"metric_name": "managed_keys",
		"metric_data": map[string]interface{}{
			"total": m.totpKeys + m.kmseKeys,
			"metric_details": []map[string]interface{}{
				{"type": "totp", "count": m.totpKeys},
				{"type": "kmse", "count": m.kmseKeys},
			},
		},
	})

	usageMetrics = append(usageMetrics, buildSSHMetric(m.sshCertUnits, m.sshOTPUnits))

	usageMetrics = append(usageMetrics, buildIdTokenUnitsMetric(m.oidcTokenUnits, m.spiffeTokenUnits))

	usageMetrics = append(usageMetrics, buildExternalCaMetric(m.externalCAUnits))

	return usageMetrics
}

// roundUsageMetrics rounds all float64 values in the usage metrics to 4 decimal places
func roundUsageMetrics(metrics []map[string]interface{}) {
	for _, metric := range metrics {
		if metricData, ok := metric["metric_data"].(map[string]interface{}); ok {
			// Round the total if it's a float64
			if total, ok := metricData["total"].(float64); ok {
				metricData["total"] = roundToFour(total)
			}

			// Round values in metric_details if present
			if details, ok := metricData["metric_details"].([]map[string]interface{}); ok {
				for _, detail := range details {
					if count, ok := detail["count"].(float64); ok {
						detail["count"] = roundToFour(count)
					}
				}
			}
		}
	}
}

// roundToFour takes a float64 and rounds it to 4 decimal places.
func roundToFour(val float64) float64 {
	ratio := math.Pow(10, 4)
	return math.Round(val*ratio) / ratio
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
		// Check presence of a stored metrics timestamp for the requested month.
		// If present, return the canonical end-of-month for the requested
		// `month`. The stored timestamp acts strictly as a
		// presence indicator.
		requestedMonthStart := timeutil.StartOfMonth(month)
		requestedMonthTimestamp, err := c.GetMetricsLastUpdateTime(ctx, requestedMonthStart)

		// The requested month has not been updated yet.
		if err != nil || requestedMonthTimestamp.IsZero() {
			return time.Time{}
		}

		// Use requested month's canonical end-of-month.
		dataUpdatedAt = timeutil.EndOfMonth(month).UTC()
	}

	return dataUpdatedAt
}

// dynamicRolesTotal returns the total number of dynamic roles, rolesets and library sets.
func dynamicRolesTotal(counts *RoleCounts) int {
	if counts == nil {
		return 0
	}
	return counts.AWSDynamicRoles + counts.AzureDynamicRoles + counts.DatabaseDynamicRoles +
		counts.GCPRolesets + counts.LDAPDynamicRoles + counts.LDAPLibrarySets +
		counts.OpenLDAPDynamicRoles + counts.OpenLDAPLibrarySets + counts.AlicloudDynamicRoles +
		counts.RabbitMQDynamicRoles + counts.ConsulDynamicRoles + counts.NomadDynamicRoles +
		counts.KubernetesDynamicRoles + counts.MongoDBAtlasDynamicRoles + counts.TerraformCloudDynamicRoles
}

// buildDynamicRolesMetric creates the dynamic_roles metric from role counts.
func buildDynamicRolesMetric(counts *RoleCounts) map[string]interface{} {
	total := dynamicRolesTotal(counts)
	awsCount := 0
	azureCount := 0
	databaseCount := 0
	gcpCount := 0
	ldapCount := 0
	ldapLibraryCount := 0
	openldapCount := 0
	openldapLibraryCount := 0
	alicloudCount := 0
	rabbitmqCount := 0
	consulCount := 0
	nomadCount := 0
	kubernetesCount := 0
	mongodbatlasCount := 0
	terraformCount := 0

	if counts != nil {
		awsCount = counts.AWSDynamicRoles
		azureCount = counts.AzureDynamicRoles
		databaseCount = counts.DatabaseDynamicRoles
		gcpCount = counts.GCPRolesets
		ldapCount = counts.LDAPDynamicRoles
		ldapLibraryCount = counts.LDAPLibrarySets
		openldapCount = counts.OpenLDAPDynamicRoles
		openldapLibraryCount = counts.OpenLDAPLibrarySets
		alicloudCount = counts.AlicloudDynamicRoles
		rabbitmqCount = counts.RabbitMQDynamicRoles
		consulCount = counts.ConsulDynamicRoles
		nomadCount = counts.NomadDynamicRoles
		kubernetesCount = counts.KubernetesDynamicRoles
		mongodbatlasCount = counts.MongoDBAtlasDynamicRoles
		terraformCount = counts.TerraformCloudDynamicRoles
	}

	details := []map[string]interface{}{
		{"type": "aws_dynamic", "count": awsCount},
		{"type": "azure_dynamic", "count": azureCount},
		{"type": "database_dynamic", "count": databaseCount},
		{"type": "gcp_dynamic", "count": gcpCount},
		{"type": "ldap_dynamic", "count": ldapCount},
		{"type": "ldap_library", "count": ldapLibraryCount},
		{"type": "openldap_dynamic", "count": openldapCount},
		{"type": "openldap_library", "count": openldapLibraryCount},
		{"type": "alicloud_dynamic", "count": alicloudCount},
		{"type": "rabbitmq_dynamic", "count": rabbitmqCount},
		{"type": "consul_dynamic", "count": consulCount},
		{"type": "nomad_dynamic", "count": nomadCount},
		{"type": "kubernetes_dynamic", "count": kubernetesCount},
		{"type": "mongodbatlas_dynamic", "count": mongodbatlasCount},
		{"type": "terraform_dynamic", "count": terraformCount},
	}

	return map[string]interface{}{
		"metric_name": "dynamic_roles",
		"metric_data": map[string]interface{}{
			"total":          total,
			"metric_details": details,
		},
	}
}

// autoRotatedRolesTotal returns the total number of auto-rotated (static) roles.
func autoRotatedRolesTotal(counts *RoleCounts) int {
	if counts == nil {
		return 0
	}
	return counts.AWSStaticRoles + counts.AzureStaticRoles + counts.DatabaseStaticRoles +
		counts.GCPStaticAccounts + counts.GCPImpersonatedAccounts + counts.LDAPStaticRoles +
		counts.OpenLDAPStaticRoles + counts.OSLocalAccountRoles
}

// buildAutoRotatedRolesMetric creates the auto_rotated_roles metric from role counts.
func buildAutoRotatedRolesMetric(counts *RoleCounts) map[string]interface{} {
	total := autoRotatedRolesTotal(counts)
	awsCount := 0
	azureCount := 0
	databaseCount := 0
	gcpStaticCount := 0
	gcpImpersonatedCount := 0
	ldapCount := 0
	openldapCount := 0
	osLocalAccountCount := 0

	if counts != nil {
		awsCount = counts.AWSStaticRoles
		azureCount = counts.AzureStaticRoles
		databaseCount = counts.DatabaseStaticRoles
		gcpStaticCount = counts.GCPStaticAccounts
		gcpImpersonatedCount = counts.GCPImpersonatedAccounts
		ldapCount = counts.LDAPStaticRoles
		openldapCount = counts.OpenLDAPStaticRoles
		osLocalAccountCount = counts.OSLocalAccountRoles
	}

	details := []map[string]interface{}{
		{"type": "aws_static", "count": awsCount},
		{"type": "azure_static", "count": azureCount},
		{"type": "database_static", "count": databaseCount},
		{"type": "gcp_static", "count": gcpStaticCount},
		{"type": "gcp_impersonated", "count": gcpImpersonatedCount},
		{"type": "ldap_static", "count": ldapCount},
		{"type": "openldap_static", "count": openldapCount},
		{"type": "os_local_account_static", "count": osLocalAccountCount},
	}

	return map[string]interface{}{
		"metric_name": "auto_rotated_roles",
		"metric_data": map[string]interface{}{
			"total":          total,
			"metric_details": details,
		},
	}
}

// formatPkiMetric creates the pki_units metric from the stored PKI duration-adjusted certificate count.
func formatPkiMetric(count float64) map[string]interface{} {
	return map[string]interface{}{
		"metric_name": "pki_units",
		"metric_data": map[string]interface{}{
			"total": count,
		},
	}
}

// buildPkiBillingMetric reads the stored PKI count for the month and formats it as the pki_units metric.
func (b *SystemBackend) buildPkiBillingMetric(ctx context.Context, month time.Time) (map[string]interface{}, error) {
	count, err := b.Core.GetStoredPkiDurationAdjustedCount(ctx, month)
	if err != nil {
		return nil, fmt.Errorf("error retrieving PKI duration-adjusted count for month: %w", err)
	}

	return formatPkiMetric(count), nil
}

// buildIdTokenUnitsMetric creates the id_token_units metric from the stored OIDC and SPIFFE JWT units.
func buildIdTokenUnitsMetric(oidcTokenCount, spiffeJwtUnits float64) map[string]interface{} {
	return map[string]interface{}{
		"metric_name": "id_token_units",
		"metric_data": map[string]interface{}{
			"total": oidcTokenCount + spiffeJwtUnits,
			"metric_details": []map[string]interface{}{
				{"type": "oidc", "count": oidcTokenCount},
				{"type": "spiffe", "count": spiffeJwtUnits},
			},
		},
	}
}

// buildExternalCaMetric creates the external_ca_pki_units metric from the stored external CA certificate units.
func buildExternalCaMetric(count float64) map[string]interface{} {
	return map[string]interface{}{
		"metric_name": "external_ca_pki_units",
		"metric_data": map[string]interface{}{
			"total": count,
		},
	}
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

// getDataProtectionCounts retrieves Transit, Transform, and GCP KMS call counts
// Data protection call counts are stored at local path only
// Each cluster tracks its own total requests to avoid double counting
func (c *Core) getDataProtectionCounts(ctx context.Context, month time.Time) (uint64, uint64, uint64, error) {
	transitCounts, err := c.GetStoredTransitCallCounts(ctx, month)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("error retrieving local transit call counts: %w", err)
	}
	transformCounts, err := c.GetStoredTransformCallCounts(ctx, month)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("error retrieving local transform call counts: %w", err)
	}
	gcpKmsCounts, err := c.GetStoredGcpKmsCallCounts(ctx, month)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("error retrieving local GCP KMS call counts: %w", err)
	}

	return transitCounts, transformCounts, gcpKmsCounts, nil
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

// buildSSHMetric creates the ssh_units metric from the stored SSH certificate and OTP units.
func buildSSHMetric(certCounts, otpCounts float64) map[string]interface{} {
	return map[string]interface{}{
		"metric_name": "ssh_units",
		"metric_data": map[string]interface{}{
			"total": certCounts + otpCounts,
			"metric_details": []map[string]interface{}{
				{
					"type":  "otp_units",
					"count": otpCounts,
				},
				{
					"type":  "certificate_units",
					"count": certCounts,
				},
			},
		},
	}
}
