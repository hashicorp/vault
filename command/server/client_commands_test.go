// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"maps"
	"slices"
	"strings"
	"testing"

	"github.com/hashicorp/cli"
	base "github.com/hashicorp/vault/command/base"
	clientcmd "github.com/hashicorp/vault/command/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestInitCommands_ClientBinaryParity checks that vault-client has every server
// command, either working or as a placeholder. We also check for stale aliases
// and make sure vault keeps the real commands.
func TestInitCommands_ClientBinaryParity(t *testing.T) {
	t.Parallel()

	// Keep these lists separate from the command maps so we don't test them against themselves.
	expectedReal := []string{
		"agent",
		"agent generate-config",
		"auth",
		"auth disable",
		"auth enable",
		"auth help",
		"auth list",
		"auth move",
		"auth tune",
		"delete",
		"hcp connect",
		"hcp disconnect",
		"kv",
		"kv delete",
		"kv destroy",
		"kv enable-versioning",
		"kv get",
		"kv list",
		"kv metadata",
		"kv metadata delete",
		"kv metadata get",
		"kv metadata patch",
		"kv metadata put",
		"kv patch",
		"kv put",
		"kv rollback",
		"kv undelete",
		"lease",
		"lease lookup",
		"lease renew",
		"lease revoke",
		"list",
		"login",
		"namespace",
		"namespace create",
		"namespace delete",
		"namespace list",
		"namespace lock",
		"namespace lookup",
		"namespace patch",
		"namespace unlock",
		"operator",
		"operator key-status",
		"operator rotate",
		"operator unseal",
		"patch",
		"path-help",
		"print",
		"print token",
		"proxy",
		"read",
		"secrets",
		"secrets disable",
		"secrets enable",
		"secrets list",
		"secrets move",
		"secrets tune",
		"ssh",
		"status",
		"token",
		"token capabilities",
		"token create",
		"token lookup",
		"token renew",
		"token revoke",
		"tpm",
		"tpm attest",
		"tpm cleanup",
		"tpm ek",
		"tpm enroll",
		"tpm inspect",
		"tpm reattest",
		"unwrap",
		"version",
		"write",
	}
	expectedReal = append(expectedReal, entClientBinaryRealCommands...)
	expectedUnavailable := []string{
		"audit",
		"audit disable",
		"audit enable",
		"audit list",
		"debug",
		"events",
		"events subscribe",
		"events subscriptions",
		"events subscriptions create",
		"events subscriptions delete",
		"events subscriptions get",
		"events subscriptions list",
		"license",
		"license get",
		"license inspect",
		"monitor",
		"operator diagnose",
		"operator generate-root",
		"operator import",
		"operator init",
		"operator members",
		"operator migrate",
		"operator raft",
		"operator raft autopilot",
		"operator raft autopilot get-config",
		"operator raft autopilot set-config",
		"operator raft autopilot state",
		"operator raft join",
		"operator raft list-peers",
		"operator raft remove-peer",
		"operator raft snapshot",
		"operator raft snapshot inspect",
		"operator raft snapshot load",
		"operator raft snapshot restore",
		"operator raft snapshot save",
		"operator raft snapshot unload",
		"operator rekey",
		"operator seal",
		"operator step-down",
		"operator usage",
		"operator utilization",
		"pki",
		"pki health-check",
		"pki issue",
		"pki list-intermediates",
		"pki reissue",
		"pki verify-sign",
		"plugin",
		"plugin deregister",
		"plugin info",
		"plugin list",
		"plugin register",
		"plugin reload",
		"plugin reload-status",
		"plugin runtime",
		"plugin runtime deregister",
		"plugin runtime info",
		"plugin runtime list",
		"plugin runtime register",
		"policy",
		"policy delete",
		"policy fmt",
		"policy list",
		"policy read",
		"policy write",
		"server",
		"transform",
		"transform import",
		"transform import-version",
		"transit",
		"transit envelope",
		"transit envelope decrypt",
		"transit envelope encrypt",
		"transit envelope header",
		"transit import",
		"transit import-version",
		"version-history",
	}
	expectedClasses := make(map[string]bool)
	for _, partition := range []struct {
		names       []string
		unavailable bool
	}{
		{expectedReal, false},
		{expectedUnavailable, true},
	} {
		for _, name := range partition.names {
			require.NotContains(t, expectedClasses, name, "duplicate or overlapping policy key")
			expectedClasses[name] = partition.unavailable
		}
	}

	serverUI, clientUI := cli.NewMockUi(), cli.NewMockUi()
	server := initCommands(serverUI, serverUI, &base.RunOptions{})
	client := clientcmd.InitClientBinaryCommands(clientUI, clientUI, &base.RunOptions{})
	for name := range server {
		assert.Contains(t, client, name, "every actual server key needs an exact client registration")
	}
	assert.ElementsMatch(t, slices.Collect(maps.Keys(expectedClasses)), slices.Collect(maps.Keys(client)), "client policy must be exhaustive")

	var real, unavailable, extra []string
	for name, factory := range client {
		wantUnavailable, reviewed := expectedClasses[name]
		require.True(t, reviewed, "unreviewed client key %q", name)
		cmd, err := factory()
		require.NoError(t, err, name)
		require.NotNil(t, cmd, name)
		_, isUnavailable := cmd.(*clientcmd.UnavailableCommand)
		if isUnavailable {
			unavailable = append(unavailable, name)
		} else {
			real = append(real, name)
		}
		require.Equal(t, wantUnavailable, isUnavailable, "wrong factory class for %q", name)
		if serverFactory, present := server[name]; !present {
			extra = append(extra, name)
		} else if isUnavailable {
			require.NotNil(t, serverFactory, "nil server factory for %q", name)
			serverCmd, serverErr := serverFactory()
			require.NoError(t, serverErr, "server factory for %q", name)
			require.NotNil(t, serverCmd, "nil server command for %q", name)
			_, serverUnavailable := serverCmd.(*clientcmd.UnavailableCommand)
			require.False(t, serverUnavailable, "server factory returned a client unavailable shim for %q", name)
		}
	}
	assert.ElementsMatch(t, expectedReal, real, "real client commands must match reviewed policy")
	assert.ElementsMatch(t, expectedUnavailable, unavailable, "unavailable shims must match reviewed policy")

	// The CLI creates these parent commands from their subcommands, so they don't
	// have their own factories in the server command map.
	parentAliases := []string{"events", "operator raft autopilot"}
	expectedExtra := append(parentAliases, entClientBinaryExtraAliases...)
	require.ElementsMatch(t, expectedExtra, extra, "only reviewed aliases may lack an exact server key")
	for _, parent := range parentAliases {
		hasDescendant := false
		for name := range server {
			if strings.HasPrefix(name, parent+" ") {
				hasDescendant = true
				break
			}
		}
		require.True(t, hasDescendant, "stale intermediate alias %q", parent)
	}
}
