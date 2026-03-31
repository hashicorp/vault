// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import (
	"context"

	"github.com/hashicorp/vault/sdk/framework"
)

func (i *IdentityStore) loadSCIMClients(ctx context.Context) error {
	return nil
}

func (i *IdentityStore) invalidateSCIMClient(ctx context.Context, key string) {
}

func (i *IdentityStore) startSCIMDeletingClientCleanup(ctx context.Context, isActive bool) {
}

func (i *IdentityStore) stopSCIMDeletingClientCleanup() {
}

func (i *IdentityStore) enqueueSCIMCleanup(clientID string, namespaceID string) {
}

func scimPaths(_ *IdentityStore) []*framework.Path {
	return []*framework.Path{}
}
