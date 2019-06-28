package storagepacker

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/helper/locksutil"
	"github.com/hashicorp/vault/sdk/logical"
)

// MatchingStorage returns true if the given storage path is
// controlled by this StoragePacker.
func (s *StoragePackerV2) MatchingStorage(path string) bool {
	return strings.HasPrefix(path, s.BucketStorageView.Prefix())
}

func (s *StoragePackerV2) bucketKeyFromPath(path string) (string, error) {
	bucketRoot := s.BucketStorageView.Prefix()

	if strings.HasPrefix(path, bucketRoot) {
		return strings.TrimPrefix(path, bucketRoot), nil
	} else {
		return "", fmt.Errorf("path %q doesn't match bucket storage %q", path, bucketRoot)
	}
}

// InvalidateItems handles invalidation of a StoragePackerV2 bucket during replication.
//
// Given a storage path (to a bucket) and new contents, return
// "present": all items present in the new bucket (new, changed, or unchanged)
// "deleted": all items that were in the cache but are now absent
//
// For normal WAL operation calls to this function are serialized based on the order PutItem
// wrote the buckets; for recovery from the Merkle tree this will not be the case.
//
// For each bucket B we receive, process the items in it, compared with items in
// the previous version of B in the cache (if any, call it C) and elsewhere in the current cache.
//
// If an item exists in B, report it in "present" even if it also exists elsewere in the cache
//
// If an item exists in C, but not in B, then it might have been deleted.
//    * If the item exists elsewhere in the cache (might be a child, might even be a parent)
//      don't mark it as deleted. Either we will do it later (when we process that bucket)
//      or the item is now in its correct location.
//         * One problem here is that we can't rely upon the longest suffix match to find the bucket.
//         * Another problem is that the local side might have a newer (or older) version in its
//           memory compared to storage. If we can't be confident that InvalidateItems() will be
//           later called on the bucket that *does* have the consistent version, we have to report any
//           such item as "present" to give the upper layer a chance to look at it.
//    * If the item does not exist elsewhere in the cache, report it as deleted.
//
// If an item has been moved by sharding, it's possible this will cause it to be reported
// as deleted and re-created.
//
// Either B or C might be nonexistent; just treat those as empty.
//
// If we're the secondary, we shouldn't have to worry about concurrent access.
func (s *StoragePackerV2) InvalidateItems(ctx context.Context, path string, newValue []byte) (present []*Item, deleted []*Item, err error) {
	// Double-check that it's on our path, and figure out which bucket it should be
	bucketKey, err := s.bucketKeyFromPath(path)
	if err != nil {
		return nil, nil, err
	}

	var replacementBucket *LockedBucket
	bucketRemoved := false

	if newValue == nil {
		bucketRemoved = true
		// Create a placeholder so we can run the rest of the logic
		// against an empty collection
		replacementBucket = s.newEmptyBucket(bucketKey)
	} else {
		// Fake up a storage entry for the bucket and decode it
		storage := &logical.StorageEntry{
			Value: newValue,
		}
		replacementBucket, err = s.DecodeBucket(storage)
		if err != nil {
			return nil, nil, errwrap.Wrapf("unable to decode replicated bucket: {{err}}", err)
		}
		if replacementBucket.Key != bucketKey {
			return nil, nil, fmt.Errorf("decoded bucket key %q doesn't match path component %q", replacementBucket.Key, bucketKey)
		}
	}

	present = make([]*Item, 0)

	originalBucket, err := s.GetBucket(ctx, bucketKey)
	if err != nil {
		return nil, nil, errwrap.Wrapf("problem finding original bucket: {{err}}", err)
	}
	originalBucket.Lock()
	defer originalBucket.Unlock()

	// Any item included in the new bucket should be reported present.
	if !replacementBucket.HasShards && replacementBucket.ItemMap != nil {
		for k, v := range replacementBucket.ItemMap {
			present = append(present, &Item{
				ID:   k,
				Data: v,
			})
		}
	}

	// Is the new bucket a replacement, or new?  If it's new, nothing to compare against.
	if originalBucket.Key == replacementBucket.Key {
		// Any item not present in the new bucket *might* be deleted, or
		// it might be elsewhere (in which case we want to report that value instead,
		// to make sure the correct version is being used.)
		maybeDeleted := itemDifferenceBetween(originalBucket, replacementBucket)
		var revisit []*Item
		deleted, revisit, err = s.identifyItemsAbsentOrShadowed(bucketKey, maybeDeleted)
		if err != nil {
			// Return partial information?
			return present, nil, err
		}
		present = append(present, revisit...)
	}

	// Swap the replacement bucket in to the cache (or delete)
	cacheKey := s.GetCacheKey(bucketKey)
	s.bucketsCacheLock.Lock()
	if bucketRemoved {
		s.bucketsCache.Delete(cacheKey)
	} else {
		s.bucketsCache.Insert(cacheKey, replacementBucket)
	}
	s.bucketsCacheLock.Unlock()

	return present, deleted, nil

}

