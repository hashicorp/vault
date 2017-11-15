package vault

import (
	"reflect"
	"testing"

	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/logical"
)

func TestIdentityStore_ListAlias(t *testing.T) {
	var err error
	var resp *logical.Response

	is, githubAccessor, _ := testIdentityStoreWithGithubAuth(t)

	entityReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "entity",
	}
	resp, err = is.HandleRequest(entityReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	if resp == nil {
		t.Fatalf("expected a non-nil response")
	}
	entityID := resp.Data["id"].(string)

	// Create an alias
	aliasData := map[string]interface{}{
		"name":           "testaliasname",
		"mount_accessor": githubAccessor,
	}
	aliasReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "entity-alias",
		Data:      aliasData,
	}
	resp, err = is.HandleRequest(aliasReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	aliasData["name"] = "entityalias"
	aliasData["entity_id"] = entityID
	resp, err = is.HandleRequest(aliasReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	listReq := &logical.Request{
		Operation: logical.ListOperation,
		Path:      "entity-alias/id",
	}
	resp, err = is.HandleRequest(listReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	keys := resp.Data["keys"].([]string)
	if len(keys) != 2 {
		t.Fatalf("bad: lengh of alias IDs listed; expected: 2, actual: %d", len(keys))
	}
}

// This test is required because MemDB does not take care of ensuring
// uniqueness of indexes that are marked unique.
func TestIdentityStore_AliasSameAliasNames(t *testing.T) {
	var err error
	var resp *logical.Response
	is, githubAccessor, _ := testIdentityStoreWithGithubAuth(t)

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
	resp, err = is.HandleRequest(aliasReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	// Register another alias with same name
	resp, err = is.HandleRequest(aliasReq)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || !resp.IsError() {
		t.Fatalf("expected an error due to alias name not being unique")
	}
}

func TestIdentityStore_MemDBAliasIndexes(t *testing.T) {
	var err error

	is, githubAccessor, _ := testIdentityStoreWithGithubAuth(t)
	if is == nil {
		t.Fatal("failed to create test identity store")
	}

	validateMountResp := is.validateMountAccessorFunc(githubAccessor)
	if validateMountResp == nil {
		t.Fatal("failed to validate github auth mount")
	}

	entity := &identity.Entity{
		ID:   "testentityid",
		Name: "testentityname",
	}

	entity.BucketKeyHash = is.entityPacker.BucketKeyHashByItemID(entity.ID)

	err = is.MemDBUpsertEntity(entity)
	if err != nil {
		t.Fatal(err)
	}

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

	err = is.MemDBUpsertAlias(alias, false)
	if err != nil {
		t.Fatal(err)
	}

	aliasFetched, err := is.MemDBAliasByID("testaliasid", false, false)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(alias, aliasFetched) {
		t.Fatalf("bad: mismatched aliases; expected: %#v\n actual: %#v\n", alias, aliasFetched)
	}

	aliasFetched, err = is.MemDBAliasByCanonicalID(entity.ID, false, false)
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

	aliasesFetched, err := is.MemDBAliasesByMetadata(map[string]string{
		"testkey1": "testmetadatavalue1",
	}, false, false)
	if err != nil {
		t.Fatal(err)
	}

	if len(aliasesFetched) != 1 {
		t.Fatalf("bad: length of aliases; expected: 1, actual: %d", len(aliasesFetched))
	}

	if !reflect.DeepEqual(alias, aliasesFetched[0]) {
		t.Fatalf("bad: mismatched aliases; expected: %#v\n actual: %#v\n", alias, aliasFetched)
	}

	aliasesFetched, err = is.MemDBAliasesByMetadata(map[string]string{
		"testkey2": "testmetadatavalue2",
	}, false, false)
	if err != nil {
		t.Fatal(err)
	}

	if len(aliasesFetched) != 1 {
		t.Fatalf("bad: length of aliases; expected: 1, actual: %d", len(aliasesFetched))
	}

	if !reflect.DeepEqual(alias, aliasesFetched[0]) {
		t.Fatalf("bad: mismatched aliases; expected: %#v\n actual: %#v\n", alias, aliasFetched)
	}

	aliasesFetched, err = is.MemDBAliasesByMetadata(map[string]string{
		"testkey1": "testmetadatavalue1",
		"testkey2": "testmetadatavalue2",
	}, false, false)
	if err != nil {
		t.Fatal(err)
	}

	if len(aliasesFetched) != 1 {
		t.Fatalf("bad: length of aliases; expected: 1, actual: %d", len(aliasesFetched))
	}

	if !reflect.DeepEqual(alias, aliasesFetched[0]) {
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

	err = is.MemDBUpsertAlias(alias2, false)
	if err != nil {
		t.Fatal(err)
	}

	aliasesFetched, err = is.MemDBAliasesByMetadata(map[string]string{
		"testkey1": "testmetadatavalue1",
	}, false, false)
	if err != nil {
		t.Fatal(err)
	}

	if len(aliasesFetched) != 2 {
		t.Fatalf("bad: length of aliases; expected: 2, actual: %d", len(aliasesFetched))
	}

	aliasesFetched, err = is.MemDBAliasesByMetadata(map[string]string{
		"testkey3": "testmetadatavalue3",
	}, false, false)
	if err != nil {
		t.Fatal(err)
	}

	if len(aliasesFetched) != 1 {
		t.Fatalf("bad: length of aliases; expected: 1, actual: %d", len(aliasesFetched))
	}

	err = is.MemDBDeleteAliasByID("testaliasid", false)
	if err != nil {
		t.Fatal(err)
	}

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

	is, githubAccessor, _ := testIdentityStoreWithGithubAuth(t)

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
	resp, err = is.HandleRequest(aliasReq)
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
	is, githubAccessor, _ := testIdentityStoreWithGithubAuth(t)

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

	// This will create an alias and a corresponding entity
	resp, err = is.HandleRequest(aliasReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	aliasID := resp.Data["id"].(string)

	updateData := map[string]interface{}{
		"name":           "updatedaliasname",
		"mount_accessor": githubAccessor,
		"metadata":       []string{"organization=updatedorganization", "team=updatedteam"},
	}

	aliasReq.Data = updateData
	aliasReq.Path = "entity-alias/id/" + aliasID
	resp, err = is.HandleRequest(aliasReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	aliasReq.Operation = logical.ReadOperation
	resp, err = is.HandleRequest(aliasReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	aliasMetadata := resp.Data["metadata"].(map[string]string)
	updatedOrg := aliasMetadata["organization"]
	updatedTeam := aliasMetadata["team"]

	if resp.Data["name"] != "updatedaliasname" || updatedOrg != "updatedorganization" || updatedTeam != "updatedteam" {
		t.Fatalf("failed to update alias information; \n response data: %#v\n", resp.Data)
	}
}

func TestIdentityStore_AliasUpdate_ByID(t *testing.T) {
	var err error
	var resp *logical.Response
	is, githubAccessor, _ := testIdentityStoreWithGithubAuth(t)

	updateData := map[string]interface{}{
		"name":           "updatedaliasname",
		"mount_accessor": githubAccessor,
		"metadata":       []string{"organization=updatedorganization", "team=updatedteam"},
	}

	updateReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "entity-alias/id/invalidaliasid",
		Data:      updateData,
	}

	// Try to update an non-existent alias
	resp, err = is.HandleRequest(updateReq)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || !resp.IsError() {
		t.Fatalf("expected an error due to invalid alias id")
	}

	registerData := map[string]interface{}{
		"name":           "testaliasname",
		"mount_accessor": githubAccessor,
		"metadata":       []string{"organization=hashicorp", "team=vault"},
	}

	registerReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "entity-alias",
		Data:      registerData,
	}

	resp, err = is.HandleRequest(registerReq)
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
	resp, err = is.HandleRequest(updateReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	readReq := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      updateReq.Path,
	}
	resp, err = is.HandleRequest(readReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	aliasMetadata := resp.Data["metadata"].(map[string]string)
	updatedOrg := aliasMetadata["organization"]
	updatedTeam := aliasMetadata["team"]

	if resp.Data["name"] != "updatedaliasname" || updatedOrg != "updatedorganization" || updatedTeam != "updatedteam" {
		t.Fatalf("failed to update alias information; \n response data: %#v\n", resp.Data)
	}

	delete(registerReq.Data, "name")

	resp, err = is.HandleRequest(registerReq)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || !resp.IsError() {
		t.Fatalf("expected error due to missing alias name")
	}

	registerReq.Data["name"] = "testaliasname"
	delete(registerReq.Data, "mount_accessor")

	resp, err = is.HandleRequest(registerReq)
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

	is, githubAccessor, _ := testIdentityStoreWithGithubAuth(t)

	registerData := map[string]interface{}{
		"name":           "testaliasname",
		"mount_accessor": githubAccessor,
		"metadata":       []string{"organization=hashicorp", "team=vault"},
	}

	registerReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "entity-alias",
		Data:      registerData,
	}

	resp, err = is.HandleRequest(registerReq)
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
	resp, err = is.HandleRequest(aliasReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	if resp.Data["id"].(string) == "" ||
		resp.Data["canonical_id"].(string) == "" ||
		resp.Data["name"].(string) != registerData["name"] ||
		resp.Data["mount_type"].(string) != "github" {
		t.Fatalf("bad: alias read response; \nexpected: %#v \nactual: %#v\n", registerData, resp.Data)
	}

	aliasReq.Operation = logical.DeleteOperation
	resp, err = is.HandleRequest(aliasReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	aliasReq.Operation = logical.ReadOperation
	resp, err = is.HandleRequest(aliasReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	if resp != nil {
		t.Fatalf("bad: alias read response; expected: nil, actual: %#v\n", resp)
	}
}
