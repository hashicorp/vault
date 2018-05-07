package storagepacker

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"testing"

	"github.com/golang/protobuf/ptypes"
	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/helper/logging"
	"github.com/hashicorp/vault/logical"

	log "github.com/hashicorp/go-hclog"
)

const (
	testIterationCount   = 5000
	testBucketBaseCount  = defaultBucketBaseCount
	testBucketShardCount = defaultBucketShardCount
	testBucketMaxSize    = defaultBucketMaxSize
)

func TestStoragePackerV2_Inmem(t *testing.T) {
	sp, err := NewStoragePackerV2(&Config{
		BucketBaseCount:  testBucketBaseCount,
		BucketShardCount: testBucketShardCount,
		BucketMaxSize:    testBucketMaxSize,
		View:             &logical.InmemStorage{},
		Logger:           logging.NewVaultLogger(log.Trace),
	})
	if err != nil {
		t.Fatal(err)
	}

	entity := &identity.Entity{
		Metadata: map[string]string{
			"samplekey1": "samplevalue1",
			"samplekey2": "samplevalue2",
			"samplekey3": "samplevalue3",
			"samplekey4": "samplevalue4",
			"samplekey5": "samplevalue5",
		},
	}
	testPutItem(t, sp, entity)
	testGetItem(t, sp, false)
	testDeleteItem(t, sp)
	testGetItem(t, sp, true)
}

func TestStoragePackerV2_File(t *testing.T) {
	filePath, err := ioutil.TempDir("", "vault")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	//fmt.Printf("filePath: %q\n", filePath)
	defer os.RemoveAll(filePath)

	logger := logging.NewVaultLogger(log.Trace)

	config := map[string]string{
		"path": filePath,
	}

	storage, err := logical.NewLogicalStorage(logical.LogicalTypeFile, config, logger)
	if err != nil {
		t.Fatal(err)
	}

	sp, err := NewStoragePackerV2(&Config{
		BucketBaseCount:  testBucketBaseCount,
		BucketShardCount: testBucketShardCount,
		BucketMaxSize:    testBucketMaxSize,
		View:             storage,
		Logger:           logger,
	})
	if err != nil {
		t.Fatal(err)
	}

	entity := &identity.Entity{
		Metadata: map[string]string{
			"samplekey1": "samplevalue1",
			"samplekey2": "samplevalue2",
			"samplekey3": "samplevalue3",
			"samplekey4": "samplevalue4",
			"samplekey5": "samplevalue5",
		},
	}

	testPutItem(t, sp, entity)
	testGetItem(t, sp, false)
	testDeleteItem(t, sp)
	testGetItem(t, sp, true)
}

func TestStoragePackerV2_isPowerOfTwo(t *testing.T) {
	powersOfTwo := []int{1, 2, 4, 1024, 4096}
	notPowersOfTwo := []int{0, 3, 5, 1000, 1023, 4095, 4097, 10000}
	for _, val := range powersOfTwo {
		if !isPowerOfTwo(val) {
			t.Fatalf("%d is a power of two", val)
		}
	}
	for _, val := range notPowersOfTwo {
		if isPowerOfTwo(val) {
			t.Fatalf("%d is not a power of two", val)
		}
	}
}

func testPutItem(t *testing.T, sp *StoragePackerV2, entity *identity.Entity) {
	t.Helper()
	for i := 1; i <= testIterationCount; i++ {
		if i%500 == 0 {
			fmt.Printf("put item iteration: %d\n", i)
		}
		id := strconv.Itoa(i)
		entity.ID = id

		marshaledMessage, err := ptypes.MarshalAny(entity)
		if err != nil {
			t.Fatal(err)
		}

		item := &Item{
			ID:      id,
			Message: marshaledMessage,
		}
		if err != nil {
			t.Fatal(err)
		}

		_, err = sp.PutItem(item)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func testGetItem(t *testing.T, sp *StoragePackerV2, expectNil bool) {
	t.Helper()
	for i := 1; i <= testIterationCount; i++ {
		if i%500 == 0 {
			fmt.Printf("get item iteration: %d\n", i)
		}
		id := strconv.Itoa(i)

		itemFetched, err := sp.GetItem(id)
		if err != nil {
			t.Fatal(err)
		}

		switch expectNil {
		case itemFetched == nil:
			continue
		default:
			t.Fatalf("expected nil for item %q\n", id)
		}

		if itemFetched == nil {
			t.Fatalf("failed to read the inserted item %q", id)
		}

		var fetchedMessage identity.Entity
		err = ptypes.UnmarshalAny(itemFetched.Message, &fetchedMessage)
		if err != nil {
			t.Fatal(err)
		}

		if fetchedMessage.ID != id {
			t.Fatalf("failed to fetch item ID: %q\n", id)
		}
	}
}

func testDeleteItem(t *testing.T, sp *StoragePackerV2) {
	t.Helper()
	for i := 1; i <= testIterationCount; i++ {
		if i%500 == 0 {
			fmt.Printf("delete item iteration: %d\n", i)
		}
		id := strconv.Itoa(i)
		err := sp.DeleteItem(id)
		if err != nil {
			t.Fatal(err)
		}
	}
}
