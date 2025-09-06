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
	return NewSnapshotTestCaseWithStorages(t, backend, &logical.InmemStorage{}, &logical.InmemStorage{})
}

// NewSnapshotTestCaseWithStorages is used to create a snapshot test case for a
// particular backend, using the provided storage instances. The test case is
// used to ensure that the backend behaves correctly when it receives snapshot
// operations, without having to do the end-to-end setup of creating a raft
// cluster, taking a snapshot, and loading it.
func NewSnapshotTestCaseWithStorages(t testing.TB, backend logical.Backend, regularStorage, snapshotStorage logical.Storage) *SnapshotTestCase {
	s := &SnapshotTestCase{
		backend:         backend,
		regularStorage:  regularStorage,
		snapshotStorage: snapshotStorage,
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

func (s *SnapshotTestCase) runCase(t testing.TB, path string, op logical.Operation, additionalOpts ...Option) {
	opts := options{}
	for _, opt := range additionalOpts {
		opt(&opts)
	}
	ctx := context.Background()
	req1 := &logical.Request{
		Path:      path,
		Operation: op,
		Storage:   s.storageRouter,
	}
	if opts.modifyRequests != nil {
		opts.modifyRequests(req1)
	}
	normalResp, err := s.backend.HandleRequest(ctx, req1)
	require.NoError(t, err)

	req2 := &logical.Request{
		Path:               path,
		Operation:          op,
		Storage:            s.storageRouter,
		RequiresSnapshotID: "snapshot_id",
	}
	if opts.modifyRequests != nil {
		opts.modifyRequests(req2)
	}
	_, err = s.backend.HandleRequest(logical.CreateContextWithSnapshotID(ctx, "snapshot_id"), req2)
	require.NoError(t, err)

	req3 := &logical.Request{
		Path:      path,
		Operation: op,
		Storage:   s.storageRouter,
	}
	if opts.modifyRequests != nil {
		opts.modifyRequests(req3)
	}
	normalResp2, err := s.backend.HandleRequest(ctx, req3)
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
func (s *SnapshotTestCase) RunList(t testing.TB, path string, additionalOpts ...Option) {
	s.runCase(t, path, logical.ListOperation, additionalOpts...)
}

// RunRead runs a read operation without a snapshot, a read operation from a
// snapshot, and then another read operation without a snapshot. The test
// verifies that the read operation from the snapshot does not cause the results
// to change
func (s *SnapshotTestCase) RunRead(t testing.TB, path string, additionalOpts ...Option) {
	s.runCase(t, path, logical.ReadOperation, additionalOpts...)
}

// DoRecover performs a read operation from a snapshot, and then a recover. The
// test returns the results of the recover operation
func (s *SnapshotTestCase) DoRecover(t testing.TB, path string, additionalOpts ...Option) (*logical.Response, error) {
	opts := options{}
	for _, opt := range additionalOpts {
		opt(&opts)
	}

	ctx := context.Background()
	readPath := path
	if opts.recoverSourcePath != "" {
		readPath = opts.recoverSourcePath
	}
	readReq := &logical.Request{
		Path:               readPath,
		Operation:          logical.ReadOperation,
		Storage:            s.storageRouter,
		RequiresSnapshotID: "snapshot_id",
	}
	if opts.modifyRequests != nil {
		opts.modifyRequests(readReq)
	}
	readResp, err := s.backend.HandleRequest(logical.CreateContextWithSnapshotID(ctx, "snapshot_id"), readReq)
	require.NoError(t, err)
	require.NotNil(t, readResp)

	recoverSourcePath := ""
	if opts.recoverSourcePath != "" {
		recoverSourcePath = opts.recoverSourcePath
	}
	recoverReq := &logical.Request{
		Path:               path,
		Operation:          logical.RecoverOperation,
		Storage:            s.storageRouter,
		Data:               readResp.Data,
		RequiresSnapshotID: "snapshot_id",
		RecoverSourcePath:  recoverSourcePath,
	}
	if opts.modifyRequests != nil {
		opts.modifyRequests(recoverReq)
	}
	return s.backend.HandleRequest(ctx, recoverReq)
}

type options struct {
	modifyRequests    func(req *logical.Request)
	recoverSourcePath string
}
type Option func(o *options)

// WithModifyRequests allows you to modify the request before it is sent to the
// backend. This is useful for setting the client token or other request
// data
func WithModifyRequests(modify func(req *logical.Request)) Option {
	return func(o *options) {
		o.modifyRequests = modify
	}
}

// WithRecoverSourcePath allows you to specify a different path to read from when
// performing a recover operation.
func WithRecoverSourcePath(recoverSourcePath string) Option {
	return func(o *options) {
		o.recoverSourcePath = recoverSourcePath
	}
}
