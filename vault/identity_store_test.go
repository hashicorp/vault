package vault

import (
	"testing"
	"time"

	credGithub "github.com/hashicorp/vault/builtin/credential/github"
	"github.com/hashicorp/vault/logical"
)

func TestIdentityStore_CreateEntity(t *testing.T) {
	is, ghAccessor, _ := testIdentityStoreWithGithubAuth(t)
	alias := &logical.Alias{
		MountType:     "github",
		MountAccessor: ghAccessor,
		Name:          "githubuser",
	}

	entity, err := is.CreateEntity(alias)
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

	// Try recreating an entity with the same alias details. It should fail.
	entity, err = is.CreateEntity(alias)
	if err == nil {
		t.Fatalf("expected an error")
	}
}

func TestIdentityStore_EntityByAliasFactors(t *testing.T) {
	var err error
	var resp *logical.Response

	is, ghAccessor, _ := testIdentityStoreWithGithubAuth(t)

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
	resp, err = is.HandleRequest(registerReq)
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

	resp, err = is.HandleRequest(aliasReq)
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

	core, is, ts, _ := testCoreWithIdentityTokenGithub(t)

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
	resp, err = is.HandleRequest(registerReq)
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
	te := &TokenEntry{
		Path:     "test",
		Policies: []string{"default", responseWrappingPolicyName},
		EntityID: entityID,
	}

	if err := ts.create(te); err != nil {
		t.Fatal(err)
	}

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

	resp, err = core.HandleRequest(wrapReq)
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
	_, ts, _, _ := TestCoreWithTokenStore(t)

	// Create a token which has EntityID set
	te := &TokenEntry{
		Path:     "test",
		Policies: []string{"dev", "prod"},
		EntityID: "testentityid",
	}

	if err := ts.create(te); err != nil {
		t.Fatal(err)
	}

	// Create a child token; this should inherit the EntityID
	tokenReq := &logical.Request{
		Operation:   logical.UpdateOperation,
		Path:        "create",
		ClientToken: te.ID,
	}

	resp, err := ts.HandleRequest(tokenReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v err: %v", err, resp)
	}

	if resp.Auth.EntityID != te.EntityID {
		t.Fatalf("bad: entity ID; expected: %v, actual: %v", te.EntityID, resp.Auth.EntityID)
	}

	// Create an orphan token; this should not inherit the EntityID
	tokenReq.Path = "create-orphan"
	resp, err = ts.HandleRequest(tokenReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v err: %v", err, resp)
	}

	if resp.Auth.EntityID != "" {
		t.Fatalf("expected entity ID to be not set")
	}
}

func testCoreWithIdentityTokenGithub(t *testing.T) (*Core, *IdentityStore, *TokenStore, string) {
	is, ghAccessor, core := testIdentityStoreWithGithubAuth(t)
	ts := testTokenStore(t, core)
	return core, is, ts, ghAccessor
}

func testCoreWithIdentityTokenGithubRoot(t *testing.T) (*Core, *IdentityStore, *TokenStore, string, string) {
	is, ghAccessor, core, root := testIdentityStoreWithGithubAuthRoot(t)
	ts := testTokenStore(t, core)
	return core, is, ts, ghAccessor, root
}

func testIdentityStoreWithGithubAuth(t *testing.T) (*IdentityStore, string, *Core) {
	is, ghA, c, _ := testIdentityStoreWithGithubAuthRoot(t)
	return is, ghA, c
}

// testIdentityStoreWithGithubAuth returns an instance of identity store which
// is mounted by default. This function also enables the github auth backend to
// assist with testing aliases and entities that require an valid mount
// accessor of an auth backend.
func testIdentityStoreWithGithubAuthRoot(t *testing.T) (*IdentityStore, string, *Core, string) {
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

	err = c.enableCredential(meGH)
	if err != nil {
		t.Fatal(err)
	}

	// Identity store will be mounted by now, just fetch it from router
	identitystore := c.router.MatchingBackend("identity/")
	if identitystore == nil {
		t.Fatalf("failed to fetch identity store from router")
	}

	return identitystore.(*IdentityStore), meGH.Accessor, c, root
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
