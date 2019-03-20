package token

import (
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/credential/approle"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
)

func TestBatchTokens(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"kv": vault.LeasedPassthroughBackendFactory,
		},
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
	rootToken := client.Token()
	var err error

	// Set up a KV path
	err = client.Sys().Mount("kv", &api.MountInput{
		Type: "kv",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Write("kv/foo", map[string]interface{}{
		"foo": "bar",
		"ttl": "5m",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Write the test policy
	err = client.Sys().PutPolicy("test", `
path "kv/*" {
	capabilities = ["read"]
}`)
	if err != nil {
		t.Fatal(err)
	}

	// Mount the auth backend
	err = client.Sys().EnableAuthWithOptions("approle", &api.EnableAuthOptions{
		Type: "approle",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Tune the mount
	if err = client.Sys().TuneMount("auth/approle", api.MountConfigInput{
		DefaultLeaseTTL: "5s",
		MaxLeaseTTL:     "5s",
	}); err != nil {
		t.Fatal(err)
	}

	// Create role
	resp, err := client.Logical().Write("auth/approle/role/test", map[string]interface{}{
		"policies": "test",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Get role_id
	resp, err = client.Logical().Read("auth/approle/role/test/role-id")
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("expected a response for fetching the role-id")
	}
	roleID := resp.Data["role_id"]

	// Get secret_id
	resp, err = client.Logical().Write("auth/approle/role/test/secret-id", map[string]interface{}{})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("expected a response for fetching the secret-id")
	}
	secretID := resp.Data["secret_id"]

	// Login
	testLogin := func(mountTuneType, roleType string, batch bool) string {
		t.Helper()
		if err = client.Sys().TuneMount("auth/approle", api.MountConfigInput{
			TokenType: mountTuneType,
		}); err != nil {
			t.Fatal(err)
		}
		_, err = client.Logical().Write("auth/approle/role/test", map[string]interface{}{
			"token_type": roleType,
		})
		if err != nil {
			t.Fatal(err)
		}

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
		if batch && !strings.HasPrefix(resp.Auth.ClientToken, "b.") {
			t.Fatal("expected a batch token")
		}
		if !batch && strings.HasPrefix(resp.Auth.ClientToken, "b.") {
			t.Fatal("expected a non-batch token")
		}
		return resp.Auth.ClientToken
	}
	testLogin("service", "default", false)
	testLogin("service", "batch", false)
	testLogin("service", "service", false)
	testLogin("batch", "default", true)
	testLogin("batch", "batch", true)
	testLogin("batch", "service", true)
	testLogin("default-service", "default", false)
	testLogin("default-service", "batch", true)
	testLogin("default-service", "service", false)
	testLogin("default-batch", "default", true)
	testLogin("default-batch", "batch", true)
	testLogin("default-batch", "service", false)

	finalToken := testLogin("batch", "batch", true)

	client.SetToken(finalToken)
	resp, err = client.Logical().Read("kv/foo")
	if err != nil {
		t.Fatal(err)
	}
	if resp.Data["foo"].(string) != "bar" {
		t.Fatal("bad")
	}
	if resp.LeaseID == "" {
		t.Fatal("expected lease")
	}
	if !resp.Renewable {
		t.Fatal("expected renewable")
	}
	if resp.LeaseDuration > 5 {
		t.Fatalf("lease duration too big: %d", resp.LeaseDuration)
	}
	leaseID := resp.LeaseID

	lastDuration := resp.LeaseDuration
	for i := 0; i < 3; i++ {
		time.Sleep(time.Second)
		resp, err = client.Sys().Renew(leaseID, 0)
		if err != nil {
			t.Fatal(err)
		}
		if resp.LeaseDuration >= lastDuration {
			t.Fatal("expected duration to go down")
		}
		lastDuration = resp.LeaseDuration
	}

	client.SetToken(rootToken)
	time.Sleep(2 * time.Second)
	resp, err = client.Logical().Write("sys/leases/lookup", map[string]interface{}{
		"lease_id": leaseID,
	})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestBatchToken_ParentLeaseRevoke(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"kv": vault.LeasedPassthroughBackendFactory,
		},
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
	rootToken := client.Token()
	var err error

	// Set up a KV path
	err = client.Sys().Mount("kv", &api.MountInput{
		Type: "kv",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Write("kv/foo", map[string]interface{}{
		"foo": "bar",
		"ttl": "5m",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Write the test policy
	err = client.Sys().PutPolicy("test", `
path "kv/*" {
	capabilities = ["read"]
}`)
	if err != nil {
		t.Fatal(err)
	}

	// Create a second root token
	secret, err := client.Auth().Token().Create(&api.TokenCreateRequest{
		Policies: []string{"root"},
	})
	if err != nil {
		t.Fatal(err)
	}
	rootToken2 := secret.Auth.ClientToken

	// Use this new token to create a batch token
	client.SetToken(rootToken2)
	secret, err = client.Auth().Token().Create(&api.TokenCreateRequest{
		Policies: []string{"test"},
		Type:     "batch",
	})
	if err != nil {
		t.Fatal(err)
	}
	batchToken := secret.Auth.ClientToken
	client.SetToken(batchToken)
	_, err = client.Auth().Token().LookupSelf()
	if err != nil {
		t.Fatal(err)
	}
	if secret.Auth.ClientToken[0:2] != "b." {
		t.Fatal(secret.Auth.ClientToken)
	}

	// Get a lease with the batch token
	resp, err := client.Logical().Read("kv/foo")
	if err != nil {
		t.Fatal(err)
	}
	if resp.Data["foo"].(string) != "bar" {
		t.Fatal("bad")
	}
	if resp.LeaseID == "" {
		t.Fatal("expected lease")
	}
	leaseID := resp.LeaseID

	// Check the lease
	resp, err = client.Logical().Write("sys/leases/lookup", map[string]interface{}{
		"lease_id": leaseID,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Revoke the parent
	client.SetToken(rootToken2)
	err = client.Auth().Token().RevokeSelf("")
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(1 * time.Second)

	// Verify the batch token is not usable anymore
	client.SetToken(rootToken)
	_, err = client.Auth().Token().Lookup(batchToken)
	if err == nil {
		t.Fatal("expected error")
	}

	// Verify the lease has been revoked
	resp, err = client.Logical().Write("sys/leases/lookup", map[string]interface{}{
		"lease_id": leaseID,
	})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestTokenStore_Roles_Batch(t *testing.T) {
	cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	core := cluster.Cores[0].Core
	vault.TestWaitActive(t, core)
	client := cluster.Cores[0].Client
	rootToken := client.Token()

	var err error
	var secret *api.Secret

	// Test service
	{
		_, err = client.Logical().Write("auth/token/roles/testrole", map[string]interface{}{
			"bound_cidrs": []string{},
			"token_type":  "service",
		})
		if err != nil {
			t.Fatal(err)
		}
		secret, err = client.Auth().Token().CreateWithRole(&api.TokenCreateRequest{
			Policies: []string{"default"},
			Type:     "batch",
		}, "testrole")
		if err != nil {
			t.Fatal(err)
		}
		client.SetToken(secret.Auth.ClientToken)
		_, err = client.Auth().Token().LookupSelf()
		if err != nil {
			t.Fatal(err)
		}
		if secret.Auth.ClientToken[0:2] != "s." {
			t.Fatal(secret.Auth.ClientToken)
		}
	}

	// Test batch
	{
		client.SetToken(rootToken)
		_, err = client.Logical().Write("auth/token/roles/testrole", map[string]interface{}{
			"token_type": "batch",
		})
		// Orphan not set so we should error
		if err == nil {
			t.Fatal("expected error")
		}
		_, err = client.Logical().Write("auth/token/roles/testrole", map[string]interface{}{
			"token_type": "batch",
			"orphan":     true,
		})
		// Renewable set so we should error
		if err == nil {
			t.Fatal("expected error")
		}
		_, err = client.Logical().Write("auth/token/roles/testrole", map[string]interface{}{
			"token_type": "batch",
			"orphan":     true,
			"renewable":  false,
		})
		if err != nil {
			t.Fatal(err)
		}
		secret, err = client.Auth().Token().CreateWithRole(&api.TokenCreateRequest{
			Policies: []string{"default"},
			Type:     "service",
		}, "testrole")
		if err != nil {
			t.Fatal(err)
		}
		client.SetToken(secret.Auth.ClientToken)
		_, err = client.Auth().Token().LookupSelf()
		if err != nil {
			t.Fatal(err)
		}
		if secret.Auth.ClientToken[0:2] != "b." {
			t.Fatal(secret.Auth.ClientToken)
		}
	}

	// Test default-service
	{
		client.SetToken(rootToken)
		_, err = client.Logical().Write("auth/token/roles/testrole", map[string]interface{}{
			"token_type": "default-service",
		})
		if err != nil {
			t.Fatal(err)
		}
		// Client specifies batch
		secret, err = client.Auth().Token().CreateWithRole(&api.TokenCreateRequest{
			Policies: []string{"default"},
			Type:     "batch",
		}, "testrole")
		if err != nil {
			t.Fatal(err)
		}
		client.SetToken(secret.Auth.ClientToken)
		_, err = client.Auth().Token().LookupSelf()
		if err != nil {
			t.Fatal(err)
		}
		if secret.Auth.ClientToken[0:2] != "b." {
			t.Fatal(secret.Auth.ClientToken)
		}
		// Client specifies service
		client.SetToken(rootToken)
		secret, err = client.Auth().Token().CreateWithRole(&api.TokenCreateRequest{
			Policies: []string{"default"},
			Type:     "service",
		}, "testrole")
		if err != nil {
			t.Fatal(err)
		}
		client.SetToken(secret.Auth.ClientToken)
		_, err = client.Auth().Token().LookupSelf()
		if err != nil {
			t.Fatal(err)
		}
		if secret.Auth.ClientToken[0:2] != "s." {
			t.Fatal(secret.Auth.ClientToken)
		}
		// Client doesn't specify
		client.SetToken(rootToken)
		secret, err = client.Auth().Token().CreateWithRole(&api.TokenCreateRequest{
			Policies: []string{"default"},
		}, "testrole")
		if err != nil {
			t.Fatal(err)
		}
		client.SetToken(secret.Auth.ClientToken)
		_, err = client.Auth().Token().LookupSelf()
		if err != nil {
			t.Fatal(err)
		}
		if secret.Auth.ClientToken[0:2] != "s." {
			t.Fatal(secret.Auth.ClientToken)
		}
	}

	// Test default-batch
	{
		client.SetToken(rootToken)
		_, err = client.Logical().Write("auth/token/roles/testrole", map[string]interface{}{
			"token_type": "default-batch",
		})
		if err != nil {
			t.Fatal(err)
		}
		// Client specifies batch
		secret, err = client.Auth().Token().CreateWithRole(&api.TokenCreateRequest{
			Policies: []string{"default"},
			Type:     "batch",
		}, "testrole")
		if err != nil {
			t.Fatal(err)
		}
		client.SetToken(secret.Auth.ClientToken)
		_, err = client.Auth().Token().LookupSelf()
		if err != nil {
			t.Fatal(err)
		}
		if secret.Auth.ClientToken[0:2] != "b." {
			t.Fatal(secret.Auth.ClientToken)
		}
		// Client specifies service
		client.SetToken(rootToken)
		secret, err = client.Auth().Token().CreateWithRole(&api.TokenCreateRequest{
			Policies: []string{"default"},
			Type:     "service",
		}, "testrole")
		if err != nil {
			t.Fatal(err)
		}
		client.SetToken(secret.Auth.ClientToken)
		_, err = client.Auth().Token().LookupSelf()
		if err != nil {
			t.Fatal(err)
		}
		if secret.Auth.ClientToken[0:2] != "s." {
			t.Fatal(secret.Auth.ClientToken)
		}
		// Client doesn't specify
		client.SetToken(rootToken)
		secret, err = client.Auth().Token().CreateWithRole(&api.TokenCreateRequest{
			Policies: []string{"default"},
		}, "testrole")
		if err != nil {
			t.Fatal(err)
		}
		client.SetToken(secret.Auth.ClientToken)
		_, err = client.Auth().Token().LookupSelf()
		if err != nil {
			t.Fatal(err)
		}
		if secret.Auth.ClientToken[0:2] != "b." {
			t.Fatal(secret.Auth.ClientToken)
		}
	}
}
