// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package vault

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/hashicorp/vault/helper/timeutil"
	"github.com/hashicorp/vault/vault/billing"
)

var (
	ErrCouldNotGetBillingSubView        = fmt.Errorf("could not get billing sub view")
	ErrConsumptionBillingNotInitialized = fmt.Errorf("consumption billing is not initialized")
)

func (c *Core) setupConsumptionBilling(ctx context.Context) error {
	// We need replication (post unseal) to start before we run the consumption billing metrics worker
	// This is because there is primary/secondary cluster specific logic
	c.consumptionBillingLock.Lock()
	logger := c.baseLogger.Named("billing")
	c.AddLogger(logger)
	c.consumptionBilling = &billing.ConsumptionBilling{
		BillingConfig: c.billingConfig,
		DataProtectionCallCounts: billing.DataProtectionCallCounts{
			Transit:   &atomic.Uint64{},
			Transform: &atomic.Uint64{},
		},
		Logger: logger,
	}
	c.consumptionBillingLock.Unlock()
	c.postUnsealFuncs = append(c.postUnsealFuncs, func() {
		c.consumptionBillingMetricsWorker(ctx)
		// Start the perf standby plugin counts worker if this is a perf standby
		// Access perfStandby field directly to avoid deadlock during post-unseal
		if c.perfStandby {
			go c.perfStandbyPluginCountsWorker(ctx)
		}
		// Active nodes don't need a separate worker - they flush counts via
		// the existing consumptionBillingMetricsWorker -> updateBillingMetrics path
	})

	return nil
}

func (c *Core) consumptionBillingMetricsWorker(ctx context.Context) {
	go func() {
		c.consumptionBillingLock.RLock()
		// Check if the clock has been overridden for testing purposes
		clock := c.consumptionBilling.BillingConfig.TestOverrideClock
		metricsCadence := c.consumptionBilling.BillingConfig.MetricsUpdateCadence
		c.consumptionBillingLock.RUnlock()
		if clock == nil {
			clock = timeutil.DefaultClock{}
		}

		ticker := clock.NewTicker(billing.BillingWriteInterval)
		if metricsCadence > 0 {
			// For testing purposes
			ticker = clock.NewTicker(metricsCadence)
		}
		defer ticker.Stop()

		untilNextMonth := func(now time.Time) time.Duration {
			// IMPORTANT: Do not use time.Until() here; it uses the real clock.
			// We need the injected clock (for tests) to control time math.
			d := timeutil.StartOfNextMonth(now.UTC()).Sub(now.UTC())
			if d < 0 {
				return 0
			}
			return d
		}
		endOfMonth := clock.NewTimer(untilNextMonth(clock.Now()))
		for {
			select {
			case <-ticker.C:
				if err := c.updateBillingMetrics(ctx, clock.Now()); err != nil {
					c.logger.Error("error updating billing metrics", "error", err)
				}
			case <-ctx.Done():
				return
			case <-endOfMonth.C:
				// Reset the timer for the next month
				currentMonth := clock.Now()
				c.logger.Debug("reached end of month, resetting timer", "currentMonth", currentMonth)
				previousMonth := timeutil.StartOfPreviousMonth(currentMonth)
				// On month boundary, we need to flush the current in-memory counts to storage
				if err := c.updateBillingMetrics(ctx, previousMonth); err != nil {
					c.logger.Error("error updating billing metrics at month boundary", "error", err)
				}
				c.HandleStartOfMonth(ctx, currentMonth)
				endOfMonth.Reset(untilNextMonth(currentMonth))

			}
		}
	}()
}

// HandleStartOfMonth cleans up monthly billing data from
// n-2 months ago, and also resets all in memory billing metrics when the start of the month is reached.
func (c *Core) HandleStartOfMonth(ctx context.Context, currentMonth time.Time) {
	c.logger.Info("handling start of month operations", "currentMonth", currentMonth)
	// We only delete n-2 month billing metrics on the active node
	if standby, _ := c.Standby(); !standby && !c.PerfStandby() {
		if err := c.deletePreviousMonthBillingMetrics(ctx, currentMonth); err != nil {
			c.logger.Error("error deleting historical month billing metrics", "error", err)
		}
	}
	if err := c.resetInMemoryBillingMetrics(); err != nil {
		c.logger.Error("error resetting in memory billing metrics", "error", err)
	}
}

