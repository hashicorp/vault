package vault_test

import (
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/credential/approle"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
)

func TestExpiration_RenewToken_TestCluster(t *testing.T) {
	// Use a TestCluster and the approle backend to test renewal
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
		DefaultLeaseTTL: "5s",
		MaxLeaseTTL:     "5s",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create role
	resp, err := client.Logical().Write("auth/approle/role/role-period", map[string]interface{}{
		"period": "5s",
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
	// Wait 3 seconds
	time.Sleep(3 * time.Second)

	// Renew
	resp, err = client.Logical().Write("auth/token/renew", map[string]interface{}{
		"token": roleToken,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("expected a response for renew")
	}

	// Wait 3 seconds
	time.Sleep(3 * time.Second)

	// Do a second renewal to ensure that period can be renewed past sys/mount max_ttl
	resp, err = client.Logical().Write("auth/token/renew", map[string]interface{}{
		"token": roleToken,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("expected a response for renew")
	}
}
