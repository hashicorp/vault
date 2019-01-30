package storagepacker

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"sync"

	radix "github.com/armon/go-radix"
	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/compressutil"
	"github.com/hashicorp/vault/helper/locksutil"
	"github.com/hashicorp/vault/logical"
)

type HashType uint

const (
	HashTypeBlake2b256 HashType = iota
	HashTypeMD5
)

const (
	defaultBucketBaseCount  = 256
	defaultBucketShardCount = 16
	// Larger size of the bucket size adversely affects the performance of the
	// storage packer. Also, some of the backends impose a maximum size limit
	// on the objects that gets persisted. For example, Consul imposes 256KB if using transactions
	// and DynamoDB imposes 400KB. Going forward, if there exists storage
	// backends that has more constrained limits, this will have to become more
	// flexible. For now, 240KB seems like a decent value.
	defaultBucketMaxSize = 240 * 1024

	DefaultStoragePackerBucketsPrefix = "packer/buckets/"
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

	// The hash type to use at the base bucket level. Shards always use blake.
	// For backwards compat.
	BaseHashType HashType
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
}

// LockedBucket embeds a bucket and its corresponding lock to ensure thread
// safety
type LockedBucket struct {
	*Bucket
	lock sync.RWMutex
}

// BucketPath returns the storage entry key for a given bucket key
func (s *StoragePackerV1) BucketPath(bucketKey string) string {
	return s.ViewPrefix + bucketKey
}

// BucketKeyHash returns the MD5 hash of the bucket storage key in which
// the item will be stored. The choice of MD5 is only for hash performance
// reasons since its value is not used for any security sensitive operation.
func (s *StoragePackerV1) BucketKeyHashByItemID(itemID string) string {
	return s.BucketKeyHashByKey(s.BucketPath(s.BucketKey(itemID)))
}

// BucketKeyHashByKey returns the MD5 hash of the bucket storage key
func (s *StoragePackerV1) BucketKeyHashByKey(bucketKey string) string {
	hf := md5.New()
	hf.Write([]byte(bucketKey))
	return hex.EncodeToString(hf.Sum(nil))
}

// View returns the storage view configured to be used by the packer
func (s *StoragePackerV1) StorageView() logical.Storage {
	return s.View
}

