package command

import (
	"strings"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/consts"
	"github.com/mitchellh/cli"
)

func testPluginDeregisterCommand(tb testing.TB) (*cli.MockUi, *PluginDeregisterCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &PluginDeregisterCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func TestPluginDeregisterCommand_Run(t *testing.T) {
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
			[]string{"foo", "bar", "fizz"},
			"Too many arguments",
			1,
		},
		{
			"not_a_plugin",
			[]string{consts.PluginTypeCredential.String(), "nope_definitely_never_a_plugin_nope"},
			"",
			0,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			client, closer := testVaultServer(t)
			defer closer()

			ui, cmd := testPluginDeregisterCommand(t)
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

		pluginDir, cleanup := testPluginDir(t)
		defer cleanup(t)

		client, _, closer := testVaultServerPluginDir(t, pluginDir)
		defer closer()

		pluginName := "my-plugin"
		_, sha256Sum := testPluginCreateAndRegister(t, client, pluginDir, pluginName, consts.PluginTypeCredential)

		ui, cmd := testPluginDeregisterCommand(t)
		cmd.client = client

		if err := client.Sys().RegisterPlugin(&api.RegisterPluginInput{
			Name:    pluginName,
			Type:    consts.PluginTypeCredential,
			Command: pluginName,
			SHA256:  sha256Sum,
		}); err != nil {
			t.Fatal(err)
		}

		code := cmd.Run([]string{
			consts.PluginTypeCredential.String(),
			pluginName,
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Success! Deregistered plugin (if it was registered): "
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}

		resp, err := client.Sys().ListPlugins(&api.ListPluginsInput{
			Type: consts.PluginTypeCredential,
		})
		if err != nil {
			t.Fatal(err)
		}

		found := false
		for _, plugins := range resp.PluginsByType {
			for _, p := range plugins {
				if p == pluginName {
					found = true
				}
			}
		}
		if found {
			t.Errorf("expected %q to not be in %q", pluginName, resp.PluginsByType)
		}
	})

	t.Run("communication_failure", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServerBad(t)
		defer closer()

		ui, cmd := testPluginDeregisterCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			consts.PluginTypeCredential.String(),
			"my-plugin",
		})
		if exp := 2; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Error deregistering plugin named my-plugin: "
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}
	})

	t.Run("no_tabs", func(t *testing.T) {
		t.Parallel()

		_, cmd := testPluginDeregisterCommand(t)
		assertNoTabs(t, cmd)
	})
}
