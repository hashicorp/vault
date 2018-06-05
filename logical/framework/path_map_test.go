package framework

import (
	"context"
	"testing"

	saltpkg "github.com/hashicorp/vault/helper/salt"
	"github.com/hashicorp/vault/logical"
)

func TestPathMap(t *testing.T) {
	p := &PathMap{Name: "foo"}
	storage := new(logical.InmemStorage)
	var b logical.Backend = &Backend{Paths: p.Paths()}

	ctx := context.Background()

	// Write via HTTP
	_, err := b.HandleRequest(ctx, &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "map/foo/a",
		Data: map[string]interface{}{
			"value": "bar",
		},
		Storage: storage,
	})
	if err != nil {
		t.Fatalf("bad: %#v", err)
	}

	// Read via HTTP
	resp, err := b.HandleRequest(ctx, &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "map/foo/a",
		Storage:   storage,
	})
	if err != nil {
		t.Fatalf("bad: %#v", err)
	}
	if resp.Data["value"] != "bar" {
		t.Fatalf("bad: %#v", resp)
	}

	// Read via API
	v, err := p.Get(ctx, storage, "a")
	if err != nil {
		t.Fatalf("bad: %#v", err)
	}
	if v["value"] != "bar" {
		t.Fatalf("bad: %#v", v)
	}

	// Read via API with other casing
	v, err = p.Get(ctx, storage, "A")
	if err != nil {
		t.Fatalf("bad: %#v", err)
	}
	if v["value"] != "bar" {
		t.Fatalf("bad: %#v", v)
	}

	// Verify List
	keys, err := p.List(ctx, storage, "")
	if err != nil {
		t.Fatalf("bad: %#v", err)
	}
	if len(keys) != 1 || keys[0] != "a" {
		t.Fatalf("bad: %#v", keys)
	}

	// LIST via HTTP
	resp, err = b.HandleRequest(ctx, &logical.Request{
		Operation: logical.ListOperation,
		Path:      "map/foo/",
		Storage:   storage,
	})
	if err != nil {
		t.Fatalf("bad: %#v", err)
	}
	if len(resp.Data) != 1 || len(resp.Data["keys"].([]string)) != 1 ||
		resp.Data["keys"].([]string)[0] != "a" {
		t.Fatalf("bad: %#v", resp)
	}

	// Delete via HTTP
	resp, err = b.HandleRequest(ctx, &logical.Request{
		Operation: logical.DeleteOperation,
		Path:      "map/foo/a",
		Storage:   storage,
	})
	if err != nil {
		t.Fatalf("bad: %#v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %#v", resp)
	}

	// Re-read via HTTP
	resp, err = b.HandleRequest(ctx, &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "map/foo/a",
		Storage:   storage,
	})
	if err != nil {
		t.Fatalf("bad: %#v", err)
	}
	if _, ok := resp.Data["value"]; ok {
		t.Fatalf("bad: %#v", resp)
	}

	// Re-read via API
	v, err = p.Get(ctx, storage, "a")
	if err != nil {
		t.Fatalf("bad: %#v", err)
	}
	if v != nil {
		t.Fatalf("bad: %#v", v)
	}
}

func TestPathMap_getInvalid(t *testing.T) {
	p := &PathMap{Name: "foo"}
	storage := new(logical.InmemStorage)

	v, err := p.Get(context.Background(), storage, "nope")
	if err != nil {
		t.Fatalf("bad: %#v", err)
	}
	if v != nil {
		t.Fatalf("bad: %#v", v)
	}
}

func TestPathMap_routes(t *testing.T) {
	p := &PathMap{Name: "foo"}
	TestBackendRoutes(t, &Backend{Paths: p.Paths()}, []string{
		"map/foo",         // Normal
		"map/foo/bar",     // Normal
		"map/foo/bar-baz", // Hyphen key
	})
}

