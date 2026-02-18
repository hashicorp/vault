// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package vault

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault/billing"
)

func (c *Core) storeThirdPartyPluginCountsLocked(ctx context.Context, localPathPrefix string, currentMonth time.Time, thirdPartyPluginCounts int) error {
	billingPath := billing.GetMonthlyBillingMetricPath(localPathPrefix, currentMonth, billing.ThirdPartyPluginsPrefix)
	entry := &logical.StorageEntry{
		Key:   billingPath,
		Value: []byte(strconv.Itoa(thirdPartyPluginCounts)),
	}
	view, ok := c.GetBillingSubView()
	if !ok {
		return nil
	}
	return view.Put(ctx, entry)
}

func (c *Core) getStoredThirdPartyPluginCountsLocked(ctx context.Context, localPathPrefix string, currentMonth time.Time) (int, error) {
	billingPath := billing.GetMonthlyBillingMetricPath(localPathPrefix, currentMonth, billing.ThirdPartyPluginsPrefix)
	view, ok := c.GetBillingSubView()
	if !ok {
		return 0, nil
	}
	entry, err := view.Get(ctx, billingPath)
	if err != nil {
		return 0, err
	}
	if entry == nil {
		return 0, nil
	}
	thirdPartyPluginCounts, err := strconv.Atoi(string(entry.Value))
	if err != nil {
		return 0, err
	}
	return thirdPartyPluginCounts, nil
}

// UpdateMaxThirdPartyPlugins updates the max number of third-party plugins for the given month.
// Note that this count is per cluster. It does NOT de-duplicate across clusters. For that reason,
// we will always store the count at the "local" prefix.
func (c *Core) UpdateMaxThirdPartyPluginCounts(ctx context.Context, currentMonth time.Time) (int, error) {
	c.consumptionBillingLock.RLock()
	cb := c.consumptionBilling
	c.consumptionBillingLock.RUnlock()

	if cb == nil {
		return 0, ErrConsumptionBillingNotInitialized
	}

	cb.BillingStorageLock.Lock()
	defer cb.BillingStorageLock.Unlock()

	previousThirdPartyPluginCounts, err := c.getStoredThirdPartyPluginCountsLocked(ctx, billing.LocalPrefix, currentMonth)
	if err != nil {
		return 0, err
	}
	currentThirdPartyPluginCounts, err := c.ListDeduplicatedExternalSecretPlugins(ctx)
	if err != nil {
		return 0, err
	}
	maxCount := c.compareCounts(previousThirdPartyPluginCounts, len(currentThirdPartyPluginCounts), "Third-Party Plugins")
	err = c.storeThirdPartyPluginCountsLocked(ctx, billing.LocalPrefix, currentMonth, maxCount)
	if err != nil {
		return 0, err
	}
	return maxCount, nil
}

func (c *Core) GetStoredThirdPartyPluginCounts(ctx context.Context, month time.Time) (int, error) {
	c.consumptionBillingLock.RLock()
	cb := c.consumptionBilling
	c.consumptionBillingLock.RUnlock()

	if cb == nil {
		return 0, ErrConsumptionBillingNotInitialized
	}

	cb.BillingStorageLock.RLock()
	defer cb.BillingStorageLock.RUnlock()
	return c.getStoredThirdPartyPluginCountsLocked(ctx, billing.LocalPrefix, month)
}

func combineRoleCounts(a, b *RoleCounts) *RoleCounts {
	if a == nil && b == nil {
		return &RoleCounts{}
	}
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	return &RoleCounts{
		a.AWSDynamicRoles + b.AWSDynamicRoles,
		a.AWSStaticRoles + b.AWSStaticRoles,
		a.AzureDynamicRoles + b.AzureDynamicRoles,
		a.AzureStaticRoles + b.AzureStaticRoles,
		a.DatabaseDynamicRoles + b.DatabaseDynamicRoles,
		a.DatabaseStaticRoles + b.DatabaseStaticRoles,
		a.GCPRolesets + b.GCPRolesets,
		a.GCPStaticAccounts + b.GCPStaticAccounts,
		a.GCPImpersonatedAccounts + b.GCPImpersonatedAccounts,
		a.LDAPDynamicRoles + b.LDAPDynamicRoles,
		a.LDAPStaticRoles + b.LDAPStaticRoles,
		a.OpenLDAPDynamicRoles + b.OpenLDAPDynamicRoles,
		a.OpenLDAPStaticRoles + b.OpenLDAPStaticRoles,
		a.AlicloudDynamicRoles + b.AlicloudDynamicRoles,
		a.RabbitMQDynamicRoles + b.RabbitMQDynamicRoles,
		a.ConsulDynamicRoles + b.ConsulDynamicRoles,
		a.NomadDynamicRoles + b.NomadDynamicRoles,
		a.KubernetesDynamicRoles + b.KubernetesDynamicRoles,
		a.MongoDBAtlasDynamicRoles + b.MongoDBAtlasDynamicRoles,
		a.TerraformCloudDynamicRoles + b.TerraformCloudDynamicRoles,
	}
}

