package vault

import (
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/go-test/deep"
	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/logical"
)

func TestIdentityStore_GroupEntityMembershipUpgrade(t *testing.T) {
	c, keys, rootToken := TestCoreUnsealed(t)

	// Create a group
	resp, err := c.identityStore.HandleRequest(namespace.RootContext(nil), &logical.Request{
		Path:      "group",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"name": "testgroup",
		},
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err:%v\nresp: %#v", err, resp)
	}

	// Create a memdb transaction
	txn := c.identityStore.db.Txn(true)
	defer txn.Abort()

	// Fetch the above created group
	group, err := c.identityStore.MemDBGroupByNameInTxn(namespace.RootContext(nil), txn, "testgroup", true)
	if err != nil {
		t.Fatal(err)
	}

	// Manually add an invalid entity as the group's member
	group.MemberEntityIDs = []string{"invalidentityid"}

	// Persist the group
	err = c.identityStore.UpsertGroupInTxn(txn, group, true)
	if err != nil {
		t.Fatal(err)
	}

	txn.Commit()

	// Perform seal and unseal forcing an upgrade
	err = c.Seal(rootToken)
	if err != nil {
		t.Fatal(err)
	}
	for i, key := range keys {
		unseal, err := TestCoreUnseal(c, key)
		if err != nil {
			t.Fatal(err)
		}
		if i+1 == len(keys) && !unseal {
			t.Fatalf("failed to unseal")
		}
	}

	// Read the group and ensure that invalid entity id is cleaned up
	group, err = c.identityStore.MemDBGroupByName(namespace.RootContext(nil), "testgroup", false)
	if err != nil {
		t.Fatal(err)
	}

	if len(group.MemberEntityIDs) != 0 {
		t.Fatalf("bad: member entity IDs; expected: none, actual: %#v", group.MemberEntityIDs)
	}
}

