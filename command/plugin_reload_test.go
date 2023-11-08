// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"strings"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	"github.com/mitchellh/cli"
)

func testPluginReloadCommand(tb testing.TB) (*cli.MockUi, *PluginReloadCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &PluginReloadCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func testPluginReloadStatusCommand(tb testing.TB) (*cli.MockUi, *PluginReloadStatusCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &PluginReloadStatusCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func TestPluginReloadCommand_Run(t *testing.T) {
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
			[]string{"-plugin", "foo", "-mounts", "bar"},
			"Too many arguments",
			1,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			client, closer := testVaultServer(t)
			defer closer()

			ui, cmd := testPluginReloadCommand(t)
			cmd.ApiClient = client

			args := append([]string{}, tc.args...)
			code := cmd.Run(args)
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

		pluginDir, cleanup := corehelpers.MakeTestPluginDir(t)
		defer cleanup(t)

		client, _, closer := testVaultServerPluginDir(t, pluginDir)
		defer closer()

		pluginName := "my-plugin"
		_, sha256Sum := testPluginCreateAndRegister(t, client, pluginDir, pluginName, api.PluginTypeCredential, "")

		ui, cmd := testPluginReloadCommand(t)
		cmd.ApiClient = client

		if err := client.Sys().RegisterPlugin(&api.RegisterPluginInput{
			Name:    pluginName,
			Type:    api.PluginTypeCredential,
			Command: pluginName,
			SHA256:  sha256Sum,
		}); err != nil {
			t.Fatal(err)
		}

		code := cmd.Run([]string{
			"-plugin", pluginName,
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Success! Reloaded plugin: "
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}
	})
}

func TestPluginReloadStatusCommand_Run(t *testing.T) {
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
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			client, closer := testVaultServer(t)
			defer closer()

			ui, cmd := testPluginReloadCommand(t)
			cmd.ApiClient = client

			args := append([]string{}, tc.args...)
			code := cmd.Run(args)
			if code != tc.code {
				t.Errorf("expected %d to be %d", code, tc.code)
			}

			combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
			if !strings.Contains(combined, tc.out) {
				t.Errorf("expected %q to contain %q", combined, tc.out)
			}
		})
	}
}
