// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import (
	"context"
	"time"
)

func (m *ExpirationManager) deleteIrrevocableLease(ctx context.Context, le *leaseEntry) {}

func (m *ExpirationManager) removeIrrevocableLeasesEnabled(removeIrrevocableLeaseAfter time.Duration) bool {
	return false
}
