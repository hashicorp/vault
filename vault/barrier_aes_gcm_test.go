package vault

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/logging"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/physical"
	"github.com/hashicorp/vault/physical/inmem"
)

var (
	logger = logging.NewVaultLogger(log.Trace)
)

// mockBarrier returns a physical backend, security barrier, and master key
func mockBarrier(t testing.TB) (physical.Backend, SecurityBarrier, []byte) {
	inm, err := inmem.NewInmem(nil, logger)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	b, err := NewAESGCMBarrier(inm)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Initialize and unseal
	key, _ := b.GenerateKey()
	b.Initialize(context.Background(), key)
	b.Unseal(context.Background(), key)
	return inm, b, key
}

func TestAESGCMBarrier_Basic(t *testing.T) {
	inm, err := inmem.NewInmem(nil, logger)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	b, err := NewAESGCMBarrier(inm)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	testBarrier(t, b)
}

func TestAESGCMBarrier_Rotate(t *testing.T) {
	inm, err := inmem.NewInmem(nil, logger)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	b, err := NewAESGCMBarrier(inm)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	testBarrier_Rotate(t, b)
}

func TestAESGCMBarrier_Upgrade(t *testing.T) {
	inm, err := inmem.NewInmem(nil, logger)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	b1, err := NewAESGCMBarrier(inm)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	b2, err := NewAESGCMBarrier(inm)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	testBarrier_Upgrade(t, b1, b2)
}

func TestAESGCMBarrier_Upgrade_Rekey(t *testing.T) {
	inm, err := inmem.NewInmem(nil, logger)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	b1, err := NewAESGCMBarrier(inm)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	b2, err := NewAESGCMBarrier(inm)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	testBarrier_Upgrade_Rekey(t, b1, b2)
}

func TestAESGCMBarrier_Rekey(t *testing.T) {
	inm, err := inmem.NewInmem(nil, logger)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	b, err := NewAESGCMBarrier(inm)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	testBarrier_Rekey(t, b)
}

// Test an upgrade from the old (0.1) barrier/init to the new
// core/keyring style
func TestAESGCMBarrier_BackwardsCompatible(t *testing.T) {
	inm, err := inmem.NewInmem(nil, logger)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	b, err := NewAESGCMBarrier(inm)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Generate a barrier/init entry
	encrypt, _ := b.GenerateKey()
	init := &barrierInit{
		Version: 1,
		Key:     encrypt,
	}
	buf, _ := json.Marshal(init)

	// Protect with master key
	master, _ := b.GenerateKey()
	gcm, _ := b.aeadFromKey(master)
	value, err := b.encrypt(barrierInitPath, initialKeyTerm, gcm, buf)
	if err != nil {
		t.Fatal(err)
	}

	// Write to the physical backend
	pe := &physical.Entry{
		Key:   barrierInitPath,
		Value: value,
	}
	inm.Put(context.Background(), pe)

	// Create a fake key
	gcm, _ = b.aeadFromKey(encrypt)
	value, err = b.encrypt("test/foo", initialKeyTerm, gcm, []byte("test"))
	if err != nil {
		t.Fatal(err)
	}
	pe = &physical.Entry{
		Key:   "test/foo",
		Value: value,
	}
	inm.Put(context.Background(), pe)

	// Should still be initialized
	isInit, err := b.Initialized(context.Background())
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !isInit {
		t.Fatalf("should be initialized")
	}

	// Unseal should work and migrate online
	err = b.Unseal(context.Background(), master)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Check for migration
	out, err := inm.Get(context.Background(), barrierInitPath)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out != nil {
		t.Fatalf("should delete old barrier init")
	}

	// Should have keyring
	out, err = inm.Get(context.Background(), keyringPath)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out == nil {
		t.Fatalf("should have keyring file")
	}

	// Attempt to read encrypted key
	entry, err := b.Get(context.Background(), "test/foo")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if string(entry.Value) != "test" {
		t.Fatalf("bad: %#v", entry)
	}
}

// Verify data sent through is encrypted
func TestAESGCMBarrier_Confidential(t *testing.T) {
	inm, err := inmem.NewInmem(nil, logger)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	b, err := NewAESGCMBarrier(inm)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Initialize and unseal
	key, _ := b.GenerateKey()
	b.Initialize(context.Background(), key)
	b.Unseal(context.Background(), key)

	// Put a logical entry
	entry := &logical.StorageEntry{Key: "test", Value: []byte("test")}
	err = b.Put(context.Background(), entry)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Check the physical entry
	pe, err := inm.Get(context.Background(), "test")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if pe == nil {
		t.Fatalf("missing physical entry")
	}

	if pe.Key != "test" {
		t.Fatalf("bad: %#v", pe)
	}
	if bytes.Equal(pe.Value, entry.Value) {
		t.Fatalf("bad: %#v", pe)
	}
}

// Verify data sent through cannot be tampered with
func TestAESGCMBarrier_Integrity(t *testing.T) {
	inm, err := inmem.NewInmem(nil, logger)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	b, err := NewAESGCMBarrier(inm)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Initialize and unseal
	key, _ := b.GenerateKey()
	b.Initialize(context.Background(), key)
	b.Unseal(context.Background(), key)

	// Put a logical entry
	entry := &logical.StorageEntry{Key: "test", Value: []byte("test")}
	err = b.Put(context.Background(), entry)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Change a byte in the underlying physical entry
	pe, _ := inm.Get(context.Background(), "test")
	pe.Value[15]++
	err = inm.Put(context.Background(), pe)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Read from the barrier
	_, err = b.Get(context.Background(), "test")
	if err == nil {
		t.Fatalf("should fail!")
	}
}

