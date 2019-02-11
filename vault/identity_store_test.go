package vault

import (
	"context"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/golang/protobuf/ptypes"
	uuid "github.com/hashicorp/go-uuid"
	credGithub "github.com/hashicorp/vault/builtin/credential/github"
	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/storagepacker"
	"github.com/hashicorp/vault/logical"
)

func TestIdentityStore_UnsealingWhenConflictingAliasNames(t *testing.T) {
	err := AddTestCredentialBackend("github", credGithub.Factory)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	c, unsealKey, root := TestCoreUnsealed(t)

	meGH := &MountEntry{
		Table:       credentialTableType,
		Path:        "github/",
		Type:        "github",
		Description: "github auth",
	}

	err = c.enableCredential(namespace.RootContext(nil), meGH)
	if err != nil {
		t.Fatal(err)
	}

	alias := &identity.Alias{
		ID:            "alias1",
		CanonicalID:   "entity1",
		MountType:     "github",
		MountAccessor: meGH.Accessor,
		Name:          "githubuser",
	}
	entity := &identity.Entity{
		ID:       "entity1",
		Name:     "name1",
		Policies: []string{"foo", "bar"},
		Aliases: []*identity.Alias{
			alias,
		},
		NamespaceID: namespace.RootNamespaceID,
	}
	entity.BucketKeyHash = c.identityStore.entityPacker.BucketKeyHashByItemID(entity.ID)

	err = c.identityStore.upsertEntity(namespace.RootContext(nil), entity, nil, true)
	if err != nil {
		t.Fatal(err)
	}

	alias2 := &identity.Alias{
		ID:            "alias2",
		CanonicalID:   "entity2",
		MountType:     "github",
		MountAccessor: meGH.Accessor,
		Name:          "GITHUBUSER",
	}
	entity2 := &identity.Entity{
		ID:       "entity2",
		Name:     "name2",
		Policies: []string{"foo", "bar"},
		Aliases: []*identity.Alias{
			alias2,
		},
		NamespaceID: namespace.RootNamespaceID,
	}
	entity2.BucketKeyHash = c.identityStore.entityPacker.BucketKeyHashByItemID(entity2.ID)

	// Persist the second entity directly without the regular flow. This will skip
	// merging of these enties.
	entity2Any, err := ptypes.MarshalAny(entity2)
	if err != nil {
		t.Fatal(err)
	}
	item := &storagepacker.Item{
		ID:      entity2.ID,
		Message: entity2Any,
	}
	if err = c.identityStore.entityPacker.PutItem(item); err != nil {
		t.Fatal(err)
	}

	// Seal and ensure that unseal works
	if err = c.Seal(root); err != nil {
		t.Fatal(err)
	}

	var unsealed bool
	for i := 0; i < 3; i++ {
		unsealed, err = c.Unseal(unsealKey[i])
		if err != nil {
			t.Fatal(err)
		}
	}
	if !unsealed {
		t.Fatal("still sealed")
	}
}

func TestIdentityStore_EntityIDPassthrough(t *testing.T) {
	// Enable GitHub auth and initialize
	ctx := namespace.RootContext(nil)
	is, ghAccessor, core := testIdentityStoreWithGithubAuth(ctx, t)
	alias := &logical.Alias{
		MountType:     "github",
		MountAccessor: ghAccessor,
		Name:          "githubuser",
	}

	// Create an entity with GitHub alias
	entity, err := is.CreateOrFetchEntity(ctx, alias)
	if err != nil {
		t.Fatal(err)
	}
	if entity == nil {
		t.Fatalf("expected a non-nil entity")
	}

	// Create a token with the above created entity set on it
	ent := &logical.TokenEntry{
		ID:           "testtokenid",
		Path:         "test",
		Policies:     []string{"root"},
		CreationTime: time.Now().Unix(),
		EntityID:     entity.ID,
		NamespaceID:  namespace.RootNamespaceID,
	}
	if err := core.tokenStore.create(ctx, ent); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Set a request handler to the noop backend which responds with the entity
	// ID received in the request object
	requestHandler := func(ctx context.Context, req *logical.Request) (*logical.Response, error) {
		return &logical.Response{
			Data: map[string]interface{}{
				"entity_id": req.EntityID,
			},
		}, nil
	}

	noop := &NoopBackend{
		RequestHandler: requestHandler,
	}

	// Mount the noop backend
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "logical/")
	meUUID, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	err = core.router.Mount(noop, "test/backend/", &MountEntry{Path: "test/backend/", Type: "noop", UUID: meUUID, Accessor: "noop-accessor", namespace: namespace.RootNamespace}, view)
	if err != nil {
		t.Fatal(err)
	}

	// Make the request with the above created token
	resp, err := core.HandleRequest(ctx, &logical.Request{
		ClientToken: "testtokenid",
		Operation:   logical.ReadOperation,
		Path:        "test/backend/foo",
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\n err: %v", resp, err)
	}

	// Expected entity ID to be in the response
	if resp.Data["entity_id"] != entity.ID {
		t.Fatalf("expected entity ID to be passed through to the backend")
	}
}

