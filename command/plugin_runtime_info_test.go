// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"regexp"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func testPluginRuntimeInfoCommand(tb testing.TB) (*cli.MockUi, *PluginRuntimeInfoCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &PluginRuntimeInfoCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func TestPluginRuntimeInfoCommand_Run(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		args []string
		out  string
		code int
	}{
		{
			"not_enough_args",
			[]string{"-type=container"},
			"Not enough arguments",
			1,
		},
		{
			"too_many_args",
			[]string{"-type=container", "bar", "baz"},
			"Too many arguments",
			1,
		},
		{
			"invalid_runtime_type",
			[]string{"-type=foo", "bar"},
			"\"foo\" is not a supported plugin runtime type",
			2,
		},
		{
			"info_container_on_empty_plugin_runtime_catalog",
			[]string{"-type=container", "my-plugin-runtime"},
			"Error reading plugin runtime named my-plugin-runtime",
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

				ui, cmd := testPluginRuntimeInfoCommand(t)
				cmd.client = client

				code := cmd.Run(tc.args)
				if code != tc.code {
					t.Errorf("expected %d to be %d", code, tc.code)
				}

				combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
				matcher := regexp.MustCompile(tc.out)
				if !matcher.MatchString(combined) {
					t.Errorf("expected %q to contain %q", combined, tc.out)
				}
			})
		}
	})

	t.Run("communication_failure", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServerBad(t)
		defer closer()

		ui, cmd := testPluginRuntimeInfoCommand(t)
		cmd.client = client

		code := cmd.Run([]string{"-type=container", "my-plugin-runtime"})
		if exp := 2; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Error reading plugin runtime named my-plugin-runtime"
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}
	})

	t.Run("no_tabs", func(t *testing.T) {
		t.Parallel()

		_, cmd := testPluginRuntimeInfoCommand(t)
		assertNoTabs(t, cmd)
	})
}
