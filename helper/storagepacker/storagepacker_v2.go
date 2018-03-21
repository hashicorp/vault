package storagepacker

import (
	"context"
	"fmt"
	"math"
	"strings"
	"sync"

	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/helper/strutil"

	"github.com/hashicorp/vault/helper/cryptoutil"
	"github.com/hashicorp/vault/logical"
	log "github.com/mgutz/logxi/v1"
)

const (
	defaultBucketCount      = 256
	defaultBucketShardCount = 32
	defaultBucketMaxSize    = 256 * 1024
)

const (
	shared    = false
	exclusive = true
)

type Config struct {
	// View is the storage to be used by all the buckets
	View logical.Storage

	// ViewPrefix is the prefix to be used for the bucket keys in the view
	ViewPrefix string

	// Logger for output
	Logger log.Logger

	// BucketCount is the number of buckets to create at the base level
	BucketCount int

	// BucketShardCount is the number of buckets any given bucket can get
	// sharded into when it exceeds the maximum allowed size
	BucketShardCount int

	// BucketMaxSize (in bytes) is the maximum allowed size per bucket. When
	// the size of the bucket exceeds this limit, it gets sharded into the
	// configured number of pieces incrementally.
	BucketMaxSize int64
}

// StoragePackerV2 packs many items into abstractions called buckets. The goal
// is to employ a reduced number of storage entries for a relatively huge
// number of items. This is the second version of the utility which supports
// indefinitely expanding the capacity of the storage by sharding the buckets
// when they exceed the imposed limit.
type StoragePackerV2 struct {
	config           *Config
	bucketLocksCache *sync.Map
}

// Clone creates a replica of the bucket
func (b *BucketV2) Clone() (*BucketV2, error) {
	if b == nil {
		return nil, fmt.Errorf("nil bucket")
	}

	marshaledBucket, err := proto.Marshal(b)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal bucket: %v", err)
	}

	var clonedBucket BucketV2
	err = proto.Unmarshal(marshaledBucket, &clonedBucket)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal bucket: %v", err)
	}

	return &clonedBucket, nil
}

