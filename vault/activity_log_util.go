// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import (
	"context"
	"time"
)

// sendCurrentFragment is a no-op on OSS
func (a *ActivityLog) sendCurrentFragment(ctx context.Context) error {
	return nil
}

// setupClientIDsUsageInfo is a no-op on OSS
func (c *Core) setupClientIDsUsageInfo(ctx context.Context) {
}

// handleClientIDsInMemoryEndOfMonth is a no-op on OSS
func (a *ActivityLog) handleClientIDsInMemoryEndOfMonth(ctx context.Context, currentTime time.Time) {
}