// storeMaxKvCountsLocked must be called with BillingStorageLock held
func (c *Core) storeMaxKvCountsLocked(ctx context.Context, maxKvCounts int, localPathPrefix string, month time.Time) error {
	billingPath := billing.GetMonthlyBillingMetricPath(localPathPrefix, month, billing.KvHWMCountsHWM)
	entry := &logical.StorageEntry{
		Key:   billingPath,
		Value: []byte(strconv.Itoa(maxKvCounts)),
	}
	view, ok := c.GetBillingSubView()
	if !ok {
		return nil
	}
	return view.Put(ctx, entry)
}

// getStoredMaxKvCountsLocked must be called with BillingStorageLock held
func (c *Core) getStoredMaxKvCountsLocked(ctx context.Context, localPathPrefix string, month time.Time) (int, error) {
	billingPath := billing.GetMonthlyBillingMetricPath(localPathPrefix, month, billing.KvHWMCountsHWM)
	view, ok := c.GetBillingSubView()
	if !ok {
		return 0, nil
	}
	entry, err := view.Get(ctx, billingPath)
	if err != nil {
		return 0, err
	}
	if entry == nil {
		return 0, nil
	}
	maxKvCounts, err := strconv.Atoi(string(entry.Value))
	if err != nil {
		return 0, err
	}
	return maxKvCounts, nil
}

func (c *Core) GetStoredHWMKvCounts(ctx context.Context, localPathPrefix string, month time.Time) (int, error) {
	c.consumptionBillingLock.RLock()
	cb := c.consumptionBilling
	c.consumptionBillingLock.RUnlock()

	if cb == nil {
		return 0, ErrConsumptionBillingNotInitialized
	}

	cb.BillingStorageLock.RLock()
	defer cb.BillingStorageLock.RUnlock()
	return c.getStoredMaxKvCountsLocked(ctx, localPathPrefix, month)
}

// UpdateMaxKvCounts updates the HWM kv counts for the given month, and returns the value that was stored.
func (c *Core) UpdateMaxKvCounts(ctx context.Context, localPathPrefix string, currentMonth time.Time) (int, error) {
	c.consumptionBillingLock.RLock()
	cb := c.consumptionBilling
	c.consumptionBillingLock.RUnlock()

	if cb == nil {
		return 0, ErrConsumptionBillingNotInitialized
	}

	cb.BillingStorageLock.Lock()
	defer cb.BillingStorageLock.Unlock()

	local := localPathPrefix == billing.LocalPrefix

	// Get the current count of kv version 1 secrets
	currentKvCounts, err := c.GetKvUsageMetricsByNamespace(ctx, "1", "", local, !local, false)
	if err != nil {
		c.logger.Error("error getting count of kv version 1 secrets", "error", err)
		return 0, err
	}
	totalKvCounts := getTotalSecretsAcrossAllNamespaces(currentKvCounts)

	// Get the current count of kv version 2 secrets
	currentKvCounts, err = c.GetKvUsageMetricsByNamespace(ctx, "2", "", local, !local, false)
	if err != nil {
		c.logger.Error("error getting current count of kv version 2 secrets", "error", err)
		return 0, err
	}
	totalKvCounts += getTotalSecretsAcrossAllNamespaces(currentKvCounts)

	// Get the stored max kv counts
	maxKvCounts, err := c.getStoredMaxKvCountsLocked(ctx, localPathPrefix, currentMonth)
	if err != nil {
		c.logger.Error("error getting stored max kv counts", "error", err)
		return 0, err
	}
	if maxKvCounts == 0 {
		maxKvCounts = totalKvCounts
	}
	if totalKvCounts > maxKvCounts {
		c.logger.Info("updating max kv counts", "totalKvCounts", totalKvCounts, "maxKvCounts", maxKvCounts)
		maxKvCounts = totalKvCounts
	}
	err = c.storeMaxKvCountsLocked(ctx, maxKvCounts, localPathPrefix, currentMonth)
	if err != nil {
		c.logger.Error("error storing max kv counts", "error", err)
		return 0, err
	}
	return maxKvCounts, nil
}

