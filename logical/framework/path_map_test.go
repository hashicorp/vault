package framework

import (
	"testing"

	"github.com/hashicorp/vault/logical"
)

func TestPathMap(t *testing.T) {
	p := &PathMap{Name: "foo"}
	storage := new(logical.InmemStorage)
	var b logical.Backend = &Backend{Paths: p.Paths()}

	// Write via HTTP
	_, err := b.HandleRequest(&logical.Request{
		Operation: logical.WriteOperation,
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
	if v != "bar" {
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
	if v != "" {
		t.Fatalf("bad: %#v", v)
	}
}
