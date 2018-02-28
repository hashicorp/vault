package storagepacker

import (
	"reflect"
	"testing"

	"github.com/golang/protobuf/ptypes"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/logical"
	log "github.com/mgutz/logxi/v1"
)

func BenchmarkStoragePacker(b *testing.B) {
	storagePacker, err := NewStoragePacker(&logical.InmemStorage{}, log.New("storagepackertest"), "")
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		itemID, err := uuid.GenerateUUID()
		if err != nil {
			b.Fatal(err)
		}

		item := &Item{
			ID: itemID,
		}

		err = storagePacker.PutItem(item)
		if err != nil {
			b.Fatal(err)
		}

		fetchedItem, err := storagePacker.GetItem(itemID)
		if err != nil {
			b.Fatal(err)
		}

		if fetchedItem == nil {
			b.Fatalf("failed to read stored item with ID: %q, iteration: %d", item.ID, i)
		}

		if fetchedItem.ID != item.ID {
			b.Fatalf("bad: item ID; expected: %q\n actual: %q", item.ID, fetchedItem.ID)
		}

		err = storagePacker.DeleteItem(item.ID)
		if err != nil {
			b.Fatal(err)
		}

		fetchedItem, err = storagePacker.GetItem(item.ID)
		if err != nil {
			b.Fatal(err)
		}
		if fetchedItem != nil {
			b.Fatalf("failed to delete item")
		}
	}
}

func TestStoragePacker(t *testing.T) {
	storagePacker, err := NewStoragePacker(&logical.InmemStorage{}, log.New("storagepackertest"), "")
	if err != nil {
		t.Fatal(err)
	}

	// Persist a storage entry
	item1 := &Item{
		ID: "item1",
	}

	err = storagePacker.PutItem(item1)
	if err != nil {
		t.Fatal(err)
	}

	// Verify that it can be read
	fetchedItem, err := storagePacker.GetItem(item1.ID)
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
	err = storagePacker.DeleteItem(item1.ID)
	if err != nil {
		t.Fatal(err)
	}

	// Check that the deletion was successful
	fetchedItem, err = storagePacker.GetItem(item1.ID)
	if err != nil {
		t.Fatal(err)
	}

	if fetchedItem != nil {
		t.Fatalf("failed to delete item")
	}
}

func TestStoragePacker_SerializeDeserializeComplexItem(t *testing.T) {
	storagePacker, err := NewStoragePacker(&logical.InmemStorage{}, log.New("storagepackertest"), "")
	if err != nil {
		t.Fatal(err)
	}

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
		BucketKeyHash:   "entity_hash",
		MergedEntityIDs: []string{"merged_entity_id1", "merged_entity_id2"},
		Policies:        []string{"policy1", "policy2"},
	}

	marshaledEntity, err := ptypes.MarshalAny(entity)
	if err != nil {
		t.Fatal(err)
	}
	err = storagePacker.PutItem(&Item{
		ID:      entity.ID,
		Message: marshaledEntity,
	})
	if err != nil {
		t.Fatal(err)
	}

	itemFetched, err := storagePacker.GetItem(entity.ID)
	if err != nil {
		t.Fatal(err)
	}

	var itemDecoded identity.Entity
	err = ptypes.UnmarshalAny(itemFetched.Message, &itemDecoded)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(&itemDecoded, entity) {
		t.Fatalf("bad: expected: %#v\nactual: %#v\n", entity, itemDecoded)
	}
}