func TestPathMap_Salted(t *testing.T) {
	storage := new(logical.InmemStorage)

	salt, err := saltpkg.NewSalt(context.Background(), storage, &saltpkg.Config{
		HashFunc: saltpkg.SHA1Hash,
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	testSalting(t, context.Background(), storage, salt, &PathMap{Name: "foo", Salt: salt})
}

func testSalting(t *testing.T, ctx context.Context, storage logical.Storage, salt *saltpkg.Salt, p *PathMap) {
	var b logical.Backend = &Backend{Paths: p.Paths()}
	var err error

	// Write via HTTP
	_, err = b.HandleRequest(ctx, &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "map/foo/a",
		Data: map[string]interface{}{
			"value": "bar",
		},
		Storage: storage,
	})
	if err != nil {
		t.Fatalf("bad: %#v", err)
	}

	// Non-salted version should not be there
	out, err := storage.Get(ctx, "struct/map/foo/a")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out != nil {
		t.Fatalf("non-salted key found")
	}

	// Ensure the path is salted
	expect := "s" + salt.SaltIDHashFunc("a", saltpkg.SHA256Hash)
	out, err = storage.Get(ctx, "struct/map/foo/"+expect)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out == nil {
		t.Fatalf("missing salted key")
	}

	// Read via HTTP
	resp, err := b.HandleRequest(ctx, &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "map/foo/a",
		Storage:   storage,
	})
	if err != nil {
		t.Fatalf("bad: %#v", err)
	}
	if resp.Data["value"] != "bar" {
		t.Fatalf("bad: %#v", resp)
	}

	// Read via API
	v, err := p.Get(ctx, storage, "a")
	if err != nil {
		t.Fatalf("bad: %#v", err)
	}
	if v["value"] != "bar" {
		t.Fatalf("bad: %#v", v)
	}

	// Read via API with other casing
	v, err = p.Get(ctx, storage, "A")
	if err != nil {
		t.Fatalf("bad: %#v", err)
	}
	if v["value"] != "bar" {
		t.Fatalf("bad: %#v", v)
	}

	// Verify List
	keys, err := p.List(ctx, storage, "")
	if err != nil {
		t.Fatalf("bad: %#v", err)
	}
	if len(keys) != 1 || keys[0] != expect {
		t.Fatalf("bad: %#v", keys)
	}

	// Delete via HTTP
	resp, err = b.HandleRequest(ctx, &logical.Request{
		Operation: logical.DeleteOperation,
		Path:      "map/foo/a",
		Storage:   storage,
	})
	if err != nil {
		t.Fatalf("bad: %#v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %#v", resp)
	}

	// Re-read via HTTP
	resp, err = b.HandleRequest(ctx, &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "map/foo/a",
		Storage:   storage,
	})
	if err != nil {
		t.Fatalf("bad: %#v", err)
	}
	if _, ok := resp.Data["value"]; ok {
		t.Fatalf("bad: %#v", resp)
	}

	// Re-read via API
	v, err = p.Get(ctx, storage, "a")
	if err != nil {
		t.Fatalf("bad: %#v", err)
	}
	if v != nil {
		t.Fatalf("bad: %#v", v)
	}

	// Put in a non-salted version and make sure that after reading it's been
	// upgraded
	err = storage.Put(ctx, &logical.StorageEntry{
		Key:   "struct/map/foo/b",
		Value: []byte(`{"foo": "bar"}`),
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	// A read should transparently upgrade
	resp, err = b.HandleRequest(ctx, &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "map/foo/b",
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	list, _ := storage.List(ctx, "struct/map/foo/")
	if len(list) != 1 {
		t.Fatalf("unexpected number of entries left after upgrade; expected 1, got %d", len(list))
	}
	found := false
	for _, v := range list {
		if v == "s"+salt.SaltIDHashFunc("b", saltpkg.SHA256Hash) {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("did not find upgraded value")
	}

	// Put in a SHA1 salted version and make sure that after reading its been
	// upgraded
	err = storage.Put(ctx, &logical.StorageEntry{
		Key:   "struct/map/foo/" + salt.SaltID("b"),
		Value: []byte(`{"foo": "bar"}`),
	})
	if err != nil {
		t.Fatal(err)
	}

	// A read should transparently upgrade
	resp, err = b.HandleRequest(ctx, &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "map/foo/b",
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	list, _ = storage.List(ctx, "struct/map/foo/")
	if len(list) != 1 {
		t.Fatalf("unexpected number of entries left after upgrade; expected 1, got %d", len(list))
	}
	found = false
	for _, v := range list {
		if v == "s"+salt.SaltIDHashFunc("b", saltpkg.SHA256Hash) {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("did not find upgraded value")
	}
}

func TestPathMap_SaltFunc(t *testing.T) {
	storage := new(logical.InmemStorage)

	salt, err := saltpkg.NewSalt(context.Background(), storage, &saltpkg.Config{
		HashFunc: saltpkg.SHA1Hash,
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	saltFunc := func(context.Context) (*saltpkg.Salt, error) {
		return salt, nil
	}

	testSalting(t, context.Background(), storage, salt, &PathMap{Name: "foo", SaltFunc: saltFunc})
}
