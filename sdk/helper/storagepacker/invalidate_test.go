package storagepacker

import (
	"reflect"
	"strings"
	"testing"

	"github.com/golang/protobuf/proto"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/helper/compressutil"
	"github.com/hashicorp/vault/sdk/logical"
)

func TestStoragePackerV2_MatchingStorage(t *testing.T) {
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}
	bucketStorageView := logical.NewStorageView(storage, "test/foo/buckets/")
	configStorageView := logical.NewStorageView(storage, "dontcare/config")
	logger := log.New(&log.LoggerOptions{Name: "storagepackertest"})

	sp, err := NewStoragePackerV2(ctx, &Config{
		BucketStorageView: bucketStorageView,
		ConfigStorageView: configStorageView,
		Logger:            logger,
	})

	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		Path     string
		Expected bool
	}{
		{"test/foo/buckets/v2/00", true},
		{"test/foo/buckets/00", false},
		{"test/foo/unrelated/v2/00", false},
		{"test/foo/buckets/v2/fa/c/e", true},
		{"test/foo/", false},
		{"test/bar/buckets/v2/00", false},
	}

	for _, tc := range cases {
		actual := sp.MatchingStorage(tc.Path)
		if actual != tc.Expected {
			t.Errorf("MatchingStorage returned wrong result on %q", tc.Path)
		}
	}
}

func TestStoragePackerV2_ParentBuckets(t *testing.T) {
	testCases := []struct {
		BucketsToAdd   []string
		StartBucketKey string
		ExpectedKeys   []string
	}{
		{[]string{"00/1", "00/2"},
			"00/1",
			[]string{"00", "00/1"}},
		{[]string{"00/1", "00/2"},
			"00",
			[]string{"00"}},
		{[]string{"00/0", "00/0/0"},
			"00/0",
			[]string{"00", "00/0"}},
		{[]string{"ab/c", "ab/c/d", "ab/c/d/e"},
			"ab/c/d/e",
			[]string{"ab", "ab/c", "ab/c/d", "ab/c/d/e"}},
	}

	for _, tc := range testCases {
		storagePacker := getStoragePacker(t)
		ctx := namespace.RootContext(nil)
		for _, bucketKey := range tc.BucketsToAdd {
			emptyBucket := &LockedBucket{
				Bucket:    createBucket(bucketKey, []*Item{}),
				LockEntry: nil,
			}
			cacheKey := storagePacker.GetCacheKey(bucketKey)
			storagePacker.bucketsCache.Insert(cacheKey, emptyBucket)
		}

		startBucket, err := storagePacker.GetBucket(ctx, tc.StartBucketKey)
		if err != nil {
			t.Fatal(err)
		}
		if startBucket.Key != tc.StartBucketKey {
			t.Fatalf("mismatched bucket key %q != %q", startBucket.Key, tc.StartBucketKey)
		}

		path := storagePacker.bucketAndAllParents(startBucket)
		resultKeys := make([]string, 0, len(tc.ExpectedKeys))
		for _, bucket := range path {
			resultKeys = append(resultKeys, bucket.Key)
		}
		if !reflect.DeepEqual(resultKeys, tc.ExpectedKeys) {
			t.Fatalf("mismatched bucket path %v != %v", resultKeys, tc.ExpectedKeys)
		}
	}
}

func marshalledBucket(bucket *Bucket) ([]byte, error) {
	marshaledBucket, err := proto.Marshal(bucket)
	if err != nil {
		return nil, err
	}

	compressedBucket, err := compressutil.Compress(marshaledBucket, &compressutil.CompressionConfig{
		Type: compressutil.CompressionTypeSnappy,
	})
	return compressedBucket, err
}

func createItemsForBucket(key string, n int) []*Item {
	cacheKey := strings.Replace(key, "/", "", -1)
	ids := generateLotsOfCollidingIDs(n, cacheKey)

	allItems := make([]*Item, n)
	for i, _ := range allItems {
		allItems[i] = &Item{
			ID:   ids[i],
			Data: incompressibleData(1000),
		}
	}
	return allItems
}

func createBucket(key string, items []*Item) *Bucket {
	b := &Bucket{
		Key:     key,
		ItemMap: make(map[string][]byte, len(items)),
	}
	for _, item := range items {
		b.ItemMap[item.ID] = item.Data
	}
	return b
}

