// +build !race

package command

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/vault/vault/diagnose"
	"github.com/mitchellh/cli"
)

func testOperatorDiagnoseCommand(tb testing.TB) *OperatorDiagnoseCommand {
	tb.Helper()

	ui := cli.NewMockUi()
	return &OperatorDiagnoseCommand{
		diagnose: diagnose.New(),
		BaseCommand: &BaseCommand{
			UI: ui,
		},
		skipEndEnd: true,
	}
}

func TestOperatorDiagnoseCommand_Run(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name     string
		args     []string
		expected []*diagnose.Result
	}{
		{
			"diagnose_ok",
			[]string{
				"-config", "./server/test-fixtures/config_diagnose_ok.hcl",
			},
			[]*diagnose.Result{
				{
					Name:   "parse-config",
					Status: diagnose.OkStatus,
				},
				{
					Name:   "setup-telemetry",
					Status: diagnose.OkStatus,
				},
				{
					Name:   "storage",
					Status: diagnose.OkStatus,
					Children: []*diagnose.Result{
						{
							Name:   "create-storage-backend",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "test-storage-tls-consul",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "test-consul-direct-access-storage",
							Status: diagnose.OkStatus,
						},
					},
				},
				{
					Name:   "service-discovery",
					Status: diagnose.OkStatus,
					Children: []*diagnose.Result{
						{
							Name:   "test-serviceregistration-tls-consul",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "test-consul-direct-access-service-discovery",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "init-consul-serviceregistration",
							Status: diagnose.OkStatus,
						},
					},
				},
				{
					Name:   "create-seal",
					Status: diagnose.OkStatus,
				},
				{
					Name:   "setup-core",
					Status: diagnose.OkStatus,
					Children: []*diagnose.Result{
						{
							Name:   "init-randreader",
							Status: diagnose.OkStatus,
						},
					},
				},
				{
					Name:   "setup-ha-storage",
					Status: diagnose.OkStatus,
					Children: []*diagnose.Result{
						{
							Name:   "create-ha-storage-backend",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "test-consul-direct-access-storage",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "test-storage-tls-consul",
							Status: diagnose.OkStatus,
						},
					},
				},
				{
					Name:   "determine-redirect",
					Status: diagnose.OkStatus,
				},
				{
					Name:   "find-cluster-addr",
					Status: diagnose.OkStatus,
				},
				{
					Name:   "init-core",
					Status: diagnose.ErrorStatus,
				},
				{
					Name:   "init-listeners",
					Status: diagnose.OkStatus,
					Children: []*diagnose.Result{
						{
							Name:   "create-listeners",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "check-listener-tls",
							Status: diagnose.WarningStatus,
							Warnings: []string{
								"TLS is disabled in a Listener config stanza.",
							},
						},
					},
				},
				{
					Name:    "unseal",
					Status:  diagnose.ErrorStatus,
					Message: "core could not be initialized",
				},
				{
					Name:   "run-listeners",
					Status: diagnose.OkStatus,
				},
				{
					Name:    "start-servers",
					Status:  diagnose.ErrorStatus,
					Message: CoreUninitializedErr,
				},
			},
		},
		{
			"diagnose_invalid_storage",
			[]string{
				"-config", "./server/test-fixtures/nostore_config.hcl",
			},
			[]*diagnose.Result{
				{
					Name:   "parse-config",
					Status: diagnose.OkStatus,
				},
				{
					Name:   "setup-telemetry",
					Status: diagnose.OkStatus,
				},
				{
					Name:    "storage",
					Status:  diagnose.ErrorStatus,
					Message: "no storage stanza found in config",
					Children: []*diagnose.Result{
						{
							Name:   "create-storage-backend",
							Status: diagnose.ErrorStatus,
						},
					},
				},
				{
					Name:   "service-discovery",
					Status: diagnose.ErrorStatus,
				},
				{
					Name:   "create-seal",
					Status: diagnose.OkStatus,
				},
				{
					Name:    "setup-core",
					Status:  diagnose.ErrorStatus,
					Message: BackendUninitializedErr,
					Children: []*diagnose.Result{
						{
							Name:   "init-randreader",
							Status: diagnose.OkStatus,
						},
					},
				},
				{
					Name:    "setup-ha-storage",
					Status:  diagnose.ErrorStatus,
					Message: BackendUninitializedErr,
				},
				{
					Name:   "determine-redirect",
					Status: diagnose.OkStatus,
				},
				{
					Name:   "find-cluster-addr",
					Status: diagnose.OkStatus,
				},
				{
					Name:   "init-core",
					Status: diagnose.ErrorStatus,
				},
				{
					Name:   "init-listeners",
					Status: diagnose.OkStatus,
					Children: []*diagnose.Result{
						{
							Name:   "create-listeners",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "check-listener-tls",
							Status: diagnose.WarningStatus,
							Warnings: []string{
								"TLS is disabled in a Listener config stanza.",
							},
						},
					},
				},
				{
					Name:    "unseal",
					Status:  diagnose.ErrorStatus,
					Message: "core could not be initialized",
				},
				{
					Name:   "run-listeners",
					Status: diagnose.OkStatus,
				},
				{
					Name:    "start-servers",
					Status:  diagnose.ErrorStatus,
					Message: CoreUninitializedErr,
				},
			},
		},
		{
			"diagnose_listener_config_ok",
			[]string{
				"-config", "./server/test-fixtures/tls_config_ok.hcl",
			},
			[]*diagnose.Result{
				{
					Name:   "parse-config",
					Status: diagnose.OkStatus,
				},
				{
					Name:   "setup-telemetry",
					Status: diagnose.OkStatus,
				},
				{
					Name:   "storage",
					Status: diagnose.OkStatus,
					Children: []*diagnose.Result{
						{
							Name:   "create-storage-backend",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "test-storage-tls-consul",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "test-consul-direct-access-storage",
							Status: diagnose.OkStatus,
						},
					},
				},
				{
					Name:   "service-discovery",
					Status: diagnose.OkStatus,
					Children: []*diagnose.Result{
						{
							Name:   "test-serviceregistration-tls-consul",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "test-consul-direct-access-service-discovery",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "init-consul-serviceregistration",
							Status: diagnose.OkStatus,
						},
					},
				},
				{
					Name:   "create-seal",
					Status: diagnose.OkStatus,
				},
				{
					Name:   "setup-core",
					Status: diagnose.OkStatus,
					Children: []*diagnose.Result{
						{
							Name:   "init-randreader",
							Status: diagnose.OkStatus,
						},
					},
				},
				{
					Name:   "setup-ha-storage",
					Status: diagnose.OkStatus,
					Children: []*diagnose.Result{
						{
							Name:   "create-ha-storage-backend",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "test-consul-direct-access-storage",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "test-storage-tls-consul",
							Status: diagnose.OkStatus,
						},
					},
				},
				{
					Name:   "determine-redirect",
					Status: diagnose.OkStatus,
				},
				{
					Name:   "find-cluster-addr",
					Status: diagnose.OkStatus,
				},
				{
					Name:   "init-core",
					Status: diagnose.ErrorStatus,
				},
				{
					Name:   "init-listeners",
					Status: diagnose.OkStatus,
					Children: []*diagnose.Result{
						{
							Name:   "create-listeners",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "check-listener-tls",
							Status: diagnose.WarningStatus,
							Warnings: []string{
								"TLS is disabled in a Listener config stanza.",
							},
						},
					},
				},
				{
					Name:    "unseal",
					Status:  diagnose.ErrorStatus,
					Message: "core could not be initialized",
				},
				{
					Name:   "run-listeners",
					Status: diagnose.OkStatus,
				},
				{
					Name:    "start-servers",
					Status:  diagnose.ErrorStatus,
					Message: CoreUninitializedErr,
				},
			},
		},
		{
			"diagnose_invalid_https_storage",
			[]string{
				"-config", "./server/test-fixtures/config_bad_https_storage.hcl",
			},
			[]*diagnose.Result{
				{
					Name:   "parse-config",
					Status: diagnose.OkStatus,
				},
				{
					Name:   "setup-telemetry",
					Status: diagnose.OkStatus,
				},
				{
					Name:   "storage",
					Status: diagnose.ErrorStatus,
					Children: []*diagnose.Result{
						{
							Name:   "create-storage-backend",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "test-storage-tls-consul",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "test-consul-direct-access-storage",
							Status: diagnose.OkStatus,
						},
					},
				},
				{
					Name:   "service-discovery",
					Status: diagnose.OkStatus,
					Children: []*diagnose.Result{
						{
							Name:   "test-serviceregistration-tls-consul",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "test-consul-direct-access-service-discovery",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "init-consul-serviceregistration",
							Status: diagnose.OkStatus,
						},
					},
				},
				{
					Name:   "create-seal",
					Status: diagnose.OkStatus,
				},
				{
					Name:   "setup-core",
					Status: diagnose.OkStatus,
					Children: []*diagnose.Result{
						{
							Name:   "init-randreader",
							Status: diagnose.OkStatus,
						},
					},
				},
				{
					Name:   "setup-ha-storage",
					Status: diagnose.OkStatus,
					Children: []*diagnose.Result{
						{
							Name:   "create-ha-storage-backend",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "test-consul-direct-access-storage",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "test-storage-tls-consul",
							Status: diagnose.OkStatus,
						},
					},
				},
			},
		},
		{
			"diagnose_invalid_https_hastorage",
			[]string{
				"-config", "./server/test-fixtures/config_diagnose_hastorage_bad_https.hcl",
			},
			[]*diagnose.Result{
				{
					Name:   "parse-config",
					Status: diagnose.OkStatus,
				},
				{
					Name:   "setup-telemetry",
					Status: diagnose.OkStatus,
				},
				{
					Name:   "storage",
					Status: diagnose.ErrorStatus,
					Children: []*diagnose.Result{
						{
							Name:   "create-storage-backend",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "test-storage-tls-consul",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "test-consul-direct-access-storage",
							Status: diagnose.OkStatus,
						},
					},
				},
			},
		},
		{
			"diagnose_invalid_https_sr",
			[]string{
				"-config", "./server/test-fixtures/diagnose_bad_https_consul_sr.hcl",
			},
			[]*diagnose.Result{
				{
					Name:   "parse-config",
					Status: diagnose.OkStatus,
				},
				{
					Name:   "setup-telemetry",
					Status: diagnose.OkStatus,
				},
				{
					Name:   "storage",
					Status: diagnose.OkStatus,
					Children: []*diagnose.Result{
						{
							Name:   "create-storage-backend",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "test-storage-tls-consul",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "test-consul-direct-access-storage",
							Status: diagnose.OkStatus,
						},
					},
				},
				{
					Name:    "service-discovery",
					Status:  diagnose.ErrorStatus,
					Message: "failed to verify certificate: x509: certificate has expired or is not yet valid",
					Children: []*diagnose.Result{
						{
							Name:    "test-serviceregistration-tls-consul",
							Status:  diagnose.ErrorStatus,
							Message: "failed to verify certificate: x509: certificate has expired or is not yet valid",
						},
						{
							Name:   "test-consul-direct-access-service-discovery",
							Status: diagnose.WarningStatus,
							Warnings: []string{
								diagnose.DirAccessErr,
							},
						},
						{
							Name:   "init-consul-serviceregistration",
							Status: diagnose.OkStatus,
						},
					},
				},
			},
		},
		{
			"diagnose_direct_storage_access",
			[]string{
				"-config", "./server/test-fixtures/diagnose_ok_storage_direct_access.hcl",
			},
			[]*diagnose.Result{
				{
					Name:   "parse-config",
					Status: diagnose.OkStatus,
				},
				{
					Name:   "setup-telemetry",
					Status: diagnose.OkStatus,
				},
				{
					Name:   "storage",
					Status: diagnose.OkStatus,
					Children: []*diagnose.Result{
						{
							Name:   "create-storage-backend",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "test-storage-tls-consul",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "test-consul-direct-access-storage",
							Status: diagnose.WarningStatus,
							Warnings: []string{
								diagnose.DirAccessErr,
							},
						},
					},
				},
				{
					Name:   "service-discovery",
					Status: diagnose.OkStatus,
					Children: []*diagnose.Result{
						{
							Name:   "test-serviceregistration-tls-consul",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "test-consul-direct-access-service-discovery",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "init-consul-serviceregistration",
							Status: diagnose.OkStatus,
						},
					},
				},
				{
					Name:   "create-seal",
					Status: diagnose.OkStatus,
				},
				{
					Name:   "setup-core",
					Status: diagnose.OkStatus,
					Children: []*diagnose.Result{
						{
							Name:   "init-randreader",
							Status: diagnose.OkStatus,
						},
					},
				},
				{
					Name:   "setup-ha-storage",
					Status: diagnose.OkStatus,
					Children: []*diagnose.Result{
						{
							Name:   "create-ha-storage-backend",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "test-consul-direct-access-storage",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "test-storage-tls-consul",
							Status: diagnose.OkStatus,
						},
					},
				},
				{
					Name:   "determine-redirect",
					Status: diagnose.OkStatus,
				},
				{
					Name:   "find-cluster-addr",
					Status: diagnose.OkStatus,
				},
				{
					Name:   "init-core",
					Status: diagnose.ErrorStatus,
				},
				{
					Name:   "init-listeners",
					Status: diagnose.OkStatus,
					Children: []*diagnose.Result{
						{
							Name:   "create-listeners",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "check-listener-tls",
							Status: diagnose.WarningStatus,
							Warnings: []string{
								"TLS is disabled in a Listener config stanza.",
							},
						},
					},
				},
				{
					Name:    "unseal",
					Status:  diagnose.ErrorStatus,
					Message: "core could not be initialized",
				},
				{
					Name:   "run-listeners",
					Status: diagnose.OkStatus,
				},
				{
					Name:    "start-servers",
					Status:  diagnose.ErrorStatus,
					Message: CoreUninitializedErr,
				},
			},
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

				cmd := testOperatorDiagnoseCommand(t)
				cmd.client = client

				cmd.Run(tc.args)
				result := cmd.diagnose.Finalize(context.Background())

				for i, exp := range tc.expected {
					if i >= len(result.Children) {
						t.Fatalf("there are at least %d test cases, but fewer actual results", i)
					}
					act := result.Children[i]
					if err := compareResult(t, exp, act); err != nil {
						t.Fatalf("%v", err)
					}
				}
			})
		}
	})
}

