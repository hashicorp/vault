// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"context"
	"crypto/rand"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/constants"
	"github.com/hashicorp/vault/sdk/helper/cryptoutil"
	"github.com/stretchr/testify/require"
)

func TestAcmeConfig(t *testing.T) {
	t.Parallel()

	cluster, client, _ := setupAcmeBackend(t)
	defer cluster.Cleanup()

	cases := []struct {
		name        string
		AcmeConfig  map[string]interface{}
		prefixUrl   string
		validConfig bool
		works       bool
	}{
		{"unspecified-root", map[string]interface{}{
			"enabled":         true,
			"allowed_issuers": "*",
			"allowed_roles":   "*",
			"dns_resolver":    "",
			"eab_policy_name": "",
		}, "acme/", true, true},
		{"bad-policy-root", map[string]interface{}{
			"enabled":                  true,
			"allowed_issuers":          "*",
			"allowed_roles":            "*",
			"default_directory_policy": "bad",
			"dns_resolver":             "",
			"eab_policy_name":          "",
		}, "acme/", false, false},
		{"forbid-root", map[string]interface{}{
			"enabled":                  true,
			"allowed_issuers":          "*",
			"allowed_roles":            "*",
			"default_directory_policy": "forbid",
			"dns_resolver":             "",
			"eab_policy_name":          "",
		}, "acme/", true, false},
		{"sign-verbatim-root", map[string]interface{}{
			"enabled":                  true,
			"allowed_issuers":          "*",
			"allowed_roles":            "*",
			"default_directory_policy": "sign-verbatim",
			"dns_resolver":             "",
			"eab_policy_name":          "",
		}, "acme/", true, true},
		{"role-root", map[string]interface{}{
			"enabled":                  true,
			"allowed_issuers":          "*",
			"allowed_roles":            "*",
			"default_directory_policy": "role:exists",
			"dns_resolver":             "",
			"eab_policy_name":          "",
		}, "acme/", true, true},
		{"bad-role-root", map[string]interface{}{
			"enabled":                  true,
			"allowed_issuers":          "*",
			"allowed_roles":            "*",
			"default_directory_policy": "role:notgood",
			"dns_resolver":             "",
			"eab_policy_name":          "",
		}, "acme/", false, true},
		{"disallowed-role-root", map[string]interface{}{
			"enabled":                  true,
			"allowed_issuers":          "*",
			"allowed_roles":            "good",
			"default_directory_policy": "role:exists",
			"dns_resolver":             "",
			"eab_policy_name":          "",
		}, "acme/", false, false},
	}

	roleConfig := map[string]interface{}{
		"issuer_ref":       "default",
		"allowed_domains":  "example.com",
		"allow_subdomains": true,
		"max_ttl":          "720h",
	}

	testCtx := context.Background()

	for _, tc := range cases {
		deadline := time.Now().Add(1 * time.Minute)
		subTestCtx, _ := context.WithDeadline(testCtx, deadline)

		_, err := client.Logical().WriteWithContext(subTestCtx, "pki/roles/exists", roleConfig)
		require.NoError(t, err)
		_, err = client.Logical().WriteWithContext(subTestCtx, "pki/roles/good", roleConfig)
		require.NoError(t, err)

		t.Run(tc.name, func(t *testing.T) {
			_, err := client.Logical().WriteWithContext(subTestCtx, "pki/config/acme", tc.AcmeConfig)

			if tc.validConfig {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				return
			}

			_, err = client.Logical().ReadWithContext(subTestCtx, "pki/acme/directory")
			if tc.works {
				require.NoError(t, err)

				baseAcmeURL := "/v1/pki/" + tc.prefixUrl
				accountKey, err := cryptoutil.GenerateRSAKey(rand.Reader, 2048)
				require.NoError(t, err, "failed creating rsa key")

				acmeClient := getAcmeClientForCluster(t, cluster, baseAcmeURL, accountKey)

				// Create new account
				_, err = acmeClient.Discover(subTestCtx)
				require.NoError(t, err, "failed acme discovery call")
			} else {
				require.Error(t, err, "Acme Configuration should prevent usage")
			}
		})
	}
}

// TestAcmeExternalPolicyOss make sure setting external-policy on OSS within acme configuration fails
func TestAcmeExternalPolicyOss(t *testing.T) {
	if constants.IsEnterprise {
		t.Skip("this test is only valid on OSS")
	}

	t.Parallel()
	b, s := CreateBackendWithStorage(t)

	values := []string{"external-policy", "external-policy:", "external-policy:test"}
	for _, value := range values {
		t.Run(value, func(st *testing.T) {
			_, err := CBWrite(b, s, "config/acme", map[string]interface{}{
				"enabled":                  true,
				"default_directory_policy": value,
			})

			require.Error(st, err, "should have failed setting acme config")
		})
	}
}

