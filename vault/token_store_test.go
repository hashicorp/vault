// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"reflect"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/go-sockaddr"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/helper/metricsutil"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/locksutil"
	"github.com/hashicorp/vault/sdk/helper/tokenutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/mapstructure"
)

func TestTokenStore_CreateOrphanResponse(t *testing.T) {
	c, _, root := TestCoreUnsealed(t)
	resp, err := c.HandleRequest(namespace.RootContext(nil), &logical.Request{
		Operation:   logical.UpdateOperation,
		Path:        "auth/token/create-orphan",
		ClientToken: root,
		Data: map[string]interface{}{
			"policies": "default",
		},
	})
	if err != nil && (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v, resp: %#v", err, resp)
	}
	if !resp.Auth.Orphan {
		t.Fatalf("failed to set orphan as true in the response")
	}
}

// TestTokenStore_CubbyholeDeletion tests that a token's cubbyhole
// can be used and that the cubbyhole is removed after the token is revoked.
func TestTokenStore_CubbyholeDeletion(t *testing.T) {
	c, _, root := TestCoreUnsealed(t)
	testTokenStore_CubbyholeDeletion(t, c, root)
}

// TestTokenStore_CubbyholeDeletionSSCTokensDisabled tests that a legacy token's
// cubbyhole can be used, and that the cubbyhole is removed after the token is revoked.
func TestTokenStore_CubbyholeDeletionSSCTokensDisabled(t *testing.T) {
	c, _, root := TestCoreUnsealed(t)
	c.disableSSCTokens = true
	testTokenStore_CubbyholeDeletion(t, c, root)
}

func testTokenStore_CubbyholeDeletion(t *testing.T, c *Core, root string) {
	ts := c.tokenStore

	for i := 0; i < 10; i++ {
		// Create a token
		tokenReq := &logical.Request{
			Operation:   logical.UpdateOperation,
			Path:        "create",
			ClientToken: root,
			Data: map[string]interface{}{
				"ttl": "600s",
			},
		}
		// Supplying token ID forces SHA1 hashing to be used
		if i%2 == 0 {
			tokenReq.Data = map[string]interface{}{
				"id":  "testroot",
				"ttl": "600s",
			}
		}
		resp := testMakeTokenViaRequest(t, ts, tokenReq)
		token := resp.Auth.ClientToken

		// Write data in the token's cubbyhole
		resp, err := c.HandleRequest(namespace.RootContext(nil), &logical.Request{
			ClientToken: token,
			Operation:   logical.UpdateOperation,
			Path:        "cubbyhole/sample/data",
			Data: map[string]interface{}{
				"foo": "bar",
			},
		})
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
		}

		// Revoke the token
		resp, err = ts.HandleRequest(namespace.RootContext(nil), &logical.Request{
			ClientToken: token,
			Path:        "revoke-self",
			Operation:   logical.UpdateOperation,
		})
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
		}
	}

	// List the cubbyhole keys
	cubbyholeKeys, err := ts.cubbyholeBackend.storageView.List(namespace.RootContext(nil), "")
	if err != nil {
		t.Fatal(err)
	}

	// There should be no entries
	if len(cubbyholeKeys) != 0 {
		t.Fatalf("bad: len(cubbyholeKeys); expected: 0, actual: %d", len(cubbyholeKeys))
	}
}

func TestTokenStore_CubbyholeTidy(t *testing.T) {
	c, _, root := TestCoreUnsealed(t)
	testTokenStore_CubbyholeTidy(t, c, root, namespace.RootContext(nil))
}

func testTokenStore_CubbyholeTidy(t *testing.T, c *Core, root string, nsCtx context.Context) {
	ts := c.tokenStore

	backend := c.router.MatchingBackend(nsCtx, mountPathCubbyhole)
	view := c.router.MatchingStorageByAPIPath(nsCtx, mountPathCubbyhole)

	for i := 1; i <= 20; i++ {
		// Create 20 tokens
		tokenReq := &logical.Request{
			Operation:   logical.UpdateOperation,
			Path:        "create",
			ClientToken: root,
			Data: map[string]interface{}{
				"ttl": "600s",
			},
		}

		resp := testMakeTokenViaRequestContext(t, nsCtx, ts, tokenReq)
		token := resp.Auth.ClientToken

		// Supplying token ID forces SHA1 hashing to be used
		if i%3 == 0 {
			tokenReq.Data = map[string]interface{}{
				"id":  "testroot",
				"ttl": "600s",
			}
		}

		// Create 4 junk cubbyhole entries
		if i%5 == 0 {
			invalidToken, err := uuid.GenerateUUID()
			if err != nil {
				t.Fatal(err)
			}

			resp, err := backend.HandleRequest(nsCtx, &logical.Request{
				ClientToken: invalidToken,
				Operation:   logical.UpdateOperation,
				Path:        "cubbyhole/sample/data",
				Data: map[string]interface{}{
					"foo": "bar",
				},
				Storage: view,
			})
			if err != nil || (resp != nil && resp.IsError()) {
				t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
			}
		}

		// Write into cubbyholes of 10 tokens
		if i%2 == 0 {
			continue
		}
		resp, err := c.HandleRequest(nsCtx, &logical.Request{
			ClientToken: token,
			Operation:   logical.UpdateOperation,
			Path:        "cubbyhole/sample/data",
			Data: map[string]interface{}{
				"foo": "bar",
			},
		})
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
		}
	}

	// List all the cubbyhole storage keys
	cubbyholeKeys, err := view.List(nsCtx, "")
	if err != nil {
		t.Fatal(err)
	}

	if len(cubbyholeKeys) != 14 {
		t.Fatalf("bad: len(cubbyholeKeys); expected: 14, actual: %d", len(cubbyholeKeys))
	}

	// Tidy cubbyhole storage
	resp, err := ts.HandleRequest(nsCtx, &logical.Request{
		Path:      "tidy",
		Operation: logical.UpdateOperation,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
	}

	// Wait for tidy operation to complete
	time.Sleep(2 * time.Second)

	// List all the cubbyhole storage keys
	cubbyholeKeys, err = view.List(nsCtx, "")
	if err != nil {
		t.Fatal(err)
	}

	// The junk entries must have been cleaned up
	if len(cubbyholeKeys) != 10 {
		t.Fatalf("bad: len(cubbyholeKeys); expected: 10, actual: %d", len(cubbyholeKeys))
	}
}

func TestTokenStore_Salting(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ts := c.tokenStore

	saltedID, err := ts.SaltID(namespace.RootContext(nil), "foo")
	if err != nil {
		t.Fatal(err)
	}
	if strings.HasPrefix(saltedID, "h") {
		t.Fatalf("expected sha1 hash; got sha2-256 hmac")
	}

	saltedID, err = ts.SaltID(namespace.RootContext(nil), "hvs.foo")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(saltedID, "h") {
		t.Fatalf("expected sha2-256 hmac; got sha1 hash")
	}

	nsCtx := namespace.ContextWithNamespace(context.Background(), &namespace.Namespace{ID: "testid", Path: "ns1"})
	saltedID, err = ts.SaltID(nsCtx, "foo")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(saltedID, "h") {
		t.Fatalf("expected sha2-256 hmac; got sha1 hash")
	}

	saltedID, err = ts.SaltID(nsCtx, "hvs.foo")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(saltedID, "h") {
		t.Fatalf("expected sha2-256 hmac; got sha1 hash")
	}
}

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
	c, _, _ := TestCoreUnsealed(t)
	ts := c.tokenStore

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

	saltedID, err := ts.SaltID(namespace.RootContext(nil), entry.ID)
	if err != nil {
		t.Fatal(err)
	}
	le := &logical.StorageEntry{
		Key:   saltedID,
		Value: enc,
	}

	if err := ts.idView(namespace.RootNamespace).Put(namespace.RootContext(nil), le); err != nil {
		t.Fatal(err)
	}

	// Register with exp manager so lookup works
	auth := &logical.Auth{
		DisplayName:    entry.DisplayName,
		CreationPath:   entry.Path,
		Policies:       entry.Policies,
		ExplicitMaxTTL: entry.ExplicitMaxTTL,
		NumUses:        entry.NumUses,
		LeaseOptions: logical.LeaseOptions{
			TTL: time.Hour,
		},
		ClientToken: entry.ID,
	}
	// Same as entry from TokenEntryOld, but used for RegisterAuth
	registryEntry := &logical.TokenEntry{
		DisplayName:    entry.DisplayName,
		Path:           entry.Path,
		Policies:       entry.Policies,
		CreationTime:   entry.CreationTime,
		ExplicitMaxTTL: entry.ExplicitMaxTTL,
		NumUses:        entry.NumUses,
		NamespaceID:    namespace.RootNamespaceID,
	}

	if err := ts.expiration.RegisterAuth(namespace.RootContext(nil), registryEntry, auth, ""); err != nil {
		t.Fatal(err)
	}

	out, err := ts.Lookup(namespace.RootContext(nil), entry.ID)
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
	ent := &logical.TokenEntry{
		DisplayName:    "test-display-name",
		Path:           "test",
		Policies:       []string{"dev", "ops"},
		CreationTime:   time.Now().Unix(),
		ExplicitMaxTTL: 100,
		NumUses:        10,
		NamespaceID:    namespace.RootNamespaceID,
	}
	if err := ts.create(namespace.RootContext(nil), ent); err != nil {
		t.Fatalf("err: %s", err)
	}
	auth = &logical.Auth{
		DisplayName:    ent.DisplayName,
		CreationPath:   ent.Path,
		Policies:       ent.Policies,
		ExplicitMaxTTL: ent.ExplicitMaxTTL,
		NumUses:        ent.NumUses,
		LeaseOptions: logical.LeaseOptions{
			TTL: time.Hour,
		},
		ClientToken: ent.ID,
	}
	if err := ts.expiration.RegisterAuth(namespace.RootContext(nil), ent, auth, ""); err != nil {
		t.Fatal(err)
	}

	out, err = ts.Lookup(namespace.RootContext(nil), ent.ID)
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
	ent = &logical.TokenEntry{
		Path:                     "test",
		Policies:                 []string{"dev", "ops"},
		DisplayNameDeprecated:    "test-display-name",
		CreationTimeDeprecated:   time.Now().Unix(),
		ExplicitMaxTTLDeprecated: 100,
		NumUsesDeprecated:        10,
		NamespaceID:              namespace.RootNamespaceID,
	}
	if err := ts.create(namespace.RootContext(nil), ent); err != nil {
		t.Fatalf("err: %s", err)
	}
	auth = &logical.Auth{
		DisplayName:    ent.DisplayName,
		CreationPath:   ent.Path,
		Policies:       ent.Policies,
		ExplicitMaxTTL: ent.ExplicitMaxTTL,
		NumUses:        ent.NumUses,
		LeaseOptions: logical.LeaseOptions{
			TTL: time.Hour,
		},
		ClientToken: ent.ID,
	}
	if err := ts.expiration.RegisterAuth(namespace.RootContext(nil), ent, auth, ""); err != nil {
		t.Fatal(err)
	}

	out, err = ts.Lookup(namespace.RootContext(nil), ent.ID)
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
	ent = &logical.TokenEntry{
		Path:              "test",
		NumUses:           5,
		NumUsesDeprecated: 10,
		NamespaceID:       namespace.RootNamespaceID,
	}
	if err := ts.create(namespace.RootContext(nil), ent); err != nil {
		t.Fatalf("err: %s", err)
	}
	auth = &logical.Auth{
		DisplayName:    ent.DisplayName,
		CreationPath:   ent.Path,
		Policies:       ent.Policies,
		ExplicitMaxTTL: ent.ExplicitMaxTTL,
		NumUses:        ent.NumUses,
		LeaseOptions: logical.LeaseOptions{
			TTL: time.Hour,
		},
		ClientToken: ent.ID,
	}
	if err := ts.expiration.RegisterAuth(namespace.RootContext(nil), ent, auth, ""); err != nil {
		t.Fatal(err)
	}

	out, err = ts.Lookup(namespace.RootContext(nil), ent.ID)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if out.NumUses != 5 {
		t.Fatalf("bad: num_uses: expected: 5, actual: %d", out.NumUses)
	}

	// Switch the values from deprecated and proper field and check if the
	// lower value is still getting picked up
	ent = &logical.TokenEntry{
		Path:              "test",
		NumUses:           10,
		NumUsesDeprecated: 5,
		NamespaceID:       namespace.RootNamespaceID,
	}
	if err := ts.create(namespace.RootContext(nil), ent); err != nil {
		t.Fatalf("err: %s", err)
	}
	auth = &logical.Auth{
		DisplayName:    ent.DisplayName,
		CreationPath:   ent.Path,
		Policies:       ent.Policies,
		ExplicitMaxTTL: ent.ExplicitMaxTTL,
		NumUses:        ent.NumUses,
		LeaseOptions: logical.LeaseOptions{
			TTL: time.Hour,
		},
		ClientToken: ent.ID,
	}
	if err := ts.expiration.RegisterAuth(namespace.RootContext(nil), ent, auth, ""); err != nil {
		t.Fatal(err)
	}

	out, err = ts.Lookup(namespace.RootContext(nil), ent.ID)
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

func testMakeServiceTokenViaBackend(t testing.TB, ts *TokenStore, root, client, ttl string, policy []string) {
	t.Helper()
	testMakeTokenViaBackend(t, ts, root, client, ttl, policy, false)
}

func testMakeTokenViaBackend(t testing.TB, ts *TokenStore, root, client, ttl string, policy []string, batch bool) {
	t.Helper()
	req := logical.TestRequest(t, logical.UpdateOperation, "create")
	req.ClientToken = root
	if batch {
		req.Data["type"] = "batch"
	} else {
		req.Data["id"] = client
	}
	req.Data["policies"] = policy
	req.Data["ttl"] = ttl
	resp := testMakeTokenViaRequest(t, ts, req)

	if resp.Auth.ClientToken != client {
		t.Fatalf("bad: %#v", resp)
	}
}

func testMakeTokenViaRequest(t testing.TB, ts *TokenStore, req *logical.Request) *logical.Response {
	t.Helper()
	return testMakeTokenViaRequestContext(t, namespace.RootContext(nil), ts, req)
}

func testMakeTokenViaRequestContext(t testing.TB, ctx context.Context, ts *TokenStore, req *logical.Request) *logical.Response {
	t.Helper()
	resp, err := ts.HandleRequest(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatalf("got nil token from create call")
	}
	// Let the caller handle the error
	if resp.IsError() {
		return resp
	}

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		t.Fatal(err)
	}

	te := &logical.TokenEntry{
		Path:        resp.Auth.CreationPath,
		NamespaceID: ns.ID,
	}

	if resp.Auth.TokenType != logical.TokenTypeBatch {
		if err := ts.expiration.RegisterAuth(ctx, te, resp.Auth, ""); err != nil {
			t.Fatal(err)
		}
	}

	te, err = ts.Lookup(ctx, resp.Auth.ClientToken)
	if err != nil {
		t.Fatal(err)
	}
	if te == nil {
		t.Fatal("token entry was nil")
	}

	return resp
}

func testMakeTokenDirectly(t testing.TB, ts *TokenStore, te *logical.TokenEntry) {
	if te.NamespaceID == "" {
		te.NamespaceID = namespace.RootNamespaceID
	}
	if te.CreationTime == 0 {
		te.CreationTime = time.Now().Unix()
	}
	if err := ts.create(namespace.RootContext(nil), te); err != nil {
		t.Fatal(err)
	}
	if te.Type == logical.TokenTypeDefault {
		te.Type = logical.TokenTypeService
	}
	auth := &logical.Auth{
		NumUses:     te.NumUses,
		DisplayName: te.DisplayName,
		Policies:    te.Policies,
		Metadata:    te.Meta,
		LeaseOptions: logical.LeaseOptions{
			TTL:       te.TTL,
			Renewable: te.TTL > 0,
		},
		ClientToken:    te.ID,
		Accessor:       te.Accessor,
		EntityID:       te.EntityID,
		Period:         te.Period,
		ExplicitMaxTTL: te.ExplicitMaxTTL,
		CreationPath:   te.Path,
		TokenType:      te.Type,
	}
	err := ts.expiration.RegisterAuth(namespace.RootContext(nil), te, auth, "")
	switch err {
	case nil:
		if te.Type == logical.TokenTypeBatch {
			t.Fatal("expected error from trying to register auth with batch token")
		}
	default:
		if te.Type != logical.TokenTypeBatch {
			t.Fatal(err)
		}
	}
}

func testMakeServiceTokenViaCore(t testing.TB, c *Core, root, client, ttl string, policy []string) {
	testMakeTokenViaCore(t, c, root, client, ttl, "", policy, false, nil)
}

func testMakeTokenViaCore(t testing.TB, c *Core, root, client, ttl, period string, policy []string, batch bool, outAuth *logical.Auth) {
	req := logical.TestRequest(t, logical.UpdateOperation, "auth/token/create")
	req.ClientToken = root
	if batch {
		req.Data["type"] = "batch"
	} else {
		req.Data["id"] = client
	}
	req.Data["policies"] = policy
	req.Data["ttl"] = ttl

	if len(period) != 0 {
		req.Data["period"] = period
	}

	resp, err := c.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v\nresp: %#v", err, resp)
	}
	if !batch {
		if resp.Auth.ClientToken != client {
			t.Fatalf("bad: %#v", *resp)
		}
	}
	if outAuth != nil && resp != nil && resp.Auth != nil {
		*outAuth = *resp.Auth
	}
}

func TestTokenStore_AccessorIndex(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ts := c.tokenStore

	ent := &logical.TokenEntry{
		Path:        "test",
		Policies:    []string{"dev", "ops"},
		TTL:         time.Hour,
		NamespaceID: namespace.RootNamespaceID,
	}
	testMakeTokenDirectly(t, ts, ent)

	out, err := ts.Lookup(namespace.RootContext(nil), ent.ID)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Ensure that accessor is created
	if out == nil || out.Accessor == "" {
		t.Fatalf("bad: %#v", out)
	}

	aEntry, err := ts.lookupByAccessor(namespace.RootContext(nil), out.Accessor, false, false)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Verify that the value returned from the index matches the token ID
	if aEntry.TokenID != ent.ID {
		t.Fatalf("bad: got\n%s\nexpected\n%s\n", aEntry.TokenID, ent.ID)
	}

	// Make sure a batch token doesn't get an accessor
	ent.Type = logical.TokenTypeBatch
	testMakeTokenDirectly(t, ts, ent)

	out, err = ts.Lookup(namespace.RootContext(nil), ent.ID)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Ensure that accessor is created
	if out == nil || out.Accessor != "" {
		t.Fatalf("bad: %#v", out)
	}
}

func TestTokenStore_HandleRequest_LookupAccessor(t *testing.T) {
	c, _, root := TestCoreUnsealed(t)
	ts := c.tokenStore

	testMakeServiceTokenViaBackend(t, ts, root, "tokenid", "60s", []string{"foo"})
	out, err := ts.Lookup(namespace.RootContext(nil), "tokenid")
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

	resp, err := ts.HandleRequest(namespace.RootContext(nil), req)
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
	c, _, root := TestCoreUnsealed(t)
	ts := c.tokenStore

	testKeys := []string{"token1", "token2", "token3", "token4"}
	for _, key := range testKeys {
		testMakeServiceTokenViaBackend(t, ts, root, key, "60s", []string{"foo"})
	}

	// Revoke root to make the number of accessors match
	internalRoot, _ := c.DecodeSSCToken(root)
	salted, err := ts.SaltID(namespace.RootContext(nil), internalRoot)
	if err != nil {
		t.Fatal(err)
	}
	ts.revokeInternal(namespace.RootContext(nil), salted, false)

	req := logical.TestRequest(t, logical.ListOperation, "accessors/")

	resp, err := ts.HandleRequest(namespace.RootContext(nil), req)
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
		aEntry, err := ts.lookupByAccessor(namespace.RootContext(nil), accessor, false, false)
		if err != nil {
			t.Fatal(err)
		}
		if aEntry.TokenID == "" || aEntry.AccessorID == "" {
			t.Fatalf("error, accessor entry looked up is empty, but no error thrown")
		}
		saltID, err := ts.SaltID(namespace.RootContext(nil), accessor)
		if err != nil {
			t.Fatal(err)
		}
		le := &logical.StorageEntry{Key: saltID, Value: []byte(aEntry.TokenID)}
		if err := ts.accessorView(namespace.RootNamespace).Put(namespace.RootContext(nil), le); err != nil {
			t.Fatalf("failed to persist accessor index entry: %v", err)
		}
	}

	// Do the lookup again, should get same result
	resp, err = ts.HandleRequest(namespace.RootContext(nil), req)
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
		aEntry, err := ts.lookupByAccessor(namespace.RootContext(nil), accessor, false, false)
		if err != nil {
			t.Fatal(err)
		}
		if aEntry.TokenID == "" || aEntry.AccessorID == "" {
			t.Fatalf("error, accessor entry looked up is empty, but no error thrown")
		}
	}
}

