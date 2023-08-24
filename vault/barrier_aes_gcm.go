// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/subtle"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/physical"
	"go.uber.org/atomic"
)

const (
	// initialKeyTerm is the hard coded initial key term. This is
	// used only for values that are not encrypted with the keyring.
	initialKeyTerm = 1

	// termSize the number of bytes used for the key term.
	termSize = 4

	autoRotateCheckInterval = 5 * time.Minute
	legacyRotateReason      = "legacy rotation"
)

// Versions of the AESGCM storage methodology
const (
	AESGCMVersion1 = 0x1
	AESGCMVersion2 = 0x2
)

// barrierInit is the JSON encoded value stored
type barrierInit struct {
	Version int    // Version is the current format version
	Key     []byte // Key is the primary encryption key
}

// Validate AESGCMBarrier satisfies SecurityBarrier interface
var (
	_                      SecurityBarrier = &AESGCMBarrier{}
	barrierEncryptsMetric                  = []string{"barrier", "estimated_encryptions"}
	barrierRotationsMetric                 = []string{"barrier", "auto_rotation"}
)

// AESGCMBarrier is a SecurityBarrier implementation that uses the AES
// cipher core and the Galois Counter Mode block mode. It defaults to
// the golang NONCE default value of 12 and a key size of 256
// bit. AES-GCM is high performance, and provides both confidentiality
// and integrity.
type AESGCMBarrier struct {
	backend physical.Backend

	l      sync.RWMutex
	sealed bool

	// keyring is used to maintain all of the encryption keys, including
	// the active key used for encryption, but also prior keys to allow
	// decryption of keys encrypted under previous terms.
	keyring *Keyring

	// cache is used to reduce the number of AEAD constructions we do
	cache     map[uint32]cipher.AEAD
	cacheLock sync.RWMutex

	// currentAESGCMVersionByte is prefixed to a message to allow for
	// future versioning of barrier implementations. It's var instead
	// of const to allow for testing
	currentAESGCMVersionByte byte

	initialized atomic.Bool

	UnaccountedEncryptions *atomic.Int64
	// Used only for testing
	RemoteEncryptions     *atomic.Int64
	totalLocalEncryptions *atomic.Int64
}

func (b *AESGCMBarrier) RotationConfig() (kc KeyRotationConfig, err error) {
	if b.keyring == nil {
		return kc, errors.New("keyring not yet present")
	}
	return b.keyring.rotationConfig.Clone(), nil
}

func (b *AESGCMBarrier) SetRotationConfig(ctx context.Context, rotConfig KeyRotationConfig) error {
	b.l.Lock()
	defer b.l.Unlock()
	rotConfig.Sanitize()
	if !rotConfig.Equals(b.keyring.rotationConfig) {
		b.keyring.rotationConfig = rotConfig

		return b.persistKeyring(ctx, b.keyring)
	}
	return nil
}

// NewAESGCMBarrier is used to construct a new barrier that uses
// the provided physical backend for storage.
func NewAESGCMBarrier(physical physical.Backend) (*AESGCMBarrier, error) {
	b := &AESGCMBarrier{
		backend:                  physical,
		sealed:                   true,
		cache:                    make(map[uint32]cipher.AEAD),
		currentAESGCMVersionByte: byte(AESGCMVersion2),
		UnaccountedEncryptions:   atomic.NewInt64(0),
		RemoteEncryptions:        atomic.NewInt64(0),
		totalLocalEncryptions:    atomic.NewInt64(0),
	}
	return b, nil
}

// Initialized checks if the barrier has been initialized
// and has a root key set.
func (b *AESGCMBarrier) Initialized(ctx context.Context) (bool, error) {
	if b.initialized.Load() {
		return true, nil
	}

	// Read the keyring file
	keys, err := b.backend.List(ctx, keyringPrefix)
	if err != nil {
		return false, fmt.Errorf("failed to check for initialization: %w", err)
	}
	if strutil.StrListContains(keys, "keyring") {
		b.initialized.Store(true)
		return true, nil
	}

	// Fallback, check for the old sentinel file
	out, err := b.backend.Get(ctx, barrierInitPath)
	if err != nil {
		return false, fmt.Errorf("failed to check for initialization: %w", err)
	}
	b.initialized.Store(out != nil)
	return out != nil, nil
}

