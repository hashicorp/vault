// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/mitchellh/cli"
)

func testPluginRegisterCommand(tb testing.TB) (*cli.MockUi, *PluginRegisterCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &PluginRegisterCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func TestPluginRegisterCommand_Run(t *testing.T) {
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
			2,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			client, closer := testVaultServer(t)
			defer closer()

			ui, cmd := testPluginRegisterCommand(t)
			cmd.client = client

			args := append([]string{"-sha256", "abcd1234"}, tc.args...)
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
		_, sha256Sum := testPluginCreate(t, pluginDir, pluginName)

		ui, cmd := testPluginRegisterCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"-sha256", sha256Sum,
			consts.PluginTypeCredential.String(), pluginName,
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Success! Registered plugin: my-plugin"
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}

		resp, err := client.Sys().ListPlugins(&api.ListPluginsInput{
			Type: api.PluginTypeCredential,
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
		if !found {
			t.Errorf("expected %q to be in %q", pluginName, resp.PluginsByType)
		}
	})

	t.Run("integration with version", func(t *testing.T) {
		t.Parallel()

		pluginDir, cleanup := corehelpers.MakeTestPluginDir(t)
		defer cleanup(t)

		client, _, closer := testVaultServerPluginDir(t, pluginDir)
		defer closer()

		const pluginName = "my-plugin"
		versions := []string{"v1.0.0", "v2.0.1"}
		_, sha256Sum := testPluginCreate(t, pluginDir, pluginName)
		types := []api.PluginType{api.PluginTypeCredential, api.PluginTypeDatabase, api.PluginTypeSecrets}

		for _, typ := range types {
			for _, version := range versions {
				ui, cmd := testPluginRegisterCommand(t)
				cmd.client = client

				code := cmd.Run([]string{
					"-version=" + version,
					"-sha256=" + sha256Sum,
					typ.String(),
					pluginName,
				})
				if exp := 0; code != exp {
					t.Errorf("expected %d to be %d", code, exp)
				}

				expected := "Success! Registered plugin: my-plugin"
				combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
				if !strings.Contains(combined, expected) {
					t.Errorf("expected %q to contain %q", combined, expected)
				}
			}
		}

		resp, err := client.Sys().ListPlugins(&api.ListPluginsInput{
			Type: api.PluginTypeUnknown,
		})
		if err != nil {
			t.Fatal(err)
		}

		found := make(map[api.PluginType]int)
		versionsFound := make(map[api.PluginType][]string)
		for _, p := range resp.Details {
			if p.Name == pluginName {
				typ, err := api.ParsePluginType(p.Type)
				if err != nil {
					t.Fatal(err)
				}
				found[typ]++
				versionsFound[typ] = append(versionsFound[typ], p.Version)
			}
		}

		for _, typ := range types {
			if found[typ] != 2 {
				t.Fatalf("expected %q to be found 2 times, but found it %d times for %s type in %#v", pluginName, found[typ], typ.String(), resp.Details)
			}
			sort.Strings(versions)
			sort.Strings(versionsFound[typ])
			if !reflect.DeepEqual(versions, versionsFound[typ]) {
				t.Fatalf("expected %v versions but got %v", versions, versionsFound[typ])
			}
		}
	})

	t.Run("communication_failure", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServerBad(t)
		defer closer()

		ui, cmd := testPluginRegisterCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"-sha256", "abcd1234",
			consts.PluginTypeCredential.String(), "my-plugin",
		})
		if exp := 2; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Error registering plugin my-plugin:"
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}
	})

	t.Run("no_tabs", func(t *testing.T) {
		t.Parallel()

		_, cmd := testPluginRegisterCommand(t)
		assertNoTabs(t, cmd)
	})
}
