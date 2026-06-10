// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package plugincatalog

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPluginCatalog_SetupWithCustomKey tests SetupPluginCatalog with a custom PGP key
func TestPluginCatalog_SetupWithCustomKey(t *testing.T) {
	t.Parallel()

	// Generate a test PGP key pair
	_, pubKeyArmored := generatePGPKeyPair(t)

	// Create temporary directories
	tmpDir := t.TempDir()
	pluginDir := filepath.Join(tmpDir, "plugins")
	err := os.MkdirAll(pluginDir, 0o755)
	require.NoError(t, err)

	keyPath := filepath.Join(tmpDir, "custom-key.asc")
	err = os.WriteFile(keyPath, []byte(pubKeyArmored), 0o644)
	require.NoError(t, err)

	// Setup plugin catalog with custom key
	catalog, err := SetupPluginCatalog(context.Background(), &PluginCatalogInput{
		BuiltinRegistry: corehelpers.NewMockBuiltinRegistry(),
		CatalogView:     &logical.InmemStorage{},
		PluginDirectory: pluginDir,
		Logger:          log.NewNullLogger(),
		PluginPGPKey:    keyPath,
	})

	require.NoError(t, err)
	assert.NotNil(t, catalog)
	assert.Equal(t, keyPath, catalog.pluginPGPKey, "custom key path should be set")

	// Verify that getVerifyFunc returns a function that uses the custom key
	verifyFunc := catalog.getVerifyFunc()
	assert.NotNil(t, verifyFunc)
}

// TestPluginCatalog_verifyOfficialPlugins_WithCustomKey tests verifyOfficialPlugins with custom key
func TestPluginCatalog_verifyOfficialPlugins_WithCustomKey(t *testing.T) {
	t.Parallel()

	// Generate a test PGP key pair
	privKey, pubKeyArmored := generatePGPKeyPair(t)

	// Create temporary directories
	tmpDir := t.TempDir()
	pluginDir := filepath.Join(tmpDir, "plugins")
	err := os.MkdirAll(pluginDir, 0o755)
	require.NoError(t, err)

	keyPath := filepath.Join(tmpDir, "custom-key.asc")
	err = os.WriteFile(keyPath, []byte(pubKeyArmored), 0o644)
	require.NoError(t, err)

	// Create a plugin artifact with proper signatures
	pluginName := "vault-plugin-test"
	pluginVersion := "1.0.0"
	artifactDir := filepath.Join(pluginDir, GetExtractedArtifactDir(pluginName, pluginVersion))
	err = os.MkdirAll(artifactDir, 0o755)
	require.NoError(t, err)

	contents := generatePluginArtifactContents(t, pluginName, pluginVersion, consts.PluginTypeSecrets, true, privKey)
	for filename, data := range contents {
		err := os.WriteFile(filepath.Join(artifactDir, filename), data, 0o644)
		require.NoError(t, err)
	}

	storage := &logical.InmemStorage{}
	// Setup plugin catalog with custom key
	catalog, err := SetupPluginCatalog(context.Background(), &PluginCatalogInput{
		BuiltinRegistry: corehelpers.NewMockBuiltinRegistry(),
		CatalogView:     storage,
		PluginDirectory: pluginDir,
		Logger:          log.NewNullLogger(),
		PluginPGPKey:    keyPath,
	})
	require.NoError(t, err)

	// Register the plugin as official
	pluginEntry := &pluginutil.PluginRunner{
		Name:    pluginName,
		Type:    consts.PluginTypeSecrets,
		Version: pluginVersion,
		Command: filepath.Join(GetExtractedArtifactDir(pluginName, pluginVersion), pluginName),
		Builtin: false,
	}

	// Store the plugin in catalog
	entry, err := logical.StorageEntryJSON(pluginEntry.Name, pluginEntry)
	require.NoError(t, err)
	err = storage.Put(context.Background(), entry)
	require.NoError(t, err)

	// Verify official plugins - should succeed with custom key
	err = catalog.verifyOfficialPlugins(context.Background())
	require.NoError(t, err)
}

