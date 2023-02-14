package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func testEventsSubscribeCommand(tb testing.TB) (*cli.MockUi, *EventsSubscribeCommands) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &EventsSubscribeCommands{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

// TestEventsSubscribeCommand_Run tests that the command argument parsing is working as expected.
func TestEventsSubscribeCommand_Run(t *testing.T) {
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
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			client, closer := testVaultServer(t)
			defer closer()

			ui, cmd := testEventsSubscribeCommand(t)
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
}
