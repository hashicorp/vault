// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package diagnose

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/hashicorp/go-secure-stdlib/tlsutil"
	"github.com/hashicorp/vault/internalshared/configutil"
)

const (
	minVersionError = "'tls_min_version' value %q not supported, please specify one of [tls10,tls11,tls12,tls13]"
	maxVersionError = "'tls_max_version' value %q not supported, please specify one of [tls10,tls11,tls12,tls13]"
)

// ListenerChecks diagnoses warnings and the first encountered error for the listener
// configuration stanzas.
func ListenerChecks(ctx context.Context, listeners []*configutil.Listener) ([]string, []error) {
	testName := "Check Listener TLS"
	ctx, span := StartSpan(ctx, testName)
	defer span.End()

	// These aggregated warnings and errors are returned purely for testing purposes.
	// The errors and warnings will report in this function itself.
	var listenerWarnings []string
	var listenerErrors []error

	for _, l := range listeners {
		listenerID := l.Address

		if l.TLSDisable {
			Warn(ctx, fmt.Sprintf("Listener at address %s: TLS is disabled in a listener config stanza.", listenerID))
			continue
		}
		if l.TLSDisableClientCerts {
			Warn(ctx, fmt.Sprintf("Listener at address %s: TLS for a listener is turned on without requiring client certificates.", listenerID))
		}
		status, warning := TLSMutualExclusionCertCheck(l)
		if status == 1 {
			Warn(ctx, warning)
		}

		// Perform the TLS version check for listeners.
		if l.TLSMinVersion == "" {
			l.TLSMinVersion = "tls12"
		}
		if l.TLSMaxVersion == "" {
			l.TLSMaxVersion = "tls13"
		}
		_, ok := tlsutil.TLSLookup[l.TLSMinVersion]
		if !ok {
			err := fmt.Errorf("Listener at address %s: %s.", listenerID, fmt.Sprintf(minVersionError, l.TLSMinVersion))
			listenerErrors = append(listenerErrors, err)
			Fail(ctx, err.Error())
		}
		_, ok = tlsutil.TLSLookup[l.TLSMaxVersion]
		if !ok {
			err := fmt.Errorf("Listener at address %s: %s.", listenerID, fmt.Sprintf(maxVersionError, l.TLSMaxVersion))
			listenerErrors = append(listenerErrors, err)
			Fail(ctx, err.Error())
		}

		// Perform checks on the TLS Cryptographic Information.
		warnings, err := TLSFileChecks(l.TLSCertFile, l.TLSKeyFile)
		listenerWarnings, listenerErrors = outputError(ctx, warnings, listenerWarnings, err, listenerErrors, listenerID)

		// Perform checks on the Client CA Cert
		warnings, err = TLSClientCAFileCheck(l)
		listenerWarnings, listenerErrors = outputError(ctx, warnings, listenerWarnings, err, listenerErrors, listenerID)
		// TODO: Use listenerutil.TLSConfig to warn on incorrect protocol specified
		// Alternatively, use tlsutil.SetupTLSConfig.
	}
	return listenerWarnings, listenerErrors
}

func outputError(ctx context.Context, newWarnings, listenerWarnings []string, newErr error, listenerErrors []error, listenerID string) ([]string, []error) {
	for _, warning := range newWarnings {
		warning = listenerID + ": " + warning
		listenerWarnings = append(listenerWarnings, warning)
		Warn(ctx, warning)
	}
	if newErr != nil {
		errMsg := listenerID + ": " + newErr.Error()
		listenerErrors = append(listenerErrors, errors.New(errMsg))
		Fail(ctx, errMsg)
	}
	return listenerWarnings, listenerErrors
}

// TLSFileChecks returns an error and warnings after checking TLS information
func TLSFileChecks(certpath, keypath string) ([]string, error) {
	warnings, err := TLSCertCheck(certpath)
	if err != nil {
		return warnings, err
	}

	// Utilize the native TLS Loading mechanism to ensure we have missed no errors
	_, err = tls.LoadX509KeyPair(certpath, keypath)
	return warnings, err
}

// TLSCertCheck returns an error and warning after checking TLS information on the given cert
func TLSCertCheck(certpath string) ([]string, error) {
	// Parse TLS Certs from the certpath
	leafCerts, interCerts, rootCerts, err := ParseTLSInformation(certpath)
	if err != nil {
		return nil, err
	}

	// Check for TLS Warnings
	warnings, err := TLSFileWarningChecks(leafCerts, interCerts, rootCerts)
	if err != nil {
		return warnings, err
	}

	// Check for TLS Errors
	if err = TLSErrorChecks(leafCerts, interCerts, rootCerts); err != nil {
		return warnings, err
	}
	return warnings, err
}

