// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"maps"
	"slices"
	"testing"

	"github.com/hashicorp/cli"
	"github.com/stretchr/testify/require"
)

// TestLoginHandlers pins the exact set of auth methods that `login` and
// `auth help` understand. vault-client once shipped with no login handlers
// at all, so every `login -method=...` failed with "Unknown auth method".
func TestLoginHandlers(t *testing.T) {
	t.Parallel()

	expected := append([]string{
		"alicloud",
		"aws",
		"cert",
		"cf",
		"gcp",
		"github",
		"kerberos",
		"ldap",
		"oci",
		"oidc",
		"okta",
		"pcf",
		"radius",
		"spiffe",
		"token",
		"tpm",
		"userpass",
	}, entLoginHandlers...)

	require.ElementsMatch(t, expected, slices.Collect(maps.Keys(LoginHandlers())))
}

// TestInitClientCommands_LoginHandlers verifies that the `login` and
// `auth help` commands registered by InitClientCommands carry the full
// LoginHandlers set rather than an empty map.
func TestInitClientCommands_LoginHandlers(t *testing.T) {
	t.Parallel()

	ui := cli.NewMockUi()
	commands := InitClientCommands(ui, ui, nil)
	want := slices.Sorted(maps.Keys(LoginHandlers()))

	login, err := commands["login"]()
	require.NoError(t, err)
	require.Equal(t, want, slices.Sorted(maps.Keys(login.(*LoginCommand).Handlers)))

	authHelp, err := commands["auth help"]()
	require.NoError(t, err)
	require.Equal(t, want, slices.Sorted(maps.Keys(authHelp.(*AuthHelpCommand).Handlers)))
}
