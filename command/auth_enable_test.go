// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/go-test/deep"
	"github.com/hashicorp/vault/helper/builtinplugins"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/mitchellh/cli"
)

// credentialBackendAdjustmentFactor allows for adjusting test assertions for
// credential backends. Add 1 to account for the "token" backend, which is visible
// when you walk the filesystem but is treated as special and excluded from the registry.
// Subtract 1 to account for "oidc" which is an alias of "jwt" and not a separate plugin.
var credentialBackendAdjustmentFactor = 1 - 1

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
				t.Errorf("expected command return code to be %d, got %d", tc.code, code)
			}

			combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
			if !strings.Contains(combined, tc.out) {
				t.Errorf("expected %q in response\n got: %+v", tc.out, combined)
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
			"-audit-non-hmac-request-keys", "foo,bar",
			"-audit-non-hmac-response-keys", "foo,bar",
			"-passthrough-request-headers", "authorization,authentication",
			"-passthrough-request-headers", "www-authentication",
			"-allowed-response-headers", "authorization",
			"-listing-visibility", "unauth",
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
		if diff := deep.Equal([]string{"authorization,authentication", "www-authentication"}, authInfo.Config.PassthroughRequestHeaders); len(diff) > 0 {
			t.Errorf("Failed to find expected values in PassthroughRequestHeaders. Difference is: %v", diff)
		}
		if diff := deep.Equal([]string{"authorization"}, authInfo.Config.AllowedResponseHeaders); len(diff) > 0 {
			t.Errorf("Failed to find expected values in AllowedResponseHeaders. Difference is: %v", diff)
		}
		if diff := deep.Equal([]string{"foo,bar"}, authInfo.Config.AuditNonHMACRequestKeys); len(diff) > 0 {
			t.Errorf("Failed to find expected values in AuditNonHMACRequestKeys. Difference is: %v", diff)
		}
		if diff := deep.Equal([]string{"foo,bar"}, authInfo.Config.AuditNonHMACResponseKeys); len(diff) > 0 {
			t.Errorf("Failed to find expected values in AuditNonHMACResponseKeys. Difference is: %v", diff)
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

		modFile, err := ioutil.ReadFile("../go.mod")
		if err != nil {
			t.Fatal(err)
		}
		modLines := strings.Split(string(modFile), "\n")
		for _, p := range modLines {
			splitLine := strings.Split(strings.TrimSpace(p), " ")
			if len(splitLine) == 0 {
				continue
			}
			potPlug := strings.TrimPrefix(splitLine[0], "github.com/hashicorp/")
			if strings.HasPrefix(potPlug, "vault-plugin-auth-") {
				backends = append(backends, strings.TrimPrefix(potPlug, "vault-plugin-auth-"))
			}
		}
		// Since "pcf" plugin in the Vault registry is also pointed at the "vault-plugin-auth-cf"
		// repository, we need to manually append it here so it'll tie out with our expected number
		// of credential backends.
		backends = append(backends, "pcf")

		expected := len(builtinplugins.Registry.Keys(consts.PluginTypeCredential)) + credentialBackendAdjustmentFactor
		if len(backends) != expected {
			t.Fatalf("expected %d credential backends, got %d", expected, len(backends))
		}

		for _, b := range backends {
			var expectedResult int = 0

			// Not a builtin
			if b == "token" {
				continue
			}

			ui, cmd := testAuthEnableCommand(t)
			cmd.client = client

			actualResult := cmd.Run([]string{
				b,
			})

			// Need to handle deprecated builtins specially
			status, _ := builtinplugins.Registry.DeprecationStatus(b, consts.PluginTypeCredential)
			if status == consts.PendingRemoval || status == consts.Removed {
				expectedResult = 2
			}

			if actualResult != expectedResult {
				t.Errorf("type: %s - got: %d, expected: %d - %s", b, actualResult, expectedResult, ui.OutputWriter.String()+ui.ErrorWriter.String())
			}
		}
	})
}
