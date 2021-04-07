package diagnose

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/vault/internalshared/listenerutil"
	"github.com/hashicorp/vault/sdk/helper/tlsutil"
	"github.com/hashicorp/vault/vault"
)

const minVersionError = "'tls_min_version' value %q not supported, please specify one of [tls10,tls11,tls12,tls13]"
const maxVersionError = "'tls_max_version' value %q not supported, please specify one of [tls10,tls11,tls12,tls13]"

func ListenerChecks(listeners []listenerutil.Listener) error {
	for _, listener := range listeners {
		l := listener.Config
		var err error

		// Perform the TLS version check for listeners
		if l.TLSMinVersion == "" {
			l.TLSMinVersion = "tls12"
		}
		if l.TLSMaxVersion == "" {
			l.TLSMaxVersion = "tls13"
		}
		_, ok := tlsutil.TLSLookup[l.TLSMinVersion]
		if !ok {
			return fmt.Errorf(minVersionError, l.TLSMinVersion)
		}
		_, ok = tlsutil.TLSLookup[l.TLSMaxVersion]
		if !ok {
			return fmt.Errorf(maxVersionError, l.TLSMaxVersion)
		}

		// Perform checks on the TLS Cryptographic Information.
		if l.TLSRequireAndVerifyClientCert {
			err = TLSFileChecks(l.TLSCertFile, l.TLSKeyFile, l.TLSClientCAFile, true)
		} else {
			err = TLSFileChecks(l.TLSCertFile, l.TLSKeyFile, "", false)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// TLSFileChecks contains manual error checks against the TLS configuration
func TLSFileChecks(certFilePath, keyFilePath, rootFilePath string, checkRoot bool) error {

	// LoadX509KeyPair will check if the key/cert information can be loaded from files,
	// if they exist with keys and certs of the same algorithm type, if there
	// is an unknown algorithm type being used, and if the files have trailing
	// data.
	cert, err := tls.LoadX509KeyPair(certFilePath, keyFilePath)
	if err != nil {
		return err
	}

	// LoadX509KeyPair has a nil leaf certificate because it does not retain the
	// parsed form, so we have to manually create it ourselves.

	x509Cert, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return err
	}
	cert.Leaf = x509Cert

	if !checkRoot {
		// Check that certificate isn't expired and is of correct usage type
		if _, err = cert.Leaf.Verify(x509.VerifyOptions{
			KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		}); err != nil {
			return fmt.Errorf("failed to verify certificate: " + err.Error())
		}
	} else {
		if rootFilePath == "" {
			return fmt.Errorf("TLS Root CA file not specified, but TLS chain verification is required")
		}
		caPool := x509.NewCertPool()
		data, err := ioutil.ReadFile(rootFilePath)
		if err != nil {
			return fmt.Errorf("failed to read tls_client_ca_file: %s", err.Error())
		}

		if !caPool.AppendCertsFromPEM(data) {
			return fmt.Errorf("failed to parse CA certificate in tls_client_ca_file")
		}
		// Check that certificate isn't expired and is of correct usage type
		if _, err = cert.Leaf.Verify(x509.VerifyOptions{
			Roots:         caPool,
			Intermediates: x509.NewCertPool(),
			KeyUsages:     []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		}); err != nil {
			return fmt.Errorf("failed to verify certificate: " + err.Error())
		}
	}

	return nil
}

// ServerListenerActiveProbe attempts to use TLS information to set up a TLS server with each listener
// and generate a successful request through to the server.
// TODO
func ServerListenerActiveProbe(core *vault.Core) error {
	return fmt.Errorf("Method not implemented")
}
