package pluginutil

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"errors"
	"net/url"
	"os"
	"time"

	"github.com/SermoDigital/jose/jws"
	"github.com/hashicorp/errwrap"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/certutil"
)

var (
	// PluginUnwrapTokenEnv is the ENV name used to pass unwrap tokens to the
	// plugin.
	PluginUnwrapTokenEnv = "VAULT_UNWRAP_TOKEN"

	// PluginCACertPEMEnv is an ENV name used for holding a CA PEM-encoded
	// string. Used for testing.
	PluginCACertPEMEnv = "VAULT_TESTING_PLUGIN_CA_PEM"
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

	tlsConfig.BuildNameToCertificate()

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

// VaultPluginTLSProvider is run inside a plugin and retrieves the response
// wrapped TLS certificate from vault. It returns a configured TLS Config.
func VaultPluginTLSProvider(apiTLSConfig *api.TLSConfig) func() (*tls.Config, error) {
	if os.Getenv(PluginMetadataModeEnv) == "true" {
		return nil
	}

	return func() (*tls.Config, error) {
		unwrapToken := os.Getenv(PluginUnwrapTokenEnv)

		// Parse the JWT and retrieve the vault address
		wt, err := jws.ParseJWT([]byte(unwrapToken))
		if err != nil {
			return nil, errwrap.Wrapf("error decoding token: {{err}}", err)
		}
		if wt == nil {
			return nil, errors.New("nil decoded token")
		}

		addrRaw := wt.Claims().Get("addr")
		if addrRaw == nil {
			return nil, errors.New("decoded token does not contain the active node's api_addr")
		}
		vaultAddr, ok := addrRaw.(string)
		if !ok {
			return nil, errors.New("decoded token's api_addr not valid")
		}
		if vaultAddr == "" {
			return nil, errors.New(`no vault api_addr found`)
		}

		// Sanity check the value
		if _, err := url.Parse(vaultAddr); err != nil {
			return nil, errwrap.Wrapf("error parsing the vault api_addr: {{err}}", err)
		}

		// Unwrap the token
		clientConf := api.DefaultConfig()
		clientConf.Address = vaultAddr
		if apiTLSConfig != nil {
			err := clientConf.ConfigureTLS(apiTLSConfig)
			if err != nil {
				return nil, errwrap.Wrapf("error configuring api client {{err}}", err)
			}
		}
		client, err := api.NewClient(clientConf)
		if err != nil {
			return nil, errwrap.Wrapf("error during api client creation: {{err}}", err)
		}

		secret, err := client.Logical().Unwrap(unwrapToken)
		if err != nil {
			return nil, errwrap.Wrapf("error during token unwrap request: {{err}}", err)
		}
		if secret == nil {
			return nil, errors.New("error during token unwrap request: secret is nil")
		}

		// Retrieve and parse the server's certificate
		serverCertBytesRaw, ok := secret.Data["ServerCert"].(string)
		if !ok {
			return nil, errors.New("error unmarshalling certificate")
		}

		serverCertBytes, err := base64.StdEncoding.DecodeString(serverCertBytesRaw)
		if err != nil {
			return nil, errwrap.Wrapf("error parsing certificate: {{err}}", err)
		}

		serverCert, err := x509.ParseCertificate(serverCertBytes)
		if err != nil {
			return nil, errwrap.Wrapf("error parsing certificate: {{err}}", err)
		}

		// Retrieve and parse the server's private key
		serverKeyB64, ok := secret.Data["ServerKey"].(string)
		if !ok {
			return nil, errors.New("error unmarshalling certificate")
		}

		serverKeyRaw, err := base64.StdEncoding.DecodeString(serverKeyB64)
		if err != nil {
			return nil, errwrap.Wrapf("error parsing certificate: {{err}}", err)
		}

		serverKey, err := x509.ParseECPrivateKey(serverKeyRaw)
		if err != nil {
			return nil, errwrap.Wrapf("error parsing certificate: {{err}}", err)
		}

		// Add CA cert to the cert pool
		caCertPool := x509.NewCertPool()
		caCertPool.AddCert(serverCert)

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
			ServerName:   serverCert.Subject.CommonName,
		}
		tlsConfig.BuildNameToCertificate()

		return tlsConfig, nil
	}
}
