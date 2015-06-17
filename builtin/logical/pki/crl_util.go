package pki

import (
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"sync"
	"time"

	"github.com/hashicorp/vault/logical"
)

type revocationInfo struct {
	CertificateBytes []byte `json:"certificate_bytes"`
	RevocationTime   int64  `json:"revocation_time"`
}

var (
	crlLifetime       = time.Hour * 72
	revokeStorageLock = &sync.Mutex{}
)

func revokeCert(req *logical.Request, serial string) (*logical.Response, error) {
	alreadyRevoked := false
	var err error

	revInfo := revocationInfo{}

	certEntry, userErr, intErr := fetchCertBySerial(req, "revoked/", serial)
	if certEntry != nil {
		// Verify that it is also deleted from certs/
		// in case of partial failure from an earlier run.
		certEntry, _, _ = fetchCertBySerial(req, "certs/", serial)
		if certEntry != nil {
			alreadyRevoked = true

			revEntry, err := req.Storage.Get("revoked/" + serial)
			if err != nil {
				return nil, fmt.Errorf("Error getting existing revocation info")
			}

			err = revEntry.DecodeJSON(&revInfo)
			if err != nil {
				return nil, fmt.Errorf("Error decoding existing revocation info")
			}
		} else {
			return nil, nil
		}
	}

	if !alreadyRevoked {
		certEntry, userErr, intErr = fetchCertBySerial(req, "certs/", serial)
		switch {
		case userErr != nil:
			return logical.ErrorResponse(userErr.Error()), nil
		case intErr != nil:
			return nil, intErr
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

	userErr, intErr = buildCRL(req)
	switch {
	case userErr != nil:
		return logical.ErrorResponse(fmt.Sprintf("Error during CRL building: %s", userErr)), nil
	case intErr != nil:
		return nil, fmt.Errorf("Error encountered during CRL building: %s", intErr)
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

func buildCRL(req *logical.Request) (error, error) {
	revokedSerials, err := req.Storage.List("revoked/")
	if err != nil {
		return nil, fmt.Errorf("Error fetching list of revoked certs: %s", err)
	}

	revokedCerts := []pkix.RevokedCertificate{}
	var revInfo revocationInfo
	for _, serial := range revokedSerials {
		revokedEntry, err := req.Storage.Get("revoked/" + serial)
		if err != nil {
			return nil, fmt.Errorf("Unable to fetch revoked cert with serial %s: %s", serial, err)
		}
		if revokedEntry == nil {
			return nil, fmt.Errorf("Revoked certificate entry for serial %s is nil", serial)
		}
		if revokedEntry.Value == nil || len(revokedEntry.Value) == 0 {
			// TODO: In this case, remove it and continue? How likely is this to
			// happen? Alternately, could skip it entirely, or could implement a
			// delete function so that there is a way to remove these
			return nil, fmt.Errorf("Found revoked serial but actual certificate is empty")
		}

		err = revokedEntry.DecodeJSON(&revInfo)
		if err != nil {
			return nil, fmt.Errorf("Error decoding revocation entry for serial %s: %s", serial, err)
		}

		revokedCert, err := x509.ParseCertificate(revInfo.CertificateBytes)
		if err != nil {
			return nil, fmt.Errorf("Unable to parse stored revoked certificate with serial %s: %s", serial, err)
		}

		if revokedCert.NotAfter.Before(time.Now()) {
			err = req.Storage.Delete(serial)
			if err != nil {
				return nil, fmt.Errorf("Unable to delete revoked, expired certificate with serial %s: %s", serial, err)
			}
			continue
		}

		revokedCerts = append(revokedCerts, pkix.RevokedCertificate{
			SerialNumber:   revokedCert.SerialNumber,
			RevocationTime: time.Unix(revInfo.RevocationTime, 0),
		})
	}

	signingBundle, caCert, userErr, intErr := fetchCAInfo(req)
	switch {
	case userErr != nil:
		return fmt.Errorf("Could not fetch the CA certificate: %s", userErr), nil
	case intErr != nil:
		return nil, fmt.Errorf("Error fetching CA certificate: %s", intErr)
	}

	// TODO: Make expiry configurable
	crlBytes, err := caCert.CreateCRL(rand.Reader, signingBundle.PrivateKey, revokedCerts, time.Now(), time.Now().Add(crlLifetime))
	if err != nil {
		return nil, fmt.Errorf("Error creating new CRL: %s", err)
	}

	err = req.Storage.Put(&logical.StorageEntry{
		Key:   "crl",
		Value: crlBytes,
	})
	if err != nil {
		return nil, fmt.Errorf("Error storing CRL: %s", err)
	}

	return nil, nil
}
