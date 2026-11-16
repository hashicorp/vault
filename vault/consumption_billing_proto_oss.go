// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: MPL-2.0

//go:build !enterprise

package vault

import (
	"context"
	"errors"
	"time"
)

// errSecureHubClientNotReady is the CE stub of the enterprise sentinel. It is
// never returned by the CE sendBillingMetrics, but must be defined so that
// consumption_billing.go (a shared file) compiles in both editions.
var errSecureHubClientNotReady = errors.New("secure hub client is not ready")

func (c *Core) sendBillingMetrics(ctx context.Context, currentTime time.Time) error {
	return nil
}
