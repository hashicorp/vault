package diagnose

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
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

		var err error
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
	data, err := ioutil.ReadFile(certFilePath)
	if err != nil {
		return fmt.Errorf("failed to read tls_client_ca_file: %w", err)
	}

	certBlocks := []*pem.Block{}
	leafCerts := []*x509.Certificate{}
	rootPool := x509.NewCertPool()
	interPool := x509.NewCertPool()
	rst := []byte(data)
	for len(rst) != 0 {
		block, rest := pem.Decode(rst)
		if block == nil {
			return fmt.Errorf("could not decode cert")
		}
		certBlocks = append(certBlocks, block)
		rst = rest
	}

	if len(certBlocks) == 0 {
		return fmt.Errorf("no certificates found in cert file")
	}

	for _, certBlock := range certBlocks {
		cert, err := x509.ParseCertificate(certBlock.Bytes)
		if err != nil {
			return fmt.Errorf("A pem block does not parse to a certificate: %w", err)
		}

		// Detect if the certificate is a root, leaf, or intermediate
		if cert.IsCA && bytes.Equal(cert.RawIssuer, cert.RawSubject) {
			// It's a root
			rootPool.AddCert(cert)
		} else if cert.IsCA {
			// It's not a root but it's a CA, so it's an inter
			interPool.AddCert(cert)
		} else {
			// It's gotta be a leaf
			leafCerts = append(leafCerts, cert)
		}
	}

	// Make sure there's only one leaf. If there are multiple, it's a bad pem file.
	if len(leafCerts) != 1 {
		return fmt.Errorf("Number of leaf certificates detected is not one. Instead, it is: %d", len(leafCerts))
	}

	// Check that certificate isn't expired, is of correct usage type, and has an appropriate
	// chain.
	if chains, err := leafCerts[0].Verify(x509.VerifyOptions{
		Roots:         rootPool,
		Intermediates: interPool,
		KeyUsages:     []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}); err != nil {
		return fmt.Errorf("failed to verify certificate: %w", err)
	} else {
		fmt.Printf("%+v", chains)
	}

	// LoadX509KeyPair has a nil leaf certificate because it does not retain the
	// parsed form, so we have to manually create it ourselves.

	// cert, err := tls.LoadX509KeyPair(certFilePath, keyFilePath)
	// if err != nil {
	// 	return err
	// }
	// x509Cert, err := x509.ParseCertificate(cert.Certificate[0])
	// if err != nil {
	// 	return err
	// }
	// cert.Leaf = x509Cert

	return nil
}

// ServerListenerActiveProbe attempts to use TLS information to set up a TLS server with each listener
// and generate a successful request through to the server.
// TODO
func ServerListenerActiveProbe(core *vault.Core) error {
	return fmt.Errorf("Method not implemented")
}
