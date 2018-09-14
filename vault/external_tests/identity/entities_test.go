package identity

import (
	"strings"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/credential/approle"
	"github.com/hashicorp/vault/helper/strutil"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
)

func TestIdentityStore_EntityDisabled(t *testing.T) {
	// Use a TestCluster and the approle backend to get a token and entity for testing
	coreConfig := &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"approle": approle.Factory,
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
	// Use a TestCluster and the approle backend to get a token and entity for testing
	coreConfig := &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"approle": approle.Factory,
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
