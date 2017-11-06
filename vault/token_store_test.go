package vault

import (
	"encoding/json"
	"fmt"
	"path"
	"reflect"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/locksutil"
	"github.com/hashicorp/vault/logical"
)

type TokenEntryOld struct {
	ID             string
	Accessor       string
	Parent         string
	Policies       []string
	Path           string
	Meta           map[string]string
	DisplayName    string
	NumUses        int
	CreationTime   int64
	TTL            time.Duration
	ExplicitMaxTTL time.Duration
	Role           string
	Period         time.Duration
}

func TestTokenStore_TokenEntryUpgrade(t *testing.T) {
	var err error
	_, ts, _, _ := TestCoreWithTokenStore(t)

	// Use a struct that does not have struct tags to store the items and
	// check if the lookup code handles them properly while reading back
	entry := &TokenEntryOld{
		DisplayName:    "test-display-name",
		Path:           "test",
		Policies:       []string{"dev", "ops"},
		CreationTime:   time.Now().Unix(),
		ExplicitMaxTTL: 100,
		NumUses:        10,
	}
	entry.ID, err = uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}

	enc, err := json.Marshal(entry)
	if err != nil {
		t.Fatal(err)
	}

	saltedId, err := ts.SaltID(entry.ID)
	if err != nil {
		t.Fatal(err)
	}
	path := lookupPrefix + saltedId
	le := &logical.StorageEntry{
		Key:   path,
		Value: enc,
	}

	if err := ts.view.Put(le); err != nil {
		t.Fatal(err)
	}

	out, err := ts.Lookup(entry.ID)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if out.DisplayName != "test-display-name" {
		t.Fatalf("bad: display_name: expected: test-display-name, actual: %s", out.DisplayName)
	}
	if out.CreationTime == 0 {
		t.Fatal("bad: expected a non-zero creation time")
	}
	if out.ExplicitMaxTTL != 100 {
		t.Fatalf("bad: explicit_max_ttl: expected: 100, actual: %d", out.ExplicitMaxTTL)
	}
	if out.NumUses != 10 {
		t.Fatalf("bad: num_uses: expected: 10, actual: %d", out.NumUses)
	}

	// Test the default case to ensure there are no regressions
	ent := &TokenEntry{
		DisplayName:    "test-display-name",
		Path:           "test",
		Policies:       []string{"dev", "ops"},
		CreationTime:   time.Now().Unix(),
		ExplicitMaxTTL: 100,
		NumUses:        10,
	}
	if err := ts.create(ent); err != nil {
		t.Fatalf("err: %s", err)
	}

	out, err = ts.Lookup(ent.ID)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if out.DisplayName != "test-display-name" {
		t.Fatalf("bad: display_name: expected: test-display-name, actual: %s", out.DisplayName)
	}
	if out.CreationTime == 0 {
		t.Fatal("bad: expected a non-zero creation time")
	}
	if out.ExplicitMaxTTL != 100 {
		t.Fatalf("bad: explicit_max_ttl: expected: 100, actual: %d", out.ExplicitMaxTTL)
	}
	if out.NumUses != 10 {
		t.Fatalf("bad: num_uses: expected: 10, actual: %d", out.NumUses)
	}

	// Fill in the deprecated fields and read out from proper fields
	ent = &TokenEntry{
		Path:                     "test",
		Policies:                 []string{"dev", "ops"},
		DisplayNameDeprecated:    "test-display-name",
		CreationTimeDeprecated:   time.Now().Unix(),
		ExplicitMaxTTLDeprecated: 100,
		NumUsesDeprecated:        10,
	}
	if err := ts.create(ent); err != nil {
		t.Fatalf("err: %s", err)
	}

	out, err = ts.Lookup(ent.ID)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if out.DisplayName != "test-display-name" {
		t.Fatalf("bad: display_name: expected: test-display-name, actual: %s", out.DisplayName)
	}
	if out.CreationTime == 0 {
		t.Fatal("bad: expected a non-zero creation time")
	}
	if out.ExplicitMaxTTL != 100 {
		t.Fatalf("bad: explicit_max_ttl: expected: 100, actual: %d", out.ExplicitMaxTTL)
	}
	if out.NumUses != 10 {
		t.Fatalf("bad: num_uses: expected: 10, actual: %d", out.NumUses)
	}

	// Check if NumUses picks up a lower value
	ent = &TokenEntry{
		Path:              "test",
		NumUses:           5,
		NumUsesDeprecated: 10,
	}
	if err := ts.create(ent); err != nil {
		t.Fatalf("err: %s", err)
	}

	out, err = ts.Lookup(ent.ID)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if out.NumUses != 5 {
		t.Fatalf("bad: num_uses: expected: 5, actual: %d", out.NumUses)
	}

	// Switch the values from deprecated and proper field and check if the
	// lower value is still getting picked up
	ent = &TokenEntry{
		Path:              "test",
		NumUses:           10,
		NumUsesDeprecated: 5,
	}
	if err := ts.create(ent); err != nil {
		t.Fatalf("err: %s", err)
	}

	out, err = ts.Lookup(ent.ID)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if out.NumUses != 5 {
		t.Fatalf("bad: num_uses: expected: 5, actual: %d", out.NumUses)
	}
}

func getBackendConfig(c *Core) *logical.BackendConfig {
	return &logical.BackendConfig{
		Logger: c.logger,
		System: logical.StaticSystemView{
			DefaultLeaseTTLVal: time.Hour * 24,
			MaxLeaseTTLVal:     time.Hour * 24 * 32,
		},
	}
}

func testMakeToken(t *testing.T, ts *TokenStore, root, client, ttl string, policy []string) {
	req := logical.TestRequest(t, logical.UpdateOperation, "create")
	req.ClientToken = root
	req.Data["id"] = client
	req.Data["policies"] = policy
	req.Data["ttl"] = ttl

	resp, err := ts.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp.Auth.ClientToken != client {
		t.Fatalf("bad: %#v", resp)
	}
}

func testCoreMakeToken(t *testing.T, c *Core, root, client, ttl string, policy []string) {
	req := logical.TestRequest(t, logical.UpdateOperation, "auth/token/create")
	req.ClientToken = root
	req.Data["id"] = client
	req.Data["policies"] = policy
	req.Data["ttl"] = ttl

	resp, err := c.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp.IsError() {
		t.Fatalf("err: %v %v", err, *resp)
	}
	if resp.Auth.ClientToken != client {
		t.Fatalf("bad: %#v", *resp)
	}
}

func TestTokenStore_AccessorIndex(t *testing.T) {
	_, ts, _, _ := TestCoreWithTokenStore(t)

	ent := &TokenEntry{Path: "test", Policies: []string{"dev", "ops"}}
	if err := ts.create(ent); err != nil {
		t.Fatalf("err: %s", err)
	}

	out, err := ts.Lookup(ent.ID)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Ensure that accessor is created
	if out == nil || out.Accessor == "" {
		t.Fatalf("bad: %#v", out)
	}

	aEntry, err := ts.lookupByAccessor(out.Accessor, false)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Verify that the value returned from the index matches the token ID
	if aEntry.TokenID != ent.ID {
		t.Fatalf("bad: got\n%s\nexpected\n%s\n", aEntry.TokenID, ent.ID)
	}
}

func TestTokenStore_HandleRequest_LookupAccessor(t *testing.T) {
	_, ts, _, root := TestCoreWithTokenStore(t)
	testMakeToken(t, ts, root, "tokenid", "", []string{"foo"})
	out, err := ts.Lookup("tokenid")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if out == nil {
		t.Fatalf("err: %s", err)
	}

	req := logical.TestRequest(t, logical.UpdateOperation, "lookup-accessor")
	req.Data = map[string]interface{}{
		"accessor": out.Accessor,
	}

	resp, err := ts.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if resp.Data == nil {
		t.Fatalf("response should contain data")
	}

	if resp.Data["accessor"].(string) == "" {
		t.Fatalf("accessor should not be empty")
	}

	// Verify that the lookup-accessor operation does not return the token ID
	if resp.Data["id"].(string) != "" {
		t.Fatalf("token ID should not be returned")
	}
}

func TestTokenStore_HandleRequest_ListAccessors(t *testing.T) {
	_, ts, _, root := TestCoreWithTokenStore(t)

	testKeys := []string{"token1", "token2", "token3", "token4"}
	for _, key := range testKeys {
		testMakeToken(t, ts, root, key, "", []string{"foo"})
	}

	// Revoke root to make the number of accessors match
	salted, err := ts.SaltID(root)
	if err != nil {
		t.Fatal(err)
	}
	ts.revokeSalted(salted)

	req := logical.TestRequest(t, logical.ListOperation, "accessors")

	resp, err := ts.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if resp.Data == nil {
		t.Fatalf("response should contain data")
	}
	if resp.Data["keys"] == nil {
		t.Fatalf("keys should not be empty")
	}
	keys := resp.Data["keys"].([]string)
	if len(keys) != len(testKeys) {
		t.Fatalf("wrong number of accessors found")
	}
	if len(resp.Warnings) != 0 {
		t.Fatalf("got warnings:\n%#v", resp.Warnings)
	}

	// Test upgrade from old struct method of accessor storage (of token id)
	for _, accessor := range keys {
		aEntry, err := ts.lookupByAccessor(accessor, false)
		if err != nil {
			t.Fatal(err)
		}
		if aEntry.TokenID == "" || aEntry.AccessorID == "" {
			t.Fatalf("error, accessor entry looked up is empty, but no error thrown")
		}
		salted, err := ts.SaltID(accessor)
		if err != nil {
			t.Fatal(err)
		}
		path := accessorPrefix + salted
		le := &logical.StorageEntry{Key: path, Value: []byte(aEntry.TokenID)}
		if err := ts.view.Put(le); err != nil {
			t.Fatalf("failed to persist accessor index entry: %v", err)
		}
	}

	// Do the lookup again, should get same result
	resp, err = ts.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if resp.Data == nil {
		t.Fatalf("response should contain data")
	}
	if resp.Data["keys"] == nil {
		t.Fatalf("keys should not be empty")
	}
	keys2 := resp.Data["keys"].([]string)
	if len(keys) != len(testKeys) {
		t.Fatalf("wrong number of accessors found")
	}
	if len(resp.Warnings) != 0 {
		t.Fatalf("got warnings:\n%#v", resp.Warnings)
	}

	for _, accessor := range keys2 {
		aEntry, err := ts.lookupByAccessor(accessor, false)
		if err != nil {
			t.Fatal(err)
		}
		if aEntry.TokenID == "" || aEntry.AccessorID == "" {
			t.Fatalf("error, accessor entry looked up is empty, but no error thrown")
		}
	}
}

func TestTokenStore_HandleRequest_RevokeAccessor(t *testing.T) {
	_, ts, _, root := TestCoreWithTokenStore(t)
	testMakeToken(t, ts, root, "tokenid", "", []string{"foo"})
	out, err := ts.Lookup("tokenid")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if out == nil {
		t.Fatalf("err: %s", err)
	}

	req := logical.TestRequest(t, logical.UpdateOperation, "revoke-accessor")
	req.Data = map[string]interface{}{
		"accessor": out.Accessor,
	}

	_, err = ts.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	out, err = ts.Lookup("tokenid")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if out != nil {
		t.Fatalf("bad:\ngot %#v\nexpected: nil\n", out)
	}
}

func TestTokenStore_RootToken(t *testing.T) {
	_, ts, _, _ := TestCoreWithTokenStore(t)

	te, err := ts.rootToken()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if te.ID == "" {
		t.Fatalf("missing ID")
	}

	out, err := ts.Lookup(te.ID)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !reflect.DeepEqual(out, te) {
		t.Fatalf("bad: expected:%#v\nactual:%#v", te, out)
	}
}

func TestTokenStore_CreateLookup(t *testing.T) {
	c, ts, _, _ := TestCoreWithTokenStore(t)

	ent := &TokenEntry{Path: "test", Policies: []string{"dev", "ops"}}
	if err := ts.create(ent); err != nil {
		t.Fatalf("err: %v", err)
	}
	if ent.ID == "" {
		t.Fatalf("missing ID")
	}

	out, err := ts.Lookup(ent.ID)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !reflect.DeepEqual(out, ent) {
		t.Fatalf("bad: expected:%#v\nactual:%#v", ent, out)
	}

	// New store should share the salt
	ts2, err := NewTokenStore(c, getBackendConfig(c))
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	ts2.SetExpirationManager(c.expiration)

	if err := ts2.Initialize(); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Should still match
	out, err = ts2.Lookup(ent.ID)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !reflect.DeepEqual(out, ent) {
		t.Fatalf("bad: expected:%#v\nactual:%#v", ent, out)
	}
}