// putItemIntoBucket is a recursive function that finds the appropriate bucket
// to store the item based on the storage space available in the buckets. This
// method avoids race conditions by acquiring locks on each bucket key. Note
// that this function is optimized for a write-heavy workload. Locks are
// directly held for writing. Acquiring the lock for reading first and then
// switching the locks for writing provides less value as the readers are
// scarce. Acquiring a read lock first and then switching to a write lock would
// mean that the code gets complex and that the bucket gets deserialized twice
// adversely affecting the performance.
func (s *StoragePackerV2) putItemIntoBucket(bucket *BucketV2, item *Item) (string, error) {
	var lock *sync.RWMutex
	if bucket == nil {
		// Compute the index at which the primary bucket should reside
		primaryIndex, err := s.primaryBucketIndex(item.ID)
		if err != nil {
			return "", err
		}

		// Prepend the index by the prefix
		primaryKey := s.config.ViewPrefix + primaryIndex

		// Read the primary bucket
		bucket, lock, err = s.GetBucket(primaryKey, exclusive)
		if err != nil {
			s.unlockBucket(lock, exclusive)
			return "", err
		}

		// If the primary bucket does not exist, create one
		if bucket == nil {
			bucket = s.newBucket(primaryKey, 0)
			lockRaw, ok := s.bucketLocksCache.Load(primaryKey)
			if !ok {
				// For safety
				s.unlockBucket(lock, exclusive)
				return "", fmt.Errorf("unable to acquire lock for key %q", primaryKey)
			}

			lock = lockRaw.(*sync.RWMutex)
			lock.Lock()
		}
	}

	// For sanity
	if bucket == nil {
		s.unlockBucket(lock, exclusive)
		return "", fmt.Errorf("bucket is nil")
	}

	if lock == nil {
		lockRaw, ok := s.bucketLocksCache.Load(bucket.Key)
		if !ok {
			return "", fmt.Errorf("unable to acquire lock for key %q", bucket.Key)
		}
		lock = lockRaw.(*sync.RWMutex)
	}

	// Serializing and deserializing a proto message with empty map translates
	// to a nil. Hence, an already initialized field can potentially be nil
	// when its read from storage.
	if bucket.Buckets == nil {
		bucket.Buckets = make(map[string]*BucketV2)
	}

	// Compute the shard index to which the item belongs
	shardIndex, err := shardBucketIndex(item.ID, int(bucket.Depth), int(s.config.BucketCount), int(s.config.BucketShardCount))
	if err != nil {
		s.unlockBucket(lock, exclusive)
		return "", errwrap.Wrapf("failed to compute the bucket shard index: {{err}}", err)
	}

	// Check if the bucket shard to hold the item already exists
	bucketShard, ok := bucket.Buckets[shardIndex]

	// If the bucket shard is not present, create one
	if !ok {
		// The key to the bucket shard relative to its parent bucket
		shardKey := bucket.Key + "/" + shardIndex

		// Create the bucket shard to hold the item with an incremented depth
		bucketShard = s.newBucket(shardKey, bucket.Depth+1)

		// Add the newly created bucket shard to the parent bucket
		bucket.Buckets[shardIndex] = bucketShard
	}

	// For safety
	if bucketShard == nil {
		s.unlockBucket(lock, exclusive)
		return "", fmt.Errorf("bucket shard is nil")
	}

	// If the bucket shard is already pushed out, continue the operation in the
	// external bucket
	if bucketShard.External {
		externalBucket, _, err := s.GetBucket(bucketShard.Key, exclusive)

		// By now, the lock on the external bucket will be held. Release the
		// lock on the current bucket.
		s.unlockBucket(lock, exclusive)

		if err != nil {
			return "", err
		}

		if externalBucket == nil {
			return "", fmt.Errorf("failed to read the pushed out bucket shard: %q\n", bucketShard.Key)
		}

		return s.putItemIntoBucket(externalBucket, item)
	}

	//
	// Bucket shard is local to the parent bucket
	//

	// At this point only 2 things can happen. The current bucket accommodates
	// the item or the respective bucket shard gets pushed out as an
	// independent bucket to hold the item. In both the cases the lock on the
	// current bucket should be held until the operation completes.
	// Essentially, we can defer the lock release.

	// Ensure that the lock on the current bucket eventually gets released
	defer s.unlockBucket(lock, exclusive)

	// Update the item in the bucket shard
	bucketShard.Items[item.ID] = item

	// Check if the bucket exceeds the size limit after the item addition
	limitExceeded, err := s.bucketExceedsSizeLimit(bucket, item)
	if err != nil {
		return "", err
	}

	// If the bucket size is within the limit, persist the bucket and return
	if !limitExceeded {
		return bucketShard.Key, s.PutBucket(bucket)
	}

	//
	// The bucket size has exceeded the limit, push the bucket shard out as an
	// independent bucket and insert the item in the pushed out bucket.
	//

	// Mark the bucket shard as external, indicating that it doesn't
	// reside in its parent bucket
	bucketShard.External = true

	// Clone the bucket shard and use the clone as the pushed out bucket
	externalBucket, err := bucketShard.Clone()
	if err != nil {
		return "", err
	}

	// Clear the items in the bucket shard
	bucketShard.Items = nil

	// Split the items in the external bucket, among their respective bucket
	// shards
	err = s.splitItemsInBucket(externalBucket)
	if err != nil {
		return "", err
	}

	lockRaw, ok := s.bucketLocksCache.Load(externalBucket.Key)
	if !ok {
		return "", fmt.Errorf("failed to acquire lock for external bucket key %q", externalBucket.Key)
	}
	externalBucketLock := lockRaw.(*sync.RWMutex)
	externalBucketLock.Lock()

	// Insert the item in the bucket that got pushed out. Note that the lock on
	// the current bucket is still not released. When a bucket is getting
	// pushed out, holding the lock on the parent bucket ensures safety during
	// the externalization of a bucket shard.
	bucketKey, err := s.putItemIntoBucket(externalBucket, item)
	if err != nil {
		return "", err
	}

	// By now the item is already persisted in the external bucket, all the
	// while lock on its parent bucket was held. The parent bucket now needs to
	// persist the fact that a shard got moved out, thus reducing its size; and
	// release the lock.
	return bucketKey, s.PutBucket(bucket)
}

