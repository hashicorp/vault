package pki

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/errutil"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathRevoke(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: `revoke`,
		Fields: map[string]*framework.FieldSchema{
			"serial_number": {
				Type: framework.TypeString,
				Description: `Certificate serial number, in colon- or
hyphen-separated octal`,
			},
			"certificate": {
				Type: framework.TypeString,
				Description: `Certificate to revoke in PEM format; must be
signed by an issuer in this mount.`,
			},
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.metricsWrap("revoke", noRole, b.pathRevokeWrite),
				// This should never be forwarded. See backend.go for more information.
				// If this needs to write, the entire request will be forwarded to the
				// active node of the current performance cluster, but we don't want to
				// forward invalid revoke requests there.
			},
		},

		HelpSynopsis:    pathRevokeHelpSyn,
		HelpDescription: pathRevokeHelpDesc,
	}
}

func pathRotateCRL(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: `crl/rotate`,

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathRotateCRLRead,
				// See backend.go; we will read a lot of data prior to calling write,
				// so this request should be forwarded when it is first seen, not
				// when it is ready to write.
				ForwardPerformanceStandby: true,
			},
		},

		HelpSynopsis:    pathRotateCRLHelpSyn,
		HelpDescription: pathRotateCRLHelpDesc,
	}
}

func (b *backend) pathRevokeWriteHandleCertificate(ctx context.Context, req *logical.Request, data *framework.FieldData, certPem string) (string, []byte, error) {
	// This function handles just the verification of the certificate against
	// the global issuer set, checking whether or not it is importable.
	//
	// We return the parsed serial number, an optionally-nil byte array to
	// write out to disk, and an error if one occurred.
	//
	// First start by parsing the certificate.
	pemBlock, _ := pem.Decode([]byte(certPem))
	if pemBlock == nil {
		return "", nil, errutil.UserError{Err: "certificate contains no PEM data"}
	}

	certReference, err := x509.ParseCertificate(pemBlock.Bytes)
	if err != nil {
		return "", nil, errutil.UserError{Err: fmt.Sprintf("certificate could not be parsed: %v", err)}
	}

	// Ensure we have a well-formed serial number before continuing.
	serial := serialFromCert(certReference)
	if len(serial) == 0 {
		return "", nil, errutil.UserError{Err: fmt.Sprintf("invalid serial number on presented certificate")}
	}

	// We have two approaches here: we could start verifying against issuers
	// (which involves fetching and parsing them), or we could see if, by
	// some chance we've already imported it (cheap). The latter tells us
	// if we happen to have a serial number collision (which shouldn't
	// happen in practice) versus an already-imported cert (which might
	// happen and its fine to handle safely).
	//
	// Start with the latter since its cheaper. Fetch the cert (by serial)
	// and if it exists, compare the contents.
	certEntry, err := fetchCertBySerial(ctx, b, req, req.Path, serial)
	if err != nil {
		return serial, nil, err
	}

	if certEntry != nil {
		// As seen with importing issuers, it is best to parse the certificate
		// and compare parsed values, rather than attempting to infer equality
		// from the raw data.
		certReferenceStored, err := x509.ParseCertificate(certEntry.Value)
		if err != nil {
			return serial, nil, err
		}

		if !areCertificatesEqual(certReference, certReferenceStored) {
			// Here we refuse the import with an error because the two certs
			// are unequal but we would've otherwise overwritten the existing
			// copy.
			return serial, nil, fmt.Errorf("certificate with same serial but unequal value already present in this cluster's storage; refusing to revoke")
		} else {
			// Otherwise, we can return without an error as we've already
			// imported this certificate, likely when we issued it. We don't
			// need to re-verify the signature as we assume it was already
			// verified when it was imported.
			return serial, nil, nil
		}
	}

	// Otherwise, we must not have a stored copy. From here on out, the second
	// parameter (except in error cases) should be the copy to write out.
	//
	// Fetch and iterate through each issuer.
	sc := b.makeStorageContext(ctx, req.Storage)
	issuers, err := sc.listIssuers()
	if err != nil {
		return serial, nil, err
	}

	foundMatchingIssuer := false
	for _, issuerId := range issuers {
		issuer, err := sc.fetchIssuerById(issuerId)
		if err != nil {
			return serial, nil, err
		}

		issuerCert, err := issuer.GetCertificate()
		if err != nil {
			return serial, nil, err
		}

		if err := certReference.CheckSignatureFrom(issuerCert); err == nil {
			// If the signature was valid, we found our match and can safely
			// exit.
			foundMatchingIssuer = true
			break
		}
	}

	if foundMatchingIssuer {
		return serial, certReference.Raw, nil
	}

	return serial, nil, errutil.UserError{Err: fmt.Sprintf("unable to verify signature on presented cert from any present issuer in this mount; certificates from previous CAs will need to have their issuing CA and key re-imported if revocation is necessary")}
}