func TestIdentityStore_CreateOrFetchEntity(t *testing.T) {
	ctx := namespace.RootContext(nil)
	is, ghAccessor, _ := testIdentityStoreWithGithubAuth(ctx, t)
	alias := &logical.Alias{
		MountType:     "github",
		MountAccessor: ghAccessor,
		Name:          "githubuser",
		Metadata: map[string]string{
			"foo": "a",
		},
	}

	entity, err := is.CreateOrFetchEntity(ctx, alias)
	if err != nil {
		t.Fatal(err)
	}
	if entity == nil {
		t.Fatalf("expected a non-nil entity")
	}

	if len(entity.Aliases) != 1 {
		t.Fatalf("bad: length of aliases; expected: 1, actual: %d", len(entity.Aliases))
	}

	if entity.Aliases[0].Name != alias.Name {
		t.Fatalf("bad: alias name; expected: %q, actual: %q", alias.Name, entity.Aliases[0].Name)
	}

	entity, err = is.CreateOrFetchEntity(ctx, alias)
	if err != nil {
		t.Fatal(err)
	}
	if entity == nil {
		t.Fatalf("expected a non-nil entity")
	}

	if len(entity.Aliases) != 1 {
		t.Fatalf("bad: length of aliases; expected: 1, actual: %d", len(entity.Aliases))
	}

	if entity.Aliases[0].Name != alias.Name {
		t.Fatalf("bad: alias name; expected: %q, actual: %q", alias.Name, entity.Aliases[0].Name)
	}
	if diff := deep.Equal(entity.Aliases[0].Metadata, map[string]string{"foo": "a"}); diff != nil {
		t.Fatal(diff)
	}

	// Add a new alias to the entity and verify its existence
	registerReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "entity-alias",
		Data: map[string]interface{}{
			"name":           "githubuser2",
			"canonical_id":   entity.ID,
			"mount_accessor": ghAccessor,
		},
	}

	resp, err := is.HandleRequest(ctx, registerReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	entity, err = is.CreateOrFetchEntity(ctx, alias)
	if err != nil {
		t.Fatal(err)
	}
	if entity == nil {
		t.Fatalf("expected a non-nil entity")
	}

	if len(entity.Aliases) != 2 {
		t.Fatalf("bad: length of aliases; expected: 2, actual: %d", len(entity.Aliases))
	}

	if entity.Aliases[1].Name != "githubuser2" {
		t.Fatalf("bad: alias name; expected: %q, actual: %q", alias.Name, "githubuser2")
	}

	if diff := deep.Equal(entity.Aliases[1].Metadata, map[string]string(nil)); diff != nil {
		t.Fatal(diff)
	}

	// Change the metadata of an existing alias and verify that
	// a the change takes effect only for the target alias.
	alias.Metadata = map[string]string{
		"foo": "zzzz",
	}

	entity, err = is.CreateOrFetchEntity(ctx, alias)
	if err != nil {
		t.Fatal(err)
	}
	if entity == nil {
		t.Fatalf("expected a non-nil entity")
	}

	if len(entity.Aliases) != 2 {
		t.Fatalf("bad: length of aliases; expected: 2, actual: %d", len(entity.Aliases))
	}

	if diff := deep.Equal(entity.Aliases[0].Metadata, map[string]string{"foo": "zzzz"}); diff != nil {
		t.Fatal(diff)
	}

	if diff := deep.Equal(entity.Aliases[1].Metadata, map[string]string(nil)); diff != nil {
		t.Fatal(diff)
	}
}

