package command

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/hashicorp/vault/helper/builtinplugins"
	"github.com/hashicorp/vault/helper/consts"
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

	t.Run("mount_all", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServerAllBackends(t)
		defer closer()

		files, err := ioutil.ReadDir("../builtin/credential")
		if err != nil {
			t.Fatal(err)
		}

		var backends []string
		for _, f := range files {
			if f.IsDir() {
				backends = append(backends, f.Name())
			}
		}

		plugins, err := ioutil.ReadDir("../vendor/github.com/hashicorp")
		if err != nil {
			t.Fatal(err)
		}
		for _, p := range plugins {
			if p.IsDir() && strings.HasPrefix(p.Name(), "vault-plugin-auth-") {
				backends = append(backends, strings.TrimPrefix(p.Name(), "vault-plugin-auth-"))
			}
		}

		// Add 1 to account for the "token" backend, which is visible when you walk the filesystem but
		// is treated as special and excluded from the registry.
		// Subtract 1 to account for "oidc" which is an alias of "jwt" and not a separate plugin.
		expected := len(builtinplugins.Registry.Keys(consts.PluginTypeCredential))
		if len(backends) != expected {
			t.Fatalf("expected %d credential backends, got %d", expected, len(backends))
		}

		for _, b := range backends {
			if b == "token" {
				continue
			}

			ui, cmd := testAuthEnableCommand(t)
			cmd.client = client

			code := cmd.Run([]string{
				b,
			})
			if exp := 0; code != exp {
				t.Errorf("type %s, expected %d to be %d - %s", b, code, exp, ui.OutputWriter.String()+ui.ErrorWriter.String())
			}
		}
	})
}
