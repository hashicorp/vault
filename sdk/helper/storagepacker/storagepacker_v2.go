package storagepacker

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"strings"
	"sync"
	"sync/atomic"

	radix "github.com/armon/go-radix"
	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/sdk/helper/compressutil"
	"github.com/hashicorp/vault/sdk/helper/cryptoutil"
	"github.com/hashicorp/vault/sdk/helper/locksutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/physical"
)

const (
	defaultBaseBucketBits  = 8
	defaultBucketShardBits = 4
)

var (
	shardLocks = make(map[*locksutil.LockEntry]struct{}, 32)
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
}

// StoragePacker packs many items into abstractions called buckets. The goal
// is to employ a reduced number of storage entries for a relatively huge
// number of items. This is the second version of the utility which supports
// indefinitely expanding the capacity of the storage by sharding the buckets
// when they exceed the imposed limit.
type StoragePackerV2 struct {
	*Config
	storageLocks []*locksutil.LockEntry

	bucketsCache     *radix.Tree
	bucketsCacheLock sync.RWMutex

	queueMode     uint32
	queuedBuckets sync.Map

	// disableSharding is used for tests
	disableSharding bool

	// Ensures we only process one sharding at a time, so that we can't grab
	// the same locks for the child buckets
	shardLock sync.RWMutex

	prewarmCache sync.Once
}

// LockedBucket embeds a bucket and its corresponding lock to ensure thread
// safety
type LockedBucket struct {
	*locksutil.LockEntry
	*Bucket
}

func (s *StoragePackerV2) BucketsView() *logical.StorageView {
	return s.BucketStorageView
}

func (s *StoragePackerV2) BucketStorageKeyForItemID(itemID string) string {
	hexVal := GetItemIDHash(itemID)

	s.bucketsCacheLock.RLock()
	_, bucketRaw, found := s.bucketsCache.LongestPrefix(hexVal)
	s.bucketsCacheLock.RUnlock()

	if found {
		return bucketRaw.(*LockedBucket).Key
	}

	// If we have existing buckets we'd have parsed them in on startup so this
	// is a fresh storagepacker, so we use the root bits to return a proper
	// number of chars.
	return hexVal[0 : s.BaseBucketBits/4]
}

func (s *StoragePackerV2) BucketKey(itemID string) string {
	return s.GetCacheKey(s.BucketStorageKeyForItemID(itemID))
}

func (s *StoragePackerV2) GetCacheKey(key string) string {
	return strings.Replace(key, "/", "", -1)
}

func GetItemIDHash(itemID string) string {
	return hex.EncodeToString(cryptoutil.Blake2b256Hash(itemID))
}

func (s *StoragePackerV2) BucketKeys(ctx context.Context) ([]string, error) {
	var retErr error
	s.prewarmCache.Do(func() {
		diskBuckets, err := logical.CollectKeys(ctx, s.BucketStorageView)
		if err != nil {
			retErr = err
			return
		}
		for _, key := range diskBuckets {
			// Read from the underlying view
			storageEntry, err := s.BucketStorageView.Get(ctx, key)
			if err != nil {
				retErr = errwrap.Wrapf("failed to read packed storage entry: {{err}}", err)
				return
			}
			if storageEntry == nil {
				retErr = fmt.Errorf("no data found at bucket %s", key)
				return
			}

			bucket, err := s.DecodeBucket(storageEntry)
			if err != nil {
				retErr = err
				return
			}

			s.bucketsCacheLock.Lock()
			s.bucketsCache.Insert(s.GetCacheKey(bucket.Key), bucket)
			s.bucketsCacheLock.Unlock()
		}
	})
	if retErr != nil {
		return nil, retErr
	}

	ret := make([]string, 0, 256)
	s.bucketsCacheLock.RLock()
	s.bucketsCache.Walk(func(s string, _ interface{}) bool {
		ret = append(ret, s)
		return false
	})
	s.bucketsCacheLock.RUnlock()

	return ret, nil
}

