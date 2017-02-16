package vault

import (
	"bytes"
	"reflect"
	"testing"
	"time"
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
	key1 := &Key{Term: 1, Version: 1, Value: testKey, InstallTime: time.Now()}
	k, err := k.AddKey(key1)
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
	k, err = k.AddKey(key1)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Should not allow conficting set
	testConflict := []byte("nope")
	key1Conf := &Key{Term: 1, Version: 1, Value: testConflict, InstallTime: time.Now()}
	_, err = k.AddKey(key1Conf)
	if err == nil {
		t.Fatalf("err: %v", err)
	}

	// Add a new key
	testSecond := []byte("second")
	key2 := &Key{Term: 2, Version: 1, Value: testSecond, InstallTime: time.Now()}
	k, err = k.AddKey(key2)
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
	k, err = k.RemoveKey(1)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Read of old key should not work
	if tKey := k.TermKey(1); tKey != nil {
		t.Fatalf("bad: %v", tKey)
	}

	// Remove the active key should fail
	k, err = k.RemoveKey(2)
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
	k = k.SetMasterKey(master)
	out = k.MasterKey()
	if !bytes.Equal(out, master) {
		t.Fatalf("bad: %v", out)
	}

	// Update master
	k = k.SetMasterKey(master2)
	out = k.MasterKey()
	if !bytes.Equal(out, master2) {
		t.Fatalf("bad: %v", out)
	}
}

func TestKeyring_Serialize(t *testing.T) {
	k := NewKeyring()
	master := []byte("test")
	k = k.SetMasterKey(master)

	now := time.Now()
	testKey := []byte("testing")
	testSecond := []byte("second")
	k, _ = k.AddKey(&Key{Term: 1, Version: 1, Value: testKey, InstallTime: now})
	k, _ = k.AddKey(&Key{Term: 2, Version: 1, Value: testSecond, InstallTime: now})

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
		// Work around timezone bug due to DeepEqual using == for comparison
		if !key1.InstallTime.Equal(key2.InstallTime) {
			t.Fatalf("bad: key 1:\n%#v\nkey 2:\n%#v", key1, key2)
		}
		key1.InstallTime = key2.InstallTime
		if !reflect.DeepEqual(key1, key2) {
			t.Fatalf("bad: key 1:\n%#v\nkey 2:\n%#v", key1, key2)
		}
	}
}

func TestKey_Serialize(t *testing.T) {
	k := &Key{
		Term:        10,
		Version:     1,
		Value:       []byte("foobarbaz"),
		InstallTime: time.Now(),
	}

	buf, err := k.Serialize()
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	out, err := DeserializeKey(buf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Work around timezone bug due to DeepEqual using == for comparison
	if !k.InstallTime.Equal(out.InstallTime) {
		t.Fatalf("bad: expected:\n%#v\nactual:\n%#v", k, out)
	}
	k.InstallTime = out.InstallTime

	if !reflect.DeepEqual(k, out) {
		t.Fatalf("bad: %#v", out)
	}
}
