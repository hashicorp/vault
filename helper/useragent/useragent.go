package useragent

import (
	"fmt"
	"runtime"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/version"
)

var (
	// projectURL is the project URL.
	projectURL = "https://www.vaultproject.io/"

	// rt is the runtime - variable for tests.
	rt = runtime.Version()

	// versionFunc is the func that returns the current version. This is a
	// function to take into account the different build processes and distinguish
	// between enterprise and oss builds.
	versionFunc = func() string {
		return version.GetVersion().VersionNumber()
	}
)

// String returns the consistent user-agent string for Vault.
//
// e.g. Vault/0.10.4 (+https://www.vaultproject.io/; go1.10.1)
func String() string {
	return fmt.Sprintf("Vault/%s (+%s; %s)",
		versionFunc(), projectURL, rt)
}

// PluginString is usable by plugins to return a user-agent string reflecting
// the running Vault version and an optional plugin name.
//
// e.g. Vault/0.10.4 (+https://www.vaultproject.io/; azure-auth; go1.10.1)
func PluginString(env *logical.PluginEnvironment, pluginName string) string {
	var name string

	if pluginName != "" {
		name = pluginName + "; "
	}

	return fmt.Sprintf("Vault/%s (+%s; %s%s)",
		env.VaultVersion, projectURL, name, rt)
}
