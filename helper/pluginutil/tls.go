package pluginutil

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	mathrand "math/rand"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/SermoDigital/jose/jws"
	"github.com/hashicorp/errwrap"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/api"
)

var (
	// PluginUnwrapTokenEnv is the ENV name used to pass unwrap tokens to the
	// plugin.
	PluginUnwrapTokenEnv = "VAULT_UNWRAP_TOKEN"
)

type Wrapper interface {
	ResponseWrapData(data map[string]interface{}, ttl time.Duration, jwt bool) (string, error)
}

// GenerateCACert returns a CA cert used to later sign the certificates for the
// plugin client and server.
func GenerateCACert() ([]byte, *x509.Certificate, *ecdsa.PrivateKey, error) {
	key, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		return nil, nil, nil, err
	}

	host, err := uuid.GenerateUUID()
	if err != nil {
		return nil, nil, nil, err
	}
	host = "localhost"
	template := &x509.Certificate{
		Subject: pkix.Name{
			CommonName: host,
		},
		DNSNames: []string{host},
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
			x509.ExtKeyUsageClientAuth,
		},
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageKeyAgreement | x509.KeyUsageCertSign,
		SerialNumber: big.NewInt(mathrand.Int63()),
		NotBefore:    time.Now().Add(-30 * time.Second),
		// 30 years of single-active uptime ought to be enough for anybody
		NotAfter:              time.Now().Add(262980 * time.Hour),
		BasicConstraintsValid: true,
		IsCA: true,
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, template, template, key.Public(), key)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("unable to generate replicated cluster certificate: %v", err)
	}

	caCert, err := x509.ParseCertificate(certBytes)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error parsing generated replication certificate: %v", err)
	}

	return certBytes, caCert, key, nil
}

// generateSignedCert is used internally to create certificates for the plugin
// client and server. These certs are signed by the given CA Cert and Key.
func generateSignedCert(CACert *x509.Certificate, CAKey *ecdsa.PrivateKey) ([]byte, *x509.Certificate, *ecdsa.PrivateKey, error) {
	host, err := uuid.GenerateUUID()
	if err != nil {
		return nil, nil, nil, err
	}
	host = "localhost"
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
		SerialNumber: big.NewInt(mathrand.Int63()),
		NotBefore:    time.Now().Add(-30 * time.Second),
		NotAfter:     time.Now().Add(262980 * time.Hour),
	}

	clientKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		return nil, nil, nil, errwrap.Wrapf("error generating client key: {{err}}", err)
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, template, CACert, clientKey.Public(), CAKey)
	if err != nil {
		return nil, nil, nil, errwrap.Wrapf("unable to generate client certificate: {{err}}", err)
	}

	clientCert, err := x509.ParseCertificate(certBytes)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error parsing generated replication certificate: %v", err)
	}

	return certBytes, clientCert, clientKey, nil
}

// CreateClientTLSConfig creates a signed certificate and returns a configured
// TLS config.
func CreateClientTLSConfig(CACert *x509.Certificate, CAKey *ecdsa.PrivateKey) (*tls.Config, error) {
	clientCertBytes, clientCert, clientKey, err := generateSignedCert(CACert, CAKey)
	if err != nil {
		return nil, err
	}

	cert := tls.Certificate{
		Certificate: [][]byte{clientCertBytes},
		PrivateKey:  clientKey,
		Leaf:        clientCert,
	}

	clientCertPool := x509.NewCertPool()
	clientCertPool.AddCert(CACert)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      clientCertPool,
		ClientCAs:    clientCertPool,
		ServerName:   CACert.Subject.CommonName,
		MinVersion:   tls.VersionTLS12,
	}

	tlsConfig.BuildNameToCertificate()

	return tlsConfig, nil
}

// WrapServerConfig is used to create a server certificate and private key, then
// wrap them in an unwrap token for later retrieval by the plugin.
func WrapServerConfig(sys Wrapper, CACertBytes []byte, CACert *x509.Certificate, CAKey *ecdsa.PrivateKey) (string, error) {
	serverCertBytes, _, serverKey, err := generateSignedCert(CACert, CAKey)
	if err != nil {
		return "", err
	}
	rawKey, err := x509.MarshalECPrivateKey(serverKey)
	if err != nil {
		return "", err
	}

	wrapToken, err := sys.ResponseWrapData(map[string]interface{}{
		"CACert":     CACertBytes,
		"ServerCert": serverCertBytes,
		"ServerKey":  rawKey,
	}, time.Second*10, true)

	return wrapToken, err
}

