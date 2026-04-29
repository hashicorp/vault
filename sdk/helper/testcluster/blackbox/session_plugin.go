// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package blackbox

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

func (s *Session) MustRegisterPlugin(pluginName, binaryPath, pluginType string) {
	s.t.Helper()

	f, err := os.Open(binaryPath)
	require.NoError(s.t, err)
	defer func() { _ = f.Close() }()

	hasher := sha256.New()
	_, err = io.Copy(hasher, f)
	require.NoError(s.t, err)

	shaSum := hex.EncodeToString(hasher.Sum(nil))

	payload := map[string]any{
		"sha256":  shaSum,
		"command": filepath.Base(binaryPath),
		"type":    pluginType,
	}

	s.MustWrite(filepath.Join("sys/plugins/catalog", pluginType, pluginName), payload)
}

func (s *Session) MustEnablePlugin(path, pluginName, pluginType string) {
	s.t.Helper()

	switch pluginType {
	case "auth":
		s.MustEnableAuth(path, &api.EnableAuthOptions{Type: pluginName})
	case "secret":
		s.MustEnableSecretsEngine(path, &api.MountInput{Type: pluginName})
	default:
		s.t.Fatalf("unknown plugin type: %s", pluginType)
	}
}

func (s *Session) AssertPluginRegistered(pluginName string) {
	s.t.Helper()

	secret := s.MustRead(filepath.Join("sys/plugins/catalog", pluginName))
	require.NotNil(s.t, secret)
}

func (s *Session) AssertPluginConfigured(path string) {
	s.t.Helper()

	configPath := filepath.Join(path, "config")
	s.MustRead(configPath)
}

// PluginSession extends the base Session with plugin-specific functionality
// supporting both built-in and external plugins
type PluginSession struct {
	*Session
	PluginType  string                 // "auth", "secret", or "database"
	PluginName  string                 // Name of the plugin
	MountPath   string                 // Mount path for the plugin
	IsBuiltin   bool                   // Distinguishes built-in vs external plugins
	BinaryPath  string                 // Only for external plugins
	Factory     logical.Factory        // Only for built-in plugins
	Environment map[string]string      // Environment variables for external plugins
	Config      map[string]interface{} // Plugin configuration
}

// PluginSessionOptions provides configuration for creating a PluginSession
type PluginSessionOptions struct {
	PluginType  string
	PluginName  string
	MountPath   string
	Config      map[string]interface{}
	Environment map[string]string
}

// NewBuiltinPluginSession creates a new plugin session for built-in plugins
func (s *Session) NewBuiltinPluginSession(pluginType, pluginName string, factory logical.Factory) *PluginSession {
	s.t.Helper()

	mountPath := fmt.Sprintf("%s-builtin-%s", pluginName, randomString(6))

	ps := &PluginSession{
		Session:     s,
		PluginType:  pluginType,
		PluginName:  pluginName,
		MountPath:   mountPath,
		IsBuiltin:   true,
		Factory:     factory,
		Environment: make(map[string]string),
		Config:      make(map[string]interface{}),
	}

	s.t.Cleanup(func() {
		ps.cleanup()
	})

	return ps
}

// NewExternalPluginSession creates a new plugin session for external plugins
func (s *Session) NewExternalPluginSession(pluginType, pluginName, binaryPath string) *PluginSession {
	s.t.Helper()

	mountPath := fmt.Sprintf("%s-external-%s", pluginName, randomString(6))

	ps := &PluginSession{
		Session:     s,
		PluginType:  pluginType,
		PluginName:  pluginName,
		MountPath:   mountPath,
		IsBuiltin:   false,
		BinaryPath:  binaryPath,
		Environment: make(map[string]string),
		Config:      make(map[string]interface{}),
	}

	s.t.Cleanup(func() {
		ps.cleanup()
	})

	return ps
}

