package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func testStatusCommand(tb testing.TB) (*cli.MockUi, *StatusCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &StatusCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func TestStatusCommand_Run(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		args   []string
		sealed bool
		out    string
		code   int
	}{
		{
			"unsealed",
			nil,
			false,
			"Sealed          false",
			0,
		},
		{
			"sealed",
			nil,
			true,
			"Sealed             true",
			2,
		},
		{
			"args",
			[]string{"foo"},
			false,
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

				if tc.sealed {
					if err := client.Sys().Seal(); err != nil {
						t.Fatal(err)
					}
				}

				ui, cmd := testStatusCommand(t)
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

		ui, cmd := testStatusCommand(t)
		cmd.client = client

		code := cmd.Run([]string{})
		if exp := 1; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Error checking seal status: "
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}
	})

	t.Run("no_tabs", func(t *testing.T) {
		t.Parallel()

		_, cmd := testStatusCommand(t)
		assertNoTabs(t, cmd)
	})
}
