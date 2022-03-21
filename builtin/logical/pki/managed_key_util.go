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

func generateCABundle(_ context.Context, _ *backend, input *inputBundle, data *certutil.CreationBundle, randomSource io.Reader) (*certutil.ParsedCertBundle, error) {
	if kmsRequested(input) {
		return nil, errEntOnly
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

func withManagedPKIKey(_ context.Context, _ *backend, _ keyId, _ string, _ logical.ManagedSigningKeyConsumer) error {
	return errEntOnly
}
