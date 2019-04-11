package cert

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/pem"
	mathrand "math/rand"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/hashicorp/go-sockaddr"

	"golang.org/x/net/http2"

	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"net"
	"os"
	"reflect"
	"testing"
	"time"

	cleanhttp "github.com/hashicorp/go-cleanhttp"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	vaulthttp "github.com/hashicorp/vault/http"

	rootcerts "github.com/hashicorp/go-rootcerts"
	"github.com/hashicorp/vault/builtin/logical/pki"
	"github.com/hashicorp/vault/helper/certutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	logicaltest "github.com/hashicorp/vault/logical/testing"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/mapstructure"
)

const (
	serverCertPath = "test-fixtures/cacert.pem"
	serverKeyPath  = "test-fixtures/cakey.pem"
	serverCAPath   = serverCertPath

	testRootCACertPath1 = "test-fixtures/testcacert1.pem"
	testRootCAKeyPath1  = "test-fixtures/testcakey1.pem"
	testCertPath1       = "test-fixtures/testissuedcert4.pem"
	testKeyPath1        = "test-fixtures/testissuedkey4.pem"
	testIssuedCertCRL   = "test-fixtures/issuedcertcrl"

	testRootCACertPath2 = "test-fixtures/testcacert2.pem"
	testRootCAKeyPath2  = "test-fixtures/testcakey2.pem"
	testRootCertCRL     = "test-fixtures/cacert2crl"
)

func generateTestCertAndConnState(t *testing.T, template *x509.Certificate) (string, tls.ConnectionState, error) {
	t.Helper()
	tempDir, err := ioutil.TempDir("", "vault-cert-auth-test-")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("test %s, temp dir %s", t.Name(), tempDir)
	caCertTemplate := &x509.Certificate{
		Subject: pkix.Name{
			CommonName: "localhost",
		},
		DNSNames:              []string{"localhost"},
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
		KeyUsage:              x509.KeyUsage(x509.KeyUsageCertSign | x509.KeyUsageCRLSign),
		SerialNumber:          big.NewInt(mathrand.Int63()),
		NotBefore:             time.Now().Add(-30 * time.Second),
		NotAfter:              time.Now().Add(262980 * time.Hour),
		BasicConstraintsValid: true,
		IsCA:                  true,
	}
	caKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	caBytes, err := x509.CreateCertificate(rand.Reader, caCertTemplate, caCertTemplate, caKey.Public(), caKey)
	if err != nil {
		t.Fatal(err)
	}
	caCert, err := x509.ParseCertificate(caBytes)
	if err != nil {
		t.Fatal(err)
	}
	caCertPEMBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	}
	err = ioutil.WriteFile(filepath.Join(tempDir, "ca_cert.pem"), pem.EncodeToMemory(caCertPEMBlock), 0755)
	if err != nil {
		t.Fatal(err)
	}
	marshaledCAKey, err := x509.MarshalECPrivateKey(caKey)
	if err != nil {
		t.Fatal(err)
	}
	caKeyPEMBlock := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: marshaledCAKey,
	}
	err = ioutil.WriteFile(filepath.Join(tempDir, "ca_key.pem"), pem.EncodeToMemory(caKeyPEMBlock), 0755)
	if err != nil {
		t.Fatal(err)
	}

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	certBytes, err := x509.CreateCertificate(rand.Reader, template, caCert, key.Public(), caKey)
	if err != nil {
		t.Fatal(err)
	}
	certPEMBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	}
	err = ioutil.WriteFile(filepath.Join(tempDir, "cert.pem"), pem.EncodeToMemory(certPEMBlock), 0755)
	if err != nil {
		t.Fatal(err)
	}
	marshaledKey, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		t.Fatal(err)
	}
	keyPEMBlock := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: marshaledKey,
	}
	err = ioutil.WriteFile(filepath.Join(tempDir, "key.pem"), pem.EncodeToMemory(keyPEMBlock), 0755)
	if err != nil {
		t.Fatal(err)
	}
	connInfo, err := testConnState(filepath.Join(tempDir, "cert.pem"), filepath.Join(tempDir, "key.pem"), filepath.Join(tempDir, "ca_cert.pem"))
	return tempDir, connInfo, err
}

// Unlike testConnState, this method does not use the same 'tls.Config' objects for
// both dialing and listening. Instead, it runs the server without specifying its CA.
// But the client, presents the CA cert of the server to trust the server.
// The client can present a cert and key which is completely independent of server's CA.
// The connection state returned will contain the certificate presented by the client.
func connectionState(serverCAPath, serverCertPath, serverKeyPath, clientCertPath, clientKeyPath string) (tls.ConnectionState, error) {
	serverKeyPair, err := tls.LoadX509KeyPair(serverCertPath, serverKeyPath)
	if err != nil {
		return tls.ConnectionState{}, err
	}
	// Prepare the listener configuration with server's key pair
	listenConf := &tls.Config{
		Certificates: []tls.Certificate{serverKeyPair},
		ClientAuth:   tls.RequestClientCert,
	}

	clientKeyPair, err := tls.LoadX509KeyPair(clientCertPath, clientKeyPath)
	if err != nil {
		return tls.ConnectionState{}, err
	}
	// Load the CA cert required by the client to authenticate the server.
	rootConfig := &rootcerts.Config{
		CAFile: serverCAPath,
	}
	serverCAs, err := rootcerts.LoadCACerts(rootConfig)
	if err != nil {
		return tls.ConnectionState{}, err
	}
	// Prepare the dial configuration that the client uses to establish the connection.
	dialConf := &tls.Config{
		Certificates: []tls.Certificate{clientKeyPair},
		RootCAs:      serverCAs,
	}

	// Start the server.
	list, err := tls.Listen("tcp", "127.0.0.1:0", listenConf)
	if err != nil {
		return tls.ConnectionState{}, err
	}
	defer list.Close()

	// Accept connections.
	serverErrors := make(chan error, 1)
	connState := make(chan tls.ConnectionState)
	go func() {
		defer close(connState)
		serverConn, err := list.Accept()
		if err != nil {
			serverErrors <- err
			close(serverErrors)
			return
		}
		defer serverConn.Close()

		// Read the ping
		buf := make([]byte, 4)
		_, err = serverConn.Read(buf)
		if (err != nil) && (err != io.EOF) {
			serverErrors <- err
			close(serverErrors)
			return
		}
		close(serverErrors)
		connState <- serverConn.(*tls.Conn).ConnectionState()
	}()

	// Establish a connection from the client side and write a few bytes.
	clientErrors := make(chan error, 1)
	go func() {
		addr := list.Addr().String()
		conn, err := tls.Dial("tcp", addr, dialConf)
		if err != nil {
			clientErrors <- err
			close(clientErrors)
			return
		}
		defer conn.Close()

		// Write ping
		_, err = conn.Write([]byte("ping"))
		if err != nil {
			clientErrors <- err
		}
		close(clientErrors)
	}()

	for err = range clientErrors {
		if err != nil {
			return tls.ConnectionState{}, fmt.Errorf("error in client goroutine:%v", err)
		}
	}

	for err = range serverErrors {
		if err != nil {
			return tls.ConnectionState{}, fmt.Errorf("error in server goroutine:%v", err)
		}
	}
	// Grab the current state
	return <-connState, nil
}

