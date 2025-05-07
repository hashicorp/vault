// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package storagepacker

import (
	"context"
	"crypto/md5"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/armon/go-metrics"
	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/go-hclog"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/compressutil"
	"github.com/hashicorp/vault/sdk/helper/locksutil"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	bucketCount = 256
	// StoragePackerBucketsPrefix is the default storage key prefix under which
	// bucket data will be stored.
	StoragePackerBucketsPrefix = "packer/buckets/"
)

// StoragePacker packs items into a specific number of buckets by hashing
// its identifier and indexing on it. Currently this supports only 256 bucket entries and
// hence relies on the first byte of the hash value for indexing.
type StoragePacker struct {
	view         logical.Storage
	logger       log.Logger
	storageLocks []*locksutil.LockEntry
	viewPrefix   string
}

// View returns the storage view configured to be used by the packer
func (s *StoragePacker) View() logical.Storage {
	return s.view
}

// GetBucket returns a bucket for a given key
func (s *StoragePacker) GetBucket(ctx context.Context, key string) (*Bucket, error) {
	if key == "" {
		return nil, fmt.Errorf("missing bucket key")
	}

	lock := locksutil.LockForKey(s.storageLocks, key)
	lock.RLock()
	defer lock.RUnlock()

	// Read from storage
	storageEntry, err := s.view.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to read packed storage entry: %w", err)
	}
	if storageEntry == nil {
		return nil, nil
	}

	uncompressedData, notCompressed, err := compressutil.Decompress(storageEntry.Value)
	if err != nil {
		return nil, fmt.Errorf("failed to decompress packed storage entry: %w", err)
	}
	if notCompressed {
		uncompressedData = storageEntry.Value
	}

	var bucket Bucket
	err = proto.Unmarshal(uncompressedData, &bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to decode packed storage entry: %w", err)
	}

	return &bucket, nil
}

// upsert either inserts a new item into the bucket or updates an existing one
// if an item with a matching key is already present.
func (s *Bucket) upsert(item *Item) error {
	if s == nil {
		return fmt.Errorf("nil storage bucket")
	}

	if item == nil {
		return fmt.Errorf("nil item")
	}

	if item.ID == "" {
		return fmt.Errorf("missing item ID")
	}

	// Look for an item with matching key and don't modify the collection while
	// iterating
	foundIdx := -1
	for itemIdx, bucketItems := range s.Items {
		if bucketItems.ID == item.ID {
			foundIdx = itemIdx
			break
		}
	}

	// If there is no match, append the item, otherwise update it
	if foundIdx == -1 {
		s.Items = append(s.Items, item)
	} else {
		s.Items[foundIdx] = item
	}

	return nil
}

// BucketKey returns the storage key of the bucket where the given item will be
// stored.
func (s *StoragePacker) BucketKey(itemID string) string {
	hf := md5.New()
	input := []byte(itemID)
	n, err := hf.Write(input)
	// Make linter happy
	if err != nil || n != len(input) {
		return ""
	}
	index := uint8(hf.Sum(nil)[0])
	return s.viewPrefix + strconv.Itoa(int(index))
}

// DeleteItem removes the item from the respective bucket
func (s *StoragePacker) DeleteItem(ctx context.Context, itemID string) error {
	return s.DeleteMultipleItems(ctx, nil, []string{itemID})
}

func (s *StoragePacker) DeleteMultipleItems(ctx context.Context, logger hclog.Logger, itemIDs []string) error {
	defer metrics.MeasureSince([]string{"storage_packer", "delete_items"}, time.Now())
	if len(itemIDs) == 0 {
		return nil
	}

	if logger == nil {
		logger = hclog.NewNullLogger()
	}

	// Sort the ids by the bucket they will be deleted from
	lockKeys := make([]string, 0)
	byBucket := make(map[string]map[string]struct{})
	for _, id := range itemIDs {
		bucketKey := s.BucketKey(id)
		bucket, ok := byBucket[bucketKey]
		if !ok {
			bucket = make(map[string]struct{})
			byBucket[bucketKey] = bucket

			// Add the lock key once
			lockKeys = append(lockKeys, bucketKey)
		}

		bucket[id] = struct{}{}
	}

	locks := locksutil.LocksForKeys(s.storageLocks, lockKeys)
	for _, lock := range locks {
		lock.Lock()
		defer lock.Unlock()
	}

	logger.Debug("deleting multiple items from storagepacker; caching and deleting from buckets", "total_items", len(itemIDs))

	// For each bucket, load from storage, remove the necessary items, and add
	// write it back out to storage
	pctDone := 0
	idx := 0
	for bucketKey, itemsToRemove := range byBucket {
		// Read bucket from storage
		storageEntry, err := s.view.Get(ctx, bucketKey)
		if err != nil {
			return fmt.Errorf("failed to read packed storage value: %w", err)
		}
		if storageEntry == nil {
			logger.Warn("could not find bucket", "bucket", bucketKey)
			continue
		}

		uncompressedData, notCompressed, err := compressutil.Decompress(storageEntry.Value)
		if err != nil {
			return fmt.Errorf("failed to decompress packed storage value: %w", err)
		}
		if notCompressed {
			uncompressedData = storageEntry.Value
		}

		bucket := new(Bucket)
		err = proto.Unmarshal(uncompressedData, bucket)
		if err != nil {
			return fmt.Errorf("failed decoding packed storage entry: %w", err)
		}

		// Look for a matching storage entries and delete them from the list.
		for i := 0; i < len(bucket.Items); i++ {
			if _, ok := itemsToRemove[bucket.Items[i].ID]; ok {
				copy(bucket.Items[i:], bucket.Items[i+1:])
				bucket.Items[len(bucket.Items)-1] = nil // allow GC
				bucket.Items = bucket.Items[:len(bucket.Items)-1]

				// Since we just moved a value to position i we need to
				// decrement i so we replay this position
				i--
			}
		}

		// Fail if the context is canceled, the storage calls will fail anyways
		if ctx.Err() != nil {
			return ctx.Err()
		}

		err = s.putBucket(ctx, bucket)
		if err != nil {
			return err
		}

		newPctDone := idx * 100.0 / len(byBucket)
		if int(newPctDone) > pctDone {
			pctDone = int(newPctDone)
			logger.Trace("bucket persistence progress", "percent", pctDone, "buckets_persisted", idx)
		}

		idx++
	}

	return nil
}

