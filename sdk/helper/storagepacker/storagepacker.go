package storagepacker

import (
	"context"

	"github.com/hashicorp/vault/sdk/logical"
)

type StoragePackerFactory func(context.Context, *Config) (StoragePacker, error)

type StoragePacker interface {
	BucketsView() *logical.StorageView
	BucketKey(string) string
	GetCacheKey(string) string
	BucketKeys(context.Context) ([]string, error)
	GetBucket(context.Context, string, bool) (*LockedBucket, error)
	DecodeBucket(*logical.StorageEntry) (*LockedBucket, error)
	PutBucket(context.Context, *LockedBucket) error
	DeleteBucket(context.Context, string) error
	DeleteItem(context.Context, string) error
	GetItem(context.Context, string) (*Item, error)
	PutItem(context.Context, *Item) error
}

/*

// Is the given storage path controlled by this StoragePacker?
func (s *StoragePackerV2) MatchingStorage(path string) bool

// Set the value of multiple items, as an efficient but non-atomic operation.
// Each affected bucket will be written just once.
// A multierror will be returned if some of the writes fail.
func (s *StoragePackerV2) PutItem(context.Context, items ...*Item) error

// Delete the specified set of itemIDs
func (s *StoragePackerV2) DeleteItem(context.Context, ids ...string) error

// Retrieve a single item
func (s *StoragePackerV2) GetItem(context.Context, id string) (*Item, error)

// Retrieve a set of items identified by ID; missing items signal by nil's in
// the returned slice so that the the request and response lengths are the same.
// The retrieval is non-atomic so may reflect partial changes made by PutItem
func (s *StoragePackerV2) GetItems(context.Context, ids ...string) ([]*Item, error)

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
