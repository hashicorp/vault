// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package snapshots

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

type SnapshotTestCase struct {
	backend         logical.Backend
	regularStorage  logical.Storage
	snapshotStorage logical.Storage
	storageRouter   logical.Storage
}

type storageProvider struct {
	storage logical.Storage
}

func (s *storageProvider) SnapshotStorage(ctx context.Context, id string) (logical.Storage, error) {
	return s.storage, nil
}

var _ logical.SnapshotStorageProvider = (*storageProvider)(nil)

// NewSnapshotTestCase is used to create a snapshot test case for a particular
// backend. The test case is used to ensure that the backend behaves correctly
// when it receives snapshot operations, without having to do the end-to-end
// setup of creating a raft cluster, taking a snapshot, and loading it.
func NewSnapshotTestCase(t testing.TB, backend logical.Backend) *SnapshotTestCase {
	s := &SnapshotTestCase{
		backend:         backend,
		regularStorage:  &logical.InmemStorage{},
		snapshotStorage: &logical.InmemStorage{},
	}

	s.storageRouter = logical.NewSnapshotStorageRouter(s.regularStorage, &storageProvider{s.snapshotStorage})
	return s
}

func (s *SnapshotTestCase) SnapshotStorage() logical.Storage {
	return s.snapshotStorage
}

func (s *SnapshotTestCase) RegularStorage() logical.Storage {
	return s.regularStorage
}

func (s *SnapshotTestCase) runCase(t testing.TB, path string, op logical.Operation) {
	ctx := context.Background()
	normalResp, err := s.backend.HandleRequest(ctx, &logical.Request{
		Path:      path,
		Operation: op,
		Storage:   s.storageRouter,
	})
	require.NoError(t, err)

	_, err = s.backend.HandleRequest(logical.CreateContextWithSnapshotID(ctx, "snapshot_id"), &logical.Request{
		Path:               path,
		Operation:          op,
		Storage:            s.storageRouter,
		RequiresSnapshotID: "snapshot_id",
	})
	require.NoError(t, err)

	normalResp2, err := s.backend.HandleRequest(ctx, &logical.Request{
		Path:      path,
		Operation: op,
		Storage:   s.storageRouter,
	})
	require.NoError(t, err)
	if normalResp == nil || normalResp2 == nil {
		require.Equal(t, normalResp, normalResp2)
	} else {
		require.Equal(t, normalResp.Data, normalResp2.Data)
	}
}

// RunList runs a list operation without a snapshot, a list operation from a
// snapshot, and then another list operation without a snapshot. The test
// verifies that the list operation from the snapshot does not cause the results
// to change
func (s *SnapshotTestCase) RunList(t testing.TB, path string) {
	s.runCase(t, path, logical.ListOperation)
}

// RunRead runs a read operation without a snapshot, a read operation from a
// snapshot, and then another read operation without a snapshot. The test
// verifies that the read operation from the snapshot does not cause the results
// to change
func (s *SnapshotTestCase) RunRead(t testing.TB, path string) {
	s.runCase(t, path, logical.ReadOperation)
}

// DoRecover performs a read operation from a snapshot, and then a recover. The
// test returns the results of the recover operation
func (s *SnapshotTestCase) DoRecover(t testing.TB, path string) (*logical.Response, error) {
	ctx := context.Background()
	readResp, err := s.backend.HandleRequest(logical.CreateContextWithSnapshotID(ctx, "snapshot_id"), &logical.Request{
		Path:               path,
		Operation:          logical.ReadOperation,
		Storage:            s.storageRouter,
		RequiresSnapshotID: "snapshot_id",
	})
	require.NoError(t, err)
	require.NotNil(t, readResp)

	return s.backend.HandleRequest(ctx, &logical.Request{
		Path:               path,
		Operation:          logical.RecoverOperation,
		Storage:            s.storageRouter,
		Data:               readResp.Data,
		RequiresSnapshotID: "snapshot_id",
	})
}