func TestIdentityStore_EntityByAliasFactors(t *testing.T) {
	var err error
	var resp *logical.Response

	ctx := namespace.RootContext(nil)
	is, ghAccessor, _ := testIdentityStoreWithGithubAuth(ctx, t)

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
	entityID := idRaw.(string)
	if entityID == "" {
		t.Fatalf("invalid entity id")
	}

	aliasData := map[string]interface{}{
		"entity_id":      entityID,
		"name":           "alias_name",
		"mount_accessor": ghAccessor,
	}
	aliasReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "alias",
		Data:      aliasData,
	}

	resp, err = is.HandleRequest(ctx, aliasReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	if resp == nil {
		t.Fatalf("expected a non-nil response")
	}

	entity, err := is.entityByAliasFactors(ghAccessor, "alias_name", false)
	if err != nil {
		t.Fatal(err)
	}
	if entity == nil {
		t.Fatalf("expected a non-nil entity")
	}
	if entity.ID != entityID {
		t.Fatalf("bad: entity ID; expected: %q actual: %q", entityID, entity.ID)
	}
}

func TestIdentityStore_WrapInfoInheritance(t *testing.T) {
	var err error
	var resp *logical.Response

	ctx := namespace.RootContext(nil)
	core, is, ts, _ := testCoreWithIdentityTokenGithub(ctx, t)

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
	entityID := idRaw.(string)
	if entityID == "" {
		t.Fatalf("invalid entity id")
	}

	// Create a token which has EntityID set and has update permissions to
	// sys/wrapping/wrap
	te := &logical.TokenEntry{
		Path:     "test",
		Policies: []string{"default", responseWrappingPolicyName},
		EntityID: entityID,
		TTL:      time.Hour,
	}
	testMakeTokenDirectly(t, ts, te)

	wrapReq := &logical.Request{
		Path:        "sys/wrapping/wrap",
		ClientToken: te.ID,
		Operation:   logical.UpdateOperation,
		Data: map[string]interface{}{
			"foo": "bar",
		},
		WrapInfo: &logical.RequestWrapInfo{
			TTL: time.Duration(5 * time.Second),
		},
	}

	resp, err = core.HandleRequest(ctx, wrapReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v, err: %v", resp, err)
	}

	if resp.WrapInfo == nil {
		t.Fatalf("expected a non-nil WrapInfo")
	}

	if resp.WrapInfo.WrappedEntityID != entityID {
		t.Fatalf("bad: WrapInfo in response not having proper entity ID set; expected: %q, actual:%q", entityID, resp.WrapInfo.WrappedEntityID)
	}
}

func TestIdentityStore_TokenEntityInheritance(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ts := c.tokenStore

	// Create a token which has EntityID set
	te := &logical.TokenEntry{
		Path:     "test",
		Policies: []string{"dev", "prod"},
		EntityID: "testentityid",
		TTL:      time.Hour,
	}
	testMakeTokenDirectly(t, ts, te)

	// Create a child token; this should inherit the EntityID
	tokenReq := &logical.Request{
		Operation:   logical.UpdateOperation,
		Path:        "create",
		ClientToken: te.ID,
	}

	ctx := namespace.RootContext(nil)
	resp, err := ts.HandleRequest(ctx, tokenReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v err: %v", err, resp)
	}

	if resp.Auth.EntityID != te.EntityID {
		t.Fatalf("bad: entity ID; expected: %v, actual: %v", te.EntityID, resp.Auth.EntityID)
	}

	// Create an orphan token; this should not inherit the EntityID
	tokenReq.Path = "create-orphan"
	resp, err = ts.HandleRequest(ctx, tokenReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v err: %v", err, resp)
	}

	if resp.Auth.EntityID != "" {
		t.Fatalf("expected entity ID to be not set")
	}
}