// Get returns a bucket for a given key
func (s *StoragePackerV2) GetBucket(ctx context.Context, key string, skipCache bool) (*LockedBucket, error) {
	cacheKey := s.GetCacheKey(key)

	if key == "" {
		return nil, fmt.Errorf("missing bucket key")
	}

	lock := locksutil.LockForKey(s.storageLocks, cacheKey)
	lock.RLock()

	s.bucketsCacheLock.RLock()
	_, bucketRaw, found := s.bucketsCache.LongestPrefix(cacheKey)
	s.bucketsCacheLock.RUnlock()

	if found && !skipCache {
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

	if found && !skipCache {
		ret := bucketRaw.(*LockedBucket)
		return ret, nil
	}

	// Read from the underlying view
	storageEntry, err := s.BucketStorageView.Get(ctx, key)
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
func (s *StoragePackerV2) DecodeBucket(storageEntry *logical.StorageEntry) (*LockedBucket, error) {
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

	cacheKey := s.GetCacheKey(storageEntry.Key)
	lock := locksutil.LockForKey(s.storageLocks, cacheKey)
	lb := &LockedBucket{
		LockEntry: lock,
		Bucket:    &bucket,
	}
	lb.Key = storageEntry.Key

	return lb, nil
}

// Put stores a packed bucket entry
func (s *StoragePackerV2) PutBucket(ctx context.Context, bucket *LockedBucket) error {
	if bucket == nil {
		return fmt.Errorf("nil bucket entry")
	}

	if bucket.Key == "" {
		return fmt.Errorf("missing key")
	}

	bucket.Lock()
	defer bucket.Unlock()

	if err := s.storeBucket(ctx, bucket, true); err != nil {
		return errwrap.Wrapf("failed at high level bucket put: {{err}}", err)
	}

	s.bucketsCacheLock.Lock()
	s.bucketsCache.Insert(s.GetCacheKey(bucket.Key), bucket)
	s.bucketsCacheLock.Unlock()

	return nil
}

func (s *StoragePackerV2) shardBucket(ctx context.Context, bucket *LockedBucket, cacheKey string, allowLocking bool) error {
	if allowLocking {
		s.shardLock.Lock()
		defer s.shardLock.Unlock()

		bucketLock := bucket.LockEntry
		defer func() {
			for lock := range shardLocks {
				// Let the initial calling function take care of the highest level lock
				if lock != bucketLock {
					lock.Unlock()
				}
			}

			// Empty the map
			shardLocks = make(map[*locksutil.LockEntry]struct{}, 32)
		}()
	}

	numShards := int(math.Pow(2.0, float64(s.BucketShardBits)))

	// Create the shards and lock them
	s.Logger.Info("sharding bucket", "bucket_key", bucket.Key, "num_shards", numShards)
	defer s.Logger.Info("sharding bucket process exited", "bucket_key", bucket.Key)

	shardLocks[bucket.LockEntry] = struct{}{}

	shards := make(map[string]*LockedBucket, numShards)
	for i := 0; i < numShards; i++ {
		shardKey := fmt.Sprintf("%x", i)
		lock := locksutil.LockForKey(s.storageLocks, cacheKey+shardKey)
		shardedBucket := &LockedBucket{
			LockEntry: lock,
			Bucket: &Bucket{
				Key:     fmt.Sprintf("%s/%s", bucket.Key, shardKey),
				ItemMap: make(map[string][]byte),
			},
		}
		shards[shardKey] = shardedBucket
		// If it was equal we'd be locked already
		s.Logger.Trace("created shard", "shard_key", shardKey)
		// Don't try to lock the same lock twice in case it hashes that way
		if _, ok := shardLocks[lock]; !ok {
			s.Logger.Trace("locking lock", "shard_key", shardKey)
			lock.Lock()
			shardLocks[lock] = struct{}{}
		}
	}

	s.Logger.Debug("resilvering items")

	parentPrefix := s.GetCacheKey(bucket.Key)
	// Resilver the items
	for k, v := range bucket.ItemMap {
		itemKey := strings.TrimPrefix(k, parentPrefix)[0 : s.BucketShardBits/4]
		s.Logger.Trace("resilvering item", "parent_prefix", parentPrefix, "item_id", k, "item_key", itemKey)
		// Sanity check
		childBucket, ok := shards[itemKey]
		if !ok {
			// We didn't complete sharding so don't make other parts of the
			// code think that it completed
			s.Logger.Error("failed to find sharded storagepacker bucket", "bucket_key", bucket.Key, "item_key", itemKey)
			return errors.New("failed to shard storagepacker bucket")
		}
		childBucket.ItemMap[k] = v
	}

	s.Logger.Debug("storing sharded buckets")

	// Ensure we can write all of these buckets. Create a cleanup function if not.
	retErr := new(multierror.Error)
	cleanupStorage := func() {
		for _, v := range shards {
			if err := s.BucketStorageView.Delete(ctx, v.Key); err != nil {
				retErr = multierror.Append(retErr, err)
				// Don't exit out, clean up as much as possible
			}
		}
	}
	for k, v := range shards {
		s.Logger.Trace("storing bucket", "shard", k)
		if err := s.storeBucket(ctx, v, false); err != nil {
			s.Logger.Debug("encountered error", "shard", k)
			retErr = multierror.Append(retErr, err)
			cleanupStorage()
			return retErr
		}
	}

	cleanupCache := func() {
		for _, v := range shards {
			s.bucketsCache.Delete(s.GetCacheKey(v.Key))
		}
	}
	// Add to the cache. It's not too late to back out, via the cleanup cache
	// function. We hold the lock while storing the updated original bucket so
	// that nobody accesses it in an inconsistent state.
	s.Logger.Debug("updating cache")
	s.bucketsCacheLock.Lock()
	{
		for _, v := range shards {
			s.bucketsCache.Insert(s.GetCacheKey(v.Key), v)
		}

		// Finally, update the original and persist
		origBucketItemMap := bucket.ItemMap
		bucket.ItemMap = nil
		if err := s.storeBucket(ctx, bucket, false); err != nil {
			retErr = multierror.Append(retErr, err)
			bucket.ItemMap = origBucketItemMap
			cleanupStorage()
			cleanupCache()
		}
	}

	s.bucketsCacheLock.Unlock()

	return retErr.ErrorOrNil()
}

// storeBucket actually stores the bucket. It expects that it's already locked.
func (s *StoragePackerV2) storeBucket(ctx context.Context, bucket *LockedBucket, allowLocking bool) error {
	if atomic.LoadUint32(&s.queueMode) == 1 {
		s.queuedBuckets.Store(bucket.Key, bucket)
		return nil
	}

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
	err = s.BucketStorageView.Put(ctx, &logical.StorageEntry{
		Key:   bucket.Key,
		Value: compressedBucket,
	})
	if err != nil {
		if strings.Contains(err.Error(), physical.ErrValueTooLarge) && !s.disableSharding {
			err = s.shardBucket(ctx, bucket, s.GetCacheKey(bucket.Key), allowLocking)
		}
		if err != nil {
			return errwrap.Wrapf("failed to persist packed storage entry: {{err}}", err)
		}
	}

	return nil
}

// DeleteBucket deletes an entire bucket entry
func (s *StoragePackerV2) DeleteBucket(ctx context.Context, key string) error {
	if key == "" {
		return fmt.Errorf("missing key")
	}

	cacheKey := s.GetCacheKey(key)

	lock := locksutil.LockForKey(s.storageLocks, cacheKey)
	lock.Lock()
	defer lock.Unlock()

	if err := s.BucketStorageView.Delete(ctx, key); err != nil {
		return errwrap.Wrapf("failed to delete packed storage entry: {{err}}", err)
	}

	s.bucketsCacheLock.Lock()
	s.bucketsCache.Delete(cacheKey)
	s.bucketsCacheLock.Unlock()

	return nil
}

// upsert either inserts a new item into the bucket or updates an existing one
// if an item with a matching key is already present.
func (s *LockedBucket) upsert(item *Item) error {
	if s == nil {
		return fmt.Errorf("nil storage bucket")
	}

	if item == nil {
		return fmt.Errorf("nil item")
	}

	if item.ID == "" {
		return fmt.Errorf("missing item ID")
	}

	if s.ItemMap == nil {
		s.ItemMap = make(map[string][]byte)
	}

	itemHash := GetItemIDHash(item.ID)

	s.ItemMap[itemHash] = item.Data
	return nil
}

// DeleteItem removes the storage entry which the given key refers to from its
// corresponding bucket.
func (s *StoragePackerV2) DeleteItem(ctx context.Context, itemID string) error {
	if itemID == "" {
		return fmt.Errorf("empty item ID")
	}

	// Get the bucket key
	bucketKey := s.BucketStorageKeyForItemID(itemID)
	cacheKey := s.GetCacheKey(bucketKey)

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
		storageEntry, err := s.BucketStorageView.Get(ctx, bucketKey)
		if err != nil {
			return errwrap.Wrapf("failed to read packed storage value: {{err}}", err)
		}
		if storageEntry == nil {
			return nil
		}

		bucket, err = s.DecodeBucket(storageEntry)
		if err != nil {
			return errwrap.Wrapf("error decoding existing storage entry for deletion: {{err}}", err)
		}

		bucket.LockEntry = lock

		s.bucketsCacheLock.Lock()
		s.bucketsCache.Insert(cacheKey, bucket)
		s.bucketsCacheLock.Unlock()
	}

	if len(bucket.ItemMap) == 0 {
		return nil
	}

	itemHash := GetItemIDHash(itemID)

	_, ok := bucket.ItemMap[itemHash]
	if !ok {
		return nil
	}

	delete(bucket.ItemMap, itemHash)
	return s.storeBucket(ctx, bucket, true)
}

// GetItem fetches the storage entry for a given key from its corresponding
// bucket.
func (s *StoragePackerV2) GetItem(ctx context.Context, itemID string) (*Item, error) {
	if itemID == "" {
		return nil, fmt.Errorf("empty item ID")
	}

	bucketKey := s.BucketStorageKeyForItemID(itemID)
	cacheKey := s.GetCacheKey(bucketKey)

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
		storageEntry, err := s.BucketStorageView.Get(ctx, bucketKey)
		if err != nil {
			return nil, errwrap.Wrapf("failed to read packed storage value: {{err}}", err)
		}
		if storageEntry == nil {
			return nil, nil
		}

		bucket, err = s.DecodeBucket(storageEntry)
		if err != nil {
			return nil, errwrap.Wrapf("error decoding existing storage entry: {{err}}", err)
		}

		bucket.LockEntry = lock

		s.bucketsCacheLock.Lock()
		s.bucketsCache.Insert(cacheKey, bucket)
		s.bucketsCacheLock.Unlock()
	}

	if len(bucket.ItemMap) == 0 {
		return nil, nil
	}

	itemHash := GetItemIDHash(itemID)

	data, ok := bucket.ItemMap[itemHash]
	if !ok {
		return nil, nil
	}

	return &Item{
		ID:   itemID,
		Data: data,
	}, nil
}

