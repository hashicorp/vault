package diagnose

import (
	"crypto/tls"
	"crypto/x509"

	"github.com/hashicorp/vault/internalshared/listenerutil"
	"github.com/hashicorp/vault/vault"
)

func ListenerChecks(listeners []listenerutil.Listener) error {
	for _, listener := range listeners {
		l := listener.Config
		err := TLSFileChecks(l.TLSCertFile, l.TLSKeyFile)
		if err != nil {
			return err
		}
	}
	return nil
}

// TLSChecks contains manual error checks against the TLS configuration
func TLSFileChecks(certFilePath, keyFilePath string) error {

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

	// TODO: Check root as well via l.TLSClientCAFile

	// Check that certificate isn't expired and is of correct usage type
	cert.Leaf.Verify(x509.VerifyOptions{
		KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	})
	return nil
}

// ServerListenerActiveProbe attempts to use TLS information to set up a TLS server with each listener
// and generate a successful request through to the server.
// TODO
func ServerListenerActiveProbe(core *vault.Core) {}
