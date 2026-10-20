// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package command

import (
	"slices"
	"strings"
	"testing"

	"github.com/hashicorp/cli"
	base "github.com/hashicorp/vault/command/base"
	"github.com/stretchr/testify/require"
)

// entServerOnlyCommands is empty in CE: enterprise-only commands are not built.
var entServerOnlyCommands []string

// TestInitCommands_CE verifies that the CE vault server binary registers
// exactly the commands it registered before server-only commands moved out
// of command/client, so the move did not drop or add any command.
func TestInitCommands_CE(t *testing.T) {
	t.Parallel()

	expected := []string{
		"agent",
		"agent generate-config",
		"audit",
		"audit disable",
		"audit enable",
		"audit list",
		"auth",
		"auth disable",
		"auth enable",
		"auth help",
		"auth list",
		"auth move",
		"auth tune",
		"debug",
		"delete",
		"events subscribe",
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
		"monitor",
		"namespace",
		"namespace create",
		"namespace delete",
		"namespace list",
		"namespace lock",
		"namespace lookup",
		"namespace patch",
		"namespace unlock",
		"operator",
		"operator diagnose",
		"operator generate-root",
		"operator init",
		"operator key-status",
		"operator members",
		"operator migrate",
		"operator raft",
		"operator raft autopilot get-config",
		"operator raft autopilot set-config",
		"operator raft autopilot state",
		"operator raft join",
		"operator raft list-peers",
		"operator raft remove-peer",
		"operator raft snapshot",
		"operator raft snapshot inspect",
		"operator raft snapshot restore",
		"operator raft snapshot save",
		"operator rekey",
		"operator rotate",
		"operator seal",
		"operator step-down",
		"operator unseal",
		"operator usage",
		"operator utilization",
		"patch",
		"path-help",
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
		"server",
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
		"unwrap",
		"version",
		"version-history",
		"write",
	}

	ui := cli.NewMockUi()
	stdin := strings.NewReader("")
	commands := initCommands(ui, ui, &base.RunOptions{Stdin: stdin})

	var names []string
	for name := range commands {
		names = append(names, name)
	}
	slices.Sort(names)
	require.Equal(t, expected, names)
}
