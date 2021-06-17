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

	"github.com/hashicorp/vault/internalshared/listenerutil"
	"github.com/hashicorp/vault/sdk/helper/tlsutil"
)

const minVersionError = "'tls_min_version' value %q not supported, please specify one of [tls10,tls11,tls12,tls13]"
const maxVersionError = "'tls_max_version' value %q not supported, please specify one of [tls10,tls11,tls12,tls13]"

// ListenerChecks diagnoses warnings and the first encountered error for the listener
// configuration stanzas.
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
		listenerWarnings, listenerErrors = outputErrors(ctx, warnings, listenerWarnings, err, listenerErrors, listenerID)

		// TODO: Use listenerutil.TLSConfig to warn on incorrect protocol specified
		// Alternatively, use tlsutil.SetupTLSConfig.
	}
	return listenerWarnings, listenerErrors
}

func outputErrors(ctx context.Context, newWarnings, listenerWarnings []string, newErr error, listenerErrors []error, listenerID string) ([]string, []error) {
	for _, warning := range newWarnings {
		warning = listenerID + ": " + warning
		if !containsString(listenerWarnings, warning) {
			listenerWarnings = append(listenerWarnings, warning)
			Warn(ctx, warning)
		}
	}
	if newErr != nil {
		errMsg := listenerID + ": " + newErr.Error()
		listenerErrors = append(listenerErrors, fmt.Errorf(errMsg))
		Error(ctx, fmt.Errorf(errMsg))
	}
	return listenerWarnings, listenerErrors
}

func containsString(messages []string, lookup_msg string) bool {
	//TODO: change messages to a map for performance reasons
	for _, msg := range messages {
		if msg == lookup_msg {
			return true
		}
	}
	return false

}

// GetCertsSubjects returns a list of the DER-encoded subjects of
// all of the certificates in the slice.
func GetCertsSubjects(certs []*x509.Certificate) [][]byte {
	res := make([][]byte, len(certs))
	for i, lc := range certs {
		res[i] = lc.RawSubject
	}
	return res
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
		return leafCerts, interCerts, rootCerts, fmt.Errorf("failed to read certificate file: %w", err)
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

	// Verify checks that certificate isn't expired, is of correct usage type, and has an appropriate
	// chain. We start with Root
	for _, root := range rootCerts {
		_, err := root.Verify(x509.VerifyOptions{
			Roots:     rootPool,
			KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
		})
		if err != nil {
			return fmt.Errorf("failed to verify root certificate: %w", err)
		}
	}

	// Verifying intermediate certs
	for _, inter := range interCerts {
		_, err := inter.Verify(x509.VerifyOptions{
			Roots:     rootPool,
			KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
		})
		if err != nil {
			return fmt.Errorf("failed to verify intermediate certificate: %w", err)
		}
	}

	rootSubjects := GetCertsSubjects(rootCerts)
	if len(rootSubjects) == 0 {
		// This is either a self signed cert or the root is just not provided. In any
		// case, we need to bypass the root verification step by adding the leaf itself to the
		// root pool.
		rootCerts = append(rootCerts, leafCerts[0])
	}

	// Verifying leaf cert
	for _, leaf := range leafCerts {
		_, err := leaf.Verify(x509.VerifyOptions{
			Roots:         rootPool,
			Intermediates: interPool,
			KeyUsages:     []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
		})
		if err != nil {
			return fmt.Errorf("failed to verify primary provided leaf certificate: %w", err)
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
		warnings = append(warnings, fmt.Sprintf("More than one leaf certificate detected. Please ensure that there is one unique leaf certificate being supplied to vault in the vault server config file."))
	}

	for _, c := range leafCerts {
		if willExpire, timeToExpiry := NearExpiration(c); willExpire {
			warnings = append(warnings, fmt.Sprintf("leaf certificate %d is expired or near expiry. Time to expire is: %s", c.SerialNumber, timeToExpiry))
		}
	}
	for _, c := range interCerts {
		if willExpire, timeToExpiry := NearExpiration(c); willExpire {
			warnings = append(warnings, fmt.Sprintf("intermediate certificate %d is expired or near expiry. Time to expire is: %s", c.SerialNumber, timeToExpiry))
		}
	}
	for _, c := range rootCerts {
		if willExpire, timeToExpiry := NearExpiration(c); willExpire {
			warnings = append(warnings, fmt.Sprintf("root certificate %d is expired or near expiry. Time to expire is: %s", c.SerialNumber, timeToExpiry))
		}
	}

	return warnings, nil
}

// NearExpiration returns a true if a certficate will expire in a month and false otherwise
func NearExpiration(c *x509.Certificate) (bool, time.Duration) {
	oneMonthFromNow := time.Now().Add(30 * 24 * time.Hour)
	var timeToExpiry time.Duration
	if oneMonthFromNow.After(c.NotAfter) {
		timeToExpiry := oneMonthFromNow.Sub(c.NotAfter)
		return true, timeToExpiry
	}
	return false, timeToExpiry
}
