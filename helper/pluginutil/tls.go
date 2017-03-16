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

func GenerateX509Cert() ([]byte, *x509.Certificate, *ecdsa.PrivateKey, error) {
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

func GenerateClientCert(CACert *x509.Certificate, CAKey *ecdsa.PrivateKey) ([]byte, *x509.Certificate, []byte, error) {
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

	keyBytes, err := x509.MarshalECPrivateKey(clientKey)
	if err != nil {
		return nil, nil, nil, err
	}

	return certBytes, clientCert, keyBytes, nil
}

// VaultPluginTLSProvider is run inside a plugin and retrives the response
// wrapped TLS certificate from vault. It returns a configured tlsConfig.
func VaultPluginTLSProvider() (*tls.Config, error) {
	unwrapToken := os.Getenv("VAULT_WRAP_TOKEN")
	if strings.Count(unwrapToken, ".") != 2 {
		return nil, errors.New("Could not parse unwraptoken")
	}

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

	CABytesRaw, ok := secret.Data["CACert"].(string)
	if !ok {
		return nil, errors.New("error unmarshalling certificate")
	}

	CABytes, err := base64.StdEncoding.DecodeString(CABytesRaw)
	if err != nil {
		return nil, fmt.Errorf("error parsing certificate: %v", err)
	}

	CACert, err := x509.ParseCertificate(CABytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing certificate: %v", err)
	}

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

	serverKeyRaw, ok := secret.Data["ServerKey"].(string)
	if !ok {
		return nil, errors.New("error unmarshalling certificate")
	}

	serverKey, err := base64.StdEncoding.DecodeString(serverKeyRaw)
	if err != nil {
		return nil, fmt.Errorf("error parsing certificate: %v", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AddCert(CACert)

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
