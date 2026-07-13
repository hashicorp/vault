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

// Must be called with BillingStorageLock held
func (c *Core) storeAttributionDataLocked(ctx context.Context, localPathPrefix string, month time.Time, attributionMetricName string, data *logical.MetricTypeAttribution) error {
	billingPath := billing.GetAttributionMaxPath(localPathPrefix, month, attributionMetricName)

	entry, err := logical.StorageEntryJSON(billingPath, data)
	if err != nil {
		return fmt.Errorf("failed to create storage entry for attribution data: %w", err)
	}

	view, ok := c.GetBillingSubView()
	if !ok {
		return nil
	}

	return view.Put(ctx, entry)
}

// Must be called with BillingStorageLock held
func (c *Core) getStoredAttributionDataLocked(ctx context.Context, localPathPrefix string, month time.Time, attributionMetricName string) (*logical.MetricTypeAttribution, error) {
	billingPath := billing.GetAttributionMaxPath(localPathPrefix, month, attributionMetricName)

	view, ok := c.GetBillingSubView()
	if !ok {
		return &logical.MetricTypeAttribution{}, errors.New("error reading attribution data: billing subview not available")
	}

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

	cb.BillingStorageLock.RLock()
	defer cb.BillingStorageLock.RUnlock()

	return c.getStoredAttributionDataLocked(ctx, localPathPrefix, month, attributionMetricName)
}
