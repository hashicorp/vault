package pki

import (
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"time"

	"github.com/hashicorp/vault/helper/certutil"
	"github.com/hashicorp/vault/logical"
)

type revocationInfo struct {
	CertificateBytes []byte `json:"certificate_bytes"`
	RevocationTime   int64  `json:"revocation_time"`
}

// Revokes a cert, and tries to be smart about error recovery
func revokeCert(b *backend, req *logical.Request, serial string) (*logical.Response, error) {
	alreadyRevoked := false
	var revInfo revocationInfo

	certEntry, err := fetchCertBySerial(req, "revoked/", serial)
	// Don't check error because it's expected that it may fail here;
	// just check for existence
	if certEntry != nil {
		// Verify that it is also deleted from certs/
		// in case of partial failure from an earlier run.
		certEntry, _ = fetchCertBySerial(req, "certs/", serial)
		if certEntry == nil {
			// Everything seems sane, so don't rebuild the CRL
			return nil, nil
		}

		// Still exists in certs/; set the revocation info, below it will
		// be removed from certs/ and the CRL rotated
		alreadyRevoked = true

		revEntry, err := req.Storage.Get("revoked/" + serial)
		if revEntry == nil || err != nil {
			return nil, fmt.Errorf("Error getting existing revocation info")
		}

		err = revEntry.DecodeJSON(&revInfo)
		if err != nil {
			return nil, fmt.Errorf("Error decoding existing revocation info")
		}
	}

	if !alreadyRevoked {
		certEntry, err = fetchCertBySerial(req, "certs/", serial)
		switch err.(type) {
		case certutil.UserError:
			return logical.ErrorResponse(err.Error()), nil
		case certutil.InternalError:
			return nil, err
		}

		cert, err := x509.ParseCertificate(certEntry.Value)
		if err != nil {
			return nil, fmt.Errorf("Error parsing certificate")
		}
		if cert == nil {
			return nil, fmt.Errorf("Got a nil certificate")
		}

		if cert.NotAfter.Before(time.Now()) {
			return nil, nil
		}

		revInfo.CertificateBytes = certEntry.Value
		revInfo.RevocationTime = time.Now().Unix()

		certEntry, err = logical.StorageEntryJSON("revoked/"+serial, revInfo)
		if err != nil {
			return nil, fmt.Errorf("Error creating revocation entry")
		}

		err = req.Storage.Put(certEntry)
		if err != nil {
			return nil, fmt.Errorf("Error saving revoked certificate to new location")
		}

	}

	crlErr := buildCRL(b, req)
	switch crlErr.(type) {
	case certutil.UserError:
		return logical.ErrorResponse(fmt.Sprintf("Error during CRL building: %s", crlErr)), nil
	case certutil.InternalError:
		return nil, fmt.Errorf("Error encountered during CRL building: %s", crlErr)
	}

	err = req.Storage.Delete("certs/" + serial)

	if err != nil {
		return nil, fmt.Errorf("Error deleting cert from valid-certs location")
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"revocation_time": revInfo.RevocationTime,
		},
	}, nil
}

// Builds a CRL by going through the list of revoked certificates and building
// a new CRL with the stored revocation times and serial numbers.
//
// If a certificate has already expired, it will be removed entirely rather than
// become part of the new CRL.
func buildCRL(b *backend, req *logical.Request) error {
	revokedSerials, err := req.Storage.List("revoked/")
	if err != nil {
		return certutil.InternalError{Err: fmt.Sprintf("Error fetching list of revoked certs: %s", err)}
	}

	revokedCerts := []pkix.RevokedCertificate{}
	var revInfo revocationInfo
	for _, serial := range revokedSerials {
		revokedEntry, err := req.Storage.Get("revoked/" + serial)
		if err != nil {
			return certutil.InternalError{Err: fmt.Sprintf("Unable to fetch revoked cert with serial %s: %s", serial, err)}
		}
		if revokedEntry == nil {
			return certutil.InternalError{Err: fmt.Sprintf("Revoked certificate entry for serial %s is nil", serial)}
		}
		if revokedEntry.Value == nil || len(revokedEntry.Value) == 0 {
			// TODO: In this case, remove it and continue? How likely is this to
			// happen? Alternately, could skip it entirely, or could implement a
			// delete function so that there is a way to remove these
			return certutil.InternalError{Err: fmt.Sprintf("Found revoked serial but actual certificate is empty")}
		}

		err = revokedEntry.DecodeJSON(&revInfo)
		if err != nil {
			return certutil.InternalError{Err: fmt.Sprintf("Error decoding revocation entry for serial %s: %s", serial, err)}
		}

		revokedCert, err := x509.ParseCertificate(revInfo.CertificateBytes)
		if err != nil {
			return certutil.InternalError{Err: fmt.Sprintf("Unable to parse stored revoked certificate with serial %s: %s", serial, err)}
		}

		if revokedCert.NotAfter.Before(time.Now()) {
			err = req.Storage.Delete(serial)
			if err != nil {
				return certutil.InternalError{Err: fmt.Sprintf("Unable to delete revoked, expired certificate with serial %s: %s", serial, err)}
			}
			continue
		}

		revokedCerts = append(revokedCerts, pkix.RevokedCertificate{
			SerialNumber:   revokedCert.SerialNumber,
			RevocationTime: time.Unix(revInfo.RevocationTime, 0),
		})
	}

	signingBundle, caErr := fetchCAInfo(req)
	switch caErr.(type) {
	case certutil.UserError:
		return certutil.UserError{Err: fmt.Sprintf("Could not fetch the CA certificate: %s", caErr)}
	case certutil.InternalError:
		return certutil.InternalError{Err: fmt.Sprintf("Error fetching CA certificate: %s", caErr)}
	}

	crlLifetime := b.crlLifetime
	crlInfo, err := b.CRL(req.Storage)
	if err != nil {
		return certutil.InternalError{Err: fmt.Sprintf("Error fetching CRL config information: %s", err)}
	}
	if crlInfo != nil {
		crlDur, err := time.ParseDuration(crlInfo.Expiry)
		if err != nil {
			return certutil.InternalError{Err: fmt.Sprintf("Error parsing CRL duration of %s", crlInfo.Expiry)}
		}
		crlLifetime = crlDur
	}

	crlBytes, err := signingBundle.Certificate.CreateCRL(rand.Reader, signingBundle.PrivateKey, revokedCerts, time.Now(), time.Now().Add(crlLifetime))
	if err != nil {
		return certutil.InternalError{Err: fmt.Sprintf("Error creating new CRL: %s", err)}
	}

	err = req.Storage.Put(&logical.StorageEntry{
		Key:   "crl",
		Value: crlBytes,
	})
	if err != nil {
		return certutil.InternalError{Err: fmt.Sprintf("Error storing CRL: %s", err)}
	}

	return nil
}
