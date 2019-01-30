package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func testPluginListCommand(tb testing.TB) (*cli.MockUi, *PluginListCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &PluginListCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func TestPluginListCommand_Run(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		args []string
		out  string
		code int
	}{
		{
			"too_many_args",
			[]string{"foo", "fizz"},
			"Too many arguments",
			1,
		},
		{
			"lists",
			nil,
			"Plugins",
			0,
		},
	}

	t.Run("validations", func(t *testing.T) {
		t.Parallel()

		for _, tc := range cases {
			tc := tc

			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				client, closer := testVaultServer(t)
				defer closer()

				ui, cmd := testPluginListCommand(t)
				cmd.client = client

				code := cmd.Run(tc.args)
				if code != tc.code {
					t.Errorf("expected %d to be %d", code, tc.code)
				}

				combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
				if !strings.Contains(combined, tc.out) {
					t.Errorf("expected %q to contain %q", combined, tc.out)
				}
			})
		}
	})

	t.Run("communication_failure", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServerBad(t)
		defer closer()

		ui, cmd := testPluginListCommand(t)
		cmd.client = client

		code := cmd.Run([]string{"database"})
		if exp := 2; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Error listing available plugins: "
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}
	})

	t.Run("no_tabs", func(t *testing.T) {
		t.Parallel()

		_, cmd := testPluginListCommand(t)
		assertNoTabs(t, cmd)
	})
}
