// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package issuing

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	ctx509 "github.com/google/certificate-transparency-go/x509"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/logical"
)

// disableVerifyCertificateEnvVar is an environment variable that can be used to disable the
// verification done when issuing or signing certificates that was added by VAULT-22013. It
// is meant as a scape hatch to avoid breaking deployments that the new verification would
// break.
const disableVerifyCertificateEnvVar = "VAULT_DISABLE_PKI_CONSTRAINTS_VERIFICATION"

func isCertificateVerificationDisabled() (bool, error) {
	disableRaw, ok := os.LookupEnv(disableVerifyCertificateEnvVar)
	if !ok {
		return false, nil
	}

	disable, err := strconv.ParseBool(disableRaw)
	if err != nil {
		return false, fmt.Errorf("failed parsing environment variable %s: %w", disableVerifyCertificateEnvVar, err)
	}

	return disable, nil
}

func VerifyCertificate(ctx context.Context, storage logical.Storage, issuerId IssuerID, parsedBundle *certutil.ParsedCertBundle) error {
	if verificationDisabled, err := isCertificateVerificationDisabled(); err != nil {
		return err
	} else if verificationDisabled {
		return nil
	}

	rootCertPool := ctx509.NewCertPool()
	intermediateCertPool := ctx509.NewCertPool()

	for _, certificate := range parsedBundle.CAChain {
		cert, err := convertCertificate(certificate.Bytes)
		if err != nil {
			return err
		}
		if bytes.Equal(cert.RawIssuer, cert.RawSubject) {
			rootCertPool.AddCert(cert)
		} else {
			intermediateCertPool.AddCert(cert)
		}
	}
	if len(rootCertPool.Subjects()) < 1 {
		// Alright, this is weird, since we don't have the root CA, we'll treat the intermediate as
		// the root, otherwise we'll get a "x509: certificate signed by unknown authority" error.
		rootCertPool, intermediateCertPool = intermediateCertPool, rootCertPool
	}

	// Note that we use github.com/google/certificate-transparency-go/x509 to perform certificate verification,
	// since that library provides options to disable checks that the standard library does not.

	options := ctx509.VerifyOptions{
		Roots:                          rootCertPool,
		Intermediates:                  intermediateCertPool,
		CurrentTime:                    time.Time{},
		KeyUsages:                      nil,
		MaxConstraintComparisions:      0, // Use the library's 'sensible default'
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
