package physical

import (
	"reflect"
	"sort"
	"testing"
	"time"
)

func ExerciseBackend(t *testing.T, b Backend) {
	t.Helper()
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

func ExerciseBackend_ListPrefix(t *testing.T, b Backend) {
	t.Helper()
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

func ExerciseHABackend(t *testing.T, b HABackend, b2 HABackend) {
	t.Helper()
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

func ExerciseTransactionalBackend(t *testing.T, b Backend) {
	t.Helper()
	tb, ok := b.(Transactional)
	if !ok {
		t.Fatal("Not a transactional backend")
	}

	txns := SetupTestingTransactions(t, b)

	if err := tb.Transaction(txns); err != nil {
		t.Fatal(err)
	}

	keys, err := b.List("")
	if err != nil {
		t.Fatal(err)
	}

	expected := []string{"foo", "zip"}

	sort.Strings(keys)
	sort.Strings(expected)
	if !reflect.DeepEqual(keys, expected) {
		t.Fatalf("mismatch: expected\n%#v\ngot\n%#v\n", expected, keys)
	}

	entry, err := b.Get("foo")
	if err != nil {
		t.Fatal(err)
	}
	if entry == nil {
		t.Fatal("got nil entry")
	}
	if entry.Value == nil {
		t.Fatal("got nil value")
	}
	if string(entry.Value) != "bar3" {
		t.Fatal("updates did not apply correctly")
	}

	entry, err = b.Get("zip")
	if err != nil {
		t.Fatal(err)
	}
	if entry == nil {
		t.Fatal("got nil entry")
	}
	if entry.Value == nil {
		t.Fatal("got nil value")
	}
	if string(entry.Value) != "zap3" {
		t.Fatal("updates did not apply correctly")
	}
}

func SetupTestingTransactions(t *testing.T, b Backend) []*TxnEntry {
	t.Helper()
	// Add a few keys so that we test rollback with deletion
	if err := b.Put(&Entry{
		Key:   "foo",
		Value: []byte("bar"),
	}); err != nil {
		t.Fatal(err)
	}
	if err := b.Put(&Entry{
		Key:   "zip",
		Value: []byte("zap"),
	}); err != nil {
		t.Fatal(err)
	}
	if err := b.Put(&Entry{
		Key: "deleteme",
	}); err != nil {
		t.Fatal(err)
	}
	if err := b.Put(&Entry{
		Key: "deleteme2",
	}); err != nil {
		t.Fatal(err)
	}

	txns := []*TxnEntry{
		&TxnEntry{
			Operation: PutOperation,
			Entry: &Entry{
				Key:   "foo",
				Value: []byte("bar2"),
			},
		},
		&TxnEntry{
			Operation: DeleteOperation,
			Entry: &Entry{
				Key: "deleteme",
			},
		},
		&TxnEntry{
			Operation: PutOperation,
			Entry: &Entry{
				Key:   "foo",
				Value: []byte("bar3"),
			},
		},
		&TxnEntry{
			Operation: DeleteOperation,
			Entry: &Entry{
				Key: "deleteme2",
			},
		},
		&TxnEntry{
			Operation: PutOperation,
			Entry: &Entry{
				Key:   "zip",
				Value: []byte("zap3"),
			},
		},
	}

	return txns
}
