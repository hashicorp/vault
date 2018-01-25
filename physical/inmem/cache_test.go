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
	if err != nil {
		t.Fatal(err)
	}

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

func TestCache_Disable(t *testing.T) {
	logger := logformat.NewVaultLogger(log.LevelTrace)

	inm, err := NewInmem(nil, logger)
	if err != nil {
		t.Fatal(err)
	}
	cache := physical.NewCache(inm, 0, logger)

	disabledTests := func() {
		ent := &physical.Entry{
			Key:   "foo",
			Value: []byte("bar"),
		}
		err = inm.Put(context.Background(), ent)
		if err != nil {
			t.Fatalf("err: %v", err)
		}

		// Read should work
		out, err := cache.Get(context.Background(), "foo")
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if out == nil {
			t.Fatalf("should have key")
		}

		err = inm.Delete(context.Background(), ent.Key)
		if err != nil {
			t.Fatal(err)
		}

		// Should not work
		out, err = cache.Get(context.Background(), "foo")
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if out != nil {
			t.Fatalf("should not have key")
		}

		// Put through the cache and try again
		err = cache.Put(context.Background(), ent)
		if err != nil {
			t.Fatalf("err: %v", err)
		}

		// Read should work in both
		out, err = inm.Get(context.Background(), "foo")
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if out == nil {
			t.Fatalf("should have key")
		}
		out, err = cache.Get(context.Background(), "foo")
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if out == nil {
			t.Fatalf("should have key")
		}

		err = inm.Delete(context.Background(), ent.Key)
		if err != nil {
			t.Fatal(err)
		}

		// Should not work
		out, err = cache.Get(context.Background(), "foo")
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if out != nil {
			t.Fatalf("should not have key")
		}
	}

	enabledTests := func() {
		ent := &physical.Entry{
			Key:   "foo",
			Value: []byte("bar"),
		}
		err = inm.Put(context.Background(), ent)
		if err != nil {
			t.Fatalf("err: %v", err)
		}

		// Read should work
		out, err := cache.Get(context.Background(), "foo")
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if out == nil {
			t.Fatalf("should have key")
		}

		err = inm.Delete(context.Background(), ent.Key)
		if err != nil {
			t.Fatal(err)
		}

		// Should work
		out, err = cache.Get(context.Background(), "foo")
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if out == nil {
			t.Fatalf("should have key")
		}

		// Put through the cache and try again
		err = cache.Put(context.Background(), ent)
		if err != nil {
			t.Fatalf("err: %v", err)
		}

		// Read should work for both
		out, err = inm.Get(context.Background(), "foo")
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if out == nil {
			t.Fatalf("should have key")
		}
		out, err = cache.Get(context.Background(), "foo")
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if out == nil {
			t.Fatalf("should have key")
		}

		err = inm.Delete(context.Background(), ent.Key)
		if err != nil {
			t.Fatal(err)
		}

		// Should work
		out, err = cache.Get(context.Background(), "foo")
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if out == nil {
			t.Fatalf("should have key")
		}

		// Put through the cache
		err = cache.Put(context.Background(), ent)
		if err != nil {
			t.Fatalf("err: %v", err)
		}

		// Read should work for both
		out, err = inm.Get(context.Background(), "foo")
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if out == nil {
			t.Fatalf("should have key")
		}
		out, err = cache.Get(context.Background(), "foo")
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if out == nil {
			t.Fatalf("should have key")
		}

		// Delete via cache
		err = cache.Delete(context.Background(), ent.Key)
		if err != nil {
			t.Fatal(err)
		}

		// Read should not work for either
		out, err = inm.Get(context.Background(), "foo")
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if out != nil {
			t.Fatalf("should not have key")
		}
		out, err = cache.Get(context.Background(), "foo")
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if out != nil {
			t.Fatalf("should not have key")
		}
	}

	disabledTests()
	cache.SetEnabled(true)
	enabledTests()
	cache.SetEnabled(false)
	disabledTests()
}
