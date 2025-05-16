// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logical

import (
	"context"
	"errors"
)

type snapshotStorageView struct {
	storage    Storage
	snapshotID string
}

var readOnlyErr = errors.New("read-only storage view")

func (s *snapshotStorageView) List(ctx context.Context, prefix string) ([]string, error) {
	return s.storage.List(CreateContextWithSnapshotID(ctx, s.snapshotID), prefix)
}

func (s *snapshotStorageView) Get(ctx context.Context, key string) (*StorageEntry, error) {
	return s.storage.Get(CreateContextWithSnapshotID(ctx, s.snapshotID), key)
}

func (s *snapshotStorageView) Put(_ context.Context, _ *StorageEntry) error {
	return readOnlyErr
}

func (s *snapshotStorageView) Delete(_ context.Context, _ string) error {
	return readOnlyErr
}

// NewSnapshotStorageView creates a storage view that provides read-only access
// to the given snapshot's storage.
func NewSnapshotStorageView(request *Request) (Storage, error) {
	if request.RequiresSnapshotID == "" {
		return nil, errors.New("no snapshot in request")
	}
	return &snapshotStorageView{
		storage:    request.Storage,
		snapshotID: request.RequiresSnapshotID,
	}, nil
}

// SnapshotStorageProvider is an interface that provides a method to retrieve
// the snapshot by ID
type SnapshotStorageProvider interface {
	SnapshotStorage(ctx context.Context, id string) (Storage, error)
}

type SnapshotStorageRouter struct {
	underlying Storage
	manager    SnapshotStorageProvider
}

// NewSnapshotStorageRouter creates a new storage instance that routes to either
// a snapshot (given the snapshot's ID in the context) or to the underlying
// storage
func NewSnapshotStorageRouter(underlying Storage, manager SnapshotStorageProvider) Storage {
	return &SnapshotStorageRouter{
		underlying: underlying,
		manager:    manager,
	}
}

var _ Storage = (*SnapshotStorageRouter)(nil)

func (s *SnapshotStorageRouter) getStorageRead(ctx context.Context) (Storage, error) {
	snapID, ok := ContextSnapshotIDValue(ctx)
	if ok && snapID != "" {
		snapshotStorage, err := s.manager.SnapshotStorage(ctx, snapID)
		if err != nil {
			return nil, err
		}
		return snapshotStorage, nil
	}
	return s.underlying, nil
}

func (s *SnapshotStorageRouter) List(ctx context.Context, key string) ([]string, error) {
	storage, err := s.getStorageRead(ctx)
	if err != nil {
		return nil, err
	}
	return storage.List(ctx, key)
}

func (s *SnapshotStorageRouter) Get(ctx context.Context, key string) (*StorageEntry, error) {
	storage, err := s.getStorageRead(ctx)
	if err != nil {
		return nil, err
	}
	return storage.Get(ctx, key)
}

var WriteErr = errors.New("attempted write operation on snapshot")

func (s *SnapshotStorageRouter) checkWrite(ctx context.Context) error {
	snapID, ok := ContextSnapshotIDValue(ctx)
	if ok && snapID != "" {
		return WriteErr
	}
	return nil
}

func (s *SnapshotStorageRouter) Put(ctx context.Context, entry *StorageEntry) error {
	if err := s.checkWrite(ctx); err != nil {
		return err
	}
	return s.underlying.Put(ctx, entry)
}

func (s *SnapshotStorageRouter) Delete(ctx context.Context, key string) error {
	if err := s.checkWrite(ctx); err != nil {
		return err
	}
	return s.underlying.Delete(ctx, key)
}
