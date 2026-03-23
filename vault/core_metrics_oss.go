// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import (
	"context"
)

// IsKMIPEnabled is a stub for OSS. KMIP is an enterprise feature.
func (c *Core) IsKMIPEnabled(ctx context.Context) (bool, error) {
	return false, nil
}
