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
		a.DatabaseDynamicRoles + b.DatabaseDynamicRoles,
		a.DatabaseStaticRoles + b.DatabaseStaticRoles,
		a.GCPRolesets + b.GCPRolesets,
		a.GCPStaticAccounts + b.GCPStaticAccounts,
		a.GCPImpersonatedAccounts + b.GCPImpersonatedAccounts,
		a.LDAPDynamicRoles + b.LDAPDynamicRoles,
		a.LDAPStaticRoles + b.LDAPStaticRoles,
		a.OpenLDAPDynamicRoles + b.OpenLDAPDynamicRoles,
		a.OpenLDAPStaticRoles + b.OpenLDAPStaticRoles,
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
	if currentRoleCounts.AWSDynamicRoles > maxRoleCounts.AWSDynamicRoles {
		maxRoleCounts.AWSDynamicRoles = currentRoleCounts.AWSDynamicRoles
	}
	if currentRoleCounts.AzureDynamicRoles > maxRoleCounts.AzureDynamicRoles {
		maxRoleCounts.AzureDynamicRoles = currentRoleCounts.AzureDynamicRoles
	}
	if currentRoleCounts.GCPRolesets > maxRoleCounts.GCPRolesets {
		maxRoleCounts.GCPRolesets = currentRoleCounts.GCPRolesets
	}
	if currentRoleCounts.AWSStaticRoles > maxRoleCounts.AWSStaticRoles {
		maxRoleCounts.AWSStaticRoles = currentRoleCounts.AWSStaticRoles
	}
	if currentRoleCounts.DatabaseDynamicRoles > maxRoleCounts.DatabaseDynamicRoles {
		maxRoleCounts.DatabaseDynamicRoles = currentRoleCounts.DatabaseDynamicRoles
	}
	if currentRoleCounts.OpenLDAPStaticRoles > maxRoleCounts.OpenLDAPStaticRoles {
		maxRoleCounts.OpenLDAPStaticRoles = currentRoleCounts.OpenLDAPStaticRoles
	}
	if currentRoleCounts.OpenLDAPDynamicRoles > maxRoleCounts.OpenLDAPDynamicRoles {
		maxRoleCounts.OpenLDAPDynamicRoles = currentRoleCounts.OpenLDAPDynamicRoles
	}
	if currentRoleCounts.LDAPDynamicRoles > maxRoleCounts.LDAPDynamicRoles {
		maxRoleCounts.LDAPDynamicRoles = currentRoleCounts.LDAPDynamicRoles
	}
	if currentRoleCounts.LDAPStaticRoles > maxRoleCounts.LDAPStaticRoles {
		maxRoleCounts.LDAPStaticRoles = currentRoleCounts.LDAPStaticRoles
	}
	if currentRoleCounts.DatabaseStaticRoles > maxRoleCounts.DatabaseStaticRoles {
		maxRoleCounts.DatabaseStaticRoles = currentRoleCounts.DatabaseStaticRoles
	}
	if currentRoleCounts.GCPImpersonatedAccounts > maxRoleCounts.GCPImpersonatedAccounts {
		maxRoleCounts.GCPImpersonatedAccounts = currentRoleCounts.GCPImpersonatedAccounts
	}
	if currentRoleCounts.GCPStaticAccounts > maxRoleCounts.GCPStaticAccounts {
		maxRoleCounts.GCPStaticAccounts = currentRoleCounts.GCPStaticAccounts
	}
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

func (c *Core) GetBillingSubView() *BarrierView {
	return c.systemBarrierView.SubView(billing.BillingSubPath)
}