// Initialize works only if the barrier has not been initialized
// and makes use of the given root key.
func (b *AESGCMBarrier) Initialize(ctx context.Context, key []byte, sealKey []byte, reader io.Reader) error {
	// Verify the key size
	min, max := b.KeyLength()
	if len(key) < min || len(key) > max {
		return fmt.Errorf("key size must be %d or %d", min, max)
	}

	// Check if already initialized
	if alreadyInit, err := b.Initialized(ctx); err != nil {
		return err
	} else if alreadyInit {
		return ErrBarrierAlreadyInit
	}

	// Generate encryption key
	encryptionKey, err := b.GenerateKey(reader)
	if err != nil {
		return fmt.Errorf("failed to generate encryption key: %w", err)
	}

	// Create a new keyring, install the keys
	keyring := NewKeyring()
	keyring = keyring.SetRootKey(key)
	keyring, err = keyring.AddKey(&Key{
		Term:    1,
		Version: 1,
		Value:   encryptionKey,
	})
	if err != nil {
		return fmt.Errorf("failed to create keyring: %w", err)
	}

	err = b.persistKeyring(ctx, keyring)
	if err != nil {
		return err
	}

	if len(sealKey) > 0 {
		primary, err := b.aeadFromKey(encryptionKey)
		if err != nil {
			return err
		}

		err = b.putInternal(ctx, 1, primary, &logical.StorageEntry{
			Key:   shamirKekPath,
			Value: sealKey,
		})
		if err != nil {
			return fmt.Errorf("failed to store new seal key: %w", err)
		}
	}

	return nil
}

// persistKeyring is used to write out the keyring using the
// root key to encrypt it.
func (b *AESGCMBarrier) persistKeyring(ctx context.Context, keyring *Keyring) error {
	// Create the keyring entry
	keyringBuf, err := keyring.Serialize()
	defer memzero(keyringBuf)
	if err != nil {
		return fmt.Errorf("failed to serialize keyring: %w", err)
	}

	// Create the AES-GCM
	gcm, err := b.aeadFromKey(keyring.RootKey())
	if err != nil {
		return err
	}

	// Encrypt the barrier init value
	value, err := b.encrypt(keyringPath, initialKeyTerm, gcm, keyringBuf)
	if err != nil {
		return err
	}

	// Create the keyring physical entry
	pe := &physical.Entry{
		Key:   keyringPath,
		Value: value,
	}
	if err := b.backend.Put(ctx, pe); err != nil {
		return fmt.Errorf("failed to persist keyring: %w", err)
	}

	// Serialize the root key value
	key := &Key{
		Term:    1,
		Version: 1,
		Value:   keyring.RootKey(),
	}
	keyBuf, err := key.Serialize()
	defer memzero(keyBuf)
	if err != nil {
		return fmt.Errorf("failed to serialize root key: %w", err)
	}

	// Encrypt the root key
	activeKey := keyring.ActiveKey()
	aead, err := b.aeadFromKey(activeKey.Value)
	if err != nil {
		return err
	}
	value, err = b.encryptTracked(rootKeyPath, activeKey.Term, aead, keyBuf)
	if err != nil {
		return err
	}

	// Update the rootKeyPath for standby instances
	pe = &physical.Entry{
		Key:   rootKeyPath,
		Value: value,
	}
	if err := b.backend.Put(ctx, pe); err != nil {
		return fmt.Errorf("failed to persist root key: %w", err)
	}
	return nil
}

// GenerateKey is used to generate a new key
func (b *AESGCMBarrier) GenerateKey(reader io.Reader) ([]byte, error) {
	// Generate a 256bit key
	buf := make([]byte, 2*aes.BlockSize)
	_, err := reader.Read(buf)

	return buf, err
}

// KeyLength is used to sanity check a key
func (b *AESGCMBarrier) KeyLength() (int, int) {
	return aes.BlockSize, 2 * aes.BlockSize
}

// Sealed checks if the barrier has been unlocked yet. The Barrier
// is not expected to be able to perform any CRUD until it is unsealed.
func (b *AESGCMBarrier) Sealed() (bool, error) {
	b.l.RLock()
	sealed := b.sealed
	b.l.RUnlock()
	return sealed, nil
}

// VerifyRoot is used to check if the given key matches the root key
func (b *AESGCMBarrier) VerifyRoot(key []byte) error {
	b.l.RLock()
	defer b.l.RUnlock()
	if b.sealed {
		return ErrBarrierSealed
	}
	if subtle.ConstantTimeCompare(key, b.keyring.RootKey()) != 1 {
		return ErrBarrierInvalidKey
	}
	return nil
}