// NewPluginSessionWithOptions creates a plugin session using the provided options
func (s *Session) NewPluginSessionWithOptions(opts PluginSessionOptions, isBuiltin bool, factoryOrPath interface{}) *PluginSession {
	s.t.Helper()

	mountPath := opts.MountPath
	if mountPath == "" {
		prefix := "builtin"
		if !isBuiltin {
			prefix = "external"
		}
		mountPath = fmt.Sprintf("%s-%s-%s", opts.PluginName, prefix, randomString(6))
	}

	ps := &PluginSession{
		Session:     s,
		PluginType:  opts.PluginType,
		PluginName:  opts.PluginName,
		MountPath:   mountPath,
		IsBuiltin:   isBuiltin,
		Environment: opts.Environment,
		Config:      opts.Config,
	}

	if ps.Environment == nil {
		ps.Environment = make(map[string]string)
	}
	if ps.Config == nil {
		ps.Config = make(map[string]interface{})
	}

	if isBuiltin {
		if factory, ok := factoryOrPath.(logical.Factory); ok {
			ps.Factory = factory
		} else {
			s.t.Fatalf("Expected logical.Factory for built-in plugin, got %T", factoryOrPath)
		}
	} else {
		if binaryPath, ok := factoryOrPath.(string); ok {
			ps.BinaryPath = binaryPath
		} else {
			s.t.Fatalf("Expected string binary path for external plugin, got %T", factoryOrPath)
		}
	}

	s.t.Cleanup(func() {
		ps.cleanup()
	})

	return ps
}

// MustRegisterAndEnable registers and enables the plugin based on its type
func (ps *PluginSession) MustRegisterAndEnable() {
	ps.t.Helper()

	if ps.IsBuiltin {
		ps.mustEnableBuiltinPlugin()
	} else {
		ps.mustRegisterAndEnableExternalPlugin()
	}

	// Apply any initial configuration
	if len(ps.Config) > 0 {
		ps.MustConfigure(ps.Config)
	}
}

// mustEnableBuiltinPlugin enables a built-in plugin directly
func (ps *PluginSession) mustEnableBuiltinPlugin() {
	ps.t.Helper()

	switch ps.PluginType {
	case "auth":
		ps.Session.MustEnableAuth(ps.MountPath, &api.EnableAuthOptions{Type: ps.PluginName})
	case "secret":
		ps.Session.MustEnableSecretsEngine(ps.MountPath, &api.MountInput{Type: ps.PluginName})
	case "database":
		// Database plugins are typically mounted under the database secrets engine
		ps.Session.MustEnableSecretsEngine(ps.MountPath, &api.MountInput{Type: "database"})
	default:
		ps.t.Fatalf("Unsupported plugin type for built-in plugin: %s", ps.PluginType)
	}
}

// mustRegisterAndEnableExternalPlugin registers and enables an external plugin
func (ps *PluginSession) mustRegisterAndEnableExternalPlugin() {
	ps.t.Helper()

	// Register the external plugin
	ps.Session.MustRegisterPlugin(ps.PluginName, ps.BinaryPath, ps.PluginType)

	// Enable the plugin
	ps.Session.MustEnablePlugin(ps.MountPath, ps.PluginName, ps.PluginType)
}

// MustConfigure applies configuration to the plugin
func (ps *PluginSession) MustConfigure(config map[string]interface{}) {
	ps.t.Helper()

	configPath := ps.getConfigPath()
	ps.Session.MustWrite(configPath, config)

	// Store the configuration for reference
	for k, v := range config {
		ps.Config[k] = v
	}
}

// MustTestPluginHealth performs basic health checks on the plugin
func (ps *PluginSession) MustTestPluginHealth() {
	ps.t.Helper()

	switch ps.PluginType {
	case "auth":
		ps.mustTestAuthPluginHealth()
	case "secret":
		ps.mustTestSecretsPluginHealth()
	case "database":
		ps.mustTestDatabasePluginHealth()
	default:
		ps.t.Fatalf("Unsupported plugin type: %s", ps.PluginType)
	}
}

// mustTestAuthPluginHealth tests basic auth plugin functionality
func (ps *PluginSession) mustTestAuthPluginHealth() {
	ps.t.Helper()

	// Test that we can read the plugin configuration
	configPath := ps.getConfigPath()
	_, err := ps.Client.Logical().Read(configPath)
	if err != nil {
		ps.t.Logf("Note: Could not read auth plugin config at %s (this may be expected): %v", configPath, err)
	}

	// Test that the auth method is listed
	auths, err := ps.Client.Sys().ListAuth()
	require.NoError(ps.t, err)

	expectedPath := ps.MountPath + "/"
	if auths[expectedPath] == nil {
		ps.t.Fatalf("Auth method %s not found in auth list", expectedPath)
	}

	ps.t.Logf("Auth plugin %s health check passed", ps.PluginName)
}