func TestTokenStore_HandleRequest_Renew_Revoke_Accessor(t *testing.T) {
	exp := mockExpiration(t)
	ts := exp.tokenStore

	rootToken, err := ts.rootToken(namespace.RootContext(nil))
	if err != nil {
		t.Fatal(err)
	}
	root := rootToken.ID

	testMakeServiceTokenViaBackend(t, ts, root, "tokenid", "", []string{"foo"})

	auth := &logical.Auth{
		ClientToken: "tokenid",
		LeaseOptions: logical.LeaseOptions{
			TTL:       time.Hour,
			Renewable: true,
		},
	}

	te, err := ts.Lookup(namespace.RootContext(nil), "tokenid")
	if err != nil {
		t.Fatal(err)
	}
	if te == nil {
		t.Fatal("token entry was nil")
	}

	err = exp.RegisterAuth(namespace.RootContext(nil), te, auth, "")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	out, err := ts.Lookup(namespace.RootContext(nil), "tokenid")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if out == nil {
		t.Fatalf("err: %s", err)
	}

	lookupAccessor := func() int64 {
		req := logical.TestRequest(t, logical.UpdateOperation, "lookup-accessor")
		req.Data = map[string]interface{}{
			"accessor": out.Accessor,
		}
		resp, err := ts.HandleRequest(namespace.RootContext(nil), req)
		if err != nil {
			t.Fatalf("err: %s", err)
		}
		if resp.Data == nil {
			t.Fatal("response should contain data")
		}
		if resp.Data["accessor"].(string) == "" {
			t.Fatal("accessor should not be empty")
		}
		ttl := resp.Data["ttl"].(int64)
		if ttl == 0 {
			t.Fatal("ttl was zero")
		}
		return ttl
	}

	firstTTL := lookupAccessor()

	// Sleep so we can verify renew behavior
	time.Sleep(10 * time.Second)

	secondTTL := lookupAccessor()
	if secondTTL >= firstTTL {
		t.Fatalf("second TTL %d was greater than first %d", secondTTL, firstTTL)
	}

	req := logical.TestRequest(t, logical.UpdateOperation, "renew-accessor")
	req.Data = map[string]interface{}{
		"accessor": out.Accessor,
	}
	resp, err := ts.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if resp.Auth == nil {
		t.Fatal("resp auth is nil")
	}
	if resp.Auth.ClientToken != "" {
		t.Fatal("client token found in response")
	}

	thirdTTL := lookupAccessor()
	if thirdTTL <= secondTTL {
		t.Fatalf("third TTL %d was not greater than second %d", thirdTTL, secondTTL)
	}

	req = logical.TestRequest(t, logical.UpdateOperation, "revoke-accessor")
	req.Data = map[string]interface{}{
		"accessor": out.Accessor,
	}

	_, err = ts.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Revoking the token using the accessor should be idempotent since
	// auth/token/revoke is
	resp, err = ts.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if len(resp.Warnings) != 1 {
		t.Fatalf("Was expecting 1 warning, got %d", len(resp.Warnings))
	}

	time.Sleep(200 * time.Millisecond)

	out, err = ts.Lookup(namespace.RootContext(nil), "tokenid")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if out != nil {
		t.Fatalf("bad:\ngot %#v\nexpected: nil\n", out)
	}

	// Now test without registering the token through the expiration manager
	testMakeServiceTokenViaBackend(t, ts, root, "tokenid", "", []string{"foo"})
	out, err = ts.Lookup(namespace.RootContext(nil), "tokenid")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if out == nil {
		t.Fatalf("err: %s", err)
	}

	req = logical.TestRequest(t, logical.UpdateOperation, "revoke-accessor")
	req.Data = map[string]interface{}{
		"accessor": out.Accessor,
	}

	_, err = ts.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	time.Sleep(200 * time.Millisecond)

	out, err = ts.Lookup(namespace.RootContext(nil), "tokenid")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if out != nil {
		t.Fatalf("bad:\ngot %#v\nexpected: nil\n", out)
	}
}

func TestTokenStore_RootToken(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ts := c.tokenStore

	te, err := ts.rootToken(namespace.RootContext(nil))
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if te.ID == "" {
		t.Fatalf("missing ID")
	}

	out, err := ts.Lookup(namespace.RootContext(nil), te.ID)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	deepEqualTokenEntries(t, out, te)
}

func TestTokenStore_NoRootBatch(t *testing.T) {
	c, _, root := TestCoreUnsealed(t)

	req := logical.TestRequest(t, logical.UpdateOperation, "auth/token/create")
	req.ClientToken = root
	req.Data["type"] = "batch"
	req.Data["policies"] = "root"
	req.Data["ttl"] = "5m"

	resp, err := c.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("expected response")
	}
	if !resp.IsError() {
		t.Fatalf("expected error, got %#v", *resp)
	}
}

func TestTokenStore_CreateLookup(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ts := c.tokenStore

	ent := &logical.TokenEntry{
		NamespaceID: namespace.RootNamespaceID,
		Path:        "test",
		Policies:    []string{"dev", "ops"},
		TTL:         time.Hour,
	}
	testMakeTokenDirectly(t, ts, ent)
	if ent.ID == "" {
		t.Fatalf("missing ID")
	}

	out, err := ts.Lookup(namespace.RootContext(nil), ent.ID)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	deepEqualTokenEntries(t, out, ent)

	// New store should share the salt
	ts2, err := NewTokenStore(namespace.RootContext(nil), hclog.New(&hclog.LoggerOptions{}), c, getBackendConfig(c))
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	ts2.SetExpirationManager(c.expiration)

	// Should still match
	out, err = ts2.Lookup(namespace.RootContext(nil), ent.ID)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	deepEqualTokenEntries(t, out, ent)
}

func TestTokenStore_CreateLookup_ProvidedID(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ts := c.tokenStore

	ent := &logical.TokenEntry{
		ID:          "foobarbaz",
		NamespaceID: namespace.RootNamespaceID,
		Path:        "test",
		Policies:    []string{"dev", "ops"},
		TTL:         time.Hour,
	}
	testMakeTokenDirectly(t, ts, ent)
	if ent.ID != "foobarbaz" {
		t.Fatalf("bad: ent.ID: expected:\"foobarbaz\"\n actual:%s", ent.ID)
	}
	if err := ts.create(namespace.RootContext(nil), ent); err == nil {
		t.Fatal("expected error creating token with the same ID")
	}

	out, err := ts.Lookup(namespace.RootContext(nil), ent.ID)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	deepEqualTokenEntries(t, out, ent)

	// New store should share the salt
	ts2, err := NewTokenStore(namespace.RootContext(nil), hclog.New(&hclog.LoggerOptions{}), c, getBackendConfig(c))
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	ts2.SetExpirationManager(c.expiration)

	// Should still match
	out, err = ts2.Lookup(namespace.RootContext(nil), ent.ID)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	deepEqualTokenEntries(t, out, ent)
}

func TestTokenStore_CreateLookup_ExpirationInRestoreMode(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ts := c.tokenStore

	ent := &logical.TokenEntry{
		NamespaceID: namespace.RootNamespaceID,
		Path:        "test",
		Policies:    []string{"dev", "ops"},
	}
	if err := ts.create(namespace.RootContext(nil), ent); err != nil {
		t.Fatalf("err: %v", err)
	}
	if ent.ID == "" {
		t.Fatalf("missing ID")
	}

	// Replace the lease with a lease with an expire time in the past
	saltedID, err := ts.SaltID(namespace.RootContext(nil), ent.ID)
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
		namespace:   namespace.RootNamespace,
	}
	if err := ts.expiration.persistEntry(namespace.RootContext(nil), le); err != nil {
		t.Fatalf("err: %v", err)
	}

	out, err := ts.Lookup(namespace.RootContext(nil), ent.ID)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	deepEqualTokenEntries(t, out, ent)

	// Set to expired lease time
	le.ExpireTime = time.Now().Add(-1 * time.Hour)
	if err := ts.expiration.persistEntry(namespace.RootContext(nil), le); err != nil {
		t.Fatalf("err: %v", err)
	}

	err = c.stopExpiration()
	if err != nil {
		t.Fatal(err)
	}

	// Reset expiration manager to restore mode
	ts.expiration.restoreModeLock.Lock()
	atomic.StoreInt32(ts.expiration.restoreMode, 1)
	ts.expiration.restoreLocks = locksutil.CreateLocks()
	ts.expiration.restoreModeLock.Unlock()

	// Test that the token lookup does not return the token entry due to the
	// expired lease
	out, err = ts.Lookup(namespace.RootContext(nil), ent.ID)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out != nil {
		t.Fatalf("lease expired, no token expected: %#v", out)
	}
}

func TestTokenStore_UseToken(t *testing.T) {
	c, _, root := TestCoreUnsealed(t)
	ts := c.tokenStore

	// Lookup the root token
	ent, err := ts.Lookup(namespace.RootContext(nil), root)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Root is an unlimited use token, should be a no-op
	te, err := ts.UseToken(namespace.RootContext(nil), ent)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if te == nil {
		t.Fatalf("token entry after use was nil")
	}

	// Lookup the root token again
	ent2, err := ts.Lookup(namespace.RootContext(nil), root)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if !reflect.DeepEqual(ent, ent2) {
		t.Fatalf("bad: ent:%#v ent2:%#v", ent, ent2)
	}

	// Create a restricted token
	ent = &logical.TokenEntry{
		Path:        "test",
		Policies:    []string{"dev", "ops"},
		NumUses:     2,
		TTL:         time.Hour,
		NamespaceID: namespace.RootNamespaceID,
	}
	testMakeTokenDirectly(t, ts, ent)

	// Use the token
	te, err = ts.UseToken(namespace.RootContext(nil), ent)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if te == nil {
		t.Fatalf("token entry for use #1 was nil")
	}

	// Lookup the token
	ent2, err = ts.Lookup(namespace.RootContext(nil), ent.ID)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Should be reduced
	if ent2.NumUses != 1 {
		t.Fatalf("bad: %#v", ent2)
	}

	// Use the token
	te, err = ts.UseToken(namespace.RootContext(nil), ent)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if te == nil {
		t.Fatalf("token entry for use #2 was nil")
	}
	if te.NumUses != tokenRevocationPending {
		t.Fatalf("token entry after use #2 did not have revoke flag")
	}
	ts.revokeOrphan(namespace.RootContext(nil), te.ID)

	// Lookup the token
	ent2, err = ts.Lookup(namespace.RootContext(nil), ent.ID)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Should be revoked
	if ent2 != nil {
		t.Fatalf("bad: %#v", ent2)
	}
}

func TestTokenStore_Revoke(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ts := c.tokenStore

	ent := &logical.TokenEntry{Path: "test", Policies: []string{"dev", "ops"}, NamespaceID: namespace.RootNamespaceID}
	if err := ts.create(namespace.RootContext(nil), ent); err != nil {
		t.Fatalf("err: %v", err)
	}

	err := ts.revokeOrphan(namespace.RootContext(nil), "")
	if err.Error() != "cannot revoke blank token" {
		t.Fatalf("err: %v", err)
	}
	err = ts.revokeOrphan(namespace.RootContext(nil), ent.ID)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	out, err := ts.Lookup(namespace.RootContext(nil), ent.ID)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out != nil {
		t.Fatalf("bad: %#v", out)
	}
}

func TestTokenStore_Revoke_Leases(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ts := c.tokenStore

	view := NewBarrierView(c.barrier, "noop/")

	// Mount a noop backend
	noop := &NoopBackend{}
	err := ts.expiration.router.Mount(noop, "noop/", &MountEntry{UUID: "noopuuid", Accessor: "noopaccessor", namespace: namespace.RootNamespace}, view)
	if err != nil {
		t.Fatal(err)
	}

	ent := &logical.TokenEntry{Path: "test", Policies: []string{"dev", "ops"}, NamespaceID: namespace.RootNamespaceID}
	if err := ts.create(namespace.RootContext(nil), ent); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Register a lease
	req := &logical.Request{
		Operation:   logical.ReadOperation,
		Path:        "noop/foo",
		ClientToken: ent.ID,
	}
	req.SetTokenEntry(ent)
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
	leaseID, err := ts.expiration.Register(namespace.RootContext(nil), req, resp, "")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Revoke the token
	err = ts.revokeOrphan(namespace.RootContext(nil), ent.ID)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	time.Sleep(200 * time.Millisecond)

	// Verify the lease is gone
	out, err := ts.expiration.loadEntry(namespace.RootContext(nil), leaseID)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out != nil {
		t.Fatalf("bad: %#v", out)
	}
}

func TestTokenStore_Revoke_Orphan(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ts := c.tokenStore

	ent := &logical.TokenEntry{
		NamespaceID: namespace.RootNamespaceID,
		Path:        "test",
		Policies:    []string{"dev", "ops"},
		TTL:         time.Hour,
	}
	testMakeTokenDirectly(t, ts, ent)

	ent2 := &logical.TokenEntry{
		NamespaceID: namespace.RootNamespaceID,
		Parent:      ent.ID,
		TTL:         time.Hour,
	}
	testMakeTokenDirectly(t, ts, ent2)

	err := ts.revokeOrphan(namespace.RootContext(nil), ent.ID)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	out, err := ts.Lookup(namespace.RootContext(nil), ent2.ID)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Unset the expected token parent's ID
	ent2.Parent = ""

	deepEqualTokenEntries(t, out, ent2)
}

// This was the original function name, and now it just calls
// the non recursive version for a variety of depths.
func TestTokenStore_RevokeTree(t *testing.T) {
	testTokenStore_RevokeTree_NonRecursive(t, 1, false)
	testTokenStore_RevokeTree_NonRecursive(t, 2, false)
	testTokenStore_RevokeTree_NonRecursive(t, 10, false)

	// corrupted trees with cycles
	testTokenStore_RevokeTree_NonRecursive(t, 1, true)
	testTokenStore_RevokeTree_NonRecursive(t, 10, true)
}

// Revokes a given Token Store tree non recursively.
// The second parameter refers to the depth of the tree.
func testTokenStore_RevokeTree_NonRecursive(t testing.TB, depth uint64, injectCycles bool) {
	c, _, _ := TestCoreUnsealed(t)
	ts := c.tokenStore
	root, children := buildTokenTree(t, ts, depth)

	var cyclePaths []string
	if injectCycles {
		// Make the root the parent of itself
		saltedRoot, _ := ts.SaltID(namespace.RootContext(nil), root.ID)
		key := fmt.Sprintf("%s/%s", saltedRoot, saltedRoot)
		cyclePaths = append(cyclePaths, key)
		le := &logical.StorageEntry{Key: key}

		if err := ts.parentView(namespace.RootNamespace).Put(namespace.RootContext(nil), le); err != nil {
			t.Fatalf("err: %v", err)
		}

		// Make a deep child the parent of a shallow child
		shallow, _ := ts.SaltID(namespace.RootContext(nil), children[0].ID)
		deep, _ := ts.SaltID(namespace.RootContext(nil), children[len(children)-1].ID)
		key = fmt.Sprintf("%s/%s", deep, shallow)
		cyclePaths = append(cyclePaths, key)
		le = &logical.StorageEntry{Key: key}

		if err := ts.parentView(namespace.RootNamespace).Put(namespace.RootContext(nil), le); err != nil {
			t.Fatalf("err: %v", err)
		}
	}

	err := ts.revokeTree(namespace.RootContext(nil), &leaseEntry{})
	if err.Error() != "cannot tree-revoke blank token" {
		t.Fatal(err)
	}

	saltCtx := namespace.RootContext(nil)
	saltedID, err := c.tokenStore.SaltID(saltCtx, root.ID)
	if err != nil {
		t.Fatal(err)
	}
	tokenLeaseID := path.Join(root.Path, saltedID)

	tokenLease, err := ts.expiration.loadEntry(namespace.RootContext(nil), tokenLeaseID)
	if err != nil || tokenLease == nil {
		t.Fatalf("err: %v, tokenLease: %#v", err, tokenLease)
	}

	// Nuke tree non recursively.
	err = ts.revokeTree(namespace.RootContext(nil), tokenLease)
	if err != nil {
		t.Fatal(err)
	}
	// Append the root to ensure it was successfully
	// deleted.
	children = append(children, root)
	for _, entry := range children {
		out, err := ts.Lookup(namespace.RootContext(nil), entry.ID)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if out != nil {
			t.Fatalf("bad: %#v", out)
		}
	}

	for _, path := range cyclePaths {
		entry, err := ts.parentView(namespace.RootNamespace).Get(namespace.RootContext(nil), path)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if entry != nil {
			t.Fatalf("expected reference to be deleted: %v", entry)
		}
	}
}

// A benchmark function that tests testTokenStore_RevokeTree_NonRecursive
// for a variety of different depths.
func BenchmarkTokenStore_RevokeTree(b *testing.B) {
	benchmarks := []uint64{0, 1, 2, 4, 8, 16, 20}
	for _, depth := range benchmarks {
		b.Run(fmt.Sprintf("Tree of Depth %d", depth), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				testTokenStore_RevokeTree_NonRecursive(b, depth, false)
			}
		})
	}
}

// Builds a TokenTree of a specified depth, so that
// we may run revoke tests on it.
func buildTokenTree(t testing.TB, ts *TokenStore, depth uint64) (root *logical.TokenEntry, children []*logical.TokenEntry) {
	root = &logical.TokenEntry{
		NamespaceID: namespace.RootNamespaceID,
		TTL:         time.Hour,
	}
	testMakeTokenDirectly(t, ts, root)

	frontier := []*logical.TokenEntry{root}
	current := uint64(0)
	for current < depth {
		next := make([]*logical.TokenEntry, 0, 2*len(frontier))
		for _, node := range frontier {
			left := &logical.TokenEntry{
				Parent:      node.ID,
				TTL:         time.Hour,
				NamespaceID: namespace.RootNamespaceID,
			}
			testMakeTokenDirectly(t, ts, left)

			right := &logical.TokenEntry{
				Parent:      node.ID,
				TTL:         time.Hour,
				NamespaceID: namespace.RootNamespaceID,
			}
			testMakeTokenDirectly(t, ts, right)

			children = append(children, left, right)
			next = append(next, left, right)
		}
		frontier = next
		current++
	}

	return root, children
}

func TestTokenStore_RevokeSelf(t *testing.T) {
	exp := mockExpiration(t)
	ts := exp.tokenStore

	ent1 := &logical.TokenEntry{
		TTL:         time.Hour,
		NamespaceID: namespace.RootNamespaceID,
	}
	testMakeTokenDirectly(t, ts, ent1)

	ent2 := &logical.TokenEntry{
		Parent:      ent1.ID,
		TTL:         time.Hour,
		NamespaceID: namespace.RootNamespaceID,
	}
	testMakeTokenDirectly(t, ts, ent2)

	ent3 := &logical.TokenEntry{
		Parent:      ent2.ID,
		TTL:         time.Hour,
		NamespaceID: namespace.RootNamespaceID,
	}
	testMakeTokenDirectly(t, ts, ent3)

	ent4 := &logical.TokenEntry{
		Parent:      ent2.ID,
		TTL:         time.Hour,
		NamespaceID: namespace.RootNamespaceID,
	}
	testMakeTokenDirectly(t, ts, ent4)

	req := logical.TestRequest(t, logical.UpdateOperation, "revoke-self")
	req.ClientToken = ent1.ID

	resp, err := ts.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v\nresp: %#v", err, resp)
	}

	lookup := []string{ent1.ID, ent2.ID, ent3.ID, ent4.ID}
	var out *logical.TokenEntry
	for _, id := range lookup {
		var found bool
		for i := 0; i < 10; i++ {
			out, err = ts.Lookup(namespace.RootContext(nil), id)
			if err != nil {
				t.Fatalf("err: %v", err)
			}
			if out == nil {
				found = true
				break
			}
			time.Sleep(1000 * time.Millisecond)
		}
		if !found {
			t.Fatalf("bad: %#v", out)
		}
	}
}

