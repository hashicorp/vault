// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package builtinplugins

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

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
