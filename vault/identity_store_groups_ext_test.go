package vault_test

import (
	"testing"

	"github.com/hashicorp/vault/api"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"

	credLdap "github.com/hashicorp/vault/builtin/credential/ldap"
)

// Testing the fix for GH-4351
func TestIdentityStore_ExternalGroupMembershipsAcrossMounts(t *testing.T) {
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

	// Enable the first LDAP auth
	err := client.Sys().EnableAuthWithOptions("ldap", &api.EnableAuthOptions{
		Type: "ldap",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Extract out the mount accessor for LDAP auth
	auths, err := client.Sys().ListAuth()
	if err != nil {
		t.Fatal(err)
	}
	ldapMountAccessor1 := auths["ldap/"].Accessor

	// Configure LDAP auth
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

	// Create a group in LDAP auth
	_, err = client.Logical().Write("auth/ldap/groups/testgroup1", map[string]interface{}{
		"policies": "testgroup1-policy",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Tie the group to a user
	_, err = client.Logical().Write("auth/ldap/users/tesla", map[string]interface{}{
		"policies": "default",
		"groups":   "testgroup1",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create an external group
	secret, err := client.Logical().Write("identity/group", map[string]interface{}{
		"type": "external",
	})
	if err != nil {
		t.Fatal(err)
	}
	ldapExtGroupID1 := secret.Data["id"].(string)

	// Associate a group from LDAP auth as a group-alias in the external group
	_, err = client.Logical().Write("identity/group-alias", map[string]interface{}{
		"name":           "testgroup1",
		"mount_accessor": ldapMountAccessor1,
		"canonical_id":   ldapExtGroupID1,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Login using LDAP
	secret, err = client.Logical().Write("auth/ldap/login/tesla", map[string]interface{}{
		"password": "password",
	})
	if err != nil {
		t.Fatal(err)
	}
	ldapClientToken := secret.Auth.ClientToken

	//
	// By now, the entity ID of the token should be automatically added as a
	// member in the external group.
	//

	// Extract the entity ID of the token
	secret, err = client.Logical().Read("auth/token/lookup/" + ldapClientToken)
	if err != nil {
		t.Fatal(err)
	}
	entityID := secret.Data["entity_id"].(string)

	// Enable another LDAP auth mount
	err = client.Sys().EnableAuthWithOptions("ldap2", &api.EnableAuthOptions{
		Type: "ldap",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Extract the mount accessor
	auths, err = client.Sys().ListAuth()
	if err != nil {
		t.Fatal(err)
	}
	ldapMountAccessor2 := auths["ldap2/"].Accessor

	// Create an entity-alias asserting that the user "tesla" from the first
	// and second LDAP mounts as the same.
	_, err = client.Logical().Write("identity/entity-alias", map[string]interface{}{
		"name":           "tesla",
		"mount_accessor": ldapMountAccessor2,
		"canonical_id":   entityID,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Configure second LDAP auth
	_, err = client.Logical().Write("auth/ldap2/config", map[string]interface{}{
		"url":      "ldap://ldap.forumsys.com",
		"userattr": "uid",
		"userdn":   "dc=example,dc=com",
		"groupdn":  "dc=example,dc=com",
		"binddn":   "cn=read-only-admin,dc=example,dc=com",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create a group in second LDAP auth
	_, err = client.Logical().Write("auth/ldap2/groups/testgroup2", map[string]interface{}{
		"policies": "testgroup2-policy",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create a user in second LDAP auth
	_, err = client.Logical().Write("auth/ldap2/users/tesla", map[string]interface{}{
		"policies": "default",
		"groups":   "testgroup2",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create another external group
	secret, err = client.Logical().Write("identity/group", map[string]interface{}{
		"type": "external",
	})
	if err != nil {
		t.Fatal(err)
	}
	ldapExtGroupID2 := secret.Data["id"].(string)

	// Create a group-alias tying the external group to "testgroup2" group in second LDAP
	_, err = client.Logical().Write("identity/group-alias", map[string]interface{}{
		"name":           "testgroup2",
		"mount_accessor": ldapMountAccessor2,
		"canonical_id":   ldapExtGroupID2,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Login using second LDAP
	_, err = client.Logical().Write("auth/ldap2/login/tesla", map[string]interface{}{
		"password": "password",
	})
	if err != nil {
		t.Fatal(err)
	}

	//
	// By now the same entity ID of the token from first LDAP should have been
	// added as a member of the second external group.
	//

	// Check that entityID is present in both the external groups
	secret, err = client.Logical().Read("identity/group/id/" + ldapExtGroupID1)
	if err != nil {
		t.Fatal(err)
	}
	extGroup1Entities := secret.Data["member_entity_ids"].([]interface{})

	found := false
	for _, item := range extGroup1Entities {
		if item.(string) == entityID {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("missing entity ID %q first external group with ID %q", entityID, ldapExtGroupID1)
	}

	secret, err = client.Logical().Read("identity/group/id/" + ldapExtGroupID2)
	if err != nil {
		t.Fatal(err)
	}
	extGroup2Entities := secret.Data["member_entity_ids"].([]interface{})
	found = false
	for _, item := range extGroup2Entities {
		if item.(string) == entityID {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("missing entity ID %q first external group with ID %q", entityID, ldapExtGroupID2)
	}
}
