package server

import (
	"net"
	"testing"
)

func TestTCPListener(t *testing.T) {
	ln, _, err := tcpListenerFactory(map[string]string{
		"address":     "127.0.0.1:0",
		"tls_disable": "1",
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	connFn := func(lnReal net.Listener) (net.Conn, error) {
		return net.Dial("tcp", ln.Addr().String())
	}

	testListenerImpl(t, ln, connFn)
}

func TestTCPListener_tls(t *testing.T) {
	// TODO
	t.Skip()

	ln, _, err := tcpListenerFactory(map[string]string{
		"address":     "127.0.0.1:0",
		"tls_disable": "1",
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	connFn := func(lnReal net.Listener) (net.Conn, error) {
		return net.Dial("tcp", ln.Addr().String())
	}

	testListenerImpl(t, ln, connFn)
}
