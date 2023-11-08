// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"strings"
	"testing"

	"github.com/hashicorp/vault/version"
	"github.com/mitchellh/cli"
)

func testVersionCommand(tb testing.TB) (*cli.MockUi, *VersionCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &VersionCommand{
		VersionInfo: &version.VersionInfo{},
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func TestVersionCommand_Run(t *testing.T) {
	t.Parallel()

	t.Run("output", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		ui, cmd := testVersionCommand(t)
		cmd.ApiClient = client

		code := cmd.Run(nil)
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Vault"
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to equal %q", combined, expected)
		}
	})

	t.Run("no_tabs", func(t *testing.T) {
		t.Parallel()

		_, cmd := testVersionCommand(t)
		assertNoTabs(t, cmd)
	})
}
