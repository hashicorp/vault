package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"

	credUserpass "github.com/hashicorp/vault/builtin/credential/userpass"
)

func testAuthHelpCommand(tb testing.TB) (*cli.MockUi, *AuthHelpCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &AuthHelpCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
		Handlers: map[string]LoginHandler{
			"userpass": &credUserpass.CLIHandler{
				DefaultMount: "userpass",
			},
		},
	}
}

func TestAuthHelpCommand_Run(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		args []string
		out  string
		code int
	}{
		{
			"too_many_args",
			[]string{"foo", "bar"},
			"Too many arguments",
			1,
		},
		{
			"not_enough_args",
			nil,
			"Not enough arguments",
			1,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			client, closer := testVaultServer(t)
			defer closer()

			ui, cmd := testAuthHelpCommand(t)
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

	t.Run("path", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		if err := client.Sys().EnableAuth("foo", "userpass", ""); err != nil {
			t.Fatal(err)
		}

		ui, cmd := testAuthHelpCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"foo/",
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Usage: vault login -method=userpass"
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}
	})

	t.Run("type", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		// No mounted auth methods

		ui, cmd := testAuthHelpCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"userpass",
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Usage: vault login -method=userpass"
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}
	})

	t.Run("communication_failure", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServerBad(t)
		defer closer()

		ui, cmd := testAuthHelpCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"sys/mounts",
		})
		if exp := 2; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Error listing auth methods: "
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}
	})

	t.Run("no_tabs", func(t *testing.T) {
		t.Parallel()

		_, cmd := testAuthHelpCommand(t)
		assertNoTabs(t, cmd)
	})
}
