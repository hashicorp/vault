package cert

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/logical"
	logicaltest "github.com/hashicorp/vault/logical/testing"
	"github.com/mitchellh/mapstructure"
)

func testFactory(t *testing.T) logical.Backend {
	b, err := Factory(&logical.BackendConfig{
		System: &logical.StaticSystemView{
			DefaultLeaseTTLVal: 300 * time.Second,
			MaxLeaseTTLVal:     1800 * time.Second,
		},
		StorageView: &logical.InmemStorage{},
	})
	if err != nil {
		t.Fatal("error: %s", err)
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
			testAccStepCert(t, "aaa", ca1, "foo", false),
			testAccStepCert(t, "bbb", ca2, "foo", false),
			testAccStepCert(t, "ccc", ca3, "foo", true),
		},
	}
	tc.Steps = append(tc.Steps, testAccStepListCerts(t, []string{"aaa", "bbb"})...)
	logicaltest.Test(t, tc)
}

// Test a client trusted by a CA
func TestBackend_basic_CA(t *testing.T) {
	connState := testConnState(t, "test-fixtures/keys/cert.pem",
		"test-fixtures/keys/key.pem", "test-fixtures/root/rootcacert.pem")
	ca, err := ioutil.ReadFile("test-fixtures/root/rootcacert.pem")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	logicaltest.Test(t, logicaltest.TestCase{
		Backend: testFactory(t),
		Steps: []logicaltest.TestStep{
			testAccStepCert(t, "web", ca, "foo", false),
			testAccStepLogin(t, connState),
			testAccStepCertLease(t, "web", ca, "foo"),
			testAccStepCertTTL(t, "web", ca, "foo"),
			testAccStepLogin(t, connState),
			testAccStepCertNoLease(t, "web", ca, "foo"),
			testAccStepLoginDefaultLease(t, connState),
		},
	})
}

// Test CRL behavior
func TestBackend_CRLs(t *testing.T) {
	connState := testConnState(t, "test-fixtures/keys/cert.pem",
		"test-fixtures/keys/key.pem", "test-fixtures/root/rootcacert.pem")
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

// Test a self-signed client that is trusted
func TestBackend_basic_singleCert(t *testing.T) {
	connState := testConnState(t, "test-fixtures/keys/cert.pem",
		"test-fixtures/keys/key.pem", "test-fixtures/root/rootcacert.pem")
	ca, err := ioutil.ReadFile("test-fixtures/root/rootcacert.pem")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	logicaltest.Test(t, logicaltest.TestCase{
		Backend: testFactory(t),
		Steps: []logicaltest.TestStep{
			testAccStepCert(t, "web", ca, "foo", false),
			testAccStepLogin(t, connState),
		},
	})
}

// Test an untrusted self-signed client
func TestBackend_untrusted(t *testing.T) {
	connState := testConnState(t, "test-fixtures/keys/cert.pem",
		"test-fixtures/keys/key.pem", "test-fixtures/root/rootcacert.pem")
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

				return nil
			},
		},
	}
}

func testAccStepCert(
	t *testing.T, name string, cert []byte, policies string, expectError bool) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "certs/" + name,
		ErrorOk:   expectError,
		Data: map[string]interface{}{
			"certificate":  string(cert),
			"policies":     policies,
			"display_name": name,
			"lease":        1000,
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

func testConnState(t *testing.T, certPath, keyPath, rootCertPath string) tls.ConnectionState {
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	rootCAs, err := api.LoadCACert(rootCertPath)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	listenConf := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		ClientAuth:         tls.RequestClientCert,
		InsecureSkipVerify: false,
		RootCAs:            rootCAs,
	}
	dialConf := new(tls.Config)
	*dialConf = *listenConf
	list, err := tls.Listen("tcp", "127.0.0.1:0", listenConf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	defer list.Close()

	go func() {
		addr := list.Addr().String()
		conn, err := tls.Dial("tcp", addr, dialConf)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		defer conn.Close()

		// Write ping
		conn.Write([]byte("ping"))
	}()

	serverConn, err := list.Accept()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	defer serverConn.Close()

	// Read the pign
	buf := make([]byte, 4)
	serverConn.Read(buf)

	// Grab the current state
	connState := serverConn.(*tls.Conn).ConnectionState()
	return connState
}