// ReloadKeyring is used to re-read the underlying keyring.
// This is used for HA deployments to ensure the latest keyring
// is present in the leader.
func (b *AESGCMBarrier) ReloadKeyring(ctx context.Context) error {
	b.l.Lock()
	defer b.l.Unlock()

	// Create the AES-GCM
	gcm, err := b.aeadFromKey(b.keyring.RootKey())
	if err != nil {
		return err
	}

	// Read in the keyring
	out, err := b.backend.Get(ctx, keyringPath)
	if err != nil {
		return fmt.Errorf("failed to check for keyring: %w", err)
	}

	// Ensure that the keyring exists. This should never happen,
	// and indicates something really bad has happened.
	if out == nil {
		return errors.New("keyring unexpectedly missing")
	}

	// Verify the term is always just one
	term := binary.BigEndian.Uint32(out.Value[:4])
	if term != initialKeyTerm {
		return errors.New("term mis-match")
	}

	// Decrypt the barrier init key
	plain, err := b.decrypt(keyringPath, gcm, out.Value)
	defer memzero(plain)
	if err != nil {
		if strings.Contains(err.Error(), "message authentication failed") {
			return ErrBarrierInvalidKey
		}
		return err
	}

	// Reset enc. counters, this may be a leadership change
	b.totalLocalEncryptions.Store(0)
	b.totalLocalEncryptions.Store(0)
	b.UnaccountedEncryptions.Store(0)
	b.RemoteEncryptions.Store(0)

	return b.recoverKeyring(plain)
}

func (b *AESGCMBarrier) recoverKeyring(plaintext []byte) error {
	keyring, err := DeserializeKeyring(plaintext)
	if err != nil {
		return fmt.Errorf("keyring deserialization failed: %w", err)
	}

	// Setup the keyring and finish
	b.cache = make(map[uint32]cipher.AEAD)
	b.keyring = keyring
	return nil
}

// ReloadRootKey is used to re-read the underlying root key.
// This is used for HA deployments to ensure the latest root key
// is available for keyring reloading.
func (b *AESGCMBarrier) ReloadRootKey(ctx context.Context) error {
	// Read the rootKeyPath upgrade
	out, err := b.Get(ctx, rootKeyPath)
	if err != nil {
		return fmt.Errorf("failed to read root key path: %w", err)
	}

	// The rootKeyPath could be missing (backwards incompatible),
	// we can ignore this and attempt to make progress with the current
	// root key.
	if out == nil {
		return nil
	}

	// Grab write lock and refetch
	b.l.Lock()
	defer b.l.Unlock()

	out, err = b.lockSwitchedGet(ctx, rootKeyPath, false)
	if err != nil {
		return fmt.Errorf("failed to read root key path: %w", err)
	}

	if out == nil {
		return nil
	}

	// Deserialize the root key
	key, err := DeserializeKey(out.Value)
	memzero(out.Value)
	if err != nil {
		return fmt.Errorf("failed to deserialize key: %w", err)
	}

	// Check if the root key is the same
	if subtle.ConstantTimeCompare(b.keyring.RootKey(), key.Value) == 1 {
		return nil
	}

	// Update the root key
	oldKeyring := b.keyring
	b.keyring = b.keyring.SetRootKey(key.Value)
	oldKeyring.Zeroize(false)
	return nil
}

