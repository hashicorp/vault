// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"fmt"
	"math/rand"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/armon/go-metrics"
	"github.com/go-test/deep"
	"github.com/hashicorp/go-hclog"
	uuid "github.com/hashicorp/go-uuid"
	credGithub "github.com/hashicorp/vault/builtin/credential/github"
	"github.com/hashicorp/vault/builtin/credential/userpass"
	credUserpass "github.com/hashicorp/vault/builtin/credential/userpass"
	"github.com/hashicorp/vault/helper/activationflags"
	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/storagepacker"
	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/hashicorp/vault/sdk/physical/inmem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func TestIdentityStore_DeleteEntityAlias(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	txn := c.identityStore.db.Txn(true)
	defer txn.Abort()

	alias := &identity.Alias{
		ID:             "testAliasID1",
		CanonicalID:    "testEntityID",
		MountType:      "testMountType",
		MountAccessor:  "testMountAccessor",
		Name:           "testAliasName",
		LocalBucketKey: c.identityStore.localAliasPacker.BucketKey("testEntityID"),
	}
	alias2 := &identity.Alias{
		ID:             "testAliasID2",
		CanonicalID:    "testEntityID",
		MountType:      "testMountType",
		MountAccessor:  "testMountAccessor2",
		Name:           "testAliasName2",
		LocalBucketKey: c.identityStore.localAliasPacker.BucketKey("testEntityID"),
	}
	entity := &identity.Entity{
		ID:       "testEntityID",
		Name:     "testEntityName",
		Policies: []string{"foo", "bar"},
		Aliases: []*identity.Alias{
			alias,
			alias2,
		},
		NamespaceID: namespace.RootNamespaceID,
		BucketKey:   c.identityStore.entityPacker.BucketKey("testEntityID"),
	}

	err := c.identityStore.upsertEntityInTxn(context.Background(), txn, entity, nil, false)
	require.NoError(t, err)

	err = c.identityStore.deleteAliasesInEntityInTxn(txn, entity, []*identity.Alias{alias, alias2})
	require.NoError(t, err)

	txn.Commit()

	alias, err = c.identityStore.MemDBAliasByID("testAliasID1", false, false)
	require.NoError(t, err)
	require.Nil(t, alias)

	alias, err = c.identityStore.MemDBAliasByID("testAliasID2", false, false)
	require.NoError(t, err)
	require.Nil(t, alias)

	entity, err = c.identityStore.MemDBEntityByID("testEntityID", false)
	require.NoError(t, err)

	require.Len(t, entity.Aliases, 0)
}

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
		ID:             "alias1",
		CanonicalID:    "entity1",
		MountType:      "github",
		MountAccessor:  meGH.Accessor,
		Name:           "githubuser",
		LocalBucketKey: c.identityStore.localAliasPacker.BucketKey("entity1"),
	}
	entity := &identity.Entity{
		ID:       "entity1",
		Name:     "name1",
		Policies: []string{"foo", "bar"},
		Aliases: []*identity.Alias{
			alias,
		},
		NamespaceID: namespace.RootNamespaceID,
		BucketKey:   c.identityStore.entityPacker.BucketKey("entity1"),
	}

	err = c.identityStore.upsertEntity(namespace.RootContext(nil), entity, nil, true)
	if err != nil {
		t.Fatal(err)
	}

	alias2 := &identity.Alias{
		ID:             "alias2",
		CanonicalID:    "entity2",
		MountType:      "github",
		MountAccessor:  meGH.Accessor,
		Name:           "GITHUBUSER",
		LocalBucketKey: c.identityStore.localAliasPacker.BucketKey("entity2"),
	}
	entity2 := &identity.Entity{
		ID:       "entity2",
		Name:     "name2",
		Policies: []string{"foo", "bar"},
		Aliases: []*identity.Alias{
			alias2,
		},
		NamespaceID: namespace.RootNamespaceID,
		BucketKey:   c.identityStore.entityPacker.BucketKey("entity2"),
	}

	// Persist the second entity directly without the regular flow. This will skip
	// merging of these enties.
	entity2Any, err := anypb.New(entity2)
	if err != nil {
		t.Fatal(err)
	}
	item := &storagepacker.Item{
		ID:      entity2.ID,
		Message: entity2Any,
	}

	ctx := namespace.RootContext(nil)
	if err = c.identityStore.entityPacker.PutItem(ctx, item); err != nil {
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
	entity, _, err := is.CreateOrFetchEntity(ctx, alias)
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
	is, ghAccessor, upAccessor, _ := testIdentityStoreWithGithubUserpassAuth(ctx, t)

	alias := &logical.Alias{
		MountType:     "github",
		MountAccessor: ghAccessor,
		Name:          "githubuser",
		Metadata: map[string]string{
			"foo": "a",
		},
	}

	entity, _, err := is.CreateOrFetchEntity(ctx, alias)
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

	entity, _, err = is.CreateOrFetchEntity(ctx, alias)
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
			"mount_accessor": upAccessor,
		},
	}

	resp, err := is.HandleRequest(ctx, registerReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	entity, _, err = is.CreateOrFetchEntity(ctx, alias)
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

	entity, _, err = is.CreateOrFetchEntity(ctx, alias)
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
		ID:             "alias1",
		CanonicalID:    "entity1",
		MountType:      "github",
		MountAccessor:  meGH.Accessor,
		Name:           "githubuser",
		LocalBucketKey: c.identityStore.localAliasPacker.BucketKey("entity1"),
	}
	entity := &identity.Entity{
		ID:       "entity1",
		Name:     "name1",
		Policies: []string{"foo", "bar"},
		Aliases: []*identity.Alias{
			alias,
		},
		NamespaceID: namespace.RootNamespaceID,
		BucketKey:   c.identityStore.entityPacker.BucketKey("entity1"),
	}
	err = c.identityStore.upsertEntity(namespace.RootContext(nil), entity, nil, true)
	if err != nil {
		t.Fatal(err)
	}

	alias2 := &identity.Alias{
		ID:             "alias2",
		CanonicalID:    "entity2",
		MountType:      "github",
		MountAccessor:  meGH.Accessor,
		Name:           "githubuser",
		LocalBucketKey: c.identityStore.localAliasPacker.BucketKey("entity2"),
	}
	entity2 := &identity.Entity{
		ID:       "entity2",
		Name:     "name2",
		Policies: []string{"bar", "baz"},
		Aliases: []*identity.Alias{
			alias2,
		},
		NamespaceID: namespace.RootNamespaceID,
		BucketKey:   c.identityStore.entityPacker.BucketKey("entity2"),
	}

	err = c.identityStore.upsertEntity(namespace.RootContext(nil), entity2, nil, true)
	if err != nil {
		t.Fatal(err)
	}

	newEntity, _, err := c.identityStore.CreateOrFetchEntity(namespace.RootContext(nil), &logical.Alias{
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

func testIdentityStoreWithGithubUserpassAuth(ctx context.Context, t *testing.T) (*IdentityStore, string, string, *Core) {
	// Setup 2 auth backends, github and userpass
	err := AddTestCredentialBackend("github", credGithub.Factory)
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

	return c.identityStore, githubMe.Accessor, userpassMe.Accessor, c
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

func expectSingleCount(t *testing.T, sink *metrics.InmemSink, keyPrefix string) {
	t.Helper()

	intervals := sink.Data()
	// Test crossed an interval boundary, don't try to deal with it.
	if len(intervals) > 1 {
		t.Skip("Detected interval crossing.")
	}

	var counter *metrics.SampledValue = nil

	for _, c := range intervals[0].Counters {
		if strings.HasPrefix(c.Name, keyPrefix) {
			counter = &c
			break
		}
	}
	if counter == nil {
		t.Fatalf("No %q counter found.", keyPrefix)
	}

	if counter.Count != 1 {
		t.Errorf("Counter number of samples %v is not 1.", counter.Count)
	}

	if counter.Sum != 1.0 {
		t.Errorf("Counter sum %v is not 1.", counter.Sum)
	}
}

func TestIdentityStore_NewEntityCounter(t *testing.T) {
	// Add github credential factory to core config
	err := AddTestCredentialBackend("github", credGithub.Factory)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	c, _, _, sink := TestCoreUnsealedWithMetrics(t)

	meGH := &MountEntry{
		Table:       credentialTableType,
		Path:        "github/",
		Type:        "github",
		Description: "github auth",
	}

	ctx := namespace.RootContext(nil)
	err = c.enableCredential(ctx, meGH)
	if err != nil {
		t.Fatal(err)
	}

	is := c.identityStore
	ghAccessor := meGH.Accessor

	alias := &logical.Alias{
		MountType:     "github",
		MountAccessor: ghAccessor,
		Name:          "githubuser",
		Metadata: map[string]string{
			"foo": "a",
		},
	}

	_, _, err = is.CreateOrFetchEntity(ctx, alias)
	if err != nil {
		t.Fatal(err)
	}

	expectSingleCount(t, sink, "identity.entity.creation")

	_, _, err = is.CreateOrFetchEntity(ctx, alias)
	if err != nil {
		t.Fatal(err)
	}

	expectSingleCount(t, sink, "identity.entity.creation")
}

func TestIdentityStore_UpdateAliasMetadataPerAccessor(t *testing.T) {
	entity := &identity.Entity{
		ID:       "testEntityID",
		Name:     "testEntityName",
		Policies: []string{"foo", "bar"},
		Aliases: []*identity.Alias{
			{
				ID:            "testAliasID1",
				CanonicalID:   "testEntityID",
				MountType:     "testMountType",
				MountAccessor: "testMountAccessor",
				Name:          "sameAliasName",
			},
			{
				ID:            "testAliasID2",
				CanonicalID:   "testEntityID",
				MountType:     "testMountType",
				MountAccessor: "testMountAccessor2",
				Name:          "sameAliasName",
			},
		},
		NamespaceID: namespace.RootNamespaceID,
	}

	login := &logical.Alias{
		MountType:     "testMountType",
		MountAccessor: "testMountAccessor",
		Name:          "sameAliasName",
		ID:            "testAliasID",
		Metadata:      map[string]string{"foo": "bar"},
	}

	if i := changedAliasIndex(entity, login); i != 0 {
		t.Fatalf("wrong alias index changed. Expected 0, got %d", i)
	}

	login2 := &logical.Alias{
		MountType:     "testMountType",
		MountAccessor: "testMountAccessor2",
		Name:          "sameAliasName",
		ID:            "testAliasID2",
		Metadata:      map[string]string{"bar": "foo"},
	}

	if i := changedAliasIndex(entity, login2); i != 1 {
		t.Fatalf("wrong alias index changed. Expected 1, got %d", i)
	}
}

// TestIdentityStore_DeleteCaseSensitivityKey tests that
// casesensitivity key gets removed from storage if it exists upon
// initializing identity store.
func TestIdentityStore_DeleteCaseSensitivityKey(t *testing.T) {
	c, unsealKey, root := TestCoreUnsealed(t)
	ctx := context.Background()

	// add caseSensitivityKey to storage
	entry, err := logical.StorageEntryJSON(caseSensitivityKey, &casesensitivity{
		DisableLowerCasedNames: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	err = c.identityStore.view.Put(ctx, entry)
	if err != nil {
		t.Fatal(err)
	}

	// check if the value is stored in storage
	storageEntry, err := c.identityStore.view.Get(ctx, caseSensitivityKey)
	if err != nil {
		t.Fatal(err)
	}

	if storageEntry == nil {
		t.Fatalf("bad: expected a non-nil entry for casesensitivity key")
	}

	// Seal and unseal to trigger identityStore initialize
	if err = c.Seal(root); err != nil {
		t.Fatal(err)
	}

	var unsealed bool
	for i := 0; i < len(unsealKey); i++ {
		unsealed, err = c.Unseal(unsealKey[i])
		if err != nil {
			t.Fatal(err)
		}
	}
	if !unsealed {
		t.Fatal("still sealed")
	}

	// check if caseSensitivityKey exists after initialize
	storageEntry, err = c.identityStore.view.Get(ctx, caseSensitivityKey)
	if err != nil {
		t.Fatal(err)
	}

	if storageEntry != nil {
		t.Fatalf("bad: expected no entry for casesensitivity key")
	}
}

// TestIdentityStoreInvalidate_Entities verifies the proper handling of
// entities in the Invalidate method.
func TestIdentityStoreInvalidate_Entities(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)

	// Create an entity in storage then call the Invalidate function
	//
	id, err := uuid.GenerateUUID()
	require.NoError(t, err)

	entity := &identity.Entity{
		Name:        "test",
		NamespaceID: namespace.RootNamespaceID,
		ID:          id,
		Aliases:     []*identity.Alias{},
		BucketKey:   c.identityStore.entityPacker.BucketKey(id),
	}

	p := c.identityStore.entityPacker

	// Persist the entity which we are merging to
	entityAsAny, err := anypb.New(entity)
	require.NoError(t, err)

	item := &storagepacker.Item{
		ID:      id,
		Message: entityAsAny,
	}

	err = p.PutItem(context.Background(), item)
	require.NoError(t, err)

	c.identityStore.Invalidate(context.Background(), p.BucketKey(id))

	txn := c.identityStore.db.Txn(true)

	memEntity, err := c.identityStore.MemDBEntityByIDInTxn(txn, id, true)
	assert.NoError(t, err)
	assert.NotNil(t, memEntity)

	txn.Commit()

	// Modify the entity in storage then call the Invalidate function
	entity.Metadata = make(map[string]string)
	entity.Metadata["foo"] = "bar"

	entityAsAny, err = anypb.New(entity)
	require.NoError(t, err)

	item.Message = entityAsAny

	p.PutItem(context.Background(), item)

	c.identityStore.Invalidate(context.Background(), p.BucketKey(id))

	txn = c.identityStore.db.Txn(true)

	memEntity, err = c.identityStore.MemDBEntityByIDInTxn(txn, id, true)
	assert.NoError(t, err)
	assert.Contains(t, memEntity.Metadata, "foo")

	txn.Commit()

	// Delete the entity in storage then call the Invalidate function
	err = p.DeleteItem(context.Background(), id)
	require.NoError(t, err)

	c.identityStore.Invalidate(context.Background(), p.BucketKey(id))

	txn = c.identityStore.db.Txn(true)

	memEntity, err = c.identityStore.MemDBEntityByIDInTxn(txn, id, true)
	assert.NoError(t, err)
	assert.Nil(t, memEntity)

	txn.Commit()
}

// TestIdentityStoreInvalidate_EntityAliasDelete verifies that the
// invalidateEntityBucket method properly cleans up aliases from
// MemDB that are no longer associated with the entity in the
// storage bucket.
func TestIdentityStoreInvalidate_EntityAliasDelete(t *testing.T) {
	ctx := namespace.ContextWithNamespace(context.Background(), namespace.RootNamespace)
	c, _, root := TestCoreUnsealed(t)

	// Enable a No-Op auth method
	c.credentialBackends["noop"] = func(context.Context, *logical.BackendConfig) (logical.Backend, error) {
		return &NoopBackend{
			BackendType: logical.TypeCredential,
		}, nil
	}
	mountAccessor1 := "noop-accessor1"
	mountAccessor2 := "noop-accessor2"
	mountAccessor3 := "noon-accessor3"

	createMountEntry := func(path, uuid, mountAccessor string, local bool) *MountEntry {
		return &MountEntry{
			Table:            credentialTableType,
			Path:             path,
			Type:             "noop",
			UUID:             uuid,
			Accessor:         mountAccessor,
			BackendAwareUUID: uuid + "backend",
			NamespaceID:      namespace.RootNamespaceID,
			namespace:        namespace.RootNamespace,
			Local:            local,
		}
	}

	c.auth = &MountTable{
		Type: credentialTableType,
		Entries: []*MountEntry{
			createMountEntry("/noop1", "abcd", mountAccessor1, false),
			createMountEntry("/noop2", "ghij", mountAccessor2, false),
			createMountEntry("/noop3", "mnop", mountAccessor3, true),
		},
	}

	require.NoError(t, c.setupCredentials(context.Background()))

	// Create an entity
	req := &logical.Request{
		ClientToken: root,
		Operation:   logical.UpdateOperation,
		Path:        "entity",
		Data: map[string]interface{}{
			"name": "alice",
		},
	}

	resp, err := c.identityStore.HandleRequest(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Contains(t, resp.Data, "id")

	entityID := resp.Data["id"].(string)

	createEntityAlias := func(name, mountAccessor string) string {
		req = &logical.Request{
			ClientToken: root,
			Operation:   logical.UpdateOperation,
			Path:        "entity-alias",
			Data: map[string]interface{}{
				"name":           name,
				"canonical_id":   entityID,
				"mount_accessor": mountAccessor,
			},
		}

		resp, err = c.identityStore.HandleRequest(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Contains(t, resp.Data, "id")

		return resp.Data["id"].(string)
	}

	alias1ID := createEntityAlias("alias1", mountAccessor1)
	alias2ID := createEntityAlias("alias2", mountAccessor2)
	alias3ID := createEntityAlias("alias3", mountAccessor3)

	// Update the entity in storage only to remove alias2 then call invalidate
	bucketKey := c.identityStore.entityPacker.BucketKey(entityID)
	bucket, err := c.identityStore.entityPacker.GetBucket(context.Background(), bucketKey)
	require.NoError(t, err)
	require.NotNil(t, bucket)

	bucketEntityItem := bucket.Items[0] // since there's only 1 entity
	bucketEntity, err := c.identityStore.parseEntityFromBucketItem(context.Background(), bucketEntityItem)
	require.NoError(t, err)
	require.NotNil(t, bucketEntity)

	replacementAliases := make([]*identity.Alias, 1)
	for _, a := range bucketEntity.Aliases {
		if a.ID != alias2ID {
			replacementAliases[0] = a
			break
		}
	}

	bucketEntity.Aliases = replacementAliases

	bucketEntityItem.Message, err = anypb.New(bucketEntity)
	require.NoError(t, err)

	require.NoError(t, c.identityStore.entityPacker.PutItem(context.Background(), bucketEntityItem))

	c.identityStore.Invalidate(context.Background(), bucketKey)

	alias1, err := c.identityStore.MemDBAliasByID(alias1ID, false, false)
	assert.NoError(t, err)
	assert.NotNil(t, alias1)

	alias2, err := c.identityStore.MemDBAliasByID(alias2ID, false, false)
	assert.NoError(t, err)
	assert.Nil(t, alias2)

	alias3, err := c.identityStore.MemDBAliasByID(alias3ID, false, false)
	assert.NoError(t, err)
	assert.NotNil(t, alias3)
}

// TestIdentityStoreInvalidate_EntityLocalAliasDelete verifies that the
// invalidateLocalAliasesBucket method properly cleans up aliases from
// MemDB that are no longer associated with the entity in the
// storage bucket.
func TestIdentityStoreInvalidate_EntityLocalAliasDelete(t *testing.T) {
	ctx := namespace.ContextWithNamespace(context.Background(), namespace.RootNamespace)
	c, _, root := TestCoreUnsealed(t)

	// Enable a No-Op auth method
	c.credentialBackends["noop"] = func(context.Context, *logical.BackendConfig) (logical.Backend, error) {
		return &NoopBackend{
			BackendType: logical.TypeCredential,
		}, nil
	}
	mountAccessor1 := "noop-accessor1"
	mountAccessor2 := "noop-accessor2"
	mountAccessor3 := "noon-accessor3"

	createMountEntry := func(path, uuid, mountAccessor string, local bool) *MountEntry {
		return &MountEntry{
			Table:            credentialTableType,
			Path:             path,
			Type:             "noop",
			UUID:             uuid,
			Accessor:         mountAccessor,
			BackendAwareUUID: uuid + "backend",
			NamespaceID:      namespace.RootNamespaceID,
			namespace:        namespace.RootNamespace,
			Local:            local,
		}
	}

	c.auth = &MountTable{
		Type: credentialTableType,
		Entries: []*MountEntry{
			createMountEntry("/noop1", "abcd", mountAccessor1, true),
			createMountEntry("/noop2", "ghij", mountAccessor2, true),
			createMountEntry("/noop3", "mnop", mountAccessor3, true),
		},
	}

	require.NoError(t, c.setupCredentials(context.Background()))

	// Create an entity
	req := &logical.Request{
		ClientToken: root,
		Operation:   logical.UpdateOperation,
		Path:        "entity",
		Data: map[string]interface{}{
			"name": "alice",
		},
	}

	resp, err := c.identityStore.HandleRequest(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Contains(t, resp.Data, "id")

	entityID := resp.Data["id"].(string)

	createEntityAlias := func(name, mountAccessor string) string {
		req = &logical.Request{
			ClientToken: root,
			Operation:   logical.UpdateOperation,
			Path:        "entity-alias",
			Data: map[string]interface{}{
				"name":           name,
				"canonical_id":   entityID,
				"mount_accessor": mountAccessor,
			},
		}

		resp, err = c.identityStore.HandleRequest(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Contains(t, resp.Data, "id")

		return resp.Data["id"].(string)
	}

	alias1ID := createEntityAlias("alias1", mountAccessor1)
	alias2ID := createEntityAlias("alias2", mountAccessor2)
	alias3ID := createEntityAlias("alias3", mountAccessor3)

	for i, aliasID := range []string{alias1ID, alias2ID, alias3ID} {
		alias, err := c.identityStore.MemDBAliasByID(aliasID, false, false)
		require.NoError(t, err, i)
		require.NotNil(t, alias, i)
	}

	// // Update the entity in storage only to remove alias2 then call invalidate
	bucketKey := c.identityStore.entityPacker.BucketKey(entityID)
	bucket, err := c.identityStore.entityPacker.GetBucket(context.Background(), bucketKey)
	require.NoError(t, err)
	require.NotNil(t, bucket)

	bucketEntityItem := bucket.Items[0] // since there's only 1 entity
	bucketEntity, err := c.identityStore.parseEntityFromBucketItem(context.Background(), bucketEntityItem)
	require.NoError(t, err)
	require.NotNil(t, bucketEntity)

	bucketKey = c.identityStore.localAliasPacker.BucketKey(entityID)
	bucketLocalAlias, err := c.identityStore.localAliasPacker.GetBucket(context.Background(), bucketKey)
	require.NoError(t, err)
	require.NotNil(t, bucketLocalAlias)

	bucketLocalAliasItem := bucketLocalAlias.Items[0]
	require.Equal(t, entityID, bucketLocalAliasItem.ID)

	var localAliases identity.LocalAliases

	err = anypb.UnmarshalTo(bucketLocalAliasItem.Message, &localAliases, proto.UnmarshalOptions{})
	require.NoError(t, err)

	memDBEntity, err := c.identityStore.MemDBEntityByID(entityID, false)
	require.NoError(t, err)
	require.NotNil(t, memDBEntity)

	replacementAliases := make([]*identity.Alias, 0)
	for _, a := range memDBEntity.Aliases {
		if a.ID != alias2ID {
			replacementAliases = append(replacementAliases, a)
		}
	}

	localAliases.Aliases = replacementAliases

	bucketLocalAliasItem.Message, err = anypb.New(&localAliases)
	require.NoError(t, err)

	require.NoError(t, c.identityStore.localAliasPacker.PutItem(context.Background(), bucketLocalAliasItem))

	c.identityStore.Invalidate(context.Background(), bucketKey)

	alias1, err := c.identityStore.MemDBAliasByID(alias1ID, false, false)
	assert.NoError(t, err)
	assert.NotNil(t, alias1)

	alias2, err := c.identityStore.MemDBAliasByID(alias2ID, false, false)
	assert.NoError(t, err)
	assert.Nil(t, alias2)

	alias3, err := c.identityStore.MemDBAliasByID(alias3ID, false, false)
	assert.NoError(t, err)
	assert.NotNil(t, alias3)
}

// TestIdentityStoreInvalidate_LocalAliasesWithEntity verifies the correct
// handling of local aliases in the Invalidate method.
func TestIdentityStoreInvalidate_LocalAliasesWithEntity(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)

	// Create an entity in storage then call the Invalidate function
	//
	entityID, err := uuid.GenerateUUID()
	require.NoError(t, err)

	entity := &identity.Entity{
		Name:        "test",
		NamespaceID: namespace.RootNamespaceID,
		ID:          entityID,
		Aliases:     []*identity.Alias{},
		BucketKey:   c.identityStore.entityPacker.BucketKey(entityID),
	}

	aliasID, err := uuid.GenerateUUID()
	require.NoError(t, err)

	localAliases := &identity.LocalAliases{
		Aliases: []*identity.Alias{
			{
				ID:            aliasID,
				Name:          "test",
				NamespaceID:   namespace.RootNamespaceID,
				CanonicalID:   entityID,
				MountAccessor: "userpass-000000",
			},
		},
	}

	ep := c.identityStore.entityPacker

	// Persist the entity which we are merging to
	entityAsAny, err := anypb.New(entity)
	require.NoError(t, err)

	entityItem := &storagepacker.Item{
		ID:      entityID,
		Message: entityAsAny,
	}

	err = ep.PutItem(context.Background(), entityItem)
	require.NoError(t, err)

	c.identityStore.Invalidate(context.Background(), ep.BucketKey(entityID))

	lap := c.identityStore.localAliasPacker

	localAliasesAsAny, err := anypb.New(localAliases)
	require.NoError(t, err)

	localAliasesItem := &storagepacker.Item{
		ID:      entityID,
		Message: localAliasesAsAny,
	}

	err = lap.PutItem(context.Background(), localAliasesItem)
	require.NoError(t, err)

	c.identityStore.Invalidate(context.Background(), lap.BucketKey(entityID))

	txn := c.identityStore.db.Txn(true)

	memDBEntity, err := c.identityStore.MemDBEntityByIDInTxn(txn, entityID, true)
	assert.NoError(t, err)
	assert.NotNil(t, memDBEntity)

	memDBLocalAlias, err := c.identityStore.MemDBAliasByIDInTxn(txn, aliasID, true, false)
	assert.NoError(t, err)
	assert.NotNil(t, memDBLocalAlias)
	assert.Equal(t, 1, len(memDBEntity.Aliases))
	assert.NotNil(t, memDBEntity.Aliases[0])
	assert.Equal(t, memDBEntity.Aliases[0].ID, memDBLocalAlias.ID)

	txn.Commit()
}

// TestIdentityStoreInvalidate_TemporaryEntity verifies the proper handling of
// temporary entities in the Invalidate method.
func TestIdentityStoreInvalidate_TemporaryEntity(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)

	// Create an entity in storage then call the Invalidate function
	//
	entityID, err := uuid.GenerateUUID()
	require.NoError(t, err)

	tempEntity := &identity.Entity{
		Name:        "test",
		NamespaceID: namespace.RootNamespaceID,
		ID:          entityID,
		Aliases:     []*identity.Alias{},
		BucketKey:   c.identityStore.entityPacker.BucketKey(entityID),
	}

	lap := c.identityStore.localAliasPacker
	ep := c.identityStore.entityPacker

	// Persist the entity which we are merging to
	tempEntityAsAny, err := anypb.New(tempEntity)
	require.NoError(t, err)

	tempEntityItem := &storagepacker.Item{
		ID:      entityID + tmpSuffix,
		Message: tempEntityAsAny,
	}

	err = lap.PutItem(context.Background(), tempEntityItem)
	require.NoError(t, err)

	entityAsAny := tempEntityAsAny

	entityItem := &storagepacker.Item{
		ID:      entityID,
		Message: entityAsAny,
	}

	err = ep.PutItem(context.Background(), entityItem)
	require.NoError(t, err)

	c.identityStore.Invalidate(context.Background(), ep.BucketKey(entityID))

	txn := c.identityStore.db.Txn(true)

	memDBEntity, err := c.identityStore.MemDBEntityByIDInTxn(txn, entityID, true)
	assert.NoError(t, err)
	assert.NotNil(t, memDBEntity)

	item, err := lap.GetItem(lap.BucketKey(entityID) + tmpSuffix)
	assert.NoError(t, err)
	assert.Nil(t, item)
}

// TestIdentityStoreLoadingIsDeterministic tests the default error resolver and
// the identity cleanup rename resolver to ensure that loading is deterministic
// for both.
func TestIdentityStoreLoadingIsDeterministic(t *testing.T) {
	t.Run(t.Name()+"error-resolver", func(t *testing.T) {
		identityStoreLoadingIsDeterministic(t, false)
	})
	t.Run(t.Name()+"identity-cleanup", func(t *testing.T) {
		identityStoreLoadingIsDeterministic(t, true)
	})
}

// identityStoreLoadingIsDeterministic is a property-based test helper that
// ensures the loading logic of the entity store is deterministic. This is
// important because we perform certain merges and corrections of duplicates on
// load and non-deterministic order can cause divergence between different nodes
// or even after seal/unseal cycles on one node. Loading _should_ be
// deterministic anyway if all data in storage was correct see comments inline
// for examples of ways storage can be corrupt with respect to the expected
// schema invariants.
func identityStoreLoadingIsDeterministic(t *testing.T, identityDeduplication bool) {
	// Create some state in store that could trigger non-deterministic behavior.
	// The nature of the identity store schema is such that the order of loading
	// entities etc shouldn't matter even if it was non-deterministic, however due
	// to many and varied historical (and possibly current/future) bugs, we have
	// seen many cases where storage ends up with duplicates persisted. This is
	// not ideal of course and our code attempts to "fix" on the fly with merges
	// on load. But it's hampered by the fact that the current implementation does
	// not load entities in a deterministic order. which means that different
	// nodes potentially resolve merges differently. This test proves that that
	// happens and should hopefully provide some confidence that we don't
	// introduce non-determinism in the future somehow. It's a bit odd we have to
	// inject essentially invalid data into storage to trigger the issue but
	// that's what we get in real life sometimes!
	logger := corehelpers.NewTestLogger(t)
	ims, err := inmem.NewTransactionalInmemHA(nil, logger)
	require.NoError(t, err)

	cfg := &CoreConfig{
		Physical:        ims,
		HAPhysical:      ims.(physical.HABackend),
		Logger:          logger,
		BuiltinRegistry: corehelpers.NewMockBuiltinRegistry(),
		CredentialBackends: map[string]logical.Factory{
			"userpass": userpass.Factory,
		},
	}

	c, sealKeys, rootToken := TestCoreUnsealedWithConfig(t, cfg)

	// Inject values into storage
	upme, err := TestUserpassMount(c, false)
	require.NoError(t, err)
	localMe, err := TestUserpassMount(c, true)
	require.NoError(t, err)

	ctx := context.Background()

	if identityDeduplication {
		err = c.FeatureActivationFlags.Write(ctx, activationflags.IdentityDeduplication, true)
		require.NoError(t, err)
	}

	// We create 100 entities each with 1 non-local alias and 1 local alias. We
	// then randomly create duplicate alias or local alias entries with a
	// probability that is unrealistic but ensures we have duplicates on every
	// test run with high probability and more than 1 duplicate often.
	for i := 0; i <= 100; i++ {
		name := fmt.Sprintf("entity-%d", i)
		alias := fmt.Sprintf("alias-%d", i)
		localAlias := fmt.Sprintf("localalias-%d", i)
		e := makeEntityForPacker(t, name, c.identityStore.entityPacker)
		attachAlias(t, e, alias, upme)
		attachAlias(t, e, localAlias, localMe)
		err = TestHelperWriteToStoragePacker(ctx, c.identityStore.entityPacker, e.ID, e)
		require.NoError(t, err)

		// Subset of entities get a duplicate alias and/or duplicate local alias.
		// We'll use a probability of 0.3 for each dup so that we expect at least a
		// few double and maybe triple duplicates of each type every few test runs
		// and may have duplicates of both types or neither etc.
		pDup := 0.3
		rnd := rand.Float64()
		dupeNum := 1
		for rnd < pDup && dupeNum < 10 {
			e := makeEntityForPacker(t, fmt.Sprintf("entity-%d-dup-%d", i, dupeNum), c.identityStore.entityPacker)
			attachAlias(t, e, alias, upme)
			err = TestHelperWriteToStoragePacker(ctx, c.identityStore.entityPacker, e.ID, e)
			require.NoError(t, err)
			// Toss again to see if we continue
			rnd = rand.Float64()
			dupeNum++
		}
		// Toss the coin again to see if there are any local dupes
		dupeNum = 1
		rnd = rand.Float64()
		for rnd < pDup && dupeNum < 10 {
			e := makeEntityForPacker(t, fmt.Sprintf("entity-%d-localdup-%d", i, dupeNum), c.identityStore.entityPacker)
			attachAlias(t, e, localAlias, localMe)
			err = TestHelperWriteToStoragePacker(ctx, c.identityStore.entityPacker, e.ID, e)
			require.NoError(t, err)
			rnd = rand.Float64()
			dupeNum++
		}
		// See if we should add entity _name_ duplicates too (with no aliases)
		rnd = rand.Float64()
		for rnd < pDup {
			e := makeEntityForPacker(t, name, c.identityStore.entityPacker)
			err = TestHelperWriteToStoragePacker(ctx, c.identityStore.entityPacker, e.ID, e)
			require.NoError(t, err)
			rnd = rand.Float64()
		}
		// One more edge case is that it's currently possible as of the time of
		// writing for a failure during entity invalidation to result in a permanent
		// "cached" entity in the local alias packer even though we do have the
		// replicated entity in the entity packer too. This is a bug and will
		// hopefully be fixed at some point soon, but even after it is it's
		// important that we still test for it since existing clusters may still
		// have this persistent state. Pick a low probability but one we're very
		// likely to hit in 100 iterations and write the entity to the local alias
		// table too (this mimics the behavior of cacheTemporaryEntity).
		pFailedLocalAliasInvalidation := 0.02
		if rand.Float64() < pFailedLocalAliasInvalidation {
			err = TestHelperWriteToStoragePacker(ctx, c.identityStore.localAliasPacker, e.ID+tmpSuffix, e)
			require.NoError(t, err)
		}
	}

	// Create some groups
	for i := 0; i <= 100; i++ {
		name := fmt.Sprintf("group-%d", i)
		// Add an alias to every other group
		alias := ""
		if i%2 == 0 {
			alias = fmt.Sprintf("groupalias-%d", i)
		}
		e := makeGroupWithNameAndAlias(t, name, alias, c.identityStore.groupPacker, upme)
		err = TestHelperWriteToStoragePacker(ctx, c.identityStore.groupPacker, e.ID, e)
		require.NoError(t, err)
	}
	// Now add 10 groups with the same alias to ensure duplicates don't cause
	// non-deterministic behavior.
	for i := 0; i <= 10; i++ {
		name := fmt.Sprintf("group-dup-%d", i)
		e := makeGroupWithNameAndAlias(t, name, "groupalias-dup", c.identityStore.groupPacker, upme)
		err = TestHelperWriteToStoragePacker(ctx, c.identityStore.groupPacker, e.ID, e)
		require.NoError(t, err)
	}
	// Add a second and third groups with duplicate names too.
	for _, name := range []string{"group-0", "group-1", "group-1"} {
		e := makeGroupWithNameAndAlias(t, name, "", c.identityStore.groupPacker, upme)
		err = TestHelperWriteToStoragePacker(ctx, c.identityStore.groupPacker, e.ID, e)
		require.NoError(t, err)
	}

	entIdentityStoreDeterminismTestSetup(t, ctx, c, upme, localMe)

	// Storage is now primed for the test.

	// To test that this is deterministic we need to load from storage a bunch of
	// times and make sure we get the same result. For easier debugging we'll
	// build a list of human readable ids that we can compare.
	prevLoadedNames := []string{}
	for i := 0; i < 10; i++ {
		// Seal and unseal to reload the identity store
		require.NoError(t, c.Seal(rootToken))
		require.True(t, c.Sealed())
		for _, key := range sealKeys {
			unsealed, err := c.Unseal(key)
			require.NoError(t, err, "failed unseal on attempt %d", i)
			if unsealed {
				break
			}
		}
		require.False(t, c.Sealed())

		// Identity store should be loaded now. Check it's contents.
		loadedNames := []string{}

		tx := c.identityStore.db.Txn(false)

		// Entities + their aliases
		iter, err := tx.LowerBound(entitiesTable, "id", "")
		require.NoError(t, err)
		for item := iter.Next(); item != nil; item = iter.Next() {
			e := item.(*identity.Entity)
			loadedNames = append(loadedNames, e.Name)
			for _, a := range e.Aliases {
				loadedNames = append(loadedNames, a.Name)
			}
		}
		// This is a non-triviality check to make sure we actually loaded stuff and
		// are not just passing because of a bug in the test.
		numLoaded := len(loadedNames)
		require.Greater(t, numLoaded, 300, "not enough entities and aliases loaded on attempt %d", i)

		// Standalone alias query
		iter, err = tx.LowerBound(entityAliasesTable, "id", "")
		require.NoError(t, err)
		for item := iter.Next(); item != nil; item = iter.Next() {
			a := item.(*identity.Alias)
			loadedNames = append(loadedNames, a.Name)
		}

		// Groups
		iter, err = tx.LowerBound(groupsTable, "id", "")
		require.NoError(t, err)
		for item := iter.Next(); item != nil; item = iter.Next() {
			g := item.(*identity.Group)
			loadedNames = append(loadedNames, g.Name)
			if g.Alias != nil {
				loadedNames = append(loadedNames, g.Alias.Name)
			}
		}
		// This is a non-triviality check to make sure we actually loaded stuff and
		// are not just passing because of a bug in the test.
		groupsLoaded := len(loadedNames) - numLoaded
		require.Greater(t, groupsLoaded, 140, "not enough groups and aliases loaded on attempt %d", i)

		// note `lastIDs` argument is not needed any more but we can't change the
		// signature without breaking enterprise. It's simpler to keep it unused for
		// now until both parts of this merge.
		entIdentityStoreDeterminismAssert(t, i, loadedNames, nil)

		if i > 0 {
			// Should be in the same order if we are deterministic since MemDB has strong ordering.
			require.Equal(t, prevLoadedNames, loadedNames, "different result on attempt %d", i)
		}

		prevLoadedNames = loadedNames
	}
}

// TestIdentityStoreLoadingDuplicateReporting tests the reporting of different
// types of duplicates during unseal when in case-sensitive mode.
func TestIdentityStoreLoadingDuplicateReporting(t *testing.T) {
	logger := corehelpers.NewTestLogger(t)
	ims, err := inmem.NewTransactionalInmemHA(nil, logger)
	require.NoError(t, err)

	cfg := &CoreConfig{
		Physical:        ims,
		HAPhysical:      ims.(physical.HABackend),
		Logger:          logger,
		BuiltinRegistry: corehelpers.NewMockBuiltinRegistry(),
		CredentialBackends: map[string]logical.Factory{
			"userpass": userpass.Factory,
		},
	}

	c, sealKeys, rootToken := TestCoreUnsealedWithConfig(t, cfg)

	// Inject values into storage
	upme, err := TestUserpassMount(c, false)
	require.NoError(t, err)
	localMe, err := TestUserpassMount(c, true)
	require.NoError(t, err)

	ctx := namespace.RootContext(nil)

	identityCreateCaseDuplicates(t, ctx, c, upme, localMe)

	entIdentityStoreDuplicateReportTestSetup(t, ctx, c, rootToken)

	// Storage is now primed for the test.

	// Seal and unseal to reload the identity store
	require.NoError(t, c.Seal(rootToken))
	require.True(t, c.Sealed())

	// Setup a logger we can use to capture unseal logs
	var unsealLogs []string
	unsealLogger := &logFn{
		fn: func(msg string, args []interface{}) {
			pairs := make([]string, 0, len(args)/2)
			for pair := range slices.Chunk(args, 2) {
				// Yes this will panic if we didn't log an even number of args but thats
				// OK because that's a bug!
				pairs = append(pairs, fmt.Sprintf("%s=%s", pair[0], pair[1]))
			}
			unsealLogs = append(unsealLogs, fmt.Sprintf("%s: %s", msg, strings.Join(pairs, " ")))
		},
	}
	logger.RegisterSink(unsealLogger)

	for _, key := range sealKeys {
		unsealed, err := c.Unseal(key)
		require.NoError(t, err)
		if unsealed {
			break
		}
	}
	require.False(t, c.Sealed())
	logger.DeregisterSink(unsealLogger)

	// Identity store should be loaded now. Check it's contents.

	// We don't expect any actual behavior change just logs reporting duplicates.
	// We could assert the current "expected" behavior but it's actually broken in
	// many of these cases and seems strange to encode in a test that we want
	// broken behavior!
	numDupes := make(map[string]int)
	duplicateCountRe := regexp.MustCompile(`(\d+) (different-case( local)? entity alias|entity|group) duplicates found`)
	for _, log := range unsealLogs {
		if matches := duplicateCountRe.FindStringSubmatch(log); len(matches) >= 3 {
			num, _ := strconv.Atoi(matches[1])
			numDupes[matches[2]] = num
		}
	}
	t.Logf("numDupes: %v", numDupes)
	wantAliases, wantLocalAliases, wantEntities, wantGroups := identityStoreDuplicateReportTestWantDuplicateCounts()
	require.Equal(t, wantLocalAliases, numDupes["different-case local entity alias"])
	require.Equal(t, wantAliases, numDupes["different-case entity alias"])
	require.Equal(t, wantEntities, numDupes["entity"])
	require.Equal(t, wantGroups, numDupes["group"])
}

type logFn struct {
	fn func(msg string, args []interface{})
}

// Accept implements hclog.SinkAdapter
func (f *logFn) Accept(name string, level hclog.Level, msg string, args ...interface{}) {
	f.fn(msg, args)
}
