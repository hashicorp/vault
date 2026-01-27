// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package vault

import (
	"context"
	"time"

	"github.com/hashicorp/vault/helper/timeutil"
	"github.com/hashicorp/vault/vault/billing"
)

func (c *Core) setupConsumptionBilling(ctx context.Context) error {
	// We need replication (post unseal) to start before we run the consumption billing metrics worker
	// This is because there is primary/secondary cluster specific logic
	c.consumptionBillingLock.Lock()
	c.consumptionBilling = &billing.ConsumptionBilling{
		BillingConfig: c.billingConfig,
	}
	c.consumptionBillingLock.Unlock()
	c.postUnsealFuncs = append(c.postUnsealFuncs, func() {
		c.consumptionBillingMetricsWorker(ctx)
	})

	return nil
}

func (c *Core) consumptionBillingMetricsWorker(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		c.consumptionBillingLock.RLock()
		if c.consumptionBilling.BillingConfig.MetricsUpdateCadence > 0 {
			// For testing purposes
			ticker = time.NewTicker(c.consumptionBilling.BillingConfig.MetricsUpdateCadence)
		}
		c.consumptionBillingLock.RUnlock()
		defer ticker.Stop()

		endOfMonth := time.NewTimer(time.Until(timeutil.StartOfNextMonth(time.Now())))
		for {
			select {
			case <-ticker.C:
				if err := c.updateBillingMetrics(ctx); err != nil {
					c.logger.Error("error updating billing metrics", "error", err)
				}
			case <-ctx.Done():
				return
			case <-endOfMonth.C:
				// Reset the timer for the next month
				endOfMonth.Reset(time.Until(timeutil.StartOfNextMonth(time.Now())))
				// On month boundary, we need to flush the current in-memory counts to storage
				if err := c.updateBillingMetrics(ctx); err != nil {
					c.logger.Error("error updating billing metrics at month boundary", "error", err)
				}
			}
		}
	}()
}

func (c *Core) updateBillingMetrics(ctx context.Context) error {
	if c.PerfStandby() {
		// We do not update billing metrics on performance standbys
		// Instead we send any in memory counts to the primary. This doesn't apply
		// to role counts, but will be used for other metrics
	} else if standby, _ := c.Standby(); standby {
		// Do nothing if we are a standby. All requests get forwarded anyway
	} else {
		// The active node will need to flush max role counts to storage
		currentMonth := time.Now()
		if c.isPrimary() {
			c.UpdateReplicatedHWMMetrics(ctx, currentMonth)
		}
		c.UpdateLocalHWMMetrics(ctx, currentMonth)
	}
	return nil
}

func (c *Core) UpdateReplicatedHWMMetrics(ctx context.Context, currentMonth time.Time) error {
	_, err := c.UpdateMaxRoleCounts(ctx, billing.ReplicatedPrefix, currentMonth)
	if err != nil {
		c.logger.Error("error updating replicated max role counts", "error", err)
		// We won't return an error. Instead we will log the errors and attempt to continue
	} else {
		c.logger.Info("updated replicated hwm role counts", "prefix", billing.ReplicatedPrefix, "currentMonth", currentMonth)
	}
	if _, err = c.UpdateMaxKvCounts(ctx, billing.ReplicatedPrefix, currentMonth); err != nil {
		// We won't return an error. Instead we will log the errors and attempt to continue
		c.logger.Error("error updating replicated max kv counts", "error", err)
	} else {
		c.logger.Info("updated replicated max kv counts", "prefix", billing.ReplicatedPrefix, "currentMonth", currentMonth)
	}
	return nil
}

func (c *Core) UpdateLocalHWMMetrics(ctx context.Context, currentMonth time.Time) error {
	if _, err := c.UpdateMaxRoleCounts(ctx, billing.LocalPrefix, currentMonth); err != nil {
		c.logger.Error("error updating local max role counts", "error", err)
	} else {
		c.logger.Info("updated local max role counts", "prefix", billing.LocalPrefix, "currentMonth", currentMonth)
	}
	if _, err := c.UpdateMaxKvCounts(ctx, billing.LocalPrefix, currentMonth); err != nil {
		c.logger.Error("error updating local max kv counts", "error", err)
	} else {
		c.logger.Info("updated local max kv counts", "prefix", billing.LocalPrefix, "currentMonth", currentMonth)
	}
	// The count of external plugins is per cluster, and we do not de-duplicate across clusters.
	// For that reason, we will always store the count at the "local" prefix, so that the count does not
	// get replicated.
	if _, err := c.UpdateMaxThirdPartyPluginCounts(ctx, currentMonth); err != nil {
		c.logger.Error("error updating local max external plugin counts", "error", err)
	} else {
		c.logger.Info("updated local max external plugin counts", "prefix", billing.LocalPrefix, "currentMonth", currentMonth)
	}
	return nil
}
