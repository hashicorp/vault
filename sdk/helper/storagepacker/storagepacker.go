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

	// Retrieve a set of items identified by ID; missing items signal by nil's in
	// the returned slice so that the the request and response lengths are the same.
	GetItems(context.Context, ...string) ([]*Item, error)

	// Special case single item to avoid the slice
	GetItem(context.Context, string) (*Item, error)

	// Retrieve the entire set of items as a single slice
	// This *is* an atomic operation.
	AllItems(context.Context) ([]*Item, error)
}

/*

// Is the given storage path controlled by this StoragePacker?
func (s *StoragePackerV2) MatchingStorage(path string) bool

// Retrieve the entire set of items as a single slice
// This *is* an atomic operation.
func (s *StoragePackerV2) AllItems(context.Context) ([]*Item, error)

// Given a storage path (to a bucket) and new contents, return
// "present": all items present in the new bucket
// "deleted": all items that were in the cache but are now absent
// For normal WAL operation this is serialized based on the order PutItem
// wrote them; for recovery from the Merkle tree this will not be the case.
func (s *StoragePackerV2) InvalidateItems(context.Context, path string, newValu []byte) (present []*Item, deleted []*Item, error)

// Defer writes until FlushQueue is called
func (s *StoragePackerV2) SetQueueMode(enabled bool)
func (s *StoragePackerV2) FlushQueue(context.Context) error

*/
