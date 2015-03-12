package vault

import (
	"reflect"
	"testing"
)

func testBarrier(t *testing.T, b SecurityBarrier) {
	// Should not be initialized
	init, err := b.Initialized()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if init {
		t.Fatalf("should not be initialized")
	}

	// Should start sealed
	sealed, err := b.Sealed()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !sealed {
		t.Fatalf("should be sealed")
	}

	// Sealing should be a no-op
	if err := b.Seal(); err != nil {
		t.Fatalf("err: %v", err)
	}

	// All operations should fail
	e := &Entry{Key: "test", Value: []byte("test")}
	if err := b.Put(e); err != ErrBarrierSealed {
		t.Fatalf("err: %v", err)
	}
	if _, err := b.Get("test"); err != ErrBarrierSealed {
		t.Fatalf("err: %v", err)
	}
	if err := b.Delete("test"); err != ErrBarrierSealed {
		t.Fatalf("err: %v", err)
	}
	if _, err := b.List(""); err != ErrBarrierSealed {
		t.Fatalf("err: %v", err)
	}

	// Get a new key
	key, err := b.GenerateKey()
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Validate minimum key length
	min, max := b.KeyLength()
	if min < 16 {
		t.Fatalf("minimum key size too small: %d", min)
	}
	if max < min {
		t.Fatalf("maximum key size smaller than min")
	}

	// Unseal should not work
	if err := b.Unseal(key); err != ErrBarrierNotInit {
		t.Fatalf("err: %v", err)
	}

	// Initialize the vault
	if err := b.Initialize(key); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Double Initialize should fail
	if err := b.Initialize(key); err != ErrBarrierAlreadyInit {
		t.Fatalf("err: %v", err)
	}

	// Should be initialized
	init, err = b.Initialized()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !init {
		t.Fatalf("should be initialized")
	}

	// Should still be sealed
	sealed, err = b.Sealed()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !sealed {
		t.Fatalf("should sealed")
	}

	// Unseal should work
	if err := b.Unseal(key); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Unseal should no-op when done twice
	if err := b.Unseal(key); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Should no longer be sealed
	sealed, err = b.Sealed()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if sealed {
		t.Fatalf("should be unsealed")
	}

	// Operations should work
	out, err := b.Get("test")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out != nil {
		t.Fatalf("bad: %v", out)
	}

	// List should have only "barrier/"
	keys, err := b.List("")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(keys) != 1 || keys[0] != "barrier/" {
		t.Fatalf("bad: %v", keys)
	}

	// Try to write
	if err := b.Put(e); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Should be equal
	out, err = b.Get("test")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !reflect.DeepEqual(out, e) {
		t.Fatalf("bad: %v exp: %v", out, e)
	}

	// List should show the items
	keys, err = b.List("")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(keys) != 2 {
		t.Fatalf("bad: %v", keys)
	}
	if keys[0] != "barrier/" || keys[1] != "test" {
		t.Fatalf("bad: %v", keys)
	}

	// Delete should clear
	err = b.Delete("test")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Double Delete is fine
	err = b.Delete("test")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Should be nil
	out, err = b.Get("test")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out != nil {
		t.Fatalf("bad: %v", out)
	}

	// List should have nothing
	keys, err = b.List("")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(keys) != 1 || keys[0] != "barrier/" {
		t.Fatalf("bad: %v", keys)
	}

	// Add the item back
	if err := b.Put(e); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Reseal should prevent any updates
	if err := b.Seal(); err != nil {
		t.Fatalf("err: %v", err)
	}

	// No access allowed
	if _, err := b.Get("test"); err != ErrBarrierSealed {
		t.Fatalf("err: %v", err)
	}

	// Unseal should work
	if err := b.Unseal(key); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Should be equal
	out, err = b.Get("test")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !reflect.DeepEqual(out, e) {
		t.Fatalf("bad: %v exp: %v", out, e)
	}

	// Final cleanup
	err = b.Delete("test")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Reseal should prevent any updates
	if err := b.Seal(); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Modify the key
	key[0]++

	// Unseal should fail
	if err := b.Unseal(key); err != ErrBarrierInvalidKey {
		t.Fatalf("err: %v", err)
	}
}
