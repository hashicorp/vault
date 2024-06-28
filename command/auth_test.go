// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"testing"

	"github.com/hashicorp/cli"
	"github.com/hashicorp/vault/command/token"
)

func testAuthCommand(tb testing.TB) (*cli.MockUi, *AuthCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &AuthCommand{
		BaseCommand: &BaseCommand{
			UI: ui,

			// Override to our own token helper
			tokenHelper: token.NewTestingTokenHelper(),
		},
	}
}

func TestAuthCommand_Run(t *testing.T) {
	t.Parallel()

	t.Run("no_tabs", func(t *testing.T) {
		t.Parallel()

		_, cmd := testAuthCommand(t)
		assertNoTabs(t, cmd)
	})
}