func TestTokenStore_HandleRequest_NonAssignable(t *testing.T) {
	c, _, root := TestCoreUnsealed(t)
	ts := c.tokenStore

	req := logical.TestRequest(t, logical.UpdateOperation, "create")
	req.ClientToken = root
	req.Data["policies"] = []string{"default", "foo"}

	resp, err := ts.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v\nresp: %#v", err, resp)
	}

	req.Data["policies"] = []string{"default", "foo", responseWrappingPolicyName}

	resp, err = ts.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("got a nil response")
	}
	if !resp.IsError() {
		t.Fatalf("expected error; response is %#v", *resp)
	}

	// Batch tokens too
	req.Data["type"] = "batch"
	resp, err = ts.HandleRequest(namespace.RootContext(nil), req)
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
	c, _, root := TestCoreUnsealed(t)
	ts := c.tokenStore

	req := logical.TestRequest(t, logical.UpdateOperation, "create")
	req.ClientToken = root
	req.Data["display_name"] = "foo_bar.baz!"

	resp, err := ts.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v\nresp: %#v", err, resp)
	}

	expected := &logical.TokenEntry{
		ID:          resp.Auth.ClientToken,
		NamespaceID: namespace.RootNamespaceID,
		Accessor:    resp.Auth.Accessor,
		Parent:      root,
		Policies:    []string{"root"},
		Path:        "auth/token/create",
		DisplayName: "token-foo-bar-baz",
		TTL:         0,
		Type:        logical.TokenTypeService,
	}
	out, err := ts.Lookup(namespace.RootContext(nil), resp.Auth.ClientToken)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	expected.CreationTime = out.CreationTime
	expected.CubbyholeID = out.CubbyholeID
	if !reflect.DeepEqual(out, expected) {
		t.Fatalf("bad: expected:%#v\nactual:%#v", expected, out)
	}
}

func deepEqualTokenEntries(t *testing.T, a *logical.TokenEntry, b *logical.TokenEntry) {
	if diff := cmp.Diff(a, b, cmpopts.IgnoreFields(logical.TokenEntry{}, "ExternalID")); diff != "" {
		t.Fatalf("bad diff in token entries: %s", diff)
	}
}

func TestTokenStore_HandleRequest_CreateToken_NumUses(t *testing.T) {
	c, _, root := TestCoreUnsealed(t)
	ts := c.tokenStore

	req := logical.TestRequest(t, logical.UpdateOperation, "create")
	req.ClientToken = root
	req.Data["num_uses"] = "1"

	// Make sure batch tokens can't do limited use counts
	req.Data["type"] = "batch"
	resp, err := ts.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("error handling request: %v", err)
	}
	if resp == nil || !resp.IsError() {
		t.Fatalf("expected response error: resp: %#v ", resp)
	}

	delete(req.Data, "type")
	resp, err = ts.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v\nresp: %#v", err, resp)
	}

	expected := &logical.TokenEntry{
		ID:          resp.Auth.ClientToken,
		NamespaceID: namespace.RootNamespaceID,
		Accessor:    resp.Auth.Accessor,
		Parent:      root,
		Policies:    []string{"root"},
		Path:        "auth/token/create",
		DisplayName: "token",
		NumUses:     1,
		TTL:         0,
		Type:        logical.TokenTypeService,
	}
	out, err := ts.Lookup(namespace.RootContext(nil), resp.Auth.ClientToken)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	expected.CreationTime = out.CreationTime
	expected.CubbyholeID = out.CubbyholeID
	if !reflect.DeepEqual(out, expected) {
		t.Fatalf("bad: expected:%#v\nactual:%#v", expected, out)
	}
}

func TestTokenStore_HandleRequest_CreateToken_NumUses_Invalid(t *testing.T) {
	c, _, root := TestCoreUnsealed(t)
	ts := c.tokenStore

	req := logical.TestRequest(t, logical.UpdateOperation, "create")
	req.ClientToken = root
	req.Data["num_uses"] = "-1"

	resp, err := ts.HandleRequest(namespace.RootContext(nil), req)
	if err != logical.ErrInvalidRequest {
		t.Fatalf("err: %v resp: %#v", err, resp)
	}
}

func TestTokenStore_HandleRequest_CreateToken_NumUses_Restricted(t *testing.T) {
	c, _, root := TestCoreUnsealed(t)
	ts := c.tokenStore

	req := logical.TestRequest(t, logical.UpdateOperation, "create")
	req.ClientToken = root
	req.Data["num_uses"] = "1"

	resp, err := ts.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v\nresp: %#v", err, resp)
	}

	// We should NOT be able to use the restricted token to create a new token
	req.ClientToken = resp.Auth.ClientToken
	_, err = ts.HandleRequest(namespace.RootContext(nil), req)
	if err != logical.ErrInvalidRequest {
		t.Fatalf("err: %v resp: %#v", err, resp)
	}
}

func TestTokenStore_HandleRequest_CreateToken_NoPolicy(t *testing.T) {
	c, _, root := TestCoreUnsealed(t)
	ts := c.tokenStore

	req := logical.TestRequest(t, logical.UpdateOperation, "create")
	req.ClientToken = root

	// Make sure batch tokens won't automatically assign root
	req.Data["type"] = "batch"
	resp, err := ts.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("error handling request: %v", err)
	}
	if resp == nil || !resp.IsError() {
		t.Fatalf("expected response error: resp: %#v", resp)
	}

	delete(req.Data, "type")
	resp, err = ts.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v\nresp: %#v", err, resp)
	}

	expected := &logical.TokenEntry{
		ID:          resp.Auth.ClientToken,
		NamespaceID: namespace.RootNamespaceID,
		Accessor:    resp.Auth.Accessor,
		Parent:      root,
		Policies:    []string{"root"},
		Path:        "auth/token/create",
		DisplayName: "token",
		TTL:         0,
		Type:        logical.TokenTypeService,
	}
	out, err := ts.Lookup(namespace.RootContext(nil), resp.Auth.ClientToken)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	expected.CreationTime = out.CreationTime
	expected.CubbyholeID = out.CubbyholeID
	if !reflect.DeepEqual(out, expected) {
		t.Fatalf("bad: expected:%#v\nactual:%#v", expected, out)
	}
}

func TestTokenStore_HandleRequest_CreateToken_BadParent(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ts := c.tokenStore

	req := logical.TestRequest(t, logical.UpdateOperation, "create")
	req.ClientToken = "random"

	resp, err := ts.HandleRequest(namespace.RootContext(nil), req)
	if err != logical.ErrInvalidRequest {
		t.Fatalf("err: %v resp: %#v", err, resp)
	}
	if resp.Data["error"] != "parent token lookup failed: no parent found" {
		t.Fatalf("bad: %#v", resp)
	}
}

func TestTokenStore_HandleRequest_CreateToken(t *testing.T) {
	c, _, root := TestCoreUnsealed(t)
	ts := c.tokenStore

	req := logical.TestRequest(t, logical.UpdateOperation, "create")
	req.ClientToken = root
	req.Data["policies"] = []string{"foo"}

	resp, err := ts.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v\nresp: %#v", err, resp)
	}
	if resp.Auth.ClientToken == "" {
		t.Fatalf("bad: %#v", resp)
	}
}

func TestTokenStore_HandleRequest_CreateToken_RootID(t *testing.T) {
	c, _, root := TestCoreUnsealed(t)
	ts := c.tokenStore

	req := logical.TestRequest(t, logical.UpdateOperation, "create")
	req.ClientToken = root
	req.Data["id"] = "foobar"
	req.Data["policies"] = []string{"foo"}

	resp, err := ts.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v\nresp: %#v", err, resp)
	}
	if resp.Auth.ClientToken != "foobar" {
		t.Fatalf("bad: %#v", resp)
	}

	// Retry with batch; batch should not actually accept a custom ID
	req.Data["type"] = "batch"
	resp, err = ts.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v\nresp: %#v", err, resp)
	}
	out, _ := ts.Lookup(namespace.RootContext(nil), resp.Auth.ClientToken)
	if out.ID == "foobar" {
		t.Fatalf("bad: %#v", out)
	}
}

func TestTokenStore_HandleRequest_CreateToken_NonRootID(t *testing.T) {
	c, _, root := TestCoreUnsealed(t)
	ts := c.tokenStore
	testMakeServiceTokenViaBackend(t, ts, root, "client", "", []string{"foo"})

	req := logical.TestRequest(t, logical.UpdateOperation, "create")
	req.ClientToken = "client"
	req.Data["id"] = "foobar"
	req.Data["policies"] = []string{"foo"}

	resp, err := ts.HandleRequest(namespace.RootContext(nil), req)
	if err != logical.ErrInvalidRequest {
		t.Fatalf("err: %v resp: %#v", err, resp)
	}
	if resp.Data["error"] != "root or sudo privileges required to specify token id" {
		t.Fatalf("bad: %#v", resp)
	}

	// Retry with batch
	req.Data["type"] = "batch"
	resp, err = ts.HandleRequest(namespace.RootContext(nil), req)
	if err != logical.ErrInvalidRequest {
		t.Fatalf("err: %v resp: %#v", err, resp)
	}
	if resp.Data["error"] != "root or sudo privileges required to specify token id" {
		t.Fatalf("bad: %#v", resp)
	}
}

func TestTokenStore_HandleRequest_CreateToken_NonRoot_Subset(t *testing.T) {
	c, _, root := TestCoreUnsealed(t)
	ts := c.tokenStore
	testMakeServiceTokenViaBackend(t, ts, root, "client", "", []string{"foo", "bar"})

	req := logical.TestRequest(t, logical.UpdateOperation, "create")
	req.ClientToken = "client"
	req.Data["policies"] = []string{"foo"}

	resp, err := ts.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v\nresp: %#v", err, resp)
	}
	if resp.Auth.ClientToken == "" {
		t.Fatalf("bad: %#v", resp)
	}

	ent := &logical.TokenEntry{
		NamespaceID: namespace.RootNamespaceID,
		Path:        "test",
		Policies:    []string{"foo", "bar"},
		TTL:         time.Hour,
	}
	testMakeTokenDirectly(t, ts, ent)
	req.ClientToken = ent.ID
	resp, err = ts.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v\nresp: %#v", err, resp)
	}
	if resp.Auth.ClientToken == "" {
		t.Fatalf("bad: %#v", resp)
	}
}

func TestTokenStore_HandleRequest_CreateToken_NonRoot_InvalidSubset(t *testing.T) {
	c, _, root := TestCoreUnsealed(t)
	ts := c.tokenStore
	testMakeServiceTokenViaBackend(t, ts, root, "client", "", []string{"foo", "bar"})

	req := logical.TestRequest(t, logical.UpdateOperation, "create")
	req.ClientToken = "client"
	req.Data["policies"] = []string{"foo", "bar", "baz"}

	resp, err := ts.HandleRequest(namespace.RootContext(nil), req)
	if err != logical.ErrInvalidRequest {
		t.Fatalf("err: %v resp: %#v", err, resp)
	}
	if resp.Data["error"] != "child policies must be subset of parent" {
		t.Fatalf("bad: %#v", resp)
	}
}

func TestTokenStore_HandleRequest_CreateToken_NonRoot_RootChild(t *testing.T) {
	core, _, root := TestCoreUnsealed(t)
	ts := core.tokenStore
	ps := core.policyStore

	policy, _ := ParseACLPolicy(namespace.RootNamespace, tokenCreationPolicy)
	policy.Name = "test1"
	if err := ps.SetPolicy(namespace.RootContext(nil), policy); err != nil {
		t.Fatal(err)
	}

	testMakeServiceTokenViaBackend(t, ts, root, "sudoClient", "", []string{"test1"})

	req := logical.TestRequest(t, logical.UpdateOperation, "create")
	req.ClientToken = "sudoClient"
	req.MountPoint = "auth/token/"
	req.Data["policies"] = []string{"root"}

	resp, err := ts.HandleRequest(namespace.RootContext(nil), req)
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
	c, _, root := TestCoreUnsealed(t)
	ts := c.tokenStore

	req := logical.TestRequest(t, logical.UpdateOperation, "create")
	req.ClientToken = root
	req.Data = map[string]interface{}{
		"ttl":      "5m",
		"policies": "root",
	}
	resp := testMakeTokenViaRequest(t, ts, req)

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
	resp, err := ts.HandleRequest(namespace.RootContext(nil), req)
	if err == nil {
		t.Fatalf("expected error")
	}
	expectedErr := "expiring root tokens cannot create non-expiring root tokens"
	if resp.Data["error"].(string) != expectedErr {
		t.Fatalf("got err=%v, expected=%v", resp.Data["error"], expectedErr)
	}
}

func TestTokenStore_HandleRequest_CreateToken_Root_RootChild(t *testing.T) {
	c, _, root := TestCoreUnsealed(t)
	ts := c.tokenStore

	req := logical.TestRequest(t, logical.UpdateOperation, "create")
	req.ClientToken = root

	resp, err := ts.HandleRequest(namespace.RootContext(nil), req)
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
	c, _, root := TestCoreUnsealed(t)
	ts := c.tokenStore
	testMakeServiceTokenViaBackend(t, ts, root, "client", "", []string{"foo"})

	req := logical.TestRequest(t, logical.UpdateOperation, "create")
	req.ClientToken = "client"
	req.Data["no_parent"] = true
	req.Data["policies"] = []string{"foo"}

	resp, err := ts.HandleRequest(namespace.RootContext(nil), req)
	if err != logical.ErrInvalidRequest {
		t.Fatalf("err: %v resp: %#v", err, resp)
	}
	if resp.Data["error"] != "root or sudo privileges required to create orphan token" {
		t.Fatalf("bad: %#v", resp)
	}
}

func TestTokenStore_HandleRequest_CreateToken_Root_NoParent(t *testing.T) {
	c, _, root := TestCoreUnsealed(t)
	ts := c.tokenStore

	req := logical.TestRequest(t, logical.UpdateOperation, "create")
	req.ClientToken = root
	req.Data["no_parent"] = true
	req.Data["policies"] = []string{"foo"}

	resp := testMakeTokenViaRequest(t, ts, req)
	if resp.Auth.ClientToken == "" {
		t.Fatalf("bad: %#v", resp)
	}

	out, _ := ts.Lookup(namespace.RootContext(nil), resp.Auth.ClientToken)
	if out.Parent != "" {
		t.Fatalf("bad: %#v", out)
	}
}

func TestTokenStore_HandleRequest_CreateToken_PathBased_NoParent(t *testing.T) {
	c, _, root := TestCoreUnsealed(t)
	ts := c.tokenStore

	req := logical.TestRequest(t, logical.UpdateOperation, "create-orphan")
	req.ClientToken = root
	req.Data["policies"] = []string{"foo"}

	resp := testMakeTokenViaRequest(t, ts, req)

	if resp.Auth.ClientToken == "" {
		t.Fatalf("bad: %#v", resp)
	}

	out, _ := ts.Lookup(namespace.RootContext(nil), resp.Auth.ClientToken)
	if out.Parent != "" {
		t.Fatalf("bad: %#v", out)
	}
}

func TestTokenStore_HandleRequest_CreateToken_Metadata(t *testing.T) {
	c, _, root := TestCoreUnsealed(t)
	ts := c.tokenStore

	req := logical.TestRequest(t, logical.UpdateOperation, "create")
	req.ClientToken = root
	req.Data["policies"] = []string{"foo"}
	meta := map[string]string{
		"user":   "armon",
		"source": "github",
	}
	req.Data["meta"] = meta

	resp := testMakeTokenViaRequest(t, ts, req)
	if resp.Auth.ClientToken == "" {
		t.Fatalf("bad: %#v", resp)
	}

	out, _ := ts.Lookup(namespace.RootContext(nil), resp.Auth.ClientToken)
	if !reflect.DeepEqual(out.Meta, meta) {
		t.Fatalf("bad: expected:%#v\nactual:%#v", meta, out.Meta)
	}

	// Test with batch tokens
	req.Data["type"] = "batch"
	resp = testMakeTokenViaRequest(t, ts, req)
	if resp.Auth.ClientToken == "" {
		t.Fatalf("bad: %#v", resp)
	}

	out, _ = ts.Lookup(namespace.RootContext(nil), resp.Auth.ClientToken)
	if !reflect.DeepEqual(out.Meta, meta) {
		t.Fatalf("bad: expected:%#v\nactual:%#v", meta, out.Meta)
	}
}

func TestTokenStore_HandleRequest_CreateToken_Lease(t *testing.T) {
	c, _, root := TestCoreUnsealed(t)
	ts := c.tokenStore

	req := logical.TestRequest(t, logical.UpdateOperation, "create")
	req.ClientToken = root
	req.Data["policies"] = []string{"foo"}
	req.Data["lease"] = "1h"

	resp := testMakeTokenViaRequest(t, ts, req)
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
	c, _, root := TestCoreUnsealed(t)
	ts := c.tokenStore

	req := logical.TestRequest(t, logical.UpdateOperation, "create")
	req.ClientToken = root
	req.Data["policies"] = []string{"foo"}
	req.Data["ttl"] = "1h"

	resp := testMakeTokenViaRequest(t, ts, req)
	if resp.Auth.ClientToken == "" {
		t.Fatalf("bad: %#v", resp)
	}
	if resp.Auth.TTL != time.Hour {
		t.Fatalf("bad: %#v", resp)
	}
	if !resp.Auth.Renewable {
		t.Fatalf("bad: %#v", resp)
	}

	// Test batch tokens
	req.Data["type"] = "batch"
	resp = testMakeTokenViaRequest(t, ts, req)
	if resp.Auth.ClientToken == "" {
		t.Fatalf("bad: %#v", resp)
	}
	if resp.Auth.TTL != time.Hour {
		t.Fatalf("bad: %#v", resp)
	}
	if resp.Auth.Renewable {
		t.Fatalf("bad: %#v", resp)
	}
}

func TestTokenStore_HandleRequest_CreateToken_Metric(t *testing.T) {
	c, _, root, sink := TestCoreUnsealedWithMetrics(t)
	ts := c.tokenStore

	req := logical.TestRequest(t, logical.UpdateOperation, "create")
	req.ClientToken = root
	req.Data["ttl"] = "3h"
	req.Data["policies"] = []string{"foo"}
	req.MountPoint = "test/mount"

	resp := testMakeTokenViaRequest(t, ts, req)
	if resp.Auth.ClientToken == "" {
		t.Fatalf("bad: %#v", resp)
	}

	expectSingleCount(t, sink, "token.creation")
}

func TestTokenStore_HandleRequest_Revoke(t *testing.T) {
	exp := mockExpiration(t)
	ts := exp.tokenStore

	rootToken, err := ts.rootToken(namespace.RootContext(nil))
	if err != nil {
		t.Fatal(err)
	}
	root := rootToken.ID

	testMakeServiceTokenViaBackend(t, ts, root, "child", "60s", []string{"root", "foo"})

	te, err := ts.Lookup(namespace.RootContext(nil), "child")
	if err != nil {
		t.Fatal(err)
	}
	if te == nil {
		t.Fatal("token entry was nil")
	}

	auth := &logical.Auth{
		ClientToken: "child",
		LeaseOptions: logical.LeaseOptions{
			TTL:       time.Hour,
			Renewable: true,
		},
	}
	err = exp.RegisterAuth(namespace.RootContext(nil), te, auth, "")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	testMakeServiceTokenViaBackend(t, ts, "child", "sub-child", "50s", []string{"foo"})

	te, err = ts.Lookup(namespace.RootContext(nil), "sub-child")
	if err != nil {
		t.Fatal(err)
	}
	if te == nil {
		t.Fatal("token entry was nil")
	}

	auth = &logical.Auth{
		ClientToken: "sub-child",
		LeaseOptions: logical.LeaseOptions{
			TTL:       time.Hour,
			Renewable: true,
		},
	}
	err = exp.RegisterAuth(namespace.RootContext(nil), te, auth, "")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	req := logical.TestRequest(t, logical.UpdateOperation, "revoke")
	req.Data = map[string]interface{}{
		"token": "child",
	}
	resp, err := ts.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v\nresp: %#v", err, resp)
	}
	if resp != nil {
		t.Fatalf("bad: %#v", resp)
	}

	time.Sleep(200 * time.Millisecond)

	out, err := ts.Lookup(namespace.RootContext(nil), "child")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out != nil {
		t.Fatalf("bad: %#v", out)
	}

	// Sub-child should not exist
	out, err = ts.Lookup(namespace.RootContext(nil), "sub-child")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out != nil {
		t.Fatalf("bad: %v", out)
	}

	// Now test without registering the tokens through the expiration manager
	testMakeServiceTokenViaBackend(t, ts, root, "child", "60s", []string{"root", "foo"})
	testMakeServiceTokenViaBackend(t, ts, "child", "sub-child", "50s", []string{"foo"})

	req = logical.TestRequest(t, logical.UpdateOperation, "revoke")
	req.Data = map[string]interface{}{
		"token": "child",
	}
	resp, err = ts.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v\nresp: %#v", err, resp)
	}
	if resp != nil {
		t.Fatalf("bad: %#v", resp)
	}

	time.Sleep(200 * time.Millisecond)

	out, err = ts.Lookup(namespace.RootContext(nil), "child")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out != nil {
		t.Fatalf("bad: %#v", out)
	}

	// Sub-child should not exist
	out, err = ts.Lookup(namespace.RootContext(nil), "sub-child")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out != nil {
		t.Fatalf("bad: %v", out)
	}
}