// storeMaxRoleCountsLocked must be called with BillingStorageLock held
func (c *Core) storeMaxRoleCountsLocked(ctx context.Context, maxRoleCounts *RoleCounts, localPathPrefix string, month time.Time) error {
	billingPath := billing.GetMonthlyBillingMetricPath(localPathPrefix, month, billing.RoleHWMCountsHWM)
	entry, err := logical.StorageEntryJSON(billingPath, maxRoleCounts)
	if err != nil {
		return err
	}
	view, ok := c.GetBillingSubView()
	if !ok {
		return nil
	}
	return view.Put(ctx, entry)
}

func (c *Core) UpdateMaxRoleCounts(ctx context.Context, localPathPrefix string, currentMonth time.Time) (*RoleCounts, error) {
	c.consumptionBillingLock.RLock()
	cb := c.consumptionBilling
	c.consumptionBillingLock.RUnlock()

	if cb == nil {
		return nil, ErrConsumptionBillingNotInitialized
	}

	cb.BillingStorageLock.Lock()
	defer cb.BillingStorageLock.Unlock()

	local := localPathPrefix == billing.LocalPrefix
	currentRoleCounts := c.getRoleCountsInternal(local, !local, true)

	maxRoleCounts, err := c.getStoredRoleCountsLocked(ctx, localPathPrefix, currentMonth)
	if maxRoleCounts == nil {
		maxRoleCounts = &RoleCounts{}
	}
	if currentRoleCounts == nil {
		currentRoleCounts = &RoleCounts{}
	}
	maxRoleCounts.AWSDynamicRoles = c.compareCounts(currentRoleCounts.AWSDynamicRoles, maxRoleCounts.AWSDynamicRoles, "AWS Dynamic Roles")
	maxRoleCounts.AzureDynamicRoles = c.compareCounts(currentRoleCounts.AzureDynamicRoles, maxRoleCounts.AzureDynamicRoles, "Azure Dynamic Roles")
	maxRoleCounts.AzureStaticRoles = c.compareCounts(currentRoleCounts.AzureStaticRoles, maxRoleCounts.AzureStaticRoles, "Azure Static Roles")
	maxRoleCounts.GCPRolesets = c.compareCounts(currentRoleCounts.GCPRolesets, maxRoleCounts.GCPRolesets, "GCP Rolesets")
	maxRoleCounts.AWSStaticRoles = c.compareCounts(currentRoleCounts.AWSStaticRoles, maxRoleCounts.AWSStaticRoles, "AWS Static Roles")
	maxRoleCounts.DatabaseDynamicRoles = c.compareCounts(currentRoleCounts.DatabaseDynamicRoles, maxRoleCounts.DatabaseDynamicRoles, "Database Dynamic Roles")
	maxRoleCounts.OpenLDAPStaticRoles = c.compareCounts(currentRoleCounts.OpenLDAPStaticRoles, maxRoleCounts.OpenLDAPStaticRoles, "OpenLDAP Static Roles")
	maxRoleCounts.OpenLDAPDynamicRoles = c.compareCounts(currentRoleCounts.OpenLDAPDynamicRoles, maxRoleCounts.OpenLDAPDynamicRoles, "OpenLDAP Dynamic Roles")
	maxRoleCounts.LDAPDynamicRoles = c.compareCounts(currentRoleCounts.LDAPDynamicRoles, maxRoleCounts.LDAPDynamicRoles, "LDAP Dynamic Roles")
	maxRoleCounts.LDAPStaticRoles = c.compareCounts(currentRoleCounts.LDAPStaticRoles, maxRoleCounts.LDAPStaticRoles, "LDAP Static Roles")
	maxRoleCounts.DatabaseStaticRoles = c.compareCounts(currentRoleCounts.DatabaseStaticRoles, maxRoleCounts.DatabaseStaticRoles, "Database Static Roles")
	maxRoleCounts.GCPImpersonatedAccounts = c.compareCounts(currentRoleCounts.GCPImpersonatedAccounts, maxRoleCounts.GCPImpersonatedAccounts, "GCPImpersonated Accounts")
	maxRoleCounts.GCPStaticAccounts = c.compareCounts(currentRoleCounts.GCPStaticAccounts, maxRoleCounts.GCPStaticAccounts, "GCP Static Accounts")
	maxRoleCounts.AlicloudDynamicRoles = c.compareCounts(currentRoleCounts.AlicloudDynamicRoles, maxRoleCounts.AlicloudDynamicRoles, "Alicloud Dynamic Roles")
	maxRoleCounts.RabbitMQDynamicRoles = c.compareCounts(currentRoleCounts.RabbitMQDynamicRoles, maxRoleCounts.RabbitMQDynamicRoles, "RabbitMQ Dynamic Roles")
	maxRoleCounts.ConsulDynamicRoles = c.compareCounts(currentRoleCounts.ConsulDynamicRoles, maxRoleCounts.ConsulDynamicRoles, "Consul Dynamic Roles")
	maxRoleCounts.NomadDynamicRoles = c.compareCounts(currentRoleCounts.NomadDynamicRoles, maxRoleCounts.NomadDynamicRoles, "Nomad Dynamic Roles")
	maxRoleCounts.KubernetesDynamicRoles = c.compareCounts(currentRoleCounts.KubernetesDynamicRoles, maxRoleCounts.KubernetesDynamicRoles, "Kubernetes Dynamic Roles")
	maxRoleCounts.MongoDBAtlasDynamicRoles = c.compareCounts(currentRoleCounts.MongoDBAtlasDynamicRoles, maxRoleCounts.MongoDBAtlasDynamicRoles, "MongoDB Atlas Dynamic Roles")
	maxRoleCounts.TerraformCloudDynamicRoles = c.compareCounts(currentRoleCounts.TerraformCloudDynamicRoles, maxRoleCounts.TerraformCloudDynamicRoles, "Terraform Cloud Dynamic Roles")

	err = c.storeMaxRoleCountsLocked(ctx, maxRoleCounts, localPathPrefix, currentMonth)
	if err != nil {
		return nil, err
	}

	return maxRoleCounts, nil
}

