package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func testTokenLookupCommand(tb testing.TB) (*cli.MockUi, *TokenLookupCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &TokenLookupCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func TestTokenLookupCommand_Run(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		args []string
		out  string
		code int
	}{
		{
			"accessor_no_args",
			[]string{"-accessor"},
			"Not enough arguments",
			1,
		},
		{
			"accessor_too_many_args",
			[]string{"-accessor", "abcd1234", "efgh5678"},
			"Too many arguments",
			1,
		},
		{
			"too_many_args",
			[]string{"abcd1234", "efgh5678"},
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

				ui, cmd := testTokenLookupCommand(t)
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

	t.Run("token", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		token, _ := testTokenAndAccessor(t, client)

		ui, cmd := testTokenLookupCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			token,
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := token
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}
	})

	t.Run("self", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		ui, cmd := testTokenLookupCommand(t)
		cmd.client = client

		code := cmd.Run([]string{})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "display_name"
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}
	})

	t.Run("accessor", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		_, accessor := testTokenAndAccessor(t, client)

		ui, cmd := testTokenLookupCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"-accessor",
			accessor,
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := accessor
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}
	})

	t.Run("communication_failure", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServerBad(t)
		defer closer()

		ui, cmd := testTokenLookupCommand(t)
		cmd.client = client

		code := cmd.Run([]string{})
		if exp := 2; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Error looking up token: "
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}
	})

	t.Run("no_tabs", func(t *testing.T) {
		t.Parallel()

		_, cmd := testTokenLookupCommand(t)
		assertNoTabs(t, cmd)
	})
}
