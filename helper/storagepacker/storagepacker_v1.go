package storagepacker

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"sync"

	radix "github.com/armon/go-radix"
	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/compressutil"
	"github.com/hashicorp/vault/helper/cryptoutil"
	"github.com/hashicorp/vault/helper/locksutil"
	"github.com/hashicorp/vault/logical"
)

const (
	defaultBaseBucketBits  = 8
	defaultBucketShardBits = 4
	// Larger size of the bucket size adversely affects the performance of the
	// storage packer. Also, some of the backends impose a maximum size limit
	// on the objects that gets persisted. For example, Consul imposes 256KB if using transactions
	// and DynamoDB imposes 400KB. Going forward, if there exists storage
	// backends that has more constrained limits, this will have to become more
	// flexible. For now, 240KB seems like a decent value.
	defaultBucketMaxSize = 240 * 1024
)

type Config struct {
	// BucketStorageView is the storage to be used by all the buckets
	BucketStorageView *logical.StorageView `json:"-"`

	// ConfigStorageView is the storage to store config info
	ConfigStorageView *logical.StorageView `json:"-"`

	// Logger for output
	Logger log.Logger `json:"-"`

	// BaseBucketBits is the number of bits to use for buckets at the base level
	BaseBucketBits int `json:"base_bucket_bits"`

	// BucketShardBits is the number of bits to use for sub-buckets a bucket
	// gets sharded into when it reaches the maximum threshold.
	BucketShardBits int `json:"-"`

	// BucketMaxSize (in bytes) is the maximum allowed size per bucket. When
	// the size of the bucket reaches a threshold relative to this limit, it
	// gets sharded into the configured number of pieces incrementally.
	BucketMaxSize int64 `json:"-"`
}

// StoragePacker packs many items into abstractions called buckets. The goal
// is to employ a reduced number of storage entries for a relatively huge
// number of items. This is the second version of the utility which supports
// indefinitely expanding the capacity of the storage by sharding the buckets
// when they exceed the imposed limit.
type StoragePackerV1 struct {
	*Config
	storageLocks []*locksutil.LockEntry
	bucketsCache *radix.Tree

	// Note that we're slightly loosy-goosy with this lock. The reason is that
	// outside of an identity store upgrade case, only PutItem will ever write
	// a bucket, and that will always fetch a lock on the bucket first. This
	// will also cover the sharding case since you'd get the parent lock first.
	// So we can get away with only locking just when modifying, because we
	// should already be locked in terms of an entry overwriting itself.
	bucketsCacheLock sync.RWMutex
}

// LockedBucket embeds a bucket and its corresponding lock to ensure thread
// safety
type LockedBucket struct {
	sync.RWMutex
	*Bucket
}

func (s *StoragePackerV1) BucketsView() *logical.StorageView {
	return s.BucketStorageView
}

func (s *StoragePackerV1) BucketStorageKeyForItemID(itemID string) string {
	hexVal := hex.EncodeToString(cryptoutil.Blake2b256Hash(itemID))

	s.bucketsCacheLock.RLock()
	_, bucketRaw, found := s.bucketsCache.LongestPrefix(hexVal)
	s.bucketsCacheLock.RUnlock()

	if found {
		return bucketRaw.(*LockedBucket).Key
	}

	// If we have existing buckets we'd have parsed them in on startup
	// (assuming that all users load all entries on startup), so this is a
	// fresh storagepacker, so we use the root bits to return a proper number
	// of chars. But first do that, lock, and try again to ensure nothing
	// changed without holding a lock.
	cacheKey := hexVal[0 : s.BaseBucketBits/4]
	lock := locksutil.LockForKey(s.storageLocks, cacheKey)
	lock.RLock()

	s.bucketsCacheLock.RLock()
	_, bucketRaw, found = s.bucketsCache.LongestPrefix(hexVal)
	s.bucketsCacheLock.RUnlock()

	lock.RUnlock()

	if found {
		return bucketRaw.(*LockedBucket).Key
	}

	return cacheKey
}

func (s *StoragePackerV1) BucketHashKeyForItemID(itemID string) string {
	return GetCacheKey(s.BucketStorageKeyForItemID(itemID))
}