func TestTokenStore_CreateLookup_ProvidedID(t *testing.T) {
	c, ts, _, _ := TestCoreWithTokenStore(t)

	ent := &TokenEntry{
		ID:       "foobarbaz",
		Path:     "test",
		Policies: []string{"dev", "ops"},
	}
	if err := ts.create(ent); err != nil {
		t.Fatalf("err: %v", err)
	}
	if ent.ID != "foobarbaz" {
		t.Fatalf("bad: ent.ID: expected:\"foobarbaz\"\n actual:%s", ent.ID)
	}
	if err := ts.create(ent); err == nil {
		t.Fatal("expected error creating token with the same ID")
	}

	out, err := ts.Lookup(ent.ID)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !reflect.DeepEqual(out, ent) {
		t.Fatalf("bad: expected:%#v\nactual:%#v", ent, out)
	}

	// New store should share the salt
	ts2, err := NewTokenStore(c, getBackendConfig(c))
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	ts2.SetExpirationManager(c.expiration)

	if err := ts2.Initialize(); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Should still match
	out, err = ts2.Lookup(ent.ID)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !reflect.DeepEqual(out, ent) {
		t.Fatalf("bad: expected:%#v\nactual:%#v", ent, out)
	}
}

func TestTokenStore_CreateLookup_ExpirationInRestoreMode(t *testing.T) {
	_, ts, _, _ := TestCoreWithTokenStore(t)

	ent := &TokenEntry{Path: "test", Policies: []string{"dev", "ops"}}
	if err := ts.create(ent); err != nil {
		t.Fatalf("err: %v", err)
	}
	if ent.ID == "" {
		t.Fatalf("missing ID")
	}

	// Replace the lease with a lease with an expire time in the past
	saltedID, err := ts.SaltID(ent.ID)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Create a lease entry
	leaseID := path.Join(ent.Path, saltedID)
	le := &leaseEntry{
		LeaseID:     leaseID,
		ClientToken: ent.ID,
		Path:        ent.Path,
		IssueTime:   time.Now(),
		ExpireTime:  time.Now().Add(1 * time.Hour),
	}
	if err := ts.expiration.persistEntry(le); err != nil {
		t.Fatalf("err: %v", err)
	}

	out, err := ts.Lookup(ent.ID)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !reflect.DeepEqual(out, ent) {
		t.Fatalf("bad: expected:%#v\nactual:%#v", ent, out)
	}

	// Set to expired lease time
	le.ExpireTime = time.Now().Add(-1 * time.Hour)
	if err := ts.expiration.persistEntry(le); err != nil {
		t.Fatalf("err: %v", err)
	}

	err = ts.expiration.Stop()
	if err != nil {
		t.Fatal(err)
	}

	// Reset expiration manager to restore mode
	ts.expiration.restoreModeLock.Lock()
	atomic.StoreInt32(&ts.expiration.restoreMode, 1)
	ts.expiration.restoreLocks = locksutil.CreateLocks()
	ts.expiration.restoreModeLock.Unlock()

	// Test that the token lookup does not return the token entry due to the
	// expired lease
	out, err = ts.Lookup(ent.ID)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out != nil {
		t.Fatalf("lease expired, no token expected: %#v", out)
	}
}

func TestTokenStore_UseToken(t *testing.T) {
	_, ts, _, root := TestCoreWithTokenStore(t)

	// Lookup the root token
	ent, err := ts.Lookup(root)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Root is an unlimited use token, should be a no-op
	te, err := ts.UseToken(ent)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if te == nil {
		t.Fatalf("token entry after use was nil")
	}

	// Lookup the root token again
	ent2, err := ts.Lookup(root)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if !reflect.DeepEqual(ent, ent2) {
		t.Fatalf("bad: ent:%#v ent2:%#v", ent, ent2)
	}

	// Create a retstricted token
	ent = &TokenEntry{Path: "test", Policies: []string{"dev", "ops"}, NumUses: 2}
	if err := ts.create(ent); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Use the token
	te, err = ts.UseToken(ent)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if te == nil {
		t.Fatalf("token entry for use #1 was nil")
	}

	// Lookup the token
	ent2, err = ts.Lookup(ent.ID)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Should be reduced
	if ent2.NumUses != 1 {
		t.Fatalf("bad: %#v", ent2)
	}

	// Use the token
	te, err = ts.UseToken(ent)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if te == nil {
		t.Fatalf("token entry for use #2 was nil")
	}
	if te.NumUses != tokenRevocationDeferred {
		t.Fatalf("token entry after use #2 did not have revoke flag")
	}
	ts.Revoke(te.ID)

	// Lookup the token
	ent2, err = ts.Lookup(ent.ID)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Should be revoked
	if ent2 != nil {
		t.Fatalf("bad: %#v", ent2)
	}
}

func TestTokenStore_Revoke(t *testing.T) {
	_, ts, _, _ := TestCoreWithTokenStore(t)

	ent := &TokenEntry{Path: "test", Policies: []string{"dev", "ops"}}
	if err := ts.create(ent); err != nil {
		t.Fatalf("err: %v", err)
	}

	err := ts.Revoke("")
	if err.Error() != "cannot revoke blank token" {
		t.Fatalf("err: %v", err)
	}
	err = ts.Revoke(ent.ID)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	out, err := ts.Lookup(ent.ID)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out != nil {
		t.Fatalf("bad: %#v", out)
	}
}

func TestTokenStore_Revoke_Leases(t *testing.T) {
	c, ts, _, _ := TestCoreWithTokenStore(t)

	view := NewBarrierView(c.barrier, "noop/")

	// Mount a noop backend
	noop := &NoopBackend{}
	err := ts.expiration.router.Mount(noop, "noop/", &MountEntry{UUID: "noopuuid", Accessor: "noopaccessor"}, view)
	if err != nil {
		t.Fatal(err)
	}

	ent := &TokenEntry{Path: "test", Policies: []string{"dev", "ops"}}
	if err := ts.create(ent); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Register a lease
	req := &logical.Request{
		Operation:   logical.ReadOperation,
		Path:        "noop/foo",
		ClientToken: ent.ID,
	}
	resp := &logical.Response{
		Secret: &logical.Secret{
			LeaseOptions: logical.LeaseOptions{
				TTL: 20 * time.Millisecond,
			},
		},
		Data: map[string]interface{}{
			"access_key": "xyz",
			"secret_key": "abcd",
		},
	}
	leaseID, err := ts.expiration.Register(req, resp)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Revoke the token
	err = ts.Revoke(ent.ID)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Verify the lease is gone
	out, err := ts.expiration.loadEntry(leaseID)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out != nil {
		t.Fatalf("bad: %#v", out)
	}
}

func TestTokenStore_Revoke_Orphan(t *testing.T) {
	_, ts, _, _ := TestCoreWithTokenStore(t)

	ent := &TokenEntry{Path: "test", Policies: []string{"dev", "ops"}}
	if err := ts.create(ent); err != nil {
		t.Fatalf("err: %v", err)
	}

	ent2 := &TokenEntry{Parent: ent.ID}
	if err := ts.create(ent2); err != nil {
		t.Fatalf("err: %v", err)
	}

	err := ts.Revoke(ent.ID)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	out, err := ts.Lookup(ent2.ID)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !reflect.DeepEqual(out, ent2) {
		t.Fatalf("bad: expected:%#v\nactual:%#v", ent2, out)
	}
}

func TestTokenStore_RevokeTree(t *testing.T) {
	_, ts, _, _ := TestCoreWithTokenStore(t)

	ent1 := &TokenEntry{}
	if err := ts.create(ent1); err != nil {
		t.Fatalf("err: %v", err)
	}

	ent2 := &TokenEntry{Parent: ent1.ID}
	if err := ts.create(ent2); err != nil {
		t.Fatalf("err: %v", err)
	}

	ent3 := &TokenEntry{Parent: ent2.ID}
	if err := ts.create(ent3); err != nil {
		t.Fatalf("err: %v", err)
	}

	ent4 := &TokenEntry{Parent: ent2.ID}
	if err := ts.create(ent4); err != nil {
		t.Fatalf("err: %v", err)
	}

	err := ts.RevokeTree("")
	if err.Error() != "cannot tree-revoke blank token" {
		t.Fatalf("err: %v", err)
	}
	err = ts.RevokeTree(ent1.ID)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	lookup := []string{ent1.ID, ent2.ID, ent3.ID, ent4.ID}
	for _, id := range lookup {
		out, err := ts.Lookup(id)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if out != nil {
			t.Fatalf("bad: %#v", out)
		}
	}
}

func TestTokenStore_RevokeSelf(t *testing.T) {
	_, ts, _, _ := TestCoreWithTokenStore(t)

	ent1 := &TokenEntry{}
	if err := ts.create(ent1); err != nil {
		t.Fatalf("err: %v", err)
	}

	ent2 := &TokenEntry{Parent: ent1.ID}
	if err := ts.create(ent2); err != nil {
		t.Fatalf("err: %v", err)
	}

	ent3 := &TokenEntry{Parent: ent2.ID}
	if err := ts.create(ent3); err != nil {
		t.Fatalf("err: %v", err)
	}

	ent4 := &TokenEntry{Parent: ent2.ID}
	if err := ts.create(ent4); err != nil {
		t.Fatalf("err: %v", err)
	}

	req := logical.TestRequest(t, logical.UpdateOperation, "revoke-self")
	req.ClientToken = ent1.ID

	resp, err := ts.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}

	lookup := []string{ent1.ID, ent2.ID, ent3.ID, ent4.ID}
	for _, id := range lookup {
		out, err := ts.Lookup(id)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if out != nil {
			t.Fatalf("bad: %#v", out)
		}
	}
}

func TestTokenStore_HandleRequest_NonAssignable(t *testing.T) {
	_, ts, _, root := TestCoreWithTokenStore(t)

	req := logical.TestRequest(t, logical.UpdateOperation, "create")
	req.ClientToken = root
	req.Data["policies"] = []string{"default", "foo"}

	resp, err := ts.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}

	req.Data["policies"] = []string{"default", "foo", responseWrappingPolicyName}

	resp, err = ts.HandleRequest(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("got a nil response")
	}
	if !resp.IsError() {
		t.Fatalf("expected error; response is %#v", *resp)
	}
}

func TestTokenStore_HandleRequest_CreateToken_DisplayName(t *testing.T) {
	_, ts, _, root := TestCoreWithTokenStore(t)

	req := logical.TestRequest(t, logical.UpdateOperation, "create")
	req.ClientToken = root
	req.Data["display_name"] = "foo_bar.baz!"

	resp, err := ts.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}

	expected := &TokenEntry{
		ID:          resp.Auth.ClientToken,
		Accessor:    resp.Auth.Accessor,
		Parent:      root,
		Policies:    []string{"root"},
		Path:        "auth/token/create",
		DisplayName: "token-foo-bar-baz",
		TTL:         0,
	}
	out, err := ts.Lookup(resp.Auth.ClientToken)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	expected.CreationTime = out.CreationTime
	if !reflect.DeepEqual(out, expected) {
		t.Fatalf("bad: expected:%#v\nactual:%#v", expected, out)
	}
}

func TestTokenStore_HandleRequest_CreateToken_NumUses(t *testing.T) {
	_, ts, _, root := TestCoreWithTokenStore(t)

	req := logical.TestRequest(t, logical.UpdateOperation, "create")
	req.ClientToken = root
	req.Data["num_uses"] = "1"

	resp, err := ts.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}

	expected := &TokenEntry{
		ID:          resp.Auth.ClientToken,
		Accessor:    resp.Auth.Accessor,
		Parent:      root,
		Policies:    []string{"root"},
		Path:        "auth/token/create",
		DisplayName: "token",
		NumUses:     1,
		TTL:         0,
	}
	out, err := ts.Lookup(resp.Auth.ClientToken)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	expected.CreationTime = out.CreationTime
	if !reflect.DeepEqual(out, expected) {
		t.Fatalf("bad: expected:%#v\nactual:%#v", expected, out)
	}
}

