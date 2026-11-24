// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"testing"

	"github.com/hashicorp/cli"
	base "github.com/hashicorp/vault/command/base"
	clientcmd "github.com/hashicorp/vault/command/client"
	"github.com/hashicorp/vault/command/token"
)

func testAuthCommand(tb testing.TB) (*cli.MockUi, *clientcmd.AuthCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	cmd := &clientcmd.AuthCommand{
		BaseCommand: &base.BaseCommand{
			UI: ui,
		},
	}

	// Override to our own token helper
	cmd.SetTokenHelper(token.NewTestingTokenHelper())
	return ui, cmd
}

func TestAuthCommand_Run(t *testing.T) {
	t.Parallel()

	t.Run("no_tabs", func(t *testing.T) {
		t.Parallel()

		_, cmd := testAuthCommand(t)
		assertNoTabs(t, cmd)
	})
}
