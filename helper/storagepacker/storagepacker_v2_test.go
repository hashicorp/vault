package storagepacker

import (
	"io/ioutil"
	"os"
	"strconv"
	"testing"

	"github.com/golang/protobuf/ptypes"
	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/helper/logformat"
	"github.com/hashicorp/vault/logical"
	log "github.com/mgutz/logxi/v1"
)

func testPutItem(t *testing.T, sp *StoragePackerV2, count int, entity *identity.Entity) {
	t.Helper()
	for i := 1; i <= count; i++ {
		/*
			if i%500 == 0 {
				fmt.Printf("put item iteration: %d\n", i)
			}
		*/
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

func testGetItem(t *testing.T, sp *StoragePackerV2, count int, expectNil bool) {
	t.Helper()
	for i := 1; i <= count; i++ {
		/*
			if i%500 == 0 {
				fmt.Printf("get item iteration: %d\n", i)
			}
		*/
		id := strconv.Itoa(i)

		itemFetched, err := sp.GetItem(id)
		if err != nil {
			t.Fatal(err)
		}

		switch expectNil {
		case itemFetched == nil:
			return
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

func testDeleteItem(t *testing.T, sp *StoragePackerV2, count int) {
	t.Helper()
	for i := 1; i <= count; i++ {
		/*
			if i%500 == 0 {
				fmt.Printf("delete item iteration: %d\n", i)
			}
		*/
		id := strconv.Itoa(i)
		err := sp.DeleteItem(id)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestStoragePackerV2_PutGetDeleteInmem(t *testing.T) {
	sp, err := NewStoragePackerV2(&Config{
		View:   &logical.InmemStorage{},
		Logger: log.New("storagepackertest"),
	})
	if err != nil {
		t.Fatal(err)
	}

	entity := &identity.Entity{
		Metadata: map[string]string{
			"samplekey": "samplevalue",
		},
	}
	count := 1000
	testPutItem(t, sp, count, entity)
	testGetItem(t, sp, count, false)
	testDeleteItem(t, sp, count)
	testGetItem(t, sp, count, true)
}

func TestStoragePackerV2_PutGetDelete_File(t *testing.T) {
	filePath, err := ioutil.TempDir(".", "vault")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.RemoveAll(filePath)

	logger := logformat.NewVaultLogger(log.LevelTrace)

	config := map[string]string{
		"path": filePath,
	}

	storage, err := logical.NewLogicalStorage(logical.LogicalTypeFile, config, logger)
	if err != nil {
		t.Fatal(err)
	}

	sp, err := NewStoragePackerV2(&Config{
		View:             storage,
		Logger:           log.New("storagepackertest"),
		BucketCount:      256,
		BucketShardCount: 32,
		BucketMaxSize:    256 * 1024,
	})
	if err != nil {
		t.Fatal(err)
	}

	count := 100
	entity := &identity.Entity{
		Metadata: map[string]string{
			"samplekey1": "samplevalue1",
			"samplekey2": "samplevalue2",
			"samplekey3": "samplevalue3",
			"samplekey4": "samplevalue4",
			"samplekey5": "samplevalue5",
		},
	}

	testPutItem(t, sp, count, entity)
	testGetItem(t, sp, count, false)
	testDeleteItem(t, sp, count)
	testGetItem(t, sp, count, true)
}
