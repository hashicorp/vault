// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: MPL-2.0

package vault

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault/billing"
)

// storeAttributionDataLocked writes attribution data to the given view.
// Must be called with BillingStorageLock held.
func storeAttributionDataLocked(ctx context.Context, view logical.Storage, localPathPrefix string, month time.Time, attributionMetricName string, data *logical.MetricTypeAttribution) error {
	billingPath := billing.GetAttributionMaxPath(localPathPrefix, month, attributionMetricName)

	entry, err := logical.StorageEntryJSON(billingPath, data)
	if err != nil {
		return fmt.Errorf("failed to create storage entry for attribution data: %w", err)
	}

	return view.Put(ctx, entry)
}

// getStoredAttributionDataLocked reads attribution data from the given view.
// Must be called with BillingStorageLock held.
func getStoredAttributionDataLocked(ctx context.Context, view logical.Storage, localPathPrefix string, month time.Time, attributionMetricName string) (*logical.MetricTypeAttribution, error) {
	billingPath := billing.GetAttributionMaxPath(localPathPrefix, month, attributionMetricName)

	entry, err := view.Get(ctx, billingPath)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve attribution data: %w", err)
	}

	if entry == nil {
		return &logical.MetricTypeAttribution{}, nil
	}

	var data logical.MetricTypeAttribution
	if err := jsonutil.DecodeJSON(entry.Value, &data); err != nil {
		return nil, fmt.Errorf("failed to decode attribution data: %w", err)
	}

	return &data, nil
}

func (c *Core) GetStoredAttributionData(ctx context.Context, localPathPrefix string, month time.Time, attributionMetricName string) (*logical.MetricTypeAttribution, error) {
	c.consumptionBillingLock.RLock()
	cb := c.consumptionBilling
	c.consumptionBillingLock.RUnlock()

	if cb == nil {
		return nil, ErrConsumptionBillingNotInitialized
	}

	view, ok := c.GetBillingSubView()
	if !ok {
		return nil, errors.New("error reading attribution data: billing subview not available")
	}

	cb.BillingStorageLock.RLock()
	defer cb.BillingStorageLock.RUnlock()

	return getStoredAttributionDataLocked(ctx, view, localPathPrefix, month, attributionMetricName)
}

// StoreAttributionData stores attribution data for the given metric. It acquires the billing
// storage lock internally; internal callers that already hold the lock should use storeAttributionDataLocked.
func (c *Core) StoreAttributionData(ctx context.Context, localPathPrefix string, month time.Time, attributionMetricName string, data *logical.MetricTypeAttribution) error {
	c.consumptionBillingLock.RLock()
	cb := c.consumptionBilling
	c.consumptionBillingLock.RUnlock()

	if cb == nil {
		return ErrConsumptionBillingNotInitialized
	}

	view, ok := c.GetBillingSubView()
	if !ok {
		return errors.New("billing subview not available")
	}

	cb.BillingStorageLock.Lock()
	defer cb.BillingStorageLock.Unlock()

	return storeAttributionDataLocked(ctx, view, localPathPrefix, month, attributionMetricName, data)
}

// StoreCertAttribution stores mount/namespace attribution data for the given certificate metric
// (e.g. PKI, SSH cert, SSH OTP) by merging the incoming mount-level deltas into any
// pre-existing MetricTypeAttribution for that metric stored in the current month.
//
// countDelta is the cluster-wide increment being flushed (e.g. inc.PkiDurationAdjustedCerts).
// It is added to the stored total and written as MetricTypeAttribution.Count.
// incomingMounts are the per-mount deltas from the current batch (keyed by mount accessor).
func (c *Core) StoreCertAttribution(ctx context.Context, metricName string, countDelta float64, incomingMounts map[string]logical.MountAttribution, currentMonth time.Time) error {
	c.consumptionBillingLock.RLock()
	cb := c.consumptionBilling
	c.consumptionBillingLock.RUnlock()

	if cb == nil {
		return ErrConsumptionBillingNotInitialized
	}

	view, ok := c.GetBillingSubView()
	if !ok {
		return errors.New("billing subview not available")
	}

	cb.BillingStorageLock.Lock()
	defer cb.BillingStorageLock.Unlock()

	return storeCertAttributionLocked(ctx, view, billing.LocalPrefix, metricName, countDelta, incomingMounts, currentMonth)
}

// storeCertAttributionLocked merges incomingMounts into the existing MetricTypeAttribution
// for metricName, then writes the result. Must be called with BillingStorageLock held.
func storeCertAttributionLocked(ctx context.Context, view logical.Storage, localPathPrefix string, metricName string, countDelta float64, incomingMounts map[string]logical.MountAttribution, currentMonth time.Time) error {
	existing, err := getStoredAttributionDataLocked(ctx, view, localPathPrefix, currentMonth, metricName)
	if err != nil {
		return fmt.Errorf("failed to read existing attribution for %s: %w", metricName, err)
	}

	if existing.Mounts == nil {
		existing.Mounts = make(map[string]logical.MountAttribution)
	}

	// Merge per-mount deltas from the incoming batch into the stored per-mount totals.
	// Always take metadata (path, namespace, type, UUID) from the incoming entry —
	// it reflects the mount's current state. Only the count is accumulated from
	// storage so that totals are not lost across flushes.
	for accessor, attr := range incomingMounts {
		if prev, ok := existing.Mounts[accessor]; ok {
			attr.Count = toFloat64(prev.Count) + toFloat64(attr.Count)
		}
		existing.Mounts[accessor] = attr
	}

	// Accumulate the cluster-wide total and stamp with the worker-run time so
	// all metrics updated in the same flush cycle share the same timestamp.
	existing.Count = toFloat64(existing.Count) + countDelta
	existing.LastUpdated = currentMonth

	return storeAttributionDataLocked(ctx, view, localPathPrefix, currentMonth, metricName, existing)
}

