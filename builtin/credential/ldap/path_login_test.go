package ldap

import (
	"testing"

	"github.com/hashicorp/vault/api"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
)

func TestLdapAuthBackend_UsernameAliasCaseSensitivity(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"ldap": Factory,
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

	// Create an ldap mount
	err := client.Sys().EnableAuthWithOptions("ldap", &api.EnableAuthOptions{
		Type: "ldap",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Write("auth/ldap/config", map[string]interface{}{
		"url":      "ldap://ldap.forumsys.com",
		"userattr": "uid",
		"userdn":   "dc=example,dc=com",
		"groupdn":  "dc=example,dc=com",
		"binddn":   "cn=read-only-admin,dc=example,dc=com",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Write("auth/ldap/groups/testgroup", map[string]interface{}{
		"policies": "testgrouppolicy",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Write("auth/ldap/users/tesla", map[string]interface{}{
		"policies": "default",
		"groups":   "testgroup",
	})
	if err != nil {
		t.Fatal(err)
	}

	usernames := []string{
		"tesla",
		"Tesla",
		"teSlA",
	}

	// Login using different cases for usernames and ensure that only one
	// entity is getting created
	for _, username := range usernames {
		secret, err := client.Logical().Write("auth/ldap/login/"+username, map[string]interface{}{
			"password": "password",
		})
		if err != nil {
			t.Fatal(err)
		}
		if secret.Auth.ClientToken == "" {
			t.Fatalf("failed to perform login")
		}
		resp, err := client.Logical().List("identity/entity/id")
		if err != nil {
			t.Fatal(err)
		}
		if len(resp.Data["keys"].([]interface{})) != 1 {
			t.Fatalf("failed to list entities")
		}
	}
}
