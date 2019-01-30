package vault

import (
	"strings"
	"testing"

	credLdap "github.com/hashicorp/vault/builtin/credential/ldap"
	credUserpass "github.com/hashicorp/vault/builtin/credential/userpass"
	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/logical"
)

func TestIdentityStore_CaseInsensitiveGroupAliasName(t *testing.T) {
	ctx := namespace.RootContext(nil)
	i, accessor, _ := testIdentityStoreWithGithubAuth(ctx, t)

	// Create a group
	resp, err := i.HandleRequest(ctx, &logical.Request{
		Path:      "group",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"type": "external",
		},
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v\nresp: %#v", err, resp)
	}
	groupID := resp.Data["id"].(string)

	testAliasName := "testAliasName"

	// Create a case sensitive alias name
	resp, err = i.HandleRequest(ctx, &logical.Request{
		Path:      "group-alias",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"mount_accessor": accessor,
			"canonical_id":   groupID,
			"name":           testAliasName,
		},
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err:%v\nresp: %#v", err, resp)
	}
	aliasID := resp.Data["id"].(string)

	// Ensure that reading the alias returns case sensitive alias name
	resp, err = i.HandleRequest(ctx, &logical.Request{
		Path:      "group-alias/id/" + aliasID,
		Operation: logical.ReadOperation,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err:%v\nresp: %#v", err, resp)
	}
	aliasName := resp.Data["name"].(string)
	if aliasName != testAliasName {
		t.Fatalf("bad alias name; expected: %q, actual: %q", testAliasName, aliasName)
	}

	// Overwrite the alias using lower cased alias name. This shouldn't error.
	resp, err = i.HandleRequest(ctx, &logical.Request{
		Path:      "group-alias/id/" + aliasID,
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"mount_accessor": accessor,
			"canonical_id":   groupID,
			"name":           strings.ToLower(testAliasName),
		},
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err:%v\nresp: %#v", err, resp)
	}

	// Ensure that reading the alias returns lower cased alias name
	resp, err = i.HandleRequest(ctx, &logical.Request{
		Path:      "group-alias/id/" + aliasID,
		Operation: logical.ReadOperation,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err:%v\nresp: %#v", err, resp)
	}
	aliasName = resp.Data["name"].(string)
	if aliasName != strings.ToLower(testAliasName) {
		t.Fatalf("bad alias name; expected: %q, actual: %q", testAliasName, aliasName)
	}
}

func TestIdentityStore_EnsureNoDanglingGroupAlias(t *testing.T) {
	err := AddTestCredentialBackend("userpass", credUserpass.Factory)
	if err != nil {
		t.Fatal(err)
	}

	err = AddTestCredentialBackend("ldap", credLdap.Factory)
	if err != nil {
		t.Fatal(err)
	}

	c, _, _ := TestCoreUnsealed(t)

	ctx := namespace.RootContext(nil)

	userpassMe := &MountEntry{
		Table:       credentialTableType,
		Path:        "userpass/",
		Type:        "userpass",
		Description: "userpass",
	}
	err = c.enableCredential(ctx, userpassMe)
	if err != nil {
		t.Fatal(err)
	}

	ldapMe := &MountEntry{
		Table:       credentialTableType,
		Path:        "ldap/",
		Type:        "ldap",
		Description: "ldap",
	}
	err = c.enableCredential(ctx, ldapMe)
	if err != nil {
		t.Fatal(err)
	}

	// Create a group
	resp, err := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "group",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"type": "external",
		},
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v\n", resp, err)
	}
	groupID := resp.Data["id"].(string)

	// Add an alias to the group from the userpass auth method
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "group-alias",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"name":           "testgroupalias",
			"mount_accessor": userpassMe.Accessor,
			"canonical_id":   groupID,
		},
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v\n", resp, err)
	}
	userpassGroupAliasID := resp.Data["id"].(string)

	// Ensure that the alias is readable
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "group-alias/id/" + userpassGroupAliasID,
		Operation: logical.ReadOperation,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v\n", resp, err)
	}
	if resp == nil || resp.Data["id"].(string) != userpassGroupAliasID {
		t.Fatalf("failed to read userpass group alias")
	}

	// Attach a different alias to the same group, overriding the previous one
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "group-alias",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"name":           "testgroupalias",
			"mount_accessor": ldapMe.Accessor,
			"canonical_id":   groupID,
		},
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v\n", resp, err)
	}
	ldapGroupAliasID := resp.Data["id"].(string)

	// Ensure that the new alias is readable
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "group-alias/id/" + ldapGroupAliasID,
		Operation: logical.ReadOperation,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v\n", resp, err)
	}
	if resp == nil || resp.Data["id"].(string) != ldapGroupAliasID {
		t.Fatalf("failed to read ldap group alias")
	}

	// Ensure previous alias is gone
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "group-alias/id/" + userpassGroupAliasID,
		Operation: logical.ReadOperation,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v\n", resp, err)
	}
	if resp != nil {
		t.Fatalf("expected a nil response")
	}
}

