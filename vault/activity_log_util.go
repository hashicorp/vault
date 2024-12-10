// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import (
	"context"

	"github.com/hashicorp/vault/vault/activity"
)

// sendCurrentFragment is a no-op on OSS
func (a *ActivityLog) sendCurrentFragment(ctx context.Context) error {
	return nil
}

// receiveSecondaryPreviousMonthGlobalData is a no-op on OSS
func (a *ActivityLog) receiveSecondaryPreviousMonthGlobalData(ctx context.Context, month int64, clients *activity.LogFragment) error {
	return nil
}

// sendPreviousMonthGlobalClientsWorker is a no-op on OSS
func (a *ActivityLog) sendPreviousMonthGlobalClientsWorker(ctx context.Context) error {
	return nil
}
