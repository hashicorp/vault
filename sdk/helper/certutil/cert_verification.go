// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package certutil

import (
	"bytes"
	"crypto/x509"
	"fmt"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/helper/certutil/x509verify"
)

// VerifyCertificateChain validates the certificate chain in parsedBundle using
// the supplied options. When rootChainOnly is true, at least one root CA must
// be present in the chain; otherwise intermediates may serve as trust anchors
// when no root is available.
func VerifyCertificateChain(parsedBundle *ParsedCertBundle, options x509verify.VerifyOptions, rootChainOnly bool) error {
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

	rootCertPool := x509verify.NewCertPool()
	intermediateCertPool := x509verify.NewCertPool()

	for index, certBlock := range parsedBundle.CAChain {
		cert := certBlock.Certificate
		if cert == nil {
			var err error
			cert, err = x509.ParseCertificate(certBlock.Bytes)
			if err != nil {
				return fmt.Errorf("could not parse certificate number %v in chain: %w", index, err)
			}
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

	if !rootChainOnly && rootCertPool.Len() < 1 {
		// In this case, we don't have the root CA.  In some cases systems do trust an intermediate
		// directly, and this will work.  To accommodate those cases, we'll treat the intermediate
		// as the root.
		//
		// In other cases, such as in Common Criteria mode, we are required to return an error if
		// no root certificate is present.
		//
		// If there's no root, and we don't treat the intermediates as root certificates, we'd get
		// a "x509: certificate signed by unknown authority" error.
		rootCertPool, intermediateCertPool = intermediateCertPool, rootCertPool
	}

	options.Roots = rootCertPool
	options.Intermediates = intermediateCertPool
	options.CurrentTime = time.Now()

	cert := parsedBundle.Certificate
	if cert == nil {
		var err error
		cert, err = x509.ParseCertificate(parsedBundle.CertificateBytes)
		if err != nil {
			return fmt.Errorf("cannot parse certificate for validation: %w", err)
		}
	}

	_, err := x509verify.Verify(cert, options)
	return err
}

// VerifyCertificate validates the certificate chain in parsedBundle using the
// supplied options, treating intermediates as trust anchors when no root is
// present.
func VerifyCertificate(parsedBundle *ParsedCertBundle, options x509verify.VerifyOptions) error {
	return VerifyCertificateChain(parsedBundle, options, false)
}
