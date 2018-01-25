package inmem

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/helper/logformat"
	"github.com/hashicorp/vault/physical"
	log "github.com/mgutz/logxi/v1"
)

func TestCache(t *testing.T) {
	logger := logformat.NewVaultLogger(log.LevelTrace)

	inm, err := NewInmem(nil, logger)
	if err != nil {
		t.Fatal(err)
	}
	cache := physical.NewCache(inm, 0, logger)
	physical.ExerciseBackend(t, cache)
	physical.ExerciseBackend_ListPrefix(t, cache)
}

func TestCache_Purge(t *testing.T) {
	logger := logformat.NewVaultLogger(log.LevelTrace)

	inm, err := NewInmem(nil, logger)
	if err != nil {
		t.Fatal(err)
	}
	cache := physical.NewCache(inm, 0, logger)
	cache.SetEnabled(true)

	ent := &physical.Entry{
		Key:   "foo",
		Value: []byte("bar"),
	}
	err = cache.Put(context.Background(), ent)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Delete from under
	inm.Delete(context.Background(), "foo")

	// Read should work
	out, err := cache.Get(context.Background(), "foo")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out == nil {
		t.Fatalf("should have key")
	}

	// Clear the cache
	cache.Purge(context.Background())

	// Read should fail
	out, err = cache.Get(context.Background(), "foo")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out != nil {
		t.Fatalf("should not have key")
	}
}
