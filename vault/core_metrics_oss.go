// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import (
	"context"

	"github.com/hashicorp/vault/sdk/logical"
)

// IsKMIPEnabled is a stub for OSS. KMIP is an enterprise feature.
func (c *Core) IsKMIPEnabled(ctx context.Context) (bool, error) {
	return false, nil
}

// CollectKmipMounts is a stub for OSS. KMIP is an enterprise feature.
func (c *Core) CollectKmipMounts() (map[string]logical.MountAttribution, error) {
	return nil, nil
}
