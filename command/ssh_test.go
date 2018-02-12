package command

import (
	"testing"

	"github.com/mitchellh/cli"
)

func testSSHCommand(tb testing.TB) (*cli.MockUi, *SSHCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &SSHCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func TestSSHCommand_Run(t *testing.T) {
	t.Parallel()
	t.Skip("Need a way to setup target infrastructure")
}