func TestIdentityStore_GroupAliasDeletionOnGroupDeletion(t *testing.T) {
	var resp *logical.Response
	var err error

	ctx := namespace.RootContext(nil)
	i, accessor, _ := testIdentityStoreWithGithubAuth(ctx, t)

	resp, err = i.HandleRequest(ctx, &logical.Request{
		Path:      "group",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"type": "external",
		},
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v\n", resp, err)
	}
	groupID := resp.Data["id"].(string)

	resp, err = i.HandleRequest(ctx, &logical.Request{
		Path:      "group-alias",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"name":           "testgroupalias",
			"mount_accessor": accessor,
			"canonical_id":   groupID,
		},
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v\n", resp, err)
	}
	groupAliasID := resp.Data["id"].(string)

	resp, err = i.HandleRequest(ctx, &logical.Request{
		Path:      "group/id/" + groupID,
		Operation: logical.DeleteOperation,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
	}

	resp, err = i.HandleRequest(ctx, &logical.Request{
		Path:      "group-alias/id/" + groupAliasID,
		Operation: logical.ReadOperation,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
	}
	if resp != nil {
		t.Fatalf("expected a nil response")
	}
}

func TestIdentityStore_GroupAliases_CRUD(t *testing.T) {
	var resp *logical.Response
	var err error
	ctx := namespace.RootContext(nil)
	i, accessor, _ := testIdentityStoreWithGithubAuth(ctx, t)

	groupReq := &logical.Request{
		Path:      "group",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"type": "external",
		},
	}
	resp, err = i.HandleRequest(ctx, groupReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v\n", resp, err)
	}
	groupID := resp.Data["id"].(string)

	groupAliasReq := &logical.Request{
		Path:      "group-alias",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"name":           "testgroupalias",
			"mount_accessor": accessor,
			"canonical_id":   groupID,
			"mount_type":     "ldap",
		},
	}
	resp, err = i.HandleRequest(ctx, groupAliasReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v\n", resp, err)
	}
	groupAliasID := resp.Data["id"].(string)

	groupAliasReq.Path = "group-alias/id/" + groupAliasID
	groupAliasReq.Operation = logical.ReadOperation
	resp, err = i.HandleRequest(ctx, groupAliasReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v\n", resp, err)
	}

	if resp.Data["id"].(string) != groupAliasID {
		t.Fatalf("bad: group alias: %#v\n", resp.Data)
	}

	resp, err = i.HandleRequest(ctx, &logical.Request{
		Path:      "group-alias/id/" + groupAliasID,
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"name":           "testupdatedgroupaliasname",
			"mount_accessor": accessor,
			"canonical_id":   groupID,
			"mount_type":     "ldap",
		},
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v; resp: %#v", err, resp)
	}
	if resp.Data["id"].(string) != groupAliasID {
		t.Fatalf("bad: group alias: %#v\n", resp.Data)
	}

	groupAliasReq.Operation = logical.DeleteOperation
	resp, err = i.HandleRequest(ctx, groupAliasReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v\n", resp, err)
	}

	groupAliasReq.Operation = logical.ReadOperation
	resp, err = i.HandleRequest(ctx, groupAliasReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v\n", resp, err)
	}

	if resp != nil {
		t.Fatalf("failed to delete group alias")
	}
}

func TestIdentityStore_GroupAliases_MemDBIndexes(t *testing.T) {
	var err error
	ctx := namespace.RootContext(nil)
	i, accessor, _ := testIdentityStoreWithGithubAuth(ctx, t)

	group := &identity.Group{
		ID:   "testgroupid",
		Name: "testgroupname",
		Metadata: map[string]string{
			"testmetadatakey1": "testmetadatavalue1",
			"testmetadatakey2": "testmetadatavalue2",
		},
		Alias: &identity.Alias{
			ID:            "testgroupaliasid",
			Name:          "testalias",
			MountAccessor: accessor,
			CanonicalID:   "testgroupid",
			MountType:     "ldap",
		},
		ParentGroupIDs:  []string{"testparentgroupid1", "testparentgroupid2"},
		MemberEntityIDs: []string{"testentityid1", "testentityid2"},
		Policies:        []string{"testpolicy1", "testpolicy2"},
		BucketKeyHash:   i.groupPacker.BucketKeyHashByItemID("testgroupid"),
	}

	txn := i.db.Txn(true)
	defer txn.Abort()
	err = i.MemDBUpsertAliasInTxn(txn, group.Alias, true)
	if err != nil {
		t.Fatal(err)
	}
	err = i.MemDBUpsertGroupInTxn(txn, group)
	if err != nil {
		t.Fatal(err)
	}
	txn.Commit()

	alias, err := i.MemDBAliasByID("testgroupaliasid", false, true)
	if err != nil {
		t.Fatal(err)
	}
	if alias.ID != "testgroupaliasid" {
		t.Fatalf("bad: group alias: %#v\n", alias)
	}

	group, err = i.MemDBGroupByAliasID("testgroupaliasid", false)
	if err != nil {
		t.Fatal(err)
	}
	if group.ID != "testgroupid" {
		t.Fatalf("bad: group: %#v\n", group)
	}

	aliasByFactors, err := i.MemDBAliasByFactors(group.Alias.MountAccessor, group.Alias.Name, false, true)
	if err != nil {
		t.Fatal(err)
	}
	if aliasByFactors.ID != "testgroupaliasid" {
		t.Fatalf("bad: group alias: %#v\n", aliasByFactors)
	}
}

func TestIdentityStore_GroupAliases_AliasOnInternalGroup(t *testing.T) {
	var err error
	var resp *logical.Response

	ctx := namespace.RootContext(nil)
	i, accessor, _ := testIdentityStoreWithGithubAuth(ctx, t)

	groupReq := &logical.Request{
		Path:      "group",
		Operation: logical.UpdateOperation,
	}
	resp, err = i.HandleRequest(ctx, groupReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v; err: %v", resp, err)
	}
	groupID := resp.Data["id"].(string)

	aliasReq := &logical.Request{
		Path:      "group-alias",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"name":           "testname",
			"mount_accessor": accessor,
			"canonical_id":   groupID,
		},
	}
	resp, err = i.HandleRequest(ctx, aliasReq)
	if err != nil {
		t.Fatal(err)
	}
	if !resp.IsError() {
		t.Fatalf("expected an error")
	}
}
