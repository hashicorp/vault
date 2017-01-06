package physical

import (
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/logformat"
	log "github.com/mgutz/logxi/v1"
)

func testNewBackend(t *testing.T) {
	logger := logformat.NewVaultLogger(log.LevelTrace)

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

	keys, err = b.List("")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(keys) != 2 {
		t.Fatalf("bad: %v", keys)
	}
	sort.Strings(keys)
	if keys[0] != "foo" || keys[1] != "foo/" {
		t.Fatalf("bad: %v", keys)
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

	// Removal of nested secret should not leave artifacts
	e = &Entry{Key: "foo/nested1/nested2/nested3", Value: []byte("baz")}
	err = b.Put(e)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	err = b.Delete("foo/nested1/nested2/nested3")
	if err != nil {
		t.Fatalf("failed to remove nested secret: %v", err)
	}

	keys, err = b.List("foo/")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if len(keys) != 1 {
		t.Fatalf("there should be only one key left after deleting nested "+
			"secret: %v", keys)
	}

	if keys[0] != "bar" {
		t.Fatalf("bad keys after deleting nested: %v", keys)
	}

	// Make a second nested entry to test prefix removal
	e = &Entry{Key: "foo/zip", Value: []byte("zap")}
	err = b.Put(e)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Delete should not remove the prefix
	err = b.Delete("foo/bar")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	keys, err = b.List("")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(keys) != 1 {
		t.Fatalf("bad: %v", keys)
	}
	if keys[0] != "foo/" {
		t.Fatalf("bad: %v", keys)
	}

	// Delete should remove the prefix
	err = b.Delete("foo/zip")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	keys, err = b.List("")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(keys) != 0 {
		t.Fatalf("bad: %v", keys)
	}
}

func testBackend_ListPrefix(t *testing.T, b Backend) {
	e1 := &Entry{Key: "foo", Value: []byte("test")}
	e2 := &Entry{Key: "foo/bar", Value: []byte("test")}
	e3 := &Entry{Key: "foo/bar/baz", Value: []byte("test")}

	defer func() {
		b.Delete("foo")
		b.Delete("foo/bar")
		b.Delete("foo/bar/baz")
	}()

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

type delays struct {
	beforeGet  time.Duration
	beforeList time.Duration
}

func testEventuallyConsistentBackend(t *testing.T, b Backend, d delays) {

	// no delay required: nothing written to bucket
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

	// no delay required: nothing written to bucket
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
	time.Sleep(d.beforeGet)
	out, err = b.Get("foo")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !reflect.DeepEqual(out, e) {
		t.Fatalf("bad: %v expected: %v", out, e)
	}

	// List should not be empty
	time.Sleep(d.beforeList)
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
	time.Sleep(d.beforeList)
	keys, err = b.List("")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(keys) != 0 {
		t.Fatalf("bad: %v", keys)
	}

	// Get should fail
	time.Sleep(d.beforeGet)
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

	time.Sleep(d.beforeList)
	keys, err = b.List("")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(keys) != 2 {
		t.Fatalf("bad: %v", keys)
	}
	sort.Strings(keys)
	if keys[0] != "foo" || keys[1] != "foo/" {
		t.Fatalf("bad: %v", keys)
	}

	// Delete with children should work
	err = b.Delete("foo")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Get should return the child
	time.Sleep(d.beforeGet)
	out, err = b.Get("foo/bar")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out == nil {
		t.Fatalf("missing child")
	}

	// Removal of nested secret should not leave artifacts
	e = &Entry{Key: "foo/nested1/nested2/nested3", Value: []byte("baz")}
	err = b.Put(e)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	err = b.Delete("foo/nested1/nested2/nested3")
	if err != nil {
		t.Fatalf("failed to remove nested secret: %v", err)
	}

	time.Sleep(d.beforeList)
	keys, err = b.List("foo/")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if len(keys) != 1 {
		t.Fatalf("there should be only one key left after deleting nested "+
			"secret: %v", keys)
	}

	if keys[0] != "bar" {
		t.Fatalf("bad keys after deleting nested: %v", keys)
	}

	// Make a second nested entry to test prefix removal
	e = &Entry{Key: "foo/zip", Value: []byte("zap")}
	err = b.Put(e)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Delete should not remove the prefix
	err = b.Delete("foo/bar")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	time.Sleep(d.beforeList)
	keys, err = b.List("")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(keys) != 1 {
		t.Fatalf("bad: %v", keys)
	}
	if keys[0] != "foo/" {
		t.Fatalf("bad: %v", keys)
	}

	// Delete should remove the prefix
	err = b.Delete("foo/zip")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	time.Sleep(d.beforeList)
	keys, err = b.List("")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(keys) != 0 {
		t.Fatalf("bad: %v", keys)
	}
}

func testEventuallyConsistentBackend_ListPrefix(t *testing.T, b Backend, d delays) {
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
	time.Sleep(d.beforeList)
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
	time.Sleep(d.beforeList)
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
	time.Sleep(d.beforeList)
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
