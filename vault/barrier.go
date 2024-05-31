// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/hashicorp/vault/sdk/logical"
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

	// ErrPlaintextTooLarge is returned if a plaintext is offered for encryption
	// that is too large to encrypt in memory
	ErrPlaintextTooLarge = errors.New("plaintext value too large")
)

const (
	// barrierInitPath is the path used to store our init sentinel file
	barrierInitPath = "barrier/init"

	// keyringPath is the location of the keyring data. This is encrypted
	// by the root key.
	keyringPath   = "core/keyring"
	keyringPrefix = "core/"

	// keyringUpgradePrefix is the path used to store keyring update entries.
	// When running in HA mode, the active instance will install the new key
	// and re-write the keyring. For standby instances, they need an upgrade
	// path from key N to N+1. They cannot just use the root key because
	// in the event of a rekey, that root key can no longer decrypt the keyring.
	// When key N+1 is installed, we create an entry at "prefix/N" which uses
	// encryption key N to provide the N+1 key. The standby instances scan
	// for this periodically and refresh their keyring. The upgrade keys
	// are deleted after a few minutes, but this provides enough time for the
	// standby instances to upgrade without causing any disruption.
	keyringUpgradePrefix = "core/upgrade/"

	// rootKeyPath is the location of the root key. This is encrypted
	// by the latest key in the keyring. This is only used by standby instances
	// to handle the case of a rekey. If the active instance does a rekey,
	// the standby instances can no longer reload the keyring since they
	// have the old root key. This key can be decrypted if you have the
	// keyring to discover the new root key. The new root key is then
	// used to reload the keyring itself.
	rootKeyPath = "core/master"

	// shamirKekPath is used with Shamir in v1.3+ to store a copy of the
	// unseal key behind the barrier.  As with rootKeyPath this is primarily
	// used by standbys to handle rekeys.  It also comes into play when restoring
	// raft snapshots.
	shamirKekPath = "core/shamir-kek"
)

// SecurityBarrier is a critical component of Vault. It is used to wrap
// an untrusted physical backend and provide a single point of encryption,
// decryption and checksum verification. The goal is to ensure that any
// data written to the barrier is confidential and that integrity is preserved.
// As a real-world analogy, this is the steel and concrete wrapper around
// a Vault. The barrier should only be Unlockable given its key.
type SecurityBarrier interface {
	// Initialized checks if the barrier has been initialized
	// and has a root key set.
	Initialized(ctx context.Context) (bool, error)

	// Initialize works only if the barrier has not been initialized
	// and makes use of the given root key.  When sealKey is provided
	// it's because we're using a new-style Shamir seal, and rootKey
	// is to be stored using sealKey to encrypt it.
	Initialize(ctx context.Context, rootKey []byte, sealKey []byte, random io.Reader) error

	// GenerateKey is used to generate a new key
	GenerateKey(io.Reader) ([]byte, error)

	// KeyLength is used to sanity check a key
	KeyLength() (int, int)

	// Sealed checks if the barrier has been unlocked yet. The Barrier
	// is not expected to be able to perform any CRUD until it is unsealed.
	Sealed() (bool, error)

	// Unseal is used to provide the unseal key which permits the barrier
	// to be unsealed. If the key is not correct, the barrier remains sealed.
	Unseal(ctx context.Context, key []byte) error

	// VerifyRoot is used to check if the given key matches the root key
	VerifyRoot(key []byte) error

	// SetRootKey is used to directly set a new root key. This is used in
	// replicated scenarios due to the chicken and egg problem of reloading the
	// keyring from disk before we have the root key to decrypt it.
	SetRootKey(key []byte) error

	// ReloadKeyring is used to re-read the underlying keyring.
	// This is used for HA deployments to ensure the latest keyring
	// is present in the leader.
	ReloadKeyring(ctx context.Context) error

	// ReloadRootKey is used to re-read the underlying root key.
	// This is used for HA deployments to ensure the latest root key
	// is available for keyring reloading.
	ReloadRootKey(ctx context.Context) error

	// Seal is used to re-seal the barrier. This requires the barrier to
	// be unsealed again to perform any further operations.
	Seal() error

	// Rotate is used to create a new encryption key. All future writes
	// should use the new key, while old values should still be decryptable.
	Rotate(ctx context.Context, reader io.Reader) (uint32, error)

	// CreateUpgrade creates an upgrade path key to the given term from the previous term
	CreateUpgrade(ctx context.Context, term uint32) error

	// DestroyUpgrade destroys the upgrade path key to the given term
	DestroyUpgrade(ctx context.Context, term uint32) error

	// CheckUpgrade looks for an upgrade to the current term and installs it
	CheckUpgrade(ctx context.Context) (bool, uint32, error)

	// ActiveKeyInfo is used to inform details about the active key
	ActiveKeyInfo() (*KeyInfo, error)

	// RotationConfig returns the auto-rotation config for the barrier key
	RotationConfig() (KeyRotationConfig, error)

	// SetRotationConfig updates the auto-rotation config for the barrier key
	SetRotationConfig(ctx context.Context, config KeyRotationConfig) error

	// Rekey is used to change the master key used to protect the keyring
	Rekey(context.Context, []byte) error

	// For replication we must send over the keyring, so this must be available
	Keyring() (*Keyring, error)

	// For encryption count shipping, a function which handles updating local encryption counts if the consumer succeeds.
	// This isolates the barrier code from the replication system
	ConsumeEncryptionCount(consumer func(int64) error) error

	// Add encryption counts from a remote source (downstream cluster node)
	AddRemoteEncryptions(encryptions int64)

	// Check whether an automatic rotation is due
	CheckBarrierAutoRotate(ctx context.Context) (string, error)

	// SecurityBarrier must provide the storage APIs
	logical.Storage

	// SecurityBarrier must provide the encryption APIs
	BarrierEncryptor

	DetectDeadlocks() bool
}

// BarrierStorage is the storage only interface required for a Barrier.
type BarrierStorage interface {
	// Put is used to insert or update an entry
	Put(ctx context.Context, entry *logical.StorageEntry) error

	// Get is used to fetch an entry
	Get(ctx context.Context, key string) (*logical.StorageEntry, error)

	// Delete is used to permanently delete an entry
	Delete(ctx context.Context, key string) error

	// List is used ot list all the keys under a given
	// prefix, up to the next prefix.
	List(ctx context.Context, prefix string) ([]string, error)
}

// BarrierEncryptor is the in memory only interface that does not actually
// use the underlying barrier. It is used for lower level modules like the
// Write-Ahead-Log and Merkle index to allow them to use the barrier.
type BarrierEncryptor interface {
	Encrypt(ctx context.Context, key string, plaintext []byte) ([]byte, error)
	Decrypt(ctx context.Context, key string, ciphertext []byte) ([]byte, error)
}

// KeyInfo is used to convey information about the encryption key
type KeyInfo struct {
	Term        int
	InstallTime time.Time
	Encryptions int64
}
