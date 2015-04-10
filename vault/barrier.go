package vault

import (
	"errors"

	"github.com/hashicorp/vault/logical"
)

var (
	// ErrBarrierSealed is returned if an operation is performed on
	// a sealed barrier. No operation is expected to succeed before unsealing
	ErrBarrierSealed = errors.New("Vault is sealed")

	// ErrBarrierAlreadyInit is returned if the barrier is already
	// initialized. This prevents a re-initialization.
	ErrBarrierAlreadyInit = errors.New("Vault is already initialized")

	// ErrBarrierNotInit is returned if a non-initialized barrier
	// is attempted to be unsealed.
	ErrBarrierNotInit = errors.New("Vault is not initialized")

	// ErrBarrierInvalidKey is returned if the Unseal key is invalid
	ErrBarrierInvalidKey = errors.New("Unseal failed, invalid key")
)

const (
	// barrierInitPath is the path used to store our init sentinel file
	barrierInitPath = "barrier/init"
)

// SecurityBarrier is a critical component of Vault. It is used to wrap
// an untrusted physical backend and provide a single point of encryption,
// decryption and checksum verification. The goal is to ensure that any
// data written to the barrier is confidential and that integrity is preserved.
// As a real-world analogy, this is the steel and concrete wrapper around
// a Vault. The barrier should only be Unlockable given its key.
type SecurityBarrier interface {
	// Initialized checks if the barrier has been initialized
	// and has a master key set.
	Initialized() (bool, error)

	// Initialize works only if the barrier has not been initialized
	// and makes use of the given master key.
	Initialize([]byte) error

	// GenerateKey is used to generate a new key
	GenerateKey() ([]byte, error)

	// KeyLength is used to sanity check a key
	KeyLength() (int, int)

	// Sealed checks if the barrier has been unlocked yet. The Barrier
	// is not expected to be able to perform any CRUD until it is unsealed.
	Sealed() (bool, error)

	// Unseal is used to provide the master key which permits the barrier
	// to be unsealed. If the key is not correct, the barrier remains sealed.
	Unseal(key []byte) error

	// Seal is used to re-seal the barrier. This requires the barrier to
	// be unsealed again to perform any further operations.
	Seal() error

	// SecurityBarrier must provide the storage APIs
	BarrierStorage
}

// BarrierStorage is the storage only interface required for a Barrier.
type BarrierStorage interface {
	// Put is used to insert or update an entry
	Put(entry *Entry) error

	// Get is used to fetch an entry
	Get(key string) (*Entry, error)

	// Delete is used to permanently delete an entry
	Delete(key string) error

	// List is used ot list all the keys under a given
	// prefix, up to the next prefix.
	List(prefix string) ([]string, error)
}

// Entry is used to represent data stored by the security barrier
type Entry struct {
	Key   string
	Value []byte
}

// Logical turns the Entry into a logical storage entry.
func (e *Entry) Logical() *logical.StorageEntry {
	return &logical.StorageEntry{
		Key:   e.Key,
		Value: e.Value,
	}
}
