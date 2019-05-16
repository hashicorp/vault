package storagepacker

import (
	"context"
	"testing"

	"bytes"
	"github.com/go-test/deep"
	"github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	log "github.com/hashicorp/go-hclog"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
	"math/rand"
	"strings"
)

func createStoragePacker(tb testing.TB, storage logical.Storage) *StoragePackerV2 {
	storageView := logical.NewStorageView(storage, "packer/buckets/v2")
	storagePacker, err := NewStoragePackerV2(context.Background(), &Config{
		BucketStorageView: storageView,
		ConfigStorageView: logical.NewStorageView(storage, "packer/config"),
		Logger:            log.New(&log.LoggerOptions{Name: "storagepackertest"}),
	})
	if err != nil {
		tb.Fatal(err)
	}
	return storagePacker.(*StoragePackerV2)
}

func getStoragePacker(tb testing.TB) *StoragePackerV2 {
	storage := &logical.InmemStorage{}
	return createStoragePacker(tb, storage)
}

func BenchmarkStoragePackerV2(b *testing.B) {
	storagePacker := getStoragePacker(b)

	ctx := namespace.RootContext(nil)

	ctx = context.Background()

	for i := 0; i < b.N; i++ {
		itemID, err := uuid.GenerateUUID()
		if err != nil {
			b.Fatal(err)
		}

		item := &Item{
			ID: itemID,
		}

		err = storagePacker.PutItem(ctx, item)
		if err != nil {
			b.Fatal(err)
		}

		fetchedItem, err := storagePacker.GetItem(ctx, itemID)
		if err != nil {
			b.Fatal(err)
		}

		if fetchedItem == nil {
			b.Fatalf("failed to read stored item with ID: %q, iteration: %d", item.ID, i)
		}

		if fetchedItem.ID != item.ID {
			b.Fatalf("bad: item ID; expected: %q\n actual: %q", item.ID, fetchedItem.ID)
		}

		err = storagePacker.DeleteItem(ctx, item.ID)
		if err != nil {
			b.Fatal(err)
		}

		fetchedItem, err = storagePacker.GetItem(ctx, item.ID)
		if err != nil {
			b.Fatal(err)
		}
		if fetchedItem != nil {
			b.Fatalf("failed to delete item")
		}
	}
}

func TestStoragePackerV2(t *testing.T) {
	storagePacker := getStoragePacker(t)

	// Persist a storage entry
	item1 := &Item{
		ID:   "item1",
		Data: []byte("data1"),
	}

	ctx := namespace.RootContext(nil)

	err := storagePacker.PutItem(ctx, item1)
	if err != nil {
		t.Fatal(err)
	}

	// Verify that it can be read
	fetchedItem, err := storagePacker.GetItem(ctx, item1.ID)
	if err != nil {
		t.Fatal(err)
	}
	if fetchedItem == nil {
		t.Fatalf("failed to read the stored item")
	}

	if item1.ID != fetchedItem.ID {
		t.Fatalf("bad: item ID; expected: %q\n actual: %q\n", item1.ID, fetchedItem.ID)
	}

	// Delete item1
	err = storagePacker.DeleteItem(ctx, item1.ID)
	if err != nil {
		t.Fatal(err)
	}

	// Check that the deletion was successful
	fetchedItem, err = storagePacker.GetItem(ctx, item1.ID)
	if err != nil {
		t.Fatal(err)
	}

	if fetchedItem != nil {
		t.Fatalf("failed to delete item")
	}
}

