package framework

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/logical"
)

func TestPathStruct(t *testing.T) {
	p := &PathStruct{
		Name: "foo",
		Path: "bar",
		Schema: map[string]*FieldSchema{
			"value": &FieldSchema{Type: TypeString},
		},
		Read: true,
	}

	storage := new(logical.InmemStorage)
	var b logical.Backend = &Backend{Paths: p.Paths()}

	ctx := context.Background()

	// Write via HTTP
	_, err := b.HandleRequest(ctx, &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "bar",
		Data: map[string]interface{}{
			"value": "baz",
		},
		Storage: storage,
	})
	if err != nil {
		t.Fatalf("bad: %#v", err)
	}

	// Read via HTTP
	resp, err := b.HandleRequest(ctx, &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "bar",
		Storage:   storage,
	})
	if err != nil {
		t.Fatalf("bad: %#v", err)
	}
	if resp.Data["value"] != "baz" {
		t.Fatalf("bad: %#v", resp)
	}

	// Read via API
	v, err := p.Get(ctx, storage)
	if err != nil {
		t.Fatalf("bad: %#v", err)
	}
	if v["value"] != "baz" {
		t.Fatalf("bad: %#v", v)
	}

	// Delete via HTTP
	resp, err = b.HandleRequest(ctx, &logical.Request{
		Operation: logical.DeleteOperation,
		Path:      "bar",
		Data:      nil,
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
		Path:      "bar",
		Storage:   storage,
	})
	if err != nil {
		t.Fatalf("bad: %#v", err)
	}
	if _, ok := resp.Data["value"]; ok {
		t.Fatalf("bad: %#v", resp)
	}

	// Re-read via API
	v, err = p.Get(ctx, storage)
	if err != nil {
		t.Fatalf("bad: %#v", err)
	}
	if v != nil {
		t.Fatalf("bad: %#v", v)
	}
}