func (c *Core) deletePreviousMonthBillingMetrics(ctx context.Context, currentMonth time.Time) error {
	twoMonthsAgo := timeutil.StartOfPreviousMonth(currentMonth).AddDate(0, -1, 0)
	// Delete billing metrics from both replicated and local prefixes
	for _, pathPrefix := range []string{billing.ReplicatedPrefix, billing.LocalPrefix} {
		// If we are not the primary, then do not delete replicate metrics
		if !c.isPrimary() && pathPrefix == billing.ReplicatedPrefix {
			continue
		}
		billingPath := billing.GetMonthlyBillingPath(pathPrefix, twoMonthsAgo)
		view, ok := c.GetBillingSubView()
		if !ok {
			return ErrCouldNotGetBillingSubView
		}
		metricPaths, err := view.List(ctx, billingPath)
		if err != nil {
			return err
		}
		for _, segment := range metricPaths {
			err = view.Delete(ctx, billingPath+segment)
			if err != nil {
				c.logger.Error("error deleting previous month billing metric", "error", err, "metricPath", billingPath+segment)
			}
		}
	}
	return nil
}

func (c *Core) resetInMemoryBillingMetrics() error {
	// Reset Transit/Tranform DP counts
	c.logger.Info("resetting in memory billing metrics")
	c.consumptionBillingLock.Lock()
	defer c.consumptionBillingLock.Unlock()
	c.consumptionBilling.DataProtectionCallCounts.Transit.Store(0)
	c.consumptionBilling.DataProtectionCallCounts.Transform.Store(0)
	c.consumptionBilling.KmipSeenEnabledThisMonth.Store(false)
	return nil
}

func (c *Core) updateBillingMetrics(ctx context.Context, currentMonth time.Time) error {
	// Check if systemBarrierView is initialized
	c.mountsLock.RLock()
	initialized := c.systemBarrierView != nil
	c.mountsLock.RUnlock()

	if !initialized {
		return nil
	}
	if c.PerfStandby() {
		// We do not update billing metrics on performance standbys
		// Instead we send any in memory counts to the primary. This doesn't apply
		// to role counts, but will be used for other metrics
	} else if standby, _ := c.Standby(); standby {
		// Do nothing if we are a standby. All requests get forwarded anyway
	} else {
		// The active node will need to flush max role counts to storage
		if c.isPrimary() {
			c.UpdateReplicatedHWMMetrics(ctx, currentMonth)
		}
		c.UpdateLocalHWMMetrics(ctx, currentMonth)
		if err := c.UpdateLocalAggregatedMetrics(ctx, currentMonth); err != nil {
			c.logger.Error("error updating cluster data protection call counts", "error", err)
		} else {
			c.logger.Info("updated cluster data protection call counts", "prefix", billing.LocalPrefix, "currentMonth", currentMonth)
		}

	}
	return nil
}

func (c *Core) UpdateReplicatedHWMMetrics(ctx context.Context, currentMonth time.Time) error {
	_, _, err := c.UpdateMaxRoleAndManagedKeyCounts(ctx, billing.ReplicatedPrefix, currentMonth)
	if err != nil {
		c.logger.Error("error updating replicated max role and managed key counts", "error", err)
		// We won't return an error. Instead we will log the errors and attempt to continue
	} else {
		c.logger.Info("updated replicated hwm role and managed key counts", "prefix", billing.ReplicatedPrefix, "currentMonth", currentMonth)
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
	if _, _, err := c.UpdateMaxRoleAndManagedKeyCounts(ctx, billing.LocalPrefix, currentMonth); err != nil {
		c.logger.Error("error updating local max role and managed key counts", "error", err)
	} else {
		c.logger.Info("updated local max role and managed key counts", "prefix", billing.LocalPrefix, "currentMonth", currentMonth)
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
	if _, err := c.UpdateKmipEnabled(ctx, currentMonth); err != nil {
		c.logger.Error("error updating local kmip enabled", "error", err)
	} else {
		c.logger.Info("updated local kmip enabled", "prefix", billing.LocalPrefix, "currentMonth", currentMonth)
	}

	return nil
}

// UpdateLocalAggregatedMetrics updates local metrics that are aggregated across all replicated clusters
func (c *Core) UpdateLocalAggregatedMetrics(ctx context.Context, currentMonth time.Time) error {
	if _, err := c.UpdateTransitCallCounts(ctx, currentMonth); err != nil {
		return fmt.Errorf("could not store transit data protection call counts: %w", err)
	}
	if _, err := c.UpdateTransformCallCounts(ctx, currentMonth); err != nil {
		return fmt.Errorf("could not store transform data protection call counts: %w", err)
	}
	return nil
}