func (c *Core) GetStoredHWMRoleCounts(ctx context.Context, localPathPrefix string, month time.Time) (*RoleCounts, error) {
	c.consumptionBillingLock.RLock()
	cb := c.consumptionBilling
	c.consumptionBillingLock.RUnlock()

	if cb == nil {
		return nil, ErrConsumptionBillingNotInitialized
	}

	cb.BillingStorageLock.RLock()
	defer cb.BillingStorageLock.RUnlock()
	return c.getStoredRoleCountsLocked(ctx, localPathPrefix, month)
}

func (c *Core) getStoredRoleCountsLocked(ctx context.Context, localPathPrefix string, month time.Time) (*RoleCounts, error) {
	billingPath := billing.GetMonthlyBillingMetricPath(localPathPrefix, month, billing.RoleHWMCountsHWM)
	var maxRoleCounts *RoleCounts
	view, ok := c.GetBillingSubView()
	if !ok {
		return &RoleCounts{}, nil
	}
	maxRoleCountsRaw, err := view.Get(ctx, billingPath)
	if err != nil {
		return nil, err
	}
	if maxRoleCountsRaw == nil {
		return &RoleCounts{}, nil
	}
	if err := maxRoleCountsRaw.DecodeJSON(&maxRoleCounts); err != nil {
		return nil, err
	}
	return maxRoleCounts, nil
}

func (c *Core) compareCounts(current, previous int, metricName string) int {
	if previous > current {
		return previous
	}
	c.logger.Debug("updating max counts", "metricName", metricName, "previous", previous, "current", current)
	return current
}

func (c *Core) GetBillingSubView() (*BarrierView, bool) {
	c.mountsLock.RLock()
	view := c.systemBarrierView
	c.mountsLock.RUnlock()

	if view == nil {
		return nil, false
	}
	return view.SubView(billing.BillingSubPath), true
}

// storeTransitCallCountsLocked must be called with BillingStorageLock held
func (c *Core) storeTransitCallCountsLocked(ctx context.Context, transitCount uint64, localPathPrefix string, month time.Time) error {
	// Store count for each data protection type separately because they are atomic counters
	billingPath := billing.GetMonthlyBillingMetricPath(localPathPrefix, month, billing.TransitDataProtectionCallCountsPrefix)
	entry := &logical.StorageEntry{
		Key:   billingPath,
		Value: []byte(strconv.FormatUint(transitCount, 10)),
	}
	view, ok := c.GetBillingSubView()
	if !ok {
		return nil
	}
	return view.Put(ctx, entry)
}