func (b *backend) pathRevokeWrite(ctx context.Context, req *logical.Request, data *framework.FieldData, _ *roleEntry) (*logical.Response, error) {
	rawSerial, haveSerial := data.GetOk("serial_number")
	rawCertificate, haveCert := data.GetOk("certificate")

	if !haveSerial && !haveCert {
		return logical.ErrorResponse("The serial number or certificate to revoke must be provided."), nil
	} else if haveSerial && haveCert {
		return logical.ErrorResponse("Must provide either the certificate or the serial to revoke; not both."), nil
	}

	var serial string
	if haveSerial {
		// Easy case: this cert should be in storage already.
		serial = rawSerial.(string)
	} else {
		// Otherwise, we've gotta parse the certificate from the request and
		// then import it into cluster-local storage. Before writing the
		// certificate (and forwarding), we want to verify this certificate
		// was actually signed by one of our present issuers.
		var err error
		var certBytes []byte
		serial, certBytes, err = b.pathRevokeWriteHandleCertificate(ctx, req, data, rawCertificate.(string))
		if err != nil {
			return nil, err
		}

		// At this point, a forward operation will occur if we're on a standby
		// node as we're now attempting to write the bytes of the cert out to
		// disk.
		if certBytes != nil {
			err = req.Storage.Put(ctx, &logical.StorageEntry{
				Key:   "certs/" + serial,
				Value: certBytes,
			})
			if err != nil {
				return nil, err
			}
		}

		// Finally, we have a valid serial number to use for BYOC revocation!
	}

	if len(serial) == 0 {
		return logical.ErrorResponse("The serial number must be provided"), nil
	}

	// Assumption: this check is cheap. Call this twice, in the cert-import
	// case, to allow cert verification to get rejected on the standby node,
	// but we still need it to protect the serial number case.
	if b.System().ReplicationState().HasState(consts.ReplicationPerformanceStandby) {
		return nil, logical.ErrReadOnly
	}

	// We store and identify by lowercase colon-separated hex, but other
	// utilities use dashes and/or uppercase, so normalize
	serial = strings.ReplaceAll(strings.ToLower(serial), "-", ":")

	b.revokeStorageLock.Lock()
	defer b.revokeStorageLock.Unlock()

	return revokeCert(ctx, b, req, serial, false)
}

func (b *backend) pathRotateCRLRead(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	b.revokeStorageLock.RLock()
	defer b.revokeStorageLock.RUnlock()

	crlErr := b.crlBuilder.rebuild(ctx, b, req, false)
	if crlErr != nil {
		switch crlErr.(type) {
		case errutil.UserError:
			return logical.ErrorResponse(fmt.Sprintf("Error during CRL building: %s", crlErr)), nil
		default:
			return nil, fmt.Errorf("error encountered during CRL building: %w", crlErr)
		}
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"success": true,
		},
	}, nil
}

const pathRevokeHelpSyn = `
Revoke a certificate by serial number.
`

const pathRevokeHelpDesc = `
This allows certificates to be revoked using its serial number. A root token is required.
`

const pathRotateCRLHelpSyn = `
Force a rebuild of the CRL.
`

const pathRotateCRLHelpDesc = `
Force a rebuild of the CRL. This can be used to remove expired certificates from it if no certificates have been revoked. A root token is required.
`
