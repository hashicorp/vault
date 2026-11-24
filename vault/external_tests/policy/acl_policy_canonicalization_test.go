// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package policy

import (
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/testhelpers/minimal"
	"github.com/stretchr/testify/require"
)

// deniedSuperAdminPolicy is an ACL policy for a delegated administrator that
// may manage principals but is explicitly forbidden from assigning the
// high-privilege "super-admin" policy, via a value-specific denied_parameters
// restriction on every policy-bearing parameter.
const deniedSuperAdminPolicy = `
path "auth/token/create" {
	capabilities = ["create", "update"]
	denied_parameters = {
		"policies" = ["super-admin"]
	}
}

path "auth/token/roles/*" {
	capabilities = ["create", "update"]
	denied_parameters = {
		"allowed_policies" = ["super-admin"]
		"disallowed_policies" = ["super-admin"]
		"allowed_policies_glob" = ["super-admin"]
		"disallowed_policies_glob" = ["super-admin"]
	}
}

path "identity/entity" {
	capabilities = ["create", "update"]
	denied_parameters = {
		"policies" = ["super-admin"]
	}
}

path "identity/group" {
	capabilities = ["create", "update"]
	denied_parameters = {
		"policies" = ["super-admin"]
	}
}

path "auth/approle/role/*" {
	capabilities = ["create", "update"]
	denied_parameters = {
		"token_policies" = ["super-admin"]
	}
}

path "auth/userpass/users/*" {
	capabilities = ["create", "update"]
	denied_parameters = {
		"token_policies" = ["super-admin"]
	}
}
`

// superAdminPolicy stands in for an existing, higher-privilege non-root policy
// that already exists in the namespace and that the delegated administrator
// must not be able to assign.
const superAdminPolicy = `
path "secret/*" {
	capabilities = ["create", "read", "update", "delete", "list"]
}
`

// policyNameVariants returns representations of the denied policy name
// "super-admin" that are not byte-identical to the denied value but that all
// canonicalize to it. Vault lowercases policy names and trims surrounding
// whitespace before storing or enforcing them, so each of these resolves to the
// denied policy and must be rejected by denied_parameters.
func policyNameVariants() map[string]string {
	return map[string]string{
		"mixed_case":           "Super-Admin",
		"upper_case":           "SUPER-ADMIN",
		"leading_space":        " super-admin",
		"trailing_space":       "super-admin ",
		"surrounding_spaces":   "  super-admin  ",
		"tabs":                 "\tsuper-admin\t",
		"newlines":             "\nsuper-admin\n",
		"unicode_nbsp":         "\u00a0super-admin\u00a0",
		"whitespace_plus_case": " Super-Admin\t",
		"exact_denied_value":   "super-admin",
	}
}

// setupDelegatedAdmin creates a single-node cluster with a "super-admin" policy
// and a delegated administrator token restricted by deniedSuperAdminPolicy. It
// returns the root client and the delegated administrator token.
func setupDelegatedAdmin(t *testing.T) (client *api.Client, delegatedAdminToken string) {
	t.Helper()

	cluster := minimal.NewTestSoloCluster(t, nil)
	client = cluster.Cores[0].Client

	require.NoError(t, client.Sys().PutPolicy("super-admin", superAdminPolicy))
	require.NoError(t, client.Sys().PutPolicy("delegated-admin", deniedSuperAdminPolicy))

	secret, err := client.Auth().Token().Create(&api.TokenCreateRequest{
		Policies: []string{"delegated-admin"},
	})
	require.NoError(t, err)
	require.NotNil(t, secret)
	require.NotNil(t, secret.Auth)

	return client, secret.Auth.ClientToken
}