func TestTokenStore_HandleRequest_RevokeOrphan(t *testing.T) {
	c, _, root := TestCoreUnsealed(t)
	ts := c.tokenStore
	testMakeServiceTokenViaBackend(t, ts, root, "child", "60s", []string{"root", "foo"})
	testMakeServiceTokenViaBackend(t, ts, "child", "sub-child", "50s", []string{"foo"})

	req := logical.TestRequest(t, logical.UpdateOperation, "auth/token/revoke-orphan")
	req.Data = map[string]interface{}{
		"token": "child",
	}
	req.ClientToken = root
	resp, err := c.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v\nresp: %#v", err, resp)
	}
	if resp != nil {
		t.Fatalf("bad: %#v", resp)
	}

	time.Sleep(200 * time.Millisecond)

	out, err := ts.Lookup(namespace.RootContext(nil), "child")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out != nil {
		t.Fatalf("bad: %v", out)
	}

	// Check that the parent entry is properly cleaned up
	saltedID, err := ts.SaltID(namespace.RootContext(nil), "child")
	if err != nil {
		t.Fatal(err)
	}
	children, err := ts.idView(namespace.RootNamespace).List(namespace.RootContext(nil), parentPrefix+saltedID+"/")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(children) != 0 {
		t.Fatalf("bad: %v", children)
	}

	// Sub-child should exist!
	out, err = ts.Lookup(namespace.RootContext(nil), "sub-child")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out == nil {
		t.Fatalf("bad: %v", out)
	}
}

func TestTokenStore_HandleRequest_RevokeOrphan_NonRoot(t *testing.T) {
	c, _, root := TestCoreUnsealed(t)
	ts := c.tokenStore
	testMakeServiceTokenViaBackend(t, ts, root, "child", "60s", []string{"foo"})

	out, err := ts.Lookup(namespace.RootContext(nil), "child")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out == nil {
		t.Fatalf("bad: %v", out)
	}

	req := logical.TestRequest(t, logical.UpdateOperation, "auth/token/revoke-orphan")
	req.Data = map[string]interface{}{
		"token": "child",
	}
	req.ClientToken = "child"
	resp, err := c.HandleRequest(namespace.RootContext(nil), req)
	if !errors.Is(err, logical.ErrPermissionDenied) {
		t.Fatalf("did not get expected error when non-root revoking itself with orphan flag; resp is %#v; err is %#v", resp, err)
	}

	time.Sleep(200 * time.Millisecond)

	// Should still exist
	out, err = ts.Lookup(namespace.RootContext(nil), "child")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out == nil {
		t.Fatalf("bad: %v", out)
	}
}

func TestTokenStore_HandleRequest_Lookup(t *testing.T) {
	t.Run("service_token", func(t *testing.T) {
		testTokenStoreHandleRequestLookup(t, false, false)
	})

	t.Run("service_token_periodic", func(t *testing.T) {
		testTokenStoreHandleRequestLookup(t, false, true)
	})

	t.Run("batch_token", func(t *testing.T) {
		testTokenStoreHandleRequestLookup(t, true, false)
	})
}

func testTokenStoreHandleRequestLookup(t *testing.T, batch, periodic bool) {
	c, _, root := TestCoreUnsealed(t)
	ts := c.tokenStore
	req := logical.TestRequest(t, logical.UpdateOperation, "lookup")
	req.Data = map[string]interface{}{
		"token": root,
	}
	resp, err := ts.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v\nresp: %#v", err, resp)
	}
	if resp == nil {
		t.Fatalf("bad: %#v", resp)
	}

	internalRoot, _ := c.DecodeSSCToken(root)

	exp := map[string]interface{}{
		"id":               internalRoot,
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
		"type":             "service",
	}

	if resp.Data["creation_time"].(int64) == 0 {
		t.Fatalf("creation time was zero")
	}
	delete(resp.Data, "creation_time")

	if !reflect.DeepEqual(resp.Data, exp) {
		t.Fatalf("bad: expected:%#v\nactual:%#v", exp, resp.Data)
	}

	period := ""
	if periodic {
		period = "3600s"
	}

	outAuth := new(logical.Auth)
	testMakeTokenViaCore(t, c, root, "client", "3600s", period, []string{"foo"}, batch, outAuth)

	tokenType := "service"
	expID := "client"
	if batch {
		tokenType = "batch"
		expID = outAuth.ClientToken
	}

	// Test via POST
	req = logical.TestRequest(t, logical.UpdateOperation, "lookup")
	req.Data = map[string]interface{}{
		"token": expID,
	}
	resp, err = ts.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v\nresp: %#v", err, resp)
	}
	if resp == nil {
		t.Fatalf("bad: %#v", resp)
	}

	exp = map[string]interface{}{
		"id":               expID,
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
		"renewable":        !batch,
		"entity_id":        "",
		"type":             tokenType,
	}
	if periodic {
		exp["period"] = int64(3600)
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

	if diff := deep.Equal(resp.Data, exp); diff != nil {
		t.Fatal(diff)
	}

	// Test last_renewal_time functionality
	req = logical.TestRequest(t, logical.UpdateOperation, "renew")
	req.Data = map[string]interface{}{
		"token": expID,
	}
	resp, err = ts.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatalf("bad: %#v", resp)
	}
	if batch && !resp.IsError() || !batch && resp.IsError() {
		t.Fatalf("err: %v\nresp: %#v", err, resp)
	}

	req.Path = "lookup"
	resp, err = ts.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v\nresp: %#v", err, resp)
	}
	if resp == nil {
		t.Fatalf("bad: %#v", resp)
	}

	if !batch {
		if resp.Data["last_renewal_time"].(int64) == 0 {
			t.Fatalf("last_renewal_time was zero")
		}
	} else if _, ok := resp.Data["last_renewal_time"]; ok {
		t.Fatal("expected zero last renewal time")
	}
}

func TestTokenStore_HandleRequest_LookupSelf(t *testing.T) {
	c, _, root := TestCoreUnsealed(t)
	ts := c.tokenStore
	testMakeServiceTokenViaCore(t, c, root, "client", "3600s", []string{"foo"})

	req := logical.TestRequest(t, logical.ReadOperation, "lookup-self")
	req.ClientToken = "client"
	resp, err := ts.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v\nresp: %#v", err, resp)
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
		"type":             "service",
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
	root, err := ts.rootToken(namespace.RootContext(nil))
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
	err = exp.RegisterAuth(namespace.RootContext(nil), root, auth, "")
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
	resp, err := ts.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v\nresp: %#v", err, resp)
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

func TestTokenStore_HandleRequest_CreateToken_ExistingEntityAlias(t *testing.T) {
	core, _, root := TestCoreUnsealed(t)
	i := core.identityStore
	ctx := namespace.RootContext(nil)
	testPolicyName := "testpolicy"
	entityAliasName := "testentityalias"
	testRoleName := "test"

	// Create manually an entity
	resp, err := i.HandleRequest(ctx, &logical.Request{
		Path:      "entity",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"name":     "testentity",
			"policies": []string{testPolicyName},
		},
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v\nresp: %#v", err, resp)
	}
	entityID := resp.Data["id"].(string)

	// Find mount accessor
	resp, err = core.systemBackend.HandleRequest(namespace.RootContext(nil), &logical.Request{
		Path:      "auth",
		Operation: logical.ReadOperation,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
	}
	tokenMountAccessor := resp.Data["token/"].(map[string]interface{})["accessor"].(string)

	// Create manually an entity alias
	resp, err = i.HandleRequest(ctx, &logical.Request{
		Path:      "entity-alias",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"name":           entityAliasName,
			"canonical_id":   entityID,
			"mount_accessor": tokenMountAccessor,
		},
	})
	if err != nil {
		t.Fatalf("error handling request: %v", err)
	}

	// Create token role
	resp, err = core.HandleRequest(ctx, &logical.Request{
		Path:        "auth/token/roles/" + testRoleName,
		ClientToken: root,
		Operation:   logical.CreateOperation,
		Data: map[string]interface{}{
			"orphan":                 true,
			"period":                 "72h",
			"path_suffix":            "happenin",
			"bound_cidrs":            []string{"0.0.0.0/0"},
			"allowed_entity_aliases": []string{"test1", "test2", entityAliasName},
		},
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v\nresp: %#v", err, resp)
	}

	resp, err = core.HandleRequest(ctx, &logical.Request{
		Path:        "auth/token/create/" + testRoleName,
		Operation:   logical.UpdateOperation,
		ClientToken: root,
		Data: map[string]interface{}{
			"entity_alias": entityAliasName,
		},
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
	}
	if resp == nil {
		t.Fatal("expected a response")
	}
	if resp.Auth.EntityID != entityID {
		t.Fatalf("expected %q got %q", entityID, resp.Auth.EntityID)
	}

	policyFound := false
	for _, policy := range resp.Auth.IdentityPolicies {
		if policy == testPolicyName {
			policyFound = true
		}
	}
	if !policyFound {
		t.Fatalf("Policy %q not derived by entity but should be. Auth %#v", testPolicyName, resp.Auth)
	}
}

func TestTokenStore_HandleRequest_CreateToken_ExistingEntityAliasMixedCase(t *testing.T) {
	core, _, root := TestCoreUnsealed(t)
	i := core.identityStore
	ctx := namespace.RootContext(nil)
	testPolicyName := "testpolicy"
	entityAliasName := "tEStEntityAliaS"
	entityAliasNameLower := "testentityalias"
	testRoleName := "test"

	// Create manually an entity
	resp, err := i.HandleRequest(ctx, &logical.Request{
		Path:      "entity",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"name":     "testentity",
			"policies": []string{testPolicyName},
		},
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v\nresp: %#v", err, resp)
	}
	entityID := resp.Data["id"].(string)

	// Find mount accessor
	resp, err = core.systemBackend.HandleRequest(namespace.RootContext(nil), &logical.Request{
		Path:      "auth",
		Operation: logical.ReadOperation,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
	}
	tokenMountAccessor := resp.Data["token/"].(map[string]interface{})["accessor"].(string)

	// Create manually an entity alias
	resp, err = i.HandleRequest(ctx, &logical.Request{
		Path:      "entity-alias",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"name":           entityAliasName,
			"canonical_id":   entityID,
			"mount_accessor": tokenMountAccessor,
		},
	})
	if err != nil {
		t.Fatalf("error handling request: %v", err)
	}

	// Create token role
	resp, err = core.HandleRequest(ctx, &logical.Request{
		Path:        "auth/token/roles/" + testRoleName,
		ClientToken: root,
		Operation:   logical.CreateOperation,
		Data: map[string]interface{}{
			"orphan":                 true,
			"period":                 "72h",
			"path_suffix":            "happenin",
			"bound_cidrs":            []string{"0.0.0.0/0"},
			"allowed_entity_aliases": []string{"test1", "test2", entityAliasName},
		},
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v\nresp: %#v", err, resp)
	}

	respMixedCase, err := core.HandleRequest(ctx, &logical.Request{
		Path:        "auth/token/create/" + testRoleName,
		Operation:   logical.UpdateOperation,
		ClientToken: root,
		Data: map[string]interface{}{
			"entity_alias": entityAliasName,
		},
	})
	if err != nil {
		t.Fatalf("error handling request: %v", err)
	}
	if respMixedCase.Auth.EntityID != entityID {
		t.Fatalf("expected %q got %q", entityID, respMixedCase.Auth.EntityID)
	}

	// lowercase entity alias should match a mixed case alias
	respLowerCase, err := core.HandleRequest(ctx, &logical.Request{
		Path:        "auth/token/create/" + testRoleName,
		Operation:   logical.UpdateOperation,
		ClientToken: root,
		Data: map[string]interface{}{
			"entity_alias": entityAliasNameLower,
		},
	})
	if err != nil {
		t.Fatalf("error handling request: %v", err)
	}

	// A token created with the mixed case alias should return the same entity
	// id as the normal case response.
	if respLowerCase.Auth.EntityID != entityID {
		t.Fatalf("expected %q got %q", entityID, respLowerCase.Auth.EntityID)
	}
}

func TestTokenStore_HandleRequest_CreateToken_NonExistingEntityAlias(t *testing.T) {
	core, _, root := TestCoreUnsealed(t)
	i := core.identityStore
	ctx := namespace.RootContext(nil)
	entityAliasName := "testentityalias"
	testRoleName := "test"

	// Create token role
	resp, err := core.HandleRequest(ctx, &logical.Request{
		Path:        "auth/token/roles/" + testRoleName,
		ClientToken: root,
		Operation:   logical.CreateOperation,
		Data: map[string]interface{}{
			"period":                 "72h",
			"path_suffix":            "happenin",
			"bound_cidrs":            []string{"0.0.0.0/0"},
			"allowed_entity_aliases": []string{"test1", "test2", entityAliasName},
		},
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v\nresp: %#v", err, resp)
	}

	// Create token with non existing entity alias
	resp, err = core.HandleRequest(ctx, &logical.Request{
		Path:        "auth/token/create/" + testRoleName,
		Operation:   logical.UpdateOperation,
		ClientToken: root,
		Data: map[string]interface{}{
			"entity_alias": entityAliasName,
		},
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
	}
	if resp == nil {
		t.Fatal("expected a response")
	}

	// Read the new entity
	resp, err = i.HandleRequest(ctx, &logical.Request{
		Path:      "entity/id/" + resp.Auth.EntityID,
		Operation: logical.ReadOperation,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
	}

	// Get the attached alias information
	aliases := resp.Data["aliases"].([]interface{})
	if len(aliases) != 1 {
		t.Fatalf("expected only one alias but got %d; Aliases: %#v", len(aliases), aliases)
	}
	alias := &identity.Alias{}
	if err := mapstructure.Decode(aliases[0], alias); err != nil {
		t.Fatal(err)
	}

	// Validate
	if alias.Name != entityAliasName {
		t.Fatalf("alias name should be %q but is %q", entityAliasName, alias.Name)
	}
}

func TestTokenStore_HandleRequest_CreateToken_GlobPatternWildcardEntityAlias(t *testing.T) {
	core, _, root := TestCoreUnsealed(t)
	i := core.identityStore
	ctx := namespace.RootContext(nil)
	testRoleName := "test"

	tests := []struct {
		name        string
		globPattern string
		aliasName   string
	}{
		{
			name:        "prefix-asterisk",
			globPattern: "*-web",
			aliasName:   "department-web",
		},
		{
			name:        "suffix-asterisk",
			globPattern: "web-*",
			aliasName:   "web-department",
		},
		{
			name:        "middle-asterisk",
			globPattern: "web-*-web",
			aliasName:   "web-department-web",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create token role
			resp, err := core.HandleRequest(ctx, &logical.Request{
				Path:        "auth/token/roles/" + testRoleName,
				ClientToken: root,
				Operation:   logical.CreateOperation,
				Data: map[string]interface{}{
					"period":                 "72h",
					"path_suffix":            "happening",
					"bound_cidrs":            []string{"0.0.0.0/0"},
					"allowed_entity_aliases": []string{"test1", "test2", test.globPattern},
				},
			})
			if err != nil || (resp != nil && resp.IsError()) {
				t.Fatalf("err: %v\nresp: %#v", err, resp)
			}

			// Create token with non existing entity alias
			resp, err = core.HandleRequest(ctx, &logical.Request{
				Path:        "auth/token/create/" + testRoleName,
				Operation:   logical.UpdateOperation,
				ClientToken: root,
				Data: map[string]interface{}{
					"entity_alias": test.aliasName,
				},
			})
			if err != nil || (resp != nil && resp.IsError()) {
				t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
			}
			if resp == nil {
				t.Fatal("expected a response")
			}

			// Read the new entity
			resp, err = i.HandleRequest(ctx, &logical.Request{
				Path:      "entity/id/" + resp.Auth.EntityID,
				Operation: logical.ReadOperation,
			})
			if err != nil || (resp != nil && resp.IsError()) {
				t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
			}

			// Get the attached alias information
			aliases := resp.Data["aliases"].([]interface{})
			if len(aliases) != 1 {
				t.Fatalf("expected only one alias but got %d; Aliases: %#v", len(aliases), aliases)
			}
			alias := &identity.Alias{}
			if err := mapstructure.Decode(aliases[0], alias); err != nil {
				t.Fatal(err)
			}

			// Validate
			if alias.Name != test.aliasName {
				t.Fatalf("alias name should be %q but is %q", test.aliasName, alias.Name)
			}
		})
	}
}

func TestTokenStore_HandleRequest_CreateToken_NotAllowedEntityAlias(t *testing.T) {
	core, _, root := TestCoreUnsealed(t)
	i := core.identityStore
	ctx := namespace.RootContext(nil)
	testPolicyName := "testpolicy"
	entityAliasName := "testentityalias"
	testRoleName := "test"

	// Create manually an entity
	resp, err := i.HandleRequest(ctx, &logical.Request{
		Path:      "entity",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"name":     "testentity",
			"policies": []string{testPolicyName},
		},
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v\nresp: %#v", err, resp)
	}
	entityID := resp.Data["id"].(string)

	// Find mount accessor
	resp, err = core.systemBackend.HandleRequest(ctx, &logical.Request{
		Path:      "auth",
		Operation: logical.ReadOperation,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
	}
	tokenMountAccessor := resp.Data["token/"].(map[string]interface{})["accessor"].(string)

	// Create manually an entity alias
	resp, err = i.HandleRequest(ctx, &logical.Request{
		Path:      "entity-alias",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"name":           entityAliasName,
			"canonical_id":   entityID,
			"mount_accessor": tokenMountAccessor,
		},
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v\nresp: %#v", err, resp)
	}

	// Create token role
	resp, err = core.HandleRequest(ctx, &logical.Request{
		Path:        "auth/token/roles/" + testRoleName,
		ClientToken: root,
		Operation:   logical.CreateOperation,
		Data: map[string]interface{}{
			"period":                 "72h",
			"allowed_entity_aliases": []string{"test1", "test2", "testentityaliasn"},
		},
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v\nresp: %#v", err, resp)
	}

	resp, _ = core.HandleRequest(ctx, &logical.Request{
		Path:        "auth/token/create/" + testRoleName,
		Operation:   logical.UpdateOperation,
		ClientToken: root,
		Data: map[string]interface{}{
			"entity_alias": entityAliasName,
		},
	})
	if resp == nil || resp.Data == nil {
		t.Fatal("expected a response")
	}
	if resp.Data["error"] != "invalid 'entity_alias' value" {
		t.Fatalf("wrong error returned. Err: %s", resp.Data["error"])
	}
}

func TestTokenStore_HandleRequest_CreateToken_NoRoleEntityAlias(t *testing.T) {
	core, _, root := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	entityAliasName := "testentityalias"

	resp, _ := core.HandleRequest(ctx, &logical.Request{
		Path:        "auth/token/create",
		Operation:   logical.UpdateOperation,
		ClientToken: root,
		Data: map[string]interface{}{
			"entity_alias": entityAliasName,
		},
	})
	if resp == nil || resp.Data == nil {
		t.Fatal("expected a response")
	}
	if resp.Data["error"] != "'entity_alias' is only allowed in combination with token role" {
		t.Fatalf("wrong error returned. Err: %s", resp.Data["error"])
	}
}

func TestTokenStore_HandleRequest_RenewSelf(t *testing.T) {
	exp := mockExpiration(t)
	ts := exp.tokenStore

	// Create new token
	root, err := ts.rootToken(namespace.RootContext(nil))
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
	err = exp.RegisterAuth(namespace.RootContext(nil), root, auth, "")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Get the original expire time to compare
	originalExpire := auth.ExpirationTime()

	beforeRenew := time.Now()
	req := logical.TestRequest(t, logical.UpdateOperation, "renew-self")
	req.ClientToken = auth.ClientToken
	req.Data["increment"] = "3600s"
	resp, err := ts.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v\nresp: %#v", err, resp)
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
	core, _, root := TestCoreUnsealed(t)

	req := logical.TestRequest(t, logical.ReadOperation, "auth/token/roles/test")
	req.ClientToken = root

	resp, err := core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v\nresp: %#v", err, resp)
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
		"bound_cidrs":      []string{"0.0.0.0/0"},
		"explicit_max_ttl": "2h",
		"token_num_uses":   123,
	}

	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v\nresp: %#v", err, resp)
	}

	if resp != nil {
		t.Fatalf("expected a nil response")
	}

	req.Operation = logical.ReadOperation
	req.Data = map[string]interface{}{}

	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v\nresp: %#v", err, resp)
	}

	if resp == nil {
		t.Fatalf("got a nil response")
	}

	expected := map[string]interface{}{
		"name":                     "test",
		"orphan":                   true,
		"token_period":             int64(259200),
		"period":                   int64(259200),
		"allowed_policies":         []string{"test1", "test2"},
		"disallowed_policies":      []string{},
		"allowed_policies_glob":    []string{},
		"disallowed_policies_glob": []string{},
		"path_suffix":              "happenin",
		"explicit_max_ttl":         int64(7200),
		"token_explicit_max_ttl":   int64(7200),
		"renewable":                true,
		"token_type":               "default-service",
		"token_num_uses":           123,
		"allowed_entity_aliases":   []string(nil),
		"token_no_default_policy":  false,
	}

	if resp.Data["bound_cidrs"].([]*sockaddr.SockAddrMarshaler)[0].String() != "0.0.0.0/0" {
		t.Fatal("unexpected bound cidrs")
	}
	delete(resp.Data, "bound_cidrs")
	if resp.Data["token_bound_cidrs"].([]*sockaddr.SockAddrMarshaler)[0].String() != "0.0.0.0/0" {
		t.Fatal("unexpected bound cidrs")
	}
	delete(resp.Data, "token_bound_cidrs")

	if diff := deep.Equal(expected, resp.Data); diff != nil {
		t.Fatal(diff)
	}

	// Now test updating; this should be set to an UpdateOperation
	// automatically due to the existence check
	req.Operation = logical.CreateOperation
	req.Data = map[string]interface{}{
		"period":                  "79h",
		"allowed_policies":        "test3",
		"path_suffix":             "happenin",
		"renewable":               false,
		"explicit_max_ttl":        "80h",
		"token_num_uses":          0,
		"token_no_default_policy": true,
	}

	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v\nresp: %#v", err, resp)
	}

	if resp != nil {
		t.Fatalf("expected a nil response")
	}

	req.Operation = logical.ReadOperation
	req.Data = map[string]interface{}{}

	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v\nresp: %#v", err, resp)
	}
	if resp == nil {
		t.Fatalf("got a nil response")
	}

	expected = map[string]interface{}{
		"name":                     "test",
		"orphan":                   true,
		"period":                   int64(284400),
		"token_period":             int64(284400),
		"allowed_policies":         []string{"test3"},
		"disallowed_policies":      []string{},
		"allowed_policies_glob":    []string{},
		"disallowed_policies_glob": []string{},
		"path_suffix":              "happenin",
		"token_explicit_max_ttl":   int64(288000),
		"explicit_max_ttl":         int64(288000),
		"renewable":                false,
		"token_type":               "default-service",
		"allowed_entity_aliases":   []string(nil),
		"token_no_default_policy":  true,
	}

	if resp.Data["bound_cidrs"].([]*sockaddr.SockAddrMarshaler)[0].String() != "0.0.0.0/0" {
		t.Fatal("unexpected bound cidrs")
	}
	delete(resp.Data, "bound_cidrs")
	if resp.Data["token_bound_cidrs"].([]*sockaddr.SockAddrMarshaler)[0].String() != "0.0.0.0/0" {
		t.Fatal("unexpected bound cidrs")
	}
	delete(resp.Data, "token_bound_cidrs")

	if diff := deep.Equal(expected, resp.Data); diff != nil {
		t.Fatal(diff)
	}

	// Now set explicit max ttl and clear the period
	req.Operation = logical.CreateOperation
	req.Data = map[string]interface{}{
		"explicit_max_ttl": "5",
		"period":           "0s",
	}
	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v\nresp: %#v", err, resp)
	}

	req.Operation = logical.ReadOperation
	req.Data = map[string]interface{}{}

	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v\nresp: %#v", err, resp)
	}
	if resp == nil {
		t.Fatalf("got a nil response")
	}

	expected = map[string]interface{}{
		"name":                     "test",
		"orphan":                   true,
		"explicit_max_ttl":         int64(5),
		"token_explicit_max_ttl":   int64(5),
		"allowed_policies":         []string{"test3"},
		"disallowed_policies":      []string{},
		"allowed_policies_glob":    []string{},
		"disallowed_policies_glob": []string{},
		"path_suffix":              "happenin",
		"period":                   int64(0),
		"token_period":             int64(0),
		"renewable":                false,
		"token_type":               "default-service",
		"allowed_entity_aliases":   []string(nil),
		"token_no_default_policy":  true,
	}

	if resp.Data["bound_cidrs"].([]*sockaddr.SockAddrMarshaler)[0].String() != "0.0.0.0/0" {
		t.Fatal("unexpected bound cidrs")
	}
	delete(resp.Data, "bound_cidrs")
	if resp.Data["token_bound_cidrs"].([]*sockaddr.SockAddrMarshaler)[0].String() != "0.0.0.0/0" {
		t.Fatal("unexpected bound cidrs")
	}
	delete(resp.Data, "token_bound_cidrs")

	if diff := deep.Equal(expected, resp.Data); diff != nil {
		t.Fatal(diff)
	}

	// Update path_suffix and bound_cidrs with empty values
	req.Operation = logical.CreateOperation
	req.Data = map[string]interface{}{
		"path_suffix":             "",
		"bound_cidrs":             []string{},
		"token_no_default_policy": false,
	}
	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v\nresp: %#v", err, resp)
	}

	req.Operation = logical.ReadOperation
	req.Data = map[string]interface{}{}

	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v\nresp: %#v", err, resp)
	}
	if resp == nil {
		t.Fatalf("got a nil response")
	}

	expected = map[string]interface{}{
		"name":                     "test",
		"orphan":                   true,
		"token_explicit_max_ttl":   int64(5),
		"explicit_max_ttl":         int64(5),
		"allowed_policies":         []string{"test3"},
		"disallowed_policies":      []string{},
		"allowed_policies_glob":    []string{},
		"disallowed_policies_glob": []string{},
		"path_suffix":              "",
		"period":                   int64(0),
		"token_period":             int64(0),
		"renewable":                false,
		"token_type":               "default-service",
		"allowed_entity_aliases":   []string(nil),
		"token_no_default_policy":  false,
	}

	if diff := deep.Equal(expected, resp.Data); diff != nil {
		t.Fatal(diff)
	}

	req.Operation = logical.ListOperation
	req.Path = "auth/token/roles"
	req.Data = map[string]interface{}{}
	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v\nresp: %#v", err, resp)
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
		t.Fatalf("expected \"test\", got %q", keys[0])
	}

	req.Operation = logical.DeleteOperation
	req.Path = "auth/token/roles/test"
	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v\nresp: %#v", err, resp)
	}
	if resp != nil {
		t.Fatalf("expected a nil response")
	}

	req.Operation = logical.ReadOperation
	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v\nresp: %#v", err, resp)
	}
	if resp != nil {
		t.Fatalf("expected a nil response")
	}
}

