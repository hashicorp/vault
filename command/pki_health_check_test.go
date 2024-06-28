// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/cli"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/healthcheck"
	"github.com/stretchr/testify/require"
)

func TestPKIHC_AllGood(t *testing.T) {
	t.Parallel()

	client, closer := testVaultServer(t)
	defer closer()

	if err := client.Sys().Mount("pki", &api.MountInput{
		Type: "pki",
		Config: api.MountConfigInput{
			AuditNonHMACRequestKeys:   healthcheck.VisibleReqParams,
			AuditNonHMACResponseKeys:  healthcheck.VisibleRespParams,
			PassthroughRequestHeaders: []string{"If-Modified-Since"},
			AllowedResponseHeaders:    []string{"Last-Modified", "Replay-Nonce", "Link", "Location"},
			MaxLeaseTTL:               "36500d",
		},
	}); err != nil {
		t.Fatalf("pki mount error: %#v", err)
	}

	if resp, err := client.Logical().Write("pki/root/generate/internal", map[string]interface{}{
		"key_type":    "ec",
		"common_name": "Root X1",
		"ttl":         "3650d",
	}); err != nil || resp == nil {
		t.Fatalf("failed to prime CA: %v", err)
	}

	if _, err := client.Logical().Read("pki/crl/rotate"); err != nil {
		t.Fatalf("failed to rotate CRLs: %v", err)
	}

	if _, err := client.Logical().Write("pki/roles/testing", map[string]interface{}{
		"allow_any_name": true,
		"no_store":       true,
	}); err != nil {
		t.Fatalf("failed to write role: %v", err)
	}

	if _, err := client.Logical().Write("pki/config/auto-tidy", map[string]interface{}{
		"enabled":         true,
		"tidy_cert_store": true,
	}); err != nil {
		t.Fatalf("failed to write auto-tidy config: %v", err)
	}

	if _, err := client.Logical().Write("pki/tidy", map[string]interface{}{
		"tidy_cert_store": true,
	}); err != nil {
		t.Fatalf("failed to run tidy: %v", err)
	}

	path, err := url.Parse(client.Address())
	require.NoError(t, err, "failed parsing client address")

	if _, err := client.Logical().Write("pki/config/cluster", map[string]interface{}{
		"path": path.JoinPath("/v1/", "pki/").String(),
	}); err != nil {
		t.Fatalf("failed to update local cluster: %v", err)
	}

	if _, err := client.Logical().Write("pki/config/acme", map[string]interface{}{
		"enabled": "true",
	}); err != nil {
		t.Fatalf("failed to update acme config: %v", err)
	}

	_, _, results := execPKIHC(t, client, true)

	validateExpectedPKIHC(t, expectedAllGood, results)
}

func TestPKIHC_AllBad(t *testing.T) {
	t.Parallel()

	client, closer := testVaultServer(t)
	defer closer()

	if err := client.Sys().Mount("pki", &api.MountInput{
		Type: "pki",
	}); err != nil {
		t.Fatalf("pki mount error: %#v", err)
	}

	if resp, err := client.Logical().Write("pki/root/generate/internal", map[string]interface{}{
		"key_type":    "ec",
		"common_name": "Root X1",
		"ttl":         "35d",
	}); err != nil || resp == nil {
		t.Fatalf("failed to prime CA: %v", err)
	}

	if _, err := client.Logical().Write("pki/config/crl", map[string]interface{}{
		"expiry": "5s",
	}); err != nil {
		t.Fatalf("failed to issue leaf cert: %v", err)
	}

	if _, err := client.Logical().Read("pki/crl/rotate"); err != nil {
		t.Fatalf("failed to rotate CRLs: %v", err)
	}

	time.Sleep(5 * time.Second)

	if _, err := client.Logical().Write("pki/roles/testing", map[string]interface{}{
		"allow_localhost":             true,
		"allowed_domains":             "*.example.com",
		"allow_glob_domains":          true,
		"allow_wildcard_certificates": true,
		"no_store":                    false,
		"key_type":                    "ec",
		"ttl":                         "30d",
	}); err != nil {
		t.Fatalf("failed to write role: %v", err)
	}

	if _, err := client.Logical().Write("pki/issue/testing", map[string]interface{}{
		"common_name": "something.example.com",
	}); err != nil {
		t.Fatalf("failed to issue leaf cert: %v", err)
	}

	if _, err := client.Logical().Write("pki/config/auto-tidy", map[string]interface{}{
		"enabled":         false,
		"tidy_cert_store": false,
	}); err != nil {
		t.Fatalf("failed to write auto-tidy config: %v", err)
	}

	_, _, results := execPKIHC(t, client, true)

	validateExpectedPKIHC(t, expectedAllBad, results)
}