func TestTokenStore_HandleRequest_CreateToken_NumUses_Invalid(t *testing.T) {
	_, ts, _, root := TestCoreWithTokenStore(t)

	req := logical.TestRequest(t, logical.UpdateOperation, "create")
	req.ClientToken = root
	req.Data["num_uses"] = "-1"

	resp, err := ts.HandleRequest(req)
	if err != logical.ErrInvalidRequest {
		t.Fatalf("err: %v %v", err, resp)
	}
}

func TestTokenStore_HandleRequest_CreateToken_NumUses_Restricted(t *testing.T) {
	_, ts, _, root := TestCoreWithTokenStore(t)

	req := logical.TestRequest(t, logical.UpdateOperation, "create")
	req.ClientToken = root
	req.Data["num_uses"] = "1"

	resp, err := ts.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}

	// We should NOT be able to use the restricted token to create a new token
	req.ClientToken = resp.Auth.ClientToken
	_, err = ts.HandleRequest(req)
	if err != logical.ErrInvalidRequest {
		t.Fatalf("err: %v %v", err, resp)
	}
}

func TestTokenStore_HandleRequest_CreateToken_NoPolicy(t *testing.T) {
	_, ts, _, root := TestCoreWithTokenStore(t)

	req := logical.TestRequest(t, logical.UpdateOperation, "create")
	req.ClientToken = root

	resp, err := ts.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}

	expected := &TokenEntry{
		ID:          resp.Auth.ClientToken,
		Accessor:    resp.Auth.Accessor,
		Parent:      root,
		Policies:    []string{"root"},
		Path:        "auth/token/create",
		DisplayName: "token",
		TTL:         0,
	}
	out, err := ts.Lookup(resp.Auth.ClientToken)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	expected.CreationTime = out.CreationTime
	if !reflect.DeepEqual(out, expected) {
		t.Fatalf("bad: expected:%#v\nactual:%#v", expected, out)
	}
}

