package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/vault/physical"
)

const (
	// aesgcmVersionByte is prefixed to a message to allow for
	// future versioning of barrier implementations.
	aesgcmVersionByte = 0x1
)

// barrierInit is the JSON encoded value stored
type barrierInit struct {
	Version int    // Version is the current format version
	Key     []byte // Key is the primary encryption key
}

// AESGCMBarrier is a SecurityBarrier implementation that
// uses a 128bit AES encryption cipher with the Galois Counter Mode.
// AES-GCM is high performance, and provides both confidentiality
// and integrity.
type AESGCMBarrier struct {
	backend physical.Backend

	l      sync.RWMutex
	sealed bool

	// primary is the AEAD keyed from the encryption key.
	// This is the cipher that should be used to encrypt/decrypt
	// all the underlying values. It will be available if the
	// barrier is unsealed.
	primary cipher.AEAD
}

// NewAESGCMBarrier is used to construct a new barrier that uses
// the provided physical backend for storage.
func NewAESGCMBarrier(physical physical.Backend) (*AESGCMBarrier, error) {
	b := &AESGCMBarrier{
		backend: physical,
		sealed:  true,
	}
	return b, nil
}

// Initialized checks if the barrier has been initialized
// and has a master key set.
func (b *AESGCMBarrier) Initialized() (bool, error) {
	// Read the init sentinel file
	out, err := b.backend.Get(barrierInitPath)
	if err != nil {
		return false, fmt.Errorf("failed to check for initialization: %v", err)
	}
	return out != nil, nil
}

// Initialize works only if the barrier has not been initialized
// and makes use of the given master key.
func (b *AESGCMBarrier) Initialize(key []byte) error {
	// Verify the key size
	if len(key) != aes.BlockSize {
		return fmt.Errorf("Key size must be %d", aes.BlockSize)
	}

	// Check if already initialized
	if alreadyInit, err := b.Initialized(); err != nil {
		return err
	} else if alreadyInit {
		return ErrBarrierAlreadyInit
	}

	// Create the AES-GCM
	gcm, err := b.aeadFromKey(key)
	if err != nil {
		return err
	}

	// Generate encryption key
	encrypt, err := b.GenerateKey()
	if err != nil {
		return fmt.Errorf("failed to generate encryption key: %v", err)
	}
	defer memzero(encrypt)

	// Create the barrier init entry
	init := barrierInit{
		Version: 1,
		Key:     encrypt,
	}
	buf, err := json.Marshal(init)
	if err != nil {
		return fmt.Errorf("failed to create barrier entry: %v", err)
	}
	defer memzero(buf)

	// Encrypt the barrier init value
	value := b.encrypt(gcm, buf)

	// Create the barrierInitPath
	pe := &physical.Entry{
		Key:   barrierInitPath,
		Value: value,
	}
	if err := b.backend.Put(pe); err != nil {
		return fmt.Errorf("failed to create initialization key: %v", err)
	}
	return nil
}

// GenerateKey is used to generate a new key
func (b *AESGCMBarrier) GenerateKey() ([]byte, error) {
	buf := make([]byte, aes.BlockSize)
	_, err := rand.Read(buf)
	return buf, err
}

// KeyLength is used to sanity check a key
func (b *AESGCMBarrier) KeyLength() (int, int) {
	return aes.BlockSize, aes.BlockSize
}

// Sealed checks if the barrier has been unlocked yet. The Barrier
// is not expected to be able to perform any CRUD until it is unsealed.
func (b *AESGCMBarrier) Sealed() (bool, error) {
	b.l.RLock()
	defer b.l.RUnlock()
	return b.sealed, nil
}

// Unseal is used to provide the master key which permits the barrier
// to be unsealed. If the key is not correct, the barrier remains sealed.
func (b *AESGCMBarrier) Unseal(key []byte) error {
	b.l.Lock()
	defer b.l.Unlock()

	// Do nothing if already unsealed
	if !b.sealed {
		return nil
	}

	// Read the barrier initialization key
	out, err := b.backend.Get(barrierInitPath)
	if err != nil {
		return fmt.Errorf("failed to check for initialization: %v", err)
	}
	if out == nil {
		return ErrBarrierNotInit
	}

	// Create the AES-GCM
	gcm, err := b.aeadFromKey(key)
	if err != nil {
		return err
	}

	// Decrypt the barrier init key
	plain, err := b.decrypt(gcm, out.Value)
	if err != nil {
		if strings.Contains(err.Error(), "message authentication failed") {
			return ErrBarrierInvalidKey
		}
		return err
	}
	defer memzero(plain)

	// Unmarshal the barrier init
	var init barrierInit
	if err := json.Unmarshal(plain, &init); err != nil {
		return fmt.Errorf("failed to unmarshal barrier init file")
	}
	defer memzero(init.Key)

	// Initialize the master encryption key
	b.primary, err = b.aeadFromKey(init.Key)
	if err != nil {
		return err
	}

	// Set the vault as unsealed
	b.sealed = false
	return nil
}

