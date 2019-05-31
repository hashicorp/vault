package storagepacker

import (
	"fmt"
	"github.com/hashicorp/vault/sdk/helper/locksutil"
	"math"
	"testing"
	"testing/quick"
)

// Is every element of the partition a different bucket?
func partitionBucketsAreUnique(t *testing.T, partition []*partitionedRequests) bool {
	bucketsSeen := make(map[*LockedBucket]bool, len(partition))
	for _, p := range partition {
		if _, found := bucketsSeen[p.Bucket]; found {
			t.Logf("non-unique bucket found")
			return false
		}
		bucketsSeen[p.Bucket] = true
	}
	return true
}

// Are the buckets in a partition arranged in increasing order?
func partitionKeysInOrder(t *testing.T, partition []*partitionedRequests) bool {
	for i, p := range partition {
		if i > 0 {
			prevP := partition[i-1]
			if p.Bucket.Key <= prevP.Bucket.Key {
				t.Logf("buckets not ordered in partition: %v <= %v", p.Bucket.Key, prevP.Bucket.Key)
				return false
			}
		}
	}
	return true
}

// Is every key present in the partition?
func partitionHasAllItems(t *testing.T, partition []*partitionedRequests, ids []string) bool {
	idsRequired := make(map[string]bool, len(ids))
	for _, id := range ids {
		idsRequired[id] = true
	}

	for _, p := range partition {
		if p == nil || p.Requests == nil || len(p.Requests) == 0 {
			t.Logf("nil or empty partition")
			return false
		}
		for _, r := range p.Requests {
			_, present := idsRequired[r.ID]
			if !present {
				t.Logf("extra or duplicated ID: '%v'", r.ID)
				return false
			}
			delete(idsRequired, r.ID)
		}
	}
	if len(idsRequired) != 0 {
		t.Logf("IDs not found in partition")
		return false
	}
	return true
}

// FIXME: this will move into the constructor
func insertAllBuckets(s *StoragePackerV2) {
	numBuckets := int(math.Pow(2.0, float64(s.BaseBucketBits)))
	for i := 0; i < numBuckets; i++ {
		bucketKey := fmt.Sprintf("%0x", i)
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

func TestPartitionProperties(t *testing.T) {
	checkIds := func(ids []string) bool {
		// Higher-level function should probably check this.
		if dup, _ := checkForDuplicateIds(ids); dup {
			return true
		}

		s := getStoragePacker(t)
		insertAllBuckets(s)
		requests := s.keysForIDs(ids)
		partition, err := s.partitionRequests(requests)
		if err != nil {
			t.Logf("error in partitionRequests: %v", err)
			return false
		}
		if partitionBucketsAreUnique(t, partition) &&
			partitionKeysInOrder(t, partition) &&
			partitionHasAllItems(t, partition, ids) {
			retry := s.lockBuckets(partition, true)
			if retry {
				t.Logf("shareded bucket found")
				return false
			}
			s.unlockBuckets(partition, true)
			return true
		} else {
			return false
		}
	}
	// Highly artificial use case
	ids01 := generateLotsOfCollidingIDs(100, "01")
	if !checkIds(ids01) {
		t.Error("Failed colliding IDs test.")
	}

	// Random testing
	if err := quick.Check(checkIds, nil); err != nil {
		t.Error(err)
	}
}
