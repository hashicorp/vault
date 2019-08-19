package storagepacker

import (
	"context"
)

type StoragePackerFactory func(context.Context, *Config) (StoragePacker, error)

type StoragePacker interface {
	// Set the value of multiple items, as an efficient but non-atomic operation.
	// Each affected bucket will be written just once.
	// A multierror will be returned if some of the writes fail.
	PutItem(context.Context, ...*Item) error

	// Delete the specified set of itemIDs
	DeleteItem(context.Context, ...string) error

	// Retrieve a set of items identified by ID; missing items have nil Value,
	// so that the the request and response lengths are the same.
	GetItems(context.Context, ...string) ([]*Item, error)

	// Special case single item to avoid the slice
	GetItem(context.Context, string) (*Item, error)

	// Retrieve the entire set of items as a single slice
	// This *is* an atomic operation.
	AllItems(context.Context) ([]*Item, error)

	// Defer writes until FlushQueue is called
	// For correctness, the caller must prevent concurrent access
	// during flushing or while turing queue mode back off.
	SetQueueMode(bool)
	FlushQueue(context.Context) error

	// Is the given storage path controlled by this StoragePacker?
	MatchingStorage(string) bool

	// Given a storage path (to a bucket) and new contents, invalidate the existing bucket and
	// reload it.
	InvalidateItems(context.Context, string, []byte) ([]*Item, []*Item, error)
}
