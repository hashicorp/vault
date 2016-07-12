package vault

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/logical"
)

func getBackendConfig(c *Core) *logical.BackendConfig {
	return &logical.BackendConfig{
		Logger: c.logger,
		System: logical.StaticSystemView{
			DefaultLeaseTTLVal: time.Hour * 24,
			MaxLeaseTTLVal:     time.Hour * 24 * 30,
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
	if resp.Auth.ClientToken != client {
		t.Fatalf("bad: %#v", resp)
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

	token, err := ts.lookupByAccessor(out.Accessor)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Verify that the value returned from the index matches the token ID
	if token != ent.ID {
		t.Fatalf("bad: got\n%s\nexpected\n%s\n", token, ent.ID)
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

	req := logical.TestRequest(t, logical.UpdateOperation, "lookup-accessor/"+out.Accessor)

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

	req := logical.TestRequest(t, logical.UpdateOperation, "revoke-accessor/"+out.Accessor)

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

	// Should still match
	out, err = ts2.Lookup(ent.ID)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !reflect.DeepEqual(out, ent) {
		t.Fatalf("bad: expected:%#v\nactual:%#v", ent, out)
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
	if te.NumUses != -1 {
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
	_, ts, _, _ := TestCoreWithTokenStore(t)

	// Mount a noop backend
	noop := &NoopBackend{}
	ts.expiration.router.Mount(noop, "", &MountEntry{UUID: ""}, nil)

	ent := &TokenEntry{Path: "test", Policies: []string{"dev", "ops"}}
	if err := ts.create(ent); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Register a lease
	req := &logical.Request{
		Operation:   logical.ReadOperation,
		Path:        "secret/foo",
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
	if err.Error() != "cannot revoke blank token" {
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

	req := logical.TestRequest(t, logical.UpdateOperation, "revoke/child")
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

	req := logical.TestRequest(t, logical.UpdateOperation, "revoke-orphan/child")
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

	req := logical.TestRequest(t, logical.UpdateOperation, "revoke-orphan/child")
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
	req := logical.TestRequest(t, logical.ReadOperation, "lookup/"+root)
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
		"role":             "",
		"explicit_max_ttl": int64(0),
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
	req = logical.TestRequest(t, logical.ReadOperation, "lookup/client")
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
		"role":             "",
		"explicit_max_ttl": int64(0),
		"renewable":        true,
	}

	if resp.Data["creation_time"].(int64) == 0 {
		t.Fatalf("creation time was zero")
	}
	delete(resp.Data, "creation_time")

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
		"role":             "",
		"explicit_max_ttl": int64(0),
		"renewable":        true,
	}

	if resp.Data["creation_time"].(int64) == 0 {
		t.Fatalf("creation time was zero")
	}
	delete(resp.Data, "creation_time")

	// Depending on timing of the test this may have ticked down, so accept 3599
	if resp.Data["ttl"].(int64) == 3599 {
		resp.Data["ttl"] = int64(3600)
	}

	if !reflect.DeepEqual(resp.Data, exp) {
		t.Fatalf("bad: expected:%#v\nactual:%#v", exp, resp.Data)
	}

	// Test last_renewal_time functionality
	req = logical.TestRequest(t, logical.UpdateOperation, "renew/client")
	resp, err = ts.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp == nil {
		t.Fatalf("bad: %#v", resp)
	}

	req = logical.TestRequest(t, logical.ReadOperation, "lookup/client")
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
	_, ts, _, root := TestCoreWithTokenStore(t)
	req := logical.TestRequest(t, logical.ReadOperation, "lookup-self")
	req.ClientToken = root
	resp, err := ts.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp == nil {
		t.Fatalf("bad: %#v", resp)
	}

	exp := map[string]interface{}{
		"id":               root,
		"accessor":         resp.Data["accessor"],
		"policies":         []string{"root"},
		"path":             "auth/token/root",
		"meta":             map[string]string(nil),
		"display_name":     "root",
		"orphan":           true,
		"num_uses":         0,
		"creation_ttl":     int64(0),
		"ttl":              int64(0),
		"role":             "",
		"explicit_max_ttl": int64(0),
	}

	if resp.Data["creation_time"].(int64) == 0 {
		t.Fatalf("creation time was zero")
	}
	delete(resp.Data, "creation_time")

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
	req := logical.TestRequest(t, logical.UpdateOperation, "renew/"+root.ID)
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
		"name":             "test",
		"orphan":           true,
		"period":           int64(259200),
		"allowed_policies": []string{"default", "test1", "test2"},
		"path_suffix":      "happenin",
		"explicit_max_ttl": int64(0),
		"renewable":        true,
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
		"name":             "test",
		"orphan":           true,
		"period":           int64(284400),
		"allowed_policies": []string{"default", "test3"},
		"path_suffix":      "happenin",
		"explicit_max_ttl": int64(0),
		"renewable":        false,
	}

	if !reflect.DeepEqual(expected, resp.Data) {
		t.Fatalf("bad: expected:%#v\nactual:%#v", expected, resp.Data)
	}

	// Now test setting explicit max ttl at the same time as period, which
	// should be an error
	req.Operation = logical.CreateOperation
	req.Data = map[string]interface{}{
		"explicit_max_ttl": "5",
	}

	resp, err = core.HandleRequest(req)
	if err == nil {
		t.Fatalf("expected error")
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
		"name":             "test",
		"orphan":           true,
		"explicit_max_ttl": int64(5),
		"allowed_policies": []string{"default", "test3"},
		"path_suffix":      "happenin",
		"period":           int64(0),
		"renewable":        false,
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

func TestTokenStore_RoleAllowedRoles(t *testing.T) {
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
			t.Fatal(fmt.Sprintf("response auth was nil, resp is %#v", *resp))
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
	if len(resp.Warnings()) == 0 {
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
			t.Fatal(fmt.Sprintf("response auth was nil, resp is %#v", *resp))
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

		req.Operation = logical.ReadOperation
		req.Path = "auth/token/lookup-self"
		resp, err = core.HandleRequest(req)
		if err == nil {
			t.Fatalf("expected error")
		}
	}
}
