// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package identity

import (
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/testhelpers/minimal"
	"github.com/stretchr/testify/require"
)

// TestAppRole_MalformedLoginNoPanic verifies that sending a malformed AppRole
// login request — where secret_id is a JSON object instead of a string —
// does not cause a nil pointer dereference panic in aliasNameFromLoginRequest
// (vault/core.go). Previously, the framework would return a logical.ErrorResponse
// with resp.Auth == nil, and the unguarded access to resp.Auth.Alias caused a panic.
//
// Regression test for VAULT-23235.
func TestAppRole_MalformedLoginNoPanic(t *testing.T) {
	t.Parallel()

	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	// Enable AppRole auth method.
	err := client.Sys().EnableAuthWithOptions("approle", &api.EnableAuthOptions{
		Type: "approle",
	})
	require.NoError(t, err)

	// Create a basic role.
	_, err = client.Logical().Write("auth/approle/role/test-role", map[string]interface{}{
		"token_ttl":     "1h",
		"token_max_ttl": "4h",
	})
	require.NoError(t, err)

	// Read the role-id so we can send a structurally valid role_id.
	secret, err := client.Logical().Read("auth/approle/role/test-role/role-id")
	require.NoError(t, err)
	require.NotNil(t, secret)
	roleID, ok := secret.Data["role_id"].(string)
	require.True(t, ok, "role_id should be a string")

	// Send a malformed login where secret_id is a map (object) instead of a
	// plain string. The framework field validator should reject this and return
	// a logical.ErrorResponse; previously aliasNameFromLoginRequest would then
	// panic on resp.Auth.Alias because resp.Auth was nil.
	//
	// The expected outcome is a non-nil error from Vault (a 400-level response)
	// with no panic — the test process surviving this call is the key assertion.
	_, err = client.Logical().Write("auth/approle/login", map[string]interface{}{
		"role_id": roleID,
		"secret_id": map[string]interface{}{
			"Length": "58bd8e99-84ed-1920-ad03-6325b2b69380",
		},
	})
	require.Error(t, err, "expected an error for malformed secret_id, not a successful login")

	// Vault must still be operational after receiving the malformed request.
	health, err := client.Sys().Health()
	require.NoError(t, err)
	require.True(t, health.Initialized, "vault should still be initialized after malformed login")
	require.False(t, health.Sealed, "vault should not be sealed after malformed login")
}
