package vault

import (
	"log"
	"os"
	"reflect"
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

func mockTokenStore(t *testing.T) (*Core, *TokenStore, string) {
	logger := log.New(os.Stderr, "", log.LstdFlags)

	c, _, root := TestCoreUnsealed(t)

	ts, err := NewTokenStore(c, getBackendConfig(c))
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	router := NewRouter()
	router.Mount(ts, "auth/token/", &MountEntry{UUID: ""}, ts.view)

	view := c.systemBarrierView.SubView(expirationSubPath)
	exp := NewExpirationManager(router, view, ts, logger)
	ts.SetExpirationManager(exp)
	return c, ts, root
}

func TestTokenStore_RootToken(t *testing.T) {
	_, ts, _ := mockTokenStore(t)

	te, err := ts.RootToken()
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
		t.Fatalf("bad: %#v", out)
	}
}

func TestTokenStore_CreateLookup(t *testing.T) {
	c, ts, _ := mockTokenStore(t)

	ent := &TokenEntry{Path: "test", Policies: []string{"dev", "ops"}}
	if err := ts.Create(ent); err != nil {
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
		t.Fatalf("bad: %#v", out)
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
		t.Fatalf("bad: %#v", out)
	}
}

func TestTokenStore_CreateLookup_ProvidedID(t *testing.T) {
	c, ts, _ := mockTokenStore(t)

	ent := &TokenEntry{
		ID:       "foobarbaz",
		Path:     "test",
		Policies: []string{"dev", "ops"},
	}
	if err := ts.Create(ent); err != nil {
		t.Fatalf("err: %v", err)
	}
	if ent.ID != "foobarbaz" {
		t.Fatalf("bad: %#v", ent)
	}

	out, err := ts.Lookup(ent.ID)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !reflect.DeepEqual(out, ent) {
		t.Fatalf("bad: %#v", out)
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
		t.Fatalf("bad: %#v", out)
	}
}

