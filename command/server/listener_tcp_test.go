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

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
)

func TestTCPListener(t *testing.T) {
	ln, _, _, err := tcpListenerFactory(map[string]string{
		"address":     "127.0.0.1:0",
		"tls_disable": "1",
	})
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
	inBytes, _ := ioutil.ReadFile(wd + "reload_foo.pem")
	ioutil.WriteFile(td+"reload_curr.pem", inBytes, 0777)
	inBytes, _ = ioutil.ReadFile(wd + "reload_foo.key")
	ioutil.WriteFile(td+"reload_curr.key", inBytes, 0777)
	inBytes, _ = ioutil.ReadFile(wd + "reload_ca.pem")
	certPool := x509.NewCertPool()
	ok := certPool.AppendCertsFromPEM(inBytes)
	if !ok {
		t.Fatal("not ok when appending CA cert")
	}

	ln, _, reloadFunc, err := tcpListenerFactory(map[string]string{
		"address":       "127.0.0.1:0",
		"tls_cert_file": td + "reload_curr.pem",
		"tls_key_file":  td + "reload_curr.key",
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	core, _, root := vault.TestCoreUnsealed(t)
	bc := &logical.BackendConfig{
		Logger: nil,
		System: logical.StaticSystemView{
			DefaultLeaseTTLVal: time.Hour * 24,
			MaxLeaseTTLVal:     time.Hour * 24 * 30,
		},
	}

	vault.NewSystemBackend(core, bc)
	core.AddReloadFunc("listentest", reloadFunc)

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

	inBytes, _ = ioutil.ReadFile(wd + "reload_bar.pem")
	ioutil.WriteFile(td+"reload_curr.pem", inBytes, 0777)
	inBytes, _ = ioutil.ReadFile(wd + "reload_bar.key")
	ioutil.WriteFile(td+"reload_curr.key", inBytes, 0777)

	req := logical.TestRequest(t, logical.UpdateOperation, "sys/reload")
	req.ClientToken = root
	resp, err := core.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if resp != nil {
		t.Fatal("expected nil response")
	}

	testListenerImpl(t, ln, connFn, "bar.example.com")
}
