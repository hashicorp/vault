// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"strings"
	"testing"

	"github.com/hashicorp/cli"
)

func testEventsSubscriptionsCommand(tb testing.TB) (*cli.MockUi, *EventsSubscriptionCreateCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &EventsSubscriptionCreateCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

// TestEventsSubscriptionsCommand_Run tests that the command argument parsing is working as expected.
func TestEventsSubscriptionsCommand_Run(t *testing.T) {
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
			"not_enough_args",
			[]string{"foo"},
			"Not enough arguments",
			1,
		},
		{
			"too_many_args",
			[]string{"foo", "bar", "baz"},
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

			ui, cmd := testEventsSubscriptionsCommand(t)
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
}
