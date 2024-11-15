// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package issuing

import (
	"context"
	"fmt"
	ctx509 "github.com/google/certificate-transparency-go/x509"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/logical"
	"time"
)

func VerifyCertificate(ctx context.Context, storage logical.Storage, issuerId IssuerID, parsedBundle *certutil.ParsedCertBundle, data *framework.FieldData) error {
	certChainPool := ctx509.NewCertPool()
	for _, certificate := range parsedBundle.CAChain {
		cert, err := convertCertificate(certificate.Bytes)
		if err != nil {
			return err
		}
		certChainPool.AddCert(cert)
	}

	// Validation Code, assuming we need to validate the entire chain of constraints

	// Note that we use github.com/google/certificate-transparency-go/x509 to perform certificate verification,
	// since that library provides options to disable checks that the standard library does not.

	options := ctx509.VerifyOptions{
		Intermediates:                  nil, // We aren't verifying the chain here, this would do more work
		Roots:                          certChainPool,
		CurrentTime:                    time.Time{},
		KeyUsages:                      nil,
		MaxConstraintComparisions:      0, // This means infinite
		DisableTimeChecks:              true,
		DisableEKUChecks:               true,
		DisableCriticalExtensionChecks: false,
		DisableNameChecks:              false,
		DisablePathLenChecks:           false,
		DisableNameConstraintChecks:    false,
	}

	if err := entSetCertVerifyOptions(ctx, storage, issuerId, &options); err != nil {
		return err
	}

	certificate, err := convertCertificate(parsedBundle.CertificateBytes)
	if err != nil {
		return err
	}

	_, err = certificate.Verify(options)
	return err
}

func convertCertificate(certBytes []byte) (*ctx509.Certificate, error) {
	ret, err := ctx509.ParseCertificate(certBytes)
	if err != nil {
		return nil, fmt.Errorf("cannot convert certificate for validation: %w", err)
	}
	return ret, nil
}
