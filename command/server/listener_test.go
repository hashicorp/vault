package server

import (
	"bytes"
	"io"
	"net"
	"testing"
	"time"
)

type testListenerConnFn func(net.Listener) (net.Conn, error)

func testListenerImpl(t *testing.T, ln net.Listener, connFn testListenerConnFn) {
	serverCh := make(chan net.Conn, 1)
	go func() {
		server, err := ln.Accept()
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		serverCh <- server
	}()

	client, err := connFn(ln)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	server := <-serverCh
	defer client.Close()
	defer server.Close()

	var buf bytes.Buffer
	go io.Copy(&buf, server)

	if _, err := client.Write([]byte("foo")); err != nil {
		t.Fatalf("err: %s", err)
	}

	time.Sleep(100 * time.Millisecond)

	if buf.String() != "foo" {
		t.Fatalf("bad: %v", buf.String())
	}
}