func TestTokenStore_HandleRequest_CreateToken_BadParent(t *testing.T) {
	_, ts, _, _ := TestCoreWithTokenStore(t)

	req := logical.TestRequest(t, logical.UpdateOperation, "create")
	req.ClientToken = "random"

	resp, err := ts.HandleRequest(req)
	if err != logical.ErrInvalidRequest {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp.Data["error"] != "parent token lookup failed" {
		t.Fatalf("bad: %#v", resp)
	}
}

func TestTokenStore_HandleRequest_CreateToken(t *testing.T) {
	_, ts, _, root := TestCoreWithTokenStore(t)

	req := logical.TestRequest(t, logical.UpdateOperation, "create")
	req.ClientToken = root
	req.Data["policies"] = []string{"foo"}

	resp, err := ts.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp.Auth.ClientToken == "" {
		t.Fatalf("bad: %#v", resp)
	}
}

func TestTokenStore_HandleRequest_CreateToken_RootID(t *testing.T) {
	_, ts, _, root := TestCoreWithTokenStore(t)

	req := logical.TestRequest(t, logical.UpdateOperation, "create")
	req.ClientToken = root
	req.Data["id"] = "foobar"
	req.Data["policies"] = []string{"foo"}

	resp, err := ts.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp.Auth.ClientToken != "foobar" {
		t.Fatalf("bad: %#v", resp)
	}
}

func TestTokenStore_HandleRequest_CreateToken_NonRootID(t *testing.T) {
	_, ts, _, root := TestCoreWithTokenStore(t)
	testMakeToken(t, ts, root, "client", "", []string{"foo"})

	req := logical.TestRequest(t, logical.UpdateOperation, "create")
	req.ClientToken = "client"
	req.Data["id"] = "foobar"
	req.Data["policies"] = []string{"foo"}

	resp, err := ts.HandleRequest(req)
	if err != logical.ErrInvalidRequest {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp.Data["error"] != "root or sudo privileges required to specify token id" {
		t.Fatalf("bad: %#v", resp)
	}
}

func TestTokenStore_HandleRequest_CreateToken_NonRoot_Subset(t *testing.T) {
	_, ts, _, root := TestCoreWithTokenStore(t)
	testMakeToken(t, ts, root, "client", "", []string{"foo", "bar"})

	req := logical.TestRequest(t, logical.UpdateOperation, "create")
	req.ClientToken = "client"
	req.Data["policies"] = []string{"foo"}

	resp, err := ts.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp.Auth.ClientToken == "" {
		t.Fatalf("bad: %#v", resp)
	}
}

func TestTokenStore_HandleRequest_CreateToken_NonRoot_InvalidSubset(t *testing.T) {
	_, ts, _, root := TestCoreWithTokenStore(t)
	testMakeToken(t, ts, root, "client", "", []string{"foo", "bar"})

	req := logical.TestRequest(t, logical.UpdateOperation, "create")
	req.ClientToken = "client"
	req.Data["policies"] = []string{"foo", "bar", "baz"}

	resp, err := ts.HandleRequest(req)
	if err != logical.ErrInvalidRequest {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp.Data["error"] != "child policies must be subset of parent" {
		t.Fatalf("bad: %#v", resp)
	}
}

func TestTokenStore_HandleRequest_CreateToken_NonRoot_RootChild(t *testing.T) {
	core, ts, _, root := TestCoreWithTokenStore(t)
	ps := core.policyStore

	policy, _ := ParseACLPolicy(tokenCreationPolicy)
	policy.Name = "test1"
	if err := ps.SetPolicy(policy); err != nil {
		t.Fatal(err)
	}

	testMakeToken(t, ts, root, "sudoClient", "", []string{"test1"})

	req := logical.TestRequest(t, logical.UpdateOperation, "create")
	req.ClientToken = "sudoClient"
	req.MountPoint = "auth/token/"
	req.Data["policies"] = []string{"root"}

	resp, err := ts.HandleRequest(req)
	if err != logical.ErrInvalidRequest {
		t.Fatalf("err: %v; resp: %#v", err, resp)
	}
	if resp == nil || resp.Data == nil {
		t.Fatalf("expected a response")
	}
	if resp.Data["error"].(string) != "root tokens may not be created without parent token being root" {
		t.Fatalf("bad: %#v", resp)
	}
}

func TestTokenStore_HandleRequest_CreateToken_Root_RootChild_NoExpiry_Expiry(t *testing.T) {
	_, ts, _, root := TestCoreWithTokenStore(t)

	req := logical.TestRequest(t, logical.UpdateOperation, "create")
	req.ClientToken = root
	req.Data = map[string]interface{}{
		"ttl": "5m",
	}

	resp, err := ts.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v; resp: %#v", err, resp)
	}
	if resp == nil || resp.Auth == nil {
		t.Fatalf("failed to create a root token using another root token")
	}
	if !reflect.DeepEqual(resp.Auth.Policies, []string{"root"}) {
		t.Fatalf("bad: policies: expected: root; actual: %s", resp.Auth.Policies)
	}
	if resp.Auth.TTL.Seconds() != 300 {
		t.Fatalf("bad: expected 300 second ttl, got %v", resp.Auth.TTL.Seconds())
	}

	req.ClientToken = resp.Auth.ClientToken
	req.Data = map[string]interface{}{
		"ttl": "0",
	}
	resp, err = ts.HandleRequest(req)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestTokenStore_HandleRequest_CreateToken_Root_RootChild(t *testing.T) {
	_, ts, _, root := TestCoreWithTokenStore(t)

	req := logical.TestRequest(t, logical.UpdateOperation, "create")
	req.ClientToken = root

	resp, err := ts.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v; resp: %#v", err, resp)
	}
	if resp == nil || resp.Auth == nil {
		t.Fatalf("failed to create a root token using another root token")
	}
	if !reflect.DeepEqual(resp.Auth.Policies, []string{"root"}) {
		t.Fatalf("bad: policies: expected: root; actual: %s", resp.Auth.Policies)
	}
}

func TestTokenStore_HandleRequest_CreateToken_NonRoot_NoParent(t *testing.T) {
	_, ts, _, root := TestCoreWithTokenStore(t)
	testMakeToken(t, ts, root, "client", "", []string{"foo"})

	req := logical.TestRequest(t, logical.UpdateOperation, "create")
	req.ClientToken = "client"
	req.Data["no_parent"] = true
	req.Data["policies"] = []string{"foo"}

	resp, err := ts.HandleRequest(req)
	if err != logical.ErrInvalidRequest {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp.Data["error"] != "root or sudo privileges required to create orphan token" {
		t.Fatalf("bad: %#v", resp)
	}
}

func TestTokenStore_HandleRequest_CreateToken_Root_NoParent(t *testing.T) {
	_, ts, _, root := TestCoreWithTokenStore(t)

	req := logical.TestRequest(t, logical.UpdateOperation, "create")
	req.ClientToken = root
	req.Data["no_parent"] = true
	req.Data["policies"] = []string{"foo"}

	resp, err := ts.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp.Auth.ClientToken == "" {
		t.Fatalf("bad: %#v", resp)
	}

	out, _ := ts.Lookup(resp.Auth.ClientToken)
	if out.Parent != "" {
		t.Fatalf("bad: %#v", out)
	}
}

func TestTokenStore_HandleRequest_CreateToken_PathBased_NoParent(t *testing.T) {
	_, ts, _, root := TestCoreWithTokenStore(t)

	req := logical.TestRequest(t, logical.UpdateOperation, "create-orphan")
	req.ClientToken = root
	req.Data["policies"] = []string{"foo"}

	resp, err := ts.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp.Auth.ClientToken == "" {
		t.Fatalf("bad: %#v", resp)
	}

	out, _ := ts.Lookup(resp.Auth.ClientToken)
	if out.Parent != "" {
		t.Fatalf("bad: %#v", out)
	}
}

func TestTokenStore_HandleRequest_CreateToken_Metadata(t *testing.T) {
	_, ts, _, root := TestCoreWithTokenStore(t)

	req := logical.TestRequest(t, logical.UpdateOperation, "create")
	req.ClientToken = root
	req.Data["policies"] = []string{"foo"}
	meta := map[string]string{
		"user":   "armon",
		"source": "github",
	}
	req.Data["meta"] = meta

	resp, err := ts.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp.Auth.ClientToken == "" {
		t.Fatalf("bad: %#v", resp)
	}

	out, _ := ts.Lookup(resp.Auth.ClientToken)
	if !reflect.DeepEqual(out.Meta, meta) {
		t.Fatalf("bad: expected:%#v\nactual:%#v", meta, out.Meta)
	}
}

func TestTokenStore_HandleRequest_CreateToken_Lease(t *testing.T) {
	_, ts, _, root := TestCoreWithTokenStore(t)

	req := logical.TestRequest(t, logical.UpdateOperation, "create")
	req.ClientToken = root
	req.Data["policies"] = []string{"foo"}
	req.Data["lease"] = "1h"

	resp, err := ts.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp.Auth.ClientToken == "" {
		t.Fatalf("bad: %#v", resp)
	}
	if resp.Auth.TTL != time.Hour {
		t.Fatalf("bad: %#v", resp)
	}
	if !resp.Auth.Renewable {
		t.Fatalf("bad: %#v", resp)
	}
}

func TestTokenStore_HandleRequest_CreateToken_TTL(t *testing.T) {
	_, ts, _, root := TestCoreWithTokenStore(t)

	req := logical.TestRequest(t, logical.UpdateOperation, "create")
	req.ClientToken = root
	req.Data["policies"] = []string{"foo"}
	req.Data["ttl"] = "1h"

	resp, err := ts.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp.Auth.ClientToken == "" {
		t.Fatalf("bad: %#v", resp)
	}
	if resp.Auth.TTL != time.Hour {
		t.Fatalf("bad: %#v", resp)
	}
	if !resp.Auth.Renewable {
		t.Fatalf("bad: %#v", resp)
	}
}

func TestTokenStore_HandleRequest_Revoke(t *testing.T) {
	_, ts, _, root := TestCoreWithTokenStore(t)
	testMakeToken(t, ts, root, "child", "", []string{"root", "foo"})
	testMakeToken(t, ts, "child", "sub-child", "", []string{"foo"})

	req := logical.TestRequest(t, logical.UpdateOperation, "revoke")
	req.Data = map[string]interface{}{
		"token": "child",
	}
	resp, err := ts.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp != nil {
		t.Fatalf("bad: %#v", resp)
	}

	out, err := ts.Lookup("child")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out != nil {
		t.Fatalf("bad: %v", out)
	}

	// Sub-child should not exist
	out, err = ts.Lookup("sub-child")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out != nil {
		t.Fatalf("bad: %v", out)
	}
}

func TestTokenStore_HandleRequest_RevokeOrphan(t *testing.T) {
	_, ts, _, root := TestCoreWithTokenStore(t)
	testMakeToken(t, ts, root, "child", "", []string{"root", "foo"})
	testMakeToken(t, ts, "child", "sub-child", "", []string{"foo"})

	req := logical.TestRequest(t, logical.UpdateOperation, "revoke-orphan")
	req.Data = map[string]interface{}{
		"token": "child",
	}
	req.ClientToken = root
	resp, err := ts.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp != nil {
		t.Fatalf("bad: %#v", resp)
	}

	out, err := ts.Lookup("child")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out != nil {
		t.Fatalf("bad: %v", out)
	}

	// Sub-child should exist!
	out, err = ts.Lookup("sub-child")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out == nil {
		t.Fatalf("bad: %v", out)
	}
}

func TestTokenStore_HandleRequest_RevokeOrphan_NonRoot(t *testing.T) {
	_, ts, _, root := TestCoreWithTokenStore(t)
	testMakeToken(t, ts, root, "child", "", []string{"foo"})

	out, err := ts.Lookup("child")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out == nil {
		t.Fatalf("bad: %v", out)
	}

	req := logical.TestRequest(t, logical.UpdateOperation, "revoke-orphan")
	req.Data = map[string]interface{}{
		"token": "child",
	}
	req.ClientToken = "child"
	resp, err := ts.HandleRequest(req)
	if err != logical.ErrInvalidRequest {
		t.Fatalf("did not get error when non-root revoking itself with orphan flag; resp is %#v", resp)
	}

	// Should still exist
	out, err = ts.Lookup("child")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out == nil {
		t.Fatalf("bad: %v", out)
	}
}

func TestTokenStore_HandleRequest_Lookup(t *testing.T) {
	c, ts, _, root := TestCoreWithTokenStore(t)
	req := logical.TestRequest(t, logical.UpdateOperation, "lookup")
	req.Data = map[string]interface{}{
		"token": root,
	}
	resp, err := ts.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp == nil {
		t.Fatalf("bad: %#v", resp)
	}

	exp := map[string]interface{}{
		"id":               root,
		"accessor":         resp.Data["accessor"].(string),
		"policies":         []string{"root"},
		"path":             "auth/token/root",
		"meta":             map[string]string(nil),
		"display_name":     "root",
		"orphan":           true,
		"num_uses":         0,
		"creation_ttl":     int64(0),
		"ttl":              int64(0),
		"explicit_max_ttl": int64(0),
		"expire_time":      nil,
		"entity_id":        "",
	}

	if resp.Data["creation_time"].(int64) == 0 {
		t.Fatalf("creation time was zero")
	}
	delete(resp.Data, "creation_time")

	if !reflect.DeepEqual(resp.Data, exp) {
		t.Fatalf("bad: expected:%#v\nactual:%#v", exp, resp.Data)
	}

	testCoreMakeToken(t, c, root, "client", "3600s", []string{"foo"})

	// Test via GET
	req = logical.TestRequest(t, logical.UpdateOperation, "lookup")
	req.Data = map[string]interface{}{
		"token": "client",
	}
	resp, err = ts.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp == nil {
		t.Fatalf("bad: %#v", resp)
	}

	exp = map[string]interface{}{
		"id":               "client",
		"accessor":         resp.Data["accessor"],
		"policies":         []string{"default", "foo"},
		"path":             "auth/token/create",
		"meta":             map[string]string(nil),
		"display_name":     "token",
		"orphan":           false,
		"num_uses":         0,
		"creation_ttl":     int64(3600),
		"ttl":              int64(3600),
		"explicit_max_ttl": int64(0),
		"renewable":        true,
		"entity_id":        "",
	}

	if resp.Data["creation_time"].(int64) == 0 {
		t.Fatalf("creation time was zero")
	}
	delete(resp.Data, "creation_time")
	if resp.Data["issue_time"].(time.Time).IsZero() {
		t.Fatal("issue time is default time")
	}
	delete(resp.Data, "issue_time")
	if resp.Data["expire_time"].(time.Time).IsZero() {
		t.Fatal("expire time is default time")
	}
	delete(resp.Data, "expire_time")

	// Depending on timing of the test this may have ticked down, so accept 3599
	if resp.Data["ttl"].(int64) == 3599 {
		resp.Data["ttl"] = int64(3600)
	}

	if !reflect.DeepEqual(resp.Data, exp) {
		t.Fatalf("bad: expected:%#v\nactual:%#v", exp, resp.Data)
	}

	// Test via POST
	req = logical.TestRequest(t, logical.UpdateOperation, "lookup")
	req.Data = map[string]interface{}{
		"token": "client",
	}
	resp, err = ts.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp == nil {
		t.Fatalf("bad: %#v", resp)
	}

	exp = map[string]interface{}{
		"id":               "client",
		"accessor":         resp.Data["accessor"],
		"policies":         []string{"default", "foo"},
		"path":             "auth/token/create",
		"meta":             map[string]string(nil),
		"display_name":     "token",
		"orphan":           false,
		"num_uses":         0,
		"creation_ttl":     int64(3600),
		"ttl":              int64(3600),
		"explicit_max_ttl": int64(0),
		"renewable":        true,
		"entity_id":        "",
	}

	if resp.Data["creation_time"].(int64) == 0 {
		t.Fatalf("creation time was zero")
	}
	delete(resp.Data, "creation_time")
	if resp.Data["issue_time"].(time.Time).IsZero() {
		t.Fatal("issue time is default time")
	}
	delete(resp.Data, "issue_time")
	if resp.Data["expire_time"].(time.Time).IsZero() {
		t.Fatal("expire time is default time")
	}
	delete(resp.Data, "expire_time")

	// Depending on timing of the test this may have ticked down, so accept 3599
	if resp.Data["ttl"].(int64) == 3599 {
		resp.Data["ttl"] = int64(3600)
	}

	if !reflect.DeepEqual(resp.Data, exp) {
		t.Fatalf("bad: expected:%#v\nactual:%#v", exp, resp.Data)
	}

	// Test last_renewal_time functionality
	req = logical.TestRequest(t, logical.UpdateOperation, "renew")
	req.Data = map[string]interface{}{
		"token": "client",
	}
	resp, err = ts.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp == nil {
		t.Fatalf("bad: %#v", resp)
	}

	req = logical.TestRequest(t, logical.UpdateOperation, "lookup")
	req.Data = map[string]interface{}{
		"token": "client",
	}
	resp, err = ts.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp == nil {
		t.Fatalf("bad: %#v", resp)
	}

	if resp.Data["last_renewal_time"].(int64) == 0 {
		t.Fatalf("last_renewal_time was zero")
	}
}

func TestTokenStore_HandleRequest_LookupSelf(t *testing.T) {
	c, ts, _, root := TestCoreWithTokenStore(t)
	testCoreMakeToken(t, c, root, "client", "3600s", []string{"foo"})

	req := logical.TestRequest(t, logical.ReadOperation, "lookup-self")
	req.ClientToken = "client"
	resp, err := ts.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp == nil {
		t.Fatalf("bad: %#v", resp)
	}

	exp := map[string]interface{}{
		"id":               "client",
		"accessor":         resp.Data["accessor"],
		"policies":         []string{"default", "foo"},
		"path":             "auth/token/create",
		"meta":             map[string]string(nil),
		"display_name":     "token",
		"orphan":           false,
		"renewable":        true,
		"num_uses":         0,
		"creation_ttl":     int64(3600),
		"ttl":              int64(3600),
		"explicit_max_ttl": int64(0),
		"entity_id":        "",
	}

	if resp.Data["creation_time"].(int64) == 0 {
		t.Fatalf("creation time was zero")
	}
	delete(resp.Data, "creation_time")
	if resp.Data["issue_time"].(time.Time).IsZero() {
		t.Fatalf("creation time was zero")
	}
	delete(resp.Data, "issue_time")
	if resp.Data["expire_time"].(time.Time).IsZero() {
		t.Fatalf("expire time was zero")
	}
	delete(resp.Data, "expire_time")

	// Depending on timing of the test this may have ticked down, so accept 3599
	if resp.Data["ttl"].(int64) == 3599 {
		resp.Data["ttl"] = int64(3600)
	}

	if !reflect.DeepEqual(resp.Data, exp) {
		t.Fatalf("bad: expected:%#v\nactual:%#v", exp, resp.Data)
	}
}

func TestTokenStore_HandleRequest_Renew(t *testing.T) {
	exp := mockExpiration(t)
	ts := exp.tokenStore

	// Create new token
	root, err := ts.rootToken()
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Create a new token
	auth := &logical.Auth{
		ClientToken: root.ID,
		LeaseOptions: logical.LeaseOptions{
			TTL:       time.Hour,
			Renewable: true,
		},
	}
	err = exp.RegisterAuth("auth/token/root", auth)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Get the original expire time to compare
	originalExpire := auth.ExpirationTime()

	beforeRenew := time.Now()
	req := logical.TestRequest(t, logical.UpdateOperation, "renew")
	req.Data = map[string]interface{}{
		"token":     root.ID,
		"increment": "3600s",
	}

	req.Data["increment"] = "3600s"
	resp, err := ts.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}

	// Get the new expire time
	newExpire := resp.Auth.ExpirationTime()
	if newExpire.Before(originalExpire) {
		t.Fatalf("should expire later: %s %s", newExpire, originalExpire)
	}
	if newExpire.Before(beforeRenew.Add(time.Hour)) {
		t.Fatalf("should have at least an hour: %s %s", newExpire, beforeRenew)
	}
}

func TestTokenStore_HandleRequest_RenewSelf(t *testing.T) {
	exp := mockExpiration(t)
	ts := exp.tokenStore

	// Create new token
	root, err := ts.rootToken()
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Create a new token
	auth := &logical.Auth{
		ClientToken: root.ID,
		LeaseOptions: logical.LeaseOptions{
			TTL:       time.Hour,
			Renewable: true,
		},
	}
	err = exp.RegisterAuth("auth/token/root", auth)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Get the original expire time to compare
	originalExpire := auth.ExpirationTime()

	beforeRenew := time.Now()
	req := logical.TestRequest(t, logical.UpdateOperation, "renew-self")
	req.ClientToken = auth.ClientToken
	req.Data["increment"] = "3600s"
	resp, err := ts.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}

	// Get the new expire time
	newExpire := resp.Auth.ExpirationTime()
	if newExpire.Before(originalExpire) {
		t.Fatalf("should expire later: %s %s", newExpire, originalExpire)
	}
	if newExpire.Before(beforeRenew.Add(time.Hour)) {
		t.Fatalf("should have at least an hour: %s %s", newExpire, beforeRenew)
	}
}

func TestTokenStore_RoleCRUD(t *testing.T) {
	core, _, _, root := TestCoreWithTokenStore(t)

	req := logical.TestRequest(t, logical.ReadOperation, "auth/token/roles/test")
	req.ClientToken = root

	resp, err := core.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp != nil {
		t.Fatalf("should not see a role")
	}

	// First test creation
	req.Operation = logical.CreateOperation
	req.Data = map[string]interface{}{
		"orphan":           true,
		"period":           "72h",
		"allowed_policies": "test1,test2",
		"path_suffix":      "happenin",
	}

	resp, err = core.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp != nil {
		t.Fatalf("expected a nil response")
	}

	req.Operation = logical.ReadOperation
	req.Data = map[string]interface{}{}

	resp, err = core.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp == nil {
		t.Fatalf("got a nil response")
	}

	expected := map[string]interface{}{
		"name":                "test",
		"orphan":              true,
		"period":              int64(259200),
		"allowed_policies":    []string{"test1", "test2"},
		"disallowed_policies": []string{},
		"path_suffix":         "happenin",
		"explicit_max_ttl":    int64(0),
		"renewable":           true,
	}

	if !reflect.DeepEqual(expected, resp.Data) {
		t.Fatalf("bad: expected:%#v\nactual:%#v", expected, resp.Data)
	}

	// Now test updating; this should be set to an UpdateOperation
	// automatically due to the existence check
	req.Operation = logical.CreateOperation
	req.Data = map[string]interface{}{
		"period":           "79h",
		"allowed_policies": "test3",
		"path_suffix":      "happenin",
		"renewable":        false,
	}

	resp, err = core.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp != nil {
		t.Fatalf("expected a nil response")
	}

	req.Operation = logical.ReadOperation
	req.Data = map[string]interface{}{}

	resp, err = core.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp == nil {
		t.Fatalf("got a nil response")
	}

	expected = map[string]interface{}{
		"name":                "test",
		"orphan":              true,
		"period":              int64(284400),
		"allowed_policies":    []string{"test3"},
		"disallowed_policies": []string{},
		"path_suffix":         "happenin",
		"explicit_max_ttl":    int64(0),
		"renewable":           false,
	}

	if !reflect.DeepEqual(expected, resp.Data) {
		t.Fatalf("bad: expected:%#v\nactual:%#v", expected, resp.Data)
	}

	// Now set explicit max ttl and clear the period
	req.Operation = logical.CreateOperation
	req.Data = map[string]interface{}{
		"explicit_max_ttl": "5",
		"period":           "0s",
	}
	resp, err = core.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}

	req.Operation = logical.ReadOperation
	req.Data = map[string]interface{}{}

	resp, err = core.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp == nil {
		t.Fatalf("got a nil response")
	}

	expected = map[string]interface{}{
		"name":                "test",
		"orphan":              true,
		"explicit_max_ttl":    int64(5),
		"allowed_policies":    []string{"test3"},
		"disallowed_policies": []string{},
		"path_suffix":         "happenin",
		"period":              int64(0),
		"renewable":           false,
	}

	if !reflect.DeepEqual(expected, resp.Data) {
		t.Fatalf("bad: expected:%#v\nactual:%#v", expected, resp.Data)
	}

	req.Operation = logical.ListOperation
	req.Path = "auth/token/roles"
	req.Data = map[string]interface{}{}
	resp, err = core.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp == nil {
		t.Fatalf("got a nil response")
	}
	keysInt, ok := resp.Data["keys"]
	if !ok {
		t.Fatalf("did not find keys in response")
	}
	keys, ok := keysInt.([]string)
	if !ok {
		t.Fatalf("could not convert keys interface to key list")
	}
	if len(keys) != 1 {
		t.Fatalf("unexpected number of keys: %d", len(keys))
	}
	if keys[0] != "test" {
		t.Fatalf("expected \"test\", got \"%s\"", keys[0])
	}

	req.Operation = logical.DeleteOperation
	req.Path = "auth/token/roles/test"
	resp, err = core.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp != nil {
		t.Fatalf("expected a nil response")
	}

	req.Operation = logical.ReadOperation
	resp, err = core.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp != nil {
		t.Fatalf("expected a nil response")
	}
}