// toFloat64 converts an interface{} count value to float64.
// Count fields are stored as float64 in memory but may be deserialised as
// json.Number after a storage round-trip, so all cases are handled here.
func toFloat64(v interface{}) float64 {
	switch n := v.(type) {
	case float64:
		return n
	case float32:
		return float64(n)
	case int:
		return float64(n)
	case int64:
		return float64(n)
	case uint:
		return float64(n)
	case uint64:
		return float64(n)
	case interface{ Float64() (float64, error) }:
		f, _ := n.Float64()
		return f
	}
	return 0
}

func (c *Core) UpdateMountAttribution(ctx context.Context, tracker *billing.AttributionTracker, mountTypePrefix string, currentMonth time.Time) error {
	c.consumptionBillingLock.RLock()
	cb := c.consumptionBilling
	c.consumptionBillingLock.RUnlock()
	if cb == nil {
		return ErrConsumptionBillingNotInitialized
	}

	view, ok := c.GetBillingSubView()
	if !ok {
		return errors.New("error updating mount attribution: billing subview not available")
	}

	cb.BillingStorageLock.Lock()
	defer cb.BillingStorageLock.Unlock()

	// Retrieve the attributions already stored for this month.
	stored, err := getStoredAttributionDataLocked(ctx, view, billing.LocalPrefix, currentMonth, mountTypePrefix)
	if err != nil {
		return err
	}
	// Initialize the mounts map if nil (e.g. first flush of the month).
	if stored.Mounts == nil {
		stored.Mounts = make(map[string]logical.MountAttribution)
	}

	// Swap in-memory attributions into stored, then clear the in-memory map.
	// Always take metadata (path, namespace, type, UUID) from the in-memory
	// entry — it reflects the mount's current state. Only the count is
	// accumulated from storage so that totals are not lost across flushes.
	tracker.MountAttributionLock.Lock()
	for mountAccessor, inMem := range tracker.MountAttribution {
		if existing, ok := stored.Mounts[mountAccessor]; ok {
			inMem.Count = toFloat64(existing.Count) + toFloat64(inMem.Count)
		}
		stored.Mounts[mountAccessor] = inMem
		delete(tracker.MountAttribution, mountAccessor)
	}
	tracker.MountAttributionLock.Unlock()

	// Recompute the top-level total count from the per-mount breakdown.
	var total float64
	for _, m := range stored.Mounts {
		total += toFloat64(m.Count)
	}
	stored.Count = total
	stored.LastUpdated = currentMonth

	return storeAttributionDataLocked(ctx, view, billing.LocalPrefix, currentMonth, mountTypePrefix, stored)
}

func (c *Core) UpdateTransitAttribution(ctx context.Context, currentMonth time.Time) error {
	c.consumptionBillingLock.RLock()
	cb := c.consumptionBilling
	c.consumptionBillingLock.RUnlock()
	if cb == nil {
		return ErrConsumptionBillingNotInitialized
	}
	return c.UpdateMountAttribution(ctx, &cb.SecretEngineCounts.Transit.AttributionTracker, billing.TransitDataProtectionCallCountsPrefix, currentMonth)
}

func (c *Core) UpdateTransformAttribution(ctx context.Context, currentMonth time.Time) error {
	c.consumptionBillingLock.RLock()
	cb := c.consumptionBilling
	c.consumptionBillingLock.RUnlock()
	if cb == nil {
		return ErrConsumptionBillingNotInitialized
	}
	return c.UpdateMountAttribution(ctx, &cb.SecretEngineCounts.Transform.AttributionTracker, billing.TransformDataProtectionCallCountsPrefix, currentMonth)
}

func (c *Core) UpdateGcpKmsAttribution(ctx context.Context, currentMonth time.Time) error {
	c.consumptionBillingLock.RLock()
	cb := c.consumptionBilling
	c.consumptionBillingLock.RUnlock()
	if cb == nil {
		return ErrConsumptionBillingNotInitialized
	}
	return c.UpdateMountAttribution(ctx, &cb.SecretEngineCounts.GcpKms.AttributionTracker, billing.GcpKmsDataProtectionCallCountsPrefix, currentMonth)
}

func (c *Core) UpdateSpiffeAttribution(ctx context.Context, currentMonth time.Time) error {
	c.consumptionBillingLock.RLock()
	cb := c.consumptionBilling
	c.consumptionBillingLock.RUnlock()
	if cb == nil {
		return ErrConsumptionBillingNotInitialized
	}
	return c.UpdateMountAttribution(ctx, &cb.SecretEngineCounts.Spiffe.AttributionTracker, billing.SpiffeJwtNormalizedTokenUnits, currentMonth)
}