// Unseal is used to provide the root key which permits the barrier
// to be unsealed. If the key is not correct, the barrier remains sealed.
func (b *AESGCMBarrier) Unseal(ctx context.Context, key []byte) error {
	b.l.Lock()
	defer b.l.Unlock()

	// Do nothing if already unsealed
	if !b.sealed {
		return nil
	}

	// Create the AES-GCM
	gcm, err := b.aeadFromKey(key)
	if err != nil {
		return err
	}

	// Read in the keyring
	out, err := b.backend.Get(ctx, keyringPath)
	if err != nil {
		return fmt.Errorf("failed to check for keyring: %w", err)
	}
	if out != nil {
		// Verify the term is always just one
		term := binary.BigEndian.Uint32(out.Value[:4])
		if term != initialKeyTerm {
			return errors.New("term mis-match")
		}

		// Decrypt the barrier init key
		plain, err := b.decrypt(keyringPath, gcm, out.Value)
		defer memzero(plain)
		if err != nil {
			if strings.Contains(err.Error(), "message authentication failed") {
				return ErrBarrierInvalidKey
			}
			return err
		}

		// Recover the keyring
		err = b.recoverKeyring(plain)
		if err != nil {
			return fmt.Errorf("keyring deserialization failed: %w", err)
		}

		b.sealed = false

		return nil
	}

	// Read the barrier initialization key
	out, err = b.backend.Get(ctx, barrierInitPath)
	if err != nil {
		return fmt.Errorf("failed to check for initialization: %w", err)
	}
	if out == nil {
		return ErrBarrierNotInit
	}

	// Verify the term is always just one
	term := binary.BigEndian.Uint32(out.Value[:4])
	if term != initialKeyTerm {
		return errors.New("term mis-match")
	}

	// Decrypt the barrier init key
	plain, err := b.decrypt(barrierInitPath, gcm, out.Value)
	if err != nil {
		if strings.Contains(err.Error(), "message authentication failed") {
			return ErrBarrierInvalidKey
		}
		return err
	}
	defer memzero(plain)

	// Unmarshal the barrier init
	var init barrierInit
	if err := jsonutil.DecodeJSON(plain, &init); err != nil {
		return fmt.Errorf("failed to unmarshal barrier init file")
	}

	// Setup a new keyring, this is for backwards compatibility
	keyringNew := NewKeyring()
	keyring := keyringNew.SetRootKey(key)

	// AddKey reuses the root, so we are only zeroizing after this call
	defer keyringNew.Zeroize(false)

	keyring, err = keyring.AddKey(&Key{
		Term:    1,
		Version: 1,
		Value:   init.Key,
	})
	if err != nil {
		return fmt.Errorf("failed to create keyring: %w", err)
	}
	if err := b.persistKeyring(ctx, keyring); err != nil {
		return err
	}

	// Delete the old barrier entry
	if err := b.backend.Delete(ctx, barrierInitPath); err != nil {
		return fmt.Errorf("failed to delete barrier init file: %w", err)
	}

	// Set the vault as unsealed
	b.keyring = keyring
	b.sealed = false

	return nil
}

// Seal is used to re-seal the barrier. This requires the barrier to
// be unsealed again to perform any further operations.
func (b *AESGCMBarrier) Seal() error {
	b.l.Lock()
	defer b.l.Unlock()

	// Remove the primary key, and seal the vault
	b.cache = make(map[uint32]cipher.AEAD)
	b.keyring.Zeroize(true)
	b.keyring = nil
	b.sealed = true
	return nil
}