// VaultPluginTLSProvider is run inside a plugin and retrives the response
// wrapped TLS certificate from vault. It returns a configured TLS Config.
func VaultPluginTLSProvider() (*tls.Config, error) {
	unwrapToken := os.Getenv(PluginUnwrapTokenEnv)

	// Ensure unwrap token is a JWT
	if strings.Count(unwrapToken, ".") != 2 {
		return nil, errors.New("Could not parse unwraptoken")
	}

	// Parse the JWT and retrieve the vault address
	wt, err := jws.ParseJWT([]byte(unwrapToken))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error decoding token: %s", err))
	}
	if wt == nil {
		return nil, errors.New("nil decoded token")
	}

	addrRaw := wt.Claims().Get("addr")
	if addrRaw == nil {
		return nil, errors.New("decoded token does not contain primary cluster address")
	}
	vaultAddr, ok := addrRaw.(string)
	if !ok {
		return nil, errors.New("decoded token's address not valid")
	}
	if vaultAddr == "" {
		return nil, errors.New(`no address for the vault found`)
	}

	// Sanity check the value
	if _, err := url.Parse(vaultAddr); err != nil {
		return nil, errors.New(fmt.Sprintf("error parsing the vault address: %s", err))
	}

	// Unwrap the token
	clientConf := api.DefaultConfig()
	clientConf.Address = vaultAddr
	client, err := api.NewClient(clientConf)
	if err != nil {
		return nil, errwrap.Wrapf("error during token unwrap request: {{err}}", err)
	}

	secret, err := client.Logical().Unwrap(unwrapToken)
	if err != nil {
		return nil, errwrap.Wrapf("error during token unwrap request: {{err}}", err)
	}

	// Retrieve and parse the CA Certificate
	CABytesRaw, ok := secret.Data["CACert"].(string)
	if !ok {
		return nil, errors.New("error unmarshalling CA certificate")
	}

	CABytes, err := base64.StdEncoding.DecodeString(CABytesRaw)
	if err != nil {
		return nil, fmt.Errorf("error parsing certificate: %v", err)
	}

	CACert, err := x509.ParseCertificate(CABytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing certificate: %v", err)
	}

	// Retrieve and parse the server's certificate
	serverCertBytesRaw, ok := secret.Data["ServerCert"].(string)
	if !ok {
		return nil, errors.New("error unmarshalling certificate")
	}

	serverCertBytes, err := base64.StdEncoding.DecodeString(serverCertBytesRaw)
	if err != nil {
		return nil, fmt.Errorf("error parsing certificate: %v", err)
	}

	serverCert, err := x509.ParseCertificate(serverCertBytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing certificate: %v", err)
	}

	// Retrieve and parse the server's private key
	serverKeyB64, ok := secret.Data["ServerKey"].(string)
	if !ok {
		return nil, errors.New("error unmarshalling certificate")
	}

	serverKeyRaw, err := base64.StdEncoding.DecodeString(serverKeyB64)
	if err != nil {
		return nil, fmt.Errorf("error parsing certificate: %v", err)
	}

	serverKey, err := x509.ParseECPrivateKey(serverKeyRaw)
	if err != nil {
		return nil, fmt.Errorf("error parsing certificate: %v", err)
	}

	// Add CA cert to the cert pool
	caCertPool := x509.NewCertPool()
	caCertPool.AddCert(CACert)

	// Build a certificate object out of the server's cert and private key.
	cert := tls.Certificate{
		Certificate: [][]byte{serverCertBytes},
		PrivateKey:  serverKey,
		Leaf:        serverCert,
	}

	// Setup TLS config
	tlsConfig := &tls.Config{
		ClientCAs:  caCertPool,
		RootCAs:    caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
		// TLS 1.2 minimum
		MinVersion:   tls.VersionTLS12,
		Certificates: []tls.Certificate{cert},
	}
	tlsConfig.BuildNameToCertificate()

	return tlsConfig, nil
}
