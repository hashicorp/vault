package storagepacker

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/compressutil"
	"github.com/hashicorp/vault/helper/locksutil"
	"github.com/hashicorp/vault/helper/storagepacker"
	sp2 "github.com/hashicorp/vault/helper/storagepacker"
	"github.com/hashicorp/vault/logical"
)

const (
	LegacyStoragePackerBucketsPrefix = "packer/buckets/"
)

// LegacyStoragePacker packs the objects into a specific number of buckets by hashing
// its ID and indexing it. Currently this supports only 256 bucket entries and
// hence relies on the first byte of the hash value for indexing. The items
// that gets inserted into the packer should implement StorageBucketItem
// interface.
type LegacyStoragePacker struct {
	view         *logical.StorageView
	logger       log.Logger
	storageLocks []*locksutil.LockEntry
	viewPrefix   string
}

// BucketPath returns the storage entry key for a given bucket key
func (s *LegacyStoragePacker) BucketPath(bucketKey string) string {
	return s.viewPrefix + bucketKey
}

// BucketKeyHash returns the MD5 hash of the bucket storage key in which
// the item will be stored. The choice of MD5 is only for hash performance
// reasons since its value is not used for any security sensitive operation.
func (s *LegacyStoragePacker) BucketKeyHashByItemID(itemID string) string {
	return s.BucketPath(s.BucketKey(itemID))
}

func (s *LegacyStoragePacker) GetCacheKey(key string) string {
	return key
}

// BucketKeyHashByKey returns the MD5 hash of the bucket storage key
func (s *LegacyStoragePacker) BucketKeyHashByKey(bucketKey string) string {
	hf := md5.New()
	hf.Write([]byte(bucketKey))
	return hex.EncodeToString(hf.Sum(nil))
}

// View returns the storage view configured to be used by the packer
func (s *LegacyStoragePacker) BucketsView() *logical.StorageView {
	return s.view
}

// Get returns a bucket for a given key
func (s *LegacyStoragePacker) GetBucket(ctx context.Context, key string) (*sp2.LockedBucket, error) {
	if key == "" {
		return nil, fmt.Errorf("missing bucket key")
	}

	lock := locksutil.LockForKey(s.storageLocks, key)
	lock.RLock()
	defer lock.RUnlock()

	// Read from the underlying view
	storageEntry, err := s.view.Get(ctx, key)
	if err != nil {
		return nil, errwrap.Wrapf("failed to read packed storage entry: {{err}}", err)
	}

	lb, err := s.DecodeBucket(storageEntry)
	if err != nil {
		return nil, err
	}

	return lb, nil
}

func (s *LegacyStoragePacker) DecodeBucket(storageEntry *logical.StorageEntry) (*sp2.LockedBucket, error) {
	if storageEntry == nil {
		return nil, nil
	}

	uncompressedData, notCompressed, err := compressutil.Decompress(storageEntry.Value)
	if err != nil {
		return nil, errwrap.Wrapf("failed to decompress packed storage entry: {{err}}", err)
	}
	if notCompressed {
		uncompressedData = storageEntry.Value
	}

	var bucket sp2.Bucket
	err = proto.Unmarshal(uncompressedData, &bucket)
	if err != nil {
		return nil, errwrap.Wrapf("failed to decode packed storage entry: {{err}}", err)
	}

	return &sp2.LockedBucket{Bucket: &bucket}, nil
}

func (s *LegacyStoragePacker) DeleteBucket(ctx context.Context, key string) error {
	if key == "" {
		return fmt.Errorf("missing bucket key")
	}

	lock := locksutil.LockForKey(s.storageLocks, key)
	lock.Lock()
	defer lock.Unlock()

	return s.view.Delete(ctx, key)
}

