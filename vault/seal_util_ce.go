// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import (
	"context"

	"github.com/hashicorp/vault/sdk/physical"
)

func HasPartiallyWrappedPaths(_ context.Context, _ physical.Backend) (bool, error) {
	return false, nil
}
