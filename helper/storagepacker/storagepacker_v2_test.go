package storagepacker

import (
	"context"
	"testing"

	"github.com/go-test/deep"
	"github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	log "github.com/hashicorp/go-hclog"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
)

func getStoragePacker(tb testing.TB) *StoragePackerV2 {
	storage := &logical.InmemStorage{}
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

func BenchmarkStoragePackerV2(b *testing.B) {
	storagePacker := getStoragePacker(b)

	ctx := namespace.RootContext(nil)

	ctx := context.Background()

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
		ID: "item1",
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

	marshaledEntity, err := ptypes.MarshalAny(entity)
	if err != nil {
		t.Fatal(err)
	}

	ctx := namespace.RootContext(nil)
	err = storagePacker.PutItem(ctx, &Item{
		ID:      entity.ID,
		Message: marshaledEntity,
	})
	if err != nil {
		t.Fatal(err)
	}

	itemFetched, err := storagePacker.GetItem(ctx, entity.ID)
	if err != nil {
		t.Fatal(err)
	}

	var itemDecoded identity.Entity
	err = ptypes.UnmarshalAny(itemFetched.Message, &itemDecoded)
	if err != nil {
		t.Fatal(err)
	}

	if !proto.Equal(&itemDecoded, entity) {
		diff := deep.Equal(&itemDecoded, entity)
		t.Fatal(diff)
	}
}