func compareResult(t *testing.T, exp *diagnose.Result, act *diagnose.Result) error {
	if exp.Name != act.Name {
		return fmt.Errorf("names mismatch: %s vs %s", exp.Name, act.Name)
	}
	if exp.Status != act.Status {
		if act.Status != diagnose.OkStatus {
			return fmt.Errorf("section %s, status mismatch: %s vs %s, got error %s", exp.Name, exp.Status, act.Status, act.Message)

		}
		return fmt.Errorf("section %s, status mismatch: %s vs %s", exp.Name, exp.Status, act.Status)
	}
	if exp.Message != "" && exp.Message != act.Message && !strings.Contains(act.Message, exp.Message) {
		return fmt.Errorf("section %s, message not found: %s in %s", exp.Name, exp.Message, act.Message)
	}
	if len(exp.Warnings) != len(act.Warnings) {
		return fmt.Errorf("section %s, warning count mismatch: %d vs %d", exp.Name, len(exp.Warnings), len(act.Warnings))
	}
	for j := range exp.Warnings {
		if !strings.Contains(act.Warnings[j], exp.Warnings[j]) {
			return fmt.Errorf("section %s, warning message not found: %s in %s", exp.Name, exp.Warnings[j], act.Warnings[j])
		}
	}
	if len(exp.Children) != len(act.Children) {
		return fmt.Errorf("section %s, child count mismatch: %d vs %d", exp.Name, len(exp.Children), len(act.Children))
	}
	return nil
}
