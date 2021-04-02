// +build !race

package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func testOperatorDiagnoseCommand(tb testing.TB) (*cli.MockUi, *OperatorDiagnoseCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &OperatorDiagnoseCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func TestOperatorDiagnoseCommand_Run(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name         string
		args         []string
		outFragments []string
		code         int
	}{
		{
			"diagnose_ok",
			[]string{
				"-config", "./server/test-fixtures/config.hcl",
			},
			[]string{"Parse configuration\n\x1b[F\x1b[32m[  ok  ]\x1b[0m Parse configuration\n[      ] Access storage\n\x1b[F\x1b[32m[  ok  ]\x1b[0m Access storage\n"},
			0,
		},
		{
			"diagnose_invalid_storage",
			[]string{
				"-config", "./server/test-fixtures/nostore_config.hcl",
			},
			[]string{"Parse configuration\n\x1b[F\x1b[32m[  ok  ]\x1b[0m Parse configuration\n[      ] Access storage\n\x1b[F\x1b[31m[failed]\x1b[0m Access storage\nA storage backend must be specified\n"},
			1,
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

				ui, cmd := testOperatorDiagnoseCommand(t)
				cmd.client = client

				code := cmd.Run(tc.args)
				if code != tc.code {
					t.Errorf("%s: expected %d to be %d", tc.name, code, tc.code)
				}

				combined := ui.OutputWriter.String() + ui.ErrorWriter.String()

				for _, outputFragment := range tc.outFragments {
					if !strings.Contains(combined, outputFragment) {
						t.Errorf("%s: expected %q to contain %q", tc.name, combined, outputFragment)
					}
				}
			})
		}
	})
}
