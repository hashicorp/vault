package pluginutil

import (
	"os"

	"github.com/hashicorp/go-version"
)

var (
	// PluginVaultVersionEnv is the ENV name used to pass the version of the
	// vault server to the plugin
	PluginVaultVersionEnv = "VAULT_VERSION"
)

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

		constraint, err := version.NewConstraint(">= 0.9.2")
		if err != nil {
			return true
		}

		return constraint.Check(ver)
	}

	return true
}
