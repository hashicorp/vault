package storagepacker

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"

	radix "github.com/armon/go-radix"
	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/helper/strutil"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/cryptoutil"
	"github.com/hashicorp/vault/logical"
)

const (
	defaultBucketBaseCount  = 256
	defaultBucketShardCount = 16
	// Larger size of the bucket size adversely affects the performance of the
	// storage packer. Also, some of the backends impose a maximum size limit
	// on the objects that gets persisted. For example, Consul imposes 512KB
	// and DynamoDB imposes 400KB. Going forward, if there exists storage
	// backends that has more constrained limits, this will have to become more
	// flexible. For now, 380KB seems like a decent bargain.
	defaultBucketMaxSize = 380 * 1024
)

type Config struct {
	// View is the storage to be used by all the buckets
	View logical.Storage

	// ViewPrefix is the prefix to be used for the buckets in the view
	ViewPrefix string

	// Logger for output
	Logger log.Logger

	// BucketBaseCount is the number of buckets to create at the base level.
	// The value should be a power of 2.
	BucketBaseCount int

	// BucketShardCount is the number of sub-buckets a bucket gets sharded into
	// when it reaches the maximum threshold. The value should be a power of 2.
	BucketShardCount int

	// BucketMaxSize (in bytes) is the maximum allowed size per bucket. When
	// the size of the bucket reaches a threshold relative to this limit, it
	// gets sharded into the configured number of pieces incrementally.
	BucketMaxSize int64
}

// StoragePackerV2 packs many items into abstractions called buckets. The goal
// is to employ a reduced number of storage entries for a relatively huge
// number of items. This is the second version of the utility which supports
// indefinitely expanding the capacity of the storage by sharding the buckets
// when they exceed the imposed limit.
type StoragePackerV2 struct {
	config       *Config
	bucketsCache *radix.Tree
}

