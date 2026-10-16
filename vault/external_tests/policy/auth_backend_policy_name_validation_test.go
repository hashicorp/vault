// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package policy

import (
	"fmt"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/credential/approle"
	awsauth "github.com/hashicorp/vault/builtin/credential/aws"
	"github.com/hashicorp/vault/builtin/credential/ldap"
	"github.com/hashicorp/vault/builtin/credential/radius"
	credUserpass "github.com/hashicorp/vault/builtin/credential/userpass"
	"github.com/hashicorp/vault/helper/testhelpers/minimal"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/require"
)

// TestAuthBackendPolicyNameValidation_AppRole verifies AppRole policy updates
// reject malformed policy names for both legacy and token policy fields.
func TestAuthBackendPolicyNameValidation_AppRole(t *testing.T) {
	t.Parallel()

	cluster := minimal.NewTestSoloCluster(t, &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"approle": approle.Factory,
		},
	})
	client := cluster.Cores[0].Client

	const mountPath = "approle-policy-validation"
	err := client.Sys().EnableAuthWithOptions(mountPath, &api.EnableAuthOptions{Type: "approle"})
	require.NoError(t, err)

	rolePath := fmt.Sprintf("auth/%s/role/role-invalid-policy", mountPath)
	_, err = client.Logical().Write(rolePath, map[string]interface{}{
		"policies": "default",
	})
	require.NoError(t, err)

	policyPath := rolePath + "/policies"
	for _, data := range []map[string]interface{}{
		{"policies": "../bad"},
		{"token_policies": "../bad"},
	} {
		_, err = client.Logical().Write(policyPath, data)
		require.Error(t, err)
		require.ErrorContains(t, err, "invalid policy name")
		require.ErrorContains(t, err, "../bad")
	}
}

// TestAuthBackendPolicyNameValidation_AWS verifies AWS role-tag policy
// overrides reject malformed policy names.
func TestAuthBackendPolicyNameValidation_AWS(t *testing.T) {
	t.Parallel()

	cluster := minimal.NewTestSoloCluster(t, &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"aws": awsauth.Factory,
		},
	})
	client := cluster.Cores[0].Client

	const mountPath = "aws-policy-validation"
	err := client.Sys().EnableAuthWithOptions(mountPath, &api.EnableAuthOptions{Type: "aws"})
	require.NoError(t, err)

	rolePath := fmt.Sprintf("auth/%s/role/invalid-policy-role", mountPath)
	_, err = client.Logical().Write(rolePath, map[string]interface{}{
		"auth_type":    "ec2",
		"policies":     "p,q,r,s",
		"role_tag":     "VaultRole",
		"bound_ami_id": "ami-abcd123",
	})
	require.NoError(t, err)

	_, err = client.Logical().Write(rolePath+"/tag", map[string]interface{}{
		"policies": "../bad",
	})
	require.Error(t, err)
	require.ErrorContains(t, err, "invalid policy name")
	require.ErrorContains(t, err, "../bad")
}

// TestAuthBackendPolicyNameValidation_LDAP verifies LDAP user and group writes
// reject malformed policy names before entries are persisted.
func TestAuthBackendPolicyNameValidation_LDAP(t *testing.T) {
	t.Parallel()

	cluster := minimal.NewTestSoloCluster(t, &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"ldap": ldap.Factory,
		},
	})
	client := cluster.Cores[0].Client

	const mountPath = "ldap-policy-validation"
	err := client.Sys().EnableAuthWithOptions(mountPath, &api.EnableAuthOptions{Type: "ldap"})
	require.NoError(t, err)

	_, err = client.Logical().Write(fmt.Sprintf("auth/%s/groups/engineers", mountPath), map[string]interface{}{
		"policies": "../bad",
	})
	require.Error(t, err)
	require.ErrorContains(t, err, "invalid policy name")
	require.ErrorContains(t, err, "../bad")

	_, err = client.Logical().Write(fmt.Sprintf("auth/%s/users/hermes-conrad", mountPath), map[string]interface{}{
		"groups":   []string{"engineers"},
		"policies": "../bad",
	})
	require.Error(t, err)
	require.ErrorContains(t, err, "invalid policy name")
	require.ErrorContains(t, err, "../bad")
}

// TestAuthBackendPolicyNameValidation_Radius verifies Radius user writes reject
// malformed policy names before entries are persisted.
func TestAuthBackendPolicyNameValidation_Radius(t *testing.T) {
	t.Parallel()

	cluster := minimal.NewTestSoloCluster(t, &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"radius": radius.Factory,
		},
	})
	client := cluster.Cores[0].Client

	const mountPath = "radius-policy-validation"
	err := client.Sys().EnableAuthWithOptions(mountPath, &api.EnableAuthOptions{Type: "radius"})
	require.NoError(t, err)

	_, err = client.Logical().Write(fmt.Sprintf("auth/%s/users/admin", mountPath), map[string]interface{}{
		"policies": "../bad",
	})
	require.Error(t, err)
	require.ErrorContains(t, err, "invalid policy name")
	require.ErrorContains(t, err, "../bad")
}

// TestAuthBackendPolicyNameValidation_Userpass verifies userpass policy updates
// reject malformed policy names for both legacy and token policy fields.
func TestAuthBackendPolicyNameValidation_Userpass(t *testing.T) {
	t.Parallel()

	cluster := minimal.NewTestSoloCluster(t, &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"userpass": credUserpass.Factory,
		},
	})
	client := cluster.Cores[0].Client

	const mountPath = "userpass-policy-validation"
	err := client.Sys().EnableAuthWithOptions(mountPath, &api.EnableAuthOptions{Type: "userpass"})
	require.NoError(t, err)

	userPath := fmt.Sprintf("auth/%s/users/testuser", mountPath)
	_, err = client.Logical().Write(userPath, map[string]interface{}{
		"password": "testpassword",
		"policies": "default",
	})
	require.NoError(t, err)

	policyPath := userPath + "/policies"
	for _, data := range []map[string]interface{}{
		{"policies": "../bad"},
		{"token_policies": "../bad"},
	} {
		_, err = client.Logical().Write(policyPath, data)
		require.Error(t, err)
		require.ErrorContains(t, err, "invalid policy name")
		require.ErrorContains(t, err, "../bad")
	}
}