func TestStoragePackerV2_SerializeDeserializeComplexItem_Version1(t *testing.T) {
	storagePacker := getStoragePacker(t)

	ctx := context.Background()

	timeNow := ptypes.TimestampNow()

	alias1 := &identity.Alias{
		ID:            "alias_id",
		CanonicalID:   "canonical_id",
		MountType:     "mount_type",
		MountAccessor: "mount_accessor",
		Metadata: map[string]string{
			"aliasmkey": "aliasmvalue",
		},
		Name:                   "alias_name",
		CreationTime:           timeNow,
		LastUpdateTime:         timeNow,
		MergedFromCanonicalIDs: []string{"merged_from_canonical_id"},
	}

	entity := &identity.Entity{
		Aliases: []*identity.Alias{alias1},
		ID:      "entity_id",
		Name:    "entity_name",
		Metadata: map[string]string{
			"testkey1": "testvalue1",
			"testkey2": "testvalue2",
		},
		CreationTime:    timeNow,
		LastUpdateTime:  timeNow,
		BucketKey:       "entity_hash",
		MergedEntityIDs: []string{"merged_entity_id1", "merged_entity_id2"},
		Policies:        []string{"policy1", "policy2"},
	}

	marshaledBytes, err := proto.Marshal(entity)
	if err != nil {
		t.Fatal(err)
	}

	ctx = namespace.RootContext(nil)
	err = storagePacker.PutItem(ctx, &Item{
		ID:   entity.ID,
		Data: marshaledBytes,
	})
	if err != nil {
		t.Fatal(err)
	}

	itemFetched, err := storagePacker.GetItem(ctx, entity.ID)
	if err != nil {
		t.Fatal(err)
	}

	var itemDecoded identity.Entity
	err = proto.Unmarshal(itemFetched.Data, &itemDecoded)
	if err != nil {
		t.Fatal(err)
	}

	if !proto.Equal(&itemDecoded, entity) {
		diff := deep.Equal(&itemDecoded, entity)
		t.Fatal(diff)
	}
}

func TestStoragePackerV2_Recovery(t *testing.T) {
	storage := &logical.InmemStorage{}
	storagePacker1 := createStoragePacker(t, storage)

	// Persist a storage entry
	item1 := &Item{
		ID:   "item1",
		Data: []byte("data1"),
	}

	ctx := namespace.RootContext(nil)
	err := storagePacker1.PutItem(ctx, item1)
	if err != nil {
		t.Fatal(err)
	}

	// Verify that it can be read back from a second packer backed by the same storage
	storagePacker2 := createStoragePacker(t, storage)

	fetchedItem, err := storagePacker2.GetItem(ctx, item1.ID)
	if err != nil {
		t.Fatal(err)
	}
	if fetchedItem == nil {
		t.Fatalf("failed to read the stored item")
	}
	if item1.ID != fetchedItem.ID {
		t.Fatalf("bad: item ID; expected: %q\n actual: %q\n", item1.ID, fetchedItem.ID)
	}
}

const idLetters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func generateRandomID(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = idLetters[rand.Intn(len(idLetters))]
	}
	return string(b)
}

func generateLotsOfCollidingIDs(numDesired int, prefix string) []string {
	count := 0
	ids := make([]string, 0, numDesired)
	for count < numDesired {
		// We might as well use the desired output prefix as
		// an input prefix as well, so that they can be visually
		// distinguished.
		id := prefix + generateRandomID(40-len(prefix))
		key := GetItemIDHash(id)
		if strings.HasPrefix(key, prefix) {
			ids = append(ids, id)
			count += 1
		}
	}
	return ids
}

func incompressibleData(length int) (out []byte) {
	out = make([]byte, length)
	rand.Read(out)
	return
}