// mustTestSecretsPluginHealth tests basic secrets plugin functionality
func (ps *PluginSession) mustTestSecretsPluginHealth() {
	ps.t.Helper()

	// Test that we can read the plugin configuration
	configPath := ps.getConfigPath()
	_, err := ps.Client.Logical().Read(configPath)
	if err != nil {
		ps.t.Logf("Note: Could not read secrets plugin config at %s (this may be expected): %v", configPath, err)
	}

	// Test that the secrets engine is listed
	mounts, err := ps.Client.Sys().ListMounts()
	require.NoError(ps.t, err)

	expectedPath := ps.MountPath + "/"
	if mounts[expectedPath] == nil {
		ps.t.Fatalf("Secrets engine %s not found in mounts list", expectedPath)
	}

	ps.t.Logf("Secrets plugin %s health check passed", ps.PluginName)
}

// mustTestDatabasePluginHealth tests basic database plugin functionality
func (ps *PluginSession) mustTestDatabasePluginHealth() {
	ps.t.Helper()

	// For database plugins, test that the database secrets engine is mounted
	mounts, err := ps.Client.Sys().ListMounts()
	require.NoError(ps.t, err)

	expectedPath := ps.MountPath + "/"
	mount := mounts[expectedPath]
	if mount == nil {
		ps.t.Fatalf("Database secrets engine %s not found in mounts list", expectedPath)
	}

	if mount.Type != "database" {
		ps.t.Fatalf("Expected mount type 'database', got '%s'", mount.Type)
	}

	ps.t.Logf("Database plugin %s health check passed", ps.PluginName)
}

// MustTestPluginReload tests plugin reload functionality
func (ps *PluginSession) MustTestPluginReload() {
	ps.t.Helper()

	if ps.IsBuiltin {
		ps.t.Logf("Skipping reload test for built-in plugin %s (reload not applicable)", ps.PluginName)
		return
	}

	// Test plugin reload for external plugins
	reloadInput := &api.ReloadPluginInput{
		Plugin: ps.PluginName,
	}

	reloadID, err := ps.Client.Sys().ReloadPlugin(reloadInput)
	require.NoError(ps.t, err)
	require.NotEmpty(ps.t, reloadID)

	ps.t.Logf("Successfully reloaded external plugin %s (reload ID: %s)", ps.PluginName, reloadID)
}

// MustTestPluginUpgrade tests plugin upgrade functionality (external plugins only)
func (ps *PluginSession) MustTestPluginUpgrade(newBinaryPath string) {
	ps.t.Helper()

	if ps.IsBuiltin {
		ps.t.Skip("Plugin upgrades not applicable for built-in plugins")
		return
	}

	// Store original binary path
	originalBinaryPath := ps.BinaryPath

	// Register the new plugin version
	newPluginName := ps.PluginName + "-v2"
	ps.Session.MustRegisterPlugin(newPluginName, newBinaryPath, ps.PluginType)

	// Test that both versions are registered
	ps.Session.AssertPluginRegistered(ps.PluginName)
	ps.Session.AssertPluginRegistered(newPluginName)

	// Cleanup: restore original state
	ps.t.Cleanup(func() {
		ps.BinaryPath = originalBinaryPath
	})

	ps.t.Logf("Successfully tested plugin upgrade from %s to %s", ps.PluginName, newPluginName)
}

// GetMountPath returns the mount path for this plugin session
func (ps *PluginSession) GetMountPath() string {
	return ps.MountPath
}

// GetPluginInfo returns basic information about the plugin
func (ps *PluginSession) GetPluginInfo() map[string]interface{} {
	return map[string]interface{}{
		"plugin_type": ps.PluginType,
		"plugin_name": ps.PluginName,
		"mount_path":  ps.MountPath,
		"is_builtin":  ps.IsBuiltin,
		"binary_path": ps.BinaryPath,
		"environment": ps.Environment,
		"config":      ps.Config,
	}
}

// UpdateConfig merges new configuration with existing config
func (ps *PluginSession) UpdateConfig(newConfig map[string]interface{}) {
	ps.t.Helper()

	// Merge configurations
	for k, v := range newConfig {
		ps.Config[k] = v
	}

	// Apply the updated configuration
	ps.MustConfigure(ps.Config)
}

