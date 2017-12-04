package framework

import (
	"testing"

	"github.com/hashicorp/vault/helper/salt"
	"github.com/hashicorp/vault/logical"
)

func TestPathMap(t *testing.T) {
	p := &PathMap{Name: "foo"}
	storage := new(logical.InmemStorage)
	var b logical.Backend = &Backend{Paths: p.Paths()}

	// Write via HTTP
	_, err := b.HandleRequest(&logical.Request{
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
	resp, err := b.HandleRequest(&logical.Request{
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
	v, err := p.Get(storage, "a")
	if err != nil {
		t.Fatalf("bad: %#v", err)
	}
	if v["value"] != "bar" {
		t.Fatalf("bad: %#v", v)
	}

	// Read via API with other casing
	v, err = p.Get(storage, "A")
	if err != nil {
		t.Fatalf("bad: %#v", err)
	}
	if v["value"] != "bar" {
		t.Fatalf("bad: %#v", v)
	}

	// Verify List
	keys, err := p.List(storage, "")
	if err != nil {
		t.Fatalf("bad: %#v", err)
	}
	if len(keys) != 1 || keys[0] != "a" {
		t.Fatalf("bad: %#v", keys)
	}

	// LIST via HTTP
	resp, err = b.HandleRequest(&logical.Request{
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
	resp, err = b.HandleRequest(&logical.Request{
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
	resp, err = b.HandleRequest(&logical.Request{
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
	v, err = p.Get(storage, "a")
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

	v, err := p.Get(storage, "nope")
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
	salt, err := salt.NewSalt(storage, &salt.Config{
		HashFunc: salt.SHA1Hash,
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	p := &PathMap{Name: "foo", Salt: salt}
	var b logical.Backend = &Backend{Paths: p.Paths()}

	// Write via HTTP
	_, err = b.HandleRequest(&logical.Request{
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
	out, err := storage.Get("struct/map/foo/a")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out != nil {
		t.Fatalf("non-salted key found")
	}

	// Ensure the path is salted
	expect := salt.SaltID("a")
	out, err = storage.Get("struct/map/foo/" + expect)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out == nil {
		t.Fatalf("missing salted key")
	}

	// Read via HTTP
	resp, err := b.HandleRequest(&logical.Request{
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
	v, err := p.Get(storage, "a")
	if err != nil {
		t.Fatalf("bad: %#v", err)
	}
	if v["value"] != "bar" {
		t.Fatalf("bad: %#v", v)
	}

	// Read via API with other casing
	v, err = p.Get(storage, "A")
	if err != nil {
		t.Fatalf("bad: %#v", err)
	}
	if v["value"] != "bar" {
		t.Fatalf("bad: %#v", v)
	}

	// Verify List
	keys, err := p.List(storage, "")
	if err != nil {
		t.Fatalf("bad: %#v", err)
	}
	if len(keys) != 1 || keys[0] != expect {
		t.Fatalf("bad: %#v", keys)
	}

	// Delete via HTTP
	resp, err = b.HandleRequest(&logical.Request{
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
	resp, err = b.HandleRequest(&logical.Request{
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
	v, err = p.Get(storage, "a")
	if err != nil {
		t.Fatalf("bad: %#v", err)
	}
	if v != nil {
		t.Fatalf("bad: %#v", v)
	}

	// Put in a non-salted version and make sure that after reading it's been
	// upgraded
	err = storage.Put(&logical.StorageEntry{
		Key:   "struct/map/foo/b",
		Value: []byte(`{"foo": "bar"}`),
	})
	if err != nil {
		t.Fatal("err: %v", err)
	}
	// A read should transparently upgrade
	resp, err = b.HandleRequest(&logical.Request{
		Operation: logical.ReadOperation,
		Path:      "map/foo/b",
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	list, _ := storage.List("struct/map/foo/")
	if len(list) != 1 {
		t.Fatalf("unexpected number of entries left after upgrade; expected 1, got %d", len(list))
	}
	found := false
	for _, v := range list {
		if v == salt.SaltID("b") {
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
	locSalt, err := salt.NewSalt(storage, &salt.Config{
		HashFunc: salt.SHA1Hash,
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	saltFunc := func() (*salt.Salt, error) {
		return locSalt, nil
	}
	p := &PathMap{Name: "foo", SaltFunc: saltFunc}
	var b logical.Backend = &Backend{Paths: p.Paths()}

	// Write via HTTP
	_, err = b.HandleRequest(&logical.Request{
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
	out, err := storage.Get("struct/map/foo/a")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out != nil {
		t.Fatalf("non-salted key found")
	}

	// Ensure the path is salted
	expect := locSalt.SaltID("a")
	out, err = storage.Get("struct/map/foo/" + expect)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out == nil {
		t.Fatalf("missing salted key")
	}

	// Read via HTTP
	resp, err := b.HandleRequest(&logical.Request{
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
	v, err := p.Get(storage, "a")
	if err != nil {
		t.Fatalf("bad: %#v", err)
	}
	if v["value"] != "bar" {
		t.Fatalf("bad: %#v", v)
	}

	// Read via API with other casing
	v, err = p.Get(storage, "A")
	if err != nil {
		t.Fatalf("bad: %#v", err)
	}
	if v["value"] != "bar" {
		t.Fatalf("bad: %#v", v)
	}

	// Verify List
	keys, err := p.List(storage, "")
	if err != nil {
		t.Fatalf("bad: %#v", err)
	}
	if len(keys) != 1 || keys[0] != expect {
		t.Fatalf("bad: %#v", keys)
	}

	// Delete via HTTP
	resp, err = b.HandleRequest(&logical.Request{
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
	resp, err = b.HandleRequest(&logical.Request{
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
	v, err = p.Get(storage, "a")
	if err != nil {
		t.Fatalf("bad: %#v", err)
	}
	if v != nil {
		t.Fatalf("bad: %#v", v)
	}

	// Put in a non-salted version and make sure that after reading it's been
	// upgraded
	err = storage.Put(&logical.StorageEntry{
		Key:   "struct/map/foo/b",
		Value: []byte(`{"foo": "bar"}`),
	})
	if err != nil {
		t.Fatal("err: %v", err)
	}
	// A read should transparently upgrade
	resp, err = b.HandleRequest(&logical.Request{
		Operation: logical.ReadOperation,
		Path:      "map/foo/b",
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	list, _ := storage.List("struct/map/foo/")
	if len(list) != 1 {
		t.Fatalf("unexpected number of entries left after upgrade; expected 1, got %d", len(list))
	}
	found := false
	for _, v := range list {
		if v == locSalt.SaltID("b") {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("did not find upgraded value")
	}
}