func TestTokenStore_RoleDisallowedPoliciesWithRoot(t *testing.T) {
	var resp *logical.Response
	var err error

	_, ts, _, root := TestCoreWithTokenStore(t)

	// Don't set disallowed_policies. Verify that a read on the role does return a non-nil value.
	roleReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "roles/role1",
		Data: map[string]interface{}{
			"disallowed_policies": "root,testpolicy",
		},
		ClientToken: root,
	}
	resp, err = ts.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%v", err, resp)
	}

	roleReq.Operation = logical.ReadOperation
	resp, err = ts.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%v", err, resp)
	}

	expected := []string{"root", "testpolicy"}
	if !reflect.DeepEqual(resp.Data["disallowed_policies"], expected) {
		t.Fatalf("bad: expected: %#v, actual: %#v", expected, resp.Data["disallowed_policies"])
	}
}

func TestTokenStore_RoleDisallowedPolicies(t *testing.T) {
	var req *logical.Request
	var resp *logical.Response
	var err error

	core, ts, _, root := TestCoreWithTokenStore(t)
	ps := core.policyStore

	// Create 3 different policies
	policy, _ := ParseACLPolicy(tokenCreationPolicy)
	policy.Name = "test1"
	if err := ps.SetPolicy(policy); err != nil {
		t.Fatal(err)
	}

	policy, _ = ParseACLPolicy(tokenCreationPolicy)
	policy.Name = "test2"
	if err := ps.SetPolicy(policy); err != nil {
		t.Fatal(err)
	}

	policy, _ = ParseACLPolicy(tokenCreationPolicy)
	policy.Name = "test3"
	if err := ps.SetPolicy(policy); err != nil {
		t.Fatal(err)
	}

	// Create roles with different disallowed_policies configuration
	req = logical.TestRequest(t, logical.UpdateOperation, "roles/test1")
	req.ClientToken = root
	req.Data = map[string]interface{}{
		"disallowed_policies": "test1",
	}
	resp, err = ts.HandleRequest(req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%v", err, resp)
	}

	req = logical.TestRequest(t, logical.UpdateOperation, "roles/test23")
	req.ClientToken = root
	req.Data = map[string]interface{}{
		"disallowed_policies": "test2,test3",
	}
	resp, err = ts.HandleRequest(req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%v", err, resp)
	}

	req = logical.TestRequest(t, logical.UpdateOperation, "roles/test123")
	req.ClientToken = root
	req.Data = map[string]interface{}{
		"disallowed_policies": "test1,test2,test3",
	}
	resp, err = ts.HandleRequest(req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%v", err, resp)
	}

	// Create a token that has all the policies defined above
	req = logical.TestRequest(t, logical.UpdateOperation, "create")
	req.ClientToken = root
	req.Data["policies"] = []string{"test1", "test2", "test3"}
	resp, err = ts.HandleRequest(req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%v", err, resp)
	}
	if resp == nil || resp.Auth == nil {
		t.Fatal("got nil response")
	}
	if resp.Auth.ClientToken == "" {
		t.Fatalf("bad: ClientToken; resp:%#v", resp)
	}
	parentToken := resp.Auth.ClientToken

	req = logical.TestRequest(t, logical.UpdateOperation, "create/test1")
	req.ClientToken = parentToken
	resp, err = ts.HandleRequest(req)
	if err == nil || resp != nil && !resp.IsError() {
		t.Fatalf("expected an error response, got %#v", resp)
	}

	req = logical.TestRequest(t, logical.UpdateOperation, "create/test23")
	req.ClientToken = parentToken
	resp, err = ts.HandleRequest(req)
	if err == nil || resp != nil && !resp.IsError() {
		t.Fatal("expected an error response")
	}

	req = logical.TestRequest(t, logical.UpdateOperation, "create/test123")
	req.ClientToken = parentToken
	resp, err = ts.HandleRequest(req)
	if err == nil || resp != nil && !resp.IsError() {
		t.Fatal("expected an error response")
	}

	// Disallowed should act as a blacklist so make sure we can still make
	// something with other policies in the request
	req = logical.TestRequest(t, logical.UpdateOperation, "create/test123")
	req.Data["policies"] = []string{"foo", "bar"}
	req.ClientToken = parentToken
	resp, err = ts.HandleRequest(req)
	if err != nil || resp == nil || resp.IsError() {
		t.Fatalf("err:%v resp:%v", err, resp)
	}

	// Create a role to have 'default' policy disallowed
	req = logical.TestRequest(t, logical.UpdateOperation, "roles/default")
	req.ClientToken = root
	req.Data = map[string]interface{}{
		"disallowed_policies": "default",
	}
	resp, err = ts.HandleRequest(req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%v", err, resp)
	}

	req = logical.TestRequest(t, logical.UpdateOperation, "create/default")
	req.ClientToken = parentToken
	resp, err = ts.HandleRequest(req)
	if err == nil || resp != nil && !resp.IsError() {
		t.Fatal("expected an error response")
	}
}

func TestTokenStore_RoleAllowedPolicies(t *testing.T) {
	_, ts, _, root := TestCoreWithTokenStore(t)

	req := logical.TestRequest(t, logical.UpdateOperation, "roles/test")
	req.ClientToken = root
	req.Data = map[string]interface{}{
		"allowed_policies": "test1,test2",
	}

	resp, err := ts.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp != nil {
		t.Fatalf("expected a nil response")
	}

	req.Data = map[string]interface{}{}

	req.Path = "create/test"
	req.Data["policies"] = []string{"foo"}
	resp, err = ts.HandleRequest(req)
	if err == nil {
		t.Fatalf("expected error")
	}

	req.Data["policies"] = []string{"test2"}
	resp, err = ts.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp.Auth.ClientToken == "" {
		t.Fatalf("bad: %#v", resp)
	}

	// When allowed_policies is blank, should fall back to a subset of the parent policies
	req = logical.TestRequest(t, logical.UpdateOperation, "roles/test")
	req.ClientToken = root
	req.Data = map[string]interface{}{
		"allowed_policies": "",
	}
	resp, err = ts.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp != nil {
		t.Fatalf("expected a nil response")
	}

	req = logical.TestRequest(t, logical.UpdateOperation, "create")
	req.ClientToken = root
	req.Data["policies"] = []string{"test1", "test2", "test3"}
	resp, err = ts.HandleRequest(req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%v", err, resp)
	}
	if resp == nil || resp.Auth == nil {
		t.Fatal("got nil response")
	}
	if resp.Auth.ClientToken == "" {
		t.Fatalf("bad: ClientToken; resp:%#v", resp)
	}
	if !reflect.DeepEqual(resp.Auth.Policies, []string{"default", "test1", "test2", "test3"}) {
		t.Fatalf("bad: %#v", resp.Auth.Policies)
	}
	parentToken := resp.Auth.ClientToken

	req.Data = map[string]interface{}{}
	req.ClientToken = parentToken

	req.Path = "create/test"
	req.Data["policies"] = []string{"foo"}
	resp, err = ts.HandleRequest(req)
	if err == nil {
		t.Fatalf("expected error")
	}

	req.Data["policies"] = []string{"test2"}
	resp, err = ts.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp.Auth.ClientToken == "" {
		t.Fatalf("bad: %#v", resp)
	}

	delete(req.Data, "policies")
	resp, err = ts.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp.Auth.ClientToken == "" {
		t.Fatalf("bad: %#v", resp)
	}
	if !reflect.DeepEqual(resp.Auth.Policies, []string{"default", "test1", "test2", "test3"}) {
		t.Fatalf("bad: %#v", resp.Auth.Policies)
	}
}

func TestTokenStore_RoleOrphan(t *testing.T) {
	_, ts, _, root := TestCoreWithTokenStore(t)

	req := logical.TestRequest(t, logical.UpdateOperation, "roles/test")
	req.ClientToken = root
	req.Data = map[string]interface{}{
		"orphan": true,
	}

	resp, err := ts.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp != nil {
		t.Fatalf("expected a nil response")
	}

	req.Path = "create/test"
	resp, err = ts.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp.Auth.ClientToken == "" {
		t.Fatalf("bad: %#v", resp)
	}

	out, err := ts.Lookup(resp.Auth.ClientToken)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if out.Parent != "" {
		t.Fatalf("expected orphan token, but found a parent")
	}

	if !strings.HasPrefix(out.Path, "auth/token/create/test") {
		t.Fatalf("expected role in path but did not find it")
	}
}

func TestTokenStore_RolePathSuffix(t *testing.T) {
	_, ts, _, root := TestCoreWithTokenStore(t)

	req := logical.TestRequest(t, logical.UpdateOperation, "roles/test")
	req.ClientToken = root
	req.Data = map[string]interface{}{
		"path_suffix": "happenin",
	}

	resp, err := ts.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp != nil {
		t.Fatalf("expected a nil response")
	}

	req.Path = "create/test"
	resp, err = ts.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp.Auth.ClientToken == "" {
		t.Fatalf("bad: %#v", resp)
	}

	out, err := ts.Lookup(resp.Auth.ClientToken)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if out.Path != "auth/token/create/test/happenin" {
		t.Fatalf("expected role in path but did not find it")
	}
}