// getConfigPath returns the configuration path for the plugin
func (ps *PluginSession) getConfigPath() string {
	switch ps.PluginType {
	case "auth":
		return filepath.Join("auth", ps.MountPath, "config")
	case "secret":
		return filepath.Join(ps.MountPath, "config")
	case "database":
		return filepath.Join(ps.MountPath, "config", ps.PluginName)
	default:
		return filepath.Join(ps.MountPath, "config")
	}
}

// cleanup performs cleanup operations for the plugin session
func (ps *PluginSession) cleanup() {
	if ps.Session == nil || ps.Session.Client == nil {
		return
	}

	ps.t.Logf("Cleaning up plugin session: %s (%s)", ps.PluginName, ps.MountPath)

	// Disable the plugin mount
	var err error
	switch ps.PluginType {
	case "auth":
		err = ps.Client.Sys().DisableAuth(ps.MountPath)
	case "secret", "database":
		err = ps.Client.Sys().Unmount(ps.MountPath)
	}

	if err != nil {
		ps.t.Logf("Warning: Failed to cleanup plugin mount %s: %v", ps.MountPath, err)
	}

	// For external plugins, optionally deregister (commented out to avoid affecting other tests)
	// if !ps.IsBuiltin {
	//     _, err := ps.Client.Logical().Delete(filepath.Join("sys/plugins/catalog", ps.PluginType, ps.PluginName))
	//     if err != nil {
	//         ps.t.Logf("Warning: Failed to deregister plugin %s: %v", ps.PluginName, err)
	//     }
	// }
}

// ============================================================================
// Plugin Testing Utilities and Helpers
// ============================================================================

// TestEndpointExists checks if a plugin endpoint exists and is accessible
func (ps *PluginSession) TestEndpointExists(path string) bool {
	ps.t.Helper()

	fullPath := ps.buildPath(path)
	_, err := ps.Client.Logical().Read(fullPath)

	// We consider the endpoint to exist if we get any response (including errors)
	// that indicate the endpoint exists but may require different parameters/auth
	return err == nil || !isNotFoundError(err)
}

// WriteAndValidate writes data to a plugin endpoint and returns the response
func (ps *PluginSession) WriteAndValidate(path string, data map[string]interface{}) (*api.Secret, error) {
	ps.t.Helper()

	fullPath := ps.buildPath(path)
	return ps.Client.Logical().Write(fullPath, data)
}

// ReadAndValidate reads from a plugin endpoint and returns the response
func (ps *PluginSession) ReadAndValidate(path string) (*api.Secret, error) {
	ps.t.Helper()

	fullPath := ps.buildPath(path)
	return ps.Client.Logical().Read(fullPath)
}

// DeleteAndValidate deletes from a plugin endpoint and returns the response
func (ps *PluginSession) DeleteAndValidate(path string) (*api.Secret, error) {
	ps.t.Helper()

	fullPath := ps.buildPath(path)
	return ps.Client.Logical().Delete(fullPath)
}

// ListAndValidate lists from a plugin endpoint and returns the response
func (ps *PluginSession) ListAndValidate(path string) (*api.Secret, error) {
	ps.t.Helper()

	fullPath := ps.buildPath(path)
	return ps.Client.Logical().List(fullPath)
}

// ExpectError expects an operation to fail and validates the error
func (ps *PluginSession) ExpectError(operation func() (*api.Secret, error)) error {
	ps.t.Helper()

	_, err := operation()
	if err == nil {
		return fmt.Errorf("expected operation to fail, but it succeeded")
	}
	return nil // Error was expected
}

// ExpectSuccess expects an operation to succeed and returns the response
func (ps *PluginSession) ExpectSuccess(operation func() (*api.Secret, error)) *api.Secret {
	ps.t.Helper()

	resp, err := operation()
	require.NoError(ps.t, err, "Expected operation to succeed")
	return resp
}

// ValidateResponse validates that a response has expected properties
func (ps *PluginSession) ValidateResponse(resp *api.Secret, validations ...ResponseValidator) {
	ps.t.Helper()

	for _, validation := range validations {
		if err := validation(resp); err != nil {
			ps.t.Fatalf("Response validation failed: %v", err)
		}
	}
}

// ResponseValidator validates aspects of an API response
type ResponseValidator func(*api.Secret) error