// Get reads a bucket from the storage while holding the respective lock on the
// bucket key
func (s *StoragePackerV2) GetBucket(key string, lockType bool) (*BucketV2, *sync.RWMutex, error) {
	if key == "" {
		return nil, nil, fmt.Errorf("missing bucket key")
	}

	// Check if there exists a lock for the bucket key
	_, ok := s.bucketLocksCache.Load(key)

	// If not, create one
	if !ok {
		s.bucketLocksCache.Store(key, &sync.RWMutex{})
	}

	// Fetch the lock for the key
	lockRaw, ok := s.bucketLocksCache.Load(key)
	// A lock should have been created. If not, error out.
	if !ok {
		return nil, nil, fmt.Errorf("failed to acquire lock for bucket key %q", key)
	}

	lock := lockRaw.(*sync.RWMutex)

	// Acquire the lock on the bucket key
	if lockType == exclusive {
		lock.Lock()
	} else {
		lock.RLock()
	}

	// Read the bucket from the underlying view
	entry, err := s.config.View.Get(context.Background(), key)
	if err != nil {
		return nil, lock, errwrap.Wrapf("failed to read bucket: {{err}}", err)
	}
	if entry == nil {
		// If the bucket key is invalid, there shouldn't be a corresponding
		// lock for it.
		s.bucketLocksCache.Delete(key)

		// Release the lock
		s.unlockBucket(lock, lockType)
		return nil, nil, nil
	}

	var bucketWrapper BucketWrapper
	err = proto.Unmarshal(entry.Value, &bucketWrapper)
	if err != nil {
		return nil, lock, errwrap.Wrapf("failed to decode bucket wrapper: {{err}}", err)
	}

	var bucket BucketV2
	err = proto.Unmarshal(bucketWrapper.Data, &bucket)
	if err != nil {
		return nil, lock, errwrap.Wrapf("failed to decode bucket: {{err}}", err)
	}

	// Set the size of the bucket within itself
	bucket.Size = bucketWrapper.Size

	return &bucket, lock, nil
}

// Put stores a bucket in storage. A write lock on the bucket key should be
// held by the caller of this function.
func (s *StoragePackerV2) PutBucket(bucket *BucketV2) error {
	if bucket == nil {
		return fmt.Errorf("nil bucket entry")
	}

	if bucket.Key == "" {
		return fmt.Errorf("missing bucket key")
	}

	if !strings.HasPrefix(bucket.Key, s.config.ViewPrefix) {
		return fmt.Errorf("bucket entry key should have %q prefix", s.config.ViewPrefix)
	}

	marshaledBucket, err := proto.Marshal(bucket)
	if err != nil {
		return err
	}

	// Store the marshaled bucket and its size in the wrapper
	bucketWrapper := &BucketWrapper{
		Data: marshaledBucket,
		Size: int64(len(marshaledBucket)),
	}

	marshaledWrapper, err := proto.Marshal(bucketWrapper)
	if err != nil {
		return err
	}

	// Persist the bucket wrapper
	return s.config.View.Put(context.Background(), &logical.StorageEntry{
		Key:   bucket.Key,
		Value: marshaledWrapper,
	})
}

// getItemFromBucket is a recursive function that fetches the given item ID in
// the bucket hierarchy
func (s *StoragePackerV2) getItemFromBucket(bucket *BucketV2, itemID string) (*Item, error) {
	var lock *sync.RWMutex
	if bucket == nil {
		primaryIndex, err := s.primaryBucketIndex(itemID)
		if err != nil {
			return nil, err
		}

		bucket, lock, err = s.GetBucket(s.config.ViewPrefix+primaryIndex, shared)
		if err != nil {
			s.unlockBucket(lock, shared)
			return nil, errwrap.Wrapf("failed to read packed storage item: {{err}}", err)
		}
	}

	if bucket == nil {
		s.unlockBucket(lock, shared)
		return nil, nil
	}

	if lock == nil {
		lockRaw, ok := s.bucketLocksCache.Load(bucket.Key)
		if !ok {
			return nil, fmt.Errorf("unable to acquire lock for key %q", bucket.Key)
		}
		lock = lockRaw.(*sync.RWMutex)
	}

	shardIndex, err := shardBucketIndex(itemID, int(bucket.Depth), int(s.config.BucketCount), int(s.config.BucketShardCount))
	if err != nil {
		s.unlockBucket(lock, shared)
		return nil, errwrap.Wrapf("failed to compute the bucket shard index: {{err}}", err)
	}

	bucketShard, ok := bucket.Buckets[shardIndex]
	if !ok {
		s.unlockBucket(lock, shared)
		return nil, nil
	}

	// If the bucket shard is already pushed out, continue the operation in the
	// external bucket
	if bucketShard.External {
		externalBucket, _, err := s.GetBucket(bucketShard.Key, shared)

		// By now, the lock on the external bucket will be held. Release the
		// lock on the current bucket.
		s.unlockBucket(lock, shared)

		if err != nil {
			return nil, err
		}

		if externalBucket == nil {
			return nil, fmt.Errorf("failed to read external bucket: %q\n", bucketShard.Key)
		}

		return s.getItemFromBucket(externalBucket, itemID)
	}

	// At this point the item either has to be local to the bucket or it
	// doesn't exist.

	// Ensure that the lock on the current bucket eventually gets released
	defer s.unlockBucket(lock, shared)

	return bucketShard.Items[itemID], nil
}

