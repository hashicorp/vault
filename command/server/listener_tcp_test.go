// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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

	"github.com/hashicorp/go-sockaddr"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/mitchellh/cli"
	"github.com/pires/go-proxyproto"
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

	testListenerImpl(t, ln, connFn, "", 0, "127.0.0.1", false)
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

	testListenerImpl(t, ln, connFn(true), "foo.example.com", 0, "127.0.0.1", false)

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

	testListenerImpl(t, ln, connFn(false), "foo.example.com", 0, "127.0.0.1", false)
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

	testListenerImpl(t, ln, connFn(true), "foo.example.com", tls.VersionTLS13, "127.0.0.1", false)

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

	testListenerImpl(t, ln, connFn(false), "foo.example.com", tls.VersionTLS13, "127.0.0.1", false)

	ln, _, _, err = tcpListenerFactory(&configutil.Listener{
		Address:               "127.0.0.1:0",
		TLSCertFile:           wd + "reload_foo.pem",
		TLSKeyFile:            wd + "reload_foo.key",
		TLSDisableClientCerts: true,
		TLSClientCAFile:       wd + "reload_ca.pem",
		TLSMaxVersion:         "tls12",
	}, nil, cli.NewMockUi())
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	testListenerImpl(t, ln, connFn(false), "foo.example.com", tls.VersionTLS12, "127.0.0.1", false)
}