func TestIdentityStore_MergeConflictingAliases(t *testing.T) {
	err := AddTestCredentialBackend("github", credGithub.Factory)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	c, _, _ := TestCoreUnsealed(t)

	meGH := &MountEntry{
		Table:       credentialTableType,
		Path:        "github/",
		Type:        "github",
		Description: "github auth",
	}

	err = c.enableCredential(namespace.RootContext(nil), meGH)
	if err != nil {
		t.Fatal(err)
	}

	alias := &identity.Alias{
		ID:            "alias1",
		CanonicalID:   "entity1",
		MountType:     "github",
		MountAccessor: meGH.Accessor,
		Name:          "githubuser",
	}
	entity := &identity.Entity{
		ID:       "entity1",
		Name:     "name1",
		Policies: []string{"foo", "bar"},
		Aliases: []*identity.Alias{
			alias,
		},
		NamespaceID: namespace.RootNamespaceID,
	}
	entity.BucketKeyHash = c.identityStore.entityPacker.BucketKeyHashByItemID(entity.ID)
	err = c.identityStore.upsertEntity(namespace.RootContext(nil), entity, nil, true)
	if err != nil {
		t.Fatal(err)
	}

	alias2 := &identity.Alias{
		ID:            "alias2",
		CanonicalID:   "entity2",
		MountType:     "github",
		MountAccessor: meGH.Accessor,
		Name:          "githubuser",
	}
	entity2 := &identity.Entity{
		ID:       "entity2",
		Name:     "name2",
		Policies: []string{"bar", "baz"},
		Aliases: []*identity.Alias{
			alias2,
		},
		NamespaceID: namespace.RootNamespaceID,
	}

	entity2.BucketKeyHash = c.identityStore.entityPacker.BucketKeyHashByItemID(entity2.ID)

	err = c.identityStore.upsertEntity(namespace.RootContext(nil), entity2, nil, true)
	if err != nil {
		t.Fatal(err)
	}

	newEntity, err := c.identityStore.CreateOrFetchEntity(namespace.RootContext(nil), &logical.Alias{
		MountAccessor: meGH.Accessor,
		MountType:     "github",
		Name:          "githubuser",
	})
	if err != nil {
		t.Fatal(err)
	}
	if newEntity == nil {
		t.Fatal("nil new entity")
	}

	entityToUse := "entity1"
	if newEntity.ID == "entity1" {
		entityToUse = "entity2"
	}
	if len(newEntity.MergedEntityIDs) != 1 || newEntity.MergedEntityIDs[0] != entityToUse {
		t.Fatalf("bad merged entity ids: %v", newEntity.MergedEntityIDs)
	}
	if diff := deep.Equal(newEntity.Policies, []string{"bar", "baz", "foo"}); diff != nil {
		t.Fatal(diff)
	}

	newEntity, err = c.identityStore.MemDBEntityByID(entityToUse, false)
	if err != nil {
		t.Fatal(err)
	}
	if newEntity != nil {
		t.Fatal("got a non-nil entity")
	}
}

func testCoreWithIdentityTokenGithub(ctx context.Context, t *testing.T) (*Core, *IdentityStore, *TokenStore, string) {
	is, ghAccessor, core := testIdentityStoreWithGithubAuth(ctx, t)
	return core, is, core.tokenStore, ghAccessor
}

func testCoreWithIdentityTokenGithubRoot(ctx context.Context, t *testing.T) (*Core, *IdentityStore, *TokenStore, string, string) {
	is, ghAccessor, core, root := testIdentityStoreWithGithubAuthRoot(ctx, t)
	return core, is, core.tokenStore, ghAccessor, root
}

func testIdentityStoreWithGithubAuth(ctx context.Context, t *testing.T) (*IdentityStore, string, *Core) {
	is, ghA, c, _ := testIdentityStoreWithGithubAuthRoot(ctx, t)
	return is, ghA, c
}

// testIdentityStoreWithGithubAuthRoot returns an instance of identity store
// which is mounted by default. This function also enables the github auth
// backend to assist with testing aliases and entities that require an valid
// mount accessor of an auth backend.
func testIdentityStoreWithGithubAuthRoot(ctx context.Context, t *testing.T) (*IdentityStore, string, *Core, string) {
	// Add github credential factory to core config
	err := AddTestCredentialBackend("github", credGithub.Factory)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	c, _, root := TestCoreUnsealed(t)

	meGH := &MountEntry{
		Table:       credentialTableType,
		Path:        "github/",
		Type:        "github",
		Description: "github auth",
	}

	err = c.enableCredential(ctx, meGH)
	if err != nil {
		t.Fatal(err)
	}

	return c.identityStore, meGH.Accessor, c, root
}

func TestIdentityStore_MetadataKeyRegex(t *testing.T) {
	key := "validVALID012_-=+/"

	if !metaKeyFormatRegEx(key) {
		t.Fatal("failed to accept valid metadata key")
	}

	key = "a:b"
	if metaKeyFormatRegEx(key) {
		t.Fatal("accepted invalid metadata key")
	}
}