// TestPluginCatalog_SetupWithRawCustomKey tests SetupPluginCatalog with a raw PGP key (not a file path)
func TestPluginCatalog_SetupWithRawCustomKey(t *testing.T) {
	t.Parallel()

	// Generate a test PGP key pair
	_, pubKeyArmored := generatePGPKeyPair(t)

	// Create temporary directories
	tmpDir := t.TempDir()
	pluginDir := filepath.Join(tmpDir, "plugins")
	err := os.MkdirAll(pluginDir, 0o755)
	require.NoError(t, err)

	// Setup plugin catalog with raw PGP key (not a file path)
	catalog, err := SetupPluginCatalog(context.Background(), &PluginCatalogInput{
		BuiltinRegistry: corehelpers.NewMockBuiltinRegistry(),
		CatalogView:     &logical.InmemStorage{},
		PluginDirectory: pluginDir,
		Logger:          log.NewNullLogger(),
		PluginPGPKey:    pubKeyArmored, // Pass raw key directly
	})

	require.NoError(t, err)
	assert.NotNil(t, catalog)
	assert.Equal(t, pubKeyArmored, catalog.pluginPGPKey, "raw custom key should be set")

	// Verify that getVerifyFunc returns a function that uses the custom key
	verifyFunc := catalog.getVerifyFunc()
	assert.NotNil(t, verifyFunc)
}

// TestPluginCatalog_verifyOfficialPlugins_WithRawCustomKey tests verifyOfficialPlugins with raw custom key
func TestPluginCatalog_verifyOfficialPlugins_WithRawCustomKey(t *testing.T) {
	t.Parallel()

	// Generate a test PGP key pair
	privKey, pubKeyArmored := generatePGPKeyPair(t)

	// Create temporary directories
	tmpDir := t.TempDir()
	pluginDir := filepath.Join(tmpDir, "plugins")
	err := os.MkdirAll(pluginDir, 0o755)
	require.NoError(t, err)

	// Create a plugin artifact with proper signatures
	pluginName := "vault-plugin-test-raw"
	pluginVersion := "1.0.0"
	artifactDir := filepath.Join(pluginDir, GetExtractedArtifactDir(pluginName, pluginVersion))
	err = os.MkdirAll(artifactDir, 0o755)
	require.NoError(t, err)

	contents := generatePluginArtifactContents(t, pluginName, pluginVersion, consts.PluginTypeSecrets, true, privKey)
	for filename, data := range contents {
		err := os.WriteFile(filepath.Join(artifactDir, filename), data, 0o644)
		require.NoError(t, err)
	}

	storage := &logical.InmemStorage{}
	// Setup plugin catalog with raw PGP key (not a file path)
	catalog, err := SetupPluginCatalog(context.Background(), &PluginCatalogInput{
		BuiltinRegistry: corehelpers.NewMockBuiltinRegistry(),
		CatalogView:     storage,
		PluginDirectory: pluginDir,
		Logger:          log.NewNullLogger(),
		PluginPGPKey:    pubKeyArmored, // Pass raw key directly
	})
	require.NoError(t, err)

	// Register the plugin as official
	pluginEntry := &pluginutil.PluginRunner{
		Name:    pluginName,
		Type:    consts.PluginTypeSecrets,
		Version: pluginVersion,
		Command: filepath.Join(GetExtractedArtifactDir(pluginName, pluginVersion), pluginName),
		Builtin: false,
	}

	// Store the plugin in catalog
	entry, err := logical.StorageEntryJSON(pluginEntry.Name, pluginEntry)
	require.NoError(t, err)
	err = storage.Put(context.Background(), entry)
	require.NoError(t, err)

	// Verify official plugins - should succeed with raw custom key
	err = catalog.verifyOfficialPlugins(context.Background())
	require.NoError(t, err)
}