// PutItem stores a storage entry in its corresponding bucket
func (s *StoragePackerV2) PutItem(ctx context.Context, item *Item) error {
	if item == nil {
		return fmt.Errorf("nil item")
	}

	if item.ID == "" {
		return fmt.Errorf("missing ID in item")
	}

	if item.Data == nil {
		return fmt.Errorf("missing data in item")
	}

	if item.Message != nil {
		return fmt.Errorf("'Message' is deprecated; use 'Data' instead")
	}

	// Get the bucket key
	bucketKey := s.BucketStorageKeyForItemID(item.ID)
	cacheKey := s.GetCacheKey(bucketKey)

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
		storageEntry, err := s.BucketStorageView.Get(ctx, bucketKey)
		if err != nil {
			return errwrap.Wrapf("failed to read packed storage value: {{err}}", err)
		}

		if storageEntry == nil {
			bucket = &LockedBucket{
				LockEntry: lock,
				Bucket: &Bucket{
					Key: bucketKey,
				},
			}
		} else {
			bucket, err = s.DecodeBucket(storageEntry)
			if err != nil {
				return errwrap.Wrapf("error decoding existing storage entry for upsert: {{err}}", err)
			}

			bucket.LockEntry = lock
		}

		s.bucketsCacheLock.Lock()
		s.bucketsCache.Insert(cacheKey, bucket)
		s.bucketsCacheLock.Unlock()
	}

	if err := bucket.upsert(item); err != nil {
		return errwrap.Wrapf("failed to update entry in packed storage entry: {{err}}", err)
	}

	// Persist the result
	return s.storeBucket(ctx, bucket, true)
}