// TestAcmeIPRangeConfiguration tests setting and modifying challenge IP range configuration
func TestAcmeIPRangeConfiguration(t *testing.T) {
	t.Parallel()

	cluster, client, _ := setupAcmeBackend(t)
	defer cluster.Cleanup()

	testCtx := context.Background()

	_, err := client.Logical().WriteWithContext(testCtx, "pki/config/acme", map[string]interface{}{
		"enabled": true,
	})
	require.NoError(t, err)

	cases := []struct {
		name                    string
		writeConfig             map[string]interface{}
		expectErrorContains     string
		expectedPermittedRanges []string
		expectedExcludedRanges  []string
	}{
		{
			name: "set_permitted_ip_ranges",
			writeConfig: map[string]interface{}{
				"challenge_permitted_ip_ranges": []string{"192.168.1.0/24", "10.0.0.0/8"},
			},
			expectedPermittedRanges: []string{"192.168.1.0/24", "10.0.0.0/8"},
			expectedExcludedRanges:  []string{},
		},
		{
			name: "set_excluded_ip_ranges",
			writeConfig: map[string]interface{}{
				"challenge_permitted_ip_ranges": []string{},
				"challenge_excluded_ip_ranges":  []string{"127.0.0.0/8", "::1/128"},
			},
			expectedPermittedRanges: []string{},
			expectedExcludedRanges:  []string{"127.0.0.0/8", "::1/128"},
		},
		{
			name: "set_both_ranges",
			writeConfig: map[string]interface{}{
				"challenge_permitted_ip_ranges": []string{"192.168.0.0/16"},
				"challenge_excluded_ip_ranges":  []string{"192.168.1.0/24"},
			},
			expectedPermittedRanges: []string{"192.168.0.0/16"},
			expectedExcludedRanges:  []string{"192.168.1.0/24"},
		},
		{
			name: "set_individual_ips",
			writeConfig: map[string]interface{}{
				"challenge_permitted_ip_ranges": []string{"192.168.1.100"},
				"challenge_excluded_ip_ranges":  []string{"10.0.0.1"},
			},
			expectedPermittedRanges: []string{"192.168.1.100"},
			expectedExcludedRanges:  []string{"10.0.0.1"},
		},
		{
			name: "invalid_cidr_notation",
			writeConfig: map[string]interface{}{
				"challenge_permitted_ip_ranges": []string{"invalid-cidr"},
			},
			expectErrorContains: "invalid CIDR or IP address",
		},
		{
			name: "invalid_excluded_ip",
			writeConfig: map[string]interface{}{
				"challenge_excluded_ip_ranges": []string{"999.999.999.999"},
			},
			expectErrorContains: "invalid CIDR or IP address",
		},
		{
			name: "modify_existing_config",
			writeConfig: map[string]interface{}{
				"challenge_excluded_ip_ranges": []string{"10.0.0.0/8"},
			},
			expectedPermittedRanges: []string{"192.168.1.100"},
			expectedExcludedRanges:  []string{"10.0.0.0/8"},
		},
		{
			name: "clear_ranges",
			writeConfig: map[string]interface{}{
				"challenge_permitted_ip_ranges": []string{},
				"challenge_excluded_ip_ranges":  []string{},
			},
			expectedPermittedRanges: []string{},
			expectedExcludedRanges:  []string{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(st *testing.T) {
			deadline := time.Now().Add(1 * time.Minute)
			subTestCtx, _ := context.WithDeadline(testCtx, deadline)

			_, err := client.Logical().WriteWithContext(subTestCtx, "pki/config/acme", tc.writeConfig)

			if tc.expectErrorContains != "" {
				require.Contains(st, err.Error(), tc.expectErrorContains)
				return
			}

			require.NoError(st, err)

			// Read back and verify
			resp, err := client.Logical().ReadWithContext(subTestCtx, "pki/config/acme")
			require.NoError(st, err)
			require.NotNil(st, resp)

			var permittedRanges []interface{}
			if resp.Data["challenge_permitted_ip_ranges"] != nil {
				permittedRanges = resp.Data["challenge_permitted_ip_ranges"].([]interface{})
			}

			var excludedRanges []interface{}
			if resp.Data["challenge_excluded_ip_ranges"] != nil {
				excludedRanges = resp.Data["challenge_excluded_ip_ranges"].([]interface{})
			}

			// Verify permitted ranges
			require.Len(st, permittedRanges, len(tc.expectedPermittedRanges))
			for _, expected := range tc.expectedPermittedRanges {
				require.Contains(st, permittedRanges, expected)
			}

			// Verify excluded ranges
			require.Len(st, excludedRanges, len(tc.expectedExcludedRanges))
			for _, expected := range tc.expectedExcludedRanges {
				require.Contains(st, excludedRanges, expected)
			}
		})
	}
}
