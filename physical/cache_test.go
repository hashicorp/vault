package physical

import (
	"log"
	"os"
	"testing"
)

func TestCache(t *testing.T) {
	logger := log.New(os.Stderr, "", log.LstdFlags)
	inm := NewInmem(logger)
	cache := NewCache(inm, 0)
	testBackend(t, cache)
	testBackend_ListPrefix(t, cache)
}

func TestCache_Purge(t *testing.T) {
	logger := log.New(os.Stderr, "", log.LstdFlags)
	inm := NewInmem(logger)
	cache := NewCache(inm, 0)

	ent := &Entry{
		Key:   "foo",
		Value: []byte("bar"),
	}
	err := cache.Put(ent)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Delete from under
	inm.Delete("foo")

	// Read should work
	out, err := cache.Get("foo")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out == nil {
		t.Fatalf("should have key")
	}

	// Clear the cache
	cache.Purge()

	// Read should fail
	out, err = cache.Get("foo")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out != nil {
		t.Fatalf("should not have key")
	}
}
