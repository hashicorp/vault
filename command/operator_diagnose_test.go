// +build !race

package command

import (
	"github.com/go-test/deep"
	"github.com/hashicorp/vault/vault/diagnose"
	"reflect"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func testOperatorDiagnoseCommand(tb testing.TB) *OperatorDiagnoseCommand {
	tb.Helper()

	ui := cli.NewMockUi()
	return &OperatorDiagnoseCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func TestOperatorDiagnoseCommand_Run(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name     string
		args     []string
		expected diagnose.Result
		code     int
	}{
		{
			"diagnose_ok",
			[]string{
				"-config", "./server/test-fixtures/config_diagnose_ok.hcl",
			},
			diagnose.Result{
				Name:   "initialization",
				Status: diagnose.OkStatus,
				Children: []*diagnose.Result{
					{
						Name:   "parse-config",
						Status: diagnose.OkStatus,
					},
					{
						Name:   "init-listeners",
						Status: diagnose.WarningStatus,
						Warnings: []string{
							"TLS is disabled in a Listener config stanza.",
						},
					},
					{
						Name:   "storage",
						Status: diagnose.OkStatus,
					},
				},
			},
			0,
		},
		{
			"diagnose_invalid_storage",
			[]string{
				"-config", "./server/test-fixtures/nostore_config.hcl",
			},
			diagnose.Result{
				Name:   "initialization",
				Status: diagnose.OkStatus,
				Children: []*diagnose.Result{
					{
						Name:   "parse-config",
						Status: diagnose.OkStatus,
					},
					{
						Name:   "init-listeners",
						Status: diagnose.WarningStatus,
						Warnings: []string{
							"TLS is disabled in a Listener config stanza.",
						},
					},
					{
						Name:    "storage",
						Status:  diagnose.ErrorStatus,
						Message: "A storage backend must be specified",
					},
				},
			},
			1,
		},
		{
			"diagnose_listener_config_ok",
			[]string{
				"-config", "./server/test-fixtures/tls_config_ok.hcl",
			},
			diagnose.Result{
				Name:   "initialization",
				Status: diagnose.OkStatus,
				Children: []*diagnose.Result{
					{
						Name:   "parse-config",
						Status: diagnose.OkStatus,
					},
					{
						Name:   "init-listeners",
						Status: diagnose.OkStatus,
					},
					{
						Name:   "storage",
						Status: diagnose.OkStatus,
					},
				},
			},
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

				diagnose.Init()
				cmd := testOperatorDiagnoseCommand(t)
				cmd.client = client

				code := cmd.Run(tc.args)
				if code != tc.code {
					t.Errorf("%s: expected %d to be %d", tc.name, code, tc.code)
				}

				result := diagnose.Shutdown()
				result.ZeroTimes()
				if !reflect.DeepEqual(tc.expected, *result) {
					t.Fatalf("result mismatch: %s", strings.Join(deep.Equal(tc.expected, *result), "\n"))
				}
			})
		}
	})
}
