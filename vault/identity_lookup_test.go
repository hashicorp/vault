package vault

import (
	"testing"

	"github.com/hashicorp/vault/logical"
)

func TestIdentityStore_Lookup_EntityAlias(t *testing.T) {
	var err error
	var resp *logical.Response

	i, accessor, _ := testIdentityStoreWithGithubAuth(t)

	entityReq := &logical.Request{
		Path:      "entity",
		Operation: logical.UpdateOperation,
	}

	resp, err = i.HandleRequest(entityReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\n err: %#v\n", resp, err)
	}
	entityID := resp.Data["id"].(string)

	entityAliasReq := &logical.Request{
		Path:      "entity-alias",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"parent_id":      entityID,
			"name":           "testentityaliasname",
			"mount_type":     "ldap",
			"mount_accessor": accessor,
		},
	}

	resp, err = i.HandleRequest(entityAliasReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\n err: %#v\n", resp, err)
	}
	entityAliasID := resp.Data["id"].(string)

	lookupReq := &logical.Request{
		Path:      "lookup/entity-alias",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"type": "id",
			"id":   entityAliasID,
		},
	}
	resp, err = i.HandleRequest(lookupReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\n err: %#v\n", resp, err)
	}
	if resp.Data["id"].(string) != entityAliasID {
		t.Fatalf("bad: group alias: %#v\n", resp.Data)
	}

	lookupReq.Data = map[string]interface{}{
		"type":           "factors",
		"mount_accessor": accessor,
		"name":           "testentityaliasname",
	}
	resp, err = i.HandleRequest(lookupReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\n err: %#v\n", resp, err)
	}
	if resp.Data["id"].(string) != entityAliasID {
		t.Fatalf("bad: entity alias: %#v\n", resp.Data)
	}

	entityReq = &logical.Request{
		Path:      "entity/id/" + entityID,
		Operation: logical.ReadOperation,
	}
	resp, err = i.HandleRequest(entityReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\n err: %#v\n", resp, err)
	}

	lookupReq.Data = map[string]interface{}{
		"type":      "parent_id",
		"parent_id": entityID,
	}
	resp, err = i.HandleRequest(lookupReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\n err: %#v\n", resp, err)
	}
	if resp.Data["id"].(string) != entityAliasID {
		t.Fatalf("bad: entity alias: %#v\n", resp.Data)
	}
}

func TestIdentityStore_Lookup_GroupAlias(t *testing.T) {
	var err error
	var resp *logical.Response

	i, accessor, _ := testIdentityStoreWithGithubAuth(t)

	groupReq := &logical.Request{
		Path:      "group",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"type": "external",
		},
	}

	resp, err = i.HandleRequest(groupReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\n err: %#v\n", resp, err)
	}
	groupID := resp.Data["id"].(string)

	groupAliasReq := &logical.Request{
		Path:      "group-alias",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"parent_id":      groupID,
			"name":           "testgroupaliasname",
			"mount_type":     "ldap",
			"mount_accessor": accessor,
		},
	}

	resp, err = i.HandleRequest(groupAliasReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\n err: %#v\n", resp, err)
	}
	groupAliasID := resp.Data["id"].(string)

	lookupReq := &logical.Request{
		Path:      "lookup/group-alias",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"type": "id",
			"id":   groupAliasID,
		},
	}
	resp, err = i.HandleRequest(lookupReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\n err: %#v\n", resp, err)
	}
	if resp.Data["id"].(string) != groupAliasID {
		t.Fatalf("bad: group alias: %#v\n", resp.Data)
	}

	lookupReq.Data = map[string]interface{}{
		"type":           "factors",
		"mount_accessor": accessor,
		"name":           "testgroupaliasname",
	}
	resp, err = i.HandleRequest(lookupReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\n err: %#v\n", resp, err)
	}
	if resp.Data["id"].(string) != groupAliasID {
		t.Fatalf("bad: group alias: %#v\n", resp.Data)
	}

	lookupReq.Data = map[string]interface{}{
		"type":      "parent_id",
		"parent_id": groupID,
	}
	resp, err = i.HandleRequest(lookupReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\n err: %#v\n", resp, err)
	}
	if resp.Data["id"].(string) != groupAliasID {
		t.Fatalf("bad: group alias: %#v\n", resp.Data)
	}
}

func TestIdentityStore_Lookup_Group(t *testing.T) {
	var err error
	var resp *logical.Response

	i, _, _ := testIdentityStoreWithGithubAuth(t)

	groupReq := &logical.Request{
		Path:      "group",
		Operation: logical.UpdateOperation,
	}
	resp, err = i.HandleRequest(groupReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\n err: %#v\n", resp, err)
	}
	groupID := resp.Data["id"].(string)
	groupName := resp.Data["name"].(string)

	lookupGroupReq := &logical.Request{
		Path:      "lookup/group",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"type": "id",
			"id":   groupID,
		},
	}

	resp, err = i.HandleRequest(lookupGroupReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\n err: %#v\n", resp, err)
	}
	if resp.Data["id"].(string) != groupID {
		t.Fatalf("failed to lookup group")
	}

	lookupGroupReq.Data = map[string]interface{}{
		"type": "name",
		"name": groupName,
	}

	resp, err = i.HandleRequest(lookupGroupReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\n err: %#v\n", resp, err)
	}
	if resp.Data["id"].(string) != groupID {
		t.Fatalf("failed to lookup group")
	}
}