// getStoredTransitCallCountsLocked must be called with BillingStorageLock held
func (c *Core) getStoredTransitCallCountsLocked(ctx context.Context, localPathPrefix string, month time.Time) (uint64, error) {
	// Retrieve count for each data protection type separately because they are atomic counters
	billingPath := billing.GetMonthlyBillingMetricPath(localPathPrefix, month, billing.TransitDataProtectionCallCountsPrefix)
	view, ok := c.GetBillingSubView()
	if !ok {
		return 0, nil
	}
	entry, err := view.Get(ctx, billingPath)
	if err != nil {
		return 0, err
	}
	if entry == nil {
		return 0, nil
	}
	transitCount, err := strconv.ParseUint(string(entry.Value), 10, 64)
	if err != nil {
		return 0, err
	}
	return transitCount, nil
}

func (c *Core) GetStoredTransitCallCounts(ctx context.Context, month time.Time) (uint64, error) {
	c.consumptionBillingLock.RLock()
	cb := c.consumptionBilling
	c.consumptionBillingLock.RUnlock()

	if cb == nil {
		return 0, ErrConsumptionBillingNotInitialized
	}

	cb.BillingStorageLock.RLock()
	defer cb.BillingStorageLock.RUnlock()
	return c.getStoredTransitCallCountsLocked(ctx, billing.LocalPrefix, month)
}

func (c *Core) UpdateTransitCallCounts(ctx context.Context, currentMonth time.Time) (uint64, error) {
	c.consumptionBillingLock.RLock()
	cb := c.consumptionBilling
	c.consumptionBillingLock.RUnlock()

	if cb == nil {
		return 0, ErrConsumptionBillingNotInitialized
	}
	cb.BillingStorageLock.Lock()
	defer cb.BillingStorageLock.Unlock()
	storedTransitCount, err := c.getStoredTransitCallCountsLocked(ctx, billing.LocalPrefix, currentMonth)
	if err != nil {
		return 0, err
	}

	// Sum the current count with the stored count
	transitCount := cb.DataProtectionCallCounts.Transit.Swap(0) + storedTransitCount

	err = c.storeTransitCallCountsLocked(ctx, transitCount, billing.LocalPrefix, currentMonth)
	if err != nil {
		return 0, err
	}

	return transitCount, nil
}

func (c *Core) storeKmipEnabledLocked(ctx context.Context, localPathPrefix string, currentMonth time.Time, kmipEnabled bool) error {
	billingPath := billing.GetMonthlyBillingMetricPath(localPathPrefix, currentMonth, billing.KmipEnabledPrefix)
	entry, err := logical.StorageEntryJSON(billingPath, kmipEnabled)
	if err != nil {
		return err
	}
	view, ok := c.GetBillingSubView()
	if !ok {
		return nil
	}
	return view.Put(ctx, entry)
}

func (c *Core) getStoredKmipEnabledLocked(ctx context.Context, localPathPrefix string, currentMonth time.Time) (bool, error) {
	billingPath := billing.GetMonthlyBillingMetricPath(localPathPrefix, currentMonth, billing.KmipEnabledPrefix)
	view, ok := c.GetBillingSubView()
	if !ok {
		return false, nil
	}
	entry, err := view.Get(ctx, billingPath)
	if err != nil {
		return false, err
	}
	if entry == nil {
		return false, nil
	}
	var kmipEnabled bool
	if err := entry.DecodeJSON(&kmipEnabled); err != nil {
		return false, err
	}
	return kmipEnabled, nil
}

func (c *Core) GetStoredKmipEnabled(ctx context.Context, currentMonth time.Time) (bool, error) {
	c.consumptionBillingLock.RLock()
	cb := c.consumptionBilling
	c.consumptionBillingLock.RUnlock()

	if cb == nil {
		return false, ErrConsumptionBillingNotInitialized
	}

	cb.BillingStorageLock.RLock()
	defer cb.BillingStorageLock.RUnlock()
	return c.getStoredKmipEnabledLocked(ctx, billing.LocalPrefix, currentMonth)
}

