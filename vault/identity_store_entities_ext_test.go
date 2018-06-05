package vault_test

import (
	"strings"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/credential/approle"
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