func TestStoragePackerV2_RecoveryAfterCrash(t *testing.T) {
	// Set up the following scenario:
	//    Bucket  00 has to be split
	//    Buckets 00/0, 00/1, 00/2 have been written by the splitter, then we crash
	// We can do this by saving bucket 00, then replacing it and deleting
	// storage items 00/3 .. 00/F
	//
	// Recover and verify that data with keys 000... and 003... are still accessible.
	//
	// Force bucket 00 to be split again
	// Verify old data with keys 000... and 003.. are still accessible.
	//

	// FIXME: compile these in, intead of expensively generating them
	// each time?
	ids := [][]string{
		generateLotsOfCollidingIDs(10, "000"),
		generateLotsOfCollidingIDs(10, "001"),
		generateLotsOfCollidingIDs(10, "003"),
	}

	// We'll use 1000 for our size and 10000 for the maximum, that
	// should mean we can store 9 but the 10th overflows.
	storage := logical.InmemStorageWithMaxSize(10000)
	storagePacker1 := createStoragePacker(t, storage)

	allItems := make([]*Item, 30)
	for i, _ := range allItems {
		allItems[i] = &Item{
			ID:   ids[i%3][i/3],
			Data: incompressibleData(1000),
		}
	}

	// Insert first 9 items, no split yet
	ctx := namespace.RootContext(nil)
	for _, item := range allItems[:9] {
		storagePacker1.PutItem(ctx, item)
	}

	// Verify no split yet, key should be 00, check it?
	bucketKey00 := storagePacker1.BucketStorageKeyForItemID(ids[0][0])
	bucket00, err := storagePacker1.GetBucket(ctx, bucketKey00, true)
	if err != nil {
		t.Fatalf("Key %v error %v", bucketKey00, err)
	}
	if len(bucket00.Bucket.ItemMap) != 9 {
		t.Fatalf("Pre-split bucket %v contains %v items.",
			bucket00.Bucket.Key,
			len(bucket00.Bucket.ItemMap))
	}

	// Save bucket 00
	storageEntry00, err := storagePacker1.BucketStorageView.Get(ctx, bucketKey00)
	t.Logf("Saved entry key=%v length=%v",
		storageEntry00.Key,
		len(storageEntry00.Value))

	// Add the 10th item, force a split
	storagePacker1.PutItem(ctx, allItems[9])

	// Verify split occurred, get new key
	bucketKey000 := storagePacker1.BucketStorageKeyForItemID(ids[0][0])
	if bucketKey000 != "00/0" {
		t.Fatalf("Unexpected key for shard %v", bucketKey000)
	}

	bucket000, err := storagePacker1.GetBucket(ctx, bucketKey000, false)
	if err != nil {
		t.Fatalf("Key %v error %v", bucketKey000, err)
	}
	if bucket000.Key != "00/0" {
		t.Fatalf("Unexpected key for bucket %v", bucket000.Key)
	}
	if len(bucket000.Bucket.ItemMap) != 4 {
		t.Errorf("Post-split bucket %v contains %v items.",
			bucket000.Bucket.Key,
			len(bucket000.Bucket.ItemMap))
	}

	// Try bypassing cache too.
	bucket000, err = storagePacker1.GetBucket(ctx, bucketKey000, true)
	if err != nil {
		t.Fatalf("Key %v error %v", bucketKey000, err)
	}
	if bucket000.Key != "00/0" {
		t.Fatalf("Unexpected key for bucket %v", bucket000.Key)
	}
	if len(bucket000.Bucket.ItemMap) != 4 {
		t.Errorf("Post-split bucket %v contains %v items.",
			bucket000.Bucket.Key,
			len(bucket000.Bucket.ItemMap))
	}

	// Delete storage 3 through F, and replace the new storage entry
	// with the old one
	for _, shard := range []string{"3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f"} {
		storagePacker1.BucketStorageView.Delete(ctx, "00/"+shard)
	}
	storagePacker1.BucketStorageView.Put(ctx, storageEntry00)

	// Now start a new storage packer using the "corrupted" storage
	storagePacker2 := createStoragePacker(t, storage)

	// Verify that items 0 through 8 still OK
	for idx, item := range allItems[:9] {
		storedItem, err := storagePacker2.GetItem(ctx, item.ID)
		if err != nil {
			t.Errorf("recovery: Error retrieving item %v ID %v: %v",
				idx, item.ID, err)
		} else if storedItem == nil {
			t.Errorf("recovery: Nil item %v ID %v",
				idx, item.ID)
		} else if !bytes.Equal(storedItem.Data, item.Data) {
			t.Errorf("recovery: Item %v ID %v: data mismatch",
				idx, item.ID)
		}
	}

	// Put in a new 000, this will trigger a data mismatch because
	// it goes in 00/0 but sharding 00 will overwrite it with the old
	// version.
	allItems[0] = &Item{
		ID:   ids[0][0],
		Data: incompressibleData(1000),
	}
	storagePacker2.PutItem(ctx, allItems[0])

	// Insert remaining keys (including the one that may have gotten
	// lost because we crashed while inserting it.)
	for _, item := range allItems[9:] {
		storagePacker2.PutItem(ctx, item)
	}

	// Verify that all are present
	for idx, item := range allItems {
		storedItem, err := storagePacker2.GetItem(ctx, item.ID)
		if err != nil {
			t.Errorf("postinsert: Error retrieving item %v ID %v: %v",
				idx, item.ID, err)
		} else if storedItem == nil {
			t.Errorf("postinsert: Nil item %v ID %v",
				idx, item.ID)
		} else if !bytes.Equal(storedItem.Data, item.Data) {
			t.Errorf("postinsert: Item %v ID %v: data mismatch",
				idx, item.ID)
		}
	}
}

