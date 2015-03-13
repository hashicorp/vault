package server

import (
	"net"
	"testing"
)

func TestTCPListener(t *testing.T) {
	ln, err := tcpListenerFactory(map[string]string{
		"address": "127.0.0.1:0",
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	connFn := func(lnReal net.Listener) (net.Conn, error) {
		return net.Dial("tcp", ln.Addr().String())
	}

	testListenerImpl(t, ln, connFn)
}
