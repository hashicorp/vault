package vault

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/logical"
)

func TestIdentityStore_ListGroupAlias(t *testing.T) {
	var err error
	var resp *logical.Response

	is, githubAccessor, _ := testIdentityStoreWithGithubAuth(t)

	groupReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "group",
		Data: map[string]interface{}{
			"type": "external",
		},
	}
	resp, err = is.HandleRequest(context.Background(), groupReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	if resp == nil {
		t.Fatalf("expected a non-nil response")
	}
	groupID := resp.Data["id"].(string)

	// Create an alias
	aliasData := map[string]interface{}{
		"name":           "groupalias",
		"mount_accessor": githubAccessor,
		"canonical_id":   groupID,
	}
	aliasReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "group-alias",
		Data:      aliasData,
	}
	resp, err = is.HandleRequest(context.Background(), aliasReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	aliasID := resp.Data["id"].(string)

	listReq := &logical.Request{
		Operation: logical.ListOperation,
		Path:      "group-alias/id",
	}
	resp, err = is.HandleRequest(context.Background(), listReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	keys := resp.Data["keys"].([]string)
	if len(keys) != 1 {
		t.Fatalf("bad: length of alias IDs listed; expected: 1, actual: %d", len(keys))
	}

	// Do some due diligence on the key info
	aliasInfoRaw, ok := resp.Data["key_info"]
	if !ok {
		t.Fatal("expected key_info map in response")
	}
	aliasInfo := aliasInfoRaw.(map[string]interface{})
	for _, key := range keys {
		infoRaw, ok := aliasInfo[key]
		if !ok {
			t.Fatal("expected key info")
		}
		info := infoRaw.(map[string]interface{})
		t.Logf("alias info: %#v", info)
		switch {
		case info["name"].(string) != "groupalias":
			t.Fatalf("bad name: %v", info["name"].(string))
		case info["mount_accessor"].(string) != githubAccessor:
			t.Fatalf("bad mount_accessor: %v", info["mount_accessor"].(string))
		}
	}

	// Now do the same with entity info
	listReq = &logical.Request{
		Operation: logical.ListOperation,
		Path:      "group/id",
	}
	resp, err = is.HandleRequest(context.Background(), listReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	keys = resp.Data["keys"].([]string)
	if len(keys) != 1 {
		t.Fatalf("bad: length of entity IDs listed; expected: 1, actual: %d", len(keys))
	}

	groupInfoRaw, ok := resp.Data["key_info"]
	if !ok {
		t.Fatal("expected key_info map in response")
	}

	// This is basically verifying that the group has the alias in key_info
	// that we expect to be tied to it, plus tests a value further down in it
	// for fun
	groupInfo := groupInfoRaw.(map[string]interface{})
	for _, key := range keys {
		infoRaw, ok := groupInfo[key]
		if !ok {
			t.Fatal("expected key info")
		}
		info := infoRaw.(map[string]interface{})
		t.Logf("group info: %#v", info)
		alias := info["alias"].(map[string]interface{})
		switch {
		case alias["id"].(string) != aliasID:
			t.Fatalf("bad alias id: %v", alias["id"])
		case alias["mount_accessor"].(string) != githubAccessor:
			t.Fatalf("bad mount accessor: %v", alias["mount_accessor"])
		}
	}
}

func TestIdentityStore_GroupAliasDeletionOnGroupDeletion(t *testing.T) {
	var resp *logical.Response
	var err error

	i, accessor, _ := testIdentityStoreWithGithubAuth(t)

	resp, err = i.HandleRequest(context.Background(), &logical.Request{
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

	resp, err = i.HandleRequest(context.Background(), &logical.Request{
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

	resp, err = i.HandleRequest(context.Background(), &logical.Request{
		Path:      "group/id/" + groupID,
		Operation: logical.DeleteOperation,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
	}

	resp, err = i.HandleRequest(context.Background(), &logical.Request{
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
	i, accessor, _ := testIdentityStoreWithGithubAuth(t)

	groupReq := &logical.Request{
		Path:      "group",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"type": "external",
		},
	}
	resp, err = i.HandleRequest(context.Background(), groupReq)
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
	resp, err = i.HandleRequest(context.Background(), groupAliasReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v\n", resp, err)
	}
	groupAliasID := resp.Data["id"].(string)

	groupAliasReq.Path = "group-alias/id/" + groupAliasID
	groupAliasReq.Operation = logical.ReadOperation
	resp, err = i.HandleRequest(context.Background(), groupAliasReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v\n", resp, err)
	}

	if resp.Data["id"].(string) != groupAliasID {
		t.Fatalf("bad: group alias: %#v\n", resp.Data)
	}

	resp, err = i.HandleRequest(context.Background(), &logical.Request{
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
	resp, err = i.HandleRequest(context.Background(), groupAliasReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v\n", resp, err)
	}

	groupAliasReq.Operation = logical.ReadOperation
	resp, err = i.HandleRequest(context.Background(), groupAliasReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v\n", resp, err)
	}

	if resp != nil {
		t.Fatalf("failed to delete group alias")
	}
}

func TestIdentityStore_GroupAliases_MemDBIndexes(t *testing.T) {
	var err error
	i, accessor, _ := testIdentityStoreWithGithubAuth(t)

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

	err = i.MemDBUpsertAlias(group.Alias, true)
	if err != nil {
		t.Fatal(err)
	}

	err = i.MemDBUpsertGroup(group)
	if err != nil {
		t.Fatal(err)
	}

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

	i, accessor, _ := testIdentityStoreWithGithubAuth(t)

	groupReq := &logical.Request{
		Path:      "group",
		Operation: logical.UpdateOperation,
	}
	resp, err = i.HandleRequest(context.Background(), groupReq)
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
	resp, err = i.HandleRequest(context.Background(), aliasReq)
	if err != nil {
		t.Fatal(err)
	}
	if !resp.IsError() {
		t.Fatalf("expected an error")
	}
}
