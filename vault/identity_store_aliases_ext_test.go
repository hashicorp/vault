package vault_test

import (
	"testing"

	"github.com/hashicorp/vault/api"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"

	credLdap "github.com/hashicorp/vault/builtin/credential/ldap"
)

func TestIdentityStore_EntityAliasLocalMount(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"ldap": credLdap.Factory,
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

	// Create a local auth mount
	err := client.Sys().EnableAuthWithOptions("ldap", &api.EnableAuthOptions{
		Type:  "ldap",
		Local: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Extract out the mount accessor for LDAP auth
	auths, err := client.Sys().ListAuth()
	if err != nil {
		t.Fatal(err)
	}
	ldapMountAccessor := auths["ldap/"].Accessor

	// Create an entity
	secret, err := client.Logical().Write("identity/entity", nil)
	if err != nil {
		t.Fatal(err)
	}
	entityID := secret.Data["id"].(string)

	// Attempt to create an entity alias against a local mount should fail
	secret, err = client.Logical().Write("identity/entity-alias", map[string]interface{}{
		"name":           "testuser",
		"mount_accessor": ldapMountAccessor,
		"canonical_id":   entityID,
	})
	if err == nil {
		t.Fatalf("expected error since mount is local")
	}
}
