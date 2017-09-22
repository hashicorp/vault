package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func testPathHelpCommand(tb testing.TB) (*cli.MockUi, *PathHelpCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &PathHelpCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func TestPathHelpCommand_Run(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		args []string
		out  string
		code int
	}{
		{
			"not_enough_args",
			[]string{},
			"Not enough arguments",
			1,
		},
		{
			"too_many_args",
			[]string{"foo", "bar"},
			"Too many arguments",
			1,
		},
		{
			"not_found",
			[]string{"nope/not/once/never"},
			"",
			2,
		},
		{
			"kv",
			[]string{"secret/"},
			"The kv backend",
			0,
		},
		{
			"sys",
			[]string{"sys/mounts"},
			"currently mounted backends",
			0,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			client, closer := testVaultServer(t)
			defer closer()

			ui, cmd := testPathHelpCommand(t)
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

	t.Run("communication_failure", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServerBad(t)
		defer closer()

		ui, cmd := testPathHelpCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"sys/mounts",
		})
		if exp := 2; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Error retrieving help: "
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}
	})

	t.Run("no_tabs", func(t *testing.T) {
		t.Parallel()

		_, cmd := testPathHelpCommand(t)
		assertNoTabs(t, cmd)
	})
}
