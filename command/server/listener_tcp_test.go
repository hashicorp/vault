package server

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"testing"
	"time"
)

func TestTCPListener(t *testing.T) {
	ln, _, _, err := tcpListenerFactory(map[string]string{
		"address":     "127.0.0.1:0",
		"tls_disable": "1",
	}, nil)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	connFn := func(lnReal net.Listener) (net.Conn, error) {
		return net.Dial("tcp", ln.Addr().String())
	}

	testListenerImpl(t, ln, connFn, "")
}

// TestTCPListener_tls tests both TLS generally and also the reload capability
// of core, system backend, and the listener logic
func TestTCPListener_tls(t *testing.T) {
	wd, _ := os.Getwd()
	wd += "/test-fixtures/reload/"

	td, err := ioutil.TempDir("", fmt.Sprintf("vault-test-%d", rand.New(rand.NewSource(time.Now().Unix())).Int63))
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(td)

	// Setup initial certs
	inBytes, _ := ioutil.ReadFile(wd + "reload_ca.pem")
	certPool := x509.NewCertPool()
	ok := certPool.AppendCertsFromPEM(inBytes)
	if !ok {
		t.Fatal("not ok when appending CA cert")
	}

	ln, _, _, err := tcpListenerFactory(map[string]string{
		"address":       "127.0.0.1:0",
		"tls_cert_file": wd + "reload_foo.pem",
		"tls_key_file":  wd + "reload_foo.key",
	}, nil)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	connFn := func(lnReal net.Listener) (net.Conn, error) {
		conn, err := tls.Dial("tcp", ln.Addr().String(), &tls.Config{
			RootCAs: certPool,
		})
		if err != nil {
			return nil, err
		}
		if err = conn.Handshake(); err != nil {
			return nil, err
		}
		return conn, nil
	}

	testListenerImpl(t, ln, connFn, "foo.example.com")
}