// LockedBucket embeds a bucket and its corresponding lock to ensure thread
// safety
type LockedBucket struct {
	*BucketV2
	lock sync.RWMutex
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

	if config.BucketBaseCount == 0 {
		config.BucketBaseCount = defaultBucketBaseCount
	}

	if config.BucketShardCount == 0 {
		config.BucketShardCount = defaultBucketShardCount
	}

	if config.BucketMaxSize == 0 {
		config.BucketMaxSize = defaultBucketMaxSize
	}

	if !isPowerOfTwo(config.BucketBaseCount) {
		return nil, fmt.Errorf("bucket base count of %d is not a power of two", config.BucketBaseCount)
	}

	if !isPowerOfTwo(config.BucketShardCount) {
		return nil, fmt.Errorf("bucket shard count of %d is not a power of two", config.BucketShardCount)
	}

	if config.BucketShardCount < 2 {
		return nil, fmt.Errorf("bucket shard count should at least be 2")
	}

	// Create a new packer object for the given view
	packer := &StoragePackerV2{
		config:       config,
		bucketsCache: radix.New(),
	}

	return packer, nil
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

// Get reads a bucket from the storage
func (s *StoragePackerV2) GetBucket(key string) (*LockedBucket, error) {
	if key == "" {
		return nil, fmt.Errorf("missing bucket key")
	}

	raw, exists := s.bucketsCache.Get(key)
	if exists {
		return raw.(*LockedBucket), nil
	}

	// Read from the underlying view
	entry, err := s.config.View.Get(context.Background(), key)
	if err != nil {
		return nil, errwrap.Wrapf("failed to read bucket: {{err}}", err)
	}
	if entry == nil {
		return nil, nil
	}

	var bucket BucketV2
	err = proto.Unmarshal(entry.Value, &bucket)
	if err != nil {
		return nil, errwrap.Wrapf("failed to decode bucket: {{err}}", err)
	}

	// Serializing and deserializing a proto message with empty map translates
	// to a nil. Ensure that the required fields are initialized properly.
	if bucket.Buckets == nil {
		bucket.Buckets = make(map[string]*BucketV2)
	}
	if bucket.Items == nil {
		bucket.Items = make(map[string]*Item)
	}

	// Update the unencrypted size of the bucket
	bucket.Size = int64(len(entry.Value))

	lb := &LockedBucket{
		BucketV2: &bucket,
	}
	s.bucketsCache.Insert(bucket.Key, lb)

	return lb, nil
}

// Put stores a bucket in storage
func (s *StoragePackerV2) PutBucket(bucket *LockedBucket) error {
	if bucket == nil {
		return fmt.Errorf("nil bucket entry")
	}

	if bucket.Key == "" {
		return fmt.Errorf("missing bucket key")
	}

	if !strings.HasPrefix(bucket.Key, s.config.ViewPrefix) {
		return fmt.Errorf("bucket entry key should have %q prefix", s.config.ViewPrefix)
	}

	marshaledBucket, err := proto.Marshal(bucket.BucketV2)
	if err != nil {
		return err
	}

	err = s.config.View.Put(context.Background(), &logical.StorageEntry{
		Key:   bucket.Key,
		Value: marshaledBucket,
	})
	if err != nil {
		return err
	}

	bucket.Size = int64(len(marshaledBucket))

	s.bucketsCache.Insert(bucket.Key, bucket)

	return nil
}

// putItem is a recursive function that finds the appropriate bucket
// to store the item based on the storage space available in the buckets.
func (s *StoragePackerV2) putItem(bucket *LockedBucket, item *Item, depth int) (string, error) {
	// Bucket will be nil for the first time when its not known which base
	// level bucket the item belongs to.
	if bucket == nil {
		// Enforce zero depth
		depth = 0

		// Compute the index of the base bucket
		baseIndex, err := s.baseBucketIndex(item.ID)
		if err != nil {
			return "", err
		}

		// Prepend the index with the prefix
		baseKey := s.config.ViewPrefix + baseIndex

		// Check if the base bucket exists
		bucket, err = s.GetBucket(baseKey)
		if err != nil {
			return "", err
		}

		// If the base bucket does not exist, create one
		if bucket == nil {
			bucket = s.newBucket(baseKey)
		}
	}

	// Compute the shard index to which the item belongs
	shardIndex, err := s.shardBucketIndex(item.ID, depth)
	if err != nil {
		return "", errwrap.Wrapf("failed to compute the bucket shard index: {{err}}", err)
	}
	shardKey := bucket.Key + "/" + shardIndex

	// Acquire lock on the bucket
	bucket.lock.Lock()

	if bucket.Sharded {
		// If the bucket is already sharded out, release the lock and continue
		// insertion at the next level.
		bucket.lock.Unlock()
		shardedBucket, err := s.GetBucket(shardKey)
		if err != nil {
			return "", err
		}
		if shardedBucket == nil {
			shardedBucket = s.newBucket(shardKey)
		}
		return s.putItem(shardedBucket, item, depth+1)
	}

	// From this point on, the item may get inserted either in the current
	// bucket or at its next level. In both cases, there will be a need to
	// persist the current bucket. Hence the lock on the current bucket is
	// deferred.
	defer bucket.lock.Unlock()

	// Check if a bucket shard is already present for the shard index. If not,
	// create one.
	bucketShard, ok := bucket.Buckets[shardIndex]
	if !ok {
		bucketShard = s.newBucket(shardKey).BucketV2
		bucket.Buckets[shardIndex] = bucketShard
	}

	// Check if the insertion of the item makes the bucket size exceed the
	// limit.
	exceedsLimit, err := s.bucketExceedsSizeLimit(bucket, item)
	if err != nil {
		return "", err
	}

	// If the bucket size after addition of the item doesn't exceed the limit,
	// insert the item persist the bucket.
	if !exceedsLimit {
		bucketShard.Items[item.ID] = item
		return bucket.Key, s.PutBucket(bucket)
	}

	// The bucket size after addition of the item exceeds the size limit. Split
	// the bucket into shards.
	err = s.splitBucket(bucket, depth)
	if err != nil {
		return "", err
	}

	shardedBucket, err := s.GetBucket(bucketShard.Key)
	if err != nil {
		return "", err
	}

	bucketKey, err := s.putItem(shardedBucket, item, depth+1)
	if err != nil {
		return "", err
	}

	return bucketKey, s.PutBucket(bucket)
}

// getItem is a recursive function that fetches the given item ID in
// the bucket hierarchy
func (s *StoragePackerV2) getItem(bucket *LockedBucket, itemID string, depth int) (*Item, error) {
	if bucket == nil {
		// Enforce zero depth
		depth = 0

		baseIndex, err := s.baseBucketIndex(itemID)
		if err != nil {
			return nil, err
		}

		bucket, err = s.GetBucket(s.config.ViewPrefix + baseIndex)
		if err != nil {
			return nil, errwrap.Wrapf("failed to read packed storage item: {{err}}", err)
		}
	}

	if bucket == nil {
		return nil, nil
	}

	shardIndex, err := s.shardBucketIndex(itemID, depth)
	if err != nil {
		return nil, errwrap.Wrapf("failed to compute the bucket shard index: {{err}}", err)
	}

	shardKey := bucket.Key + "/" + shardIndex

	bucket.lock.RLock()

	if bucket.Sharded {
		bucket.lock.RUnlock()
		shardedBucket, err := s.GetBucket(shardKey)
		if err != nil {
			return nil, err
		}
		if shardedBucket == nil {
			return nil, nil
		}
		return s.getItem(shardedBucket, itemID, depth+1)
	}

	defer bucket.lock.RUnlock()

	bucketShard, ok := bucket.Buckets[shardIndex]
	if !ok {
		return nil, nil
	}

	if bucketShard == nil {
		return nil, nil
	}

	return bucketShard.Items[itemID], nil
}

// deleteItem is a recursive function that finds the bucket holding
// the item and removes the item from it
func (s *StoragePackerV2) deleteItem(bucket *LockedBucket, itemID string, depth int) error {
	if bucket == nil {
		// Enforce zero depth
		depth = 0

		baseIndex, err := s.baseBucketIndex(itemID)
		if err != nil {
			return err
		}

		bucket, err = s.GetBucket(s.config.ViewPrefix + baseIndex)
		if err != nil {
			return errwrap.Wrapf("failed to read packed storage item: {{err}}", err)
		}
	}

	if bucket == nil {
		return nil
	}

	shardIndex, err := s.shardBucketIndex(itemID, depth)
	if err != nil {
		return errwrap.Wrapf("failed to compute the bucket shard index: {{err}}", err)
	}

	shardKey := bucket.Key + "/" + shardIndex

	bucket.lock.Lock()

	if bucket.Sharded {
		bucket.lock.Unlock()
		shardedBucket, err := s.GetBucket(shardKey)
		if err != nil {
			return err
		}
		if shardedBucket == nil {
			return nil
		}
		return s.deleteItem(shardedBucket, itemID, depth+1)
	}

	defer bucket.lock.Unlock()

	bucketShard, ok := bucket.Buckets[shardIndex]
	if !ok {
		return nil
	}

	if bucketShard == nil {
		return nil
	}

	delete(bucketShard.Items, itemID)

	return s.PutBucket(bucket)
}

// GetItem fetches the item using the given item identifier
func (s *StoragePackerV2) GetItem(itemID string) (*Item, error) {
	if itemID == "" {
		return nil, fmt.Errorf("empty item ID")
	}

	return s.getItem(nil, itemID, 0)
}

// PutItem persists the given item
func (s *StoragePackerV2) PutItem(item *Item) (string, error) {
	if item == nil {
		return "", fmt.Errorf("nil item")
	}

	if item.ID == "" {
		return "", fmt.Errorf("missing ID in item")
	}

	bucketKey, err := s.putItem(nil, item, 0)
	if err != nil {
		return "", err
	}

	return bucketKey, nil
}

// DeleteItem removes the item using the given item identifier
func (s *StoragePackerV2) DeleteItem(itemID string) error {
	if itemID == "" {
		return fmt.Errorf("empty item ID")
	}

	return s.deleteItem(nil, itemID, 0)
}

// bucketExceedsSizeLimit computes if the given bucket is exceeding the
// configured size limit on the storage packer
func (s *StoragePackerV2) bucketExceedsSizeLimit(bucket *LockedBucket, item *Item) (bool, error) {
	marshaledItem, err := proto.Marshal(item)
	if err != nil {
		return false, fmt.Errorf("failed to marshal item: %v", err)
	}

	expectedBucketSize := bucket.Size + int64(len(marshaledItem))

	// The objects that leave storage packer to get persisted get inflated due
	// to extra bits coming off of encryption. So, we consider the bucket to be
	// full much earlier to compensate for the encryption overhead. Testing
	// with the threshold of 70% of the max size resulted in object sizes
	// coming dangerously close to the actual limit. Hence, setting 60% as the
	// cut-off value. This is purely a heuristic threshold.
	max := math.Ceil((float64(s.config.BucketMaxSize) * float64(60)) / float64(100))

	return float64(expectedBucketSize) > max, nil
}

func (s *StoragePackerV2) splitBucket(bucket *LockedBucket, depth int) error {
	for _, shard := range bucket.Buckets {
		for itemID, item := range shard.Items {
			if shard.Buckets == nil {
				shard.Buckets = make(map[string]*BucketV2)
			}
			subShardIndex, err := s.shardBucketIndex(itemID, depth+1)
			if err != nil {
				return err
			}
			subShard, ok := shard.Buckets[subShardIndex]
			if !ok {
				subShardKey := shard.Key + "/" + subShardIndex
				subShard = s.newBucket(subShardKey).BucketV2
				shard.Buckets[subShardIndex] = subShard
			}
			subShard.Items[itemID] = item
		}

		shard.Items = nil
		err := s.PutBucket(&LockedBucket{BucketV2: shard})
		if err != nil {
			return err
		}
	}
	bucket.Buckets = nil
	bucket.Sharded = true
	return nil
}

// baseBucketIndex returns the index of the base bucket to which the
// given item belongs
func (s *StoragePackerV2) baseBucketIndex(itemID string) (string, error) {
	// Hash the item ID
	hashVal, err := cryptoutil.Blake2b256Hash(itemID)
	if err != nil {
		return "", err
	}

	// Extract the index value of the base bucket from the hash of the item ID
	return strutil.BitMaskedIndexHex(hashVal, bitsNeeded(s.config.BucketBaseCount))
}

// shardBucketIndex returns the index of the bucket shard to which the given
// item belongs at a particular depth.
func (s *StoragePackerV2) shardBucketIndex(itemID string, depth int) (string, error) {
	// Hash the item ID
	hashVal, err := cryptoutil.Blake2b256Hash(itemID)
	if err != nil {
		return "", err
	}

	// Compute the bits required to enumerate base buckets
	shardsBitCount := bitsNeeded(s.config.BucketShardCount)

	// Compute the bits that are already consumed by the base bucket and the
	// shards at previous levels.
	ignoreBits := bitsNeeded(s.config.BucketBaseCount) + depth*shardsBitCount

	// Extract the index value of the bucket shard from the hash of the item ID
	return strutil.BitMaskedIndexHex(hashVal[ignoreBits:], shardsBitCount)
}

// bitsNeeded returns the minimum number of bits required to enumerate the
// natural numbers below the given value
func bitsNeeded(value int) int {
	if value < 2 {
		return 1
	}
	bitCount := int(math.Ceil(math.Log2(float64(value))))
	if isPowerOfTwo(value) {
		bitCount++
	}
	return bitCount
}

func isPowerOfTwo(val int) bool {
	return val != 0 && (val&(val-1) == 0)
}

func (s *StoragePackerV2) newBucket(key string) *LockedBucket {
	return &LockedBucket{
		BucketV2: &BucketV2{
			Key:     key,
			Buckets: make(map[string]*BucketV2),
			Items:   make(map[string]*Item),
		},
	}
}

type WalkFunc func(item *Item) error

// Walk traverses through all the buckets and all the items in each bucket and
// invokes the given function on each item.
func (s *StoragePackerV2) Walk(fn WalkFunc) error {
	var err error
	for base := 0; base < s.config.BucketBaseCount; base++ {
		baseKey := s.config.ViewPrefix + strconv.FormatInt(int64(base), 16)
		err = s.bucketWalk(baseKey, fn)
		if err != nil {
			return err
		}
	}
	return nil
}

// bucketWalk is a pre-order traversal of the bucket hierarchy starting from
// the bucket corresponding to the given key. The function fn will be called on
// all the items in the hierarchy.
func (s *StoragePackerV2) bucketWalk(key string, fn WalkFunc) error {
	bucket, err := s.GetBucket(key)
	if err != nil {
		return err
	}
	if bucket == nil {
		return nil
	}

	if !bucket.Sharded {
		for _, b := range bucket.Buckets {
			for _, item := range b.Items {
				err := fn(item)
				if err != nil {
					return err
				}
			}
		}
		return nil
	}

	for i := 0; i < s.config.BucketShardCount; i++ {
		shardKey := bucket.Key + "/" + strconv.FormatInt(int64(i), 16)
		err = s.bucketWalk(shardKey, fn)
		if err != nil {
			return err
		}
	}

	return nil
}
