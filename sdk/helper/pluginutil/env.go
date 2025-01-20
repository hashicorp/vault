// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package pluginutil

import (
	"os"

	"github.com/hashicorp/go-secure-stdlib/mlock"
	version "github.com/hashicorp/go-version"
)

const (
	// PluginAutoMTLSEnv is used to ensure AutoMTLS is used. This will override
	// setting a TLSProviderFunc for a plugin.
	PluginAutoMTLSEnv = "VAULT_PLUGIN_AUTOMTLS_ENABLED"

	// PluginMlockEnabled is the ENV name used to pass the configuration for
	// enabling mlock
	PluginMlockEnabled = "VAULT_PLUGIN_MLOCK_ENABLED"

	// PluginVaultVersionEnv is the ENV name used to pass the version of the
	// vault server to the plugin
	PluginVaultVersionEnv = "VAULT_VERSION"

	// PluginMetadataModeEnv is an ENV name used to disable TLS communication
	// to bootstrap mounting plugins.
	PluginMetadataModeEnv = "VAULT_PLUGIN_METADATA_MODE"

	// PluginUnwrapTokenEnv is the ENV name used to pass unwrap tokens to the
	// plugin.
	PluginUnwrapTokenEnv = "VAULT_UNWRAP_TOKEN"

	// PluginCACertPEMEnv is an ENV name used for holding a CA PEM-encoded
	// string. Used for testing.
	PluginCACertPEMEnv = "VAULT_TESTING_PLUGIN_CA_PEM"

	// PluginMultiplexingOptOut is an ENV name used to define a comma separated list of plugin names
	// opted-out of the multiplexing feature; for emergencies if multiplexing ever causes issues
	PluginMultiplexingOptOut = "VAULT_PLUGIN_MULTIPLEXING_OPT_OUT"

	// PluginUseLegacyEnvLayering opts out of new environment variable precedence.
	// If set to true, Vault process environment variables take precedence over any
	// colliding plugin-specific environment variables. Otherwise, plugin-specific
	// environment variables take precedence over Vault process environment variables.
	PluginUseLegacyEnvLayering = "VAULT_PLUGIN_USE_LEGACY_ENV_LAYERING"
)

// OptionallyEnableMlock determines if mlock should be called, and if so enables
// mlock.
func OptionallyEnableMlock() error {
	if os.Getenv(PluginMlockEnabled) == "true" {
		return mlock.LockMemory()
	}

	return nil
}

// GRPCSupport defaults to returning true, unless VAULT_VERSION is missing or
// it fails to meet the version constraint.
func GRPCSupport() bool {
	verString := os.Getenv(PluginVaultVersionEnv)
	// If the env var is empty, we fall back to netrpc for backward compatibility.
	if verString == "" {
		return false
	}
	if verString != "unknown" {
		ver, err := version.NewVersion(verString)
		if err != nil {
			return true
		}
		// Due to some regressions on 0.9.2 & 0.9.3 we now require version 0.9.4
		// to allow the plugin framework to default to gRPC.
		constraint, err := version.NewConstraint(">= 0.9.4")
		if err != nil {
			return true
		}
		return constraint.Check(ver)
	}
	return true
}

// InMetadataMode returns true if the plugin calling this function is running in metadata mode.
func InMetadataMode() bool {
	return os.Getenv(PluginMetadataModeEnv) == "true"
}
