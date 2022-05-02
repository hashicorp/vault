//go:build !enterprise

package pki

import (
	"context"
	"errors"
	"io"

	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/logical"
)

var errEntOnly = errors.New("managed keys are supported within enterprise edition only")

func generateManagedKeyCABundle(ctx context.Context, b *backend, input *inputBundle, keyId managedKeyId, data *certutil.CreationBundle, randomSource io.Reader) (bundle *certutil.ParsedCertBundle, err error) {
	return nil, errEntOnly
}

func generateManagedKeyCSRBundle(ctx context.Context, b *backend, input *inputBundle, keyId managedKeyId, data *certutil.CreationBundle, addBasicConstraints bool, randomSource io.Reader) (bundle *certutil.ParsedCSRBundle, err error) {
	return nil, errEntOnly
}

func parseManagedKeyCABundle(ctx context.Context, b *backend, req *logical.Request, bundle *certutil.CertBundle) (*certutil.ParsedCertBundle, error) {
	return nil, errEntOnly
}

func withManagedPKIKey(ctx context.Context, b *backend, keyId managedKeyId, mountPoint string, f logical.ManagedSigningKeyConsumer) error {
	return errEntOnly
}

func extractManagedKeyId(privateKeyBytes []byte) (UUIDKey, error) {
	return "", errEntOnly
}

func createKmsKeyBundle(mkc managedKeyContext, keyId managedKeyId) (certutil.KeyBundle, certutil.PrivateKeyType, error) {
	return certutil.KeyBundle{}, certutil.UnknownPrivateKey, errEntOnly
}
