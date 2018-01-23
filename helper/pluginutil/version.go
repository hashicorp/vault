package pluginutil

import (
	"os"

	gversion "github.com/hashicorp/go-version"
)

var (
	// PluginVaultVersionEnv is the ENV name used to pass the version of the
	// vault server to the plugin
	PluginVaultVersionEnv = "VAULT_VERSION"
)

// GRPCSupport returns true if Vault can support GRPC transport, false otherwise
func GRPCSupport() bool {
	verString := os.Getenv(PluginVaultVersionEnv)

	if verString != "" && verString != "unknown" {
		ver, err := gversion.NewVersion(verString)
		if err != nil {
			return false
		}

		constraint, err := gversion.NewConstraint(">= 0.9.2")
		if err != nil {
			return false
		}

		return constraint.Check(ver)
	}

	return false
}
