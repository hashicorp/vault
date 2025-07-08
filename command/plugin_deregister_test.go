// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"strings"
	"testing"

	"github.com/hashicorp/cli"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	"github.com/hashicorp/vault/sdk/helper/consts"
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
			[]string{"foo"},
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

		pluginDir := corehelpers.MakeTestPluginDir(t)

		client, _, closer := testVaultServerPluginDir(t, pluginDir)
		defer closer()

		pluginName := "my-plugin"
		_, sha256Sum := testPluginCreateAndRegister(t, client, pluginDir, pluginName, api.PluginTypeCredential, "")

		ui, cmd := testPluginDeregisterCommand(t)
		cmd.client = client

		registerResp, err := client.Sys().RegisterPluginDetailed(&api.RegisterPluginInput{
			Name:    pluginName,
			Type:    api.PluginTypeCredential,
			Command: pluginName,
			SHA256:  sha256Sum,
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(registerResp.Warnings) > 0 {
			t.Errorf("expected no warnings, got %q", registerResp.Warnings)
		}

		code := cmd.Run([]string{
			consts.PluginTypeCredential.String(),
			pluginName,
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Success! Deregistered auth plugin: "
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}

		listResp, err := client.Sys().ListPlugins(&api.ListPluginsInput{
			Type: api.PluginTypeCredential,
		})
		if err != nil {
			t.Fatal(err)
		}

		found := false
		for _, plugins := range listResp.PluginsByType {
			for _, p := range plugins {
				if p == pluginName {
					found = true
				}
			}
		}
		if found {
			t.Errorf("expected %q to not be in %q", pluginName, listResp.PluginsByType)
		}
	})

	t.Run("integration with version", func(t *testing.T) {
		t.Parallel()

		pluginDir := corehelpers.MakeTestPluginDir(t)

		client, _, closer := testVaultServerPluginDir(t, pluginDir)
		defer closer()

		pluginName := "my-plugin"
		_, _, version := testPluginCreateAndRegisterVersioned(t, client, pluginDir, pluginName, api.PluginTypeCredential)

		ui, cmd := testPluginDeregisterCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"-version=" + version,
			consts.PluginTypeCredential.String(),
			pluginName,
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Success! Deregistered auth plugin: "
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}

		resp, err := client.Sys().ListPlugins(&api.ListPluginsInput{
			Type: api.PluginTypeUnknown,
		})
		if err != nil {
			t.Fatal(err)
		}

		found := false
		for _, p := range resp.Details {
			if p.Name == pluginName {
				found = true
			}
		}
		if found {
			t.Errorf("expected %q to not be in %#v", pluginName, resp.Details)
		}
	})

	t.Run("integration with missing version", func(t *testing.T) {
		t.Parallel()

		pluginDir := corehelpers.MakeTestPluginDir(t)

		client, _, closer := testVaultServerPluginDir(t, pluginDir)
		defer closer()

		pluginName := "my-plugin"
		testPluginCreateAndRegisterVersioned(t, client, pluginDir, pluginName, api.PluginTypeCredential)

		ui, cmd := testPluginDeregisterCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			consts.PluginTypeCredential.String(),
			pluginName,
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "does not exist in the catalog"
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}

		resp, err := client.Sys().ListPlugins(&api.ListPluginsInput{
			Type: api.PluginTypeUnknown,
		})
		if err != nil {
			t.Fatal(err)
		}

		found := false
		for _, p := range resp.Details {
			if p.Name == pluginName {
				found = true
			}
		}
		if !found {
			t.Errorf("expected %q to be in %#v", pluginName, resp.Details)
		}
	})

	t.Run("deregister builtin", func(t *testing.T) {
		t.Parallel()

		pluginDir := corehelpers.MakeTestPluginDir(t)

		client, _, closer := testVaultServerPluginDir(t, pluginDir)
		defer closer()

		ui, cmd := testPluginDeregisterCommand(t)
		cmd.client = client

		expected := "is a builtin plugin"
		if code := cmd.Run([]string{
			consts.PluginTypeCredential.String(),
			"github",
		}); code != 2 {
			t.Errorf("expected %d to be %d", code, 2)
		} else if !strings.Contains(ui.ErrorWriter.String(), expected) {
			t.Errorf("expected %q to contain %q", ui.ErrorWriter.String(), expected)
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
