package storagepacker

import (
	"context"
	"crypto/md5"
	"fmt"
	"strconv"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/errwrap"
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
func (s *StoragePacker) GetBucket(key string) (*Bucket, error) {
	if key == "" {
		return nil, fmt.Errorf("missing bucket key")
	}

	lock := locksutil.LockForKey(s.storageLocks, key)
	lock.RLock()
	defer lock.RUnlock()

	// Read from storage
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
func (s *StoragePacker) DeleteItem(_ context.Context, itemID string) error {
	return s.DeleteMultipleItems(context.Background(), nil, itemID)
}

func (s *StoragePacker) DeleteMultipleItems(ctx context.Context, logger hclog.Logger, itemIDs ...string) error {
	var err error
	switch len(itemIDs) {
	case 0:
		// Nothing
		return nil

	case 1:
		logger = hclog.NewNullLogger()
		fallthrough

	default:
		lockIndexes := make(map[string]struct{}, len(s.storageLocks))
		for _, itemID := range itemIDs {
			bucketKey := s.BucketKey(itemID)
			if _, ok := lockIndexes[bucketKey]; !ok {
				lockIndexes[bucketKey] = struct{}{}
			}
		}

		lockKeys := make([]string, 0, len(lockIndexes))
		for k := range lockIndexes {
			lockKeys = append(lockKeys, k)
		}

		locks := locksutil.LocksForKeys(s.storageLocks, lockKeys)
		for _, lock := range locks {
			lock.Lock()
			defer lock.Unlock()
		}
	}

	if logger == nil {
		logger = hclog.NewNullLogger()
	}

	bucketCache := make(map[string]*Bucket, len(s.storageLocks))

	logger.Debug("deleting multiple items from storagepacker; caching and deleting from buckets", "total_items", len(itemIDs))

	var pctDone int
	for idx, itemID := range itemIDs {
		bucketKey := s.BucketKey(itemID)

		bucket, bucketFound := bucketCache[bucketKey]
		if !bucketFound {
			// Read from storage
			storageEntry, err := s.view.Get(context.Background(), bucketKey)
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

			bucket = new(Bucket)
			err = proto.Unmarshal(uncompressedData, bucket)
			if err != nil {
				return errwrap.Wrapf("failed decoding packed storage entry: {{err}}", err)
			}
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
			bucket.Items[foundIdx] = bucket.Items[len(bucket.Items)-1]
			bucket.Items = bucket.Items[:len(bucket.Items)-1]
			if !bucketFound {
				bucketCache[bucketKey] = bucket
			}
		}

		newPctDone := idx * 100.0 / len(itemIDs)
		if int(newPctDone) > pctDone {
			pctDone = int(newPctDone)
			logger.Trace("bucket item removal progress", "percent", pctDone, "items_removed", idx)
		}
	}

	logger.Debug("persisting buckets", "total_buckets", len(bucketCache))

	// Persist all buckets in the cache; these will be the ones that had
	// deletions
	pctDone = 0
	idx := 0
	for _, bucket := range bucketCache {
		// Fail if the context is canceled, the storage calls will fail anyways
		if ctx.Err() != nil {
			return ctx.Err()
		}

		err = s.putBucket(ctx, bucket)
		if err != nil {
			return err
		}

		newPctDone := idx * 100.0 / len(bucketCache)
		if int(newPctDone) > pctDone {
			pctDone = int(newPctDone)
			logger.Trace("bucket persistence progress", "percent", pctDone, "buckets_persisted", idx)
		}

		idx++
	}

	return nil
}

func (s *StoragePacker) putBucket(ctx context.Context, bucket *Bucket) error {
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
func (s *StoragePacker) GetItem(itemID string) (*Item, error) {
	if itemID == "" {
		return nil, fmt.Errorf("empty item ID")
	}

	bucketKey := s.BucketKey(itemID)

	// Fetch the bucket entry
	bucket, err := s.GetBucket(bucketKey)
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

// PutItem stores the given item in its respective bucket
func (s *StoragePacker) PutItem(_ context.Context, item *Item) error {
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
	storageEntry, err := s.view.Get(context.Background(), bucketKey)
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

	return s.putBucket(context.Background(), bucket)
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