func TestIdentityStore_MemberGroupIDDelete(t *testing.T) {
	ctx := namespace.RootContext(nil)
	i, _, _ := testIdentityStoreWithGithubAuth(ctx, t)

	// Create a child group
	resp, err := i.HandleRequest(ctx, &logical.Request{
		Path:      "group",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"name": "child",
		},
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
	}
	childGroupID := resp.Data["id"].(string)

	// Create a parent group with the above group ID as its child
	resp, err = i.HandleRequest(ctx, &logical.Request{
		Path:      "group",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"name":             "parent",
			"member_group_ids": []string{childGroupID},
		},
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
	}

	// Ensure that member group ID is properly updated
	resp, err = i.HandleRequest(ctx, &logical.Request{
		Path:      "group/name/parent",
		Operation: logical.ReadOperation,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
	}
	memberGroupIDs := resp.Data["member_group_ids"].([]string)
	if len(memberGroupIDs) != 1 && memberGroupIDs[0] != childGroupID {
		t.Fatalf("bad: member group ids; expected: %#v, actual: %#v", []string{childGroupID}, memberGroupIDs)
	}

	// Clear the member group IDs from the parent group
	resp, err = i.HandleRequest(ctx, &logical.Request{
		Path:      "group/name/parent",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"member_group_ids": []string{},
		},
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
	}

	// Ensure that member group ID is properly deleted
	resp, err = i.HandleRequest(ctx, &logical.Request{
		Path:      "group/name/parent",
		Operation: logical.ReadOperation,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
	}
	memberGroupIDs = resp.Data["member_group_ids"].([]string)
	if len(memberGroupIDs) != 0 {
		t.Fatalf("bad: length of member group ids; expected: %d, actual: %d", 0, len(memberGroupIDs))
	}
}

func TestIdentityStore_CaseInsensitiveGroupName(t *testing.T) {
	ctx := namespace.RootContext(nil)
	i, _, _ := testIdentityStoreWithGithubAuth(ctx, t)

	testGroupName := "testGroupName"

	// Create an group with case sensitive name
	resp, err := i.HandleRequest(ctx, &logical.Request{
		Path:      "group",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"name": testGroupName,
		},
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err:%v\nresp: %#v", err, resp)
	}
	groupID := resp.Data["id"].(string)

	// Lookup the group by ID and check that name returned is case sensitive
	resp, err = i.HandleRequest(ctx, &logical.Request{
		Path:      "group/id/" + groupID,
		Operation: logical.ReadOperation,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err:%v\nresp: %#v", err, resp)
	}
	groupName := resp.Data["name"].(string)
	if groupName != testGroupName {
		t.Fatalf("bad group name; expected: %q, actual: %q", testGroupName, groupName)
	}

	// Lookup the group by case sensitive name
	resp, err = i.HandleRequest(ctx, &logical.Request{
		Path:      "group/name/" + testGroupName,
		Operation: logical.ReadOperation,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v\nresp: %#v", err, resp)
	}
	groupName = resp.Data["name"].(string)
	if groupName != testGroupName {
		t.Fatalf("bad group name; expected: %q, actual: %q", testGroupName, groupName)
	}

	// Lookup the group by case insensitive name
	resp, err = i.HandleRequest(ctx, &logical.Request{
		Path:      "group/name/" + strings.ToLower(testGroupName),
		Operation: logical.ReadOperation,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v\nresp: %#v", err, resp)
	}
	groupName = resp.Data["name"].(string)
	if groupName != testGroupName {
		t.Fatalf("bad group name; expected: %q, actual: %q", testGroupName, groupName)
	}

	// Ensure that there is only one group
	resp, err = i.HandleRequest(ctx, &logical.Request{
		Path:      "group/name",
		Operation: logical.ListOperation,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v\nresp: %#v", err, resp)
	}
	if len(resp.Data["keys"].([]string)) != 1 {
		t.Fatalf("bad length of groups; expected: 1, actual: %d", len(resp.Data["keys"].([]string)))
	}
}

func TestIdentityStore_GroupByName(t *testing.T) {
	ctx := namespace.RootContext(nil)
	i, _, _ := testIdentityStoreWithGithubAuth(ctx, t)

	// Create an entity using the "name" endpoint
	resp, err := i.HandleRequest(ctx, &logical.Request{
		Path:      "group/name/testgroupname",
		Operation: logical.UpdateOperation,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
	}
	if resp == nil {
		t.Fatalf("expected a non-nil response")
	}

	// Test the read by name endpoint
	resp, err = i.HandleRequest(ctx, &logical.Request{
		Path:      "group/name/testgroupname",
		Operation: logical.ReadOperation,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
	}
	if resp == nil || resp.Data["name"].(string) != "testgroupname" {
		t.Fatalf("bad entity response: %#v", resp)
	}

	// Update group metadata using the name endpoint
	groupMetadata := map[string]string{
		"foo": "bar",
	}
	resp, err = i.HandleRequest(ctx, &logical.Request{
		Path:      "group/name/testgroupname",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"metadata": groupMetadata,
		},
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
	}

	// Check the updated result
	resp, err = i.HandleRequest(ctx, &logical.Request{
		Path:      "group/name/testgroupname",
		Operation: logical.ReadOperation,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
	}
	if resp == nil || !reflect.DeepEqual(resp.Data["metadata"].(map[string]string), groupMetadata) {
		t.Fatalf("bad group response: %#v", resp)
	}

	// Delete the group using the name endpoint
	resp, err = i.HandleRequest(ctx, &logical.Request{
		Path:      "group/name/testgroupname",
		Operation: logical.DeleteOperation,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
	}

	// Check if deletion was successful
	resp, err = i.HandleRequest(ctx, &logical.Request{
		Path:      "group/name/testgroupname",
		Operation: logical.ReadOperation,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
	}
	if resp != nil {
		t.Fatalf("expected a nil response")
	}

	// Create 2 entities
	resp, err = i.HandleRequest(ctx, &logical.Request{
		Path:      "group/name/testgroupname",
		Operation: logical.UpdateOperation,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
	}
	if resp == nil {
		t.Fatalf("expected a non-nil response")
	}
	resp, err = i.HandleRequest(ctx, &logical.Request{
		Path:      "group/name/testgroupname2",
		Operation: logical.UpdateOperation,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
	}
	if resp == nil {
		t.Fatalf("expected a non-nil response")
	}

	// List the entities by name
	resp, err = i.HandleRequest(ctx, &logical.Request{
		Path:      "group/name",
		Operation: logical.ListOperation,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
	}

	expected := []string{"testgroupname2", "testgroupname"}
	sort.Strings(expected)
	actual := resp.Data["keys"].([]string)
	sort.Strings(actual)
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("bad: group list response; expected: %#v\nactual: %#v", expected, actual)
	}
}

func TestIdentityStore_Groups_TypeMembershipAdditions(t *testing.T) {
	var err error
	var resp *logical.Response

	ctx := namespace.RootContext(nil)
	i, _, _ := testIdentityStoreWithGithubAuth(ctx, t)
	groupReq := &logical.Request{
		Path:      "group",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"type":              "external",
			"member_entity_ids": "sampleentityid",
		},
	}

	resp, err = i.HandleRequest(ctx, groupReq)
	if err != nil {
		t.Fatal(err)
	}
	if !resp.IsError() {
		t.Fatalf("expected an error")
	}

	groupReq.Data = map[string]interface{}{
		"type":             "external",
		"member_group_ids": "samplegroupid",
	}

	resp, err = i.HandleRequest(ctx, groupReq)
	if err != nil {
		t.Fatal(err)
	}
	if !resp.IsError() {
		t.Fatalf("expected an error")
	}
}

func TestIdentityStore_Groups_TypeImmutability(t *testing.T) {
	var err error
	var resp *logical.Response

	ctx := namespace.RootContext(nil)
	i, _, _ := testIdentityStoreWithGithubAuth(ctx, t)
	groupReq := &logical.Request{
		Path:      "group",
		Operation: logical.UpdateOperation,
	}

	resp, err = i.HandleRequest(ctx, groupReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v, err: %v", resp, err)
	}
	internalGroupID := resp.Data["id"].(string)

	groupReq.Data = map[string]interface{}{
		"type": "external",
	}
	resp, err = i.HandleRequest(ctx, groupReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v, err: %v", resp, err)
	}
	externalGroupID := resp.Data["id"].(string)

	// Try to mark internal group as external
	groupReq.Data = map[string]interface{}{
		"type": "external",
	}
	groupReq.Path = "group/id/" + internalGroupID
	resp, err = i.HandleRequest(ctx, groupReq)
	if err != nil {
		t.Fatal(err)
	}
	if !resp.IsError() {
		t.Fatalf("expected an error")
	}

	// Try to mark internal group as external
	groupReq.Data = map[string]interface{}{
		"type": "internal",
	}
	groupReq.Path = "group/id/" + externalGroupID
	resp, err = i.HandleRequest(ctx, groupReq)
	if err != nil {
		t.Fatal(err)
	}
	if !resp.IsError() {
		t.Fatalf("expected an error")
	}
}

func TestIdentityStore_MemDBGroupIndexes(t *testing.T) {
	var err error
	ctx := namespace.RootContext(nil)
	i, _, _ := testIdentityStoreWithGithubAuth(ctx, t)

	// Create a dummy group
	group := &identity.Group{
		ID:   "testgroupid",
		Name: "testgroupname",
		Metadata: map[string]string{
			"testmetadatakey1": "testmetadatavalue1",
			"testmetadatakey2": "testmetadatavalue2",
		},
		ParentGroupIDs:  []string{"testparentgroupid1", "testparentgroupid2"},
		MemberEntityIDs: []string{"testentityid1", "testentityid2"},
		Policies:        []string{"testpolicy1", "testpolicy2"},
		BucketKeyHash:   i.groupPacker.BucketKeyHashByItemID("testgroupid"),
	}

	// Insert it into memdb
	txn := i.db.Txn(true)
	defer txn.Abort()
	err = i.MemDBUpsertGroupInTxn(txn, group)
	if err != nil {
		t.Fatal(err)
	}
	txn.Commit()

	// Insert another dummy group
	group = &identity.Group{
		ID:   "testgroupid2",
		Name: "testgroupname2",
		Metadata: map[string]string{
			"testmetadatakey2": "testmetadatavalue2",
			"testmetadatakey3": "testmetadatavalue3",
		},
		ParentGroupIDs:  []string{"testparentgroupid2", "testparentgroupid3"},
		MemberEntityIDs: []string{"testentityid2", "testentityid3"},
		Policies:        []string{"testpolicy2", "testpolicy3"},
		BucketKeyHash:   i.groupPacker.BucketKeyHashByItemID("testgroupid2"),
	}

	// Insert it into memdb

	txn = i.db.Txn(true)
	defer txn.Abort()
	err = i.MemDBUpsertGroupInTxn(txn, group)
	if err != nil {
		t.Fatal(err)
	}
	txn.Commit()

	var fetchedGroup *identity.Group

	// Fetch group given the name
	fetchedGroup, err = i.MemDBGroupByName(namespace.RootContext(nil), "testgroupname", false)
	if err != nil {
		t.Fatal(err)
	}
	if fetchedGroup == nil || fetchedGroup.Name != "testgroupname" {
		t.Fatalf("failed to fetch an indexed group")
	}

	// Fetch group given the ID
	fetchedGroup, err = i.MemDBGroupByID("testgroupid", false)
	if err != nil {
		t.Fatal(err)
	}
	if fetchedGroup == nil || fetchedGroup.Name != "testgroupname" {
		t.Fatalf("failed to fetch an indexed group")
	}

	var fetchedGroups []*identity.Group
	// Fetch the subgroups of a given group ID
	fetchedGroups, err = i.MemDBGroupsByParentGroupID("testparentgroupid1", false)
	if err != nil {
		t.Fatal(err)
	}
	if len(fetchedGroups) != 1 || fetchedGroups[0].Name != "testgroupname" {
		t.Fatalf("failed to fetch an indexed group")
	}

	fetchedGroups, err = i.MemDBGroupsByParentGroupID("testparentgroupid2", false)
	if err != nil {
		t.Fatal(err)
	}
	if len(fetchedGroups) != 2 {
		t.Fatalf("failed to fetch a indexed groups")
	}

	// Fetch groups based on member entity ID
	fetchedGroups, err = i.MemDBGroupsByMemberEntityID("testentityid1", false, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(fetchedGroups) != 1 || fetchedGroups[0].Name != "testgroupname" {
		t.Fatalf("failed to fetch an indexed group")
	}

	fetchedGroups, err = i.MemDBGroupsByMemberEntityID("testentityid2", false, false)
	if err != nil {
		t.Fatal(err)
	}

	if len(fetchedGroups) != 2 {
		t.Fatalf("failed to fetch groups by entity ID")
	}
}

func TestIdentityStore_GroupsCreateUpdate(t *testing.T) {
	var resp *logical.Response
	var err error

	ctx := namespace.RootContext(nil)
	is, _, _ := testIdentityStoreWithGithubAuth(ctx, t)

	// Create an entity and get its ID
	entityRegisterReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "entity",
	}
	resp, err = is.HandleRequest(ctx, entityRegisterReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v, err: %v", resp, err)
	}
	entityID1 := resp.Data["id"].(string)

	// Create another entity and get its ID
	resp, err = is.HandleRequest(ctx, entityRegisterReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v, err: %v", resp, err)
	}
	entityID2 := resp.Data["id"].(string)

	// Create a group with the above created 2 entities as its members
	groupData := map[string]interface{}{
		"policies":          "testpolicy1,testpolicy2",
		"metadata":          []string{"testkey1=testvalue1", "testkey2=testvalue2"},
		"member_entity_ids": []string{entityID1, entityID2},
	}

	// Create a group and get its ID
	groupReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "group",
		Data:      groupData,
	}
	resp, err = is.HandleRequest(ctx, groupReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v, err: %v", resp, err)
	}
	memberGroupID1 := resp.Data["id"].(string)

	// Create another group and get its ID
	resp, err = is.HandleRequest(ctx, groupReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v, err: %v", resp, err)
	}
	memberGroupID2 := resp.Data["id"].(string)

	// Create a group with the above 2 groups as its members
	groupData["member_group_ids"] = []string{memberGroupID1, memberGroupID2}
	resp, err = is.HandleRequest(ctx, groupReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v, err: %v", resp, err)
	}
	groupID := resp.Data["id"].(string)

	// Read the group using its iD and check if all the fields are properly
	// set
	groupReq = &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "group/id/" + groupID,
	}
	resp, err = is.HandleRequest(ctx, groupReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v, err: %v", resp, err)
	}

	expectedData := map[string]interface{}{
		"policies": []string{"testpolicy1", "testpolicy2"},
		"metadata": map[string]string{
			"testkey1": "testvalue1",
			"testkey2": "testvalue2",
		},
		"parent_group_ids": []string(nil),
	}
	expectedData["id"] = resp.Data["id"]
	expectedData["type"] = resp.Data["type"]
	expectedData["name"] = resp.Data["name"]
	expectedData["member_group_ids"] = resp.Data["member_group_ids"]
	expectedData["member_entity_ids"] = resp.Data["member_entity_ids"]
	expectedData["creation_time"] = resp.Data["creation_time"]
	expectedData["last_update_time"] = resp.Data["last_update_time"]
	expectedData["modify_index"] = resp.Data["modify_index"]
	expectedData["alias"] = resp.Data["alias"]

	if diff := deep.Equal(expectedData, resp.Data); diff != nil {
		t.Fatal(diff)
	}

	// Update the policies and metadata in the group
	groupReq.Operation = logical.UpdateOperation
	groupReq.Data = groupData

	// Update by setting ID in the param
	groupData["id"] = groupID
	groupData["policies"] = "updatedpolicy1,updatedpolicy2"
	groupData["metadata"] = []string{"updatedkey=updatedvalue"}
	resp, err = is.HandleRequest(ctx, groupReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v, err: %v", resp, err)
	}

	// Check if updates are reflected
	groupReq.Operation = logical.ReadOperation
	resp, err = is.HandleRequest(ctx, groupReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v, err: %v", resp, err)
	}

	expectedData["policies"] = []string{"updatedpolicy1", "updatedpolicy2"}
	expectedData["metadata"] = map[string]string{
		"updatedkey": "updatedvalue",
	}
	expectedData["last_update_time"] = resp.Data["last_update_time"]
	expectedData["modify_index"] = resp.Data["modify_index"]
	if !reflect.DeepEqual(expectedData, resp.Data) {
		t.Fatalf("bad: group data; expected: %#v\n actual: %#v\n", expectedData, resp.Data)
	}
}

func TestIdentityStore_GroupsCRUD_ByID(t *testing.T) {
	var resp *logical.Response
	var err error
	ctx := namespace.RootContext(nil)
	is, _, _ := testIdentityStoreWithGithubAuth(ctx, t)

	// Create an entity and get its ID
	entityRegisterReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "entity",
	}
	resp, err = is.HandleRequest(ctx, entityRegisterReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v, err: %v", resp, err)
	}
	entityID1 := resp.Data["id"].(string)

	// Create another entity and get its ID
	resp, err = is.HandleRequest(ctx, entityRegisterReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v, err: %v", resp, err)
	}
	entityID2 := resp.Data["id"].(string)

	// Create a group with the above created 2 entities as its members
	groupData := map[string]interface{}{
		"policies":          "testpolicy1,testpolicy2",
		"metadata":          []string{"testkey1=testvalue1", "testkey2=testvalue2"},
		"member_entity_ids": []string{entityID1, entityID2},
	}

	// Create a group and get its ID
	groupRegisterReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "group",
		Data:      groupData,
	}
	resp, err = is.HandleRequest(ctx, groupRegisterReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v, err: %v", resp, err)
	}
	memberGroupID1 := resp.Data["id"].(string)

	// Create another group and get its ID
	resp, err = is.HandleRequest(ctx, groupRegisterReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v, err: %v", resp, err)
	}
	memberGroupID2 := resp.Data["id"].(string)

	// Create a group with the above 2 groups as its members
	groupData["member_group_ids"] = []string{memberGroupID1, memberGroupID2}
	resp, err = is.HandleRequest(ctx, groupRegisterReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v, err: %v", resp, err)
	}
	groupID := resp.Data["id"].(string)

	// Read the group using its name and check if all the fields are properly
	// set
	groupReq := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "group/id/" + groupID,
	}
	resp, err = is.HandleRequest(ctx, groupReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v, err: %v", resp, err)
	}

	expectedData := map[string]interface{}{
		"policies": []string{"testpolicy1", "testpolicy2"},
		"metadata": map[string]string{
			"testkey1": "testvalue1",
			"testkey2": "testvalue2",
		},
		"parent_group_ids": []string(nil),
	}
	expectedData["id"] = resp.Data["id"]
	expectedData["type"] = resp.Data["type"]
	expectedData["name"] = resp.Data["name"]
	expectedData["member_group_ids"] = resp.Data["member_group_ids"]
	expectedData["member_entity_ids"] = resp.Data["member_entity_ids"]
	expectedData["creation_time"] = resp.Data["creation_time"]
	expectedData["last_update_time"] = resp.Data["last_update_time"]
	expectedData["modify_index"] = resp.Data["modify_index"]
	expectedData["alias"] = resp.Data["alias"]

	if diff := deep.Equal(expectedData, resp.Data); diff != nil {
		t.Fatal(diff)
	}

	// Update the policies and metadata in the group
	groupReq.Operation = logical.UpdateOperation
	groupReq.Data = groupData
	groupData["policies"] = "updatedpolicy1,updatedpolicy2"
	groupData["metadata"] = []string{"updatedkey=updatedvalue"}
	resp, err = is.HandleRequest(ctx, groupReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v, err: %v", resp, err)
	}

	// Check if updates are reflected
	groupReq.Operation = logical.ReadOperation
	resp, err = is.HandleRequest(ctx, groupReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v, err: %v", resp, err)
	}

	expectedData["policies"] = []string{"updatedpolicy1", "updatedpolicy2"}
	expectedData["metadata"] = map[string]string{
		"updatedkey": "updatedvalue",
	}
	expectedData["last_update_time"] = resp.Data["last_update_time"]
	expectedData["modify_index"] = resp.Data["modify_index"]
	if !reflect.DeepEqual(expectedData, resp.Data) {
		t.Fatalf("bad: group data; expected: %#v\n actual: %#v\n", expectedData, resp.Data)
	}

	// Check if delete is working properly
	groupReq.Operation = logical.DeleteOperation
	resp, err = is.HandleRequest(ctx, groupReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v, err: %v", resp, err)
	}

	groupReq.Operation = logical.ReadOperation
	resp, err = is.HandleRequest(ctx, groupReq)
	if err != nil {
		t.Fatal(err)
	}
	if resp != nil {
		t.Fatalf("expected a nil response")
	}
}

func TestIdentityStore_GroupMultiCase(t *testing.T) {
	var resp *logical.Response
	var err error
	ctx := namespace.RootContext(nil)
	is, _, _ := testIdentityStoreWithGithubAuth(ctx, t)
	groupRegisterReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "group",
	}

	// Create 'build' group
	buildGroupData := map[string]interface{}{
		"name":     "build",
		"policies": "buildpolicy",
	}
	groupRegisterReq.Data = buildGroupData
	resp, err = is.HandleRequest(ctx, groupRegisterReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v, err: %v", resp, err)
	}
	buildGroupID := resp.Data["id"].(string)

	// Create 'deploy' group
	deployGroupData := map[string]interface{}{
		"name":     "deploy",
		"policies": "deploypolicy",
	}
	groupRegisterReq.Data = deployGroupData
	resp, err = is.HandleRequest(ctx, groupRegisterReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v, err: %v", resp, err)
	}
	deployGroupID := resp.Data["id"].(string)

	// Create an entity ID
	entityRegisterReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "entity",
	}
	resp, err = is.HandleRequest(ctx, entityRegisterReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v, err: %v", resp, err)
	}
	entityID1 := resp.Data["id"].(string)

	// Add the entity as a member of 'build' group
	entityIDReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "group/id/" + buildGroupID,
		Data: map[string]interface{}{
			"member_entity_ids": []string{entityID1},
		},
	}
	resp, err = is.HandleRequest(ctx, entityIDReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v, err: %v", resp, err)
	}

	// Add the entity as a member of the 'deploy` group
	entityIDReq.Path = "group/id/" + deployGroupID
	resp, err = is.HandleRequest(ctx, entityIDReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v, err: %v", resp, err)
	}

	policiesResult, err := is.groupPoliciesByEntityID(entityID1)
	if err != nil {
		t.Fatal(err)
	}

	policies := []string{}
	for _, nsPolicies := range policiesResult {
		policies = append(policies, nsPolicies...)
	}
	sort.Strings(policies)
	expected := []string{"deploypolicy", "buildpolicy"}
	sort.Strings(expected)
	if !reflect.DeepEqual(expected, policies) {
		t.Fatalf("bad: policies; expected: %#v\nactual:%#v", expected, policies)
	}
}

/*
Test groups hierarchy:
                ------- eng(entityID3) -------
                |                            |
         ----- vault -----        -- ops(entityID2) --
         |               |        |                  |
   kube(entityID1)    identity    build            deploy
*/
func TestIdentityStore_GroupHierarchyCases(t *testing.T) {
	var resp *logical.Response
	var err error
	ctx := namespace.RootContext(nil)
	is, _, _ := testIdentityStoreWithGithubAuth(ctx, t)
	groupRegisterReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "group",
	}

	// Create 'kube' group
	kubeGroupData := map[string]interface{}{
		"name":     "kube",
		"policies": "kubepolicy",
	}
	groupRegisterReq.Data = kubeGroupData
	resp, err = is.HandleRequest(ctx, groupRegisterReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v, err: %v", resp, err)
	}
	kubeGroupID := resp.Data["id"].(string)

	// Create 'identity' group
	identityGroupData := map[string]interface{}{
		"name":     "identity",
		"policies": "identitypolicy",
	}
	groupRegisterReq.Data = identityGroupData
	resp, err = is.HandleRequest(ctx, groupRegisterReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v, err: %v", resp, err)
	}
	identityGroupID := resp.Data["id"].(string)

	// Create 'build' group
	buildGroupData := map[string]interface{}{
		"name":     "build",
		"policies": "buildpolicy",
	}
	groupRegisterReq.Data = buildGroupData
	resp, err = is.HandleRequest(ctx, groupRegisterReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v, err: %v", resp, err)
	}
	buildGroupID := resp.Data["id"].(string)

	// Create 'deploy' group
	deployGroupData := map[string]interface{}{
		"name":     "deploy",
		"policies": "deploypolicy",
	}
	groupRegisterReq.Data = deployGroupData
	resp, err = is.HandleRequest(ctx, groupRegisterReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v, err: %v", resp, err)
	}
	deployGroupID := resp.Data["id"].(string)

	// Create 'vault' with 'kube' and 'identity' as member groups
	vaultMemberGroupIDs := []string{kubeGroupID, identityGroupID}
	vaultGroupData := map[string]interface{}{
		"name":             "vault",
		"policies":         "vaultpolicy",
		"member_group_ids": vaultMemberGroupIDs,
	}
	groupRegisterReq.Data = vaultGroupData
	resp, err = is.HandleRequest(ctx, groupRegisterReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v, err: %v", resp, err)
	}
	vaultGroupID := resp.Data["id"].(string)

	// Create 'ops' group with 'build' and 'deploy' as member groups
	opsMemberGroupIDs := []string{buildGroupID, deployGroupID}
	opsGroupData := map[string]interface{}{
		"name":             "ops",
		"policies":         "opspolicy",
		"member_group_ids": opsMemberGroupIDs,
	}
	groupRegisterReq.Data = opsGroupData
	resp, err = is.HandleRequest(ctx, groupRegisterReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v, err: %v", resp, err)
	}
	opsGroupID := resp.Data["id"].(string)

	// Create 'eng' group with 'vault' and 'ops' as member groups
	engMemberGroupIDs := []string{vaultGroupID, opsGroupID}
	engGroupData := map[string]interface{}{
		"name":             "eng",
		"policies":         "engpolicy",
		"member_group_ids": engMemberGroupIDs,
	}

	groupRegisterReq.Data = engGroupData
	resp, err = is.HandleRequest(ctx, groupRegisterReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v, err: %v", resp, err)
	}
	engGroupID := resp.Data["id"].(string)

	/*
		fmt.Printf("engGroupID: %#v\n", engGroupID)
		fmt.Printf("vaultGroupID: %#v\n", vaultGroupID)
		fmt.Printf("opsGroupID: %#v\n", opsGroupID)
		fmt.Printf("kubeGroupID: %#v\n", kubeGroupID)
		fmt.Printf("identityGroupID: %#v\n", identityGroupID)
		fmt.Printf("buildGroupID: %#v\n", buildGroupID)
		fmt.Printf("deployGroupID: %#v\n", deployGroupID)
	*/

	var memberGroupIDs []string
	// Fetch 'eng' group
	engGroup, err := is.MemDBGroupByID(engGroupID, false)
	if err != nil {
		t.Fatal(err)
	}
	memberGroupIDs, err = is.memberGroupIDsByID(engGroup.ID)
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(memberGroupIDs)
	sort.Strings(engMemberGroupIDs)
	if !reflect.DeepEqual(engMemberGroupIDs, memberGroupIDs) {
		t.Fatalf("bad: group membership IDs; expected: %#v\n actual: %#v\n", engMemberGroupIDs, memberGroupIDs)
	}

	vaultGroup, err := is.MemDBGroupByID(vaultGroupID, false)
	if err != nil {
		t.Fatal(err)
	}
	memberGroupIDs, err = is.memberGroupIDsByID(vaultGroup.ID)
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(memberGroupIDs)
	sort.Strings(vaultMemberGroupIDs)
	if !reflect.DeepEqual(vaultMemberGroupIDs, memberGroupIDs) {
		t.Fatalf("bad: group membership IDs; expected: %#v\n actual: %#v\n", vaultMemberGroupIDs, memberGroupIDs)
	}

	opsGroup, err := is.MemDBGroupByID(opsGroupID, false)
	if err != nil {
		t.Fatal(err)
	}
	memberGroupIDs, err = is.memberGroupIDsByID(opsGroup.ID)
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(memberGroupIDs)
	sort.Strings(opsMemberGroupIDs)
	if !reflect.DeepEqual(opsMemberGroupIDs, memberGroupIDs) {
		t.Fatalf("bad: group membership IDs; expected: %#v\n actual: %#v\n", opsMemberGroupIDs, memberGroupIDs)
	}

	groupUpdateReq := &logical.Request{
		Operation: logical.UpdateOperation,
	}

	// Adding 'engGroupID' under 'kubeGroupID' should fail
	groupUpdateReq.Path = "group/name/kube"
	groupUpdateReq.Data = kubeGroupData
	kubeGroupData["member_group_ids"] = []string{engGroupID}
	resp, err = is.HandleRequest(ctx, groupUpdateReq)
	if err == nil {
		t.Fatalf("expected an error response")
	}

	// Create an entity ID
	entityRegisterReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "entity",
	}
	resp, err = is.HandleRequest(ctx, entityRegisterReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v, err: %v", resp, err)
	}
	entityID1 := resp.Data["id"].(string)

	// Add the entity as a member of 'kube' group
	entityIDReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "group/id/" + kubeGroupID,
		Data: map[string]interface{}{
			"member_entity_ids": []string{entityID1},
		},
	}
	resp, err = is.HandleRequest(ctx, entityIDReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v, err: %v", resp, err)
	}

	// Create a second entity ID
	resp, err = is.HandleRequest(ctx, entityRegisterReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v, err: %v", resp, err)
	}
	entityID2 := resp.Data["id"].(string)

	// Add the entity as a member of 'ops' group
	entityIDReq.Path = "group/id/" + opsGroupID
	entityIDReq.Data = map[string]interface{}{
		"member_entity_ids": []string{entityID2},
	}
	resp, err = is.HandleRequest(ctx, entityIDReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v, err: %v", resp, err)
	}

	// Create a third entity ID
	resp, err = is.HandleRequest(ctx, entityRegisterReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v, err: %v", resp, err)
	}
	entityID3 := resp.Data["id"].(string)

	// Add the entity as a member of 'eng' group
	entityIDReq.Path = "group/id/" + engGroupID
	entityIDReq.Data = map[string]interface{}{
		"member_entity_ids": []string{entityID3},
		"member_group_ids":  engMemberGroupIDs,
	}
	resp, err = is.HandleRequest(ctx, entityIDReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v, err: %v", resp, err)
	}

	policiesResult, err := is.groupPoliciesByEntityID(entityID1)
	if err != nil {
		t.Fatal(err)
	}
	var policies []string
	for _, nsPolicies := range policiesResult {
		policies = append(policies, nsPolicies...)
	}
	sort.Strings(policies)
	expected := []string{"kubepolicy", "vaultpolicy", "engpolicy"}
	sort.Strings(expected)
	if !reflect.DeepEqual(expected, policies) {
		t.Fatalf("bad: policies; expected: %#v\nactual:%#v", expected, policies)
	}

	policiesResult, err = is.groupPoliciesByEntityID(entityID2)
	if err != nil {
		t.Fatal(err)
	}
	policies = nil
	for _, nsPolicies := range policiesResult {
		policies = append(policies, nsPolicies...)
	}
	sort.Strings(policies)
	expected = []string{"opspolicy", "engpolicy"}
	sort.Strings(expected)
	if !reflect.DeepEqual(expected, policies) {
		t.Fatalf("bad: policies; expected: %#v\nactual:%#v", expected, policies)
	}

	policiesResult, err = is.groupPoliciesByEntityID(entityID3)
	if err != nil {
		t.Fatal(err)
	}
	policies = nil
	for _, nsPolicies := range policiesResult {
		policies = append(policies, nsPolicies...)
	}

	if len(policies) != 1 && policies[0] != "engpolicy" {
		t.Fatalf("bad: policies; expected: 'engpolicy'\nactual:%#v", policies)
	}

	groups, inheritedGroups, err := is.groupsByEntityID(entityID1)
	if err != nil {
		t.Fatal(err)
	}
	if len(groups) != 1 {
		t.Fatalf("bad: length of groups; expected: 1, actual: %d", len(groups))
	}
	if len(inheritedGroups) != 2 {
		t.Fatalf("bad: length of inheritedGroups; expected: 2, actual: %d", len(inheritedGroups))
	}

	groups, inheritedGroups, err = is.groupsByEntityID(entityID2)
	if err != nil {
		t.Fatal(err)
	}
	if len(groups) != 1 {
		t.Fatalf("bad: length of groups; expected: 1, actual: %d", len(groups))
	}
	if len(inheritedGroups) != 1 {
		t.Fatalf("bad: length of inheritedGroups; expected: 1, actual: %d", len(inheritedGroups))
	}

	groups, inheritedGroups, err = is.groupsByEntityID(entityID3)
	if err != nil {
		t.Fatal(err)
	}
	if len(groups) != 1 {
		t.Fatalf("bad: length of groups; expected: 1, actual: %d", len(groups))
	}
	if len(inheritedGroups) != 0 {
		t.Fatalf("bad: length of inheritedGroups; expected: 0, actual: %d", len(inheritedGroups))
	}
}
