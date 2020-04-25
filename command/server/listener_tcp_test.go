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

	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/mitchellh/cli"
)

func TestTCPListener(t *testing.T) {
	ln, _, _, err := tcpListenerFactory(&configutil.Listener{
		Address:    "127.0.0.1:0",
		TLSDisable: true,
	}, nil, cli.NewMockUi())
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	connFn := func(lnReal net.Listener) (net.Conn, error) {
		return net.Dial("tcp", ln.Addr().String())
	}

	testListenerImpl(t, ln, connFn, "")
}

// TestTCPListener_tls tests TLS generally
func TestTCPListener_tls(t *testing.T) {
	wd, _ := os.Getwd()
	wd += "/test-fixtures/reload/"

	td, err := ioutil.TempDir("", fmt.Sprintf("vault-test-%d", rand.New(rand.NewSource(time.Now().Unix())).Int63()))
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

	ln, _, _, err := tcpListenerFactory(&configutil.Listener{
		Address:                       "127.0.0.1:0",
		TLSCertFile:                   wd + "reload_foo.pem",
		TLSKeyFile:                    wd + "reload_foo.key",
		TLSRequireAndVerifyClientCert: true,
		TLSClientCAFile:               wd + "reload_ca.pem",
	}, nil, cli.NewMockUi())
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	cwd, _ := os.Getwd()

	clientCert, _ := tls.LoadX509KeyPair(
		cwd+"/test-fixtures/reload/reload_foo.pem",
		cwd+"/test-fixtures/reload/reload_foo.key")

	connFn := func(clientCerts bool) func(net.Listener) (net.Conn, error) {
		return func(lnReal net.Listener) (net.Conn, error) {
			conf := &tls.Config{
				RootCAs: certPool,
			}
			if clientCerts {
				conf.Certificates = []tls.Certificate{clientCert}
			}
			conn, err := tls.Dial("tcp", ln.Addr().String(), conf)

			if err != nil {
				return nil, err
			}
			if err = conn.Handshake(); err != nil {
				return nil, err
			}
			return conn, nil
		}
	}

	testListenerImpl(t, ln, connFn(true), "foo.example.com")

	ln, _, _, err = tcpListenerFactory(&configutil.Listener{
		Address:                       "127.0.0.1:0",
		TLSCertFile:                   wd + "reload_foo.pem",
		TLSKeyFile:                    wd + "reload_foo.key",
		TLSRequireAndVerifyClientCert: true,
		TLSDisableClientCerts:         true,
		TLSClientCAFile:               wd + "reload_ca.pem",
	}, nil, cli.NewMockUi())
	if err == nil {
		t.Fatal("expected error due to mutually exclusive client cert options")
	}

	ln, _, _, err = tcpListenerFactory(&configutil.Listener{
		Address:               "127.0.0.1:0",
		TLSCertFile:           wd + "reload_foo.pem",
		TLSKeyFile:            wd + "reload_foo.key",
		TLSDisableClientCerts: true,
		TLSClientCAFile:       wd + "reload_ca.pem",
	}, nil, cli.NewMockUi())
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	testListenerImpl(t, ln, connFn(false), "foo.example.com")
}

func TestTCPListener_tls13(t *testing.T) {
	wd, _ := os.Getwd()
	wd += "/test-fixtures/reload/"

	td, err := ioutil.TempDir("", fmt.Sprintf("vault-test-%d", rand.New(rand.NewSource(time.Now().Unix())).Int63()))
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

	ln, _, _, err := tcpListenerFactory(&configutil.Listener{
		Address:                       "127.0.0.1:0",
		TLSCertFile:                   wd + "reload_foo.pem",
		TLSKeyFile:                    wd + "reload_foo.key",
		TLSRequireAndVerifyClientCert: true,
		TLSClientCAFile:               wd + "reload_ca.pem",
		TLSMinVersion:                 "tls13",
	}, nil, cli.NewMockUi())
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	cwd, _ := os.Getwd()

	clientCert, _ := tls.LoadX509KeyPair(
		cwd+"/test-fixtures/reload/reload_foo.pem",
		cwd+"/test-fixtures/reload/reload_foo.key")

	connFn := func(clientCerts bool) func(net.Listener) (net.Conn, error) {
		return func(lnReal net.Listener) (net.Conn, error) {
			conf := &tls.Config{
				RootCAs: certPool,
			}
			if clientCerts {
				conf.Certificates = []tls.Certificate{clientCert}
			}
			conn, err := tls.Dial("tcp", ln.Addr().String(), conf)

			if err != nil {
				return nil, err
			}
			if err = conn.Handshake(); err != nil {
				return nil, err
			}
			return conn, nil
		}
	}

	testListenerImpl(t, ln, connFn(true), "foo.example.com")

	ln, _, _, err = tcpListenerFactory(&configutil.Listener{
		Address:                       "127.0.0.1:0",
		TLSCertFile:                   wd + "reload_foo.pem",
		TLSKeyFile:                    wd + "reload_foo.key",
		TLSRequireAndVerifyClientCert: true,
		TLSDisableClientCerts:         true,
		TLSClientCAFile:               wd + "reload_ca.pem",
		TLSMinVersion:                 "tls13",
	}, nil, cli.NewMockUi())
	if err == nil {
		t.Fatal("expected error due to mutually exclusive client cert options")
	}

	ln, _, _, err = tcpListenerFactory(&configutil.Listener{
		Address:               "127.0.0.1:0",
		TLSCertFile:           wd + "reload_foo.pem",
		TLSKeyFile:            wd + "reload_foo.key",
		TLSDisableClientCerts: true,
		TLSClientCAFile:       wd + "reload_ca.pem",
		TLSMinVersion:         "tls13",
	}, nil, cli.NewMockUi())
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	testListenerImpl(t, ln, connFn(false), "foo.example.com")
}
