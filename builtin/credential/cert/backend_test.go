package cert

import (
	"crypto/rand"
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

	"github.com/hashicorp/go-rootcerts"
	"github.com/hashicorp/vault/helper/certutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	logicaltest "github.com/hashicorp/vault/logical/testing"
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

	b, err := Factory(config)
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

	resp, err = b.HandleRequest(certReq)
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
	resp, err = b.HandleRequest(loginReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	// Wait until the certificate expires
	time.Sleep(5 * time.Second)

	// Login attempt after certificate expiry should fail
	resp, err = b.HandleRequest(loginReq)
	if err == nil {
		t.Fatalf("expected error due to expired certificate")
	}
}

func TestBackend_RegisteredNonCA_CRL(t *testing.T) {
	config := logical.TestBackendConfig()
	storage := &logical.InmemStorage{}
	config.StorageView = storage

	b, err := Factory(config)
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

	resp, err := b.HandleRequest(certReq)
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
	resp, err = b.HandleRequest(loginReq)
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
	resp, err = b.HandleRequest(crlReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	// Attempt login with the same connection state but with the CRL registered
	resp, err = b.HandleRequest(loginReq)
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

	b, err := Factory(config)
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

	resp, err := b.HandleRequest(certReq)
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
	resp, err = b.HandleRequest(loginReq)
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
	resp, err = b.HandleRequest(loginReq)
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
	resp, err = b.HandleRequest(crlReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	// Attempt login with the revoked certificate.
	resp, err = b.HandleRequest(loginReq)
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
	resp, err = b.HandleRequest(certReq)
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
	resp, err = b.HandleRequest(loginReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	// Register a CRL containing the root CA certificate used above.
	rootCRL, err := ioutil.ReadFile(testRootCertCRL)
	if err != nil {
		t.Fatal(err)
	}
	crlData["crl"] = rootCRL
	resp, err = b.HandleRequest(crlReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	// Attempt login with the same connection state but with the CRL registered
	resp, err = b.HandleRequest(loginReq)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || !resp.IsError() {
		t.Fatalf("expected failure due to revoked certificate")
	}
}

func testFactory(t *testing.T) logical.Backend {
	b, err := Factory(&logical.BackendConfig{
		System: &logical.StaticSystemView{
			DefaultLeaseTTLVal: 300 * time.Second,
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
		Backend: testFactory(t),
		Steps: []logicaltest.TestStep{
			testAccStepCert(t, "aaa", ca1, "foo", "", false),
			testAccStepCert(t, "bbb", ca2, "foo", "", false),
			testAccStepCert(t, "ccc", ca3, "foo", "", true),
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
		Backend: testFactory(t),
		Steps: []logicaltest.TestStep{
			testAccStepCert(t, "web", ca, "foo", "", false),
			testAccStepLogin(t, connState),
			testAccStepCertLease(t, "web", ca, "foo"),
			testAccStepCertTTL(t, "web", ca, "foo"),
			testAccStepLogin(t, connState),
			testAccStepCertNoLease(t, "web", ca, "foo"),
			testAccStepLoginDefaultLease(t, connState),
			testAccStepCert(t, "web", ca, "foo", "*.example.com", false),
			testAccStepLogin(t, connState),
			testAccStepCert(t, "web", ca, "foo", "*.invalid.com", false),
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
		Backend: testFactory(t),
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
		Backend: testFactory(t),
		Steps: []logicaltest.TestStep{
			testAccStepCert(t, "web", ca, "foo", "", false),
			testAccStepLogin(t, connState),
			testAccStepCert(t, "web", ca, "foo", "example.com", false),
			testAccStepLogin(t, connState),
			testAccStepCert(t, "web", ca, "foo", "invalid", false),
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
		Backend: testFactory(t),
		Steps: []logicaltest.TestStep{
			testAccStepCert(t, "1unconstrained", ca, "foo", "", false),
			testAccStepCert(t, "2matching", ca, "foo", "*.example.com,whatever", false),
			testAccStepCert(t, "3invalid", ca, "foo", "invalid", false),
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
		Backend: testFactory(t),
		Steps: []logicaltest.TestStep{
			testAccStepLoginInvalid(t, connState),
		},
	})
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
			if resp.Auth.TTL != 300*time.Second {
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

func testAccStepCert(
	t *testing.T, name string, cert []byte, policies string, allowedNames string, expectError bool) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "certs/" + name,
		ErrorOk:   expectError,
		Data: map[string]interface{}{
			"certificate":   string(cert),
			"policies":      policies,
			"display_name":  name,
			"allowed_names": allowedNames,
			"lease":         1000,
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
	dialConf := new(tls.Config)
	*dialConf = *listenConf
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

	lb, err := Factory(&logical.BackendConfig{
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

	resp, err := b.pathCertWrite(req, fd)
	if err != nil {
		t.Fatal(err)
	}

	empty_login_fd := &framework.FieldData{
		Raw:    map[string]interface{}{},
		Schema: pathLogin(b).Fields,
	}
	resp, err = b.pathLogin(req, empty_login_fd)
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
	req.Auth.IssueTime = time.Now()

	// Normal renewal
	resp, err = b.pathLoginRenew(req, empty_login_fd)
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
	resp, err = b.pathCertWrite(req, fd)
	if err != nil {
		t.Fatal(err)
	}

	resp, err = b.pathLoginRenew(req, empty_login_fd)
	if err == nil {
		t.Fatal("expected error")
	}

	// Put the policies back, this shold be okay
	fd.Raw["policies"] = "bar,foo"
	resp, err = b.pathCertWrite(req, fd)
	if err != nil {
		t.Fatal(err)
	}

	resp, err = b.pathLoginRenew(req, empty_login_fd)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("got nil response from renew")
	}
	if resp.IsError() {
		t.Fatalf("got error: %#v", *resp)
	}

	// Delete CA, make sure we can't renew
	resp, err = b.pathCertDelete(req, fd)
	if err != nil {
		t.Fatal(err)
	}

	resp, err = b.pathLoginRenew(req, empty_login_fd)
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