func checkAllItems(t *testing.T, storagePacker *StoragePackerV2, ctx context.Context, allItems []*Item) {
	for idx, item := range allItems {
		storedItem, err := storagePacker.GetItem(ctx, item.ID)
		if err != nil {
			t.Errorf("Error retrieving item %v ID %v: %v",
				idx, item.ID, err)
		} else if storedItem == nil {
			t.Errorf("Nil item %v ID %v",
				idx, item.ID)
		} else if !bytes.Equal(storedItem.Data, item.Data) {
			t.Errorf("Item %v ID %v: data mismatch",
				idx, item.ID)
		}
	}
}

func TestStoragePackerV2_Recurse(t *testing.T) {
	// Created from generateLotsOfCollidingIDs(10, "0000")
	ids := []string{
		"0000rZdeLoTTuouIGnsMpqYZqgwWzUpOLKXgRuEV",
		"0000qToiiioByTRsjmBpfkAbsQtSdRTPfVqDGohJ",
		"0000evHgpGImPYdBHIDmSrBGIHafnjgiMCVJypSJ",
		"0000HUdnRbChiiPYrUBUEhwDvRoDeeyjEKaEVifJ",
		"0000fUPXeNVBIIUhzzalccCzzFMTjohWoCjMFbLq",
		"0000YZZGZSdHTLCnqXkhHytkMKrsFidIQRTBywWr",
		"0000OosvkBZTwMPHQkPWbvHFhQwhDAJKUEVjUaAg",
		"0000NHPWdQSmtkbKKfFEazXbPorbFgSidSrubKcK",
		"0000llbJsqIorLwuqjWruVtiLEPYTyZXttzjNLLL",
		"0000IkgtiWeZAzjFUalTJoYfeSeLGIILKoIuaiJV",
	}

	storage := logical.InmemStorageWithMaxSize(10000)
	storagePacker1 := createStoragePacker(t, storage)
	ctx := namespace.RootContext(nil)

	allItems := make([]*Item, 10)
	for i, _ := range allItems[:9] {
		allItems[i] = &Item{
			ID:   ids[i],
			Data: incompressibleData(1000),
		}
	}

	allItems[9] = &Item{
		ID:   ids[9],
		Data: incompressibleData(9000),
	}

	for idx, item := range allItems {
		err := storagePacker1.PutItem(ctx, item)
		if err != nil {
			t.Fatalf("Error inserting key %v %v: %v", idx, item.ID, err)
		}
	}

	checkAllItems(t, storagePacker1, ctx, allItems)

	storagePacker2 := createStoragePacker(t, storage)
	checkAllItems(t, storagePacker2, ctx, allItems)
}

func TestStoragePackerV2_Race(t *testing.T) {
	n := 200
	ids := [][]string{
		generateLotsOfCollidingIDs(200, "00"),
		generateLotsOfCollidingIDs(200, "10"),
	}

	storage := logical.InmemStorageWithMaxSize(10000)
	storagePacker1 := createStoragePacker(t, storage)
	ctx := namespace.RootContext(nil)

	allItems := make([]*Item, n*2)
	for i, _ := range allItems {
		allItems[i] = &Item{
			ID:   ids[i/n][i%n],
			Data: incompressibleData(1000),
		}
	}

	done := make(chan bool)

	writer := func(items []*Item) {
		for _, item := range items {
			err := storagePacker1.PutItem(ctx, item)
			if err != nil {
				t.Errorf("Error inserting key %v: %v", item.ID, err)
			}
		}
		done <- true
	}
	go writer(allItems[:n])
	go writer(allItems[n:])

	_ = <-done
	_ = <-done

	checkAllItems(t, storagePacker1, ctx, allItems)
}