// ParseTLSInformation parses certficate information and returns it from a cert path.
func ParseTLSInformation(certFilePath string) ([]*x509.Certificate, []*x509.Certificate, []*x509.Certificate, error) {
	leafCerts := []*x509.Certificate{}
	interCerts := []*x509.Certificate{}
	rootCerts := []*x509.Certificate{}
	data, err := ioutil.ReadFile(certFilePath)
	if err != nil {
		return leafCerts, interCerts, rootCerts, fmt.Errorf("Failed to read certificate file: %w.", err)
	}

	certBlocks := []*pem.Block{}
	rst := []byte(data)
	for len(rst) != 0 {
		block, rest := pem.Decode(rst)
		if block == nil {
			return leafCerts, interCerts, rootCerts, fmt.Errorf("Could not decode certificate in certificate file.")
		}
		certBlocks = append(certBlocks, block)
		rst = rest
	}

	if len(certBlocks) == 0 {
		return leafCerts, interCerts, rootCerts, fmt.Errorf("No certificates found in certificate file.")
	}

	for _, certBlock := range certBlocks {
		cert, err := x509.ParseCertificate(certBlock.Bytes)
		if err != nil {
			return leafCerts, interCerts, rootCerts, fmt.Errorf("A PEM block does not parse to a certificate: %w.", err)
		}

		// Detect if the certificate is a root, leaf, or intermediate
		if cert.IsCA && bytes.Equal(cert.RawIssuer, cert.RawSubject) {
			// It's a root
			rootCerts = append(rootCerts, cert)
		} else if cert.IsCA {
			// It's not a root but it's a CA, so it's an inter
			interCerts = append(interCerts, cert)
		} else {
			// It's gotta be a leaf
			leafCerts = append(leafCerts, cert)
		}
	}

	return leafCerts, interCerts, rootCerts, nil
}

// TLSErrorChecks contains manual error checks against the TLS configuration
func TLSErrorChecks(leafCerts, interCerts, rootCerts []*x509.Certificate) error {
	// Make sure there's the proper number of leafCerts. If there are multiple, it's a bad pem file.
	if len(leafCerts) == 0 {
		return fmt.Errorf("No leaf certificates detected.")
	}

	// First, create root pools and interPools from the root and inter certs lists
	rootPool := x509.NewCertPool()
	interPool := x509.NewCertPool()

	for _, root := range rootCerts {
		rootPool.AddCert(root)
	}
	for _, inter := range interCerts {
		interPool.AddCert(inter)
	}

	var err error
	// Verify checks that certificate isn't expired, is of correct usage type, and has an appropriate
	// chain. We start with Root
	for _, root := range rootCerts {
		_, err = root.Verify(x509.VerifyOptions{
			Roots:     rootPool,
			KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
		})
		if err != nil {
			return fmt.Errorf("Failed to verify root certificate: %w.", err)
		}
	}

	// Verifying intermediate certs
	for _, inter := range interCerts {
		_, err = inter.Verify(x509.VerifyOptions{
			Roots:     rootPool,
			KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
		})
		if err != nil {
			return fmt.Errorf("Failed to verify intermediate certificate: %w.", err)
		}
	}

	rootSubjs := rootPool.Subjects()
	if len(rootSubjs) == 0 && len(leafCerts) > 0 {
		// this is a self signed server certificate, or the root is just not provided. In any
		// case, we need to bypass the root verification step by adding the leaf itself to the
		// root pool.
		rootPool.AddCert(leafCerts[0])
	}

	// Verifying leaf cert
	for _, leaf := range leafCerts {
		_, err = leaf.Verify(x509.VerifyOptions{
			Roots:         rootPool,
			Intermediates: interPool,
			KeyUsages:     []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
		})
		if err != nil {
			return fmt.Errorf("Failed to verify primary provided leaf certificate: %w.", err)
		}
	}

	return nil
}

