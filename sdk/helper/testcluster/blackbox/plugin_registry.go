// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package blackbox

import (
	"fmt"

	"github.com/hashicorp/vault/sdk/logical"
)

// BuiltinPluginRegistry manages built-in plugin factories
type BuiltinPluginRegistry struct {
	AuthFactories     map[string]logical.Factory
	SecretsFactories  map[string]logical.Factory
	DatabaseFactories map[string]logical.Factory
}

// NewBuiltinPluginRegistry creates a new built-in plugin registry
func NewBuiltinPluginRegistry() *BuiltinPluginRegistry {
	return &BuiltinPluginRegistry{
		AuthFactories:     make(map[string]logical.Factory),
		SecretsFactories:  make(map[string]logical.Factory),
		DatabaseFactories: make(map[string]logical.Factory),
	}
}

// RegisterAuthPlugin registers a built-in auth plugin factory
func (r *BuiltinPluginRegistry) RegisterAuthPlugin(name string, factory logical.Factory) {
	r.AuthFactories[name] = factory
}

// RegisterSecretsPlugin registers a built-in secrets plugin factory
func (r *BuiltinPluginRegistry) RegisterSecretsPlugin(name string, factory logical.Factory) {
	r.SecretsFactories[name] = factory
}

// RegisterDatabasePlugin registers a built-in database plugin factory
func (r *BuiltinPluginRegistry) RegisterDatabasePlugin(name string, factory logical.Factory) {
	r.DatabaseFactories[name] = factory
}

// GetBuiltinPluginFactory retrieves a built-in plugin factory by type and name
func (r *BuiltinPluginRegistry) GetBuiltinPluginFactory(pluginType, pluginName string) (logical.Factory, error) {
	switch pluginType {
	case "auth":
		if factory, exists := r.AuthFactories[pluginName]; exists {
			return factory, nil
		}
		return nil, fmt.Errorf("auth plugin %s not found in registry", pluginName)
	case "secret":
		if factory, exists := r.SecretsFactories[pluginName]; exists {
			return factory, nil
		}
		return nil, fmt.Errorf("secrets plugin %s not found in registry", pluginName)
	case "database":
		if factory, exists := r.DatabaseFactories[pluginName]; exists {
			return factory, nil
		}
		return nil, fmt.Errorf("database plugin %s not found in registry", pluginName)
	default:
		return nil, fmt.Errorf("unsupported plugin type: %s", pluginType)
	}
}

// ListBuiltinPlugins lists all built-in plugins of a given type
func (r *BuiltinPluginRegistry) ListBuiltinPlugins(pluginType string) []string {
	var plugins []string

	switch pluginType {
	case "auth":
		for name := range r.AuthFactories {
			plugins = append(plugins, name)
		}
	case "secret":
		for name := range r.SecretsFactories {
			plugins = append(plugins, name)
		}
	case "database":
		for name := range r.DatabaseFactories {
			plugins = append(plugins, name)
		}
	}

	return plugins
}

// GetAllBuiltinPlugins returns all registered built-in plugins
func (r *BuiltinPluginRegistry) GetAllBuiltinPlugins() map[string][]string {
	return map[string][]string{
		"auth":     r.ListBuiltinPlugins("auth"),
		"secret":   r.ListBuiltinPlugins("secret"),
		"database": r.ListBuiltinPlugins("database"),
	}
}

// DefaultBuiltinPluginRegistry returns a registry with common built-in plugins
func DefaultBuiltinPluginRegistry() *BuiltinPluginRegistry {
	registry := NewBuiltinPluginRegistry()

	// Teams can register actual plugin factories here using the provided methods:
	// registry.RegisterSecretsPlugin("kv", kvFactory)
	// registry.RegisterSecretsPlugin("transit", transitFactory)
	// registry.RegisterAuthPlugin("userpass", userpassFactory)
	// registry.RegisterDatabasePlugin("postgres", postgresFactory)

	return registry
}

// ============================================================================
// Global Registry and Convenience Functions
// ============================================================================

// Global registry instance for convenience
var DefaultRegistry = DefaultBuiltinPluginRegistry()

// Convenience functions for global registry access
func GetBuiltinPluginFactory(pluginType, pluginName string) (logical.Factory, error) {
	return DefaultRegistry.GetBuiltinPluginFactory(pluginType, pluginName)
}

func ListBuiltinPlugins(pluginType string) []string {
	return DefaultRegistry.ListBuiltinPlugins(pluginType)
}

func RegisterBuiltinAuthPlugin(name string, factory logical.Factory) {
	DefaultRegistry.RegisterAuthPlugin(name, factory)
}

func RegisterBuiltinSecretsPlugin(name string, factory logical.Factory) {
	DefaultRegistry.RegisterSecretsPlugin(name, factory)
}

func RegisterBuiltinDatabasePlugin(name string, factory logical.Factory) {
	DefaultRegistry.RegisterDatabasePlugin(name, factory)
}
