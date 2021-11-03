package server

import (
	"bytes"
	"crypto/tls"
	"io"
	"net"
	"testing"

)

type testListenerConnFn func(net.Listener) (net.Conn, error)

func testListenerImpl(t *testing.T, ln net.Listener, connFn testListenerConnFn, certName string, expectedVersion uint16) {
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

	server := <-serverCh
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
		t.Fatalf("bad: %v", buf.String())
	}
}


func TestProfilingUnauthenticatedInFlightAccess(t *testing.T) {

	config, err := LoadConfigFile("./test-fixtures/profiling_unauth_in_flight_access.hcl")
	if err != nil {
		t.Fatalf("Error encountered when loading config %+v", err)
	}
	if !config.Listeners[0].Profiling.UnauthenticatedInFlightAccess {
		t.Fatalf("failed to read UnauthenticatedInFlightAccess")
	}
}