// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestExpiration_IrrevocableLeaseRemovalDisabled verifies that the irrevocable lease removal is disabled on Vault CE
func TestExpiration_IrrevocableLeaseRemovalDisabled(t *testing.T) {
	exp := mockExpiration(t)
	require.Equal(t, false, exp.irrevocableLeaseRemovalEnabled)
}
