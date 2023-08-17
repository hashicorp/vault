// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package useragent

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/hashicorp/vault/sdk/logical"
)

var (
	// projectURL is the project URL.
	projectURL = "https://www.vaultproject.io/"

	// rt is the runtime - variable for tests.
	rt = runtime.Version()
)

// String returns the consistent user-agent string for Vault.
// Deprecated: use PluginString instead.
//
// Example output:
//
//	Vault (+https://www.vaultproject.io/; go1.19.5)
//
// Given comments will be appended to the semicolon-delimited comment section:
//
//	Vault (+https://www.vaultproject.io/; go1.19.5; comment-0; comment-1)
//
// At one point the user-agent string returned contained the Vault
// version hardcoded into the vault/sdk/version/ package.  This worked for builtin
// plugins that are compiled into the `vault` binary, in that it correctly described
// the version of that Vault binary.  It did not work for external plugins: for them,
// the version will be based on the version stored in the sdk based on the
// contents of the external plugin's go.mod.  We've kept the String method around
// to avoid breaking builds, but you should be using PluginString.
func String(comments ...string) string {
	c := append([]string{"+" + projectURL, rt}, comments...)
	return fmt.Sprintf("Vault (%s)", strings.Join(c, "; "))
}

// PluginString is usable by plugins to return a user-agent string reflecting
// the running Vault version and an optional plugin name.
//
// e.g. Vault/0.10.4 (+https://www.vaultproject.io/; azure-auth; go1.10.1)
//
// Given comments will be appended to the semicolon-delimited comment section.
//
// e.g. Vault/0.10.4 (+https://www.vaultproject.io/; azure-auth; go1.10.1; comment-0; comment-1)
//
// Returns an empty string if the given env is nil.
func PluginString(env *logical.PluginEnvironment, pluginName string, comments ...string) string {
	if env == nil {
		return ""
	}

	// Construct comments
	c := []string{"+" + projectURL}
	if pluginName != "" {
		c = append(c, pluginName)
	}
	c = append(c, rt)
	c = append(c, comments...)

	// Construct version string
	v := env.VaultVersion
	if env.VaultVersionPrerelease != "" {
		v = fmt.Sprintf("%s-%s", v, env.VaultVersionPrerelease)
	}
	if env.VaultVersionMetadata != "" {
		v = fmt.Sprintf("%s+%s", v, env.VaultVersionMetadata)
	}

	return fmt.Sprintf("Vault/%s (%s)", v, strings.Join(c, "; "))
}
