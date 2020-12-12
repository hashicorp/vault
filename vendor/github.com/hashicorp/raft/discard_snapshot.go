package raft

import (
	"fmt"
	"io"
)

// DiscardSnapshotStore is used to successfully snapshot while
// always discarding the snapshot. This is useful for when the
// log should be truncated but no snapshot should be retained.
// This should never be used for production use, and is only
// suitable for testing.
type DiscardSnapshotStore struct{}

// DiscardSnapshotSink is used to fulfill the SnapshotSink interface
// while always discarding the . This is useful for when the log
// should be truncated but no snapshot should be retained. This
// should never be used for production use, and is only suitable
// for testing.
type DiscardSnapshotSink struct{}

// NewDiscardSnapshotStore is used to create a new DiscardSnapshotStore.
func NewDiscardSnapshotStore() *DiscardSnapshotStore {
	return &DiscardSnapshotStore{}
}

// Create returns a valid type implementing the SnapshotSink which
// always discards the snapshot.
func (d *DiscardSnapshotStore) Create(version SnapshotVersion, index, term uint64,
	configuration Configuration, configurationIndex uint64, trans Transport) (SnapshotSink, error) {
	return &DiscardSnapshotSink{}, nil
}

// List returns successfully with a nil for []*SnapshotMeta.
func (d *DiscardSnapshotStore) List() ([]*SnapshotMeta, error) {
	return nil, nil
}

// Open returns an error since the DiscardSnapshotStore does not
// support opening snapshots.
func (d *DiscardSnapshotStore) Open(id string) (*SnapshotMeta, io.ReadCloser, error) {
	return nil, nil, fmt.Errorf("open is not supported")
}

// Write returns successfully with the length of the input byte slice
// to satisfy the WriteCloser interface
func (d *DiscardSnapshotSink) Write(b []byte) (int, error) {
	return len(b), nil
}

// Close returns a nil error
func (d *DiscardSnapshotSink) Close() error {
	return nil
}

// ID returns "discard" for DiscardSnapshotSink
func (d *DiscardSnapshotSink) ID() string {
	return "discard"
}

// Cancel returns successfully with a nil error
func (d *DiscardSnapshotSink) Cancel() error {
	return nil
}