// deleteItemFromBucket is a recursive function that finds the bucket holding
// the item corresponding to the given item ID, and removes the item from it.
func (s *StoragePackerV2) deleteItemFromBucket(bucket *BucketV2, itemID string) error {
	var lock *sync.RWMutex
	if bucket == nil {
		primaryIndex, err := s.primaryBucketIndex(itemID)
		if err != nil {
			return err
		}

		bucket, lock, err = s.GetBucket(s.config.ViewPrefix+primaryIndex, exclusive)
		if err != nil {
			s.unlockBucket(lock, exclusive)
			return errwrap.Wrapf("failed to read packed storage item: {{err}}", err)
		}
	}

	if bucket == nil {
		// For safety
		s.unlockBucket(lock, exclusive)
		return nil
	}

	if lock == nil {
		lockRaw, ok := s.bucketLocksCache.Load(bucket.Key)
		if !ok {
			return fmt.Errorf("unable to acquire lock for key %q", bucket.Key)
		}
		lock = lockRaw.(*sync.RWMutex)
	}

	shardIndex, err := shardBucketIndex(itemID, int(bucket.Depth), int(s.config.BucketCount), int(s.config.BucketShardCount))
	if err != nil {
		s.unlockBucket(lock, exclusive)
		return errwrap.Wrapf("failed to compute the bucket shard index: {{err}}", err)
	}

	bucketShard, ok := bucket.Buckets[shardIndex]
	if !ok {
		s.unlockBucket(lock, exclusive)
		return nil
	}

	// If the bucket shard is already pushed out, continue the operation in the
	// pushed out bucket
	if bucketShard.External {
		externalBucket, _, err := s.GetBucket(bucketShard.Key, exclusive)

		// By now, the lock on the external bucket will be held. Release the
		// lock on the current bucket.
		s.unlockBucket(lock, exclusive)

		if err != nil {
			return err
		}

		if externalBucket == nil {
			return fmt.Errorf("failed to read external bucket: %q\n", bucketShard.Key)
		}

		return s.deleteItemFromBucket(externalBucket, itemID)
	}

	// Ensure that the lock on the current bucket eventually gets released
	defer s.unlockBucket(lock, exclusive)

	// Delete the item from the respective shard
	delete(bucketShard.Items, itemID)

	// Persist the change
	return s.PutBucket(bucket)
}

// GetItem fetches the item using the given item identifier
func (s *StoragePackerV2) GetItem(itemID string) (*Item, error) {
	if itemID == "" {
		return nil, fmt.Errorf("empty item ID")
	}

	return s.getItemFromBucket(nil, itemID)
}

// PutItem persists the given item
func (s *StoragePackerV2) PutItem(item *Item) (string, error) {
	if item == nil {
		return "", fmt.Errorf("nil item")
	}

	if item.ID == "" {
		return "", fmt.Errorf("missing ID in item")
	}

	return s.putItemIntoBucket(nil, item)
}

// DeleteItem removes the item using the given item identifier
func (s *StoragePackerV2) DeleteItem(itemID string) error {
	if itemID == "" {
		return fmt.Errorf("empty item ID")
	}

	return s.deleteItemFromBucket(nil, itemID)
}

// bucketExceedsSizeLimit indicates if the given bucket is exceeding the
// configured size limit on the storage packer
func (s *StoragePackerV2) bucketExceedsSizeLimit(bucket *BucketV2, item *Item) (bool, error) {
	marshaledItem, err := proto.Marshal(item)
	if err != nil {
		return false, fmt.Errorf("failed to marshal item: %v", err)
	}

	size := bucket.Size + int64(len(marshaledItem))

	// Sharding of buckets begins when the size of the bucket reaches 90% of
	// the maximum allowed size. Hopefully, this compensates for data structure
	// adjustment costs and also avoids edge cases with respect to the limit
	// imposed by the underlying physical backend.
	max := math.Ceil(float64(s.config.BucketMaxSize) * 0.9)

	return float64(size) > max, nil
}

