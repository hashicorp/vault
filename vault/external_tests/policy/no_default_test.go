package policy

import (
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/credential/ldap"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

func TestPolicy_NoDefaultPolicy(t *testing.T) {
	var err error
	coreConfig := &vault.CoreConfig{
		DisableMlock: true,
		DisableCache: true,
		Logger:       hclog.NewNullLogger(),
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
	secret, err := client.Logical().Write("auth/ldap/config", map[string]interface{}{
		"url":                     "ldap://ldap.forumsys.com",
		"userattr":                "uid",
		"userdn":                  "dc=example,dc=com",
		"groupdn":                 "dc=example,dc=com",
		"binddn":                  "cn=read-only-admin,dc=example,dc=com",
		"token_no_default_policy": true,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create a local user in LDAP
	secret, err = client.Logical().Write("auth/ldap/users/tesla", map[string]interface{}{
		"policies": "foo",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Login with LDAP and create a token
	secret, err = client.Logical().Write("auth/ldap/login/tesla", map[string]interface{}{
		"password": "password",
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
	coreConfig := &vault.CoreConfig{
		DisableMlock: true,
		DisableCache: true,
		Logger:       hclog.NewNullLogger(),
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
	secret, err := client.Logical().Write("auth/ldap/config", map[string]interface{}{
		"url":       "ldap://ldap.forumsys.com",
		"userattr":  "uid",
		"userdn":    "dc=example,dc=com",
		"groupdn":   "dc=example,dc=com",
		"binddn":    "cn=read-only-admin,dc=example,dc=com",
		"token_ttl": "24h",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create a local user in LDAP without any policies configured
	secret, err = client.Logical().Write("auth/ldap/users/tesla", map[string]interface{}{})
	if err != nil {
		t.Fatal(err)
	}

	// Login with LDAP and create a token
	secret, err = client.Logical().Write("auth/ldap/login/tesla", map[string]interface{}{
		"password": "password",
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

	// Sleep a bit to let the lease elapse
	time.Sleep(3 * time.Second)

	// Renew the token
	secret, err = client.Logical().Write("auth/token/renew", map[string]interface{}{
		"token": token,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Verify that the lease renewal extended the duration properly. We give it a one
	// second leeway to prevent test failure in case the response is delayed.
	if secret.Auth.LeaseDuration <= 86399 {
		t.Fatalf("failed to renew lease, got: %v", secret.Auth.LeaseDuration)
	}
}
