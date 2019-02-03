package storagepacker

import (
	"context"

	"github.com/hashicorp/vault/logical"
)

type StoragePackerFactory func(context.Context, *Config) (StoragePacker, error)

type StoragePacker interface {
	BucketsView() *logical.StorageView
	BucketKeyHashByItemID(string) string
	GetCacheKey(string) string
	GetBucket(context.Context, string) (*LockedBucket, error)
	DecodeBucket(*logical.StorageEntry) (*LockedBucket, error)
	PutBucket(context.Context, *LockedBucket) error
	DeleteBucket(context.Context, string) error
	DeleteItem(context.Context, string) error
	GetItem(context.Context, string) (*Item, error)
	PutItem(context.Context, *Item) error
	SetQueueMode(enabled bool)
	FlushQueue(context.Context) error
}
