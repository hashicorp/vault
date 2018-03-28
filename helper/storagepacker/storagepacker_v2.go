package storagepacker

import (
	"context"
	"fmt"
	"math"
	"strings"
	"sync"

	radix "github.com/armon/go-radix"
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
	config       *Config
	bucketsCache *radix.Tree
}

type PackedBucket struct {
	sync.RWMutex
	Data *BucketV2
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
// to store the item based on the storage space available in the buckets.
func (s *StoragePackerV2) putItemIntoBucket(bucket *PackedBucket, item *Item) (string, error) {
	if bucket == nil {
		// Compute the index at which the primary bucket should reside
		primaryIndex, err := s.primaryBucketIndex(item.ID)
		if err != nil {
			return "", err
		}

		// Prepend the index with the prefix
		primaryKey := s.config.ViewPrefix + primaryIndex

		// Check if the primary bucket exists
		bucket, err = s.GetBucket(primaryKey)
		if err != nil {
			return "", err
		}

		// If the primary bucket does not exist, create one
		if bucket == nil {
			bucket = &PackedBucket{
				Data: s.newBucket(primaryKey, 0),
			}
		}
	}

	// For sanity
	if bucket == nil {
		return "", fmt.Errorf("bucket is nil")
	}

	// Compute the shard index to which the item belongs
	shardIndex, err := shardBucketIndex(item.ID, int(bucket.Data.Depth), int(s.config.BucketCount), int(s.config.BucketShardCount))
	if err != nil {
		return "", errwrap.Wrapf("failed to compute the bucket shard index: {{err}}", err)
	}

	bucket.Lock()

	// Check if the bucket shard to hold the item already exists
	bucketShard, ok := bucket.Data.Buckets[shardIndex]

	// If the bucket shard is not present, create one
	if !ok {
		// The key to the bucket shard relative to its parent bucket
		shardKey := bucket.Data.Key + "/" + shardIndex

		// Create the bucket shard to hold the item with an incremented depth
		bucketShard = s.newBucket(shardKey, bucket.Data.Depth+1)

		// Add the newly created bucket shard to the parent bucket
		bucket.Data.Buckets[shardIndex] = bucketShard
	}

	if bucketShard == nil {
		bucket.Unlock()
		return "", fmt.Errorf("bucket shard is nil")
	}

	// If the bucket shard is already pushed out, continue the operation in the
	// external bucket
	if bucketShard.External {
		bucket.Unlock()
		externalBucket, err := s.GetBucket(bucketShard.Key)
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
	defer bucket.Unlock()

	// Update the item in the bucket shard
	bucketShard.Items[item.ID] = item

	// Check if the bucket exceeds the size limit after the addition
	limitExceeded, err := s.bucketExceedsSizeLimit(bucket, item)
	if err != nil {
		return "", err
	}

	// If the bucket size is within the limit, return the updated bucket
	if !limitExceeded {
		return bucketShard.Key, s.PutBucket(bucket)
	}

	//
	// If the bucket size has exceeded the limit, push the bucket shard out as
	// an independent bucket and insert the item in the pushed out bucket.
	//

	// Mark the bucket shard as external, indicating that it doesn't
	// reside in its parent bucket
	bucketShard.External = true

	// Clone the bucket and use the clone as the pushed out bucket
	clone, err := bucketShard.Clone()
	if err != nil {
		return "", err
	}

	externalBucket := &PackedBucket{
		Data: clone,
	}

	// Clear the items in the pushed out bucket shard
	bucketShard.Items = nil

	// Split the items in the bucket that gets pushed out, among their
	// respective bucket shards
	err = s.splitItemsInBucket(externalBucket)
	if err != nil {
		return "", err
	}

	// Insert the item in the bucket that got pushed out. Note that the lock on
	// the current bucket is still not released. When a bucket is getting
	// pushed out, holding the lock on the parent bucket ensures safety during
	// the externalization of a bucket shard.
	bucketKey, err := s.putItemIntoBucket(externalBucket, item)
	if err != nil {
		return "", err
	}

	return bucketKey, s.PutBucket(bucket)
}

// Get reads a bucket from the storage
func (s *StoragePackerV2) GetBucket(key string) (*PackedBucket, error) {
	if key == "" {
		return nil, fmt.Errorf("missing bucket key")
	}

	_, raw, exists := s.bucketsCache.LongestPrefix(key)
	if exists {
		return raw.(*PackedBucket), nil
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

	bucket.Size = int64(len(entry.Value))

	return s.UpdateBucketCache(&PackedBucket{Data: &bucket}), nil
}

func (s *StoragePackerV2) UpdateBucketCache(bucket *PackedBucket) *PackedBucket {
	var packedBucket *PackedBucket
	_, raw, exists := s.bucketsCache.LongestPrefix(bucket.Data.Key)
	if exists {
		packedBucket = raw.(*PackedBucket)
	}
	if packedBucket == nil {
		packedBucket = bucket
	}
	s.bucketsCache.Insert(packedBucket.Data.Key, packedBucket)
	return packedBucket
}

// Put stores a bucket in storage
func (s *StoragePackerV2) PutBucket(bucket *PackedBucket) error {
	if bucket == nil {
		return fmt.Errorf("nil bucket entry")
	}

	if bucket.Data.Key == "" {
		return fmt.Errorf("missing bucket key")
	}

	if !strings.HasPrefix(bucket.Data.Key, s.config.ViewPrefix) {
		return fmt.Errorf("bucket entry key should have %q prefix", s.config.ViewPrefix)
	}

	marshaledBucket, err := proto.Marshal(bucket.Data)
	if err != nil {
		return err
	}

	err = s.config.View.Put(context.Background(), &logical.StorageEntry{
		Key:   bucket.Data.Key,
		Value: marshaledBucket,
	})
	if err != nil {
		return err
	}

	s.UpdateBucketCache(bucket)

	return nil
}

// getItemFromBucket is a recursive function that fetches the given item ID in
// the bucket hierarchy
func (s *StoragePackerV2) getItemFromBucket(bucket *PackedBucket, itemID string) (*Item, error) {
	if bucket == nil {
		primaryIndex, err := s.primaryBucketIndex(itemID)
		if err != nil {
			return nil, err
		}

		bucket, err = s.GetBucket(s.config.ViewPrefix + primaryIndex)
		if err != nil {
			return nil, errwrap.Wrapf("failed to read packed storage item: {{err}}", err)
		}
	}

	if bucket == nil {
		return nil, nil
	}

	shardIndex, err := shardBucketIndex(itemID, int(bucket.Data.Depth), int(s.config.BucketCount), int(s.config.BucketShardCount))
	if err != nil {
		return nil, errwrap.Wrapf("failed to compute the bucket shard index: {{err}}", err)
	}

	bucket.RLock()

	bucketShard, ok := bucket.Data.Buckets[shardIndex]
	if !ok {
		bucket.RUnlock()
		return nil, nil
	}

	// If the bucket shard is already pushed out, continue the operation in the
	// external bucket
	if bucketShard.External {
		bucket.RUnlock()
		externalBucket, err := s.GetBucket(bucketShard.Key)
		if err != nil {
			return nil, err
		}

		if externalBucket == nil {
			return nil, fmt.Errorf("failed to read external bucket: %q\n", bucketShard.Key)
		}

		return s.getItemFromBucket(externalBucket, itemID)
	}

	defer bucket.RUnlock()

	return bucketShard.Items[itemID], nil
}

// deleteItemFromBucket is a recursive function that finds the bucket holding
// the item and removes the item from it
func (s *StoragePackerV2) deleteItemFromBucket(bucket *PackedBucket, itemID string) error {
	if bucket == nil {
		primaryIndex, err := s.primaryBucketIndex(itemID)
		if err != nil {
			return err
		}

		bucket, err = s.GetBucket(s.config.ViewPrefix + primaryIndex)
		if err != nil {
			return errwrap.Wrapf("failed to read packed storage item: {{err}}", err)
		}
	}

	if bucket == nil {
		return nil
	}

	shardIndex, err := shardBucketIndex(itemID, int(bucket.Data.Depth), int(s.config.BucketCount), int(s.config.BucketShardCount))
	if err != nil {
		return errwrap.Wrapf("failed to compute the bucket shard index: {{err}}", err)
	}

	bucket.Lock()

	bucketShard, ok := bucket.Data.Buckets[shardIndex]
	if !ok {
		bucket.Unlock()
		return nil
	}

	// If the bucket shard is already pushed out, continue the operation in the
	// pushed out bucket
	if bucketShard.External {
		bucket.Unlock()
		externalBucket, err := s.GetBucket(bucketShard.Key)

		if err != nil {
			return err
		}

		if externalBucket == nil {
			return fmt.Errorf("failed to read external bucket: %q\n", bucketShard.Key)
		}

		return s.deleteItemFromBucket(externalBucket, itemID)
	}

	defer bucket.Unlock()

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
func (s *StoragePackerV2) bucketExceedsSizeLimit(bucket *PackedBucket, item *Item) (bool, error) {
	marshaledItem, err := proto.Marshal(item)
	if err != nil {
		return false, fmt.Errorf("failed to marshal item: %v", err)
	}

	size := bucket.Data.Size + int64(len(marshaledItem))

	// Sharding of buckets begins when the size of the bucket reaches 90% of
	// the maximum allowed size. Hopefully, this compensates for data structure
	// adjustment costs and also avoids edge cases with respect to the limit
	// imposed by the underlying physical backend.
	max := math.Ceil((float64(s.config.BucketMaxSize) * float64(90)) / float64(100))

	return float64(size) > max, nil
}

// splitItemsInBucket breaks the list of items in the bucket and divides them
// such that they belong to their respective bucket shards
func (s *StoragePackerV2) splitItemsInBucket(bucket *PackedBucket) error {
	for itemID, item := range bucket.Data.Items {
		shardIndex, err := shardBucketIndex(itemID, int(bucket.Data.Depth), int(s.config.BucketCount), int(s.config.BucketShardCount))
		if err != nil {
			return err
		}
		bucketShard, ok := bucket.Data.Buckets[shardIndex]
		if !ok {
			shardKey := bucket.Data.Key + "/" + shardIndex
			bucketShard = s.newBucket(shardKey, bucket.Data.Depth+1)
			bucket.Data.Buckets[shardIndex] = bucketShard
		}
		bucketShard.Items[itemID] = item
	}

	bucket.Data.Items = nil

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
		config:       config,
		bucketsCache: radix.New(),
	}

	return packer, nil
}
