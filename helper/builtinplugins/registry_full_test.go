// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package builtinplugins

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_newRegistry(t *testing.T) {
	actual := newRegistry()
	expMinimal := newMinimalRegistry()
	expFullAddon := newFullAddonRegistry()

	require.Equal(t, len(expMinimal.credentialBackends)+len(expFullAddon.credentialBackends), len(actual.credentialBackends),
		"newRegistry() total auth backends mismatch total of minimal and full addon registries")
	require.Equal(t, len(expMinimal.databasePlugins)+len(expFullAddon.databasePlugins), len(actual.databasePlugins),
		"newRegistry() total database plugins mismatch total of common and full addon registries")
	require.Equal(t, len(expMinimal.logicalBackends)+len(expFullAddon.logicalBackends), len(actual.logicalBackends),
		"newRegistry() total logical backends mismatch total of common and full addon registries")

	assertRegistrySubset(t, actual, expMinimal, "common")
	assertRegistrySubset(t, actual, expFullAddon, "full addon")
}

func assertRegistrySubset(t *testing.T, r, subset *registry, subsetName string) {
	t.Helper()

	for k := range subset.credentialBackends {
		require.Contains(t, r.credentialBackends, k, fmt.Sprintf("expected to contain %s auth backend", subsetName))
	}

	for k := range subset.databasePlugins {
		require.Contains(t, r.databasePlugins, k, fmt.Sprintf("expected to contain %s database plugin", subsetName))
	}

	for k := range subset.logicalBackends {
		require.Contains(t, r.logicalBackends, k, fmt.Sprintf("expected to contain %s logical backend", subsetName))
	}
}
