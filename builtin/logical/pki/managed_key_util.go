//go:build !enterprise

package pki

import (
	"context"
	"crypto"
	"errors"
	"io"

	"github.com/hashicorp/vault/sdk/helper/certutil"
)

var errEntOnly = errors.New("managed keys are supported within enterprise edition only")

func generateManagedKeyCABundle(ctx context.Context, b *backend, keyId managedKeyId, data *certutil.CreationBundle, randomSource io.Reader) (bundle *certutil.ParsedCertBundle, err error) {
	return nil, errEntOnly
}

func generateManagedKeyCSRBundle(ctx context.Context, b *backend, keyId managedKeyId, data *certutil.CreationBundle, addBasicConstraints bool, randomSource io.Reader) (bundle *certutil.ParsedCSRBundle, err error) {
	return nil, errEntOnly
}

func getManagedKeyPublicKey(ctx context.Context, b *backend, keyId managedKeyId) (crypto.PublicKey, error) {
	return nil, errEntOnly
}

func parseManagedKeyCABundle(ctx context.Context, b *backend, bundle *certutil.CertBundle) (*certutil.ParsedCertBundle, error) {
	return nil, errEntOnly
}

func extractManagedKeyId(privateKeyBytes []byte) (UUIDKey, error) {
	return "", errEntOnly
}

func createKmsKeyBundle(ctx context.Context, b *backend, keyId managedKeyId) (certutil.KeyBundle, certutil.PrivateKeyType, error) {
	return certutil.KeyBundle{}, certutil.UnknownPrivateKey, errEntOnly
}

func getManagedKeyInfo(ctx context.Context, b *backend, keyId managedKeyId) (*managedKeyInfo, error) {
	return nil, errEntOnly
}
