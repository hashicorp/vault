// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package policy

import (
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/credential/ldap"
	credUserpass "github.com/hashicorp/vault/builtin/credential/userpass"
	ldaphelper "github.com/hashicorp/vault/helper/testhelpers/ldap"
	"github.com/hashicorp/vault/helper/testhelpers/minimal"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

func TestPolicy_NoDefaultPolicy(t *testing.T) {
	var err error
	coreConfig := &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"ldap": ldap.Factory,
		},
	}

	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})

	cluster.Start()
	defer cluster.Cleanup()

	cores := cluster.Cores

	vault.TestWaitActive(t, cores[0].Core)

	client := cores[0].Client

	err = client.Sys().EnableAuthWithOptions("ldap", &api.EnableAuthOptions{
		Type: "ldap",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Configure LDAP auth backend
	cleanup, cfg := ldaphelper.PrepareTestContainer(t, "master")
	defer cleanup()

	_, err = client.Logical().Write("auth/ldap/config", map[string]interface{}{
		"url":                     cfg.Url,
		"userattr":                cfg.UserAttr,
		"userdn":                  cfg.UserDN,
		"groupdn":                 cfg.GroupDN,
		"groupattr":               cfg.GroupAttr,
		"binddn":                  cfg.BindDN,
		"bindpass":                cfg.BindPassword,
		"token_no_default_policy": true,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create a local user in LDAP
	secret, err := client.Logical().Write("auth/ldap/users/hermes conrad", map[string]interface{}{
		"policies": "foo",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Login with LDAP and create a token
	secret, err = client.Logical().Write("auth/ldap/login/hermes conrad", map[string]interface{}{
		"password": "hermes",
	})
	if err != nil {
		t.Fatal(err)
	}
	token := secret.Auth.ClientToken

	// Lookup the token to get the entity ID
	secret, err = client.Auth().Token().Lookup(token)
	if err != nil {
		t.Fatal(err)
	}

	if diff := deep.Equal(secret.Data["policies"], []interface{}{"foo"}); diff != nil {
		t.Fatal(diff)
	}
}

func TestPolicy_NoConfiguredPolicy(t *testing.T) {
	var err error
	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	err = client.Sys().EnableAuthWithOptions("ldap", &api.EnableAuthOptions{
		Type: "ldap",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Configure LDAP auth backend
	cleanup, cfg := ldaphelper.PrepareTestContainer(t, "master")
	defer cleanup()

	_, err = client.Logical().Write("auth/ldap/config", map[string]interface{}{
		"url":       cfg.Url,
		"userattr":  cfg.UserAttr,
		"userdn":    cfg.UserDN,
		"groupdn":   cfg.GroupDN,
		"groupattr": cfg.GroupAttr,
		"binddn":    cfg.BindDN,
		"bindpass":  cfg.BindPassword,
		"token_ttl": "24h",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create a local user in LDAP without any policies configured
	secret, err := client.Logical().Write("auth/ldap/users/hermes conrad", map[string]interface{}{})
	if err != nil {
		t.Fatal(err)
	}

	// Login with LDAP and create a token
	secret, err = client.Logical().Write("auth/ldap/login/hermes conrad", map[string]interface{}{
		"password": "hermes",
	})
	if err != nil {
		t.Fatal(err)
	}
	token := secret.Auth.ClientToken

	// Lookup the token to get the entity ID
	secret, err = client.Auth().Token().Lookup(token)
	if err != nil {
		t.Fatal(err)
	}

	if diff := deep.Equal(secret.Data["policies"], []interface{}{"default"}); diff != nil {
		t.Fatal(diff)
	}

	// Renew the token with an increment of 2 hours to ensure that lease renewal
	// occurred and can be checked against the default lease duration with a
	// big enough delta.
	secret, err = client.Logical().Write("auth/token/renew", map[string]interface{}{
		"token":     token,
		"increment": "2h",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Verify that the lease renewal extended the duration properly.
	if float64(secret.Auth.LeaseDuration) < (1 * time.Hour).Seconds() {
		t.Fatalf("failed to renew lease, got: %v", secret.Auth.LeaseDuration)
	}
}

func TestPolicy_TokenRenewal(t *testing.T) {
	cases := []struct {
		name             string
		tokenPolicies    []string
		identityPolicies []string
	}{
		{
			"default only",
			nil,
			nil,
		},
		{
			"with token policies",
			[]string{"token-policy"},
			nil,
		},
		{
			"with identity policies",
			nil,
			[]string{"identity-policy"},
		},
		{
			"with token and identity policies",
			[]string{"token-policy"},
			[]string{"identity-policy"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			coreConfig := &vault.CoreConfig{
				CredentialBackends: map[string]logical.Factory{
					"userpass": credUserpass.Factory,
				},
			}
			cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
				HandlerFunc: vaulthttp.Handler,
			})
			cluster.Start()
			defer cluster.Cleanup()

			core := cluster.Cores[0].Core
			vault.TestWaitActive(t, core)
			client := cluster.Cores[0].Client

			// Enable userpass auth
			err := client.Sys().EnableAuthWithOptions("userpass", &api.EnableAuthOptions{
				Type: "userpass",
			})
			if err != nil {
				t.Fatal(err)
			}

			// Add a user to userpass backend
			data := map[string]interface{}{
				"password": "testpassword",
			}
			if len(tc.tokenPolicies) > 0 {
				data["token_policies"] = tc.tokenPolicies
			}
			_, err = client.Logical().Write("auth/userpass/users/testuser", data)
			if err != nil {
				t.Fatal(err)
			}

			// Set up entity if we're testing against an identity_policies
			if len(tc.identityPolicies) > 0 {
				auths, err := client.Sys().ListAuth()
				if err != nil {
					t.Fatal(err)
				}
				userpassAccessor := auths["userpass/"].Accessor

				resp, err := client.Logical().Write("identity/entity", map[string]interface{}{
					"name":     "test-entity",
					"policies": tc.identityPolicies,
				})
				if err != nil {
					t.Fatal(err)
				}
				entityID := resp.Data["id"].(string)

				// Create an alias
				resp, err = client.Logical().Write("identity/entity-alias", map[string]interface{}{
					"name":           "testuser",
					"mount_accessor": userpassAccessor,
					"canonical_id":   entityID,
				})
				if err != nil {
					t.Fatal(err)
				}
			}

			// Authenticate
			secret, err := client.Logical().Write("auth/userpass/login/testuser", map[string]interface{}{
				"password": "testpassword",
			})
			if err != nil {
				t.Fatal(err)
			}
			clientToken := secret.Auth.ClientToken

			// Verify the policies exist in the login response
			expectedTokenPolicies := append([]string{"default"}, tc.tokenPolicies...)
			if !strutil.EquivalentSlices(secret.Auth.TokenPolicies, expectedTokenPolicies) {
				t.Fatalf("token policy mismatch:\nexpected: %v\ngot: %v", expectedTokenPolicies, secret.Auth.TokenPolicies)
			}

			if !strutil.EquivalentSlices(secret.Auth.IdentityPolicies, tc.identityPolicies) {
				t.Fatalf("identity policy mismatch:\nexpected: %v\ngot: %v", tc.identityPolicies, secret.Auth.IdentityPolicies)
			}

			expectedPolicies := append(expectedTokenPolicies, tc.identityPolicies...)
			if !strutil.EquivalentSlices(secret.Auth.Policies, expectedPolicies) {
				t.Fatalf("policy mismatch:\nexpected: %v\ngot: %v", expectedPolicies, secret.Auth.Policies)
			}

			// Renew token
			secret, err = client.Logical().Write("auth/token/renew", map[string]interface{}{
				"token": clientToken,
			})
			if err != nil {
				t.Fatal(err)
			}

			// Verify the policies exist in the renewal response
			if !strutil.EquivalentSlices(secret.Auth.TokenPolicies, expectedTokenPolicies) {
				t.Fatalf("policy mismatch:\nexpected: %v\ngot: %v", expectedTokenPolicies, secret.Auth.TokenPolicies)
			}

			if !strutil.EquivalentSlices(secret.Auth.IdentityPolicies, tc.identityPolicies) {
				t.Fatalf("identity policy mismatch:\nexpected: %v\ngot: %v", tc.identityPolicies, secret.Auth.IdentityPolicies)
			}

			if !strutil.EquivalentSlices(secret.Auth.Policies, expectedPolicies) {
				t.Fatalf("policy mismatch:\nexpected: %v\ngot: %v", expectedPolicies, secret.Auth.Policies)
			}
		})
	}
}
