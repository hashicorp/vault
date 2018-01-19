package framework

import (
	"context"
	"reflect"
	"testing"

	"github.com/hashicorp/vault/logical"
)

func TestWAL(t *testing.T) {
	s := new(logical.InmemStorage)

	ctx := context.Background()

	// WAL should be empty to start
	keys, err := ListWAL(ctx, s)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if len(keys) > 0 {
		t.Fatalf("bad: %#v", keys)
	}

	// Write an entry to the WAL
	id, err := PutWAL(ctx, s, "foo", "bar")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// The key should be in the WAL
	keys, err = ListWAL(ctx, s)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if !reflect.DeepEqual(keys, []string{id}) {
		t.Fatalf("bad: %#v", keys)
	}

	// Should be able to get the value
	entry, err := GetWAL(ctx, s, id)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if entry.Kind != "foo" {
		t.Fatalf("bad: %#v", entry)
	}
	if entry.Data != "bar" {
		t.Fatalf("bad: %#v", entry)
	}

	// Should be able to delete the value
	if err := DeleteWAL(ctx, s, id); err != nil {
		t.Fatalf("err: %s", err)
	}
	entry, err = GetWAL(ctx, s, id)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if entry != nil {
		t.Fatalf("bad: %#v", entry)
	}
}
