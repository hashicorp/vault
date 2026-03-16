// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package blackbox

import (
	"fmt"
	"testing"
)

// PluginConfig represents a plugin configuration with validation
type PluginConfig struct {
	Type        string                     `json:"type"`
	Name        string                     `json:"name"`
	Config      map[string]interface{}     `json:"config"`
	Policies    []string                   `json:"policies,omitempty"`
	Environment map[string]string          `json:"environment,omitempty"`
	Required    []string                   `json:"required,omitempty"` // Required config fields
	Optional    []string                   `json:"optional,omitempty"` // Optional config fields
	Validators  map[string]ConfigValidator `json:"-"`                  // Field validators
}

// ConfigValidator validates a configuration field
type ConfigValidator func(value interface{}) error

// NewPluginConfig creates a new plugin configuration with validation
func NewPluginConfig(pluginType, pluginName string) *PluginConfig {
	return &PluginConfig{
		Type:        pluginType,
		Name:        pluginName,
		Config:      make(map[string]interface{}),
		Environment: make(map[string]string),
		Validators:  make(map[string]ConfigValidator),
	}
}

// SetRequired sets the required configuration fields
func (pc *PluginConfig) SetRequired(fields ...string) *PluginConfig {
	pc.Required = fields
	return pc
}

// SetOptional sets the optional configuration fields
func (pc *PluginConfig) SetOptional(fields ...string) *PluginConfig {
	pc.Optional = fields
	return pc
}

// AddValidator adds a validator for a specific field
func (pc *PluginConfig) AddValidator(field string, validator ConfigValidator) *PluginConfig {
	pc.Validators[field] = validator
	return pc
}

// SetConfig sets a configuration value
func (pc *PluginConfig) SetConfig(key string, value interface{}) *PluginConfig {
	pc.Config[key] = value
	return pc
}

// SetEnvironment sets an environment variable
func (pc *PluginConfig) SetEnvironment(key, value string) *PluginConfig {
	pc.Environment[key] = value
	return pc
}

// ValidateConfig validates the plugin configuration
func (pc *PluginConfig) ValidateConfig() error {
	// Check required fields
	for _, field := range pc.Required {
		if _, exists := pc.Config[field]; !exists {
			return fmt.Errorf("required configuration field '%s' is missing", field)
		}
	}

	// Run field validators
	for field, validator := range pc.Validators {
		if value, exists := pc.Config[field]; exists {
			if err := validator(value); err != nil {
				return fmt.Errorf("validation failed for field '%s': %w", field, err)
			}
		}
	}

	return nil
}

// ============================================================================
// Common Configuration Validators
// ============================================================================

// StringValidator validates that a value is a string
func StringValidator(value interface{}) error {
	if _, ok := value.(string); !ok {
		return fmt.Errorf("expected string, got %T", value)
	}
	return nil
}

// IntValidator validates that a value is an integer
func IntValidator(value interface{}) error {
	switch value.(type) {
	case int, int64, int32:
		return nil
	case float64:
		// JSON numbers are often decoded as float64
		return nil
	default:
		return fmt.Errorf("expected integer, got %T", value)
	}
}

// BoolValidator validates that a value is a boolean
func BoolValidator(value interface{}) error {
	if _, ok := value.(bool); !ok {
		return fmt.Errorf("expected boolean, got %T", value)
	}
	return nil
}

// NonEmptyStringValidator validates that a value is a non-empty string
func NonEmptyStringValidator(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("expected string, got %T", value)
	}
	if str == "" {
		return fmt.Errorf("string cannot be empty")
	}
	return nil
}

// ============================================================================
// Configuration Management Functions
// ============================================================================

// LoadPluginConfig loads a plugin configuration from a file
func LoadPluginConfig(t *testing.T, configPath string) (*PluginConfig, error) {
	t.Helper()

	// Framework teams can implement JSON/YAML/etc. loading here
	return &PluginConfig{
		Type:        "unknown",
		Name:        "unknown",
		Config:      make(map[string]interface{}),
		Environment: make(map[string]string),
		Validators:  make(map[string]ConfigValidator),
	}, nil
}

// ApplyPluginConfig applies a plugin configuration to a session
func ApplyPluginConfig(ps *PluginSession, config *PluginConfig) error {
	// Validate the configuration
	if err := config.ValidateConfig(); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	// Apply configuration
	for key, value := range config.Config {
		ps.Config[key] = value
	}

	// Apply environment variables
	for key, value := range config.Environment {
		ps.Environment[key] = value
	}

	// Apply the configuration to the plugin
	if len(ps.Config) > 0 {
		ps.MustConfigure(ps.Config)
	}

	return nil
}
