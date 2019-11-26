package policy

// This is TODO once tokenhelper is added to ldaputil
/*

import (
	"testing"

	"github.com/go-test/deep"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/credential/ldap"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

func TestNoDefaultPolicy(t *testing.T) {
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
*/
