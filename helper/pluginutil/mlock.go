package pluginutil

import (
	"os"

	"github.com/hashicorp/vault/helper/mlock"
)

var (
	// PluginMlockEnabled is the ENV name used to pass the configuration for
	// enabling mlock
	PluginMlockEnabled = "VAULT_PLUGIN_MLOCK_ENABLED"
)

// OptionallyEnableMlock determines if mlock should be called, and if so enables
// mlock.
func OptionallyEnableMlock() error {
	if os.Getenv(PluginMlockEnabled) == "true" {
		return mlock.LockMemory()
	}

	return nil
}
