package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func testPoliciesDeprecatedCommand(tb testing.TB) (*cli.MockUi, *PoliciesDeprecatedCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &PoliciesDeprecatedCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func TestPoliciesDeprecatedCommand_Run(t *testing.T) {
	t.Parallel()

	// TODO: remove in 0.9.0
	t.Run("deprecated_arg", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		ui, cmd := testPoliciesDeprecatedCommand(t)
		cmd.client = client

		// vault policies ARG -> vault policy read ARG
		code := cmd.Run([]string{"default"})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d: %s", code, exp, ui.ErrorWriter.String())
		}
		stdout := ui.OutputWriter.String()

		if expected := "token/"; !strings.Contains(stdout, expected) {
			t.Errorf("expected %q to contain %q", stdout, expected)
		}
	})

	t.Run("deprecated_no_args", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		ui, cmd := testPoliciesDeprecatedCommand(t)
		cmd.client = client

		// vault policies -> vault policy list
		code := cmd.Run([]string{})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d: %s", code, exp, ui.ErrorWriter.String())
		}
		stdout := ui.OutputWriter.String()

		if expected := "root"; !strings.Contains(stdout, expected) {
			t.Errorf("expected %q to contain %q", stdout, expected)
		}
	})

	t.Run("deprecated_with_flags", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		ui, cmd := testPoliciesDeprecatedCommand(t)
		cmd.client = client

		// vault policies -flag -> vault policy list
		code := cmd.Run([]string{
			"-address", client.Address(),
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d: %s", code, exp, ui.ErrorWriter.String())
		}
		stdout := ui.OutputWriter.String()

		if expected := "root"; !strings.Contains(stdout, expected) {
			t.Errorf("expected %q to contain %q", stdout, expected)
		}
	})

	t.Run("no_tabs", func(t *testing.T) {
		t.Parallel()

		_, cmd := testPoliciesDeprecatedCommand(t)
		assertNoTabs(t, cmd)
	})
}
