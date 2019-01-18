package vault

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/vault/helper/jsonutil"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/logical"
)

func TestAuth_ReadOnlyViewDuringMount(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	c.credentialBackends["noop"] = func(ctx context.Context, config *logical.BackendConfig) (logical.Backend, error) {
		err := config.StorageView.Put(ctx, &logical.StorageEntry{
			Key:   "bar",
			Value: []byte("baz"),
		})
		if err == nil || !strings.Contains(err.Error(), logical.ErrSetupReadOnly.Error()) {
			t.Fatalf("expected a read-only error")
		}
		return &NoopBackend{
			BackendType: logical.TypeCredential,
		}, nil
	}

	me := &MountEntry{
		Table: credentialTableType,
		Path:  "foo",
		Type:  "noop",
	}
	err := c.enableCredential(namespace.RootContext(nil), me)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
}

func TestCore_DefaultAuthTable(t *testing.T) {
	c, keys, _ := TestCoreUnsealed(t)
	verifyDefaultAuthTable(t, c.auth)

	// Start a second core with same physical
	conf := &CoreConfig{
		Physical:     c.physical,
		DisableMlock: true,
	}
	c2, err := NewCore(conf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	for i, key := range keys {
		unseal, err := TestCoreUnseal(c2, key)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if i+1 == len(keys) && !unseal {
			t.Fatalf("should be unsealed")
		}
	}

	// Verify matching mount tables
	if !reflect.DeepEqual(c.auth, c2.auth) {
		t.Fatalf("mismatch: %v %v", c.auth, c2.auth)
	}
}

func TestCore_EnableCredential(t *testing.T) {
	c, keys, _ := TestCoreUnsealed(t)
	c.credentialBackends["noop"] = func(context.Context, *logical.BackendConfig) (logical.Backend, error) {
		return &NoopBackend{
			BackendType: logical.TypeCredential,
		}, nil
	}

	me := &MountEntry{
		Table: credentialTableType,
		Path:  "foo",
		Type:  "noop",
	}
	err := c.enableCredential(namespace.RootContext(nil), me)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	match := c.router.MatchingMount(namespace.RootContext(nil), "auth/foo/bar")
	if match != "auth/foo/" {
		t.Fatalf("missing mount, match: %q", match)
	}

	conf := &CoreConfig{
		Physical:     c.physical,
		DisableMlock: true,
	}
	c2, err := NewCore(conf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	c2.credentialBackends["noop"] = func(context.Context, *logical.BackendConfig) (logical.Backend, error) {
		return &NoopBackend{
			BackendType: logical.TypeCredential,
		}, nil
	}
	for i, key := range keys {
		unseal, err := TestCoreUnseal(c2, key)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if i+1 == len(keys) && !unseal {
			t.Fatalf("should be unsealed")
		}
	}

	// Verify matching auth tables
	if !reflect.DeepEqual(c.auth, c2.auth) {
		t.Fatalf("mismatch: %v %v", c.auth, c2.auth)
	}
}

// Test that the local table actually gets populated as expected with local
// entries, and that upon reading the entries from both are recombined
// correctly
func TestCore_EnableCredential_Local(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	c.credentialBackends["noop"] = func(context.Context, *logical.BackendConfig) (logical.Backend, error) {
		return &NoopBackend{
			BackendType: logical.TypeCredential,
		}, nil
	}

	c.auth = &MountTable{
		Type: credentialTableType,
		Entries: []*MountEntry{
			&MountEntry{
				Table:            credentialTableType,
				Path:             "noop/",
				Type:             "noop",
				UUID:             "abcd",
				Accessor:         "noop-abcd",
				BackendAwareUUID: "abcde",
				NamespaceID:      namespace.RootNamespaceID,
				namespace:        namespace.RootNamespace,
			},
			&MountEntry{
				Table:            credentialTableType,
				Path:             "noop2/",
				Type:             "noop",
				UUID:             "bcde",
				Accessor:         "noop-bcde",
				BackendAwareUUID: "bcdea",
				NamespaceID:      namespace.RootNamespaceID,
				namespace:        namespace.RootNamespace,
			},
		},
	}

	// Both should set up successfully
	err := c.setupCredentials(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	rawLocal, err := c.barrier.Get(context.Background(), coreLocalAuthConfigPath)
	if err != nil {
		t.Fatal(err)
	}
	if rawLocal == nil {
		t.Fatal("expected non-nil local credential")
	}
	localCredentialTable := &MountTable{}
	if err := jsonutil.DecodeJSON(rawLocal.Value, localCredentialTable); err != nil {
		t.Fatal(err)
	}
	if len(localCredentialTable.Entries) > 0 {
		t.Fatalf("expected no entries in local credential table, got %#v", localCredentialTable)
	}

	c.auth.Entries[1].Local = true
	if err := c.persistAuth(context.Background(), c.auth, nil); err != nil {
		t.Fatal(err)
	}

	rawLocal, err = c.barrier.Get(context.Background(), coreLocalAuthConfigPath)
	if err != nil {
		t.Fatal(err)
	}
	if rawLocal == nil {
		t.Fatal("expected non-nil local credential")
	}
	localCredentialTable = &MountTable{}
	if err := jsonutil.DecodeJSON(rawLocal.Value, localCredentialTable); err != nil {
		t.Fatal(err)
	}
	if len(localCredentialTable.Entries) != 1 {
		t.Fatalf("expected one entry in local credential table, got %#v", localCredentialTable)
	}

	oldCredential := c.auth
	if err := c.loadCredentials(context.Background()); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(oldCredential, c.auth) {
		t.Fatalf("expected\n%#v\ngot\n%#v\n", oldCredential, c.auth)
	}

	if len(c.auth.Entries) != 2 {
		t.Fatalf("expected two credential entries, got %#v", localCredentialTable)
	}
}

func TestCore_EnableCredential_twice_409(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	c.credentialBackends["noop"] = func(context.Context, *logical.BackendConfig) (logical.Backend, error) {
		return &NoopBackend{
			BackendType: logical.TypeCredential,
		}, nil
	}

	me := &MountEntry{
		Table: credentialTableType,
		Path:  "foo",
		Type:  "noop",
	}
	err := c.enableCredential(namespace.RootContext(nil), me)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// 2nd should be a 409 error
	err2 := c.enableCredential(namespace.RootContext(nil), me)
	switch err2.(type) {
	case logical.HTTPCodedError:
		if err2.(logical.HTTPCodedError).Code() != 409 {
			t.Fatalf("invalid code given")
		}
	default:
		t.Fatalf("expected a different error type")
	}
}

func TestCore_EnableCredential_Token(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	me := &MountEntry{
		Table: credentialTableType,
		Path:  "foo",
		Type:  "token",
	}
	err := c.enableCredential(namespace.RootContext(nil), me)
	if err.Error() != "token credential backend cannot be instantiated" {
		t.Fatalf("err: %v", err)
	}
}

func TestCore_DisableCredential(t *testing.T) {
	c, keys, _ := TestCoreUnsealed(t)
	c.credentialBackends["noop"] = func(context.Context, *logical.BackendConfig) (logical.Backend, error) {
		return &NoopBackend{
			BackendType: logical.TypeCredential,
		}, nil
	}

	err := c.disableCredential(namespace.RootContext(nil), "foo")
	if err != nil && !strings.HasPrefix(err.Error(), "no matching mount") {
		t.Fatal(err)
	}

	me := &MountEntry{
		Table: credentialTableType,
		Path:  "foo",
		Type:  "noop",
	}
	err = c.enableCredential(namespace.RootContext(nil), me)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	err = c.disableCredential(namespace.RootContext(nil), "foo")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	match := c.router.MatchingMount(namespace.RootContext(nil), "auth/foo/bar")
	if match != "" {
		t.Fatalf("backend present")
	}

	conf := &CoreConfig{
		Physical:     c.physical,
		DisableMlock: true,
	}
	c2, err := NewCore(conf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	for i, key := range keys {
		unseal, err := TestCoreUnseal(c2, key)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if i+1 == len(keys) && !unseal {
			t.Fatalf("should be unsealed")
		}
	}

	// Verify matching mount tables
	if !reflect.DeepEqual(c.auth, c2.auth) {
		t.Fatalf("mismatch: %v %v", c.auth, c2.auth)
	}
}

func TestCore_DisableCredential_Protected(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	err := c.disableCredential(namespace.RootContext(nil), "token")
	if err.Error() != "token credential backend cannot be disabled" {
		t.Fatalf("err: %v", err)
	}
}

func TestCore_DisableCredential_Cleanup(t *testing.T) {
	noop := &NoopBackend{
		Login:       []string{"login"},
		BackendType: logical.TypeCredential,
	}
	c, _, _ := TestCoreUnsealed(t)
	c.credentialBackends["noop"] = func(context.Context, *logical.BackendConfig) (logical.Backend, error) {
		return noop, nil
	}

	me := &MountEntry{
		Table: credentialTableType,
		Path:  "foo",
		Type:  "noop",
	}
	err := c.enableCredential(namespace.RootContext(nil), me)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Store the view
	view := c.router.MatchingStorageByAPIPath(namespace.RootContext(nil), "auth/foo/")

	// Inject data
	se := &logical.StorageEntry{
		Key:   "plstodelete",
		Value: []byte("test"),
	}
	if err := view.Put(context.Background(), se); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Generate a new token auth
	noop.Response = &logical.Response{
		Auth: &logical.Auth{
			Policies: []string{"foo"},
		},
	}
	r := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "auth/foo/login",
	}
	resp, err := c.HandleRequest(namespace.RootContext(nil), r)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp.Auth.ClientToken == "" {
		t.Fatalf("bad: %#v", resp)
	}

	// Disable should cleanup
	err = c.disableCredential(namespace.RootContext(nil), "foo")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Token should be revoked
	te, err := c.tokenStore.Lookup(namespace.RootContext(nil), resp.Auth.ClientToken)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if te != nil {
		t.Fatalf("bad: %#v", te)
	}

	// View should be empty
	out, err := logical.CollectKeys(context.Background(), view)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(out) != 0 {
		t.Fatalf("bad: %#v", out)
	}
}

func TestDefaultAuthTable(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	table := c.defaultAuthTable()
	verifyDefaultAuthTable(t, table)
}

func verifyDefaultAuthTable(t *testing.T, table *MountTable) {
	if len(table.Entries) != 1 {
		t.Fatalf("bad: %v", table.Entries)
	}
	if table.Type != credentialTableType {
		t.Fatalf("bad: %v", *table)
	}
	for idx, entry := range table.Entries {
		switch idx {
		case 0:
			if entry.Path != "token/" {
				t.Fatalf("bad: %v", entry)
			}
			if entry.Type != "token" {
				t.Fatalf("bad: %v", entry)
			}
		}
		if entry.Description == "" {
			t.Fatalf("bad: %v", entry)
		}
		if entry.UUID == "" {
			t.Fatalf("bad: %v", entry)
		}
	}
}
