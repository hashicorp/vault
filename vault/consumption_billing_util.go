// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package vault

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault/billing"
)

func (c *Core) storeThirdPartyPluginCountsLocked(ctx context.Context, localPathPrefix string, currentMonth time.Time, thirdPartyPluginCounts int) error {
	billingPath := billing.GetMonthlyBillingPath(localPathPrefix, currentMonth, billing.ThirdPartyPluginsPrefix)
	entry := &logical.StorageEntry{
		Key:   billingPath,
		Value: []byte(strconv.Itoa(thirdPartyPluginCounts)),
	}
	return c.GetBillingSubView().Put(ctx, entry)
}

func (c *Core) getStoredThirdPartyPluginCountsLocked(ctx context.Context, localPathPrefix string, currentMonth time.Time) (int, error) {
	billingPath := billing.GetMonthlyBillingPath(localPathPrefix, currentMonth, billing.ThirdPartyPluginsPrefix)
	entry, err := c.GetBillingSubView().Get(ctx, billingPath)
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
	c.consumptionBilling.BillingStorageLock.Lock()
	defer c.consumptionBilling.BillingStorageLock.Unlock()

	previousThirdPartyPluginCounts, err := c.getStoredThirdPartyPluginCountsLocked(ctx, billing.LocalPrefix, currentMonth)
	if err != nil {
		return 0, err
	}
	currentThirdPartyPluginCounts, err := c.ListExternalSecretPlugins(ctx)
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
	c.consumptionBilling.BillingStorageLock.RLock()
	defer c.consumptionBilling.BillingStorageLock.RUnlock()
	return c.getStoredThirdPartyPluginCountsLocked(ctx, billing.LocalPrefix, month)
}

func combineRoleCounts(ctx context.Context, a, b *RoleCounts) *RoleCounts {
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
	billingPath := billing.GetMonthlyBillingPath(localPathPrefix, month, billing.KvHWMCountsHWM)
	entry := &logical.StorageEntry{
		Key:   billingPath,
		Value: []byte(strconv.Itoa(maxKvCounts)),
	}
	return c.GetBillingSubView().Put(ctx, entry)
}

// getStoredMaxKvCountsLocked must be called with BillingStorageLock held
func (c *Core) getStoredMaxKvCountsLocked(ctx context.Context, localPathPrefix string, month time.Time) (int, error) {
	billingPath := billing.GetMonthlyBillingPath(localPathPrefix, month, billing.KvHWMCountsHWM)
	entry, err := c.GetBillingSubView().Get(ctx, billingPath)
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
	c.consumptionBilling.BillingStorageLock.RLock()
	defer c.consumptionBilling.BillingStorageLock.RUnlock()
	return c.getStoredMaxKvCountsLocked(ctx, localPathPrefix, month)
}

// UpdateMaxKvCounts updates the HWM kv counts for the given month, and returns the value that was stored.
func (c *Core) UpdateMaxKvCounts(ctx context.Context, localPathPrefix string, currentMonth time.Time) (int, error) {
	c.consumptionBilling.BillingStorageLock.Lock()
	defer c.consumptionBilling.BillingStorageLock.Unlock()

	local := localPathPrefix == billing.LocalPrefix

	// Get the current count of kv version 1 secrets
	currentKvCounts, err := c.GetKvUsageMetricsByNamespace(ctx, "1", "", local, !local)
	if err != nil {
		c.logger.Error("error getting count of kv version 1 secrets", "error", err)
		return 0, err
	}
	totalKvCounts := getTotalSecretsAcrossAllNamespaces(currentKvCounts)

	// Get the current count of kv version 2 secrets
	currentKvCounts, err = c.GetKvUsageMetricsByNamespace(ctx, "2", "", local, !local)
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
	billingPath := billing.GetMonthlyBillingPath(localPathPrefix, month, billing.RoleHWMCountsHWM)
	entry, err := logical.StorageEntryJSON(billingPath, maxRoleCounts)
	if err != nil {
		return err
	}
	return c.GetBillingSubView().Put(ctx, entry)
}

func (c *Core) UpdateMaxRoleCounts(ctx context.Context, localPathPrefix string, currentMonth time.Time) (*RoleCounts, error) {
	c.consumptionBilling.BillingStorageLock.Lock()
	defer c.consumptionBilling.BillingStorageLock.Unlock()

	local := localPathPrefix == billing.LocalPrefix
	currentRoleCounts := c.getRoleCountsInternal(local, !local)

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
	c.consumptionBilling.BillingStorageLock.RLock()
	defer c.consumptionBilling.BillingStorageLock.RUnlock()
	return c.getStoredRoleCountsLocked(ctx, localPathPrefix, month)
}

func (c *Core) getStoredRoleCountsLocked(ctx context.Context, localPathPrefix string, month time.Time) (*RoleCounts, error) {
	billingPath := billing.GetMonthlyBillingPath(localPathPrefix, month, billing.RoleHWMCountsHWM)
	var maxRoleCounts *RoleCounts
	maxRoleCountsRaw, err := c.GetBillingSubView().Get(ctx, billingPath)
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

func (c *Core) GetBillingSubView() *BarrierView {
	return c.systemBarrierView.SubView(billing.BillingSubPath)
}
