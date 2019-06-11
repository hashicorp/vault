package storagepacker

import (
	"bytes"
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

func TestStoragePackerV2_Invalidate_Creation(t *testing.T) {
	// Create the []bytes version of a bucket by hand
	n := 10
	ids := generateLotsOfCollidingIDs(n, "00")

	allItems := make([]*Item, n)
	for i, _ := range allItems {
		allItems[i] = &Item{
			ID:   ids[i],
			Data: incompressibleData(1000),
		}
	}

	b := &Bucket{
		Key:     "00",
		ItemMap: make(map[string][]byte, n),
	}
	for _, item := range allItems {
		b.ItemMap[item.ID] = item.Data
	}

	rawBucket, err := marshalledBucket(b)
	if err != nil {
		t.Fatal(err)
	}
	bucketPath := "packer/buckets/v2/00"

	// Invalidate an empty storage packer to see if all
	// items are reported.
	storagePacker := getStoragePacker(t)
	ctx := namespace.RootContext(nil)
	present, deleted, err := storagePacker.InvalidateItems(ctx, bucketPath, rawBucket)
	if err != nil {
		t.Fatal(err)
	}
	if len(present) != n {
		t.Fatalf("%d elements ppresent", len(present))
	}
	if len(deleted) != 0 {
		t.Fatalf("%d elements deleted", len(deleted))
	}

	// Check return value items
	checkItems := make(map[string][]byte, n)
	for _, item := range allItems {
		checkItems[item.ID] = item.Data
	}
	for idx, item := range present {
		if check, found := checkItems[item.ID]; found {
			if !bytes.Equal(check, item.Data) {
				t.Errorf("mismatched item %d ID %q in present", idx, item.ID)
			}
			delete(checkItems, item.ID)
		} else {
			t.Errorf("unknown or duplicate item %d ID %q in present", idx, item.ID)
		}
	}

	// Also check that SP has been modified
	checkAllItems(t, storagePacker, ctx, allItems)

}
