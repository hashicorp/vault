// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !race

package command

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hashicorp/cli"
	"github.com/hashicorp/vault/helper/constants"
	pkihelper "github.com/hashicorp/vault/helper/testhelpers/pki"
	"github.com/hashicorp/vault/vault/diagnose"
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

func generateTLSConfigOk(t *testing.T, ca pkihelper.LeafWithIntermediary) string {
	t.Helper()
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "tls_config_ok.hcl")

	templateFile := "./server/test-fixtures/tls_config_ok.hcl"
	contents, err := os.ReadFile(templateFile)
	if err != nil {
		t.Fatalf("failed to read file %s: %v", templateFile, err)
	}
	contents = []byte(strings.ReplaceAll(string(contents), "{REPLACE_LEAF_CERT_FILE}", ca.Leaf.CertFile))
	contents = []byte(strings.ReplaceAll(string(contents), "{REPLACE_LEAF_KEY_FILE}", ca.Leaf.KeyFile))

	err = os.WriteFile(configPath, contents, 0o644)
	if err != nil {
		t.Fatalf("failed to write file %s: %v", configPath, err)
	}

	return configPath
}

func generateTransitTLSCheck(t *testing.T, ca pkihelper.LeafWithIntermediary) string {
	t.Helper()
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "diagnose_seal_transit_tls_check.hcl")

	templateFile := "./server/test-fixtures/diagnose_seal_transit_tls_check.hcl"
	contents, err := os.ReadFile(templateFile)
	if err != nil {
		t.Fatalf("failed to read file %s: %v", templateFile, err)
	}
	contents = []byte(strings.ReplaceAll(string(contents), "{REPLACE_LEAF_CERT_FILE}", ca.Leaf.CertFile))
	contents = []byte(strings.ReplaceAll(string(contents), "{REPLACE_LEAF_KEY_FILE}", ca.Leaf.KeyFile))
	contents = []byte(strings.ReplaceAll(string(contents), "{REPLACE_COMBINED_CA_CHAIN_FILE}", ca.CombinedCaFile))

	err = os.WriteFile(configPath, contents, 0o644)
	if err != nil {
		t.Fatalf("failed to write file %s: %v", configPath, err)
	}

	return configPath
}

