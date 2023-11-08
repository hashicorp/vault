// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"regexp"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func testPluginRuntimeDeregisterCommand(tb testing.TB) (*cli.MockUi, *PluginRuntimeDeregisterCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &PluginRuntimeDeregisterCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func TestPluginRuntimeDeregisterCommand_Run(t *testing.T) {
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
			[]string{"-type=container", "foo", "baz"},
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
			"Error deregistering plugin runtime named my-plugin-runtime",
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

				ui, cmd := testPluginRuntimeDeregisterCommand(t)
				cmd.ApiClient = client

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

		ui, cmd := testPluginRuntimeDeregisterCommand(t)
		cmd.ApiClient = client

		code := cmd.Run([]string{"-type=container", "my-plugin-runtime"})
		if exp := 2; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Error deregistering plugin runtime named my-plugin-runtime"
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}
	})

	t.Run("no_tabs", func(t *testing.T) {
		t.Parallel()

		_, cmd := testPluginRuntimeDeregisterCommand(t)
		assertNoTabs(t, cmd)
	})
}