func TestTokenStore_RolePeriod(t *testing.T) {
	core, _, _, root := TestCoreWithTokenStore(t)

	core.defaultLeaseTTL = 10 * time.Second
	core.maxLeaseTTL = 10 * time.Second

	// Note: these requests are sent to Core since Core handles registration
	// with the expiration manager and we need the storage to be consistent

	req := logical.TestRequest(t, logical.UpdateOperation, "auth/token/roles/test")
	req.ClientToken = root
	req.Data = map[string]interface{}{
		"period": 300,
	}

	resp, err := core.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp != nil {
		t.Fatalf("expected a nil response")
	}

	// This first set of logic is to verify that a normal non-root token will
	// be given a TTL of 10 seconds, and that renewing will not cause the TTL to
	// increase since that's the configured backend max. Then we verify that
	// increment works.
	{
		req.Path = "auth/token/create"
		req.Data = map[string]interface{}{
			"policies": []string{"default"},
		}
		resp, err = core.HandleRequest(req)
		if err != nil {
			t.Fatalf("err: %v %v", err, resp)
		}
		if resp.Auth.ClientToken == "" {
			t.Fatalf("bad: %#v", resp)
		}

		req.ClientToken = resp.Auth.ClientToken
		req.Operation = logical.ReadOperation
		req.Path = "auth/token/lookup-self"
		resp, err = core.HandleRequest(req)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		ttl := resp.Data["ttl"].(int64)
		if ttl > 10 {
			t.Fatalf("TTL too large")
		}

		// Let the TTL go down a bit to 8 seconds
		time.Sleep(2 * time.Second)

		req.Operation = logical.UpdateOperation
		req.Path = "auth/token/renew-self"
		resp, err = core.HandleRequest(req)
		if err != nil {
			t.Fatalf("err: %v %v", err, resp)
		}

		req.Operation = logical.ReadOperation
		req.Path = "auth/token/lookup-self"
		resp, err = core.HandleRequest(req)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		ttl = resp.Data["ttl"].(int64)
		if ttl > 8 {
			t.Fatalf("TTL too large")
		}

		// Renewing should not have the increment increase since we've hit the
		// max
		req.Operation = logical.UpdateOperation
		req.Path = "auth/token/renew-self"
		req.Data = map[string]interface{}{
			"increment": 1,
		}
		resp, err = core.HandleRequest(req)
		if err != nil {
			t.Fatalf("err: %v %v", err, resp)
		}

		req.Operation = logical.ReadOperation
		req.Path = "auth/token/lookup-self"
		resp, err = core.HandleRequest(req)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		ttl = resp.Data["ttl"].(int64)
		if ttl > 8 {
			t.Fatalf("TTL too large")
		}
	}

	// Now we create a token against the role. We should be able to renew;
	// increment should be ignored as well.
	{
		req.ClientToken = root
		req.Operation = logical.UpdateOperation
		req.Path = "auth/token/create/test"
		resp, err = core.HandleRequest(req)
		if err != nil {
			t.Fatalf("err: %v %v", err, resp)
		}
		if resp == nil {
			t.Fatal("response was nil")
		}
		if resp.Auth == nil {
			t.Fatalf(fmt.Sprintf("response auth was nil, resp is %#v", *resp))
		}
		if resp.Auth.ClientToken == "" {
			t.Fatalf("bad: %#v", resp)
		}

		req.ClientToken = resp.Auth.ClientToken
		req.Operation = logical.ReadOperation
		req.Path = "auth/token/lookup-self"
		resp, err = core.HandleRequest(req)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		ttl := resp.Data["ttl"].(int64)
		if ttl < 299 {
			t.Fatalf("TTL too small (expected %d, got %d", 299, ttl)
		}

		// Let the TTL go down a bit to 3 seconds
		time.Sleep(3 * time.Second)

		req.Operation = logical.UpdateOperation
		req.Path = "auth/token/renew-self"
		req.Data = map[string]interface{}{
			"increment": 1,
		}
		resp, err = core.HandleRequest(req)
		if err != nil {
			t.Fatalf("err: %v %v", err, resp)
		}

		req.Operation = logical.ReadOperation
		req.Path = "auth/token/lookup-self"
		resp, err = core.HandleRequest(req)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		ttl = resp.Data["ttl"].(int64)
		if ttl < 299 {
			t.Fatalf("TTL too small (expected %d, got %d", 299, ttl)
		}
	}
}

func TestTokenStore_RoleExplicitMaxTTL(t *testing.T) {
	core, _, _, root := TestCoreWithTokenStore(t)

	core.defaultLeaseTTL = 5 * time.Second
	core.maxLeaseTTL = 5 * time.Hour

	// Note: these requests are sent to Core since Core handles registration
	// with the expiration manager and we need the storage to be consistent

	// Make sure we can't make it larger than the system/mount max; we should get a warning on role write and an error on token creation
	req := logical.TestRequest(t, logical.UpdateOperation, "auth/token/roles/test")
	req.ClientToken = root
	req.Data = map[string]interface{}{
		"explicit_max_ttl": "100h",
	}

	resp, err := core.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp == nil {
		t.Fatalf("expected a warning")
	}

	req.Operation = logical.UpdateOperation
	req.Path = "auth/token/create/test"
	resp, err = core.HandleRequest(req)
	if err != nil {
		t.Fatalf("expected an error")
	}
	if len(resp.Warnings) == 0 {
		t.Fatalf("expected a warning")
	}

	// Reset to a good explicit max
	req = logical.TestRequest(t, logical.UpdateOperation, "auth/token/roles/test")
	req.ClientToken = root
	req.Data = map[string]interface{}{
		"explicit_max_ttl": "10s",
	}

	resp, err = core.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp != nil {
		t.Fatalf("expected a nil response")
	}

	// This first set of logic is to verify that a normal non-root token will
	// be given a TTL of 5 seconds, and that renewing will cause the TTL to
	// increase
	{
		req.Path = "auth/token/create"
		req.Data = map[string]interface{}{
			"policies": []string{"default"},
		}
		resp, err = core.HandleRequest(req)
		if err != nil {
			t.Fatalf("err: %v %v", err, resp)
		}
		if resp.Auth.ClientToken == "" {
			t.Fatalf("bad: %#v", resp)
		}

		req.ClientToken = resp.Auth.ClientToken
		req.Operation = logical.ReadOperation
		req.Path = "auth/token/lookup-self"
		resp, err = core.HandleRequest(req)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		ttl := resp.Data["ttl"].(int64)
		if ttl > 5 {
			t.Fatalf("TTL too large")
		}

		// Let the TTL go down a bit to 3 seconds
		time.Sleep(2 * time.Second)

		req.Operation = logical.UpdateOperation
		req.Path = "auth/token/renew-self"
		resp, err = core.HandleRequest(req)
		if err != nil {
			t.Fatalf("err: %v %v", err, resp)
		}

		req.Operation = logical.ReadOperation
		req.Path = "auth/token/lookup-self"
		resp, err = core.HandleRequest(req)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		ttl = resp.Data["ttl"].(int64)
		if ttl < 4 {
			t.Fatalf("TTL too small after renewal")
		}
	}

	// Now we create a token against the role. After renew our max should still
	// be the same.
	{
		req.ClientToken = root
		req.Operation = logical.UpdateOperation
		req.Path = "auth/token/create/test"
		resp, err = core.HandleRequest(req)
		if err != nil {
			t.Fatalf("err: %v %v", err, resp)
		}
		if resp == nil {
			t.Fatal("response was nil")
		}
		if resp.Auth == nil {
			t.Fatalf(fmt.Sprintf("response auth was nil, resp is %#v", *resp))
		}
		if resp.Auth.ClientToken == "" {
			t.Fatalf("bad: %#v", resp)
		}

		req.ClientToken = resp.Auth.ClientToken
		req.Operation = logical.ReadOperation
		req.Path = "auth/token/lookup-self"
		resp, err = core.HandleRequest(req)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		ttl := resp.Data["ttl"].(int64)
		if ttl > 10 {
			t.Fatalf("TTL too big")
		}
		maxTTL := resp.Data["explicit_max_ttl"].(int64)
		if maxTTL != 10 {
			t.Fatalf("expected 6 for explicit max TTL, got %d", maxTTL)
		}

		// Let the TTL go down a bit to ~7 seconds (8 against explicit max)
		time.Sleep(2 * time.Second)

		req.Operation = logical.UpdateOperation
		req.Path = "auth/token/renew-self"
		req.Data = map[string]interface{}{
			"increment": 300,
		}
		resp, err = core.HandleRequest(req)
		if err != nil {
			t.Fatalf("err: %v %v", err, resp)
		}

		req.Operation = logical.ReadOperation
		req.Path = "auth/token/lookup-self"
		resp, err = core.HandleRequest(req)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		ttl = resp.Data["ttl"].(int64)
		if ttl > 8 {
			t.Fatalf("TTL too big")
		}

		// Let the TTL go down a bit more to ~5 seconds (6 against explicit max)
		time.Sleep(2 * time.Second)

		req.Operation = logical.UpdateOperation
		req.Path = "auth/token/renew-self"
		req.Data = map[string]interface{}{
			"increment": 300,
		}
		resp, err = core.HandleRequest(req)
		if err != nil {
			t.Fatalf("err: %v %v", err, resp)
		}

		req.Operation = logical.ReadOperation
		req.Path = "auth/token/lookup-self"
		resp, err = core.HandleRequest(req)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		ttl = resp.Data["ttl"].(int64)
		if ttl > 6 {
			t.Fatalf("TTL too big")
		}

		// It should expire
		time.Sleep(8 * time.Second)

		req.Operation = logical.UpdateOperation
		req.Path = "auth/token/renew-self"
		req.Data = map[string]interface{}{
			"increment": 300,
		}
		resp, err = core.HandleRequest(req)
		if err == nil {
			t.Fatalf("expected error")
		}

		time.Sleep(2 * time.Second)

		req.Operation = logical.ReadOperation
		req.Path = "auth/token/lookup-self"
		resp, err = core.HandleRequest(req)
		if resp != nil && err == nil {
			t.Fatalf("expected error, response is %#v", *resp)
		}
		if err == nil {
			t.Fatalf("expected error")
		}
	}
}

