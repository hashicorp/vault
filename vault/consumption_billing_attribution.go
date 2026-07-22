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
	for accessor, attr := range incomingMounts {
		if prev, ok := existing.Mounts[accessor]; ok {
			prev.Count = toFloat64(prev.Count) + toFloat64(attr.Count)
			existing.Mounts[accessor] = prev
		} else {
			existing.Mounts[accessor] = attr
		}
	}

	// Accumulate the cluster-wide total and stamp with the worker-run time so
	// all metrics updated in the same flush cycle share the same timestamp.
	existing.Count = toFloat64(existing.Count) + countDelta
	existing.LastUpdated = currentMonth

	return storeAttributionDataLocked(ctx, view, localPathPrefix, currentMonth, metricName, existing)
}

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