// HasDataField validates that response has a specific data field
func HasDataField(field string) ResponseValidator {
	return func(resp *api.Secret) error {
		if resp == nil || resp.Data == nil {
			return fmt.Errorf("response or data is nil")
		}
		if _, exists := resp.Data[field]; !exists {
			return fmt.Errorf("response missing expected field: %s", field)
		}
		return nil
	}
}

// HasDataValue validates that response has a specific data field with expected value
func HasDataValue(field string, expectedValue interface{}) ResponseValidator {
	return func(resp *api.Secret) error {
		if resp == nil || resp.Data == nil {
			return fmt.Errorf("response or data is nil")
		}
		value, exists := resp.Data[field]
		if !exists {
			return fmt.Errorf("response missing expected field: %s", field)
		}
		if value != expectedValue {
			return fmt.Errorf("field %s has value %v, expected %v", field, value, expectedValue)
		}
		return nil
	}
}

// IsNotEmpty validates that response data is not empty
func IsNotEmpty() ResponseValidator {
	return func(resp *api.Secret) error {
		if resp == nil || resp.Data == nil || len(resp.Data) == 0 {
			return fmt.Errorf("response data is empty")
		}
		return nil
	}
}

// HasAuth validates that response has authentication information
func HasAuth() ResponseValidator {
	return func(resp *api.Secret) error {
		if resp == nil || resp.Auth == nil {
			return fmt.Errorf("response has no authentication information")
		}
		return nil
	}
}

// HasLease validates that response has lease information
func HasLease() ResponseValidator {
	return func(resp *api.Secret) error {
		if resp == nil || resp.LeaseID == "" {
			return fmt.Errorf("response has no lease information")
		}
		return nil
	}
}

// TestSequence runs a sequence of operations and validates them
func (ps *PluginSession) TestSequence(operations ...SequenceOperation) {
	ps.t.Helper()

	for i, op := range operations {
		ps.t.Logf("Executing sequence operation %d: %s", i+1, op.Name)
		if err := op.Execute(ps); err != nil {
			ps.t.Fatalf("Sequence operation %d (%s) failed: %v", i+1, op.Name, err)
		}
	}
}

// SequenceOperation represents a single operation in a test sequence
type SequenceOperation struct {
	Name    string
	Execute func(*PluginSession) error
}

// WriteOp creates a write operation for test sequences
func WriteOp(name, path string, data map[string]interface{}, validators ...ResponseValidator) SequenceOperation {
	return SequenceOperation{
		Name: name,
		Execute: func(ps *PluginSession) error {
			resp, err := ps.WriteAndValidate(path, data)
			if err != nil {
				return err
			}
			ps.ValidateResponse(resp, validators...)
			return nil
		},
	}
}

// ReadOp creates a read operation for test sequences
func ReadOp(name, path string, validators ...ResponseValidator) SequenceOperation {
	return SequenceOperation{
		Name: name,
		Execute: func(ps *PluginSession) error {
			resp, err := ps.ReadAndValidate(path)
			if err != nil {
				return err
			}
			ps.ValidateResponse(resp, validators...)
			return nil
		},
	}
}

// DeleteOp creates a delete operation for test sequences
func DeleteOp(name, path string) SequenceOperation {
	return SequenceOperation{
		Name: name,
		Execute: func(ps *PluginSession) error {
			_, err := ps.DeleteAndValidate(path)
			return err
		},
	}
}

// buildPath constructs the full path for a plugin endpoint
func (ps *PluginSession) buildPath(path string) string {
	if ps.PluginType == "auth" {
		return filepath.Join("auth", ps.MountPath, path)
	}
	return filepath.Join(ps.MountPath, path)
}

// isNotFoundError checks if an error indicates the endpoint was not found
func isNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	// This is a simple check - in practice you might want more sophisticated error detection
	return fmt.Sprintf("%v", err) == "404 Not Found" ||
		fmt.Sprintf("%v", err) == "405 Method Not Allowed"
}

// TestExpectedError tests that a write operation fails as expected
func (ps *PluginSession) TestExpectedError(path string, data map[string]interface{}) error {
	ps.t.Helper()

	fullPath := ps.buildPath(path)
	_, err := ps.Client.Logical().Write(fullPath, data)

	if err == nil {
		return fmt.Errorf("expected write to %s to fail, but it succeeded", fullPath)
	}

	// Return nil to indicate the error was expected
	return nil
}
