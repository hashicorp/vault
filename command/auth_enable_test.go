package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func testAuthEnableCommand(tb testing.TB) (*cli.MockUi, *AuthEnableCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &AuthEnableCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func TestAuthEnableCommand_Run(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		args []string
		out  string
		code int
	}{
		{
			"not_enough_args",
			nil,
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
			"not_a_valid_auth",
			[]string{"nope_definitely_not_a_valid_mount_like_ever"},
			"",
			2,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			client, closer := testVaultServer(t)
			defer closer()

			ui, cmd := testAuthEnableCommand(t)
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

	t.Run("integration", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		ui, cmd := testAuthEnableCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"-path", "auth_integration/",
			"-description", "The best kind of test",
			"userpass",
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Success! Enabled userpass auth method at:"
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}

		auths, err := client.Sys().ListAuth()
		if err != nil {
			t.Fatal(err)
		}

		authInfo, ok := auths["auth_integration/"]
		if !ok {
			t.Fatalf("expected mount to exist")
		}
		if exp := "userpass"; authInfo.Type != exp {
			t.Errorf("expected %q to be %q", authInfo.Type, exp)
		}
		if exp := "The best kind of test"; authInfo.Description != exp {
			t.Errorf("expected %q to be %q", authInfo.Description, exp)
		}
	})

	t.Run("communication_failure", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServerBad(t)
		defer closer()

		ui, cmd := testAuthEnableCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"userpass",
		})
		if exp := 2; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Error enabling userpass auth: "
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}
	})

	t.Run("no_tabs", func(t *testing.T) {
		t.Parallel()

		_, cmd := testAuthEnableCommand(t)
		assertNoTabs(t, cmd)
	})
}
