// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

//go:build !enterprise

package vault

import (
	"context"
	"time"
)

// sendPluginCounts is a no-op on OSS
func (c *Core) sendPluginCounts(ctx context.Context) error {
	return nil
}

// perfStandbyPluginCountsWorker is a no-op on OSS
func (c *Core) perfStandbyPluginCountsWorker(ctx context.Context) {
	// No-op: performance standby plugin counts worker is enterprise-only
}

// updateMaxKmseKeyCounts is a no-op on OSS
func (c *Core) updateMaxKmseKeyCounts(ctx context.Context, currentKeyCounts int, localPathPrefix string, currentMonth time.Time) (int, error) {
	return 0, nil
}

// GetStoredHWMKmseCounts is a no-op on OSS
func (c *Core) GetStoredHWMKmseCounts(ctx context.Context, localPathPrefix string, month time.Time) (int, error) {
	return 0, nil
}
