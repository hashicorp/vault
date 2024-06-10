// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package builtinplugins

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_newRegistry(t *testing.T) {
	actual := newRegistry()
	expMinimal := newMinimalRegistry()
	expFullAddon := newFullAddonRegistry()

	require.Equal(t, len(expMinimal.credentialBackends)+len(expFullAddon.credentialBackends), len(actual.credentialBackends),
		"newRegistry() total auth backends mismatch total of common and full addon registries")
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
		if !assert.Contains(t, r.credentialBackends, k) {
			t.Fatalf("missing %s auth backend=%v, newRegistry()=%v", subsetName, k, r.credentialBackends)
		}
	}

	for k := range subset.databasePlugins {
		if !assert.Contains(t, r.databasePlugins, k) {
			t.Fatalf("missing %s database plugin=%v, newRegistry()=%v", subsetName, k, r.databasePlugins)
		}
	}

	for k := range subset.logicalBackends {
		if !assert.Contains(t, r.logicalBackends, k) {
			t.Fatalf("missing %s logical backend=%v, newRegistry()=%v", subsetName, k, r.logicalBackends)
		}
	}
}
