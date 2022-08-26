package pluginutil

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/helper/certutil"
)

// generateCert is used internally to create certificates for the plugin
// client and server.
func generateCert() ([]byte, *ecdsa.PrivateKey, error) {
	key, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	host, err := uuid.GenerateUUID()
	if err != nil {
		return nil, nil, err
	}

	sn, err := certutil.GenerateSerialNumber()
	if err != nil {
		return nil, nil, err
	}

	template := &x509.Certificate{
		Subject: pkix.Name{
			CommonName: host,
		},
		DNSNames: []string{host},
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageClientAuth,
			x509.ExtKeyUsageServerAuth,
		},
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageKeyAgreement,
		SerialNumber: sn,
		NotBefore:    time.Now().Add(-30 * time.Second),
		NotAfter:     time.Now().Add(262980 * time.Hour),
		IsCA:         true,
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, template, template, key.Public(), key)
	if err != nil {
		return nil, nil, errwrap.Wrapf("unable to generate client certificate: {{err}}", err)
	}

	return certBytes, key, nil
}

// createClientTLSConfig creates a signed certificate and returns a configured
// TLS config.
func createClientTLSConfig(certBytes []byte, key *ecdsa.PrivateKey) (*tls.Config, error) {
	clientCert, err := x509.ParseCertificate(certBytes)
	if err != nil {
		return nil, errwrap.Wrapf("error parsing generated plugin certificate: {{err}}", err)
	}

	cert := tls.Certificate{
		Certificate: [][]byte{certBytes},
		PrivateKey:  key,
		Leaf:        clientCert,
	}

	clientCertPool := x509.NewCertPool()
	clientCertPool.AddCert(clientCert)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      clientCertPool,
		ClientCAs:    clientCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ServerName:   clientCert.Subject.CommonName,
		MinVersion:   tls.VersionTLS12,
	}

	return tlsConfig, nil
}

// wrapServerConfig is used to create a server certificate and private key, then
// wrap them in an unwrap token for later retrieval by the plugin.
func wrapServerConfig(ctx context.Context, sys RunnerUtil, certBytes []byte, key *ecdsa.PrivateKey) (string, error) {
	rawKey, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return "", err
	}

	wrapInfo, err := sys.ResponseWrapData(ctx, map[string]interface{}{
		"ServerCert": certBytes,
		"ServerKey":  rawKey,
	}, time.Second*60, true)
	if err != nil {
		return "", err
	}

	return wrapInfo.Token, nil
}
