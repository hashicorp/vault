package command

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/hashicorp/vault/helper/builtinplugins"
	"github.com/hashicorp/vault/helper/consts"
	"github.com/mitchellh/cli"
)

func testSecretsEnableCommand(tb testing.TB) (*cli.MockUi, *SecretsEnableCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &SecretsEnableCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func TestSecretsEnableCommand_Run(t *testing.T) {
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
			"not_a_valid_mount",
			[]string{"nope_definitely_not_a_valid_mount_like_ever"},
			"",
			2,
		},
		{
			"mount",
			[]string{"transit"},
			"Success! Enabled the transit secrets engine at: transit/",
			0,
		},
		{
			"mount_path",
			[]string{
				"-path", "transit_mount_point",
				"transit",
			},
			"Success! Enabled the transit secrets engine at: transit_mount_point/",
			0,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			client, closer := testVaultServer(t)
			defer closer()

			ui, cmd := testSecretsEnableCommand(t)
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

		ui, cmd := testSecretsEnableCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"-path", "mount_integration/",
			"-description", "The best kind of test",
			"-default-lease-ttl", "30m",
			"-max-lease-ttl", "1h",
			"-force-no-cache",
			"pki",
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Success! Enabled the pki secrets engine at: mount_integration/"
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}

		mounts, err := client.Sys().ListMounts()
		if err != nil {
			t.Fatal(err)
		}

		mountInfo, ok := mounts["mount_integration/"]
		if !ok {
			t.Fatalf("expected mount to exist")
		}
		if exp := "pki"; mountInfo.Type != exp {
			t.Errorf("expected %q to be %q", mountInfo.Type, exp)
		}
		if exp := "The best kind of test"; mountInfo.Description != exp {
			t.Errorf("expected %q to be %q", mountInfo.Description, exp)
		}
		if exp := 1800; mountInfo.Config.DefaultLeaseTTL != exp {
			t.Errorf("expected %d to be %d", mountInfo.Config.DefaultLeaseTTL, exp)
		}
		if exp := 3600; mountInfo.Config.MaxLeaseTTL != exp {
			t.Errorf("expected %d to be %d", mountInfo.Config.MaxLeaseTTL, exp)
		}
		if exp := true; mountInfo.Config.ForceNoCache != exp {
			t.Errorf("expected %t to be %t", mountInfo.Config.ForceNoCache, exp)
		}
	})

	t.Run("communication_failure", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServerBad(t)
		defer closer()

		ui, cmd := testSecretsEnableCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"pki",
		})
		if exp := 2; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Error enabling: "
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}
	})

	t.Run("no_tabs", func(t *testing.T) {
		t.Parallel()

		_, cmd := testSecretsEnableCommand(t)
		assertNoTabs(t, cmd)
	})

	t.Run("mount_all", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServerAllBackends(t)
		defer closer()

		files, err := ioutil.ReadDir("../builtin/logical")
		if err != nil {
			t.Fatal(err)
		}

		var backends []string
		for _, f := range files {
			if f.IsDir() {
				if f.Name() == "plugin" {
					continue
				}
				backends = append(backends, f.Name())
			}
		}

		plugins, err := ioutil.ReadDir("../vendor/github.com/hashicorp")
		if err != nil {
			t.Fatal(err)
		}
		for _, p := range plugins {
			if p.IsDir() && strings.HasPrefix(p.Name(), "vault-plugin-secrets-") {
				backends = append(backends, strings.TrimPrefix(p.Name(), "vault-plugin-secrets-"))
			}
		}

		// backends are found by walking the directory, which includes the database backend,
		// however, the plugins registry omits that one
		if len(backends) != len(builtinplugins.Registry.Keys(consts.PluginTypeSecrets))+1 {
			t.Fatalf("expected %d logical backends, got %d", len(builtinplugins.Registry.Keys(consts.PluginTypeSecrets))+1, len(backends))
		}

		for _, b := range backends {
			ui, cmd := testSecretsEnableCommand(t)
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