// Rotate is used to create a new encryption key. All future writes
// should use the new key, while old values should still be decryptable.
func (b *AESGCMBarrier) Rotate(ctx context.Context, randomSource io.Reader) (uint32, error) {
	b.l.Lock()
	defer b.l.Unlock()
	if b.sealed {
		return 0, ErrBarrierSealed
	}

	// Generate a new key
	encrypt, err := b.GenerateKey(randomSource)
	if err != nil {
		return 0, fmt.Errorf("failed to generate encryption key: %w", err)
	}

	// Get the next term
	term := b.keyring.ActiveTerm()
	newTerm := term + 1

	// Add a new encryption key
	newKeyring, err := b.keyring.AddKey(&Key{
		Term:    newTerm,
		Version: 1,
		Value:   encrypt,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to add new encryption key: %w", err)
	}

	// Persist the new keyring
	if err := b.persistKeyring(ctx, newKeyring); err != nil {
		return 0, err
	}

	// Clear encryption tracking
	b.RemoteEncryptions.Store(0)
	b.totalLocalEncryptions.Store(0)
	b.UnaccountedEncryptions.Store(0)

	// Swap the keyrings
	b.keyring = newKeyring

	return newTerm, nil
}

// CreateUpgrade creates an upgrade path key to the given term from the previous term
func (b *AESGCMBarrier) CreateUpgrade(ctx context.Context, term uint32) error {
	b.l.RLock()
	if b.sealed {
		b.l.RUnlock()
		return ErrBarrierSealed
	}

	// Get the key for this term
	termKey := b.keyring.TermKey(term)
	buf, err := termKey.Serialize()
	defer memzero(buf)
	if err != nil {
		b.l.RUnlock()
		return err
	}

	// Get the AEAD for the previous term
	prevTerm := term - 1
	primary, err := b.aeadForTerm(prevTerm)
	if err != nil {
		b.l.RUnlock()
		return err
	}

	key := fmt.Sprintf("%s%d", keyringUpgradePrefix, prevTerm)
	value, err := b.encryptTracked(key, prevTerm, primary, buf)
	b.l.RUnlock()
	if err != nil {
		return err
	}
	// Create upgrade key
	pe := &physical.Entry{
		Key:   key,
		Value: value,
	}
	return b.backend.Put(ctx, pe)
}

// DestroyUpgrade destroys the upgrade path key to the given term
func (b *AESGCMBarrier) DestroyUpgrade(ctx context.Context, term uint32) error {
	path := fmt.Sprintf("%s%d", keyringUpgradePrefix, term-1)
	return b.Delete(ctx, path)
}

// CheckUpgrade looks for an upgrade to the current term and installs it
func (b *AESGCMBarrier) CheckUpgrade(ctx context.Context) (bool, uint32, error) {
	b.l.RLock()
	if b.sealed {
		b.l.RUnlock()
		return false, 0, ErrBarrierSealed
	}

	// Get the current term
	activeTerm := b.keyring.ActiveTerm()

	// Check for an upgrade key
	upgrade := fmt.Sprintf("%s%d", keyringUpgradePrefix, activeTerm)
	entry, err := b.lockSwitchedGet(ctx, upgrade, false)
	if err != nil {
		b.l.RUnlock()
		return false, 0, err
	}

	// Nothing to do if no upgrade
	if entry == nil {
		b.l.RUnlock()
		return false, 0, nil
	}

	// Upgrade from read lock to write lock
	b.l.RUnlock()
	b.l.Lock()
	defer b.l.Unlock()

	// Validate base cases and refetch values again

	if b.sealed {
		return false, 0, ErrBarrierSealed
	}

	activeTerm = b.keyring.ActiveTerm()

	upgrade = fmt.Sprintf("%s%d", keyringUpgradePrefix, activeTerm)
	entry, err = b.lockSwitchedGet(ctx, upgrade, false)
	if err != nil {
		return false, 0, err
	}

	if entry == nil {
		return false, 0, nil
	}

	// Deserialize the key
	key, err := DeserializeKey(entry.Value)
	memzero(entry.Value)
	if err != nil {
		return false, 0, err
	}

	// Update the keyring
	newKeyring, err := b.keyring.AddKey(key)
	if err != nil {
		return false, 0, fmt.Errorf("failed to add new encryption key: %w", err)
	}
	b.keyring = newKeyring

	// Done!
	return true, key.Term, nil
}

// ActiveKeyInfo is used to inform details about the active key
func (b *AESGCMBarrier) ActiveKeyInfo() (*KeyInfo, error) {
	b.l.RLock()
	defer b.l.RUnlock()
	if b.sealed {
		return nil, ErrBarrierSealed
	}

	// Determine the key install time
	term := b.keyring.ActiveTerm()
	key := b.keyring.TermKey(term)

	// Return the key info
	info := &KeyInfo{
		Term:        int(term),
		InstallTime: key.InstallTime,
		Encryptions: b.encryptions(),
	}
	return info, nil
}

// Rekey is used to change the root key used to protect the keyring
func (b *AESGCMBarrier) Rekey(ctx context.Context, key []byte) error {
	b.l.Lock()
	defer b.l.Unlock()

	newKeyring, err := b.updateRootKeyCommon(key)
	if err != nil {
		return err
	}

	// Persist the new keyring
	if err := b.persistKeyring(ctx, newKeyring); err != nil {
		return err
	}

	// Swap the keyrings
	oldKeyring := b.keyring
	b.keyring = newKeyring
	oldKeyring.Zeroize(false)
	return nil
}

// SetRootKey updates the keyring's in-memory root key but does not persist
// anything to storage
func (b *AESGCMBarrier) SetRootKey(key []byte) error {
	b.l.Lock()
	defer b.l.Unlock()

	newKeyring, err := b.updateRootKeyCommon(key)
	if err != nil {
		return err
	}

	// Swap the keyrings
	oldKeyring := b.keyring
	b.keyring = newKeyring
	oldKeyring.Zeroize(false)
	return nil
}

// Performs common tasks related to updating the root key; note that the lock
// must be held before calling this function
func (b *AESGCMBarrier) updateRootKeyCommon(key []byte) (*Keyring, error) {
	if b.sealed {
		return nil, ErrBarrierSealed
	}

	// Verify the key size
	min, max := b.KeyLength()
	if len(key) < min || len(key) > max {
		return nil, fmt.Errorf("key size must be %d or %d", min, max)
	}

	return b.keyring.SetRootKey(key), nil
}

// Put is used to insert or update an entry
func (b *AESGCMBarrier) Put(ctx context.Context, entry *logical.StorageEntry) error {
	defer metrics.MeasureSince([]string{"barrier", "put"}, time.Now())
	b.l.RLock()
	if b.sealed {
		b.l.RUnlock()
		return ErrBarrierSealed
	}

	term := b.keyring.ActiveTerm()
	primary, err := b.aeadForTerm(term)
	b.l.RUnlock()
	if err != nil {
		return err
	}

	return b.putInternal(ctx, term, primary, entry)
}

func (b *AESGCMBarrier) putInternal(ctx context.Context, term uint32, primary cipher.AEAD, entry *logical.StorageEntry) error {
	value, err := b.encryptTracked(entry.Key, term, primary, entry.Value)
	if err != nil {
		return err
	}
	pe := &physical.Entry{
		Key:      entry.Key,
		Value:    value,
		SealWrap: entry.SealWrap,
	}
	return b.backend.Put(ctx, pe)
}

// Get is used to fetch an entry
func (b *AESGCMBarrier) Get(ctx context.Context, key string) (*logical.StorageEntry, error) {
	return b.lockSwitchedGet(ctx, key, true)
}

func (b *AESGCMBarrier) lockSwitchedGet(ctx context.Context, key string, getLock bool) (*logical.StorageEntry, error) {
	defer metrics.MeasureSince([]string{"barrier", "get"}, time.Now())
	if getLock {
		b.l.RLock()
	}
	if b.sealed {
		if getLock {
			b.l.RUnlock()
		}
		return nil, ErrBarrierSealed
	}

	// Read the key from the backend
	pe, err := b.backend.Get(ctx, key)
	if err != nil {
		if getLock {
			b.l.RUnlock()
		}
		return nil, err
	} else if pe == nil {
		if getLock {
			b.l.RUnlock()
		}
		return nil, nil
	}

	if len(pe.Value) < 4 {
		if getLock {
			b.l.RUnlock()
		}
		return nil, errors.New("invalid value")
	}

	// Verify the term
	term := binary.BigEndian.Uint32(pe.Value[:4])

	// Get the GCM by term
	// It is expensive to do this first but it is not a
	// normal case that this won't match
	gcm, err := b.aeadForTerm(term)
	if getLock {
		b.l.RUnlock()
	}
	if err != nil {
		return nil, err
	}
	if gcm == nil {
		return nil, fmt.Errorf("no decryption key available for term %d", term)
	}

	// Decrypt the ciphertext
	plain, err := b.decrypt(key, gcm, pe.Value)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %w", err)
	}

	// Wrap in a logical entry
	entry := &logical.StorageEntry{
		Key:      key,
		Value:    plain,
		SealWrap: pe.SealWrap,
	}
	return entry, nil
}