func TestTokenStore_RoleDisallowedPoliciesWithRoot(t *testing.T) {
	var resp *logical.Response
	var err error

	c, _, root := TestCoreUnsealed(t)
	ts := c.tokenStore

	// Don't set disallowed_policies. Verify that a read on the role does return a non-nil value.
	roleReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "roles/role1",
		Data: map[string]interface{}{
			"disallowed_policies": "root,testpolicy",
		},
		ClientToken: root,
	}
	resp, err = ts.HandleRequest(namespace.RootContext(nil), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%v", err, resp)
	}

	roleReq.Operation = logical.ReadOperation
	resp, err = ts.HandleRequest(namespace.RootContext(nil), roleReq)
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

	core, _, root := TestCoreUnsealed(t)
	ts := core.tokenStore
	ps := core.policyStore

	// Create 3 different policies
	policy, _ := ParseACLPolicy(namespace.RootNamespace, tokenCreationPolicy)
	policy.Name = "test1"
	if err := ps.SetPolicy(namespace.RootContext(nil), policy); err != nil {
		t.Fatal(err)
	}

	policy, _ = ParseACLPolicy(namespace.RootNamespace, tokenCreationPolicy)
	policy.Name = "test2"
	if err := ps.SetPolicy(namespace.RootContext(nil), policy); err != nil {
		t.Fatal(err)
	}

	policy, _ = ParseACLPolicy(namespace.RootNamespace, tokenCreationPolicy)
	policy.Name = "test3"
	if err := ps.SetPolicy(namespace.RootContext(nil), policy); err != nil {
		t.Fatal(err)
	}

	// Create roles with different disallowed_policies configuration
	req = logical.TestRequest(t, logical.UpdateOperation, "roles/test1")
	req.ClientToken = root
	req.Data = map[string]interface{}{
		"disallowed_policies": "test1",
	}
	resp, err = ts.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%v", err, resp)
	}

	req = logical.TestRequest(t, logical.UpdateOperation, "roles/test23")
	req.ClientToken = root
	req.Data = map[string]interface{}{
		"disallowed_policies": "test2,test3",
	}
	resp, err = ts.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%v", err, resp)
	}

	req = logical.TestRequest(t, logical.UpdateOperation, "roles/test123")
	req.ClientToken = root
	req.Data = map[string]interface{}{
		"disallowed_policies": "test1,test2,test3",
	}
	resp, err = ts.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%v", err, resp)
	}

	// policy containing a glob character in the non-glob disallowed_policies field
	req = logical.TestRequest(t, logical.UpdateOperation, "roles/testglobdisabled")
	req.ClientToken = root
	req.Data = map[string]interface{}{
		"disallowed_policies": "test*",
	}
	resp, err = ts.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%v", err, resp)
	}

	// Create a token that has all the policies defined above
	req = logical.TestRequest(t, logical.UpdateOperation, "create")
	req.ClientToken = root
	req.Data["policies"] = []string{"test1", "test2", "test3"}
	resp = testMakeTokenViaRequest(t, ts, req)
	if resp == nil || resp.Auth == nil {
		t.Fatal("got nil response")
	}
	if resp.Auth.ClientToken == "" {
		t.Fatalf("bad: ClientToken; resp:%#v", resp)
	}
	parentToken := resp.Auth.ClientToken

	// Test that the parent token's policies are rejected by disallowed_policies
	req = logical.TestRequest(t, logical.UpdateOperation, "create/test1")
	req.ClientToken = parentToken
	resp, err = ts.HandleRequest(namespace.RootContext(nil), req)
	if err == nil || resp != nil && !resp.IsError() {
		t.Fatalf("expected an error response, got %#v", resp)
	}

	req = logical.TestRequest(t, logical.UpdateOperation, "create/test23")
	req.ClientToken = parentToken
	resp, err = ts.HandleRequest(namespace.RootContext(nil), req)
	if err == nil || resp != nil && !resp.IsError() {
		t.Fatalf("expected an error response, got %#v", resp)
	}

	req = logical.TestRequest(t, logical.UpdateOperation, "create/test123")
	req.ClientToken = parentToken
	resp, err = ts.HandleRequest(namespace.RootContext(nil), req)
	if err == nil || resp != nil && !resp.IsError() {
		t.Fatalf("expected an error response, got %#v", resp)
	}

	// Disallowed should act as a blacklist so make sure we can still make
	// something with other policies in the request
	req = logical.TestRequest(t, logical.UpdateOperation, "create/test123")
	req.Data["policies"] = []string{"foo", "bar"}
	req.ClientToken = parentToken
	testMakeTokenViaRequest(t, ts, req)

	// Check to be sure 'test*' without globbing matches 'test*'
	req = logical.TestRequest(t, logical.UpdateOperation, "create/testglobdisabled")
	req.Data["policies"] = []string{"test*"}
	req.ClientToken = parentToken
	resp, err = ts.HandleRequest(namespace.RootContext(nil), req)
	if err == nil || resp != nil && !resp.IsError() {
		t.Fatalf("expected an error response, got %#v", resp)
	}

	// Check to be sure 'test*' without globbing doesn't match 'test1' or 'test'
	req = logical.TestRequest(t, logical.UpdateOperation, "create/testglobdisabled")
	req.Data["policies"] = []string{"test1", "test"}
	req.ClientToken = parentToken
	testMakeTokenViaRequest(t, ts, req)

	// Create a role to have 'default' policy disallowed
	req = logical.TestRequest(t, logical.UpdateOperation, "roles/default")
	req.ClientToken = root
	req.Data = map[string]interface{}{
		"disallowed_policies": "default",
	}
	resp, err = ts.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%v", err, resp)
	}

	req = logical.TestRequest(t, logical.UpdateOperation, "create/default")
	req.ClientToken = parentToken
	resp, err = ts.HandleRequest(namespace.RootContext(nil), req)
	if err == nil || resp != nil && !resp.IsError() {
		t.Fatal("expected an error response")
	}
}

func TestTokenStore_RoleAllowedPolicies(t *testing.T) {
	c, _, root := TestCoreUnsealed(t)
	ts := c.tokenStore

	req := logical.TestRequest(t, logical.UpdateOperation, "roles/test")
	req.ClientToken = root
	req.Data = map[string]interface{}{
		"allowed_policies": "test1,test2",
	}

	resp, err := ts.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v\nresp: %#v", err, resp)
	}
	if resp != nil {
		t.Fatalf("expected a nil response")
	}

	req.Data = map[string]interface{}{}

	req.Path = "create/test"
	req.Data["policies"] = []string{"foo"}
	resp, err = ts.HandleRequest(namespace.RootContext(nil), req)
	if err == nil {
		t.Fatalf("expected error")
	}

	req.Data["policies"] = []string{"test2"}
	resp = testMakeTokenViaRequest(t, ts, req)
	if resp.Auth.ClientToken == "" {
		t.Fatalf("bad: %#v", resp)
	}

	// test not glob matching when using allowed_policies instead of allowed_policies_glob
	req = logical.TestRequest(t, logical.UpdateOperation, "roles/testnoglob")
	req.ClientToken = root
	req.Data = map[string]interface{}{
		"allowed_policies": "test*",
	}

	resp, err = ts.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v\nresp: %#v", err, resp)
	}
	if resp != nil {
		t.Fatalf("expected a nil response")
	}

	req.Path = "create/testnoglob"
	req.Data["policies"] = []string{"test"}
	resp, err = ts.HandleRequest(namespace.RootContext(nil), req)
	if err == nil {
		t.Fatalf("expected error")
	}

	req.Data["policies"] = []string{"testfoo"}
	resp, err = ts.HandleRequest(namespace.RootContext(nil), req)
	if err == nil {
		t.Fatalf("expected error")
	}

	req.Data["policies"] = []string{"test*"}
	resp = testMakeTokenViaRequest(t, ts, req)
	if resp.Auth.ClientToken == "" {
		t.Fatalf("bad: %#v", resp)
	}

	// When allowed_policies is blank, should fall back to a subset of the parent policies
	req = logical.TestRequest(t, logical.UpdateOperation, "roles/test")
	req.ClientToken = root
	req.Data = map[string]interface{}{
		"allowed_policies": "",
	}
	resp, err = ts.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v\nresp: %#v", err, resp)
	}
	if resp != nil {
		t.Fatalf("expected a nil response")
	}

	req = logical.TestRequest(t, logical.UpdateOperation, "create")
	req.ClientToken = root
	req.Data["policies"] = []string{"test1", "test2", "test3"}
	resp = testMakeTokenViaRequest(t, ts, req)
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
	resp, err = ts.HandleRequest(namespace.RootContext(nil), req)
	if err == nil {
		t.Fatalf("expected error")
	}

	req.Data["policies"] = []string{"test2"}
	resp = testMakeTokenViaRequest(t, ts, req)
	if resp.Auth.ClientToken == "" {
		t.Fatalf("bad: %#v", resp)
	}

	delete(req.Data, "policies")
	resp = testMakeTokenViaRequest(t, ts, req)
	if resp.Auth.ClientToken == "" {
		t.Fatalf("bad: %#v", resp)
	}
	if !reflect.DeepEqual(resp.Auth.Policies, []string{"default", "test1", "test2", "test3"}) {
		t.Fatalf("bad: %#v", resp.Auth.Policies)
	}
}

func TestTokenStore_RoleDisallowedPoliciesGlob(t *testing.T) {
	var req *logical.Request
	var resp *logical.Response
	var err error

	core, _, root := TestCoreUnsealed(t)
	ts := core.tokenStore
	ps := core.policyStore

	// Create 4 different policies
	policy, _ := ParseACLPolicy(namespace.RootNamespace, tokenCreationPolicy)
	policy.Name = "test1"
	if err := ps.SetPolicy(namespace.RootContext(nil), policy); err != nil {
		t.Fatal(err)
	}

	policy, _ = ParseACLPolicy(namespace.RootNamespace, tokenCreationPolicy)
	policy.Name = "test2"
	if err := ps.SetPolicy(namespace.RootContext(nil), policy); err != nil {
		t.Fatal(err)
	}

	policy, _ = ParseACLPolicy(namespace.RootNamespace, tokenCreationPolicy)
	policy.Name = "test3"
	if err := ps.SetPolicy(namespace.RootContext(nil), policy); err != nil {
		t.Fatal(err)
	}

	policy, _ = ParseACLPolicy(namespace.RootNamespace, tokenCreationPolicy)
	policy.Name = "test3b"
	if err := ps.SetPolicy(namespace.RootContext(nil), policy); err != nil {
		t.Fatal(err)
	}

	// Create roles with different disallowed_policies configuration
	req = logical.TestRequest(t, logical.UpdateOperation, "roles/test1")
	req.ClientToken = root
	req.Data = map[string]interface{}{
		"disallowed_policies_glob": "test1",
	}
	resp, err = ts.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%v", err, resp)
	}

	req = logical.TestRequest(t, logical.UpdateOperation, "roles/testnot23")
	req.ClientToken = root
	req.Data = map[string]interface{}{
		"disallowed_policies_glob": "test2,test3*",
	}
	resp, err = ts.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%v", err, resp)
	}

	// Create a token that has all the policies defined above
	req = logical.TestRequest(t, logical.UpdateOperation, "create")
	req.ClientToken = root
	req.Data["policies"] = []string{"test1", "test2", "test3", "test3b"}
	resp = testMakeTokenViaRequest(t, ts, req)
	if resp == nil || resp.Auth == nil {
		t.Fatal("got nil response")
	}
	if resp.Auth.ClientToken == "" {
		t.Fatalf("bad: ClientToken; resp:%#v", resp)
	}
	parentToken := resp.Auth.ClientToken

	// Test that the parent token's policies are rejected by disallowed_policies
	req = logical.TestRequest(t, logical.UpdateOperation, "create/test1")
	req.ClientToken = parentToken
	resp, err = ts.HandleRequest(namespace.RootContext(nil), req)
	if err == nil || resp != nil && !resp.IsError() {
		t.Fatalf("expected an error response, got %#v", resp)
	}
	req = logical.TestRequest(t, logical.UpdateOperation, "create/testnot23")
	req.ClientToken = parentToken
	resp, err = ts.HandleRequest(namespace.RootContext(nil), req)
	if err == nil || resp != nil && !resp.IsError() {
		t.Fatalf("expected an error response, got %#v", resp)
	}

	// Disallowed should act as a blacklist so make sure we can still make
	// something with other policies in the request
	req = logical.TestRequest(t, logical.UpdateOperation, "create/test1")
	req.Data["policies"] = []string{"foo", "bar"}
	req.ClientToken = parentToken
	testMakeTokenViaRequest(t, ts, req)

	// Check to be sure 'test3*' matches 'test3'
	req = logical.TestRequest(t, logical.UpdateOperation, "create/testnot23")
	req.Data["policies"] = []string{"test3"}
	req.ClientToken = parentToken
	resp, err = ts.HandleRequest(namespace.RootContext(nil), req)
	if err == nil || resp != nil && !resp.IsError() {
		t.Fatalf("expected an error response, got %#v", resp)
	}

	// Check to be sure 'test3*' matches 'test3b'
	req = logical.TestRequest(t, logical.UpdateOperation, "create/testnot23")
	req.Data["policies"] = []string{"test3b"}
	req.ClientToken = parentToken
	resp, err = ts.HandleRequest(namespace.RootContext(nil), req)
	if err == nil || resp != nil && !resp.IsError() {
		t.Fatalf("expected an error response, got %#v", resp)
	}

	// Check that non-blacklisted policies still work
	req = logical.TestRequest(t, logical.UpdateOperation, "create/testnot23")
	req.Data["policies"] = []string{"test1"}
	req.ClientToken = parentToken
	testMakeTokenViaRequest(t, ts, req)

	// Create a role to have 'default' policy disallowed
	req = logical.TestRequest(t, logical.UpdateOperation, "roles/default")
	req.ClientToken = root
	req.Data = map[string]interface{}{
		"disallowed_policies_glob": "default",
	}
	resp, err = ts.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%v", err, resp)
	}

	req = logical.TestRequest(t, logical.UpdateOperation, "create/default")
	req.ClientToken = parentToken
	resp, err = ts.HandleRequest(namespace.RootContext(nil), req)
	if err == nil || resp != nil && !resp.IsError() {
		t.Fatal("expected an error response")
	}
}

