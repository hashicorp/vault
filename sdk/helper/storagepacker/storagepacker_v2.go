// Implements a storage layer that uses a radix tree and multiple buckets
// to provide arbitrarily-scalable storage on top of fixed-maximum-size objects.
package storagepacker

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"
	"sync"
	"sync/atomic"

	radix "github.com/armon/go-radix"
	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/sdk/helper/compressutil"
	"github.com/hashicorp/vault/sdk/helper/locksutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/physical"
)

const (
	defaultBaseBucketBits  = 8
	defaultBucketShardBits = 4
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
//
// The locking discipline for StoragePackerV2:
// * Acquire locks in the order storageLocks < bucketsCacheLock to avoid
//   deadlock.
// * Hold the storageLock for the duration of any read, write, or reshard operation on a given
//   bucket. The storageLock is determined by the "cache key" of the bucket, which
//   omits all of the '/'s.
// * Acquire locks for multiple buckets in the order defined by locksutil.
type StoragePackerV2 struct {
	*Config
	storageLocks []*locksutil.LockEntry

	// The bucketsCache stores LockedBucket objects.
	// For simplicity, the radix tree should cover the entire
	// range of BaseBucketBits.
	bucketsCache     *radix.Tree
	bucketsCacheLock sync.RWMutex

	queueMode     uint32
	queuedBuckets sync.Map

	// disableSharding is used for tests
	disableSharding bool

	prewarmCache sync.Once
}

type LockedBucket struct {
	*Bucket
	*locksutil.LockEntry
}

func (s *StoragePackerV2) preloadBucketsFromDisk(ctx context.Context, keys []string) error {
	// If we crash while writing out new buckets created by sharding,
	// we will have version N of the bucket in storage, and version N+1
	// represented in the sub-buckets.
	//
	// To bring these two into alignment, we can roll back to version N.
	// It might be possible to foll forward in some cases (particularly if
	// only one item can be added at a time?) but the crash must have
	// occurred before PutBucket() returned, so the new data need not be saved.
	//
	// If we load a bucket that has sub-buckets but is non-empty,
	// delete all of the sub-buckets.  (The current implementation recurses
	// so there may be more than one level of sub-buckets created by shardBucket.)

	// Visit the buckets using inorder traversal, parents before children.
	sort.Strings(keys)

	nonemptyParent := "NOT_A_PREFIX"

	s.bucketsCacheLock.Lock()
	defer s.bucketsCacheLock.Unlock()
	for _, key := range keys {
		if strings.HasPrefix(key, nonemptyParent) {
			s.Logger.Warn("detected shadowed bucket, removing", "key", key, "parent", nonemptyParent)
			s.BucketStorageView.Delete(ctx, key)
			continue
		}

		// Read from the underlying view
		storageEntry, err := s.BucketStorageView.Get(ctx, key)
		if err != nil {
			return errwrap.Wrapf("failed to read packed storage entry: {{err}}", err)
		}
		if storageEntry == nil {
			return fmt.Errorf("no data found at bucket %s", key)
		}
		bucket, err := s.DecodeBucket(storageEntry)
		if err != nil {
			return err
		}
		if !bucket.HasShards {
			nonemptyParent = key
		}

		s.bucketsCache.Insert(s.GetCacheKey(bucket.Key), bucket)
	}

	s.addAllBaseBucketsToCache()
	return nil
}

// Ensure all base buckets are present on the cache, even if they
// are not present in storage.  This simplifies acquiring locks for
// multiple buckets because every necessary bucket is already present
// in bucketsCache. We don't have to worry about adding them on-demand
// in the middle of that process, and if any object is added to these
// placeholders they'll get saved to stoarge.
func (s *StoragePackerV2) addAllBaseBucketsToCache() {
	for _, bucketKey := range s.getAllBaseBucketKeys() {
		if _, present := s.bucketsCache.Get(bucketKey); !present {
			lock := locksutil.LockForKey(s.storageLocks, bucketKey)
			s.bucketsCache.Insert(bucketKey,
				&LockedBucket{
					LockEntry: lock,
					Bucket: &Bucket{
						Key:       bucketKey,
						ItemMap:   make(map[string][]byte),
						HasShards: false,
					},
				})
		}
	}
}

func (s *StoragePackerV2) BucketKeys(ctx context.Context) ([]string, error) {
	var retErr error
	s.prewarmCache.Do(func() {
		diskBuckets, err := logical.CollectKeys(ctx, s.BucketStorageView)
		if err != nil {
			retErr = err
			return
		}
		retErr = s.preloadBucketsFromDisk(ctx, diskBuckets)
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

// DecodeBucket parses a bucket from storage, but doesn't add it to the cache.
// (We may want to read a v1 bucket and upgrade it?)
func (s *StoragePackerV2) DecodeBucket(storageEntry *logical.StorageEntry) (*LockedBucket, error) {
	if storageEntry == nil || storageEntry.Value == nil {
		return nil, errors.New("decoding nil storageEntry")
	}

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
	if lb.Key == "" {
		lb.Key = storageEntry.Key
	}

	return lb, nil
}

// GetBucket provides raw access to one of the storage buckets.
// The caller is responsible for acquiring the lock and verifying that the bucket
// has not been sharded in the meantime.
func (s *StoragePackerV2) GetBucket(ctx context.Context, bucketKey string) (*LockedBucket, error) {
	cacheKey := s.GetCacheKey(bucketKey)

	s.bucketsCacheLock.RLock()
	_, bucketRaw, found := s.bucketsCache.LongestPrefix(cacheKey)
	s.bucketsCacheLock.RUnlock()

	if found {
		return bucketRaw.(*LockedBucket), nil
	} else {
		return nil, fmt.Errorf("key %q not found in cache", cacheKey)
	}

}

// shardBucket splits an overly-large bucket into multiple children.
func (s *StoragePackerV2) shardBucket(ctx context.Context, bucket *LockedBucket, cacheKey string) error {
	numShards := int(math.Pow(2.0, float64(s.BucketShardBits)))

	if len(cacheKey)+s.BucketShardBits/4 > KeyLength {
		// Maximum depth of the radix tree exceeded --- extremely unlikely to happen
		// naturally because of the difficulty creating enough collisions, but a
		// wierd configuration (very large BucketShardBits) could potentially cause this to occur.
		return errors.New("attempting to shard past the end of the key")
	}

	// Create the shards
	s.Logger.Info("sharding bucket", "bucket_key", bucket.Key, "num_shards", numShards)
	defer s.Logger.Info("sharding bucket process exited", "bucket_key", bucket.Key)

	shards := make(map[string]*LockedBucket, numShards)
	for i := 0; i < numShards; i++ {
		shardKey := fmt.Sprintf("%x", i)
		lock := locksutil.LockForKey(s.storageLocks, cacheKey+shardKey)
		shardedBucket := &LockedBucket{
			LockEntry: lock,
			Bucket: &Bucket{
				Key:       fmt.Sprintf("%s/%s", bucket.Key, shardKey),
				ItemMap:   make(map[string][]byte),
				HasShards: false,
			},
		}
		shards[shardKey] = shardedBucket
		s.Logger.Trace("created shard", "shard_key", shardKey)
	}

	s.Logger.Debug("resilvering items")

	parentPrefix := s.GetCacheKey(bucket.Key)
	// Resilver the items
	for id, v := range bucket.ItemMap {
		key := GetItemIDHash(id)
		shardKey := strings.TrimPrefix(key, parentPrefix)[0 : s.BucketShardBits/4]
		s.Logger.Trace("resilvering item", "parent_prefix", parentPrefix, "item_id", id, "shard", shardKey)
		// Sanity check
		childBucket, ok := shards[shardKey]
		if !ok {
			// We didn't complete sharding so don't make other parts of the
			// code think that it completed
			s.Logger.Error("failed to find sharded storagepacker bucket", "bucket_key", bucket.Key, "shard", shardKey)
			return errors.New("failed to shard storagepacker bucket")
		}
		childBucket.ItemMap[id] = v
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
		// Avoid queuing mode, go directly to "store" not "persist"
		if err := s.storeBucket(ctx, v, 1); err != nil {
			s.Logger.Debug("encountered error", "shard", k)
			retErr = multierror.Append(retErr, err)
			cleanupStorage()
			return retErr
		}
	}

	// Update the original and persist it.
	// Prior to this point, version N is in storage and we revert back to it.
	// After this point, storage is in version N+1 but the cache has not yet
	// been updated; anybody trying to access one of the keys will hit
	// this bucket and block on the lock.
	origBucketItemMap := bucket.ItemMap
	bucket.ItemMap = nil
	bucket.HasShards = true
	if err := s.storeBucket(ctx, bucket, 1); err != nil {
		retErr = multierror.Append(retErr, err)
		bucket.ItemMap = origBucketItemMap
		bucket.HasShards = false
		cleanupStorage()
		return retErr
	}

	// Add to the cache.
	s.Logger.Debug("updating cache")
	s.bucketsCacheLock.Lock()
	for _, v := range shards {
		s.bucketsCache.Insert(s.GetCacheKey(v.Key), v)
	}
	s.bucketsCacheLock.Unlock()

	return retErr.ErrorOrNil()
}

// persistBuckets puts a bucket into persistent storage, or at least enqueues it.
// Expect that the bucket is already locked at this point.
func (s *StoragePackerV2) persistBucket(ctx context.Context, bucket *LockedBucket) error {
	if atomic.LoadUint32(&s.queueMode) == 1 {
		s.queuedBuckets.Store(bucket.Key, bucket)
		return nil
	}
	return s.storeBucket(ctx, bucket, 0)
}

// storeBucket actually stores the bucket.
// "depth" is used to prevent recursive operation, if one of the shards
// is itself too big, we fail and report an error.
func (s *StoragePackerV2) storeBucket(ctx context.Context, bucket *LockedBucket, depth int) error {
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
		if strings.Contains(err.Error(), physical.ErrValueTooLarge) {
			if depth > 0 {
				return errwrap.Wrapf("recursive sharding detected: {{err}}", err)
			} else if !s.disableSharding {
				err = s.shardBucket(ctx, bucket, s.GetCacheKey(bucket.Key))
			}
		}
		if err != nil {
			return errwrap.Wrapf("failed to persist packed storage entry: {{err}}", err)
		}
	}

	return nil
}

// upsertItems either inserts a new item into the bucket or updates an existing one
// if an item with a matching key is already present.
func (p *partitionedRequests) upsertItems() error {
	if p == nil {
		return fmt.Errorf("nil partition")
	}
	bucket := p.Bucket
	if bucket == nil {
		return fmt.Errorf("nil bucket")
	}

	if bucket.HasShards {
		return fmt.Errorf("upserting item into sharded bucket")
	}
	if bucket.ItemMap == nil {
		bucket.ItemMap = make(map[string][]byte)
	}

	for _, r := range p.Requests {
		// Nil data checked earlier
		bucket.ItemMap[r.ID] = r.Value.Data
	}
	return nil
}

// PutItem sets the value of multiple items, as an efficient but non-atomic operation.
// Each affected bucket will be written just once.
// A multierror will be returned if some of the writes fail.
func (s *StoragePackerV2) PutItem(ctx context.Context, items ...*Item) error {
	idsSeen := make(map[string]bool, len(items))
	for idx, item := range items {
		if item == nil {
			return fmt.Errorf("nil item at index %v", idx)
		}
		if item.ID == "" {
			return fmt.Errorf("missing ID in item at index %v", idx)
		}
		if item.Message != nil {
			return fmt.Errorf("'Message' is deprecated; use 'Data' instead")
		}
		if item.Data == nil {
			return fmt.Errorf("missing data in item at index %v", idx)
		}
		if _, found := idsSeen[item.ID]; found {
			return fmt.Errorf("duplicate ID %v in item at index %v", item.ID, idx)
		}
		idsSeen[item.ID] = true
	}

	// Data flow:
	//     (id, data)+
	// 1. compute hashes
	//     (id, key=hash(id), data)+
	// 2. sort and partition by bucket
	//     bucket NN => (id, key=NN*, data)+
	// 3. acquire storage locks in storage lock order
	// 4. process items in each bucket

	requests := s.keysForItems(items)

	// Identify the buckets and acquire their corresponding write-locks
	retryRequired := true
	var partition []*partitionedRequests
	for retryRequired {
		var err error
		partition, err = s.partitionRequests(requests)
		if err != nil {
			return err
		}
		retryRequired = s.lockBuckets(partition, false)
	}
	defer s.unlockBuckets(partition, false)

	// Update all buckets first
	for _, p := range partition {
		if err := p.upsertItems(); err != nil {
			return errwrap.Wrapf("failed to update storage bucket: {{err}}", err)
		}
	}

	// Persist the result (partial application is OK)
	var merr *multierror.Error
	for _, p := range partition {
		if err := s.persistBucket(ctx, p.Bucket); err != nil {
			merr = multierror.Append(merr, err)
		}
	}
	return merr.ErrorOrNil()
}

// DeleteItem removes the items identified by the set of IDs from storage.
// Deletion of a non-existent ID is not signalled as an error.
func (s *StoragePackerV2) DeleteItem(ctx context.Context, itemIDs ...string) error {
	requests := s.keysForIDs(itemIDs)

	// Identify the buckets and acquire their corresponding write-locks
	retryRequired := true
	var partition []*partitionedRequests
	for retryRequired {
		var err error
		partition, err = s.partitionRequests(requests)
		if err != nil {
			return err
		}
		retryRequired = s.lockBuckets(partition, false)
	}
	defer s.unlockBuckets(partition, false)

	// Update all buckets first
	for _, p := range partition {
		if p.Bucket.ItemMap == nil {
			continue
		}
		for _, k := range p.Requests {
			delete(p.Bucket.ItemMap, k.ID)
		}
	}

	// Persist the result (partial application is OK)
	var merr *multierror.Error
	for _, p := range partition {
		if err := s.persistBucket(ctx, p.Bucket); err != nil {
			merr = multierror.Append(merr, err)
		}
	}
	return merr.ErrorOrNil()
}

// GetsItems retrieves a set of items identified by ID.
// Missing items signal by nil's in the returned slice so that the the request
// and response lengths are the same.
func (s *StoragePackerV2) GetItems(ctx context.Context, ids ...string) ([]*Item, error) {
	requests := s.keysForIDs(ids)

	// Identify the buckets and acquire their corresponding read-locks
	// If we wanted to increase parallelism, perhaps could lock just one bucket
	// at a time, but for now it's simpler to just follow the model of
	// PutItem.
	retryRequired := true
	var partition []*partitionedRequests
	for retryRequired {
		var err error
		partition, err = s.partitionRequests(requests)
		if err != nil {
			return nil, err
		}
		retryRequired = s.lockBuckets(partition, true)
	}
	defer s.unlockBuckets(partition, true)

	// Walk each bucket and look for the corresponding keys.
	// The Value in each request is initialized to nil so no update
	// is needed if not present.
	for _, p := range partition {
		if p.Bucket.ItemMap == nil {
			continue
		}
		for _, req := range p.Requests {
			data, ok := p.Bucket.ItemMap[req.ID]
			if ok {
				req.Value = &Item{
					ID:   req.ID,
					Data: data,
				}
			}
		}
	}

	// Copy just the values to the output
	items := make([]*Item, len(requests))
	for i, req := range requests {
		items[i] = req.Value
	}
	return items, nil
}

// GetItem fetches a single item by ID.
func (s *StoragePackerV2) GetItem(ctx context.Context, itemID string) (*Item, error) {
	if itemID == "" {
		return nil, fmt.Errorf("empty item ID")
	}

	singleItem, err := s.GetItems(ctx, itemID)
	if err != nil {
		return nil, err
	}
	return singleItem[0], nil
}

// AllItems retrieve the entire set of items as a single slice.
// This *is* an atomic operation so that it provides a point-in-time
// view of the state.
func (s *StoragePackerV2) AllItems(context.Context) ([]*Item, error) {
	// Acquire all the locks ahead of time (and in order)
	for _, l := range s.storageLocks {
		l.RLock()
	}
	defer func() {
		for _, l := range s.storageLocks {
			l.RUnlock()
		}
	}()

	items := make([]*Item, 0)

	s.bucketsCacheLock.RLock()
	defer s.bucketsCacheLock.RUnlock()

	s.bucketsCache.Walk(func(s string, v interface{}) bool {
		bucket := v.(*LockedBucket)
		if bucket.HasShards || bucket.ItemMap == nil {
			return false
		}
		for id, data := range bucket.ItemMap {
			items = append(items, &Item{
				ID:   id,
				Data: data,
			})
		}
		return false
	})

	return items, nil
}

// NewStoragePackerV2 creates a new storage packer for a given view
func NewStoragePackerV2(ctx context.Context, config *Config) (StoragePacker, error) {
	if config.BucketStorageView == nil {
		return nil, fmt.Errorf("nil buckets view")
	}

	// Should we check if the view's prefix ends with a /?
	config.BucketStorageView = config.BucketStorageView.SubView("v2/")

	if config.ConfigStorageView == nil {
		return nil, fmt.Errorf("nil config view")
	}

	if config.Logger == nil {
		return nil, fmt.Errorf("nil logger")
	}

	if config.BaseBucketBits == 0 {
		config.BaseBucketBits = defaultBaseBucketBits
	}

	// At this point, look for an existing saved configuration
	needPersist := true
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
		config.BaseBucketBits = exist.BaseBucketBits

		// The shard count can change, but it might be confusing if we
		// do so.
		if config.BucketShardBits != exist.BucketShardBits {
			config.Logger.Info("BucketShardBits changed to %d from its original value of %d",
				config.BucketShardBits, exist.BucketShardBits)
		}
	}

	if config.BucketShardBits == 0 {
		config.BucketShardBits = defaultBucketShardBits
	}

	if config.BaseBucketBits%4 != 0 {
		return nil, fmt.Errorf("bucket base bits of %d is not a multiple of 4", config.BaseBucketBits)
	}

	if config.BucketShardBits%4 != 0 {
		return nil, fmt.Errorf("bucket shard bits of %d is not a multiple of 4", config.BucketShardBits)
	}

	if config.BaseBucketBits < 4 {
		return nil, errors.New("bucket base bits should be at least 4")
	}
	if config.BucketShardBits < 4 {
		return nil, errors.New("bucket shard count should be at least 4")
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

// SetQueueMode enables or disabled deferred writes of buckets to storage.
// The client must use FlushQueue to commit the accumulated writes.
func (s *StoragePackerV2) SetQueueMode(enabled bool) {
	if enabled {
		atomic.StoreUint32(&s.queueMode, 1)
	} else {
		atomic.StoreUint32(&s.queueMode, 0)
	}
}

// FlushQueue writes out accumulated buckets to storage.
// This mechanism doesn't guarantee a consistent ordering of the writes,
// so if sharding has occurred the resulting storage could be
// incorrectly recovered.
func (s *StoragePackerV2) FlushQueue(ctx context.Context) error {
	var err *multierror.Error
	s.queuedBuckets.Range(func(key, value interface{}) bool {
		bucket := value.(*LockedBucket)
		lErr := s.storeBucket(ctx, bucket, 0)
		if lErr != nil {
			err = multierror.Append(err, lErr)
		}
		s.queuedBuckets.Delete(key)
		return true
	})

	return err.ErrorOrNil()
}
