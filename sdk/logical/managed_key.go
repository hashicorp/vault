package logical

import (
	"context"
	"crypto"
	"io"
)

type ManagedKey interface {
	Name() string
	// Present returns true if the key is established in the KMS.  This may return false if for example
	// an HSM library is not configured on all cluster nodes.
	Present(ctx context.Context) (bool, error)
	Finalize(context.Context) error
}

type ManagedKeySystemView interface {
	GetManagedKey(ctx context.Context, name string) (ManagedKey, error)
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
