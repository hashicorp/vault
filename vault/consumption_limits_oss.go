// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import "context"

// ListBillingLimits is a no-op on OSS
func (c *Core) ListBillingLimits(ctx context.Context) (map[string]float64, error) {
	return nil, nil
}

// GetBillingLimit is a no-op on OSS
func (c *Core) GetBillingLimit(ctx context.Context, metric string) (float64, bool, error) {
	return 0, false, nil
}

// SetBillingLimit is a no-op on OSS
func (c *Core) SetBillingLimit(ctx context.Context, metric string, limit float64) error {
	return nil
}

// DeleteBillingLimit is a no-op on OSS
func (c *Core) DeleteBillingLimit(ctx context.Context, metric string) error {
	return nil
}

// evaluateConsumptionLimits is a no-op on OSS
func evaluateConsumptionLimits(totals map[string]float64, limits map[string]float64) []consumptionLimitState {
	return nil
}
