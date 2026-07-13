// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import "context"

// loadTPMGroups is a no-op for OSS builds
func (i *IdentityStore) loadTPMGroups(ctx context.Context) error {
	return nil
}

// invalidateTPMGroupBucket is a no-op for OSS builds
func (i *IdentityStore) invalidateTPMGroupBucket(ctx context.Context, key string) {}
