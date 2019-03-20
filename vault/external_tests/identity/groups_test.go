package identity

import (
	"testing"

	"github.com/hashicorp/vault/api"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"

	"github.com/hashicorp/vault/builtin/credential/github"
	credLdap "github.com/hashicorp/vault/builtin/credential/ldap"
)

func TestIdentityStore_ListGroupAlias(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"github": github.Factory,
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

	err := client.Sys().EnableAuthWithOptions("github", &api.EnableAuthOptions{
		Type: "github",
	})
	if err != nil {
		t.Fatal(err)
	}

	mounts, err := client.Sys().ListAuth()
	if err != nil {
		t.Fatal(err)
	}
	var githubAccessor string
	for k, v := range mounts {
		t.Logf("key: %v\nmount: %#v", k, *v)
		if k == "github/" {
			githubAccessor = v.Accessor
			break
		}
	}
	if githubAccessor == "" {
		t.Fatal("did not find github accessor")
	}

	resp, err := client.Logical().Write("identity/group", map[string]interface{}{
		"type": "external",
	})
	if err != nil {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	groupID := resp.Data["id"].(string)

	resp, err = client.Logical().Write("identity/group-alias", map[string]interface{}{
		"name":           "groupalias",
		"mount_accessor": githubAccessor,
		"canonical_id":   groupID,
	})
	if err != nil {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	aliasID := resp.Data["id"].(string)

	resp, err = client.Logical().List("identity/group-alias/id")
	if err != nil {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	keys := resp.Data["keys"].([]interface{})
	if len(keys) != 1 {
		t.Fatalf("bad: length of alias IDs listed; expected: 1, actual: %d", len(keys))
	}

	// Do some due diligence on the key info
	aliasInfoRaw, ok := resp.Data["key_info"]
	if !ok {
		t.Fatal("expected key_info map in response")
	}
	aliasInfo := aliasInfoRaw.(map[string]interface{})
	if len(aliasInfo) != 1 {
		t.Fatalf("bad: length of alias ID key info; expected: 1, actual: %d", len(aliasInfo))
	}

	infoRaw, ok := aliasInfo[aliasID]
	if !ok {
		t.Fatal("expected to find alias ID in key info map")
	}
	info := infoRaw.(map[string]interface{})
	t.Logf("alias info: %#v", info)
	switch {
	case info["name"].(string) != "groupalias":
		t.Fatalf("bad name: %v", info["name"].(string))
	case info["mount_accessor"].(string) != githubAccessor:
		t.Fatalf("bad mount_accessor: %v", info["mount_accessor"].(string))
	}

	// Now do the same with group info
	resp, err = client.Logical().List("identity/group/id")
	if err != nil {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	keys = resp.Data["keys"].([]interface{})
	if len(keys) != 1 {
		t.Fatalf("bad: length of group IDs listed; expected: 1, actual: %d", len(keys))
	}

	groupInfoRaw, ok := resp.Data["key_info"]
	if !ok {
		t.Fatal("expected key_info map in response")
	}

	// This is basically verifying that the group has the alias in key_info
	// that we expect to be tied to it, plus tests a value further down in it
	// for fun
	groupInfo := groupInfoRaw.(map[string]interface{})
	if len(groupInfo) != 1 {
		t.Fatalf("bad: length of group ID key info; expected: 1, actual: %d", len(groupInfo))
	}

	infoRaw, ok = groupInfo[groupID]
	if !ok {
		t.Fatal("expected key info")
	}
	info = infoRaw.(map[string]interface{})
	t.Logf("group info: %#v", info)
	alias := info["alias"].(map[string]interface{})
	switch {
	case alias["id"].(string) != aliasID:
		t.Fatalf("bad alias id: %v", alias["id"])
	case alias["mount_accessor"].(string) != githubAccessor:
		t.Fatalf("bad mount accessor: %v", alias["mount_accessor"])
	case alias["mount_path"].(string) != "auth/github/":
		t.Fatalf("bad mount path: %v", alias["mount_path"])
	case alias["mount_type"].(string) != "github":
		t.Fatalf("bad mount type: %v", alias["mount_type"])
	}
}

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
	secret, err = client.Logical().Write("auth/token/lookup", map[string]interface{}{
		"token": ldapClientToken,
	})
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