func GetCacheKey(key string) string {
	return strings.Replace(key, "/", "", -1)
}

// Get returns a bucket for a given key
func (s *StoragePackerV1) GetBucket(key string) (*LockedBucket, error) {
	cacheKey := GetCacheKey(key)

	if key == "" {
		return nil, fmt.Errorf("missing bucket key")
	}

	lock := locksutil.LockForKey(s.storageLocks, cacheKey)
	lock.RLock()

	s.bucketsCacheLock.RLock()
	_, bucketRaw, found := s.bucketsCache.LongestPrefix(cacheKey)
	s.bucketsCacheLock.RUnlock()

	if found {
		ret := bucketRaw.(*LockedBucket)
		lock.RUnlock()
		return ret, nil
	}

	// Swap out for a write lock
	lock.RUnlock()
	lock.Lock()
	defer lock.Unlock()

	// Check for it to have been added
	s.bucketsCacheLock.RLock()
	_, bucketRaw, found = s.bucketsCache.LongestPrefix(cacheKey)
	s.bucketsCacheLock.RUnlock()

	if found {
		ret := bucketRaw.(*LockedBucket)
		return ret, nil
	}

	// Read from the underlying view
	storageEntry, err := s.BucketStorageView.Get(context.Background(), key)
	if err != nil {
		return nil, errwrap.Wrapf("failed to read packed storage entry: {{err}}", err)
	}
	if storageEntry == nil {
		return nil, nil
	}

	bucket, err := s.DecodeBucket(storageEntry)
	if err != nil {
		return nil, err
	}

	s.bucketsCacheLock.Lock()
	s.bucketsCache.Insert(cacheKey, bucket)
	s.bucketsCacheLock.Unlock()

	return bucket, nil
}

// NOTE: Don't put inserting into the cache here, as that will mess with
// upgrade cases for the identity store as we want to keep the bucket out of
// the cache until we actually re-store it.
func (s *StoragePackerV1) DecodeBucket(storageEntry *logical.StorageEntry) (*LockedBucket, error) {
	uncompressedData, notCompressed, err := compressutil.Decompress(storageEntry.Value)
	if err != nil {
		return nil, errwrap.Wrapf("failed to decompress packed storage entry: {{err}}", err)
	}
	if notCompressed {
		uncompressedData = storageEntry.Value
	}

	var bucket Bucket
	err = proto.Unmarshal(uncompressedData, &bucket)
	if err != nil {
		return nil, errwrap.Wrapf("failed to decode packed storage entry: {{err}}", err)
	}

	lb := &LockedBucket{
		Bucket: &bucket,
	}
	lb.Key = storageEntry.Key

	return lb, nil
}

// Put stores a packed bucket entry
func (s *StoragePackerV1) PutBucket(bucket *LockedBucket) error {
	if bucket == nil {
		return fmt.Errorf("nil bucket entry")
	}

	if bucket.Key == "" {
		return fmt.Errorf("missing key")
	}

	cacheKey := GetCacheKey(bucket.Key)

	lock := locksutil.LockForKey(s.storageLocks, cacheKey)
	lock.Lock()

	bucket.Lock()
	err := s.storeBucket(bucket)
	bucket.Unlock()

	lock.Unlock()

	return err
}

// storeBucket actually stores the bucket. It expects that it's already locked.
func (s *StoragePackerV1) storeBucket(bucket *LockedBucket) error {
	marshaledBucket, err := proto.Marshal(bucket.Bucket)
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
	err = s.BucketStorageView.Put(context.Background(), &logical.StorageEntry{
		Key:   bucket.Key,
		Value: compressedBucket,
	})
	if err != nil {
		return errwrap.Wrapf("failed to persist packed storage entry: {{err}}", err)
	}

	s.bucketsCacheLock.Lock()
	s.bucketsCache.Insert(GetCacheKey(bucket.Key), bucket)
	s.bucketsCacheLock.Unlock()

	return nil
}

