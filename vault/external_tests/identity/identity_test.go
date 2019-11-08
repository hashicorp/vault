package identity

import (
	"testing"

	"github.com/go-ldap/ldap/v3"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	ldapcred "github.com/hashicorp/vault/builtin/credential/ldap"
	"github.com/hashicorp/vault/helper/namespace"
	ldaphelper "github.com/hashicorp/vault/helper/testhelpers/ldap"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/helper/ldaputil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

func TestIdentityStore_Integ_GroupAliases(t *testing.T) {
	var err error
	coreConfig := &vault.CoreConfig{
		DisableMlock: true,
		DisableCache: true,
		Logger:       log.NewNullLogger(),
		CredentialBackends: map[string]logical.Factory{
			"ldap": ldapcred.Factory,
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

	cleanup, cfg := ldaphelper.PrepareTestContainer(t, "latest")
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
	assertMember := func(groupName, groupID string, expectFound bool) {
		secret, err = client.Logical().Read("identity/group/id/" + groupID)
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
		if found != expectFound {
			negation := ""
			if !expectFound {
				negation = "not "
			}
			t.Fatalf("expected entity ID %q to %sbe part of %q group", entityID, negation, groupName)
		}
	}

	assertMember("ship_crew", shipCrewGroupID, false)
	assertMember("admin_staff", adminStaffGroupID, true)
	assertMember("devops", devopsGroupID, true)
	assertMember("engineer", devopsGroupID, true)

	// Now add Hermes to ship_crew
	{
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

		hermesDn := "cn=Hermes Conrad,ou=people,dc=planetexpress,dc=com"
		shipCrewDn := "cn=ship_crew,ou=people,dc=planetexpress,dc=com"
		ldapreq := ldap.ModifyRequest{DN: shipCrewDn}
		ldapreq.Add("member", []string{hermesDn})
		err = conn.Modify(&ldapreq)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Re-login with LDAP
	secret, err = client.Logical().Write("auth/ldap/login/hermes conrad", map[string]interface{}{
		"password": "hermes",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Hermes should now be in ship_crew external group
	assertMember("ship_crew", shipCrewGroupID, true)
	assertMember("admin_staff", adminStaffGroupID, true)
	assertMember("devops", devopsGroupID, true)
	assertMember("engineer", devopsGroupID, true)

	identityStore := cores[0].IdentityStore()

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

	assertMember("ship_crew", shipCrewGroupID, true)
	assertMember("admin_staff", adminStaffGroupID, true)
	assertMember("devops", devopsGroupID, true)
	assertMember("engineer", devopsGroupID, true)

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
