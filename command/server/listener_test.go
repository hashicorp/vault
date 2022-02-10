package server

import (
	"bytes"
	"crypto/tls"
	"io"
	"net"
	"testing"
	"time"
)

type testListenerConnFn func(net.Listener) (net.Conn, error)

func testListenerImpl(t *testing.T, ln net.Listener, connFn testListenerConnFn, certName string, expectedVersion uint16, expectedAddr string, expectedPort int, expectError bool) {
	serverCh := make(chan net.Conn, 1)
	go func() {
		server, err := ln.Accept()
		if err != nil {
			t.Errorf("err: %s", err)
			return
		}
		if certName != "" {
			tlsConn := server.(*tls.Conn)
			tlsConn.Handshake()
		}
		serverCh <- server

		addr := server.RemoteAddr().(*net.TCPAddr)
		if addr.IP.String() != expectedAddr {
			t.Errorf("bad: %v", addr)
		}
		if expectedPort != 0 && addr.Port != expectedPort {
			t.Errorf("bad: %v", addr)
		}
	}()

	client, err := connFn(ln)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if certName != "" {
		tlsConn := client.(*tls.Conn)
		if expectedVersion != 0 && tlsConn.ConnectionState().Version != expectedVersion {
			t.Fatalf("expected version %d, got %d", expectedVersion, tlsConn.ConnectionState().Version)
		}
		if len(tlsConn.ConnectionState().PeerCertificates) != 1 {
			t.Fatalf("err: number of certs too long")
		}
		peerName := tlsConn.ConnectionState().PeerCertificates[0].Subject.CommonName
		if peerName != certName {
			t.Fatalf("err: bad cert name %s, expected %s", peerName, certName)
		}
	}

	var server net.Conn
	ticker := time.NewTicker(10 * time.Second)
	select {
	case <-ticker.C:
		break
	case server = <-serverCh:
	}

	if server == nil {
		if client != nil {
			client.Close()
		}
		if !expectError {
			// Something failed already so we abort the test early
			t.Fatal("aborting test because the server did not accept the connection")
		}
		return
	}
	defer client.Close()
	defer server.Close()

	var buf bytes.Buffer
	copyCh := make(chan struct{})
	go func() {
		io.Copy(&buf, server)
		close(copyCh)
	}()

	if _, err := client.Write([]byte("foo")); err != nil {
		t.Fatalf("err: %s", err)
	}

	client.Close()

	<-copyCh
	if buf.String() != "foo" {
		if !expectError {
			t.Fatalf("bad: %v", buf.String())
		}
		return
	}
}

func TestProfilingUnauthenticatedInFlightAccess(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/unauth_in_flight_access.hcl")
	if err != nil {
		t.Fatalf("Error encountered when loading config %+v", err)
	}
	if !config.Listeners[0].InFlightRequestLogging.UnauthenticatedInFlightAccess {
		t.Fatalf("failed to read UnauthenticatedInFlightAccess")
	}
}