// upsert either inserts a new item into the bucket or updates an existing one
// if an item with a matching key is already present.
func legacyUpsert(s *sp2.Bucket, item *sp2.Item) error {
	if s == nil {
		return fmt.Errorf("nil storage bucket")
	}

	if item == nil {
		return fmt.Errorf("nil item")
	}

	if item.ID == "" {
		return fmt.Errorf("missing item ID")
	}

	// Look for an item with matching key and don't modify the collection
	// while iterating
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

// BucketIndex returns the bucket key index for a given storage key
func (s *LegacyStoragePacker) BucketIndex(key string) uint8 {
	hf := md5.New()
	hf.Write([]byte(key))
	return uint8(hf.Sum(nil)[0])
}

// BucketKey returns the bucket key for a given item ID
func (s *LegacyStoragePacker) BucketKey(itemID string) string {
	return strconv.Itoa(int(s.BucketIndex(itemID)))
}

// DeleteItem removes the storage entry which the given key refers to from its
// corresponding bucket.
func (s *LegacyStoragePacker) DeleteItem(ctx context.Context, itemID string) error {

	if itemID == "" {
		return fmt.Errorf("empty item ID")
	}

	// Get the bucket key
	bucketKey := s.BucketKey(itemID)

	// Prepend the view prefix
	bucketPath := s.BucketPath(bucketKey)

	// Read from underlying view
	storageEntry, err := s.view.Get(ctx, bucketPath)
	if err != nil {
		return errwrap.Wrapf("failed to read packed storage value: {{err}}", err)
	}
	if storageEntry == nil {
		return nil
	}

	uncompressedData, notCompressed, err := compressutil.Decompress(storageEntry.Value)
	if err != nil {
		return errwrap.Wrapf("failed to decompress packed storage value: {{err}}", err)
	}
	if notCompressed {
		uncompressedData = storageEntry.Value
	}

	var bucket sp2.Bucket
	err = proto.Unmarshal(uncompressedData, &bucket)
	if err != nil {
		return errwrap.Wrapf("failed decoding packed storage entry: {{err}}", err)
	}

	// Look for a matching storage entry
	foundIdx := -1
	for itemIdx, item := range bucket.Items {
		if item.ID == itemID {
			foundIdx = itemIdx
			break
		}
	}

	// If there is a match, remove it from the collection and persist the
	// resulting collection
	if foundIdx != -1 {
		bucket.Items = append(bucket.Items[:foundIdx], bucket.Items[foundIdx+1:]...)

		// Persist bucket entry only if there is an update
		err = s.PutBucket(ctx, &sp2.LockedBucket{Bucket: &bucket})
		if err != nil {
			return err
		}
	}

	return nil
}

// Put stores a packed bucket entry
func (s *LegacyStoragePacker) PutBucket(ctx context.Context, bucket *sp2.LockedBucket) error {
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
		return errwrap.Wrapf("failed to marshal bucket: {{err}}", err)
	}

	compressedBucket, err := compressutil.Compress(marshaledBucket, &compressutil.CompressionConfig{
		Type: compressutil.CompressionTypeSnappy,
	})
	if err != nil {
		return errwrap.Wrapf("failed to compress packed bucket: {{err}}", err)
	}

	// Store the compressed value
	err = s.view.Put(ctx, &logical.StorageEntry{
		Key:   bucket.Key,
		Value: compressedBucket,
	})
	if err != nil {
		return errwrap.Wrapf("failed to persist packed storage entry: {{err}}", err)
	}

	return nil
}

// GetItem fetches the storage entry for a given key from its corresponding
// bucket.
func (s *LegacyStoragePacker) GetItem(ctx context.Context, itemID string) (*sp2.Item, error) {
	if itemID == "" {
		return nil, fmt.Errorf("empty item ID")
	}

	bucketKey := s.BucketKey(itemID)
	bucketPath := s.BucketPath(bucketKey)

	// Fetch the bucket entry
	bucket, err := s.GetBucket(ctx, bucketPath)
	if err != nil {
		return nil, errwrap.Wrapf("failed to read packed storage item: {{err}}", err)
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

// PutItem stores a storage entry in its corresponding bucket
func (s *LegacyStoragePacker) PutItem(ctx context.Context, item *sp2.Item) error {
	if item == nil {
		return fmt.Errorf("nil item")
	}

	if item.ID == "" {
		return fmt.Errorf("missing ID in item")
	}

	var err error
	bucketKey := s.BucketKey(item.ID)
	bucketPath := s.BucketPath(bucketKey)

	bucket := &sp2.Bucket{
		Key: bucketPath,
	}

	// In this case, we persist the storage entry regardless of the read
	// storageEntry below is nil or not. Hence, directly acquire write lock
	// even to read the entry.
	lock := locksutil.LockForKey(s.storageLocks, bucketPath)
	lock.Lock()
	defer lock.Unlock()

	// Check if there is an existing bucket for a given key
	storageEntry, err := s.view.Get(context.Background(), bucketPath)
	if err != nil {
		return errwrap.Wrapf("failed to read packed storage bucket entry: {{err}}", err)
	}

	if storageEntry == nil {
		// If the bucket entry does not exist, this will be the only item the
		// bucket that is going to be persisted.
		bucket.Items = []*sp2.Item{
			item,
		}
	} else {
		uncompressedData, notCompressed, err := compressutil.Decompress(storageEntry.Value)
		if err != nil {
			return errwrap.Wrapf("failed to decompress packed storage entry: {{err}}", err)
		}
		if notCompressed {
			uncompressedData = storageEntry.Value
		}

		err = proto.Unmarshal(uncompressedData, bucket)
		if err != nil {
			return errwrap.Wrapf("failed to decode packed storage entry: {{err}}", err)
		}

		err = legacyUpsert(bucket, item)
		if err != nil {
			return errwrap.Wrapf("failed to update entry in packed storage entry: {{err}}", err)
		}
	}

	// Persist the result
	return s.PutBucket(ctx, &sp2.LockedBucket{Bucket: bucket})
}

// NewLegacyStoragePacker creates a new storage packer for a given view
func NewLegacyStoragePacker(ctx context.Context, config *storagepacker.Config) (storagepacker.StoragePacker, error) {
	if config.BucketStorageView == nil {
		return nil, fmt.Errorf("nil view")
	}

	// Create a new packer object for the given view
	packer := &LegacyStoragePacker{
		view:         config.BucketStorageView,
		logger:       config.Logger,
		storageLocks: locksutil.CreateLocks(),
	}

	return packer, nil
}

func (s *LegacyStoragePacker) SetQueueMode(bool)                {}
func (s *LegacyStoragePacker) FlushQueue(context.Context) error { return nil }
