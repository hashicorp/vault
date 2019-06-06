package storagepacker

import (
	"context"
	"testing"

	"bytes"
	"fmt"
	"github.com/go-test/deep"
	"github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	log "github.com/hashicorp/go-hclog"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/hashicorp/vault/sdk/physical/inmem"
	"math/rand"
	"strings"
)

func createStoragePacker(tb testing.TB, storage logical.Storage) *StoragePackerV2 {
	storageView := logical.NewStorageView(storage, "packer/buckets/")
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

func getEntity() *identity.Entity {
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
	return entity
}

func TestStoragePackerV2_SerializeDeserializeComplexItem(t *testing.T) {
	storagePacker := getStoragePacker(t)

	ctx := context.Background()

	entity := getEntity()
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

	// Did configuration get saved?
	entry, err := storagePacker1.ConfigStorageView.Get(ctx, "config")
	if err != nil {
		t.Fatal(err)
	}
	if entry == nil {
		t.Fatal("No config saved.")
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
	bucket00, err := storagePacker1.GetBucket(ctx, ids[0][0])
	if err != nil {
		t.Fatalf("Key %v error %v", ids[0][0], err)
	}
	bucketKey00 := bucket00.Key
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
	bucket000, err := storagePacker1.GetBucket(ctx, ids[0][0])
	if err != nil {
		t.Fatalf("Key %v error %v", ids[0][0], err)
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
			t.Errorf("Nil item %v ID %v Key %v",
				idx, item.ID, GetItemIDHash(item.ID))
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

func run_concurrentTest(t *testing.T, n int, ids [][]string) {
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

func TestStoragePackerV2_PartitionedInserts(t *testing.T) {
	n := 200
	ids := [][]string{
		generateLotsOfCollidingIDs(n, "00"),
		generateLotsOfCollidingIDs(n, "10"),
	}
	run_concurrentTest(t, n, ids)
}

func TestStoragePackerV2_ConcurrentInserts(t *testing.T) {
	n := 200
	ids := [][]string{
		generateLotsOfCollidingIDs(n, "00"),
		generateLotsOfCollidingIDs(n, "00"),
	}
	run_concurrentTest(t, n, ids)
}

func TestStoragePackerV2_Variadic(t *testing.T) {
	storage := logical.InmemStorageWithMaxSize(10000)
	storagePacker := createStoragePacker(t, storage)
	ctx := namespace.RootContext(nil)
	n := 257 // Ensure at least one collision

	allItems := make([]*Item, n)
	ids := make([]string, n)
	for i, _ := range allItems {
		ids[i] = fmt.Sprintf("test-%d", i)
		allItems[i] = &Item{
			ID:   ids[i],
			Data: incompressibleData(1000),
		}
	}

	err := storagePacker.PutItem(ctx, allItems...)
	if err != nil {
		t.Fatalf("error on PutItem: %v", err)
	}

	requestedIDs := ids[100:150]
	items, err := storagePacker.GetItems(ctx, requestedIDs...)
	if err != nil {
		t.Fatalf("error on GetItems: %v", err)
	}
	if len(items) != len(requestedIDs) {
		t.Fatalf("mismatched response length %v request length %v", len(items), len(requestedIDs))
	}
	for idx, item := range items {
		if item == nil {
			t.Errorf("nil item %v ID %v", idx, requestedIDs[idx])
		} else if item.ID != requestedIDs[idx] {
			t.Errorf("mismatched return value on item %v requested %v got %v", idx, requestedIDs[idx], item.ID)
		} else if !bytes.Equal(allItems[100+idx].Data, item.Data) {
			t.Errorf("item %v ID %v: data mismatch", idx, item.ID)
		}
	}

	toDelete := ids[:100]
	err = storagePacker.DeleteItem(ctx, toDelete...)
	if err != nil {
		t.Fatalf("error on DeleteItem: %v", err)
	}
	checkAllItems(t, storagePacker, ctx, allItems[100:])

	requestedIDs = []string{
		"test-1",            // deleted
		"test-120",          // still present
		"test-doesnotexist", // never there to begin with
	}
	items, err = storagePacker.GetItems(ctx, requestedIDs...)
	if err != nil {
		t.Fatalf("error on GetItems: %v", err)
	}
	if len(items) != 3 {
		t.Fatalf("mismatched response length %v request length %v", len(items), 3)
	}
	if items[0] != nil {
		t.Errorf("test-1 sill present")
	}
	if items[1] == nil {
		t.Errorf("test-120 not present")
	} else if items[1].ID != "test-120" {
		t.Errorf("test-120 ID mismatch, got %v", items[1].ID)
	}
	if items[2] != nil {
		t.Errorf("test-doesnotexist reported present")
	}
}

func checkReturnedItems(t *testing.T, allItems []*Item, returnedItems []*Item) {
	itemsById := make(map[string]*Item, len(allItems))

	for _, item := range allItems {
		itemsById[item.ID] = item
	}

	for _, item := range returnedItems {
		if item == nil {
			t.Fatalf("nil item in list")
			return
		}
		orig, found := itemsById[item.ID]
		if !found {
			t.Errorf("item %q unexpectedly present or duplicated", item.ID)
		} else if !bytes.Equal(orig.Data, item.Data) {
			t.Errorf("item %q: data mismatch", item.ID)
		}
		delete(itemsById, item.ID)
	}

	// should be empty
	for k, _ := range itemsById {
		t.Errorf("item %q not found in returned list", k)
	}
}

func TestStoragePackerV2_AllItems(t *testing.T) {
	storage := logical.InmemStorageWithMaxSize(10000)
	storagePacker := createStoragePacker(t, storage)
	ctx := namespace.RootContext(nil)

	ids := append(generateLotsOfCollidingIDs(12, "00"),
		generateLotsOfCollidingIDs(20, "1")...)
	n := len(ids)

	allItems := make([]*Item, n)
	for i, _ := range allItems {
		allItems[i] = &Item{
			ID:   ids[i],
			Data: incompressibleData(1000),
		}
	}

	err := storagePacker.PutItem(ctx, allItems...)
	if err != nil {
		t.Fatalf("error on PutItem: %+v", err)
	}

	stored, err := storagePacker.AllItems(ctx)
	if err != nil {
		t.Fatalf("error on AllItems: %+v", err)
	}

	checkReturnedItems(t, allItems, stored)

	storagePacker2 := createStoragePacker(t, storage)
	stored, err = storagePacker2.AllItems(ctx)
	if err != nil {
		t.Fatalf("error on recovered AllItems: %+v", err)
	}

	checkReturnedItems(t, allItems, stored)
}

func TestStoragePackerV2_Queued(t *testing.T) {
	storage := logical.InmemStorageWithMaxSize(10000)
	storagePacker := createStoragePacker(t, storage)
	ctx := namespace.RootContext(nil)

	numBatches := 2
	batchSize := 80
	n := numBatches * batchSize

	allItems := make([]*Item, n)
	ids := make([]string, n)
	for i, _ := range allItems {
		ids[i] = fmt.Sprintf("test-%d", i)
		allItems[i] = &Item{
			ID:   ids[i],
			Data: incompressibleData(100),
		}
	}

	storagePacker.SetQueueMode(true)

	for b := 1; b < numBatches; b++ {
		startIdx := b * batchSize
		endIdx := (b + 1) * batchSize
		err := storagePacker.PutItem(ctx, allItems[startIdx:endIdx]...)
		if err != nil {
			t.Fatal(err)
		}
	}

	diskBuckets, err := logical.CollectKeys(ctx, storagePacker.BucketStorageView)
	if err != nil {
		t.Fatal(err)
	}
	if len(diskBuckets) != 0 {
		t.Fatal("storage buckets exist even though queue mode is off")
	}

	err = storagePacker.FlushQueue(ctx)
	if err != nil {
		t.Fatal(err)
	}

	storagePacker.SetQueueMode(false)
	storagePacker.PutItem(ctx, allItems[:batchSize]...)
	checkAllItems(t, storagePacker, ctx, allItems)

	// Recover into second object to verify persistence
	storagePacker2 := createStoragePacker(t, storage)
	checkAllItems(t, storagePacker2, ctx, allItems)

}

func TestStoragePackerV2_CreationErrors(t *testing.T) {
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}
	bucketStorageView := logical.NewStorageView(storage, "packer/buckets/v2")
	configStorageView := logical.NewStorageView(storage, "packer/config")
	logger := log.New(&log.LoggerOptions{Name: "storagepackertest"})

	cases := []struct {
		Config *Config
		Error  string
	}{
		{&Config{
			BucketStorageView: nil,
			ConfigStorageView: configStorageView,
			Logger:            logger,
		}, "nil buckets view"},
		{&Config{
			BucketStorageView: bucketStorageView,
			ConfigStorageView: nil,
			Logger:            logger,
		}, "nil config view"},
		{&Config{
			BucketStorageView: bucketStorageView,
			ConfigStorageView: configStorageView,
			Logger:            nil,
		}, "nil logger"},
		{&Config{
			BucketStorageView: bucketStorageView,
			ConfigStorageView: configStorageView,
			Logger:            logger,
			BaseBucketBits:    -8,
			BucketShardBits:   4,
		}, "should be at least 4"},
		{&Config{
			BucketStorageView: bucketStorageView,
			ConfigStorageView: configStorageView,
			Logger:            logger,
			BaseBucketBits:    8,
			BucketShardBits:   -8,
		}, "should be at least 4"},
		{&Config{
			BucketStorageView: bucketStorageView,
			ConfigStorageView: configStorageView,
			Logger:            logger,
			BaseBucketBits:    7,
			BucketShardBits:   4,
		}, "is not a multiple of 4"},
		{&Config{
			BucketStorageView: bucketStorageView,
			ConfigStorageView: configStorageView,
			Logger:            logger,
			BaseBucketBits:    12,
			BucketShardBits:   2,
		}, "is not a multiple of 4"},
	}

	for _, tc := range cases {
		_, err := NewStoragePackerV2(ctx, tc.Config)
		if err == nil {
			t.Fatalf("No error, expected %q", tc.Error)
		}
		if !strings.Contains(err.Error(), tc.Error) {
			t.Fatalf("Error %q didn't match %q", err.Error(), tc.Error)
		}
	}

}

func TestStoragePackerV2_PutItemErrors(t *testing.T) {
	storage := logical.InmemStorageWithMaxSize(10000)
	storagePacker := createStoragePacker(t, storage)
	ctx := namespace.RootContext(nil)

	n := 10
	allItems := make([]*Item, n)
	for i, _ := range allItems {
		allItems[i] = &Item{
			ID:   fmt.Sprintf("test-%d", i),
			Data: []byte("dontcare"),
		}
	}

	entity := getEntity()
	message, err := ptypes.MarshalAny(entity)
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		Which       int
		Error       string
		Replacement *Item
	}{
		{0, "nil item", nil},
		{1, "missing ID", &Item{
			ID:   "",
			Data: []byte("dontcare")}},
		{2, "missing data", &Item{
			ID:   "present",
			Data: nil}},
		{3, "deprecated", &Item{
			ID:      "present",
			Message: message}},
		{4, "duplicate", &Item{
			ID:   "test-1",
			Data: []byte("dontcare")}},
	}

	for _, tc := range cases {
		tmp := allItems[tc.Which]
		allItems[tc.Which] = tc.Replacement
		err := storagePacker.PutItem(ctx, allItems...)
		if err == nil {
			t.Fatalf("missing error check for %q", tc.Error)
		}
		if !strings.Contains(err.Error(), tc.Error) {
			t.Fatalf("Error %q didn't match %q", err.Error(), tc.Error)
		}
		allItems[tc.Which] = tmp
	}
}

func TestStoragePackerV2_UpsertErrors(t *testing.T) {
	id := "123456"

	b := &LockedBucket{
		Bucket: &Bucket{
			Key:       "00",
			Items:     []*Item{},
			ItemMap:   nil,
			HasShards: false,
		},
	}
	item := &itemRequest{
		ID:  id,
		Key: GetItemIDHash(id),
		Value: &Item{
			ID:      id,
			Message: nil,
			Data:    []byte{0, 1, 2, 3, 4},
		},
		Bucket: b.Bucket,
	}
	p := &partitionedRequests{
		Bucket:   b,
		Requests: []*itemRequest{item},
	}

	var nilPR *partitionedRequests
	nilPR = nil
	err := nilPR.upsertItems()
	if err == nil {
		t.Fatalf("no error on nil receiver")
	}

	p2 := &partitionedRequests{
		Bucket:   nil,
		Requests: []*itemRequest{item},
	}
	err = p2.upsertItems()
	if err == nil {
		t.Fatalf("no error on nil bucket")
	}

	b.HasShards = true
	err = p.upsertItems()
	if err == nil {
		t.Fatalf("no error sharded bucket")
	}

	b.HasShards = false
	err = p.upsertItems()
	if err != nil {
		t.Fatalf("upsert expected to succeed on nil ItemMap: %v", err)
	}
	if b.ItemMap == nil {
		t.Fatalf("upsert didn't create ItemMap")
	}
}

// A clone of InMemStorage that fails on particular paths, for testing sharding
type FailingInmemStorage struct {
	underlying   physical.Backend
	suffixToFail string
	t            *testing.T
}

func NewFailingStorage(t *testing.T) *FailingInmemStorage {
	conf := make(map[string]string)
	conf["max_value_size"] = "10000"
	phy, _ := inmem.NewInmem(conf, nil)
	return &FailingInmemStorage{
		underlying:   phy,
		suffixToFail: "", // default: everything
		t:            t,
	}
}

func (s *FailingInmemStorage) Get(ctx context.Context, key string) (*logical.StorageEntry, error) {
	if strings.Contains(key, s.suffixToFail) {
		return nil, fmt.Errorf("Failing for test purposes.")
	}

	entry, err := s.underlying.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}
	return &logical.StorageEntry{
		Key:      entry.Key,
		Value:    entry.Value,
		SealWrap: entry.SealWrap,
	}, nil
}

func (s *FailingInmemStorage) Put(ctx context.Context, entry *logical.StorageEntry) error {
	s.t.Logf("Put %v\n", entry.Key)
	// Don't fail the put that would fail with a too large exception
	if len(entry.Value) <= 10000 && strings.HasSuffix(entry.Key, s.suffixToFail) {
		return fmt.Errorf("Failing for test purposes.")
	}

	return s.underlying.Put(ctx, &physical.Entry{
		Key:      entry.Key,
		Value:    entry.Value,
		SealWrap: entry.SealWrap,
	})
}

func (s *FailingInmemStorage) Delete(ctx context.Context, key string) error {
	s.t.Logf("Delete %v\n", key)
	if strings.HasSuffix(key, s.suffixToFail) {
		return fmt.Errorf("Failing for test purposes.")
	}

	return s.underlying.Delete(ctx, key)
}

func (s *FailingInmemStorage) List(ctx context.Context, prefix string) ([]string, error) {
	return s.underlying.List(ctx, prefix)
}

func testFailedPutWhileSharding(t *testing.T, storagePacker *StoragePackerV2, allItems []*Item) {
	ctx := namespace.RootContext(nil)
	err := storagePacker.PutItem(ctx, allItems...)
	if err == nil {
		t.Fatalf("error not signalled to PutItem")
	} else {
		t.Logf("PutItem error: %+v", err)
	}

	// Check that no intermediate buckets are hanging around in cache or storage
	bucket, err := storagePacker.GetBucket(ctx, "0000")
	if err != nil {
		t.Fatal(err)
	}
	if bucket.Key != "00" {
		t.Fatalf("unexpected bucket key %q", bucket.Key)
	}

	list, err := storagePacker.BucketStorageView.List(ctx, "00/")
	if err != nil {
		t.Fatal(err)
	}
	if len(list) > 0 {
		t.Fatalf("expected empty list of storage objects, got %v", list)
	}
}

func TestStoragePackerV2_StorageErrors(t *testing.T) {
	ctx := namespace.RootContext(nil)
	storage := NewFailingStorage(t)
	storage.suffixToFail = "DOESNTMATCH"
	storagePacker := createStoragePacker(t, storage)

	backend := storage.underlying.(*inmem.InmemBackend)
	backend.FailPut(true)
	err := storagePacker.PutItem(ctx, &Item{ID: "test", Data: []byte("test")})
	if err == nil {
		t.Fatalf("Error not signalled to PutItem")
	} else {
		t.Logf("PutItem error: %+v", err)
	}

	backend.FailPut(false)
	err = storagePacker.PutItem(ctx, &Item{ID: "test", Data: []byte("test")})
	if err != nil {
		t.Fatal(err)
	}

	backend.FailPut(true)
	err = storagePacker.DeleteItem(ctx, "test")
	if err == nil {
		t.Fatalf("Error not signalled to DeleteItem")
	} else {
		t.Logf("DeleteItem error: %+v", err)
	}

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

	allItems := make([]*Item, 10)
	for i, _ := range allItems {
		allItems[i] = &Item{
			ID:   ids[i],
			Data: incompressibleData(1000),
		}
	}

	backend.FailPut(false)
	// Induce failure on one of the leaf shards created
	storage.suffixToFail = "00/0/0/c"
	testFailedPutWhileSharding(t, storagePacker, allItems)

	// Induce failure on the root bucket instead
	// This causes the test to fail because cleanup only does
	// the immediate shards, not the ones below those.
	/*
		storage.suffixToFail = "00"
		testFailedPutWhileSharding(t, storagePacker, allItems)
	*/
}