func TestTokenStore_Periodic(t *testing.T) {
	core, _, _, root := TestCoreWithTokenStore(t)

	core.defaultLeaseTTL = 10 * time.Second
	core.maxLeaseTTL = 10 * time.Second

	// Note: these requests are sent to Core since Core handles registration
	// with the expiration manager and we need the storage to be consistent

	req := logical.TestRequest(t, logical.UpdateOperation, "auth/token/roles/test")
	req.ClientToken = root
	req.Data = map[string]interface{}{
		"period": 300,
	}

	resp, err := core.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp != nil {
		t.Fatalf("expected a nil response")
	}

	// First make one directly and verify on renew it uses the period.
	{
		req.ClientToken = root
		req.Operation = logical.UpdateOperation
		req.Path = "auth/token/create"
		resp, err = core.HandleRequest(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp == nil {
			t.Fatal("response was nil")
		}
		if resp.Auth == nil {
			t.Fatalf(fmt.Sprintf("response auth was nil, resp is %#v", *resp))
		}
		if resp.Auth.ClientToken == "" {
			t.Fatalf("bad: %#v", resp)
		}

		req.ClientToken = resp.Auth.ClientToken
		req.Operation = logical.ReadOperation
		req.Path = "auth/token/lookup-self"
		resp, err = core.HandleRequest(req)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		ttl := resp.Data["ttl"].(int64)
		if ttl < 299 {
			t.Fatalf("TTL too small (expected %d, got %d)", 299, ttl)
		}

		// Let the TTL go down a bit
		time.Sleep(2 * time.Second)

		req.Operation = logical.UpdateOperation
		req.Path = "auth/token/renew-self"
		req.Data = map[string]interface{}{
			"increment": 1,
		}
		resp, err = core.HandleRequest(req)
		if err != nil {
			t.Fatalf("err: %v %v", err, resp)
		}

		req.Operation = logical.ReadOperation
		req.Path = "auth/token/lookup-self"
		resp, err = core.HandleRequest(req)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		ttl = resp.Data["ttl"].(int64)
		if ttl < 299 {
			t.Fatalf("TTL too small (expected %d, got %d)", 299, ttl)
		}
	}

	// Do the same with an explicit max TTL
	{
		req.ClientToken = root
		req.Operation = logical.UpdateOperation
		req.Path = "auth/token/create"
		req.Data = map[string]interface{}{
			"period":           300,
			"explicit_max_ttl": 150,
		}
		resp, err = core.HandleRequest(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp == nil {
			t.Fatal("response was nil")
		}
		if resp.Auth == nil {
			t.Fatalf(fmt.Sprintf("response auth was nil, resp is %#v", *resp))
		}
		if resp.Auth.ClientToken == "" {
			t.Fatalf("bad: %#v", resp)
		}

		req.ClientToken = resp.Auth.ClientToken
		req.Operation = logical.ReadOperation
		req.Path = "auth/token/lookup-self"
		resp, err = core.HandleRequest(req)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		ttl := resp.Data["ttl"].(int64)
		if ttl < 149 || ttl > 150 {
			t.Fatalf("TTL bad (expected %d, got %d)", 149, ttl)
		}

		// Let the TTL go down a bit
		time.Sleep(2 * time.Second)

		req.Operation = logical.UpdateOperation
		req.Path = "auth/token/renew-self"
		req.Data = map[string]interface{}{
			"increment": 76,
		}
		resp, err = core.HandleRequest(req)
		if err != nil {
			t.Fatalf("err: %v %v", err, resp)
		}

		req.Operation = logical.ReadOperation
		req.Path = "auth/token/lookup-self"
		resp, err = core.HandleRequest(req)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		ttl = resp.Data["ttl"].(int64)
		if ttl < 140 || ttl > 150 {
			t.Fatalf("TTL bad (expected around %d, got %d)", 145, ttl)
		}
	}

	// Now we create a token against the role and also set the te value
	// directly. We should use the smaller of the two and be able to renew;
	// increment should be ignored as well.
	{
		req.ClientToken = root
		req.Operation = logical.UpdateOperation
		req.Path = "auth/token/create/test"
		req.Data = map[string]interface{}{
			"period": 150,
		}
		resp, err = core.HandleRequest(req)
		if err != nil {
			t.Fatalf("err: %v %v", err, resp)
		}
		if resp == nil {
			t.Fatal("response was nil")
		}
		if resp.Auth == nil {
			t.Fatalf(fmt.Sprintf("response auth was nil, resp is %#v", *resp))
		}
		if resp.Auth.ClientToken == "" {
			t.Fatalf("bad: %#v", resp)
		}

		req.ClientToken = resp.Auth.ClientToken
		req.Operation = logical.ReadOperation
		req.Path = "auth/token/lookup-self"
		resp, err = core.HandleRequest(req)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		ttl := resp.Data["ttl"].(int64)
		if ttl < 149 || ttl > 150 {
			t.Fatalf("TTL bad (expected %d, got %d)", 149, ttl)
		}

		// Let the TTL go down a bit
		time.Sleep(2 * time.Second)

		req.Operation = logical.UpdateOperation
		req.Path = "auth/token/renew-self"
		req.Data = map[string]interface{}{
			"increment": 1,
		}
		resp, err = core.HandleRequest(req)
		if err != nil {
			t.Fatalf("err: %v %v", err, resp)
		}

		req.Operation = logical.ReadOperation
		req.Path = "auth/token/lookup-self"
		resp, err = core.HandleRequest(req)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		ttl = resp.Data["ttl"].(int64)
		if ttl < 149 {
			t.Fatalf("TTL bad (expected %d, got %d)", 149, ttl)
		}
	}

	// Now do the same, also using an explicit max in the role
	{
		req.Path = "auth/token/roles/test"
		req.ClientToken = root
		req.Data = map[string]interface{}{
			"period":           300,
			"explicit_max_ttl": 150,
		}

		req.ClientToken = root
		req.Operation = logical.UpdateOperation
		req.Path = "auth/token/create/test"
		req.Data = map[string]interface{}{
			"period":           150,
			"explicit_max_ttl": 130,
		}
		resp, err = core.HandleRequest(req)
		if err != nil {
			t.Fatalf("err: %v %v", err, resp)
		}
		if resp == nil {
			t.Fatal("response was nil")
		}
		if resp.Auth == nil {
			t.Fatalf(fmt.Sprintf("response auth was nil, resp is %#v", *resp))
		}
		if resp.Auth.ClientToken == "" {
			t.Fatalf("bad: %#v", resp)
		}

		req.ClientToken = resp.Auth.ClientToken
		req.Operation = logical.ReadOperation
		req.Path = "auth/token/lookup-self"
		resp, err = core.HandleRequest(req)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		ttl := resp.Data["ttl"].(int64)
		if ttl < 129 || ttl > 130 {
			t.Fatalf("TTL bad (expected %d, got %d)", 129, ttl)
		}

		// Let the TTL go down a bit
		time.Sleep(4 * time.Second)

		req.Operation = logical.UpdateOperation
		req.Path = "auth/token/renew-self"
		req.Data = map[string]interface{}{
			"increment": 1,
		}
		resp, err = core.HandleRequest(req)
		if err != nil {
			t.Fatalf("err: %v %v", err, resp)
		}

		req.Operation = logical.ReadOperation
		req.Path = "auth/token/lookup-self"
		resp, err = core.HandleRequest(req)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		ttl = resp.Data["ttl"].(int64)
		if ttl > 127 {
			t.Fatalf("TTL bad (expected < %d, got %d)", 128, ttl)
		}
	}
}

func TestTokenStore_NoDefaultPolicy(t *testing.T) {
	var resp *logical.Response
	var err error

	core, ts, _, root := TestCoreWithTokenStore(t)
	ps := core.policyStore
	policy, _ := ParseACLPolicy(tokenCreationPolicy)
	policy.Name = "policy1"
	if err := ps.SetPolicy(policy); err != nil {
		t.Fatal(err)
	}

	// Root token creates a token with desired policy. The created token
	// should also have 'default' attached to it.
	tokenData := map[string]interface{}{
		"policies": []string{"policy1"},
	}
	tokenReq := &logical.Request{
		Path:        "create",
		ClientToken: root,
		Operation:   logical.UpdateOperation,
		Data:        tokenData,
	}
	resp, err = ts.HandleRequest(tokenReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v, resp: %v", err, resp)
	}
	if !reflect.DeepEqual(resp.Auth.Policies, []string{"default", "policy1"}) {
		t.Fatalf("bad: policies: expected: [policy, default]; actual: %s", resp.Auth.Policies)
	}

	newToken := resp.Auth.ClientToken

	// Root token creates a token with desired policy, but also requests
	// that the token to not have 'default' policy. The resulting token
	// should not have 'default' policy on it.
	tokenData["no_default_policy"] = true
	resp, err = ts.HandleRequest(tokenReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v, resp: %v", err, resp)
	}

	if !reflect.DeepEqual(resp.Auth.Policies, []string{"policy1"}) {
		t.Fatalf("bad: policies: expected: [policy1]; actual: %s", resp.Auth.Policies)
	}

	// A non-root token which has 'default' policy attached requests for a
	// child token. Child token should also have 'default' policy attached.
	tokenReq.ClientToken = newToken
	tokenReq.Data = nil
	resp, err = ts.HandleRequest(tokenReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v, resp: %v", err, resp)
	}

	if !reflect.DeepEqual(resp.Auth.Policies, []string{"default", "policy1"}) {
		t.Fatalf("bad: policies: expected: [default policy1]; actual: %s", resp.Auth.Policies)
	}

	// A non-root token which has 'default' policy attached, request for a
	// child token to not have 'default' policy while not sending a list
	tokenReq.Data = map[string]interface{}{
		"no_default_policy": true,
	}
	resp, err = ts.HandleRequest(tokenReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v, resp: %v", err, resp)
	}

	if !reflect.DeepEqual(resp.Auth.Policies, []string{"policy1"}) {
		t.Fatalf("bad: policies: expected: [policy1]; actual: %s", resp.Auth.Policies)
	}

	// In this case "default" shouldn't exist because we are not inheriting
	// parent policies
	tokenReq.Data = map[string]interface{}{
		"policies":          []string{"policy1"},
		"no_default_policy": true,
	}
	resp, err = ts.HandleRequest(tokenReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v, resp: %v", err, resp)
	}

	if !reflect.DeepEqual(resp.Auth.Policies, []string{"policy1"}) {
		t.Fatalf("bad: policies: expected: [policy1]; actual: %s", resp.Auth.Policies)
	}

	// This is a non-root token which does not have 'default' policy
	// attached
	newToken = resp.Auth.ClientToken
	tokenReq.Data = nil
	tokenReq.ClientToken = newToken
	resp, err = ts.HandleRequest(tokenReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v, resp: %v", err, resp)
	}

	if !reflect.DeepEqual(resp.Auth.Policies, []string{"policy1"}) {
		t.Fatalf("bad: policies: expected: [policy1]; actual: %s", resp.Auth.Policies)
	}

	roleReq := &logical.Request{
		ClientToken: root,
		Path:        "roles/role1",
		Operation:   logical.CreateOperation,
	}
	resp, err = ts.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v, resp: %v", err, resp)
	}
	tokenReq.Path = "create/role1"
	tokenReq.Data = map[string]interface{}{
		"policies": []string{"policy1"},
	}
	resp, err = ts.HandleRequest(tokenReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v, resp: %v", err, resp)
	}
	if !reflect.DeepEqual(resp.Auth.Policies, []string{"policy1"}) {
		t.Fatalf("bad: policies: expected: [policy1]; actual: %s", resp.Auth.Policies)
	}

	// If 'allowed_policies' in role does not have 'default' in it, the
	// tokens generated using that role should still have the 'default' policy
	// attached to them.
	roleReq.Operation = logical.UpdateOperation
	roleReq.Data = map[string]interface{}{
		"allowed_policies": "policy1",
	}
	resp, err = ts.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v, resp: %v", err, resp)
	}

	resp, err = ts.HandleRequest(tokenReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v, resp: %v", err, resp)
	}
	if !reflect.DeepEqual(resp.Auth.Policies, []string{"default", "policy1"}) {
		t.Fatalf("bad: policies: expected: [default policy1]; actual: %s", resp.Auth.Policies)
	}

	// If 'allowed_policies' in role does not have 'default' in it, the
	// tokens generated using that role should not have 'default' policy
	// attached to them if disallowed_policies contains "default"
	roleReq.Operation = logical.UpdateOperation
	roleReq.Data = map[string]interface{}{
		"allowed_policies":    "policy1",
		"disallowed_policies": "default",
	}
	resp, err = ts.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v, resp: %v", err, resp)
	}

	resp, err = ts.HandleRequest(tokenReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v, resp: %v", err, resp)
	}
	if !reflect.DeepEqual(resp.Auth.Policies, []string{"policy1"}) {
		t.Fatalf("bad: policies: expected: [policy1]; actual: %s", resp.Auth.Policies)
	}

	roleReq.Data = map[string]interface{}{
		"allowed_policies":    "",
		"disallowed_policies": "default",
	}
	resp, err = ts.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v, resp: %v", err, resp)
	}

	resp, err = ts.HandleRequest(tokenReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v, resp: %v", err, resp)
	}
	if !reflect.DeepEqual(resp.Auth.Policies, []string{"policy1"}) {
		t.Fatalf("bad: policies: expected: [policy1]; actual: %s", resp.Auth.Policies)
	}

	// Ensure that if default is in both allowed and disallowed, disallowed wins
	roleReq.Data = map[string]interface{}{
		"allowed_policies":    "default",
		"disallowed_policies": "default",
	}
	resp, err = ts.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v, resp: %v", err, resp)
	}

	delete(tokenReq.Data, "policies")
	resp, err = ts.HandleRequest(tokenReq)
	if err == nil || (resp != nil && !resp.IsError()) {
		t.Fatalf("err: %v, resp: %v", err, resp)
	}
}

