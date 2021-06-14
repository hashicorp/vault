package diagnose

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/internalshared/listenerutil"
	"github.com/hashicorp/vault/sdk/helper/tlsutil"
)

const minVersionError = "'tls_min_version' value %q not supported, please specify one of [tls10,tls11,tls12,tls13]"
const maxVersionError = "'tls_max_version' value %q not supported, please specify one of [tls10,tls11,tls12,tls13]"

func ListenerChecks(ctx context.Context, listeners []listenerutil.Listener) ([]string, []error) {

	// These aggregated warnings and errors are returned purely for testing purposes.
	// The errors and warnings will report in this function itself.
	var listenerWarnings []string
	var listenerErrors []error

	for _, listener := range listeners {
		l := listener.Config
		listenerID := l.Address

		// Perform the TLS version check for listeners.
		if l.TLSMinVersion == "" {
			l.TLSMinVersion = "tls12"
		}
		if l.TLSMaxVersion == "" {
			l.TLSMaxVersion = "tls13"
		}
		_, ok := tlsutil.TLSLookup[l.TLSMinVersion]
		if !ok {
			err := fmt.Errorf("listener at address: %s has error %s: ", listenerID, fmt.Sprintf(minVersionError, l.TLSMinVersion))
			listenerErrors = append(listenerErrors, err)
			Error(ctx, err)
		}
		_, ok = tlsutil.TLSLookup[l.TLSMaxVersion]
		if !ok {
			err := fmt.Errorf("listener at address: %s has error %s: ", listenerID, fmt.Sprintf(maxVersionError, l.TLSMaxVersion))
			listenerErrors = append(listenerErrors, err)
			Error(ctx, err)
		}

		// Perform checks on the TLS Cryptographic Information.
		warnings, err := TLSFileChecks(l.TLSCertFile, l.TLSKeyFile)
		for _, warning := range warnings {
			warning = listenerID + ": " + warning
			listenerWarnings = append(listenerWarnings, warning)
			Warn(ctx, warning)
		}
		if err != nil {
			errMsg := listenerID + ": " + err.Error()
			listenerErrors = append(listenerErrors, fmt.Errorf(errMsg))
			Error(ctx, fmt.Errorf(errMsg))
		}

		// TODO: Use listenerutil.TLSConfig to warn on incorrect protocol specified
		// Alternatively, use tlsutil.SetupTLSConfig.
	}
	return listenerWarnings, listenerErrors
}

// TLSFileChecks returns an error and warnings after checking TLS information
func TLSFileChecks(certpath, keypath string) ([]string, error) {
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

	// Utilize the native TLS Loading mechanism to ensure we have missed no errors
	_, err = tls.LoadX509KeyPair(certpath, keypath)
	return warnings, err
}

// ParseTLSInformation parses certficate information and returns it from a cert path.
func ParseTLSInformation(certFilePath string) ([]*x509.Certificate, []*x509.Certificate, []*x509.Certificate, error) {
	leafCerts := []*x509.Certificate{}
	interCerts := []*x509.Certificate{}
	rootCerts := []*x509.Certificate{}
	data, err := ioutil.ReadFile(certFilePath)
	if err != nil {
		return leafCerts, interCerts, rootCerts, fmt.Errorf("failed to read tls_client_ca_file: %w", err)
	}

	certBlocks := []*pem.Block{}
	rst := []byte(data)
	for len(rst) != 0 {
		block, rest := pem.Decode(rst)
		if block == nil {
			return leafCerts, interCerts, rootCerts, fmt.Errorf("could not decode cert")
		}
		certBlocks = append(certBlocks, block)
		rst = rest
	}

	if len(certBlocks) == 0 {
		return leafCerts, interCerts, rootCerts, fmt.Errorf("no certificates found in cert file")
	}

	for _, certBlock := range certBlocks {
		cert, err := x509.ParseCertificate(certBlock.Bytes)
		if err != nil {
			return leafCerts, interCerts, rootCerts, fmt.Errorf("A pem block does not parse to a certificate: %w", err)
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
	// First, create root pools and interPools from the root and inter certs lists

	rootPool := x509.NewCertPool()
	interPool := x509.NewCertPool()

	for _, root := range rootCerts {
		rootPool.AddCert(root)
	}
	for _, inter := range interCerts {
		interPool.AddCert(inter)
	}

	// Make sure there's only one leaf. If there are multiple, it's a bad pem file.
	if len(leafCerts) != 1 {
		return fmt.Errorf("Number of leaf certificates detected is not one. Instead, it is: %d", len(leafCerts))
	}

	rootSubjs := rootPool.Subjects()
	if len(rootSubjs) == 0 {
		// this is a self signed server certificate, or the root is just not provided. In any
		// case, we need to bypass the root verification step by adding the leaf itself to the
		// root pool.
		rootPool.AddCert(leafCerts[0])
	}

	// Verify checks that certificate isn't expired, is of correct usage type, and has an appropriate
	// chain.
	_, err := leafCerts[0].Verify(x509.VerifyOptions{
		Roots:         rootPool,
		Intermediates: interPool,
		KeyUsages:     []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	})

	if err != nil {
		return fmt.Errorf("failed to verify primary provided leaf certificate: %w", err)
	}

	return nil
}

func TLSFileWarningChecks(leafcerts, interCerts, rootCerts []*x509.Certificate) ([]string, error) {
	var warnings []string
	for _, c := range leafcerts {
		if NearExpiration(c) {
			warnings = append(warnings, fmt.Sprintf("certificate %d is expired or near expiry", c.SerialNumber))
		}
	}
	for _, c := range interCerts {
		if NearExpiration(c) {
			warnings = append(warnings, fmt.Sprintf("certificate %d is expired or near expiry", c.SerialNumber))
		}
	}
	for _, c := range rootCerts {
		if NearExpiration(c) {
			warnings = append(warnings, fmt.Sprintf("certificate %d is expired or near expiry", c.SerialNumber))
		}
	}

	return warnings, nil
}

// NearExpiration returns a true if a certficate will expire in a week and false otherwise
func NearExpiration(c *x509.Certificate) bool {
	oneWeekFromNow := time.Now().Add(7 * 24 * time.Hour)
	if oneWeekFromNow.After(c.NotAfter) {
		return true
	}
	return false
}

// TLSMutualExclusionCertCheck returns error if both TLSDisableClientCerts and TLSRequireAndVerifyClientCert are set
func TLSMutualExclusionCertCheck(l *configutil.Listener) error {

	if l.TLSDisableClientCerts {
		if l.TLSRequireAndVerifyClientCert {
			return fmt.Errorf("'tls_disable_client_certs' and 'tls_require_and_verify_client_cert' are mutually exclusive")
		}
	}
	return nil
}