// Delete is used to permanently delete an entry
func (b *AESGCMBarrier) Delete(ctx context.Context, key string) error {
	defer metrics.MeasureSince([]string{"barrier", "delete"}, time.Now())
	b.l.RLock()
	sealed := b.sealed
	b.l.RUnlock()
	if sealed {
		return ErrBarrierSealed
	}

	return b.backend.Delete(ctx, key)
}

// List is used ot list all the keys under a given
// prefix, up to the next prefix.
func (b *AESGCMBarrier) List(ctx context.Context, prefix string) ([]string, error) {
	defer metrics.MeasureSince([]string{"barrier", "list"}, time.Now())
	b.l.RLock()
	sealed := b.sealed
	b.l.RUnlock()
	if sealed {
		return nil, ErrBarrierSealed
	}

	return b.backend.List(ctx, prefix)
}

// aeadForTerm returns the AES-GCM AEAD for the given term
func (b *AESGCMBarrier) aeadForTerm(term uint32) (cipher.AEAD, error) {
	// Check for the keyring
	keyring := b.keyring
	if keyring == nil {
		return nil, nil
	}

	// Check the cache for the aead
	b.cacheLock.RLock()
	aead, ok := b.cache[term]
	b.cacheLock.RUnlock()
	if ok {
		return aead, nil
	}

	// Read the underlying key
	key := keyring.TermKey(term)
	if key == nil {
		return nil, nil
	}

	// Create a new aead
	aead, err := b.aeadFromKey(key.Value)
	if err != nil {
		return nil, err
	}

	// Update the cache
	b.cacheLock.Lock()
	b.cache[term] = aead
	b.cacheLock.Unlock()
	return aead, nil
}

