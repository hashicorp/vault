package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func testTokenRevokeCommand(tb testing.TB) (*cli.MockUi, *TokenRevokeCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &TokenRevokeCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func TestTokenRevokeCommand_Run(t *testing.T) {
	t.Parallel()

	validations := []struct {
		name string
		args []string
		out  string
		code int
	}{
		{
			"bad_mode",
			[]string{"-mode=banana"},
			"Invalid mode",
			1,
		},
		{
			"empty",
			nil,
			"Not enough arguments",
			1,
		},
		{
			"args_with_self",
			[]string{"-self", "abcd1234"},
			"Too many arguments",
			1,
		},
		{
			"too_many_args",
			[]string{"abcd1234", "efgh5678"},
			"Too many arguments",
			1,
		},
		{
			"self_and_accessor",
			[]string{"-self", "-accessor"},
			"Cannot use -self with -accessor",
			1,
		},
		{
			"self_and_mode",
			[]string{"-self", "-mode=orphan"},
			"Cannot use -self with -mode",
			1,
		},
		{
			"accessor_and_mode_orphan",
			[]string{"-accessor", "-mode=orphan", "abcd1234"},
			"Cannot use -accessor with -mode=orphan",
			1,
		},
		{
			"accessor_and_mode_path",
			[]string{"-accessor", "-mode=path", "abcd1234"},
			"Cannot use -accessor with -mode=path",
			1,
		},
	}

	t.Run("validations", func(t *testing.T) {
		t.Parallel()

		for _, tc := range validations {
			tc := tc

			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				client, closer := testVaultServer(t)
				defer closer()

				ui, cmd := testTokenRevokeCommand(t)
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

		ui, cmd := testTokenRevokeCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			token,
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Success! Revoked token"
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}

		secret, err := client.Auth().Token().Lookup(token)
		if secret != nil || err == nil {
			t.Errorf("expected token to be revoked: %#v", secret)
		}
	})

	t.Run("self", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		ui, cmd := testTokenRevokeCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"-self",
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Success! Revoked token"
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}

		secret, err := client.Auth().Token().LookupSelf()
		if secret != nil || err == nil {
			t.Errorf("expected token to be revoked: %#v", secret)
		}
	})

	t.Run("accessor", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		token, accessor := testTokenAndAccessor(t, client)

		ui, cmd := testTokenRevokeCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"-accessor",
			accessor,
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Success! Revoked token"
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}

		secret, err := client.Auth().Token().Lookup(token)
		if secret != nil || err == nil {
			t.Errorf("expected token to be revoked: %#v", secret)
		}
	})

	t.Run("communication_failure", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServerBad(t)
		defer closer()

		ui, cmd := testTokenRevokeCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"abcd1234",
		})
		if exp := 2; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Error revoking token: "
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}
	})

	t.Run("no_tabs", func(t *testing.T) {
		t.Parallel()

		_, cmd := testTokenRevokeCommand(t)
		assertNoTabs(t, cmd)
	})
}
