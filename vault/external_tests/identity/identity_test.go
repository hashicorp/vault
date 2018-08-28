package identity

import (
	"testing"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/credential/ldap"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
)

func TestIdentityStore_Integ_GroupAliases(t *testing.T) {
	var err error
	coreConfig := &vault.CoreConfig{
		DisableMlock: true,
		DisableCache: true,
		Logger:       log.NewNullLogger(),
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

	auth, err := client.Sys().ListAuth()
	if err != nil {
		t.Fatal(err)
	}

	accessor := auth["ldap/"].Accessor

	secret, err := client.Logical().Write("identity/group", map[string]interface{}{
		"type": "external",
		"name": "ldap_Italians",
	})
	if err != nil {
		t.Fatal(err)
	}
	italiansGroupID := secret.Data["id"].(string)

	secret, err = client.Logical().Write("identity/group", map[string]interface{}{
		"type": "external",
		"name": "ldap_Scientists",
	})
	if err != nil {
		t.Fatal(err)
	}
	scientistsGroupID := secret.Data["id"].(string)

	secret, err = client.Logical().Write("identity/group", map[string]interface{}{
		"type": "external",
		"name": "ldap_devops",
	})
	if err != nil {
		t.Fatal(err)
	}
	devopsGroupID := secret.Data["id"].(string)

	secret, err = client.Logical().Write("identity/group-alias", map[string]interface{}{
		"name":           "Italians",
		"canonical_id":   italiansGroupID,
		"mount_accessor": accessor,
	})
	if err != nil {
		t.Fatal(err)
	}

	secret, err = client.Logical().Write("identity/group-alias", map[string]interface{}{
		"name":           "Scientists",
		"canonical_id":   scientistsGroupID,
		"mount_accessor": accessor,
	})
	if err != nil {
		t.Fatal(err)
	}

	secret, err = client.Logical().Write("identity/group-alias", map[string]interface{}{
		"name":           "devops",
		"canonical_id":   devopsGroupID,
		"mount_accessor": accessor,
	})
	if err != nil {
		t.Fatal(err)
	}

	secret, err = client.Logical().Read("identity/group/id/" + italiansGroupID)
	if err != nil {
		t.Fatal(err)
	}
	aliasMap := secret.Data["alias"].(map[string]interface{})
	if aliasMap["canonical_id"] != italiansGroupID ||
		aliasMap["name"] != "Italians" ||
		aliasMap["mount_accessor"] != accessor {
		t.Fatalf("bad: group alias: %#v\n", aliasMap)
	}

	secret, err = client.Logical().Read("identity/group/id/" + scientistsGroupID)
	if err != nil {
		t.Fatal(err)
	}
	aliasMap = secret.Data["alias"].(map[string]interface{})
	if aliasMap["canonical_id"] != scientistsGroupID ||
		aliasMap["name"] != "Scientists" ||
		aliasMap["mount_accessor"] != accessor {
		t.Fatalf("bad: group alias: %#v\n", aliasMap)
	}

	// Configure LDAP auth backend
	secret, err = client.Logical().Write("auth/ldap/config", map[string]interface{}{
		"url":      "ldap://ldap.forumsys.com",
		"userattr": "uid",
		"userdn":   "dc=example,dc=com",
		"groupdn":  "dc=example,dc=com",
		"binddn":   "cn=read-only-admin,dc=example,dc=com",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create a local group in LDAP backend
	secret, err = client.Logical().Write("auth/ldap/groups/devops", map[string]interface{}{
		"policies": "default",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create a local group in LDAP backend
	secret, err = client.Logical().Write("auth/ldap/groups/engineers", map[string]interface{}{
		"policies": "default",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create a local user in LDAP
	secret, err = client.Logical().Write("auth/ldap/users/tesla", map[string]interface{}{
		"policies": "default",
		"groups":   "engineers,devops",
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
	entityID := secret.Data["entity_id"].(string)

	// Re-read the Scientists, Italians and devops group. This entity ID should have
	// been added to both of these groups by now.
	secret, err = client.Logical().Read("identity/group/id/" + italiansGroupID)
	if err != nil {
		t.Fatal(err)
	}
	groupMap := secret.Data
	found := false
	for _, entityIDRaw := range groupMap["member_entity_ids"].([]interface{}) {
		if entityIDRaw.(string) == entityID {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected entity ID %q to be part of Italians group", entityID)
	}

	secret, err = client.Logical().Read("identity/group/id/" + scientistsGroupID)
	if err != nil {
		t.Fatal(err)
	}
	groupMap = secret.Data
	found = false
	for _, entityIDRaw := range groupMap["member_entity_ids"].([]interface{}) {
		if entityIDRaw.(string) == entityID {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected entity ID %q to be part of Scientists group", entityID)
	}

	secret, err = client.Logical().Read("identity/group/id/" + devopsGroupID)
	if err != nil {
		t.Fatal(err)
	}
	groupMap = secret.Data
	found = false
	for _, entityIDRaw := range groupMap["member_entity_ids"].([]interface{}) {
		if entityIDRaw.(string) == entityID {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected entity ID %q to be part of devops group", entityID)
	}

	identityStore := cores[0].IdentityStore()

	group, err := identityStore.MemDBGroupByID(italiansGroupID, true)
	if err != nil {
		t.Fatal(err)
	}

	// Remove its member entities
	group.MemberEntityIDs = nil

	err = identityStore.UpsertGroup(group, true)
	if err != nil {
		t.Fatal(err)
	}

	group, err = identityStore.MemDBGroupByID(italiansGroupID, true)
	if err != nil {
		t.Fatal(err)
	}
	if group.MemberEntityIDs != nil {
		t.Fatalf("failed to remove entity ID from the group")
	}

	group, err = identityStore.MemDBGroupByID(scientistsGroupID, true)
	if err != nil {
		t.Fatal(err)
	}

	// Remove its member entities
	group.MemberEntityIDs = nil

	err = identityStore.UpsertGroup(group, true)
	if err != nil {
		t.Fatal(err)
	}

	group, err = identityStore.MemDBGroupByID(scientistsGroupID, true)
	if err != nil {
		t.Fatal(err)
	}
	if group.MemberEntityIDs != nil {
		t.Fatalf("failed to remove entity ID from the group")
	}

	group, err = identityStore.MemDBGroupByID(devopsGroupID, true)
	if err != nil {
		t.Fatal(err)
	}

	// Remove its member entities
	group.MemberEntityIDs = nil

	err = identityStore.UpsertGroup(group, true)
	if err != nil {
		t.Fatal(err)
	}

	group, err = identityStore.MemDBGroupByID(devopsGroupID, true)
	if err != nil {
		t.Fatal(err)
	}
	if group.MemberEntityIDs != nil {
		t.Fatalf("failed to remove entity ID from the group")
	}

	_, err = client.Auth().Token().Renew(token, 0)
	if err != nil {
		t.Fatal(err)
	}

	// EntityIDs should have been added to the groups again during renewal
	secret, err = client.Logical().Read("identity/group/id/" + italiansGroupID)
	if err != nil {
		t.Fatal(err)
	}
	groupMap = secret.Data
	found = false
	for _, entityIDRaw := range groupMap["member_entity_ids"].([]interface{}) {
		if entityIDRaw.(string) == entityID {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected entity ID %q to be part of Italians group", entityID)
	}

	secret, err = client.Logical().Read("identity/group/id/" + scientistsGroupID)
	if err != nil {
		t.Fatal(err)
	}
	groupMap = secret.Data
	found = false
	for _, entityIDRaw := range groupMap["member_entity_ids"].([]interface{}) {
		if entityIDRaw.(string) == entityID {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected entity ID %q to be part of scientists group", entityID)
	}

	secret, err = client.Logical().Read("identity/group/id/" + devopsGroupID)
	if err != nil {
		t.Fatal(err)
	}

	groupMap = secret.Data
	found = false
	for _, entityIDRaw := range groupMap["member_entity_ids"].([]interface{}) {
		if entityIDRaw.(string) == entityID {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected entity ID %q to be part of devops group", entityID)
	}

	// Remove user tesla from the devops group in LDAP backend
	secret, err = client.Logical().Write("auth/ldap/users/tesla", map[string]interface{}{
		"policies": "default",
		"groups":   "engineers",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Renewing the token now should remove its entity ID from the devops
	// group
	_, err = client.Auth().Token().Renew(token, 0)
	if err != nil {
		t.Fatal(err)
	}

	group, err = identityStore.MemDBGroupByID(devopsGroupID, true)
	if err != nil {
		t.Fatal(err)
	}
	if group.MemberEntityIDs != nil {
		t.Fatalf("failed to remove entity ID from the group")
	}
}
