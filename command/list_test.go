package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func testListCommand(tb testing.TB) (*cli.MockUi, *ListCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &ListCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func TestListCommand_Run(t *testing.T) {
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
			"default",
			[]string{"secret/list"},
			"bar\nbaz\nfoo",
			0,
		},
		{
			"default_slash",
			[]string{"secret/list/"},
			"bar\nbaz\nfoo",
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

				keys := []string{
					"secret/list/foo",
					"secret/list/bar",
					"secret/list/baz",
				}
				for _, k := range keys {
					if _, err := client.Logical().Write(k, map[string]interface{}{
						"foo": "bar",
					}); err != nil {
						t.Fatal(err)
					}
				}

				ui, cmd := testListCommand(t)
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

		ui, cmd := testListCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"secret/list",
		})
		if exp := 2; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Error listing secret/list/: "
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}
	})

	t.Run("no_tabs", func(t *testing.T) {
		t.Parallel()

		_, cmd := testListCommand(t)
		assertNoTabs(t, cmd)
	})
}
