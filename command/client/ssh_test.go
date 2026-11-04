// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/cli"
	"github.com/hashicorp/vault/api"
	base "github.com/hashicorp/vault/command/base"
	"github.com/stretchr/testify/require"
)

func testSSHCommand(tb testing.TB) (*cli.MockUi, *SSHCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &SSHCommand{
		BaseCommand: &base.BaseCommand{
			UI: ui,
		},
	}
}

// TestParseSSHCommand preserves SSH username/port precedence and stops parsing
// at the remote command, whose flags must not affect Vault's credential request.
func TestParseSSHCommand(t *testing.T) {
	t.Parallel()

	tests := []struct {
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
			"Flags after the ssh command are not passed because they are part of the command",
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
		{
			"Allow single args which don't have a value",
			[]string{
				"-v",
				"hostname",
			},
			"hostname",
			"",
			"",
			nil,
		},
		{
			"Allow single args before and after the hostname and command",
			[]string{
				"-v",
				"hostname",
				"-v",
				"command",
				"-v",
			},
			"hostname",
			"",
			"",
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			_, cmd := testSSHCommand(t)
			hostname, username, port, err := cmd.parseSSHCommand(test.args)
			require.Equal(t, test.err, err, "SSH argument parsing error")
			require.Equal(t, test.hostname, hostname, "SSH destination")
			require.Equal(t, test.username, username, "SSH username precedence")
			require.Equal(t, test.port, port, "SSH port precedence")
		})
	}
}

// TestIsSingleSSHArg verifies that valueless SSH options do not consume the
// destination as an option value.
func TestIsSingleSSHArg(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		arg  string
		want bool
	}{
		{
			"-v is a single ssh arg",
			"-v",
			true,
		},
		{
			"-o is NOT a single ssh arg",
			"-o",
			false,
		},
		{
			"Repeated args like -vvv is still a single ssh arg",
			"-vvv",
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			_, cmd := testSSHCommand(t)
			got := cmd.isSingleSSHArg(test.arg)
			require.Equal(t, test.want, got, "SSH option %q", test.arg)
		})
	}
}

// TestSSHCommandOmitFlagWarning verifies that remote flags are not mistaken for
// misplaced Vault flags, without reading a user's SSH keys or contacting Vault.
func TestSSHCommandOmitFlagWarning(t *testing.T) {
	t.Parallel()

	ui, cmd := testSSHCommand(t)
	client, err := api.NewClient(&api.Config{Address: "http://127.0.0.1:1"})
	require.NoError(t, err, "construct an unused client")
	cmd.SetClient(client)
	cmd.SetHCPTokenHelper(testHCPTokenHelper{})
	publicKeyPath := filepath.Join(t.TempDir(), "missing.pub")
	_, err = os.Stat(publicKeyPath)
	require.ErrorIs(t, err, os.ErrNotExist, "the failure must not depend on a user's keys")
	require.Equal(t, 1, cmd.Run([]string{
		"-mode=ca", "-role=test-role", "-public-key-path=" + publicKeyPath,
		"user@127.0.0.1", "-extraFlag", "bug",
	}), "a missing test-local public key must fail before any API request")

	combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
	require.Contains(t, combined, "failed to read public key", "confirm the controlled failure")
	require.NotContains(t, combined, "Command flags must be provided before positional arguments", "remote flags are not Vault flags")
}
