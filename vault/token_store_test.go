package vault

import (
	"reflect"
	"testing"
)

func mockTokenStore(t *testing.T) *TokenStore {
	c, _ := TestCoreUnsealed(t)
	ts, err := NewTokenStore(c)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	return ts
}

func TestTokenStore_CreateLookup(t *testing.T) {
	ts := mockTokenStore(t)

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
	ts2, err := NewTokenStore(ts.core)
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

func TestTokenStore_Revoke(t *testing.T) {
	ts := mockTokenStore(t)

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

func TestTokenStore_Revoke_Orphan(t *testing.T) {
	ts := mockTokenStore(t)

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
	ts := mockTokenStore(t)

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

func TestTokenStore_RevokeAll(t *testing.T) {
	ts := mockTokenStore(t)

	ent1 := &TokenEntry{}
	if err := ts.Create(ent1); err != nil {
		t.Fatalf("err: %v", err)
	}

	ent2 := &TokenEntry{Parent: ent1.ID}
	if err := ts.Create(ent2); err != nil {
		t.Fatalf("err: %v", err)
	}

	ent3 := &TokenEntry{}
	if err := ts.Create(ent3); err != nil {
		t.Fatalf("err: %v", err)
	}

	ent4 := &TokenEntry{Parent: ent3.ID}
	if err := ts.Create(ent4); err != nil {
		t.Fatalf("err: %v", err)
	}

	err := ts.RevokeAll()
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
