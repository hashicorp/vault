// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package identity

import (
	"fmt"
	"testing"

	"github.com/go-ldap/ldap/v3"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/namespace"
	ldaphelper "github.com/hashicorp/vault/helper/testhelpers/ldap"
	"github.com/hashicorp/vault/helper/testhelpers/minimal"
	"github.com/hashicorp/vault/sdk/helper/ldaputil"
	"github.com/stretchr/testify/require"
)

func TestIdentityStore_ExternalGroupMemberships_DifferentMounts(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	// Create a entity
	secret, err := client.Logical().Write("identity/entity", map[string]interface{}{
		"name": "testentityname",
	})
	require.NoError(t, err)
	entityID := secret.Data["id"].(string)

	cleanup, config1 := ldaphelper.PrepareTestContainer(t, "master")
	defer cleanup()

	cleanup2, config2 := ldaphelper.PrepareTestContainer(t, "master")
	defer cleanup2()

	setupFunc := func(path string, cfg *ldaputil.ConfigEntry) string {
		// Create an external group
		resp, err := client.Logical().Write("identity/group", map[string]interface{}{
			"type":     "external",
			"name":     path + "ldap_admin_staff",
			"policies": []string{"admin-policy"},
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotNil(t, resp.Data)
		groupID := resp.Data["id"].(string)

		// Enable LDAP mount in Vault
		err = client.Sys().EnableAuthWithOptions(path, &api.EnableAuthOptions{
			Type: "ldap",
		})
		require.NoError(t, err)

		// Take out its accessor
		auth, err := client.Sys().ListAuth()
		require.NoError(t, err)
		accessor := auth[path+"/"].Accessor
		require.NotEmpty(t, accessor)

		// Create an external group alias
		resp, err = client.Logical().Write("identity/group-alias", map[string]interface{}{
			"name":           "admin_staff",
			"canonical_id":   groupID,
			"mount_accessor": accessor,
		})
		require.NoError(t, err)

		// Create a user in Vault
		_, err = client.Logical().Write("auth/"+path+"/users/hermes conrad", map[string]interface{}{
			"password": "hermes",
		})
		require.NoError(t, err)

		// Create an entity alias
		client.Logical().Write("identity/entity-alias", map[string]interface{}{
			"name":           "hermes conrad",
			"canonical_id":   entityID,
			"mount_accessor": accessor,
		})

		// Configure LDAP auth
		secret, err = client.Logical().Write("auth/"+path+"/config", map[string]interface{}{
			"url":       cfg.Url,
			"userattr":  cfg.UserAttr,
			"userdn":    cfg.UserDN,
			"groupdn":   cfg.GroupDN,
			"groupattr": cfg.GroupAttr,
			"binddn":    cfg.BindDN,
			"bindpass":  cfg.BindPassword,
		})
		require.NoError(t, err)

		secret, err = client.Logical().Write("auth/"+path+"/login/hermes conrad", map[string]interface{}{
			"password": "hermes",
		})
		require.NoError(t, err)

		policies, err := secret.TokenPolicies()
		require.NoError(t, err)
		require.Contains(t, policies, "admin-policy")

		secret, err = client.Logical().Read("identity/group/id/" + groupID)
		require.NoError(t, err)
		require.Contains(t, secret.Data["member_entity_ids"], entityID)

		return groupID
	}
	groupID1 := setupFunc("ldap", config1)
	groupID2 := setupFunc("ldap2", config2)

	// Remove hermes conrad from admin_staff group
	removeLdapGroupMember(t, config1, "admin_staff", "hermes conrad")
	secret, err = client.Logical().Write("auth/ldap/login/hermes conrad", map[string]interface{}{
		"password": "hermes",
	})
	require.NoError(t, err)

	secret, err = client.Logical().Read("identity/group/id/" + groupID1)
	require.NoError(t, err)
	require.NotContains(t, secret.Data["member_entity_ids"], entityID)

	secret, err = client.Logical().Read("identity/group/id/" + groupID2)
	require.NoError(t, err)
	require.Contains(t, secret.Data["member_entity_ids"], entityID)
}

func TestIdentityStore_Integ_GroupAliases(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	err := client.Sys().EnableAuthWithOptions("ldap", &api.EnableAuthOptions{
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
		"name": "ldap_ship_crew",
	})
	if err != nil {
		t.Fatal(err)
	}
	shipCrewGroupID := secret.Data["id"].(string)

	secret, err = client.Logical().Write("identity/group", map[string]interface{}{
		"type": "external",
		"name": "ldap_admin_staff",
	})
	if err != nil {
		t.Fatal(err)
	}
	adminStaffGroupID := secret.Data["id"].(string)

	secret, err = client.Logical().Write("identity/group", map[string]interface{}{
		"type": "external",
		"name": "ldap_devops",
	})
	if err != nil {
		t.Fatal(err)
	}
	devopsGroupID := secret.Data["id"].(string)

	secret, err = client.Logical().Write("identity/group-alias", map[string]interface{}{
		"name":           "ship_crew",
		"canonical_id":   shipCrewGroupID,
		"mount_accessor": accessor,
	})
	if err != nil {
		t.Fatal(err)
	}

	secret, err = client.Logical().Write("identity/group-alias", map[string]interface{}{
		"name":           "admin_staff",
		"canonical_id":   adminStaffGroupID,
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

	secret, err = client.Logical().Read("identity/group/id/" + shipCrewGroupID)
	if err != nil {
		t.Fatal(err)
	}
	aliasMap := secret.Data["alias"].(map[string]interface{})
	if aliasMap["canonical_id"] != shipCrewGroupID ||
		aliasMap["name"] != "ship_crew" ||
		aliasMap["mount_accessor"] != accessor {
		t.Fatalf("bad: group alias: %#v\n", aliasMap)
	}

	secret, err = client.Logical().Read("identity/group/id/" + adminStaffGroupID)
	if err != nil {
		t.Fatal(err)
	}
	aliasMap = secret.Data["alias"].(map[string]interface{})
	if aliasMap["canonical_id"] != adminStaffGroupID ||
		aliasMap["name"] != "admin_staff" ||
		aliasMap["mount_accessor"] != accessor {
		t.Fatalf("bad: group alias: %#v\n", aliasMap)
	}

	cleanup, cfg := ldaphelper.PrepareTestContainer(t, "master")
	defer cleanup()

	// Configure LDAP auth
	secret, err = client.Logical().Write("auth/ldap/config", map[string]interface{}{
		"url":       cfg.Url,
		"userattr":  cfg.UserAttr,
		"userdn":    cfg.UserDN,
		"groupdn":   cfg.GroupDN,
		"groupattr": cfg.GroupAttr,
		"binddn":    cfg.BindDN,
		"bindpass":  cfg.BindPassword,
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
	secret, err = client.Logical().Write("auth/ldap/users/hermes conrad", map[string]interface{}{
		"policies": "default",
		"groups":   "engineers,devops",
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
	entityID := secret.Data["entity_id"].(string)

	// Re-read the admin_staff, ship_crew and devops group. This entity ID should have
	// been added to admin_staff but not ship_crew.
	assertMember(t, client, entityID, "ship_crew", shipCrewGroupID, false)
	assertMember(t, client, entityID, "admin_staff", adminStaffGroupID, true)
	assertMember(t, client, entityID, "devops", devopsGroupID, true)
	assertMember(t, client, entityID, "engineer", devopsGroupID, true)

	// Now add Hermes to ship_crew
	addLdapGroupMember(t, cfg, "ship_crew", "hermes conrad")

	// Re-login with LDAP
	secret, err = client.Logical().Write("auth/ldap/login/hermes conrad", map[string]interface{}{
		"password": "hermes",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Hermes should now be in ship_crew external group
	assertMember(t, client, entityID, "ship_crew", shipCrewGroupID, true)
	assertMember(t, client, entityID, "admin_staff", adminStaffGroupID, true)
	assertMember(t, client, entityID, "devops", devopsGroupID, true)
	assertMember(t, client, entityID, "engineer", devopsGroupID, true)

	identityStore := cluster.Cores[0].IdentityStore()

	group, err := identityStore.MemDBGroupByID(shipCrewGroupID, true)
	if err != nil {
		t.Fatal(err)
	}

	// Remove its member entities
	group.MemberEntityIDs = nil

	ctx := namespace.RootContext(nil)

	err = identityStore.UpsertGroup(ctx, group, true)
	if err != nil {
		t.Fatal(err)
	}

	group, err = identityStore.MemDBGroupByID(shipCrewGroupID, true)
	if err != nil {
		t.Fatal(err)
	}
	if group.MemberEntityIDs != nil {
		t.Fatalf("failed to remove entity ID from the group")
	}

	group, err = identityStore.MemDBGroupByID(adminStaffGroupID, true)
	if err != nil {
		t.Fatal(err)
	}

	// Remove its member entities
	group.MemberEntityIDs = nil

	err = identityStore.UpsertGroup(ctx, group, true)
	if err != nil {
		t.Fatal(err)
	}

	group, err = identityStore.MemDBGroupByID(adminStaffGroupID, true)
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

	err = identityStore.UpsertGroup(ctx, group, true)
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

	assertMember(t, client, entityID, "ship_crew", shipCrewGroupID, true)
	assertMember(t, client, entityID, "admin_staff", adminStaffGroupID, true)
	assertMember(t, client, entityID, "devops", devopsGroupID, true)
	assertMember(t, client, entityID, "engineer", devopsGroupID, true)

	// Remove user hermes conrad from the devops group in LDAP backend
	secret, err = client.Logical().Write("auth/ldap/users/hermes conrad", map[string]interface{}{
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

func TestIdentityStore_Integ_RemoveFromExternalGroup(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	err := client.Sys().EnableAuthWithOptions("ldap", &api.EnableAuthOptions{
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

	adminPolicy := "admin_policy"
	secret, err := client.Logical().Write("identity/group", map[string]interface{}{
		"type":     "external",
		"name":     "ldap_admin_staff",
		"policies": []string{adminPolicy},
	})
	if err != nil {
		t.Fatal(err)
	}
	adminStaffGroupID := secret.Data["id"].(string)
	adminGroupName := "admin_staff"

	secret, err = client.Logical().Write("identity/group-alias", map[string]interface{}{
		"name":           adminGroupName,
		"canonical_id":   adminStaffGroupID,
		"mount_accessor": accessor,
	})
	if err != nil {
		t.Fatal(err)
	}

	secret, err = client.Logical().Read("identity/group/id/" + adminStaffGroupID)
	if err != nil {
		t.Fatal(err)
	}
	aliasMap := secret.Data["alias"].(map[string]interface{})
	if aliasMap["canonical_id"] != adminStaffGroupID ||
		aliasMap["name"] != adminGroupName ||
		aliasMap["mount_accessor"] != accessor {
		t.Fatalf("bad: group alias: %#v\n", aliasMap)
	}

	cleanup, cfg := ldaphelper.PrepareTestContainer(t, "master")
	defer cleanup()

	// Configure LDAP auth
	secret, err = client.Logical().Write("auth/ldap/config", map[string]interface{}{
		"url":       cfg.Url,
		"userattr":  cfg.UserAttr,
		"userdn":    cfg.UserDN,
		"groupdn":   cfg.GroupDN,
		"groupattr": cfg.GroupAttr,
		"binddn":    cfg.BindDN,
		"bindpass":  cfg.BindPassword,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create a local user in LDAP
	secret, err = client.Logical().Write("auth/ldap/users/hermes conrad", map[string]interface{}{
		"policies": "default",
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
	tokenPolicies, err := secret.TokenPolicies()
	if err != nil {
		t.Fatal(err)
	}
	if !strutil.StrListContains(tokenPolicies, adminPolicy) {
		t.Fatalf("expected token policies to contain %s, got: %v", adminPolicy, tokenPolicies)
	}

	// Lookup the token to get the entity ID
	secret, err = client.Auth().Token().Lookup(token)
	if err != nil {
		t.Fatal(err)
	}
	entityID := secret.Data["entity_id"].(string)

	assertMember(t, client, entityID, adminGroupName, adminStaffGroupID, true)

	// Now remove Hermes from admin_staff
	removeLdapGroupMember(t, cfg, adminGroupName, "hermes conrad")

	// Re-login with LDAP
	secret, err = client.Logical().Write("auth/ldap/login/hermes conrad", map[string]interface{}{
		"password": "hermes",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Hermes should now be out of admin_staff group
	assertMember(t, client, entityID, adminGroupName, adminStaffGroupID, false)
	tokenPolicies, err = secret.TokenPolicies()
	if err != nil {
		t.Fatal(err)
	}
	if strutil.StrListContains(tokenPolicies, adminPolicy) {
		t.Fatalf("expected token policies to not contain %s, got: %v", adminPolicy, tokenPolicies)
	}

	// Add Hermes back to admin_staff
	addLdapGroupMember(t, cfg, adminGroupName, "hermes conrad")

	// Re-login with LDAP
	secret, err = client.Logical().Write("auth/ldap/login/hermes conrad", map[string]interface{}{
		"password": "hermes",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Hermes should now be back in admin_staff group
	assertMember(t, client, entityID, adminGroupName, adminStaffGroupID, true)
	tokenPolicies, err = secret.TokenPolicies()
	if err != nil {
		t.Fatal(err)
	}
	if !strutil.StrListContains(tokenPolicies, adminPolicy) {
		t.Fatalf("expected token policies to contain %s, got: %v", adminPolicy, tokenPolicies)
	}

	// Remove Hermes from admin_staff once again
	removeLdapGroupMember(t, cfg, adminGroupName, "hermes conrad")

	oldToken := client.Token()
	client.SetToken(secret.Auth.ClientToken)
	secret, err = client.Auth().Token().RenewSelf(1)
	if err != nil {
		t.Fatal(err)
	}
	client.SetToken(oldToken)
	assertMember(t, client, entityID, adminGroupName, adminStaffGroupID, false)
	tokenPolicies, err = secret.TokenPolicies()
	if err != nil {
		t.Fatal(err)
	}
	if strutil.StrListContains(tokenPolicies, adminPolicy) {
		t.Fatalf("expected token policies to not contain %s, got: %v", adminPolicy, tokenPolicies)
	}
}

func assertMember(t *testing.T, client *api.Client, entityID, groupName, groupID string, expectFound bool) {
	t.Helper()
	secret, err := client.Logical().Read("identity/group/id/" + groupID)
	if err != nil {
		t.Fatal(err)
	}
	groupMap := secret.Data

	groupEntityMembers, ok := groupMap["member_entity_ids"].([]interface{})
	if !ok && expectFound {
		t.Fatalf("expected member_entity_ids not to be nil")
	}

	// if type assertion fails and expectFound is false, groupEntityMembers
	// is nil, then let's just return, nothing to be done!
	if !ok && !expectFound {
		return
	}

	found := false
	for _, entityIDRaw := range groupEntityMembers {
		if entityIDRaw.(string) == entityID {
			found = true
		}
	}
	if found != expectFound {
		negation := ""
		if !expectFound {
			negation = "not "
		}
		t.Fatalf("expected entity ID %q to %sbe part of %q group", entityID, negation, groupName)
	}
}

func removeLdapGroupMember(t *testing.T, cfg *ldaputil.ConfigEntry, groupCN, userCN string) {
	userDN := fmt.Sprintf("cn=%s,ou=people,dc=planetexpress,dc=com", userCN)
	groupDN := fmt.Sprintf("cn=%s,ou=people,dc=planetexpress,dc=com", groupCN)
	ldapreq := ldap.ModifyRequest{DN: groupDN}
	ldapreq.Delete("member", []string{userDN})
	addRemoveLdapGroupMember(t, cfg, userCN, &ldapreq)
}

func addLdapGroupMember(t *testing.T, cfg *ldaputil.ConfigEntry, groupCN, userCN string) {
	userDN := fmt.Sprintf("cn=%s,ou=people,dc=planetexpress,dc=com", userCN)
	groupDN := fmt.Sprintf("cn=%s,ou=people,dc=planetexpress,dc=com", groupCN)
	ldapreq := ldap.ModifyRequest{DN: groupDN}
	ldapreq.Add("member", []string{userDN})
	addRemoveLdapGroupMember(t, cfg, userCN, &ldapreq)
}

func addRemoveLdapGroupMember(t *testing.T, cfg *ldaputil.ConfigEntry, userCN string, req *ldap.ModifyRequest) {
	logger := log.New(nil)
	ldapClient := ldaputil.Client{LDAP: ldaputil.NewLDAP(), Logger: logger}
	// LDAP server won't accept changes unless we connect with TLS.  This
	// isn't the default config returned by PrepareTestContainer because
	// the Vault LDAP backend won't work with it, even with InsecureTLS,
	// because the ServerName should be planetexpress.com and not localhost.
	conn, err := ldapClient.DialLDAP(cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	err = conn.Bind(cfg.BindDN, cfg.BindPassword)
	if err != nil {
		t.Fatal(err)
	}

	err = conn.Modify(req)
	if err != nil {
		t.Fatal(err)
	}
}