// TestACLPolicyCanonicalization_TokenCreate verifies that a delegated
// administrator cannot bypass a value-specific denied_parameters restriction on
// the "policies" parameter of token creation by submitting a non-canonical
// representation of the denied policy name, such as a mixed-case or
// whitespace-padded variant. Vault canonicalizes policy names before storing
// them, so without canonicalization at the ACL layer such a request would be
// permitted and would mint a token holding the denied policy.
func TestACLPolicyCanonicalization_TokenCreate(t *testing.T) {
	t.Parallel()

	client, delegatedAdminToken := setupDelegatedAdmin(t)
	client.SetToken(delegatedAdminToken)

	for name, variant := range policyNameVariants() {
		t.Run(name, func(t *testing.T) {
			resp, err := client.Auth().Token().Create(&api.TokenCreateRequest{
				Policies: []string{variant},
			})
			require.Error(t, err, "token creation with policy %q must be denied", variant)
			require.Contains(t, err.Error(), "permission denied")

			if resp != nil && resp.Auth != nil {
				require.NotContains(t, resp.Auth.Policies, "super-admin",
					"issued token must not carry the denied policy")
			}
		})
	}
}

// TestACLPolicyCanonicalization_TokenRole verifies that the token-role policy
// fields are protected by denied_parameters against non-canonical policy names.
// These fields were not covered by the earlier partial remediation. They are
// the strongest non-sudo escalation path because the policies configured on a
// role determine the policies of tokens issued through it, and those may be
// disjoint from the creating caller's own policies.
func TestACLPolicyCanonicalization_TokenRole(t *testing.T) {
	t.Parallel()

	client, delegatedAdminToken := setupDelegatedAdmin(t)
	client.SetToken(delegatedAdminToken)

	fields := []string{
		"allowed_policies",
		"disallowed_policies",
		"allowed_policies_glob",
		"disallowed_policies_glob",
	}

	for _, field := range fields {
		for name, variant := range policyNameVariants() {
			t.Run(field+"/"+name, func(t *testing.T) {
				// Each subtest writes to a distinct role name so that subtests
				// share no state and can run independently.
				rolePath := "auth/token/roles/role-" + field + "-" + name

				_, err := client.Logical().Write(rolePath, map[string]interface{}{
					field: []string{variant},
				})
				require.Error(t, err, "writing %s=%q must be denied", field, variant)
				require.Contains(t, err.Error(), "permission denied")
			})
		}
	}
}

// TestACLPolicyCanonicalization_IdentityEntity verifies that identity entity
// policy assignment honors denied_parameters for non-canonical policy names.
// Identity policies are applied to any token whose entity matches, so this path
// can escalate privileges without the parent-policy ceiling that constrains
// direct token creation.
func TestACLPolicyCanonicalization_IdentityEntity(t *testing.T) {
	t.Parallel()

	client, delegatedAdminToken := setupDelegatedAdmin(t)
	client.SetToken(delegatedAdminToken)

	for name, variant := range policyNameVariants() {
		t.Run(name, func(t *testing.T) {
			_, err := client.Logical().Write("identity/entity", map[string]interface{}{
				"name":     "entity-" + name,
				"policies": []string{variant},
			})
			require.Error(t, err, "entity policy %q must be denied", variant)
			require.Contains(t, err.Error(), "permission denied")
		})
	}
}

// TestACLPolicyCanonicalization_IdentityGroup verifies that identity group
// policy assignment honors denied_parameters for non-canonical policy names.
// Group policies apply to every member entity, making this a broad escalation
// path if the ACL check can be evaded.
func TestACLPolicyCanonicalization_IdentityGroup(t *testing.T) {
	t.Parallel()

	client, delegatedAdminToken := setupDelegatedAdmin(t)
	client.SetToken(delegatedAdminToken)

	for name, variant := range policyNameVariants() {
		t.Run(name, func(t *testing.T) {
			_, err := client.Logical().Write("identity/group", map[string]interface{}{
				"name":     "group-" + name,
				"policies": []string{variant},
			})
			require.Error(t, err, "group policy %q must be denied", variant)
			require.Contains(t, err.Error(), "permission denied")
		})
	}
}

