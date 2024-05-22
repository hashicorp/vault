// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logical

import (
	"context"
	"crypto"
	"io"

	wrapping "github.com/hashicorp/go-kms-wrapping/v2"
)

//go:generate enumer -type=KeyUsage -trimprefix=KeyUsage -transform=snake
type KeyUsage int

const (
	KeyUsageEncrypt KeyUsage = 1 + iota
	KeyUsageDecrypt
	KeyUsageSign
	KeyUsageVerify
	KeyUsageWrap
	KeyUsageUnwrap
	KeyUsageGenerateRandom
)

type ManagedKey interface {
	// Name is a human-readable identifier for a managed key that may change/renamed. Use Uuid if a
	// long term consistent identifier is needed.
	Name() string
	// UUID is a unique identifier for a managed key that is guaranteed to remain
	// consistent even if a key is migrated or renamed.
	UUID() string
	// Present returns true if the key is established in the KMS.  This may return false if for example
	// an HSM library is not configured on all cluster nodes.
	Present(ctx context.Context) (bool, error)

	// AllowsAll returns true if all the requested usages are supported by the managed key.
	AllowsAll(usages []KeyUsage) bool
}

type (
	ManagedKeyConsumer             func(context.Context, ManagedKey) error
	ManagedSigningKeyConsumer      func(context.Context, ManagedSigningKey) error
	ManagedEncryptingKeyConsumer   func(context.Context, ManagedEncryptingKey) error
	ManagedMACKeyConsumer          func(context.Context, ManagedMACKey) error
	ManagedKeyRandomSourceConsumer func(context.Context, ManagedKeyRandomSource) error
)

type ManagedKeySystemView interface {
	// WithManagedKeyByName retrieves an instantiated managed key for consumption by the given function.  The
	// provided key can only be used within the scope of that function call
	WithManagedKeyByName(ctx context.Context, keyName, backendUUID string, f ManagedKeyConsumer) error
	// WithManagedKeyByUUID retrieves an instantiated managed key for consumption by the given function.  The
	// provided key can only be used within the scope of that function call
	WithManagedKeyByUUID(ctx context.Context, keyUuid, backendUUID string, f ManagedKeyConsumer) error

	// WithManagedSigningKeyByName retrieves an instantiated managed signing key for consumption by the given function,
	// with the same semantics as WithManagedKeyByName
	WithManagedSigningKeyByName(ctx context.Context, keyName, backendUUID string, f ManagedSigningKeyConsumer) error
	// WithManagedSigningKeyByUUID retrieves an instantiated managed signing key for consumption by the given function,
	// with the same semantics as WithManagedKeyByUUID
	WithManagedSigningKeyByUUID(ctx context.Context, keyUuid, backendUUID string, f ManagedSigningKeyConsumer) error
	// WithManagedSigningKeyByName retrieves an instantiated managed signing key for consumption by the given function,
	// with the same semantics as WithManagedKeyByName
	WithManagedEncryptingKeyByName(ctx context.Context, keyName, backendUUID string, f ManagedEncryptingKeyConsumer) error
	// WithManagedSigningKeyByUUID retrieves an instantiated managed signing key for consumption by the given function,
	// with the same semantics as WithManagedKeyByUUID
	WithManagedEncryptingKeyByUUID(ctx context.Context, keyUuid, backendUUID string, f ManagedEncryptingKeyConsumer) error
	// WithManagedMACKeyByName retrieves an instantiated managed MAC key by name for consumption by the given function,
	// with the same semantics as WithManagedKeyByName.
	WithManagedMACKeyByName(ctx context.Context, keyName, backendUUID string, f ManagedMACKeyConsumer) error
	// WithManagedMACKeyByUUID retrieves an instantiated managed MAC key by UUID for consumption by the given function,
	// with the same semantics as WithManagedKeyByUUID.
	WithManagedMACKeyByUUID(ctx context.Context, keyUUID, backendUUID string, f ManagedMACKeyConsumer) error
}

type ManagedAsymmetricKey interface {
	ManagedKey
	GetPublicKey(ctx context.Context) (crypto.PublicKey, error)
}

type ManagedKeyLifecycle interface {
	// GenerateKey generates a key in the KMS if it didn't yet exist, returning the id.
	// If it already existed, returns the existing id.  KMSKey's key material is ignored if present.
	GenerateKey(ctx context.Context) (string, error)
}

type ManagedSigningKey interface {
	ManagedAsymmetricKey

	// Sign returns a digital signature of the provided value.  The SignerOpts param must provide the hash function
	// that generated the value (if any).
	// The optional randomSource specifies the source of random values and may be ignored by the implementation
	// (such as on HSMs with their own internal RNG)
	Sign(ctx context.Context, value []byte, randomSource io.Reader, opts crypto.SignerOpts) ([]byte, error)

	// Verify verifies the provided signature against the value.  The SignerOpts param must provide the hash function
	// that generated the value (if any).
	// If true is returned the signature is correct, false otherwise.
	Verify(ctx context.Context, signature, value []byte, opts crypto.SignerOpts) (bool, error)

	// GetSigner returns an implementation of crypto.Signer backed by the managed key.  This should be called
	// as needed so as to use per request contexts.
	GetSigner(context.Context) (crypto.Signer, error)
}

type ManagedEncryptingKey interface {
	ManagedKey
	Encrypt(ctx context.Context, plaintext []byte, options ...wrapping.Option) ([]byte, error)
	Decrypt(ctx context.Context, ciphertext []byte, options ...wrapping.Option) ([]byte, error)
}

type ManagedMACKey interface {
	ManagedKey

	// MAC generates a MAC tag using the provided algorithm for the provided value.
	MAC(ctx context.Context, algorithm string, data []byte) ([]byte, error)
}

type ManagedKeyRandomSource interface {
	ManagedKey

	// GetRandomBytes returns a number (specified by the count parameter) of random bytes sourced from the target managed key.
	GetRandomBytes(count int) ([]byte, error)
}
