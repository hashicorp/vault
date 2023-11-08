// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"regexp"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func testPluginRuntimeListCommand(tb testing.TB) (*cli.MockUi, *PluginRuntimeListCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &PluginRuntimeListCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func TestPluginRuntimeListCommand_Run(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		args []string
		out  string
		code int
	}{
		{
			"too_many_args",
			[]string{"foo"},
			"Too many arguments",
			1,
		},
		{
			"invalid_runtime_type",
			[]string{"-type=foo"},
			"\"foo\" is not a supported plugin runtime type",
			2,
		},
		{
			"list container on empty plugin runtime catalog",
			[]string{"-type=container"},
			"Error listing available plugin runtimes:",
			2,
		},
		{
			"list on empty plugin runtime catalog",
			nil,
			"Error listing available plugin runtimes:",
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

				ui, cmd := testPluginRuntimeListCommand(t)
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

		ui, cmd := testPluginRuntimeListCommand(t)
		cmd.ApiClient = client

		code := cmd.Run([]string{"-type=container"})
		if exp := 2; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Error listing available plugin runtimes: "
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}
	})

	t.Run("no_tabs", func(t *testing.T) {
		t.Parallel()

		_, cmd := testPluginRuntimeListCommand(t)
		assertNoTabs(t, cmd)
	})
}
