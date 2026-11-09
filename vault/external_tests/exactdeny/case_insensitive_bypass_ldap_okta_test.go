// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package exactdeny

import (
	"fmt"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/testhelpers/minimal"
	"github.com/stretchr/testify/require"
)

// TestLDAP_GroupExactDenyNotBypassedByCaseVariant verifies that a delegated
// token with wildcard access to auth/<mount>/groups/* cannot bypass an exact
// deny on a specific LDAP group by requesting a case variant of its name.
// With case_sensitive_names=false (the default) the LDAP backend lowercases
// the group name before every storage access, so every variant addresses the
// protected group. The backend is mounted at a non-default path to exercise
// mount-relative matching.
func TestLDAP_GroupExactDenyNotBypassedByCaseVariant(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	rootClient := cluster.Cores[0].Client

	const mount = "custom-ldap"
	require.NoError(t, rootClient.Sys().EnableAuthWithOptions(mount, &api.EnableAuthOptions{Type: "ldap"}))
	_, err := rootClient.Logical().Write(fmt.Sprintf("auth/%s/config", mount), map[string]interface{}{
		"case_sensitive_names": false,
	})
	require.NoError(t, err)

	_, err = rootClient.Logical().Write(fmt.Sprintf("auth/%s/groups/admin", mount), map[string]interface{}{"policies": "protected"})
	require.NoError(t, err)
	_, err = rootClient.Logical().Write(fmt.Sprintf("auth/%s/groups/other", mount), map[string]interface{}{"policies": "unprotected"})
	require.NoError(t, err)

	policy := fmt.Sprintf(`
path "auth/%[1]s/groups/*" {
  capabilities = ["read", "create", "update", "list", "delete"]
}
path "auth/%[1]s/groups/admin" {
  capabilities = ["deny"]
}
`, mount)
	scopedClient := newScopedClient(t, rootClient, "ldap-group-wildcard-deny-admin", policy)

	// Baseline: the exact deny protects the canonical (lowercase) group name.
	_, err = scopedClient.Logical().Write(fmt.Sprintf("auth/%s/groups/admin", mount), map[string]interface{}{"policies": "attacker"})
	require.ErrorContains(t, err, "permission denied")

	// The bypass: read, overwrite, and delete via a case variant must all be
	// denied.
	for _, variant := range []string{"ADMIN", "Admin", "aDmIn"} {
		path := fmt.Sprintf("auth/%s/groups/%s", mount, variant)

		_, err = scopedClient.Logical().Read(path)
		require.ErrorContains(t, err, "permission denied")

		_, err = scopedClient.Logical().Write(path, map[string]interface{}{"policies": "attacker"})
		require.ErrorContains(t, err, "permission denied")

		_, err = scopedClient.Logical().Delete(path)
		require.ErrorContains(t, err, "permission denied")
	}

	// The protected group must be unchanged.
	resp, err := rootClient.Logical().Read(fmt.Sprintf("auth/%s/groups/admin", mount))
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, []interface{}{"protected"}, resp.Data["policies"])

	// Wildcard access to non-protected groups must continue to work.
	_, err = scopedClient.Logical().Write(fmt.Sprintf("auth/%s/groups/other", mount), map[string]interface{}{"policies": "updated"})
	require.NoError(t, err)
}

