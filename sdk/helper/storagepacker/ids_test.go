package storagepacker

import (
	"context"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/logical"
	"testing"
)

func getStoragePackerBits(tb testing.TB, baseBucketBits int, bucketShardBits int) *StoragePackerV2 {
	storage := &logical.InmemStorage{}
	storageView := logical.NewStorageView(storage, "packer/buckets/v2")
	storagePacker, err := NewStoragePackerV2(context.Background(), &Config{
		BucketStorageView: storageView,
		ConfigStorageView: logical.NewStorageView(storage, "packer/config"),
		Logger:            log.New(&log.LoggerOptions{Name: "storagepackertest"}),
		BaseBucketBits:    baseBucketBits,
		BucketShardBits:   bucketShardBits,
	})
	if err != nil {
		tb.Fatal(err)
	}
	return storagePacker.(*StoragePackerV2)
}

func TestStoragePackerV2_FirstBucketKey(t *testing.T) {
	// Default settings
	s_8_4 := getStoragePackerBits(t, 8, 4)

	key, err := s_8_4.firstKey("12345678")
	if key != "12" || err != nil {
		t.Fatalf("first key should be 12, got %q, %v", key, err)
	}
	// TODO: is there a way to check later elements of the storage key?
	// Maybe requires actual sharding as the code is now?

	key, err = s_8_4.firstKey("1")
	if err == nil {
		t.Fatalf("should get key too short error")
	}

	s_12_8 := getStoragePackerBits(t, 12, 8)
	key, err = s_12_8.firstKey("12345678")
	if key != "123" || err != nil {
		t.Fatalf("first key should be 123, got %q, %v", key, err)
	}

	s_4_4 := getStoragePackerBits(t, 4, 4)
	key, err = s_4_4.firstKey("fedc")
	if key != "f" || err != nil {
		t.Fatalf("first key should be f, got %q, %v", key, err)
	}

}
