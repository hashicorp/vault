package policy

import (
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/credential/ldap"
	ldaphelper "github.com/hashicorp/vault/helper/testhelpers/ldap"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

func TestPolicy_NoDefaultPolicy(t *testing.T) {
	var err error
	coreConfig := &vault.CoreConfig{
		DisableMlock: true,
		DisableCache: true,
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
	cleanup, cfg := ldaphelper.PrepareTestContainer(t, "latest")
	defer cleanup()

	_, err = client.Logical().Write("auth/ldap/config", map[string]interface{}{
		"url":                     cfg.Url,
		"userattr":                cfg.UserAttr,
		"userdn":                  cfg.UserDN,
		"groupdn":                 cfg.GroupDN,
		"groupattr":               cfg.GroupAttr,
		"binddn":                  cfg.BindDN,
		"bindpass":                cfg.BindPassword,
		"token_no_default_policy": true,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create a local user in LDAP
	secret, err := client.Logical().Write("auth/ldap/users/hermes conrad", map[string]interface{}{
		"policies": "foo",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Login with LDAP and create a token
	secret, err = client.Logical().Write("auth/ldap/login/hermes conrad", map[string]interface{}{
		"password": "hermes",
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
	cleanup, cfg := ldaphelper.PrepareTestContainer(t, "latest")
	defer cleanup()

	_, err = client.Logical().Write("auth/ldap/config", map[string]interface{}{
		"url":       cfg.Url,
		"userattr":  cfg.UserAttr,
		"userdn":    cfg.UserDN,
		"groupdn":   cfg.GroupDN,
		"groupattr": cfg.GroupAttr,
		"binddn":    cfg.BindDN,
		"bindpass":  cfg.BindPassword,
		"token_ttl": "24h",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create a local user in LDAP without any policies configured
	secret, err := client.Logical().Write("auth/ldap/users/hermes conrad", map[string]interface{}{})
	if err != nil {
		t.Fatal(err)
	}

	// Login with LDAP and create a token
	secret, err = client.Logical().Write("auth/ldap/login/hermes conrad", map[string]interface{}{
		"password": "hermes",
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

	// Renew the token with an increment of 2 hours to ensure that lease renewal
	// occurred and can be checked against the default lease duration with a
	// big enough delta.
	secret, err = client.Logical().Write("auth/token/renew", map[string]interface{}{
		"token":     token,
		"increment": "2h",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Verify that the lease renewal extended the duration properly.
	if float64(secret.Auth.LeaseDuration) < (1 * time.Hour).Seconds() {
		t.Fatalf("failed to renew lease, got: %v", secret.Auth.LeaseDuration)
	}
}
