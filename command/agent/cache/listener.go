package cache

import (
	"crypto/tls"
	"fmt"
	"net"

	"strconv"
	"strings"

	"github.com/hashicorp/vault/command/agent/config"
	"github.com/hashicorp/vault/command/server"
	"github.com/hashicorp/vault/helper/listenerutil"
)

func StartListener(lnConfig *config.Listener, unixSocketsConfig *config.UnixSockets) (net.Listener, *tls.Config, error) {
	addr, ok := lnConfig.Config["address"].(string)
	if !ok {
		return nil, nil, fmt.Errorf("invalid address")
	}

	bindProto := "tcp"
	switch lnConfig.Type {
	case "tcp":
		if addr == "" {
			addr = "127.0.0.1:8007"
		}

		// If they've passed 0.0.0.0, we only want to bind on IPv4
		// rather than golang's dual stack default
		if strings.HasPrefix(addr, "0.0.0.0:") {
			bindProto = "tcp4"
		}

	case "unix":
		addr = "unix://" + addr
	default:
		return nil, nil, fmt.Errorf("invalid listener type: %q", lnConfig.Type)
	}

	var netAddr net.Addr
	switch {
	case strings.HasPrefix(addr, "unix://"):
		netAddr = &net.UnixAddr{
			Name: addr[len("unix://"):],
			Net:  "unix",
		}
	default:
		host, port, err := net.SplitHostPort(addr)
		if err != nil {
			return nil, nil, err
		}

		nPort, err := strconv.Atoi(port)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid tcp port %q", port)
		}

		ip := net.ParseIP(host)
		if ip == nil {
			return nil, nil, fmt.Errorf("invalid ip address %q", addr)
		}
		netAddr = &net.TCPAddr{
			IP:   ip,
			Port: nPort,
		}
	}

	var ln net.Listener
	var err error
	switch addrType := netAddr.(type) {
	case *net.UnixAddr:
		ln, err = listenerutil.UnixSocketListener(addrType.Name, unixSocketsConfig)
		if err != nil {
			return nil, nil, err
		}

	case *net.TCPAddr:
		ln, err = net.Listen(bindProto, addrType.String())
		if err != nil {
			return nil, nil, err
		}
		ln = &server.TCPKeepAliveListener{ln.(*net.TCPListener)}
	}

	props := map[string]string{"addr": ln.Addr().String()}
	ln, props, _, tlsConf, err := listenerutil.WrapTLS(ln, props, lnConfig.Config, nil)
	if err != nil {
		return nil, nil, err
	}

	return ln, tlsConf, nil
}
