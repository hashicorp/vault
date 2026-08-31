// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: MPL-2.0

//go:build !enterprise

package vault

import (
	"context"
	"time"
)

func (c *Core) sendBillingMetrics(ctx context.Context, currentTime time.Time) error {
	return nil
}