// TestACLPolicyCanonicalization_AuthRoles verifies that auth-method role policy
// assignment honors denied_parameters for non-canonical policy names, using
// AppRole and userpass as representative consumers. Both canonicalize the
// submitted policy values after ACL evaluation, so a mismatch between the two
// layers would let a delegated administrator mint principals that authenticate
// with the denied policy.
func TestACLPolicyCanonicalization_AuthRoles(t *testing.T) {
	t.Parallel()

	client, delegatedAdminToken := setupDelegatedAdmin(t)

	require.NoError(t, client.Sys().EnableAuthWithOptions("approle", &api.EnableAuthOptions{Type: "approle"}))
	require.NoError(t, client.Sys().EnableAuthWithOptions("userpass", &api.EnableAuthOptions{Type: "userpass"}))
	client.SetToken(delegatedAdminToken)

	for name, variant := range policyNameVariants() {
		t.Run("approle/"+name, func(t *testing.T) {
			_, err := client.Logical().Write("auth/approle/role/approle-"+name, map[string]interface{}{
				"token_policies": []string{variant},
			})
			require.Error(t, err, "approle token_policies %q must be denied", variant)
			require.Contains(t, err.Error(), "permission denied")
		})

		t.Run("userpass/"+name, func(t *testing.T) {
			_, err := client.Logical().Write("auth/userpass/users/user-"+name, map[string]interface{}{
				"password":       "test-password-1234",
				"token_policies": []string{variant},
			})
			require.Error(t, err, "userpass token_policies %q must be denied", variant)
			require.Contains(t, err.Error(), "permission denied")
		})
	}
}

// TestACLPolicyCanonicalization_NonDeniedPolicyStillAllowed verifies that
// canonicalization only blocks values that actually resolve to the denied
// policy. A policy name that is merely similar, or that differs by interior
// whitespace, is a genuinely different policy and must remain assignable, so
// that the fix does not over-block legitimate delegated administration.
func TestACLPolicyCanonicalization_NonDeniedPolicyStillAllowed(t *testing.T) {
	t.Parallel()

	cluster := minimal.NewTestSoloCluster(t, nil)
	root := cluster.Cores[0].Client

	require.NoError(t, root.Sys().PutPolicy("super-admin", superAdminPolicy))
	require.NoError(t, root.Sys().PutPolicy("delegated-admin", deniedSuperAdminPolicy))
	require.NoError(t, root.Sys().PutPolicy("read-only", `
path "secret/data/*" {
	capabilities = ["read"]
}
`))
	require.NoError(t, root.Sys().PutPolicy("super-admin-readonly", `
path "secret/data/readonly/*" {
	capabilities = ["read"]
}
`))

	adminToken, err := root.Auth().Token().Create(&api.TokenCreateRequest{
		Policies: []string{"delegated-admin", "read-only", "super-admin-readonly"},
	})
	require.NoError(t, err)
	require.NotNil(t, adminToken)
	require.NotNil(t, adminToken.Auth)

	delegated, err := root.Clone()
	require.NoError(t, err)
	delegated.SetToken(adminToken.Auth.ClientToken)

	allowed := map[string]string{
		"plain_other_policy":   "read-only",
		"mixed_case_other":     "Read-Only",
		"padded_other":         "  read-only  ",
		"similar_but_distinct": "super-admin-readonly",
	}

	for name, value := range allowed {
		t.Run(name, func(t *testing.T) {
			secret, err := delegated.Auth().Token().Create(&api.TokenCreateRequest{
				Policies: []string{value},
			})
			require.NoError(t, err, "policy %q should remain assignable", value)
			require.NotNil(t, secret)
			require.NotNil(t, secret.Auth)
			require.NotContains(t, secret.Auth.Policies, "super-admin")
		})
	}
}