func TestTokenStore_UseToken(t *testing.T) {
	_, ts, root := mockTokenStore(t)

	// Lookup the root token
	ent, err := ts.Lookup(root)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Root is an unlimited use token, should be a no-op
	err = ts.UseToken(ent)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Lookup the root token again
	ent2, err := ts.Lookup(root)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if !reflect.DeepEqual(ent, ent2) {
		t.Fatalf("bad: %#v %#v", ent, ent2)
	}

	// Create a retstricted token
	ent = &TokenEntry{Path: "test", Policies: []string{"dev", "ops"}, NumUses: 2}
	if err := ts.Create(ent); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Use the token
	err = ts.UseToken(ent)
	if err != nil {
		t.Fatalf("err: %v", err)
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
	err = ts.UseToken(ent)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

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
	_, ts, _ := mockTokenStore(t)

	ent := &TokenEntry{Path: "test", Policies: []string{"dev", "ops"}}
	if err := ts.Create(ent); err != nil {
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
	_, ts, _ := mockTokenStore(t)

	// Mount a noop backend
	noop := &NoopBackend{}
	ts.expiration.router.Mount(noop, "", &MountEntry{UUID: ""}, nil)

	ent := &TokenEntry{Path: "test", Policies: []string{"dev", "ops"}}
	if err := ts.Create(ent); err != nil {
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
	_, ts, _ := mockTokenStore(t)

	ent := &TokenEntry{Path: "test", Policies: []string{"dev", "ops"}}
	if err := ts.Create(ent); err != nil {
		t.Fatalf("err: %v", err)
	}

	ent2 := &TokenEntry{Parent: ent.ID}
	if err := ts.Create(ent2); err != nil {
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
		t.Fatalf("bad: %#v", out)
	}
}

func TestTokenStore_RevokeTree(t *testing.T) {
	_, ts, _ := mockTokenStore(t)

	ent1 := &TokenEntry{}
	if err := ts.Create(ent1); err != nil {
		t.Fatalf("err: %v", err)
	}

	ent2 := &TokenEntry{Parent: ent1.ID}
	if err := ts.Create(ent2); err != nil {
		t.Fatalf("err: %v", err)
	}

	ent3 := &TokenEntry{Parent: ent2.ID}
	if err := ts.Create(ent3); err != nil {
		t.Fatalf("err: %v", err)
	}

	ent4 := &TokenEntry{Parent: ent2.ID}
	if err := ts.Create(ent4); err != nil {
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

func TestTokenStore_HandleRequest_CreateToken_DisplayName(t *testing.T) {
	_, ts, root := mockTokenStore(t)

	req := logical.TestRequest(t, logical.WriteOperation, "create")
	req.ClientToken = root
	req.Data["display_name"] = "foo_bar.baz!"

	resp, err := ts.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}

	expected := &TokenEntry{
		ID:          resp.Auth.ClientToken,
		Parent:      root,
		Policies:    []string{"root"},
		Path:        "auth/token/create",
		DisplayName: "token-foo-bar-baz",
	}
	out, err := ts.Lookup(resp.Auth.ClientToken)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !reflect.DeepEqual(out, expected) {
		t.Fatalf("bad: %#v", out)
	}
}

func TestTokenStore_HandleRequest_CreateToken_NumUses(t *testing.T) {
	_, ts, root := mockTokenStore(t)

	req := logical.TestRequest(t, logical.WriteOperation, "create")
	req.ClientToken = root
	req.Data["num_uses"] = "1"

	resp, err := ts.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}

	expected := &TokenEntry{
		ID:          resp.Auth.ClientToken,
		Parent:      root,
		Policies:    []string{"root"},
		Path:        "auth/token/create",
		DisplayName: "token",
		NumUses:     1,
	}
	out, err := ts.Lookup(resp.Auth.ClientToken)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !reflect.DeepEqual(out, expected) {
		t.Fatalf("bad: %#v", out)
	}
}

func TestTokenStore_HandleRequest_CreateToken_NumUses_Invalid(t *testing.T) {
	_, ts, root := mockTokenStore(t)

	req := logical.TestRequest(t, logical.WriteOperation, "create")
	req.ClientToken = root
	req.Data["num_uses"] = "-1"

	resp, err := ts.HandleRequest(req)
	if err != logical.ErrInvalidRequest {
		t.Fatalf("err: %v %v", err, resp)
	}
}

func TestTokenStore_HandleRequest_CreateToken_NumUses_Restricted(t *testing.T) {
	_, ts, root := mockTokenStore(t)

	req := logical.TestRequest(t, logical.WriteOperation, "create")
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
	_, ts, root := mockTokenStore(t)

	req := logical.TestRequest(t, logical.WriteOperation, "create")
	req.ClientToken = root

	resp, err := ts.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}

	expected := &TokenEntry{
		ID:          resp.Auth.ClientToken,
		Parent:      root,
		Policies:    []string{"root"},
		Path:        "auth/token/create",
		DisplayName: "token",
	}
	out, err := ts.Lookup(resp.Auth.ClientToken)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !reflect.DeepEqual(out, expected) {
		t.Fatalf("bad: %#v", out)
	}
}

func TestTokenStore_HandleRequest_CreateToken_BadParent(t *testing.T) {
	_, ts, _ := mockTokenStore(t)

	req := logical.TestRequest(t, logical.WriteOperation, "create")
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
	_, ts, root := mockTokenStore(t)

	req := logical.TestRequest(t, logical.WriteOperation, "create")
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
	_, ts, root := mockTokenStore(t)

	req := logical.TestRequest(t, logical.WriteOperation, "create")
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
	_, ts, root := mockTokenStore(t)
	testMakeToken(t, ts, root, "client", []string{"foo"})

	req := logical.TestRequest(t, logical.WriteOperation, "create")
	req.ClientToken = "client"
	req.Data["id"] = "foobar"
	req.Data["policies"] = []string{"foo"}

	resp, err := ts.HandleRequest(req)
	if err != logical.ErrInvalidRequest {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp.Data["error"] != "root required to specify token id" {
		t.Fatalf("bad: %#v", resp)
	}
}

func TestTokenStore_HandleRequest_CreateToken_NonRoot_Subset(t *testing.T) {
	_, ts, root := mockTokenStore(t)
	testMakeToken(t, ts, root, "client", []string{"foo", "bar"})

	req := logical.TestRequest(t, logical.WriteOperation, "create")
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
	_, ts, root := mockTokenStore(t)
	testMakeToken(t, ts, root, "client", []string{"foo", "bar"})

	req := logical.TestRequest(t, logical.WriteOperation, "create")
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
	_, ts, root := mockTokenStore(t)
	testMakeToken(t, ts, root, "client", []string{"foo"})

	req := logical.TestRequest(t, logical.WriteOperation, "create")
	req.ClientToken = "client"
	req.Data["no_parent"] = true
	req.Data["policies"] = []string{"foo"}

	resp, err := ts.HandleRequest(req)
	if err != logical.ErrInvalidRequest {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp.Data["error"] != "root required to create orphan token" {
		t.Fatalf("bad: %#v", resp)
	}
}

func TestTokenStore_HandleRequest_CreateToken_Root_NoParent(t *testing.T) {
	_, ts, root := mockTokenStore(t)

	req := logical.TestRequest(t, logical.WriteOperation, "create")
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

func TestTokenStore_HandleRequest_CreateToken_Metadata(t *testing.T) {
	_, ts, root := mockTokenStore(t)

	req := logical.TestRequest(t, logical.WriteOperation, "create")
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
		t.Fatalf("bad: %#v", out)
	}
}

func TestTokenStore_HandleRequest_CreateToken_Lease(t *testing.T) {
	_, ts, root := mockTokenStore(t)

	req := logical.TestRequest(t, logical.WriteOperation, "create")
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

func TestTokenStore_HandleRequest_Revoke(t *testing.T) {
	_, ts, root := mockTokenStore(t)
	testMakeToken(t, ts, root, "child", []string{"root", "foo"})
	testMakeToken(t, ts, "child", "sub-child", []string{"foo"})

	req := logical.TestRequest(t, logical.WriteOperation, "revoke/child")
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
	_, ts, root := mockTokenStore(t)
	testMakeToken(t, ts, root, "child", []string{"root", "foo"})
	testMakeToken(t, ts, "child", "sub-child", []string{"foo"})

	req := logical.TestRequest(t, logical.WriteOperation, "revoke-orphan/child")
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

func TestTokenStore_HandleRequest_Lookup(t *testing.T) {
	_, ts, root := mockTokenStore(t)
	req := logical.TestRequest(t, logical.ReadOperation, "lookup/"+root)
	resp, err := ts.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp == nil {
		t.Fatalf("bad: %#v", resp)
	}

	exp := map[string]interface{}{
		"id":           root,
		"policies":     []string{"root"},
		"path":         "auth/token/root",
		"meta":         map[string]string(nil),
		"display_name": "root",
		"num_uses":     0,
	}
	if !reflect.DeepEqual(resp.Data, exp) {
		t.Fatalf("bad: %#v exp: %#v", resp.Data, exp)
	}
}

func TestTokenStore_HandleRequest_RevokePrefix(t *testing.T) {
	exp := mockExpiration(t)
	ts := exp.tokenStore

	// Create new token
	root, err := ts.RootToken()
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Create a new token
	auth := &logical.Auth{
		ClientToken: root.ID,
		LeaseOptions: logical.LeaseOptions{
			TTL: time.Hour,
		},
	}
	err = exp.RegisterAuth("auth/github/login", auth)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	req := logical.TestRequest(t, logical.WriteOperation, "revoke-prefix/auth/github/")
	resp, err := ts.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp != nil {
		t.Fatalf("bad: %#v", resp)
	}

	out, err := ts.Lookup(root.ID)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out != nil {
		t.Fatalf("bad: %v", out)
	}
}

func TestTokenStore_HandleRequest_LookupSelf(t *testing.T) {
	_, ts, root := mockTokenStore(t)
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
		"id":           root,
		"policies":     []string{"root"},
		"path":         "auth/token/root",
		"meta":         map[string]string(nil),
		"display_name": "root",
		"num_uses":     0,
	}
	if !reflect.DeepEqual(resp.Data, exp) {
		t.Fatalf("bad: %#v exp: %#v", resp.Data, exp)
	}
}

func TestTokenStore_HandleRequest_Renew(t *testing.T) {
	exp := mockExpiration(t)
	ts := exp.tokenStore

	// Create new token
	root, err := ts.RootToken()
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

	beforeRenew := time.Now().UTC()
	req := logical.TestRequest(t, logical.WriteOperation, "renew/"+root.ID)
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

func testMakeToken(t *testing.T, ts *TokenStore, root, client string, policy []string) {
	req := logical.TestRequest(t, logical.WriteOperation, "create")
	req.ClientToken = root
	req.Data["id"] = client
	req.Data["policies"] = policy

	resp, err := ts.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp.Auth.ClientToken != client {
		t.Fatalf("bad: %#v", resp)
	}
}

func testCoreMakeToken(t *testing.T, c *Core, root, client string, policy []string) {
	req := logical.TestRequest(t, logical.WriteOperation, "auth/token/create")
	req.ClientToken = root
	req.Data["id"] = client
	req.Data["policies"] = policy

	resp, err := c.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp.Auth.ClientToken != client {
		t.Fatalf("bad: %#v", resp)
	}
}
