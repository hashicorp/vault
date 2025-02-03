// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package certutil

import (
	"bytes"
	"fmt"
	"time"

	ctx509 "github.com/google/certificate-transparency-go/x509"
	"github.com/hashicorp/errwrap"
)

func VerifyCertificate(parsedBundle *ParsedCertBundle, options ctx509.VerifyOptions) error {
	// If private key exists, check if it matches the public key of cert
	if parsedBundle.PrivateKey != nil && parsedBundle.Certificate != nil {
		equal, err := ComparePublicKeys(parsedBundle.Certificate.PublicKey, parsedBundle.PrivateKey.Public())
		if err != nil {
			return errwrap.Wrapf("could not compare public and private keys: {{err}}", err)
		}
		if !equal {
			return fmt.Errorf("public key of certificate does not match private key")
		}
	}

	rootCertPool := ctx509.NewCertPool()
	intermediateCertPool := ctx509.NewCertPool()

	for index, certificate := range parsedBundle.CAChain {
		cert, err := convertCertificate(certificate.Bytes)
		if err != nil {
			return fmt.Errorf("could not parse certificate number %v in chain: %w", index, err)
		}
		if index > 0 && !cert.IsCA {
			// Sometimes the leaf certificate is contained inside the bundle
			return fmt.Errorf("certificate %v is not a CA certificate", index)
		}
		if bytes.Equal(cert.RawIssuer, cert.RawSubject) {
			// Occasionally verify is called with a self-signed certificate that is not a CA;
			// We don't break that use case here
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

	options.Roots = rootCertPool
	options.Intermediates = intermediateCertPool
	options.CurrentTime = time.Now()

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