func TestTokenStore_RoleAllowedPoliciesGlob(t *testing.T) {
	c, _, root := TestCoreUnsealed(t)
	ts := c.tokenStore

	// test literal matching works in allowed_policies_glob
	req := logical.TestRequest(t, logical.UpdateOperation, "roles/test")
	req.ClientToken = root
	req.Data = map[string]interface{}{
		"allowed_policies_glob": "test1,test2",
	}

	resp, err := ts.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v\nresp: %#v", err, resp)
	}
	if resp != nil {
		t.Fatalf("expected a nil response")
	}

	req.Data = map[string]interface{}{}

	req.Path = "create/test"
	req.Data["policies"] = []string{"foo"}
	resp, err = ts.HandleRequest(namespace.RootContext(nil), req)
	if err == nil {
		t.Fatalf("expected error")
	}

	req.Data["policies"] = []string{"test2"}
	resp = testMakeTokenViaRequest(t, ts, req)
	if resp.Auth.ClientToken == "" {
		t.Fatalf("bad: %#v", resp)
	}

	// test glob matching in allowed_policies_glob
	req = logical.TestRequest(t, logical.UpdateOperation, "roles/test")
	req.ClientToken = root
	req.Data = map[string]interface{}{
		"allowed_policies_glob": "test*",
	}

	resp, err = ts.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v\nresp: %#v", err, resp)
	}
	if resp != nil {
		t.Fatalf("expected a nil response")
	}

	req.Path = "create/test"
	req.Data["policies"] = []string{"footest"}
	resp, err = ts.HandleRequest(namespace.RootContext(nil), req)
	if err == nil {
		t.Fatalf("expected error")
	}

	req.Data["policies"] = []string{"testfoo", "test2", "test"}
	resp = testMakeTokenViaRequest(t, ts, req)
	if resp.Auth.ClientToken == "" {
		t.Fatalf("bad: %#v", resp)
	}
}

func TestTokenStore_RoleOrphan(t *testing.T) {
	c, _, root := TestCoreUnsealed(t)
	ts := c.tokenStore

	req := logical.TestRequest(t, logical.UpdateOperation, "roles/test")
	req.ClientToken = root
	req.Data = map[string]interface{}{
		"orphan": true,
	}

	resp, err := ts.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v\nresp: %#v", err, resp)
	}
	if resp != nil {
		t.Fatalf("expected a nil response")
	}

	req.Path = "create/test"
	resp, err = ts.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v\nresp: %#v", err, resp)
	}
	if resp.Auth.ClientToken == "" {
		t.Fatalf("bad: %#v", resp)
	}

	out, err := ts.Lookup(namespace.RootContext(nil), resp.Auth.ClientToken)
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
	c, _, root := TestCoreUnsealed(t)
	ts := c.tokenStore

	req := logical.TestRequest(t, logical.CreateOperation, "roles/test")
	req.ClientToken = root
	req.Data = map[string]interface{}{
		"path_suffix": "happenin",
	}

	resp, err := ts.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v\nresp: %#v", err, resp)
	}
	if resp != nil {
		t.Fatalf("expected a nil response")
	}

	req.Path = "create/test"
	req.Operation = logical.UpdateOperation
	resp, err = ts.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v\nresp: %#v", err, resp)
	}
	if resp.Auth.ClientToken == "" {
		t.Fatalf("bad: %#v", resp)
	}

	out, err := ts.Lookup(namespace.RootContext(nil), resp.Auth.ClientToken)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if out.Path != "auth/token/create/test/happenin" {
		t.Fatalf("expected role in path but did not find it")
	}
}

func TestTokenStore_RolePeriod(t *testing.T) {
	core, _, root := TestCoreUnsealed(t)

	core.defaultLeaseTTL = 10 * time.Second
	core.maxLeaseTTL = 10 * time.Second

	// Note: these requests are sent to Core since Core handles registration
	// with the expiration manager and we need the storage to be consistent

	req := logical.TestRequest(t, logical.UpdateOperation, "auth/token/roles/test")
	req.ClientToken = root
	req.Data = map[string]interface{}{
		"period": 5,
	}

	resp, err := core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v\nresp: %#v", err, resp)
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
		resp, err = core.HandleRequest(namespace.RootContext(nil), req)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("err: %v\nresp: %#v", err, resp)
		}
		if resp.Auth.ClientToken == "" {
			t.Fatalf("bad: %#v", resp)
		}

		req.ClientToken = resp.Auth.ClientToken
		req.Operation = logical.ReadOperation
		req.Path = "auth/token/lookup-self"
		resp, err = core.HandleRequest(namespace.RootContext(nil), req)
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
		resp, err = core.HandleRequest(namespace.RootContext(nil), req)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("err: %v\nresp: %#v", err, resp)
		}

		req.Operation = logical.ReadOperation
		req.Path = "auth/token/lookup-self"
		resp, err = core.HandleRequest(namespace.RootContext(nil), req)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		ttl = resp.Data["ttl"].(int64)
		if ttl > 8 {
			t.Fatalf("TTL too large: %d", ttl)
		}

		// Renewing should not have the increment increase since we've hit the
		// max
		req.Operation = logical.UpdateOperation
		req.Path = "auth/token/renew-self"
		req.Data = map[string]interface{}{
			"increment": 2,
		}
		resp, err = core.HandleRequest(namespace.RootContext(nil), req)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("err: %v\nresp: %#v", err, resp)
		}

		req.Operation = logical.ReadOperation
		req.Path = "auth/token/lookup-self"
		resp, err = core.HandleRequest(namespace.RootContext(nil), req)
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
		resp, err = core.HandleRequest(namespace.RootContext(nil), req)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("err: %v\nresp: %#v", err, resp)
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
		resp, err = core.HandleRequest(namespace.RootContext(nil), req)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		ttl := resp.Data["ttl"].(int64)
		if ttl > 5 {
			t.Fatalf("TTL too large (expected %d, got %d", 5, ttl)
		}

		// Let the TTL go down a bit to 3 seconds
		time.Sleep(3 * time.Second)

		req.Operation = logical.UpdateOperation
		req.Path = "auth/token/renew-self"
		req.Data = map[string]interface{}{
			"increment": 2,
		}
		resp, err = core.HandleRequest(namespace.RootContext(nil), req)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("err: %v\nresp: %#v", err, resp)
		}

		req.Operation = logical.ReadOperation
		req.Path = "auth/token/lookup-self"
		resp, err = core.HandleRequest(namespace.RootContext(nil), req)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		ttl = resp.Data["ttl"].(int64)
		if ttl > 5 {
			t.Fatalf("TTL too large (expected %d, got %d", 5, ttl)
		}
	}
}

func TestTokenStore_RoleExplicitMaxTTL(t *testing.T) {
	core, _, root := TestCoreUnsealed(t)

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

	resp, err := core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v\nresp: %#v", err, resp)
	}
	if resp == nil {
		t.Fatalf("expected a warning")
	}

	req.Operation = logical.UpdateOperation
	req.Path = "auth/token/create/test"
	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
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

	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v\nresp: %#v", err, resp)
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
		resp, err = core.HandleRequest(namespace.RootContext(nil), req)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("err: %v\nresp: %#v", err, resp)
		}
		if resp.Auth.ClientToken == "" {
			t.Fatalf("bad: %#v", resp)
		}

		req.ClientToken = resp.Auth.ClientToken
		req.Operation = logical.ReadOperation
		req.Path = "auth/token/lookup-self"
		resp, err = core.HandleRequest(namespace.RootContext(nil), req)
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
		resp, err = core.HandleRequest(namespace.RootContext(nil), req)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("err: %v\nresp: %#v", err, resp)
		}

		req.Operation = logical.ReadOperation
		req.Path = "auth/token/lookup-self"
		resp, err = core.HandleRequest(namespace.RootContext(nil), req)
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
		resp, err = core.HandleRequest(namespace.RootContext(nil), req)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("err: %v\nresp: %#v", err, resp)
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
		resp, err = core.HandleRequest(namespace.RootContext(nil), req)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		ttl := resp.Data["ttl"].(int64)
		if ttl > 10 {
			t.Fatalf("TTL too big")
		}
		// explicit max ttl is stored in the role so not returned here
		maxTTL := resp.Data["explicit_max_ttl"].(int64)
		if maxTTL != 0 {
			t.Fatalf("expected 0 for explicit max TTL, got %d", maxTTL)
		}

		// Let the TTL go down a bit to ~7 seconds (8 against explicit max)
		time.Sleep(2 * time.Second)

		req.Operation = logical.UpdateOperation
		req.Path = "auth/token/renew-self"
		req.Data = map[string]interface{}{
			"increment": 300,
		}
		resp, err = core.HandleRequest(namespace.RootContext(nil), req)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("err: %v\nresp: %#v", err, resp)
		}

		req.Operation = logical.ReadOperation
		req.Path = "auth/token/lookup-self"
		resp, err = core.HandleRequest(namespace.RootContext(nil), req)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		ttl = resp.Data["ttl"].(int64)
		if ttl > 8 {
			t.Fatalf("TTL too big: %d", ttl)
		}

		// Let the TTL go down a bit more to ~5 seconds (6 against explicit max)
		time.Sleep(2 * time.Second)

		req.Operation = logical.UpdateOperation
		req.Path = "auth/token/renew-self"
		req.Data = map[string]interface{}{
			"increment": 300,
		}
		resp, err = core.HandleRequest(namespace.RootContext(nil), req)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("err: %v\nresp: %#v", err, resp)
		}

		req.Operation = logical.ReadOperation
		req.Path = "auth/token/lookup-self"
		resp, err = core.HandleRequest(namespace.RootContext(nil), req)
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
		resp, err = core.HandleRequest(namespace.RootContext(nil), req)
		if err == nil {
			t.Fatalf("expected error")
		}

		time.Sleep(2 * time.Second)

		req.Operation = logical.ReadOperation
		req.Path = "auth/token/lookup-self"
		resp, err = core.HandleRequest(namespace.RootContext(nil), req)
		if resp != nil && err == nil {
			t.Fatalf("expected error, response is %#v", *resp)
		}
		if err == nil {
			t.Fatalf("expected error")
		}
	}
}

func TestTokenStore_RoleTokenFields(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	// c, _, root := TestCoreUnsealed(t)
	ts := c.tokenStore
	rootContext := namespace.RootContext(context.Background())

	boundCIDRs, err := parseutil.ParseAddrs([]string{"127.0.0.1/32"})
	if err != nil {
		t.Fatal(err)
	}

	// First test the upgrade case. Create a role with values and ensure they
	// are reflected properly on read.
	{
		roleEntry := &tsRoleEntry{
			Name: "test",
			TokenParams: tokenutil.TokenParams{
				TokenType: logical.TokenTypeBatch,
			},
			Period:         time.Second,
			ExplicitMaxTTL: time.Hour,
		}
		roleEntry.BoundCIDRs = boundCIDRs
		ns := namespace.RootNamespace
		jsonEntry, err := logical.StorageEntryJSON("test", roleEntry)
		if err != nil {
			t.Fatal(err)
		}
		if err := ts.rolesView(ns).Put(rootContext, jsonEntry); err != nil {
			t.Fatal(err)
		}
		// Read it back
		roleEntry, err = ts.tokenStoreRole(rootContext, "test")
		if err != nil {
			t.Fatal(err)
		}
		expRoleEntry := &tsRoleEntry{
			Name: "test",
			TokenParams: tokenutil.TokenParams{
				TokenPeriod:         time.Second,
				TokenExplicitMaxTTL: time.Hour,
				TokenBoundCIDRs:     boundCIDRs,
				TokenType:           logical.TokenTypeBatch,
			},
			Period:         time.Second,
			ExplicitMaxTTL: time.Hour,
			BoundCIDRs:     boundCIDRs,
		}
		if diff := deep.Equal(expRoleEntry, roleEntry); diff != nil {
			t.Fatal(diff)
		}
	}

	// Now, read that back through the API and verify we see what we expect
	{
		req := logical.TestRequest(t, logical.ReadOperation, "roles/test")
		resp, err := ts.HandleRequest(rootContext, req)
		if err != nil {
			t.Fatalf("err: %v", err)
		}

		expected := map[string]interface{}{
			"name":                     "test",
			"orphan":                   false,
			"period":                   int64(1),
			"token_period":             int64(1),
			"allowed_policies":         []string(nil),
			"disallowed_policies":      []string(nil),
			"allowed_policies_glob":    []string(nil),
			"disallowed_policies_glob": []string(nil),
			"path_suffix":              "",
			"token_explicit_max_ttl":   int64(3600),
			"explicit_max_ttl":         int64(3600),
			"renewable":                false,
			"token_type":               "batch",
			"allowed_entity_aliases":   []string(nil),
			"token_no_default_policy":  false,
		}

		if resp.Data["bound_cidrs"].([]*sockaddr.SockAddrMarshaler)[0].String() != "127.0.0.1" {
			t.Fatalf("unexpected bound cidrs: %s", resp.Data["bound_cidrs"].([]*sockaddr.SockAddrMarshaler)[0].String())
		}
		delete(resp.Data, "bound_cidrs")
		if resp.Data["token_bound_cidrs"].([]*sockaddr.SockAddrMarshaler)[0].String() != "127.0.0.1" {
			t.Fatalf("unexpected token bound cidrs: %s", resp.Data["token_bound_cidrs"].([]*sockaddr.SockAddrMarshaler)[0].String())
		}
		delete(resp.Data, "token_bound_cidrs")

		if diff := deep.Equal(expected, resp.Data); diff != nil {
			t.Fatal(diff)
		}
	}

	// Put values in just the old locations, but through the API
	{
		req := logical.TestRequest(t, logical.UpdateOperation, "roles/test")
		req.Data = map[string]interface{}{
			"explicit_max_ttl": 7200,
			"token_type":       "default-service",
			"period":           5,
			"bound_cidrs":      boundCIDRs[0].String(),
		}

		resp, err := ts.HandleRequest(rootContext, req)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("err: %v\nresp: %#v", err, resp)
		}
		if resp != nil {
			t.Fatalf("expected a nil response")
		}

		req = logical.TestRequest(t, logical.ReadOperation, "roles/test")
		resp, err = ts.HandleRequest(rootContext, req)
		if err != nil {
			t.Fatalf("err: %v", err)
		}

		expected := map[string]interface{}{
			"name":                     "test",
			"orphan":                   false,
			"period":                   int64(5),
			"token_period":             int64(5),
			"allowed_policies":         []string(nil),
			"disallowed_policies":      []string(nil),
			"allowed_policies_glob":    []string(nil),
			"disallowed_policies_glob": []string(nil),
			"path_suffix":              "",
			"token_explicit_max_ttl":   int64(7200),
			"explicit_max_ttl":         int64(7200),
			"renewable":                false,
			"token_type":               "default-service",
			"allowed_entity_aliases":   []string(nil),
			"token_no_default_policy":  false,
		}

		if resp.Data["bound_cidrs"].([]*sockaddr.SockAddrMarshaler)[0].String() != "127.0.0.1" {
			t.Fatalf("unexpected bound cidrs: %s", resp.Data["bound_cidrs"].([]*sockaddr.SockAddrMarshaler)[0].String())
		}
		delete(resp.Data, "bound_cidrs")
		if resp.Data["token_bound_cidrs"].([]*sockaddr.SockAddrMarshaler)[0].String() != "127.0.0.1" {
			t.Fatalf("unexpected token bound cidrs: %s", resp.Data["token_bound_cidrs"].([]*sockaddr.SockAddrMarshaler)[0].String())
		}
		delete(resp.Data, "token_bound_cidrs")

		if diff := deep.Equal(expected, resp.Data); diff != nil {
			t.Fatal(diff)
		}
	}
	// Same thing for just the new locations
	{
		req := logical.TestRequest(t, logical.UpdateOperation, "roles/test")
		req.Data = map[string]interface{}{
			"token_explicit_max_ttl": 5200,
			"token_type":             "default-service",
			"token_period":           7,
			"token_bound_cidrs":      boundCIDRs[0].String(),
		}

		resp, err := ts.HandleRequest(rootContext, req)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("err: %v\nresp: %#v", err, resp)
		}
		if resp != nil {
			t.Fatalf("expected a nil response")
		}

		req = logical.TestRequest(t, logical.ReadOperation, "roles/test")
		resp, err = ts.HandleRequest(rootContext, req)
		if err != nil {
			t.Fatalf("err: %v", err)
		}

		expected := map[string]interface{}{
			"name":                     "test",
			"orphan":                   false,
			"period":                   int64(0),
			"token_period":             int64(7),
			"allowed_policies":         []string(nil),
			"disallowed_policies":      []string(nil),
			"allowed_policies_glob":    []string(nil),
			"disallowed_policies_glob": []string(nil),
			"path_suffix":              "",
			"token_explicit_max_ttl":   int64(5200),
			"explicit_max_ttl":         int64(0),
			"renewable":                false,
			"token_type":               "default-service",
			"allowed_entity_aliases":   []string(nil),
			"token_no_default_policy":  false,
		}

		if resp.Data["token_bound_cidrs"].([]*sockaddr.SockAddrMarshaler)[0].String() != "127.0.0.1" {
			t.Fatalf("unexpected token bound cidrs: %s", resp.Data["token_bound_cidrs"].([]*sockaddr.SockAddrMarshaler)[0].String())
		}
		delete(resp.Data, "token_bound_cidrs")

		if diff := deep.Equal(expected, resp.Data); diff != nil {
			t.Fatal(diff)
		}
	}
	// Put values in both locations
	{
		req := logical.TestRequest(t, logical.UpdateOperation, "roles/test")
		req.Data = map[string]interface{}{
			"token_explicit_max_ttl": 7200,
			"explicit_max_ttl":       5200,
			"token_type":             "service",
			"token_period":           5,
			"period":                 1,
			"token_bound_cidrs":      boundCIDRs[0].String(),
			"bound_cidrs":            boundCIDRs[0].String(),
		}

		resp, err := ts.HandleRequest(rootContext, req)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("err: %v\nresp: %#v", err, resp)
		}
		if resp == nil {
			t.Fatalf("expected a non-nil response")
		}
		if len(resp.Warnings) != 3 {
			t.Fatalf("expected 3 warnings, got %#v", resp.Warnings)
		}

		req = logical.TestRequest(t, logical.ReadOperation, "roles/test")
		resp, err = ts.HandleRequest(rootContext, req)
		if err != nil {
			t.Fatalf("err: %v", err)
		}

		expected := map[string]interface{}{
			"name":                     "test",
			"orphan":                   false,
			"period":                   int64(0),
			"token_period":             int64(5),
			"allowed_policies":         []string(nil),
			"disallowed_policies":      []string(nil),
			"allowed_policies_glob":    []string(nil),
			"disallowed_policies_glob": []string(nil),
			"path_suffix":              "",
			"token_explicit_max_ttl":   int64(7200),
			"explicit_max_ttl":         int64(0),
			"renewable":                false,
			"token_type":               "service",
			"allowed_entity_aliases":   []string(nil),
			"token_no_default_policy":  false,
		}

		if resp.Data["token_bound_cidrs"].([]*sockaddr.SockAddrMarshaler)[0].String() != "127.0.0.1" {
			t.Fatalf("unexpected token bound cidrs: %s", resp.Data["token_bound_cidrs"].([]*sockaddr.SockAddrMarshaler)[0].String())
		}
		delete(resp.Data, "token_bound_cidrs")

		if diff := deep.Equal(expected, resp.Data); diff != nil {
			t.Fatal(diff)
		}
	}
}

