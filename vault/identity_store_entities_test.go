package vault

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"

	uuid "github.com/hashicorp/go-uuid"
	credGithub "github.com/hashicorp/vault/builtin/credential/github"
	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/logical"
)

func TestIdentityStore_EntityDeleteGroupMembershipUpdate(t *testing.T) {
	i, _, _ := testIdentityStoreWithGithubAuth(namespace.RootContext(nil), t)

	// Create an entity
	resp, err := i.HandleRequest(namespace.RootContext(nil), &logical.Request{
		Path:      "entity",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"name": "testentity",
		},
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err:%v\nresp: %#v", err, resp)
	}
	entityID := resp.Data["id"].(string)

	// Create a group
	resp, err = i.HandleRequest(namespace.RootContext(nil), &logical.Request{
		Path:      "group",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"name":              "testgroup",
			"member_entity_ids": []string{entityID},
		},
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err:%v\nresp: %#v", err, resp)
	}

	// Ensure that the group has entity ID as its member
	resp, err = i.HandleRequest(namespace.RootContext(nil), &logical.Request{
		Path:      "group/name/testgroup",
		Operation: logical.ReadOperation,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err:%v\nresp: %#v", err, resp)
	}
	expected := []string{entityID}
	actual := resp.Data["member_entity_ids"].([]string)
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("bad: member entity ids; expected: %#v\nactual: %#v", expected, actual)
	}

	// Delete the entity
	resp, err = i.HandleRequest(namespace.RootContext(nil), &logical.Request{
		Path:      "entity/name/testentity",
		Operation: logical.DeleteOperation,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err:%v\nresp: %#v", err, resp)
	}

	// Ensure that the group does not have entity ID as it's member anymore
	resp, err = i.HandleRequest(namespace.RootContext(nil), &logical.Request{
		Path:      "group/name/testgroup",
		Operation: logical.ReadOperation,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err:%v\nresp: %#v", err, resp)
	}
	expected = []string{}
	actual = resp.Data["member_entity_ids"].([]string)
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("bad: member entity ids; expected: %#v\nactual: %#v", expected, actual)
	}
}

func TestIdentityStore_CaseInsensitiveEntityName(t *testing.T) {
	ctx := namespace.RootContext(nil)
	i, _, _ := testIdentityStoreWithGithubAuth(ctx, t)

	testEntityName := "testEntityName"

	// Create an entity with case sensitive name
	resp, err := i.HandleRequest(ctx, &logical.Request{
		Path:      "entity",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"name": testEntityName,
		},
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err:%v\nresp: %#v", err, resp)
	}
	entityID := resp.Data["id"].(string)

	// Lookup the entity by ID and check that name returned is case sensitive
	resp, err = i.HandleRequest(ctx, &logical.Request{
		Path:      "entity/id/" + entityID,
		Operation: logical.ReadOperation,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err:%v\nresp: %#v", err, resp)
	}
	entityName := resp.Data["name"].(string)
	if entityName != testEntityName {
		t.Fatalf("bad entity name; expected: %q, actual: %q", testEntityName, entityName)
	}

	// Lookup the entity by case sensitive name
	resp, err = i.HandleRequest(ctx, &logical.Request{
		Path:      "entity/name/" + testEntityName,
		Operation: logical.ReadOperation,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v\nresp: %#v", err, resp)
	}
	entityName = resp.Data["name"].(string)
	if entityName != testEntityName {
		t.Fatalf("bad entity name; expected: %q, actual: %q", testEntityName, entityName)
	}

	// Lookup the entity by case insensitive name
	resp, err = i.HandleRequest(ctx, &logical.Request{
		Path:      "entity/name/" + strings.ToLower(testEntityName),
		Operation: logical.ReadOperation,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v\nresp: %#v", err, resp)
	}
	entityName = resp.Data["name"].(string)
	if entityName != testEntityName {
		t.Fatalf("bad entity name; expected: %q, actual: %q", testEntityName, entityName)
	}

	// Ensure that there is only one entity
	resp, err = i.HandleRequest(ctx, &logical.Request{
		Path:      "entity/name",
		Operation: logical.ListOperation,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v\nresp: %#v", err, resp)
	}
	if len(resp.Data["keys"].([]string)) != 1 {
		t.Fatalf("bad length of entities; expected: 1, actual: %d", len(resp.Data["keys"].([]string)))
	}
}

func TestIdentityStore_EntityByName(t *testing.T) {
	ctx := namespace.RootContext(nil)
	i, _, _ := testIdentityStoreWithGithubAuth(ctx, t)

	// Create an entity using the "name" endpoint
	resp, err := i.HandleRequest(ctx, &logical.Request{
		Path:      "entity/name/testentityname",
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
		Path:      "entity/name/testentityname",
		Operation: logical.ReadOperation,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
	}
	if resp == nil || resp.Data["name"].(string) != "testentityname" {
		t.Fatalf("bad entity response: %#v", resp)
	}

	// Update entity metadata using the name endpoint
	entityMetadata := map[string]string{
		"foo": "bar",
	}
	resp, err = i.HandleRequest(ctx, &logical.Request{
		Path:      "entity/name/testentityname",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"metadata": entityMetadata,
		},
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
	}

	// Check the updated result
	resp, err = i.HandleRequest(ctx, &logical.Request{
		Path:      "entity/name/testentityname",
		Operation: logical.ReadOperation,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
	}
	if resp == nil || !reflect.DeepEqual(resp.Data["metadata"].(map[string]string), entityMetadata) {
		t.Fatalf("bad entity response: %#v", resp)
	}

	// Delete the entity using the name endpoint
	resp, err = i.HandleRequest(ctx, &logical.Request{
		Path:      "entity/name/testentityname",
		Operation: logical.DeleteOperation,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
	}

	// Check if deletion was successful
	resp, err = i.HandleRequest(ctx, &logical.Request{
		Path:      "entity/name/testentityname",
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
		Path:      "entity/name/testentityname",
		Operation: logical.UpdateOperation,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
	}
	if resp == nil {
		t.Fatalf("expected a non-nil response")
	}
	resp, err = i.HandleRequest(ctx, &logical.Request{
		Path:      "entity/name/testentityname2",
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
		Path:      "entity/name",
		Operation: logical.ListOperation,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
	}

	expected := []string{"testentityname2", "testentityname"}
	sort.Strings(expected)
	actual := resp.Data["keys"].([]string)
	sort.Strings(actual)
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("bad: entity list response; expected: %#v\nactual: %#v", expected, actual)
	}
}

func TestIdentityStore_EntityReadGroupIDs(t *testing.T) {
	var err error
	var resp *logical.Response

	ctx := namespace.RootContext(nil)
	i, _, _ := testIdentityStoreWithGithubAuth(ctx, t)

	entityReq := &logical.Request{
		Path:      "entity",
		Operation: logical.UpdateOperation,
	}

	resp, err = i.HandleRequest(ctx, entityReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
	}

	entityID := resp.Data["id"].(string)

	groupReq := &logical.Request{
		Path:      "group",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"member_entity_ids": []string{
				entityID,
			},
		},
	}

	resp, err = i.HandleRequest(ctx, groupReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
	}

	groupID := resp.Data["id"].(string)

	// Create another group with the above created group as its subgroup

	groupReq.Data = map[string]interface{}{
		"member_group_ids": []string{groupID},
	}
	resp, err = i.HandleRequest(ctx, groupReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
	}

	inheritedGroupID := resp.Data["id"].(string)

	lookupReq := &logical.Request{
		Path:      "lookup/entity",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"type": "id",
			"id":   entityID,
		},
	}

	resp, err = i.HandleRequest(ctx, lookupReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
	}

	expected := []string{groupID, inheritedGroupID}
	actual := resp.Data["group_ids"].([]string)
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("bad: group_ids; expected: %#v\nactual: %#v\n", expected, actual)
	}

	expected = []string{groupID}
	actual = resp.Data["direct_group_ids"].([]string)
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("bad: direct_group_ids; expected: %#v\nactual: %#v\n", expected, actual)
	}

	expected = []string{inheritedGroupID}
	actual = resp.Data["inherited_group_ids"].([]string)
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("bad: inherited_group_ids; expected: %#v\nactual: %#v\n", expected, actual)
	}
}

func TestIdentityStore_EntityCreateUpdate(t *testing.T) {
	var err error
	var resp *logical.Response

	ctx := namespace.RootContext(nil)
	is, _, _ := testIdentityStoreWithGithubAuth(ctx, t)

	entityData := map[string]interface{}{
		"name":     "testentityname",
		"metadata": []string{"someusefulkey=someusefulvalue"},
		"policies": []string{"testpolicy1", "testpolicy2"},
	}

	entityReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "entity",
		Data:      entityData,
	}

	// Create the entity
	resp, err = is.HandleRequest(ctx, entityReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	entityID := resp.Data["id"].(string)

	updateData := map[string]interface{}{
		// Set the entity ID here
		"id":       entityID,
		"name":     "updatedentityname",
		"metadata": []string{"updatedkey=updatedvalue"},
		"policies": []string{"updatedpolicy1", "updatedpolicy2"},
	}
	entityReq.Data = updateData

	// Update the entity
	resp, err = is.HandleRequest(ctx, entityReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	entityReq.Path = "entity/id/" + entityID
	entityReq.Operation = logical.ReadOperation

	// Read the entity
	resp, err = is.HandleRequest(ctx, entityReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	if resp.Data["id"] != entityID ||
		resp.Data["name"] != updateData["name"] ||
		!reflect.DeepEqual(resp.Data["policies"], updateData["policies"]) {
		t.Fatalf("bad: entity response after update; resp: %#v\n updateData: %#v\n", resp.Data, updateData)
	}
}

func TestIdentityStore_CloneImmutability(t *testing.T) {
	alias := &identity.Alias{
		ID:                     "testaliasid",
		Name:                   "testaliasname",
		MergedFromCanonicalIDs: []string{"entityid1"},
	}

	entity := &identity.Entity{
		ID:   "testentityid",
		Name: "testentityname",
		Aliases: []*identity.Alias{
			alias,
		},
	}

	clonedEntity, err := entity.Clone()
	if err != nil {
		t.Fatal(err)
	}

	// Modify entity
	entity.Aliases[0].ID = "invalidid"

	if clonedEntity.Aliases[0].ID == "invalidid" {
		t.Fatalf("cloned entity is mutated")
	}

	clonedAlias, err := alias.Clone()
	if err != nil {
		t.Fatal(err)
	}

	alias.MergedFromCanonicalIDs[0] = "invalidid"

	if clonedAlias.MergedFromCanonicalIDs[0] == "invalidid" {
		t.Fatalf("cloned alias is mutated")
	}
}

func TestIdentityStore_MemDBImmutability(t *testing.T) {
	var err error
	ctx := namespace.RootContext(nil)
	is, githubAccessor, _ := testIdentityStoreWithGithubAuth(ctx, t)

	validateMountResp := is.core.router.validateMountByAccessor(githubAccessor)
	if validateMountResp == nil {
		t.Fatal("failed to validate github auth mount")
	}

	alias1 := &identity.Alias{
		CanonicalID:   "testentityid",
		ID:            "testaliasid",
		MountAccessor: githubAccessor,
		MountType:     validateMountResp.MountType,
		Name:          "testaliasname",
		Metadata: map[string]string{
			"testkey1": "testmetadatavalue1",
			"testkey2": "testmetadatavalue2",
		},
	}

	entity := &identity.Entity{
		ID:   "testentityid",
		Name: "testentityname",
		Metadata: map[string]string{
			"someusefulkey": "someusefulvalue",
		},
		Aliases: []*identity.Alias{
			alias1,
		},
	}

	entity.BucketKeyHash = is.entityPacker.BucketKeyHashByItemID(entity.ID)

	txn := is.db.Txn(true)
	defer txn.Abort()

	err = is.MemDBUpsertEntityInTxn(txn, entity)
	if err != nil {
		t.Fatal(err)
	}

	txn.Commit()

	entityFetched, err := is.MemDBEntityByID(entity.ID, true)
	if err != nil {
		t.Fatal(err)
	}

	// Modify the fetched entity outside of a transaction
	entityFetched.Aliases[0].ID = "invalidaliasid"

	entityFetched, err = is.MemDBEntityByID(entity.ID, false)
	if err != nil {
		t.Fatal(err)
	}

	if entityFetched.Aliases[0].ID == "invalidaliasid" {
		t.Fatal("memdb item is mutable outside of transaction")
	}
}

func TestIdentityStore_ListEntities(t *testing.T) {
	var err error
	var resp *logical.Response

	ctx := namespace.RootContext(nil)
	is, _, _ := testIdentityStoreWithGithubAuth(ctx, t)

	entityReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "entity",
	}

	expected := []string{}
	for i := 0; i < 10; i++ {
		resp, err = is.HandleRequest(ctx, entityReq)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("err:%v resp:%#v", err, resp)
		}
		expected = append(expected, resp.Data["id"].(string))
	}

	listReq := &logical.Request{
		Operation: logical.ListOperation,
		Path:      "entity/id",
	}

	resp, err = is.HandleRequest(ctx, listReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	actual := resp.Data["keys"].([]string)

	// Sort the operands for DeepEqual to work
	sort.Strings(actual)
	sort.Strings(expected)
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("bad: listed entity IDs; expected: %#v\n actual: %#v\n", expected, actual)
	}
}

func TestIdentityStore_LoadingEntities(t *testing.T) {
	var resp *logical.Response
	// Add github credential factory to core config
	err := AddTestCredentialBackend("github", credGithub.Factory)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	c := TestCore(t)
	unsealKeys, token := TestCoreInit(t, c)
	for _, key := range unsealKeys {
		if _, err := TestCoreUnseal(c, TestKeyCopy(key)); err != nil {
			t.Fatalf("unseal err: %s", err)
		}
	}

	if c.Sealed() {
		t.Fatal("should not be sealed")
	}

	meGH := &MountEntry{
		Table:       credentialTableType,
		Path:        "github/",
		Type:        "github",
		Description: "github auth",
		namespace:   namespace.RootNamespace,
	}

	// Mount UUID for github auth
	meGHUUID, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	meGH.UUID = meGHUUID

	// Mount accessor for github auth
	githubAccessor, err := c.generateMountAccessor("github")
	if err != nil {
		panic(fmt.Sprintf("could not generate github accessor: %v", err))
	}
	meGH.Accessor = githubAccessor

	// Storage view for github auth
	ghView := NewBarrierView(c.barrier, credentialBarrierPrefix+meGH.UUID+"/")

	// Sysview for github auth
	ghSysview := c.mountEntrySysView(meGH)

	// Create new github auth credential backend
	ghAuth, err := c.newCredentialBackend(context.Background(), meGH, ghSysview, ghView)
	if err != nil {
		t.Fatal(err)
	}

	// Mount github auth
	err = c.router.Mount(ghAuth, "auth/github", meGH, ghView)
	if err != nil {
		t.Fatal(err)
	}

	// Identity store will be mounted by now, just fetch it from router
	identitystore := c.router.MatchingBackend(namespace.RootContext(nil), "identity/")
	if identitystore == nil {
		t.Fatalf("failed to fetch identity store from router")
	}

	is := identitystore.(*IdentityStore)

	registerData := map[string]interface{}{
		"name":     "testentityname",
		"metadata": []string{"someusefulkey=someusefulvalue"},
		"policies": []string{"testpolicy1", "testpolicy2"},
	}

	registerReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "entity",
		Data:      registerData,
	}

	ctx := namespace.RootContext(nil)

	// Register the entity
	resp, err = is.HandleRequest(ctx, registerReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	entityID := resp.Data["id"].(string)

	readReq := &logical.Request{
		Path:      "entity/id/" + entityID,
		Operation: logical.ReadOperation,
	}

	// Ensure that entity is created
	resp, err = is.HandleRequest(ctx, readReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	if resp.Data["id"] != entityID {
		t.Fatalf("failed to read the created entity")
	}

	// Perform a seal/unseal cycle
	err = c.Seal(token)
	if err != nil {
		t.Fatalf("failed to seal core: %v", err)
	}

	if !c.Sealed() {
		t.Fatal("should be sealed")
	}

	for _, key := range unsealKeys {
		if _, err := TestCoreUnseal(c, TestKeyCopy(key)); err != nil {
			t.Fatalf("unseal err: %s", err)
		}
	}

	if c.Sealed() {
		t.Fatal("should not be sealed")
	}

	// Check if the entity is restored
	resp, err = is.HandleRequest(ctx, readReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	if resp.Data["id"] != entityID {
		t.Fatalf("failed to read the created entity after a seal/unseal cycle")
	}
}

func TestIdentityStore_MemDBEntityIndexes(t *testing.T) {
	var err error

	ctx := namespace.RootContext(nil)
	is, githubAccessor, _ := testIdentityStoreWithGithubAuth(ctx, t)

	validateMountResp := is.core.router.validateMountByAccessor(githubAccessor)
	if validateMountResp == nil {
		t.Fatal("failed to validate github auth mount")
	}

	alias1 := &identity.Alias{
		CanonicalID:   "testentityid",
		ID:            "testaliasid",
		MountAccessor: githubAccessor,
		MountType:     validateMountResp.MountType,
		Name:          "testaliasname",
		Metadata: map[string]string{
			"testkey1": "testmetadatavalue1",
			"testkey2": "testmetadatavalue2",
		},
	}

	alias2 := &identity.Alias{
		CanonicalID:   "testentityid",
		ID:            "testaliasid2",
		MountAccessor: validateMountResp.MountAccessor,
		MountType:     validateMountResp.MountType,
		Name:          "testaliasname2",
		Metadata: map[string]string{
			"testkey2": "testmetadatavalue2",
			"testkey3": "testmetadatavalue3",
		},
	}

	entity := &identity.Entity{
		ID:   "testentityid",
		Name: "testentityname",
		Metadata: map[string]string{
			"someusefulkey": "someusefulvalue",
		},
		Aliases: []*identity.Alias{
			alias1,
			alias2,
		},
	}

	entity.BucketKeyHash = is.entityPacker.BucketKeyHashByItemID(entity.ID)

	txn := is.db.Txn(true)
	defer txn.Abort()
	err = is.MemDBUpsertEntityInTxn(txn, entity)
	if err != nil {
		t.Fatal(err)
	}
	txn.Commit()

	// Fetch the entity using its ID
	entityFetched, err := is.MemDBEntityByID(entity.ID, false)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(entity, entityFetched) {
		t.Fatalf("bad: mismatched entities; expected: %#v\n actual: %#v\n", entity, entityFetched)
	}

	// Fetch the entity using its name
	entityFetched, err = is.MemDBEntityByName(namespace.RootContext(nil), entity.Name, false)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(entity, entityFetched) {
		t.Fatalf("entity mismatched entities; expected: %#v\n actual: %#v\n", entity, entityFetched)
	}

	txn = is.db.Txn(false)
	entitiesFetched, err := is.MemDBEntitiesByBucketEntryKeyHashInTxn(txn, entity.BucketKeyHash)
	if err != nil {
		t.Fatal(err)
	}

	if len(entitiesFetched) != 1 {
		t.Fatalf("bad: length of entities; expected: 1, actual: %d", len(entitiesFetched))
	}

	err = is.MemDBDeleteEntityByID(entity.ID)
	if err != nil {
		t.Fatal(err)
	}

	entityFetched, err = is.MemDBEntityByID(entity.ID, false)
	if err != nil {
		t.Fatal(err)
	}

	if entityFetched != nil {
		t.Fatalf("bad: entity; expected: nil, actual: %#v\n", entityFetched)
	}

	entityFetched, err = is.MemDBEntityByName(namespace.RootContext(nil), entity.Name, false)
	if err != nil {
		t.Fatal(err)
	}

	if entityFetched != nil {
		t.Fatalf("bad: entity; expected: nil, actual: %#v\n", entityFetched)
	}

}

func TestIdentityStore_EntityCRUD(t *testing.T) {
	var err error
	var resp *logical.Response

	ctx := namespace.RootContext(nil)
	is, _, _ := testIdentityStoreWithGithubAuth(ctx, t)

	registerData := map[string]interface{}{
		"name":     "testentityname",
		"metadata": []string{"someusefulkey=someusefulvalue"},
		"policies": []string{"testpolicy1", "testpolicy2"},
	}

	registerReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "entity",
		Data:      registerData,
	}

	// Register the entity
	resp, err = is.HandleRequest(ctx, registerReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	idRaw, ok := resp.Data["id"]
	if !ok {
		t.Fatalf("entity id not present in response")
	}
	id := idRaw.(string)
	if id == "" {
		t.Fatalf("invalid entity id")
	}

	readReq := &logical.Request{
		Path:      "entity/id/" + id,
		Operation: logical.ReadOperation,
	}

	resp, err = is.HandleRequest(ctx, readReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	if resp.Data["id"] != id ||
		resp.Data["name"] != registerData["name"] ||
		!reflect.DeepEqual(resp.Data["policies"], registerData["policies"]) {
		t.Fatalf("bad: entity response")
	}

	updateData := map[string]interface{}{
		"name":     "updatedentityname",
		"metadata": []string{"updatedkey=updatedvalue"},
		"policies": []string{"updatedpolicy1", "updatedpolicy2"},
	}

	updateReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "entity/id/" + id,
		Data:      updateData,
	}

	resp, err = is.HandleRequest(ctx, updateReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	resp, err = is.HandleRequest(ctx, readReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	if resp.Data["id"] != id ||
		resp.Data["name"] != updateData["name"] ||
		!reflect.DeepEqual(resp.Data["policies"], updateData["policies"]) {
		t.Fatalf("bad: entity response after update; resp: %#v\n updateData: %#v\n", resp.Data, updateData)
	}

	deleteReq := &logical.Request{
		Path:      "entity/id/" + id,
		Operation: logical.DeleteOperation,
	}

	resp, err = is.HandleRequest(ctx, deleteReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	resp, err = is.HandleRequest(ctx, readReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	if resp != nil {
		t.Fatalf("expected a nil response; actual: %#v\n", resp)
	}
}

func TestIdentityStore_MergeEntitiesByID(t *testing.T) {
	var err error
	var resp *logical.Response

	ctx := namespace.RootContext(nil)
	is, githubAccessor, _ := testIdentityStoreWithGithubAuth(ctx, t)

	registerData := map[string]interface{}{
		"name":     "testentityname2",
		"metadata": []string{"someusefulkey=someusefulvalue"},
	}

	registerData2 := map[string]interface{}{
		"name":     "testentityname",
		"metadata": []string{"someusefulkey=someusefulvalue"},
	}

	aliasRegisterData1 := map[string]interface{}{
		"name":           "testaliasname1",
		"mount_accessor": githubAccessor,
		"metadata":       []string{"organization=hashicorp", "team=vault"},
	}

	aliasRegisterData2 := map[string]interface{}{
		"name":           "testaliasname2",
		"mount_accessor": githubAccessor,
		"metadata":       []string{"organization=hashicorp", "team=vault"},
	}

	aliasRegisterData3 := map[string]interface{}{
		"name":           "testaliasname3",
		"mount_accessor": githubAccessor,
		"metadata":       []string{"organization=hashicorp", "team=vault"},
	}

	aliasRegisterData4 := map[string]interface{}{
		"name":           "testaliasname4",
		"mount_accessor": githubAccessor,
		"metadata":       []string{"organization=hashicorp", "team=vault"},
	}

	registerReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "entity",
		Data:      registerData,
	}

	// Register the entity
	resp, err = is.HandleRequest(ctx, registerReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	entityID1 := resp.Data["id"].(string)

	// Set entity ID in alias registration data and register alias
	aliasRegisterData1["entity_id"] = entityID1
	aliasRegisterData2["entity_id"] = entityID1

	aliasReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "alias",
		Data:      aliasRegisterData1,
	}

	// Register the alias
	resp, err = is.HandleRequest(ctx, aliasReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	// Register the alias
	aliasReq.Data = aliasRegisterData2
	resp, err = is.HandleRequest(ctx, aliasReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	entity1, err := is.MemDBEntityByID(entityID1, false)
	if err != nil {
		t.Fatal(err)
	}
	if entity1 == nil {
		t.Fatalf("failed to create entity: %v", err)
	}
	if len(entity1.Aliases) != 2 {
		t.Fatalf("bad: number of aliases in entity; expected: 2, actual: %d", len(entity1.Aliases))
	}

	registerReq.Data = registerData2
	// Register another entity
	resp, err = is.HandleRequest(ctx, registerReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	entityID2 := resp.Data["id"].(string)
	// Set entity ID in alias registration data and register alias
	aliasRegisterData3["entity_id"] = entityID2
	aliasRegisterData4["entity_id"] = entityID2

	aliasReq = &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "alias",
		Data:      aliasRegisterData3,
	}

	// Register the alias
	resp, err = is.HandleRequest(ctx, aliasReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	// Register the alias
	aliasReq.Data = aliasRegisterData4
	resp, err = is.HandleRequest(ctx, aliasReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	entity2, err := is.MemDBEntityByID(entityID2, false)
	if err != nil {
		t.Fatal(err)
	}
	if entity2 == nil {
		t.Fatalf("failed to create entity: %v", err)
	}

	if len(entity2.Aliases) != 2 {
		t.Fatalf("bad: number of aliases in entity; expected: 2, actual: %d", len(entity2.Aliases))
	}

	mergeData := map[string]interface{}{
		"to_entity_id":    entityID1,
		"from_entity_ids": []string{entityID2},
	}
	mergeReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "entity/merge",
		Data:      mergeData,
	}

	resp, err = is.HandleRequest(ctx, mergeReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	entityReq := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "entity/id/" + entityID2,
	}
	resp, err = is.HandleRequest(ctx, entityReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	if resp != nil {
		t.Fatalf("entity should have been deleted")
	}

	entityReq.Path = "entity/id/" + entityID1
	resp, err = is.HandleRequest(ctx, entityReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	entity2Aliases := resp.Data["aliases"].([]interface{})
	if len(entity2Aliases) != 4 {
		t.Fatalf("bad: number of aliases in entity; expected: 4, actual: %d", len(entity2Aliases))
	}

	for _, aliasRaw := range entity2Aliases {
		alias := aliasRaw.(map[string]interface{})
		aliasLookedUp, err := is.MemDBAliasByID(alias["id"].(string), false, false)
		if err != nil {
			t.Fatal(err)
		}
		if aliasLookedUp == nil {
			t.Fatalf("index for alias id %q is not updated", alias["id"].(string))
		}
	}
}