// aeadFromKey returns an AES-GCM AEAD using the given key.
func (b *AESGCMBarrier) aeadFromKey(key []byte) (cipher.AEAD, error) {
	// Create the AES cipher
	aesCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create the GCM mode AEAD
	gcm, err := cipher.NewGCM(aesCipher)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize GCM mode")
	}
	return gcm, nil
}

// encrypt is used to encrypt a value
func (b *AESGCMBarrier) encrypt(path string, term uint32, gcm cipher.AEAD, plain []byte) ([]byte, error) {
	// Allocate the output buffer with room for tern, version byte,
	// nonce, GCM tag and the plaintext

	extra := termSize + 1 + gcm.NonceSize() + gcm.Overhead()
	if len(plain) > math.MaxInt-extra {
		return nil, ErrPlaintextTooLarge
	}

	capacity := len(plain) + extra
	size := termSize + 1 + gcm.NonceSize()
	out := make([]byte, size, capacity)

	// Set the key term
	binary.BigEndian.PutUint32(out[:4], term)

	// Set the version byte
	out[4] = b.currentAESGCMVersionByte

	// Generate a random nonce
	nonce := out[5 : 5+gcm.NonceSize()]
	n, err := rand.Read(nonce)
	if err != nil {
		return nil, err
	}
	if n != len(nonce) {
		return nil, errors.New("unable to read enough random bytes to fill gcm nonce")
	}

	// Seal the output
	switch b.currentAESGCMVersionByte {
	case AESGCMVersion1:
		out = gcm.Seal(out, nonce, plain, nil)
	case AESGCMVersion2:
		aad := []byte(nil)
		if path != "" {
			aad = []byte(path)
		}
		out = gcm.Seal(out, nonce, plain, aad)
	default:
		panic("Unknown AESGCM version")
	}

	return out, nil
}

func termLabel(term uint32) []metrics.Label {
	return []metrics.Label{
		{
			Name:  "term",
			Value: strconv.FormatUint(uint64(term), 10),
		},
	}
}

// decrypt is used to decrypt a value using the keyring
func (b *AESGCMBarrier) decrypt(path string, gcm cipher.AEAD, cipher []byte) ([]byte, error) {
	if len(cipher) < 5+gcm.NonceSize() {
		return nil, fmt.Errorf("invalid cipher length")
	}
	// Capture the parts
	nonce := cipher[5 : 5+gcm.NonceSize()]
	raw := cipher[5+gcm.NonceSize():]
	out := make([]byte, 0, len(raw)-gcm.NonceSize())

	// Attempt to open
	switch cipher[4] {
	case AESGCMVersion1:
		return gcm.Open(out, nonce, raw, nil)
	case AESGCMVersion2:
		aad := []byte(nil)
		if path != "" {
			aad = []byte(path)
		}
		return gcm.Open(out, nonce, raw, aad)
	default:
		return nil, fmt.Errorf("version bytes mis-match")
	}
}

// Encrypt is used to encrypt in-memory for the BarrierEncryptor interface
func (b *AESGCMBarrier) Encrypt(ctx context.Context, key string, plaintext []byte) ([]byte, error) {
	b.l.RLock()
	if b.sealed {
		b.l.RUnlock()
		return nil, ErrBarrierSealed
	}

	term := b.keyring.ActiveTerm()
	primary, err := b.aeadForTerm(term)
	b.l.RUnlock()
	if err != nil {
		return nil, err
	}

	ciphertext, err := b.encryptTracked(key, term, primary, plaintext)
	if err != nil {
		return nil, err
	}

	return ciphertext, nil
}

// Decrypt is used to decrypt in-memory for the BarrierEncryptor interface
func (b *AESGCMBarrier) Decrypt(_ context.Context, key string, ciphertext []byte) ([]byte, error) {
	b.l.RLock()
	if b.sealed {
		b.l.RUnlock()
		return nil, ErrBarrierSealed
	}

	if len(ciphertext) == 0 {
		b.l.RUnlock()
		return nil, fmt.Errorf("empty ciphertext")
	}

	// Verify the term
	if len(ciphertext) < 4 {
		b.l.RUnlock()
		return nil, fmt.Errorf("invalid ciphertext term")
	}
	term := binary.BigEndian.Uint32(ciphertext[:4])

	// Get the GCM by term
	// It is expensive to do this first but it is not a
	// normal case that this won't match
	gcm, err := b.aeadForTerm(term)
	b.l.RUnlock()
	if err != nil {
		return nil, err
	}
	if gcm == nil {
		return nil, fmt.Errorf("no decryption key available for term %d", term)
	}

	// Decrypt the ciphertext
	plain, err := b.decrypt(key, gcm, ciphertext)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %w", err)
	}

	return plain, nil
}