func TestPKIHC_OnlyIssuer(t *testing.T) {
	t.Parallel()

	client, closer := testVaultServer(t)
	defer closer()

	if err := client.Sys().Mount("pki", &api.MountInput{
		Type: "pki",
	}); err != nil {
		t.Fatalf("pki mount error: %#v", err)
	}

	if resp, err := client.Logical().Write("pki/root/generate/internal", map[string]interface{}{
		"key_type":    "ec",
		"common_name": "Root X1",
		"ttl":         "35d",
	}); err != nil || resp == nil {
		t.Fatalf("failed to prime CA: %v", err)
	}

	_, _, results := execPKIHC(t, client, true)
	validateExpectedPKIHC(t, expectedEmptyWithIssuer, results)
}

func TestPKIHC_NoMount(t *testing.T) {
	t.Parallel()

	client, closer := testVaultServer(t)
	defer closer()

	code, message, _ := execPKIHC(t, client, false)
	if code != 1 {
		t.Fatalf("Expected return code 1 from invocation on non-existent mount, got %v\nOutput: %v", code, message)
	}

	if !strings.Contains(message, "route entry not found") {
		t.Fatalf("Expected failure to talk about missing route entry, got exit code %v\nOutput: %v", code, message)
	}
}

func TestPKIHC_ExpectedEmptyMount(t *testing.T) {
	t.Parallel()

	client, closer := testVaultServer(t)
	defer closer()

	if err := client.Sys().Mount("pki", &api.MountInput{
		Type: "pki",
	}); err != nil {
		t.Fatalf("pki mount error: %#v", err)
	}

	code, message, _ := execPKIHC(t, client, false)
	if code != 1 {
		t.Fatalf("Expected return code 1 from invocation on empty mount, got %v\nOutput: %v", code, message)
	}

	if !strings.Contains(message, "lacks any configured issuers,") {
		t.Fatalf("Expected failure to talk about no issuers, got exit code %v\nOutput: %v", code, message)
	}
}

func TestPKIHC_NoPerm(t *testing.T) {
	t.Parallel()

	client, closer := testVaultServer(t)
	defer closer()

	if err := client.Sys().Mount("pki", &api.MountInput{
		Type: "pki",
	}); err != nil {
		t.Fatalf("pki mount error: %#v", err)
	}

	if resp, err := client.Logical().Write("pki/root/generate/internal", map[string]interface{}{
		"key_type":    "ec",
		"common_name": "Root X1",
		"ttl":         "35d",
	}); err != nil || resp == nil {
		t.Fatalf("failed to prime CA: %v", err)
	}

	if _, err := client.Logical().Write("pki/config/crl", map[string]interface{}{
		"expiry": "5s",
	}); err != nil {
		t.Fatalf("failed to issue leaf cert: %v", err)
	}

	if _, err := client.Logical().Read("pki/crl/rotate"); err != nil {
		t.Fatalf("failed to rotate CRLs: %v", err)
	}

	time.Sleep(5 * time.Second)

	if _, err := client.Logical().Write("pki/roles/testing", map[string]interface{}{
		"allow_localhost":             true,
		"allowed_domains":             "*.example.com",
		"allow_glob_domains":          true,
		"allow_wildcard_certificates": true,
		"no_store":                    false,
		"key_type":                    "ec",
		"ttl":                         "30d",
	}); err != nil {
		t.Fatalf("failed to write role: %v", err)
	}

	if _, err := client.Logical().Write("pki/issue/testing", map[string]interface{}{
		"common_name": "something.example.com",
	}); err != nil {
		t.Fatalf("failed to issue leaf cert: %v", err)
	}

	if _, err := client.Logical().Write("pki/config/auto-tidy", map[string]interface{}{
		"enabled":         false,
		"tidy_cert_store": false,
	}); err != nil {
		t.Fatalf("failed to write auto-tidy config: %v", err)
	}

	// Remove client token.
	client.ClearToken()

	_, _, results := execPKIHC(t, client, true)
	validateExpectedPKIHC(t, expectedNoPerm, results)
}

