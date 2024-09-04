// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package revocation

import (
	"bytes"
	"context"
	"crypto/x509"
	"fmt"
	"time"

	"github.com/hashicorp/vault/builtin/logical/pki/issuing"
	"github.com/hashicorp/vault/builtin/logical/pki/pki_backend"
	"github.com/hashicorp/vault/sdk/helper/errutil"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	RevokedPath = "revoked/"
)

type RevokerFactory interface {
	GetRevoker(context.Context, logical.Storage) (Revoker, error)
}

type RevokeCertInfo struct {
	RevocationTime time.Time
	Warnings       []string
}

type Revoker interface {
	RevokeCert(cert *x509.Certificate) (RevokeCertInfo, error)
	RevokeCertBySerial(serial string) (RevokeCertInfo, error)
}

type RevocationInfo struct {
	CertificateBytes  []byte           `json:"certificate_bytes"`
	RevocationTime    int64            `json:"revocation_time"`
	RevocationTimeUTC time.Time        `json:"revocation_time_utc"`
	CertificateIssuer issuing.IssuerID `json:"issuer_id"`
}

func (ri *RevocationInfo) AssociateRevokedCertWithIsssuer(revokedCert *x509.Certificate, issuerIDCertMap map[issuing.IssuerID]*x509.Certificate) bool {
	for issuerId, issuerCert := range issuerIDCertMap {
		if bytes.Equal(revokedCert.RawIssuer, issuerCert.RawSubject) {
			if err := revokedCert.CheckSignatureFrom(issuerCert); err == nil {
				// Valid mapping. Add it to the specified entry.
				ri.CertificateIssuer = issuerId
				return true
			}
		}
	}

	return false
}

// FetchIssuerMapForRevocationChecking fetches a map of IssuerID->parsed cert for revocation
// usage. Unlike other paths, this needs to handle the legacy bundle
// more gracefully than rejecting it outright.
func FetchIssuerMapForRevocationChecking(sc pki_backend.StorageContext) (map[issuing.IssuerID]*x509.Certificate, error) {
	var err error
	var issuers []issuing.IssuerID

	if !sc.UseLegacyBundleCaStorage() {
		issuers, err = issuing.ListIssuers(sc.GetContext(), sc.GetStorage())
		if err != nil {
			return nil, fmt.Errorf("could not fetch issuers list: %w", err)
		}
	} else {
		// Hack: this isn't a real IssuerID, but it works for fetchCAInfo
		// since it resolves the reference.
		issuers = []issuing.IssuerID{issuing.LegacyBundleShimID}
	}

	issuerIDCertMap := make(map[issuing.IssuerID]*x509.Certificate, len(issuers))
	for _, issuer := range issuers {
		_, bundle, caErr := issuing.FetchCertBundleByIssuerId(sc.GetContext(), sc.GetStorage(), issuer, false)
		if caErr != nil {
			return nil, fmt.Errorf("error fetching CA certificate for issuer id %v: %w", issuer, caErr)
		}

		if bundle == nil {
			return nil, fmt.Errorf("faulty reference: %v - CA info not found", issuer)
		}

		parsedBundle, err := issuing.ParseCABundle(sc.GetContext(), sc.GetPkiManagedView(), bundle)
		if err != nil {
			return nil, errutil.InternalError{Err: err.Error()}
		}

		if parsedBundle.Certificate == nil {
			return nil, errutil.InternalError{Err: "stored CA information not able to be parsed"}
		}

		issuerIDCertMap[issuer] = parsedBundle.Certificate
	}

	return issuerIDCertMap, nil
}
