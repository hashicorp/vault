package cache

import (
	"crypto/tls"
	"fmt"
	"net"
	"strings"

	"github.com/hashicorp/go-secure-stdlib/reloadutil"
	"github.com/hashicorp/vault/command/server"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/internalshared/listenerutil"
)

type CertConfig struct { // TODO: PW: Nicer name
	Config     *tls.Config
	ReloadFunc reloadutil.ReloadFunc
}

func StartListener(lnConfig *configutil.Listener) (net.Listener, *CertConfig, error) {
	addr := lnConfig.Address

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
		if lnConfig.SocketMode != "" &&
			lnConfig.SocketUser != "" &&
			lnConfig.SocketGroup != "" {
			uConfig = &listenerutil.UnixSocketsConfig{
				Mode:  lnConfig.SocketMode,
				User:  lnConfig.SocketUser,
				Group: lnConfig.SocketGroup,
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
	tlsConf, reloadFunc, err := listenerutil.TLSConfig(lnConfig, props, nil)
	if err != nil {
		return nil, nil, err
	}
	if tlsConf != nil {
		ln = tls.NewListener(ln, tlsConf)
	}

	cfg := &CertConfig{
		Config:     tlsConf,
		ReloadFunc: reloadFunc,
	}

	return ln, cfg, nil
}