// TLSFileWarningChecks returns warnings based on the leaf certificates, intermediate certificates,
// and root certificates provided.
func TLSFileWarningChecks(leafCerts, interCerts, rootCerts []*x509.Certificate) ([]string, error) {
	var warnings []string
	// add a warning for when there are more than one leaf certs
	if len(leafCerts) > 1 {
		warnings = append(warnings, fmt.Sprintf("More than one leaf certificate detected. Please ensure that there is one unique leaf certificate being supplied to Vault in the Vault server configuration file."))
	}

	for _, c := range leafCerts {
		if willExpire, timeToExpiry := NearExpiration(c); willExpire {
			warnings = append(warnings, fmt.Sprintf("Leaf certificate %d is expired or near expiry. Time to expire is: %s.", c.SerialNumber, timeToExpiry))
		}
	}
	for _, c := range interCerts {
		if willExpire, timeToExpiry := NearExpiration(c); willExpire {
			warnings = append(warnings, fmt.Sprintf("Intermediate certificate %d is expired or near expiry. Time to expire is: %s.", c.SerialNumber, timeToExpiry))
		}
	}
	for _, c := range rootCerts {
		if willExpire, timeToExpiry := NearExpiration(c); willExpire {
			warnings = append(warnings, fmt.Sprintf("Root certificate %d is expired or near expiry. Time to expire is: %s.", c.SerialNumber, timeToExpiry))
		}
	}

	return warnings, nil
}

// NearExpiration returns a true if a certificate will expire in a month
// and false otherwise, along with the duration until the certificate expires
// which can be a negative duration if the certificate has already expired.
func NearExpiration(c *x509.Certificate) (bool, time.Duration) {
	now := time.Now()
	timeToExpiry := c.NotAfter.Sub(now)

	oneMonthFromNow := now.Add(30 * 24 * time.Hour)
	isNearExpiration := oneMonthFromNow.After(c.NotAfter)

	return isNearExpiration, timeToExpiry
}

// TLSMutualExclusionCertCheck returns error if both TLSDisableClientCerts and TLSRequireAndVerifyClientCert are set
func TLSMutualExclusionCertCheck(l *configutil.Listener) (int, string) {
	if l.TLSDisableClientCerts {
		if l.TLSRequireAndVerifyClientCert {
			return 1, "The tls_disable_client_certs and tls_require_and_verify_client_cert fields in the listener stanza of the Vault server configuration are mutually exclusive fields. Please ensure they are not both set to true."
		}
	}
	return 0, ""
}

// TLSClientCAFileCheck Checks the validity of a client CA file
func TLSClientCAFileCheck(l *configutil.Listener) ([]string, error) {
	if l.TLSDisableClientCerts {
		return nil, nil
	} else if !l.TLSRequireAndVerifyClientCert {
		return nil, nil
	}
	return TLSCAFileCheck(l.TLSClientCAFile)
}

// TLSCAFileCheck checks the validity of a TLS CA file
func TLSCAFileCheck(CAFilePath string) ([]string, error) {
	var warningsSlc []string

	// Parse TLS Certs from the tls config
	leafCerts, interCerts, rootCerts, err := ParseTLSInformation(CAFilePath)
	if err != nil {
		return nil, err
	}

	if len(rootCerts) == 0 {
		return nil, fmt.Errorf("No root certificate found in CA certificate file.")
	}
	if len(rootCerts) > 1 {
		warningsSlc = append(warningsSlc, fmt.Sprintf("Found multiple root certificates in CA Certificate file instead of just one."))
	}

	// Checking for Self-Signed cert and return an explicit error about it.
	// Self-Signed certs are placed in the leafCerts slice when parsed.
	if len(leafCerts) > 0 && !leafCerts[0].IsCA && bytes.Equal(leafCerts[0].RawIssuer, leafCerts[0].RawSubject) {
		warningsSlc = append(warningsSlc, "Found a self-signed certificate in the CA certificate file.")
	}

	if len(interCerts) > 0 {
		warningsSlc = append(warningsSlc, "Found at least one intermediate certificate in the CA certificate file.")
	}

	if len(leafCerts) > 0 {
		warningsSlc = append(warningsSlc, "Found at least one leaf certificate in the CA certificate file.")
	}

	var warnings []string
	// Check for TLS Warnings
	warnings, err = TLSFileWarningChecks(leafCerts, interCerts, rootCerts)
	for i, warning := range warnings {
		warnings[i] = strings.ReplaceAll(warning, "leaf", "root")
	}
	warningsSlc = append(warningsSlc, warnings...)
	if err != nil {
		return warningsSlc, err
	}

	// Adding rootCerts to leafCert to perform verification in TLSErrorChecks
	leafCerts = append(leafCerts, rootCerts[0])

	// Check for TLS Errors
	if err = TLSErrorChecks(leafCerts, interCerts, rootCerts); err != nil {
		return warningsSlc, errors.New(strings.ReplaceAll(err.Error(), "leaf", "root"))
	}

	return warningsSlc, err
}
