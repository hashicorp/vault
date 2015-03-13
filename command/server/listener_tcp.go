package server

import (
	"fmt"
	"net"
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

	return ln, nil
}
