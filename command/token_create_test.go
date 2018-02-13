package command

import (
	"reflect"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func testTokenCreateCommand(tb testing.TB) (*cli.MockUi, *TokenCreateCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &TokenCreateCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func TestTokenCreateCommand_Run(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		args []string
		out  string
		code int
	}{
		{
			"too_many_args",
			[]string{"abcd1234"},
			"Too many arguments",
			1,
		},
		{
			"default",
			nil,
			"token",
			0,
		},
		{
			"metadata",
			[]string{"-metadata", "foo=bar", "-metadata", "zip=zap"},
			"token",
			0,
		},
		{
			"policies",
			[]string{"-policy", "foo", "-policy", "bar"},
			"token",
			0,
		},
		{
			"field",
			[]string{
				"-field", "token_renewable",
			},
			"false",
			0,
		},
		{
			"field_not_found",
			[]string{
				"-field", "not-a-real-field",
			},
			"not present in secret",
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

				ui, cmd := testTokenCreateCommand(t)
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

	t.Run("default", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		ui, cmd := testTokenCreateCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"-field", "token",
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		token := strings.TrimSpace(ui.OutputWriter.String())
		secret, err := client.Auth().Token().Lookup(token)
		if secret == nil || err != nil {
			t.Fatal(err)
		}
	})

	t.Run("metadata", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		ui, cmd := testTokenCreateCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"-metadata", "foo=bar",
			"-metadata", "zip=zap",
			"-field", "token",
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		token := strings.TrimSpace(ui.OutputWriter.String())
		secret, err := client.Auth().Token().Lookup(token)
		if secret == nil || err != nil {
			t.Fatal(err)
		}

		meta, ok := secret.Data["meta"].(map[string]interface{})
		if !ok {
			t.Fatalf("missing meta: %#v", secret)
		}
		if _, ok := meta["foo"]; !ok {
			t.Errorf("missing meta.foo: %#v", meta)
		}
		if _, ok := meta["zip"]; !ok {
			t.Errorf("missing meta.bar: %#v", meta)
		}
	})

	t.Run("policies", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		ui, cmd := testTokenCreateCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"-policy", "foo",
			"-policy", "bar",
			"-field", "token",
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		token := strings.TrimSpace(ui.OutputWriter.String())
		secret, err := client.Auth().Token().Lookup(token)
		if secret == nil || err != nil {
			t.Fatal(err)
		}

		raw, ok := secret.Data["policies"].([]interface{})
		if !ok {
			t.Fatalf("missing policies: %#v", secret)
		}

		policies := make([]string, len(raw))
		for i := range raw {
			policies[i] = raw[i].(string)
		}

		expected := []string{"bar", "default", "foo"}
		if !reflect.DeepEqual(policies, expected) {
			t.Errorf("expected %q to be %q", policies, expected)
		}
	})

	t.Run("communication_failure", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServerBad(t)
		defer closer()

		ui, cmd := testTokenCreateCommand(t)
		cmd.client = client

		code := cmd.Run([]string{})
		if exp := 2; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Error creating token: "
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}
	})

	t.Run("no_tabs", func(t *testing.T) {
		t.Parallel()

		_, cmd := testTokenCreateCommand(t)
		assertNoTabs(t, cmd)
	})
}
