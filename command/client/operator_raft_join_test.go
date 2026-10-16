// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"testing"

	"github.com/hashicorp/cli"
	base "github.com/hashicorp/vault/command/base"
	"github.com/stretchr/testify/require"
)

// TestOperatorRaftJoinCommand_AutoJoinRequiresServerConfig verifies the server
// directs auto-join users to trusted server configuration instead of sending
// discovery metadata to the unauthenticated join API.
func TestOperatorRaftJoinCommand_AutoJoinRequiresServerConfig(t *testing.T) {
	t.Parallel()

	ui := cli.NewMockUi()
	cmd := &OperatorRaftJoinCommand{BaseCommand: &base.BaseCommand{UI: ui}}

	require.Equal(t, 1, cmd.Run([]string{"provider=aws region=eu-west-1"}), "auto-join arguments should be rejected")
	require.Contains(t, ui.ErrorWriter.String(), "configure raft.retry_join in the Vault server configuration", "the error should explain the supported alternative")
}

// TestOperatorRaftJoinCommand_HelpExplainsAutoJoinConfiguration verifies the
// command help documents that auto-join belongs in server configuration.
func TestOperatorRaftJoinCommand_HelpExplainsAutoJoinConfiguration(t *testing.T) {
	t.Parallel()

	cmd := &OperatorRaftJoinCommand{BaseCommand: &base.BaseCommand{UI: cli.NewMockUi()}}

	help := cmd.Help()
	require.Contains(t, help, "raft.retry_join", "help should name the supported configuration")
	require.NotContains(t, help, "auto-join-scheme", "help should not advertise removed request flags")
	require.NotContains(t, help, "provider=aws", "help should not show request-based auto-join")
	require.Contains(t, help, "http://127.0.0.2:8200", "help should retain the direct leader-address example")
}
