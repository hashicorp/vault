//go:build !enterprise

package pki

import (
	"context"
	"encoding/pem"
	"errors"
	"io"

	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/logical"
)

var errEntOnly = errors.New("managed keys are supported within enterprise edition only")

func generateCABundle(ctx context.Context, _ *backend, input *inputBundle, data *certutil.CreationBundle, randomSource io.Reader) (*certutil.ParsedCertBundle, error) {
	if kmsRequested(input) {
		return nil, errEntOnly
	}
	if existingKeyRequested(input) {
		keyRef, err := getExistingKeyRef(input.apiData)
		if err != nil {
			return nil, err
		}
		return certutil.CreateCertificateWithKeyGenerator(data, randomSource, existingGeneratePrivateKey(ctx, input.req.Storage, keyRef))
	}
	return certutil.CreateCertificateWithRandomSource(data, randomSource)
}

func generateCSRBundle(_ context.Context, _ *backend, input *inputBundle, data *certutil.CreationBundle, addBasicConstraints bool, randomSource io.Reader) (*certutil.ParsedCSRBundle, error) {
	if kmsRequested(input) {
		return nil, errEntOnly
	}

	return certutil.CreateCSRWithRandomSource(data, addBasicConstraints, randomSource)
}

func parseCABundle(_ context.Context, _ *backend, _ *logical.Request, bundle *certutil.CertBundle) (*certutil.ParsedCertBundle, error) {
	return bundle.ToParsedCertBundle()
}

func withManagedPKIKey(_ context.Context, _ *backend, _ managedKeyId, _ string, _ logical.ManagedSigningKeyConsumer) error {
	return errEntOnly
}

func existingGeneratePrivateKey(ctx context.Context, s logical.Storage, keyRef string) certutil.KeyGenerator {
	return func(keyType string, keyBits int, container certutil.ParsedPrivateKeyContainer, _ io.Reader) error {
		keyId, err := resolveKeyReference(ctx, s, keyRef)
		if err != nil {
			return err
		}
		key, err := fetchKeyById(ctx, s, keyId)
		if err != nil {
			return err
		}
		signer, err := key.GetSigner()
		if err != nil {
			return err
		}
		privateKeyType := certutil.GetPrivateKeyTypeFromSigner(signer)
		if privateKeyType == certutil.UnknownPrivateKey {
			return errors.New("unknown private key type loaded from key id: " + keyId.String())
		}
		blk, _ := pem.Decode([]byte(key.PrivateKey))
		container.SetParsedPrivateKey(signer, privateKeyType, blk.Bytes)
		return nil
	}
}