// splitItemsInBucket breaks the list of items in the bucket and divides them
// such that they belong to their respective bucket shards
func (s *StoragePackerV2) splitItemsInBucket(bucket *BucketV2) error {
	if bucket.Buckets == nil {
		bucket.Buckets = make(map[string]*BucketV2)
	}

	for itemID, item := range bucket.Items {
		shardIndex, err := shardBucketIndex(itemID, int(bucket.Depth), int(s.config.BucketCount), int(s.config.BucketShardCount))
		if err != nil {
			return err
		}
		bucketShard, ok := bucket.Buckets[shardIndex]
		if !ok {
			shardKey := bucket.Key + "/" + shardIndex
			bucketShard = s.newBucket(shardKey, bucket.Depth+1)
			bucket.Buckets[shardIndex] = bucketShard
		}
		bucketShard.Items[itemID] = item
	}

	// All the items are moved into their respective bucket shards. Clear out
	// the outer items.
	bucket.Items = nil

	return nil
}

// primaryBucketIndex returns the index of the primary bucket to which the
// given item belongs
func (s *StoragePackerV2) primaryBucketIndex(itemID string) (string, error) {
	// Hash the item ID
	hashVal, err := cryptoutil.Blake2b256Hash(itemID)
	if err != nil {
		return "", err
	}

	// Compute the bits required to enumerate all possible primary buckets
	bitCount := bitsNeeded(s.config.BucketCount)

	// Extract the index value of the primary bucket from the hash of the item ID
	return strutil.BitMaskedIndexHex(hashVal, bitCount)
}

// shardBucketIndex returns the index of the bucket shard to which the given
// item belongs at a particular depth.
func shardBucketIndex(itemID string, depth, bucketCount, bucketShardCount int) (string, error) {
	// Hash the item ID
	hashVal, err := cryptoutil.Blake2b256Hash(itemID)
	if err != nil {
		return "", err
	}

	// Compute the bits required to enumerate all the possible bucket shards
	shardsBitCount := bitsNeeded(bucketShardCount)

	// Compute the bits that are already consumed by the primary bucket and the
	// shards at each level.
	ignoreBits := bitsNeeded(bucketCount) + depth*shardsBitCount

	// Extract the index value of the bucket shard from the hash of the item ID
	return strutil.BitMaskedIndexHex(hashVal[ignoreBits:], shardsBitCount)
}

// bitsNeeded returns the minimum number of bits required to enumerate the
// natural numbers below the given value
func bitsNeeded(value int) int {
	return int(math.Ceil(math.Log2(float64(value))))
}

// newBucket creates a bucket using the defaults set on the storage packer
// instance
func (s *StoragePackerV2) newBucket(key string, depth int32) *BucketV2 {
	bucket := &BucketV2{
		Key:     key,
		Buckets: make(map[string]*BucketV2),
		Items:   make(map[string]*Item),
		Depth:   depth,
	}

	// Create a new lock to operate on the bucket key
	s.bucketLocksCache.Store(key, &sync.RWMutex{})

	return bucket
}

// NewStoragePackerV2 creates a new storage packer for a given view
func NewStoragePackerV2(config *Config) (*StoragePackerV2, error) {
	if config.View == nil {
		return nil, fmt.Errorf("nil view")
	}

	if config.ViewPrefix == "" {
		config.ViewPrefix = DefaultStoragePackerBucketsPrefix
	}

	if !strings.HasSuffix(config.ViewPrefix, "/") {
		config.ViewPrefix = config.ViewPrefix + "/"
	}

	if config.BucketCount == 0 {
		config.BucketCount = defaultBucketCount
	}

	if config.BucketShardCount == 0 {
		config.BucketShardCount = defaultBucketShardCount
	}

	if config.BucketMaxSize == 0 {
		config.BucketMaxSize = defaultBucketMaxSize
	}

	// Create a new packer object for the given view
	packer := &StoragePackerV2{
		config:           config,
		bucketLocksCache: &sync.Map{},
	}

	return packer, nil
}

func (s *StoragePackerV2) unlockBucket(lock *sync.RWMutex, lockType bool) {
	if lock == nil {
		return
	}

	if lockType == exclusive {
		lock.Unlock()
	} else {
		lock.RUnlock()
	}
}
