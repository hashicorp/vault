// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package physical

import (
	"context"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func ExerciseBackend(t testing.TB, b Backend) {
	t.Helper()
	ctx := context.Background()

	// Should be empty
	keys, err := b.List(ctx, "")
	if err != nil {
		t.Fatalf("initial list failed: %v", err)
	}
	if len(keys) != 0 {
		t.Errorf("initial not empty: %v", keys)
	}

	// Delete should work if it does not exist
	err = b.Delete(ctx, "foo")
	if err != nil {
		t.Fatalf("idempotent delete: %v", err)
	}

	// Get should not fail, but be nil
	out, err := b.Get(ctx, "foo")
	if err != nil {
		t.Fatalf("initial get failed: %v", err)
	}
	if out != nil {
		t.Errorf("initial get was not nil: %v", out)
	}

	// Make an entry
	e := &Entry{Key: "foo", Value: []byte("test")}
	err = b.Put(ctx, e)
	if err != nil {
		t.Fatalf("put failed: %v", err)
	}

	// Get should work
	out, err = b.Get(ctx, "foo")
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if !reflect.DeepEqual(out, e) {
		t.Errorf("bad: %v expected: %v", out, e)
	}

	// List should not be empty
	keys, err = b.List(ctx, "")
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if len(keys) != 1 || keys[0] != "foo" {
		t.Errorf("keys[0] did not equal foo: %v", keys)
	}

	// Delete should work
	err = b.Delete(ctx, "foo")
	if err != nil {
		t.Fatalf("delete: %v", err)
	}

	// Should be empty
	keys, err = b.List(ctx, "")
	if err != nil {
		t.Fatalf("list after delete: %v", err)
	}
	if len(keys) != 0 {
		t.Errorf("list after delete not empty: %v", keys)
	}

	// Get should fail
	out, err = b.Get(ctx, "foo")
	if err != nil {
		t.Fatalf("get after delete: %v", err)
	}
	if out != nil {
		t.Errorf("get after delete not nil: %v", out)
	}

	// Multiple Puts should work; GH-189
	e = &Entry{Key: "foo", Value: []byte("test")}
	err = b.Put(ctx, e)
	if err != nil {
		t.Fatalf("multi put 1 failed: %v", err)
	}
	e = &Entry{Key: "foo", Value: []byte("test")}
	err = b.Put(ctx, e)
	if err != nil {
		t.Fatalf("multi put 2 failed: %v", err)
	}

	// Make a nested entry
	e = &Entry{Key: "foo/bar", Value: []byte("baz")}
	err = b.Put(ctx, e)
	if err != nil {
		t.Fatalf("nested put failed: %v", err)
	}

	// Get should work
	out, err = b.Get(ctx, "foo/bar")
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if !reflect.DeepEqual(out, e) {
		t.Errorf("bad: %v expected: %v", out, e)
	}

	keys, err = b.List(ctx, "")
	if err != nil {
		t.Fatalf("list multi failed: %v", err)
	}
	sort.Strings(keys)
	if len(keys) != 2 || keys[0] != "foo" || keys[1] != "foo/" {
		t.Errorf("expected 2 keys [foo, foo/]: %v", keys)
	}

	// Delete with children should work
	err = b.Delete(ctx, "foo")
	if err != nil {
		t.Fatalf("delete after multi: %v", err)
	}

	// Get should return the child
	out, err = b.Get(ctx, "foo/bar")
	if err != nil {
		t.Fatalf("get after multi delete: %v", err)
	}
	if out == nil {
		t.Errorf("get after multi delete not nil: %v", out)
	}

	// Removal of nested secret should not leave artifacts
	e = &Entry{Key: "foo/nested1/nested2/nested3", Value: []byte("baz")}
	err = b.Put(ctx, e)
	if err != nil {
		t.Fatalf("deep nest: %v", err)
	}

	err = b.Delete(ctx, "foo/nested1/nested2/nested3")
	if err != nil {
		t.Fatalf("failed to remove deep nest: %v", err)
	}

	keys, err = b.List(ctx, "foo/")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(keys) != 1 || keys[0] != "bar" {
		t.Errorf("should be exactly 1 key == bar: %v", keys)
	}

	// Make a second nested entry to test prefix removal
	e = &Entry{Key: "foo/zip", Value: []byte("zap")}
	err = b.Put(ctx, e)
	if err != nil {
		t.Fatalf("failed to create second nested: %v", err)
	}

	// Delete should not remove the prefix
	err = b.Delete(ctx, "foo/bar")
	if err != nil {
		t.Fatalf("failed to delete nested prefix: %v", err)
	}

	keys, err = b.List(ctx, "")
	if err != nil {
		t.Fatalf("list nested prefix: %v", err)
	}
	if len(keys) != 1 || keys[0] != "foo/" {
		t.Errorf("should be exactly 1 key == foo/: %v", keys)
	}

	// Delete should remove the prefix
	err = b.Delete(ctx, "foo/zip")
	if err != nil {
		t.Fatalf("failed to delete second prefix: %v", err)
	}

	keys, err = b.List(ctx, "")
	if err != nil {
		t.Fatalf("listing after second delete failed: %v", err)
	}
	if len(keys) != 0 {
		t.Errorf("should be empty at end: %v", keys)
	}

	// When the root path is empty, adding and removing deep nested values should not break listing
	e = &Entry{Key: "foo/nested1/nested2/value1", Value: []byte("baz")}
	err = b.Put(ctx, e)
	if err != nil {
		t.Fatalf("deep nest: %v", err)
	}

	e = &Entry{Key: "foo/nested1/nested2/value2", Value: []byte("baz")}
	err = b.Put(ctx, e)
	if err != nil {
		t.Fatalf("deep nest: %v", err)
	}

	err = b.Delete(ctx, "foo/nested1/nested2/value2")
	if err != nil {
		t.Fatalf("failed to remove deep nest: %v", err)
	}

	keys, err = b.List(ctx, "")
	if err != nil {
		t.Fatalf("listing of root failed after deletion: %v", err)
	}
	if len(keys) == 0 {
		t.Errorf("root is returning empty after deleting a single nested value, expected nested1/: %v", keys)
		keys, err = b.List(ctx, "foo/nested1")
		if err != nil {
			t.Fatalf("listing of expected nested path 'foo/nested1' failed: %v", err)
		}
		// prove that the root should not be empty and that foo/nested1 exists
		if len(keys) != 0 {
			t.Logf("  keys can still be listed from nested1/ so it's not empty, expected nested2/: %v", keys)
		}
	}

	// cleanup left over listing bug test value
	err = b.Delete(ctx, "foo/nested1/nested2/value1")
	if err != nil {
		t.Fatalf("failed to remove deep nest: %v", err)
	}

	keys, err = b.List(ctx, "")
	if err != nil {
		t.Fatalf("listing of root failed after delete of deep nest: %v", err)
	}
	if len(keys) != 0 {
		t.Errorf("should be empty at end: %v", keys)
	}
}

func ExerciseBackend_ListPrefix(t testing.TB, b Backend) {
	t.Helper()
	ctx := context.Background()

	e1 := &Entry{Key: "foo", Value: []byte("test")}
	e2 := &Entry{Key: "foo/bar", Value: []byte("test")}
	e3 := &Entry{Key: "foo/bar/baz", Value: []byte("test")}

	defer func() {
		_ = b.Delete(ctx, "foo")
		_ = b.Delete(ctx, "foo/bar")
		_ = b.Delete(ctx, "foo/bar/baz")
	}()

	err := b.Put(ctx, e1)
	if err != nil {
		t.Fatalf("failed to put entry 1: %v", err)
	}
	err = b.Put(ctx, e2)
	if err != nil {
		t.Fatalf("failed to put entry 2: %v", err)
	}
	err = b.Put(ctx, e3)
	if err != nil {
		t.Fatalf("failed to put entry 3: %v", err)
	}

	// Scan the root
	keys, err := b.List(ctx, "")
	if err != nil {
		t.Fatalf("list root: %v", err)
	}
	sort.Strings(keys)
	if len(keys) != 2 || keys[0] != "foo" || keys[1] != "foo/" {
		t.Errorf("root expected [foo foo/]: %v", keys)
	}

	// Scan foo/
	keys, err = b.List(ctx, "foo/")
	if err != nil {
		t.Fatalf("list level 1: %v", err)
	}
	sort.Strings(keys)
	if len(keys) != 2 || keys[0] != "bar" || keys[1] != "bar/" {
		t.Errorf("level 1 expected [bar bar/]: %v", keys)
	}

	// Scan foo/bar/
	keys, err = b.List(ctx, "foo/bar/")
	if err != nil {
		t.Fatalf("list level 2: %v", err)
	}
	sort.Strings(keys)
	if len(keys) != 1 || keys[0] != "baz" {
		t.Errorf("level 1 expected [baz]: %v", keys)
	}
}

func ExerciseHABackend(t testing.TB, b HABackend, b2 HABackend) {
	t.Helper()

	// Get the lock
	lock, err := b.LockWith("foo", "bar")
	if err != nil {
		t.Fatalf("initial lock: %v", err)
	}

	// Attempt to lock
	leaderCh, err := lock.Lock(nil)
	if err != nil {
		t.Fatalf("lock attempt 1: %v", err)
	}
	if leaderCh == nil {
		t.Fatalf("missing leaderCh")
	}

	// Check the value
	held, val, err := lock.Value()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !held {
		t.Errorf("should be held")
	}
	if val != "bar" {
		t.Errorf("expected value bar: %v", err)
	}

	// Check if it's fencing that we can register the lock
	if fba, ok := b.(FencingHABackend); ok {
		require.NoError(t, fba.RegisterActiveNodeLock(lock))
	}

	// Second acquisition should fail
	lock2, err := b2.LockWith("foo", "baz")
	if err != nil {
		t.Fatalf("lock 2: %v", err)
	}

	// Checking the lock from b2 should discover that the lock is held since held
	// implies only that there is _some_ leader not that b2 is leader (this was
	// not clear before so we make it explicit with this assertion).
	held2, val2, err := lock2.Value()
	require.NoError(t, err)
	require.Equal(t, "bar", val2)
	require.True(t, held2)

	// Cancel attempt in 50 msec
	stopCh := make(chan struct{})
	time.AfterFunc(50*time.Millisecond, func() {
		close(stopCh)
	})

	// Attempt to lock
	leaderCh2, err := lock2.Lock(stopCh)
	if err != nil {
		t.Fatalf("stop lock 2: %v", err)
	}
	if leaderCh2 != nil {
		t.Errorf("should not have gotten leaderCh: %v", leaderCh2)
	}

	// Release the first lock
	lock.Unlock()

	// Attempt to lock should work
	leaderCh2, err = lock2.Lock(nil)
	if err != nil {
		t.Fatalf("lock 2 lock: %v", err)
	}
	if leaderCh2 == nil {
		t.Errorf("should get leaderCh")
	}

	// Check if it's fencing that we can register the lock
	if fba2, ok := b2.(FencingHABackend); ok {
		require.NoError(t, fba2.RegisterActiveNodeLock(lock))
	}

	// Check the value
	held, val, err = lock2.Value()
	if err != nil {
		t.Fatalf("value: %v", err)
	}
	if !held {
		t.Errorf("should still be held")
	}
	if val != "baz" {
		t.Errorf("expected: baz, got: %v", val)
	}

	// Cleanup
	lock2.Unlock()
}

func ExerciseTransactionalBackend(t testing.TB, b Backend) {
	t.Helper()
	ctx := context.Background()

	tb, ok := b.(Transactional)
	if !ok {
		t.Fatal("Not a transactional backend")
	}

	txns := SetupTestingTransactions(t, b)

	if err := tb.Transaction(ctx, txns); err != nil {
		t.Fatal(err)
	}

	keys, err := b.List(ctx, "")
	if err != nil {
		t.Fatal(err)
	}

	expected := []string{"foo", "zip"}

	sort.Strings(keys)
	sort.Strings(expected)
	if !reflect.DeepEqual(keys, expected) {
		t.Fatalf("mismatch: expected\n%#v\ngot\n%#v\n", expected, keys)
	}

	entry, err := b.Get(ctx, "foo")
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

	entry, err = b.Get(ctx, "zip")
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

func SetupTestingTransactions(t testing.TB, b Backend) []*TxnEntry {
	t.Helper()
	ctx := context.Background()

	// Add a few keys so that we test rollback with deletion
	if err := b.Put(ctx, &Entry{
		Key:   "foo",
		Value: []byte("bar"),
	}); err != nil {
		t.Fatal(err)
	}
	if err := b.Put(ctx, &Entry{
		Key:   "zip",
		Value: []byte("zap"),
	}); err != nil {
		t.Fatal(err)
	}
	if err := b.Put(ctx, &Entry{
		Key: "deleteme",
	}); err != nil {
		t.Fatal(err)
	}
	if err := b.Put(ctx, &Entry{
		Key: "deleteme2",
	}); err != nil {
		t.Fatal(err)
	}

	txns := []*TxnEntry{
		{
			Operation: PutOperation,
			Entry: &Entry{
				Key:   "foo",
				Value: []byte("bar2"),
			},
		},
		{
			Operation: DeleteOperation,
			Entry: &Entry{
				Key: "deleteme",
			},
		},
		{
			Operation: PutOperation,
			Entry: &Entry{
				Key:   "foo",
				Value: []byte("bar3"),
			},
		},
		{
			Operation: DeleteOperation,
			Entry: &Entry{
				Key: "deleteme2",
			},
		},
		{
			Operation: PutOperation,
			Entry: &Entry{
				Key:   "zip",
				Value: []byte("zap3"),
			},
		},
	}

	return txns
}

// Several tests across packages have to test logic with a few variations of
// transactional backends. Make some suitable for testing limits support that
// can be re-used.

type TestTransactionalNonLimitBackend struct{}

var _ Transactional = (*TestTransactionalNonLimitBackend)(nil)

func (b *TestTransactionalNonLimitBackend) Put(ctx context.Context, entry *Entry) error {
	return nil
}

func (b *TestTransactionalNonLimitBackend) Get(ctx context.Context, key string) (*Entry, error) {
	return nil, nil
}

func (b *TestTransactionalNonLimitBackend) Delete(ctx context.Context, key string) error {
	return nil
}

func (b *TestTransactionalNonLimitBackend) List(ctx context.Context, prefix string) ([]string, error) {
	return nil, nil
}

func (b *TestTransactionalNonLimitBackend) Transaction(ctx context.Context, txns []*TxnEntry) error {
	return nil
}

type TestTransactionalLimitBackend struct {
	TestTransactionalNonLimitBackend

	MaxEntries, MaxSize int
}

var _ TransactionalLimits = (*TestTransactionalLimitBackend)(nil)

func (b *TestTransactionalLimitBackend) TransactionLimits() (int, int) {
	return b.MaxEntries, b.MaxSize
}