func TestStoragePackerV2_InvalidateCreation(t *testing.T) {
	createAndCheck := func(key string) {
		// Create the []bytes version of a bucket by hand
		n := 10
		allItems := createItemsForBucket(key, n)
		b := createBucket(key, allItems)
		rawBucket, err := marshalledBucket(b)
		if err != nil {
			t.Fatal(err)
		}
		bucketPath := "packer/buckets/v2/" + key

		// Invalidate an empty storage packer to see if all
		// items are reported.
		storagePacker := getStoragePacker(t)
		ctx := namespace.RootContext(nil)
		present, deleted, err := storagePacker.InvalidateItems(ctx, bucketPath, rawBucket)
		if err != nil {
			t.Fatal(err)
		}
		if len(present) != n {
			t.Fatalf("%d elements  present", len(present))
		}
		if len(deleted) != 0 {
			t.Fatalf("%d elements deleted", len(deleted))
		}

		// Check present items
		checkReturnedItems(t, allItems, present)
		// Also check that SP has been modified
		checkAllItems(t, storagePacker, ctx, allItems)
	}

	t.Logf("invalidating existing bucket 00")
	createAndCheck("00")

	t.Logf("invalidating new shard 00/1")
	createAndCheck("00/1")

}

func TestStoragePackerV2_InvalidateItemDeletion(t *testing.T) {
	n := 10
	allItems := createItemsForBucket("00", n)
	b := createBucket("00", allItems[:n-2])
	rawBucket, err := marshalledBucket(b)
	if err != nil {
		t.Fatal(err)
	}
	bucketPath := "packer/buckets/v2/00"

	// Create a storage packer with n of the items, ensure
	// that the last two are reported as deleted and the first
	// is reported with its new value.
	storagePacker := getStoragePacker(t)
	ctx := namespace.RootContext(nil)
	storagePacker.PutItem(ctx, allItems...)
	storagePacker.PutItem(ctx, &Item{
		ID:   allItems[0].ID,
		Data: incompressibleData(1000),
	})

	present, deleted, err := storagePacker.InvalidateItems(ctx, bucketPath, rawBucket)
	if err != nil {
		t.Fatal(err)
	}

	checkReturnedItems(t, allItems[:n-2], present)
	checkReturnedItems(t, allItems[n-2:], deleted)

	checkAllItems(t, storagePacker, ctx, allItems[:n-2])
	lookup, err := storagePacker.GetItems(ctx, allItems[n-2].ID, allItems[n-1].ID)
	if err != nil {
		t.Fatal(err)
	}
	if lookup[0] != nil || lookup[1] != nil {
		t.Fatalf("deleted values are still reported from GetItems")
	}
}

func TestStoragePackerV2_InvalidateInvalidBucket(t *testing.T) {
	bucketPath00 := "packer/buckets/v2/00"
	garbageBucket := []byte("Not a compressed bucket, not even close.")

	storagePacker := getStoragePacker(t)
	ctx := namespace.RootContext(nil)

	// Put some items in so we can verify it's unchanged.
	allItems := createItemsForBucket("00", 5)
	err := storagePacker.PutItem(ctx, allItems...)
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = storagePacker.InvalidateItems(ctx, bucketPath00, garbageBucket)
	if err == nil {
		t.Fatalf("Expected an error.")
	}

	_, _, err = storagePacker.InvalidateItems(ctx, "not a bucket", nil)
	if err == nil {
		t.Fatalf("Expected an error.")
	}

	// Verify everything still there.
	checkAllItems(t, storagePacker, ctx, allItems)
}

