// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"testing"

	"github.com/hashicorp/cli"
	base "github.com/hashicorp/vault/command/base"
	"github.com/stretchr/testify/require"
)

// TestOperatorRaftAutopilotStateCommand_Flags verifies that the autopilot
// state command defaults to the "pretty" output format, both as the -format
// flag's default and on a UI still set to the global "table" default: the
// state endpoint returns nested values the table format cannot display.
func TestOperatorRaftAutopilotStateCommand_Flags(t *testing.T) {
	t.Parallel()

	ui := &base.VaultUI{Ui: cli.NewMockUi(), Format: "table"}
	cmd := &OperatorRaftAutopilotStateCommand{BaseCommand: &base.BaseCommand{UI: ui}}

	fl := cmd.Flags().Lookup("format")
	require.NotNil(t, fl)
	require.Equal(t, "pretty", fl.DefValue)
	require.Equal(t, "pretty", ui.Format)
}