// TestACLPolicyCanonicalization_RootAndParentSafeguardsIntact verifies that the
// pre-existing safeguards still hold alongside the canonicalization fix: the
// root policy cannot be granted by a non-root caller, and token creation
// without sudo remains restricted to a subset of the caller's own policies.
// These bound the severity of the vulnerability and must not regress.
func TestACLPolicyCanonicalization_RootAndParentSafeguardsIntact(t *testing.T) {
	t.Parallel()

	client, delegatedAdminToken := setupDelegatedAdmin(t)
	client.SetToken(delegatedAdminToken)

	t.Run("root_policy_rejected", func(t *testing.T) {
		for _, variant := range []string{"root", "Root", " root ", "ROOT"} {
			_, err := client.Auth().Token().Create(&api.TokenCreateRequest{
				Policies: []string{variant},
			})
			require.Error(t, err, "root policy variant %q must be rejected", variant)
		}
	})

	t.Run("parent_policy_subset_enforced_without_sudo", func(t *testing.T) {
		// "read-only" is not held by the delegated administrator, so even
		// though it is not explicitly denied, the non-sudo parent-policy
		// subset restriction must prevent it from being granted.
		_, err := client.Auth().Token().Create(&api.TokenCreateRequest{
			Policies: []string{"some-policy-the-parent-does-not-have"},
		})
		require.Error(t, err, "non-sudo token creation must be limited to the parent's policies")
	})
}

// TestACLPolicyCanonicalization_ConfiguredDenials verifies that a delegated
// administrator cannot assign a denied policy when the configured ACL pattern
// uses case or whitespace that differs from Vault's canonical policy name.
func TestACLPolicyCanonicalization_ConfiguredDenials(t *testing.T) {
	t.Parallel()

	cluster := minimal.NewTestSoloCluster(t, nil)
	root := cluster.Cores[0].Client
	require.NoError(t, root.Sys().PutPolicy("delegated-admin", `
path "auth/token/create" {
	capabilities = ["create", "update"]
	denied_parameters = { "policies" = [" Super-Admin "] }
}
path "auth/token/roles/*" {
	capabilities = ["create", "update"]
	denied_parameters = {
		"allowed_policies" = ["Super-Admin"]
		"disallowed_policies" = [" SUPER-ADMIN "]
		"allowed_policies_glob" = ["Super-*"]
		"disallowed_policies_glob" = [" *-Admin "]
	}
}
path "identity/entity" {
	capabilities = ["create", "update"]
	denied_parameters = { "policies" = ["Super-Admin", "TEAM.*.LEAD"] }
}
path "auth/userpass/users/*" {
	capabilities = ["create", "update"]
	denied_parameters = { "token_policies" = ["SUPER-ADMIN"] }
}
`))
	require.NoError(t, root.Sys().EnableAuthWithOptions("userpass", &api.EnableAuthOptions{Type: "userpass"}))
	secret, err := root.Auth().Token().Create(&api.TokenCreateRequest{Policies: []string{"delegated-admin"}})
	require.NoError(t, err)
	require.NotNil(t, secret.Auth)
	delegated, err := root.Clone()
	require.NoError(t, err)
	delegated.SetToken(secret.Auth.ClientToken)

	_, err = delegated.Auth().Token().Create(&api.TokenCreateRequest{Policies: []string{"Super-Admin"}})
	require.ErrorContains(t, err, "permission denied", "padded HCL denial must block token creation")

	requests := []struct {
		path string
		data map[string]interface{}
	}{
		{"auth/token/roles/allowed", map[string]interface{}{"allowed_policies": "super-admin"}},
		{"auth/token/roles/disallowed", map[string]interface{}{"disallowed_policies": []string{"super-admin"}}},
		{"auth/token/roles/allowed-glob", map[string]interface{}{"allowed_policies_glob": []string{"super-admin"}}},
		{"auth/token/roles/disallowed-glob", map[string]interface{}{"disallowed_policies_glob": []string{"random-admin"}}},
		{"identity/entity", map[string]interface{}{"name": "configured-denial-entity", "policies": []string{"SUPER-ADMIN"}}},
		{"auth/userpass/users/configured-denial-user", map[string]interface{}{"password": "test-password-1234", "token_policies": []string{"super-admin"}}},
	}
	for _, request := range requests {
		_, err := delegated.Logical().Write(request.path, request.data)
		require.ErrorContains(t, err, "permission denied", "configured denial must block %s", request.path)
	}

	// Interior wildcards in policy-name denials match canonical policy names.
	_, err = delegated.Logical().Write("identity/entity", map[string]interface{}{
		"name": "interior-glob-entity", "policies": []string{"team.platform.lead"},
	})
	require.ErrorContains(t, err, "permission denied", "interior wildcard must deny matching policy names")
}