func testPKIHealthCheckCommand(tb testing.TB) (*cli.MockUi, *PKIHealthCheckCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &PKIHealthCheckCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func execPKIHC(t *testing.T, client *api.Client, ok bool) (int, string, map[string][]map[string]interface{}) {
	t.Helper()

	stdout := bytes.NewBuffer(nil)
	stderr := bytes.NewBuffer(nil)
	runOpts := &RunOptions{
		Stdout: stdout,
		Stderr: stderr,
		Client: client,
	}

	code := RunCustom([]string{"pki", "health-check", "-format=json", "pki"}, runOpts)
	combined := stdout.String() + stderr.String()

	var results map[string][]map[string]interface{}
	if err := json.Unmarshal([]byte(combined), &results); err != nil {
		if ok {
			t.Fatalf("failed to decode json (ret %v): %v\njson:\n%v", code, err, combined)
		}
	}

	t.Log(combined)

	return code, combined, results
}

func validateExpectedPKIHC(t *testing.T, expected, results map[string][]map[string]interface{}) {
	t.Helper()

	for test, subtest := range expected {
		actual, ok := results[test]
		require.True(t, ok, fmt.Sprintf("expected top-level test %v to be present", test))

		if subtest == nil {
			continue
		}

		require.NotNil(t, actual, fmt.Sprintf("expected top-level test %v to be non-empty; wanted wireframe format %v", test, subtest))
		require.Equal(t, len(subtest), len(actual), fmt.Sprintf("top-level test %v has different number of results %v in wireframe, %v in test output\nwireframe: %v\noutput: %v\n", test, len(subtest), len(actual), subtest, actual))

		for index, subset := range subtest {
			for key, value := range subset {
				a_value, present := actual[index][key]
				require.True(t, present)
				if value != nil {
					require.Equal(t, value, a_value, fmt.Sprintf("in test: %v / result %v - when validating key %v\nWanted: %v\nGot: %v", test, index, key, subset, actual[index]))
				}
			}
		}
	}

	for name := range results {
		if _, present := expected[name]; !present {
			t.Fatalf("got unexpected health check: %v\n%v", name, results[name])
		}
	}
}

var expectedAllGood = map[string][]map[string]interface{}{
	"ca_validity_period": {
		{
			"status": "ok",
		},
	},
	"crl_validity_period": {
		{
			"status": "ok",
		},
		{
			"status": "ok",
		},
	},
	"allow_acme_headers": {
		{
			"status": "ok",
		},
	},
	"allow_if_modified_since": {
		{
			"status": "ok",
		},
	},
	"audit_visibility": {
		{
			"status": "ok",
		},
	},
	"enable_acme_issuance": {
		{
			"status": "ok",
		},
	},
	"enable_auto_tidy": {
		{
			"status": "ok",
		},
	},
	"role_allows_glob_wildcards": {
		{
			"status": "ok",
		},
	},
	"role_allows_localhost": {
		{
			"status": "ok",
		},
	},
	"role_no_store_false": {
		{
			"status": "ok",
		},
	},
	"root_issued_leaves": {
		{
			"status": "ok",
		},
	},
	"tidy_last_run": {
		{
			"status": "ok",
		},
	},
	"too_many_certs": {
		{
			"status": "ok",
		},
	},
}

var expectedAllBad = map[string][]map[string]interface{}{
	"ca_validity_period": {
		{
			"status": "critical",
		},
	},
	"crl_validity_period": {
		{
			"status": "critical",
		},
		{
			"status": "critical",
		},
	},
	"allow_acme_headers": {
		{
			"status": "not_applicable",
		},
	},
	"allow_if_modified_since": {
		{
			"status": "informational",
		},
	},
	"audit_visibility": {
		{
			"status": "informational",
		},
		{
			"status": "informational",
		},
		{
			"status": "informational",
		},
		{
			"status": "informational",
		},
		{
			"status": "informational",
		},
		{
			"status": "informational",
		},
		{
			"status": "informational",
		},
		{
			"status": "informational",
		},
		{
			"status": "informational",
		},
		{
			"status": "informational",
		},
		{
			"status": "informational",
		},
		{
			"status": "informational",
		},
		{
			"status": "informational",
		},
		{
			"status": "informational",
		},
		{
			"status": "informational",
		},
		{
			"status": "informational",
		},
		{
			"status": "informational",
		},
		{
			"status": "informational",
		},
		{
			"status": "informational",
		},
		{
			"status": "informational",
		},
		{
			"status": "informational",
		},
		{
			"status": "informational",
		},
		{
			"status": "informational",
		},
		{
			"status": "informational",
		},
		{
			"status": "informational",
		},
		{
			"status": "informational",
		},
		{
			"status": "informational",
		},
		{
			"status": "informational",
		},
		{
			"status": "informational",
		},
		{
			"status": "informational",
		},
	},
	"enable_acme_issuance": {
		{
			"status": "not_applicable",
		},
	},
	"enable_auto_tidy": {
		{
			"status": "informational",
		},
	},
	"role_allows_glob_wildcards": {
		{
			"status": "warning",
		},
	},
	"role_allows_localhost": {
		{
			"status": "warning",
		},
	},
	"role_no_store_false": {
		{
			"status": "warning",
		},
	},
	"root_issued_leaves": {
		{
			"status": "warning",
		},
	},
	"tidy_last_run": {
		{
			"status": "critical",
		},
	},
	"too_many_certs": {
		{
			"status": "ok",
		},
	},
}

var expectedEmptyWithIssuer = map[string][]map[string]interface{}{
	"ca_validity_period": {
		{
			"status": "critical",
		},
	},
	"crl_validity_period": {
		{
			"status": "ok",
		},
		{
			"status": "ok",
		},
	},
	"allow_acme_headers": {
		{
			"status": "not_applicable",
		},
	},
	"allow_if_modified_since": nil,
	"audit_visibility":        nil,
	"enable_acme_issuance": {
		{
			"status": "not_applicable",
		},
	},
	"enable_auto_tidy": {
		{
			"status": "informational",
		},
	},
	"role_allows_glob_wildcards": nil,
	"role_allows_localhost":      nil,
	"role_no_store_false":        nil,
	"root_issued_leaves": {
		{
			"status": "ok",
		},
	},
	"tidy_last_run": {
		{
			"status": "critical",
		},
	},
	"too_many_certs": {
		{
			"status": "ok",
		},
	},
}

var expectedNoPerm = map[string][]map[string]interface{}{
	"ca_validity_period": {
		{
			"status": "critical",
		},
	},
	"crl_validity_period": {
		{
			"status": "insufficient_permissions",
		},
		{
			"status": "critical",
		},
		{
			"status": "critical",
		},
	},
	"allow_acme_headers": {
		{
			"status": "insufficient_permissions",
		},
	},
	"allow_if_modified_since": nil,
	"audit_visibility":        nil,
	"enable_acme_issuance": {
		{
			"status": "insufficient_permissions",
		},
	},
	"enable_auto_tidy": {
		{
			"status": "insufficient_permissions",
		},
	},
	"role_allows_glob_wildcards": {
		{
			"status": "insufficient_permissions",
		},
	},
	"role_allows_localhost": {
		{
			"status": "insufficient_permissions",
		},
	},
	"role_no_store_false": {
		{
			"status": "insufficient_permissions",
		},
	},
	"root_issued_leaves": {
		{
			"status": "insufficient_permissions",
		},
	},
	"tidy_last_run": {
		{
			"status": "insufficient_permissions",
		},
	},
	"too_many_certs": {
		{
			"status": "insufficient_permissions",
		},
	},
}
