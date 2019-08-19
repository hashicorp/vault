package storagepacker

import (
	"fmt"
	"github.com/hashicorp/vault/sdk/helper/locksutil"
)

// Requests for multi-item put or get, partitioned by bucket.
// This lets us operate bucket-by-bucket.
type partitionedRequests struct {
	Bucket   *LockedBucket
	Requests []*Item
}

// partitionRequests takes a list of requests and partition them by which
// bucket the Key belongs to.
func (s *StoragePackerV2) partitionRequests(unsortedRequests []*Item) ([]*partitionedRequests, error) {
	sortedRequests := sortRequests(unsortedRequests)

	partition := make([]*partitionedRequests, 0)
	var lastBucket *partitionedRequests

	// The items requests are sorted by key, so the same buckets end
	// up together. But, the radix tree doesn't have a way to take
	// advantage of that to do a single pass.
	//
	// We could check whether each request has the same prefix as the
	// previous one, which would technically work given that we fully
	// populate the shards, so we don't have to worry about a set of
	// radix tree entries that looks like X, X0, X2, X4 where we pop
	// up and down between two levels.
	s.bucketsCacheLock.RLock()
	defer s.bucketsCacheLock.RUnlock()

	for _, r := range sortedRequests {
		_, bucketRaw, found := s.bucketsCache.LongestPrefix(r.key)
		if !found {
			return nil, fmt.Errorf("key %s not found in bucket cache", r.key)
		}
		bucket := bucketRaw.(*LockedBucket)
		if lastBucket == nil || lastBucket.Bucket != bucket {
			// Distinct from previous bucket
			lastBucket = &partitionedRequests{
				Bucket:   bucket,
				Requests: []*Item{r},
			}
			partition = append(partition, lastBucket)
		} else {
			// Same as previous bucket
			lastBucket.Requests = append(lastBucket.Requests, r)
		}
	}
	return partition, nil
}

// lockBuckets acquires the locks for all the identified buckets, in order.
// Check that the buckets are still unsharded after the lock is acquired;
// if not, the partitioning step must be retried, but this is an uncommon
// operation.
func (s *StoragePackerV2) lockBuckets(partition []*partitionedRequests, read bool) (retryRequired bool) {
	// We have the locks already, as part of the LockedBucket structure.
	// There's no easy way to map back from that lock to its order within
	// storageLocks.
	//
	// For future work: Is there a benefit to fast-pathing the case of
	// just one bucket?

	// Which locks are requested?
	lockNeeded := make(map[*locksutil.LockEntry]bool, len(partition))
	for _, p := range partition {
		lockNeeded[p.Bucket.LockEntry] = true
	}

	// Lock them in order
	for _, l := range s.storageLocks {
		if _, ok := lockNeeded[l]; ok {
			if read {
				l.RLock()
			} else {
				l.Lock()
			}
		}
	}

	retryRequired = false
	// Check that the buckets are still leaf nodes.
	// A sharding operation may have occurred between releasing the
	// radix tree mutex and acquiring the storage locks.
	for _, p := range partition {
		if p.Bucket.HasShards {
			retryRequired = true
			break
		}
	}

	if retryRequired {
		s.unlockBuckets(partition, read)
	}
	return
}

func (s *StoragePackerV2) unlockBuckets(partition []*partitionedRequests, read bool) {
	// Unlock only once, no need to do it in order.
	lockFreed := make(map[*locksutil.LockEntry]bool, len(partition))
	for _, p := range partition {
		if _, ok := lockFreed[p.Bucket.LockEntry]; !ok {
			if read {
				p.Bucket.RUnlock()
			} else {
				p.Bucket.Unlock()
			}
			lockFreed[p.Bucket.LockEntry] = true
		}
	}
}
