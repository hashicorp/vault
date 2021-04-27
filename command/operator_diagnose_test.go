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
				"-config", "./server/test-fixtures/config_diagnose_ok.hcl",
			},
			[]string{"Parse configuration\n\x1b[F\x1b[32m[  ok  ]\x1b[0m Check Listeners\n\x1b[F\x1b[32m[  ok  ]\x1b[0m Access storage\n\x1b[F\x1b[32m[  ok  ]\x1b[0m Service discovery\n"},
			0,
		},
		{
			"diagnose_invalid_storage",
			[]string{
				"-config", "./server/test-fixtures/nostore_config.hcl",
			},
			[]string{"Parse configuration\n\x1b[F\x1b[32m[  ok  ]\x1b[0m Check Listeners\n\x1b[F\x1b[31m[failed]\x1b[0m Access storage\nA storage backend must be specified"},
			1,
		},
		{
			"diagnose_listener_config_ok",
			[]string{
				"-config", "./server/test-fixtures/tls_config_ok.hcl",
			},
			[]string{"Parse configuration\n\x1b[F\x1b[32m[  ok  ]\x1b[0m Check Listeners\n\x1b[F\x1b[32m[  ok  ]\x1b[0m Access storage\n\x1b[F\x1b[32m[  ok  ]\x1b[0m Service discovery"},
			0,
		},
		{
			"diagnose_invalid_https_storage",
			[]string{
				"-config", "./server/test-fixtures/config_bad_https_storage.hcl",
			},
			[]string{"Access storage\nfailed to verify certificate: x509: certificate has expired or is not yet valid:"},
			1,
		},
		{
			"diagnose_invalid_https_hastorage",
			[]string{
				"-config", "./server/test-fixtures/config_diagnose_hastorage_bad_https.hcl",
			},
			[]string{"Access storage\nfailed to verify certificate: x509: certificate has expired or is not yet valid:"},
			1,
		},
		{
			"diagnose_invalid_https_sr",
			[]string{
				"-config", "./server/test-fixtures/diagnose_bad_https_consul_sr.hcl",
			},
			[]string{"Service discovery\nfailed to verify certificate: x509: certificate has expired or is not yet valid:"},
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