// Verify data sent through cannot be moved
func TestAESGCMBarrier_MoveIntegrityV1(t *testing.T) {
	inm, err := inmem.NewInmem(nil, logger)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	b, err := NewAESGCMBarrier(inm)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	b.currentAESGCMVersionByte = AESGCMVersion1

	// Initialize and unseal
	key, _ := b.GenerateKey()
	err = b.Initialize(context.Background(), key)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	err = b.Unseal(context.Background(), key)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Put a logical entry
	entry := &logical.StorageEntry{Key: "test", Value: []byte("test")}
	err = b.Put(context.Background(), entry)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Change the location of the underlying physical entry
	pe, _ := inm.Get(context.Background(), "test")
	pe.Key = "moved"
	err = inm.Put(context.Background(), pe)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Read from the barrier
	_, err = b.Get(context.Background(), "moved")
	if err != nil {
		t.Fatalf("should succeed with version 1!")
	}
}

func TestAESGCMBarrier_MoveIntegrityV2(t *testing.T) {
	inm, err := inmem.NewInmem(nil, logger)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	b, err := NewAESGCMBarrier(inm)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	b.currentAESGCMVersionByte = AESGCMVersion2

	// Initialize and unseal
	key, _ := b.GenerateKey()
	err = b.Initialize(context.Background(), key)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	err = b.Unseal(context.Background(), key)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Put a logical entry
	entry := &logical.StorageEntry{Key: "test", Value: []byte("test")}
	err = b.Put(context.Background(), entry)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Change the location of the underlying physical entry
	pe, _ := inm.Get(context.Background(), "test")
	pe.Key = "moved"
	err = inm.Put(context.Background(), pe)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Read from the barrier
	_, err = b.Get(context.Background(), "moved")
	if err == nil {
		t.Fatalf("should fail with version 2!")
	}
}

func TestAESGCMBarrier_UpgradeV1toV2(t *testing.T) {
	inm, err := inmem.NewInmem(nil, logger)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	b, err := NewAESGCMBarrier(inm)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	b.currentAESGCMVersionByte = AESGCMVersion1

	// Initialize and unseal
	key, _ := b.GenerateKey()
	err = b.Initialize(context.Background(), key)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	err = b.Unseal(context.Background(), key)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Put a logical entry
	entry := &logical.StorageEntry{Key: "test", Value: []byte("test")}
	err = b.Put(context.Background(), entry)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Seal
	err = b.Seal()
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Open again as version 2
	b, err = NewAESGCMBarrier(inm)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	b.currentAESGCMVersionByte = AESGCMVersion2

	// Unseal
	err = b.Unseal(context.Background(), key)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Check successful decryption
	_, err = b.Get(context.Background(), "test")
	if err != nil {
		t.Fatalf("Upgrade unsuccessful")
	}
}

func TestEncrypt_Unique(t *testing.T) {
	inm, err := inmem.NewInmem(nil, logger)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	b, err := NewAESGCMBarrier(inm)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	key, _ := b.GenerateKey()
	b.Initialize(context.Background(), key)
	b.Unseal(context.Background(), key)

	if b.keyring == nil {
		t.Fatalf("barrier is sealed")
	}

	entry := &logical.StorageEntry{Key: "test", Value: []byte("test")}
	term := b.keyring.ActiveTerm()
	primary, _ := b.aeadForTerm(term)

	first, err := b.encrypt("test", term, primary, entry.Value)
	if err != nil {
		t.Fatal(err)
	}
	second, err := b.encrypt("test", term, primary, entry.Value)
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Equal(first, second) == true {
		t.Fatalf("improper random seeding detected")
	}
}

func TestInitialize_KeyLength(t *testing.T) {
	inm, err := inmem.NewInmem(nil, logger)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	b, err := NewAESGCMBarrier(inm)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	long := []byte("ThisKeyDoesNotHaveTheRightLength!")
	middle := []byte("ThisIsASecretKeyAndMore")
	short := []byte("Key")

	err = b.Initialize(context.Background(), long)

	if err == nil {
		t.Fatalf("key length protection failed")
	}

	err = b.Initialize(context.Background(), middle)

	if err == nil {
		t.Fatalf("key length protection failed")
	}

	err = b.Initialize(context.Background(), short)

	if err == nil {
		t.Fatalf("key length protection failed")
	}
}

func TestEncrypt_BarrierEncryptor(t *testing.T) {
	inm, err := inmem.NewInmem(nil, logger)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	b, err := NewAESGCMBarrier(inm)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Initialize and unseal
	key, _ := b.GenerateKey()
	b.Initialize(context.Background(), key)
	b.Unseal(context.Background(), key)

	cipher, err := b.Encrypt(context.Background(), "foo", []byte("quick brown fox"))
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	plain, err := b.Decrypt(context.Background(), "foo", cipher)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if string(plain) != "quick brown fox" {
		t.Fatalf("bad: %s", plain)
	}
}
