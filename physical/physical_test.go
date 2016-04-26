package physical

import (
	"log"
	"os"
	"reflect"
	"sort"
	"testing"
	"time"
)

func testNewBackend(t *testing.T) {
	logger := log.New(os.Stderr, "", log.LstdFlags)
	_, err := NewBackend("foobar", logger, nil)
	if err == nil {
		t.Fatalf("expected error")
	}

	b, err := NewBackend("inmem", logger, nil)
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

	// Multiple Puts should work; GH-189
	e = &Entry{Key: "foo", Value: []byte("test")}
	err = b.Put(e)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	e = &Entry{Key: "foo", Value: []byte("test")}
	err = b.Put(e)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Make a nested entry
	e = &Entry{Key: "foo/bar", Value: []byte("baz")}
	err = b.Put(e)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Delete with children should work
	err = b.Delete("foo")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Get should return the child
	out, err = b.Get("foo/bar")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out == nil {
		t.Fatalf("missing child")
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
	sort.Strings(keys)
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
	sort.Strings(keys)
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
	sort.Strings(keys)
	if len(keys) != 1 {
		t.Fatalf("bad: %v", keys)
	}
	if keys[0] != "baz" {
		t.Fatalf("bad: %v", keys)
	}

}

func testHABackend(t *testing.T, b HABackend, b2 HABackend) {
	// Get the lock
	lock, err := b.LockWith("foo", "bar")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Attempt to lock
	leaderCh, err := lock.Lock(nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if leaderCh == nil {
		t.Fatalf("failed to get leader ch")
	}

	// Check the value
	held, val, err := lock.Value()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !held {
		t.Fatalf("should be held")
	}
	if val != "bar" {
		t.Fatalf("bad value: %v", err)
	}

	// Second acquisition should fail
	lock2, err := b2.LockWith("foo", "baz")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Cancel attempt in 50 msec
	stopCh := make(chan struct{})
	time.AfterFunc(50*time.Millisecond, func() {
		close(stopCh)
	})

	// Attempt to lock
	leaderCh2, err := lock2.Lock(stopCh)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if leaderCh2 != nil {
		t.Fatalf("should not get leader ch")
	}

	// Release the first lock
	lock.Unlock()

	// Attempt to lock should work
	leaderCh2, err = lock2.Lock(nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if leaderCh2 == nil {
		t.Fatalf("should get leader ch")
	}

	// Check the value
	held, val, err = lock.Value()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !held {
		t.Fatalf("should be held")
	}
	if val != "baz" {
		t.Fatalf("bad value: %v", err)
	}
	// Cleanup
	lock2.Unlock()
}