func TestTokenStore_AllowedDisallowedPolicies(t *testing.T) {
	var resp *logical.Response
	var err error

	_, ts, _, root := TestCoreWithTokenStore(t)

	roleReq := &logical.Request{
		ClientToken: root,
		Path:        "roles/role1",
		Operation:   logical.CreateOperation,
		Data: map[string]interface{}{
			"allowed_policies":    "allowed1,allowed2",
			"disallowed_policies": "disallowed1,disallowed2",
		},
	}
	resp, err = ts.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v, resp: %v", err, resp)
	}

	tokenReq := &logical.Request{
		Path:        "create/role1",
		ClientToken: root,
		Operation:   logical.UpdateOperation,
		Data: map[string]interface{}{
			"policies": []string{"allowed1"},
		},
	}
	resp, err = ts.HandleRequest(tokenReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%v", err, resp)
	}

	expected := []string{"allowed1", "default"}
	if !reflect.DeepEqual(resp.Auth.Policies, expected) {
		t.Fatalf("bad: expected:%#v actual:%#v", expected, resp.Auth.Policies)
	}

	// Try again with automatic default adding turned off
	tokenReq = &logical.Request{
		Path:        "create/role1",
		ClientToken: root,
		Operation:   logical.UpdateOperation,
		Data: map[string]interface{}{
			"policies":          []string{"allowed1"},
			"no_default_policy": true,
		},
	}
	resp, err = ts.HandleRequest(tokenReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%v", err, resp)
	}

	expected = []string{"allowed1"}
	if !reflect.DeepEqual(resp.Auth.Policies, expected) {
		t.Fatalf("bad: expected:%#v actual:%#v", expected, resp.Auth.Policies)
	}

	tokenReq.Data = map[string]interface{}{
		"policies": []string{"disallowed1"},
	}
	resp, err = ts.HandleRequest(tokenReq)
	if err == nil {
		t.Fatalf("expected an error")
	}

	roleReq.Operation = logical.UpdateOperation
	roleReq.Data = map[string]interface{}{
		"allowed_policies":    "allowed1,common",
		"disallowed_policies": "disallowed1,common",
	}
	resp, err = ts.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v, resp: %v", err, resp)
	}

	tokenReq.Data = map[string]interface{}{
		"policies": []string{"allowed1", "common"},
	}
	resp, err = ts.HandleRequest(tokenReq)
	if err == nil {
		t.Fatalf("expected an error")
	}
}

// Issue 2189
func TestTokenStore_RevokeUseCountToken(t *testing.T) {
	var resp *logical.Response
	var err error
	cubbyFuncLock := &sync.RWMutex{}
	cubbyFuncLock.Lock()

	exp := mockExpiration(t)
	ts := exp.tokenStore
	root, _ := exp.tokenStore.rootToken()

	tokenReq := &logical.Request{
		Path:        "create",
		ClientToken: root.ID,
		Operation:   logical.UpdateOperation,
		Data: map[string]interface{}{
			"num_uses": 1,
		},
	}
	resp, err = ts.HandleRequest(tokenReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%v", err, resp)
	}

	tut := resp.Auth.ClientToken
	saltTut, err := ts.SaltID(tut)
	if err != nil {
		t.Fatal(err)
	}
	te, err := ts.lookupSalted(saltTut, false)
	if err != nil {
		t.Fatal(err)
	}
	if te == nil {
		t.Fatal("nil entry")
	}
	if te.NumUses != 1 {
		t.Fatalf("bad: %d", te.NumUses)
	}

	te, err = ts.UseToken(te)
	if err != nil {
		t.Fatal(err)
	}
	if te == nil {
		t.Fatal("nil entry")
	}
	if te.NumUses != tokenRevocationDeferred {
		t.Fatalf("bad: %d", te.NumUses)
	}

	// Should return no entry because it's tainted
	te, err = ts.lookupSalted(saltTut, false)
	if err != nil {
		t.Fatal(err)
	}
	if te != nil {
		t.Fatalf("%#v", te)
	}

	// But it should show up in an API lookup call
	req := &logical.Request{
		Path:        "lookup-self",
		ClientToken: tut,
		Operation:   logical.UpdateOperation,
	}
	resp, err = ts.HandleRequest(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.Data == nil || resp.Data["num_uses"] == nil {
		t.Fatal("nil resp or data")
	}
	if resp.Data["num_uses"].(int) != -1 {
		t.Fatalf("bad: %v", resp.Data["num_uses"])
	}

	// Should return tainted entries
	te, err = ts.lookupSalted(saltTut, true)
	if err != nil {
		t.Fatal(err)
	}
	if te == nil {
		t.Fatal("nil entry")
	}
	if te.NumUses != tokenRevocationDeferred {
		t.Fatalf("bad: %d", te.NumUses)
	}

	origDestroyCubbyhole := ts.cubbyholeDestroyer

	ts.cubbyholeDestroyer = func(*TokenStore, string) error {
		return fmt.Errorf("keep it frosty")
	}

	err = ts.revokeSalted(saltTut)
	if err == nil {
		t.Fatalf("expected err")
	}

	// Since revocation failed we should see the tokenRevocationFailed canary value
	te, err = ts.lookupSalted(saltTut, true)
	if err != nil {
		t.Fatal(err)
	}
	if te == nil {
		t.Fatal("nil entry")
	}
	if te.NumUses != tokenRevocationFailed {
		t.Fatalf("bad: %d", te.NumUses)
	}

	// Check the race condition situation by making the process sleep
	ts.cubbyholeDestroyer = func(*TokenStore, string) error {
		time.Sleep(1 * time.Second)
		return fmt.Errorf("keep it frosty")
	}
	cubbyFuncLock.Unlock()

	go func() {
		cubbyFuncLock.RLock()
		err := ts.revokeSalted(saltTut)
		cubbyFuncLock.RUnlock()
		if err == nil {
			t.Fatalf("expected error")
		}
	}()

	// Give time for the function to start and grab locks
	time.Sleep(200 * time.Millisecond)
	te, err = ts.lookupSalted(saltTut, true)
	if err != nil {
		t.Fatal(err)
	}
	if te == nil {
		t.Fatal("nil entry")
	}
	if te.NumUses != tokenRevocationInProgress {
		t.Fatalf("bad: %d", te.NumUses)
	}

	// Let things catch up
	time.Sleep(2 * time.Second)

	// Put back to normal
	cubbyFuncLock.Lock()
	defer cubbyFuncLock.Unlock()
	ts.cubbyholeDestroyer = origDestroyCubbyhole

	err = ts.revokeSalted(saltTut)
	if err != nil {
		t.Fatal(err)
	}

	te, err = ts.lookupSalted(saltTut, true)
	if err != nil {
		t.Fatal(err)
	}
	if te != nil {
		t.Fatal("found entry")
	}
}

// Create a token, delete the token entry while leaking accessors, invoke tidy
// and check if the dangling accessor entry is getting removed
func TestTokenStore_HandleTidyCase1(t *testing.T) {
	var resp *logical.Response
	var err error

	_, ts, _, root := TestCoreWithTokenStore(t)

	// List the number of accessors. Since there is only root token
	// present, the list operation should return only one key.
	accessorListReq := &logical.Request{
		Operation:   logical.ListOperation,
		Path:        "accessors",
		ClientToken: root,
	}
	resp, err = ts.HandleRequest(accessorListReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%v", err, resp)
	}

	numberOfAccessors := len(resp.Data["keys"].([]string))
	if numberOfAccessors != 1 {
		t.Fatalf("bad: number of accessors. Expected: 1, Actual: %d", numberOfAccessors)
	}

	for i := 1; i <= 100; i++ {
		// Create a regular token
		tokenReq := &logical.Request{
			Operation:   logical.UpdateOperation,
			Path:        "create",
			ClientToken: root,
			Data: map[string]interface{}{
				"policies": []string{"policy1"},
			},
		}
		resp, err = ts.HandleRequest(tokenReq)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("err:%v resp:%v", err, resp)
		}
		tut := resp.Auth.ClientToken

		// Creation of another token should end up with incrementing
		// the number of accessors
		// the storage
		resp, err = ts.HandleRequest(accessorListReq)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("err:%v resp:%v", err, resp)
		}

		numberOfAccessors = len(resp.Data["keys"].([]string))
		if numberOfAccessors != i+1 {
			t.Fatalf("bad: number of accessors. Expected: %d, Actual: %d", i+1, numberOfAccessors)
		}

		// Revoke the token while leaking other items associated with the
		// token. Do this by doing what revokeSalted used to do before it was
		// fixed, i.e., by deleting the storage entry for token and its
		// cubbyhole and by not deleting its secondary index, its accessor and
		// associated leases.

		saltedTut, err := ts.SaltID(tut)
		if err != nil {
			t.Fatal(err)
		}
		_, err = ts.lookupSalted(saltedTut, true)
		if err != nil {
			t.Fatalf("failed to lookup token: %v", err)
		}

		// Destroy the token index
		path := lookupPrefix + saltedTut
		if ts.view.Delete(path); err != nil {
			t.Fatalf("failed to delete token entry: %v", err)
		}

		// Destroy the cubby space
		err = ts.destroyCubbyhole(saltedTut)
		if err != nil {
			t.Fatalf("failed to destroyCubbyhole: %v", err)
		}

		// Leaking of accessor should have resulted in no change to the number
		// of accessors
		resp, err = ts.HandleRequest(accessorListReq)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("err:%v resp:%v", err, resp)
		}

		numberOfAccessors = len(resp.Data["keys"].([]string))
		if numberOfAccessors != i+1 {
			t.Fatalf("bad: number of accessors. Expected: %d, Actual: %d", i+1, numberOfAccessors)
		}
	}

	tidyReq := &logical.Request{
		Path:        "tidy",
		Operation:   logical.UpdateOperation,
		ClientToken: root,
	}
	resp, err = ts.HandleRequest(tidyReq)
	if err != nil {
		t.Fatal(err)
	}
	if resp != nil && resp.IsError() {
		t.Fatalf("resp: %#v", resp)
	}
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%v", err, resp)
	}

	// Tidy should have removed all the dangling accessor entries
	resp, err = ts.HandleRequest(accessorListReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%v", err, resp)
	}

	numberOfAccessors = len(resp.Data["keys"].([]string))
	if numberOfAccessors != 1 {
		t.Fatalf("bad: number of accessors. Expected: 1, Actual: %d", numberOfAccessors)
	}
}

func TestTokenStore_TidyLeaseRevocation(t *testing.T) {
	exp := mockExpiration(t)
	ts := exp.tokenStore

	noop := &NoopBackend{}
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "logical/")
	meUUID, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	err = exp.router.Mount(noop, "prod/aws/", &MountEntry{UUID: meUUID, Accessor: "awsaccessor"}, view)
	if err != nil {
		t.Fatal(err)
	}

	// Create new token
	root, err := ts.rootToken()
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	req := logical.TestRequest(t, logical.UpdateOperation, "create")
	req.ClientToken = root.ID
	req.Data["policies"] = []string{"default"}

	resp, err := ts.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}

	// Create a new token
	auth := &logical.Auth{
		ClientToken: resp.Auth.ClientToken,
		LeaseOptions: logical.LeaseOptions{
			TTL:       time.Hour,
			Renewable: true,
		},
	}
	err = exp.RegisterAuth("auth/token/create", auth)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	tut := resp.Auth.ClientToken

	req = &logical.Request{
		Path:        "prod/aws/foo",
		ClientToken: tut,
	}
	resp = &logical.Response{
		Secret: &logical.Secret{
			LeaseOptions: logical.LeaseOptions{
				TTL: time.Hour,
			},
		},
	}

	leases := []string{}

	for i := 0; i < 10; i++ {
		leaseId, err := exp.Register(req, resp)
		if err != nil {
			t.Fatal(err)
		}
		leases = append(leases, leaseId)
	}

	sort.Strings(leases)

	storedLeases, err := exp.lookupByToken(tut)
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(storedLeases)
	if !reflect.DeepEqual(leases, storedLeases) {
		t.Fatalf("bad: %#v vs %#v", leases, storedLeases)
	}

	// Now, delete the token entry. The leases should still exist.
	saltedTut, err := ts.SaltID(tut)
	if err != nil {
		t.Fatal(err)
	}
	te, err := ts.lookupSalted(saltedTut, true)
	if err != nil {
		t.Fatalf("failed to lookup token: %v", err)
	}
	if te == nil {
		t.Fatal("got nil token entry")
	}

	// Destroy the token index
	path := lookupPrefix + saltedTut
	if ts.view.Delete(path); err != nil {
		t.Fatalf("failed to delete token entry: %v", err)
	}
	te, err = ts.lookupSalted(saltedTut, true)
	if err != nil {
		t.Fatalf("failed to lookup token: %v", err)
	}
	if te != nil {
		t.Fatal("got token entry")
	}

	// Verify leases still exist
	storedLeases, err = exp.lookupByToken(tut)
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(storedLeases)
	if !reflect.DeepEqual(leases, storedLeases) {
		t.Fatalf("bad: %#v vs %#v", leases, storedLeases)
	}

	// Call tidy
	ts.handleTidy(nil, nil)

	// Verify leases are gone
	storedLeases, err = exp.lookupByToken(tut)
	if err != nil {
		t.Fatal(err)
	}
	if len(storedLeases) > 0 {
		t.Fatal("found leases")
	}
}
