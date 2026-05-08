// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package vault

import (
	"context"
	"fmt"
	"math"
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

	WarningStartEndMonthOutOfRetentionRange = "the specified start_month and/or end_month fall outside the range of the current billing data retention period." +
		"Months that are not covered in the retention period will show a zero updated_at timestamp and no metrics."
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
					Callback: b.handleUseCaseConsumption,
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
		},
	}
}

func (b *SystemBackend) handleUseCaseConsumption(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
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

	startMonth, endMonth, isOutOfRetention, err := parseStartEndMonths(data, currentMonth)
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

		monthData, err := b.buildMonthBillingData(ctx, month, shouldRefresh)
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
// If they are outside of the BillingRetentionMonths range, it returns a warning. If no parameter is specified,
// the start and end defaults to the start of the BillingRetentionMonths range and the current month, respectively.
func parseStartEndMonths(data *framework.FieldData, currentMonth time.Time) (time.Time, time.Time, bool, error) {
	defaultStartMonth := timeutil.StartOfMonth(currentMonth).AddDate(0, -billing.BillingRetentionMonths+1, 0)
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

	transitCounts, transformCounts, gcpKmsCounts, err := b.Core.getDataProtectionCounts(ctx, month)
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

	kvDetails := []map[string]interface{}{
		{"type": "kv", "count": combinedKvCounts},
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

	dataProtectionDetails := []map[string]interface{}{
		{"type": "transit", "count": transitCounts},
		{"type": "transform", "count": transformCounts},
		{"type": "gcpkms", "count": gcpKmsCounts},
	}

	usageMetrics = append(usageMetrics, map[string]interface{}{
		"metric_name": "data_protection_calls",
		"metric_data": map[string]interface{}{
			"total":          transitCounts + transformCounts + gcpKmsCounts,
			"metric_details": dataProtectionDetails,
		},
	})

	pkiMetric, err := b.buildPkiBillingMetric(ctx, month)
	if err != nil {
		return nil, err
	}
	usageMetrics = append(usageMetrics, pkiMetric)

	managedKeysDetails := []map[string]interface{}{
		{"type": "totp", "count": combinedManagedKeyCounts.TotpKeys},
		{"type": "kmse", "count": combinedManagedKeyCounts.KmseKeys},
	}
	usageMetrics = append(usageMetrics, map[string]interface{}{
		"metric_name": "managed_keys",
		"metric_data": map[string]interface{}{
			"total":          combinedManagedKeyCounts.TotpKeys + combinedManagedKeyCounts.KmseKeys,
			"metric_details": managedKeysDetails,
		},
	})

	sshCounts, err := b.buildSSHMetric(ctx, month)
	if err != nil {
		return nil, err
	}
	usageMetrics = append(usageMetrics, sshCounts)

	idTokenUnitsMetric, err := b.buildIdTokenUnitsBillingMetric(ctx, month)
	if err != nil {
		return nil, err
	}
	usageMetrics = append(usageMetrics, idTokenUnitsMetric)

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

// buildDynamicRolesMetric creates the dynamic_roles metric from role counts.
func buildDynamicRolesMetric(counts *RoleCounts) map[string]interface{} {
	total := 0
	awsCount := 0
	azureCount := 0
	databaseCount := 0
	gcpCount := 0
	ldapCount := 0
	openldapCount := 0
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
		openldapCount = counts.OpenLDAPDynamicRoles
		alicloudCount = counts.AlicloudDynamicRoles
		rabbitmqCount = counts.RabbitMQDynamicRoles
		consulCount = counts.ConsulDynamicRoles
		nomadCount = counts.NomadDynamicRoles
		kubernetesCount = counts.KubernetesDynamicRoles
		mongodbatlasCount = counts.MongoDBAtlasDynamicRoles
		terraformCount = counts.TerraformCloudDynamicRoles

		total = awsCount + azureCount + databaseCount + gcpCount + ldapCount +
			openldapCount + alicloudCount + rabbitmqCount + consulCount +
			nomadCount + kubernetesCount + mongodbatlasCount + terraformCount
	}

	details := []map[string]interface{}{
		{"type": "aws_dynamic", "count": awsCount},
		{"type": "azure_dynamic", "count": azureCount},
		{"type": "database_dynamic", "count": databaseCount},
		{"type": "gcp_dynamic", "count": gcpCount},
		{"type": "ldap_dynamic", "count": ldapCount},
		{"type": "openldap_dynamic", "count": openldapCount},
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

// buildAutoRotatedRolesMetric creates the auto_rotated_roles metric from role counts.
func buildAutoRotatedRolesMetric(counts *RoleCounts) map[string]interface{} {
	total := 0
	awsCount := 0
	azureCount := 0
	databaseCount := 0
	gcpStaticCount := 0
	gcpImpersonatedCount := 0
	ldapCount := 0
	openldapCount := 0

	if counts != nil {
		awsCount = counts.AWSStaticRoles
		azureCount = counts.AzureStaticRoles
		databaseCount = counts.DatabaseStaticRoles
		gcpStaticCount = counts.GCPStaticAccounts
		gcpImpersonatedCount = counts.GCPImpersonatedAccounts
		ldapCount = counts.LDAPStaticRoles
		openldapCount = counts.OpenLDAPStaticRoles

		total = awsCount + azureCount + databaseCount + gcpStaticCount +
			gcpImpersonatedCount + ldapCount + openldapCount
	}

	details := []map[string]interface{}{
		{"type": "aws_static", "count": awsCount},
		{"type": "azure_static", "count": azureCount},
		{"type": "database_static", "count": databaseCount},
		{"type": "gcp_static", "count": gcpStaticCount},
		{"type": "gcp_impersonated", "count": gcpImpersonatedCount},
		{"type": "ldap_static", "count": ldapCount},
		{"type": "openldap_static", "count": openldapCount},
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

// buildIdTokenUnitsBillingMetric creates the billing metric for id token counts.
func (b *SystemBackend) buildIdTokenUnitsBillingMetric(ctx context.Context, month time.Time) (map[string]interface{}, error) {
	var totalTokens float64

	oidcTokenCount, err := b.Core.GetStoredOidcDurationAdjustedCount(ctx, month)
	if err != nil {
		return nil, fmt.Errorf("error retrieving OIDC duration-adjusted token count for month: %w", err)
	}

	totalTokens += oidcTokenCount

	spiffeJwtUnits, err := b.Core.GetStoredSpiffeJwtTokenUnits(ctx, month)
	if err != nil {
		return nil, fmt.Errorf("error retrieving JWT Spiffe duration-adjusted token count for month: %w", err)
	}

	totalTokens += spiffeJwtUnits

	idTokenDetails := []map[string]interface{}{
		{"type": "oidc", "count": oidcTokenCount},
		{"type": "spiffe", "count": spiffeJwtUnits},
	}

	return map[string]interface{}{
		"metric_name": "id_token_units",
		"metric_data": map[string]interface{}{
			"total":          totalTokens,
			"metric_details": idTokenDetails,
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

func (b *SystemBackend) buildSSHMetric(ctx context.Context, month time.Time) (map[string]interface{}, error) {
	certCounts, err := b.Core.GetStoredSSHDurationAdjustedCertCount(ctx, month)
	if err != nil {
		return nil, fmt.Errorf("error retrieving SSH duration-adjuested cert counts for current month: %w", err)
	}

	otpCounts, err := b.Core.GetStoredSSHOTPCount(ctx, month)
	if err != nil {
		return nil, fmt.Errorf("error retrieving SSH OTP counts for current month: %w", err)
	}

	return map[string]interface{}{
		"metric_name": "ssh_units",
		"metric_data": map[string]interface{}{
			"total": certCounts + float64(otpCounts),
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
	}, nil
}
