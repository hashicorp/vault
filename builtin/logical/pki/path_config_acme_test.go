// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/constants"
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
				accountKey, err := rsa.GenerateKey(rand.Reader, 2048)
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
