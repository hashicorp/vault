package physical

import (
	"reflect"
	"testing"
)

func testNewBackend(t *testing.T) {
	_, err := NewBackend("foobar", nil)
	if err == nil {
		t.Fatalf("expected error")
	}

	b, err := NewBackend("inmem", nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if b == nil {
		t.Fatalf("expected backend")
	}
}

func testBackend(t *testing.T, b Backend) {
	// Should be empty
	keys, err := b.List("")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(keys) != 0 {
		t.Fatalf("bad: %v", keys)
	}

	// Delete should work if it does not exist
	err = b.Delete("foo")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Get should fail
	out, err := b.Get("foo")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out != nil {
		t.Fatalf("bad: %v", out)
	}

	// Make an entry
	e := &Entry{Key: "foo", Value: []byte("test")}
	err = b.Put(e)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Get should work
	out, err = b.Get("foo")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !reflect.DeepEqual(out, e) {
		t.Fatalf("bad: %v expected: %v", out, e)
	}

	// List should not be empty
	keys, err = b.List("")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(keys) != 1 {
		t.Fatalf("bad: %v", keys)
	}
	if keys[0] != "foo" {
		t.Fatalf("bad: %v", keys)
	}

	// Delete should work
	err = b.Delete("foo")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Should be empty
	keys, err = b.List("")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(keys) != 0 {
		t.Fatalf("bad: %v", keys)
	}

	// Get should fail
	out, err = b.Get("foo")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out != nil {
		t.Fatalf("bad: %v", out)
	}
}

func testBackend_ListPrefix(t *testing.T, b Backend) {
	e1 := &Entry{Key: "foo", Value: []byte("test")}
	e2 := &Entry{Key: "foo/bar", Value: []byte("test")}
	e3 := &Entry{Key: "foo/bar/baz", Value: []byte("test")}

	err := b.Put(e1)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	err = b.Put(e2)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	err = b.Put(e3)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Scan the root
	keys, err := b.List("")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(keys) != 2 {
		t.Fatalf("bad: %v", keys)
	}
	if keys[0] != "foo" {
		t.Fatalf("bad: %v", keys)
	}
	if keys[1] != "foo/" {
		t.Fatalf("bad: %v", keys)
	}

	// Scan foo/
	keys, err = b.List("foo/")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(keys) != 2 {
		t.Fatalf("bad: %v", keys)
	}
	if keys[0] != "bar" {
		t.Fatalf("bad: %v", keys)
	}
	if keys[1] != "bar/" {
		t.Fatalf("bad: %v", keys)
	}

	// Scan foo/bar/
	keys, err = b.List("foo/bar/")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(keys) != 1 {
		t.Fatalf("bad: %v", keys)
	}
	if keys[0] != "baz" {
		t.Fatalf("bad: %v", keys)
	}
}
