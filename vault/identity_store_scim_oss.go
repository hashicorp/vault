// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import (
	"context"
)

func (i *IdentityStore) loadSCIMClients(ctx context.Context) error {
	return nil
}

func (i *IdentityStore) invalidateSCIMClient(ctx context.Context, key string) {
}
