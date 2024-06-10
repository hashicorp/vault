// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise && !minimal

package builtinplugins

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Test_newRegistry tests that newRegistry() returns a registry with
// the expected minimal registry extended with full addon registry
func Test_newRegistry(t *testing.T) {
	actual := newRegistry()
	expMinimal := newMinimalRegistry()
	expFullAddon := newFullAddonRegistry()

	require.Equal(t, len(expMinimal.credentialBackends)+len(expFullAddon.credentialBackends), len(actual.credentialBackends),
		"newRegistry() total auth backends mismatch total of minimal and full addon registries")
	require.Equal(t, len(expMinimal.databasePlugins)+len(expFullAddon.databasePlugins), len(actual.databasePlugins),
		"newRegistry() total database plugins mismatch total of minimal and full addon registries")
	require.Equal(t, len(expMinimal.logicalBackends)+len(expFullAddon.logicalBackends), len(actual.logicalBackends),
		"newRegistry() total logical backends mismatch total of minimal and full addon registries")

	assertRegistrySubset(t, actual, expMinimal, "common")
	assertRegistrySubset(t, actual, expFullAddon, "full addon")
}
