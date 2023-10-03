// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func testOperatorRaftSnapshotInspectCommand(tb testing.TB) (*cli.MockUi, *OperatorRaftSnapshotInspectCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &OperatorRaftSnapshotInspectCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func TestOperatorRaftSnapshotInspectCommand_Run(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		args []string
		out  string
		code int
	}{
		{
			"too_many_args",
			[]string{"./test-fixtures/test.snap", "test"},
			"Too many arguments",
			1,
		},
		{
			"default",
			[]string{"./test-fixtures/test.snap"},
			"ID           bolt-snapshot",
			0,
		},
		{
			"all_flags",
			[]string{"-kvdetails", "-kvdepth", "10", "-kvfilter", "core", "./test-fixtures/test.snap"},
			"Key Name                                              Count",
			0,
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

				ui, cmd := testOperatorRaftSnapshotInspectCommand(t)
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
}
