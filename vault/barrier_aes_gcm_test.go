package vault

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/hashicorp/vault/physical"
)

// mockBarrier returns a physical backend, security barrier, and master key
func mockBarrier(t *testing.T) (physical.Backend, SecurityBarrier, []byte) {
	inm := physical.NewInmem()
	b, err := NewAESGCMBarrier(inm)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Initialize and unseal
	key, _ := b.GenerateKey()
	b.Initialize(key)
	b.Unseal(key)
	return inm, b, key
}

func TestAESGCMBarrier_Basic(t *testing.T) {
	inm := physical.NewInmem()
	b, err := NewAESGCMBarrier(inm)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	testBarrier(t, b)
}

func TestAESGCMBarrier_Rotate(t *testing.T) {
	inm := physical.NewInmem()
	b, err := NewAESGCMBarrier(inm)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	testBarrier_Rotate(t, b)
}

func TestAESGCMBarrier_Upgrade(t *testing.T) {
	inm := physical.NewInmem()
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
	inm := physical.NewInmem()
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
	inm := physical.NewInmem()
	b, err := NewAESGCMBarrier(inm)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	testBarrier_Rekey(t, b)
}

// Test an upgrade from the old (0.1) barrier/init to the new
// core/keyring style
func TestAESGCMBarrier_BackwardsCompatible(t *testing.T) {
	inm := physical.NewInmem()
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
	value := b.encrypt(initialKeyTerm, gcm, buf)

	// Write to the physical backend
	pe := &physical.Entry{
		Key:   barrierInitPath,
		Value: value,
	}
	inm.Put(pe)

	// Create a fake key
	gcm, _ = b.aeadFromKey(encrypt)
	pe = &physical.Entry{
		Key:   "test/foo",
		Value: b.encrypt(initialKeyTerm, gcm, []byte("test")),
	}
	inm.Put(pe)

	// Should still be initialized
	isInit, err := b.Initialized()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !isInit {
		t.Fatalf("should be initialized")
	}

	// Unseal should work and migrate online
	err = b.Unseal(master)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Check for migraiton
	out, err := inm.Get(barrierInitPath)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out != nil {
		t.Fatalf("should delete old barrier init")
	}

	// Should have keyring
	out, err = inm.Get(keyringPath)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out == nil {
		t.Fatalf("should have keyring file")
	}

	// Attempt to read encrypted key
	entry, err := b.Get("test/foo")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if string(entry.Value) != "test" {
		t.Fatalf("bad: %#v", entry)
	}
}

// Verify data sent through is encrypted
func TestAESGCMBarrier_Confidential(t *testing.T) {
	inm := physical.NewInmem()
	b, err := NewAESGCMBarrier(inm)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Initialize and unseal
	key, _ := b.GenerateKey()
	b.Initialize(key)
	b.Unseal(key)

	// Put a logical entry
	entry := &Entry{Key: "test", Value: []byte("test")}
	err = b.Put(entry)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Check the physcial entry
	pe, err := inm.Get("test")
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

// Verify data sent through is cannot be tampered
func TestAESGCMBarrier_Integrity(t *testing.T) {
	inm := physical.NewInmem()
	b, err := NewAESGCMBarrier(inm)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Initialize and unseal
	key, _ := b.GenerateKey()
	b.Initialize(key)
	b.Unseal(key)

	// Put a logical entry
	entry := &Entry{Key: "test", Value: []byte("test")}
	err = b.Put(entry)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Change a byte in the underlying physical entry
	pe, _ := inm.Get("test")
	pe.Value[15]++
	err = inm.Put(pe)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Read from the barrier
	_, err = b.Get("test")
	if err == nil {
		t.Fatalf("should fail!")
	}
}

func TestEncrypt_Unique(t *testing.T) {
	inm := physical.NewInmem()
	b, err := NewAESGCMBarrier(inm)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	key, _ := b.GenerateKey()
	b.Initialize(key)
	b.Unseal(key)

	if b.keyring == nil {
		t.Fatalf("barrier is sealed")
	}

	entry := &Entry{Key: "test", Value: []byte("test")}
	term := b.keyring.ActiveTerm()
	primary, _ := b.aeadForTerm(term)

	first := b.encrypt(term, primary, entry.Value)
	second := b.encrypt(term, primary, entry.Value)

	if bytes.Equal(first, second) == true {
		t.Fatalf("improper random seeding detected")
	}
}

func TestInitialize_KeyLength(t *testing.T) {
	inm := physical.NewInmem()
	b, err := NewAESGCMBarrier(inm)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	long := []byte("ThisKeyDoesNotHaveTheRightLength!")
	middle := []byte("ThisIsASecretKeyAndMore")
	short := []byte("Key")

	err = b.Initialize(long)

	if err == nil {
		t.Fatalf("key length protection failed")
	}

	err = b.Initialize(middle)

	if err == nil {
		t.Fatalf("key length protection failed")
	}

	err = b.Initialize(short)

	if err == nil {
		t.Fatalf("key length protection failed")
	}
}