func TestTokenStore_Periodic(t *testing.T) {
	core, _, root := TestCoreUnsealed(t)

	core.defaultLeaseTTL = 10 * time.Second
	core.maxLeaseTTL = 10 * time.Second

	// Note: these requests are sent to Core since Core handles registration
	// with the expiration manager and we need the storage to be consistent

	req := logical.TestRequest(t, logical.UpdateOperation, "auth/token/roles/test")
	req.ClientToken = root
	req.Data = map[string]interface{}{
		"period": 5,
	}

	resp, err := core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v\nresp: %#v", err, resp)
	}
	if resp != nil {
		t.Fatalf("expected a nil response")
	}

	// First make one directly and verify on renew it uses the period.
	{
		req.ClientToken = root
		req.Operation = logical.UpdateOperation
		req.Path = "auth/token/create"
		resp, err = core.HandleRequest(namespace.RootContext(nil), req)
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
		resp, err = core.HandleRequest(namespace.RootContext(nil), req)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		ttl := resp.Data["ttl"].(int64)
		if ttl > 5 {
			t.Fatalf("TTL too large (expected %d, got %d)", 5, ttl)
		}

		// Let the TTL go down a bit
		time.Sleep(2 * time.Second)

		req.Operation = logical.UpdateOperation
		req.Path = "auth/token/renew-self"
		req.Data = map[string]interface{}{
			"increment": 1,
		}
		resp, err = core.HandleRequest(namespace.RootContext(nil), req)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("err: %v\nresp: %#v", err, resp)
		}

		req.Operation = logical.ReadOperation
		req.Path = "auth/token/lookup-self"
		resp, err = core.HandleRequest(namespace.RootContext(nil), req)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		ttl = resp.Data["ttl"].(int64)
		if ttl > 5 {
			t.Fatalf("TTL too large (expected %d, got %d)", 5, ttl)
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
			"period": 5,
		}
		resp, err = core.HandleRequest(namespace.RootContext(nil), req)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("err: %v\nresp: %#v", err, resp)
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
		resp, err = core.HandleRequest(namespace.RootContext(nil), req)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		ttl := resp.Data["ttl"].(int64)
		if ttl < 4 || ttl > 5 {
			t.Fatalf("TTL bad (expected %d, got %d)", 4, ttl)
		}

		// Let the TTL go down a bit
		time.Sleep(2 * time.Second)

		req.Operation = logical.UpdateOperation
		req.Path = "auth/token/renew-self"
		req.Data = map[string]interface{}{
			"increment": 1,
		}
		resp, err = core.HandleRequest(namespace.RootContext(nil), req)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("err: %v\nresp: %#v", err, resp)
		}

		req.Operation = logical.ReadOperation
		req.Path = "auth/token/lookup-self"
		resp, err = core.HandleRequest(namespace.RootContext(nil), req)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		ttl = resp.Data["ttl"].(int64)
		if ttl > 5 {
			t.Fatalf("TTL bad (expected less than %d, got %d)", 5, ttl)
		}
	}
}

func testTokenStore_NumUses_ErrorCheckHelper(t *testing.T, resp *logical.Response, err error) {
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
}

func testTokenStore_NumUses_SelfLookupHelper(t *testing.T, core *Core, clientToken string, expectedNumUses int) {
	req := logical.TestRequest(t, logical.ReadOperation, "auth/token/lookup-self")
	req.ClientToken = clientToken
	resp, err := core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	// Just used the token, this should decrement the num_uses counter
	expectedNumUses = expectedNumUses - 1
	actualNumUses := resp.Data["num_uses"].(int)

	if actualNumUses != expectedNumUses {
		t.Fatalf("num_uses mismatch (expected %d, got %d)", expectedNumUses, actualNumUses)
	}
}

func TestTokenStore_NumUses(t *testing.T) {
	core, _, root := TestCoreUnsealed(t)
	roleNumUses := 10
	lesserNumUses := 5
	greaterNumUses := 15

	// Create a test role with limited token_num_uses
	req := logical.TestRequest(t, logical.UpdateOperation, "auth/token/roles/test-limited-uses")
	req.ClientToken = root
	req.Data = map[string]interface{}{
		"token_num_uses": roleNumUses,
	}
	resp, err := core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v\nresp: %#v", err, resp)
	}
	if resp != nil {
		t.Fatalf("expected a nil response")
	}

	// Create a test role with unlimited token_num_uses
	req.Path = "auth/token/roles/test-unlimited-uses"
	req.Data = map[string]interface{}{}
	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v\nresp: %#v", err, resp)
	}
	if resp != nil {
		t.Fatalf("expected a nil response")
	}

	// Generate some tokens from the test roles
	req.Path = "auth/token/create/test-limited-uses"

	// first token, num_uses is expected to come from the limited uses role
	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	testTokenStore_NumUses_ErrorCheckHelper(t, resp, err)
	noOverrideToken := resp.Auth.ClientToken

	// second token, override num_uses with a lesser value, this should become the value
	// applied to the token
	req.Data = map[string]interface{}{
		"num_uses": lesserNumUses,
	}
	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	testTokenStore_NumUses_ErrorCheckHelper(t, resp, err)
	lesserOverrideToken := resp.Auth.ClientToken

	// third token, override num_uses with a greater value, the value
	// applied to the token should come from the limited uses role
	req.Data = map[string]interface{}{
		"num_uses": greaterNumUses,
	}
	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	testTokenStore_NumUses_ErrorCheckHelper(t, resp, err)
	greaterOverrideToken := resp.Auth.ClientToken

	// fourth token, override num_uses with a zero value, a num_uses value of zero
	// has an internal meaning of unlimited so num_uses == 1 is actually less than
	// num_uses == 0. In this case, the lesser value of the limited-uses role should be applied.
	req.Data = map[string]interface{}{
		"num_uses": 0,
	}
	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	testTokenStore_NumUses_ErrorCheckHelper(t, resp, err)
	zeroOverrideToken := resp.Auth.ClientToken

	// fifth token, override num_uses with a value from a role that has unlimited num_uses. num_uses
	// should be the specified num_uses parameter at the create endpoint
	req.Path = "auth/token/create/test-unlimited-uses"
	req.Data = map[string]interface{}{
		"num_uses": lesserNumUses,
	}
	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	testTokenStore_NumUses_ErrorCheckHelper(t, resp, err)
	unlimitedRoleOverrideToken := resp.Auth.ClientToken

	testTokenStore_NumUses_SelfLookupHelper(t, core, noOverrideToken, roleNumUses)
	testTokenStore_NumUses_SelfLookupHelper(t, core, lesserOverrideToken, lesserNumUses)
	testTokenStore_NumUses_SelfLookupHelper(t, core, greaterOverrideToken, roleNumUses)
	testTokenStore_NumUses_SelfLookupHelper(t, core, zeroOverrideToken, roleNumUses)
	testTokenStore_NumUses_SelfLookupHelper(t, core, unlimitedRoleOverrideToken, lesserNumUses)
}

func TestTokenStore_Periodic_ExplicitMax(t *testing.T) {
	core, _, root := TestCoreUnsealed(t)

	core.defaultLeaseTTL = 10 * time.Second
	core.maxLeaseTTL = 10 * time.Second

	// Note: these requests are sent to Core since Core handles registration
	// with the expiration manager and we need the storage to be consistent

	req := logical.TestRequest(t, logical.UpdateOperation, "auth/token/roles/test")
	req.ClientToken = root
	req.Data = map[string]interface{}{
		"period": 5,
	}

	resp, err := core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v\nresp: %#v", err, resp)
	}
	if resp != nil {
		t.Fatalf("expected a nil response")
	}

	// First make one directly and verify on renew it uses the period.
	{
		req.ClientToken = root
		req.Operation = logical.UpdateOperation
		req.Path = "auth/token/create"
		req.Data = map[string]interface{}{
			"period":           5,
			"explicit_max_ttl": 4,
		}
		resp, err = core.HandleRequest(namespace.RootContext(nil), req)
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
		resp, err = core.HandleRequest(namespace.RootContext(nil), req)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		ttl := resp.Data["ttl"].(int64)
		if ttl < 3 || ttl > 4 {
			t.Fatalf("TTL bad (expected %d, got %d)", 3, ttl)
		}

		// Let the TTL go down a bit
		time.Sleep(2 * time.Second)

		req.Operation = logical.UpdateOperation
		req.Path = "auth/token/renew-self"
		req.Data = map[string]interface{}{
			"increment": 76,
		}
		resp, err = core.HandleRequest(namespace.RootContext(nil), req)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("err: %v\nresp: %#v", err, resp)
		}

		req.Operation = logical.ReadOperation
		req.Path = "auth/token/lookup-self"
		resp, err = core.HandleRequest(namespace.RootContext(nil), req)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		ttl = resp.Data["ttl"].(int64)
		if ttl > 2 {
			t.Fatalf("TTL bad (expected less than %d, got %d)", 2, ttl)
		}
	}

	// Now we create a token against the role and also set the te value
	// directly. We should use the smaller of the two and be able to renew;
	// increment should be ignored as well.
	{
		req.Path = "auth/token/roles/test"
		req.ClientToken = root
		req.Operation = logical.UpdateOperation
		req.Data = map[string]interface{}{
			"period":           5,
			"explicit_max_ttl": 4,
		}

		resp, err := core.HandleRequest(namespace.RootContext(nil), req)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("err: %v\nresp: %#v", err, resp)
		}
		if resp != nil {
			t.Fatalf("expected a nil response")
		}

		req.ClientToken = root
		req.Operation = logical.UpdateOperation
		req.Path = "auth/token/create/test"
		resp, err = core.HandleRequest(namespace.RootContext(nil), req)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("err: %v\nresp: %#v", err, resp)
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
		resp, err = core.HandleRequest(namespace.RootContext(nil), req)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		ttl := resp.Data["ttl"].(int64)
		if ttl < 3 || ttl > 4 {
			t.Fatalf("TTL bad (expected %d, got %d)", 3, ttl)
		}

		// Let the TTL go down a bit
		time.Sleep(2 * time.Second)

		req.Operation = logical.UpdateOperation
		req.Path = "auth/token/renew-self"
		req.Data = map[string]interface{}{
			"increment": 1,
		}
		resp, err = core.HandleRequest(namespace.RootContext(nil), req)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("err: %v\nresp: %#v", err, resp)
		}

		req.Operation = logical.ReadOperation
		req.Path = "auth/token/lookup-self"
		resp, err = core.HandleRequest(namespace.RootContext(nil), req)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		ttl = resp.Data["ttl"].(int64)
		if ttl > 2 {
			t.Fatalf("TTL bad (expected less than %d, got %d)", 2, ttl)
		}
	}
}

func TestTokenStore_NoDefaultPolicy(t *testing.T) {
	var resp *logical.Response
	var err error

	core, _, root := TestCoreUnsealed(t)
	ts := core.tokenStore
	ps := core.policyStore
	policy, _ := ParseACLPolicy(namespace.RootNamespace, tokenCreationPolicy)
	policy.Name = "policy1"
	if err := ps.SetPolicy(namespace.RootContext(nil), policy); err != nil {
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
	resp = testMakeTokenViaRequest(t, ts, tokenReq)
	if !reflect.DeepEqual(resp.Auth.Policies, []string{"default", "policy1"}) {
		t.Fatalf("bad: policies: expected: [policy, default]; actual: %s", resp.Auth.Policies)
	}

	newToken := resp.Auth.ClientToken

	// Root token creates a token with desired policy, but also requests
	// that the token to not have 'default' policy. The resulting token
	// should not have 'default' policy on it.
	tokenData["no_default_policy"] = true
	resp = testMakeTokenViaRequest(t, ts, tokenReq)
	if !reflect.DeepEqual(resp.Auth.Policies, []string{"policy1"}) {
		t.Fatalf("bad: policies: expected: [policy1]; actual: %s", resp.Auth.Policies)
	}

	// A non-root token which has 'default' policy attached requests for a
	// child token. Child token should also have 'default' policy attached.
	tokenReq.ClientToken = newToken
	tokenReq.Data = nil
	resp = testMakeTokenViaRequest(t, ts, tokenReq)
	if !reflect.DeepEqual(resp.Auth.Policies, []string{"default", "policy1"}) {
		t.Fatalf("bad: policies: expected: [default policy1]; actual: %s", resp.Auth.Policies)
	}

	// A non-root token which has 'default' policy attached and period explicitly
	// set to its zero value requests for a child token. Child token should be
	// successfully created and have 'default' policy attached.
	tokenReq.Data = map[string]interface{}{
		"period": "0s",
	}
	resp = testMakeTokenViaRequest(t, ts, tokenReq)
	if !reflect.DeepEqual(resp.Auth.Policies, []string{"default", "policy1"}) {
		t.Fatalf("bad: policies: expected: [default policy1]; actual: %s", resp.Auth.Policies)
	}

	// A non-root token which has 'default' policy attached, request for a
	// child token to not have 'default' policy while not sending a list
	tokenReq.Data = map[string]interface{}{
		"no_default_policy": true,
	}
	resp = testMakeTokenViaRequest(t, ts, tokenReq)
	if !reflect.DeepEqual(resp.Auth.Policies, []string{"policy1"}) {
		t.Fatalf("bad: policies: expected: [policy1]; actual: %s", resp.Auth.Policies)
	}

	// In this case "default" shouldn't exist because we are not inheriting
	// parent policies
	tokenReq.Data = map[string]interface{}{
		"policies":          []string{"policy1"},
		"no_default_policy": true,
	}
	resp = testMakeTokenViaRequest(t, ts, tokenReq)
	if !reflect.DeepEqual(resp.Auth.Policies, []string{"policy1"}) {
		t.Fatalf("bad: policies: expected: [policy1]; actual: %s", resp.Auth.Policies)
	}

	// This is a non-root token which does not have 'default' policy
	// attached
	newToken = resp.Auth.ClientToken
	tokenReq.Data = nil
	tokenReq.ClientToken = newToken
	resp = testMakeTokenViaRequest(t, ts, tokenReq)
	if !reflect.DeepEqual(resp.Auth.Policies, []string{"policy1"}) {
		t.Fatalf("bad: policies: expected: [policy1]; actual: %s", resp.Auth.Policies)
	}

	roleReq := &logical.Request{
		ClientToken: root,
		Path:        "roles/role1",
		Operation:   logical.CreateOperation,
	}
	resp, err = ts.HandleRequest(namespace.RootContext(nil), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v, resp: %v", err, resp)
	}
	tokenReq.Path = "create/role1"
	tokenReq.Data = map[string]interface{}{
		"policies": []string{"policy1"},
	}
	resp, err = ts.HandleRequest(namespace.RootContext(nil), tokenReq)
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
	resp, err = ts.HandleRequest(namespace.RootContext(nil), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v, resp: %v", err, resp)
	}

	resp = testMakeTokenViaRequest(t, ts, tokenReq)
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
	resp, err = ts.HandleRequest(namespace.RootContext(nil), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v, resp: %v", err, resp)
	}

	resp = testMakeTokenViaRequest(t, ts, tokenReq)
	if !reflect.DeepEqual(resp.Auth.Policies, []string{"policy1"}) {
		t.Fatalf("bad: policies: expected: [policy1]; actual: %s", resp.Auth.Policies)
	}

	roleReq.Data = map[string]interface{}{
		"allowed_policies":    "",
		"disallowed_policies": "default",
	}
	resp, err = ts.HandleRequest(namespace.RootContext(nil), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v, resp: %v", err, resp)
	}

	resp = testMakeTokenViaRequest(t, ts, tokenReq)
	if !reflect.DeepEqual(resp.Auth.Policies, []string{"policy1"}) {
		t.Fatalf("bad: policies: expected: [policy1]; actual: %s", resp.Auth.Policies)
	}

	// Ensure that if default is in both allowed and disallowed, disallowed wins
	roleReq.Data = map[string]interface{}{
		"allowed_policies":    "default",
		"disallowed_policies": "default",
	}
	resp, err = ts.HandleRequest(namespace.RootContext(nil), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v, resp: %v", err, resp)
	}

	delete(tokenReq.Data, "policies")
	resp, err = ts.HandleRequest(namespace.RootContext(nil), tokenReq)
	if err == nil || (resp != nil && !resp.IsError()) {
		t.Fatalf("err: %v, resp: %v", err, resp)
	}
}

func TestTokenStore_AllowedDisallowedPolicies(t *testing.T) {
	var resp *logical.Response
	var err error

	c, _, root := TestCoreUnsealed(t)
	ts := c.tokenStore

	roleReq := &logical.Request{
		ClientToken: root,
		Path:        "roles/role1",
		Operation:   logical.CreateOperation,
		Data: map[string]interface{}{
			"allowed_policies":    "allowed1,allowed2",
			"disallowed_policies": "disallowed1,disallowed2",
		},
	}
	resp, err = ts.HandleRequest(namespace.RootContext(nil), roleReq)
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
	resp, err = ts.HandleRequest(namespace.RootContext(nil), tokenReq)
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
	resp, err = ts.HandleRequest(namespace.RootContext(nil), tokenReq)
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
	resp, err = ts.HandleRequest(namespace.RootContext(nil), tokenReq)
	if err == nil {
		t.Fatalf("expected an error")
	}

	roleReq.Operation = logical.UpdateOperation
	roleReq.Data = map[string]interface{}{
		"allowed_policies":    "allowed1,common",
		"disallowed_policies": "disallowed1,common",
	}
	resp, err = ts.HandleRequest(namespace.RootContext(nil), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v, resp: %v", err, resp)
	}

	tokenReq.Data = map[string]interface{}{
		"policies": []string{"allowed1", "common"},
	}
	resp, err = ts.HandleRequest(namespace.RootContext(nil), tokenReq)
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
	root, _ := exp.tokenStore.rootToken(namespace.RootContext(nil))

	tokenReq := &logical.Request{
		Path:        "create",
		ClientToken: root.ID,
		Operation:   logical.UpdateOperation,
		Data: map[string]interface{}{
			"num_uses": 1,
		},
	}
	resp, err = ts.HandleRequest(namespace.RootContext(nil), tokenReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%v", err, resp)
	}

	tut := resp.Auth.ClientToken
	saltTut, err := ts.SaltID(namespace.RootContext(nil), tut)
	if err != nil {
		t.Fatal(err)
	}
	te, err := ts.lookupInternal(namespace.RootContext(nil), saltTut, true, false)
	if err != nil {
		t.Fatal(err)
	}
	if te == nil {
		t.Fatal("nil entry")
	}
	if te.NumUses != 1 {
		t.Fatalf("bad: %d", te.NumUses)
	}

	te, err = ts.UseToken(namespace.RootContext(nil), te)
	if err != nil {
		t.Fatal(err)
	}
	if te == nil {
		t.Fatal("nil entry")
	}
	if te.NumUses != tokenRevocationPending {
		t.Fatalf("bad: %d", te.NumUses)
	}

	// Should return no entry because it's tainted
	te, err = ts.lookupInternal(namespace.RootContext(nil), saltTut, true, false)
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
	resp, err = ts.HandleRequest(namespace.RootContext(nil), req)
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
	te, err = ts.lookupInternal(namespace.RootContext(nil), saltTut, true, true)
	if err != nil {
		t.Fatal(err)
	}
	if te == nil {
		t.Fatal("nil entry")
	}
	if te.NumUses != tokenRevocationPending {
		t.Fatalf("bad: %d", te.NumUses)
	}

	origDestroyCubbyhole := ts.cubbyholeDestroyer

	ts.cubbyholeDestroyer = func(context.Context, *TokenStore, *logical.TokenEntry) error {
		return fmt.Errorf("keep it frosty")
	}

	err = ts.revokeInternal(namespace.RootContext(nil), saltTut, false)
	if err == nil {
		t.Fatalf("expected err")
	}

	// Since revocation failed we should still be able to get a token
	te, err = ts.lookupInternal(namespace.RootContext(nil), saltTut, true, true)
	if err != nil {
		t.Fatal(err)
	}
	if te == nil {
		t.Fatal("nil token entry")
	}

	// Check the race condition situation by making the process sleep
	ts.cubbyholeDestroyer = func(context.Context, *TokenStore, *logical.TokenEntry) error {
		time.Sleep(1 * time.Second)
		return fmt.Errorf("keep it frosty")
	}
	cubbyFuncLock.Unlock()

	errCh := make(chan error)
	defer close(errCh)
	go func() {
		cubbyFuncLock.RLock()
		err := ts.revokeInternal(namespace.RootContext(nil), saltTut, false)
		cubbyFuncLock.RUnlock()
		errCh <- err
	}()

	// Give time for the function to start and grab locks
	time.Sleep(200 * time.Millisecond)

	te, err = ts.lookupInternal(namespace.RootContext(nil), saltTut, true, true)
	if err != nil {
		t.Fatal(err)
	}
	if te == nil {
		t.Fatal("nil token entry")
	}

	err = <-errCh
	if err == nil {
		t.Fatal("expected error on ts.revokeInternal() in anonymous goroutine")
	}

	// Put back to normal
	cubbyFuncLock.Lock()
	defer cubbyFuncLock.Unlock()
	ts.cubbyholeDestroyer = origDestroyCubbyhole

	err = ts.revokeInternal(namespace.RootContext(nil), saltTut, false)
	if err != nil {
		t.Fatal(err)
	}

	te, err = ts.lookupInternal(namespace.RootContext(nil), saltTut, true, true)
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

	c, _, root := TestCoreUnsealed(t)
	ts := c.tokenStore

	// List the number of accessors. Since there is only root token
	// present, the list operation should return only one key.
	accessorListReq := &logical.Request{
		Operation:   logical.ListOperation,
		Path:        "accessors/",
		ClientToken: root,
	}
	resp, err = ts.HandleRequest(namespace.RootContext(nil), accessorListReq)
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
		resp := testMakeTokenViaRequest(t, ts, tokenReq)
		tut := resp.Auth.ClientToken

		// Creation of another token should end up with incrementing
		// the number of accessors
		// the storage
		resp, err = ts.HandleRequest(namespace.RootContext(nil), accessorListReq)
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

		saltedTut, err := ts.SaltID(namespace.RootContext(nil), tut)
		if err != nil {
			t.Fatal(err)
		}
		te, err := ts.lookupInternal(namespace.RootContext(nil), saltedTut, true, true)
		if err != nil {
			t.Fatalf("failed to lookup token: %v", err)
		}

		// Destroy the token index
		if ts.idView(namespace.RootNamespace).Delete(namespace.RootContext(nil), saltedTut); err != nil {
			t.Fatalf("failed to delete token entry: %v", err)
		}

		// Destroy the cubby space
		err = ts.cubbyholeDestroyer(namespace.RootContext(nil), ts, te)
		if err != nil {
			t.Fatalf("failed to destroyCubbyhole: %v", err)
		}

		// Leaking of accessor should have resulted in no change to the number
		// of accessors
		resp, err = ts.HandleRequest(namespace.RootContext(nil), accessorListReq)
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
	resp, err = ts.HandleRequest(namespace.RootContext(nil), tidyReq)
	if err != nil {
		t.Fatal(err)
	}
	if resp != nil && resp.IsError() {
		t.Fatalf("resp: %#v", resp)
	}
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%v", err, resp)
	}

	// Tidy runs async so give it time
	time.Sleep(10 * time.Second)

	// Tidy should have removed all the dangling accessor entries
	resp, err = ts.HandleRequest(namespace.RootContext(nil), accessorListReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%v", err, resp)
	}

	numberOfAccessors = len(resp.Data["keys"].([]string))
	if numberOfAccessors != 1 {
		t.Fatalf("bad: number of accessors. Expected: 1, Actual: %d", numberOfAccessors)
	}
}

