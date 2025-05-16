// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logical

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type checkCtxStorage struct {
	Storage
	requiredSnapshotID string
}

func (c *checkCtxStorage) Get(ctx context.Context, key string) (*StorageEntry, error) {
	if sID, _ := ContextSnapshotIDValue(ctx); sID != c.requiredSnapshotID {
		return nil, errors.New("snapshot ID mismatch")
	}
	return c.Storage.Get(ctx, key)
}

// TestSnapshotStorageView_BadRequest verifies that a request without a
// RequiresSnapshotID results in an error when creating a snapshot storage view.
func TestSnapshotStorageView_BadRequest(t *testing.T) {
	req := &Request{
		Storage: &InmemStorage{},
	}
	_, err := NewSnapshotStorageView(req)
	require.Error(t, err)
}

// TestSnapshotStorageView verifies that a request with a RequiredSnapshotID
// creates a read-only storage view. This storage will error on writes, but
// succeed on reads with the given snapshot ID in the context
func TestSnapshotStorageView(t *testing.T) {
	storage := &checkCtxStorage{
		Storage:            &InmemStorage{},
		requiredSnapshotID: "abcd",
	}

	req := &Request{
		RequiresSnapshotID: "abcd",
		Storage:            storage,
	}
	view, err := NewSnapshotStorageView(req)
	require.NoError(t, err)

	// write should error
	require.Error(t, view.Put(context.Background(), &StorageEntry{
		Key:   "foo",
		Value: []byte("bar"),
	}))
	require.Error(t, view.Delete(context.Background(), "foo"))

	// reads should contain the snapshot ID
	_, err = view.Get(context.Background(), "foo")
	require.NoError(t, err)

	_, err = view.List(context.Background(), "foo")
	require.NoError(t, err)
}

type storageProvider struct {
	snapshotStorage Storage
}

func (s *storageProvider) SnapshotStorage(ctx context.Context, id string) (Storage, error) {
	return s.snapshotStorage, nil
}

// TestSnapshotStorageRouter verifies that the snapshot storage router correctly
// routes requests to the appropriate storage based on the presence of the
// snapshot ID context key
func TestSnapshotStorageRouter(t *testing.T) {
	regularStorage := &InmemStorage{}
	snapshotStorage := &InmemStorage{}

	router := NewSnapshotStorageRouter(regularStorage, &storageProvider{
		snapshotStorage: snapshotStorage,
	})

	snapshotStorage.Put(context.Background(), &StorageEntry{
		Key:   "key",
		Value: []byte("bar"),
	})
	snapshotStorage.Put(context.Background(), &StorageEntry{
		Key:   "key2",
		Value: []byte("baz"),
	})
	regularStorage.Put(context.Background(), &StorageEntry{
		Key:   "key",
		Value: []byte("foo"),
	})

	ctx := context.Background()
	snapshotCtx := CreateContextWithSnapshotID(ctx, "snapshot_id")

	// test get with and without snapshot ID
	value, err := router.Get(ctx, "key")
	require.NoError(t, err)
	require.Equal(t, []byte("foo"), value.Value)

	value, err = router.Get(snapshotCtx, "key")
	require.NoError(t, err)
	require.Equal(t, []byte("bar"), value.Value)

	// test list with and without snapshot ID
	keys, err := router.List(ctx, "key")
	require.NoError(t, err)
	require.Len(t, keys, 1)

	keys, err = router.List(snapshotCtx, "key")
	require.NoError(t, err)
	require.Len(t, keys, 2)

	// delete with a snapshot ID should fail
	require.Error(t, router.Delete(snapshotCtx, "key"))

	// delete without a snapshot ID should succeed
	require.NoError(t, router.Delete(ctx, "key"))

	// put with a snapshot ID should fail
	require.Error(t, router.Put(snapshotCtx, &StorageEntry{
		Key:   "key",
		Value: []byte("new-value"),
	}))

	// put without a snapshot ID should succeed
	require.NoError(t, router.Put(ctx, &StorageEntry{
		Key:   "key",
		Value: []byte("new-value"),
	}))
}
