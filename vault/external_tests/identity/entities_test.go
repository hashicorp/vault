// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package identity

import (
	"strings"
	"testing"

	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/testhelpers/minimal"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

func TestIdentityStore_EntityDisabled(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	// Mount the auth backend
	err := client.Sys().EnableAuthWithOptions("approle", &api.EnableAuthOptions{
		Type: "approle",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Tune the mount
	err = client.Sys().TuneMount("auth/approle", api.MountConfigInput{
		DefaultLeaseTTL: "5m",
		MaxLeaseTTL:     "5m",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create role
	resp, err := client.Logical().Write("auth/approle/role/role-period", map[string]interface{}{
		"period": "5m",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Get role_id
	resp, err = client.Logical().Read("auth/approle/role/role-period/role-id")
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("expected a response for fetching the role-id")
	}
	roleID := resp.Data["role_id"]

	// Get secret_id
	resp, err = client.Logical().Write("auth/approle/role/role-period/secret-id", map[string]interface{}{})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("expected a response for fetching the secret-id")
	}
	secretID := resp.Data["secret_id"]

	// Login
	resp, err = client.Logical().Write("auth/approle/login", map[string]interface{}{
		"role_id":   roleID,
		"secret_id": secretID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("expected a response for login")
	}
	if resp.Auth == nil {
		t.Fatal("expected auth object from response")
	}
	if resp.Auth.ClientToken == "" {
		t.Fatal("expected a client token")
	}

	roleToken := resp.Auth.ClientToken

	client.SetToken(roleToken)
	resp, err = client.Auth().Token().LookupSelf()
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("expected a response for token lookup")
	}
	entityIDRaw, ok := resp.Data["entity_id"]
	if !ok {
		t.Fatal("expected an entity ID")
	}
	entityID, ok := entityIDRaw.(string)
	if !ok {
		t.Fatal("entity_id not a string")
	}

	client.SetToken(cluster.RootToken)
	resp, err = client.Logical().Write("identity/entity/id/"+entityID, map[string]interface{}{
		"disabled": true,
	})
	if err != nil {
		t.Fatal(err)
	}

	// This call should now fail
	client.SetToken(roleToken)
	resp, err = client.Auth().Token().LookupSelf()
	if err == nil {
		t.Fatalf("expected error, got %#v", *resp)
	}
	if !strings.Contains(err.Error(), logical.ErrPermissionDenied.Error()) {
		t.Fatalf("expected to see entity disabled error, got %v", err)
	}

	// Attempting to get a new token should also now fail
	client.SetToken("")
	resp, err = client.Logical().Write("auth/approle/login", map[string]interface{}{
		"role_id":   roleID,
		"secret_id": secretID,
	})
	if err == nil {
		t.Fatalf("expected error, got %#v", *resp)
	}
	if !strings.Contains(err.Error(), logical.ErrPermissionDenied.Error()) {
		t.Fatalf("expected to see entity disabled error, got %v", err)
	}

	client.SetToken(cluster.RootToken)
	resp, err = client.Logical().Write("identity/entity/id/"+entityID, map[string]interface{}{
		"disabled": false,
	})
	if err != nil {
		t.Fatal(err)
	}

	client.SetToken(roleToken)
	resp, err = client.Auth().Token().LookupSelf()
	if err != nil {
		t.Fatal(err)
	}

	// Getting a new token should now work again too
	client.SetToken("")
	resp, err = client.Logical().Write("auth/approle/login", map[string]interface{}{
		"role_id":   roleID,
		"secret_id": secretID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("expected a response for login")
	}
	if resp.Auth == nil {
		t.Fatal("expected auth object from response")
	}
	if resp.Auth.ClientToken == "" {
		t.Fatal("expected a client token")
	}
}

func TestIdentityStore_EntityPoliciesInInitialAuth(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	// Mount the auth backend
	err := client.Sys().EnableAuthWithOptions("approle", &api.EnableAuthOptions{
		Type: "approle",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Tune the mount
	err = client.Sys().TuneMount("auth/approle", api.MountConfigInput{
		DefaultLeaseTTL: "5m",
		MaxLeaseTTL:     "5m",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create role
	resp, err := client.Logical().Write("auth/approle/role/role-period", map[string]interface{}{
		"period": "5m",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Get role_id
	resp, err = client.Logical().Read("auth/approle/role/role-period/role-id")
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("expected a response for fetching the role-id")
	}
	roleID := resp.Data["role_id"]

	// Get secret_id
	resp, err = client.Logical().Write("auth/approle/role/role-period/secret-id", map[string]interface{}{})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("expected a response for fetching the secret-id")
	}
	secretID := resp.Data["secret_id"]

	// Login
	resp, err = client.Logical().Write("auth/approle/login", map[string]interface{}{
		"role_id":   roleID,
		"secret_id": secretID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("expected a response for login")
	}
	if resp.Auth == nil {
		t.Fatal("expected auth object from response")
	}
	if resp.Auth.ClientToken == "" {
		t.Fatal("expected a client token")
	}
	if !strutil.EquivalentSlices(resp.Auth.TokenPolicies, []string{"default"}) {
		t.Fatalf("policy mismatch, got token policies: %v", resp.Auth.TokenPolicies)
	}
	if len(resp.Auth.IdentityPolicies) > 0 {
		t.Fatalf("policy mismatch, got identity policies: %v", resp.Auth.IdentityPolicies)
	}
	if !strutil.EquivalentSlices(resp.Auth.Policies, []string{"default"}) {
		t.Fatalf("policy mismatch, got policies: %v", resp.Auth.Policies)
	}

	// Check policies
	client.SetToken(resp.Auth.ClientToken)
	resp, err = client.Auth().Token().LookupSelf()
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("expected a response for token lookup")
	}
	entityIDRaw, ok := resp.Data["entity_id"]
	if !ok {
		t.Fatal("expected an entity ID")
	}
	entityID, ok := entityIDRaw.(string)
	if !ok {
		t.Fatal("entity_id not a string")
	}
	policiesRaw := resp.Data["policies"]
	if policiesRaw == nil {
		t.Fatal("expected policies, got nil")
	}
	var policies []string
	for _, v := range policiesRaw.([]interface{}) {
		policies = append(policies, v.(string))
	}
	policiesRaw = resp.Data["identity_policies"]
	if policiesRaw != nil {
		t.Fatalf("expected nil policies, got %#v", policiesRaw)
	}
	if !strutil.EquivalentSlices(policies, []string{"default"}) {
		t.Fatalf("policy mismatch, got policies: %v", resp.Auth.Policies)
	}

	// Write more policies into the entity
	client.SetToken(cluster.RootToken)
	resp, err = client.Logical().Write("identity/entity/id/"+entityID, map[string]interface{}{
		"policies": []string{"foo", "bar"},
	})
	if err != nil {
		t.Fatal(err)
	}

	// Reauthenticate to get a token with updated policies
	client.SetToken("")
	resp, err = client.Logical().Write("auth/approle/login", map[string]interface{}{
		"role_id":   roleID,
		"secret_id": secretID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("expected a response for login")
	}
	if resp.Auth == nil {
		t.Fatal("expected auth object from response")
	}
	if resp.Auth.ClientToken == "" {
		t.Fatal("expected a client token")
	}
	if !strutil.EquivalentSlices(resp.Auth.TokenPolicies, []string{"default"}) {
		t.Fatalf("policy mismatch, got token policies: %v", resp.Auth.TokenPolicies)
	}
	if !strutil.EquivalentSlices(resp.Auth.IdentityPolicies, []string{"foo", "bar"}) {
		t.Fatalf("policy mismatch, got identity policies: %v", resp.Auth.IdentityPolicies)
	}
	if !strutil.EquivalentSlices(resp.Auth.Policies, []string{"default", "foo", "bar"}) {
		t.Fatalf("policy mismatch, got policies: %v", resp.Auth.Policies)
	}

	// Validate the policies on lookup again -- this ensures that the right
	// policies were encoded on the token but all were looked up successfully
	client.SetToken(resp.Auth.ClientToken)
	resp, err = client.Auth().Token().LookupSelf()
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("expected a response for token lookup")
	}
	entityIDRaw, ok = resp.Data["entity_id"]
	if !ok {
		t.Fatal("expected an entity ID")
	}
	entityID, ok = entityIDRaw.(string)
	if !ok {
		t.Fatal("entity_id not a string")
	}
	policies = nil
	policiesRaw = resp.Data["policies"]
	if policiesRaw == nil {
		t.Fatal("expected policies, got nil")
	}
	for _, v := range policiesRaw.([]interface{}) {
		policies = append(policies, v.(string))
	}
	if !strutil.EquivalentSlices(policies, []string{"default"}) {
		t.Fatalf("policy mismatch, got policies: %v", policies)
	}
	policies = nil
	policiesRaw = resp.Data["identity_policies"]
	if policiesRaw == nil {
		t.Fatal("expected policies, got nil")
	}
	for _, v := range policiesRaw.([]interface{}) {
		policies = append(policies, v.(string))
	}
	if !strutil.EquivalentSlices(policies, []string{"foo", "bar"}) {
		t.Fatalf("policy mismatch, got policies: %v", policies)
	}
}

// TestIdentity_EntityMerge_RequiresSudo verifies the identity entity merge API endpoint requires the sudo capability (in addition to update) when called via the HTTP API.
func TestIdentity_EntityMerge_RequiresSudo(t *testing.T) {
	t.Parallel()

	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client
	rootToken := client.Token()

	// Create two entities as root to merge.
	toEnt, err := client.Logical().Write("identity/entity", nil)
	require.NoError(t, err)
	require.NotNil(t, toEnt)
	toID, ok := toEnt.Data["id"].(string)
	require.True(t, ok)
	require.NotEmpty(t, toID)

	fromEnt, err := client.Logical().Write("identity/entity", nil)
	require.NoError(t, err)
	require.NotNil(t, fromEnt)
	fromID, ok := fromEnt.Data["id"].(string)
	require.True(t, ok)
	require.NotEmpty(t, fromID)

	// Token with update but without sudo should be denied.
	require.NoError(t, client.Sys().PutPolicy("identity-merge-no-sudo", `
path "identity/entity/merge" {
  capabilities = ["update"]
}
`))

	noSudoTok, err := client.Auth().Token().Create(&api.TokenCreateRequest{
		Policies: []string{"identity-merge-no-sudo"},
	})
	require.NoError(t, err)
	require.NotNil(t, noSudoTok)
	require.NotNil(t, noSudoTok.Auth)
	require.NotEmpty(t, noSudoTok.Auth.ClientToken)

	client.SetToken(noSudoTok.Auth.ClientToken)
	_, err = client.Logical().Write("identity/entity/merge", map[string]interface{}{
		"to_entity_id":    toID,
		"from_entity_ids": []string{fromID},
	})
	require.Error(t, err)
	require.ErrorContains(t, err, logical.ErrPermissionDenied.Error())

	// Token with update+sudo should succeed.
	client.SetToken(rootToken)
	require.NoError(t, client.Sys().PutPolicy("identity-merge-with-sudo", `
path "identity/entity/merge" {
  capabilities = ["update", "sudo"]
}
`))

	sudoTok, err := client.Auth().Token().Create(&api.TokenCreateRequest{
		Policies: []string{"identity-merge-with-sudo"},
	})
	require.NoError(t, err)
	require.NotNil(t, sudoTok)
	require.NotNil(t, sudoTok.Auth)
	require.NotEmpty(t, sudoTok.Auth.ClientToken)

	client.SetToken(sudoTok.Auth.ClientToken)
	_, err = client.Logical().Write("identity/entity/merge", map[string]interface{}{
		"to_entity_id":    toID,
		"from_entity_ids": []string{fromID},
	})
	require.NoError(t, err)

	// Verify the merge happened by checking the from entity was deleted.
	client.SetToken(rootToken)
	deleted, err := client.Logical().Read("identity/entity/id/" + fromID)
	require.NoError(t, err)
	require.Nil(t, deleted)
}