// Create a set of tokens along with a child token for each of them, delete the
// token entry while leaking accessors, invoke tidy and check if the dangling
// accessor entry is getting removed and check if child tokens are still present
// and turned into orphan tokens.
func TestTokenStore_HandleTidy_parentCleanup(t *testing.T) {
	var resp *logical.Response
	var err error

	c, _, root := TestCoreUnsealed(t)
	ts := c.tokenStore

	// List the number of accessors. Since there is only root token
	// present, the list operation should return only one key.
	accessorListReq := &logical.Request{
		Operation:   logical.ListOperation,
		Path:        "accessors/",
		ClientToken: root,
	}
	resp, err = ts.HandleRequest(namespace.RootContext(nil), accessorListReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%v", err, resp)
	}

	numberOfAccessors := len(resp.Data["keys"].([]string))
	if numberOfAccessors != 1 {
		t.Fatalf("bad: number of accessors. Expected: 1, Actual: %d", numberOfAccessors)
	}

	for i := 1; i <= 100; i++ {
		// Create a token
		tokenReq := &logical.Request{
			Operation:   logical.UpdateOperation,
			Path:        "create",
			ClientToken: root,
			Data: map[string]interface{}{
				"policies": []string{"policy1"},
			},
		}
		resp := testMakeTokenViaRequest(t, ts, tokenReq)
		tut := resp.Auth.ClientToken

		// Create a child token
		tokenReq = &logical.Request{
			Operation:   logical.UpdateOperation,
			Path:        "create",
			ClientToken: tut,
			Data: map[string]interface{}{
				"policies": []string{"policy1"},
			},
		}
		testMakeTokenViaRequest(t, ts, tokenReq)

		// Creation of another token should end up with incrementing the number of
		// accessors the storage
		resp, err = ts.HandleRequest(namespace.RootContext(nil), accessorListReq)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("err:%v resp:%v", err, resp)
		}

		numberOfAccessors = len(resp.Data["keys"].([]string))
		if numberOfAccessors != (i*2)+1 {
			t.Fatalf("bad: number of accessors. Expected: %d, Actual: %d", i+1, numberOfAccessors)
		}

		// Revoke the token while leaking other items associated with the
		// token. Do this by doing what revokeSalted used to do before it was
		// fixed, i.e., by deleting the storage entry for token and its
		// cubbyhole and by not deleting its secondary index, its accessor and
		// associated leases.

		saltedTut, err := ts.SaltID(namespace.RootContext(nil), tut)
		if err != nil {
			t.Fatal(err)
		}
		te, err := ts.lookupInternal(namespace.RootContext(nil), saltedTut, true, true)
		if err != nil {
			t.Fatalf("failed to lookup token: %v", err)
		}

		// Destroy the token index
		if ts.idView(namespace.RootNamespace).Delete(namespace.RootContext(nil), saltedTut); err != nil {
			t.Fatalf("failed to delete token entry: %v", err)
		}

		// Destroy the cubby space
		err = ts.cubbyholeDestroyer(namespace.RootContext(nil), ts, te)
		if err != nil {
			t.Fatalf("failed to destroyCubbyhole: %v", err)
		}

		// Leaking of accessor should have resulted in no change to the number
		// of accessors
		resp, err = ts.HandleRequest(namespace.RootContext(nil), accessorListReq)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("err:%v resp:%v", err, resp)
		}

		numberOfAccessors = len(resp.Data["keys"].([]string))
		if numberOfAccessors != (i*2)+1 {
			t.Fatalf("bad: number of accessors. Expected: %d, Actual: %d", (i*2)+1, numberOfAccessors)
		}
	}

	tidyReq := &logical.Request{
		Path:        "tidy",
		Operation:   logical.UpdateOperation,
		ClientToken: root,
	}
	resp, err = ts.HandleRequest(namespace.RootContext(nil), tidyReq)
	if err != nil {
		t.Fatal(err)
	}
	if resp != nil && resp.IsError() {
		t.Fatalf("resp: %#v", resp)
	}
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%v", err, resp)
	}

	// Tidy runs async so give it time
	time.Sleep(10 * time.Second)

	// Tidy should have removed all the dangling accessor entries
	resp, err = ts.HandleRequest(namespace.RootContext(nil), accessorListReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%v", err, resp)
	}

	// The number of accessors should be equal to number of valid child tokens
	// (100) + the root token (1)
	keys := resp.Data["keys"].([]string)
	numberOfAccessors = len(keys)
	if numberOfAccessors != 101 {
		t.Fatalf("bad: number of accessors. Expected: 1, Actual: %d", numberOfAccessors)
	}

	req := logical.TestRequest(t, logical.UpdateOperation, "lookup-accessor")

	for _, accessor := range keys {
		req.Data = map[string]interface{}{
			"accessor": accessor,
		}

		resp, err := ts.HandleRequest(namespace.RootContext(nil), req)
		if err != nil {
			t.Fatalf("err: %s", err)
		}
		if resp.Data == nil {
			t.Fatalf("response should contain data")
		}
		// These tokens should now be orphaned
		if resp.Data["orphan"] != true {
			t.Fatalf("token should be orphan")
		}
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
	err = exp.router.Mount(noop, "prod/aws/", &MountEntry{UUID: meUUID, Accessor: "awsaccessor", namespace: namespace.RootNamespace}, view)
	if err != nil {
		t.Fatal(err)
	}

	// Create new token
	root, err := ts.rootToken(namespace.RootContext(nil))
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Create a new token
	req := logical.TestRequest(t, logical.UpdateOperation, "create")
	req.ClientToken = root.ID
	req.Data["policies"] = []string{"default"}

	resp, err := ts.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err: %v\nresp: %#v", err, resp)
	}

	auth := &logical.Auth{
		ClientToken: resp.Auth.ClientToken,
		LeaseOptions: logical.LeaseOptions{
			TTL:       time.Hour,
			Renewable: true,
		},
	}

	te := &logical.TokenEntry{
		Path:        resp.Auth.CreationPath,
		NamespaceID: namespace.RootNamespaceID,
	}

	err = exp.RegisterAuth(namespace.RootContext(nil), te, auth, "")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Verify token entry through lookup
	testTokenEntry, err := ts.Lookup(namespace.RootContext(nil), resp.Auth.ClientToken)
	if err != nil {
		t.Fatal(err)
	}
	if testTokenEntry == nil {
		t.Fatal("token entry was nil")
	}

	tut := resp.Auth.ClientToken

	req = &logical.Request{
		Path:        "prod/aws/foo",
		ClientToken: tut,
	}
	req.SetTokenEntry(testTokenEntry)

	resp = &logical.Response{
		Secret: &logical.Secret{
			LeaseOptions: logical.LeaseOptions{
				TTL: time.Hour,
			},
		},
	}

	leases := []string{}

	for i := 0; i < 10; i++ {
		leaseID, err := exp.Register(namespace.RootContext(nil), req, resp, "")
		if err != nil {
			t.Fatal(err)
		}
		leases = append(leases, leaseID)
	}

	sort.Strings(leases)

	te, err = ts.lookupInternal(namespace.RootContext(nil), tut, false, true)
	if err != nil {
		t.Fatalf("failed to lookup token: %v", err)
	}
	if te == nil {
		t.Fatal("got nil token entry")
	}

	storedLeases, err := exp.lookupLeasesByToken(namespace.RootContext(nil), te)
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(storedLeases)
	if !reflect.DeepEqual(leases, storedLeases) {
		t.Fatalf("bad: %#v vs %#v", leases, storedLeases)
	}

	// Now, delete the token entry. The leases should still exist.
	saltedTut, err := ts.SaltID(namespace.RootContext(nil), tut)
	if err != nil {
		t.Fatal(err)
	}
	te, err = ts.lookupInternal(namespace.RootContext(nil), saltedTut, true, true)
	if err != nil {
		t.Fatalf("failed to lookup token: %v", err)
	}
	if te == nil {
		t.Fatal("got nil token entry")
	}

	// Destroy the token index
	if ts.idView(namespace.RootNamespace).Delete(namespace.RootContext(nil), saltedTut); err != nil {
		t.Fatalf("failed to delete token entry: %v", err)
	}
	te, err = ts.lookupInternal(namespace.RootContext(nil), saltedTut, true, true)
	if err != nil {
		t.Fatalf("failed to lookup token: %v", err)
	}
	if te != nil {
		t.Fatal("got token entry")
	}

	// Verify leases still exist
	storedLeases, err = exp.lookupLeasesByToken(namespace.RootContext(nil), testTokenEntry)
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(storedLeases)
	if !reflect.DeepEqual(leases, storedLeases) {
		t.Fatalf("bad: %#v vs %#v", leases, storedLeases)
	}

	// Call tidy
	ts.handleTidy(namespace.RootContext(nil), &logical.Request{}, nil)

	time.Sleep(200 * time.Millisecond)

	// Verify leases are gone
	storedLeases, err = exp.lookupLeasesByToken(namespace.RootContext(nil), testTokenEntry)
	if err != nil {
		t.Fatal(err)
	}
	if len(storedLeases) > 0 {
		t.Fatal("found leases")
	}
}

func TestTokenStore_Batch_CannotCreateChildren(t *testing.T) {
	var resp *logical.Response

	core, _, root := TestCoreUnsealed(t)
	ts := core.tokenStore

	req := &logical.Request{
		Path:        "create",
		ClientToken: root,
		Operation:   logical.UpdateOperation,
		Data: map[string]interface{}{
			"policies": []string{"policy1"},
			"type":     "batch",
		},
	}
	resp = testMakeTokenViaRequest(t, ts, req)
	if !reflect.DeepEqual(resp.Auth.Policies, []string{"default", "policy1"}) {
		t.Fatalf("bad: policies: expected: [policy, default]; actual: %s", resp.Auth.Policies)
	}

	req.ClientToken = resp.Auth.ClientToken
	resp = testMakeTokenViaRequest(t, ts, req)
	if !resp.IsError() {
		t.Fatalf("bad: expected error, got %#v", *resp)
	}
}

func TestTokenStore_Batch_CannotRevoke(t *testing.T) {
	var resp *logical.Response
	var err error

	core, _, root := TestCoreUnsealed(t)
	ts := core.tokenStore

	req := &logical.Request{
		Path:        "create",
		ClientToken: root,
		Operation:   logical.UpdateOperation,
		Data: map[string]interface{}{
			"policies": []string{"policy1"},
			"type":     "batch",
		},
	}
	resp = testMakeTokenViaRequest(t, ts, req)
	if !reflect.DeepEqual(resp.Auth.Policies, []string{"default", "policy1"}) {
		t.Fatalf("bad: policies: expected: [policy, default]; actual: %s", resp.Auth.Policies)
	}

	req.Path = "revoke"
	req.Data["token"] = resp.Auth.ClientToken
	resp, err = ts.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatal(err)
	}
	if !resp.IsError() {
		t.Fatalf("bad: expected error, got %#v", *resp)
	}
}

func TestTokenStore_Batch_NoCubbyhole(t *testing.T) {
	var resp *logical.Response
	var err error

	core, _, root := TestCoreUnsealed(t)
	ts := core.tokenStore

	req := &logical.Request{
		Path:        "create",
		ClientToken: root,
		Operation:   logical.UpdateOperation,
		Data: map[string]interface{}{
			"policies": []string{"policy1"},
			"type":     "batch",
		},
	}
	resp = testMakeTokenViaRequest(t, ts, req)
	if !reflect.DeepEqual(resp.Auth.Policies, []string{"default", "policy1"}) {
		t.Fatalf("bad: policies: expected: [policy, default]; actual: %s", resp.Auth.Policies)
	}

	te, err := ts.Lookup(namespace.RootContext(nil), resp.Auth.ClientToken)
	if err != nil {
		t.Fatal(err)
	}

	req.Path = "cubbyhole/foo"
	req.Operation = logical.CreateOperation
	req.ClientToken = te.ID
	req.SetTokenEntry(te)
	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil && !errwrap.Contains(err, logical.ErrInvalidRequest.Error()) {
		t.Fatal(err)
	}
	if !resp.IsError() {
		t.Fatalf("bad: expected error, got %#v", *resp)
	}
}

func TestTokenStore_TokenID(t *testing.T) {
	t.Run("no custom ID provided", func(t *testing.T) {
		c, _, initToken := TestCoreUnsealed(t)
		ts := c.tokenStore

		// Ensure that a regular service token has a consts.ServiceTokenPrefix prefix
		resp, err := ts.HandleRequest(namespace.RootContext(nil), &logical.Request{
			ClientToken: initToken,
			Path:        "create",
			Operation:   logical.UpdateOperation,
		})
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
		}
		if !strings.HasPrefix(resp.Auth.ClientToken, consts.ServiceTokenPrefix) {
			t.Fatalf("token %q does not have a 'hvs.' prefix", resp.Auth.ClientToken)
		}
	})

	t.Run("plain custom ID", func(t *testing.T) {
		core, _, root := TestCoreUnsealed(t)
		ts := core.tokenStore

		resp, err := ts.HandleRequest(namespace.RootContext(nil), &logical.Request{
			ClientToken: root,
			Path:        "create",
			Operation:   logical.UpdateOperation,
			Data: map[string]interface{}{
				"id": "foobar",
			},
		})
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
		}

		// Ensure that using a custom token ID results in a warning
		expectedWarning := "Supplying a custom ID for the token uses the weaker SHA1 hashing instead of the more secure SHA2-256 HMAC for token obfuscation. SHA1 hashed tokens on the wire leads to less secure lookups."
		if resp.Warnings[0] != expectedWarning {
			t.Fatalf("expected warning not present")
		}
	})

	t.Run("service token prefix in custom ID", func(t *testing.T) {
		c, _, initToken := TestCoreUnsealed(t)
		ts := c.tokenStore

		// Ensure that custom token ID having a consts.ServiceTokenPrefix prefix fails
		resp, err := ts.HandleRequest(namespace.RootContext(nil), &logical.Request{
			ClientToken: initToken,
			Path:        "create",
			Operation:   logical.UpdateOperation,
			Data: map[string]interface{}{
				"id": "hvs.foobar",
			},
		})
		if err == nil {
			t.Fatalf("expected an error")
		}
		if resp.Error().Error() != "custom token ID cannot have the 'hvs.' prefix" {
			t.Fatalf("expected input error not present in error response")
		}
	})

	t.Run("period in custom ID", func(t *testing.T) {
		core, _, root := TestCoreUnsealed(t)
		ts := core.tokenStore

		resp, err := ts.HandleRequest(namespace.RootContext(nil), &logical.Request{
			ClientToken: root,
			Path:        "create",
			Operation:   logical.UpdateOperation,
			Data: map[string]interface{}{
				"id": "foobar.baz",
			},
		})
		if err == nil {
			t.Fatalf("expected an error")
		}
		if resp.Error().Error() != "custom token ID cannot have a '.' in the value" {
			t.Fatalf("expected input error not present in error response")
		}
	})
}

func expectInGaugeCollection(t *testing.T, expectedLabels map[string]string, expectedValue float32, actual []metricsutil.GaugeLabelValues) {
	t.Helper()

	for _, glv := range actual {
		actualLabels := make(map[string]string)
		for _, l := range glv.Labels {
			actualLabels[l.Name] = l.Value
		}
		if labelsMatch(actualLabels, expectedLabels) {
			if expectedValue != glv.Value {
				t.Errorf("expeced %v for %v, got %v", expectedValue, expectedLabels, glv.Value)
			}
			return
		}
	}
	t.Errorf("didn't find labels %v", expectedLabels)
}

func TestTokenStore_Collectors(t *testing.T) {
	ctx := namespace.RootContext(nil)
	exp := mockExpiration(t)
	ts := exp.tokenStore

	// This helper is defined in expiration.go
	sampleToken(t, exp, "auth/userpass/login", true, "default")
	sampleToken(t, exp, "auth/userpass/login", true, "policy23457")
	sampleToken(t, exp, "auth/token/create", false, "root")
	sampleToken(t, exp, "auth/github/login", true, "root")
	sampleToken(t, exp, "auth/github/login", false, "root")

	waitForRestore(t, exp)

	// By namespace:
	values, err := ts.gaugeCollector(ctx)
	if err != nil {
		t.Fatalf("bad collector run: %v", err)
	}
	if len(values) != 1 {
		t.Errorf("got %v values, expected 1", len(values))
	}
	expectInGaugeCollection(t,
		map[string]string{"namespace": "root"},
		5.0,
		values)

	values, err = ts.gaugeCollectorByPolicy(ctx)
	if err != nil {
		t.Fatalf("bad collector run: %v", err)
	}
	if len(values) != 3 {
		t.Errorf("got %v values, expected 3", len(values))
	}
	expectInGaugeCollection(t,
		map[string]string{"namespace": "root", "policy": "root"},
		3.0,
		values)
	expectInGaugeCollection(t,
		map[string]string{"namespace": "root", "policy": "default"},
		1.0,
		values)
	expectInGaugeCollection(t,
		map[string]string{"namespace": "root", "policy": "policy23457"},
		1.0,
		values)

	values, err = ts.gaugeCollectorByTtl(ctx)
	if err != nil {
		t.Fatalf("bad collector run: %v", err)
	}
	if len(values) != 2 {
		t.Errorf("got %v values, expected 2", len(values))
	}
	expectInGaugeCollection(t,
		map[string]string{"namespace": "root", "creation_ttl": "1h"},
		3.0,
		values)
	expectInGaugeCollection(t,
		map[string]string{"namespace": "root", "creation_ttl": "+Inf"},
		2.0,
		values)

	// Need to set up router for this to work, TODO
	// ts.gaugeCollectorByMethod( ctx )
}
