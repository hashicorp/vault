// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"strings"
	"testing"

	"github.com/hashicorp/cli"
	base "github.com/hashicorp/vault/command/base"
	clientcmd "github.com/hashicorp/vault/command/client"
	"github.com/hashicorp/vault/sdk/helper/consts"
)

func testPluginRuntimeRegisterCommand(tb testing.TB) (*cli.MockUi, *clientcmd.PluginRuntimeRegisterCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &clientcmd.PluginRuntimeRegisterCommand{
		BaseCommand: &base.BaseCommand{
			UI: ui,
		},
	}
}

func TestPluginRuntimeRegisterCommand_Run(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		flags []string
		args  []string
		out   string
		code  int
	}{
		{
			"no type specified",
			[]string{},
			[]string{"foo"},
			"-type is required for plugin runtime registration",
			1,
		},
		{
			"invalid type",
			[]string{"-type", "foo"},
			[]string{"not"},
			"\"foo\" is not a supported plugin runtime type",
			2,
		},
		{
			"not_enough_args",
			[]string{"-type", consts.PluginRuntimeTypeContainer.String()},
			[]string{},
			"Not enough arguments",
			1,
		},
		{
			"too_many_args",
			[]string{"-type", consts.PluginRuntimeTypeContainer.String()},
			[]string{"foo", "bar"},
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

			ui, cmd := testPluginRuntimeRegisterCommand(t)
			cmd.SetClient(client)

			args := append(tc.flags, tc.args...)
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

	t.Run("communication_failure", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServerBad(t)
		defer closer()

		ui, cmd := testPluginRuntimeRegisterCommand(t)
		cmd.SetClient(client)

		code := cmd.Run([]string{"-type", consts.PluginRuntimeTypeContainer.String(), "my-plugin-runtime"})
		if exp := 2; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Error registering plugin runtime my-plugin-runtime"
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}
	})

	t.Run("no_tabs", func(t *testing.T) {
		t.Parallel()

		_, cmd := testPluginRuntimeRegisterCommand(t)
		assertNoTabs(t, cmd)
	})
}
