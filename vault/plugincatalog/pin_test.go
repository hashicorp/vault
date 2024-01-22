// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package plugincatalog

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPluginCatalog_PinnedVersionCRUD tests the CRUD operations for pinned
// versions.
func TestPluginCatalog_PinnedVersionCRUD(t *testing.T) {
	catalog := testPluginCatalog(t)

	// Register a plugin in the catalog.
	file, err := os.CreateTemp(catalog.directory, "temp")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	for _, version := range []string{"1.0.0", "2.0.0"} {
		err = catalog.Set(context.Background(), pluginutil.SetPluginInput{
			Name:    "my-plugin",
			Type:    consts.PluginTypeSecrets,
			Version: version,
			Command: filepath.Base(file.Name()),
		})
		require.NoError(t, err)
	}

	// List pinned versions before creating a pin.
	pinnedVersionsBefore, err := catalog.ListPinnedVersions(context.Background())
	require.NoError(t, err)
	assert.Empty(t, pinnedVersionsBefore)

	// Create a pinned version.
	pin := pluginutil.PinnedVersion{
		Name:    "my-plugin",
		Type:    consts.PluginTypeSecrets,
		Version: "1.0.0",
	}
	err = catalog.SetPinnedVersion(context.Background(), &pin)
	require.NoError(t, err)

	// List pinned versions after creating a pin.
	pinnedVersionsAfter, err := catalog.ListPinnedVersions(context.Background())
	require.NoError(t, err)
	require.Len(t, pinnedVersionsAfter, 1)
	assert.Equal(t, pin, *pinnedVersionsAfter[0])

	// Get the pinned version.
	pinnedVersion, err := catalog.GetPinnedVersion(context.Background(), pin.Type, pin.Name)
	require.NoError(t, err)
	assert.Equal(t, pin, *pinnedVersion)

	// Update the pinned version.
	pin.Version = "2.0.0"
	err = catalog.SetPinnedVersion(context.Background(), &pin)
	require.NoError(t, err)

	// Get the updated pinned version.
	pinnedVersion, err = catalog.GetPinnedVersion(context.Background(), pin.Type, pin.Name)
	require.NoError(t, err)
	assert.Equal(t, pin, *pinnedVersion)

	// Update to a version that isn't in the catalog.
	pin.Version = "3.0.0"
	err = catalog.SetPinnedVersion(context.Background(), &pin)
	assert.Error(t, err)

	// Delete the pinned version.
	err = catalog.DeletePinnedVersion(context.Background(), pin.Type, pin.Name)
	require.NoError(t, err)

	// Delete it again, should not error (idempotent).
	err = catalog.DeletePinnedVersion(context.Background(), pin.Type, pin.Name)
	require.NoError(t, err)

	// Verify that the pinned version is deleted.
	pinnedVersion, err = catalog.GetPinnedVersion(context.Background(), pin.Type, pin.Name)
	assert.Equal(t, pluginutil.ErrPinnedVersionNotFound, err)
	assert.Nil(t, pinnedVersion)

	// List should be empty again.
	pinnedVersionsAfterDelete, err := catalog.ListPinnedVersions(context.Background())
	require.NoError(t, err)
	assert.Empty(t, pinnedVersionsAfterDelete)
}
