package vault

import (
	"testing"

	"github.com/hashicorp/vault/logical"
)

func TestBarrierView_impl(t *testing.T) {
	var _ logical.Storage = new(BarrierView)
}

func TestBarrierView_spec(t *testing.T) {
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "foo/")
	logical.TestStorage(t, view)
}

func TestBarrierView(t *testing.T) {
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "foo/")

	// Write a key outside of foo/
	entry := &Entry{Key: "test", Value: []byte("test")}
	if err := barrier.Put(entry); err != nil {
		t.Fatalf("bad: %v", err)
	}

	// List should have no visibility
	keys, err := view.List("")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(keys) != 0 {
		t.Fatalf("bad: %v", err)
	}

	// Get should have no visibility
	out, err := view.Get("test")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out != nil {
		t.Fatalf("bad: %v", out)
	}

	// Try to put the same entry via the view
	if err := view.Put(entry.Logical()); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Check it is nested
	entry, err = barrier.Get("foo/test")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if entry == nil {
		t.Fatalf("missing nested foo/test")
	}

	// Delete nested
	if err := view.Delete("test"); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Check the nested key
	entry, err = barrier.Get("foo/test")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if entry != nil {
		t.Fatalf("nested foo/test should be gone")
	}

	// Check the non-nested key
	entry, err = barrier.Get("test")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if entry == nil {
		t.Fatalf("root test missing")
	}
}

func TestBarrierView_SubView(t *testing.T) {
	_, barrier, _ := mockBarrier(t)
	root := NewBarrierView(barrier, "foo/")
	view := root.SubView("bar/")

	// List should have no visibility
	keys, err := view.List("")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(keys) != 0 {
		t.Fatalf("bad: %v", err)
	}

	// Get should have no visibility
	out, err := view.Get("test")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out != nil {
		t.Fatalf("bad: %v", out)
	}

	// Try to put the same entry via the view
	entry := &logical.StorageEntry{Key: "test", Value: []byte("test")}
	if err := view.Put(entry); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Check it is nested
	bout, err := barrier.Get("foo/bar/test")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if bout == nil {
		t.Fatalf("missing nested foo/bar/test")
	}

	// Check for visibility in root
	out, err = root.Get("bar/test")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out == nil {
		t.Fatalf("missing nested bar/test")
	}

	// Delete nested
	if err := view.Delete("test"); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Check the nested key
	bout, err = barrier.Get("foo/bar/test")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if bout != nil {
		t.Fatalf("nested foo/bar/test should be gone")
	}
}
