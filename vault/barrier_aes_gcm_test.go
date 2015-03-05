package vault

import (
	"bytes"
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
