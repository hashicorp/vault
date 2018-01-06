package pki

import (
	"crypto/x509"
	"fmt"
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"sync/atomic"
)

const defaultSafetyBufferDuration = 259200 //seconds

func pathTidy(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "tidy",
		Fields: map[string]*framework.FieldSchema{
			"tidy_cert_store": &framework.FieldSchema{
				Type: framework.TypeBool,
				Description: `Set to true to enable tidying up
the certificate store`,
				Default: false,
			},

			"tidy_revocation_list": &framework.FieldSchema{
				Type: framework.TypeBool,
				Description: `Set to true to enable tidying up
the revocation list`,
				Default: false,
			},

			"safety_buffer": &framework.FieldSchema{
				Type: framework.TypeDurationSecond,
				Description: `The amount of extra time that must have passed
beyond certificate expiration before it is removed
from the backend storage and/or revocation list.
Defaults to 72 hours.`,
				Default: defaultSafetyBufferDuration, //72h, but TypeDurationSecond currently requires defaults to be int
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathTidyWrite,
		},

		HelpSynopsis:    pathTidyHelpSyn,
		HelpDescription: pathTidyHelpDesc,
	}
}

func (b *backend) tidyPKI(
	req *logical.Request, bufferDuration time.Duration, tidyCertStore bool, tidyRevocationList bool) error {
	if !atomic.CompareAndSwapInt32(&b.tidyRunning, 0, 1) {
		return fmt.Errorf("error running tidy pki: operation is already in progress")
	}

	if tidyCertStore {
		serials, err := req.Storage.List("certs/")
		if err != nil {
			return fmt.Errorf("error fetching list of certs: %s", err)
		}

		for _, serial := range serials {
			certEntry, err := req.Storage.Get("certs/" + serial)
			if err != nil {
				return fmt.Errorf("error fetching certificate %s: %s", serial, err)
			}

			if certEntry == nil {
				return fmt.Errorf("certificate entry for serial %s is nil", serial)
			}

			if certEntry.Value == nil || len(certEntry.Value) == 0 {
				return fmt.Errorf("found entry for serial %s but actual certificate is empty", serial)
			}

			cert, err := x509.ParseCertificate(certEntry.Value)
			if err != nil {
				return fmt.Errorf("unable to parse stored certificate with serial %s: %s", serial, err)
			}

			if time.Now().After(cert.NotAfter.Add(bufferDuration)) {
				if err := req.Storage.Delete("certs/" + serial); err != nil {
					return fmt.Errorf("error deleting serial %s from storage: %s", serial, err)
				}
			}
		}
	}

	if tidyRevocationList {
		b.revokeStorageLock.Lock()
		defer b.revokeStorageLock.Unlock()

		tidiedRevoked := false

		revokedSerials, err := req.Storage.List("revoked/")
		if err != nil {
			return fmt.Errorf("error fetching list of revoked certs: %s", err)
		}

		var revInfo revocationInfo
		for _, serial := range revokedSerials {
			revokedEntry, err := req.Storage.Get("revoked/" + serial)
			if err != nil {
				return fmt.Errorf("unable to fetch revoked cert with serial %s: %s", serial, err)
			}
			if revokedEntry == nil {
				return fmt.Errorf("revoked certificate entry for serial %s is nil", serial)
			}
			if revokedEntry.Value == nil || len(revokedEntry.Value) == 0 {
				// TODO: In this case, remove it and continue? How likely is this to
				// happen? Alternately, could skip it entirely, or could implement a
				// delete function so that there is a way to remove these
				return fmt.Errorf("found revoked serial but actual certificate is empty")
			}

			err = revokedEntry.DecodeJSON(&revInfo)
			if err != nil {
				return fmt.Errorf("error decoding revocation entry for serial %s: %s", serial, err)
			}

			revokedCert, err := x509.ParseCertificate(revInfo.CertificateBytes)
			if err != nil {
				return fmt.Errorf("unable to parse stored revoked certificate with serial %s: %s", serial, err)
			}

			if time.Now().After(revokedCert.NotAfter.Add(bufferDuration)) {
				if err := req.Storage.Delete("revoked/" + serial); err != nil {
					return fmt.Errorf("error deleting serial %s from revoked list: %s", serial, err)
				}
				tidiedRevoked = true
			}
		}

		if tidiedRevoked {
			if err := buildCRL(b, req); err != nil {
				return err
			}
		}
	}

	// Unset flag
	defer atomic.StoreInt32(&b.tidyRunning, 0)

	return nil
}

func (b *backend) pathTidyWrite(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	safetyBuffer := d.Get("safety_buffer").(int)
	tidyCertStore := d.Get("tidy_cert_store").(bool)
	tidyRevocationList := d.Get("tidy_revocation_list").(bool)

	bufferDuration := time.Duration(safetyBuffer) * time.Second

	return nil, b.tidyPKI(req, bufferDuration, tidyCertStore, tidyRevocationList)
}

const pathTidyHelpSyn = `
Tidy up the backend by removing expired certificates, revocation information,
or both.
`

const pathTidyHelpDesc = `
This endpoint allows expired certificates and/or revocation information to be
removed from the backend, freeing up storage and shortening CRLs.

For safety, this function is a noop if called without parameters; cleanup from
normal certificate storage must be enabled with 'tidy_cert_store' and cleanup
from revocation information must be enabled with 'tidy_revocation_list'.

The 'safety_buffer' parameter is useful to ensure that clock skew amongst your
hosts cannot lead to a certificate being removed from the CRL while it is still
considered valid by other hosts (for instance, if their clocks are a few
minutes behind). The 'safety_buffer' parameter can be an integer number of
seconds or a string duration like "72h".

All certificates and/or revocation information currently stored in the backend
will be checked when this endpoint is hit. The expiration of the
certificate/revocation information of each certificate being held in
certificate storage or in revocation infomation will then be checked. If the
current time, minus the value of 'safety_buffer', is greater than the
expiration, it will be removed.
`
