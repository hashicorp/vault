// Copyright (c) HashiCorp, Inc
// SPDX-License-Identifier: MPL-2.0

package types

import "io"

// MetaStore is the interface we need to some persistent, crash safe backend. We
// implement it with BoltDB for real usage but the interface allows alternatives
// to be used, or tests to mock out FS access.
type MetaStore interface {
	// Load loads the existing persisted state. If there is no existing state
	// implementations are expected to create initialize new storage and return an
	// empty state.
	Load(dir string) (PersistentState, error)

	// CommitState must atomically replace all persisted metadata in the current
	// store with the set provided. It must not return until the data is persisted
	// durably and in a crash-safe way otherwise the guarantees of the WAL will be
	// compromised. The WAL will only ever call this in a single thread at one
	// time and it will never be called concurrently with Load however it may be
	// called concurrently with Get/SetStable operations.
	CommitState(PersistentState) error

	// GetStable returns a value from stable store or nil if it doesn't exist. May
	// be called concurrently by multiple threads.
	GetStable(key []byte) ([]byte, error)

	// SetStable stores a value from stable store. May be called concurrently with
	// GetStable.
	SetStable(key, value []byte) error

	io.Closer
}

// PersistentState represents the WAL file metadata we need to store reliably to
// recover on restart.
type PersistentState struct {
	NextSegmentID uint64
	Segments      []SegmentInfo
}