// Seal is used to re-seal the barrier. This requires the barrier to
// be unsealed again to perform any further operations.
func (b *AESGCMBarrier) Seal() error {
	b.l.Lock()
	defer b.l.Unlock()

	// Remove the primary key, and seal the vault
	b.primary = nil
	b.sealed = true
	return nil
}

// Put is used to insert or update an entry
func (b *AESGCMBarrier) Put(entry *Entry) error {
	defer metrics.MeasureSince([]string{"barrier", "put"}, time.Now())
	b.l.RLock()
	defer b.l.RUnlock()

	primary := b.primary
	if primary == nil {
		return ErrBarrierSealed
	}

	pe := &physical.Entry{
		Key:   entry.Key,
		Value: b.encrypt(primary, entry.Value),
	}
	return b.backend.Put(pe)
}

// Get is used to fetch an entry
func (b *AESGCMBarrier) Get(key string) (*Entry, error) {
	defer metrics.MeasureSince([]string{"barrier", "get"}, time.Now())
	b.l.RLock()
	defer b.l.RUnlock()

	primary := b.primary
	if primary == nil {
		return nil, ErrBarrierSealed
	}

	// Read the key from the backend
	pe, err := b.backend.Get(key)
	if err != nil {
		return nil, err
	} else if pe == nil {
		return nil, nil
	}

	// Decrypt the ciphertext
	plain, err := b.decrypt(primary, pe.Value)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %v", err)
	}

	// Wrap in a logical entry
	entry := &Entry{
		Key:   key,
		Value: plain,
	}
	return entry, nil
}

// Delete is used to permanently delete an entry
func (b *AESGCMBarrier) Delete(key string) error {
	defer metrics.MeasureSince([]string{"barrier", "delete"}, time.Now())
	b.l.RLock()
	defer b.l.RUnlock()
	if b.sealed {
		return ErrBarrierSealed
	}

	return b.backend.Delete(key)
}

// List is used ot list all the keys under a given
// prefix, up to the next prefix.
func (b *AESGCMBarrier) List(prefix string) ([]string, error) {
	defer metrics.MeasureSince([]string{"barrier", "list"}, time.Now())
	b.l.RLock()
	defer b.l.RUnlock()
	if b.sealed {
		return nil, ErrBarrierSealed
	}

	return b.backend.List(prefix)
}

// aeadFromKey returns an AES-GCM AEAD using the given key.
func (b *AESGCMBarrier) aeadFromKey(key []byte) (cipher.AEAD, error) {
	// Create the AES cipher
	aesCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %v", err)
	}

	// Create the GCM mode AEAD
	gcm, err := cipher.NewGCM(aesCipher)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize GCM mode")
	}
	return gcm, nil
}

// encrypt is used to encrypt a value
func (b *AESGCMBarrier) encrypt(gcm cipher.AEAD, plain []byte) []byte {
	// Allocate the output buffer with room for version byte,
	// nonce, GCM tag and the plaintext
	capacity := 1 + gcm.NonceSize() + gcm.Overhead() + len(plain)
	size := 1 + gcm.NonceSize()
	out := make([]byte, size, capacity)

	// Set the version byte
	out[0] = aesgcmVersionByte

	// Generate a random nonce
	nonce := out[1 : 1+gcm.NonceSize()]
	rand.Read(nonce)

	// Seal the output
	out = gcm.Seal(out, nonce, plain, nil)
	return out
}

// decrypt is used to decrypt a value
func (b *AESGCMBarrier) decrypt(gcm cipher.AEAD, cipher []byte) ([]byte, error) {
	// Verify the version byte
	if cipher[0] != aesgcmVersionByte {
		return nil, fmt.Errorf("version bytes mis-match")
	}

	// Capture the parts
	nonce := cipher[1 : 1+gcm.NonceSize()]
	raw := cipher[1+gcm.NonceSize():]
	out := make([]byte, 0, len(raw)-gcm.NonceSize())

	// Attempt to open
	return gcm.Open(out, nonce, raw, nil)
}
