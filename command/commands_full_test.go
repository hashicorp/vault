// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"maps"
	"testing"

	"github.com/stretchr/testify/require"
)

// Test_extendFullAddonCommands tests the extendFullAddonCommands function.
func Test_extendFullAddonCommands(t *testing.T) {
	expMinPhysicalBackends := maps.Clone(physicalBackends)
	expMinLoginHandlers := maps.Clone(loginHandlers)

	expAddonPhysicalBackends, expAddonLogicalHandlers := newFullAddonCommands()

	extendFullAddonCommands()

	require.Equal(t, len(expMinPhysicalBackends)+len(expAddonPhysicalBackends), len(physicalBackends),
		"extended physical backends mismatch total of minimal and full addon physical backends")
	require.Equal(t, len(expMinLoginHandlers)+len(expAddonLogicalHandlers), len(loginHandlers),
		"extended logical backends mismatch total of minimal and full addon logical backends")

	for k := range expMinPhysicalBackends {
		require.Contains(t, physicalBackends, k, "expected to contain minimal physical backend")
	}

	for k := range expAddonPhysicalBackends {
		require.Contains(t, physicalBackends, k, "expected to contain full addon physical backend")
	}

	for k := range expMinLoginHandlers {
		require.Contains(t, loginHandlers, k, "expected to contain minimal login handler")
	}

	for k := range expAddonLogicalHandlers {
		require.Contains(t, loginHandlers, k, "expected to contain full addon login handler")
	}
}
