// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

//go:build !enterprise

package vault

import (
	"context"
)

// sendPluginCounts is a no-op on OSS
func (c *Core) sendPluginCounts(ctx context.Context) error {
	return nil
}

// perfStandbyPluginCountsWorker is a no-op on OSS
func (c *Core) perfStandbyPluginCountsWorker(ctx context.Context) {
	// No-op: performance standby plugin counts worker is enterprise-only
}
