// +build !race

package command

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/hashicorp/vault/vault/diagnose"
	"github.com/mitchellh/cli"
)

func testOperatorDiagnoseCommand(tb testing.TB) *OperatorDiagnoseCommand {
	tb.Helper()

	ui := cli.NewMockUi()
	return &OperatorDiagnoseCommand{
		diagnose: diagnose.New(ioutil.Discard),
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
					Name:   "Parse configuration",
					Status: diagnose.OkStatus,
				},
				{
					Name:   "init-listeners",
					Status: diagnose.WarningStatus,
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
							Name:   "test-ha-storage-tls-consul",
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
					Status: diagnose.OkStatus,
				},
				{
					Name:   "init-listeners",
					Status: diagnose.WarningStatus,
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
					Message: "Diagnose could not create a barrier seal object",
				},
				{
					Name:   "start-servers",
					Status: diagnose.OkStatus,
				},
				{
					Name:   "finalize-seal-shamir",
					Status: diagnose.OkStatus,
				},
			},
		},
		{
			"diagnose_raft_problems",
			[]string{
				"-config", "./server/test-fixtures/config_raft.hcl",
			},
			[]*diagnose.Result{
				{
					Name:   "storage",
					Status: diagnose.ErrorStatus,
					Children: []*diagnose.Result{
						{
							Name:    "create-storage-backend",
							Status:  diagnose.ErrorStatus,
							Message: "failed to open bolt file",
						},
						{
							Name:    "raft folder permission checks",
							Status:  diagnose.WarningStatus,
							Message: "too many permissions",
						},
						{
							Name:    "raft quorum",
							Status:  diagnose.ErrorStatus,
							Message: "could not determine quorum status",
						},
					},
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
					Name:    "storage",
					Status:  diagnose.ErrorStatus,
					Message: "no storage stanza found in config",
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
					Name:   "init-listeners",
					Status: diagnose.OkStatus,
					Children: []*diagnose.Result{
						{
							Name:   "create-listeners",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "check-listener-tls",
							Status: diagnose.OkStatus,
						},
					},
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
					Name:   "storage",
					Status: diagnose.ErrorStatus,
					Children: []*diagnose.Result{
						{
							Name:   "create-storage-backend",
							Status: diagnose.OkStatus,
						},
						{
							Name:    "test-storage-tls-consul",
							Status:  diagnose.ErrorStatus,
							Message: "certificate has expired or is not yet valid",
							Warnings: []string{
								"expired or near expiry",
							},
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
			"diagnose_invalid_https_hastorage",
			[]string{
				"-config", "./server/test-fixtures/config_diagnose_hastorage_bad_https.hcl",
			},
			[]*diagnose.Result{
				{
					Name:   "storage",
					Status: diagnose.WarningStatus,
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
								"consul storage does not connect to local agent, but directly to server",
							},
						},
					},
				},
				{
					Name:   "setup-ha-storage",
					Status: diagnose.ErrorStatus,
					Children: []*diagnose.Result{
						{
							Name:   "create-ha-storage-backend",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "test-consul-direct-access-storage",
							Status: diagnose.WarningStatus,
							Warnings: []string{
								"consul storage does not connect to local agent, but directly to server",
							},
						},
						{
							Name:    "test-ha-storage-tls-consul",
							Status:  diagnose.ErrorStatus,
							Message: "certificate has expired or is not yet valid",
							Warnings: []string{
								"expired or near expiry",
							},
						},
					},
				},
				{
					Name:   "find-cluster-addr",
					Status: diagnose.ErrorStatus,
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
					Name:   "service-discovery",
					Status: diagnose.ErrorStatus,
					Children: []*diagnose.Result{
						{
							Name:    "test-serviceregistration-tls-consul",
							Status:  diagnose.ErrorStatus,
							Message: "certificate has expired or is not yet valid",
							Warnings: []string{
								"expired or near expiry",
							},
						},
						{
							Name:   "test-consul-direct-access-service-discovery",
							Status: diagnose.WarningStatus,
							Warnings: []string{
								diagnose.DirAccessErr,
							},
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
					Name:   "storage",
					Status: diagnose.WarningStatus,
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

				if err := compareResults(tc.expected, result.Children); err != nil {
					t.Fatalf("Did not find expected test results: %v", err)
				}
			})
		}
	})
}

func compareResults(expected []*diagnose.Result, actual []*diagnose.Result) error {
	for _, exp := range expected {
		found := false
		// Check them all so we don't have to be order specific
		for _, act := range actual {
			fmt.Printf("%+v", act)
			if exp.Name == act.Name {
				found = true
				if err := compareResult(exp, act); err != nil {
					return err
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("could not find expected test result: %s", exp.Name)
		}
	}
	return nil
}

func compareResult(exp *diagnose.Result, act *diagnose.Result) error {
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
	if len(exp.Children) > len(act.Children) {
		errStrings := []string{}
		for _, c := range act.Children {
			errStrings = append(errStrings, fmt.Sprintf("%+v", c))
		}
		return fmt.Errorf(strings.Join(errStrings, ","))
	}

	if len(exp.Children) > 0 {
		return compareResults(exp.Children, act.Children)
	}

	return nil
}
