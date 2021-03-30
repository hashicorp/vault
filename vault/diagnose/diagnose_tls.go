package diagnose

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"

	"github.com/hashicorp/vault/internalshared/listenerutil"
	"github.com/hashicorp/vault/vault"
)

// TLSConfigChecks contains manual error checks against the TLS configuration
func TLSConfigChecks(listeners []listenerutil.Listener) {
	for _, listener := range listeners {
		l := listener.Config

		// LoadX509KeyPair will check if the key/cert information can be loaded from files,
		// if they exist with keys and certs of the same algorithm type, if there
		// is an unknown algorithm type being used, and if the files have trailing
		// data.
		cert, err := tls.LoadX509KeyPair(l.TLSCertFile, l.TLSKeyFile)
		if err != nil {
			fmt.Printf("err is: %+v", err)
		}

		// QUESTION: Should we return certificate information in Diagnose?
		fmt.Printf("info found is: %+v", cert)

		// TODO: Check root as well via l.TLSClientCAFile

		// Check that certificate isn't expired and is of correct usage type
		cert.Leaf.Verify(x509.VerifyOptions{
			KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		})
	}
}

// ServerListenerActiveProbe attempts to use TLS information to set up a TLS server with each listener
// and generate a successful request through to the server.
// TODO
func ServerListenerActiveProbe(core *vault.Core) {}
