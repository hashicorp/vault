// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package blackbox

import (
	"fmt"

	"github.com/hashicorp/vault/sdk/logical"
)

// ExtendedPluginRegistration provides enhanced registration with metadata
func (s *Session) ExtendedPluginRegistration(pluginName, binaryPath, pluginType string, metadata map[string]interface{}) error {
	s.t.Helper()

	// First, validate the binary
	if err := ValidatePluginBinary(s.t, binaryPath); err != nil {
		return fmt.Errorf("plugin binary validation failed: %w", err)
	}

	// Get binary info
	binaryInfo, err := GetPluginBinaryInfo(s.t, binaryPath)
	if err != nil {
		return fmt.Errorf("failed to get plugin binary info: %w", err)
	}

	// Perform standard registration
	s.MustRegisterPlugin(pluginName, binaryPath, pluginType)

	// Store metadata (this would typically be stored in a database or file)
	s.t.Logf("Plugin %s registered with SHA256: %s", pluginName, binaryInfo.SHA256)
	if metadata != nil {
		s.t.Logf("Plugin %s metadata: %+v", pluginName, metadata)
	}

	return nil
}

// BatchPluginRegistration registers multiple plugins at once
func (s *Session) BatchPluginRegistration(plugins []PluginRegistrationRequest) error {
	s.t.Helper()

	for _, plugin := range plugins {
		if err := s.ExtendedPluginRegistration(plugin.Name, plugin.BinaryPath, plugin.Type, plugin.Metadata); err != nil {
			return fmt.Errorf("failed to register plugin %s: %w", plugin.Name, err)
		}
	}

	s.t.Logf("Successfully registered %d plugins", len(plugins))
	return nil
}

// PluginRegistrationRequest represents a plugin registration request
type PluginRegistrationRequest struct {
	Name       string
	BinaryPath string
	Type       string
	Metadata   map[string]interface{}
}

// ============================================================================
// Built-in Plugin Direct Access Utilities
// ============================================================================

// NewBuiltinPluginSessionFromRegistry creates a plugin session from the registry
func (s *Session) NewBuiltinPluginSessionFromRegistry(pluginType, pluginName string) (*PluginSession, error) {
	s.t.Helper()

	factory, err := GetBuiltinPluginFactory(pluginType, pluginName)
	if err != nil {
		return nil, fmt.Errorf("failed to get built-in plugin factory: %w", err)
	}

	return s.NewBuiltinPluginSession(pluginType, pluginName, factory), nil
}

// ============================================================================
// Setup Helper Functions for Built-in Plugins
// ============================================================================

// SetupBuiltinAuthPlugin creates and registers a built-in auth plugin
func SetupBuiltinAuthPlugin(v *Session, pluginName string, factory logical.Factory) *PluginSession {
	v.t.Helper()

	ps := v.NewBuiltinPluginSession("auth", pluginName, factory)
	ps.MustRegisterAndEnable()

	return ps
}

// SetupBuiltinSecretsPlugin creates and registers a built-in secrets plugin
func SetupBuiltinSecretsPlugin(v *Session, pluginName string, factory logical.Factory) *PluginSession {
	v.t.Helper()

	ps := v.NewBuiltinPluginSession("secret", pluginName, factory)
	ps.MustRegisterAndEnable()

	return ps
}

// SetupBuiltinDatabasePlugin creates and registers a built-in database plugin
func SetupBuiltinDatabasePlugin(v *Session, pluginName string, factory logical.Factory) *PluginSession {
	v.t.Helper()

	ps := v.NewBuiltinPluginSession("database", pluginName, factory)
	ps.MustRegisterAndEnable()

	return ps
}

// ============================================================================
// Setup Helper Functions for External Plugins
// ============================================================================

// SetupExternalAuthPlugin creates and registers an external auth plugin
func SetupExternalAuthPlugin(v *Session, pluginName, binaryPath string) *PluginSession {
	v.t.Helper()

	ps := v.NewExternalPluginSession("auth", pluginName, binaryPath)
	ps.MustRegisterAndEnable()

	return ps
}

// SetupExternalSecretsPlugin creates and registers an external secrets plugin
func SetupExternalSecretsPlugin(v *Session, pluginName, binaryPath string) *PluginSession {
	v.t.Helper()

	ps := v.NewExternalPluginSession("secret", pluginName, binaryPath)
	ps.MustRegisterAndEnable()

	return ps
}

// SetupExternalDatabasePlugin creates and registers an external database plugin
func SetupExternalDatabasePlugin(v *Session, pluginName, binaryPath string) *PluginSession {
	v.t.Helper()

	ps := v.NewExternalPluginSession("database", pluginName, binaryPath)
	ps.MustRegisterAndEnable()

	return ps
}
