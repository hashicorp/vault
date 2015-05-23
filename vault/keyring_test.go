package vault

import (
	"bytes"
	"reflect"
	"testing"
)

func TestKeyring(t *testing.T) {
	k := NewKeyring()

	// Term should be 0
	if term := k.ActiveTerm(); term != 0 {
		t.Fatalf("bad: %d", term)
	}

	// Should have no key
	if key := k.ActiveKey(); key != nil {
		t.Fatalf("bad: %v", key)
	}

	// Add a key
	testKey := []byte("testing")
	err := k.AddKey(1, testKey)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Term should be 1
	if term := k.ActiveTerm(); term != 1 {
		t.Fatalf("bad: %d", term)
	}

	// Should have key
	key := k.ActiveKey()
	if key == nil {
		t.Fatalf("bad: %v", key)
	}
	if !bytes.Equal(key.Value, testKey) {
		t.Fatalf("bad: %v", key)
	}
	if tKey := k.TermKey(1); tKey != key {
		t.Fatalf("bad: %v", tKey)
	}

	// Should handle idempotent set
	err = k.AddKey(1, testKey)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Should not allow conficting set
	testConflict := []byte("nope")
	err = k.AddKey(1, testConflict)
	if err == nil {
		t.Fatalf("err: %v", err)
	}

	// Add a new key
	testSecond := []byte("second")
	err = k.AddKey(2, testSecond)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Term should be 2
	if term := k.ActiveTerm(); term != 2 {
		t.Fatalf("bad: %d", term)
	}

	// Should have key
	newKey := k.ActiveKey()
	if newKey == nil {
		t.Fatalf("bad: %v", key)
	}
	if !bytes.Equal(newKey.Value, testSecond) {
		t.Fatalf("bad: %v", key)
	}
	if tKey := k.TermKey(2); tKey != newKey {
		t.Fatalf("bad: %v", tKey)
	}

	// Read of old key should work
	if tKey := k.TermKey(1); tKey != key {
		t.Fatalf("bad: %v", tKey)
	}

	// Remove the old key
	err = k.RemoveKey(1)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Read of old key should not work
	if tKey := k.TermKey(1); tKey != nil {
		t.Fatalf("bad: %v", tKey)
	}

	// Remove the active key should fail
	err = k.RemoveKey(2)
	if err == nil {
		t.Fatalf("err: %v", err)
	}
}

func TestKeyring_MasterKey(t *testing.T) {
	k := NewKeyring()
	master := []byte("test")
	master2 := []byte("test2")

	// Check no master
	out := k.MasterKey()
	if out != nil {
		t.Fatalf("bad: %v", out)
	}

	// Set master
	k.SetMasterKey(master)
	out = k.MasterKey()
	if !bytes.Equal(out, master) {
		t.Fatalf("bad: %v", out)
	}

	// Update master
	k.SetMasterKey(master2)
	out = k.MasterKey()
	if !bytes.Equal(out, master2) {
		t.Fatalf("bad: %v", out)
	}
}

func TestKeyring_Serialize(t *testing.T) {
	k := NewKeyring()
	master := []byte("test")
	k.SetMasterKey(master)

	testKey := []byte("testing")
	testSecond := []byte("second")
	k.AddKey(1, testKey)
	k.AddKey(2, testSecond)

	buf, err := k.Serialize()
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	k2, err := DeserializeKeyring(buf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	out := k2.MasterKey()
	if !bytes.Equal(out, master) {
		t.Fatalf("bad: %v", out)
	}

	if k2.ActiveTerm() != k.ActiveTerm() {
		t.Fatalf("Term mismatch")
	}

	var i uint32
	for i = 1; i < k.ActiveTerm(); i++ {
		key1 := k2.TermKey(i)
		key2 := k.TermKey(i)
		if !reflect.DeepEqual(key1, key2) {
			t.Fatalf("bad: %v %v", key1, key2)
		}
	}
}
