package pki

import (
	"context"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/helper/errutil"
	"github.com/hashicorp/vault/logical"
)

type revocationInfo struct {
	CertificateBytes  []byte    `json:"certificate_bytes"`
	RevocationTime    int64     `json:"revocation_time"`
	RevocationTimeUTC time.Time `json:"revocation_time_utc"`
}

// Revokes a cert, and tries to be smart about error recovery
func revokeCert(ctx context.Context, b *backend, req *logical.Request, serial string, fromLease bool) (*logical.Response, error) {
	// As this backend is self-contained and this function does not hook into
	// third parties to manage users or resources, if the mount is tainted,
	// revocation doesn't matter anyways -- the CRL that would be written will
	// be immediately blown away by the view being cleared. So we can simply
	// fast path a successful exit.
	if b.System().Tainted() {
		return nil, nil
	}

	alreadyRevoked := false
	var revInfo revocationInfo

	revEntry, err := fetchCertBySerial(ctx, req, "revoked/", serial)
	if err != nil {
		switch err.(type) {
		case errutil.UserError:
			return logical.ErrorResponse(err.Error()), nil
		case errutil.InternalError:
			return nil, err
		}
	}
	if revEntry != nil {
		// Set the revocation info to the existing values
		alreadyRevoked = true
		err = revEntry.DecodeJSON(&revInfo)
		if err != nil {
			return nil, fmt.Errorf("error decoding existing revocation info")
		}
	}

	if !alreadyRevoked {
		certEntry, err := fetchCertBySerial(ctx, req, "certs/", serial)
		if err != nil {
			switch err.(type) {
			case errutil.UserError:
				return logical.ErrorResponse(err.Error()), nil
			case errutil.InternalError:
				return nil, err
			}
		}
		if certEntry == nil {
			return logical.ErrorResponse(fmt.Sprintf("certificate with serial %s not found", serial)), nil
		}

		cert, err := x509.ParseCertificate(certEntry.Value)
		if err != nil {
			return nil, errwrap.Wrapf("error parsing certificate: {{err}}", err)
		}
		if cert == nil {
			return nil, fmt.Errorf("got a nil certificate")
		}

		if cert.NotAfter.Before(time.Now()) {
			return nil, nil
		}

		// Compatibility: Don't revoke CAs if they had leases. New CAs going
		// forward aren't issued leases.
		if cert.IsCA && fromLease {
			return nil, nil
		}

		currTime := time.Now()
		revInfo.CertificateBytes = certEntry.Value
		revInfo.RevocationTime = currTime.Unix()
		revInfo.RevocationTimeUTC = currTime.UTC()

		revEntry, err = logical.StorageEntryJSON("revoked/"+normalizeSerial(serial), revInfo)
		if err != nil {
			return nil, fmt.Errorf("error creating revocation entry")
		}

		err = req.Storage.Put(ctx, revEntry)
		if err != nil {
			return nil, fmt.Errorf("error saving revoked certificate to new location")
		}

	}

	crlErr := buildCRL(ctx, b, req)
	switch crlErr.(type) {
	case errutil.UserError:
		return logical.ErrorResponse(fmt.Sprintf("Error during CRL building: %s", crlErr)), nil
	case errutil.InternalError:
		return nil, errwrap.Wrapf("error encountered during CRL building: {{err}}", crlErr)
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"revocation_time": revInfo.RevocationTime,
		},
	}
	if !revInfo.RevocationTimeUTC.IsZero() {
		resp.Data["revocation_time_rfc3339"] = revInfo.RevocationTimeUTC.Format(time.RFC3339Nano)
	}
	return resp, nil
}

// Builds a CRL by going through the list of revoked certificates and building
// a new CRL with the stored revocation times and serial numbers.
func buildCRL(ctx context.Context, b *backend, req *logical.Request) error {
	revokedSerials, err := req.Storage.List(ctx, "revoked/")
	if err != nil {
		return errutil.InternalError{Err: fmt.Sprintf("error fetching list of revoked certs: %s", err)}
	}

	revokedCerts := []pkix.RevokedCertificate{}
	var revInfo revocationInfo
	for _, serial := range revokedSerials {
		revokedEntry, err := req.Storage.Get(ctx, "revoked/"+serial)
		if err != nil {
			return errutil.InternalError{Err: fmt.Sprintf("unable to fetch revoked cert with serial %s: %s", serial, err)}
		}
		if revokedEntry == nil {
			return errutil.InternalError{Err: fmt.Sprintf("revoked certificate entry for serial %s is nil", serial)}
		}
		if revokedEntry.Value == nil || len(revokedEntry.Value) == 0 {
			// TODO: In this case, remove it and continue? How likely is this to
			// happen? Alternately, could skip it entirely, or could implement a
			// delete function so that there is a way to remove these
			return errutil.InternalError{Err: fmt.Sprintf("found revoked serial but actual certificate is empty")}
		}

		err = revokedEntry.DecodeJSON(&revInfo)
		if err != nil {
			return errutil.InternalError{Err: fmt.Sprintf("error decoding revocation entry for serial %s: %s", serial, err)}
		}

		revokedCert, err := x509.ParseCertificate(revInfo.CertificateBytes)
		if err != nil {
			return errutil.InternalError{Err: fmt.Sprintf("unable to parse stored revoked certificate with serial %s: %s", serial, err)}
		}

		// NOTE: We have to change this to UTC time because the CRL standard
		// mandates it but Go will happily encode the CRL without this.
		newRevCert := pkix.RevokedCertificate{
			SerialNumber: revokedCert.SerialNumber,
		}
		if !revInfo.RevocationTimeUTC.IsZero() {
			newRevCert.RevocationTime = revInfo.RevocationTimeUTC
		} else {
			newRevCert.RevocationTime = time.Unix(revInfo.RevocationTime, 0).UTC()
		}
		revokedCerts = append(revokedCerts, newRevCert)
	}

	signingBundle, caErr := fetchCAInfo(ctx, req)
	switch caErr.(type) {
	case errutil.UserError:
		return errutil.UserError{Err: fmt.Sprintf("could not fetch the CA certificate: %s", caErr)}
	case errutil.InternalError:
		return errutil.InternalError{Err: fmt.Sprintf("error fetching CA certificate: %s", caErr)}
	}

	crlLifetime := b.crlLifetime
	crlInfo, err := b.CRL(ctx, req.Storage)
	if err != nil {
		return errutil.InternalError{Err: fmt.Sprintf("error fetching CRL config information: %s", err)}
	}
	if crlInfo != nil {
		crlDur, err := time.ParseDuration(crlInfo.Expiry)
		if err != nil {
			return errutil.InternalError{Err: fmt.Sprintf("error parsing CRL duration of %s", crlInfo.Expiry)}
		}
		crlLifetime = crlDur
	}

	crlBytes, err := signingBundle.Certificate.CreateCRL(rand.Reader, signingBundle.PrivateKey, revokedCerts, time.Now(), time.Now().Add(crlLifetime))
	if err != nil {
		return errutil.InternalError{Err: fmt.Sprintf("error creating new CRL: %s", err)}
	}

	err = req.Storage.Put(ctx, &logical.StorageEntry{
		Key:   "crl",
		Value: crlBytes,
	})
	if err != nil {
		return errutil.InternalError{Err: fmt.Sprintf("error storing CRL: %s", err)}
	}

	return nil
}
