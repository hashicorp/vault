package server

import (
	"io"
	"net"
	"strings"
	"time"

	"github.com/hashicorp/vault/helper/reload"
	"github.com/mitchellh/cli"
)

func tcpListenerFactory(config map[string]interface{}, _ io.Writer, ui cli.Ui) (net.Listener, map[string]string, reload.ReloadFunc, error) {
	bindProto := "tcp"
	var addr string
	addrRaw, ok := config["address"]
	if !ok {
		addr = "127.0.0.1:8200"
	} else {
		addr = addrRaw.(string)
	}

	// If they've passed 0.0.0.0, we only want to bind on IPv4
	// rather than golang's dual stack default
	if strings.HasPrefix(addr, "0.0.0.0:") {
		bindProto = "tcp4"
	}

	ln, err := net.Listen(bindProto, addr)
	if err != nil {
		return nil, nil, nil, err
	}

	ln = tcpKeepAliveListener{ln.(*net.TCPListener)}

	ln, err = listenerWrapProxy(ln, config)
	if err != nil {
		return nil, nil, nil, err
	}

	props := map[string]string{"addr": addr}
	return listenerWrapTLS(ln, props, config, ui)
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