func TestStoragePackerV2_BucketOperations(t *testing.T) {
	ids := generateLotsOfCollidingIDs(80, "00")

	storage := logical.InmemStorageWithMaxSize(10000)
	storagePacker1 := createStoragePacker(t, storage)
	ctx := namespace.RootContext(nil)

	allItems := make([]*Item, 80)
	for i, _ := range allItems {
		allItems[i] = &Item{
			ID:   ids[i],
			Data: incompressibleData(1000),
		}
	}

	// Get bucket 00 from initially empty storage
	bucket, err := storagePacker1.GetBucket(ctx, "00", false)
	if err != nil {
		t.Fatalf("GetBucket error: %v", err)
	}
	if bucket != nil {
		t.Fatalf("GetBucket created a bucket it shouldn't have.")
	}

	// Put one item in to create it
	err = storagePacker1.PutItem(ctx, allItems[0])
	if err != nil {
		t.Fatalf("Error inserting key %v: %v", allItems[0].ID, err)
	}

	// Get bucket 00 again (should hit in cache)
	bucket, err = storagePacker1.GetBucket(ctx, "00", false)
	if err != nil {
		t.Fatalf("GetBucket error: %v", err)
	}
	if bucket == nil {
		t.Fatalf("GetBucket didn't find bucket.")
	}
	if bucket.Key != "00" {
		t.Fatalf("GetBucket mismatched key.")
	}

	// Get bucket 00 again (force cache miss)
	bucket, err = storagePacker1.GetBucket(ctx, "00", true)
	if err != nil {
		t.Fatalf("GetBucket error: %v", err)
	}
	if bucket == nil {
		t.Fatalf("GetBucket didn't find bucket.")
	}
	if bucket.Key != "00" {
		t.Fatalf("GetBucket mismatched key.")
	}

	bucket.Lock()
	for _, item := range allItems {
		itemHash := GetItemIDHash(item.ID)
		bucket.ItemMap[itemHash] = item.Data
	}
	bucket.Unlock()

	err = storagePacker1.PutBucket(ctx, bucket)
	if err != nil {
		t.Fatalf("PutBucket error: %v", err)
	}

	itemHash := GetItemIDHash(allItems[0].ID)
	// Get the shard containing that first key
	bucket, err = storagePacker1.GetBucket(ctx, itemHash, false)
	if err != nil {
		t.Fatalf("GetBucket error: %v", err)
	}
	if bucket == nil {
		t.Fatalf("GetBucket didn't find bucket.")
	}
	key00x := storagePacker1.GetCacheKey(bucket.Key)
	if !strings.HasPrefix(itemHash, key00x) {
		t.Fatalf("GetBucket mismatched key, itemHash=%v bucket.Key=%v", itemHash, bucket.Key)
	}

	t.Logf("Deleting bucket %v", bucket.Key)

	// Delete the whole bucket
	err = storagePacker1.DeleteBucket(ctx, bucket.Key)
	if err != nil {
		t.Fatalf("DeleteBucket error: %v", err)
	}

	// Now, every key should still be present except those matching the deleted bucket key
	for idx, item := range allItems {
		itemHash := GetItemIDHash(item.ID)
		storedItem, err := storagePacker1.GetItem(ctx, item.ID)
		if err != nil {
			t.Errorf("Error retrieving item %v hash %v: %v",
				idx, itemHash, err)
		} else if strings.HasPrefix(itemHash, key00x) {
			// Should have been deleted
			if storedItem != nil {
				t.Errorf("Item %v hash %v is still present", idx, itemHash)
			}
		} else {
			if storedItem == nil {
				t.Errorf("Nil item %v hash %v", idx, itemHash)
			} else if !bytes.Equal(storedItem.Data, item.Data) {
				t.Errorf("Item %v ash %v: data mismatch", idx, itemHash)
			}
		}
	}

}
