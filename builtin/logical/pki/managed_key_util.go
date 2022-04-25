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

func generateManagedKeyCABundle(_ context.Context, _ *backend, _ *inputBundle, _ *certutil.CreationBundle, _ io.Reader) (*certutil.ParsedCertBundle, error) {
	return nil, errEntOnly
}

func generateManagedKeyCSRBundle(_ context.Context, _ *backend, _ *inputBundle, _ *certutil.CreationBundle, _ bool, _ io.Reader) (*certutil.ParsedCSRBundle, error) {
	return nil, errEntOnly
}

func parseManagedKeyCABundle(_ context.Context, _ *backend, _ *logical.Request, _ *certutil.CertBundle) (*certutil.ParsedCertBundle, error) {
	return nil, errEntOnly
}

func withManagedPKIKey(_ context.Context, _ *backend, _ managedKeyId, _ string, _ logical.ManagedSigningKeyConsumer) error {
	return errEntOnly
}
