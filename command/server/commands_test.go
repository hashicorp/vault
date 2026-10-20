// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"testing"

	"github.com/hashicorp/cli"
	base "github.com/hashicorp/vault/command/base"
	clientcmd "github.com/hashicorp/vault/command/client"
	"github.com/stretchr/testify/require"
)

func Test_Commands_HCPInit(t *testing.T) {
	tests := map[string]struct {
		expectError      bool
		expectedErrorMsg string
	}{
		"initialize with success": {
			expectError: false,
		},
		"initialize with error: existing commands conflict with init commands": {
			expectError:      true,
			expectedErrorMsg: "Failed to initialize HCP commands.",
		},
	}

	for n, tst := range tests {
		n := n
		tst := tst

		t.Run(n, func(t *testing.T) {
			t.Parallel()

			mockUi := cli.NewMockUi()
			commands := initCommands(mockUi, nil, nil)
			if tst.expectError {
				clientcmd.InitHCPCommands(mockUi, commands)
				errMsg := mockUi.ErrorWriter.String()
				require.NotEmpty(t, errMsg)
				require.Contains(t, errMsg, tst.expectedErrorMsg)
			} else {
				errMsg := mockUi.ErrorWriter.String()
				require.Empty(t, errMsg)
				require.NotEmpty(t, commands)
			}
		})
	}
}

// serverOnlyCommands lists the commands that need no server internals, and so
// could be vault-client commands, but that PRD VLT-2015 keeps server-only.
var serverOnlyCommands = []string{
	"audit",
	"audit disable",
	"audit enable",
	"audit list",
	"debug",
	"events subscribe",
	"monitor",
	"operator generate-root",
	"operator init",
	"operator members",
	"operator raft",
	"operator raft autopilot get-config",
	"operator raft autopilot set-config",
	"operator raft autopilot state",
	"operator raft join",
	"operator raft list-peers",
	"operator raft remove-peer",
	"operator raft snapshot",
	"operator raft snapshot restore",
	"operator raft snapshot save",
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
	"policy list",
	"policy read",
	"policy write",
	"ssh",
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

// TestInitCommands_SupersetOfClient verifies that the vault server binary
// registers every vault-client command plus every command that PRD VLT-2015
// marks as server-only, and that vault-client registers none of the latter.
func TestInitCommands_SupersetOfClient(t *testing.T) {
	t.Parallel()

	ui := cli.NewMockUi()
	server := initCommands(ui, ui, &base.RunOptions{})
	client := clientcmd.InitClientCommands(ui, ui, &base.RunOptions{})

	for name := range client {
		require.Contains(t, server, name)
	}
	for _, name := range append(serverOnlyCommands, entServerOnlyCommands...) {
		require.Contains(t, server, name)
		require.NotContains(t, client, name)
	}
}
