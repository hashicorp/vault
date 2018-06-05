package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func testDeleteCommand(tb testing.TB) (*cli.MockUi, *DeleteCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &DeleteCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func TestDeleteCommand_Run(t *testing.T) {
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

	t.Run("validations", func(t *testing.T) {
		t.Parallel()

		for _, tc := range cases {
			tc := tc

			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				client, closer := testVaultServer(t)
				defer closer()

				ui, cmd := testDeleteCommand(t)
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

		if _, err := client.Logical().Write("secret/delete/foo", map[string]interface{}{
			"foo": "bar",
		}); err != nil {
			t.Fatal(err)
		}

		ui, cmd := testDeleteCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"secret/delete/foo",
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Success! Data deleted (if it existed) at: secret/delete/foo"
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}

		secret, _ := client.Logical().Read("secret/delete/foo")
		if secret != nil {
			t.Errorf("expected deletion: %#v", secret)
		}
	})

	t.Run("communication_failure", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServerBad(t)
		defer closer()

		ui, cmd := testDeleteCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"secret/delete/foo",
		})
		if exp := 2; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Error deleting secret/delete/foo: "
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}
	})

	t.Run("no_tabs", func(t *testing.T) {
		t.Parallel()

		_, cmd := testDeleteCommand(t)
		assertNoTabs(t, cmd)
	})
}