func (b *AESGCMBarrier) Keyring() (*Keyring, error) {
	b.l.RLock()
	defer b.l.RUnlock()
	if b.sealed {
		return nil, ErrBarrierSealed
	}

	return b.keyring.Clone(), nil
}

func (b *AESGCMBarrier) ConsumeEncryptionCount(consumer func(int64) error) error {
	if b.keyring != nil {
		// Lock to prevent replacement of the key while we consume the encryptions
		b.l.RLock()
		defer b.l.RUnlock()

		c := b.UnaccountedEncryptions.Load()
		err := consumer(c)
		if err == nil && c > 0 {
			// Consumer succeeded, remove those from local encryptions
			b.UnaccountedEncryptions.Sub(c)
		}
		return err
	}
	return nil
}

func (b *AESGCMBarrier) AddRemoteEncryptions(encryptions int64) {
	// For rollup and persistence
	b.UnaccountedEncryptions.Add(encryptions)
	// For testing
	b.RemoteEncryptions.Add(encryptions)
}

func (b *AESGCMBarrier) encryptTracked(path string, term uint32, gcm cipher.AEAD, buf []byte) ([]byte, error) {
	ct, err := b.encrypt(path, term, gcm, buf)
	if err != nil {
		return nil, err
	}
	// Increment the local encryption count, and track metrics
	b.UnaccountedEncryptions.Add(1)
	b.totalLocalEncryptions.Add(1)
	metrics.IncrCounterWithLabels(barrierEncryptsMetric, 1, termLabel(term))

	return ct, nil
}

// UnaccountedEncryptions returns the number of encryptions made on the local instance only for the current key term
func (b *AESGCMBarrier) TotalLocalEncryptions() int64 {
	return b.totalLocalEncryptions.Load()
}

func (b *AESGCMBarrier) CheckBarrierAutoRotate(ctx context.Context) (string, error) {
	const oneYear = 24 * 365 * time.Hour
	reason, err := func() (string, error) {
		b.l.RLock()
		defer b.l.RUnlock()
		if b.keyring != nil {
			// Rotation Checks
			var reason string

			rc, err := b.RotationConfig()
			if err != nil {
				return "", err
			}

			if !rc.Disabled {
				activeKey := b.keyring.ActiveKey()
				ops := b.encryptions()
				switch {
				case activeKey.Encryptions == 0 && !activeKey.InstallTime.IsZero() && time.Since(activeKey.InstallTime) > oneYear:
					reason = legacyRotateReason
				case ops > rc.MaxOperations:
					reason = "reached max operations"
				case rc.Interval > 0 && time.Since(activeKey.InstallTime) > rc.Interval:
					reason = "rotation interval reached"
				}
			}
			return reason, nil
		}
		return "", nil
	}()
	if err != nil {
		return "", err
	}
	if reason != "" {
		return reason, nil
	}

	b.l.Lock()
	defer b.l.Unlock()
	if b.keyring != nil {
		err := b.persistEncryptions(ctx)
		if err != nil {
			return "", err
		}
	}
	return reason, nil
}

// Must be called with lock held
func (b *AESGCMBarrier) persistEncryptions(ctx context.Context) error {
	if !b.sealed {
		// Encryption count persistence
		upe := b.UnaccountedEncryptions.Load()
		if upe > 0 {
			activeKey := b.keyring.ActiveKey()
			// Move local (unpersisted) encryptions to the key and persist.  This prevents us from needing to persist if
			// there has been no activity. Since persistence performs an encryption, perversely we zero out after
			// persistence and add 1 to the count to avoid this operation guaranteeing we need another
			// autoRotateCheckInterval later.
			newEncs := upe + 1
			activeKey.Encryptions += uint64(newEncs)
			newKeyring := b.keyring.Clone()
			err := b.persistKeyring(ctx, newKeyring)
			if err != nil {
				return err
			}
			b.UnaccountedEncryptions.Sub(newEncs)
		}
	}
	return nil
}

// Mostly for testing, returns the total number of encryption operations performed on the active term
func (b *AESGCMBarrier) encryptions() int64 {
	if b.keyring != nil {
		activeKey := b.keyring.ActiveKey()
		if activeKey != nil {
			return b.UnaccountedEncryptions.Load() + int64(activeKey.Encryptions)
		}
	}
	return 0
}
