// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package storagepacker

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	log "github.com/hashicorp/go-hclog"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/sdk/logical"
	"google.golang.org/protobuf/types/known/anypb"
)

func BenchmarkStoragePacker(b *testing.B) {
	storagePacker, err := NewStoragePacker(&logical.InmemStorage{}, log.New(&log.LoggerOptions{Name: "storagepackertest"}), "")
	if err != nil {
		b.Fatal(err)
	}

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

		err = storagePacker.DeleteItem(ctx, item.ID)
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
	storagePacker, err := NewStoragePacker(&logical.InmemStorage{}, log.New(&log.LoggerOptions{Name: "storagepackertest"}), "")
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	// Persist a storage entry
	item1 := &Item{
		ID: "item1",
	}

	err = storagePacker.PutItem(ctx, item1)
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
	err = storagePacker.DeleteItem(ctx, item1.ID)
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
	storagePacker, err := NewStoragePacker(&logical.InmemStorage{}, log.New(&log.LoggerOptions{Name: "storagepackertest"}), "")
	if err != nil {
		t.Fatal(err)
	}

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

	marshaledEntity, err := anypb.New(entity)
	if err != nil {
		t.Fatal(err)
	}
	err = storagePacker.PutItem(ctx, &Item{
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

	if !proto.Equal(&itemDecoded, entity) {
		t.Fatalf("bad: expected: %#v\nactual: %#v\n", entity, itemDecoded)
	}
}

func TestStoragePacker_DeleteMultiple(t *testing.T) {
	storagePacker, err := NewStoragePacker(&logical.InmemStorage{}, log.New(&log.LoggerOptions{Name: "storagepackertest"}), "")
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	// Persist a storage entry
	for i := 0; i < 100; i++ {
		item := &Item{
			ID: fmt.Sprintf("item%d", i),
		}

		err = storagePacker.PutItem(ctx, item)
		if err != nil {
			t.Fatal(err)
		}

		// Verify that it can be read
		fetchedItem, err := storagePacker.GetItem(item.ID)
		if err != nil {
			t.Fatal(err)
		}
		if fetchedItem == nil {
			t.Fatalf("failed to read the stored item")
		}

		if item.ID != fetchedItem.ID {
			t.Fatalf("bad: item ID; expected: %q\n actual: %q\n", item.ID, fetchedItem.ID)
		}
	}

	itemsToDelete := make([]string, 0, 50)
	for i := 1; i < 100; i += 2 {
		itemsToDelete = append(itemsToDelete, fmt.Sprintf("item%d", i))
	}

	err = storagePacker.DeleteMultipleItems(ctx, nil, itemsToDelete)
	if err != nil {
		t.Fatal(err)
	}

	// Check that the deletion was successful
	for i := 0; i < 100; i++ {
		fetchedItem, err := storagePacker.GetItem(fmt.Sprintf("item%d", i))
		if err != nil {
			t.Fatal(err)
		}

		if i%2 == 0 && fetchedItem == nil {
			t.Fatal("expected item not found")
		}
		if i%2 == 1 && fetchedItem != nil {
			t.Fatalf("failed to delete item")
		}
	}
}

func TestStoragePacker_DeleteMultiple_ALL(t *testing.T) {
	storagePacker, err := NewStoragePacker(&logical.InmemStorage{}, log.New(&log.LoggerOptions{Name: "storagepackertest"}), "")
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	// Persist a storage entry
	itemsToDelete := make([]string, 0, 10000)
	for i := 0; i < 10000; i++ {
		item := &Item{
			ID: fmt.Sprintf("item%d", i),
		}

		err = storagePacker.PutItem(ctx, item)
		if err != nil {
			t.Fatal(err)
		}

		// Verify that it can be read
		fetchedItem, err := storagePacker.GetItem(item.ID)
		if err != nil {
			t.Fatal(err)
		}
		if fetchedItem == nil {
			t.Fatalf("failed to read the stored item")
		}

		if item.ID != fetchedItem.ID {
			t.Fatalf("bad: item ID; expected: %q\n actual: %q\n", item.ID, fetchedItem.ID)
		}

		itemsToDelete = append(itemsToDelete, fmt.Sprintf("item%d", i))
	}

	err = storagePacker.DeleteMultipleItems(ctx, nil, itemsToDelete)
	if err != nil {
		t.Fatal(err)
	}

	// Check that the deletion was successful
	for _, item := range itemsToDelete {
		fetchedItem, err := storagePacker.GetItem(item)
		if err != nil {
			t.Fatal(err)
		}
		if fetchedItem != nil {
			t.Fatal("item not deleted")
		}
	}
}
