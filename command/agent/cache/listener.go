package cache

import (
	"crypto/tls"
	"fmt"
	"net"

	"strings"

	"github.com/hashicorp/vault/command/agent/config"
	"github.com/hashicorp/vault/command/server"
	"github.com/hashicorp/vault/helper/listenerutil"
)

func StartListener(lnConfig *config.Listener) (net.Listener, *tls.Config, error) {
	addr, ok := lnConfig.Config["address"].(string)
	if !ok {
		return nil, nil, fmt.Errorf("invalid address")
	}

	var ln net.Listener
	var err error
	switch lnConfig.Type {
	case "tcp":
		if addr == "" {
			addr = "127.0.0.1:8200"
		}

		bindProto := "tcp"
		// If they've passed 0.0.0.0, we only want to bind on IPv4
		// rather than golang's dual stack default
		if strings.HasPrefix(addr, "0.0.0.0:") {
			bindProto = "tcp4"
		}

		ln, err = net.Listen(bindProto, addr)
		if err != nil {
			return nil, nil, err
		}
		ln = &server.TCPKeepAliveListener{ln.(*net.TCPListener)}

	case "unix":
		var uConfig *listenerutil.UnixSocketsConfig
		if lnConfig.Config["socket_mode"] != nil &&
			lnConfig.Config["socket_user"] != nil &&
			lnConfig.Config["socket_group"] != nil {
			uConfig = &listenerutil.UnixSocketsConfig{
				Mode:  lnConfig.Config["socket_mode"].(string),
				User:  lnConfig.Config["socket_user"].(string),
				Group: lnConfig.Config["socket_group"].(string),
			}
		}
		ln, err = listenerutil.UnixSocketListener(addr, uConfig)
		if err != nil {
			return nil, nil, err
		}

	default:
		return nil, nil, fmt.Errorf("invalid listener type: %q", lnConfig.Type)
	}

	props := map[string]string{"addr": ln.Addr().String()}
	ln, props, _, tlsConf, err := listenerutil.WrapTLS(ln, props, lnConfig.Config, nil)
	if err != nil {
		return nil, nil, err
	}

	return ln, tlsConf, nil
}