// UpdateKmipEnabled updates the KMIP enabled status for the current month.
// Note that each cluster is billed independently, so we only store the status at the local prefix.
// Additionally, KMIP usage detection covers both local and replicated mounts, meaning if primary has KMIP,
// secondary also detects it and gets charged. This is intentional, as the KMIP usage is per cluster.
// We only store true when KMIP is enabled; we never store false. This means storing true multiple times
// is idempotent and safe.
func (c *Core) UpdateKmipEnabled(ctx context.Context, currentMonth time.Time) (bool, error) {
	c.consumptionBillingLock.RLock()
	cb := c.consumptionBilling
	c.consumptionBillingLock.RUnlock()

	if cb == nil {
		return false, ErrConsumptionBillingNotInitialized
	}

	cb.BillingStorageLock.Lock()
	defer cb.BillingStorageLock.Unlock()

	// Check if KMIP is currently enabled, including replicated mounts
	kmipEnabled, err := c.IsKMIPEnabled(ctx)
	if err != nil {
		return false, err
	}

	if kmipEnabled {
		if err := c.storeKmipEnabledLocked(ctx, billing.LocalPrefix, currentMonth, true); err != nil {
			return false, err
		}
	}

	return kmipEnabled, nil
}

// GetStoredPkiDurationAdjustedCount retrieves the stored PKI duration-adjusted certificate count
// for the specified month. The count is stored as a float64 string with 4 decimal places of precision.
// Returns 0 if no count has been stored for the given month.
func (c *Core) GetStoredPkiDurationAdjustedCount(ctx context.Context, currentMonth time.Time) (float64, error) {
	c.consumptionBillingLock.RLock()
	cb := c.consumptionBilling
	c.consumptionBillingLock.RUnlock()

	if cb == nil {
		return 0, errors.New("consumption billing is not initialized")
	}

	cb.BillingStorageLock.RLock()
	defer cb.BillingStorageLock.RUnlock()

	return c.getStoredPkiDurationAdjustedCountLocked(ctx, billing.LocalPrefix, currentMonth)
}

// UpdatePkiDurationAdjustedCount increments the stored PKI duration-adjusted certificate count
// for the specified month by the given increment value. The increment must be non-negative.
// The count is stored as a float64 string with 4 decimal places of precision.
func (c *Core) UpdatePkiDurationAdjustedCount(ctx context.Context, inc float64, currentMonth time.Time) error {
	if inc < 0 {
		return fmt.Errorf("PKI duration-adjusted increment must be non-negative, got %f", inc)
	}

	if c.consumptionBilling == nil {
		return errors.New("consumption billing is not initialized")
	}

	c.consumptionBilling.BillingStorageLock.Lock()
	defer c.consumptionBilling.BillingStorageLock.Unlock()

	return c.storePkiDurationAdjustedCountLocked(ctx, billing.LocalPrefix, currentMonth, inc)
}

func (c *Core) getStoredPkiDurationAdjustedCountLocked(ctx context.Context, localPathPrefix string, currentMonth time.Time) (float64, error) {
	billingPath := billing.GetMonthlyBillingMetricPath(localPathPrefix, currentMonth, billing.PkiDurationAdjustedCountPrefix)

	view, ok := c.GetBillingSubView()
	if !ok {
		return 0, errors.New("error reading PKI duration-adjusted count: billing subview not available")
	}

	se, err := view.Get(ctx, billingPath)
	if se == nil || err != nil {
		return 0, err
	}

	currentCount, err := strconv.ParseFloat(string(se.Value), 64)
	if err != nil {
		return 0, fmt.Errorf("error decoding current PKI duration adjusted cert count: %w", err)
	}

	return currentCount, nil
}

func (c *Core) storePkiDurationAdjustedCountLocked(ctx context.Context, localPathPrefix string, currentMonth time.Time, inc float64) error {
	currentCount, err := c.getStoredPkiDurationAdjustedCountLocked(ctx, localPathPrefix, currentMonth)
	if err != nil {
		return err
	}

	billingPath := billing.GetMonthlyBillingMetricPath(localPathPrefix, currentMonth, billing.PkiDurationAdjustedCountPrefix)
	view, ok := c.GetBillingSubView()
	if !ok {
		return errors.New("error storing PKI duration-adjusted count: billing subview not available")
	}

	// Write new value
	newCount := currentCount + inc
	entry := &logical.StorageEntry{
		Key:   billingPath,
		Value: []byte(strconv.FormatFloat(newCount, 'f', 4, 64)),
	}

	if err := view.Put(ctx, entry); err != nil {
		return fmt.Errorf("error writing PKI duration adjusted cert count: %w", err)
	}

	return nil
}