// DeleteBucket deletes an entire bucket entry
func (s *StoragePackerV1) DeleteBucket(key string) error {
	if key == "" {
		return fmt.Errorf("missing key")
	}

	cacheKey := GetCacheKey(key)

	lock := locksutil.LockForKey(s.storageLocks, cacheKey)
	lock.Lock()
	defer lock.Unlock()

	if err := s.BucketStorageView.Delete(context.Background(), key); err != nil {
		return errwrap.Wrapf("failed to delete packed storage entry: {{err}}", err)
	}

	s.bucketsCacheLock.Lock()
	s.bucketsCache.Delete(cacheKey)
	s.bucketsCacheLock.Unlock()

	return nil
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

// DeleteItem removes the storage entry which the given key refers to from its
// corresponding bucket.
func (s *StoragePackerV1) DeleteItem(itemID string) error {
	if itemID == "" {
		return fmt.Errorf("empty item ID")
	}

	var err error

	// Get the bucket key
	bucketKey := s.BucketStorageKeyForItemID(itemID)
	cacheKey := GetCacheKey(bucketKey)

	lock := locksutil.LockForKey(s.storageLocks, cacheKey)
	lock.Lock()
	defer lock.Unlock()

	var bucket *LockedBucket

	s.bucketsCacheLock.RLock()
	_, bucketRaw, found := s.bucketsCache.LongestPrefix(cacheKey)
	s.bucketsCacheLock.RUnlock()

	if found {
		bucket = bucketRaw.(*LockedBucket)
	} else {
		// Read from underlying view
		storageEntry, err := s.BucketStorageView.Get(context.Background(), bucketKey)
		if err != nil {
			return errwrap.Wrapf("failed to read packed storage value: {{err}}", err)
		}
		if storageEntry == nil {
			return nil
		}

		bucket, err = s.DecodeBucket(storageEntry)
		if err != nil {
			return errwrap.Wrapf("error decoding existing storage entry for upsert: {{err}}", err)
		}

		s.bucketsCacheLock.Lock()
		s.bucketsCache.Insert(cacheKey, bucket)
		s.bucketsCacheLock.Unlock()
	}

	bucket.Lock()
	defer bucket.Unlock()

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
		err = s.storeBucket(bucket)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetItem fetches the storage entry for a given key from its corresponding
// bucket.
func (s *StoragePackerV1) GetItem(itemID string) (*Item, error) {
	if itemID == "" {
		return nil, fmt.Errorf("empty item ID")
	}

	bucketKey := s.BucketStorageKeyForItemID(itemID)
	cacheKey := GetCacheKey(bucketKey)

	lock := locksutil.LockForKey(s.storageLocks, cacheKey)
	lock.RLock()
	defer lock.RUnlock()

	var bucket *LockedBucket

	s.bucketsCacheLock.RLock()
	_, bucketRaw, found := s.bucketsCache.LongestPrefix(cacheKey)
	s.bucketsCacheLock.RUnlock()

	if found {
		bucket = bucketRaw.(*LockedBucket)
	} else {
		// Read from underlying view
		storageEntry, err := s.BucketStorageView.Get(context.Background(), bucketKey)
		if err != nil {
			return nil, errwrap.Wrapf("failed to read packed storage value: {{err}}", err)
		}
		if storageEntry == nil {
			return nil, nil
		}

		bucket, err = s.DecodeBucket(storageEntry)
		if err != nil {
			return nil, errwrap.Wrapf("error decoding existing storage entry for upsert: {{err}}", err)
		}

		s.bucketsCacheLock.Lock()
		s.bucketsCache.Insert(cacheKey, bucket)
		s.bucketsCacheLock.Unlock()
	}

	bucket.RLock()

	// Look for a matching storage entry in the bucket items
	for _, item := range bucket.Items {
		if item.ID == itemID {
			bucket.RUnlock()
			return item, nil
		}
	}

	bucket.RUnlock()
	return nil, nil
}

// PutItem stores a storage entry in its corresponding bucket
func (s *StoragePackerV1) PutItem(item *Item) error {
	if item == nil {
		return fmt.Errorf("nil item")
	}

	if item.ID == "" {
		return fmt.Errorf("missing ID in item")
	}

	// Get the bucket key
	bucketKey := s.BucketStorageKeyForItemID(item.ID)
	cacheKey := GetCacheKey(bucketKey)

	lock := locksutil.LockForKey(s.storageLocks, cacheKey)
	lock.Lock()
	defer lock.Unlock()

	var bucket *LockedBucket

	s.bucketsCacheLock.RLock()
	_, bucketRaw, found := s.bucketsCache.LongestPrefix(cacheKey)
	s.bucketsCacheLock.RUnlock()

	if found {
		bucket = bucketRaw.(*LockedBucket)
	} else {
		// Read from underlying view
		storageEntry, err := s.BucketStorageView.Get(context.Background(), bucketKey)
		if err != nil {
			return errwrap.Wrapf("failed to read packed storage value: {{err}}", err)
		}

		if storageEntry == nil {
			bucket = &LockedBucket{
				Bucket: &Bucket{
					Key: bucketKey,
				},
			}
		} else {
			bucket, err = s.DecodeBucket(storageEntry)
			if err != nil {
				return errwrap.Wrapf("error decoding existing storage entry for upsert: {{err}}", err)
			}

			s.bucketsCacheLock.Lock()
			s.bucketsCache.Insert(cacheKey, bucket)
			s.bucketsCacheLock.Unlock()
		}
	}

	bucket.Lock()
	defer bucket.Unlock()

	if err := bucket.upsert(item); err != nil {
		return errwrap.Wrapf("failed to update entry in packed storage entry: {{err}}", err)
	}

	// Persist the result
	return s.storeBucket(bucket)
}

// NewStoragePackerV1 creates a new storage packer for a given view
func NewStoragePackerV1(ctx context.Context, config *Config) (*StoragePackerV1, error) {
	if config.BucketStorageView == nil {
		return nil, fmt.Errorf("nil buckets view")
	}

	if config.ConfigStorageView == nil {
		return nil, fmt.Errorf("nil config view")
	}

	if config.BaseBucketBits == 0 {
		config.BaseBucketBits = defaultBaseBucketBits
	}

	// At this point, look for an existing saved configuration
	var needPersist bool
	entry, err := config.ConfigStorageView.Get(ctx, "config")
	if err != nil {
		return nil, errwrap.Wrapf("error checking for existing storagepacker config: {{err}}", err)
	}
	if entry != nil {
		needPersist = false
		var exist Config
		if err := entry.DecodeJSON(&exist); err != nil {
			return nil, errwrap.Wrapf("error decoding existing storagepacker config: {{err}}", err)
		}
		// If we have an existing config, we copy the only thing we need
		// constant: the bucket base count, so we know how many to expect at
		// the base level
		//
		// The rest of the values can change; the max size can change based on
		// e.g. if storage is migrated, so as long as we don't move to a new
		// location with a smaller value we're fine (and even then we're fine
		// if we can read it; otherwise storage migration would have failed
		// anyways). The shard count is recorded in each bucket at the time
		// it's sharded; if we realize it's more efficient to do some other
		// value later we can update it and use that going forward for new
		// shards.
		config.BaseBucketBits = exist.BaseBucketBits
	}

	if config.BucketShardBits == 0 {
		config.BucketShardBits = defaultBucketShardBits
	}

	if config.BucketMaxSize == 0 {
		config.BucketMaxSize = defaultBucketMaxSize
	}

	if config.BaseBucketBits%4 != 0 {
		return nil, fmt.Errorf("bucket base bits of %d is not a multiple of four", config.BaseBucketBits)
	}

	if config.BucketShardBits%4 != 0 {
		return nil, fmt.Errorf("bucket shard count of %d is not a power of two", config.BucketShardBits)
	}

	if config.BaseBucketBits < 4 {
		return nil, errors.New("bucket base bits should be at least 4")
	}
	if config.BucketShardBits < 4 {
		return nil, errors.New("bucket shard count should at least be 4")
	}

	if needPersist {
		entry, err := logical.StorageEntryJSON("config", config)
		if err != nil {
			return nil, errwrap.Wrapf("error encoding storagepacker config: {{err}}", err)
		}
		if err := config.ConfigStorageView.Put(ctx, entry); err != nil {
			return nil, errwrap.Wrapf("error storing storagepacker config: {{err}}", err)
		}
	}

	// Create a new packer object for the given view
	packer := &StoragePackerV1{
		Config:       config,
		bucketsCache: radix.New(),
		storageLocks: locksutil.CreateLocks(),
	}

	return packer, nil
}