func TestStoragePackerV2_InvalidateShardedBucket(t *testing.T) {
	// Set up
	//
	// 00 [old contents]
	//
	// 00 [sharded]
	//   00/1 [old contents]
	//   00/2 [old contents]
	//   00/3 [empty]
	// and run through updates in any order

	n := 10
	items001 := createItemsForBucket("001", n)
	items002 := createItemsForBucket("002", n)
	allItems := append(items001, items002...)

	// Empty version post-sharding
	b00prime := createBucket("00", []*Item{})
	b00prime.HasShards = true
	rawBucket00prime, err := marshalledBucket(b00prime)
	if err != nil {
		t.Fatal(err)
	}
	bucketPath00 := "packer/buckets/v2/00"

	// Sharded buckets
	b001 := createBucket("00/1", items001)
	rawBucket001, err := marshalledBucket(b001)
	if err != nil {
		t.Fatal(err)
	}
	bucketPath001 := "packer/buckets/v2/00/1"

	b002 := createBucket("00/2", items002)
	rawBucket002, err := marshalledBucket(b002)
	if err != nil {
		t.Fatal(err)
	}
	bucketPath002 := "packer/buckets/v2/00/2"

	empty := []*Item{}
	b003 := createBucket("00/3", empty)
	rawBucket003, err := marshalledBucket(b003)
	if err != nil {
		t.Fatal(err)
	}
	bucketPath003 := "packer/buckets/v2/00/3"

	testCases := []struct {
		Name            string
		Bucket          []byte
		Path            string
		PresentExpected []*Item
		DeletedExpected []*Item
		Finished        bool // Check all items and reset after this
	}{
		// Normal order under WAL operation, note that in step 4
		// all items get reported for a second time (because they're deleted.)
		{"normal order step 1", rawBucket001, bucketPath001, items001, empty, false},
		{"normal order step 2", rawBucket002, bucketPath002, items002, empty, false},
		{"normal order step 3", rawBucket003, bucketPath003, empty, empty, false},
		{"normal order step 4", rawBucket00prime, bucketPath00, allItems, empty, true},

		// Reverse of normal order, all items temporarily listed as deleted
		{"reversed step 1", rawBucket00prime, bucketPath00, empty, allItems, false},
		{"reversed step 2", rawBucket003, bucketPath003, empty, empty, false},
		{"reversed step 3", rawBucket002, bucketPath002, items002, empty, false},
		{"reversed step 4", rawBucket001, bucketPath001, items001, empty, true},

		// Parent in the middle
		{"random step 1", rawBucket002, bucketPath002, items002, empty, false},
		{"random step 2", rawBucket00prime, bucketPath00, items002, items001, false},
		{"random step 3", rawBucket003, bucketPath003, empty, empty, false},
		{"random step 4", rawBucket001, bucketPath001, items001, empty, true},
	}

	//
	// shard, then parent (normal order)
	//
	// Create an empty storage packer and fill it via normal operation
	storagePacker := getStoragePacker(t)
	ctx := namespace.RootContext(nil)
	err = storagePacker.PutItem(ctx, allItems...)
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range testCases {
		t.Logf("step %q", tc.Name)

		present, deleted, err := storagePacker.InvalidateItems(ctx, tc.Path, tc.Bucket)
		if err != nil {
			t.Fatal(err)
		}
		checkReturnedItems(t, tc.PresentExpected, present)
		checkReturnedItems(t, tc.DeletedExpected, deleted)

		if tc.Finished {
			checkAllItems(t, storagePacker, ctx, allItems)
			storagePacker = getStoragePacker(t)
			err = storagePacker.PutItem(ctx, allItems...)
			if err != nil {
				t.Fatal(err)
			}
		}
	}
}

func TestStoragePackerV2_InvalidateDeletedBucket(t *testing.T) {
	// Set up
	// 00 [old contents]
	//   00/1 [new contents, for a subset of items]
	// Then delete 00/1 and ensure we get back to just 00.

	// items001[0] is new
	// items001[1] and [2] are modified
	// items001[3:] are unchanged

	n := 10
	items001 := createItemsForBucket("001", n+1)
	items002 := createItemsForBucket("002", n)
	allItems := append(items001[1:], items002...)

	items001[1].Data = []byte("version 2")
	items001[2].Data = []byte("version 2")

	// Base bucket, contains 2n items
	b00 := createBucket("00", allItems)
	rawBucket00, err := marshalledBucket(b00)
	if err != nil {
		t.Fatal(err)
	}
	bucketPath00 := "packer/buckets/v2/00"

	// Sharded bucket, contains some updated items
	b001 := createBucket("00/1", items001)
	rawBucket001, err := marshalledBucket(b001)
	if err != nil {
		t.Fatal(err)
	}
	bucketPath001 := "packer/buckets/v2/00/1"

	// Create an empty storage packer and fill it via in validation.
	storagePacker := getStoragePacker(t)
	ctx := namespace.RootContext(nil)

	present, deleted, err := storagePacker.InvalidateItems(ctx, bucketPath00, rawBucket00)
	if err != nil {
		t.Fatal(err)
	}

	present, deleted, err = storagePacker.InvalidateItems(ctx, bucketPath001, rawBucket001)
	if err != nil {
		t.Fatal(err)
	}

	// Verify everything present, as expected, with the version 2 objects
	checkAllItems(t, storagePacker, ctx, items001[1:])
	checkAllItems(t, storagePacker, ctx, items002)

	// Invalidate 00/1 to nil (which means it has been deleted on the primary)
	present, deleted, err = storagePacker.InvalidateItems(ctx, bucketPath001, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(present) != n {
		t.Fatalf("%d elements present", len(present))
	}
	if len(deleted) != 1 {
		t.Fatalf("%d elements deleted", len(deleted))
	}
	checkReturnedItems(t, allItems[:n], present) // version 1 objects
	checkReturnedItems(t, items001[:1], deleted) // the only new object in bucket 00/1
}