func TestBackend_PermittedDNSDomainsIntermediateCA(t *testing.T) {
	// Enable PKI secret engine and Cert auth method
	coreConfig := &vault.CoreConfig{
		DisableMlock: true,
		DisableCache: true,
		Logger:       log.NewNullLogger(),
		CredentialBackends: map[string]logical.Factory{
			"cert": Factory,
		},
		LogicalBackends: map[string]logical.Factory{
			"pki": pki.Factory,
		},
	}
	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()
	cores := cluster.Cores
	vault.TestWaitActive(t, cores[0].Core)
	client := cores[0].Client

	var err error

	// Mount /pki as a root CA
	err = client.Sys().Mount("pki", &api.MountInput{
		Type: "pki",
		Config: api.MountConfigInput{
			DefaultLeaseTTL: "16h",
			MaxLeaseTTL:     "32h",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	// Set the cluster's certificate as the root CA in /pki
	pemBundleRootCA := string(cluster.CACertPEM) + string(cluster.CAKeyPEM)
	_, err = client.Logical().Write("pki/config/ca", map[string]interface{}{
		"pem_bundle": pemBundleRootCA,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Mount /pki2 to operate as an intermediate CA
	err = client.Sys().Mount("pki2", &api.MountInput{
		Type: "pki",
		Config: api.MountConfigInput{
			DefaultLeaseTTL: "16h",
			MaxLeaseTTL:     "32h",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create a CSR for the intermediate CA
	secret, err := client.Logical().Write("pki2/intermediate/generate/internal", nil)
	if err != nil {
		t.Fatal(err)
	}
	intermediateCSR := secret.Data["csr"].(string)

	// Sign the intermediate CSR using /pki
	secret, err = client.Logical().Write("pki/root/sign-intermediate", map[string]interface{}{
		"permitted_dns_domains": ".myvault.com",
		"csr":                   intermediateCSR,
	})
	if err != nil {
		t.Fatal(err)
	}
	intermediateCertPEM := secret.Data["certificate"].(string)

	// Configure the intermediate cert as the CA in /pki2
	_, err = client.Logical().Write("pki2/intermediate/set-signed", map[string]interface{}{
		"certificate": intermediateCertPEM,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create a role on the intermediate CA mount
	_, err = client.Logical().Write("pki2/roles/myvault-dot-com", map[string]interface{}{
		"allowed_domains":  "myvault.com",
		"allow_subdomains": "true",
		"max_ttl":          "5m",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Issue a leaf cert using the intermediate CA
	secret, err = client.Logical().Write("pki2/issue/myvault-dot-com", map[string]interface{}{
		"common_name": "cert.myvault.com",
		"format":      "pem",
		"ip_sans":     "127.0.0.1",
	})
	if err != nil {
		t.Fatal(err)
	}
	leafCertPEM := secret.Data["certificate"].(string)
	leafCertKeyPEM := secret.Data["private_key"].(string)

	// Enable the cert auth method
	err = client.Sys().EnableAuthWithOptions("cert", &api.EnableAuthOptions{
		Type: "cert",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Set the intermediate CA cert as a trusted certificate in the backend
	_, err = client.Logical().Write("auth/cert/certs/myvault-dot-com", map[string]interface{}{
		"display_name": "myvault.com",
		"policies":     "default",
		"certificate":  intermediateCertPEM,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create temporary files for CA cert, client cert and client cert key.
	// This is used to configure TLS in the api client.
	caCertFile, err := ioutil.TempFile("", "caCert")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(caCertFile.Name())
	if _, err := caCertFile.Write([]byte(cluster.CACertPEM)); err != nil {
		t.Fatal(err)
	}
	if err := caCertFile.Close(); err != nil {
		t.Fatal(err)
	}

	leafCertFile, err := ioutil.TempFile("", "leafCert")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(leafCertFile.Name())
	if _, err := leafCertFile.Write([]byte(leafCertPEM)); err != nil {
		t.Fatal(err)
	}
	if err := leafCertFile.Close(); err != nil {
		t.Fatal(err)
	}

	leafCertKeyFile, err := ioutil.TempFile("", "leafCertKey")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(leafCertKeyFile.Name())
	if _, err := leafCertKeyFile.Write([]byte(leafCertKeyPEM)); err != nil {
		t.Fatal(err)
	}
	if err := leafCertKeyFile.Close(); err != nil {
		t.Fatal(err)
	}

	// This function is a copy-pasta from the NewTestCluster, with the
	// modification to reconfigure the TLS on the api client with the leaf
	// certificate generated above.
	getAPIClient := func(port int, tlsConfig *tls.Config) *api.Client {
		transport := cleanhttp.DefaultPooledTransport()
		transport.TLSClientConfig = tlsConfig.Clone()
		if err := http2.ConfigureTransport(transport); err != nil {
			t.Fatal(err)
		}
		client := &http.Client{
			Transport: transport,
			CheckRedirect: func(*http.Request, []*http.Request) error {
				// This can of course be overridden per-test by using its own client
				return fmt.Errorf("redirects not allowed in these tests")
			},
		}
		config := api.DefaultConfig()
		if config.Error != nil {
			t.Fatal(config.Error)
		}
		config.Address = fmt.Sprintf("https://127.0.0.1:%d", port)
		config.HttpClient = client

		// Set the above issued certificates as the client certificates
		config.ConfigureTLS(&api.TLSConfig{
			CACert:     caCertFile.Name(),
			ClientCert: leafCertFile.Name(),
			ClientKey:  leafCertKeyFile.Name(),
		})

		apiClient, err := api.NewClient(config)
		if err != nil {
			t.Fatal(err)
		}
		return apiClient
	}

	// Create a new api client with the desired TLS configuration
	newClient := getAPIClient(cores[0].Listeners[0].Address.Port, cores[0].TLSConfig)

	secret, err = newClient.Logical().Write("auth/cert/login", map[string]interface{}{
		"name": "myvault-dot-com",
	})
	if err != nil {
		t.Fatal(err)
	}
	if secret.Auth == nil || secret.Auth.ClientToken == "" {
		t.Fatalf("expected a successful authentication")
	}
}

func TestBackend_NonCAExpiry(t *testing.T) {
	var resp *logical.Response
	var err error

	// Create a self-signed certificate and issue a leaf certificate using the
	// CA cert
	template := &x509.Certificate{
		SerialNumber: big.NewInt(1234),
		Subject: pkix.Name{
			CommonName:         "localhost",
			Organization:       []string{"hashicorp"},
			OrganizationalUnit: []string{"vault"},
		},
		BasicConstraintsValid: true,
		NotBefore:             time.Now().Add(-30 * time.Second),
		NotAfter:              time.Now().Add(50 * time.Second),
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsage(x509.KeyUsageCertSign | x509.KeyUsageCRLSign),
	}

	// Set IP SAN
	parsedIP := net.ParseIP("127.0.0.1")
	if parsedIP == nil {
		t.Fatalf("failed to create parsed IP")
	}
	template.IPAddresses = []net.IP{parsedIP}

	// Private key for CA cert
	caPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}

	// Marshalling to be able to create PEM file
	caPrivateKeyBytes := x509.MarshalPKCS1PrivateKey(caPrivateKey)

	caPublicKey := &caPrivateKey.PublicKey

	template.IsCA = true

	caCertBytes, err := x509.CreateCertificate(rand.Reader, template, template, caPublicKey, caPrivateKey)
	if err != nil {
		t.Fatal(err)
	}

	caCert, err := x509.ParseCertificate(caCertBytes)
	if err != nil {
		t.Fatal(err)
	}

	parsedCaBundle := &certutil.ParsedCertBundle{
		Certificate:      caCert,
		CertificateBytes: caCertBytes,
		PrivateKeyBytes:  caPrivateKeyBytes,
		PrivateKeyType:   certutil.RSAPrivateKey,
	}

	caCertBundle, err := parsedCaBundle.ToCertBundle()
	if err != nil {
		t.Fatal(err)
	}

	caCertFile, err := ioutil.TempFile("", "caCert")
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(caCertFile.Name())

	if _, err := caCertFile.Write([]byte(caCertBundle.Certificate)); err != nil {
		t.Fatal(err)
	}
	if err := caCertFile.Close(); err != nil {
		t.Fatal(err)
	}

	caKeyFile, err := ioutil.TempFile("", "caKey")
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(caKeyFile.Name())

	if _, err := caKeyFile.Write([]byte(caCertBundle.PrivateKey)); err != nil {
		t.Fatal(err)
	}
	if err := caKeyFile.Close(); err != nil {
		t.Fatal(err)
	}

	// Prepare template for non-CA cert

	template.IsCA = false
	template.SerialNumber = big.NewInt(5678)

	template.KeyUsage = x509.KeyUsage(x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign)
	issuedPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}

	issuedPrivateKeyBytes := x509.MarshalPKCS1PrivateKey(issuedPrivateKey)

	issuedPublicKey := &issuedPrivateKey.PublicKey

	// Keep a short certificate lifetime so logins can be tested both when
	// cert is valid and when it gets expired
	template.NotBefore = time.Now().Add(-2 * time.Second)
	template.NotAfter = time.Now().Add(3 * time.Second)

	issuedCertBytes, err := x509.CreateCertificate(rand.Reader, template, caCert, issuedPublicKey, caPrivateKey)
	if err != nil {
		t.Fatal(err)
	}

	issuedCert, err := x509.ParseCertificate(issuedCertBytes)
	if err != nil {
		t.Fatal(err)
	}

	parsedIssuedBundle := &certutil.ParsedCertBundle{
		Certificate:      issuedCert,
		CertificateBytes: issuedCertBytes,
		PrivateKeyBytes:  issuedPrivateKeyBytes,
		PrivateKeyType:   certutil.RSAPrivateKey,
	}

	issuedCertBundle, err := parsedIssuedBundle.ToCertBundle()
	if err != nil {
		t.Fatal(err)
	}

	issuedCertFile, err := ioutil.TempFile("", "issuedCert")
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(issuedCertFile.Name())

	if _, err := issuedCertFile.Write([]byte(issuedCertBundle.Certificate)); err != nil {
		t.Fatal(err)
	}
	if err := issuedCertFile.Close(); err != nil {
		t.Fatal(err)
	}

	issuedKeyFile, err := ioutil.TempFile("", "issuedKey")
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(issuedKeyFile.Name())

	if _, err := issuedKeyFile.Write([]byte(issuedCertBundle.PrivateKey)); err != nil {
		t.Fatal(err)
	}
	if err := issuedKeyFile.Close(); err != nil {
		t.Fatal(err)
	}

	config := logical.TestBackendConfig()
	storage := &logical.InmemStorage{}
	config.StorageView = storage

	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	// Register the Non-CA certificate of the client key pair
	certData := map[string]interface{}{
		"certificate":  issuedCertBundle.Certificate,
		"policies":     "abc",
		"display_name": "cert1",
		"ttl":          10000,
	}
	certReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "certs/cert1",
		Storage:   storage,
		Data:      certData,
	}

	resp, err = b.HandleRequest(context.Background(), certReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	// Create connection state using the certificates generated
	connState, err := connectionState(caCertFile.Name(), caCertFile.Name(), caKeyFile.Name(), issuedCertFile.Name(), issuedKeyFile.Name())
	if err != nil {
		t.Fatalf("error testing connection state:%v", err)
	}

	loginReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Storage:   storage,
		Path:      "login",
		Connection: &logical.Connection{
			ConnState: &connState,
		},
	}

	// Login when the certificate is still valid. Login should succeed.
	resp, err = b.HandleRequest(context.Background(), loginReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	// Wait until the certificate expires
	time.Sleep(5 * time.Second)

	// Login attempt after certificate expiry should fail
	resp, err = b.HandleRequest(context.Background(), loginReq)
	if err == nil {
		t.Fatalf("expected error due to expired certificate")
	}
}

func TestBackend_RegisteredNonCA_CRL(t *testing.T) {
	config := logical.TestBackendConfig()
	storage := &logical.InmemStorage{}
	config.StorageView = storage

	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	nonCACert, err := ioutil.ReadFile(testCertPath1)
	if err != nil {
		t.Fatal(err)
	}

	// Register the Non-CA certificate of the client key pair
	certData := map[string]interface{}{
		"certificate":  nonCACert,
		"policies":     "abc",
		"display_name": "cert1",
		"ttl":          10000,
	}
	certReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "certs/cert1",
		Storage:   storage,
		Data:      certData,
	}

	resp, err := b.HandleRequest(context.Background(), certReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	// Connection state is presenting the client Non-CA cert and its key.
	// This is exactly what is registered at the backend.
	connState, err := connectionState(serverCAPath, serverCertPath, serverKeyPath, testCertPath1, testKeyPath1)
	if err != nil {
		t.Fatalf("error testing connection state:%v", err)
	}
	loginReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Storage:   storage,
		Path:      "login",
		Connection: &logical.Connection{
			ConnState: &connState,
		},
	}
	// Login should succeed.
	resp, err = b.HandleRequest(context.Background(), loginReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	// Register a CRL containing the issued client certificate used above.
	issuedCRL, err := ioutil.ReadFile(testIssuedCertCRL)
	if err != nil {
		t.Fatal(err)
	}
	crlData := map[string]interface{}{
		"crl": issuedCRL,
	}
	crlReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Storage:   storage,
		Path:      "crls/issuedcrl",
		Data:      crlData,
	}
	resp, err = b.HandleRequest(context.Background(), crlReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	// Attempt login with the same connection state but with the CRL registered
	resp, err = b.HandleRequest(context.Background(), loginReq)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || !resp.IsError() {
		t.Fatalf("expected failure due to revoked certificate")
	}
}

func TestBackend_CRLs(t *testing.T) {
	config := logical.TestBackendConfig()
	storage := &logical.InmemStorage{}
	config.StorageView = storage

	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	clientCA1, err := ioutil.ReadFile(testRootCACertPath1)
	if err != nil {
		t.Fatal(err)
	}
	// Register the CA certificate of the client key pair
	certData := map[string]interface{}{
		"certificate":  clientCA1,
		"policies":     "abc",
		"display_name": "cert1",
		"ttl":          10000,
	}

	certReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "certs/cert1",
		Storage:   storage,
		Data:      certData,
	}

	resp, err := b.HandleRequest(context.Background(), certReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	// Connection state is presenting the client CA cert and its key.
	// This is exactly what is registered at the backend.
	connState, err := connectionState(serverCAPath, serverCertPath, serverKeyPath, testRootCACertPath1, testRootCAKeyPath1)
	if err != nil {
		t.Fatalf("error testing connection state:%v", err)
	}
	loginReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Storage:   storage,
		Path:      "login",
		Connection: &logical.Connection{
			ConnState: &connState,
		},
	}
	resp, err = b.HandleRequest(context.Background(), loginReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	// Now, without changing the registered client CA cert, present from
	// the client side, a cert issued using the registered CA.
	connState, err = connectionState(serverCAPath, serverCertPath, serverKeyPath, testCertPath1, testKeyPath1)
	if err != nil {
		t.Fatalf("error testing connection state: %v", err)
	}
	loginReq.Connection.ConnState = &connState

	// Attempt login with the updated connection
	resp, err = b.HandleRequest(context.Background(), loginReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	// Register a CRL containing the issued client certificate used above.
	issuedCRL, err := ioutil.ReadFile(testIssuedCertCRL)
	if err != nil {
		t.Fatal(err)
	}
	crlData := map[string]interface{}{
		"crl": issuedCRL,
	}

	crlReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Storage:   storage,
		Path:      "crls/issuedcrl",
		Data:      crlData,
	}
	resp, err = b.HandleRequest(context.Background(), crlReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	// Attempt login with the revoked certificate.
	resp, err = b.HandleRequest(context.Background(), loginReq)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || !resp.IsError() {
		t.Fatalf("expected failure due to revoked certificate")
	}

	// Register a different client CA certificate.
	clientCA2, err := ioutil.ReadFile(testRootCACertPath2)
	if err != nil {
		t.Fatal(err)
	}
	certData["certificate"] = clientCA2
	resp, err = b.HandleRequest(context.Background(), certReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	// Test login using a different client CA cert pair.
	connState, err = connectionState(serverCAPath, serverCertPath, serverKeyPath, testRootCACertPath2, testRootCAKeyPath2)
	if err != nil {
		t.Fatalf("error testing connection state: %v", err)
	}
	loginReq.Connection.ConnState = &connState

	// Attempt login with the updated connection
	resp, err = b.HandleRequest(context.Background(), loginReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	// Register a CRL containing the root CA certificate used above.
	rootCRL, err := ioutil.ReadFile(testRootCertCRL)
	if err != nil {
		t.Fatal(err)
	}
	crlData["crl"] = rootCRL
	resp, err = b.HandleRequest(context.Background(), crlReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	// Attempt login with the same connection state but with the CRL registered
	resp, err = b.HandleRequest(context.Background(), loginReq)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || !resp.IsError() {
		t.Fatalf("expected failure due to revoked certificate")
	}
}

func testFactory(t *testing.T) logical.Backend {
	b, err := Factory(context.Background(), &logical.BackendConfig{
		System: &logical.StaticSystemView{
			DefaultLeaseTTLVal: 1000 * time.Second,
			MaxLeaseTTLVal:     1800 * time.Second,
		},
		StorageView: &logical.InmemStorage{},
	})
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	return b
}

// Test the certificates being registered to the backend
func TestBackend_CertWrites(t *testing.T) {
	// CA cert
	ca1, err := ioutil.ReadFile("test-fixtures/root/rootcacert.pem")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	// Non CA Cert
	ca2, err := ioutil.ReadFile("test-fixtures/keys/cert.pem")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	// Non CA cert without TLS web client authentication
	ca3, err := ioutil.ReadFile("test-fixtures/noclientauthcert.pem")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	tc := logicaltest.TestCase{
		CredentialBackend: testFactory(t),
		Steps: []logicaltest.TestStep{
			testAccStepCert(t, "aaa", ca1, "foo", allowed{}, false),
			testAccStepCert(t, "bbb", ca2, "foo", allowed{}, false),
			testAccStepCert(t, "ccc", ca3, "foo", allowed{}, true),
		},
	}
	tc.Steps = append(tc.Steps, testAccStepListCerts(t, []string{"aaa", "bbb"})...)
	logicaltest.Test(t, tc)
}

// Test a client trusted by a CA
func TestBackend_basic_CA(t *testing.T) {
	connState, err := testConnState("test-fixtures/keys/cert.pem",
		"test-fixtures/keys/key.pem", "test-fixtures/root/rootcacert.pem")
	if err != nil {
		t.Fatalf("error testing connection state: %v", err)
	}
	ca, err := ioutil.ReadFile("test-fixtures/root/rootcacert.pem")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	logicaltest.Test(t, logicaltest.TestCase{
		CredentialBackend: testFactory(t),
		Steps: []logicaltest.TestStep{
			testAccStepCert(t, "web", ca, "foo", allowed{}, false),
			testAccStepLogin(t, connState),
			testAccStepCertLease(t, "web", ca, "foo"),
			testAccStepCertTTL(t, "web", ca, "foo"),
			testAccStepLogin(t, connState),
			testAccStepCertMaxTTL(t, "web", ca, "foo"),
			testAccStepLogin(t, connState),
			testAccStepCertNoLease(t, "web", ca, "foo"),
			testAccStepLoginDefaultLease(t, connState),
			testAccStepCert(t, "web", ca, "foo", allowed{names: "*.example.com"}, false),
			testAccStepLogin(t, connState),
			testAccStepCert(t, "web", ca, "foo", allowed{names: "*.invalid.com"}, false),
			testAccStepLoginInvalid(t, connState),
		},
	})
}

// Test CRL behavior
func TestBackend_Basic_CRLs(t *testing.T) {
	connState, err := testConnState("test-fixtures/keys/cert.pem",
		"test-fixtures/keys/key.pem", "test-fixtures/root/rootcacert.pem")
	if err != nil {
		t.Fatalf("error testing connection state: %v", err)
	}
	ca, err := ioutil.ReadFile("test-fixtures/root/rootcacert.pem")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	crl, err := ioutil.ReadFile("test-fixtures/root/root.crl")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	logicaltest.Test(t, logicaltest.TestCase{
		CredentialBackend: testFactory(t),
		Steps: []logicaltest.TestStep{
			testAccStepCertNoLease(t, "web", ca, "foo"),
			testAccStepLoginDefaultLease(t, connState),
			testAccStepAddCRL(t, crl, connState),
			testAccStepReadCRL(t, connState),
			testAccStepLoginInvalid(t, connState),
			testAccStepDeleteCRL(t, connState),
			testAccStepLoginDefaultLease(t, connState),
		},
	})
}

// Test a self-signed client (root CA) that is trusted
func TestBackend_basic_singleCert(t *testing.T) {
	connState, err := testConnState("test-fixtures/root/rootcacert.pem",
		"test-fixtures/root/rootcakey.pem", "test-fixtures/root/rootcacert.pem")
	if err != nil {
		t.Fatalf("error testing connection state: %v", err)
	}
	ca, err := ioutil.ReadFile("test-fixtures/root/rootcacert.pem")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	logicaltest.Test(t, logicaltest.TestCase{
		CredentialBackend: testFactory(t),
		Steps: []logicaltest.TestStep{
			testAccStepCert(t, "web", ca, "foo", allowed{}, false),
			testAccStepLogin(t, connState),
			testAccStepCert(t, "web", ca, "foo", allowed{names: "example.com"}, false),
			testAccStepLogin(t, connState),
			testAccStepCert(t, "web", ca, "foo", allowed{names: "invalid"}, false),
			testAccStepLoginInvalid(t, connState),
			testAccStepCert(t, "web", ca, "foo", allowed{ext: "1.2.3.4:invalid"}, false),
			testAccStepLoginInvalid(t, connState),
		},
	})
}

func TestBackend_common_name_singleCert(t *testing.T) {
	connState, err := testConnState("test-fixtures/root/rootcacert.pem",
		"test-fixtures/root/rootcakey.pem", "test-fixtures/root/rootcacert.pem")
	if err != nil {
		t.Fatalf("error testing connection state: %v", err)
	}
	ca, err := ioutil.ReadFile("test-fixtures/root/rootcacert.pem")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	logicaltest.Test(t, logicaltest.TestCase{
		CredentialBackend: testFactory(t),
		Steps: []logicaltest.TestStep{
			testAccStepCert(t, "web", ca, "foo", allowed{}, false),
			testAccStepLogin(t, connState),
			testAccStepCert(t, "web", ca, "foo", allowed{common_names: "example.com"}, false),
			testAccStepLogin(t, connState),
			testAccStepCert(t, "web", ca, "foo", allowed{common_names: "invalid"}, false),
			testAccStepLoginInvalid(t, connState),
			testAccStepCert(t, "web", ca, "foo", allowed{ext: "1.2.3.4:invalid"}, false),
			testAccStepLoginInvalid(t, connState),
		},
	})
}

// Test a self-signed client with custom ext (root CA) that is trusted
func TestBackend_ext_singleCert(t *testing.T) {
	connState, err := testConnState(
		"test-fixtures/root/rootcawextcert.pem",
		"test-fixtures/root/rootcawextkey.pem",
		"test-fixtures/root/rootcacert.pem",
	)
	if err != nil {
		t.Fatalf("error testing connection state: %v", err)
	}
	ca, err := ioutil.ReadFile("test-fixtures/root/rootcacert.pem")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	logicaltest.Test(t, logicaltest.TestCase{
		CredentialBackend: testFactory(t),
		Steps: []logicaltest.TestStep{
			testAccStepCert(t, "web", ca, "foo", allowed{ext: "2.1.1.1:A UTF8String Extension"}, false),
			testAccStepLogin(t, connState),
			testAccStepCert(t, "web", ca, "foo", allowed{ext: "2.1.1.1:*,2.1.1.2:A UTF8*"}, false),
			testAccStepLogin(t, connState),
			testAccStepCert(t, "web", ca, "foo", allowed{ext: "1.2.3.45:*"}, false),
			testAccStepLoginInvalid(t, connState),
			testAccStepCert(t, "web", ca, "foo", allowed{ext: "2.1.1.1:The Wrong Value"}, false),
			testAccStepLoginInvalid(t, connState),
			testAccStepCert(t, "web", ca, "foo", allowed{ext: "2.1.1.1:*,2.1.1.2:The Wrong Value"}, false),
			testAccStepLoginInvalid(t, connState),
			testAccStepCert(t, "web", ca, "foo", allowed{ext: "2.1.1.1:"}, false),
			testAccStepLoginInvalid(t, connState),
			testAccStepCert(t, "web", ca, "foo", allowed{ext: "2.1.1.1:,2.1.1.2:*"}, false),
			testAccStepLoginInvalid(t, connState),
			testAccStepCert(t, "web", ca, "foo", allowed{names: "example.com", ext: "2.1.1.1:A UTF8String Extension"}, false),
			testAccStepLogin(t, connState),
			testAccStepCert(t, "web", ca, "foo", allowed{names: "example.com", ext: "2.1.1.1:*,2.1.1.2:A UTF8*"}, false),
			testAccStepLogin(t, connState),
			testAccStepCert(t, "web", ca, "foo", allowed{names: "example.com", ext: "1.2.3.45:*"}, false),
			testAccStepLoginInvalid(t, connState),
			testAccStepCert(t, "web", ca, "foo", allowed{names: "example.com", ext: "2.1.1.1:The Wrong Value"}, false),
			testAccStepLoginInvalid(t, connState),
			testAccStepCert(t, "web", ca, "foo", allowed{names: "example.com", ext: "2.1.1.1:*,2.1.1.2:The Wrong Value"}, false),
			testAccStepLoginInvalid(t, connState),
			testAccStepCert(t, "web", ca, "foo", allowed{names: "invalid", ext: "2.1.1.1:A UTF8String Extension"}, false),
			testAccStepLoginInvalid(t, connState),
			testAccStepCert(t, "web", ca, "foo", allowed{names: "invalid", ext: "2.1.1.1:*,2.1.1.2:A UTF8*"}, false),
			testAccStepLoginInvalid(t, connState),
			testAccStepCert(t, "web", ca, "foo", allowed{names: "invalid", ext: "1.2.3.45:*"}, false),
			testAccStepLoginInvalid(t, connState),
			testAccStepCert(t, "web", ca, "foo", allowed{names: "invalid", ext: "2.1.1.1:The Wrong Value"}, false),
			testAccStepLoginInvalid(t, connState),
			testAccStepCert(t, "web", ca, "foo", allowed{names: "invalid", ext: "2.1.1.1:*,2.1.1.2:The Wrong Value"}, false),
			testAccStepLoginInvalid(t, connState),
		},
	})
}

// Test a self-signed client with URI alt names (root CA) that is trusted
func TestBackend_dns_singleCert(t *testing.T) {
	connState, err := testConnState(
		"test-fixtures/root/rootcawdnscert.pem",
		"test-fixtures/root/rootcawdnskey.pem",
		"test-fixtures/root/rootcacert.pem",
	)
	if err != nil {
		t.Fatalf("error testing connection state: %v", err)
	}
	ca, err := ioutil.ReadFile("test-fixtures/root/rootcacert.pem")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	logicaltest.Test(t, logicaltest.TestCase{
		CredentialBackend: testFactory(t),
		Steps: []logicaltest.TestStep{
			testAccStepCert(t, "web", ca, "foo", allowed{dns: "example.com"}, false),
			testAccStepLogin(t, connState),
			testAccStepCert(t, "web", ca, "foo", allowed{dns: "*ample.com"}, false),
			testAccStepLogin(t, connState),
			testAccStepCert(t, "web", ca, "foo", allowed{dns: "notincert.com"}, false),
			testAccStepLoginInvalid(t, connState),
			testAccStepCert(t, "web", ca, "foo", allowed{dns: "abc"}, false),
			testAccStepLoginInvalid(t, connState),
			testAccStepCert(t, "web", ca, "foo", allowed{dns: "*.example.com"}, false),
			testAccStepLoginInvalid(t, connState),
		},
	})
}

// Test a self-signed client with URI alt names (root CA) that is trusted
func TestBackend_email_singleCert(t *testing.T) {
	connState, err := testConnState(
		"test-fixtures/root/rootcawemailcert.pem",
		"test-fixtures/root/rootcawemailkey.pem",
		"test-fixtures/root/rootcacert.pem",
	)
	if err != nil {
		t.Fatalf("error testing connection state: %v", err)
	}
	ca, err := ioutil.ReadFile("test-fixtures/root/rootcacert.pem")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	logicaltest.Test(t, logicaltest.TestCase{
		CredentialBackend: testFactory(t),
		Steps: []logicaltest.TestStep{
			testAccStepCert(t, "web", ca, "foo", allowed{emails: "valid@example.com"}, false),
			testAccStepLogin(t, connState),
			testAccStepCert(t, "web", ca, "foo", allowed{emails: "*@example.com"}, false),
			testAccStepLogin(t, connState),
			testAccStepCert(t, "web", ca, "foo", allowed{emails: "invalid@notincert.com"}, false),
			testAccStepLoginInvalid(t, connState),
			testAccStepCert(t, "web", ca, "foo", allowed{emails: "abc"}, false),
			testAccStepLoginInvalid(t, connState),
			testAccStepCert(t, "web", ca, "foo", allowed{emails: "*.example.com"}, false),
			testAccStepLoginInvalid(t, connState),
		},
	})
}

// Test a self-signed client with OU (root CA) that is trusted
func TestBackend_organizationalUnit_singleCert(t *testing.T) {
	connState, err := testConnState(
		"test-fixtures/root/rootcawoucert.pem",
		"test-fixtures/root/rootcawoukey.pem",
		"test-fixtures/root/rootcawoucert.pem",
	)
	if err != nil {
		t.Fatalf("error testing connection state: %v", err)
	}
	ca, err := ioutil.ReadFile("test-fixtures/root/rootcawoucert.pem")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	logicaltest.Test(t, logicaltest.TestCase{
		CredentialBackend: testFactory(t),
		Steps: []logicaltest.TestStep{
			testAccStepCert(t, "web", ca, "foo", allowed{organizational_units: "engineering"}, false),
			testAccStepLogin(t, connState),
			testAccStepCert(t, "web", ca, "foo", allowed{organizational_units: "eng*"}, false),
			testAccStepLogin(t, connState),
			testAccStepCert(t, "web", ca, "foo", allowed{organizational_units: "engineering,finance"}, false),
			testAccStepLogin(t, connState),
			testAccStepCert(t, "web", ca, "foo", allowed{organizational_units: "foo"}, false),
			testAccStepLoginInvalid(t, connState),
		},
	})
}

// Test a self-signed client with URI alt names (root CA) that is trusted
func TestBackend_uri_singleCert(t *testing.T) {
	u, err := url.Parse("spiffe://example.com/host")
	if err != nil {
		t.Fatal(err)
	}
	certTemplate := &x509.Certificate{
		Subject: pkix.Name{
			CommonName: "example.com",
		},
		DNSNames:    []string{"example.com"},
		IPAddresses: []net.IP{net.ParseIP("127.0.0.1")},
		URIs:        []*url.URL{u},
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
			x509.ExtKeyUsageClientAuth,
		},
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageKeyAgreement,
		SerialNumber: big.NewInt(mathrand.Int63()),
		NotBefore:    time.Now().Add(-30 * time.Second),
		NotAfter:     time.Now().Add(262980 * time.Hour),
	}

	tempDir, connState, err := generateTestCertAndConnState(t, certTemplate)
	if tempDir != "" {
		defer os.RemoveAll(tempDir)
	}
	if err != nil {
		t.Fatalf("error testing connection state: %v", err)
	}
	ca, err := ioutil.ReadFile(filepath.Join(tempDir, "ca_cert.pem"))
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	logicaltest.Test(t, logicaltest.TestCase{
		CredentialBackend: testFactory(t),
		Steps: []logicaltest.TestStep{
			testAccStepCert(t, "web", ca, "foo", allowed{uris: "spiffe://example.com/*"}, false),
			testAccStepLogin(t, connState),
			testAccStepCert(t, "web", ca, "foo", allowed{uris: "spiffe://example.com/host"}, false),
			testAccStepLogin(t, connState),
			testAccStepCert(t, "web", ca, "foo", allowed{uris: "spiffe://example.com/invalid"}, false),
			testAccStepLoginInvalid(t, connState),
			testAccStepCert(t, "web", ca, "foo", allowed{uris: "abc"}, false),
			testAccStepLoginInvalid(t, connState),
			testAccStepCert(t, "web", ca, "foo", allowed{uris: "http://www.google.com"}, false),
			testAccStepLoginInvalid(t, connState),
		},
	})
}

// Test against a collection of matching and non-matching rules
func TestBackend_mixed_constraints(t *testing.T) {
	connState, err := testConnState("test-fixtures/keys/cert.pem",
		"test-fixtures/keys/key.pem", "test-fixtures/root/rootcacert.pem")
	if err != nil {
		t.Fatalf("error testing connection state: %v", err)
	}
	ca, err := ioutil.ReadFile("test-fixtures/root/rootcacert.pem")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	logicaltest.Test(t, logicaltest.TestCase{
		CredentialBackend: testFactory(t),
		Steps: []logicaltest.TestStep{
			testAccStepCert(t, "1unconstrained", ca, "foo", allowed{}, false),
			testAccStepCert(t, "2matching", ca, "foo", allowed{names: "*.example.com,whatever"}, false),
			testAccStepCert(t, "3invalid", ca, "foo", allowed{names: "invalid"}, false),
			testAccStepLogin(t, connState),
			// Assumes CertEntries are processed in alphabetical order (due to store.List), so we only match 2matching if 1unconstrained doesn't match
			testAccStepLoginWithName(t, connState, "2matching"),
			testAccStepLoginWithNameInvalid(t, connState, "3invalid"),
		},
	})
}

// Test an untrusted client
func TestBackend_untrusted(t *testing.T) {
	connState, err := testConnState("test-fixtures/keys/cert.pem",
		"test-fixtures/keys/key.pem", "test-fixtures/root/rootcacert.pem")
	if err != nil {
		t.Fatalf("error testing connection state: %v", err)
	}
	logicaltest.Test(t, logicaltest.TestCase{
		CredentialBackend: testFactory(t),
		Steps: []logicaltest.TestStep{
			testAccStepLoginInvalid(t, connState),
		},
	})
}

func TestBackend_validCIDR(t *testing.T) {
	config := logical.TestBackendConfig()
	storage := &logical.InmemStorage{}
	config.StorageView = storage

	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	connState, err := testConnState("test-fixtures/keys/cert.pem",
		"test-fixtures/keys/key.pem", "test-fixtures/root/rootcacert.pem")
	if err != nil {
		t.Fatalf("error testing connection state: %v", err)
	}
	ca, err := ioutil.ReadFile("test-fixtures/root/rootcacert.pem")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	name := "web"
	boundCIDRs := []string{"127.0.0.1", "128.252.0.0/16"}

	addCertReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "certs/" + name,
		Data: map[string]interface{}{
			"certificate":         string(ca),
			"policies":            "foo",
			"display_name":        name,
			"allowed_names":       "",
			"required_extensions": "",
			"lease":               1000,
			"bound_cidrs":         boundCIDRs,
		},
		Storage:    storage,
		Connection: &logical.Connection{ConnState: &connState},
	}

	_, err = b.HandleRequest(context.Background(), addCertReq)
	if err != nil {
		t.Fatal(err)
	}

	readCertReq := &logical.Request{
		Operation:  logical.ReadOperation,
		Path:       "certs/" + name,
		Storage:    storage,
		Connection: &logical.Connection{ConnState: &connState},
	}

	readResult, err := b.HandleRequest(context.Background(), readCertReq)
	cidrsResult := readResult.Data["bound_cidrs"].([]*sockaddr.SockAddrMarshaler)

	if cidrsResult[0].String() != boundCIDRs[0] ||
		cidrsResult[1].String() != boundCIDRs[1] {
		t.Fatalf("bound_cidrs couldn't be set correctly, EXPECTED: %v, ACTUAL: %v", boundCIDRs, cidrsResult)
	}

	loginReq := &logical.Request{
		Operation:       logical.UpdateOperation,
		Path:            "login",
		Unauthenticated: true,
		Data: map[string]interface{}{
			"name": name,
		},
		Storage:    storage,
		Connection: &logical.Connection{ConnState: &connState},
	}

	// override the remote address with an IPV4 that is authorized
	loginReq.Connection.RemoteAddr = "127.0.0.1/32"

	_, err = b.HandleRequest(context.Background(), loginReq)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestBackend_invalidCIDR(t *testing.T) {
	config := logical.TestBackendConfig()
	storage := &logical.InmemStorage{}
	config.StorageView = storage

	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	connState, err := testConnState("test-fixtures/keys/cert.pem",
		"test-fixtures/keys/key.pem", "test-fixtures/root/rootcacert.pem")
	if err != nil {
		t.Fatalf("error testing connection state: %v", err)
	}
	ca, err := ioutil.ReadFile("test-fixtures/root/rootcacert.pem")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	name := "web"

	addCertReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "certs/" + name,
		Data: map[string]interface{}{
			"certificate":         string(ca),
			"policies":            "foo",
			"display_name":        name,
			"allowed_names":       "",
			"required_extensions": "",
			"lease":               1000,
			"bound_cidrs":         []string{"127.0.0.1/32", "128.252.0.0/16"},
		},
		Storage:    storage,
		Connection: &logical.Connection{ConnState: &connState},
	}

	_, err = b.HandleRequest(context.Background(), addCertReq)
	if err != nil {
		t.Fatal(err)
	}

	loginReq := &logical.Request{
		Operation:       logical.UpdateOperation,
		Path:            "login",
		Unauthenticated: true,
		Data: map[string]interface{}{
			"name": name,
		},
		Storage:    storage,
		Connection: &logical.Connection{ConnState: &connState},
	}

	// override the remote address with an IPV4 that isn't authorized
	loginReq.Connection.RemoteAddr = "127.0.0.1/8"

	_, err = b.HandleRequest(context.Background(), loginReq)
	if err == nil {
		t.Fatal("expected \"ERROR: permission denied\"")
	}
}

func testAccStepAddCRL(t *testing.T, crl []byte, connState tls.ConnectionState) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "crls/test",
		ConnState: &connState,
		Data: map[string]interface{}{
			"crl": crl,
		},
	}
}

func testAccStepReadCRL(t *testing.T, connState tls.ConnectionState) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path:      "crls/test",
		ConnState: &connState,
		Check: func(resp *logical.Response) error {
			crlInfo := CRLInfo{}
			err := mapstructure.Decode(resp.Data, &crlInfo)
			if err != nil {
				t.Fatalf("err: %v", err)
			}
			if len(crlInfo.Serials) != 1 {
				t.Fatalf("bad: expected CRL with length 1, got %d", len(crlInfo.Serials))
			}
			if _, ok := crlInfo.Serials["637101449987587619778072672905061040630001617053"]; !ok {
				t.Fatalf("bad: expected serial number not found in CRL")
			}
			return nil
		},
	}
}

func testAccStepDeleteCRL(t *testing.T, connState tls.ConnectionState) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.DeleteOperation,
		Path:      "crls/test",
		ConnState: &connState,
	}
}

func testAccStepLogin(t *testing.T, connState tls.ConnectionState) logicaltest.TestStep {
	return testAccStepLoginWithName(t, connState, "")
}

func testAccStepLoginWithName(t *testing.T, connState tls.ConnectionState, certName string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation:       logical.UpdateOperation,
		Path:            "login",
		Unauthenticated: true,
		ConnState:       &connState,
		Check: func(resp *logical.Response) error {
			if resp.Auth.TTL != 1000*time.Second {
				t.Fatalf("bad lease length: %#v", resp.Auth)
			}

			if certName != "" && resp.Auth.DisplayName != ("mnt-"+certName) {
				t.Fatalf("matched the wrong cert: %#v", resp.Auth.DisplayName)
			}

			fn := logicaltest.TestCheckAuth([]string{"default", "foo"})
			return fn(resp)
		},
		Data: map[string]interface{}{
			"name": certName,
		},
	}
}

func testAccStepLoginDefaultLease(t *testing.T, connState tls.ConnectionState) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation:       logical.UpdateOperation,
		Path:            "login",
		Unauthenticated: true,
		ConnState:       &connState,
		Check: func(resp *logical.Response) error {
			if resp.Auth.TTL != 1000*time.Second {
				t.Fatalf("bad lease length: %#v", resp.Auth)
			}

			fn := logicaltest.TestCheckAuth([]string{"default", "foo"})
			return fn(resp)
		},
	}
}

func testAccStepLoginInvalid(t *testing.T, connState tls.ConnectionState) logicaltest.TestStep {
	return testAccStepLoginWithNameInvalid(t, connState, "")
}

func testAccStepLoginWithNameInvalid(t *testing.T, connState tls.ConnectionState, certName string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation:       logical.UpdateOperation,
		Path:            "login",
		Unauthenticated: true,
		ConnState:       &connState,
		Check: func(resp *logical.Response) error {
			if resp.Auth != nil {
				return fmt.Errorf("should not be authorized: %#v", resp)
			}
			return nil
		},
		Data: map[string]interface{}{
			"name": certName,
		},
		ErrorOk: true,
	}
}

func testAccStepListCerts(
	t *testing.T, certs []string) []logicaltest.TestStep {
	return []logicaltest.TestStep{
		logicaltest.TestStep{
			Operation: logical.ListOperation,
			Path:      "certs",
			Check: func(resp *logical.Response) error {
				if resp == nil {
					return fmt.Errorf("nil response")
				}
				if resp.Data == nil {
					return fmt.Errorf("nil data")
				}
				if resp.Data["keys"] == interface{}(nil) {
					return fmt.Errorf("nil keys")
				}
				keys := resp.Data["keys"].([]string)
				if !reflect.DeepEqual(keys, certs) {
					return fmt.Errorf("mismatch: keys is %#v, certs is %#v", keys, certs)
				}
				return nil
			},
		}, logicaltest.TestStep{
			Operation: logical.ListOperation,
			Path:      "certs/",
			Check: func(resp *logical.Response) error {
				if resp == nil {
					return fmt.Errorf("nil response")
				}
				if resp.Data == nil {
					return fmt.Errorf("nil data")
				}
				if resp.Data["keys"] == interface{}(nil) {
					return fmt.Errorf("nil keys")
				}
				keys := resp.Data["keys"].([]string)
				if !reflect.DeepEqual(keys, certs) {
					return fmt.Errorf("mismatch: keys is %#v, certs is %#v", keys, certs)
				}

				return nil
			},
		},
	}
}

type allowed struct {
	names                string // allowed names in the certificate, looks at common, name, dns, email [depricated]
	common_names         string // allowed common names in the certificate
	dns                  string // allowed dns names in the SAN extension of the certificate
	emails               string // allowed email names in SAN extension of the certificate
	uris                 string // allowed uris in SAN extension of the certificate
	organizational_units string // allowed OUs in the certificate
	ext                  string // required extensions in the certificate
}

func testAccStepCert(
	t *testing.T, name string, cert []byte, policies string, testData allowed, expectError bool) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "certs/" + name,
		ErrorOk:   expectError,
		Data: map[string]interface{}{
			"certificate":                  string(cert),
			"policies":                     policies,
			"display_name":                 name,
			"allowed_names":                testData.names,
			"allowed_common_names":         testData.common_names,
			"allowed_dns_sans":             testData.dns,
			"allowed_email_sans":           testData.emails,
			"allowed_uri_sans":             testData.uris,
			"allowed_organizational_units": testData.organizational_units,
			"required_extensions":          testData.ext,
			"lease":                        1000,
		},
		Check: func(resp *logical.Response) error {
			if resp == nil && expectError {
				return fmt.Errorf("expected error but received nil")
			}
			return nil
		},
	}
}

func testAccStepCertLease(
	t *testing.T, name string, cert []byte, policies string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "certs/" + name,
		Data: map[string]interface{}{
			"certificate":  string(cert),
			"policies":     policies,
			"display_name": name,
			"lease":        1000,
		},
	}
}

func testAccStepCertTTL(
	t *testing.T, name string, cert []byte, policies string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "certs/" + name,
		Data: map[string]interface{}{
			"certificate":  string(cert),
			"policies":     policies,
			"display_name": name,
			"ttl":          "1000s",
		},
	}
}

func testAccStepCertMaxTTL(
	t *testing.T, name string, cert []byte, policies string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "certs/" + name,
		Data: map[string]interface{}{
			"certificate":  string(cert),
			"policies":     policies,
			"display_name": name,
			"ttl":          "1000s",
			"max_ttl":      "1200s",
		},
	}
}

func testAccStepCertNoLease(
	t *testing.T, name string, cert []byte, policies string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "certs/" + name,
		Data: map[string]interface{}{
			"certificate":  string(cert),
			"policies":     policies,
			"display_name": name,
		},
	}
}

