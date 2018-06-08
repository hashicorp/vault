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

func TestParseSSHCommand(t *testing.T) {
	t.Parallel()

	_, cmd := testSSHCommand(t)
	var tests = []struct {
		name     string
		args     []string
		hostname string
		username string
		port     string
		err      error
	}{
		{
			"Parse just a hostname",
			[]string{
				"hostname",
			},
			"hostname",
			"",
			"",
			nil,
		},
		{
			"Parse the standard username@hostname",
			[]string{
				"username@hostname",
			},
			"hostname",
			"username",
			"",
			nil,
		},
		{
			"Parse the username out of -o User=username",
			[]string{
				"-o", "User=username",
				"hostname",
			},
			"hostname",
			"username",
			"",
			nil,
		},
		{
			"If the username is specified with -o User=username and realname@hostname prefer realname@",
			[]string{
				"-o", "User=username",
				"realname@hostname",
			},
			"hostname",
			"realname",
			"",
			nil,
		},
		{
			"Parse the port out of -o Port=2222",
			[]string{
				"-o", "Port=2222",
				"hostname",
			},
			"hostname",
			"",
			"2222",
			nil,
		},
		{
			"Parse the port out of -p 2222",
			[]string{
				"-p", "2222",
				"hostname",
			},
			"hostname",
			"",
			"2222",
			nil,
		},
		{
			"If port is defined with -o Port=2222 and -p 2244 prefer -p",
			[]string{
				"-p", "2244",
				"-o", "Port=2222",
				"hostname",
			},
			"hostname",
			"",
			"2244",
			nil,
		},
		{
			"Ssh args with a command",
			[]string{
				"hostname",
				"command",
			},
			"hostname",
			"",
			"",
			nil,
		},
		{
			"Flags after the ssh command are not pased because they are part of the command",
			[]string{
				"username@hostname",
				"command",
				"-p 22",
			},
			"hostname",
			"username",
			"",
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			hostname, username, port, err := cmd.parseSSHCommand(test.args)
			if err != test.err {
				t.Errorf("got error: %q want %q", err, test.err)
			}
			if hostname != test.hostname {
				t.Errorf("got hostname: %q want %q", hostname, test.hostname)
			}
			if username != test.username {
				t.Errorf("got username: %q want %q", username, test.username)
			}
			if port != test.port {
				t.Errorf("got port: %q want %q", port, test.port)
			}
		})
	}
}
