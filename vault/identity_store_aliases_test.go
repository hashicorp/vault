package vault

import (
	"reflect"
	"strings"
	"testing"

	credGithub "github.com/hashicorp/vault/builtin/credential/github"
	credUserpass "github.com/hashicorp/vault/builtin/credential/userpass"
	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
)

// Issue 5729
func TestIdentityStore_DuplicateAliases(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)

	resp, err := c.systemBackend.HandleRequest(namespace.RootContext(nil), &logical.Request{
		Path:      "auth",
		Operation: logical.ReadOperation,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
	}

	tokenMountAccessor := resp.Data["token/"].(map[string]interface{})["accessor"].(string)

	// Create an entity and attach an alias to it
	resp, err = c.identityStore.HandleRequest(namespace.RootContext(nil), &logical.Request{
		Path:      "entity-alias",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"mount_accessor": tokenMountAccessor,
			"name":           "testaliasname",
		},
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
	}
	aliasID := resp.Data["id"].(string)

	// Create another entity without an alias
	resp, err = c.identityStore.HandleRequest(namespace.RootContext(nil), &logical.Request{
		Path:      "entity",
		Operation: logical.UpdateOperation,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
	}
	entityID2 := resp.Data["id"].(string)

	// Set the second entity ID as the canonical ID for the previous alias,
	// initiating an alias transfer
	resp, err = c.identityStore.HandleRequest(namespace.RootContext(nil), &logical.Request{
		Path:      "entity-alias/id/" + aliasID,
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"canonical_id": entityID2,
		},
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
	}

	// Read the new entity
	resp, err = c.identityStore.HandleRequest(namespace.RootContext(nil), &logical.Request{
		Path:      "entity/id/" + entityID2,
		Operation: logical.ReadOperation,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
	}

	// Ensure that there is only one alias
	aliases := resp.Data["aliases"].([]interface{})
	if len(aliases) != 1 {
		t.Fatalf("bad: length of aliases; expected: %d, actual: %d", 1, len(aliases))
	}

	// Ensure that no merging activity has taken place
	if len(aliases[0].(map[string]interface{})["merged_from_canonical_ids"].([]string)) != 0 {
		t.Fatalf("expected no merging to take place")
	}
}

func TestIdentityStore_CaseInsensitiveEntityAliasName(t *testing.T) {
	ctx := namespace.RootContext(nil)
	i, accessor, _ := testIdentityStoreWithGithubAuth(ctx, t)

	// Create an entity
	resp, err := i.HandleRequest(ctx, &logical.Request{
		Path:      "entity",
		Operation: logical.UpdateOperation,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err:%v\nresp: %#v", err, resp)
	}
	entityID := resp.Data["id"].(string)

	testAliasName := "testAliasName"
	// Create a case sensitive alias name
	resp, err = i.HandleRequest(ctx, &logical.Request{
		Path:      "entity-alias",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"mount_accessor": accessor,
			"canonical_id":   entityID,
			"name":           testAliasName,
		},
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err:%v\nresp: %#v", err, resp)
	}
	aliasID := resp.Data["id"].(string)

	// Ensure that reading the alias returns case sensitive alias name
	resp, err = i.HandleRequest(ctx, &logical.Request{
		Path:      "entity-alias/id/" + aliasID,
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
		Path:      "entity-alias/id/" + aliasID,
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"mount_accessor": accessor,
			"canonical_id":   entityID,
			"name":           strings.ToLower(testAliasName),
		},
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err:%v\nresp: %#v", err, resp)
	}

	// Ensure that reading the alias returns lower cased alias name
	resp, err = i.HandleRequest(ctx, &logical.Request{
		Path:      "entity-alias/id/" + aliasID,
		Operation: logical.ReadOperation,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err:%v\nresp: %#v", err, resp)
	}
	aliasName = resp.Data["name"].(string)
	if aliasName != strings.ToLower(testAliasName) {
		t.Fatalf("bad alias name; expected: %q, actual: %q", testAliasName, aliasName)
	}

	// Ensure that there is one entity alias
	resp, err = i.HandleRequest(ctx, &logical.Request{
		Path:      "entity-alias/id",
		Operation: logical.ListOperation,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err:%v\nresp: %#v", err, resp)
	}
	if len(resp.Data["keys"].([]string)) != 1 {
		t.Fatalf("bad length of entity aliases; expected: 1, actual: %d", len(resp.Data["keys"].([]string)))
	}
}

// This test is required because MemDB does not take care of ensuring
// uniqueness of indexes that are marked unique.
func TestIdentityStore_AliasSameAliasNames(t *testing.T) {
	var err error
	var resp *logical.Response

	ctx := namespace.RootContext(nil)
	is, githubAccessor, _ := testIdentityStoreWithGithubAuth(ctx, t)

	aliasData := map[string]interface{}{
		"name":           "testaliasname",
		"mount_accessor": githubAccessor,
	}

	aliasReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "entity-alias",
		Data:      aliasData,
	}

	// Register an alias
	resp, err = is.HandleRequest(ctx, aliasReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	// Register another alias with same name
	resp, err = is.HandleRequest(ctx, aliasReq)
	if err != nil {
		t.Fatal(err)
	}
	if resp != nil {
		t.Fatalf("expected no response since this modification should be idempotent")
	}
}

func TestIdentityStore_MemDBAliasIndexes(t *testing.T) {
	var err error

	ctx := namespace.RootContext(nil)
	is, githubAccessor, _ := testIdentityStoreWithGithubAuth(ctx, t)
	if is == nil {
		t.Fatal("failed to create test identity store")
	}

	validateMountResp := is.router.ValidateMountByAccessor(githubAccessor)
	if validateMountResp == nil {
		t.Fatal("failed to validate github auth mount")
	}

	entity := &identity.Entity{
		ID:   "testentityid",
		Name: "testentityname",
	}

	entity.BucketKey = is.entityPacker.BucketKey(entity.ID)

	txn := is.db.Txn(true)
	defer txn.Abort()
	err = is.MemDBUpsertEntityInTxn(txn, entity)
	if err != nil {
		t.Fatal(err)
	}
	txn.Commit()

	alias := &identity.Alias{
		CanonicalID:   entity.ID,
		ID:            "testaliasid",
		MountAccessor: githubAccessor,
		MountType:     validateMountResp.MountType,
		Name:          "testaliasname",
		Metadata: map[string]string{
			"testkey1": "testmetadatavalue1",
			"testkey2": "testmetadatavalue2",
		},
	}

	txn = is.db.Txn(true)
	defer txn.Abort()
	err = is.MemDBUpsertAliasInTxn(txn, alias, false)
	if err != nil {
		t.Fatal(err)
	}
	txn.Commit()

	aliasFetched, err := is.MemDBAliasByID("testaliasid", false, false)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(alias, aliasFetched) {
		t.Fatalf("bad: mismatched aliases; expected: %#v\n actual: %#v\n", alias, aliasFetched)
	}

	aliasFetched, err = is.MemDBAliasByFactors(validateMountResp.MountAccessor, "testaliasname", false, false)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(alias, aliasFetched) {
		t.Fatalf("bad: mismatched aliases; expected: %#v\n actual: %#v\n", alias, aliasFetched)
	}

	alias2 := &identity.Alias{
		CanonicalID:   entity.ID,
		ID:            "testaliasid2",
		MountAccessor: validateMountResp.MountAccessor,
		MountType:     validateMountResp.MountType,
		Name:          "testaliasname2",
		Metadata: map[string]string{
			"testkey1": "testmetadatavalue1",
			"testkey3": "testmetadatavalue3",
		},
	}

	txn = is.db.Txn(true)
	defer txn.Abort()
	err = is.MemDBUpsertAliasInTxn(txn, alias2, false)
	if err != nil {
		t.Fatal(err)
	}
	err = is.MemDBDeleteAliasByIDInTxn(txn, "testaliasid", false)
	if err != nil {
		t.Fatal(err)
	}
	txn.Commit()

	aliasFetched, err = is.MemDBAliasByID("testaliasid", false, false)
	if err != nil {
		t.Fatal(err)
	}

	if aliasFetched != nil {
		t.Fatalf("expected a nil alias")
	}
}

func TestIdentityStore_AliasRegister(t *testing.T) {
	var err error
	var resp *logical.Response

	ctx := namespace.RootContext(nil)
	is, githubAccessor, _ := testIdentityStoreWithGithubAuth(ctx, t)

	if is == nil {
		t.Fatal("failed to create test alias store")
	}

	aliasData := map[string]interface{}{
		"name":           "testaliasname",
		"mount_accessor": githubAccessor,
		"metadata":       []string{"organization=hashicorp", "team=vault"},
	}

	aliasReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "entity-alias",
		Data:      aliasData,
	}

	// Register the alias
	resp, err = is.HandleRequest(ctx, aliasReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	idRaw, ok := resp.Data["id"]
	if !ok {
		t.Fatalf("alias id not present in alias register response")
	}

	id := idRaw.(string)
	if id == "" {
		t.Fatalf("invalid alias id in alias register response")
	}

	entityIDRaw, ok := resp.Data["canonical_id"]
	if !ok {
		t.Fatalf("entity id not present in alias register response")
	}

	entityID := entityIDRaw.(string)
	if entityID == "" {
		t.Fatalf("invalid entity id in alias register response")
	}
}

func TestIdentityStore_AliasUpdate(t *testing.T) {
	var err error
	var resp *logical.Response
	ctx := namespace.RootContext(nil)
	is, githubAccessor, _ := testIdentityStoreWithGithubAuth(ctx, t)

	aliasData := map[string]interface{}{
		"name":           "testaliasname",
		"mount_accessor": githubAccessor,
	}

	aliasReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "entity-alias",
		Data:      aliasData,
	}

	// This will create an alias and a corresponding entity
	resp, err = is.HandleRequest(ctx, aliasReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	aliasID := resp.Data["id"].(string)
	customMetadata := make(map[string]string)
	customMetadata["foo"] = "abc"

	updateData := map[string]interface{}{
		"name":            "updatedaliasname",
		"mount_accessor":  githubAccessor,
		"custom_metadata": customMetadata,
	}

	aliasReq.Data = updateData
	aliasReq.Path = "entity-alias/id/" + aliasID
	resp, err = is.HandleRequest(ctx, aliasReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	aliasReq.Operation = logical.ReadOperation
	resp, err = is.HandleRequest(ctx, aliasReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	if resp.Data["name"] != "updatedaliasname" {
		t.Fatalf("failed to update alias information; \n response data: %#v\n", resp.Data)
	}
	if !reflect.DeepEqual(resp.Data["custom_metadata"], customMetadata) {
		t.Fatalf("failed to update alias information; \n response data: %#v\n", resp.Data)
	}
}

// Test to check that the alias cannot be updated with a new entity
// which already has an alias for the mount on the alias to be updated
func TestIdentityStore_AliasUpdate_InvalidEntity(t *testing.T) {
	var err error
	var resp *logical.Response
	ctx := namespace.RootContext(nil)
	is, githubAccessor, _ := testIdentityStoreWithGithubAuth(ctx, t)

	// Create 2 entities and 1 alias on each, against the same github mount
	resp, err = is.HandleRequest(ctx, &logical.Request{
		Path:      "entity",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"name": "testentity1",
		},
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err:%v\nresp: %#v", err, resp)
	}
	entity1ID := resp.Data["id"].(string)

	alias1Data := map[string]interface{}{
		"name":           "testaliasname1",
		"mount_accessor": githubAccessor,
		"canonical_id":   entity1ID,
	}

	aliasReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "entity-alias",
		Data:      alias1Data,
	}

	// This will create an alias against the requested entity
	resp, err = is.HandleRequest(ctx, aliasReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	resp, err = is.HandleRequest(ctx, &logical.Request{
		Path:      "entity",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"name": "testentity2",
		},
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err:%v\nresp: %#v", err, resp)
	}
	entity2ID := resp.Data["id"].(string)

	alias2Data := map[string]interface{}{
		"name":           "testaliasname2",
		"mount_accessor": githubAccessor,
		"canonical_id":   entity2ID,
	}

	aliasReq.Data = alias2Data

	// This will create an alias against the requested entity
	resp, err = is.HandleRequest(ctx, aliasReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	alias2ID := resp.Data["id"].(string)

	// Attempt to update the second alias to point to the first entity
	updateData := map[string]interface{}{
		"canonical_id": entity1ID,
	}

	aliasReq.Data = updateData
	aliasReq.Path = "entity-alias/id/" + alias2ID
	resp, err = is.HandleRequest(ctx, aliasReq)

	if err != nil {
		t.Fatal(err)
	}

	if resp == nil || !resp.IsError() {
		t.Fatalf("expected an error as alias on the github accessor exists for testentity1")
	}

}

// Test that the alias cannot be changed to a mount for which
// the entity already has an alias
func TestIdentityStore_AliasUpdate__DuplicateAccessor(t *testing.T) {
	var err error
	var resp *logical.Response
	ctx := namespace.RootContext(nil)

	// Setup 2 auth backends, github and userpass
	err = AddTestCredentialBackend("github", credGithub.Factory)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	err = AddTestCredentialBackend("userpass", credUserpass.Factory)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	c, _, _ := TestCoreUnsealed(t)

	githubMe := &MountEntry{
		Table:       credentialTableType,
		Path:        "github/",
		Type:        "github",
		Description: "github auth",
	}

	err = c.enableCredential(ctx, githubMe)
	if err != nil {
		t.Fatal(err)
	}

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

	is, githubAccessor := c.identityStore, githubMe.Accessor

	// Create 1 entity and 2 aliases on it, one for each mount
	resp, err = is.HandleRequest(ctx, &logical.Request{
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

	alias1Data := map[string]interface{}{
		"name":           "testaliasname1",
		"mount_accessor": githubAccessor,
		"canonical_id":   entityID,
	}

	aliasReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "entity-alias",
		Data:      alias1Data,
	}

	// This will create an alias against the requested entity
	resp, err = is.HandleRequest(ctx, aliasReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	alias2Data := map[string]interface{}{
		"name":           "testaliasname2",
		"mount_accessor": userpassMe.Accessor,
		"canonical_id":   entityID,
	}

	aliasReq.Data = alias2Data

	// This will create an alias against the requested entity
	resp, err = is.HandleRequest(ctx, aliasReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	alias2ID := resp.Data["id"].(string)

	// Attempt to update the userpass mount to point to the github mount
	updateData := map[string]interface{}{
		"mount_accessor": githubAccessor,
	}

	aliasReq.Data = updateData
	aliasReq.Path = "entity-alias/id/" + alias2ID
	resp, err = is.HandleRequest(ctx, aliasReq)

	if err != nil {
		t.Fatal(err)
	}

	if resp == nil || !resp.IsError() {
		t.Fatalf("expected an error as an alias on the github accessor already exists for testentity")
	}

}

// Test that alias creation fails if an alias for the specified mount
// and entity has already been created
func TestIdentityStore_AliasCreate_DuplicateAccessor(t *testing.T) {
	var err error
	var resp *logical.Response
	ctx := namespace.RootContext(nil)
	is, githubAccessor, _ := testIdentityStoreWithGithubAuth(ctx, t)

	resp, err = is.HandleRequest(ctx, &logical.Request{
		Path:      "entity",
		Operation: logical.UpdateOperation,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err:%v\nresp: %#v", err, resp)
	}
	entityID := resp.Data["id"].(string)

	aliasData := map[string]interface{}{
		"name":           "testaliasname",
		"mount_accessor": githubAccessor,
		"canonical_id":   entityID,
	}

	aliasReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "entity-alias",
		Data:      aliasData,
	}

	// This will create an alias against the requested entity
	resp, err = is.HandleRequest(ctx, aliasReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	aliasData["name"] = "testaliasname2"
	aliasReq = &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "entity-alias",
		Data:      aliasData,
	}

	// This will try to create a new alias with the same accessor and entity
	resp, err = is.HandleRequest(ctx, aliasReq)

	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || !resp.IsError() {
		t.Fatalf("expected an error as alias already exists for this accessor and entity")
	}
}

func TestIdentityStore_AliasUpdate_ByID(t *testing.T) {
	var err error
	var resp *logical.Response
	ctx := namespace.RootContext(nil)
	is, githubAccessor, _ := testIdentityStoreWithGithubAuth(ctx, t)

	updateData := map[string]interface{}{
		"name":           "updatedaliasname",
		"mount_accessor": githubAccessor,
	}

	updateReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "entity-alias/id/invalidaliasid",
		Data:      updateData,
	}

	// Try to update an non-existent alias
	resp, err = is.HandleRequest(ctx, updateReq)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || !resp.IsError() {
		t.Fatalf("expected an error due to invalid alias id")
	}

	customMetadata := make(map[string]string)
	customMetadata["foo"] = "abc"
	registerData := map[string]interface{}{
		"name":            "testaliasname",
		"mount_accessor":  githubAccessor,
		"custom_metadata": customMetadata,
	}

	registerReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "entity-alias",
		Data:      registerData,
	}

	resp, err = is.HandleRequest(ctx, registerReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	idRaw, ok := resp.Data["id"]
	if !ok {
		t.Fatalf("alias id not present in response")
	}
	id := idRaw.(string)
	if id == "" {
		t.Fatalf("invalid alias id")
	}

	updateReq.Path = "entity-alias/id/" + id
	resp, err = is.HandleRequest(ctx, updateReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	readReq := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      updateReq.Path,
	}
	resp, err = is.HandleRequest(ctx, readReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	if resp.Data["name"] != "updatedaliasname" {
		t.Fatalf("failed to update alias information; \n response data: %#v\n", resp.Data)
	}
	if !reflect.DeepEqual(resp.Data["custom_metadata"], customMetadata) {
		t.Fatalf("failed to update alias information; \n response data: %#v\n", resp.Data)
	}

	delete(registerReq.Data, "name")

	resp, err = is.HandleRequest(ctx, registerReq)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || !resp.IsError() {
		t.Fatalf("expected error due to missing alias name")
	}

	registerReq.Data["name"] = "testaliasname"
	delete(registerReq.Data, "mount_accessor")

	resp, err = is.HandleRequest(ctx, registerReq)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || !resp.IsError() {
		t.Fatalf("expected error due to missing mount accessor")
	}
}

func TestIdentityStore_AliasReadDelete(t *testing.T) {
	var err error
	var resp *logical.Response

	ctx := namespace.RootContext(nil)
	is, githubAccessor, _ := testIdentityStoreWithGithubAuth(ctx, t)

	customMetadata := make(map[string]string)
	customMetadata["foo"] = "abc"
	registerData := map[string]interface{}{
		"name":            "testaliasname",
		"mount_accessor":  githubAccessor,
		"metadata":        []string{"organization=hashicorp", "team=vault"},
		"custom_metadata": customMetadata,
	}

	registerReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "entity-alias",
		Data:      registerData,
	}

	resp, err = is.HandleRequest(ctx, registerReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	idRaw, ok := resp.Data["id"]
	if !ok {
		t.Fatalf("alias id not present in response")
	}
	id := idRaw.(string)
	if id == "" {
		t.Fatalf("invalid alias id")
	}

	// Read it back using alias id
	aliasReq := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "entity-alias/id/" + id,
	}
	resp, err = is.HandleRequest(ctx, aliasReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	if resp.Data["id"].(string) == "" ||
		resp.Data["canonical_id"].(string) == "" ||
		resp.Data["name"].(string) != registerData["name"] ||
		resp.Data["mount_type"].(string) != "github" || !reflect.DeepEqual(resp.Data["custom_metadata"], customMetadata) {
		t.Fatalf("bad: alias read response; \nexpected: %#v \nactual: %#v\n", registerData, resp.Data)
	}

	aliasReq.Operation = logical.DeleteOperation
	resp, err = is.HandleRequest(ctx, aliasReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	aliasReq.Operation = logical.ReadOperation
	resp, err = is.HandleRequest(ctx, aliasReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	if resp != nil {
		t.Fatalf("bad: alias read response; expected: nil, actual: %#v\n", resp)
	}
}