func testConnState(certPath, keyPath, rootCertPath string) (tls.ConnectionState, error) {
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return tls.ConnectionState{}, err
	}
	rootConfig := &rootcerts.Config{
		CAFile: rootCertPath,
	}
	rootCAs, err := rootcerts.LoadCACerts(rootConfig)
	if err != nil {
		return tls.ConnectionState{}, err
	}
	listenConf := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		ClientAuth:         tls.RequestClientCert,
		InsecureSkipVerify: false,
		RootCAs:            rootCAs,
	}
	dialConf := listenConf.Clone()
	// start a server
	list, err := tls.Listen("tcp", "127.0.0.1:0", listenConf)
	if err != nil {
		return tls.ConnectionState{}, err
	}
	defer list.Close()

	// Accept connections.
	serverErrors := make(chan error, 1)
	connState := make(chan tls.ConnectionState)
	go func() {
		defer close(connState)
		serverConn, err := list.Accept()
		serverErrors <- err
		if err != nil {
			close(serverErrors)
			return
		}
		defer serverConn.Close()

		// Read the ping
		buf := make([]byte, 4)
		_, err = serverConn.Read(buf)
		if (err != nil) && (err != io.EOF) {
			serverErrors <- err
			close(serverErrors)
			return
		} else {
			// EOF is a reasonable error condition, so swallow it.
			serverErrors <- nil
		}
		close(serverErrors)
		connState <- serverConn.(*tls.Conn).ConnectionState()
	}()

	// Establish a connection from the client side and write a few bytes.
	clientErrors := make(chan error, 1)
	go func() {
		addr := list.Addr().String()
		conn, err := tls.Dial("tcp", addr, dialConf)
		clientErrors <- err
		if err != nil {
			close(clientErrors)
			return
		}
		defer conn.Close()

		// Write ping
		_, err = conn.Write([]byte("ping"))
		clientErrors <- err
		close(clientErrors)
	}()

	for err = range clientErrors {
		if err != nil {
			return tls.ConnectionState{}, fmt.Errorf("error in client goroutine:%v", err)
		}
	}

	for err = range serverErrors {
		if err != nil {
			return tls.ConnectionState{}, fmt.Errorf("error in server goroutine:%v", err)
		}
	}
	// Grab the current state
	return <-connState, nil
}

