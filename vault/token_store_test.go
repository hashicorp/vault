package vault

import (
	"reflect"
	"testing"
)

func mockTokenStore(t *testing.T) *TokenStore {
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "foo/")
	ts, err := NewTokenStore(view)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	return ts
}

func TestTokenStore_CreateLookup(t *testing.T) {
	ts := mockTokenStore(t)

	ent := &TokenEntry{Source: "test", Policies: []string{"dev", "ops"}}
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
	ts2, err := NewTokenStore(ts.view)
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
