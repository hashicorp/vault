package command

import (
	"strings"
	"testing"

	"github.com/hashicorp/vault/helper/consts"
	"github.com/mitchellh/cli"
)

func testPluginInfoCommand(tb testing.TB) (*cli.MockUi, *PluginInfoCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &PluginInfoCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func TestPluginInfoCommand_Run(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		args []string
		out  string
		code int
	}{
		{
			"too_many_args",
			[]string{"foo", "bar", "fizz"},
			"Too many arguments",
			1,
		},
		{
			"no_plugin_exist",
			[]string{consts.PluginTypeCredential.String(), "not-a-real-plugin-like-ever"},
			"Error reading plugin",
			2,
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

				ui, cmd := testPluginInfoCommand(t)
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

		pluginDir, cleanup := testPluginDir(t)
		defer cleanup(t)

		client, _, closer := testVaultServerPluginDir(t, pluginDir)
		defer closer()

		pluginName := "my-plugin"
		_, sha256Sum := testPluginCreateAndRegister(t, client, pluginDir, pluginName, consts.PluginTypeCredential)

		ui, cmd := testPluginInfoCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			consts.PluginTypeCredential.String(), pluginName,
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, pluginName) {
			t.Errorf("expected %q to contain %q", combined, pluginName)
		}
		if !strings.Contains(combined, sha256Sum) {
			t.Errorf("expected %q to contain %q", combined, sha256Sum)
		}
	})

	t.Run("field", func(t *testing.T) {
		t.Parallel()

		pluginDir, cleanup := testPluginDir(t)
		defer cleanup(t)

		client, _, closer := testVaultServerPluginDir(t, pluginDir)
		defer closer()

		pluginName := "my-plugin"
		testPluginCreateAndRegister(t, client, pluginDir, pluginName, consts.PluginTypeCredential)

		ui, cmd := testPluginInfoCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"-field", "builtin",
			consts.PluginTypeCredential.String(), pluginName,
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if exp := "false"; combined != exp {
			t.Errorf("expected %q to be %q", combined, exp)
		}
	})

	t.Run("communication_failure", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServerBad(t)
		defer closer()

		ui, cmd := testPluginInfoCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			consts.PluginTypeCredential.String(), "my-plugin",
		})
		if exp := 2; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Error reading plugin named my-plugin: "
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}
	})

	t.Run("no_tabs", func(t *testing.T) {
		t.Parallel()

		_, cmd := testPluginInfoCommand(t)
		assertNoTabs(t, cmd)
	})
}
