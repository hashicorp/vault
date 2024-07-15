// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"context"
	"crypto/x509"
	"fmt"

	"github.com/hashicorp/vault/builtin/logical/pki/issuing"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// SecretCertsType is the name used to identify this type
const SecretCertsType = "pki"

func secretCerts(b *backend) *framework.Secret {
	return &framework.Secret{
		Type: SecretCertsType,
		Fields: map[string]*framework.FieldSchema{
			"certificate": {
				Type: framework.TypeString,
				Description: `The PEM-encoded concatenated certificate and
issuing certificate authority`,
			},
			"private_key": {
				Type:        framework.TypeString,
				Description: "The PEM-encoded private key for the certificate",
			},
			"serial": {
				Type: framework.TypeString,
				Description: `The serial number of the certificate, for handy
reference`,
			},
		},

		Revoke: b.secretCredsRevoke,
	}
}

func (b *backend) secretCredsRevoke(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	if req.Secret == nil {
		return nil, fmt.Errorf("secret is nil in request")
	}

	serialInt, ok := req.Secret.InternalData["serial_number"]
	if !ok {
		return nil, fmt.Errorf("could not find serial in internal secret data")
	}

	b.GetRevokeStorageLock().Lock()
	defer b.GetRevokeStorageLock().Unlock()

	sc := b.makeStorageContext(ctx, req.Storage)
	serial := serialInt.(string)

	certEntry, err := fetchCertBySerial(sc, issuing.PathCerts, serial)
	if err != nil {
		return nil, err
	}
	if certEntry == nil {
		// We can't write to revoked/ or update the CRL anyway because we don't have the cert,
		// and there's no reason to expect this will work on a subsequent
		// retry.  Just give up and let the lease get deleted.
		b.Logger().Warn("expired certificate revoke failed because not found in storage, treating as success", "serial", serial)
		return nil, nil
	}

	cert, err := x509.ParseCertificate(certEntry.Value)
	if err != nil {
		return nil, fmt.Errorf("error parsing certificate: %w", err)
	}

	// Compatibility: Don't revoke CAs if they had leases. New CAs going forward aren't issued leases.
	if cert.IsCA {
		return nil, nil
	}

	config, err := sc.CrlBuilder().GetConfigWithUpdate(sc)
	if err != nil {
		return nil, fmt.Errorf("error revoking serial: %s: failed reading config: %w", serial, err)
	}

	return revokeCert(sc, config, cert)
}