func Test_Renew(t *testing.T) {
	storage := &logical.InmemStorage{}

	lb, err := Factory(context.Background(), &logical.BackendConfig{
		System: &logical.StaticSystemView{
			DefaultLeaseTTLVal: 300 * time.Second,
			MaxLeaseTTLVal:     1800 * time.Second,
		},
		StorageView: storage,
	})
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	b := lb.(*backend)
	connState, err := testConnState("test-fixtures/keys/cert.pem",
		"test-fixtures/keys/key.pem", "test-fixtures/root/rootcacert.pem")
	if err != nil {
		t.Fatalf("error testing connection state: %v", err)
	}
	ca, err := ioutil.ReadFile("test-fixtures/root/rootcacert.pem")
	if err != nil {
		t.Fatal(err)
	}

	req := &logical.Request{
		Connection: &logical.Connection{
			ConnState: &connState,
		},
		Storage: storage,
		Auth:    &logical.Auth{},
	}

	fd := &framework.FieldData{
		Raw: map[string]interface{}{
			"name":        "test",
			"certificate": ca,
			"policies":    "foo,bar",
		},
		Schema: pathCerts(b).Fields,
	}

	resp, err := b.pathCertWrite(context.Background(), req, fd)
	if err != nil {
		t.Fatal(err)
	}

	empty_login_fd := &framework.FieldData{
		Raw:    map[string]interface{}{},
		Schema: pathLogin(b).Fields,
	}
	resp, err = b.pathLogin(context.Background(), req, empty_login_fd)
	if err != nil {
		t.Fatal(err)
	}
	if resp.IsError() {
		t.Fatalf("got error: %#v", *resp)
	}
	req.Auth.InternalData = resp.Auth.InternalData
	req.Auth.Metadata = resp.Auth.Metadata
	req.Auth.LeaseOptions = resp.Auth.LeaseOptions
	req.Auth.Policies = resp.Auth.Policies
	req.Auth.TokenPolicies = req.Auth.Policies
	req.Auth.Period = resp.Auth.Period

	// Normal renewal
	resp, err = b.pathLoginRenew(context.Background(), req, empty_login_fd)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("got nil response from renew")
	}
	if resp.IsError() {
		t.Fatalf("got error: %#v", *resp)
	}

	// Change the policies -- this should fail
	fd.Raw["policies"] = "zip,zap"
	resp, err = b.pathCertWrite(context.Background(), req, fd)
	if err != nil {
		t.Fatal(err)
	}

	resp, err = b.pathLoginRenew(context.Background(), req, empty_login_fd)
	if err == nil {
		t.Fatal("expected error")
	}

	// Put the policies back, this should be okay
	fd.Raw["policies"] = "bar,foo"
	resp, err = b.pathCertWrite(context.Background(), req, fd)
	if err != nil {
		t.Fatal(err)
	}

	resp, err = b.pathLoginRenew(context.Background(), req, empty_login_fd)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("got nil response from renew")
	}
	if resp.IsError() {
		t.Fatalf("got error: %#v", *resp)
	}

	// Add period value to cert entry
	period := 350 * time.Second
	fd.Raw["period"] = period.String()
	resp, err = b.pathCertWrite(context.Background(), req, fd)
	if err != nil {
		t.Fatal(err)
	}

	resp, err = b.pathLoginRenew(context.Background(), req, empty_login_fd)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("got nil response from renew")
	}
	if resp.IsError() {
		t.Fatalf("got error: %#v", *resp)
	}

	if resp.Auth.Period != period {
		t.Fatalf("expected a period value of %s in the response, got: %s", period, resp.Auth.Period)
	}

	// Delete CA, make sure we can't renew
	resp, err = b.pathCertDelete(context.Background(), req, fd)
	if err != nil {
		t.Fatal(err)
	}

	resp, err = b.pathLoginRenew(context.Background(), req, empty_login_fd)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("got nil response from renew")
	}
	if !resp.IsError() {
		t.Fatal("expected error")
	}
}
