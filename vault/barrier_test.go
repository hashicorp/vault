package vault

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/vault/logical"
)

func testBarrier(t *testing.T, b SecurityBarrier) {
	// Should not be initialized
	init, err := b.Initialized(context.Background())
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
	e := &logical.StorageEntry{Key: "test", Value: []byte("test")}
	if err := b.Put(context.Background(), e); err != ErrBarrierSealed {
		t.Fatalf("err: %v", err)
	}
	if _, err := b.Get(context.Background(), "test"); err != ErrBarrierSealed {
		t.Fatalf("err: %v", err)
	}
	if err := b.Delete(context.Background(), "test"); err != ErrBarrierSealed {
		t.Fatalf("err: %v", err)
	}
	if _, err := b.List(context.Background(), ""); err != ErrBarrierSealed {
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
	if err := b.Unseal(context.Background(), key); err != ErrBarrierNotInit {
		t.Fatalf("err: %v", err)
	}

	// Initialize the vault
	if err := b.Initialize(context.Background(), key); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Double Initialize should fail
	if err := b.Initialize(context.Background(), key); err != ErrBarrierAlreadyInit {
		t.Fatalf("err: %v", err)
	}

	// Should be initialized
	init, err = b.Initialized(context.Background())
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
	if err := b.Unseal(context.Background(), key); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Unseal should no-op when done twice
	if err := b.Unseal(context.Background(), key); err != nil {
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

	// Verify the master key
	if err := b.VerifyMaster(key); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Operations should work
	out, err := b.Get(context.Background(), "test")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out != nil {
		t.Fatalf("bad: %v", out)
	}

	// List should have only "core/"
	keys, err := b.List(context.Background(), "")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(keys) != 1 || keys[0] != "core/" {
		t.Fatalf("bad: %v", keys)
	}

	// Try to write
	if err := b.Put(context.Background(), e); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Should be equal
	out, err = b.Get(context.Background(), "test")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !reflect.DeepEqual(out, e) {
		t.Fatalf("bad: %v exp: %v", out, e)
	}

	// List should show the items
	keys, err = b.List(context.Background(), "")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(keys) != 2 {
		t.Fatalf("bad: %v", keys)
	}
	if keys[0] != "core/" || keys[1] != "test" {
		t.Fatalf("bad: %v", keys)
	}

	// Delete should clear
	err = b.Delete(context.Background(), "test")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Double Delete is fine
	err = b.Delete(context.Background(), "test")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Should be nil
	out, err = b.Get(context.Background(), "test")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out != nil {
		t.Fatalf("bad: %v", out)
	}

	// List should have nothing
	keys, err = b.List(context.Background(), "")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(keys) != 1 || keys[0] != "core/" {
		t.Fatalf("bad: %v", keys)
	}

	// Add the item back
	if err := b.Put(context.Background(), e); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Reseal should prevent any updates
	if err := b.Seal(); err != nil {
		t.Fatalf("err: %v", err)
	}

	// No access allowed
	if _, err := b.Get(context.Background(), "test"); err != ErrBarrierSealed {
		t.Fatalf("err: %v", err)
	}

	// Unseal should work
	if err := b.Unseal(context.Background(), key); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Should be equal
	out, err = b.Get(context.Background(), "test")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !reflect.DeepEqual(out, e) {
		t.Fatalf("bad: %v exp: %v", out, e)
	}

	// Final cleanup
	err = b.Delete(context.Background(), "test")
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
	if err := b.Unseal(context.Background(), key); err != ErrBarrierInvalidKey {
		t.Fatalf("err: %v", err)
	}
}

func testBarrier_Rotate(t *testing.T, b SecurityBarrier) {
	// Initialize the barrier
	key, _ := b.GenerateKey()
	b.Initialize(context.Background(), key)
	err := b.Unseal(context.Background(), key)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Check the key info
	info, err := b.ActiveKeyInfo()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if info.Term != 1 {
		t.Fatalf("Bad term: %d", info.Term)
	}
	if time.Since(info.InstallTime) > time.Second {
		t.Fatalf("Bad install: %v", info.InstallTime)
	}
	first := info.InstallTime

	// Write a key
	e1 := &logical.StorageEntry{Key: "test", Value: []byte("test")}
	if err := b.Put(context.Background(), e1); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Rotate the encryption key
	newTerm, err := b.Rotate(context.Background())
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if newTerm != 2 {
		t.Fatalf("bad: %v", newTerm)
	}

	// Check the key info
	info, err = b.ActiveKeyInfo()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if info.Term != 2 {
		t.Fatalf("Bad term: %d", info.Term)
	}
	if !info.InstallTime.After(first) {
		t.Fatalf("Bad install: %v", info.InstallTime)
	}

	// Write another key
	e2 := &logical.StorageEntry{Key: "foo", Value: []byte("test")}
	if err := b.Put(context.Background(), e2); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Reading both should work
	out, err := b.Get(context.Background(), e1.Key)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out == nil {
		t.Fatalf("bad: %v", out)
	}

	out, err = b.Get(context.Background(), e2.Key)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out == nil {
		t.Fatalf("bad: %v", out)
	}

	// Seal and unseal
	err = b.Seal()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	err = b.Unseal(context.Background(), key)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Reading both should work
	out, err = b.Get(context.Background(), e1.Key)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out == nil {
		t.Fatalf("bad: %v", out)
	}

	out, err = b.Get(context.Background(), e2.Key)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out == nil {
		t.Fatalf("bad: %v", out)
	}

	// Should be fine to reload keyring
	err = b.ReloadKeyring(context.Background())
	if err != nil {
		t.Fatalf("err: %v", err)
	}
}

func testBarrier_Rekey(t *testing.T, b SecurityBarrier) {
	// Initialize the barrier
	key, _ := b.GenerateKey()
	b.Initialize(context.Background(), key)
	err := b.Unseal(context.Background(), key)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Write a key
	e1 := &logical.StorageEntry{Key: "test", Value: []byte("test")}
	if err := b.Put(context.Background(), e1); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Verify the master key
	if err := b.VerifyMaster(key); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Rekey to a new key
	newKey, _ := b.GenerateKey()
	err = b.Rekey(context.Background(), newKey)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Verify the old master key
	if err := b.VerifyMaster(key); err != ErrBarrierInvalidKey {
		t.Fatalf("err: %v", err)
	}

	// Verify the new master key
	if err := b.VerifyMaster(newKey); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Reading should work
	out, err := b.Get(context.Background(), e1.Key)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out == nil {
		t.Fatalf("bad: %v", out)
	}

	// Seal
	err = b.Seal()
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Unseal with old key should fail
	err = b.Unseal(context.Background(), key)
	if err == nil {
		t.Fatalf("unseal should fail")
	}

	// Unseal with new keys should work
	err = b.Unseal(context.Background(), newKey)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Reading should work
	out, err = b.Get(context.Background(), e1.Key)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out == nil {
		t.Fatalf("bad: %v", out)
	}

	// Should be fine to reload keyring
	err = b.ReloadKeyring(context.Background())
	if err != nil {
		t.Fatalf("err: %v", err)
	}
}

func testBarrier_Upgrade(t *testing.T, b1, b2 SecurityBarrier) {
	// Initialize the barrier
	key, _ := b1.GenerateKey()
	b1.Initialize(context.Background(), key)
	err := b1.Unseal(context.Background(), key)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	err = b2.Unseal(context.Background(), key)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Rotate the encryption key
	newTerm, err := b1.Rotate(context.Background())
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Create upgrade path
	err = b1.CreateUpgrade(context.Background(), newTerm)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Check for an upgrade
	did, updated, err := b2.CheckUpgrade(context.Background())
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !did || updated != newTerm {
		t.Fatalf("failed to upgrade")
	}

	// Should have no upgrades pending
	did, updated, err = b2.CheckUpgrade(context.Background())
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if did {
		t.Fatalf("should not have upgrade")
	}

	// Rotate the encryption key
	newTerm, err = b1.Rotate(context.Background())
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Create upgrade path
	err = b1.CreateUpgrade(context.Background(), newTerm)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Destroy upgrade path
	err = b1.DestroyUpgrade(context.Background(), newTerm)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Should have no upgrades pending
	did, updated, err = b2.CheckUpgrade(context.Background())
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if did {
		t.Fatalf("should not have upgrade")
	}
}

func testBarrier_Upgrade_Rekey(t *testing.T, b1, b2 SecurityBarrier) {
	// Initialize the barrier
	key, _ := b1.GenerateKey()
	b1.Initialize(context.Background(), key)
	err := b1.Unseal(context.Background(), key)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	err = b2.Unseal(context.Background(), key)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Rekey to a new key
	newKey, _ := b1.GenerateKey()
	err = b1.Rekey(context.Background(), newKey)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Reload the master key
	err = b2.ReloadMasterKey(context.Background())
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Reload the keyring
	err = b2.ReloadKeyring(context.Background())
	if err != nil {
		t.Fatalf("err: %v", err)
	}
}
