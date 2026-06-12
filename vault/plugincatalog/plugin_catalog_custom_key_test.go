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
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	// Shared test resources initialized in TestMain
	testCustomPubKey      string
	testCustomKeyFilePath string
	testTempDir           string
)

// TestMain sets up shared resources for all tests and initializes keyrings with a custom key
func TestMain(m *testing.M) {
	// Generate a shared PGP key pair for all tests using the helper
	_, testCustomPubKey, err := generatePGPKeyPair(nil)
	if err != nil {
		panic("failed to generate test PGP key: " + err.Error())
	}

	// Create a temporary directory for the key file
	testTempDir, err = os.MkdirTemp("", "vault-plugin-catalog-test-*")
	if err != nil {
		panic("failed to create temp dir: " + err.Error())
	}

	testCustomKeyFilePath = filepath.Join(testTempDir, "test-custom-key.asc")
	err = os.WriteFile(testCustomKeyFilePath, []byte(testCustomPubKey), 0o644)
	if err != nil {
		os.RemoveAll(testTempDir)
		panic("failed to write custom key file: " + err.Error())
	}

	// Initialize keyrings with the custom key so all tests can verify the singleton behavior
	err = loadWithKey(testCustomKeyFilePath)
	if err != nil {
		os.RemoveAll(testTempDir)
		panic("failed to load custom key: " + err.Error())
	}

	// Run tests
	code := m.Run()

	// Cleanup
	os.RemoveAll(testTempDir)
	os.Exit(code)
}

// setupTestPluginDir creates a temporary plugin directory for a test
func setupTestPluginDir(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	pluginDir := filepath.Join(tmpDir, "plugins")
	err := os.MkdirAll(pluginDir, 0o755)
	require.NoError(t, err)
	return pluginDir
}

// TestPluginCatalog_SetupWithCustomKey tests SetupPluginCatalog with a custom PGP key
func TestPluginCatalog_SetupWithCustomKey(t *testing.T) {
	t.Parallel()

	pluginDir := setupTestPluginDir(t)

	// Setup plugin catalog with custom key file path
	catalog, err := SetupPluginCatalog(context.Background(), &PluginCatalogInput{
		BuiltinRegistry: corehelpers.NewMockBuiltinRegistry(),
		CatalogView:     &logical.InmemStorage{},
		PluginDirectory: pluginDir,
		Logger:          log.NewNullLogger(),
		PluginPGPKey:    testCustomKeyFilePath,
	})

	require.NoError(t, err)
	assert.NotNil(t, catalog)
	assert.Equal(t, testCustomKeyFilePath, catalog.pluginPGPKey, "custom key path should be set")

	// Verify that getVerifyFunc returns a function that uses the custom key
	verifyFunc := catalog.getVerifyFunc()
	assert.NotNil(t, verifyFunc)
}

// TestPluginCatalog_SetupWithRawCustomKey tests SetupPluginCatalog with a raw PGP key (not a file path)
func TestPluginCatalog_SetupWithRawCustomKey(t *testing.T) {
	t.Parallel()

	pluginDir := setupTestPluginDir(t)

	// Setup plugin catalog with raw PGP key (not a file path)
	catalog, err := SetupPluginCatalog(context.Background(), &PluginCatalogInput{
		BuiltinRegistry: corehelpers.NewMockBuiltinRegistry(),
		CatalogView:     &logical.InmemStorage{},
		PluginDirectory: pluginDir,
		Logger:          log.NewNullLogger(),
		PluginPGPKey:    testCustomPubKey, // Pass raw key directly
	})

	require.NoError(t, err)
	assert.NotNil(t, catalog)
	assert.Equal(t, testCustomPubKey, catalog.pluginPGPKey, "raw custom key should be set")

	// Verify that getVerifyFunc returns a function that uses the custom key
	verifyFunc := catalog.getVerifyFunc()
	assert.NotNil(t, verifyFunc)
}

// TestPluginCatalog_KeysNotOverriddenByCustomKey tests that when a custom key
// is configured, it's added to the keyrings rather than replacing the HashiCorp keys.
// Since TestMain initializes the keyrings with a custom key, we can verify the counts here.
func TestPluginCatalog_KeysNotOverriddenByCustomKey(t *testing.T) {
	t.Parallel()

	require.NotNil(t, keyRing2030, "keyRing2030 should not be nil")
	require.NotNil(t, keyRing2026, "keyRing2026 should not be nil")

	assert.Equal(t, 2, keyRing2030.CountEntities(), "keyRing2030 should contain both HashiCorp 2030 key and custom key")
	assert.Equal(t, 2, keyRing2026.CountEntities(), "keyRing2026 should contain both HashiCorp 2026 key and custom key")
}