// TestLDAP_UserExactDenyNotBypassedByCaseVariant is the auth/<mount>/users/
// counterpart of TestLDAP_GroupExactDenyNotBypassedByCaseVariant. Writing a
// user lets the caller assign policies and groups to it, so a bypass here is
// a privilege-escalation primitive.
func TestLDAP_UserExactDenyNotBypassedByCaseVariant(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	rootClient := cluster.Cores[0].Client

	const mount = "custom-ldap"
	require.NoError(t, rootClient.Sys().EnableAuthWithOptions(mount, &api.EnableAuthOptions{Type: "ldap"}))
	_, err := rootClient.Logical().Write(fmt.Sprintf("auth/%s/config", mount), map[string]interface{}{
		"case_sensitive_names": false,
	})
	require.NoError(t, err)

	_, err = rootClient.Logical().Write(fmt.Sprintf("auth/%s/users/admin", mount), map[string]interface{}{"policies": "protected"})
	require.NoError(t, err)
	_, err = rootClient.Logical().Write(fmt.Sprintf("auth/%s/users/other", mount), map[string]interface{}{"policies": "unprotected"})
	require.NoError(t, err)

	policy := fmt.Sprintf(`
path "auth/%[1]s/users/*" {
  capabilities = ["read", "create", "update", "list", "delete"]
}
path "auth/%[1]s/users/admin" {
  capabilities = ["deny"]
}
`, mount)
	scopedClient := newScopedClient(t, rootClient, "ldap-user-wildcard-deny-admin", policy)

	// Baseline: the exact deny protects the canonical (lowercase) username.
	_, err = scopedClient.Logical().Write(fmt.Sprintf("auth/%s/users/admin", mount), map[string]interface{}{"policies": "attacker"})
	require.ErrorContains(t, err, "permission denied")

	// The bypass: read, overwrite, and delete via a case variant must all be
	// denied.
	for _, variant := range []string{"ADMIN", "Admin", "aDmIn"} {
		path := fmt.Sprintf("auth/%s/users/%s", mount, variant)

		_, err = scopedClient.Logical().Read(path)
		require.ErrorContains(t, err, "permission denied")

		_, err = scopedClient.Logical().Write(path, map[string]interface{}{"policies": "attacker"})
		require.ErrorContains(t, err, "permission denied")

		_, err = scopedClient.Logical().Delete(path)
		require.ErrorContains(t, err, "permission denied")
	}

	// The protected user must be unchanged.
	resp, err := rootClient.Logical().Read(fmt.Sprintf("auth/%s/users/admin", mount))
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, []interface{}{"protected"}, resp.Data["policies"])

	// Wildcard access to non-protected users must continue to work.
	_, err = scopedClient.Logical().Write(fmt.Sprintf("auth/%s/users/other", mount), map[string]interface{}{"policies": "updated"})
	require.NoError(t, err)
}

// TestLDAP_GroupCaseSensitiveNamesDenyAppliesToCaseVariants verifies that an
// exact deny applies to case variants of its path even when LDAP is
// configured with case_sensitive_names=true. A deny on "admin" is treated as
// intent to protect the group regardless of casing, so "Admin" is denied too purposefully.
func TestLDAP_GroupCaseSensitiveNamesDenyAppliesToCaseVariants(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	rootClient := cluster.Cores[0].Client

	require.NoError(t, rootClient.Sys().EnableAuthWithOptions("ldap", &api.EnableAuthOptions{Type: "ldap"}))
	_, err := rootClient.Logical().Write("auth/ldap/config", map[string]interface{}{"case_sensitive_names": true})
	require.NoError(t, err)

	_, err = rootClient.Logical().Write("auth/ldap/groups/admin", map[string]interface{}{"policies": "protected"})
	require.NoError(t, err)
	_, err = rootClient.Logical().Write("auth/ldap/groups/Admin", map[string]interface{}{"policies": "distinct"})
	require.NoError(t, err)

	policy := `
path "auth/ldap/groups/*" {
  capabilities = ["read", "create", "update", "list", "delete"]
}
path "auth/ldap/groups/admin" {
  capabilities = ["deny"]
}
`
	scopedClient := newScopedClient(t, rootClient, "ldap-group-cs-deny-admin", policy)

	_, err = scopedClient.Logical().Write("auth/ldap/groups/admin", map[string]interface{}{"policies": "attacker"})
	require.ErrorContains(t, err, "permission denied")

	// "Admin" is a distinct group in case-sensitive mode, but the deny on
	// "admin" still applies to it.
	_, err = scopedClient.Logical().Write("auth/ldap/groups/Admin", map[string]interface{}{"policies": "updated"})
	require.ErrorContains(t, err, "permission denied")

	resp, err := rootClient.Logical().Read("auth/ldap/groups/admin")
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, []interface{}{"protected"}, resp.Data["policies"])

	resp, err = rootClient.Logical().Read("auth/ldap/groups/Admin")
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, []interface{}{"distinct"}, resp.Data["policies"])
}