func TestTCPListener_proxyProtocol(t *testing.T) {
	for name, tc := range map[string]struct {
		Behavior       string
		Header         *proxyproto.Header
		AuthorizedAddr string
		ExpectedAddr   string
		ExpectError    bool
	}{
		"none-no-header": {
			Behavior:     "",
			ExpectedAddr: "127.0.0.1",
			Header:       nil,
		},
		"none-v1": {
			Behavior:     "",
			ExpectedAddr: "127.0.0.1",
			ExpectError:  true,
			Header: &proxyproto.Header{
				Version:           1,
				Command:           proxyproto.PROXY,
				TransportProtocol: proxyproto.TCPv4,
				SourceAddr: &net.TCPAddr{
					IP:   net.ParseIP("10.1.1.1"),
					Port: 1000,
				},
				DestinationAddr: &net.TCPAddr{
					IP:   net.ParseIP("20.2.2.2"),
					Port: 2000,
				},
			},
		},
		"none-v2": {
			Behavior:     "",
			ExpectedAddr: "127.0.0.1",
			ExpectError:  true,
			Header: &proxyproto.Header{
				Version:           2,
				Command:           proxyproto.PROXY,
				TransportProtocol: proxyproto.TCPv4,
				SourceAddr: &net.TCPAddr{
					IP:   net.ParseIP("10.1.1.1"),
					Port: 1000,
				},
				DestinationAddr: &net.TCPAddr{
					IP:   net.ParseIP("20.2.2.2"),
					Port: 2000,
				},
			},
		},

		// use_always makes it possible to send the PROXY header but does not
		// require it
		"use_always-no-header": {
			Behavior:     "use_always",
			ExpectedAddr: "127.0.0.1",
			Header:       nil,
		},

		"use_always-header-v1": {
			Behavior:     "use_always",
			ExpectedAddr: "10.1.1.1",
			Header: &proxyproto.Header{
				Version:           1,
				Command:           proxyproto.PROXY,
				TransportProtocol: proxyproto.TCPv4,
				SourceAddr: &net.TCPAddr{
					IP:   net.ParseIP("10.1.1.1"),
					Port: 1000,
				},
				DestinationAddr: &net.TCPAddr{
					IP:   net.ParseIP("20.2.2.2"),
					Port: 2000,
				},
			},
		},
		"use_always-header-v1-unknown": {
			Behavior:     "use_always",
			ExpectedAddr: "127.0.0.1",
			Header: &proxyproto.Header{
				Version:           1,
				Command:           proxyproto.PROXY,
				TransportProtocol: proxyproto.UNSPEC,
			},
		},
		"use_always-header-v2": {
			Behavior:     "use_always",
			ExpectedAddr: "10.1.1.1",
			Header: &proxyproto.Header{
				Version:           2,
				Command:           proxyproto.PROXY,
				TransportProtocol: proxyproto.TCPv4,
				SourceAddr: &net.TCPAddr{
					IP:   net.ParseIP("10.1.1.1"),
					Port: 1000,
				},
				DestinationAddr: &net.TCPAddr{
					IP:   net.ParseIP("20.2.2.2"),
					Port: 2000,
				},
			},
		},
		"use_always-header-v2-unknown": {
			Behavior:     "use_always",
			ExpectedAddr: "127.0.0.1",
			Header: &proxyproto.Header{
				Version:           2,
				Command:           proxyproto.LOCAL,
				TransportProtocol: proxyproto.UNSPEC,
			},
		},
		"allow_authorized-no-header-in": {
			Behavior:       "allow_authorized",
			AuthorizedAddr: "127.0.0.1/32",
			ExpectedAddr:   "127.0.0.1",
		},
		"allow_authorized-no-header-not-in": {
			Behavior:       "allow_authorized",
			AuthorizedAddr: "10.0.0.1/32",
			ExpectedAddr:   "127.0.0.1",
		},
		"allow_authorized-v1-in": {
			Behavior:       "allow_authorized",
			AuthorizedAddr: "127.0.0.1/32",
			ExpectedAddr:   "10.1.1.1",
			Header: &proxyproto.Header{
				Version:           1,
				Command:           proxyproto.PROXY,
				TransportProtocol: proxyproto.TCPv4,
				SourceAddr: &net.TCPAddr{
					IP:   net.ParseIP("10.1.1.1"),
					Port: 1000,
				},
				DestinationAddr: &net.TCPAddr{
					IP:   net.ParseIP("20.2.2.2"),
					Port: 2000,
				},
			},
		},

		// allow_authorized still accepts the PROXY header when not in the
		// authorized addresses but discards it silently
		"allow_authorized-v1-not-in": {
			Behavior:       "allow_authorized",
			AuthorizedAddr: "10.0.0.1/32",
			ExpectedAddr:   "127.0.0.1",
			Header: &proxyproto.Header{
				Version:           1,
				Command:           proxyproto.PROXY,
				TransportProtocol: proxyproto.TCPv4,
				SourceAddr: &net.TCPAddr{
					IP:   net.ParseIP("10.1.1.1"),
					Port: 1000,
				},
				DestinationAddr: &net.TCPAddr{
					IP:   net.ParseIP("20.2.2.2"),
					Port: 2000,
				},
			},
		},

		"deny_unauthorized-no-header-in": {
			Behavior:       "deny_unauthorized",
			AuthorizedAddr: "127.0.0.1/32",
			ExpectedAddr:   "127.0.0.1",
		},
		"deny_unauthorized-no-header-not-in": {
			Behavior:       "deny_unauthorized",
			AuthorizedAddr: "10.0.0.1/32",
			ExpectedAddr:   "127.0.0.1",
			ExpectError:    true,
		},
		"deny_unauthorized-v1-in": {
			Behavior:       "deny_unauthorized",
			AuthorizedAddr: "127.0.0.1/32",
			ExpectedAddr:   "10.1.1.1",
			Header: &proxyproto.Header{
				Version:           1,
				Command:           proxyproto.PROXY,
				TransportProtocol: proxyproto.TCPv4,
				SourceAddr: &net.TCPAddr{
					IP:   net.ParseIP("10.1.1.1"),
					Port: 1000,
				},
				DestinationAddr: &net.TCPAddr{
					IP:   net.ParseIP("20.2.2.2"),
					Port: 2000,
				},
			},
		},
		"deny_unauthorized-v1-not-in": {
			Behavior:       "deny_unauthorized",
			AuthorizedAddr: "10.0.0.1/32",
			ExpectedAddr:   "127.0.0.1",
			ExpectError:    true,
			Header: &proxyproto.Header{
				Version:           1,
				Command:           proxyproto.PROXY,
				TransportProtocol: proxyproto.TCPv4,
				SourceAddr: &net.TCPAddr{
					IP:   net.ParseIP("10.1.1.1"),
					Port: 1000,
				},
				DestinationAddr: &net.TCPAddr{
					IP:   net.ParseIP("20.2.2.2"),
					Port: 2000,
				},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			proxyProtocolAuthorizedAddrs := []*sockaddr.SockAddrMarshaler{}
			if tc.AuthorizedAddr != "" {
				sockAddr, err := sockaddr.NewSockAddr(tc.AuthorizedAddr)
				if err != nil {
					t.Fatal(err)
				}
				proxyProtocolAuthorizedAddrs = append(
					proxyProtocolAuthorizedAddrs,
					&sockaddr.SockAddrMarshaler{SockAddr: sockAddr},
				)
			}

			ln, _, _, err := tcpListenerFactory(&configutil.Listener{
				Address:                      "127.0.0.1:0",
				TLSDisable:                   true,
				ProxyProtocolBehavior:        tc.Behavior,
				ProxyProtocolAuthorizedAddrs: proxyProtocolAuthorizedAddrs,
			}, nil, cli.NewMockUi())
			if err != nil {
				t.Fatalf("err: %s", err)
			}

			connFn := func(lnReal net.Listener) (net.Conn, error) {
				conn, err := net.Dial("tcp", ln.Addr().String())
				if err != nil {
					return nil, err
				}

				if tc.Header != nil {
					_, err = tc.Header.WriteTo(conn)
				}
				return conn, err
			}

			testListenerImpl(t, ln, connFn, "", 0, tc.ExpectedAddr, tc.ExpectError)
		})
	}
}
