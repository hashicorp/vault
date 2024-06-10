// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package builtinplugins

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_newRegistry(t *testing.T) {
	actual := newRegistry()
	expCommon := newCommonRegistry()
	expFullAddon := newFullAddonRegistry()

	assert.Equal(t, len(actual.credentialBackends), len(expCommon.credentialBackends)+len(expFullAddon.credentialBackends),
		"newRegistry() total auth backends mismatch total of common and full addon registries")
	assert.Equal(t, len(actual.databasePlugins), len(expCommon.databasePlugins)+len(expFullAddon.databasePlugins),
		"newRegistry() total database plugins mismatch total of common and full addon registries")
	assert.Equal(t, len(actual.logicalBackends), len(expCommon.logicalBackends)+len(expFullAddon.logicalBackends),
		"newRegistry() total logical backends mismatch total of common and full addon registries")

	assertRegistrySubset(t, actual, expCommon, "common")
	assertRegistrySubset(t, actual, expFullAddon, "full addon")
}

func assertRegistrySubset(t *testing.T, r, subset *registry, subsetName string) {
	t.Helper()

	for k := range subset.credentialBackends {
		if !assert.Contains(t, r.credentialBackends, k) {
			t.Errorf("missing %s auth backend=%v, newRegistry()=%v", subsetName, k, r.credentialBackends)
		}
	}

	for k := range subset.databasePlugins {
		if !assert.Contains(t, r.databasePlugins, k) {
			t.Errorf("missing %s database plugin=%v, newRegistry()=%v", subsetName, k, r.databasePlugins)
		}
	}

	for k := range subset.logicalBackends {
		if !assert.Contains(t, r.logicalBackends, k) {
			t.Errorf("missing %s logical backend=%v, newRegistry()=%v", subsetName, k, r.logicalBackends)
		}
	}
}