// TestLDAP_UserCaseSensitiveNamesDenyAppliesToCaseVariants is the users/
// counterpart of TestLDAP_GroupCaseSensitiveNamesDenyAppliesToCaseVariants.
func TestLDAP_UserCaseSensitiveNamesDenyAppliesToCaseVariants(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	rootClient := cluster.Cores[0].Client

	require.NoError(t, rootClient.Sys().EnableAuthWithOptions("ldap", &api.EnableAuthOptions{Type: "ldap"}))
	_, err := rootClient.Logical().Write("auth/ldap/config", map[string]interface{}{"case_sensitive_names": true})
	require.NoError(t, err)

	_, err = rootClient.Logical().Write("auth/ldap/users/admin", map[string]interface{}{"policies": "protected"})
	require.NoError(t, err)
	_, err = rootClient.Logical().Write("auth/ldap/users/Admin", map[string]interface{}{"policies": "distinct"})
	require.NoError(t, err)

	policy := `
path "auth/ldap/users/*" {
  capabilities = ["read", "create", "update", "list", "delete"]
}
path "auth/ldap/users/admin" {
  capabilities = ["deny"]
}
`
	scopedClient := newScopedClient(t, rootClient, "ldap-user-cs-deny-admin", policy)

	_, err = scopedClient.Logical().Write("auth/ldap/users/admin", map[string]interface{}{"policies": "attacker"})
	require.ErrorContains(t, err, "permission denied")

	// "Admin" is a distinct user in case-sensitive mode, but the deny on
	// "admin" still applies to it.
	_, err = scopedClient.Logical().Write("auth/ldap/users/Admin", map[string]interface{}{"policies": "updated"})
	require.ErrorContains(t, err, "permission denied")

	resp, err := rootClient.Logical().Read("auth/ldap/users/admin")
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, []interface{}{"protected"}, resp.Data["policies"])

	resp, err = rootClient.Logical().Read("auth/ldap/users/Admin")
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, []interface{}{"distinct"}, resp.Data["policies"])
}

// TestOkta_GroupExactDenyNotBypassedByCaseVariant verifies that a delegated
// token with wildcard access to auth/okta/groups/* cannot bypass an exact
// deny on a specific Okta group by requesting a case variant of its name.
// The Okta backend resolves group names with strings.EqualFold, so a variant
// addresses the same stored group. New groups are always stored lowercased
// (pathGroupWrite), which is the case covered here.
func TestOkta_GroupExactDenyNotBypassedByCaseVariant(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	rootClient := cluster.Cores[0].Client

	require.NoError(t, rootClient.Sys().EnableAuthWithOptions("okta", &api.EnableAuthOptions{Type: "okta"}))

	_, err := rootClient.Logical().Write("auth/okta/groups/admin", map[string]interface{}{"policies": "protected"})
	require.NoError(t, err)
	_, err = rootClient.Logical().Write("auth/okta/groups/other", map[string]interface{}{"policies": "unprotected"})
	require.NoError(t, err)

	policy := `
path "auth/okta/groups/*" {
  capabilities = ["read", "create", "update", "list", "delete"]
}
path "auth/okta/groups/admin" {
  capabilities = ["deny"]
}
`
	scopedClient := newScopedClient(t, rootClient, "okta-group-wildcard-deny-admin", policy)

	// Baseline: the exact deny protects the canonical (lowercase) group name.
	_, err = scopedClient.Logical().Write("auth/okta/groups/admin", map[string]interface{}{"policies": "attacker"})
	require.ErrorContains(t, err, "permission denied")

	// The bypass: read, overwrite, and delete via a case variant must all be
	// denied.
	for _, variant := range []string{"ADMIN", "Admin", "aDmIn"} {
		path := fmt.Sprintf("auth/okta/groups/%s", variant)

		_, err = scopedClient.Logical().Read(path)
		require.ErrorContains(t, err, "permission denied")

		_, err = scopedClient.Logical().Write(path, map[string]interface{}{"policies": "attacker"})
		require.ErrorContains(t, err, "permission denied")

		_, err = scopedClient.Logical().Delete(path)
		require.ErrorContains(t, err, "permission denied")
	}

	// The protected group must be unchanged and accessible by rootClient.
	resp, err := rootClient.Logical().Read("auth/okta/groups/admin")
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, []interface{}{"protected"}, resp.Data["policies"])

	// Wildcard access to non-protected groups must continue to work.
	_, err = scopedClient.Logical().Write("auth/okta/groups/other", map[string]interface{}{"policies": "updated"})
	require.NoError(t, err)
}