// TestACLPolicyCanonicalization_TokenCreateInfixGlob verifies that a delegated
// sudo caller cannot create a token with a denied policy when its ACL denial
// uses an interior wildcard, even when the configured pattern has mixed case.
func TestACLPolicyCanonicalization_TokenCreateInfixGlob(t *testing.T) {
	t.Parallel()

	cluster := minimal.NewTestSoloCluster(t, nil)
	root := cluster.Cores[0].Client
	require.NoError(t, root.Sys().PutPolicy("team.platform.lead", superAdminPolicy))
	require.NoError(t, root.Sys().PutPolicy("delegated-admin", `
path "auth/token/create" {
 capabilities = ["create", "update", "sudo"]
 denied_parameters = {
  "policies" = ["TEAM.*.LEAD"]
  "token_policies" = ["TEAM.*.LEAD"]
 }
}
`))
	secret, err := root.Auth().Token().Create(&api.TokenCreateRequest{Policies: []string{"delegated-admin"}})
	require.NoError(t, err)
	require.NotNil(t, secret.Auth)
	delegated, err := root.Clone()
	require.NoError(t, err)
	delegated.SetToken(secret.Auth.ClientToken)

	for _, data := range []map[string]interface{}{
		{"policies": []string{"team.platform.lead"}},
		{"token_policies": []string{"team.platform.lead"}},
	} {
		_, err := delegated.Logical().Write("auth/token/create", data)
		require.ErrorContains(t, err, "permission denied", "denied policy in %v must not mint a token", data)
	}
}

// TestACLPolicyCanonicalization_ConfiguredAllowlists verifies that canonical
// HCL policy patterns admit the intended policy while remaining restrictive,
// without changing the case-sensitive semantics of unrelated parameters.
func TestACLPolicyCanonicalization_ConfiguredAllowlists(t *testing.T) {
	t.Parallel()

	cluster := minimal.NewTestSoloCluster(t, nil)
	root := cluster.Cores[0].Client
	require.NoError(t, root.Sys().PutPolicy("read-only", `
path "secret/data/*" { capabilities = ["read"] }
`))
	require.NoError(t, root.Sys().PutPolicy("delegated-admin", `
path "identity/group" {
	capabilities = ["create", "update"]
	allowed_parameters = {
		"name" = []
		"policies" = [" READ-* "]
	}
	denied_parameters = { "name" = ["Super-Admin"] }
}
`))
	secret, err := root.Auth().Token().Create(&api.TokenCreateRequest{Policies: []string{"delegated-admin"}})
	require.NoError(t, err)
	require.NotNil(t, secret.Auth)
	delegated, err := root.Clone()
	require.NoError(t, err)
	delegated.SetToken(secret.Auth.ClientToken)

	_, err = delegated.Logical().Write("identity/group", map[string]interface{}{
		"name": "super-admin", "policies": []string{"Read-Only"},
	})
	require.NoError(t, err, "a canonical policy matching the configured glob should be allowed")

	_, err = delegated.Logical().Write("identity/group", map[string]interface{}{
		"name": "restricted-group", "policies": []string{"Super-Admin"},
	})
	require.ErrorContains(t, err, "permission denied", "an unrelated policy must not bypass the allowlist")
}