func TestOperatorDiagnoseCommand_Run(t *testing.T) {
	t.Parallel()
	testca := pkihelper.GenerateCertWithIntermediaryRoot(t)
	tlsConfigOkConfigFile := generateTLSConfigOk(t, testca)
	transitTLSCheckConfigFile := generateTransitTLSCheck(t, testca)

	cases := []struct {
		name     string
		args     []string
		expected []*diagnose.Result
	}{
		{
			"diagnose_ok",
			[]string{
				"-config", "./server/test-fixtures/config_diagnose_ok_singleseal.hcl",
			},
			[]*diagnose.Result{
				{
					Name:   "Parse Configuration",
					Status: diagnose.OkStatus,
				},
				{
					Name:   "Start Listeners",
					Status: diagnose.WarningStatus,
					Children: []*diagnose.Result{
						{
							Name:   "Create Listeners",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "Check Listener TLS",
							Status: diagnose.WarningStatus,
							Warnings: []string{
								"TLS is disabled in a listener config stanza.",
							},
						},
					},
				},
				{
					Name:   "Check Storage",
					Status: diagnose.OkStatus,
					Children: []*diagnose.Result{
						{
							Name:   "Create Storage Backend",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "Check Consul TLS",
							Status: diagnose.SkippedStatus,
						},
						{
							Name:   "Check Consul Direct Storage Access",
							Status: diagnose.OkStatus,
						},
					},
				},
				{
					Name:   "Check Service Discovery",
					Status: diagnose.OkStatus,
					Children: []*diagnose.Result{
						{
							Name:   "Check Consul Service Discovery TLS",
							Status: diagnose.SkippedStatus,
						},
						{
							Name:   "Check Consul Direct Service Discovery",
							Status: diagnose.OkStatus,
						},
					},
				},
				{
					Name:   "Create Vault Server Configuration Seals",
					Status: diagnose.OkStatus,
				},
				{
					Name:   "Create Core Configuration",
					Status: diagnose.OkStatus,
					Children: []*diagnose.Result{
						{
							Name:   "Initialize Randomness for Core",
							Status: diagnose.OkStatus,
						},
					},
				},
				{
					Name:   "HA Storage",
					Status: diagnose.OkStatus,
					Children: []*diagnose.Result{
						{
							Name:   "Create HA Storage Backend",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "Check HA Consul Direct Storage Access",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "Check Consul TLS",
							Status: diagnose.SkippedStatus,
						},
					},
				},
				{
					Name:   "Determine Redirect Address",
					Status: diagnose.OkStatus,
				},
				{
					Name:   "Check Cluster Address",
					Status: diagnose.OkStatus,
				},
				{
					Name:   "Check Core Creation",
					Status: diagnose.OkStatus,
				},
				{
					Name:   "Start Listeners",
					Status: diagnose.WarningStatus,
					Children: []*diagnose.Result{
						{
							Name:   "Create Listeners",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "Check Listener TLS",
							Status: diagnose.WarningStatus,
							Warnings: []string{
								"TLS is disabled in a listener config stanza.",
							},
						},
					},
				},
				{
					Name:    "Check Autounseal Encryption",
					Status:  diagnose.SkippedStatus,
					Message: "Skipping barrier encryption",
				},
				{
					Name:   "Check Server Before Runtime",
					Status: diagnose.OkStatus,
				},
				{
					Name:   "Finalize Shamir Seal",
					Status: diagnose.OkStatus,
				},
			},
		},
		{
			"diagnose_ok_multiseal",
			[]string{
				"-config", "./server/test-fixtures/config_diagnose_ok.hcl",
			},
			[]*diagnose.Result{
				{
					Name:   "Parse Configuration",
					Status: diagnose.OkStatus,
				},
				{
					Name:   "Start Listeners",
					Status: diagnose.WarningStatus,
					Children: []*diagnose.Result{
						{
							Name:   "Create Listeners",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "Check Listener TLS",
							Status: diagnose.WarningStatus,
							Warnings: []string{
								"TLS is disabled in a listener config stanza.",
							},
						},
					},
				},
				{
					Name:   "Check Storage",
					Status: diagnose.OkStatus,
					Children: []*diagnose.Result{
						{
							Name:   "Create Storage Backend",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "Check Consul TLS",
							Status: diagnose.SkippedStatus,
						},
						{
							Name:   "Check Consul Direct Storage Access",
							Status: diagnose.OkStatus,
						},
					},
				},
				{
					Name:   "Check Service Discovery",
					Status: diagnose.OkStatus,
					Children: []*diagnose.Result{
						{
							Name:   "Check Consul Service Discovery TLS",
							Status: diagnose.SkippedStatus,
						},
						{
							Name:   "Check Consul Direct Service Discovery",
							Status: diagnose.OkStatus,
						},
					},
				},
				{
					Name: "Create Vault Server Configuration Seals",
					// We can't load from storage the existing seal generation info during the test, so we expect an error.
					Status: diagnose.ErrorStatus,
				},
				{
					Name:   "Create Core Configuration",
					Status: diagnose.OkStatus,
					Children: []*diagnose.Result{
						{
							Name:   "Initialize Randomness for Core",
							Status: diagnose.OkStatus,
						},
					},
				},
				{
					Name:   "HA Storage",
					Status: diagnose.OkStatus,
					Children: []*diagnose.Result{
						{
							Name:   "Create HA Storage Backend",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "Check HA Consul Direct Storage Access",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "Check Consul TLS",
							Status: diagnose.SkippedStatus,
						},
					},
				},
				{
					Name:   "Determine Redirect Address",
					Status: diagnose.OkStatus,
				},
				{
					Name:   "Check Cluster Address",
					Status: diagnose.OkStatus,
				},
				{
					Name:   "Check Core Creation",
					Status: diagnose.OkStatus,
				},
				{
					Name:   "Start Listeners",
					Status: diagnose.WarningStatus,
					Children: []*diagnose.Result{
						{
							Name:   "Create Listeners",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "Check Listener TLS",
							Status: diagnose.WarningStatus,
							Warnings: []string{
								"TLS is disabled in a listener config stanza.",
							},
						},
					},
				},
				{
					Name:    "Check Autounseal Encryption",
					Status:  diagnose.ErrorStatus,
					Message: "Diagnose could not create a barrier seal object.",
				},
				{
					Name:   "Check Server Before Runtime",
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
					Name:   "Check Storage",
					Status: diagnose.WarningStatus,
					Children: []*diagnose.Result{
						{
							Name:   "Create Storage Backend",
							Status: diagnose.OkStatus,
						},
						{
							Name:    "Check Raft Folder Permissions",
							Status:  diagnose.WarningStatus,
							Message: "too many permissions",
						},
						{
							Name:    "Check For Raft Quorum",
							Status:  diagnose.WarningStatus,
							Message: "0 voters found",
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
					Name:    "Check Storage",
					Status:  diagnose.ErrorStatus,
					Message: "No storage stanza in Vault server configuration.",
				},
			},
		},
		{
			"diagnose_listener_config_ok",
			[]string{
				"-config", tlsConfigOkConfigFile,
			},
			[]*diagnose.Result{
				{
					Name:   "Start Listeners",
					Status: diagnose.OkStatus,
					Children: []*diagnose.Result{
						{
							Name:   "Create Listeners",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "Check Listener TLS",
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
					Name:   "Check Storage",
					Status: diagnose.ErrorStatus,
					Children: []*diagnose.Result{
						{
							Name:   "Create Storage Backend",
							Status: diagnose.OkStatus,
						},
						{
							Name:    "Check Consul TLS",
							Status:  diagnose.ErrorStatus,
							Message: "certificate has expired or is not yet valid",
							Warnings: []string{
								"expired or near expiry",
							},
						},
						{
							Name:   "Check Consul Direct Storage Access",
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
					Name:   "Check Storage",
					Status: diagnose.WarningStatus,
					Children: []*diagnose.Result{
						{
							Name:   "Create Storage Backend",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "Check Consul TLS",
							Status: diagnose.SkippedStatus,
						},
						{
							Name:   "Check Consul Direct Storage Access",
							Status: diagnose.WarningStatus,
							Advice: "We recommend connecting to a local agent.",
							Warnings: []string{
								"Vault storage is directly connected to a Consul server.",
							},
						},
					},
				},
				{
					Name:   "HA Storage",
					Status: diagnose.ErrorStatus,
					Children: []*diagnose.Result{
						{
							Name:   "Create HA Storage Backend",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "Check HA Consul Direct Storage Access",
							Status: diagnose.WarningStatus,
							Advice: "We recommend connecting to a local agent.",
							Warnings: []string{
								"Vault storage is directly connected to a Consul server.",
							},
						},
						{
							Name:    "Check Consul TLS",
							Status:  diagnose.ErrorStatus,
							Message: "certificate has expired or is not yet valid",
							Warnings: []string{
								"expired or near expiry",
							},
						},
					},
				},
				{
					Name:   "Check Cluster Address",
					Status: diagnose.ErrorStatus,
				},
			},
		},
		{
			"diagnose_seal_transit_tls_check_fail",
			[]string{
				"-config", transitTLSCheckConfigFile,
			},
			[]*diagnose.Result{
				{
					Name:   "Check Transit Seal TLS",
					Status: diagnose.WarningStatus,
					Warnings: []string{
						"Found at least one intermediate certificate in the CA certificate file.",
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
					Name:   "Check Service Discovery",
					Status: diagnose.ErrorStatus,
					Children: []*diagnose.Result{
						{
							Name:    "Check Consul Service Discovery TLS",
							Status:  diagnose.ErrorStatus,
							Message: "certificate has expired or is not yet valid",
							Warnings: []string{
								"expired or near expiry",
							},
						},
						{
							Name:   "Check Consul Direct Service Discovery",
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
					Name:   "Check Storage",
					Status: diagnose.WarningStatus,
					Children: []*diagnose.Result{
						{
							Name:   "Create Storage Backend",
							Status: diagnose.OkStatus,
						},
						{
							Name:   "Check Consul TLS",
							Status: diagnose.SkippedStatus,
						},
						{
							Name:   "Check Consul Direct Storage Access",
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
			"diagnose_raft_no_folder_backend",
			[]string{
				"-config", "./server/test-fixtures/diagnose_raft_no_bolt_folder.hcl",
			},
			[]*diagnose.Result{
				{
					Name:    "Check Storage",
					Status:  diagnose.ErrorStatus,
					Message: "Diagnose could not initialize storage backend.",
					Children: []*diagnose.Result{
						{
							Name:    "Create Storage Backend",
							Status:  diagnose.ErrorStatus,
							Message: "no such file or directory",
						},
					},
				},
			},
		},
		{
			"diagnose_telemetry_partial_circonus",
			[]string{
				"-config", "./server/test-fixtures/diagnose_bad_telemetry1.hcl",
			},
			[]*diagnose.Result{
				{
					Name:    "Check Telemetry",
					Status:  diagnose.ErrorStatus,
					Message: "incomplete Circonus telemetry configuration, missing circonus_api_url",
				},
			},
		},
		{
			"diagnose_telemetry_partial_dogstats",
			[]string{
				"-config", "./server/test-fixtures/diagnose_bad_telemetry2.hcl",
			},
			[]*diagnose.Result{
				{
					Name:    "Check Telemetry",
					Status:  diagnose.ErrorStatus,
					Message: "incomplete DogStatsD telemetry configuration, missing dogstatsd_addr, while dogstatsd_tags specified",
				},
			},
		},
		{
			"diagnose_telemetry_partial_stackdriver",
			[]string{
				"-config", "./server/test-fixtures/diagnose_bad_telemetry3.hcl",
			},
			[]*diagnose.Result{
				{
					Name:    "Check Telemetry",
					Status:  diagnose.ErrorStatus,
					Message: "incomplete Stackdriver telemetry configuration, missing stackdriver_project_id",
				},
			},
		},
		{
			"diagnose_telemetry_default",
			[]string{
				"-config", "./server/test-fixtures/config4.hcl",
			},
			[]*diagnose.Result{
				{
					Name:     "Check Telemetry",
					Status:   diagnose.WarningStatus,
					Warnings: []string{"Telemetry is using default configuration"},
				},
			},
		},
	}

	t.Run("validations", func(t *testing.T) {
		t.Parallel()

		for _, tc := range cases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				if tc.name == "diagnose_ok" && constants.IsEnterprise {
					t.Skip("Test not valid in ENT")
				} else if tc.name == "diagnose_ok_multiseal" && !constants.IsEnterprise {
					t.Skip("Test not valid in community edition")
				} else {
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
	if exp.Advice != "" && exp.Advice != act.Advice && !strings.Contains(act.Advice, exp.Advice) {
		return fmt.Errorf("section %s, advice not found: %s in %s", exp.Name, exp.Advice, act.Advice)
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
		return errors.New(strings.Join(errStrings, ","))
	}

	if len(exp.Children) > 0 {
		return compareResults(exp.Children, act.Children)
	}

	// Remove raft file if it exists
	os.Remove("./server/test-fixtures/vault.db")
	os.RemoveAll("./server/test-fixtures/raft")

	return nil
}
