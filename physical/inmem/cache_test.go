package inmem

import (
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

	ent := &physical.Entry{
		Key:   "foo",
		Value: []byte("bar"),
	}
	err = cache.Put(ent)
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

func TestCache_IgnoreCore(t *testing.T) {
	logger := logformat.NewVaultLogger(log.LevelTrace)

	inm, err := NewInmem(nil, logger)
	if err != nil {
		t.Fatal(err)
	}

	cache := physical.NewCache(inm, 0, logger)

	var ent *physical.Entry

	// First try normal handling
	ent = &physical.Entry{
		Key:   "foo",
		Value: []byte("bar"),
	}
	if err := cache.Put(ent); err != nil {
		t.Fatal(err)
	}
	ent = &physical.Entry{
		Key:   "foo",
		Value: []byte("foobar"),
	}
	if err := inm.Put(ent); err != nil {
		t.Fatal(err)
	}
	ent, err = cache.Get("foo")
	if err != nil {
		t.Fatal(err)
	}
	if string(ent.Value) != "bar" {
		t.Fatal("expected cached value")
	}

	// Now try core path
	ent = &physical.Entry{
		Key:   "core/foo",
		Value: []byte("bar"),
	}
	if err := cache.Put(ent); err != nil {
		t.Fatal(err)
	}
	ent = &physical.Entry{
		Key:   "core/foo",
		Value: []byte("foobar"),
	}
	if err := inm.Put(ent); err != nil {
		t.Fatal(err)
	}
	ent, err = cache.Get("core/foo")
	if err != nil {
		t.Fatal(err)
	}
	if string(ent.Value) != "foobar" {
		t.Fatal("expected cached value")
	}

	// Now make sure looked-up values aren't added
	ent = &physical.Entry{
		Key:   "core/zip",
		Value: []byte("zap"),
	}
	if err := inm.Put(ent); err != nil {
		t.Fatal(err)
	}
	ent, err = cache.Get("core/zip")
	if err != nil {
		t.Fatal(err)
	}
	if string(ent.Value) != "zap" {
		t.Fatal("expected non-cached value")
	}
	ent = &physical.Entry{
		Key:   "core/zip",
		Value: []byte("zipzap"),
	}
	if err := inm.Put(ent); err != nil {
		t.Fatal(err)
	}
	ent, err = cache.Get("core/zip")
	if err != nil {
		t.Fatal(err)
	}
	if string(ent.Value) != "zipzap" {
		t.Fatal("expected non-cached value")
	}
}
