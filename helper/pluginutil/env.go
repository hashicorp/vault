package pluginutil

import (
	"os"

	"github.com/hashicorp/vault/helper/mlock"
)

var (
	// PluginMlockEnabled is the ENV name used to pass the configuration for
	// enabling mlock
	PluginMlockEnabled = "VAULT_PLUGIN_MLOCK_ENABLED"

	// PluginVaultVersionEnv is the ENV name used to pass the version of the
	// vault server to the plugin
	PluginVaultVersionEnv = "VAULT_VERSION"

	// PluginMetadataModeEnv is an ENV name used to disable TLS communication
	// to bootstrap mounting plugins.
	PluginMetadataModeEnv = "VAULT_PLUGIN_METADATA_MODE"
)

// OptionallyEnableMlock determines if mlock should be called, and if so enables
// mlock.
func OptionallyEnableMlock() error {
	if os.Getenv(PluginMlockEnabled) == "true" {
		return mlock.LockMemory()
	}

	return nil
}

// InMetadataMode returns true if the plugin calling this function is running in metadata mode.
func InMetadataMode() bool {
	return os.Getenv(PluginMetadataModeEnv) == "true"
}
