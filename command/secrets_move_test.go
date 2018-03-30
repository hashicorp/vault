package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func testSecretsMoveCommand(tb testing.TB) (*cli.MockUi, *SecretsMoveCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &SecretsMoveCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func TestSecretsMoveCommand_Run(t *testing.T) {
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
			[]string{"foo", "bar", "baz"},
			"Too many arguments",
			1,
		},
		{
			"non_existent",
			[]string{"not_real", "over_here"},
			"Error moving secrets engine not_real/ to over_here/",
			2,
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

				ui, cmd := testSecretsMoveCommand(t)
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

	t.Run("integration", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		ui, cmd := testSecretsMoveCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"secret/", "generic/",
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Success! Moved secrets engine secret/ to: generic/"
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}

		mounts, err := client.Sys().ListMounts()
		if err != nil {
			t.Fatal(err)
		}

		if _, ok := mounts["generic/"]; !ok {
			t.Errorf("expected mount at generic/: %#v", mounts)
		}
	})

	t.Run("communication_failure", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServerBad(t)
		defer closer()

		ui, cmd := testSecretsMoveCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"secret/", "generic/",
		})
		if exp := 2; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Error moving secrets engine secret/ to generic/:"
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}
	})

	t.Run("no_tabs", func(t *testing.T) {
		t.Parallel()

		_, cmd := testSecretsMoveCommand(t)
		assertNoTabs(t, cmd)
	})
}
