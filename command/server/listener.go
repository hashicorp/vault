package server

import (
	"crypto/tls"
	"fmt"
	"net"
)

// ListenerFactory is the factory function to create a listener.
type ListenerFactory func(map[string]string) (net.Listener, map[string]string, error)

// BuiltinListeners is the list of built-in listener types.
var BuiltinListeners = map[string]ListenerFactory{
	"tcp": tcpListenerFactory,
}

// NewListener creates a new listener of the given type with the given
// configuration. The type is looked up in the BuiltinListeners map.
func NewListener(t string, config map[string]string) (net.Listener, map[string]string, error) {
	f, ok := BuiltinListeners[t]
	if !ok {
		return nil, nil, fmt.Errorf("unknown listener type: %s", t)
	}

	return f(config)
}

func listenerWrapTLS(
	ln net.Listener,
	props map[string]string,
	config map[string]string) (net.Listener, map[string]string, error) {
	props["tls"] = "disabled"

	if v, ok := config["tls_disable"]; ok && v != "" {
		return ln, props, nil
	}

	certFile, ok := config["tls_cert_file"]
	if !ok {
		return nil, nil, fmt.Errorf("'tls_cert_file' must be set")
	}

	keyFile, ok := config["tls_key_file"]
	if !ok {
		return nil, nil, fmt.Errorf("'tls_key_file' must be set")
	}

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, nil, fmt.Errorf("error loading TLS cert: %s", err)
	}

	tlsConf := &tls.Config{}
	tlsConf.Certificates = []tls.Certificate{cert}
	tlsConf.NextProtos = []string{"http/1.1"}

	ln = tls.NewListener(ln, tlsConf)
	props["tls"] = "enabled"
	return ln, props, nil
}
