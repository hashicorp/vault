package useragent

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/version"
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
//
// Given comments will be appended to the semicolon-delimited comment section.
//
// e.g. Vault/0.10.4 (+https://www.vaultproject.io/; go1.10.1; comment-0; comment-1)
func String(comments ...string) string {
	c := append([]string{"+" + projectURL, rt}, comments...)
	return fmt.Sprintf("Vault/%s (%s)", versionFunc(), strings.Join(c, "; "))
}

// PluginString is usable by plugins to return a user-agent string reflecting
// the running Vault version and an optional plugin name.
//
// e.g. Vault/0.10.4 (+https://www.vaultproject.io/; azure-auth; go1.10.1)
//
// Given comments will be appended to the semicolon-delimited comment section.
//
// e.g. Vault/0.10.4 (+https://www.vaultproject.io/; azure-auth; go1.10.1; comment-0; comment-1)
func PluginString(env *logical.PluginEnvironment, pluginName string, comments ...string) string {
	c := []string{"+" + projectURL}
	if pluginName != "" {
		c = append(c, pluginName)
	}
	c = append(c, rt)
	c = append(c, comments...)
	return fmt.Sprintf("Vault/%s (%s)", env.VaultVersion, strings.Join(c, "; "))
}
