// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/physical/inmem"
	"github.com/stretchr/testify/require"
)

// TestMountTableAndNamespacePathRegistry tests that the mount and namespace
// path registry works.
func TestMountTableAndNamespacePathRegistry(t *testing.T) {
	// Make a copy of the registered paths so we can cleanup later otherwise we
	// leave the global state poluted in the process and may break other tests.
	originalPaths := registeredMountOrNamespaceTableKeys
	defer func() {
		registeredMountOrNamespaceTableKeys = originalPaths
	}()

	// Register some dummy paths to test here
	registerMountOrNamespaceTablePaths("mounty1", "mounty2")

	backend, err := inmem.NewInmem(nil, hclog.NewNullLogger())
	require.NoError(t, err)

	imb := backend.(*inmem.InmemBackend)

	applyMountAndNamespaceTableKeys(backend)

	// I wanted to avoid hard coding the actual paths that are registered in init
	// here in an assertion against all paths registered since the whole point of
	// the registry is to avoid hard coding those elsewhere in code! So instead we
	// test that our own dummy registrations exist and ignore anything else.

	registeredPaths := imb.GetMountTablePaths()

	require.Contains(t, registeredPaths, "mounty1")
	require.Contains(t, registeredPaths, "mounty2")
}