func (s *StoragePacker) putBucket(ctx context.Context, bucket *Bucket) error {
	defer metrics.MeasureSince([]string{"storage_packer", "put_bucket"}, time.Now())
	if bucket == nil {
		return fmt.Errorf("nil bucket entry")
	}

	if bucket.Key == "" {
		return fmt.Errorf("missing key")
	}

	if !strings.HasPrefix(bucket.Key, s.viewPrefix) {
		return fmt.Errorf("incorrect prefix; bucket entry key should have %q prefix", s.viewPrefix)
	}

	marshaledBucket, err := proto.Marshal(bucket)
	if err != nil {
		return fmt.Errorf("failed to marshal bucket: %w", err)
	}

	compressedBucket, err := compressutil.Compress(marshaledBucket, &compressutil.CompressionConfig{
		Type: compressutil.CompressionTypeSnappy,
	})
	if err != nil {
		return fmt.Errorf("failed to compress packed bucket: %w", err)
	}

	// Store the compressed value
	err = s.view.Put(ctx, &logical.StorageEntry{
		Key:   bucket.Key,
		Value: compressedBucket,
	})
	if err != nil {
		return fmt.Errorf("failed to persist packed storage entry: %w", err)
	}

	return nil
}

// GetItem fetches the storage entry for a given key from its corresponding
// bucket.
func (s *StoragePacker) GetItem(itemID string) (*Item, error) {
	defer metrics.MeasureSince([]string{"storage_packer", "get_item"}, time.Now())

	if itemID == "" {
		return nil, fmt.Errorf("empty item ID")
	}

	bucketKey := s.BucketKey(itemID)

	// Fetch the bucket entry
	bucket, err := s.GetBucket(context.Background(), bucketKey)
	if err != nil {
		return nil, fmt.Errorf("failed to read packed storage item: %w", err)
	}
	if bucket == nil {
		return nil, nil
	}

	// Look for a matching storage entry in the bucket items
	for _, item := range bucket.Items {
		if item.ID == itemID {
			return item, nil
		}
	}

	return nil, nil
}

// PutItem stores the given item in its respective bucket
func (s *StoragePacker) PutItem(ctx context.Context, item *Item) error {
	defer metrics.MeasureSince([]string{"storage_packer", "put_item"}, time.Now())

	if item == nil {
		return fmt.Errorf("nil item")
	}

	if item.ID == "" {
		return fmt.Errorf("missing ID in item")
	}

	var err error
	bucketKey := s.BucketKey(item.ID)

	bucket := &Bucket{
		Key: bucketKey,
	}

	// In this case, we persist the storage entry regardless of the read
	// storageEntry below is nil or not. Hence, directly acquire write lock
	// even to read the entry.
	lock := locksutil.LockForKey(s.storageLocks, bucketKey)
	lock.Lock()
	defer lock.Unlock()

	// Check if there is an existing bucket for a given key
	storageEntry, err := s.view.Get(ctx, bucketKey)
	if err != nil {
		return fmt.Errorf("failed to read packed storage bucket entry: %w", err)
	}

	if storageEntry == nil {
		// If the bucket entry does not exist, this will be the only item the
		// bucket that is going to be persisted.
		bucket.Items = []*Item{
			item,
		}
	} else {
		uncompressedData, notCompressed, err := compressutil.Decompress(storageEntry.Value)
		if err != nil {
			return fmt.Errorf("failed to decompress packed storage entry: %w", err)
		}
		if notCompressed {
			uncompressedData = storageEntry.Value
		}

		err = proto.Unmarshal(uncompressedData, bucket)
		if err != nil {
			return fmt.Errorf("failed to decode packed storage entry: %w", err)
		}

		err = bucket.upsert(item)
		if err != nil {
			return fmt.Errorf("failed to update entry in packed storage entry: %w", err)
		}
	}

	return s.putBucket(ctx, bucket)
}

// NewStoragePacker creates a new storage packer for a given view
func NewStoragePacker(view logical.Storage, logger log.Logger, viewPrefix string) (*StoragePacker, error) {
	if view == nil {
		return nil, fmt.Errorf("nil view")
	}

	if viewPrefix == "" {
		viewPrefix = StoragePackerBucketsPrefix
	}

	if !strings.HasSuffix(viewPrefix, "/") {
		viewPrefix = viewPrefix + "/"
	}

	// Create a new packer object for the given view
	packer := &StoragePacker{
		view:         view,
		viewPrefix:   viewPrefix,
		logger:       logger,
		storageLocks: locksutil.CreateLocks(),
	}

	return packer, nil
}
