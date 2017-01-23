package server

import (
	"bytes"
	"crypto/tls"
	"io"
	"net"
	"testing"
)

type testListenerConnFn func(net.Listener) (net.Conn, error)

func testListenerImpl(t *testing.T, ln net.Listener, connFn testListenerConnFn, certName string) {
	serverCh := make(chan net.Conn, 1)
	go func() {
		server, err := ln.Accept()
		if err != nil {
			t.Fatalf("err: %s", err)
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