// Look up each item and determine whether another copy (or version) of it exists
// in a different bucket.
func (s *StoragePackerV2) identifyItemsAbsentOrShadowed(bucketKey string, maybeDeleted []*itemRequest) ([]*Item, []*Item, error) {
	deleted := make([]*Item, 0)
	revisit := make([]*Item, 0)

	// Using the radix lookup will find child buckets, but not parent
	// buckets.  For that we'll have to walk up the tree a step at a time.
	//
	// For example, we might be processing 00/0
	//
	// 00 -- 00/0 -- 00/0/a
	//
	// but, the item with key 000a... might exist in 00/0/a or maybe even in 00,
	// so we want to check both.
	//
	// If we find it somewhere, return the version in the bucket with the longest key.
	// (That is more likely to be a more recent version.)
	partition, err := s.partitionRequests(maybeDeleted)
	if err != nil {
		return nil, nil, errwrap.Wrapf("problem searching for deleted keys: {{err}}", err)
	}

	// Looking through the other buckets without acquiring a lock
	// is unsafe, but should be OK in the context of a replica.
	for _, p := range partition {
		bucketsToCheck := s.bucketAndAllParents(p.Bucket)
		for _, request := range p.Requests {
			found := false
			for _, b := range bucketsToCheck {
				// Skip originalBucket
				if b.Key != bucketKey {
					var data []byte
					if data, found = b.ItemMap[request.ID]; found {
						revisit = append(revisit, &Item{
							ID:   request.ID,
							Data: data,
						})
						break
					}
				}
			}
			if !found {
				// Include the original value
				// though maybe this is overkill and we should just return IDs.
				deleted = append(deleted, request.Value)
			}
		}
	}

	return deleted, revisit, nil

}

// Assemble a list (from longest key to shortest) of all the parent buckets.
func (s *StoragePackerV2) bucketAndAllParents(bucket *LockedBucket) []*LockedBucket {
	s.bucketsCacheLock.RLock()
	defer s.bucketsCacheLock.RUnlock()

	parents := []*LockedBucket{bucket}
	for k := s.parentCacheKey(bucket.Key); k != ""; k = s.parentCacheKey(k) {
		bucketRaw, found := s.bucketsCache.Get(k)
		if found {
			parents = append(parents, bucketRaw.(*LockedBucket))
		}
	}
	return parents
}

// Given a bucket key, return the cache key for its parent
func (s *StoragePackerV2) parentCacheKey(bucketKey string) string {
	bucketCacheKey := s.GetCacheKey(bucketKey)
	if len(bucketCacheKey) < s.BaseBucketBits/4 {
		return ""
	} else {
		n := len(bucketCacheKey)
		shardLen := s.BucketShardBits / 4
		return bucketCacheKey[:(n - shardLen)]
	}
}

// Create a placeholder empty bucket given its key
func (s *StoragePackerV2) newEmptyBucket(bucketKey string) *LockedBucket {
	cacheKey := s.GetCacheKey(bucketKey)
	lock := locksutil.LockForKey(s.storageLocks, cacheKey)
	return &LockedBucket{
		Bucket: &Bucket{
			Key:       bucketKey,
			HasShards: false,
			ItemMap:   make(map[string][]byte),
		},
		LockEntry: lock,
	}
}

func itemDifferenceBetween(originalBucket *LockedBucket, replacementBucket *LockedBucket) []*itemRequest {
	maybeDeleted := make([]*itemRequest, 0)

	add := func(id string, val []byte) {
		maybeDeleted = append(maybeDeleted, &itemRequest{
			ID:  id,
			Key: GetItemIDHash(id),
			Value: &Item{
				ID:   id,
				Data: val,
			},
		})
	}

	if !originalBucket.HasShards && originalBucket.ItemMap != nil {
		// Candidates for deleted list
		if replacementBucket.HasShards || replacementBucket.ItemMap == nil {
			// Everything!
			for id, v := range originalBucket.ItemMap {
				add(id, v)
			}
		} else {
			// Only those in set difference
			for id, v := range originalBucket.ItemMap {
				if _, found := replacementBucket.ItemMap[id]; !found {
					add(id, v)
				}
			}
		}
	}

	return maybeDeleted
}

// Scenarios that make things complicated:
//
//
// 1) Primary had to roll back a sharding operation due to failure.
// As a result, secondary cache has some shards but the primary does not.
//
// Secondary: 00 (full, pre-sharding), 00/0, 00/a, 00/c
// Primary:   00 only
//
// We could get any of these in any order: delete 00/0, delete 00/a, delete 00/c, invalidate 00
//
// 1b) We might not get invalidate 00 *at all* if the sharding was rolled back but
// no change was made to the parent.
//    00/0: has item version 2, that's what's in-memory on the secondary
//    00: has item version 1, that's what the primary recovers to
//    at some point we need the objects in 00/0 to be reported as present, so that we use v1,
//    even if bucket 00 isn't updated.
//
// 2) Primary sharded (maybe even to multiple levels) while we were disconnected.
// As a result, we start getting leaves before the parent instead of
// in the order shardBucket wrote.
//
// Secondary: 00
// Primary:   00 (empty), 00/0 (empty), 00/0/1, 00/0/2, ..., 00/0/f, 00/1, ..., 00/f
//
// 3) both combined, not necessarily the same bucket that got sharded
//
// Secondary: 00 (full, pre-sharding), 00/0, 00/a, 00/c
// Primary:   00 (empty), 00/0 (full), 00/1 (empty), 00/1/0, 00/1/1, ..., 00/1/f, 00/2, ..., 00/f
//
// 4) StoragePacker completely deleted and then re-created.
// Unlikely for core structures like mount table, might be possible for other
// uses in the future.
