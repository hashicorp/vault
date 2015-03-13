package server

import (
	"fmt"
	"net"
	"time"
)

func tcpListenerFactory(config map[string]string) (net.Listener, error) {
	addr, ok := config["address"]
	if !ok {
		return nil, fmt.Errorf("'address' must be set")
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	ln = tcpKeepAliveListener{ln.(*net.TCPListener)}
	return listenerWrapTLS(ln, config)
}

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
//
// This is copied directly from the Go source code.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}
