//go:build !enterprise

// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package vault

import (
	"context"
	"time"
)

func (c *Core) UpdateTransformCallCounts(ctx context.Context, currentMonth time.Time) (uint64, error) {
	// No-op in OSS
	return 0, nil
}

func (c *Core) GetStoredTransformCallCounts(ctx context.Context, month time.Time) (uint64, error) {
	return 0, nil
}
