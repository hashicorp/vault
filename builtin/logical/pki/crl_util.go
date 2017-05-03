package pki

import (
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"time"

	"github.com/hashicorp/vault/helper/errutil"
	"github.com/hashicorp/vault/logical"
)

type revocationInfo struct {
	CertificateBytes  []byte    `json:"certificate_bytes"`
	RevocationTime    int64     `json:"revocation_time"`
	RevocationTimeUTC time.Time `json:"revocation_time_utc"`
}

// Revokes a cert, and tries to be smart about error recovery
func revokeCert(b *backend, req *logical.Request, serial string, fromLease bool) (*logical.Response, error) {
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

	revEntry, err := fetchCertBySerial(req, "revoked/", serial)
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
			return nil, fmt.Errorf("Error decoding existing revocation info")
		}
	}

	if !alreadyRevoked {
		certEntry, err := fetchCertBySerial(req, "certs/", serial)
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
			return nil, fmt.Errorf("Error parsing certificate: %s", err)
		}
		if cert == nil {
			return nil, fmt.Errorf("Got a nil certificate")
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
			return nil, fmt.Errorf("Error creating revocation entry")
		}

		err = req.Storage.Put(revEntry)
		if err != nil {
			return nil, fmt.Errorf("Error saving revoked certificate to new location")
		}

	}

	crlErr := buildCRL(b, req)
	switch crlErr.(type) {
	case errutil.UserError:
		return logical.ErrorResponse(fmt.Sprintf("Error during CRL building: %s", crlErr)), nil
	case errutil.InternalError:
		return nil, fmt.Errorf("Error encountered during CRL building: %s", crlErr)
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
func buildCRL(b *backend, req *logical.Request) error {
	revokedSerials, err := req.Storage.List("revoked/")
	if err != nil {
		return errutil.InternalError{Err: fmt.Sprintf("Error fetching list of revoked certs: %s", err)}
	}

	revokedCerts := []pkix.RevokedCertificate{}
	var revInfo revocationInfo
	for _, serial := range revokedSerials {
		revokedEntry, err := req.Storage.Get("revoked/" + serial)
		if err != nil {
			return errutil.InternalError{Err: fmt.Sprintf("Unable to fetch revoked cert with serial %s: %s", serial, err)}
		}
		if revokedEntry == nil {
			return errutil.InternalError{Err: fmt.Sprintf("Revoked certificate entry for serial %s is nil", serial)}
		}
		if revokedEntry.Value == nil || len(revokedEntry.Value) == 0 {
			// TODO: In this case, remove it and continue? How likely is this to
			// happen? Alternately, could skip it entirely, or could implement a
			// delete function so that there is a way to remove these
			return errutil.InternalError{Err: fmt.Sprintf("Found revoked serial but actual certificate is empty")}
		}

		err = revokedEntry.DecodeJSON(&revInfo)
		if err != nil {
			return errutil.InternalError{Err: fmt.Sprintf("Error decoding revocation entry for serial %s: %s", serial, err)}
		}

		revokedCert, err := x509.ParseCertificate(revInfo.CertificateBytes)
		if err != nil {
			return errutil.InternalError{Err: fmt.Sprintf("Unable to parse stored revoked certificate with serial %s: %s", serial, err)}
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

	signingBundle, caErr := fetchCAInfo(req)
	switch caErr.(type) {
	case errutil.UserError:
		return errutil.UserError{Err: fmt.Sprintf("Could not fetch the CA certificate: %s", caErr)}
	case errutil.InternalError:
		return errutil.InternalError{Err: fmt.Sprintf("Error fetching CA certificate: %s", caErr)}
	}

	crlLifetime := b.crlLifetime
	crlInfo, err := b.CRL(req.Storage)
	if err != nil {
		return errutil.InternalError{Err: fmt.Sprintf("Error fetching CRL config information: %s", err)}
	}
	if crlInfo != nil {
		crlDur, err := time.ParseDuration(crlInfo.Expiry)
		if err != nil {
			return errutil.InternalError{Err: fmt.Sprintf("Error parsing CRL duration of %s", crlInfo.Expiry)}
		}
		crlLifetime = crlDur
	}

	crlBytes, err := signingBundle.Certificate.CreateCRL(rand.Reader, signingBundle.PrivateKey, revokedCerts, time.Now(), time.Now().Add(crlLifetime))
	if err != nil {
		return errutil.InternalError{Err: fmt.Sprintf("Error creating new CRL: %s", err)}
	}

	err = req.Storage.Put(&logical.StorageEntry{
		Key:   "crl",
		Value: crlBytes,
	})
	if err != nil {
		return errutil.InternalError{Err: fmt.Sprintf("Error storing CRL: %s", err)}
	}

	return nil
}
