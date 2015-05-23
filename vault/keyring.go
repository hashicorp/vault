package vault

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"
)

// Keyring is used to manage multiple encryption keys used by
// the barrier. New keys can be installed and each has a sequential term.
// The term used to encrypt a key is prefixed to the key written out.
// All data is encrypted with the latest key, but storing the old keys
// allows for decryption of keys written previously. Along with the encryption
// keys, the keyring also tracks the master key. This is necessary so that
// when a new key is added to the keyring, we can encrypt with the master key
// and write out the new keyring.
type Keyring struct {
	masterKey  []byte
	keys       map[uint32]*Key
	activeTerm uint32
	l          sync.RWMutex
}

// EncodedKeyring is used for serialization of the keyring
type EncodedKeyring struct {
	MasterKey []byte
	Keys      []*Key
}

// Key represents a single term, along with the key used.
type Key struct {
	Term  uint32
	Value []byte
}

// NewKeyring creates a new keyring
func NewKeyring() *Keyring {
	k := &Keyring{
		keys:       make(map[uint32]*Key),
		activeTerm: 0,
	}
	return k
}

// AddKey adds a new key to the keyring
func (k *Keyring) AddKey(term uint32, value []byte) error {
	k.l.Lock()
	defer k.l.Unlock()

	// Ensure there is no confict
	if key, ok := k.keys[term]; ok {
		if !bytes.Equal(key.Value, value) {
			return fmt.Errorf("Conflicting key for term %d already installed", term)
		}
		return nil
	}

	// Install the new key
	key := &Key{
		Term:  term,
		Value: value,
	}
	k.keys[term] = key

	// Update the active term if newer
	if term > k.activeTerm {
		k.activeTerm = term
	}
	return nil
}

// RemoveKey removes a new key to the keyring
func (k *Keyring) RemoveKey(term uint32) error {
	k.l.Lock()
	defer k.l.Unlock()

	// Ensure this is not the active key
	if term == k.activeTerm {
		return fmt.Errorf("Cannot remove active key")
	}

	// Delete the key
	delete(k.keys, term)
	return nil
}

// ActiveTerm returns the currently active term
func (k *Keyring) ActiveTerm() uint32 {
	k.l.RLock()
	defer k.l.RUnlock()
	return k.activeTerm
}

// ActiveKey returns the active encryption key, or nil
func (k *Keyring) ActiveKey() *Key {
	k.l.RLock()
	defer k.l.RUnlock()
	return k.keys[k.activeTerm]
}

// TermKey returns the key for the given term, or nil
func (k *Keyring) TermKey(term uint32) *Key {
	k.l.RLock()
	defer k.l.RUnlock()
	return k.keys[term]
}

// SetMasterKey is used to update the master key
func (k *Keyring) SetMasterKey(val []byte) {
	k.l.Lock()
	defer k.l.Unlock()
	k.masterKey = val
}

// MasterKey returns the master key
func (k *Keyring) MasterKey() []byte {
	k.l.RLock()
	defer k.l.RUnlock()
	return k.masterKey
}

// Serialize is used to create a byte encoded keyring
func (k *Keyring) Serialize() ([]byte, error) {
	k.l.RLock()
	defer k.l.RUnlock()

	// Create the encoded entry
	enc := EncodedKeyring{
		MasterKey: k.masterKey,
	}
	for _, key := range k.keys {
		enc.Keys = append(enc.Keys, key)
	}

	// JSON encode the keyring
	buf, err := json.Marshal(enc)
	return buf, err
}

// DeserializeKeyring is used to deserialize and return a new keyring
func DeserializeKeyring(buf []byte) (*Keyring, error) {
	// Deserialize the keyring
	var enc EncodedKeyring
	if err := json.Unmarshal(buf, &enc); err != nil {
		return nil, fmt.Errorf("deserialization failed: %v", err)
	}

	// Create a new keyring
	k := NewKeyring()
	k.SetMasterKey(enc.MasterKey)
	for _, key := range enc.Keys {
		if err := k.AddKey(key.Term, key.Value); err != nil {
			return nil, fmt.Errorf("failed to add key for term %d: %v", key.Term, err)
		}
	}
	return k, nil
}
