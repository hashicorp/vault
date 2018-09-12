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
	"github.com/hashicorp/vault/logical"
)

const (
	bucketCount                = 256
	StoragePackerBucketsPrefix = "packer/buckets/"
)

// StoragePacker packs the objects into a specific number of buckets by hashing
// its ID and indexing it. Currently this supports only 256 bucket entries and
// hence relies on the first byte of the hash value for indexing. The items
// that gets inserted into the packer should implement StorageBucketItem
// interface.
type StoragePacker struct {
	view         logical.Storage
	logger       log.Logger
	storageLocks []*locksutil.LockEntry
	viewPrefix   string
}

// BucketPath returns the storage entry key for a given bucket key
func (s *StoragePacker) BucketPath(bucketKey string) string {
	return s.viewPrefix + bucketKey
}

// BucketKeyHash returns the MD5 hash of the bucket storage key in which
// the item will be stored. The choice of MD5 is only for hash performance
// reasons since its value is not used for any security sensitive operation.
func (s *StoragePacker) BucketKeyHashByItemID(itemID string) string {
	return s.BucketKeyHashByKey(s.BucketPath(s.BucketKey(itemID)))
}

// BucketKeyHashByKey returns the MD5 hash of the bucket storage key
func (s *StoragePacker) BucketKeyHashByKey(bucketKey string) string {
	hf := md5.New()
	hf.Write([]byte(bucketKey))
	return hex.EncodeToString(hf.Sum(nil))
}

// View returns the storage view configured to be used by the packer
func (s *StoragePacker) View() logical.Storage {
	return s.view
}

// Get returns a bucket for a given key
func (s *StoragePacker) GetBucket(key string) (*Bucket, error) {
	if key == "" {
		return nil, fmt.Errorf("missing bucket key")
	}

	lock := locksutil.LockForKey(s.storageLocks, key)
	lock.RLock()
	defer lock.RUnlock()

	// Read from the underlying view
	storageEntry, err := s.view.Get(context.Background(), key)
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
func (s *StoragePacker) BucketIndex(key string) uint8 {
	hf := md5.New()
	hf.Write([]byte(key))
	return uint8(hf.Sum(nil)[0])
}

// BucketKey returns the bucket key for a given item ID
func (s *StoragePacker) BucketKey(itemID string) string {
	return strconv.Itoa(int(s.BucketIndex(itemID)))
}

// DeleteItem removes the storage entry which the given key refers to from its
// corresponding bucket.
func (s *StoragePacker) DeleteItem(itemID string) error {

	if itemID == "" {
		return fmt.Errorf("empty item ID")
	}

	// Get the bucket key
	bucketKey := s.BucketKey(itemID)

	// Prepend the view prefix
	bucketPath := s.BucketPath(bucketKey)

	// Read from underlying view
	storageEntry, err := s.view.Get(context.Background(), bucketPath)
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
func (s *StoragePacker) PutBucket(bucket *Bucket) error {
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
	err = s.view.Put(context.Background(), &logical.StorageEntry{
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
func (s *StoragePacker) GetItem(itemID string) (*Item, error) {
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
func (s *StoragePacker) PutItem(item *Item) error {
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
	storageEntry, err := s.view.Get(context.Background(), bucketPath)
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
