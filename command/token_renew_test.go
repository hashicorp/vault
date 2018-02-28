package command

import (
	"encoding/json"
	"strconv"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func testTokenRenewCommand(tb testing.TB) (*cli.MockUi, *TokenRenewCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &TokenRenewCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func TestTokenRenewCommand_Run(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		args []string
		out  string
		code int
	}{
		{
			"too_many_args",
			[]string{"foo", "bar", "baz"},
			"Too many arguments",
			1,
		},
		{
			"default",
			nil,
			"",
			0,
		},
		{
			"increment",
			[]string{"-increment", "60s"},
			"",
			0,
		},
		{
			"increment_no_suffix",
			[]string{"-increment", "60"},
			"",
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

				// Login with the token so we can renew-self.
				token, _ := testTokenAndAccessor(t, client)
				client.SetToken(token)

				ui, cmd := testTokenRenewCommand(t)
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

		_, cmd := testTokenRenewCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"-increment", "30m",
			token,
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		secret, err := client.Auth().Token().Lookup(token)
		if err != nil {
			t.Fatal(err)
		}

		str := string(secret.Data["ttl"].(json.Number))
		ttl, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			t.Fatalf("bad ttl: %#v", secret.Data["ttl"])
		}
		if exp := int64(1800); ttl > exp {
			t.Errorf("expected %d to be <= to %d", ttl, exp)
		}
	})

	t.Run("self", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		token, _ := testTokenAndAccessor(t, client)

		// Get the old token and login as the new token. We need the old token
		// to query after the lookup, but we need the new token on the client.
		oldToken := client.Token()
		client.SetToken(token)

		_, cmd := testTokenRenewCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"-increment", "30m",
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		client.SetToken(oldToken)
		secret, err := client.Auth().Token().Lookup(token)
		if err != nil {
			t.Fatal(err)
		}

		str := string(secret.Data["ttl"].(json.Number))
		ttl, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			t.Fatalf("bad ttl: %#v", secret.Data["ttl"])
		}
		if exp := int64(1800); ttl > exp {
			t.Errorf("expected %d to be <= to %d", ttl, exp)
		}
	})

	t.Run("communication_failure", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServerBad(t)
		defer closer()

		ui, cmd := testTokenRenewCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"foo/bar",
		})
		if exp := 2; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Error renewing token: "
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}
	})

	t.Run("no_tabs", func(t *testing.T) {
		t.Parallel()

		_, cmd := testTokenRenewCommand(t)
		assertNoTabs(t, cmd)
	})
}