// Get returns a bucket for a given key
func (s *StoragePackerV1) GetBucket(key string) (*Bucket, error) {
	if key == "" {
		return nil, fmt.Errorf("missing bucket key")
	}

	lock := locksutil.LockForKey(s.storageLocks, key)
	lock.RLock()
	defer lock.RUnlock()

	// Read from the underlying view
	storageEntry, err := s.View.Get(context.Background(), key)
	if err != nil {
		return nil, errwrap.Wrapf("failed to read packed storage entry: {{err}}", err)
	}
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

	var bucket Bucket
	err = proto.Unmarshal(uncompressedData, &bucket)
	if err != nil {
		return nil, errwrap.Wrapf("failed to decode packed storage entry: {{err}}", err)
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
func (s *StoragePackerV1) BucketIndex(key string) uint8 {
	hf := md5.New()
	hf.Write([]byte(key))
	return uint8(hf.Sum(nil)[0])
}

// BucketKey returns the bucket key for a given item ID
func (s *StoragePackerV1) BucketKey(itemID string) string {
	return strconv.Itoa(int(s.BucketIndex(itemID)))
}

// DeleteItem removes the storage entry which the given key refers to from its
// corresponding bucket.
func (s *StoragePackerV1) DeleteItem(itemID string) error {

	if itemID == "" {
		return fmt.Errorf("empty item ID")
	}

	// Get the bucket key
	bucketKey := s.BucketKey(itemID)

	// Prepend the view prefix
	bucketPath := s.BucketPath(bucketKey)

	// Read from underlying view
	storageEntry, err := s.View.Get(context.Background(), bucketPath)
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

	var bucket Bucket
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
		err = s.PutBucket(&bucket)
		if err != nil {
			return err
		}
	}

	return nil
}

// Put stores a packed bucket entry
func (s *StoragePackerV1) PutBucket(bucket *Bucket) error {
	if bucket == nil {
		return fmt.Errorf("nil bucket entry")
	}

	if bucket.Key == "" {
		return fmt.Errorf("missing key")
	}

	if !strings.HasPrefix(bucket.Key, s.ViewPrefix) {
		return fmt.Errorf("incorrect prefix; bucket entry key should have %q prefix", s.ViewPrefix)
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
	err = s.View.Put(context.Background(), &logical.StorageEntry{
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
func (s *StoragePackerV1) GetItem(itemID string) (*Item, error) {
	if itemID == "" {
		return nil, fmt.Errorf("empty item ID")
	}

	bucketKey := s.BucketKey(itemID)
	bucketPath := s.BucketPath(bucketKey)

	// Fetch the bucket entry
	bucket, err := s.GetBucket(bucketPath)
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
func (s *StoragePackerV1) PutItem(item *Item) error {
	if item == nil {
		return fmt.Errorf("nil item")
	}

	if item.ID == "" {
		return fmt.Errorf("missing ID in item")
	}

	var err error
	bucketKey := s.BucketKey(item.ID)
	bucketPath := s.BucketPath(bucketKey)

	bucket := &Bucket{
		Key: bucketPath,
	}

	// In this case, we persist the storage entry regardless of the read
	// storageEntry below is nil or not. Hence, directly acquire write lock
	// even to read the entry.
	lock := locksutil.LockForKey(s.storageLocks, bucketPath)
	lock.Lock()
	defer lock.Unlock()

	// Check if there is an existing bucket for a given key
	storageEntry, err := s.View.Get(context.Background(), bucketPath)
	if err != nil {
		return errwrap.Wrapf("failed to read packed storage bucket entry: {{err}}", err)
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
			return errwrap.Wrapf("failed to decompress packed storage entry: {{err}}", err)
		}
		if notCompressed {
			uncompressedData = storageEntry.Value
		}

		err = proto.Unmarshal(uncompressedData, bucket)
		if err != nil {
			return errwrap.Wrapf("failed to decode packed storage entry: {{err}}", err)
		}

		err = bucket.upsert(item)
		if err != nil {
			return errwrap.Wrapf("failed to update entry in packed storage entry: {{err}}", err)
		}
	}

	// Persist the result
	return s.PutBucket(bucket)
}

// NewStoragePackerV1 creates a new storage packer for a given view
func NewStoragePackerV1(ctx context.Context, config *Config) (*StoragePackerV1, error) {
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

	// At this point, look for an existing saved configuration
	var needPersist bool
	entry, err := config.View.Get(ctx, config.ViewPrefix+"config")
	if err != nil {
		return nil, errwrap.Wrapf("error checking for existing storagepacker config: {{err}}", err)
	}
	if entry != nil {
		needPersist = false
		var exist Config
		if err := entry.DecodeJSON(&exist); err != nil {
			return nil, errwrap.Wrapf("error decoding existing storagepacker config: {{err}}", err)
		}
		// If we have an existing config, we copy the only two things we need
		// constant:
		//
		// 1. The bucket base count, so we know how many to expect
		// 2. The base hash type. We need to know how to hash at the base. All
		// shards will use Blake.
		//
		// The rest of the values can change; the max size can change based on
		// e.g. if storage is migrated, so as long as we don't move to a new
		// location with a smaller value we're fine (and even then we're fine
		// if we can read it; otherwise storage migration would have failed
		// anyways). The shard count is recorded in each bucket at the time
		// it's sharded; if we realize it's more efficient to do some other
		// value later we can update it and use that going forward for new
		// shards.
		config.BucketBaseCount = exist.BucketBaseCount
		config.BaseHashType = exist.BaseHashType
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

	if needPersist {
		entry, err := logical.StorageEntryJSON(config.ViewPrefix+"config", config)
		if err != nil {
			return nil, errwrap.Wrapf("error encoding storagepacker config: {{err}}", err)
		}
		if err := config.View.Put(ctx, entry); err != nil {
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

// isPowerOfTwo returns true if the given value is a power of two, false
// otherwise. We also return false on 1 because there'd be no point.
func isPowerOfTwo(val int) bool {
	switch val {
	case 0, 1:
		return false
	default:
		return val&(val-1) == 0
	}
	return false
}
