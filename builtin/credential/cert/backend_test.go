package cert

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/hashicorp/vault/logical"
	logicaltest "github.com/hashicorp/vault/logical/testing"
)

// Test a client trusted by a CA
func TestBackend_basic_CA(t *testing.T) {
	connState := testConnState(t, "../../../test/key/ourdomain.cer",
		"../../../test/key/ourdomain.key")
	ca, err := ioutil.ReadFile("../../../test/ca/root.cer")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	logicaltest.Test(t, logicaltest.TestCase{
		Backend: Backend(),
		Steps: []logicaltest.TestStep{
			testAccStepCert(t, "web", ca, "foo"),
			testAccStepLogin(t, connState),
		},
	})
}

// Test a self-signed client that is trusted
func TestBackend_basic_singleCert(t *testing.T) {
	connState := testConnState(t, "../../../test/unsigned/cert.pem",
		"../../../test/unsigned/key.pem")
	ca, err := ioutil.ReadFile("../../../test/unsigned/cert.pem")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	logicaltest.Test(t, logicaltest.TestCase{
		Backend: Backend(),
		Steps: []logicaltest.TestStep{
			testAccStepCert(t, "web", ca, "foo"),
			testAccStepLogin(t, connState),
		},
	})
}

// Test an untrusted self-signed client
func TestBackend_untrusted(t *testing.T) {
	connState := testConnState(t, "../../../test/unsigned/cert.pem",
		"../../../test/unsigned/key.pem")
	logicaltest.Test(t, logicaltest.TestCase{
		Backend: Backend(),
		Steps: []logicaltest.TestStep{
			testAccStepLoginInvalid(t, connState),
		},
	})
}

func testAccStepLogin(t *testing.T, connState tls.ConnectionState) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation:       logical.WriteOperation,
		Path:            "login",
		Unauthenticated: true,
		ConnState:       &connState,
		Check: func(resp *logical.Response) error {
			if resp.Auth.Lease != 1000*time.Second {
				t.Fatalf("bad lease length: %#v", resp.Auth)
			}

			fn := logicaltest.TestCheckAuth([]string{"foo"})
			return fn(resp)
		},
	}
}

func testAccStepLoginInvalid(t *testing.T, connState tls.ConnectionState) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation:       logical.WriteOperation,
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

func testAccStepCert(
	t *testing.T, name string, cert []byte, policies string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.WriteOperation,
		Path:      "certs/" + name,
		Data: map[string]interface{}{
			"certificate":  string(cert),
			"policies":     policies,
			"display_name": name,
			"lease":        1000,
		},
	}
}

func testConnState(t *testing.T, certPath, keyPath string) tls.ConnectionState {
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	conf := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		ClientAuth:         tls.RequestClientCert,
		InsecureSkipVerify: true,
	}
	list, err := tls.Listen("tcp", "127.0.0.1:0", conf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	defer list.Close()

	go func() {
		addr := list.Addr().String()
		conn, err := tls.Dial("tcp", addr, conf)
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