// NewStoragePackerV2 creates a new storage packer for a given view
func NewStoragePackerV2(ctx context.Context, config *Config) (StoragePacker, error) {
	if config.BucketStorageView == nil {
		return nil, fmt.Errorf("nil buckets view")
	}

	config.BucketStorageView = config.BucketStorageView.SubView("v2/")

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
		// The rest of the values can change
		config.BaseBucketBits = exist.BaseBucketBits
	}

	if config.BucketShardBits == 0 {
		config.BucketShardBits = defaultBucketShardBits
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
	packer := &StoragePackerV2{
		Config:       config,
		bucketsCache: radix.New(),
		storageLocks: locksutil.CreateLocks(),
	}

	// Prewarm the cache
	if _, err := packer.BucketKeys(ctx); err != nil {
		return nil, errwrap.Wrapf("error preloading storagepacker cache: {{err}}", err)
	}

	return packer, nil
}

func (s *StoragePackerV2) SetQueueMode(enabled bool) {
	if enabled {
		atomic.StoreUint32(&s.queueMode, 1)
	} else {
		atomic.StoreUint32(&s.queueMode, 0)
	}
}

func (s *StoragePackerV2) FlushQueue(ctx context.Context) error {
	var err *multierror.Error
	s.queuedBuckets.Range(func(key, value interface{}) bool {
		lErr := s.storeBucket(ctx, value.(*LockedBucket), true)
		if lErr != nil {
			err = multierror.Append(err, lErr)
		}
		s.queuedBuckets.Delete(key)
		return true
	})

	return err.ErrorOrNil()
}
